package runner

type commitDef struct {
	MultiRow  bool   `json:"multirow"`
	URL       string `json:"url"`
	URLSuffix string `json:"url_suffix"`
}

var commitDefs = map[string]commitDef{
	"server": {
		MultiRow:  false,
		URL:       "/api/servers/servers/",
		URLSuffix: "-/sync/",
	},
	"info": {
		MultiRow:  false,
		URL:       "/api/proc/info/",
		URLSuffix: "-/sync/",
	},
	"os": {
		MultiRow:  false,
		URL:       "/api/proc/os/",
		URLSuffix: "-/sync/",
	},
	"time": {
		MultiRow:  false,
		URL:       "/api/proc/time/",
		URLSuffix: "-/sync/",
	},
	"groups": {
		MultiRow:  true,
		URL:       "/api/proc/groups/",
		URLSuffix: "sync/",
	},
	"users": {
		MultiRow:  true,
		URL:       "/api/proc/users/",
		URLSuffix: "sync/",
	},
	"interfaces": {
		MultiRow:  true,
		URL:       "/api/proc/interfaces/",
		URLSuffix: "sync/",
	},
	"addresses": {
		MultiRow:  true,
		URL:       "/api/proc/addresses/",
		URLSuffix: "sync/",
	},
	"packages": {
		MultiRow:  true,
		URL:       "/api/proc/packages/",
		URLSuffix: "sync/",
	},
}

type ServerData struct {
	Version string  `json:"version"`
	Load    float64 `json:"load"`
}

type SystemData struct {
	ID               string `json:"id,omitempty"`
	UUID             string `json:"uuid"`
	CPUType          string `json:"cpu_type"`
	CPUBrand         string `json:"cpu_brand"`
	CPUPhysicalCores int    `json:"cpu_physical_cores"`
	CPULogicalCores  int    `json:"cpu_logical_cores"`
	PhysicalMemory   uint64 `json:"physical_memory"`
	HardwareVendor   string `json:"hardware_vendor"`
	HardwareModel    string `json:"hardware_model"`
	HardwareSerial   string `json:"hardware_serial"`
	ComputerName     string `json:"computer_name"`
	Hostname         string `json:"hostname"`
	LocalHostname    string `json:"local_hostname"`
}

type OSData struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name"`
	Version      string `json:"version"`
	Major        int    `json:"major"`
	Minor        int    `json:"minor"`
	Patch        int    `json:"patch"`
	Platform     string `json:"platform"`
	PlatformLike string `json:"platform_like"`
}

type TimeData struct {
	ID       string `json:"id,omitempty"`
	Datetime string `json:"datetime"`
	BootTime uint64 `json:"boot_time"`
	Timezone string `json:"timezone"`
	Uptime   uint64 `json:"uptime"`
}

type UserData struct {
	ID          string `json:"id,omitempty"`
	UID         int    `json:"uid"`
	GID         int    `json:"gid"`
	Username    string `json:"username"`
	Description string `json:"description"`
	Directory   string `json:"directory"`
	Shell       string `json:"shell"`
}

type GroupData struct {
	ID        string `json:"id,omitempty"`
	GID       int    `json:"gid"`
	GroupName string `json:"groupname"`
}

type SystemPackageData struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Source  string `json:"source"`
	Arch    string `json:"arch"`
}

type Interface struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name"`
	Mac       string `json:"mac"`
	Type      int    `json:"type"`
	Flags     int    `json:"flags"`
	MTU       int    `json:"mtu"`
	LinkSpeed int    `json:"link_speed"`
}

type Address struct {
	ID            string `json:"id,omitempty"`
	Address       string `json:"address"`
	Broadcast     string `json:"broadcast"`
	InterfaceName string `json:"interface_name,omitempty"`
	Mask          string `json:"mask"`
}

