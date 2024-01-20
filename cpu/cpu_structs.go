package cpu

type Info struct {
	CPU        int32   `json:"cpu"`
	VendorID   string  `json:"vendorId"`
	Family     string  `json:"family"`
	Model      string  `json:"model"`
	Stepping   int32   `json:"stepping"`
	CoreID     string  `json:"coreId"`
	Cores      int32   `json:"cores"`
	ModelName  string  `json:"modelName"`
	Mhz        float64 `json:"mhz"`
	CacheSize  int32   `json:"cacheSize"`
	Percentage float64 `json:"percentage"`
}

// https://man7.org/linux/man-pages/man5/proc.5.html
type PercStat struct {
	CPU     string
	User    float64
	System  float64
	Idle    float64
	Nice    float64
	Iowait  float64
	Irq     float64
	Softirq float64
}

func (p *PercStat) getCPUTotal() float64 {
	return p.User + p.System + p.Idle + p.Nice + p.Iowait + p.Irq + p.Softirq
}
