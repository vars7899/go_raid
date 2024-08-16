package utils

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Level  		int   		`yaml:"level"`
	DiskCount   uint   		`yaml:"disk_count"`
	StripeSize 	uint64    	`yaml:"stripe_size"`
	InputDir   	string 		`yaml:"input_dir"`
	SegmentDir  string 		`yaml:"segment_dir"`
	OutputDir  	string 		`yaml:"output_dir"`
	DiskPrefix 	string 		`yaml:"disk_prefix"`
	DiskSuffix	string		`yaml:"disk_suffix"`
	Rebuild		bool		`yaml:"rebuild"`
}
func LoadConfig(configPathname string) (*Config, error) {
	var config Config

	userConfig, err := os.ReadFile(configPathname)
	if err != nil {
		return nil, fmt.Errorf("configuration error: %s", err)
	}
	if err := yaml.Unmarshal(userConfig, &config); err != nil {
		return nil, fmt.Errorf("configuration error in yaml: %s", err)
	}
	return &config, nil
}
