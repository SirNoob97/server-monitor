package disk

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/SirNoob97/server-monitor/utils"
	"golang.org/x/exp/slices"
)

func Partitions() ([]PartitionInfo, error) {
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

	ret := make([]PartitionInfo, 0, len(fileData))
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

			fields := strings.Fields(sections[0])
			if !slices.Contains(filSystems, fields[0]) || fields[1] == "none" {
				continue
			}

			blockDeviceID := fields[2]
			mountPoint := fields[4]
			mountOpts := strings.Split(fields[5], ",")
			if rootDir := fields[2]; rootDir != "" && rootDir != "/" {
				mountOpts = append(mountOpts, "bind")
			}

			fields = strings.Fields(sections[1])
			fstype := fields[0]
			device := fields[1]
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
    ret = append(ret, pi)
	}

	return ret, nil
}

func getFileSystems(procDir string) ([]string, error) {
	filename := utils.JoinPath("filesystems")
	fileData, err := utils.ReadLines(filename)
	if err != nil {
		return nil, err
	}

	var ret []string
	for _, ln := range fileData {
		if !strings.HasPrefix(ln, "nodev") {
			ret = append(ret, ln)
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
		if !errors.As(err, pathErr) {
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
