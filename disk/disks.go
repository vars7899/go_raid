package disk

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	pb "github.com/schollz/progressbar/v3"
	"github.com/vars7899/go_raid/utils"
)

type DiskCollection struct {
	Disks		[]*DirDisk
	Mux			sync.RWMutex
}

func InitializeDiskCollection(config *utils.Config)(*DiskCollection, error){
	// Remove older data partials
	if err := RemoveDiskCollectionDir(config.SegmentDir); err != nil {
		return nil, err
	}
	// Create disk collection directory
	if err := CreateDiskCollectionDir(config.SegmentDir); err != nil {
		return nil, err
	}

	bar := pb.Default(int64(config.DiskCount), "○ Generating RAID Partials")

	var wg sync.WaitGroup
	var dc DiskCollection

	dc.Disks = make([]*DirDisk, config.DiskCount)


	for index := range config.DiskCount {
		wg.Add(1) 
		
		go func(idx uint){
			defer wg.Done()
			defer bar.Add(1) // update progress bar

			diskName := generateDiskName(config.SegmentDir, config.DiskPrefix, config.DiskSuffix, strconv.Itoa(int(idx)))
			disk, err := GenNewDirDisk(diskName)
			if err != nil {
				bar.Describe(fmt.Sprintf("○ Generating RAID Partial: <%s>", diskName))
				fmt.Printf("○ Failed to create <%s> disk: %v\n", diskName, err)
				return
			}

			dc.Mux.Lock()
			dc.Disks[idx] =  disk
			dc.Mux.Unlock()

			time.Sleep(time.Millisecond * 2000)
		}(index)
	}

	wg.Wait()

	return &dc, nil
}

func (dc *DiskCollection) CloseDiskCollection() error {
	var wg sync.WaitGroup
	var closeErr error

	for _, disk := range dc.Disks {
		wg.Add(1)
		go func (){
			defer wg.Done()
			if err := disk.Close(); err != nil{
				closeErr = err
			}
		}()
		wg.Wait()
	}
	fmt.Printf("○ Closing all RAID disks\n")
	return closeErr
} 

func generateDiskName(basePath string, prefix string, suffix string, name string) string{
	return fmt.Sprintf("%s/%s%s%s.bin", basePath, prefix, name, suffix)
}

func RemoveDiskCollectionDir(pathname string) error {
	_, err := os.Stat(pathname)
	// if file does not exists, create a new directory
	if err == nil {
		fmt.Printf("○ Found disk collection: <%s>\n", pathname)
		fmt.Printf("○ Removing disk collection: <%s>\n", pathname)
		if err := os.RemoveAll(pathname); err != nil {
			return err
		}
	}
	return nil
}

func CreateDiskCollectionDir(pathname string) error {
	if err := os.MkdirAll(pathname, 0766); err != nil {
		fmt.Printf("○ Failed to create disk collection: <%s>\n", pathname)
		return err
	}
	fmt.Printf("○ Created disk collection: <%s>\n", pathname)
	return nil
}