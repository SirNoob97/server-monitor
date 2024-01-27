package disk

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/SirNoob97/server-monitor/pkg/utils"
	"golang.org/x/exp/slices"
	"golang.org/x/sys/unix"
)

func PartitionsInfo() ([]PartitionInfo, error) {
	procDir := "/proc/"
	mountInfoDir := utils.JoinPath(procDir, "1")
	fileData, useMounts, filename, err := readMountFile(mountInfoDir)
	if err != nil {
		mountInfoDir = utils.JoinPath(procDir, "self")
		fileData, useMounts, filename, err = readMountFile(mountInfoDir)
		if err != nil {
			return nil, err
		}
	}

	filSystems, err := getFileSystems(procDir)
	if err != nil {
		return nil, err
	}

	partitions := make([]PartitionInfo, 0, len(fileData))
	for _, ln := range fileData {
		var pi PartitionInfo
		if useMounts {
			fields := strings.Fields(ln)
			if fields[0] == "none" || !slices.Contains(filSystems, fields[2]) {
				continue
			}

			scapedMountPoint, err := strconv.Unquote(`"` + fields[1] + `"`)
			if err != nil {
				scapedMountPoint = fields[1]
			}

			pi = PartitionInfo{
				Device:     fields[0],
				Mountpoint: scapedMountPoint,
				FileSystem: fields[2],
				Opts:       strings.Fields(fields[3]),
			}
		} else {
			sections := strings.Split(ln, " - ")
			if len(sections) != 2 {
				return nil, fmt.Errorf("Cant parse file: %s, invalid line: %s\n", filename, ln)
			}

			firstSection := strings.Fields(sections[0])
			secondSection := strings.Fields(sections[1])
			if !slices.Contains(filSystems, secondSection[0]) || secondSection[1] == "none" {
				continue
			}

			blockDeviceID := firstSection[2]
			mountPoint := firstSection[4]
			mountOpts := strings.Split(firstSection[5], ",")
			if rootDir := firstSection[2]; rootDir != "" && rootDir != "/" {
				mountOpts = append(mountOpts, "bind")
			}

			fstype := secondSection[0]
			device := secondSection[1]
			if device == "/dev/root" {
				path := "/sys/dev/block/" + blockDeviceID
				devpath, err := os.Readlink(path)
				if err == nil {
					device = strings.Replace(device, "root", filepath.Base(devpath), 1)
				}
			}

			if strings.HasPrefix(device, "/dev/mapper/") {
				path, err := filepath.EvalSymlinks(device)
				if err == nil {
					device = path
				}
			}

			scapedMountPoint, err := strconv.Unquote(`"` + mountPoint + `"`)
			if err != nil {
				scapedMountPoint = mountPoint
			}

			pi = PartitionInfo{
				Device:     device,
				Mountpoint: scapedMountPoint,
				FileSystem: fstype,
				Opts:       mountOpts,
			}
		}
		partitions = append(partitions, pi)
	}

	for i := range partitions {
		partitions[i].Usage, err = PartitionUsage(partitions[i].Mountpoint)
		if err != nil {
			fmt.Printf("Error when reading partition usage: %s - %s", partitions[i].Mountpoint, err)
		}
	}

	return partitions, nil
}

func PartitionUsage(path string) (*UsageStat, error) {
	stat := unix.Statfs_t{}
	err := unix.Statfs(path, &stat)
	if err != nil {
		return nil, err
	}

	bSize := uint64(stat.Bsize)
	ret := &UsageStat{
		Total:       stat.Blocks * bSize,
		Free:        stat.Bavail * bSize,
		Used:        (stat.Blocks - stat.Bfree) * bSize,
		InodesTotal: stat.Files,
		InodesFree:  stat.Ffree,
		InodesUsed:  stat.Files - stat.Ffree,
	}

	if (ret.Used + ret.Free) == 0 {
		ret.UsedPercent = 0
	} else {
		ret.UsedPercent = (float64(ret.Used) / float64(ret.Used+ret.Free)) * 100.0
	}

	if ret.InodesTotal == 0 {
		ret.InodesUsedPercent = 0
	} else {
		ret.InodesUsedPercent = (float64(ret.InodesUsed) / float64(ret.InodesTotal)) * 100.0
	}

	return ret, nil
}

func getFileSystems(procDir string) ([]string, error) {
	filename := utils.JoinPath(procDir, "filesystems")
	fileData, err := utils.ReadLines(filename)
	if err != nil {
		return nil, err
	}

	var ret []string
	for _, ln := range fileData {
		if !strings.HasPrefix(ln, "nodev") {
			ret = append(ret, strings.TrimSpace(ln))
			continue
		}

		t := strings.Fields(ln)
		if len(t) != 2 || t[1] != "zfs" {
			continue
		}
		ret = append(ret, strings.TrimSpace(t[1]))
	}

	return ret, nil
}

func readMountFile(mountInfoDir string) ([]string, bool, string, error) {
	filename := utils.JoinPath(mountInfoDir, "mountinfo")
	fileData, err := utils.ReadLines(filename)
	useMounts := false
	if err != nil {
		var pathErr *os.PathError
		if !errors.As(err, &pathErr) {
			return nil, useMounts, filename, err
		}

		useMounts = true
		filename = utils.JoinPath(mountInfoDir, "mounts")
		fileData, err = utils.ReadLines(filename)
		if err != nil {
			return nil, useMounts, filename, err
		}

		return fileData, useMounts, filename, nil
	}
	return fileData, useMounts, filename, nil
}
