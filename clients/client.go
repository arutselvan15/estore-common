// Package clients provides common clients
package clients

import (
	kube "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	pdt "github.com/arutselvan15/estore-product-kube-client/pkg/client/clientset/versioned"
)

// EstoreClientInterface estore client interface
type EstoreClientInterface interface {
	GetKubeClient() kube.Interface
	GetProductClient() pdt.Interface
}

type estoreClient struct {
	k8s  *kube.Clientset
	ePdt *pdt.Clientset
}

// NewEstoreClientForConfig estore client for config
func NewEstoreClientForConfig(config *rest.Config) (EstoreClientInterface, error) {
	var err error

	c := new(estoreClient)

	c.k8s, err = kube.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	c.ePdt, err = pdt.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *estoreClient) GetKubeClient() kube.Interface {
	return c.k8s
}

func (c *estoreClient) GetProductClient() pdt.Interface {
	return c.ePdt
}
