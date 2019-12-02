package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/ansiz/beego"
)

var (
	cfg = beego.AppConfig
)

const (
	devMode = "dev"
)

// RedisConfig represents Redis connect parameters.
type RedisConfig struct {
	Host           string
	Port           int
	Password       string
	Db             int
	MaxIdle        int
	MaxActive      int
	IdleTimeout    int64
	ConnectTimeout int64
}

// MailConfig represents the email configurations.
type MailConfig struct {
	Host     string
	Port     int
	Auth     string
	Username string
	Adapter  string
	Enable   bool
}

// MySQLConfig represents MySql parameters.
type MySQLConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Db       string
	Debug    bool
	Schema   string
}

// cacheConfig represents the configuration of cache module that in Redis mode.
type cacheConfig struct {
	Key      string `json:"key"`
	Conn     string `json:"conn"`
	DBNum    string `json:"dbNum"`
	Password string `json:"password"`
}

// IPMIConfig represents the IPMI related configurations.
type IPMIConfig struct {
	Interface string
	User      string
	Password  string
}

// BackupConfig represents the configuration of backup.
type BackupConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

// BackendConfig represents backend config read from app.conf.
type BackendConfig struct {
	SchedulerType   string
	PowerSeparation bool
}

// OEMInfo includes product customization info.
type OEMInfo struct {
	DescHPC     string
	DescBackup  string
	Copyright   string
	ProductName string
}

// Commit contains the code commit-id.
type Commit struct {
	Frontend string
	Backend  string
}

// FileLoggerConfig includes file logger related configurations.
type FileLoggerConfig struct {
	Filename string
	Daily    bool // logrotate daily
	MaxDays  int
}

// GetBackendConfig returns backend config read from app.conf.
func GetBackendConfig() BackendConfig {
	config := BackendConfig{
		SchedulerType:   GetSchedulerType(),
		PowerSeparation: cfg.DefaultBool("PowerSeparation", false),
	}
	return config
}

// GetRedisConfig returns default redis config load from config file.
func GetRedisConfig() *RedisConfig {
	return &RedisConfig{
		Db:             cfg.DefaultInt("redis::RedisDB", 0),
		Password:       cfg.String("redis::RedisPasswd"),
		Host:           cfg.String("redis::RedisHost"),
		Port:           cfg.DefaultInt("redis::RedisPort", 6379),
		MaxIdle:        cfg.DefaultInt("redis::RedisMaxIdle", 100),
		MaxActive:      cfg.DefaultInt("redis::RedisMaxActive", 0),
		IdleTimeout:    cfg.DefaultInt64("redis::IdleTimeout", 120),
		ConnectTimeout: cfg.DefaultInt64("redis::ConnectTimeout", 60),
	}
}

// HPCRootPath returns the HPC_ROOT environment variable.
func HPCRootPath() string {
	return os.Getenv("HPC_ROOT")
}

// MgmtRedisConfig returns the management node's redis config.
func MgmtRedisConfig() *RedisConfig {
	config := GetRedisConfig()
	config.Host = cfg.DefaultString("agent::RedisHost", "localhost")
	config.Port = cfg.DefaultInt("agent::RedisPort", 6379)
	return config
}

// HealthCheckInterval returns the interval of health check.
func HealthCheckInterval() time.Duration {
	interval := cfg.DefaultInt("agent::HealthCheckInterval", 30)
	return time.Second * time.Duration(interval)
}

// JobSubmitTimeout returns the max timeout for job submission.
func JobSubmitTimeout() time.Duration {
	interval := cfg.DefaultInt("JobSubmitTimeout", 30)
	return time.Second * time.Duration(interval)
}

// HeartbeatInterval returns the interval of heartbeat.
func HeartbeatInterval() time.Duration {
	interval := cfg.DefaultInt("agent::HeartBeatInterval", 10)
	return time.Second * time.Duration(interval)
}

// ContainerGCInterval returns the interval of container garbage collection.
func ContainerGCInterval() time.Duration {
	interval := cfg.DefaultInt("agent::ContainerGCInterval", 7200)
	return time.Second * time.Duration(interval)
}

// MinContainerAge represents the minimum age of a completed job before its
// record and related container is purged. Default value is 86400(one day).
func MinContainerAge() time.Duration {
	interval := cfg.DefaultInt("agent::MinContainerAge", 86400)
	return time.Second * time.Duration(interval)
}

// ServicesMonitScriptsDir returns the services monit scripts directory.
func ServicesMonitScriptsDir() string {
	return AbsPath(cfg.DefaultString("agent::ServicesMonitScriptsDir",
		"monit"))
}

