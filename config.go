package dc_go_config_lib

import (
	"fmt"
	"github.com/logotipiwe/dc_go_config_lib/internal"
	"log"
	"os"
	"strconv"
	"time"
)

func GetConfig(key string) string {
	return os.Getenv(key)
}

func GetConfigOr(key, defaultVal string) string {
	env, has := os.LookupEnv(key)
	if !has {
		return defaultVal
	}
	return env
}

func GetConfigBool(key string) (bool, error) {
	val := GetConfig(key)
	res, err := strconv.ParseBool(val)
	log.Println(fmt.Sprintf("error getting bool %s from config: %v", key, err))
	return res, err
}

func GetConfigInt(key string) (int, error) {
	val := GetConfig(key)
	res, err := strconv.Atoi(val)
	if err != nil {
		log.Println(fmt.Sprintf("error getting int %s from config: %v", key, err))
	}
	return res, err
}

func LoadDcConfig() {
	err := internal.LoadDcConfigWithAttempts()
	if err != nil {
		log.Fatal(err)
	}
}

func LoadDcConfigDynamically(intervalSec int) {
	err := internal.LoadDcConfigWithAttempts()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			time.Sleep(time.Duration(intervalSec) * time.Second)
			err := internal.LoadDcConfigInternal()
			if err != nil {
				log.Println(err.Error())
			}
		}
	}()
}
