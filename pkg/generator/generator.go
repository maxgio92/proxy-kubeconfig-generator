/*
Copyright © 2022 maxgio92 me@maxgio.it
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package generator

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/maxgio92/proxy-kubeconfig-generator/pkg/utils"
)

const (
	DefaultTLSecretCAKey       = "ca"
	DefaultNamespace           = "default"
	DefaultKubeconfigSecretKey = "kubeconfig"
)

// Generates a Secret containing a Kubeconfig with the specified
// Service Account's token and server URL and CA certificate from
// the specified Secret.
func GenerateProxyKubeconfigFromSA(clientset *kubernetes.Clientset, serviceAccountName string, namespace string, server string, serverTLSSecretName string, serverTLSSecretCAKey string, serverTLSSecretNamespace string, kubeconfigSecretKey string) (*clientcmdapi.Config, error) {
	// Get Tenant Service Account token
	saSecret, err := utils.GetServiceAccountSecret(clientset, serviceAccountName, namespace)
	if err != nil {
		return nil, err
	}
	if _, ok := saSecret.Data["token"]; !ok {
		return nil, fmt.Errorf("secret %s does not contain a token", saSecret.Name)
	}

	// Get Server Proxy CA certificate
	proxyCA, err := utils.GetSecretField(clientset, serverTLSSecretName, serverTLSSecretCAKey, serverTLSSecretNamespace)
	if err != nil {
		return nil, err
	}

	// Generate the client Config for the Tenant Owner
	tenantConfig, err := utils.BuildKubeconfigFromToken(saSecret.Data["token"], proxyCA, server, serverTLSSecretNamespace)
	if err != nil {
		return nil, err
	}

	return tenantConfig, err
}
