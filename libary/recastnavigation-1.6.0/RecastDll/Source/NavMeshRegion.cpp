#include "NavMeshHelper.h"
#include "NavMeshRegion.h"
#include "DetourTileCacheBuilder.h"
#include "DetourCommon.h"
#include "DetourTileCache.h"
#include "RecastAlloc.h"
#include "RecastAssert.h"
#include <cstring>
#include <float.h>
#include <cstdarg>
#include <cstdlib>

static const int TILECACHESET_MAGIC = 'T' << 24 | 'S' << 16 | 'E' << 8 | 'T'; //'TSET';
static const int TILECACHESET_VERSION = 1;

static const int MAX_POLYS = 4096;

static const int MAX_QUERY_NODE = 2048;

struct TileCacheSetHeader
{
	int magic;
	int version;
	int numTiles;
	dtNavMeshParams meshParams;
	dtTileCacheParams cacheParams;
};

struct TileCacheTileHeader
{
	dtCompressedTileRef tileRef;
	int dataSize;
};



NavMeshRegion::NavMeshRegion(int32_t id) :
	m_bHasLoad(0),
	m_NavMeshId(id),
	m_pNavMesh(0),
	m_pNavQuery(0),
	m_pTileCache(0),
	m_pTalloc(0),
	m_pTcomp(0),
	m_pTmproc(0)
{
	m_pTalloc = new MyLinearAllocator(32 * 1024);
	m_pTcomp = new MyFastLZCompressor();
	m_pTmproc = new MyMeshProcess();

	m_filter.setIncludeFlags(0xffff);
	m_filter.setExcludeFlags(0);

	m_filter.setAreaCost(POLYAREA_GROUND, 1.0f);
	m_filter.setAreaCost(POLYAREA_WATER, 10.0f);
	m_filter.setAreaCost(POLYAREA_ROAD, 1.0f);
	m_filter.setAreaCost(POLYAREA_DOOR, 1.0f);
	m_filter.setAreaCost(POLYAREA_GRASS, 2.0f);
	m_filter.setAreaCost(POLYAREA_JUMP, 1.5f);
}

NavMeshRegion::~NavMeshRegion()
{
	Clear();
	if (m_pTalloc)
	{
		delete m_pTalloc;
		m_pTalloc = nullptr;
	}
	if (m_pTcomp)
	{
		delete m_pTcomp;
		m_pTcomp = nullptr;
	}
	if (m_pTmproc)
	{
		delete m_pTmproc;
		m_pTmproc = nullptr;
	}
}

