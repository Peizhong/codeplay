package app_engine

import (
	"context"

	"github.com/peizhong/codeplay/pkg/logger"
	"gopkg.in/yaml.v3"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	Name  string `json:"name"`
	Image string `json:"image"`
	Ports []struct {
		Name          string `json:"name"`
		ContainerPort uint16 `json:"container_port"`
		Proto         string `json:"proto"`
	} `json:"ports"`
	Limits struct {
		Cpu    string
		Memory string
	} `json:"limits"`
	Command []string `json:"command"`
}

func (c *AbstractContainer) ConvertToK8sContainer() (*corev1.Container, error) {
	container := &corev1.Container{
		Name:  c.Name,
		Image: c.Image,
	}
	container.Command = append(container.Command, c.Command...)
	container.Resources.Limits = make(corev1.ResourceList)
	if c.Limits.Cpu != "" {
		if quantity, err := resource.ParseQuantity(c.Limits.Cpu); err == nil {
			container.Resources.Limits[corev1.ResourceCPU] = quantity
		}
	}
	if c.Limits.Memory != "" {
		if quantity, err := resource.ParseQuantity(c.Limits.Memory); err == nil {
			container.Resources.Limits[corev1.ResourceMemory] = quantity
		}
	}
	return container, nil
}

type AbstractApp struct {
	Namespace   string               `json:"namespace"`
	Name        string               `json:"name"`
	Labels      map[string]string    `json:"labels"`
	Annotations map[string]string    `json:"annotations"`
	Containers  []*AbstractContainer `json:"containers"`
}

func (app *AbstractApp) ConvertToK8sDeployment() (*appsv1.Deployment, error) {
	deployment := &appsv1.Deployment{}
	deployment.Name = app.Name
	deployment.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: map[string]string{
			"app": app.Name,
		},
	}
	deployment.Spec.Template.Labels = map[string]string{
		"app": app.Name,
	}
	for k, v := range app.Labels {
		deployment.Spec.Template.Labels[k] = v
	}
	deployment.Spec.Template.Annotations = map[string]string{
		"prometheus.io/scrape": "true",
	}
	for k, v := range app.Annotations {
		deployment.Spec.Template.Annotations[k] = v
	}
	deployment.Spec.Template.Spec.ServiceAccountName = "app-user"
	for _, ac := range app.Containers {
		container, err := ac.ConvertToK8sContainer()
		if err != nil {
			continue
		}
		deployment.Spec.Template.Spec.Containers = append(deployment.Spec.Template.Spec.Containers, *container)
	}
	return deployment, nil
}

func CreateApp(ctx context.Context, app *AbstractApp) error {
	deployment, err := app.ConvertToK8sDeployment()
	if err != nil {
		return err
	}
	bs, _ := yaml.Marshal(deployment)
	logger.ToFile("deploy.txt", bs)
	_, err = clientSet.AppsV1().Deployments(app.Namespace).Create(ctx, deployment, metav1.CreateOptions{})
	return err
}

func UpdateAppImage(ctx context.Context, app AbstractApp) error {
	return nil
}

func ListPod(ctx context.Context, namespace string) ([]string, error) {
	list, err := clientSet.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		Limit: 10,
	})
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(list.Items))
	for _, item := range list.Items {
		result = append(result, item.Name)
	}
	return result, nil
}
