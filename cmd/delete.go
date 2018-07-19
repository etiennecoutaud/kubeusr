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
	"fmt"

	cmderr "github.com/etiennecoutaud/kubeusr/cmd/error"
	k8sclient "github.com/etiennecoutaud/kubeusr/k8sclient"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete <user>",
	Short: "Delete Kubernetes user",
	Long: `Example:
  Delete user in namespace foo 
  $ kubeusr delete john.doe -n foo`,
	Run: deleteCmdFunc,
}

const DeleteUsrMsg = "User \"%s\" deleted\n"

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringP("namespace", "n", "default", "Namespace context")
}

func deleteCmdFunc(cmd *cobra.Command, args []string) {
	clientSet, err := k8sclient.InitKubeClient(cmd.Flag("kubeconfig").Value.String())
	if err != nil {
		cmderr.ExitWithError(1, err)
	}
	if len(args) != 1 {
		cmderr.ExitWithError(1, fmt.Errorf("%s", "Please provide a user id"))
	}
	ns := cmd.Flag("namespace").Value.String()
	sa, err := clientSet.CoreV1().ServiceAccounts(ns).Get(args[0], metav1.GetOptions{})
	if err != nil {
		cmderr.ExitWithError(1, err)
	}

	scope := sa.Labels["scope"]

	err = clientSet.CoreV1().ServiceAccounts(cmd.Flag("namespace").Value.String()).Delete(args[0], metav1.NewDeleteOptions(1))
	if err != nil {
		cmderr.ExitWithError(1, err)
	}

	if scope == "namespace" {
		err = clientSet.RbacV1().RoleBindings(ns).Delete(args[0]+"-binding-admin", metav1.NewDeleteOptions(1))
		if err != nil {
			cmderr.ExitWithError(1, err)
		}
	} else {
		err = clientSet.RbacV1().ClusterRoleBindings().Delete(args[0]+"-binding-cluster-admin", metav1.NewDeleteOptions(1))
		if err != nil {
			cmderr.ExitWithError(1, err)
		}
	}
	fmt.Printf(DeleteUsrMsg, args[0])
}
