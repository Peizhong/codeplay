package tidb_exporter

import (
	"testing"
)

func TestRegister(t *testing.T) {
	ips := map[string][]string{
		"bidgata/db": []string{
			"10.128.4.117",
			"10.128.64.43",
			"10.128.64.27",
			"10.128.64.43",
			"10.128.64.27",
			"10.128.64.30",
			"10.128.64.30",
			"10.128.64.32",
			"10.128.64.35",
		},
		"news/content": []string{
			"10.128.4.217",
			"10.128.64.41",
			"10.128.64.42",
			"10.128.64.41",
			"10.128.64.42",
			"10.128.64.23",
			"10.128.64.23",
			"10.128.64.25",
			"10.128.64.26",
		},
	}
	result := ProbeFromIpMap(ips)
	for _, p := range result {
		t.Log(p.Print())
	}
	for _, item := range result {
		RegisterToConsul(&item)
	}
}

func TestUnRegister(t *testing.T) {
	ips := map[string][]string{
		"bidgata/db": []string{
			"10.128.4.117",
			"10.128.64.43",
			"10.128.64.27",
			"10.128.64.43",
			"10.128.64.27",
			"10.128.64.30",
			"10.128.64.30",
			"10.128.64.32",
			"10.128.64.35",
		},
		"news/content": []string{
			"10.128.4.217",
			"10.128.64.41",
			"10.128.64.42",
			"10.128.64.41",
			"10.128.64.42",
			"10.128.64.23",
			"10.128.64.23",
			"10.128.64.25",
			"10.128.64.26",
		},
		"test/db": []string{
			"10.128.109.5",
			"10.128.109.6",
			"10.128.109.7",
			"10.128.8.8",
			"10.129.4.55",
			"10.129.4.56",
			"10.129.4.57",
			"10.129.4.58",
			"10.129.4.59",
		},
	}
	result := ProbeFromIpMap(ips)
	for _, p := range result {
		t.Log(p.Print())
	}
	for _, item := range result {
		UnRegisterFromConsul(&item)
	}
}
