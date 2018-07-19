package k8sclient

import (
	b64 "encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// InitKubeClient return new kubernetes Client
func InitKubeClient(kubeconfigFlags string) (*kubernetes.Clientset, error) {

	config, err := GetKubeConfig(kubeconfigFlags)
	if err != nil {
		return nil, err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to init kubernetes clientset : %s", err)
	}
	return clientset, nil
}

func GetKubeConfig(kubeconfigFlags string) (*restclient.Config, error) {
	var kubeconfig string
	if kubeconfigFlags == "" {
		kubeconfig = os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			return nil, errors.New("At least KUBECONFIG or --kubeconfig must be set")
		}
	} else {
		kubeconfig = kubeconfigFlags
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build config from kubeconfig : %s", err)
	}
	return config, err
}

func GetCertificateAuthorityData(config *restclient.Config) (string, error) {
	var caData []byte
	var err error
	if len(config.CAData) != 0 {
		caData = config.CAData
	} else {
		caData, err = ioutil.ReadFile(config.CAFile)
		if err != nil {
			return "", err
		}
	}
	return string(b64.StdEncoding.EncodeToString(caData)), nil
}
