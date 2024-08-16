package main

import (
	"github.com/vars7899/go_raid/disk"
	"github.com/vars7899/go_raid/raid"
	"github.com/vars7899/go_raid/utils"
)

// func GenDisks (basePath string, diskPrefix string, diskSuffix string, numOfDisks uint) ([]*disk.DirDisk, error) {
// 	_, err := os.Stat(basePath)
// 	// if file does not exists, create a new directory
// 	if err == nil {
// 		fmt.Printf("%s already exists\n", basePath)
// 		fmt.Printf("trying to remove %s\n", basePath)
// 		if err := os.RemoveAll(basePath); err != nil {
// 			return nil, err
// 		}
// 	}
// 	fmt.Printf("creating %s\n", basePath)
// 	if err := os.MkdirAll(basePath, 0766); err != nil {
// 			return nil, err
// 	}
// 	var disks []*disk.DirDisk
// 	for index := range numOfDisks {
// 		diskName := fmt.Sprintf("%s/%s%d%s.bin", basePath, diskPrefix, index, diskSuffix)
// 		fmt.Println(diskName)
// 		disk, err := disk.GenNewDirDisk(diskName)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to create <%s> disk", diskName)
// 		}
// 		disks = append(disks, disk)
// 	}
// 	return disks, nil
// }
// func CloseAllDisks(disks []*disk.DirDisk) error {
// 	for _, disk := range disks {
// 		if err := disk.Close(); err != nil{
// 			return err
// 		}
// 	}
// 	fmt.Printf("closing all disks\n")
// 	return nil
// }

// func ReadEntryFile(readBuffer *[]byte, entryPath string)(error){
// 	body, err := os.ReadFile(entryPath)
// 	utils.LogFatal(err)

// 	*readBuffer = body
// 	return nil
// }
// func RebuildData(data *[]byte, outputPath string)(error){
// 	fileDir := strings.Split(outputPath, "/")
// 	dirPath := strings.Join(fileDir[:len(fileDir) -1], "/")

// 	if err := os.MkdirAll(dirPath, 0766); err != nil {
// 		return err
// 	}
// 	os.WriteFile(outputPath, *data, 0766)
// 	return nil
// }

func main() {
	// 1. Load raid configuration
	config, err := utils.LoadConfig("config.yaml")
	utils.LogFatal(err)

	// 2. Load disks from configuration
	diskCollection, err := disk.InitializeDiskCollection(config)
	utils.LogFatal(err)
	defer func (){
		err = diskCollection.CloseDiskCollection()
		utils.LogFatal(err)
	}()

	// 3. Load raid with configuration level
	r0 , err := raid.CreateRAID0(int64(config.StripeSize), diskCollection.Disks)
	utils.LogFatal(err)

	// 4. Write RAID partials
	r0.BuildData(config)
	// var incomeData []byte;
	// err = ReadEntryFile(&incomeData, config.InputDir)

	// err = r0.Write(incomeData)
	// utils.LogFatal(fmt.Errorf("failed to write data: %v", err))

	// 5. Rebuild from RAID partials (optional)

	// readData, err := r0.Read(uint64(len(incomeData)))
	// utils.LogFatal(fmt.Errorf("failed to read data from raid segments: %v", err))


	if config.Rebuild {
		r0.RebuildData(config.OutputDir)
	}
}