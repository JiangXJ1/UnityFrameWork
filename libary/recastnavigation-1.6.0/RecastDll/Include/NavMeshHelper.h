#ifndef NAVMESH_HELPER_H
#define NAVMESH_HELPER_H

#include <DetourTileCacheBuilder.h>
#include <DetourTileCache.h>

enum PolyAreas
{
	POLYAREA_GROUND = 0,
	POLYAREA_WATER = 1,
	POLYAREA_ROAD = 2,
	POLYAREA_DOOR = 3,
	POLYAREA_GRASS = 4,
	POLYAREA_JUMP = 5,
};

enum PolyFlags
{
	POLYFLAGS_WALK = 0x01,      // Ability to walk (ground, grass, road)
	POLYFLAGS_SWIM = 0x02,      // Ability to swim (water).
	POLYFLAGS_DOOR = 0x04,      // Ability to move through doors.
	POLYFLAGS_JUMP = 0x08,      // Ability to jump.
	POLYFLAGS_DISABLED = 0x10,  // Disabled polygon
	POLYFLAGS_ALL = 0xffff      // All abilities.
};

struct MyLinearAllocator : public dtTileCacheAlloc
{
public:
	MyLinearAllocator(const size_t cap);
	~MyLinearAllocator();
	void resize(const size_t cap);
	void reset() override;
	void* alloc(const size_t size) override;
	void free(void* ptr) override;
private:
	unsigned char* buffer;
	size_t capacity;
	size_t top;
	size_t high;
};

struct MyFastLZCompressor : public dtTileCacheCompressor
{
	int maxCompressedSize(const int bufferSize) override;
	dtStatus compress(const unsigned char* buffer, const int bufferSize,
		unsigned char* compressed, const int maxCompressedSize, int* compressedSize) override;
	dtStatus decompress(const unsigned char* compressed, const int compressedSize,
		unsigned char* buffer, const int maxBufferSize, int* bufferSize) override;
};

struct MyMeshProcess : public dtTileCacheMeshProcess
{
	MyMeshProcess();
	void process(struct dtNavMeshCreateParams* params, unsigned char* polyAreas, unsigned short* polyFlags) override;
};

#endif
