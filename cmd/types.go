package cmd

type KubeConfig struct {
	Kind           string                    `yaml:"kind"`
	ApiVersion     string                    `yaml:"apiVersion"`
	CurrentContext string                    `yaml:"current-context"`
	Clusters       []*KubectlClusterWithName `yaml:"clusters"`
	Contexts       []*KubectlContextWithName `yaml:"contexts"`
	Users          []*KubectlUserWithName    `yaml:"users"`
}

type KubectlClusterWithName struct {
	Name    string         `yaml:"name"`
	Cluster KubectlCluster `yaml:"cluster"`
}

type KubectlCluster struct {
	Server                   string `yaml:"server,omitempty"`
	CertificateAuthorityData string `yaml:"certificate-authority-data,omitempty"`
}

type KubectlContextWithName struct {
	Name    string         `yaml:"name"`
	Context KubectlContext `yaml:"context"`
}

type KubectlContext struct {
	Cluster   string `yaml:"cluster"`
	User      string `yaml:"user"`
	Namespace string `yaml:"namespace"`
}

type KubectlUserWithName struct {
	Name string      `yaml:"name"`
	User KubectlUser `yaml:"user"`
}

type KubectlUser struct {
	Token string `yaml:"token,omitempty"`
}
