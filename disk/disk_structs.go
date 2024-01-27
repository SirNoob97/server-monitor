package disk

type UsageStat struct {
	Total             uint64  `json:"total"`
	Free              uint64  `json:"free"`
	Used              uint64  `json:"used"`
	UsedPercent       float64 `json:"usedPercent"`
	InodesTotal       uint64  `json:"inodesTotal"`
	InodesUsed        uint64  `json:"inodesUsed"`
	InodesFree        uint64  `json:"inodesFree"`
	InodesUsedPercent float64 `json:"inodesUsedPercent"`
}

type PartitionInfo struct {
	Device     string    `json:"device"`
	Mountpoint string    `json:"mountpoint"`
	FileSystem string    `json:"fileSystem"`
	Opts       []string  `json:"opts"`
	Usage      *UsageStat `json:"usage"`
}
