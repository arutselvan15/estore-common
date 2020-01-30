// Package fake provides fake clients
package fake

import (
	"k8s.io/apimachinery/pkg/runtime"
	kube "k8s.io/client-go/kubernetes"
	kubeFake "k8s.io/client-go/kubernetes/fake"

	pdt "github.com/arutselvan15/estore-product-kube-client/pkg/client/clientset/versioned"
	pdtFake "github.com/arutselvan15/estore-product-kube-client/pkg/client/clientset/versioned/fake"

	"github.com/arutselvan15/estore-common/clients"
)

type estoreFakeClient struct {
	k8s  *kubeFake.Clientset
	ePdt *pdtFake.Clientset
}

// NewEstoreFakeClientForConfig fake clients
func NewEstoreFakeClientForConfig(pdtObjects, kubeObjects []runtime.Object) (clients.EstoreClientInterface, error) {
	c := new(estoreFakeClient)
	c.k8s = kubeFake.NewSimpleClientset(kubeObjects...)
	c.ePdt = pdtFake.NewSimpleClientset(pdtObjects...)

	return c, nil
}

func (c *estoreFakeClient) GetKubeClient() kube.Interface {
	return c.k8s
}

func (c *estoreFakeClient) GetProductClient() pdt.Interface {
	return c.ePdt
}
