package main

import (
	"fmt"
	"os"

	"sigs.k8s.io/kustomize/kyaml/fn/framework"
	"sigs.k8s.io/kustomize/kyaml/fn/framework/command"
	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

type Appspec struct {
	Image  string `yaml:"image" json:"image"`
	Port   int    `yaml:"port" json:"port"`
	Domain string `yaml:"value" json:"value"`
}

type App struct {
	Metadata struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec Appspec `yaml:"spec" json:"spec"`
}

// Generate Deployment yaml
func generateDeployment(name string, image string) (*yaml.RNode, error) {
	d, err := yaml.Parse(fmt.Sprintf(`apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: %s
  name: %s
spec:
  replicas: 1
  selector:
    matchLabels:
      app: %s
  template:
    metadata:
      labels:
        app: %s
    spec:
      containers:
      - image: %s
        name: %s
`, name, name, name, name, image, image))
	if err != nil {
		return nil, err
	}
	return d, nil
}

// Generate Service yaml
func generateService(name string, sourcePort int, targetPort int) (*yaml.RNode, error) {
	serviceName := name + "-svc"
	svc, err := yaml.Parse(fmt.Sprintf(`apiVersion: v1
kind: Service
metadata:
  labels:
    app: %s
  name: %s
spec:
  selector:
    app: %s
  ports:
  - name: http
    port: %d
    protocol: TCP
    targetPort: %d
`, name, serviceName, name, sourcePort, targetPort))
	if err != nil {
		return nil, err
	}
	return svc, nil
}

func generateIngress(name string, domain string, sourcePort int) (*yaml.RNode, error) {
	serviceName := name + "-svc"
	ing, err := yaml.Parse(fmt.Sprintf(`apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: %s
  namespace: monitoring
spec:
  rules:
  - host: %s
    http:
      paths:
      - pathType: Prefix
        backend:
          service:
            name: %s
            port:
              number: %d
        path: /
`, name, domain, serviceName, sourcePort))
	if err != nil {
		return nil, err
	}
	return ing, nil
}

func main() {
	config := new(App)
	fn := func(items []*yaml.RNode) ([]*yaml.RNode, error) {

		var newNodes []*yaml.RNode
		resourceName := config.Metadata.Name

		for _ = range items {

			// get generated Deployment
			deploy, err := generateDeployment(resourceName, config.Spec.Image)
			if err != nil {
				return nil, err
			}
			newNodes = append(newNodes, deploy)

			// get generated Service
			service, err := generateService(resourceName, config.Spec.Port, config.Spec.Port)
			if err != nil {
				return nil, err
			}
			newNodes = append(newNodes, service)

			// get generated Ingress
			ingress, err := generateIngress(resourceName, config.Spec.Domain, config.Spec.Port)
			if err != nil {
				return nil, err
			}
			newNodes = append(newNodes, ingress)

		}

		items = newNodes

		return items, nil
	}
	p := framework.SimpleProcessor{Config: config, Filter: kio.FilterFunc(fn)}
	cmd := command.Build(p, command.StandaloneDisabled, false)
	command.AddGenerateDockerfile(cmd)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
