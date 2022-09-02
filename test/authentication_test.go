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
// var uri = "https://127.0.0.1:8888"
var (
	jsonCreateUser = []byte(`{"username":"fortest","password":"Qwer1234","name":"fortest","role":"k8sclusteruser","email":"111"}`)
	jsonAddUser    = []byte(`{"username":"fortest1","password":"Qwer1234","name":"fortest","role":"k8sclusteruser","email":"111"}`)
	jsonCreateRole = []byte(`{"role_name":"projectrole","rules":[
		{
			"resource_name": "menu_安全管理",
			"resource_act":["get"],
			"resource_id":"0001-1010-0000-0000",
			"resource_parent_id":"0001-0000-0000-0000"
		  },
		   {
			"resource_name": "menu_安全管理_资源项目",
			"resource_act":["get"],
			"resource_id":"0001-1010-0100-0000",
			"resource_parent_id":"0001-1010-0000-0000"
		  },
		  {
			"resource_name": "projects",
			"resource_act":["get","create","delete"],
			"resource_id":"0010-0000-0000-0000",
			"resource_parent_id":"0000-0000-0000-0000"
		  },
		  {
			"resource_name": "projects_members",
			"resource_act":["get","create"],
			"resource_id":"0010-0001-0000-0000",
			"resource_parent_id":"0010-0000-0000-0000"
		  }
	]}`)
	jsonCreateRoleBinding = []byte(`{
		"name":"fortest1",
		"rolebinding":
			{
			  "role_name": ["k8sclusteruser","projectrole"]
			} 
	}`)
)

func requestAPI(r *http.Request) *http.Response {
	token := GetLoginToken()
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
	return resp
}

func AddUser() string {
	r, _ := http.NewRequest("POST", uri+"/webase/api/v2/user/", bytes.NewBuffer(jsonAddUser))
	resp := requestAPI(r)
	defer resp.Body.Close()
	var dat map[string]string
	json.Unmarshal([]byte(jsonCreateForAll), &dat)
	return dat["username"]
}

func CreateAndDeleteUser() string {
	username := AddUser()
	user_id := GetUserID(username)
	r, _ := http.NewRequest("DELETE", uri+"/webase/api/v2/users/"+string(user_id), nil)
	resp := requestAPI(r)
	defer resp.Body.Close()
	var dat map[string]string
	json.Unmarshal([]byte(jsonCreateForAll), &dat)
	return dat["username"]
}

func TestUserAdd(t *testing.T) {
	r, _ := http.NewRequest("POST", uri+"/webase/api/v2/user", bytes.NewBuffer(jsonCreateUser))
	resp := requestAPI(r)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		if resp.StatusCode == http.StatusInternalServerError {
			t.Errorf("Response message is %v", resp.Body)
			var dat map[string]string
			_ = json.Unmarshal([]byte(jsonCreateForAll), &dat)
			user_id := GetUserID(dat["username"])
			if user_id != "" {
				t.Logf("Response code is %v,the username is exist.", resp.StatusCode)
			} else {
				t.Errorf("Response code is %v, happened error.", resp.StatusCode)
			}

		} else {
			t.Errorf("Response code is %v", resp.StatusCode)
		}
	}
}

func TestUserDelete(t *testing.T) {
	user_id := "f9f3b09d-e09f-4944-b890-8f2034e4c8e2"
	r, _ := http.NewRequest("DELETE", uri+"/webase/api/v2/users/"+string(user_id), nil)
	resp := requestAPI(r)
	defer resp.Body.Close()
	fmt.Println("user delete response Status:", resp.Status)
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Response code is %v", resp.StatusCode)
	}
}

func GetRoleName(rolename string) string {
	r, _ := http.NewRequest("GET", uri+"/webase/api/v2/roles", nil)
	resp := requestAPI(r)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var dat []map[string]string
	json.Unmarshal([]byte(body), &dat)

	var role_name string
	for _, data := range dat {
		for _, v := range data {
			if v == rolename {
				role_name = data["role_name"]
			}
		}
	}
	return role_name
}

func TestRoleAdd(t *testing.T) {
	r, _ := http.NewRequest("POST", uri+"/webase/api/v2/roles", bytes.NewBuffer(jsonCreateRole))
	resp := requestAPI(r)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		if resp.StatusCode == http.StatusInternalServerError {
			t.Errorf("Response message is %v", resp.Body)
			var dat map[string]string
			_ = json.Unmarshal([]byte(jsonCreateForAll), &dat)
			role_name := GetRoleName(dat["role_name"])
			if role_name != "" {
				t.Logf("Response code is %v,the role_name is exist.", resp.StatusCode)
			} else {
				t.Errorf("Response code is %v, happened error.", resp.StatusCode)
			}

		} else {
			t.Errorf("Response code is %v", resp.StatusCode)
		}
	}
}

func TestRoleDelete(t *testing.T) {
	role_name := "projectrole"
	r, _ := http.NewRequest("DELETE", uri+"/webase/api/v2/roles/"+string(role_name), nil)
	resp := requestAPI(r)
	defer resp.Body.Close()
	fmt.Println("user delete response Status:", resp.Status)
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Response code is %v", resp.StatusCode)
	}
}

func GetRolebindingID(username string) string {
	r, _ := http.NewRequest("GET", uri+"/webase/api/v2/rolebindings", nil)
	resp := requestAPI(r)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var dat []map[string]string
	json.Unmarshal([]byte(body), &dat)

	var binding_id string
	for _, data := range dat {
		for _, v := range data {
			if v == username {
				binding_id = data["id"]
			}
		}
	}
	return binding_id
}

func TestRoleBindingAdd(t *testing.T) {
	// 添加用户，然后将用户与角色绑定
	_ = AddUser()
	r, _ := http.NewRequest("POST", uri+"/webase/api/v2/rolebindings", bytes.NewBuffer(jsonCreateRoleBinding))
	resp := requestAPI(r)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		if resp.StatusCode == http.StatusInternalServerError {
			t.Errorf("Response message is %v", resp.Body)
			var dat map[string]string
			_ = json.Unmarshal([]byte(jsonCreateForAll), &dat)
			user_id := GetUserID(dat["name"])
			if user_id != "" {
				t.Logf("Response code is %v,the username is exist.", resp.StatusCode)
			} else {
				t.Errorf("Response code is %v, happened error.", resp.StatusCode)
			}
		} else {
			t.Errorf("Response code is %v", resp.StatusCode)
		}
	}

}

func TestRoleBindingDelete(t *testing.T) {
	// 删除绑定，删除用户
	id := GetRolebindingID("fortest1")
	fmt.Println("binding id:", id)
	// id := "6247893f-a922-409c-baa6-e4727ccaac17"
	r, _ := http.NewRequest("DELETE", uri+"/webase/api/v2/rolebindings/"+string(id), nil)
	resp := requestAPI(r)
	defer resp.Body.Close()
	fmt.Println("user delete response Status:", resp.Status)
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Response code is %v", resp.StatusCode)
	}
	user_id := GetUserID("fortest1")
	r1, _ := http.NewRequest("DELETE", uri+"/webase/api/v2/users/"+string(user_id), nil)
	resp1 := requestAPI(r1)
	defer resp1.Body.Close()
}

func TestMenuPermission(t *testing.T) {
	r, _ := http.NewRequest("GET", uri+"/webase/api/v2/rolebindings/menu/f9f3b09d-e09f-4944-b890-8f2034e4c8e2", nil)
	resp := requestAPI(r)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Response code is %v", resp.StatusCode)
	}
}
