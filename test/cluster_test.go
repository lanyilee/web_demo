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
	_ "webase-server/store/db"
)

var cluster_id = "b046c52b-65cf-4e51-a84b-1c4bc8cf14be"
var cluster_name = "开发测试集群"
var member_name = "fortest"
var namespace = "default"

func GetClusterMemberID(cluster_id string) string {
	token := GetLoginToken()
	r, _ := http.NewRequest("GET", uri+"/webase/api/v2/clusters/"+string(cluster_id)+"/members", nil)
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

	body, _ := ioutil.ReadAll(resp.Body)
	var dat []map[string]string
	err = json.Unmarshal([]byte(body), &dat)
	var member_id string
	for _, data := range dat {
		for _, v := range data {
			if v == member_name {
				member_id = data["member_id"]
			}
		}
	}
	return member_id
}

func GetNSMemberID(cluster_id string) string {
	token := GetLoginToken()
	r, _ := http.NewRequest("GET", uri+"/webase/api/v2/clusters/"+string(cluster_id)+"/namespaces/"+string(namespace)+"/members", nil)
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

	body, _ := ioutil.ReadAll(resp.Body)
	var dat []map[string]string
	err = json.Unmarshal([]byte(body), &dat)
	var member_id string
	for _, data := range dat {
		for _, v := range data {
			if v == member_name {
				member_id = data["member_id"]
			}
		}
	}
	return member_id
}

func GetClusterAndMemberID() (string, string, string) {
	token := GetLoginToken()
	r, _ := http.NewRequest("GET", uri+"/webase/api/v2/clusters/", nil)
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
	body, _ := ioutil.ReadAll(resp.Body)
	var dat []map[string]string
	err = json.Unmarshal([]byte(body), &dat)
	var cluster_id, cluster_member_id, NS_member_id string

	for _, data := range dat {
		for _, v := range data {
			if v == cluster_name {
				cluster_id = data["id"]
				cluster_member_id = GetClusterMemberID(cluster_id)
				NS_member_id = GetNSMemberID(cluster_id)
			}
		}
	}
	return cluster_id, cluster_member_id, NS_member_id
}

func TestClusterList(t *testing.T) {
	token := GetLoginToken()
	r, _ := http.NewRequest("GET", uri+"/webase/api/v2/clusters/", nil)
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
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Response code is %v", resp.StatusCode)
	}
}

func TestListMember(t *testing.T) {
	token := GetLoginToken()
	cluster_id, _, _ := GetClusterAndMemberID()
	r, _ := http.NewRequest("GET", uri+"/webase/api/v2/clusters/"+string(cluster_id)+"/members", nil)
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

	fmt.Println("response Status:", resp.Status)
	// body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("response Body:", string(body))
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Response code is %v", resp.StatusCode)
	}
}

func TestAddMember(t *testing.T) {
	token := GetLoginToken()
	cluster_id, _, _ := GetClusterAndMemberID()
	member_name := CreateUser()
	user_id := GetUserID(member_name)

	var jsonStr = []byte(`{"user_id":"` + user_id + `","cluster_role":"user"}`)
	r, _ := http.NewRequest("POST", uri+"/webase/api/v2/clusters/"+string(cluster_id)+"/member", bytes.NewBuffer(jsonStr))
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

	fmt.Println("response Status:", resp.Status)

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Response code is %v", resp.StatusCode)
	}
}

func TestDeleteMember(t *testing.T) {
	token := GetLoginToken()
	cluster_id, member_id, _ := GetClusterAndMemberID()

	r, _ := http.NewRequest("DELETE", uri+"/webase/api/v2/clusters/"+string(cluster_id)+"/members/"+string(member_id), nil)
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
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	fmt.Println("response Status:", resp.Status)

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Response code is %v", resp.StatusCode)
	}
}
