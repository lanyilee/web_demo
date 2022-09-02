package rbac

import (
	"fmt"
	"go.uber.org/zap"
	"strings"
	"webase-server/models"

	casbin "github.com/casbin/casbin/v2"
)

var filename = map[string]string{
	"model":  "",
	"policy": "",
}

// func GetCasbinFile() map[string]string {
// 	switch os := runtime.GOOS; os {
// 	case "windows":
// 		// zap.S().Info("OS windows.")
// 		filename["model"] = "../../conf/casbin/model.conf"
// 		filename["policy"] = "../../conf/casbin/policy.csv"
// 	case "linux":
// 		zap.S().Info("Linux.")
// 		filename["model"] = "/etc/webase/casbin/model.conf"
// 		filename["policy"] = "/etc/webase/casbin/policy.csv"
// 	default:
// 		zap.S().Info("os is: ", os)
// 		filename["model"] = "/etc/webase/casbin/model.conf"
// 		filename["policy"] = "/etc/webase/casbin/policy.csv"
// 	}
// 	return filename
// }

// func TruncCasbinPolicy() {
// 	filename := GetCasbinFile()
// 	f, _ := os.OpenFile(filename["policy"], os.O_WRONLY|os.O_TRUNC, 0600)
// 	defer f.Close()
// }

func AddAdminPolicy() error {
	e := Enf
	role_name := "role:admin"
	_, err := e.AddPolicy(role_name, "*", "*")
	if err != nil {
		zap.S().Error("addPolicy failed, err: ", err)
		return err
	}
	_, err = e.AddGroupingPolicy("admin", role_name)
	if err != nil {
		zap.S().Error("addGroupingPolicy failed, err: ", err)
		return err
	}
	// if err := e.SavePolicy(); err != nil {
	// 	zap.S().Error("SavePolicy failed1, err: ", err)
	// 	return err
	// }
	return nil
}

// func NewCasbinPolicy(store *models.Store) {
// 	AddAdminPolicy()
// 	ClusterPolicy(store)
// 	NamespacePolicy(store)
// 	ProjectPolicy(store)
// 	BindingPolicy(store)
// }

type CasbinPolicy struct {
	store *models.Store
}

func (s *CasbinPolicy) ClusterPolicy(store *models.Store) {
	binding_cluster_members, _ := store.BindingMember.ListClusterMember(&models.BindingClusterMember{})
	// fmt.Println(binding_cluster_members)
	for _, bcm := range binding_cluster_members {
		username := bcm.UserName
		if username == "admin" {
			continue
		}
		rolename := bcm.RoleName
		cluster_id := bcm.ClusterID
		rules, err := store.RoleTemplate.GetRuleByRoleName(rolename)
		// fmt.Println("rules:", rules)
		if err != nil {
			zap.S().Error("get rule by role name err: ", err)
			return
		}
		for _, rule := range rules {
			// 排除需要通过成员方式绑定的资源权限
			if strings.HasPrefix(rule.Name, "*") {
				continue
			}
			if strings.HasPrefix(rule.Name, "harbor") {
				err = store.Tree.AddCasbinPolicy(rule, false, "", "")
				if err != nil {
					zap.S().Error("create casbin policy1 err: ", err)
					return
				}
				err = CreateCasbinGroupingPolicy(rolename+":"+cluster_id, rule.Name)
				if err != nil {
					zap.S().Error("create casbin group err: ", err)
					return
				}
			} else if strings.HasPrefix(rule.Name, "menu") {
				continue
			} else if strings.HasPrefix(rule.Name, "k8sclusters") {
				// 创建成员角色绑定后，在casbin中创建p和g策略
				// fmt.Println("rule1 issss:")
				err = store.Tree.AddCasbinPolicy(rule, true, bcm.ClusterID, "")
				if err != nil {
					zap.S().Error("create casbin policy2 err: ", err)
					return
				}
				err = CreateCasbinGroupingPolicy(rolename+":"+cluster_id, rule.Name+":"+bcm.ClusterID)
				if err != nil {
					zap.S().Error("create casbin group err: ", err)
					return
				}
			} else if strings.HasPrefix(rule.Name, "projects") {
				continue
			} else {
				break
			}
		}
		err = CreateCasbinGroupingPolicy(username, rolename+":"+cluster_id)
		if err != nil {
			zap.S().Error("create casbin group err: ", err)
			return
		}
	}
}

