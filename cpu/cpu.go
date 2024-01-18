package cpu

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/SirNoob97/server-monitor/utils"
)

func Status() ([]Info, error) {
	fileLocation := utils.FileLocation{
		Env:           "HOST_PROC",
		EnvDefaultVal: "/proc",
		Segments:      []string{"cpuinfo"},
	}
	fileData, err := utils.ReadFile(fileLocation)
	if err != nil {
		fmt.Printf("Error reading file: cpuinfo, %s\n", err)
	}

	var ret []Info
	var processorName string

	c := Info{CPU: -1, Cores: 1}
	for _, line := range fileData {
		fields := strings.Split(line, ":")
		if len(fields) < 2 {
			continue
		}
		key := strings.TrimSpace(fields[0])
		value := strings.TrimSpace(fields[1])

		switch key {
		case "Processor":
			processorName = value
		case "processor", "cpu number":
			if c.CPU >= 0 {
				overrideInfo(&c)
				ret = append(ret, c)
			}
			c = Info{Cores: 1, ModelName: processorName}
			t, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return ret, err
			}
			c.CPU = int32(t)
		case "vendorId", "vendor_id":
			c.VendorID = value
			if strings.Contains(value, "S390") {
				processorName = "S390"
			}
		case "cpu family":
			c.Family = value
		case "model", "CPU part":
			c.Model = value
		case "Model Name", "model name", "cpu":
			c.ModelName = value
		case "stepping", "CPU revision":
			val := value

			t, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return ret, err
			}
			c.Stepping = int32(t)
		case "cpu MHz", "cpu MHz dynamic":
			// treat this as the fallback value, thus we ignore error
			if t, err := strconv.ParseFloat(strings.Replace(value, "MHz", "", 1), 64); err == nil {
				c.Mhz = t
			}
		case "cache size":
			t, err := strconv.ParseInt(strings.Replace(value, " KB", "", 1), 10, 64)
			if err != nil {
				return ret, err
			}
			c.CacheSize = int32(t)
		case "core id":
			c.CoreID = value
		}
	}
	if c.CPU >= 0 {
		overrideInfo(&c)
		ret = append(ret, c)
	}
	return ret, nil
}

func overrideInfo(c *Info) {
	fileLocation := utils.FileLocation{
		Env:           "HOST_SYS",
		EnvDefaultVal: "/sys",
		Segments:      []string{fmt.Sprintf("devices/system/cpu/cpu%d", c.CPU)},
	}

	if len(c.CoreID) == 0 {
		fileLocation.Segments = append(fileLocation.Segments, "topology/core_id")
		fileData, err := utils.ReadFile(fileLocation)
		if err != nil {
			fmt.Printf("Error reading file: %s\n", strings.Join(fileLocation.Segments, "/"))
		}
		c.CoreID = fileData[0]
	}

	// override the value of c.Mhz because I want to report the maximum
	// clock-speed of the CPU for c.Mhz
	if len(fileLocation.Segments) > 1 {
		fileLocation.Segments = fileLocation.Segments[:len(fileLocation.Segments)-1]
	}
	fileLocation.Segments = append(fileLocation.Segments, "cpufreq/cpuinfo_max_freq")
	fileData, err := utils.ReadFile(fileLocation)
	if err != nil || len(fileData) == 0 {
		fmt.Printf("Error reading file: %s\n", strings.Join(fileLocation.Segments, "/"))
		return
	}

	value, err := strconv.ParseFloat(fileData[0], 64)
	if err != nil {
		return
	}

	c.Mhz = value / 1000.0 // value is in kHz
	if c.Mhz > 9999 {
		c.Mhz = c.Mhz / 1000.0 // value in Hz
	}
}
