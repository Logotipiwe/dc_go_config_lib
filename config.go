package dc_go_config_lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

var dcConfig []DcPropertyDto

// var dcConfigMap = make(map[string]string)
var csUrl = os.Getenv("CONFIG_SERVER_URL")

func GetConfig(key string) string {
	//s, has := dcConfigMap[key]
	//if has {
	//	return s
	//}
	return os.Getenv(key)
}

func GetConfigOr(key, defaultVal string) string {
	env, has := os.LookupEnv(key)
	if !has {
		return defaultVal
	}
	return env
}

func LoadDcConfig() {
	err := loadDcConfigWithAttempts()
	if err != nil {
		log.Fatal(err)
	}
}

func loadDcConfigWithAttempts() error {
	for i := 0; i < 5; i++ {
		if i == 4 {
			return errors.New("Failed to get config, shutting down")
		}
		println(fmt.Sprintf("Getting config from %s...", csUrl))
		err := loadDcConfigInternal(csUrl)
		if err != nil {
			println(err.Error())
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
	return nil
}

func LoadDcConfigDynamically(intervalSec int) {
	err := loadDcConfigWithAttempts()
	if err != nil {
		log.Fatal(err)
	} else {
		go func() {
			for {
				time.Sleep(time.Duration(intervalSec) * time.Second)
				err := loadDcConfigInternal(csUrl)
				if err != nil {
					log.Println(err.Error())
				}
			}
		}()
	}
}

func loadDcConfigInternal(csUrl string) error {
	request, err := http.NewRequest("GET", csUrl+"/api/get-config", nil)
	params := url.Values{}
	params.Add("mToken", GetConfig("M_TOKEN"))
	params.Add("service", GetConfig("SERVICE_NAME"))
	params.Add("namespace", GetConfig("NAMESPACE"))

	request.URL.RawQuery = params.Encode()

	if err != nil {
		return err
	}
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return errors.New("Got status " + res.Status)
	}

	defer res.Body.Close()
	var answer []DcPropertyDto
	err = json.NewDecoder(res.Body).Decode(&answer)
	if err != nil {
		return err
	}
	dcConfig = answer
	for _, dto := range dcConfig {
		//dcConfigMap[dto.Name] = dto.Value
		err := os.Setenv(dto.Name, dto.Value)
		if err != nil {
			return err
		}
	}
	return nil
}
