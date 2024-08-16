package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/vars7899/go_raid/disk"
	"github.com/vars7899/go_raid/raid"
	"github.com/vars7899/go_raid/utils"
)

func GenDisks (basePath string, diskPrefix string, diskSuffix string, numOfDisks uint) ([]*disk.DirDisk, error) {
	_, err := os.Stat(basePath)
	// if file does not exists, create a new directory
	if err == nil {
		fmt.Printf("%s already exists\n", basePath)
		fmt.Printf("trying to remove %s\n", basePath)
		if err := os.RemoveAll(basePath); err != nil {
			return nil, err
		}
	}
	fmt.Printf("creating %s\n", basePath)
	if err := os.MkdirAll(basePath, 0766); err != nil {
			return nil, err
	}
	var disks []*disk.DirDisk
	for index := range numOfDisks {
		diskName := fmt.Sprintf("%s/%s%d%s.bin", basePath, diskPrefix, index, diskSuffix)
		fmt.Println(diskName)
		disk, err := disk.GenNewDirDisk(diskName)
		if err != nil {
			return nil, fmt.Errorf("failed to create <%s> disk", diskName)
		}
		disks = append(disks, disk)
	} 
	return disks, nil
}
func CloseAllDisks(disks []*disk.DirDisk) error {
	for _, disk := range disks {
		if err := disk.Close(); err != nil{
			return err
		}
	}
	fmt.Printf("closing all disks\n")
	return nil
}

func ReadEntryFile(readBuffer *[]byte, entryPath string)(error){
	body, err := os.ReadFile(entryPath)
	if err != nil {
		return err
	}
	*readBuffer = body
	return nil
}
func AggregateData(data *[]byte, outputPath string)(error){
	fileDir := strings.Split(outputPath, "/")
	dirPath := strings.Join(fileDir[:len(fileDir) -1], "/")

	if err := os.MkdirAll(dirPath, 0766); err != nil {
		return err
	}
	os.WriteFile(outputPath, *data, 0766)
	return nil
}
func LoadConfiguration(configPath string) *utils.Config{
	config, err := utils.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	return config
}
func main() {
	config := LoadConfiguration("config.yaml")

	disks, err := GenDisks(config.SegmentDir, config.DiskPrefix, config.DiskSuffix, config.DiskCount)
	if err != nil {
		log.Fatalln(err)
	}
	defer CloseAllDisks(disks)

	r1 , err := raid.CreateRAID0(int64(config.StripeSize), disks)
	if err != nil {
		log.Fatalln(err)
	}

	var incomeData []byte;

	err = ReadEntryFile(&incomeData, config.InputDir)
	fmt.Printf("\n-length -> %d\n-data ->\n%s\n\n", len(incomeData), incomeData)
	if err != nil {
		log.Fatalln(err)
	}
	if err := r1.Write(incomeData); err != nil {
		log.Fatalf("Failed to write data: %v", err)
	}
	readData, err := r1.Read(uint64(len(incomeData)))
	if err != nil {
		log.Fatalf("Failed to write data: %v", err)
	}

	if config.Rebuild {
		AggregateData(&readData, config.OutputDir)
	}
}