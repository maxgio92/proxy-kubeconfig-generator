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
package e2e

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/maxgio92/proxy-kubeconfig-generator/pkg/utils"
	"k8s.io/client-go/kubernetes"
)

const (
	ServiceAccountName  = "myapp"
	KubeconfigSecretKey = "kubeconfig"
	Namespace           = "default"
)

func main() {
	config, err := utils.BuildClientConfig()
	if err != nil {
		panic(err)
	}

	clientset := kubernetes.NewForConfigOrDie(config)

	// Retrieve the Kubeconfig secret and build a new client Config
	tenantClientConfig, err := utils.BuildClientConfigFromSecret(clientset, ServiceAccountName+"-kubeconfig", KubeconfigSecretKey, Namespace)
	if err != nil {
		panic(err)
	}

	tenantClientset := kubernetes.NewForConfigOrDie(tenantClientConfig)

	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: "001"},
	}

	_, err = tenantClientset.CoreV1().Namespaces().Create(
		context.TODO(),
		namespace,
		metav1.CreateOptions{},
	)
	if err != nil {
		panic(err)
	}

	// Get Tenant's Namesapces through its ClientSet
	tenantNamespaces, err := tenantClientset.CoreV1().Namespaces().List(
		context.Background(),
		metav1.ListOptions{},
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nThe Tenant owner %s can list only these Namespaces through the proxy:\n", ServiceAccountName)
	for _, tenantNamespace := range tenantNamespaces.Items {
		fmt.Println(tenantNamespace.Name)
	}
}
