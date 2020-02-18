package main

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/arutselvan15/estore-common/clients"
	cfg "github.com/arutselvan15/estore-common/config"
	cLog "github.com/arutselvan15/estore-common/log"
)

func main() {
	var (
		config *rest.Config
		err    error
	)

	log := cLog.GetLogger("common")

	// read kube config in case defined in environment variable
	if cfg.GetKubeConfigPath() != "" {
		config, err = clientcmd.BuildConfigFromFlags("", cfg.GetKubeConfigPath())
		if err != nil {
			log.Fatalf("error creating config using kube config path: %v", err)
		}
	} else {
		// default get current cluster config
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Fatalf("error creating config using cluster config: %v", err)
		}
	}

	estoreClients, err := clients.NewEstoreClientForConfig(config)
	if err != nil {
		log.Fatalf("error creating estore clients for config: %v", err)
	}

	pdtList, err := estoreClients.GetProductClient().EstoreV1().Products("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("error listing estore products : %v", err)
	}

	for _, pdt := range pdtList.Items {
		log.Debugf("%s - %s", pdt.Namespace, pdt.Name)
	}
}
