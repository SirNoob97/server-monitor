package memory

type VirtualMemory struct {
	// Total amount of RAM
	Total uint64 `json:"total"`

	// RAM available for programs to allocate
	Available uint64 `json:"available"`

	// RAM used by programs
	Used uint64 `json:"used"`

	// Percentage of RAM used by programs
	UsedPercent float64 `json:"usedPercent"`

	// This is the kernel's notion of free memory; RAM chips whose bits nobody
	// cares about the value of right now. For a human consumable number,
	// Available is what you really want.
	Free uint64 `json:"free"`

	// Linux values
	// https://www.kernel.org/doc/Documentation/filesystems/proc.txt

	// Relatively temporary storage for raw disk blocks
	Buffers uint64 `json:"buffers"`

	// In-memory cache for files read from the disk (the pagecache).
	// Doesn't include SwapCached
	Cached uint64 `json:"cached"`

	// Part of Slab, that might be reclaimed, such as caches
	Sreclaimable uint64 `json:"sreclaimable"`

	// Memory that once was swapped out, is swapped back in but still also is in the swapfile
	// if memory is needed it doesn't need to be swapped out AGAIN because it is already in the swapfile.
	// This saves I/O
	SwapCached uint64 `json:"swapCached"`

	// Total amount of swap space available
	SwapTotal uint64 `json:"swapTotal"`

	// Memory which has been evicted from RAM, and is temporarily on the disk
	SwapFree uint64 `json:"swapFree"`
}
