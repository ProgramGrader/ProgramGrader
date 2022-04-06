package common

import (
	"github.com/google/go-github/v41/github"
	"runtime"
)

const (
	AWSRegion = "us-east-2"

	DYTableName = "AutoGrader_Analytics"

	Sender = "no-reply-autograder@iuscsg.org"

	ConfigurationSet = "Default"

	CharSet = "UTF-8"

	ReplyToAddresses = "iuscompsec@gmail.com"

	DepHead = "jfdoyle@ius.edu" //"zac16530@gmail.com"

	GithubOrganization = "IUS-CS"

	ConfigUrl = "https://github.com/IUS-CS/AutoGraderConfig.git"

	VersionTag = "v1"

	DocOpenerURL = "https://github.com/IUS-CS/AutoGraderDocOpener/releases/download/Release/DocsOpener.jar"

	DockerTarPrefix = "rootfs/"

	OwnerPermRw = 0600

	HealthzUrlPath = "/healthz"

	ApiUrlPrefix = "/api"

	ContentUrlPrefix = ApiUrlPrefix + "/" + VersionTag + "/content/"

	MetadataUrlPath = ApiUrlPrefix + "/" + VersionTag + "/metadata"

	OpenscapUrlPath = ApiUrlPrefix + "/" + VersionTag + "/openscap"

	ChrootServePath = "/"

	OscapCveDir = "/tmp"

	Verbose = false

	BucketName = "autograder-bucket"
)

var GithubAPIKey map[string]interface{}

var AllRepos []*github.Repository

var ConfigPath = func() string {
	if runtime.GOOS == "windows" {
		return "C:\\tmp\\AssignmentConfig"
	} else {
		return "/tmp/AssignmentConfig"
	}
}

var TempPath = func() string {
	if runtime.GOOS == "windows" {
		return "C:\\tmp"
	} else {
		return "/tmp"
	}
}

var PathSeparator = func() string {
	if runtime.GOOS == "windows" {
		return "\\"
	} else {
		return "/"
	}
}
