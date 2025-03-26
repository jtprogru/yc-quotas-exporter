package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jtprogru/yc-quotas-exporter/internal/config"
	"github.com/jtprogru/yc-quotas-exporter/internal/yandexcloud"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("config inicialization is done")

	pe, err := yandexcloud.NewExporter(cfg)
	if err != nil {
		log.Fatalf("Error creating Exporter: %v", err)
	}
	log.Fatal(pe.Run())
}
