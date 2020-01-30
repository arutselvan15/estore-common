package main

import (
	"flag"
	"os"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/arutselvan15/estore-common/clients"
	gLog "github.com/arutselvan15/estore-common/log"
)

func main() {
	var (
		kubeConfig string
		config     *rest.Config
		err        error
	)

	log := gLog.GetInstance()

	flag.StringVar(&kubeConfig, "kubeConfig", kubeConfig, "kubeConfig file")
	flag.Parse()

	// default get current cluster config
	config, err = rest.InClusterConfig()
	if err != nil {
		log.Fatalf("error creating config: %v", err)
	}

	// get kube config from env when not passed in args
	if kubeConfig == "" {
		kubeConfig = os.Getenv("KUBECONFIG")
	}

	// override default config when kubeconfig passed in args or KUBECONFIG found
	if kubeConfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
		if err != nil {
			log.Fatalf("error creating config: %v", err)
		}
	}

	estoreClients, err := clients.NewEstoreClientForConfig(config)
	if err != nil {
		log.Fatalf("error creating estore clients for config: %v", err)
	}

	pdtList, err := estoreClients.GetProductClient().EstoreV1().Products("").List(metaV1.ListOptions{})
	if err != nil {
		log.Fatalf("error listing estore products : %v", err)
	}

	for _, pdt := range pdtList.Items {
		log.Debugf("%s - %s", pdt.Namespace, pdt.Name)
	}
}
