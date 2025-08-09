package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/fengziwei0826/caching_proxy_exe/internal/conf"
)

var (
	Port       int    = 3000
	ProxyAddr  string = "http://dummyjson.com"
	ConfigPath string = ""
	ClearCache bool   = false
)

func init() {
	flag.IntVar(&Port, "port", 3000, "HTTP server port")
	flag.StringVar(&ProxyAddr, "origin", "http://dummyjson.com", "Proxy origin address")
	flag.StringVar(&ConfigPath, "config", "", "Path to configuration file")
	flag.BoolVar(&ClearCache, "clear-cache", false, "Clear cache, then return")
	flag.Parse()
}

func main() {
	gConfig := conf.GetGlobalConfig()
	gConfig.HTTPConfig.Port = Port
	gConfig.HTTPConfig.ProxyAddr = ProxyAddr
	if ConfigPath != "" {
		confErr := conf.LoadConfigFromFile(ConfigPath)
		if confErr != nil {
			log.Fatalf("Failed to load configuration from file: %v", confErr)
		}
		log.Printf("Using configuration file: %s\n", ConfigPath)
	}
	log.Printf("Global configuration: %v\n", conf.GetGlobalConfig())
	if ClearCache {
		log.Printf("Clearing cache and exiting...")
		mgr := InitCacheManager()
		defer mgr.Close()
		mgr.ClearAllCache()
		return
	}
	srv := InitCacheProxyServer()
	err := srv.Start()
	if err != nil {
		log.Printf("Failed to start server: %v", err)
		return
	}
	defer srv.Stop()
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt)
	<-exitChan
}
