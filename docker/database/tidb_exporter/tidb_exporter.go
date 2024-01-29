package tidb_exporter

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/peizhong/codeplay/docker/database/register"
)

type Topology struct {
	ClusterName string // aka cluster_id
	Project     string
	Business    string

	DPInstances   []string
	TiDBInstances []string
	TiKVInstances []string
}

func (t *Topology) Print() string {
	return fmt.Sprintf("group: %s, pd: %v, tidb: %v, tikv: %v", t.ClusterName, t.DPInstances, t.TiDBInstances, t.TiKVInstances)
}

// https://docs.pingcap.com/tidb/stable/grafana-monitor-best-practices
var ports = map[string]string{
	"pd":   "2379",  // pd server, 集群管理, raft数据一致
	"tidb": "10080", // tidb status，sql请求, 无状态
	"tikv": "20180", // tikv status, 存储数据, raft数据一致
}

func ProbeFromIpMap(cluster map[string][]string) []Topology {
	var result []Topology
	for group, ips := range cluster {
		distIp := make(map[string]struct{})
		for _, ip := range ips {
			distIp[ip] = struct{}{}
		}
		validAddrs := make(map[string][]string)
		for ip := range distIp {
			for db, port := range ports {
				addr := fmt.Sprintf("http://%s:%s/metrics", ip, port)
				func(addr string) {
					timeoutCtx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
					defer cancel()
					req, _ := http.NewRequestWithContext(timeoutCtx, http.MethodGet, addr, nil)
					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						return
					}
					io.Copy(io.Discard, resp.Body)
					resp.Body.Close()
					if resp.StatusCode == http.StatusOK {
						validAddrs[db] = append(validAddrs[db], ip)
					}
				}(addr)
			}
		}
		if len(validAddrs) > 0 {
			parts := strings.Split(group, "/")
			result = append(result, Topology{
				ClusterName:   group,
				Project:       parts[0],
				Business:      parts[1],
				DPInstances:   validAddrs["pd"],
				TiDBInstances: validAddrs["tidb"],
				TiKVInstances: validAddrs["tikv"],
			})
		}
	}
	return result
}

func RegisterToConsul(topo *Topology) error {
	// https://github.com/pingcap/tidb/tree/release-7.5/pkg/metrics/grafana
	for _, instance := range topo.DPInstances {
		if err := RegisterService(instance, 2379, 2379, fmt.Sprintf("%s-%s-DP", topo.Project, topo.Business), "gz", topo.Project, topo.Business); err != nil {
			return err
		}
	}
	for _, instance := range topo.TiDBInstances {
		if err := RegisterService(instance, 10080, 10080, fmt.Sprintf("%s-%s-TiDB", topo.Project, topo.Business), "gz", topo.Project, topo.Business); err != nil {
			return err
		}
	}
	for _, instance := range topo.TiKVInstances {
		if err := RegisterService(instance, 20180, 20180, fmt.Sprintf("%s-%s-TiKV", topo.Project, topo.Business), "gz", topo.Project, topo.Business); err != nil {
			return err
		}
	}
	return nil
}

func UnRegisterFromConsul(topo *Topology) error {
	for _, instance := range topo.DPInstances {
		if err := register.UnRegister(fmt.Sprintf("%s-%s-DP-%s:%d", topo.Project, topo.Business, instance, 2379)); err != nil {
			// return err
		}
	}
	for _, instance := range topo.TiDBInstances {
		if err := register.UnRegister(fmt.Sprintf("%s-%s-TiDB-%s:%d", topo.Project, topo.Business, instance, 10080)); err != nil {
			// return err
		}
	}
	for _, instance := range topo.TiKVInstances {
		if err := register.UnRegister(fmt.Sprintf("%s-%s-TiKV-%s:%d", topo.Project, topo.Business, instance, 20180)); err != nil {
			// return err
		}
	}
	return nil
}
