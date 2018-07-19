package util

import (
	v1 "k8s.io/api/core/v1"
)

//  NamespacesListToStringList
func NamespacesListToStringList(nsList *v1.NamespaceList) []string {
	var result []string
	for _, ns := range nsList.Items {
		result = append(result, ns.Name)
	}
	return result
}
