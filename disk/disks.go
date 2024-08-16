package disk

import (
	"fmt"
	"os"
	"strconv"

	"github.com/vars7899/go_raid/utils"
)

type DiskCollection struct {
	Disks []*DirDisk
}

func InitializeDiskCollection(config *utils.Config)(*DiskCollection, error){
	_, err := os.Stat(config.SegmentDir)
	// if file does not exists, create a new directory
	if err == nil {
		fmt.Printf("%s already exists\n", config.SegmentDir)
		fmt.Printf("trying to remove %s\n", config.SegmentDir)
		if err := os.RemoveAll(config.SegmentDir); err != nil {
			return nil, err
		}
	}
	fmt.Printf("creating %s\n", config.SegmentDir)
	if err := os.MkdirAll(config.SegmentDir, 0766); err != nil {
			return nil, err
	}
	var dc DiskCollection
	for index := range config.DiskCount {
		diskName := generateDiskName(config.SegmentDir, config.DiskPrefix, config.DiskSuffix, strconv.Itoa(int(index)))
		disk, err := GenNewDirDisk(diskName)
		if err != nil {
			return nil, fmt.Errorf("failed to create <%s> disk", diskName)
		}
		dc.Disks = append(dc.Disks, disk)
	} 
	return &dc, nil
}
func (dc *DiskCollection) CloseDiskCollection() error {
	for _, disk := range dc.Disks {
		if err := disk.Close(); err != nil{
			return err
		}
	}
	fmt.Printf("closing all disks\n")
	return nil
} 
func generateDiskName(basePath string, prefix string, suffix string, name string) string{
	diskName := fmt.Sprintf("%s/%s%s%s.bin", basePath, prefix, name, suffix)
	return diskName
}