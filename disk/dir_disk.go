package disk

import "os"

type DirDisk struct {
	Path 	string
	File 	*os.File
}

func GenNewDirDisk(path string)(*DirDisk, error){
	file, err := os.Create(path);
	if err != nil {
		return nil, err
	}
	return &DirDisk{
		Path: path,
		File: file,
	}, nil
}
func (d *DirDisk) Write(offset int64, data []byte) error {
	if len(data) == 0 {
		return nil
	}
	_, err := d.File.WriteAt(data, offset)
	return err
}
func (d *DirDisk) Read(offset int64, size int) ([]byte, error){
	data := make([]byte, size)
	_, err := d.File.ReadAt(data, offset)
	return data, err
}
func (d *DirDisk) Close() error {
	return d.File.Close()
}