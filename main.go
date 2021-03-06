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
var conf Config

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
	err = json.Unmarshal(content, &conf)
	if err != nil {
		Log(LOG_ERROR, tag, "Unable to parse the config file")
		os.Exit(1)
	}
	ShouldLog = conf.LogLevel == "debug"
	//endregion
	//region Init cache
	CacheInit()
	//endregion
	//region Start In Memory DB Update
	go UpdateUserDB()
	//endregion
	//region Start Proxy
	for _, s := range conf.Servers {
		proxy := NewProxy(s.LocalAddress, s.RemoteAddress)
		go proxy.Start()
		ProxyList = append(ProxyList, proxy)
	}
	//endregion
	//region Start AMPQ Listener if available
	if conf.Ampq.Enable {
		go StartAmpq(conf.Ampq.ConnectionString, conf.Ampq.Exchange, conf.Node)
	}
	//endregion
	//region Stats
	if conf.Stats.Enable {
		go StartStatsServer(conf.Stats.ListenAddress)
	}
	//endregion
	//region Wait for interrupt
	<-interrupt
	for _, proxy := range ProxyList {
		proxy.Stop()
	}
	if conf.Ampq.Enable {
		StopAmpq()
	}
	Log(LOG_INFO, tag, "Interrupt received, shutting down...")
	//endregion
}
