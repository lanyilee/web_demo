package rbac

import (
	"fmt"
	"testing"
)

func TestHasPermissionsForUser(t *testing.T) {
	b := HasPermissionsForUser("user1", "/webase/api/v2/k8sclusters", "get")
	fmt.Println(b)
	// e, _ := ClientCasbin()
	// has, err := e.Enforce("user1", "k8sclusters", "/webase/api/v2/k8sclusters", "get")
	// if err != nil {
	// 	log.Println("Enforce failed, err: ", err)
	// }
	// fmt.Println("has", has)
}

func TestNewCasbinPolicy(t *testing.T) {
	fmt.Println("has")
}

func TestUseXormForCasbinClient(t *testing.T) {
	fmt.Println("has")
	// e, _ := UseBeegoXormForCasbinClient()
	// fmt.Println("has", e)
}
