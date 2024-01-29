package tidb_exporter

import (
	"fmt"
	"log"

	"github.com/hashicorp/consul/api"
)

const (
	CONSUL_SERVICE_NAME = "tidb"
	CONSUL_SERVICE_TAG  = "pro-tidb"
	CONSUL_SERVICE_ADDR = "api-consul.mz.com"
)

func RegisterService(ip string, registerPort, checkPort int, cluster, idc, project, business string) error {
	cfg := api.DefaultConfig()
	cfg.Address = CONSUL_SERVICE_ADDR
	cfg.Scheme = "http"
	client, err := api.NewClient(cfg)
	if err != nil {
		log.Fatal("Failed to create Consul client:", err)
	}
	serviceDef := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-%s:%d", cluster, ip, registerPort),
		Name:    CONSUL_SERVICE_NAME,
		Address: ip,
		Port:    registerPort,
		Check: &api.AgentServiceCheck{
			//	HTTP:     fmt.Sprintf("http://%s:%d/ping", ip, checkPort),
			TCP:      fmt.Sprintf("%s:%d", ip, checkPort),
			Interval: "30s",
		},
		EnableTagOverride: true,
		Tags:              []string{CONSUL_SERVICE_TAG},
		Meta: map[string]string{
			"idc":        idc,
			"os":         "kvm",
			"project":    project,
			"business":   business,
			"cluster_id": cluster,
			//"target":     fmt.Sprintf("mongodb://%s:%d", ip, port),
			"addr": fmt.Sprintf("%s:%d", ip, registerPort),
		},
	}
	// Register the service
	if err = client.Agent().ServiceRegisterOpts(serviceDef, api.ServiceRegisterOpts{
		ReplaceExistingChecks: true,
	}); err != nil {
		log.Fatal("Failed to register service:", err)
	}
	return nil
}
