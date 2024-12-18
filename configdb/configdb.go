package configdb

import (
	"fmt"
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

func GenerateConfigDB(input *values.Values) (*ConfigDB, error) {
	ports, breakouts, err := getPortsAndBreakouts(input.Ports, input.Breakouts)
	if err != nil {
		return nil, err
	}

	configdb := ConfigDB{
		ACLRules:  map[string]ACLRule{},
		ACLTables: map[string]ACLTable{},
		Breakouts: breakouts,
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
		Ports:           ports,
		VLANs:           map[string]VLAN{},
		VLANInterface:   map[string]struct{}{},
		VXLANEVPN:       VXLANEVPN{},
		VXLANTunnels:    map[string]VXLANTunnel{},
		VXLANTunnelMaps: map[string]VXLANTunnelMap{},
	}
	return &configdb, nil
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

func getPortsAndBreakouts(ports []values.Port, breakouts map[string]string) (map[string]Port, map[string]BreakoutConfig, error) {
	configPorts := make(map[string]Port)
	configBreakouts := make(map[string]BreakoutConfig)

	for portName, breakout := range breakouts {
		breakoutPorts, err := getPortsFromBreakout(portName, breakout)
		if err != nil {
			return nil, nil, err
		}
		for name, port := range breakoutPorts {
			configPorts[name] = port
		}
		configBreakouts[portName] = BreakoutConfig{
			BreakoutMode: BreakoutMode(breakout),
		}
	}

	for _, port := range ports {
		configPort, ok := configPorts[port.Name]
		if !ok {
			return nil, nil, fmt.Errorf("no breakout configuration found for port %s", port.Name)
		}

		switch configPort.Speed {
		case 100000:
			if port.Speed == 100000 || port.Speed == 40000 {
				configPort.Speed = port.Speed
			} else {
				return nil, nil, fmt.Errorf("invalid speed %d for port %s; current breakout configuration only allows values 100000 or 40000", port.Speed, port.Name)
			}
		default:
			if port.Speed != configPort.Speed {
				return nil, nil, fmt.Errorf("invalid speed %d for port %s; check breakout configuration", port.Speed, port.Name)
			}
		}

		if string(port.FECMode) != string(configPort.FEC) {
			configPort.FEC = FECMode(port.FECMode)
		}
		if port.MTU != configPort.MTU {
			configPort.MTU = port.MTU
		}
		configPorts[port.Name] = configPort
	}

	return configPorts, configBreakouts, nil
}
