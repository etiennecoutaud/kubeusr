// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	cmderr "github.com/etiennecoutaud/kubeusr/cmd/error"
	printer "github.com/etiennecoutaud/kubeusr/cmd/printer"
	k8sclient "github.com/etiennecoutaud/kubeusr/k8sclient"
	util "github.com/etiennecoutaud/kubeusr/util"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get or List Kubernetes user",
	Long: `Example:
  Get all users in cluster 
  $ kubeusr get
	
  Get all users in foo namespace
  $ kubeusr get -n foo`,
	Run: getCmdFunc,
}

const UserLabelSelector = "type=user"

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringP("namespace", "n", "", "Namespace context")
}

func getCmdFunc(cmd *cobra.Command, args []string) {
	clientSet, err := k8sclient.InitKubeClient(cmd.Flag("kubeconfig").Value.String())
	if err != nil {
		cmderr.ExitWithError(1, err)
	}

	var namespacesStr []string

	if cmd.Flag("namespace").Value.String() != "" {
		namespacesStr = []string{cmd.Flag("namespace").Value.String()}
	} else {
		namespaces, err := clientSet.CoreV1().Namespaces().List(metav1.ListOptions{})
		if err != nil {
			cmderr.ExitWithError(1, err)
		}
		namespacesStr = util.NamespacesListToStringList(namespaces)
	}

	var allUsr []v1.ServiceAccount
	for _, ns := range namespacesStr {
		usr, err := clientSet.CoreV1().ServiceAccounts(ns).List(metav1.ListOptions{LabelSelector: UserLabelSelector})
		if err != nil {
			cmderr.ExitWithError(1, err)
		}
		allUsr = append(allUsr, usr.Items...)
	}

	printer.PrintUsers(allUsr)

}
