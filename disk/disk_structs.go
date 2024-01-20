package disk

type PartitionInfo struct {
	Device     string   `json:"device"`
	Mountpoint string   `json:"mountpoint"`
	FileSystem string   `json:"fileSystem"`
	Opts       []string `json:"opts"`
}