type commitData struct {
	Version    string              `json:"version"`
	Load       float64             `json:"load"`
	Info       SystemData          `json:"info"`
	OS         OSData              `json:"os"`
	Time       TimeData            `json:"time"`
	Users      []UserData          `json:"users"`
	Groups     []GroupData         `json:"groups"`
	Interfaces []Interface         `json:"interfaces"`
	Addresses  []Address           `json:"addresses"`
	Packages   []SystemPackageData `json:"packages"`
}

// Defines the ComparableData interface for comparing different types.
// Ensures data retrieval for each key, excluding the ID field, while minimizing the use of reflection for better performance.
type ComparableData interface {
	GetID() string
	GetKey() interface{}
	GetData() ComparableData
}

func (s SystemData) GetID() string {
	return s.ID
}

func (s SystemData) GetKey() interface{} {
	return s.UUID
}

func (s SystemData) GetData() ComparableData {
	return SystemData{
		UUID:             s.UUID,
		CPUType:          s.CPUType,
		CPUBrand:         s.CPUBrand,
		CPUPhysicalCores: s.CPUPhysicalCores,
		CPULogicalCores:  s.CPULogicalCores,
		PhysicalMemory:   s.PhysicalMemory,
		HardwareVendor:   s.HardwareVendor,
		HardwareModel:    s.HardwareModel,
		HardwareSerial:   s.HardwareSerial,
		ComputerName:     s.ComputerName,
		Hostname:         s.Hostname,
		LocalHostname:    s.LocalHostname,
	}
}

func (o OSData) GetID() string {
	return o.ID
}

func (o OSData) GetKey() interface{} {
	return o.Name
}

func (o OSData) GetData() ComparableData {
	return OSData{
		Name:         o.Name,
		Version:      o.Version,
		Major:        o.Major,
		Minor:        o.Minor,
		Patch:        o.Patch,
		Platform:     o.Platform,
		PlatformLike: o.PlatformLike,
	}
}

func (t TimeData) GetID() string {
	return t.ID
}

func (t TimeData) GetKey() interface{} {
	return t.Timezone
}

func (t TimeData) GetData() ComparableData {
	return TimeData{
		Datetime: t.Datetime,
		Timezone: t.Timezone,
		Uptime:   t.Uptime,
	}
}

func (u UserData) GetID() string {
	return u.ID
}

func (u UserData) GetKey() interface{} {
	return u.Username
}

func (u UserData) GetData() ComparableData {
	return UserData{
		Username:  u.Username,
		UID:       u.UID,
		GID:       u.GID,
		Directory: u.Directory,
		Shell:     u.Shell,
	}
}

func (g GroupData) GetID() string {
	return g.ID
}

func (g GroupData) GetKey() interface{} {
	return g.GID
}

func (g GroupData) GetData() ComparableData {
	return GroupData{
		GID:       g.GID,
		GroupName: g.GroupName,
	}
}

func (i Interface) GetID() string {
	return i.ID
}

func (i Interface) GetKey() interface{} {
	return i.Name
}

func (i Interface) GetData() ComparableData {
	return Interface{
		Name:      i.Name,
		Mac:       i.Mac,
		Type:      i.Type,
		Flags:     i.Flags,
		MTU:       i.MTU,
		LinkSpeed: i.LinkSpeed,
	}
}

func (a Address) GetID() string {
	return a.ID
}

func (a Address) GetKey() interface{} {
	return a.Address
}

func (a Address) GetData() ComparableData {
	return Address{
		Address:       a.Address,
		Broadcast:     a.Broadcast,
		InterfaceName: a.InterfaceName,
		Mask:          a.Mask,
	}
}

func (sp SystemPackageData) GetID() string {
	return sp.ID
}

func (sp SystemPackageData) GetKey() interface{} {
	return sp.Name
}

func (sp SystemPackageData) GetData() ComparableData {
	return SystemPackageData{
		Name:    sp.Name,
		Version: sp.Version,
		Source:  sp.Source,
		Arch:    sp.Arch,
	}
}
