package configProcessors

import (
	"SubmissionProducer/internal/common"
	"fmt"
	"github.com/spf13/viper"
	"log"
)

func GetClassesFromConfig(path string, name string) common.Classes {

	vInstance := viper.New()

	vInstance.SetConfigName(name)   // name of config file (without extension)
	vInstance.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name
	vInstance.AutomaticEnv()

	vInstance.AddConfigPath(path) // optionally look for config in the working directory

	if err := vInstance.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Fatalf("Missing Config File. \n Expected file to be in:" + path)
		} else {
			// Config file was found but another error was produced
			log.Fatalf("Error occured when reading Config: %s", err)
		}
	}

	var config common.Classes

	err := vInstance.Unmarshal(&config)
	if err != nil {
		fmt.Printf("unable to decode into struct, %v\n", err)
	}

	return config
}