func (s *CasbinPolicy) NamespacePolicy(store *models.Store) {
	binding_ns_members, _ := store.BindingMember.ListNamespaceMember(&models.BindingNamespaceMember{})
	// fmt.Println(binding_ns_members)
	for _, bcm := range binding_ns_members {
		username := bcm.UserName
		if username == "admin" {
			continue
		}
		rolename := bcm.RoleName
		cluster_id := bcm.ClusterID
		namespace := bcm.Namespace
		rules, err := store.RoleTemplate.GetRuleByRoleName(rolename)
		// fmt.Println("rules:", rules)
		if err != nil {
			zap.S().Error("get rule by role name err: ", err)
			return
		}
		for _, rule := range rules {
			// fmt.Println("rule2 issss:")
			// 排除需要通过成员方式绑定的资源权限
			if strings.HasPrefix(rule.Name, "*") {
				continue
			}
			if strings.HasPrefix(rule.Name, "harbor") {
				err = store.Tree.AddCasbinPolicy(rule, false, "", "")
				if err != nil {
					zap.S().Error("create casbin policy3 err: ", err)
					return
				}
				err = CreateCasbinGroupingPolicy(rolename+":"+cluster_id+"-"+namespace, rule.Name)
				if err != nil {
					zap.S().Error("create casbin group err: ", err)
					return
				}
			} else if strings.HasPrefix(rule.Name, "menu") {
				continue
			} else if strings.HasPrefix(rule.Name, "k8sclusters") {
				// 创建成员角色绑定后，在casbin中创建p和g策略
				// fmt.Println("rule issss:", rule)
				err = store.Tree.AddCasbinPolicy(rule, true, bcm.ClusterID, bcm.Namespace)
				if err != nil {
					zap.S().Error("create casbin policy4 err: ", err)
					return
				}
				err = CreateCasbinGroupingPolicy(rolename+":"+cluster_id+"-"+namespace, rule.Name+":"+bcm.ClusterID+"-"+bcm.Namespace)
				if err != nil {
					zap.S().Error("create casbin group err: ", err)
					return
				}
				err = CreateCasbinGroupingPolicy(rolename+":"+cluster_id+"-"+namespace, rule.Name+":"+bcm.ClusterID)
				if err != nil {
					zap.S().Error("create casbin group err: ", err)
					return
				}
			} else if strings.HasPrefix(rule.Name, "projects") {
				continue
			} else {
				break
			}
		}
		err = CreateCasbinGroupingPolicy(username, rolename+":"+cluster_id+"-"+namespace)
		if err != nil {
			zap.S().Error("create casbin group err: ", err)
			return
		}
	}
}

type projectSpec struct {
	models.Project
	Namespaces []projectNamespace `json:"namespaces"`
}

type projectNamespace struct {
	Cluster string   `json:"cluster_id"`
	Names   []string `json:"names"`
}

