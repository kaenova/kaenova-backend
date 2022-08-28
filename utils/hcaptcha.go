package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type HCaptchaReqeust struct {
	Secret   string `json:"secret"`
	Response string `json:"response"`
}

type HCaptchaServerResponse struct {
	Success     bool      `json:"success"`
	ChallangeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
}

var (
	ErrCannotMarshalRequest = errors.New("cannot marshal request")
	ErrCannotCreateRequest  = errors.New("cannot create request")
	ErrCannotSendRequest    = errors.New("cannot send request")
)

func VerifyHcaptcha(data HCaptchaReqeust) (bool, error) {
	var client = &http.Client{}
	var resData HCaptchaServerResponse

	reqString := fmt.Sprintf("https://hcaptcha.com/siteverify?response=%s&secret=%s", data.Response, data.Secret)

	req, err := http.NewRequest("POST", reqString, nil)
	if err != nil {
		return false, ErrCannotCreateRequest
	}

	response, err := client.Do(req)
	if err != nil {
		return false, ErrCannotSendRequest
	}

	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&resData)
	if err != nil {
		return false, err
	}

	return resData.Success, nil
}
