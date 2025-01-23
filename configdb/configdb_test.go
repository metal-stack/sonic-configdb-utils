package configdb

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/metal-stack/metal-lib/pkg/testcommon"
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

func Test_getPortsAndBreakouts(t *testing.T) {
	tests := []struct {
		name           string
		ports          []values.Port
		breakouts      map[string]string
		defaultFECMode values.FECMode
		defaultMTU     int
		wantPorts      map[string]Port
		wantBreakouts  map[string]BreakoutConfig
		wantErr        error
	}{
		{
			name:  "only breakouts defined",
			ports: []values.Port{},
			breakouts: map[string]string{
				"Ethernet0": "1x100G[40G]",
			},
			wantPorts: map[string]Port{
				"Ethernet0": {
					AdminStatus: defaultAdminStatus,
					Alias:       "Eth1(Port1)",
					Autoneg:     defaultAutonegMode,
					FEC:         defaultFECMode,
					Index:       "1",
					Lanes:       "1,2,3,4",
					MTU:         fmt.Sprintf("%d", defaultMTU),
					ParentPort:  "Ethernet0",
					Speed:       "100000",
				},
			},
			wantBreakouts: map[string]BreakoutConfig{
				"Ethernet0": {
					BreakoutMode: BreakoutMode1x100G,
				},
			},
			wantErr: nil,
		},
		{
			name: "port is missing breakout config",
			breakouts: map[string]string{
				"Ethernet4": "4x25G",
			},
			ports: []values.Port{
				{
					Name: "Ethernet0",
				},
			},
			wantErr: fmt.Errorf("no breakout configuration found for port Ethernet0"),
		},
		{
			name: "child-port is not present if breakout config 'absorbs' it",
			breakouts: map[string]string{
				"Ethernet4": "1x100G[40G]",
			},
			ports: []values.Port{
				{
					Name: "Ethernet5",
				},
			},
			wantErr: fmt.Errorf("no breakout configuration found for port Ethernet5"),
		},
		{
			name: "port speed does not match breakout config",
			breakouts: map[string]string{
				"Ethernet4": "4x25G",
			},
			ports: []values.Port{
				{
					Name:  "Ethernet5",
					Speed: 10000,
				},
			},
			wantErr: fmt.Errorf("invalid speed 10000 for port Ethernet5; check breakout configuration"),
		},
		{
			name: "only port speed 40G is allowed to override preconfigured speed",
			breakouts: map[string]string{
				"Ethernet4": "1x100G[40G]",
			},
			ports: []values.Port{
				{
					Name:  "Ethernet4",
					Speed: 50000,
				},
			},
			wantErr: fmt.Errorf("invalid speed 50000 for port Ethernet4; current breakout configuration only allows values 100000 or 40000"),
		},
		{
			name: "fec, mtu and speed can be overriden",
			breakouts: map[string]string{
				"Ethernet4": "1x100G[40G]",
			},
			ports: []values.Port{
				{
					FECMode: values.FECModeRS,
					MTU:     1500,
					Name:    "Ethernet4",
					Speed:   40000,
				},
			},
			wantBreakouts: map[string]BreakoutConfig{
				"Ethernet4": {
					BreakoutMode: BreakoutMode1x100G,
				},
			},
			wantPorts: map[string]Port{
				"Ethernet4": {
					AdminStatus: AdminStatusUp,
					Alias:       "Eth2(Port2)",
					Autoneg:     AutonegModeOff,
					FEC:         FECModeRS,
					Index:       "2",
					Lanes:       "5,6,7,8",
					MTU:         "1500",
					ParentPort:  "Ethernet4",
					Speed:       "40000",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPorts, gotBreakouts, err := getPortsAndBreakouts(tt.ports, tt.breakouts, tt.defaultFECMode, tt.defaultMTU)
			if diff := cmp.Diff(tt.wantErr, err, testcommon.ErrorStringComparer()); diff != "" {
				t.Errorf("getPortsAndBreakouts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.wantPorts, gotPorts); diff != "" {
				t.Errorf("getPortsAndBreakouts() diff = %v", diff)
			}
			if diff := cmp.Diff(tt.wantBreakouts, gotBreakouts); diff != "" {
				t.Errorf("getPortsAndBreakouts() diff = %v", diff)
			}
		})
	}
}

func Test_getVRFs(t *testing.T) {
	tests := []struct {
		name          string
		interconnects map[string]values.Interconnect
		ports         []values.Port
		vlans         []values.VLAN
		want          map[string]VRF
	}{
		{
			name: "no ports or vlans to add",
			interconnects: map[string]values.Interconnect{
				"mpls": {
					VNI: "1",
					VRF: "Vrf40",
				},
			},
			ports: []values.Port{},
			vlans: []values.VLAN{},
			want: map[string]VRF{
				"Vrf40": {
					VNI: "1",
				},
			},
		},
		{
			name: "duplicates are not added",
			interconnects: map[string]values.Interconnect{
				"mpls": {
					VNI: "1",
					VRF: "Vrf40",
				},
			},
			ports: []values.Port{
				{
					VRF: "Vrf40",
				},
			},
			vlans: []values.VLAN{
				{
					VRF: "Vrf40",
				},
			},
			want: map[string]VRF{
				"Vrf40": {
					VNI: "1",
				},
			},
		},
		{
			name: "new vrfs are added",
			interconnects: map[string]values.Interconnect{
				"mpls": {
					VNI: "1",
					VRF: "Vrf40",
				},
			},
			ports: []values.Port{
				{
					VRF: "Vrf41",
				},
			},
			vlans: []values.VLAN{
				{
					VRF: "Vrf42",
				},
			},
			want: map[string]VRF{
				"Vrf40": {
					VNI: "1",
				},
				"Vrf41": {},
				"Vrf42": {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getVRFs(tt.interconnects, tt.ports, tt.vlans)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("getVRFs() diff = %s", diff)
			}
		})
	}
}