int NavMeshRegion::InitData(const char* buffer, int32_t bufferLen)
{
	Clear();
	int index = 0;
	// Read header.
	TileCacheSetHeader header;
	int count = sizeof(TileCacheSetHeader);
	if (index + count > bufferLen)
		return 1;

	memcpy(&header, buffer + index, count);
	index += count;

	if (header.magic != TILECACHESET_MAGIC)
		return 2;

	if (header.version != TILECACHESET_VERSION)
		return 3;


	m_pNavMesh = dtAllocNavMesh();
	if (!m_pNavMesh)
		return 4;

	dtStatus status = m_pNavMesh->init(&header.meshParams);
	if (dtStatusFailed(status))
		return 5;

	m_pTileCache = dtAllocTileCache();
	if (!m_pTileCache)
		return 6;

	status = m_pTileCache->init(&header.cacheParams, m_pTalloc, m_pTcomp, m_pTmproc);
	if (dtStatusFailed(status))
		return 7;

	int tileHeaderSize = sizeof(TileCacheTileHeader);
	// Read tiles.
	for (int i = 0; i < header.numTiles; ++i)
	{
		TileCacheTileHeader tileHeader;

		if (index + tileHeaderSize > bufferLen)
			return 8;
		memcpy(&tileHeader, buffer + index, tileHeaderSize);
		index += tileHeaderSize;

		if (!tileHeader.tileRef || !tileHeader.dataSize)
			break;


		unsigned char* data = (unsigned char*)dtAlloc(tileHeader.dataSize, DT_ALLOC_PERM);
		if (!data) break;
		memset(data, 0, tileHeader.dataSize);

		if (index + tileHeader.dataSize > bufferLen)
		{
			dtFree(data);
			return 9;
		}
		memcpy(data, buffer + index, tileHeader.dataSize);
		index += tileHeader.dataSize;

		dtCompressedTileRef tile = 0;
		dtStatus addTileStatus = m_pTileCache->addTile(data, tileHeader.dataSize, DT_COMPRESSEDTILE_FREE_DATA, &tile);
		if (dtStatusFailed(addTileStatus))
		{
			dtFree(data);
		}

		if (tile)
			m_pTileCache->buildNavMeshTile(tile, m_pNavMesh);
	}

	m_pNavQuery = dtAllocNavMeshQuery();
	if (!m_pNavQuery) return 10;

	status = m_pNavQuery->init(m_pNavMesh, MAX_QUERY_NODE);
	if (!dtStatusSucceed(status)) return 11;

	m_bHasLoad = true;
	return 0;
}
bool IsPointOnLine(const float* start, const float* point, const float* end)
{
	// 计算向量
	float v1[3] = { point[0] - start[0], point[1] - start[1], point[2] - start[2] };
	float v2[3] = { end[0] - start[0], end[1] - start[1], end[2] - start[2] };
	dtVnormalize(v1);
	dtVnormalize(v2);

	// 计算向量的叉积
	float cross[3] = {
		v1[1] * v2[2] - v1[2] * v2[1],
		v1[2] * v2[0] - v1[0] * v2[2],
		v1[0] * v2[1] - v1[1] * v2[0]
	};

	// 如果叉积的模为零，则点在直线上
	static constexpr float epsilon = 1e-6;
	return std::abs(cross[0]) <= epsilon && std::abs(cross[1]) <= epsilon && std::abs(cross[2]) <= epsilon;
}

int NavMeshRegion::OptimizePath(float* straightPath, const dtPolyRef* straightPolys, const int straightPathCount, const int pathOptimizeLevel)
{
	if (pathOptimizeLevel == 0 || straightPathCount <= 2)
		return straightPathCount;
	int pointCount = 1;

	dtPolyRef lastPoly = straightPolys[0];
	float last[3], current[3], next[3];
	dtRaycastHit hit;
	dtStatus status = 0;

	for (int i = 1; i < straightPathCount - 1; ++i)
	{
		dtVcopy(last, straightPath + (pointCount - 1) * 3);
		dtVcopy(current, straightPath + i * 3);
		dtVcopy(next, straightPath + (i + 1) * 3);

		if(pathOptimizeLevel >= 1 && IsPointOnLine(last, current, next)) // Check if the point is on the line
		{
			continue;
		}


		if (pathOptimizeLevel >= 2) 
		{
			hit.t = 0;
			dtVset(hit.hitNormal, 0, 0, 0);
			hit.hitEdgeIndex = 0;
			hit.path = nullptr;
			hit.pathCount = 0;
			hit.maxPath = 0;
			hit.pathCost = 0;
			status = m_pNavQuery->raycast(lastPoly, last, next, &m_filter, DT_RAYCAST_USE_COSTS, &hit);
			if (dtStatusSucceed(status) && hit.t >= 1.0f)
			{
				continue;
			}
		}

		dtVcopy(straightPath + pointCount * 3, current);
		pointCount++;
		dtVcopy(last, current);
		lastPoly = straightPolys[i];

	}


	dtVcopy(straightPath + pointCount * 3, next);
	pointCount++;


	return pointCount;
}

