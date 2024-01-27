package memory

import (
	"strconv"
	"strings"

	"github.com/SirNoob97/server-monitor/pkg/utils"
)

func Status() (*VirtualMemory, error) {
	fileLocation := utils.FileLocation{
		Env:           "HOST_PROC",
		EnvDefaultVal: "/proc",
		Segments:      []string{"meminfo"},
	}
	fileData, err := utils.ReadFile(fileLocation)
	if err != nil {
		return nil, err
	}

	memavail := false
	ret := &VirtualMemory{}

	for _, ln := range fileData {
		fields := strings.Split(ln, ":")
		if len(fields) != 2 {
			continue
		}
		key := strings.TrimSpace(fields[0])
		value := strings.TrimSpace(fields[1])
		value = strings.Replace(value, " kB", "", -1)

		switch key {
		case "MemTotal":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, err
			}
			ret.Total = t
		case "MemFree":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, err
			}
			ret.Free = t
		case "MemAvailable":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, err
			}
			memavail = true
			ret.Available = t
		case "Buffers":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, err
			}
			ret.Buffers = t
		case "Cached":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, err
			}
			ret.Cached = t
		case "SReclaimable":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, err
			}
			ret.Sreclaimable = t
		case "SwapCached":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, err
			}
			ret.SwapCached = t
		case "SwapTotal":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, err
			}
			ret.SwapTotal = t
		case "SwapFree":
			t, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return ret, err
			}
			ret.SwapFree = t
		}
	}

	ret.Cached += ret.Sreclaimable

	if !memavail {
		ret.Available = ret.Cached + ret.Free
	}

	ret.Used = ret.Total - ret.Free - ret.Buffers - ret.Cached
	ret.UsedPercent = float64(ret.Used) / float64(ret.Total) * 100.0

	return ret, nil
}
