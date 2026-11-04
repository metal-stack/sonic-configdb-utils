package values

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

type ACL struct {
	Description string    `yaml:"description"`
	Name        string    `yaml:"name"`
	Ports       []string  `yaml:"ports"`
	Rules       []ACLRule `yaml:"rules"`
	Stage       string    `yaml:"stage"`
}

type ACLRule struct {
	Action string `yaml:"action"`
	Dst    string `yaml:"dst"`
	Src    string `yaml:"src"`
}

type AutonegMode string

const (
	AutonegModeOn  AutonegMode = "on"
	AutonegModeOff AutonegMode = "off"
)

type DockerRoutingConfigMode string

const (
	DockerRoutingConfigModeSeparated    DockerRoutingConfigMode = "separated"
	DockerRoutingConfigModeSplit        DockerRoutingConfigMode = "split"
	DockerRoutingConfigModeSplitUnified DockerRoutingConfigMode = "split-unified"
	DockerRoutingConfigModeUnified      DockerRoutingConfigMode = "unified"
)

type Feature struct {
	AutoRestart bool `yaml:"auto_restart"`
	Enabled     bool `yaml:"enabled"`
}

type FECMode string

const (
	FECModeNone FECMode = "none"
	FECModeRS   FECMode = "rs"
)

type Interconnect struct {
	UnnumberedInterfaces []string `yaml:"unnumbered_interfaces"`
	VNI                  string   `yaml:"vni"`
	VRF                  string   `yaml:"vrf"`
}

type Interface struct {
	IPs []string `yaml:"ips"`
}

type MCLAG struct {
	KeepaliveVLAN      string   `yaml:"keepalive_vlan"`
	MemberPortChannels []string `yaml:"member_port_channels"`
	PeerIP             string   `yaml:"peer_ip"`
	PeerLink           string   `yaml:"peer_link"`
	SourceIP           string   `yaml:"source_ip"`
	SystemMAC          string   `yaml:"system_mac"`
}

type MgmtInterface struct {
	GatewayAddress string `yaml:"gateway_address"`
	IP             string `yaml:"ip"`
}

type NTP struct {
	Servers      []string `yaml:"servers"`
	SrcInterface string   `yaml:"src_interface"`
	VRF          string   `yaml:"vrf"`
}

type Port struct {
	Autoneg AutonegMode `yaml:"autoneg"`
	FECMode FECMode     `yaml:"fec"`
	IPs     []string    `yaml:"ips"`
	MTU     int         `yaml:"mtu"`
	Name    string      `yaml:"name"`
	Speed   int         `yaml:"speed"`
	VRF     string      `yaml:"vrf"`
}

type Ports struct {
	DefaultAutoneg AutonegMode `yaml:"default_autoneg"`
	DefaultFEC     FECMode     `yaml:"default_fec"`
	DefaultMTU     int         `yaml:"default_mtu"`
	List           []Port      `yaml:"list"`
}

type PortChannel struct {
	Number   string   `yaml:"number"`
	MTU      int      `yaml:"mtu"`
	Fallback bool     `yaml:"fallback"`
	Members  []string `yaml:"members"`
}

type PortChannels struct {
	DefaultMTU int           `yaml:"default_mtu"`
	List       []PortChannel `yaml:"list"`
}

type SAG struct {
	MAC string `yaml:"mac"`
}

type Values struct {
	ACLs                    []ACL                   `yaml:"acls"`
	BGPPorts                []string                `yaml:"bgp_ports"`
	Breakouts               map[string]string       `yaml:"breakouts"`
	DockerRoutingConfigMode DockerRoutingConfigMode `yaml:"docker_routing_config_mode"`
	Features                map[string]Feature      `yaml:"features"`
	FRRMgmtFrameworkConfig  bool                    `yaml:"frr_mgmt_framework_config"`
	Hostname                string                  `yaml:"hostname"`
	Interconnects           map[string]Interconnect `yaml:"interconnects"`
	Interfaces              map[string]Interface    `yaml:"interfaces"`
	LLDPHelloTime           int                     `yaml:"lldp_hello_time"`
	LoopbackAddress         string                  `yaml:"loopback_address"`
	MCLAG                   MCLAG                   `yaml:"mclag"`
	MgmtInterface           MgmtInterface           `yaml:"mgmt_interface"`
	MgmtVRF                 bool                    `yaml:"mgmt_vrf"`
	Nameservers             []string                `yaml:"nameservers"`
	NTP                     NTP                     `yaml:"ntp"`
	PortChannels            PortChannels            `yaml:"portchannels"`
	Ports                   *Ports                  `yaml:"ports"`
	SAG                     *SAG                    `yaml:"sag"`
	SSHSourceranges         []string                `yaml:"ssh_sourceranges"`
	VLANs                   []VLAN                  `yaml:"vlans"`
	VLANSubinterfaces       []VLANSubinterface      `yaml:"vlan_subinterfaces"`
	VTEP                    VTEP                    `yaml:"vtep"`
}