int NavMeshRegion::DoPathFind(const float* halfExtents, const float* startPos, const float* endPos, float* pathResult, const int pathOptimizeLevel)
{
	thread_local static dtPolyRef polys[MAX_POLYS];
	thread_local static dtPolyRef straightPathPolys[MAX_POLYS];

	if (!m_bHasLoad) return -1;
	if (!startPos) return -2;
	if (!endPos) return -3;
	if (!pathResult) return -4;
	if (!halfExtents) return -5;

	dtPolyRef startRef = 0, endRef = 0;

	m_pNavQuery->findNearestPoly(startPos, halfExtents, &m_filter, &startRef, nullptr);
	m_pNavQuery->findNearestPoly(endPos, halfExtents, &m_filter, &endRef, nullptr);

	if (!startRef || !endRef) {
		//std::cerr << "Invalid start or end reference!" << std::endl;
		return -6;
	}

	int npolys = 0;
	m_pNavQuery->findPath(startRef, endRef, startPos, endPos, &m_filter, polys, &npolys, MAX_POLYS);

	if (npolys == 0) {
		//std::cerr << "No path found." << std::endl;
		return -7;
	}

	float epos1[3];
	dtVcopy(epos1, endPos);
	if (polys[npolys - 1] != endRef) {
		m_pNavQuery->closestPointOnPoly(polys[npolys - 1], endPos, epos1, nullptr);
	}

	//unsigned char straightPathFlags[MAX_POLYS];
	int nstraightPath = 0;
	m_pNavQuery->findStraightPath(
		startPos, epos1, polys, npolys,
		pathResult, 0, straightPathPolys,
		&nstraightPath, MAX_POLYS,
		DT_STRAIGHTPATH_ALL_CROSSINGS);

	int pathCount = OptimizePath(pathResult, straightPathPolys, nstraightPath, pathOptimizeLevel);

	//std::cout << "Straight path contains " << nstraightPath << " points." << std::endl;
	return pathCount;

}

int NavMeshRegion::FindPath(const float* halfExtents, const float* startPos, const float* endPos, float* pathResult)
{
	return DoPathFind(halfExtents, startPos, endPos, pathResult, 0);
}

int NavMeshRegion::FindPathLineOptimized(const float* halfExtents, const float* startPos, const float* endPos, float* pathResult)
{
	return DoPathFind(halfExtents, startPos, endPos, pathResult, 1);
}

int NavMeshRegion::FindPathCrossAreaOptimized(const float* halfExtents, const float* startPos, const float* endPos, float* pathResult)
{
	return DoPathFind(halfExtents, startPos, endPos, pathResult, 2);
}

int NavMeshRegion::FindNearestPoint(const float* halfExtents, const float* center, float* nearestPos)
{
	if (!m_bHasLoad) return -1;
	if (center == nullptr) return -2;
	if (nearestPos == nullptr) return -3;
	if (halfExtents == nullptr) return -4;
	dtPolyRef centerRef = 0;

	dtQueryFilter filter;
	filter.setIncludeFlags(0xffff);
	filter.setExcludeFlags(0);

	m_pNavQuery->findNearestPoly(center, halfExtents, &filter, &centerRef, nullptr);
	if (!centerRef)
	{
		return -5;
	}
	float epos1[3];
	dtVcopy(epos1, center);
	epos1[1] -= 1.0f;
	dtVcopy(nearestPos, center);
	nearestPos[1] += 1.0f;
	dtStatus status = m_pNavQuery->closestPointOnPoly(centerRef, epos1, nearestPos, nullptr);
	if (dtStatusSucceed(status))
		return 0;
	return -6;
}

static float frand()
{
	return (float)rand() / (float)RAND_MAX;
}

bool NavMeshRegion::FindRandomPoint(float* pos)
{
	if (!m_bHasLoad) return false;
	if (pos == nullptr) return false;
	dtQueryFilter filter;
	filter.setIncludeFlags(0xffff);
	filter.setExcludeFlags(0);

	dtPolyRef startRef = 0;
	dtStatus state = m_pNavQuery->findRandomPoint(&filter, frand, &startRef, pos);
	return dtStatusSucceed(state);
}

bool NavMeshRegion::FindRandomPointAroundCircle(const float* halfExtents, const float* centerPos, const float maxRadius, float* pos)
{
	if (!m_bHasLoad) return false;
	if (pos == nullptr) return false;
	dtQueryFilter filter;
	filter.setIncludeFlags(0xffff);
	filter.setExcludeFlags(0);

	dtPolyRef startRef = 0;
	dtPolyRef randomRef = 0;
	float startNearestPt[3];
	m_pNavQuery->findNearestPoly(centerPos, halfExtents, &filter, &startRef, startNearestPt);
	return dtStatusSucceed(m_pNavQuery->findRandomPointAroundCircle(startRef, centerPos, maxRadius, &filter, frand, &randomRef, pos));
}

