package main

import (
	"fmt"

	"github.com/spf13/viper"
)

// Configurations
type Configurations struct {
	Server struct {
		Host     string
		Port     int
		User     string
		Password string
	}
	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
	}
	ParallelScans int
	Groups        string
	Step          int
}

var conf Configurations

func loadConfig() error {

	// Set the file name of the configurations file
	viper.SetConfigName("config")

	// Set the path to look for the configurations file
	viper.AddConfigPath(".")

	// Set config type to yaml
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading configuration file, %s\n", err)
		return err
	}

	if err := viper.Unmarshal(&conf); err != nil {
		fmt.Printf("Unable to decode configure structure, %v\n", err)
		return err
	}

	fmt.Println("Configuration loaded")

	return nil
}
