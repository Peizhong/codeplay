package app_engine

import (
	"context"

	"github.com/peizhong/codeplay/pkg/logger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

var clientSet kubernetes.Interface

func InitInClusterK8sClient() error {
	// inside k8s mode
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	clientSetImp, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	version, err := clientSetImp.ServerVersion()
	if err != nil {
		return err
	}
	clientSet = clientSetImp
	logger.Sugar().Infow("init k8s client", "version", version.String())
	return nil
}

type AbstractContainer struct {
	Name  string
	Image string
	Ports []struct {
		Name          string
		ContainerPort uint16
		Proto         string
	}
}

type AbstractApp struct {
	Namespace   string
	Name        string
	Labels      map[string]string
	Annotations map[string]string
	Containers  []AbstractContainer
}

func CreateApp(ctx context.Context, app AbstractApp) error {
	return nil
}

func UpdateAppImage(ctx context.Context, app AbstractApp) error {
	return nil
}
