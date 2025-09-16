#ifndef INVOKE_HELPER_H
#define INVOKE_HELPER_H

#if !RECASTNAVIGATION_STATIC && WIN32
#define RECAST_DLL _declspec(dllexport)
#else
#define RECAST_DLL
#endif
#include <cstdint>
#include <string>
#include <map>

#ifdef __cplusplus
extern "C" {
#endif

class NavMeshRegion;

RECAST_DLL NavMeshRegion* CreateNavMesh(int32_t id);

RECAST_DLL void ReleaseNavMesh(NavMeshRegion* pNavMesh);

RECAST_DLL int32_t InitData(NavMeshRegion* pNavMesh, const char* buffer, int32_t bufferLen);

RECAST_DLL int FindPath(NavMeshRegion* pNavMesh, const float* halfExtents, const float* startPos, const float* endPos, float* straightPath);
RECAST_DLL int FindPathLineOptimized(NavMeshRegion* pNavMesh, const float* halfExtents, const float* startPos, const float* endPos, float* straightPath);
RECAST_DLL int FindPathCrossAreaOptimized(NavMeshRegion* pNavMesh, const float* halfExtents, const float* startPos, const float* endPos, float* straightPath);
RECAST_DLL bool FindNearestPoint(NavMeshRegion* pNavMesh, const float* halfExtents, const float* center, float* nearestPos);
RECAST_DLL bool FindRandomPoint(NavMeshRegion* pNavMesh, float* pos);
RECAST_DLL bool FindRandomPointAroundCircle(NavMeshRegion* pNavMesh, const float* halfExtents, const float* centerPos, const float maxRadius, float* pos);
RECAST_DLL unsigned int AddCapsuleObstacle(NavMeshRegion* pNavMesh, const float* pos, const float radius, const float height);
RECAST_DLL unsigned int AddBoxObstacle(NavMeshRegion* pNavMesh, const float* bmin, const float* bmax);
RECAST_DLL unsigned int AddRotBoxObstacle(NavMeshRegion* pNavMesh, const float* center, const float* halfExtents, const float yRadians);
RECAST_DLL void RemoveObstacle(NavMeshRegion* pNavMesh, const unsigned int obstacleId);
RECAST_DLL void Update(NavMeshRegion* pNavMesh, const float dt);
RECAST_DLL void ClearData(NavMeshRegion* pNavMesh);
RECAST_DLL void GetMeshInfo(NavMeshRegion* pNavMesh, void (*AddMeshVec)(const float* ptr), void (*AddMeshConnection)(const float* ptr, int flag));


#ifdef __cplusplus
}
#endif


#endif
