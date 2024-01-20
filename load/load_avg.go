package load

import (
	"strconv"
	"strings"
	"syscall"

	"github.com/SirNoob97/server-monitor/utils"
)

func Status() (*LoadAvg, error) {
	fileLocation := utils.FileLocation{
		Env:           "HOST_PROC",
		EnvDefaultVal: "/proc",
		Segments:      []string{"loadavg"},
	}
	lines, err := utils.ReadFile(fileLocation)
	if err != nil {
		var sysInfo syscall.Sysinfo_t
		err = syscall.Sysinfo(&sysInfo)
		if err != nil {
			return nil, err
		}

		return &LoadAvg{
			Avg1:  float64(sysInfo.Loads[0]) / float64(1<<16),
			Avg5:  float64(sysInfo.Loads[0]) / float64(1<<16),
			Avg15: float64(sysInfo.Loads[0]) / float64(1<<16),
		}, nil
	}
	values := strings.Fields(lines[0])
	avg1, err := strconv.ParseFloat(values[0], 64)
	if err != nil {
		return nil, err
	}
	avg5, err := strconv.ParseFloat(values[1], 64)
	if err != nil {
		return nil, err
	}
	avg15, err := strconv.ParseFloat(values[2], 64)
	if err != nil {
		return nil, err
	}

	return &LoadAvg{
		Avg1:  avg1,
		Avg5:  avg5,
		Avg15: avg15,
	}, nil
}
