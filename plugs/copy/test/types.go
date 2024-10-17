package test

import (
	"fmt"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

//go:generate gowrap gen -g -p ./ -i Service -bt ""
func testCopy() {
	var source NginxDomain
	// source.Domain = "hello"
	var target NginxBody

	/* 	source.MapStruct = map[string]Password{
		"hello": Password{PasswordName: "hello"},
	} */

	// @copy
	stCopy(&target, &source)

	var sliceSource []NginxDomain
	var sliceTarget []NginxBody

	lo.ForEach(sliceSource, func(item NginxDomain, index int) {
		var nginxBody NginxBody
		// @copy
		nginxDomain2NginxBodyDTO(&nginxBody, &item)
		sliceTarget = append(sliceTarget, nginxBody)
	})
	// sliceStCopy(&sliceTarget, &sliceSource)

	var mapSource map[string]NginxDomain
	var mapTarget map[string]NginxBody

	// mapStCopy(&mapTarget, &mapSource)
	for k, v := range mapSource {
		var nginxBody NginxBody
		// @copy
		nginxDomain2NginxBodyDTO(&nginxBody, &v)
		mapTarget[k] = nginxBody
	}

	// var sliceLabelSource []Label
	var sliceSelectTarget []Select

	// @copy
	// sliceLabelCopy(&sliceSelectTarget, &sliceLabelSource)

	fmt.Println(target)
	fmt.Println(sliceTarget)
	fmt.Println(mapTarget)
	fmt.Println(sliceSelectTarget)
}

type Label struct {
	Name string
}

type Select struct {
	// @copy-target-name Name
	Label string
	// @copy-target-name Name
	Value string
}

type NginxBody struct {
	Domain    string `gorm:"column:domain;comment:域名"`
	Cluster   string `gorm:"column:cluster;comment:集群"`
	ProjectID int    `gorm:"column:project_id;comment:项目ID"`
	Project   string `gorm:"column:project;comment:项目"`
	// @copy-target-method GetName
	ProjectCname string
	Product      string `gorm:"column:product;comment:产品"`
	ProductCname string
	Conf         NginxServerConf `gorm:"embedded;embeddedPrefix:conf_"`
	// @copy-target-path Conf.Domain
	ConfDomain string

	// @copy-prefix Copy
	CopyServceConf

	Map map[string]string

	// @copy-name MapStruct
	MapStruct100 map[string]*Password2
	MapStruct2   map[string]Password2
	SliceStruct  []Password2
	SliceStruct2 []*Password2

	// NoSame []string
}

func (n *NginxBody) GetName() string {
	return fmt.Sprintf("%s-%s", n.Domain, n.Cluster)
}

// source
type NginxDomain struct {
	gorm.Model
	Domain     string                     `gorm:"column:domain;comment:域名"`
	Cluster    string                     `gorm:"column:cluster;comment:集群"`
	ProjectID  int                        `gorm:"column:project_id;comment:项目ID"`
	Project    string                     `gorm:"column:project;comment:项目"`
	Product    string                     `gorm:"column:product;comment:产品"`
	Conf       NginxDomainNginxServerConf `gorm:"embedded;embeddedPrefix:conf_"`
	EnableSync bool                       `gorm:"column:enable_sync;comment:是否开启同步"`
	// 健康:healthy, 不健康:unhealthy, 未知:unknown
	MirrorSyncStatus string `gorm:"column:mirror_sync_status;comment:镜像同步状态"`
	Map              map[string]string
	MapStruct        map[string]Password
	MapStruct2       map[string]*Password
	SliceStruct      []*Password
	SliceStruct2     []Password
	// NoSame           []Password
}

func (n *NginxDomain) GetName() string {
	return fmt.Sprintf("%s-%s", n.Domain, n.Cluster)
}

type CopyServceConf struct {
	CopyDomain string
	CopyServer string
}

type NginxServerConf struct {
	Domain              string         `gorm:"column:domain;comment:域名"`
	Server              string         `gorm:"column:server;comment:server"`
	Listens             []MultListens  `gorm:"column:listens;serializer:json;comment:端口"`
	ListensSSL          []MultListens  `gorm:"column:listens_ssl;serializer:json;comment:ssl端口"`
	ClientMaxBodySize   string         `gorm:"column:client_max_body_size;comment:client_max_body_size"`
	Rewrite             string         `gorm:"column:rewrite;comment:rewrite"`
	UpstreamItems       []UpstreamItem `gorm:"foreignKey:ConfID;comment:upstream"`
	Locations           []LocationItem `gorm:"foreignKey:ConfID;comment:location"`
	Enable              bool           `gorm:"column:enable;comment:是否删除"`
	SSLOn               bool           `gorm:"column:ssl_on;comment:是否开启ssl"`
	SSLCertificate      string         `gorm:"column:ssl_certificate;comment:ssl证书"`
	SSLCertificateKey   string         `gorm:"column:ssl_certificate_key;comment:ssl证书key"`
	ErrorPage           bool           `gorm:"column:error_page;comment:error_page"`
	ServiceType         string         `gorm:"column:service_type;comment:服务类型"`
	EnableLimitConnZone bool           `gorm:"column:enable_limit_conn_zone;comment:是否开启limit_conn_zone"`
	LimitConnPerServer  int            `gorm:"column:limit_conn_per_server;comment:limit_conn_per_server"`
}

type UpstreamItem struct {
	ID                   uint                 `json:"id"`
	ConfID               uint                 `gorm:"index"`
	Name                 string               `gorm:"column:name;comment:名称"`
	LoadbalanceType      string               `gorm:"column:loadbalance_type;comment:负载均衡类型"`
	LoadbalanceValue     string               `gorm:"column:loadbalance_value;comment:负载均衡值"`
	PodServers           PodAutoUpstream      `gorm:"embedded;embeddedPrefix:pod_servers_"`
	Servers              []UpstreamServerItem `gorm:"column:servers;serializer:json;comment:upstream_server"`
	CheckHTTPSend        string               `gorm:"column:check_http_send;comment:check_http_send"`
	CheckHTTPExpectAlive string               `gorm:"column:check_http_expect_alive;comment:check_http_expect_alive"`

	CheckInterval int
	CheckRise     int
	CheckFall     int
	CheckTimeout  int

	EtcdServers []UpstreamServerItem

	Effect []*UpstreamEffect `json:"EffectV2"`

	//NodePorts []UpstreamItemNodePort `gorm:"foreignKey:UpstreamItemID;comment:nodeport"`
}

type PodAutoUpstream struct {
	SyncFromPod bool `gorm:"column:sync_from_pod;comment:是否同步pod"`
	Port        int  `gorm:"column:port;comment:pod端口"`
}

type UpstreamEffect struct {
	Ip         string
	Consistent bool
	Content    []UpstreamServerItem
}

type UpstreamServerItem struct {
	Weight  int
	HP      string // host:port
	Flag    string // 标记服务是否可使用
	Healthy int
	FromPod bool // 直连pod
}

type LocationItem struct {
	ID               uint          `json:"id"`
	ConfID           uint          `gorm:"index"`
	Key              string        `gorm:"column:key;comment:location_path"`
	UpstreamName     string        `gorm:"column:upstream_name;comment:upstream_name"`
	HeaderHost       string        `gorm:"column:header_host;comment:header_host"`
	SubDirectoryPath string        `gorm:"column:sub_directory_path;comment:sub_directory_path"`
	Rewrite          string        `gorm:"column:rewrite;comment:rewrite"`
	LimitReqZone     LimitReqZone  `gorm:"embedded;embeddedPrefix:limit_req_zone_"`
	LimitConnZone    LimitConnZone `gorm:"embedded;embeddedPrefix:limit_conn_zone_"`
	Huidu            Huidu         `gorm:"embedded;embeddedPrefix:huidu_"`
}
type LimitReqZone struct {
	Enable bool   `gorm:"column:enable"`
	Zone   string `gorm:"column:zone"`
	Burst  int    `gorm:"column:burst"`
}

type LimitConnZone struct {
	Enable    bool `gorm:"column:enable"`
	PerServer int  `gorm:"column:per_server"`
}

type Huidu struct {
	Enable        bool        `gorm:"column:enable"`
	Upstream      string      `gorm:"column:upstream"`
	Upstreamhuidu string      `gorm:"column:upstreamhuidu"`
	HuiduKey      string      `gorm:"column:huidu_key"`
	HeaderHuidu   HeaderHuidu `gorm:"column:header_huidu;serializer:json;"`
	IPHuidu       IPHuidu     `gorm:"column:ip_huidu;serializer:json;"`
	ArgsHuidu     ArgsHuidu   `gorm:"column:args_huidu;serializer:json;"`
	Content       []HuiduKV   `gorm:"column:content;serializer:json;"`
}

type Password struct {
	PasswordName string
	NestSlice    []PasswordNest
	Nest         PasswordNest
	NestMap      map[string]PasswordNest
}

type PasswordNest struct {
	// @copy-name Ipone
	Ipone string `json:"ipone"`
}

type Password2 struct {
	PasswordName string
	NestSlice    []Password2Nest
	Nest         Password2Nest
	NestMap      map[string]Password2Nest
}

type Password2Nest struct {
	// @copy-name Ipone
	Ipone1 string `json:"ipone"`
}

type NginxDomainNginxServerConf struct {
	Domain              string                    `gorm:"column:domain;comment:域名"`
	Server              string                    `gorm:"column:server;comment:server"`
	Listens             []MultListens             `gorm:"column:listens;serializer:json;comment:端口"`
	ListensSSL          []MultListens             `gorm:"column:listens_ssl;serializer:json;comment:ssl端口"`
	ClientMaxBodySize   string                    `gorm:"column:client_max_body_size;comment:client_max_body_size"`
	Rewrite             string                    `gorm:"column:rewrite;comment:rewrite"`
	UpstreamItems       []NginxDomainUpstreamItem `gorm:"foreignKey:ConfID;comment:upstream"`
	Locations           []NginxDomainLocationItem `gorm:"foreignKey:ConfID;comment:location"`
	Enable              bool                      `gorm:"column:enable;comment:是否删除"`
	SSLOn               bool                      `gorm:"column:ssl_on;comment:是否开启ssl"`
	SSLCertificate      string                    `gorm:"column:ssl_certificate;comment:ssl证书"`
	SSLCertificateKey   string                    `gorm:"column:ssl_certificate_key;comment:ssl证书key"`
	ErrorPage           bool                      `gorm:"column:error_page;comment:error_page"`
	ServiceType         string                    `gorm:"column:service_type;comment:服务类型"`
	EnableLimitConnZone bool                      `gorm:"column:enable_limit_conn_zone;comment:是否开启limit_conn_zone"`
	LimitConnPerServer  int                       `gorm:"column:limit_conn_per_server;comment:limit_conn_per_server"`
}

type MultListens struct {
	Listen int
}

type NginxDomainUpstreamItem struct {
	gorm.Model
	ConfID               uint                            `gorm:"index"`
	Name                 string                          `gorm:"column:name;comment:名称"`
	LoadbalanceType      string                          `gorm:"column:loadbalance_type;comment:负载均衡类型"`
	LoadbalanceValue     string                          `gorm:"column:loadbalance_value;comment:负载均衡值"`
	PodServers           NginxDomainPodAutoUpstream      `gorm:"embedded;embeddedPrefix:pod_servers_"`
	Servers              []NginxDomainUpstreamServerItem `gorm:"column:servers;serializer:json;comment:upstream_server"`
	CheckHTTPSend        string                          `gorm:"column:check_http_send;comment:check_http_send"`
	CheckHTTPExpectAlive string                          `gorm:"column:check_http_expect_alive;comment:check_http_expect_alive"`
	CheckInterval        int                             `gorm:"column:check_interval;comment:check_interval`
	CheckRise            int                             `gorm:"column:check_rise;comment:check_rise"`
	CheckFall            int                             `gorm:"column:check_fall;comment:check_fall"`
	CheckTimeout         int                             `gorm:"column:check_timeout;comment:check_timeout"`

	NodePorts []UpstreamItemNodePort `gorm:"foreignKey:UpstreamItemID;comment:nodeport"`
}

type UpstreamItemNodePort struct {
	gorm.Model
	UpstreamItemID uint   `gorm:"index"`
	NodePort       string `gorm:"column:node_port;comment:node_port"`
}

func (UpstreamItemNodePort) TableName() string {
	return "paas_api_nginx_upstream_item_nodeport"
}

type NginxDomainLocationItem struct {
	gorm.Model
	ConfID           uint                     `gorm:"index"`
	Key              string                   `gorm:"column:key;comment:location_path"`
	UpstreamName     string                   `gorm:"column:upstream_name;comment:upstream_name"`
	HeaderHost       string                   `gorm:"column:header_host;comment:header_host"`
	SubDirectoryPath string                   `gorm:"column:sub_directory_path;comment:sub_directory_path"`
	Rewrite          string                   `gorm:"column:rewrite;comment:rewrite"`
	LimitReqZone     NginxDomainLimitReqZone  `gorm:"embedded;embeddedPrefix:limit_req_zone_"`
	LimitConnZone    NginxDomainLimitConnZone `gorm:"embedded;embeddedPrefix:limit_conn_zone_"`
	Huidu            NginxDomainHuidu         `gorm:"embedded;embeddedPrefix:huidu_"`
}

func (l LocationItem) DiffKey() uint {
	return l.ID
}

func (LocationItem) TableName() string {
	return "paas_api_nginx_location_item"
}

type NginxDomainLimitReqZone struct {
	Enable bool   `gorm:"column:enable"`
	Zone   string `gorm:"column:zone"`
	Burst  int    `gorm:"column:burst"`
}

type NginxDomainLimitConnZone struct {
	Enable    bool `gorm:"column:enable"`
	PerServer int  `gorm:"column:per_server"`
}

type NginxDomainHuidu struct {
	Enable        bool        `gorm:"column:enable"`
	Upstream      string      `gorm:"column:upstream"`
	Upstreamhuidu string      `gorm:"column:upstreamhuidu"`
	HuiduKey      string      `gorm:"column:huidu_key"`
	HeaderHuidu   HeaderHuidu `gorm:"column:header_huidu;serializer:json;"`
	IPHuidu       IPHuidu     `gorm:"column:ip_huidu;serializer:json;"`
	ArgsHuidu     ArgsHuidu   `gorm:"column:args_huidu;serializer:json;"`
	Content       []HuiduKV   `gorm:"column:content;serializer:json;"`
}

type HuiduPair struct {
	Key, Value string
}

type HeaderHuidu struct {
	Header []HuiduPair
}

type IPHuidu struct {
	Ips []string
}

type ArgsHuidu struct {
	Args []HuiduPair
}

type HuiduKV struct {
	K, V string
}

type NginxDomainPodAutoUpstream struct {
	SyncFromPod bool `gorm:"column:sync_from_pod;comment:是否同步pod"`
	Port        int  `gorm:"column:port;comment:pod端口"`
}

type NginxDomainUpstreamServerItem struct {
	Custom        bool
	Product       string
	Project       string
	ContainerPort int32
	ClusterId     string
	NodePort      int32

	Weight  int
	HP      string // host:port
	Flag    string // 标记服务是否可使用
	Healthy int
	FromPod bool // 直连pod
}
