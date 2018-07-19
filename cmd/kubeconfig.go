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
	yaml "gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// kubeconfigCmd represents the kubeconfig command
var kubeconfigCmd = &cobra.Command{
	Use:   "kubeconfig <user>",
	Short: "Get user kubeconfig",
	Long: `Example:
  Get john.doe kubeconfig in foo namespace:
  $ kubeusr kubeconfig john.doe -n foo`,
	Run: kubeconfigCmdFunc,
}

func init() {
	rootCmd.AddCommand(kubeconfigCmd)
	kubeconfigCmd.Flags().StringP("namespace", "n", "default", "Namespace context")
}

func kubeconfigCmdFunc(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		cmderr.ExitWithError(1, fmt.Errorf("%s", "Please provide a user id"))
	}
	clientSet, err := k8sclient.InitKubeClient(cmd.Flag("kubeconfig").Value.String())
	if err != nil {
		cmderr.ExitWithError(1, err)
	}
	config, err := k8sclient.GetKubeConfig(cmd.Flag("kubeconfig").Value.String())
	if err != nil {
		cmderr.ExitWithError(1, err)
	}
	caData, err := k8sclient.GetCertificateAuthorityData(config)
	if err != nil {
		cmderr.ExitWithError(1, err)
	}
	token, err := getTokenAccessForServiceAccount(args[0], cmd.Flag("namespace").Value.String(), clientSet)
	if err != nil {
		cmderr.ExitWithError(1, err)
	}
	kubeconfigSruct := buildKubeConfig(caData, token, args[0], cmd.Flag("namespace").Value.String(), config.Host, "k8s-cluster")
	kubeconfigYaml, err := yaml.Marshal(kubeconfigSruct)
	if err != nil {
		cmderr.ExitWithError(1, err)
	}
	fmt.Printf("%s", kubeconfigYaml)
}

func getTokenAccessForServiceAccount(usrName string, ns string, clientSet *kubernetes.Clientset) (string, error) {
	sa, err := clientSet.CoreV1().ServiceAccounts(ns).Get(usrName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	secret, err := clientSet.CoreV1().Secrets(ns).Get(sa.Secrets[0].Name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return string(secret.Data["token"]), nil
}

func buildKubeConfig(caData string, token string, usrName string, ns string, host string, serverName string) *KubeConfig {
	return &KubeConfig{
		Kind:           "Config",
		ApiVersion:     "v1",
		CurrentContext: usrName + "-ctx",
		Clusters: []*KubectlClusterWithName{
			{
				Name: serverName,
				Cluster: KubectlCluster{
					Server: host,
					CertificateAuthorityData: caData,
				},
			},
		},
		Users: []*KubectlUserWithName{
			{
				Name: usrName,
				User: KubectlUser{
					Token: token,
				},
			},
		},
		Contexts: []*KubectlContextWithName{
			{
				Name: usrName + "-ctx",
				Context: KubectlContext{
					Cluster:   serverName,
					User:      usrName,
					Namespace: ns,
				},
			},
		},
	}
}