func (s *CasbinPolicy) ProjectPolicy(store *models.Store) {
	binding_pro_members, _ := store.Project.FindProjectMember(&models.BindingProjectMember{})
	// fmt.Println(binding_pro_members)
	for _, bpm := range binding_pro_members {
		username := bpm.UserName
		if username == "admin" {
			continue
		}
		rolename := bpm.RoleName
		project_id := bpm.ProjectID
		rules, err := store.RoleTemplate.GetRuleByRoleName(rolename)
		// fmt.Println("rules:", rules)
		p, _, err := store.Project.Get(&models.Project{ID: bpm.ProjectID})
		if err != nil {
			return
		}
		spec := projectSpec{
			Project:    *p,
			Namespaces: make([]projectNamespace, 0),
		}
		pbs, err := store.Project.FindBind(&models.ProjectResourceBind{ProjectID: bpm.ProjectID})
		if err != nil {
			return
		}
		for _, pb := range pbs {
			if pb.Resource == "namespace" {
				exist := false
				for i, ns := range spec.Namespaces {
					if ns.Cluster == pb.Provider {
						spec.Namespaces[i].Names = append(spec.Namespaces[i].Names, pb.ResourceName)
						exist = true
					}
				}
				if !exist {
					spec.Namespaces = append(spec.Namespaces, projectNamespace{
						Cluster: pb.Provider,
						Names:   []string{pb.ResourceName},
					})
				}
			}
		}
		for _, rule := range rules {
			// fmt.Println("rule3 issss:")
			for _, ns := range spec.Namespaces {
				for _, name := range ns.Names {
					if strings.HasPrefix(rule.Name, "*") {
						continue
					}
					if strings.HasPrefix(rule.Name, "k8sclusters_resources_namespaces") || strings.HasPrefix(rule.Name, "k8sclusters.") {
						err = store.Tree.AddCasbinPolicy(rule, true, ns.Cluster, name)
						if err != nil {
							zap.S().Error("create casbin policy5 err: ", err)
							return
						}
						err = CreateCasbinGroupingPolicy(rolename+":"+project_id, rule.Name+":"+ns.Cluster+"-"+name)
						if err != nil {
							zap.S().Error("create casbin group err: ", err)
							return
						}
						err = CreateCasbinGroupingPolicy(rolename+":"+project_id, rule.Name+":"+ns.Cluster)
						if err != nil {
							zap.S().Error("create casbin group err: ", err)
							return
						}
						err = store.Tree.AddCasbinPolicy(rule, false, "", "")
						if err != nil {
							zap.S().Error("create casbin policy6 err: ", err)
							return
						}
						err = CreateCasbinGroupingPolicy(rolename+":"+project_id, rule.Name)
						if err != nil {
							zap.S().Error("create casbin group err: ", err)
							return
						}
					}
					if strings.HasPrefix(rule.Name, "projects") {
						err = store.Tree.AddCasbinPolicy(rule, false, "", "")
						if err != nil {
							zap.S().Error("create casbin policy7 err: ", err)
							return
						}
						err = CreateCasbinGroupingPolicy(rolename+":"+project_id, rule.Name)
						if err != nil {
							zap.S().Error("create casbin group err: ", err)
							return
						}
					} else {
						continue
					}
				}
			}
		}
		err = CreateCasbinGroupingPolicy(username, rolename+":"+project_id)
		if err != nil {
			zap.S().Error("create casbin group err: ", err)
			return
		}
	}
}

func (s *CasbinPolicy) BindingPolicy(store *models.Store) {
	binding_members, _ := store.Binding.FindBinding(&models.Binding{})
	for _, bm := range binding_members {
		username := bm.UserName
		if username == "admin" {
			continue
		}
		rolename := bm.RoleName
		// for _, rolename := range bm.RoleName {
		// 查询绑定的每个角色对应的权限，一个角色对应多个权限
		rules, err := store.NodeRole.GetRuleByRoleName(rolename)
		if err != nil {
			zap.S().Error("get rule by role name err: ", err)
			return
		}
		for _, rule := range rules {
			// fmt.Println("rule4 issss:")
			// 排除需要通过成员方式绑定的资源权限
			if strings.HasPrefix(rule.Name, "menu") {
				continue
			} else {
				// 创建角色绑定后，在casbin中创建p和g策略
				err = store.Tree.CreateCasbinPolicy(rule, false, "", "")
				if err != nil {
					zap.S().Error("create casbin policy8 err: ", err)
					return
				}
				err = CreateCasbinGroupingPolicy(rolename, rule.Name)
				if err != nil {
					zap.S().Error("create casbin group err: ", err)
					return
				}
			}
		}
		err = CreateCasbinGroupingPolicy(username, rolename)
		if err != nil {
			zap.S().Error("create casbin group err: ", err)
			return
		}
		// }
	}
}

var entity *casbin.Enforcer

