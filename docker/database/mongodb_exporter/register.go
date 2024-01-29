package mongodb_exporter

import (
	"fmt"
	"log"

	"github.com/hashicorp/consul/api"
)

const (
	CONSUL_SERVICE_NAME = "mongodb_exporter"
	CONSUL_SERVICE_TAG  = "dev-mongodb"
	EXPORTER_PORT       = 9216
)

var consulServiceAddr = "api-consul.mz.com"

func ListService() {
	cfg := api.DefaultConfig()
	cfg.Address = consulServiceAddr
	cfg.Scheme = "http"
	client, err := api.NewClient(cfg)
	if err != nil {
		log.Fatal("Failed to create Consul client:", err)
	}
	services, _, err := client.Catalog().Service(CONSUL_SERVICE_NAME, "", nil)
	if err != nil {
		log.Fatal("Failed to get service:", err)
	}
	for _, service := range services {
		fmt.Println(service.ServiceID, service.ServiceName, service.ServiceAddress, service.ServicePort)
	}
}

func RegisterService(ip string, port int, cluster, idc, project, business string) error {
	cfg := api.DefaultConfig()
	cfg.Address = consulServiceAddr
	cfg.Scheme = "http"
	client, err := api.NewClient(cfg)
	if err != nil {
		log.Fatal("Failed to create Consul client:", err)
	}
	serviceDef := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-%s:%d", cluster, ip, port),
		Name:    CONSUL_SERVICE_NAME,
		Address: ip,
		Port:    EXPORTER_PORT,
		Check: &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%d/ping", ip, EXPORTER_PORT),
			Interval: "30s",
		},
		Tags: []string{CONSUL_SERVICE_TAG},
		Meta: map[string]string{
			"idc":        idc,
			"os":         "kvm",
			"project":    project,
			"business":   business,
			"cluster_id": cluster,
			"target":     fmt.Sprintf("mongodb://%s:%d", ip, port),
			"addr":       fmt.Sprintf("%s:%d", ip, port),
		},
	}
	// Register the service
	if err = client.Agent().ServiceRegister(serviceDef); err != nil {
		log.Fatal("Failed to register service:", err)
	}
	return nil
}

func DeRegister(id string) error {
	cfg := api.DefaultConfig()
	cfg.Address = consulServiceAddr
	cfg.Scheme = "http"
	client, err := api.NewClient(cfg)
	if err != nil {
		log.Fatal("Failed to create Consul client:", err)
	}
	return client.Agent().ServiceDeregister(id)
}
