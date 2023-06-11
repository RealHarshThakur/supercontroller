package controllers

import (
	"fmt"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func resourceArgList(groupVersionMap map[string][]*v1.APIResourceList) []string {
	argList := make([]string, 0)
	for groupVersion, apiResourceLists := range groupVersionMap {
		for _, apiResourceList := range apiResourceLists {
			for _, apiResource := range apiResourceList.APIResources {
				groupVersionSplit := strings.Split(groupVersion, "/")
				group := groupVersionSplit[0]
				version := groupVersionSplit[1]
				argList = append(argList, fmt.Sprintf("%s.%s.%s", apiResource.Name, version, group))
			}
		}
	}
	return argList
}

func mapAPIResourcesByGroup(apiResources []*v1.APIResourceList, group string) map[string][]*v1.APIResourceList {
	groupVersionMap := make(map[string][]*v1.APIResourceList, 0)
	for _, apiResourceList := range apiResources {
		if strings.Contains(apiResourceList.GroupVersion, group) {
			groupVersionMap[apiResourceList.GroupVersion] = append(groupVersionMap[apiResourceList.GroupVersion], apiResourceList)
		}
	}

	return groupVersionMap

}
