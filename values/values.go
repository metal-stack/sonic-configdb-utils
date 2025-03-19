package values

import "gopkg.in/yaml.v3"

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
	BGPPorts                []string          `yaml:"bgp_ports,omitempty"`
	Breakouts               map[string]string `yaml:"breakouts,omitempty"`
	DockerRoutingConfigMode `yaml:"docker_routing_config_mode,omitempty"`
	FRRMgmtFrameworkConfig  bool                    `yaml:"frr_mgmt_framework_config,omitempty"`
	Hostname                string                  `yaml:"hostname,omitempty"`
	Interconnects           map[string]Interconnect `yaml:"interconnects,omitempty"`
	LLDPHelloTime           int                     `yaml:"lldp_hello_time,omitempty"`
	LoopbackAddress         string                  `yaml:"loopback_address,omitempty"`
	MCLAG                   *MCLAG                  `yaml:"mclag,omitempty"`
	MgmtIfGateway           string                  `yaml:"mgmtif_gateway,omitempty"`
	MgmtIfIP                string                  `yaml:"mgmtif_ip,omitempty"`
	MgmtVRF                 bool                    `yaml:"mgmt_vrf,omitempty"`
	Nameservers             []string                `yaml:"nameservers,omitempty"`
	NTPServers              []string                `yaml:"ntpservers,omitempty"`
	PortChannels            []PortChannel           `yaml:"portchannels,omitempty"`
	PortChannelsDefaultMTU  int                     `yaml:"portchannels_default_mtu,omitempty"`
	Ports                   []Port                  `yaml:"ports,omitempty"`
	PortsDefaultFEC         FECMode                 `yaml:"ports_default_fec,omitempty"`
	PortsDefaultMTU         int                     `yaml:"ports_default_mtu,omitempty"`
	SAG                     `yaml:"sag,omitempty"`
	SSHSourceranges         []string `yaml:"ssh_sourceranges,omitempty"`
	VLANMembers             bool     `yaml:"vlan_members,omitempty"`
	VLANs                   []VLAN   `yaml:"vlans,omitempty"`
	VTEPs                   []VTEP   `yaml:"vteps,omitempty"`
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

func UnmarshalValues(in []byte) (*Values, error) {
	var values Values
	err := yaml.Unmarshal(in, &values)
	if err != nil {
		return nil, err
	}

	return &values, nil
}
