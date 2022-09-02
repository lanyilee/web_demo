package rbac

import (
	"webase-server/models"
)

var actions = []string{"list", "get", "create", "update", "delete", "watch", "patch"}

var globalResource = []string{"userbases", "resourceroles", "projects", "clusters", "hosts"} //, "dockerimages", "helmcharts"

var projectResource = []string{"projectmembers"}

var clusterResource = []string{"clustermembers", "nodes", "namespaces", "storageClasses", "persistentvolumes"}

var namespaceResource = []string{
	"namespacemembers",
	"pods", "endpoints", "replicasets", "replicationcontrollers",
	"deployments", "statefulsets", "daemonsets", "jobs", "cronjobs",
	"services", "configmaps", "secrets", "ingresses",
	"serviceaccounts", "role",
	"persistentvolumeclaims",
	"resourcequotas", "limitranges",
}

//ListResource 资源分类
func ListResource(kind string) []models.Resource {
	resources := make([]models.Resource, 0)
	switch kind {
	case "global":
		for _, rs := range globalResource {
			resources = append(resources, models.Resource{
				Name:    rs,
				Actions: actions,
			})
		}
	case "cluster":
		for _, rs := range clusterResource {
			resources = append(resources, models.Resource{
				Name:    rs,
				Actions: actions,
			})
		}
	case "project":
		for _, rs := range projectResource {
			resources = append(resources, models.Resource{
				Name:    rs,
				Actions: actions,
			})
		}
		fallthrough
	case "namespace":
		for _, rs := range namespaceResource {
			resources = append(resources, models.Resource{
				Name:    rs,
				Actions: actions,
			})
		}
	}
	return resources
}
