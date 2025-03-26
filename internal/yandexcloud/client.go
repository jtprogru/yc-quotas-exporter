package yandexcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	quotamanager "github.com/yandex-cloud/go-genproto/yandex/cloud/quotamanager/v1"
	ycsdk "github.com/yandex-cloud/go-sdk"

	"github.com/jtprogru/yc-quotas-exporter/internal/config"
)

const (
	defaultResourceType = "resource-manager.cloud"
)

type Client struct {
	Config *config.Config
	SDK    *ycsdk.SDK
}

func NewClient(cfg *config.Config) (*Client, error) {
	ctx := context.Background()
	sdk, err := ycsdk.Build(ctx, ycsdk.Config{
		Credentials: ycsdk.NewIAMTokenCredentials(cfg.Token),
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &Client{
		Config: cfg,
		SDK:    sdk,
	}, nil
}

func (c *Client) QuotaLimitListService() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Config.Timeout)*time.Second)
	defer cancel()

	res, err := c.SDK.QuotaManager().QuotaLimit().ListServices(ctx, &quotamanager.ListServicesRequest{
		ResourceType: defaultResourceType,
	})
	if err != nil {
		fmt.Println("res, err := c.SDK.QuotaManager().QuotaLimit().ListServices(ctx, nil)")
		log.Println(err)
		return
	}
	for _, service := range res.Services {
		fmt.Println(service)
		c.QuotaLimitList(service.Id)
	}
}

func (c *Client) QuotaLimitServiceGet(serviceID string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Config.Timeout)*time.Second)
	defer cancel()

	res, err := c.SDK.QuotaManager().QuotaLimit().Get(ctx, &quotamanager.GetQuotaLimitRequest{
		Resource: &quotamanager.Resource{
			Id:   c.Config.CloudID,
			Type: defaultResourceType,
		},
		QuotaId: serviceID,
	})
	if err != nil {
		fmt.Println("c.QuotaLimitServiceGet(...)")
		log.Println(err)
		return
	}
	fmt.Println(res)
}

func (c *Client) QuotaLimitList(serviceID string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Config.Timeout)*time.Second)
	defer cancel()

	res, err := c.SDK.QuotaManager().QuotaLimit().List(ctx, &quotamanager.ListQuotaLimitsRequest{
		Resource: &quotamanager.Resource{
			Id:   c.Config.CloudID,
			Type: defaultResourceType,
		},
		Service: serviceID,
	})
	if err != nil {
		fmt.Println("c.QuotaLimitList()")
		log.Println(err)
		return
	}
	fmt.Println(res)
}
