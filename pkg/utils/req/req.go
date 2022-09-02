package req

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"webase-server/models"
	"webase-server/pkg/errpkg"
)

func ReqJSON(method, urlStr string, body interface{}, respBody interface{}, headers map[string]string) error {
	bodyData, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(method, urlStr, bytes.NewBuffer(bodyData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := models.DefaultHTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode > 399 {
		msg := errpkg.ErrorMsg{
			Code:    resp.StatusCode,
			Message: string(data),
		}
		return msg
	}
	if respBody == nil {
		return nil
	}
	err = json.Unmarshal(data, &respBody)
	if err != nil {
		return err
	}
	return nil
}

func Req(method, urlStr string, body []byte, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest(method, urlStr, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := models.DefaultHTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 399 {
		msg := errpkg.ErrorMsg{
			Code:    resp.StatusCode,
			Message: string(data),
		}
		return data, msg
	}
	return data, nil
}
