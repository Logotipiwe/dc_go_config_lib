package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

var configGot = false

func LoadDcConfigWithAttempts() error {
	for i := 0; i < 5; i++ {
		if i == 4 {
			return errors.New("failed to get config, shutting down")
		}
		err := LoadDcConfigInternal()
		if err != nil {
			println(err.Error())
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
	return nil
}

func LoadDcConfigInternal() error {
	var csUrl = os.Getenv("CONFIG_SERVER_URL")
	if !configGot {
		println(fmt.Sprintf("Getting config from %s...", csUrl))
	}

	request, err := http.NewRequest("GET", csUrl+"/api/get-config", nil)
	params := url.Values{}
	params.Add("mToken", os.Getenv("M_TOKEN"))
	params.Add("service", os.Getenv("SERVICE_NAME"))
	params.Add("namespace", os.Getenv("NAMESPACE"))

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
	for _, dto := range answer {
		err := os.Setenv(dto.Name, dto.Value)
		if err != nil {
			return err
		}
	}
	configGot = true
	return nil
}
