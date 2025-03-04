package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/jtprogru/yc-quotas-exporter/internal/config"
	"github.com/jtprogru/yc-quotas-exporter/internal/yandexcloud"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to config file")
	flag.Parse()

	if configPath == nil {
		log.Fatal("config is required")
	}

	fmt.Println("config inicialization is done")
	cfg, err := config.New(*configPath)
	if err != nil {
		fmt.Println(err)
	}

	client, err := yandexcloud.New(cfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("client inicialization is done")

	client.QuotaLimitListService()
}
