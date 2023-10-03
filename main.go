package main

import (
	"context"
	"os"
	"strings"

	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	coreV1Types "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	DOPPLER_NAMESPACE = "doppler-operator-system"
)

// Retrieves all environment variables from the system
// That end in the string DOPPLER_TOKEN
func get_environment_variables() map[string]string {
	m := make(map[string]string)
	for _, e := range os.Environ() {
		s := strings.Split(e, "=")
		if strings.HasSuffix(s[0], "DOPPLER_TOKEN") {
			m[s[0]] = s[1]
		}
	}
	return m
}

// Create a client that can be used for interacting with secrets
// in a given cluster
func initClient() coreV1Types.SecretInterface {
	// Location of kubeconfig file
	kubeconfig := os.Getenv("HOME") + "/.kube/config"

	// Create a config
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// Create an api clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	secretsClient := clientset.CoreV1().Secrets(DOPPLER_NAMESPACE)

	return secretsClient
}

func main() {
	secretsClient := initClient()
	dopplerTokens := get_environment_variables()

	for key, value := range dopplerTokens {
		secretName := strings.ToLower(key)
		secretName = strings.ReplaceAll(secretName, "_", "-")

		data := map[string][]byte{
			"serviceToken": []byte(value),
		}

		object := metaV1.ObjectMeta{Name: secretName, Namespace: DOPPLER_NAMESPACE}
		secret := &coreV1.Secret{Data: data, ObjectMeta: object}

		_, err := secretsClient.Create(context.TODO(), secret, metaV1.CreateOptions{})
		if err != nil {
			panic(err)
		}
	}
}
