package step

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

func FileCheck(config *viper.Viper) {
	for _, common := range []string{"harbor", "mysql"} {
		if ok, _ := FileExists("bin/" + config.Get("package").(map[string]interface{})[common].(string)); ok {
			log.Printf("%s [OK]", "bin/"+config.Get("package").(map[string]interface{})["mysql"].(string))
		} else {
			log.Fatalf("%s [False]", "bin/"+config.Get("package").(map[string]interface{})["mysql"].(string))
		}
	}

	for _, kube := range []string{"docker", "master", "node"} {
		for _, pkg := range config.Get("package").(map[string]interface{})[kube].([]interface{}) {
			if ok, _ := FileExists("bin/" + pkg.(string)); ok {
				log.Printf("%s [OK]", "bin/"+pkg.(string))
			} else {
				log.Printf("%s [False]", "bin/"+pkg.(string))
			}
		}
	}
}

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return true, nil
}