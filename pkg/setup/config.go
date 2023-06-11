package setup

import (
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// KubernetesConfigs has all the fields to expose Kubernetes configs
type KubernetesConfigs struct {
	Configs map[string]*rest.Config
}

// Region has all the fields to expose regions
type Region struct {
	Code string `json:"code" yaml:"code"`
	Name string `json:"name" yaml:"name"`
}

// Config has all the fields to expose config
type Config struct {
	RegionFile       string `envconfig:"REGIONS_FILE" default:"config/regions/all.yaml"`
	KubeconfigFolder string `envconfig:"KUBECONFIG_FOLDER" default:"config/kubeconfig/"`
}

// ListRegions returns a list of regions
func ListRegions(log *logrus.Entry, regionFile string) []Region {
	yamlFile, err := ioutil.ReadFile(regionFile)
	if err != nil {
		log.Error(err)
	}

	data := map[string][]Region{}
	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		log.Error(err)
	}

	return data["regions"]
}

// LoadConfig loads the configuration from the environment variables
// and returns a Config struct and an error
func LoadConfig() (*Config, error) {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		return nil, err
	}

	requiredEnvVars := map[string]string{
		"RegionFile":       "missing REGIONS_FILE environment variable",
		"KubeconfigFolder": "missing KUBECONFIG_FOLDER environment variable",
	}

	for envVar, errMsg := range requiredEnvVars {
		value := reflect.ValueOf(config).FieldByName(envVar).String()
		if value == "" {
			return nil, errors.New(errMsg)
		}
	}

	return &config, nil
}

// BuildKubernetesConfigs builds the Kubernetes configs
func BuildKubernetesConfigs(log *logrus.Entry, config *Config) KubernetesConfigs {
	regions := ListRegions(log, config.RegionFile)

	configs := KubernetesConfigs{
		Configs: make(map[string]*rest.Config),
	}
	for _, region := range regions {
		cfg, err := clientcmd.BuildConfigFromFlags("", fmt.Sprintf("%s%s.yml", config.KubeconfigFolder, region.Code))
		if err != nil {
			log.Errorf("Error building kubeconfig: %s", err.Error())
		}

		configs.Configs[region.Code] = cfg
	}

	return configs
}