// func ClientCasbin() (*casbin.Enforcer, error) {
// 	var err error
// 	if !(entity == nil) {
// 		// fmt.Println("I am already init")
// 		return entity, nil
// 	}
// 	fmt.Println("lazy init")
// 	file_map := GetCasbinFile()
// 	entity, err = casbin.NewEnforcer(file_map["model"], file_map["policy"])
// 	// a := mysqladapter.NewDBAdapter("mysql", "root:root123@tcp(188.8.5.1:33060)/?charset=utf8")
// 	// entity = casbin.NewEnforcer(file_map["model"], a)
// 	if err != nil {
// 		zap.S().Error("errors")
// 		return nil, err
// 	}
// 	return entity, nil
// }

func CreateCasbinPolicy(rule models.Rules) error {
	e := Enf
	// e, _ := UseXormForCasbinClient()
	act := strings.Join(rule.ResourceAct, "|")
	fmt.Println("actions is :", act)
	// _, err := e.AddPolicy(rule.RoleName, rule.ResourceName, rule.ResourceObj, act)
	_, err := e.AddPolicy(rule.RoleName, rule.ResourceObj, act)
	if err != nil {
		zap.S().Error("addPolicy failed, err: ", err)
		return err
	}
	// if err := e.SavePolicy(); err != nil {
	// 	zap.S().Error("SavePolicy failed2, err: ", err)
	// 	return err
	// }
	return nil
}

func DeleteCasbinPolicy(role_name string) error {
	e := Enf
	// e, _ := UseXormForCasbinClient()
	_, err := e.DeleteRole(role_name)
	if err != nil {
		zap.S().Error("delete role failed, err: ", err)
		return err
	}
	// if err := e.SavePolicy(); err != nil {
	// 	zap.S().Error("SavePolicy failed3, err: ", err)
	// 	return err
	// }
	return nil
}

// 获取用户权限
func GetPermissionsForUser(username string) ([][]string, error) {
	e := Enf
	// e, _ := UseXormForCasbinClient()
	permissions, err := e.GetImplicitPermissionsForUser(username)
	if err != nil {
		zap.S().Error("delete role failed, err: ", err)
		return nil, err
	}
	return permissions, nil
}

// 鉴权
func HasPermissionsForUser(username, url, action string) bool {
	e := Enf
	err := e.LoadPolicy()
	if err != nil {
		zap.L().Error(err.Error())
		return false
	}
	zap.S().Info("enforce client is: ", e)
	zap.S().Info("UserName: ", username)
	zap.S().Info("action: ", action)
	zap.S().Info("url: ", url)
	// has, err := e.Enforce(username, "", url, action)
	has, err := e.Enforce(username, url, action)
	if err != nil {
		zap.S().Info("Enforce failed, err: ", err)
	}
	// e.SavePolicy()
	return has
}

func ReplaceString(source string) string {
	des := strings.Replace(source, "_", "-", -1)
	return des
}

func CreateCasbinGroupingPolicy(username, role_name string) error {
	e := Enf
	username = ReplaceString(username)
	role_name = ReplaceString(role_name)
	_, err := e.AddGroupingPolicy(username, role_name)
	if err != nil {
		zap.S().Error("addGroupingPolicy failed, err: ", err)
		return err
	}
	// if err := e.SavePolicy(); err != nil {
	// 	zap.S().Error("SavePolicy failed4, err: ", err)
	// 	return err
	// }
	return nil
}

func DeleteCasbinGroupingPolicy(username, role_name string) error {
	e := Enf
	username = ReplaceString(username)
	role_name = ReplaceString(role_name)
	_, err := e.DeleteRoleForUser(username, role_name)
	if err != nil {
		zap.S().Error("delete role for user failed, err: ", err)
		return err
	}
	// if err := e.SavePolicy(); err != nil {
	// 	zap.S().Error("SavePolicy failed5, err: ", err)
	// 	return err
	// }
	return nil
}

func GetCasbinGroupingPolicy(role_name string) bool {
	e := Enf
	role_name = ReplaceString(role_name)
	casbin_groups := e.GetFilteredGroupingPolicy(1, role_name)

	if len(casbin_groups) > 0 {
		return true
	}
	return false
}
