package conf

import (
	"github.com/Tropicana33/common/config"
	"github.com/Tropicana33/common/golog"
)

func InitConfig(file string) {
	conf, err := config.LoadConfigByFile(file)
	if err != nil {
		golog.Critical(err)
	}

	var section string = "default"
	var str string
	var num int

	// image2webp section
	section = "go-curl"
	num = conf.MustInt(section, "workerCount", 1)
	g_config.workerCount = num
	golog.Debug(section, "workerCount:", num)

	num = conf.MustInt(section, "reqLimit", 1)
	g_config.reqLimit = num
	golog.Debug(section, "reqLimit:", num)

	num = conf.MustInt(section, "startIndex", 0)
	g_config.startIndex = num
	golog.Debug(section, "startIndex:", num)

	str = conf.MustString(section, "upsServer", "127.0.0.1:8080")
	g_config.upsServer = str
	golog.Debug(section, "upsServer:", str)

	str = conf.MustString(section, "baseUrl", "http://www.a.com/index/")
	g_config.baseUrl = str
	golog.Debug(section, "baseUrl:", str)

	num = conf.MustInt(section, "pprofPort", 19333)
	g_config.pprofPort = num
	golog.Debug(section, "pprofPort:", num)
}

var g_config *GoCurlConf = new(GoCurlConf)

type GoCurlConf struct {
	// jpeg2webp
	workerCount int
	pprofPort   int
	reqLimit    int
	upsServer   string
	baseUrl     string
	startIndex  int
}

func GetConfig() *GoCurlConf {
	return g_config
}

func (conf *GoCurlConf) GetWorkerCount() int {
	if conf.workerCount <= 0 {
		return 1
	}

	return conf.workerCount
}

func (conf *GoCurlConf) GetReqLimit() int {
	if conf.reqLimit <= 0 {
		return 1
	}

	return conf.reqLimit
}

func (conf *GoCurlConf) GetStartIndex() int {
	if conf.startIndex < 0 {
		return 0
	}

	return conf.startIndex
}

func (conf *GoCurlConf) GetUpsServer() string {
	return conf.upsServer
}

func (conf *GoCurlConf) GetBaseUrl() string {
	return conf.baseUrl
}

func (conf *GoCurlConf) GetPprofPort() int {
	if conf.pprofPort == 0 {
		return 19333
	}

	return conf.pprofPort
}
