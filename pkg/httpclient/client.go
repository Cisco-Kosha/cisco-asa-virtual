package httpclient

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/kosha/cisco-asa-virtual/pkg/logger"

	"io/ioutil"
	"net/http"
)

var token string

type AuthToken struct {
	Token string `json:"Token,omitempty"`
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func MakeHttpBasicAuthCall(headers map[string]string, username, password, method, url string, body interface{}, log logger.Logger, isSecure bool) (interface{}, int, error) {

	var req *http.Request
	if body != nil {
		jsonReq, _ := json.Marshal(body)
		req, _ = http.NewRequest(method, url, bytes.NewBuffer(jsonReq))
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	var response interface{}

	res, statusCode := makeHttpBasicAuthReq(username, password, req, log, isSecure)
	if string(res) == "" {
		return nil, statusCode, fmt.Errorf("nil")
	}
	// Convert response body to target struct
	err := json.Unmarshal(res, &response)
	if err != nil {
		log.Error("Unable to parse response as json")
		log.Error(err)
		return nil, 500, err
	}
	return response, statusCode, nil
}

func makeHttpBasicAuthReq(username, password string, req *http.Request, log logger.Logger, isSecure bool) ([]byte, int) {

	req.Header.Set("Authorization", "Basic "+basicAuth(username, password))
	req.Header.Set("User-Agent", "REST API Agent")

	req.Header.Set("Accept-Encoding", "identity")

	client := &http.Client{}

	if !isSecure {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.Transport = tr
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Error(err)
		return nil, 500
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}
	return bodyBytes, resp.StatusCode
}
