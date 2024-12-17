package configdb

import (
	"github.com/metal-stack/sonic-configdb-utils/values"
)

type ConfigDB struct {
	ACLRules          map[string]ACLRule        `json:"ACL_RULE"`
	ACLTables         map[string]ACLTable       `json:"ACL_TABLE"`
	Breakouts         map[string]BreakoutConfig `json:"BREAKOUT_CFG"`
	DeviceMetadata    `json:"DEVICE_METADATA"`
	Features          map[string]Feature `json:"FEATURE"`
	LLDP              `json:"LLDP"`
	LoopbackInterface map[string]interface{}   `json:"LOOPBACK_INTERFACE"`
	MgmtInterfaces    map[string]MgmtInterface `json:"MGMT_INTERFACE"`
	MgmtPorts         map[string]MgmtPort      `json:"MGMT_PORT"`
	MgmtVRFConfig     `json:"MGMT_VRF_CONFIG"`
	NTP               `json:"NTP"`
	NTPServers        map[string]interface{} `json:"NTP_SERVER"`
	Ports             map[string]Port        `json:"PORT"`
	VLANs             map[string]VLAN        `json:"VLAN"`
	VLANInterface     map[string]interface{} `json:"VLAN_INTERFACE"`
	VXLANEVPN         `json:"VXLAN_EVPN_NVO"`
	VXLANTunnels      map[string]VXLANTunnel    `json:"VXLAN_TUNNEL"`
	VXLANTunnelMaps   map[string]VXLANTunnelMap `json:"VXLAN_TUNNEL_MAP"`
}

func GenerateConfigDB(input *values.Values) *ConfigDB {
	configdb := ConfigDB{}
	return &configdb
}
