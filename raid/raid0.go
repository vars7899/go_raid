package raid

import (
	"fmt"

	"github.com/vars7899/go_raid/disk"
)

type RAID0 struct {
	StripeSize	uint64 				// custom stripe size for each set
	Disks 		[]*disk.DirDisk		// cluster of disks
}
func CreateRAID0 (stripeSize int64, disks []*disk.DirDisk) (*RAID0, error) {
	if stripeSize <= 0 {
		return nil, fmt.Errorf("invalid stripe size")
	}
	return &RAID0{
		StripeSize: uint64(stripeSize),
		Disks: disks,
	}, nil
}
func (r *RAID0) Write(data []byte) error {
	if err := r.isRaid0Valid(); err != nil {
		return err
	}
	
	numDisks := len(r.Disks)

	dataLen := uint64(len(data))
	var dataWritten uint64

	disksOffset := make([]uint64, numDisks) 

	for dataWritten < dataLen {
		for index, disk := range r.Disks {
			remaining := dataLen - dataWritten
			if remaining <= 0 {
				break;
			}

			toWrite := r.StripeSize
			if remaining < r.StripeSize {
				toWrite = remaining
			} 

			if err := disk.Write(int64(disksOffset[index]), data[dataWritten:dataWritten+toWrite]); err != nil {
				return err
			}
			dataWritten += toWrite
			disksOffset[index] = toWrite
		}
	}
	return nil
}
func (r *RAID0) Read(size uint64) ([]byte, error) {
	if err := r.isRaid0Valid(); err != nil {
		return nil, err
	}

	numDisks := len(r.Disks)
	result := make([]byte, size)
	var bytesRead uint64 = 0

	disksOffset := make([]uint64, numDisks)

	for bytesRead < size {
		for index, disk := range r.Disks {
			remaining := size - bytesRead

			if remaining <= 0 {
				break
			}

			toRead := r.StripeSize
			if remaining < r.StripeSize {
				toRead = remaining
			}

			// Calculate where to start reading from in this disk
			currentOffset := disksOffset[index]

			// Create a buffer to store the read data from the current disk
			// readBuffer := make([]byte, toRead)

			// Read the data from the current disk
			bt, err := disk.Read(int64(currentOffset), int(toRead))
			if err != nil {
				return nil, err
			}
			fmt.Println(bt)
			// Copy the data from the buffer into the result
			copy(result[bytesRead:bytesRead+toRead], bt)

			// Update how much we have read
			bytesRead += toRead
			disksOffset[index] += toRead
		} 
	}
	return result, nil
}
func (r *RAID0) isRaid0Valid() error {
	if len(r.Disks) == 0 {
		return fmt.Errorf("no disks available")
	}
	return nil
}