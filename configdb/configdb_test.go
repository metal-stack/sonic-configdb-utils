package configdb

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/metal-stack/metal-lib/pkg/testcommon"
	p "github.com/metal-stack/sonic-configdb-utils/platform"
	"github.com/metal-stack/sonic-configdb-utils/values"
	v "github.com/metal-stack/sonic-configdb-utils/version"
)

func Test_getInterfaces(t *testing.T) {
	tests := []struct {
		name          string
		ports         *values.Ports
		bgpPorts      []string
		interconnects map[string]values.Interconnect
		interfaces    map[string]values.Interface
		want          map[string]Interface
	}{
		{
			name: "empty ports",
			ports: &values.Ports{
				List: []values.Port{},
			},
			want: map[string]Interface{},
		},
		{
			name: "port not in bgp ports, with no vrf and no ips",
			ports: &values.Ports{
				List: []values.Port{
					{
						Name: "Ethernet0",
					},
				},
			},
			want: map[string]Interface{},
		},
		{
			name: "port in bgp ports",
			ports: &values.Ports{
				List: []values.Port{
					{
						Name: "Ethernet0",
					},
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
			ports: &values.Ports{
				List: []values.Port{
					{
						Name: "Ethernet0",
						VRF:  "Vrf40",
					},
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
			ports: &values.Ports{
				List: []values.Port{
					{
						Name: "Ethernet0",
						VRF:  "Vrf40",
					},
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
			name: "port not in bgp ports without vrf but with ip",
			ports: &values.Ports{
				List: []values.Port{
					{
						Name: "Ethernet0",
						IPs:  []string{"10.1.1.1"},
					},
				},
			},
			want: map[string]Interface{
				"Ethernet0":          {},
				"Ethernet0|10.1.1.1": {},
			},
		},
		{
			name: "port in bgp ports with vrf and ips",
			ports: &values.Ports{
				List: []values.Port{
					{
						Name: "Ethernet0",
						VRF:  "Vrf40",
						IPs:  []string{"10.1.1.1"},
					},
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
		{
			name:     "port in bgp ports but not in ports",
			bgpPorts: []string{"Ethernet0"},
			want: map[string]Interface{
				"Ethernet0": {
					IPv6UseLinkLocalOnly: IPv6UseLinkLocalOnlyModeEnable,
				},
			},
		},
		{
			name: "interconnect without unnumbered interfaces",
			interconnects: map[string]values.Interconnect{
				"internet": {
					UnnumberedInterfaces: []string{},
					VNI:                  "104000",
					VRF:                  "VrfInternet",
				},
			},
			want: map[string]Interface{},
		},
		{
			name: "interconnect with unnumbered interfaces",
			interconnects: map[string]values.Interconnect{
				"internet": {
					UnnumberedInterfaces: []string{
						"Ethernet0",
						"Ethernet1",
					},
					VNI: "104000",
					VRF: "VrfInternet",
				},
			},
			want: map[string]Interface{
				"Ethernet0": {
					IPv6UseLinkLocalOnly: IPv6UseLinkLocalOnlyModeEnable,
					VRFName:              "VrfInternet",
				},
				"Ethernet1": {
					IPv6UseLinkLocalOnly: IPv6UseLinkLocalOnlyModeEnable,
					VRFName:              "VrfInternet",
				},
			},
		},
		{
			name: "interconnect with unnumbered interfaces but without vrf",
			interconnects: map[string]values.Interconnect{
				"internet": {
					UnnumberedInterfaces: []string{
						"Ethernet0",
						"Ethernet1",
					},
				},
			},
			want: map[string]Interface{
				"Ethernet0": {
					IPv6UseLinkLocalOnly: IPv6UseLinkLocalOnlyModeEnable,
				},
				"Ethernet1": {
					IPv6UseLinkLocalOnly: IPv6UseLinkLocalOnlyModeEnable,
				},
			},
		},
		{
			name: "port with ip in interconnects",
			ports: &values.Ports{
				List: []values.Port{
					{
						Name: "Ethernet0",
						IPs:  []string{"10.0.0.1/32"},
					},
				},
			},
			interconnects: map[string]values.Interconnect{
				"internal": {
					UnnumberedInterfaces: []string{"Ethernet0"},
					VRF:                  "internal",
				},
			},
			want: map[string]Interface{
				"Ethernet0": {
					IPv6UseLinkLocalOnly: IPv6UseLinkLocalOnlyModeEnable,
					VRFName:              "internal",
				},
				"Ethernet0|10.0.0.1/32": {},
			},
		},
		{
			name: "additional interfaces with ips",
			interfaces: map[string]values.Interface{
				"Ethernet50.10": {
					IPs: []string{
						"172.16.0.1/32",
						"10.1.1.1/32",
					},
				},
			},
			want: map[string]Interface{
				"Ethernet50.10":               {},
				"Ethernet50.10|172.16.0.1/32": {},
				"Ethernet50.10|10.1.1.1/32":   {},
			},
		},
		{
			name: "additional ips for existing interface entry",
			ports: &values.Ports{
				List: []values.Port{
					{
						Name: "Ethernet0",
						IPs:  []string{"10.0.0.1/32"},
					},
				},
			},
			interconnects: map[string]values.Interconnect{
				"internal": {
					UnnumberedInterfaces: []string{"Ethernet0"},
					VRF:                  "internal",
				},
			},
			interfaces: map[string]values.Interface{
				"Ethernet0": {
					IPs: []string{"172.16.0.1/32"},
				},
			},
			want: map[string]Interface{
				"Ethernet0": {
					IPv6UseLinkLocalOnly: IPv6UseLinkLocalOnlyModeEnable,
					VRFName:              "internal",
				},
				"Ethernet0|10.0.0.1/32":   {},
				"Ethernet0|172.16.0.1/32": {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getInterfaces(tt.ports, tt.bgpPorts, tt.interconnects, tt.interfaces)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("getInterfaces() %v", diff)
			}
		})
	}
}

func Test_getPortsAndBreakouts(t *testing.T) {
	tests := []struct {
		name          string
		ports         *values.Ports
		breakouts     map[string]string
		platform      *p.Platform
		wantPorts     map[string]Port
		wantBreakouts map[string]BreakoutConfig
		wantErr       error
	}{
		{
			name: "only breakouts defined",
			ports: &values.Ports{
				List: []values.Port{},
			},
			breakouts: map[string]string{
				"Ethernet0": "1x100G[40G]",
			},
			platform: &p.Platform{
				Interfaces: map[string]p.Interface{
					"Ethernet0": {
						BreakoutModes: map[string][]string{
							"1x100G[40G]": {"Eth1(Port1)"},
						},
						Index: "1,1,1,1",
						Lanes: "1,2,3,4",
					},
				},
			},
			wantPorts: map[string]Port{
				"Ethernet0": {
					AdminStatus:    defaultAdminStatus,
					Alias:          "Eth1(Port1)",
					Index:          "1",
					Lanes:          "1,2,3,4",
					MTU:            fmt.Sprintf("%d", defaultMTU),
					parentBreakout: "1x100G[40G]",
					Speed:          "100000",
				},
			},
			wantBreakouts: map[string]BreakoutConfig{
				"Ethernet0": {
					BreakoutMode: "1x100G[40G]",
				},
			},
			wantErr: nil,
		},
		{
			name: "only ports defined",
			ports: &values.Ports{
				List: []values.Port{
					{
						Name:    "Ethernet0",
						Speed:   40000,
						Autoneg: values.AutonegModeOn,
					},
					{
						Name:  "Ethernet4",
						Speed: 40000,
					},
				},
			},
			platform: &p.Platform{
				Interfaces: map[string]p.Interface{
					"Ethernet0": {
						BreakoutModes: map[string][]string{
							"1x100G[40G]": {"Eth1(Port1)"},
							"2x50G":       {"Eth1/1(Port1)", "Eth1/2(Port1)"},
						},
						Index: "1,1,1,1",
						Lanes: "1,2,3,4",
					},
					"Ethernet4": {
						BreakoutModes: map[string][]string{
							"1x100G[40G]": {"Eth2(Port2)"},
							"2x50G":       {"Eth2/1(Port2)", "Eth2/2(Port2)"},
						},
						Index: "2,2,2,2",
						Lanes: "5,6,7,8",
					},
					"Ethernet8": {
						BreakoutModes: map[string][]string{
							"1x100G[40G]": {"Eth3(Port3)"},
							"2x50G":       {"Eth3/1(Port3)", "Eth3/2(Port3)"},
						},
						Index: "3,3,3,3",
						Lanes: "9,10,11,12",
					},
				},
			},
			wantPorts: map[string]Port{
				"Ethernet0": {
					AdminStatus:    AdminStatusUp,
					Alias:          "Eth1(Port1)",
					Autoneg:        AutonegModeOn,
					Index:          "1",
					Lanes:          "1,2,3,4",
					MTU:            fmt.Sprintf("%d", defaultMTU),
					parentBreakout: "1x100G[40G]",
					Speed:          "40000",
				},
				"Ethernet4": {
					AdminStatus:    AdminStatusUp,
					Alias:          "Eth2(Port2)",
					Index:          "2",
					Lanes:          "5,6,7,8",
					MTU:            fmt.Sprintf("%d", defaultMTU),
					parentBreakout: "1x100G[40G]",
					Speed:          "40000",
				},
				"Ethernet8": {
					AdminStatus:    AdminStatusUp,
					Alias:          "Eth3(Port3)",
					Index:          "3",
					Lanes:          "9,10,11,12",
					MTU:            fmt.Sprintf("%d", defaultMTU),
					parentBreakout: "1x100G[40G]",
					Speed:          "100000",
				},
			},
			wantBreakouts: map[string]BreakoutConfig{
				"Ethernet0": {
					BreakoutMode: "1x100G[40G]",
				},
				"Ethernet4": {
					BreakoutMode: "1x100G[40G]",
				},
				"Ethernet8": {
					BreakoutMode: "1x100G[40G]",
				},
			},
			wantErr: nil,
		},
		{
			name: "child-port is not present if breakout config 'absorbs' it",
			breakouts: map[string]string{
				"Ethernet4": "1x100G[40G]",
			},
			ports: &values.Ports{
				List: []values.Port{
					{
						Name: "Ethernet5",
					},
				},
			},
			platform: &p.Platform{
				Interfaces: map[string]p.Interface{
					"Ethernet4": {
						BreakoutModes: map[string][]string{
							"1x100G[40G]": {"Eth2(Port2)"},
							"4x25G":       {"Eth2/1(Port2)", "Eth2/2(Port2)", "Eth2/3(Port2)", "Eth2/4(Port2)"},
						},
						Index: "2,2,2,2",
						Lanes: "5,6,7,8",
					},
				},
			},
			wantErr: fmt.Errorf("invalid port name Ethernet5; if you think it should be available please check your breakout configuration"),
		},
		{
			name: "port speed 0 is allowed",
			breakouts: map[string]string{
				"Ethernet4": "1x100G[40G]",
			},
			ports: &values.Ports{
				DefaultFEC: values.FECModeRS,
				DefaultMTU: 1500,
				List: []values.Port{
					{
						Name:  "Ethernet4",
						Speed: 0,
					},
				}},
			platform: &p.Platform{
				Interfaces: map[string]p.Interface{
					"Ethernet4": {
						BreakoutModes: map[string][]string{
							"1x100G[40G]": {"Eth2(Port2)"},
							"4x25G":       {"Eth2/1(Port2)", "Eth2/2(Port2)", "Eth2/3(Port2)", "Eth2/4(Port2)"},
						},
						Index: "2,2,2,2",
						Lanes: "5,6,7,8",
					},
				},
			},
			wantBreakouts: map[string]BreakoutConfig{
				"Ethernet4": {
					BreakoutMode: "1x100G[40G]",
				},
			},
			wantPorts: map[string]Port{
				"Ethernet4": {
					AdminStatus:    AdminStatusUp,
					Alias:          "Eth2(Port2)",
					FEC:            FECModeRS,
					Index:          "2",
					Lanes:          "5,6,7,8",
					MTU:            "1500",
					parentBreakout: "1x100G[40G]",
					Speed:          "100000",
				},
			},
			wantErr: nil,
		},
		{
			name: "fec, mtu and speed can be overriden",
			breakouts: map[string]string{
				"Ethernet4": "1x100G[40G]",
			},
			ports: &values.Ports{
				List: []values.Port{
					{
						FECMode: values.FECModeRS,
						MTU:     1500,
						Name:    "Ethernet4",
						Speed:   40000,
					},
				},
			},
			platform: &p.Platform{
				Interfaces: map[string]p.Interface{
					"Ethernet4": {
						BreakoutModes: map[string][]string{
							"1x100G[40G]": {"Eth2(Port2)"},
							"4x25G":       {"Eth2/1(Port2)", "Eth2/2(Port2)", "Eth2/3(Port2)", "Eth2/4(Port2)"},
						},
						Index: "2,2,2,2",
						Lanes: "5,6,7,8",
					},
				},
			},
			wantBreakouts: map[string]BreakoutConfig{
				"Ethernet4": {
					BreakoutMode: "1x100G[40G]",
				},
			},
			wantPorts: map[string]Port{
				"Ethernet4": {
					AdminStatus:    AdminStatusUp,
					Alias:          "Eth2(Port2)",
					FEC:            FECModeRS,
					Index:          "2",
					Lanes:          "5,6,7,8",
					MTU:            "1500",
					parentBreakout: "1x100G[40G]",
					Speed:          "40000",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPorts, gotBreakouts, err := getPortsAndBreakouts(tt.ports, tt.breakouts, tt.platform)
			if diff := cmp.Diff(tt.wantErr, err, testcommon.ErrorStringComparer()); diff != "" {
				t.Errorf("getPortsAndBreakouts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.wantPorts, gotPorts, testcommon.IgnoreUnexported()); diff != "" {
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
		ports         *values.Ports
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
			ports: &values.Ports{
				List: []values.Port{},
			},
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
			ports: &values.Ports{
				List: []values.Port{
					{
						VRF: "Vrf40",
					},
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
			ports: &values.Ports{
				List: []values.Port{
					{
						VRF: "Vrf41",
					},
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

func Test_getSAG(t *testing.T) {
	tests := []struct {
		name    string
		sag     *values.SAG
		version v.Branch
		want    *SAG
		wantErr bool
	}{
		{
			name:    "wrong version",
			sag:     &values.SAG{},
			version: v.Branch202111,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty mac",
			sag:     &values.SAG{},
			version: v.Branch202211,
			want:    nil,
			wantErr: false,
		},
		{
			name: "valid",
			sag: &values.SAG{
				MAC: "11:11:11:11:11:11",
			},
			version: v.Branch202211,
			want: &SAG{
				SAGGlobal: SAGGlobal{
					GatewayMAC: "11:11:11:11:11:11",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getSAG(tt.sag, tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("getSAG() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("getSAG() diff = %s", diff)
			}
		})
	}
}

func Test_getSFLOW(t *testing.T) {
	bFalse := false
	version202211 := &v.Version{Branch: string(v.Branch202211)}
	version202111 := &v.Version{Branch: string(v.Branch202111)}
	tests := []struct {
		name			string
		sflow			*values.SFLOW
		version			*v.Version
		wantGlobal 		map[string]SFLOWGlobal
		wantCollector 	map[string]SFLOWCollector
		wantSessions 	map[string]SFLOWSession
		wantErr 		bool
	}{
		{
			name: "nil input",
			sflow: nil,
			version: version202211,
		},
		{
			name: "disabled",
			sflow: &values.SFLOW{Enabled: false},
			version: version202211,
		},
		{
			name: "collector with empty name",
			sflow: &values.SFLOW{
				Enabled: true,
				Collectors: []values.SFLOWCollector{{IP: "172.17.0.1" }},
			},
			version: version202211,
			wantErr: true,
		},
		{
			name: "collector with empty ip",
			sflow: &values.SFLOW{
				Enabled: true,
				Collectors: []values.SFLOWCollector{{Name: "goflow2" }},
			},
			version: version202211,
			wantErr: true,
		},
		{
			name: "enabled with default port",
			sflow: &values.SFLOW{
				Enabled: true,
				PollingInterval: 20,
				Collectors: []values.SFLOWCollector{
					{Name: "goflow2", IP: "172.17.0.1", VRF: "default" },
				},
			},
			version: version202211,
			wantGlobal: map[string]SFLOWGlobal{
				"global": {AdminStatus: AdminStatusUp, PollingInterval: "20"},
			},
			wantCollector: map[string]SFLOWCollector{
				"goflow2": {CollectorIP: "172.17.0.1", CollectorPort: "6343", CollectorVRF: "default"},
			},
			wantSessions: map[string]SFLOWSession{},
		},
		{
			name: "enabled with multiple collectors and sessions",
			sflow: &values.SFLOW{
				Enabled: true,
				PollingInterval: 20,
				Collectors: []values.SFLOWCollector{
					{Name: "primary", IP: "172.17.0.1", Port: 9999},
					{Name: "backup", IP: "172.17.0.1", Port: 6343},
				},
				Sessions: []values.SFLOWSession{
					{Interface: "Ethernet0", SampleRate: 1024},
				},
			},
			version: version202211,
			wantGlobal: map[string]SFLOWGlobal{
				"global": {AdminStatus: AdminStatusUp, PollingInterval: "20"},
			},
			wantCollector: map[string]SFLOWCollector{
				"primary": {CollectorIP: "172.17.0.1", CollectorPort: "9999"},
				"backup":  {CollectorIP: "172.17.0.1", CollectorPort: "6343"},
			},
			wantSessions: map[string]SFLOWSession{
				"Ethernet0": {AdminStatus: AdminStatusUp, SampleRate: "1024"},
			},
		},
		{
			name: "low polling interval and disabled session",
			sflow: &values.SFLOW{
				Enabled: true,
				PollingInterval: 2,
				Sessions: []values.SFLOWSession{
					{Interface: "Ethernet0", Enabled: &bFalse},
					{Interface: "Ethernet4"},
				},
			},
			version: version202211,
			wantGlobal: map[string]SFLOWGlobal{
				"global": {AdminStatus: AdminStatusUp, PollingInterval: "5"},
			},
			wantCollector: map[string]SFLOWCollector{},
			wantSessions: map[string]SFLOWSession{
				"Ethernet0": {AdminStatus: AdminStatusDown},
				"Ethernet4": {AdminStatus: AdminStatusUp},
			},
		},
		{
			name: "202111 rejects collector_vrf",
			sflow: &values.SFLOW{
				Enabled: true,
				PollingInterval: 20,
				Collectors: []values.SFLOWCollector{
					{Name: "goflow2", IP: "172.17.0.1", VRF: "default"},
				},
			},
			version: version202111,
			wantErr: true,
		},
		{
			name: "202111 without collector_vrf",
			sflow: &values.SFLOW{
				Enabled: true,
				PollingInterval: 20,
				Collectors: []values.SFLOWCollector{
					{Name: "goflow2", IP: "172.17.0.1"},
				},
				Sessions: []values.SFLOWSession{
					{Interface: "Ethernet0", SampleRate: 1024},
				},
			},
			version: version202111,
			wantGlobal: map[string]SFLOWGlobal{
				"global": {AdminStatus: AdminStatusUp, PollingInterval: "20"},
			},
			wantCollector: map[string]SFLOWCollector{
				"goflow2": {CollectorIP: "172.17.0.1", CollectorPort: "6343"},
			},
			wantSessions: map[string]SFLOWSession{
				"Ethernet0": {AdminStatus: AdminStatusUp, SampleRate: "1024"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGlobal, gotCollectors, gotSessions, err := getSFLOW(tt.sflow, tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("getSFLOW() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.wantGlobal, gotGlobal); diff != "" {
				t.Errorf("getSFLOW() global diff = %s", diff)
			}
			if diff := cmp.Diff(tt.wantCollector, gotCollectors); diff != "" {
				t.Errorf("getSFLOW() collector diff = %s", diff)
			}
			if diff := cmp.Diff(tt.wantSessions, gotSessions); diff != "" {
				t.Errorf("getSFLOW() session diff = %s", diff)
			}
		})
	}
}
