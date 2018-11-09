package common

import (
	"github.com/kylelemons/go-gypsy/yaml"
)

var ymlFile *yaml.File
//初始化ymlFile
func init()  {
	ymlFile,_ = yaml.ReadFile("./config/task_slice.yml")
}

func GetYmlFile()(*yaml.File){
	GetLog().Println("ymlFile",ymlFile)
	if ymlFile == nil {
		ymlFile,_ = yaml.ReadFile("../config/task_slice.yml")
	}
	return ymlFile
}
