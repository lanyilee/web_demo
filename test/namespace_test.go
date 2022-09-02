package test

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestNSListMember(t *testing.T) {
	token := GetLoginToken()
	cluster_id, _, _ := GetClusterAndMemberID()
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

	fmt.Println("response Status:", resp.Status)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Response code is %v", resp.StatusCode)
	}
}

func TestNSAddMember(t *testing.T) {
	token := GetLoginToken()
	cluster_id, _, _ := GetClusterAndMemberID()
	member_name := CreateUser()
	user_id := GetUserID(member_name)

	var jsonStr = []byte(`{"user_id":"` + user_id + `","namespace_role":"reader"}`)

	r, _ := http.NewRequest("POST", uri+"/webase/api/v2/clusters/"+string(cluster_id)+"/namespaces/"+string(namespace)+"/member", bytes.NewBuffer(jsonStr))
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

func TestNSDeleteMember(t *testing.T) {
	token := GetLoginToken()
	cluster_id, _, member_id := GetClusterAndMemberID()

	r, _ := http.NewRequest("DELETE", uri+"/webase/api/v2/clusters/"+string(cluster_id)+"/namespaces/"+string(namespace)+"/members/"+string(member_id), nil)
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

	// delete user fortest
	// user_id := GetUserID("fortest")
	// r, err = http.NewRequest("DELETE", uri+"/webase/api/v2/users/"+string(user_id), nil)
	// r.Header.Set("Content-Type", "application/json")
	// r.Header.Set("Authorization", string(token))

	// timeout = time.Duration(10 * time.Second)
	// client = &http.Client{
	// 	Timeout:   timeout,
	// 	Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	// }
	// resp, err = client.Do(r)
	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()
	// fmt.Println("delete user response Status:", resp.Status)

	// if resp.StatusCode != http.StatusNoContent {
	// 	t.Errorf("Response code is %v", resp.StatusCode)
	// }

}

func TestNSGet(t *testing.T) {
	token := GetLoginToken()
	cluster_id, _, _ := GetClusterAndMemberID()
	r, _ := http.NewRequest("GET", uri+"/webase/api/v2/clusters/"+string(cluster_id)+"/namespaces", nil)
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
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Response code is %v", resp.StatusCode)
	}
}
