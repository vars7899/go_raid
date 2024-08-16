package main

import (
	"github.com/vars7899/go_raid/disk"
	"github.com/vars7899/go_raid/raid"
	"github.com/vars7899/go_raid/utils"
)


func main() {
	// 1. Load raid configuration
	config, err := utils.LoadConfig("config.yaml")
	utils.LogFatal(err)

	// 2. Load disks from configuration
	diskCollection, err := disk.InitializeDiskCollection(config)
	utils.LogFatal(err)
	defer func (){
		// clean up
		err = diskCollection.CloseDiskCollection()
		utils.LogFatal(err)
	}()

	// 3. Load raid with configuration level
	r0 , err := raid.CreateRAID0(int64(config.StripeSize), diskCollection.Disks)
	utils.LogFatal(err)

	// 4. Write RAID partials
	r0.BuildData(config)

	// 5. Rebuild from RAID partials
	if config.Rebuild {
		r0.RebuildData(config.OutputDir)
	}
}