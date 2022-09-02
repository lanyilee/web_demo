package test

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	_ "webase-server/server/api/v2"
)

// var uri = "https://webase-dev.webase.cooxun.com"
var uri = "https://127.0.0.1:8888"
var jsonCreateForAll = []byte(`{"username":"fortest","password":"Qwer1234","name":"租户","role":"user"}`)

// var jsonCreateForDelete = []byte(`{"username":"fortest1","password":"Qwer1234","name":"租户","role":"user"}`)

func TestLogin(t *testing.T) {

	var jsonStr = []byte(`{"username":"admin","password":"Qwer1234"}`)
	r, _ := http.NewRequest(http.MethodPost, uri+"/webase/public/v2/login", bytes.NewBuffer(jsonStr))
	r.Header.Set("Content-Type", "application/json")

	timeout := time.Duration(10 * time.Second)
	client := &http.Client{
		Timeout:   timeout,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	resp, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)

	var dat map[string]string
	err = json.Unmarshal([]byte(body), &dat)

	var expected map[string]string
	expected = make(map[string]string)
	expected["username"] = "admin"
	expected["name"] = "管理员"
	expected["role"] = "admin"
	if dat["username"] != expected["username"] && dat["name"] != expected["name"] {
		t.Errorf("handler returned unexpected body: got %v want %v",
			string(body), expected)
	}

}

func GetLoginToken() string {
	var jsonStr = []byte(`{"username":"admin","password":"Qwer1234"}`)
	r, _ := http.NewRequest(http.MethodPost, uri+"/webase/public/v2/login", bytes.NewBuffer(jsonStr))
	r.Header.Set("Content-Type", "application/json")

	timeout := time.Duration(10 * time.Second)
	client := &http.Client{
		Timeout:   timeout,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	resp, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		fmt.Println("login Response code is:", resp.StatusCode)
	}

	var dat map[string]string
	err = json.Unmarshal([]byte(body), &dat)

	token := dat["token"]
	return token

}

func TestUserList(t *testing.T) {
	token := GetLoginToken()

	r, err := http.NewRequest("GET", uri+"/webase/api/v2/users/", nil)
	r.Header.Set("Authorization", string(token))

	timeout := time.Duration(10 * time.Second)
	client := &http.Client{
		Timeout:   timeout,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	resp, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("user list response Status:", resp.Status)
	// body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("response Body:", string(body))
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Response code is %v", resp.StatusCode)
	}
}

func TestUserCreate(t *testing.T) {
	token := GetLoginToken()

	r, err := http.NewRequest("POST", uri+"/webase/api/v2/user/", bytes.NewBuffer(jsonCreateForAll))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", string(token))

	timeout := time.Duration(10 * time.Second)
	client := &http.Client{
		Timeout:   timeout,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	resp, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("user create response Status:", resp.Status)

	if resp.StatusCode != http.StatusNoContent {
		if resp.StatusCode == http.StatusInternalServerError {
			var dat map[string]string
			err = json.Unmarshal([]byte(jsonCreateForAll), &dat)
			user_id := GetUserID(dat["username"])
			if user_id != "" {
				t.Logf("Response code is %v,the username is exist.", resp.StatusCode)
			} else {
				t.Errorf("Response code is %v,harbor happened error.", resp.StatusCode)
			}

		} else {
			t.Errorf("Response code is %v", resp.StatusCode)
		}
	}

}

func CreateUser() string {
	token := GetLoginToken()

	r, err := http.NewRequest("POST", uri+"/webase/api/v2/user/", bytes.NewBuffer(jsonCreateForAll))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", string(token))

	timeout := time.Duration(10 * time.Second)
	client := &http.Client{
		Timeout:   timeout,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	resp, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var dat map[string]string
	err = json.Unmarshal([]byte(jsonCreateForAll), &dat)

	return dat["username"]

}

func GetUserID(username string) string {
	token := GetLoginToken()
	r, err := http.NewRequest("GET", uri+"/webase/api/v2/users/", nil)
	r.Header.Set("Authorization", string(token))

	timeout := time.Duration(10 * time.Second)
	client := &http.Client{
		Timeout:   timeout,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	resp, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var dat []map[string]string
	err = json.Unmarshal([]byte(body), &dat)

	var user_id string
	for _, data := range dat {
		for _, v := range data {
			if v == username {
				user_id = data["id"]
			}
		}
	}
	return user_id
}
func TestUserGet(t *testing.T) {
	token := GetLoginToken()
	username := CreateUser()
	user_id := GetUserID(username)

	r, err := http.NewRequest("GET", uri+"/webase/api/v2/users/"+string(user_id), nil)
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", string(token))

	timeout := time.Duration(10 * time.Second)
	client := &http.Client{
		Timeout:   timeout,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	resp, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("user get response Status:", resp.Status)
	// body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("response Body:", string(body))
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Response code is %v", resp.StatusCode)
	}
}

func TestUserUpdate(t *testing.T) {
	token := GetLoginToken()

	var jsonStr = []byte(`{"phone":"123456","email":"12@123"}`)
	username := CreateUser()
	user_id := GetUserID(username)

	r, err := http.NewRequest("PUT", uri+"/webase/api/v2/users/"+string(user_id), bytes.NewBuffer(jsonStr))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", string(token))

	timeout := time.Duration(10 * time.Second)
	client := &http.Client{
		Timeout:   timeout,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	resp, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("user update response Status:", resp.Status)

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Response code is %v", resp.StatusCode)
	}
}

func TestUserCurrent(t *testing.T) {
	token := GetLoginToken()

	r, err := http.NewRequest("GET", uri+"/webase/api/v2/user", nil)
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", string(token))

	timeout := time.Duration(10 * time.Second)
	client := &http.Client{
		Timeout:   timeout,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	resp, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("user current response Status:", resp.Status)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Response code is %v", resp.StatusCode)
	}
}

func TestUserUpdateCurrent(t *testing.T) {
	token := GetLoginToken()
	var jsonStr = []byte(`{"email":"12@4321"}`)
	r, err := http.NewRequest("PUT", uri+"/webase/api/v2/user", bytes.NewBuffer(jsonStr))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", string(token))

	timeout := time.Duration(10 * time.Second)
	client := &http.Client{
		Timeout:   timeout,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
	resp, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("user current response Status:", resp.Status)

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Response code is %v", resp.StatusCode)
	}
}
