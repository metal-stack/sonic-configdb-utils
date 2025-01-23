package values

import "gopkg.in/yaml.v2"

type DockerRoutingConfigMode string

const (
	DockerRoutingConfigModeSeparated DockerRoutingConfigMode = "separated"
	DockerRoutingConfigModeSplit     DockerRoutingConfigMode = "split"
	DockerRoutingConfigModeUnified   DockerRoutingConfigMode = "unified"
)

type FECMode string

const (
	FECModeNone FECMode = "none"
	FECModeRS   FECMode = "rs"
)

type Interconnect struct {
	VNI string `yaml:"vni,omitempty"`
	VRF string `yaml:"vrf,omitempty"`
}

type MCLAG struct {
	KeepaliveVLAN      string   `yaml:"keepalive_vlan,omitempty"`
	MemberPortChannels []string `yaml:"member_port_channels,omitempty"`
	PeerIP             string   `yaml:"peer_ip,omitempty"`
	PeerLink           string   `yaml:"peer_link,omitempty"`
	SourceIP           string   `yaml:"source_ip,omitempty"`
	SystemMAC          string   `yaml:"system_mac,omitempty"`
}

type Port struct {
	IPs     []string `yaml:"ips,omitempty"`
	FECMode `yaml:"fec,omitempty"`
	MTU     int    `yaml:"mtu,omitempty"`
	Name    string `yaml:"name,omitempty"`
	Speed   int    `yaml:"speed,omitempty"`
	VRF     string `yaml:"vrf,omitempty"`
}

type PortChannel struct {
	Number   string   `yaml:"number,omitempty"`
	MTU      int      `yaml:"mtu,omitempty"`
	Fallback bool     `yaml:"fallback,omitempty"`
	Members  []string `yaml:"members,omitempty"`
}

type SAG struct {
	MAC string `yaml:"mac,omitempty"`
}

type Values struct {
	BGPPorts                []string          `yaml:"sonic_bgp_ports,omitempty"`
	Breakouts               map[string]string `yaml:"sonic_breakouts,omitempty"`
	DockerRoutingConfigMode `yaml:"sonic_docker_routing_config_mode,omitempty"`
	FRRMgmtFrameworkConfig  bool                    `yaml:"sonic_frr_mgmt_framework_config,omitempty"`
	Hostname                string                  `yaml:"sonic_hostname,omitempty"`
	Interconnects           map[string]Interconnect `yaml:"sonic_interconnects,omitempty"`
	LLDPHelloTimer          int                     `yaml:"sonic_lldp_hello_timer,omitempty"`
	LoopbackAddress         string                  `yaml:"sonic_loopback_address,omitempty"`
	MCLAG                   `yaml:"sonic_mclag,omitempty"`
	MgmtIfGateway           string        `yaml:"sonic_mgmtif_gateway,omitempty"`
	MgmtIfIP                string        `yaml:"sonic_mgmtif_ip,omitempty"`
	MgmtVRF                 bool          `yaml:"sonic_mgmt_vrf,omitempty"`
	NTPServers              []string      `yaml:"sonic_ntpservers,omitempty"`
	PortChannels            []PortChannel `yaml:"sonic_portchannels,omitempty"`
	PortChannelsDefaultMTU  int           `yaml:"sonic_portchannels_default_mtu,omitempty"`
	Ports                   []Port        `yaml:"sonic_ports,omitempty"`
	PortsDefaultFEC         FECMode       `yaml:"sonic_ports_default_fec,omitempty"`
	PortsDefaultMTU         int           `yaml:"sonic_ports_default_mtu,omitempty"`
	PortsDefaultSpeed       int           `yaml:"sonic_ports_default_speed,omitempty"`
	SAG                     `yaml:"sonic_sag,omitempty"`
	SSHSourceranges         []string `yaml:"sonic_ssh_sourceranges,omitempty"`
	VLANs                   []VLAN   `yaml:"sonic_vlans,omitempty"`
	VTEPs                   []VTEP   `yaml:"sonic_vteps,omitempty"`
}

type VLAN struct {
	DHCPServers   []string `yaml:"dhcp_servers,omitempty"`
	ID            string   `yaml:"id,omitempty"`
	IP            string   `yaml:"ip,omitempty"`
	SAG           bool     `yaml:"sag,omitempty"`
	TaggedPorts   []string `yaml:"tagged_ports,omitempty"`
	UntaggedPorts []string `yaml:"untagged_ports,omitempty"`
	VRF           string   `yaml:"vrf,omitempty"`
}

type VTEP struct {
	Comment string `yaml:"comment,omitempty"`
	VNI     string `yaml:"vni,omitempty"`
	VLAN    string `yaml:"vlan,omitempty"`
}

func UnmarshalValues(input []byte) (*Values, error) {
	values := Values{}
	err := yaml.Unmarshal(input, &values)
	if err != nil {
		return nil, err
	}

	return &values, nil
}
