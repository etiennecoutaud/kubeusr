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
	"os"

	"path/filepath"

	"io/ioutil"

	cmderr "github.com/etiennecoutaud/kubeusr/cmd/error"
	k8sclient "github.com/etiennecoutaud/kubeusr/k8sclient"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	"k8s.io/api/rbac/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Represent a new user
type User struct {
	Name      string `yaml:"name"`
	Scope     string `yaml:"scope"`
	Namespace string `yaml:"namespace"`
}

type UserList struct {
	Items []User `yaml:"list"`
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create <user>",
	Short: "Create Kubernetes user",
	Long: `Example:
  Create simple user with namespace restricted rights in foo namespace 
  $ kubeusr create john.doe -s namespace -n foo
	
  Create user from file
  $ kubeusr create --from-file=list.yml`,

	Run: createCmdFunc,
}

const NewUsrMsgScoped = "User \"%s\" created with scoped rights in namespace \"%s\"\n"
const NewUsrMsgNonScoped = "User \"%s\" created with cluster-wide rights in namespace \"%s\"\n"

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringP("namespace", "n", "default", "Namespace context")
	createCmd.Flags().StringP("scope", "s", "namespace", "Scope rights possible values: [namespace,cluster]")
	createCmd.Flags().StringP("from-file", "f", "", "File with users list to create")
}

func createCmdFunc(cmd *cobra.Command, args []string) {
	clientSet, err := k8sclient.InitKubeClient(cmd.Flag("kubeconfig").Value.String())
	if err != nil {
		cmderr.ExitWithError(1, err)
	}

	if cmd.Flag("from-file").Value.String() != "" {
		filename, err := filepath.Abs(cmd.Flag("from-file").Value.String())
		if err != nil {
			cmderr.ExitWithError(1, err)
		}
		fileContent, err := ioutil.ReadFile(filename)
		if err != nil {
			cmderr.ExitWithError(1, err)
		}
		var usrList UserList
		err = yaml.Unmarshal(fileContent, &usrList)
		if err != nil {
			cmderr.ExitWithError(1, err)
		}

		exitCode := 0
		for _, usr := range usrList.Items {
			err := createNewUser(&usr, clientSet)
			if err != nil {
				fmt.Printf("Fail to create %s : %s\n", usr.Name, err)
				exitCode = 1
			}
		}
		os.Exit(exitCode)
	}

	if len(args) != 1 {
		cmderr.ExitWithError(1, fmt.Errorf("%s", "Please provide a user id"))
	}

	usr := &User{
		Name:      args[0],
		Namespace: cmd.Flag("namespace").Value.String(),
		Scope:     cmd.Flag("scope").Value.String(),
	}

	err = createNewUser(usr, clientSet)
	if err != nil {
		cmderr.ExitWithError(1, err)
	}
}

func createRoleBinding(usr *User, clientSet *kubernetes.Clientset) error {
	var err error
	if usr.Scope == "namespace" {
		rolebinding := &v1beta1.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: usr.Name + "-binding-admin",
			},
			Subjects: []v1beta1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      usr.Name,
					Namespace: usr.Namespace,
				},
			},
			RoleRef: v1beta1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     "admin",
			},
		}
		_, err = clientSet.RbacV1beta1().RoleBindings(usr.Namespace).Create(rolebinding)
	} else {
		clusterRolebinding := &v1beta1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: usr.Name + "-binding-cluster-admin",
			},
			Subjects: []v1beta1.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      usr.Name,
					Namespace: usr.Namespace,
				},
			},
			RoleRef: v1beta1.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     "cluster-admin",
			},
		}
		_, err = clientSet.RbacV1beta1().ClusterRoleBindings().Create(clusterRolebinding)
	}
	return err
}

func checkUsrValue(usr *User) error {
	if usr.Scope != "namespace" && usr.Scope != "cluster" {
		return fmt.Errorf("%s", "scope supported value : [namespace, cluster]")
	}
	return nil
}

func createNewUser(usr *User, clientSet *kubernetes.Clientset) error {

	err := checkUsrValue(usr)
	if err != nil {
		return err
	}

	newUsr := &v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: usr.Name,
			Labels: map[string]string{
				"type":  "user",
				"scope": usr.Scope},
		},
	}
	_, err = clientSet.CoreV1().ServiceAccounts(usr.Namespace).Create(newUsr)
	if err != nil {
		return err
	}

	err = createRoleBinding(usr, clientSet)
	if err != nil {
		err2 := clientSet.CoreV1().ServiceAccounts(usr.Namespace).Delete(usr.Name, metav1.NewDeleteOptions(1))
		if err2 != nil {
			fmt.Printf("Error: %s\n", err2)
		}
		return err
	}
	if usr.Scope == "namespace" {
		fmt.Printf(NewUsrMsgScoped, usr.Name, usr.Namespace)
	} else {
		fmt.Printf(NewUsrMsgNonScoped, usr.Name, usr.Namespace)
	}
	return nil
}
