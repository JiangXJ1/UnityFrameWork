#include "InvokeHelper.h"
#include "NavMeshRegion.h"

NavMeshRegion* CreateNavMesh(int32_t id)
{
	return new NavMeshRegion(id);
}

void ReleaseNavMesh(NavMeshRegion* pNavMesh)
{
	if (pNavMesh != nullptr)
	{
		delete pNavMesh;
		pNavMesh = nullptr;
	}
}

int32_t InitData(NavMeshRegion* pNavMesh, const char* buffer, int32_t bufferLen)
{
	if (pNavMesh == nullptr)
		return -1;
	return pNavMesh->InitData(buffer, bufferLen);
}

int FindPath(NavMeshRegion* pNavMesh, const float* halfExtents, const float* startPos, const float* endPos, float* straightPath)
{
	if (pNavMesh)
		return pNavMesh->FindPath(halfExtents, startPos, endPos, straightPath);
	return 0;
}

int FindPathLineOptimized(NavMeshRegion* pNavMesh, const float* halfExtents, const float* startPos, const float* endPos, float* straightPath) 
{
	if (pNavMesh)
		return pNavMesh->FindPathLineOptimized(halfExtents, startPos, endPos, straightPath);
	return 0;
}

int FindPathCrossAreaOptimized(NavMeshRegion* pNavMesh, const float* halfExtents, const float* startPos, const float* endPos, float* straightPath) 
{
	if (pNavMesh)
		return pNavMesh->FindPathCrossAreaOptimized(halfExtents, startPos, endPos, straightPath);
	return 0;
}

bool FindNearestPoint(NavMeshRegion* pNavMesh, const float* halfExtents, const float* center, float* nearestPos)
{
	if (pNavMesh)
		return pNavMesh->FindNearestPoint(halfExtents, center, nearestPos);
	return false;
}

bool FindRandomPoint(NavMeshRegion* pNavMesh, float* pos)
{
	if (pNavMesh)
		return pNavMesh->FindRandomPoint(pos);
	return false;
}

bool FindRandomPointAroundCircle(NavMeshRegion* pNavMesh, const float* halfExtents, const float* centerPos, const float maxRadius, float* pos)
{
	if (pNavMesh)
		return pNavMesh->FindRandomPointAroundCircle(halfExtents, centerPos, maxRadius, pos);
	return false;
}

unsigned int AddCapsuleObstacle(NavMeshRegion* pNavMesh, const float* pos, const float radius, const float height)
{
	if (pNavMesh)
		return pNavMesh->AddCapsuleObstacle(pos, radius, height);
	return 0;
}

unsigned int AddBoxObstacle(NavMeshRegion* pNavMesh, const float* bmin, const float* bmax)
{
	if (pNavMesh)
		return pNavMesh->AddBoxObstacle(bmin, bmax);
	return 0;
}

unsigned int AddRotBoxObstacle(NavMeshRegion* pNavMesh, const float* center, const float* halfExtents, const float yRadians)
{
	if (pNavMesh)
		return pNavMesh->AddBoxObstacle(center, halfExtents, yRadians);
	return 0;
}

void RemoveObstacle(NavMeshRegion* pNavMesh, const unsigned int obstacleId)
{
	if (pNavMesh)
		pNavMesh->RemoveObstacle(obstacleId);
}

void Update(NavMeshRegion* pNavMesh, const float dt)
{
	if (pNavMesh)
		pNavMesh->Update(dt);
}

void ClearData(NavMeshRegion* pNavMesh)
{
	if (pNavMesh)
		pNavMesh->Clear();
}

void GetMeshInfo(NavMeshRegion* pNavMesh, void(*AddMeshVec)(const float* ptr), void(*AddMeshConnection)(const float* ptr, int flag))
{
	if (pNavMesh)
		pNavMesh->GetMeshInfo(AddMeshVec, AddMeshConnection);
}
