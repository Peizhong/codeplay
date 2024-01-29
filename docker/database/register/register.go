package register

import (
	"log"
	"time"

	"github.com/hashicorp/consul/api"
)

const (
	CONSUL_SERVICE_ADDR = "api-consul.mz.com"
)

func UnRegister(id string) error {
	cfg := api.DefaultConfig()
	cfg.Address = CONSUL_SERVICE_ADDR
	cfg.Scheme = "http"
	client, err := api.NewClient(cfg)
	if err != nil {
		log.Fatal("Failed to create Consul client:", err)
	}
	for i := 0; i < 5; i++ {
		if err = client.Agent().ServiceDeregister(id); err == nil {
			return nil
		}
		<-time.After(time.Millisecond * 100)
	}
	return err
}
