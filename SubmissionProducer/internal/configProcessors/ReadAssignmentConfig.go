package configProcessors

import (
	"SubmissionProducer/internal/common"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"log"
)

func ReadAssignmentConfig(path string) common.AssignmentConfig {

	vInstance := viper.New()

	vInstance.SetConfigName("Assignment") // name of config file (without extension)
	vInstance.SetConfigType("yaml")       // REQUIRED if the config file does not have the extension in the name
	vInstance.AutomaticEnv()

	vInstance.AddConfigPath(path)

	if err := vInstance.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Fatalf("Missing Assignment Config File:" + path + "assignment.yaml")
		} else {
			// Config file was found but another error was produced
			log.Fatalf("Error occured when reading Config: %s", err)
		}
	}

	//vInstance.SetDefault("GradeDocs", true)
	//vInstance.SetDefault("NonCodeSubmissions", false)
	//vInstance.SetDefault("StudentTestsEnabled", false)
	//vInstance.SetDefault("NumberStudentTestsRequired", 4)

	config := common.AssignmentConfig{}

	err := vInstance.Unmarshal(&config)
	common.CheckIfErrorWithMessage(err, "unable to decode into struct")

	if config.NumberStudentTestsRequired < 1 {
		config.NumberStudentTestsRequired = 1
	}

	validate := validator.New()
	if err := validate.Struct(&config); err != nil {
		log.Fatalf("Missing required attributes %v\n", err)
	}

	return config

}
