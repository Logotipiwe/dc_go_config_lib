package src

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
var dcConfigMap = make(map[string]string)

func GetConfig(key string) string {
	s, has := dcConfigMap[key]
	if has {
		return s
	}
	return os.Getenv(key)
}

func LoadDcConfig() {
	for i := 0; i < 5; i++ {
		if i == 4 {
			log.Fatal("Failed to get config, shutting down")
		}
		err := loadDcConfigInternal()
		if err != nil {
			println(err.Error())
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
}

func loadDcConfigInternal() error {
	csUrl := os.Getenv("CONFIG_SERVER_URL")
	println(fmt.Sprintf("Getting config from %s...", csUrl))
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
		dcConfigMap[dto.Name] = dto.Value
	}
	return nil
}