type VLANSubinterface struct {
	CIDRs []string `yaml:"cidrs"`
	Port  string   `yaml:"port"`
	VLAN  string   `yaml:"vlan"`
	VRF   string   `yaml:"vrf"`
}

type VLAN struct {
	DHCPRelayServerVRF string   `yaml:"dhcp_relay_server_vrf"`
	DHCPRelaySrcIntf   string   `yaml:"dhcp_relay_src_intf"`
	DHCPServers        []string `yaml:"dhcp_servers"`
	ID                 string   `yaml:"id"`
	IP                 string   `yaml:"ip"`
	SAG                *bool    `yaml:"sag"`
	TaggedPorts        []string `yaml:"tagged_ports"`
	UntaggedPorts      []string `yaml:"untagged_ports"`
	VRF                string   `yaml:"vrf"`
	VRRP               VRRP     `yaml:"vrrp"`
}

type VRRP struct {
	Group string `yaml:"group"`
	IP    string `yaml:"ip"`
}

type VTEP struct {
	Enabled         bool             `yaml:"enabled"`
	VXLANTunnelMaps []VXLANTunnelMap `yaml:"vxlan_tunnel_maps"`
}

type VXLANTunnelMap struct {
	VNI  string `yaml:"vni"`
	VLAN string `yaml:"vlan"`
}

type valuesWithStrictACLs struct {
	ACLs  []ACL                `yaml:"acls"`
	Other map[string]yaml.Node `yaml:",inline"`
}

func UnmarshalValues(in []byte) (*Values, error) {
	var values Values
	if err := yaml.Unmarshal(in, &values); err != nil {
		return nil, err
	}
	if err := checkTopLevelKeys(in); err != nil {
		return nil, withUnknownFieldHint(err)
	}
	if err := checkACLFields(in); err != nil {
		return nil, withUnknownFieldHint(err)
	}
	return &values, nil
}

func withUnknownFieldHint(err error) error {
	return fmt.Errorf("%w\nfix or remove the fields named above, the top-level keys and the acl fields are described in the README", err)
}

func checkACLFields(in []byte) error {
	decoder := yaml.NewDecoder(bytes.NewReader(in))
	decoder.KnownFields(true)
	err := decoder.Decode(&valuesWithStrictACLs{})
	if errors.Is(err, io.EOF) {
		return nil
	}
	return err
}

func checkTopLevelKeys(in []byte) error {
	var document yaml.Node
	if err := yaml.Unmarshal(in, &document); err != nil || len(document.Content) == 0 {
		return err
	}

	known := topLevelKeys()
	var unknown []string
	for _, key := range mappingKeys(document.Content[0]) {
		if !known[key.Value] {
			unknown = append(unknown, fmt.Sprintf("line %d: field %s not found in type values.Values", key.Line, key.Value))
		}
	}
	if len(unknown) > 0 {
		return &yaml.TypeError{Errors: unknown}
	}
	return nil
}

func topLevelKeys() map[string]bool {
	keys := map[string]bool{}
	valuesType := reflect.TypeOf(Values{})
	for i := 0; i < valuesType.NumField(); i++ {
		name, _, _ := strings.Cut(valuesType.Field(i).Tag.Get("yaml"), ",")
		keys[name] = true
	}
	return keys
}

func mappingKeys(node *yaml.Node) []*yaml.Node {
	switch node.Kind {
	case yaml.AliasNode:
		return mappingKeys(node.Alias)
	case yaml.SequenceNode:
		var keys []*yaml.Node
		for _, merged := range node.Content {
			keys = append(keys, mappingKeys(merged)...)
		}
		return keys
	case yaml.MappingNode:
		var keys []*yaml.Node
		for i := 0; i+1 < len(node.Content); i += 2 {
			key, value := node.Content[i], node.Content[i+1]
			if key.Tag == "!!merge" {
				keys = append(keys, mappingKeys(value)...)
				continue
			}
			keys = append(keys, key)
		}
		return keys
	}
	return nil
}
