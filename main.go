package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/vars7899/go_raid/disk"
	"github.com/vars7899/go_raid/raid"
)

func GenDisks (basePath string, diskPrefix string, numOfDisks uint) ([]*disk.DirDisk, error) {
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
	if err := os.Mkdir(basePath, 0766); err != nil {
			return nil, err
	}
	var disks []*disk.DirDisk
	for index := range numOfDisks {
		diskName := fmt.Sprintf("%s/%s%d.txt", basePath, diskPrefix, index)
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

	body = bytes.ReplaceAll(body, []byte(" "), []byte("\n"))
	
	*readBuffer = body
	return nil
}

func main() {
	disks, err := GenDisks("raid0", "disk", 2)
	if err != nil {
		log.Fatalln(err)
	}
	defer CloseAllDisks(disks)

	r1 , err := raid.CreateRAID0(4, disks)
	if err != nil {
		log.Fatalln(err)
	}

	var incomeData []byte;

	
	// data := []byte("123456789")
	err = ReadEntryFile(&incomeData, "./entryFile.txt")
	fmt.Printf("\n-length -> %d\n-data ->\n%s\n\n", len(incomeData), incomeData)
	if err != nil {
		log.Fatalln(err)
	}
	if err := r1.Write(incomeData); err != nil {
		log.Fatalf("Failed to write data: %v", err)
	}
	// readData, err := r1.Read(uint64(len(incomeData)))
	// if err != nil {
	// 	log.Fatalf("Failed to write data: %v", err)
	// }
	// fmt.Printf("Read Data: %s\n", string(readData))

	// // Read data from the disk
	// readData, err := disk.Read(0, len(data))
	// if err != nil {
	// 	log.Fatalf("Failed to read data: %v", err)
	// }
	// fmt.Printf("Read Data: %s\n", string(readData))
}