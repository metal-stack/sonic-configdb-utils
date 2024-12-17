package configdb

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/metal-stack/sonic-configdb-utils/values"
)

func Test_getInterfaces(t *testing.T) {
	tests := []struct {
		name     string
		ports    []values.Port
		bgpPorts []string
		want     map[string]Interface
	}{
		{
			name:  "empty ports",
			ports: []values.Port{},
			want:  map[string]Interface{},
		},
		{
			name: "port not in bgp ports, with no vrf and no ips",
			ports: []values.Port{
				{
					Name: "Ethernet0",
				},
			},
			want: map[string]Interface{},
		},
		{
			name: "port in bgp ports",
			ports: []values.Port{
				{
					Name: "Ethernet0",
				},
			},
			bgpPorts: []string{"Ethernet0"},
			want: map[string]Interface{
				"Ethernet0": {
					IPv6UseLinkLocalOnly: IPv6UseLinkLocalOnlyModeEnable,
				},
			},
		},
		{
			name: "port with vrf",
			ports: []values.Port{
				{
					Name: "Ethernet0",
					VRF:  "Vrf40",
				},
			},
			want: map[string]Interface{
				"Ethernet0": {
					VRFName: "Vrf40",
				},
			},
		},
		{
			name: "port in bgp ports with vrf",
			ports: []values.Port{
				{
					Name: "Ethernet0",
					VRF:  "Vrf40",
				},
			},
			bgpPorts: []string{"Ethernet0"},
			want: map[string]Interface{
				"Ethernet0": {
					IPv6UseLinkLocalOnly: IPv6UseLinkLocalOnlyModeEnable,
					VRFName:              "Vrf40",
				},
			},
		},
		{
			name: "port not in bgp ports without vrf but with and ips",
			ports: []values.Port{
				{
					Name: "Ethernet0",
					IPs:  []string{"10.1.1.1"},
				},
			},
			want: map[string]Interface{
				"Ethernet0":          {},
				"Ethernet0|10.1.1.1": {},
			},
		},
		{
			name: "port in bgp ports with vrf and ips",
			ports: []values.Port{
				{
					Name: "Ethernet0",
					VRF:  "Vrf40",
					IPs:  []string{"10.1.1.1"},
				},
			},
			bgpPorts: []string{"Ethernet0"},
			want: map[string]Interface{
				"Ethernet0": {
					IPv6UseLinkLocalOnly: IPv6UseLinkLocalOnlyModeEnable,
					VRFName:              "Vrf40",
				},
				"Ethernet0|10.1.1.1": {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getInterfaces(tt.ports, tt.bgpPorts)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("getInterfaces() %v", diff)
			}
		})
	}
}
