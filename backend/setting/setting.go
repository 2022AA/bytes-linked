package setting

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

// 存储类型(表示文件存到哪里)
type StoreType int

const (
	_ StoreType = iota
	// StoreLocal : 节点本地
	StoreLocal
	// StoreCeph : Ceph集群
	StoreCeph
	// StoreOSS : 阿里OSS
	StoreOSS
	// StoreMinio : Minio
	StoreMinio
)

type App struct {
	JwtSecret string
	PageSize  int
	PrefixUrl string

	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string

	TempLocalRootDir string
	CephRootDir      string
	CurrentStoreType StoreType
	UploadLBHost     string
	DownloadBHost    string
}

var AppSetting = &App{}

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

type Database struct {
	Type     string
	User     string
	Password string
	Host     string
	Name     string
}

var DatabaseSetting = &Database{}

type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

var RedisSetting = &Redis{}

type Ceph struct {
	CephAccessKey  string
	CephSecretKey  string
	CephGWEndpoint string
}

var CephSetting = &Ceph{}

type OSS struct {
	OSSBucket          string
	OSSEndpoint        string
	OSSAccesskeyID     string
	OSSAccessKeySecret string
}

var OSSSetting = &OSS{}

type Minio struct {
	MinioEndpoint        string
	MinioAccesskeyID     string
	MinioAccessKeySecret string
	MinioUseSSL          bool
}

var MinioSetting = &Minio{}

type IPfs struct {
	IPfsEndpoint string
}

var IPfsSetting = &IPfs{}

var cfg *ini.File

// Setup initialize the configuration instance
func Setup() {
		var err error
		cfg, err = ini.Load("conf/app.ini")
		if err != nil {
			log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
		}

		mapTo("app", AppSetting)
		mapTo("server", ServerSetting)
		mapTo("database", DatabaseSetting)
		mapTo("redis", RedisSetting)
		mapTo("ceph", CephSetting)
		mapTo("oss", OSSSetting)
		mapTo("minio", MinioSetting)
		mapTo("ipfs", IPfsSetting)

		ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
		ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
		RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
