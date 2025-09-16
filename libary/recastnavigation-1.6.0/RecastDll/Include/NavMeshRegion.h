#ifndef NAVMESH_REGION_H
#define NAVMESH_REGION_H
#include <cstdint>
#include <stdio.h>
#include "DetourNavMesh.h"
#include "DetourNavMeshQuery.h"

void AddMeshVec(float* ptr);

class NavMeshRegion
{

private:
	bool m_bHasLoad;
	int32_t m_NavMeshId;
	class dtQueryFilter m_filter;

	class dtNavMesh* m_pNavMesh;
	class dtNavMeshQuery* m_pNavQuery;
	class dtTileCache* m_pTileCache;
	struct MyLinearAllocator* m_pTalloc;
	struct MyFastLZCompressor* m_pTcomp;
	struct MyMeshProcess* m_pTmproc;

	int OptimizePath(float* straightPath, const dtPolyRef* straightPolys, const int straightPathCount, const int pathOptimizeLevel);
	int DoPathFind(const float* halfExtents, const float* startPos, const float* endPos, float* pathResult, const int doPathOptimize);

public:
	NavMeshRegion(int32_t id);
	~NavMeshRegion();

	dtNavMesh* GetNavMesh() { return m_pNavMesh; }
	dtNavMeshQuery* GetQuery(){ return m_pNavQuery; }
	dtTileCache* GetTileCache(){ return m_pTileCache; }

	int InitData(const char* buffer, int32_t bufferLen);
	int FindPath(const float* halfExtents, const float* startPos, const float* endPos, float* pathResult);
	int FindPathLineOptimized(const float* halfExtents, const float* startPos, const float* endPos, float* pathResult);
	int FindPathCrossAreaOptimized(const float* halfExtents, const float* startPos, const float* endPos, float* pathResult);
	int FindNearestPoint(const float* halfExtents, const float* center, float* nearestPos);
	bool FindRandomPoint(float* pos);
	bool FindRandomPointAroundCircle(const float* halfExtents, const float* centerPos, const float maxRadius, float* pos);
	unsigned int AddCapsuleObstacle(const float* pos, const float radius, const float height);
	unsigned int AddBoxObstacle(const float* bmin, const float* bmax);
	unsigned int AddBoxObstacle(const float* center, const float* halfExtents, const float yRadians);
	void RemoveObstacle(const unsigned int obstacleId);
	void Update(const float dt);
	void GetMeshInfo(void (*AddMeshVec)(const float* ptr), void (*AddMeshConnection)(const float* ptr, int flag));
	void Clear();
};


#endif