unsigned int NavMeshRegion::AddCapsuleObstacle(const float* pos, const float radius, const float height) {
	if (!m_bHasLoad) return 0;
	dtObstacleRef obstacleId = 0;
	dtStatus status = m_pTileCache->addObstacle(pos, radius, height, &obstacleId);
	return obstacleId;
}

unsigned int NavMeshRegion::AddBoxObstacle(const float* bmin, const float* bmax) {
	if (!m_bHasLoad) return 0;
	dtObstacleRef obstacleId = 0;
	dtStatus status = m_pTileCache->addBoxObstacle(bmin, bmax, &obstacleId);
	return obstacleId;
}

unsigned int NavMeshRegion::AddBoxObstacle(const float* center, const float* halfExtents, const float yRadians) {
	if (!m_bHasLoad) return 0;
	dtObstacleRef obstacleId = 0;
	dtStatus status = m_pTileCache->addBoxObstacle(center, halfExtents, yRadians, &obstacleId);
	return obstacleId;
}

void NavMeshRegion::RemoveObstacle(const unsigned int obstacleId) {
	if (!m_bHasLoad || obstacleId == 0) {
		return;
	}
	m_pTileCache->removeObstacle((dtObstacleRef)obstacleId);
}

void NavMeshRegion::Update(const float dt)
{
	if (m_bHasLoad)
		m_pTileCache->update(dt, m_pNavMesh);
}

void NavMeshRegion::GetMeshInfo(void (*AddMeshVec)(const float* ptr), void (*AddMeshConnection)(const float* ptr, int flag))
{
	if (m_bHasLoad)
	{
		const dtNavMesh* pNavMesh = m_pNavMesh;
		for (int i = 0; i < pNavMesh->getMaxTiles(); ++i)
		{
			const dtMeshTile* getTile = pNavMesh->getTile(i);
			if (!getTile->header) continue;
			dtPolyRef base = pNavMesh->getPolyRefBase(getTile);

			for (int j = 0; j < getTile->header->polyCount; ++j)
			{
				const dtPoly* p = &getTile->polys[j];

				const dtMeshTile* tile = 0;
				const dtPoly* poly = 0;
				const dtPolyRef ref = base | (dtPolyRef)j;
				if (dtStatusFailed(m_pNavMesh->getTileAndPolyByRef(ref, &tile, &poly)))
					return;

				const unsigned int ip = (unsigned int)(poly - tile->polys);

				if (poly->getType() == DT_POLYTYPE_OFFMESH_CONNECTION)
				{
					dtOffMeshConnection* con = &tile->offMeshCons[ip - tile->header->offMeshBase];

					// Connection arc.
					AddMeshConnection(con->pos, con->flags);
				}
				else
				{
					const dtPolyDetail* pd = &tile->detailMeshes[ip];
					for (int i = 0; i < pd->triCount; ++i)
					{
						const unsigned char* t = &tile->detailTris[(pd->triBase + i) * 4];
						for (int j = 0; j < 3; ++j)
						{
							if (t[j] < poly->vertCount)
								AddMeshVec(&tile->verts[poly->verts[t[j]] * 3]);
							else
								AddMeshVec(&tile->detailVerts[(pd->vertBase + t[j] - poly->vertCount) * 3]);
						}
					}
				}
			}
		}
	}
}

void NavMeshRegion::Clear()
{
	m_bHasLoad = false;
	if (m_pNavMesh)
	{
		dtFreeNavMesh(m_pNavMesh);
		m_pNavMesh = nullptr;
	}
	if (m_pTileCache)
	{
		dtFreeTileCache(m_pTileCache);
		m_pTileCache = nullptr;
	}
	if (m_pNavQuery)
	{
		dtFreeNavMeshQuery(m_pNavQuery);
		m_pNavQuery = nullptr;
	}
}
