package cpu

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/SirNoob97/server-monitor/utils"
)

func Status() ([]Info, error) {
	esp, err := Specifications()
	if err != nil {
		return nil, err
	}
	if len(esp) == 0 {
		return nil, err
	}

	fileLocation := utils.FileLocation{
		Env:           "HOST_PROC",
		EnvDefaultVal: "/proc",
		Segments:      []string{"stat"},
	}

	percStat1, err := parsePercentageStats(fileLocation)
	if err != nil {
		return nil, err
	}

	if err := utils.Sleep(context.Background(), time.Second); err != nil {
		return nil, err
	}

	percStat2, err := parsePercentageStats(fileLocation)
	if err != nil {
		return nil, err
	}

	return addPercentage(esp, percStat1, percStat2)
}

func Specifications() ([]Info, error) {
	fileLocation := utils.FileLocation{
		Env:           "HOST_PROC",
		EnvDefaultVal: "/proc",
		Segments:      []string{"cpuinfo"},
	}
	fileData, err := utils.ReadFile(fileLocation)
	if err != nil {
		return nil, fmt.Errorf("Error reading file: cpuinfo, %s\n", err)
	}

	var ret []Info
	var processorName string

	c := &Info{CPU: -1, Cores: 1}
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
				overrideInfo(c)
				ret = append(ret, *c)
			}
			c.Cores = 1
			c.ModelName = processorName
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
		overrideInfo(c)
		ret = append(ret, *c)
	}

	return ret, nil
}

func parsePercentageStats(fileLocation utils.FileLocation) ([]PercStat, error) {
	fileData, err := utils.ReadFile(fileLocation)
	if err != nil {
		return nil, err
	}
	lines := []string{}
	for _, line := range fileData {
		if strings.HasPrefix(line, "cpu") {
			lines = append(lines, line)
		}
	}

	ret := make([]PercStat, 0, len(lines))
	for _, ln := range lines {
		ct, err := parseStat(ln)
		if err != nil {
			continue
		}
		ret = append(ret, *ct)
	}
	return ret, nil
}

func addPercentage(infos []Info, ps1, ps2 []PercStat) ([]Info, error) {
	if len(ps1) != len(ps2) {
		return nil, fmt.Errorf("Received different CPU counts: %d != %d", len(ps1), len(ps2))
	}

	percentages := make(map[int32]float64, len(ps2))
	for i, ps := range ps2 {
		cpu := strings.Replace(ps.CPU, "cpu", "", 1)
		cpuNum, err := strconv.ParseInt(cpu, 10, 64)
		if err != nil {
			continue
		}
		percentages[int32(cpuNum)] = calculatePercentage(ps1[i], ps)
	}

	for i, info := range infos {
		infos[i].Percentage = percentages[info.CPU]
	}

	return infos, nil
}

func calculatePercentage(ps1, ps2 PercStat) float64 {
	ps1Total, ps1Busy := calculateBusyTime(ps1)
	ps2Total, ps2Busy := calculateBusyTime(ps2)
	if ps2Busy <= ps1Busy {
		return 0
	}
	if ps2Total <= ps1Total {
		return 100
	}
	return math.Min(100, math.Max(0, (ps2Busy-ps1Busy)/(ps2Total-ps1Total)*100))
}

func calculateBusyTime(ps PercStat) (float64, float64) {
	total := ps.getCPUTotal()
	busy := total - ps.Idle - ps.Iowait
	return total, busy
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

func parseStat(line string) (*PercStat, error) {
	fields := strings.Fields(line)
	if len(fields) < 8 {
		return nil, errors.New("stat does not contain cpu info")
	}

	if !strings.HasPrefix(fields[0], "cpu") {
		return nil, errors.New("not contain cpu")
	}

	cpu := fields[0]
	if cpu == "cpu" {
		cpu = "cpu-total"
	}
	user, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return nil, err
	}
	nice, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return nil, err
	}
	system, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return nil, err
	}
	idle, err := strconv.ParseFloat(fields[4], 64)
	if err != nil {
		return nil, err
	}
	iowait, err := strconv.ParseFloat(fields[5], 64)
	if err != nil {
		return nil, err
	}
	irq, err := strconv.ParseFloat(fields[6], 64)
	if err != nil {
		return nil, err
	}
	softirq, err := strconv.ParseFloat(fields[7], 64)
	if err != nil {
		return nil, err
	}

	clocksPerSec := float64(100)
	ct := &PercStat{
		CPU:     cpu,
		User:    user / clocksPerSec,
		Nice:    nice / clocksPerSec,
		System:  system / clocksPerSec,
		Idle:    idle / clocksPerSec,
		Iowait:  iowait / clocksPerSec,
		Irq:     irq / clocksPerSec,
		Softirq: softirq / clocksPerSec,
	}

	return ct, nil
}
