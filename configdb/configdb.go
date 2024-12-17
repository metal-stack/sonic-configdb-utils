package configdb

import (
	"slices"

	"github.com/metal-stack/sonic-configdb-utils/values"
)

type ConfigDB struct {
	ACLRules          map[string]ACLRule        `json:"ACL_RULE"`
	ACLTables         map[string]ACLTable       `json:"ACL_TABLE"`
	Breakouts         map[string]BreakoutConfig `json:"BREAKOUT_CFG"`
	DeviceMetadata    `json:"DEVICE_METADATA"`
	Features          map[string]Feature   `json:"FEATURE"`
	Interfaces        map[string]Interface `json:"INTERFACE"`
	LLDP              `json:"LLDP"`
	LoopbackInterface map[string]struct{}       `json:"LOOPBACK_INTERFACE"`
	MCLAGDomains      map[string]MCLAGDomain    `json:"MCLAG_DOMAIN"`
	MCLAGInterfaces   map[string]MCLAGInterface `json:"MCLAG_INTERFACE"`
	MCLAGUniqueIPs    map[string]MCLAGUniqueIP  `json:"MCLAG_UNIQUE_IP"`
	MgmtInterfaces    map[string]MgmtInterface  `json:"MGMT_INTERFACE"`
	MgmtPorts         map[string]MgmtPort       `json:"MGMT_PORT"`
	MgmtVRFConfig     `json:"MGMT_VRF_CONFIG"`
	NTP               `json:"NTP"`
	NTPServers        map[string]struct{} `json:"NTP_SERVER"`
	Ports             map[string]Port     `json:"PORT"`
	VLANs             map[string]VLAN     `json:"VLAN"`
	VLANInterface     map[string]struct{} `json:"VLAN_INTERFACE"`
	VXLANEVPN         `json:"VXLAN_EVPN_NVO"`
	VXLANTunnels      map[string]VXLANTunnel    `json:"VXLAN_TUNNEL"`
	VXLANTunnelMaps   map[string]VXLANTunnelMap `json:"VXLAN_TUNNEL_MAP"`
}

func GenerateConfigDB(input *values.Values) *ConfigDB {
	configdb := ConfigDB{
		ACLRules:  map[string]ACLRule{},
		ACLTables: map[string]ACLTable{},
		Breakouts: getBreakoutConfig(input.Breakouts),
		DeviceMetadata: DeviceMetadata{
			Localhost: Metadata{
				DockerRoutingConfigMode: DockerRoutingConfigMode(input.DockerRoutingConfigMode),
				FRRMgmtFrameworkConfig:  input.FRRMgmtFrameworkConfig,
				Hostname:                input.Hostname,
				RouterType:              "LeafRouter",
			},
		},
		Features: map[string]Feature{
			"dhcp_relay": {
				AutoRestart: FeatureModeEnabled,
				State:       FeatureModeEnabled,
			},
		},
		Interfaces: getInterfaces(input.Ports, input.BGPPorts),
		LLDP:       LLDP{},
		LoopbackInterface: map[string]struct{}{
			"Loopback0": {},
			"Loopback0|" + input.LoopbackAddress + "/32": {},
		},
		MCLAGDomains: map[string]MCLAGDomain{
			"1": {
				MCLAGSystemID: input.SystemMAC,
				PeerIP:        input.PeerIP,
				PeerLink:      input.PeerLink,
				SourceIP:      input.SourceIP,
			},
		},
		MCLAGInterfaces: getMCLAGInterfaces(input.MemberPortChannels),
		MCLAGUniqueIPs: map[string]MCLAGUniqueIP{
			"Vlan" + input.KeepaliveVLAN: {
				UniqueIP: MCLAGUniqueIPModeEnable,
			},
		},
		MgmtInterfaces: getMgmtInterfaces(input.MgmtIfIP, input.MgmtIfGateway),
		MgmtPorts: map[string]MgmtPort{
			"eth0": {
				AdminStatus: AdminStatusUp,
				Alias:       "eth0",
				Description: "Management Port",
			},
		},
		MgmtVRFConfig: MgmtVRFConfig{
			VRFGlobal: VRFGlobal{
				MgmtVRFEnabled: input.MgmtVRF,
			},
		},
		NTP: NTP{
			NTPGlobal: NTPGlobal{
				SrcIntf: "eth0",
			},
		},
		NTPServers:      getNTPServers(input.NTPServers),
		Ports:           getPorts(input.Ports, input.Breakouts),
		VLANs:           map[string]VLAN{},
		VLANInterface:   map[string]struct{}{},
		VXLANEVPN:       VXLANEVPN{},
		VXLANTunnels:    map[string]VXLANTunnel{},
		VXLANTunnelMaps: map[string]VXLANTunnelMap{},
	}
	return &configdb
}

func getBreakoutConfig(breakouts map[string]string) map[string]BreakoutConfig {
	breakoutConfig := make(map[string]BreakoutConfig)

	for port, breakout := range breakouts {
		breakoutConfig[port] = BreakoutConfig{
			BreakoutMode: BreakoutMode(breakout),
		}
	}

	return breakoutConfig
}

func getInterfaces(ports []values.Port, bgpPorts []string) map[string]Interface {
	interfaces := make(map[string]Interface)

	for _, port := range ports {
		if len(port.IPs) == 0 && port.VRF == "" && !slices.Contains(bgpPorts, port.Name) {
			continue
		}

		intf := Interface{}
		if slices.Contains(bgpPorts, port.Name) {
			intf.IPv6UseLinkLocalOnly = IPv6UseLinkLocalOnlyModeEnable
		}
		if port.VRF != "" {
			intf.VRFName = port.VRF
		}
		interfaces[port.Name] = intf

		for _, ip := range port.IPs {
			intf = Interface{}
			interfaces[port.Name+"|"+ip] = intf
		}
	}

	return interfaces
}

func getMCLAGInterfaces(memberPortChannels []string) map[string]MCLAGInterface {
	mclagInterfaces := make(map[string]MCLAGInterface)

	for _, channel := range memberPortChannels {
		mclagInterfaces["1|Portchannel"+channel] = MCLAGInterface{
			IfType: "PortChannel",
		}
	}

	return mclagInterfaces
}

func getMgmtInterfaces(mgmtIfIP, mgmtIfGateway string) map[string]MgmtInterface {
	if mgmtIfIP == "" {
		return nil
	}

	mgmtInterfaces := make(map[string]MgmtInterface)

	eth0 := MgmtInterface{}
	if mgmtIfGateway != "" {
		eth0.GWAddr = mgmtIfGateway
	}

	mgmtInterfaces["eth0|"+mgmtIfIP] = eth0

	return mgmtInterfaces
}

func getNTPServers(servers []string) map[string]struct{} {
	ntpServers := make(map[string]struct{})

	for _, server := range servers {
		ntpServers[server] = struct{}{}
	}

	return ntpServers
}

func getPorts(ports []values.Port, breakouts map[string]string) map[string]Port {
	configPorts := make(map[string]Port)

	// TODO: deduce ports from breakouts and default or deduce all other fields, then no running config should be needed

	/*
	 - use breakouts to add defaults for each port
	 - for each port check if there is specific config in ports for it and override if necessary
	*/

	for _, port := range ports {
		configPorts[port.Name] = Port{
			AdminStatus: AdminStatusUp,
			Alias:       "",
			Autoneg:     AutonegModeOff,
			FEC:         "",
			Index:       0,
			Lanes:       "",
			MTU:         0,
			ParentPort:  "",
			Speed:       0,
		}
	}

	return configPorts
}
