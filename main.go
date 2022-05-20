package main

import (
	"encoding/json"
	"fmt"
	"github.com/DavidGamba/go-getoptions"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
)

var ProxyList []*Proxy
var ShouldLog bool = false

func main() {
	tag := "[MAIN]"
	var err error
	var configFile string
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	//region Command line options
	opt := getoptions.New()
	opt.Bool("help", false, opt.Alias("h", "?"))
	opt.StringVar(&configFile, "config", "config.json", opt.Alias("c"), opt.Description("Config File"))
	_, err = opt.Parse(os.Args[1:])
	if err != nil {
		panic(err)
	}
	if opt.Called("help") {
		fmt.Fprintf(os.Stdout, opt.Help())
		os.Exit(1)
	}
	//endregion
	//region Read config file
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		Log(LOG_ERROR, tag, "Unable to read the config file")
		os.Exit(1)
	}
	var conf Config
	err = json.Unmarshal(content, &conf)
	if err != nil {
		Log(LOG_ERROR, tag, "Unable to parse the config file")
		os.Exit(1)
	}
	ShouldLog = conf.LogLevel == "debug"
	//endregion
	CacheInit()
	UpdateUserDB()
	//region Start Proxy
	for _, s := range conf.Servers {
		proxy := NewProxy(s.LocalAddress, s.RemoteAddress)
		go proxy.Start()
		ProxyList = append(ProxyList, proxy)
	}
	//endregion
	//region Start AMPQ Listener if available
	if conf.Ampq.Enable {
		StartAmpq(conf.Ampq.ConnectionString, conf.Ampq.Exchange, conf.NodeName)
	}
	//endregion

	//region Wait for interrupt
	<-interrupt
	for _, proxy := range ProxyList {
		proxy.Stop()
	}
	StopAmpq()
	Log(LOG_INFO, LOG_INFO, "Interrupt received, shutting down...")
	//endregion
}
