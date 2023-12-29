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
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"

	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/maxgio92/proxy-kubeconfig-generator/pkg/generator"
	"github.com/maxgio92/proxy-kubeconfig-generator/pkg/utils"
)

const (
	DefaultTLSecretCAKey       = "ca"
	DefaultNamespace           = "default"
	DefaultKubeconfigSecretKey = "kubeconfig"
)

func main() {
	serviceAccountName := flag.String("serviceaccount", "", "The name of the service account for which to create the kubeconfig")
	namespace := flag.String("namespace", DefaultNamespace, "(optional) The namespace of the service account and where the kubeconfig secret will be created.")
	server := flag.String("server", "", "The server url of the kubeconfig where API requests will be sent")
	serverTLSSecretNamespace := flag.String("server-tls-secret-namespace", DefaultNamespace, "(optional) The namespace of the server TLS secret.")
	serverTLSSecretName := flag.String("server-tls-secret-name", "", "The server TLS secret name")
	serverTLSSecretCAKey := flag.String("server-tls-secret-ca-key", DefaultTLSecretCAKey, "(optional) The CA key in the server TLS secret.")
	kubeconfigSecretKey := flag.String("kubeconfig-secret-key", DefaultKubeconfigSecretKey, "(optional) The key of the kubeconfig in the secret that will be created")

	flag.Parse()

	if *serviceAccountName == "" {
		flag.Usage()
		fmt.Println(errors.New("missing service account name"))
		os.Exit(1)
	}

	if *server == "" {
		fmt.Println(errors.New("missing server url"))
		os.Exit(1)
	}

	if *serverTLSSecretName == "" {
		flag.Usage()
		fmt.Println(errors.New("missing server TLS secret name"))
		os.Exit(1)
	}

	config, err := utils.BuildClientConfig()
	if err != nil {
		fmt.Println(errors.Wrap(err, "error building client config"))
		os.Exit(1)
	}

	clientset := kubernetes.NewForConfigOrDie(config)

	tenantConfig, err := generator.GenerateProxyKubeconfigFromSA(clientset, *serviceAccountName, *namespace, *server, *serverTLSSecretName, *serverTLSSecretCAKey, *serverTLSSecretNamespace, *kubeconfigSecretKey)
	if err != nil {
		fmt.Println(errors.Wrap(err, "error generating kubeconfig from service account"))
		os.Exit(1)
	}

	err = utils.CreateKubeconfigSecret(clientset, tenantConfig, *namespace, *serviceAccountName+"-kubeconfig", *kubeconfigSecretKey)
	if err != nil {
		fmt.Println(errors.Wrap(err, "error creating kubeconfig secret"))
		os.Exit(1)
	}

	fmt.Println("Proxy kubeconfig Secret created")
}