// CacheRedisConfig returns the configuration of cache module that in Redis mode.
func CacheRedisConfig() (string, error) {
	rc := GetRedisConfig()
	cc := cacheConfig{
		Key:      cfg.DefaultString("cache_key", "hpc:cache"),
		DBNum:    strconv.Itoa(rc.Db),
		Password: rc.Password,
		Conn:     fmt.Sprintf("%s:%d", rc.Host, rc.Port),
	}
	data, err := json.Marshal(cc)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// GetAuthProvider returns auth provider type.
func GetAuthProvider() string {
	return cfg.DefaultString("backend::AuthProvider", "basic")
}

// IsDevMode judges devMode.
func IsDevMode() bool {
	return cfg.String("runmode") == devMode
}

// GetSchedulerType returns scheduler type.
func GetSchedulerType() string {
	return cfg.DefaultString("backend::SchedulerType", "slurm")
}

// GetMonitorProvider returns monitor provider config in app.conf.
func GetMonitorProvider() string {
	return cfg.DefaultString("backend::MonitorProvider", "ganglia")
}

// GetOEMInfo returns OEM info for customization.
func GetOEMInfo(namespace string) OEMInfo {
	copyright, _ := ioutil.ReadFile(AboutFile())
	info := OEMInfo{
		Copyright:  string(copyright),
		DescHPC:    cfg.DefaultString("DescHPC", "HPC Intergrated Environment"),
		DescBackup: cfg.DefaultString("DescBackup", "Professional Backup System"),
	}
	if namespace == "backup" {
		info.ProductName = cfg.DefaultString("BackupProductName", "Lankup")
	} else {
		info.ProductName = cfg.DefaultString("HPCProductName", "SuperHPC")
	}
	return info
}

// AboutFile returns production introduction.
func AboutFile() string {
	return AbsPath(cfg.DefaultString("AboutFile", "ABOUT"))
}

// GetDictsPath returns the dictionary directory.
func GetDictsPath() string {
	return cfg.DefaultString("trimps::DictPath", "/trimps/dicts/")
}

// GetTempDictsPath returns the temporary dictionary directory.
func GetTempDictsPath() string {
	return cfg.DefaultString("trimps::TempDictPath", "/trimps/tmpdicts/")
}

// PluginsPath returns the plugin directory.
func PluginsPath() string {
	return cfg.DefaultString("trimps::PluginPath", "/trimps/plugins/")
}

// GetCrackFileTempPath returns the crackfile's temporary directory, the file
// will be moved to crackfiles path after the mission started.
func GetCrackFileTempPath() string {
	return cfg.DefaultString("trimps::crack_file_temp_path", "/trimps/cracktmpfiles/")
}

// CrackFilePath returns the crackfiles directory.
func CrackFilePath() string {
	return cfg.DefaultString("trimps::crack_file_temp_path", "/trimps/crackfiles/")
}

// GetGlobalFileLogger returns the global file logger's config.
func GetGlobalFileLogger() (string, error) {
	config := cfg.DefaultString("global::FileLogger",
		`{"filename":"/var/log/hpc-backend/hpc.log","daily":true,"maxdays":10}`)
	return checkFileLoggerConfig(config)
}

// checkFileLoggerConfig validates file logger configurations, the content
// should be a JSON text that can be unmarshal to FileLoggerConfig, if
// specified path is not exist, it will create it automatically.
func checkFileLoggerConfig(text string) (string, error) {
	config := FileLoggerConfig{}
	err := json.Unmarshal([]byte(text), &config)
	if err != nil {
		return "", fmt.Errorf("parse logger config failed: %v", err)
	}
	return text, os.MkdirAll(path.Dir(config.Filename), 0755)
}

// RedisLoggerOn determines whether output log to Redis, if it is true,
// global::RedisLogger configurations is required.
func RedisLoggerOn() bool {
	return cfg.DefaultBool("global::RedisLoggerOn", false)
}

// GetGlobalRedisLogger returns the global redis logger's config.
func GetGlobalRedisLogger() string {
	return cfg.DefaultString("global::RedisLogger",
		`{"key":"hpc-logs","level":7,"addr":"127.0.0.1","host":"master"}`)
}

// JobRootPath returns trimps crack file store path.
func JobRootPath() string {
	return cfg.DefaultString("trimps::JobPath", "/trimps/jobs/")
}

// SplitCommand returns split command.
func SplitCommand() string {
	return cfg.DefaultString("trimps::SplitCommand", "java -jar /trimps/utils/split.jar")
}

// AgentTTYPort returns agent's tty listened port.
func AgentTTYPort() int {
	return cfg.DefaultInt("agent::TTYPort", 8080)
}

// RedisAuth returns redis authorization token.
func RedisAuth() string {
	return cfg.String("trimps::redis_auth")
}

// SpeedInterval returns the time interval of crack speed detection.
func SpeedInterval() int {
	return cfg.DefaultInt("trimps::speed_interval", 5)
}

// MgmtHost returns the host of mgmt node.
func MgmtHost() string {
	return cfg.DefaultString("backend::MgmtHost", "master")
}

// RedisHost returns the IP of mgmt node.
func RedisHost() string {
	return cfg.String("agent::RedisHost")
}

// MgmtSubnetMask returns the subnet mask of mgmt node.
func MgmtSubnetMask() string {
	return cfg.DefaultString("agent::MgmtSubnetMask", "255.255.255.0")
}

// CrackScript returns the crack scripts path.
func CrackScript() string {
	return cfg.DefaultString("trimps::crack_script", "/trimps/utils/crack.sh")
}

// TrimpsWorkDir returns the crack program's workspace path. for Trimps only.
func TrimpsWorkDir() string {
	return cfg.DefaultString("trimps::work_dir", "/tmp/trimps")
}

// ClusterName returns cluster's name.
func ClusterName() string {
	return cfg.DefaultString("backend::ClusterName", "dev")
}

// RetryTimes returns the number of max retries.
func RetryTimes() int {
	return cfg.DefaultInt("trimps::retry_times", 5)
}

// NISServer returns the NIS server's IP or hostname.
func NISServer() string {
	return cfg.DefaultString("nis::NISServer", "localhost")
}

// NISAdmin returns the username which has privilege to update NIS data.
func NISAdmin() string {
	return cfg.DefaultString("nis::NISAdmin", "root")
}

// NISMaker returns the program/script name which be used to generate NIS data,
// default value is the command "make -C /var/yp".
func NISMaker() string {
	return cfg.DefaultString("nis::NISMaker", "make -C /var/yp")
}

// PhotoTempPath returns the photo temporary path.(For WHU only)
func PhotoTempPath() string {
	path := cfg.DefaultString("whu::photo_temp_path", "/tmp/photo")
	// Because it will be cleaned up regularly, so the path must be limited in /tmp
	if !strings.HasPrefix(path, "/tmp/") {
		fmt.Println("Warning: incorrect photo temp path has been correctted to /tmp/photo")
		path = "tmp/photo"
	}
	return path
}

// AppDataPath returns the app root path.
func AppDataPath() string {
	return cfg.DefaultString("AppData", "/usr/hpc")
}

// PhotoPath returns the photo relative path.(For WHU only)
func PhotoPath() string {
	return cfg.DefaultString("PhotoPath", path.Join(StoragePath(), photoDirName()))
}

// PhotoURLPattern returns the URL pattern of photo.
func PhotoURLPattern() string {
	return path.Join(StorageURL(), photoDirName())
}

func photoDirName() string {
	return cfg.DefaultString("PhotoDirName", "images")
}

// StorageURL returns the URL pattern of storage.
func StorageURL() string {
	return cfg.DefaultString("StorageURL", "/storage")
}

// StoragePath returns the storage path.
func StoragePath() string {
	return cfg.DefaultString("StoragePath", path.Join(AppDataPath(), "/storage"))
}

// LicensePlugin returns license plugin path.
func LicensePlugin() string {
	return path.Join(LicensePath(), cfg.DefaultString("backend::LicensePlugin", "license.so"))
}

// License returns license path.
func License() string {
	return path.Join(LicensePath(), cfg.DefaultString("backend::License", "license.lic"))
}

// LicensePath returns the license storage path.
func LicensePath() string {
	return cfg.DefaultString("backend::LicensePath", "/usr/hpc/license")
}

// Email returns mail related config.
func Email() string {
	return cfg.DefaultString("backend::Email", `{"Enable": false}`)
}

// MailCMD returns the mail sender command.
func MailCMD() string {
	return cfg.DefaultString("MailCMD", "mail")
}

// UserPrefix returns the prefix of username which will be used for adding
// Linux system user.
func UserPrefix() string {
	return cfg.DefaultString("user_prefix", "hpc")
}

// LiteRBAC controls the RBAC version,the lite RABC has users, roles and actions,
// the full version has extra permissions, default value is true.
func LiteRBAC() bool {
	return cfg.DefaultBool("lite_rbac", true)
}

// SplitPermission disable/enable admin_system,admin_audit,admin_security roles.
func SplitPermission() bool {
	return cfg.DefaultBool("SplitPermission", false)
}

// WebAppPath returns the frontend static files root path, you can specify an
// absolute path to replace the default relative path 'webapp'.
func WebAppPath() string {
	return AbsPath(cfg.DefaultString("WebAppPath", "webapp"))
}

// WebAppStaticPath rennrns the frontend static file directory.
func WebAppStaticPath() string {
	return path.Join(WebAppPath(), "static")
}

// CommitInfo reads and returns commit info from the version file.
func CommitInfo() Commit {
	var fCommitID, bCommitID string
	content, err := ioutil.ReadFile(path.Join(WebAppPath(), "version"))
	if err != nil || len(string(content)) < 8 {
		fCommitID = "unknown"
	} else {
		fCommitID = string(content)[0:7]
	}
	content, err = ioutil.ReadFile(path.Join(HPCRootPath(), "version"))
	if err != nil || len(string(content)) < 8 {
		bCommitID = "unknown"
	} else {
		bCommitID = string(content)[0:7]
	}
	return Commit{
		Frontend: fCommitID,
		Backend:  bCommitID,
	}
}

// AppDir returns the scheduler application template files directory.
func AppDir() string {
	return cfg.DefaultString("APPDir", path.Join(AppDataPath(), "applications"))
}

// WorkDir returns the directory for job script under user home.
func WorkDir() string {
	return cfg.DefaultString("WorkDir", "applications")
}

// JobCollectInterval returns the period of job info collection.
func JobCollectInterval() int {
	return cfg.DefaultInt("JobCollectInterval", 5)
}

// DBDebug returns database debug status
func DBDebug() bool {
	return cfg.DefaultBool("mysql::DBDebug", false)
}

//// Notification (notify) ////

// NotifyType returns the notification provider.
func NotifyType() string {
	return cfg.DefaultString("notify_type", "redis")
}

// NotifyCounterKey returns the notify counter key of Redis.
func NotifyCounterKey() string {
	return cfg.DefaultString("notify_counter_key", "notify:counter")
}

//// Parallel command (pcm) ////

// PcmOutputDir returns the command redirect root directory.
func PcmOutputDir() string {
	return cfg.DefaultString("PCMOutputDir", path.Join(DefaultHomeDir(), "pcm"))
}

// PcmMaxOutput returns the max size of output.
func PcmMaxOutput() int {
	return cfg.DefaultInt("PcmMaxOutput", 1024)
}

// PcmCounterKey returns the key of command counter in Redis.
func PcmCounterKey() string {
	return cfg.DefaultString("PcmCounterKey", "pcm:counter")
}

// PcmChanBufferSize returns the notify chan buffer size.
func PcmChanBufferSize() int {
	return cfg.DefaultInt("PcmChanBufferSize", 500)
}

// MailChanBufferSize returns the mail chan buffer size.
func MailChanBufferSize() int {
	return cfg.DefaultInt("MailChanBufferSize", 1024)
}

// IPMI returns the IPMI configuration as a JSON string.
func IPMI() string {
	return cfg.DefaultString("IPMI", `ipmitool -I lanplus -U root -P root`)
}

// MonitorInterval returns the interval of monitor data collection.
func MonitorInterval() int {
	return cfg.DefaultInt("backend::MonitorInterval", 30)
}

// NodeRoles returns the roles info of node.
func NodeRoles() string {
	return cfg.DefaultString("NodeRoles",
		`[{
			"Name": "M",
			"Meaning": {
				"HPC": {
					"zh-CN": "管理",
					"en-US": "Management"
				},
				"Backup": {
					"zh-CN": "管理",
					"en-US": "Management"
				}
			},
			"LoginAllow": false
		},
		{
			"Name": "C",
			"Meaning": {
				"HPC": {
					"zh-CN": "计算",
					"en-US": "Cacl"
				},
				"Backup": {
					"zh-CN": "客户端",
					"en-US": "Client"
				}
			},
			"LoginAllow": false
		},
		{
			"Name": "I",
			"Meaning": {
				"HPC": {
					"zh-CN": "存储",
					"en-US": "I/O"
				},
				"Backup": {
					"zh-CN": "存储",
					"en-US": "Storager"
				}
			},
			"LoginAllow": false
		},
		{
			"Name": "L",
			"Meaning": {
				"HPC": {
					"zh-CN": "登陆",
					"en-US": "Login"
				},
				"Backup": {
					"zh-CN": "登陆",
					"en-US": "Login"
				}
			},
			"LoginAllow": true
		}
	]`)
}

// RRDDir returns ganglia rrd files directory.
func RRDDir() string {
	return cfg.DefaultString("backend::RRDDir", "/var/lib/ganglia/rrds")
}

// NetDevicesRRDDir returns the net devices rrd files directory.
func NetDevicesRRDDir() string {
	return cfg.DefaultString("NetDevicesRRDDir",
		path.Join(RRDDir(), ClusterName(), "NetDevices"))
}

// EnableLogger returns logger status
func EnableLogger() bool {
	return cfg.DefaultBool("backend::EnableLogger", false)
}

// NetInterface returns the net adapter name.
func NetInterface() string {
	return cfg.DefaultString("backend::InterfaceHost", "eth0")
}

// InterfaceAgent returns the NIC name which connect with agent.
func InterfaceAgent() string {
	return cfg.DefaultString("backend::InterfaceAgent", "eth1")
}

// HPCGroupID returns the hpc program user's gid.
func HPCGroupID() int {
	return cfg.DefaultInt("backend::HPCGroupID", 0)
}

// JobScale returns the default report job scale.
func JobScale() string {
	return cfg.DefaultString("JobScale",
		`{"Tiny":50,"Small":250,"Middle":500,"Big":1000}`)
}

// DiskNames returns the disk names split by comma, be used for disk I/O monitor.
func DiskNames() string {
	return cfg.DefaultString("DiskNames", "sda")
}

// CORSOrigin returns the Cross-origin resource sharing allowed origin. No need
// to set this configuration usually, just keep it "*", because we can control
// access permission in our token.
func CORSOrigin() string {
	return cfg.DefaultString("CORSOrigin", "*")
}

// TTYFilesPath returns the static files of tty.
func TTYFilesPath() string {
	return cfg.DefaultString("TTYFilesPath", path.Join(HPCRootPath(), "tty"))
}

// DefaultHomeDir returns the default home directory.
func DefaultHomeDir() string {
	return cfg.DefaultString("DefaultHomeDir", "/home")
}

// HomeDirPerm return the default home directory permission.
func HomeDirPerm() os.FileMode {
	mask := cfg.DefaultInt("HomeDirPerm", 0700)
	return os.FileMode(mask)
}

// UserMinUID returns the min uid number allowed in current system.
func UserMinUID() int {
	return cfg.DefaultInt("uid::UserMinUID", 10000)
}

// UserMaxUID returns the max uid number allowed in current system.
func UserMaxUID() int {
	return cfg.DefaultInt("uid::UserMaxUID", 60000)
}

// -------------    LDAP related    --------------------

// LDAPHost returns the LDAP server IP address.
func LDAPHost() string {
	return cfg.DefaultString("ldap::LDAPHost", "192.168.1.29")
}

// LDAPPort returns the LDAP service listened port.
func LDAPPort() int {
	return cfg.DefaultInt("ldap::LDAPPort", 389)
}

// LDAPEnableSSL returns ture if enable SSL on LDAP, default value is false.
func LDAPEnableSSL() bool {
	return cfg.DefaultBool("ldap::LDAPEnableSSL", false)
}

// LDAPUserRDN returns the user RND of the LDAP.
func LDAPUserRDN() string {
	return cfg.DefaultString("ldap::LDAPUserRDN", "ou=People,dc=hpc,dc=com")
}

// LDAPServerName returns the server name of LDAP server, it is unnecessary if
// LDAPEnableSSL is false.
func LDAPServerName() string {
	return cfg.DefaultString("ldap::LDAPServerName", "ldap.hpc.com")
}

// LDAPAdminDN returns the LDAP admin distinguish name.
func LDAPAdminDN() string {
	return cfg.DefaultString("ldap::LDAPAdminDN", "cn=root,dc=hpc,dc=com")
}

// LDAPAdminPassword returns the LDAP admin's password.
func LDAPAdminPassword() string {
	return cfg.DefaultString("ldap::LDAPAdminPassword", "8ik,*IK<")
}

// LDAPAttrHomeDir returns the home dir attribute field name of LDAP
func LDAPAttrHomeDir() string {
	return cfg.DefaultString("ldap::LDAPAttrHomeDir", "homeDirectory")
}

// LDAPAttrUID returns the uid attribute field name of LDAP.
func LDAPAttrUID() string {
	return cfg.DefaultString("ldap::LDAPAttrUID", "uidNumber")
}

// LDAPAttrGID returns the gid attribute field name of LDAP.
func LDAPAttrGID() string {
	return cfg.DefaultString("ldap::LDAPAttrGID", "gidNumber")
}

// LDAPAttrLoginShell returns the login shell attribute field name of LDAP.
func LDAPAttrLoginShell() string {
	return cfg.DefaultString("ldap::LDAPAttrLoginShell", "loginShell")
}

// LDAPAttrPassword returns the password attribute field name of LDAP.
func LDAPAttrPassword() string {
	return cfg.DefaultString("ldap::LDAPAttrPassword", "userPassword")
}

// LDAPUserObjectClasses returns the user object classes of LDAP.
func LDAPUserObjectClasses() []string {
	return cfg.DefaultStrings("ldap::LDAPUserObjectClasses", []string{"inetOrgPerson", "organizationalPerson", "posixAccount"})
}

// LDAPGroupObjectClasses returns the group object classes of LDAP.
func LDAPGroupObjectClasses() []string {
	return cfg.DefaultStrings("ldap::LDAPGroupObjectClasses", []string{"posixGroup"})
}

// LDAPGroupRDN returns the group RDN of LDAP.
func LDAPGroupRDN() string {
	return cfg.DefaultString("ldap::LDAPGroupRDN", "ou=Group,dc=hpc,dc=com")
}

// LDAPAttrGroupName returns the group name attribute field name of LDAP.
func LDAPAttrGroupName() string {
	return cfg.DefaultString("ldap::LDAPAttrGroupName", "cn")
}

// LDAPAttrGroupMember returns the group member attribute field name of LDAP.
func LDAPAttrGroupMember() string {
	return cfg.DefaultString("ldap::LDAPAttrGroupMember", "memberUid")
}

// LDAPGroupFilter returns the group filter of LDAP.
func LDAPGroupFilter() string {
	return cfg.DefaultString("ldap::LDAPGroupFilter", "(objectClass=posixGroup)")
}

// -------------    SNMP related    --------------------

// SNMPCollectInterval returns the SNMP information collect interval(unit:s).
func SNMPCollectInterval() int {
	return cfg.DefaultInt("snmp::Interval", 15)
}

// IsInitialized will return status of application.
func IsInitialized() bool {
	return cfg.DefaultBool("backend::IsInitialized", false)
}

// SetConfigValue will set and save config item at app.conf file.
func SetConfigValue(data map[string]string) error {
	for k, v := range data {
		if err := cfg.Set(k, v); err != nil {
			return err
		}
	}

	if err := cfg.SaveConfigFile(path.Join(HPCRootPath(),
		"conf/app.conf")); err != nil {
		return err
	}
	return nil
}

// GetConfigValue returns config value by key.
func GetConfigValue(key string) string {
	return cfg.DefaultString(key, "")
}

// CasbinModelFile returns the Casbin's model file path.
func CasbinModelFile() string {
	return AbsPath(cfg.DefaultString("CasbinModelFile", "conf/model.conf"))
}

// CasbinPolicyFile returns the Casbin's policy file path.
func CasbinPolicyFile() string {
	return AbsPath(cfg.DefaultString("CasbinPolicyFile", "conf/policy.csv"))
}

// GetBackupConfig returns backup config from config file.
func GetBackupConfig() BackupConfig {
	return BackupConfig{
		Host:     cfg.DefaultString("backup::Host", "192.168.1.41"),
		Port:     cfg.DefaultInt("backup::Port", 9095),
		Username: cfg.DefaultString("backup::Username", "admin"),
		Password: cfg.DefaultString("backup::Password", "admin"),
	}
}

// BackupDBConfig returns backup database config.
func BackupDBConfig() MySQLConfig {
	return MySQLConfig{
		Host:     cfg.DefaultString("backup::DBHost", "192.168.1.41"),
		Port:     cfg.DefaultInt("backup::DBPort", 3306),
		User:     cfg.DefaultString("backup::DBUser", "root"),
		Password: cfg.DefaultString("backup::DBPasswd", ""),
		Db:       cfg.DefaultString("backup::DBName", "bacula"),
		Debug:    cfg.DefaultBool("backup::DBDebug", false),
	}
}

// HPCMySQLConfig returns HPC database config.
func HPCMySQLConfig() MySQLConfig {
	return MySQLConfig{
		Host:     cfg.DefaultString("mysql::DBHost", "localhost"),
		Port:     cfg.DefaultInt("mysql::DBPort", 3306),
		User:     cfg.DefaultString("mysql::DBUser", "root"),
		Password: cfg.DefaultString("mysql::DBPasswd", ""),
		Db:       cfg.DefaultString("mysql::DBName", "hpc_backend"),
		Debug:    cfg.DefaultBool("mysql::DBDebug", false),
		Schema:   cfg.DefaultString("mysql::Schema", "standard;cloud;virtual-desktop;training"),
	}
}

// ConfigVersion replaces the release version if it IS NOT empty.
func ConfigVersion() string {
	return cfg.String("version")
}

// AgentServerPort returns the agent RPC listen port.
func AgentServerPort() int {
	return cfg.DefaultInt("AgentServerPort", 6380)
}

// UDPListenPort returns the UDP listen port.
func UDPListenPort() int {
	return cfg.DefaultInt("UDPListenPort", 6382)
}

// VirtualDesktopServerPort returns the virtual desktop agent RPC listen port.
func VirtualDesktopServerPort() int {
	return cfg.DefaultInt("VirtualDesktopServerPort", 6381)
}

// AgentFileLogger returns the agent logger configuration.
func AgentFileLogger() (string, error) {
	config := cfg.DefaultString("agent::FileLogger",
		`{"filename":"/var/log/hpc-backend/hpc-agent.log","daily":true,"maxdays":10}`)
	return checkFileLoggerConfig(config)
}

// ConsulBasePath returns the agent key prefix as HPC agent's namespace in
// consul.
func ConsulBasePath() string {
	return cfg.DefaultString("ConsulBasePath", "hpc-agent")
}

// ConsulVDAgentBasePath returns the vdagent key prefix as vdagent's namespace in
// consul.
func ConsulVDAgentBasePath() string {
	return cfg.DefaultString("ConsulVDAgentBasePath", "hpc-vdagent")
}

// ConsulServicePath returns the service path(namespace) in consul.
func ConsulServicePath() string {
	return cfg.DefaultString("ServicePath", "Nodes")
}

// ConsulAddrs returns the consul servers' addrs.
func ConsulAddrs() []string {
	return cfg.Strings("ConsulAddrs")
}

// RemoteServiceAddr returns the remote connect service HTTP addr.
func RemoteServiceAddr() string {
	return cfg.DefaultString("RemoteServiceAddr", "http://192.168.2.21:8080")
}

// RemoteServiceUsername returns the remote connect service username.
func RemoteServiceUsername() string {
	return cfg.DefaultString("RemoteServiceUsername", "guacadmin")
}

// RemoteServicePassword returns the remote connect service password.
func RemoteServicePassword() string {
	return cfg.DefaultString("RemoteServicePassword", "guacadmin")
}

// BackupConfigFileByType returns backup system dir config file path.
func BackupConfigFileByType(typ string) string {
	if typ == "dir" {
		return cfg.DefaultString("backup::DirFile", "/usr/local/bacula/etc/bacula-dir.conf")
	} else if typ == "fd" {
		return cfg.DefaultString("backup::FDFile", "/usr/local/bacula/etc/bacula-fd.conf")
	} else if typ == "sd" {
		return cfg.DefaultString("backup::SDFile", "/usr/local/bacula/etc/bacula-sd.conf")
	} else if typ == "bcons" {
		return cfg.DefaultString("backup::ConsoleFile", "/usr/local/bacula/etc/bconsole.conf")
	} else {
		return ""
	}
}

// JWTSecret returns the JSON Web Token signer secret key.
func JWTSecret() string {
	return cfg.DefaultString("JWTSecret",
		"18ddf30d665538d3ab90b8e0bf6c96879be4fa6d")
}

// RBACDebugOn represents the RBAC debug switch.
func RBACDebugOn() bool {
	return cfg.DefaultBool("RBACDebugOn", false)
}

// GetCapabilities returns all agent capabilities.
func GetCapabilities() string {
	return cfg.String("backend::Capabilities")
}

// SysAdminAsRoot determines software system admin's system privillege, if it's
// true, all admin will have root's privillege. (default false)
func SysAdminAsRoot() bool {
	return cfg.DefaultBool("SysAdminAsRoot", false)
}

// DockerImages returns docker images config file path
func DockerImages() string {
	return cfg.DefaultString("backend::DockerImages", "/var/lib/socker/images.yaml")
}

// SingularityImages returns singularity images config file path
func SingularityImages() string {
	return AbsPath(cfg.DefaultString("backend::SingularityImages",
		path.Join(HPCRootPath(), "conf/singularity.yml")))
}

// CheckListFilePath returns the checklist file path
func CheckListFilePath() string {
	return AbsPath(cfg.DefaultString("backend::ChecklistFilePath", "conf/checklist.json"))
}

//AnsibleHostFile returns the ansible host file path
func AnsibleHostFile() string {
	return cfg.DefaultString("backend::AnsibleHostFile", "/etc/ansible/grouphosts")
}

//HarborServer returns the Harbor server URL
func HarborServer() string {
	return cfg.String("backend::HarborServer")
}

//HarborUserName returns the Harbor server userName of admin
func HarborUserName() string {
	return cfg.DefaultString("backend::HarborUserName", "admin")
}

//HarborPassword returns the harbor server password of admin
func HarborPassword() string {
	return cfg.DefaultString("backend::HarborPassword", "Harbor12345")
}

// AbsPath returns the path which joined with HPC_ROOT path if the input path is
// a relative path.
func AbsPath(p string) string {
	if path.IsAbs(p) {
		return p
	}
	return path.Join(HPCRootPath(), p)
}

// AlarmConfigFilePath returns the alarm config file path
func AlarmConfigFilePath() string {
	return AbsPath(cfg.DefaultString("backend::AlarmConfigFilePath",
		path.Join(HPCRootPath(), "conf/alarm.yml")))
}

// CustomImagePath returns the path of user custom image file
func CustomImagePath() string {
	return AbsPath(cfg.DefaultString("backend::CustomImagePath", "/usr/hpc/upload/images/public"))
}

// ImageMaxSpace returns the maximum space of each user can use
func ImageMaxSpace() int {
	return cfg.DefaultInt("backend::ImageMaxSpace", 500)
}

// HarborCookieExpires returns the harbor server expires time
func HarborCookieExpires() int64 {
	return cfg.DefaultInt64("backend::HarborCookieExpires", int64(1800))
}

// VDUnregisteredPath returns the virtual desktop unregistered image path
func VDUnregisteredPath() string {
	return AbsPath(cfg.DefaultString("virtual-desktop::VDUnregisteredPath", "/usr/hpc/vd/unregistered"))
}

// VDRegisteredPath returns the virtual desktop registered image path
func VDRegisteredPath() string {
	return AbsPath(cfg.DefaultString("virtual-desktop::VDRegisteredPath", "/usr/hpc/vd/registered"))
}

// AutoSeedEnabled determines whether to auto seeding.
func AutoSeedEnabled() bool {
	return cfg.DefaultBool("virtual-desktop::AutoSeedEnabled", true)
}

// TrackerURLs returns the BitTorrent tracker URLs.
func TrackerURLs() []string {
	return cfg.DefaultStrings("virtual-desktop::TrackerURLs",
		[]string{"http://192.168.1.163:6881/announce",
			"udp://192.168.1.163:6881"})
}

// VDStartupScriptPath returns the path of vdi virtual machine startup script
func VDStartupScriptPath() string {
	return AbsPath(cfg.DefaultString("virtual-desktop::VDStartupScriptPath", "/home/vd/start.sh"))
}

// UseStrongPassword determines whether to validate password security.
func UseStrongPassword() bool {
	return cfg.DefaultBool("UseStrongPassword", true)
}

// KVMImage return the kvm image path used to create vdi virtual machine
func KVMImage() string {
	return cfg.DefaultString("virtual-desktop::KVMImage", "harbor.hpc.com/images/hpc/kvm")
}

// VDICDROMFile return the CD ROM path used to create vdi virtual machine
func VDICDROMFile() string {
	return cfg.DefaultString("virtual-desktop::VDICDROMFile", "/opt/docker-winxp/mock.iso")
}

// ManagerNet returns the management network name.
func ManagerNet() string {
	return cfg.DefaultString("virtual-desktop::ManagerNet", "mgmt")
}

// DockerInsecureMode determines whether to run socker in insecure mode.
func DockerInsecureMode() bool {
	return cfg.DefaultBool("backend::DockerInsecureMode", false)
}

// DefaultShell returns the default shell program.
func DefaultShell() string {
	return cfg.DefaultString("backend::DefaultShell", "sh")
}

// DockerNetwork determines which network connect to run Docker container.
func DockerNetwork() string {
	return cfg.DefaultString("backend::DockerNetwork", "")
}

// IgnoreNetworks returns the ignored docker networks, multiple networks split
// with comma.
func IgnoreNetworks() string {
	return cfg.DefaultString("backend::IgnoreNetworks", "host")
}

// DockerEnabled determines whether to eanble Docker.
func DockerEnabled() bool {
	return cfg.DefaultBool("backend::DockerEnabled", true)
}

// SingularityEnabled determines whether to eanble singularity.
func SingularityEnabled() bool {
	return cfg.DefaultBool("backend::SingularityEnabled", true)
}

// EpilogLockFilePath represents the container epilog lock file path.
func EpilogLockFilePath() string {
	return cfg.DefaultString("backend::EpilogLockFilePath",
		"/var/lib/socker/epilog")
}

// EpilogScript returns the Slurm epilog script file path.
func EpilogScript() string {
	return cfg.DefaultString("backend::EpilogScript",
		"/usr/hpc/scripts/epilog.sh")
}

// HarborProject return the harbor project name
func HarborProject() string {
	return cfg.DefaultString("backend::HarborProject", "hpc")
}

// HarborDefaultPageSize return the harbor request default page number
func HarborDefaultPageSize() string {
	return cfg.DefaultString("backend::HarborDefaultPageSize", "500")
}

// ImageMaxNumber return the maximum number of image each user can upload
func ImageMaxNumber() int {
	return cfg.DefaultInt("backend::ImageMaxNumber", 10)
}

// HarborServerName return the harbor server url
func DockerRegistry() string {
	return cfg.DefaultString("backend::DockerRegistry", "harbor.hpc.com")
}

// HarborPublicRepo return the harbor public repository name
func HarborPublicRepo() string {
	return cfg.DefaultString("backend::HarborPublicRepo", "public")
}

// GetUploadLimit return the file size limit of uploading file
func GetUploadLimit() int64 {
	return cfg.DefaultInt64("backend::UploadLimit", 10240)
}

// TmuxWindowDisabled determines whether to disable tmux new-window function.
func TmuxWindowDisabled() bool {
	return cfg.DefaultBool("TmuxWindowDisabled", true)
}

// WolBroadcastAddr get default wol BroadcastAddr
func WolBroadcastAddr() string {
	return cfg.DefaultString("virtual-desktop::WolBroadcastAddr", "255.255.255.255:9")
}

// ISOPath get the iso file path
func ISOPath() string {
	return cfg.DefaultString("virtual-desktop::ISOPath", "/home/hpc/vd/unregistered/iso")
}

// SharedStoragePath the path of shared storage
func SharedStoragePath() string {
	return cfg.DefaultString("virtual-desktop::SharedStoragePath", "/home")
}

// ClientMountPath the path of vd client mount
func ClientMountPath() string {
	return cfg.DefaultString("virtual-desktop::SharedStoragePath", "/home")
}

// NetworkNodes returns the network nodes hostnames.
func NetworkNodes() []string {
	return cfg.DefaultStrings("NetworkNodes", []string{"network"})
}

// ProxyUser returns the network node proxy system user name.
func ProxyUser() string {
	return cfg.DefaultString("ProxyUser", "root")
}

// VncPort returns the port of vnc connect.
func VncPort() int {
	return cfg.DefaultInt("virtual-desktop::VncPort", 5900)
}

// SpicePort returns the port of spice connect.
func SpicePort() int {
	return cfg.DefaultInt("virtual-desktop::SpicePort", 5901)
}

// TelnetPort returns the port of telnet connect.
func TelnetPort() int {
	return cfg.DefaultInt("virtual-desktop::TelnetPort", 59000)
}

// TorrentBlockSize determines the torrent block size, default 4K.
func TorrentBlockSize() string {
	return cfg.DefaultString("virtual-desktop::TorrentBlockSize", "2048")
}

// MulticastIP return the ip of multicast.
func MulticastIP() string {
	return cfg.DefaultString("backend::MulticastIP", "224.0.0.200")
}

// AllowCustomizedInstanceType determines whether to allow customized instance
// conigurations.
func AllowCustomizedInstanceType() bool {
	return cfg.DefaultBool("cloud::AllowCustomizedInstanceType", true)
}

// VMReviewRequired determines the vm req whether to need review
func VMReviewRequired() bool {
	return cfg.DefaultBool("cloud::VMReviewRequired", true)
}
