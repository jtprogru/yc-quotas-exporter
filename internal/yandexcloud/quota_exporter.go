package yandexcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"

	"github.com/jtprogru/yc-quotas-exporter/internal/config"
	quotamanager "github.com/yandex-cloud/go-genproto/yandex/cloud/quotamanager/v1"
)

// QuotaLimitServiceClient is a quotamanager.QuotaLimitServiceClient with
// lazy GRPC connection initialization.
type QuotaLimitServiceClient struct {
	getConn func(ctx context.Context) (*grpc.ClientConn, error)
}

// Get implements quotamanager.QuotaLimitServiceClient
func (c *QuotaLimitServiceClient) Get(ctx context.Context, in *quotamanager.GetQuotaLimitRequest, opts ...grpc.CallOption) (*quotamanager.QuotaLimit, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return quotamanager.NewQuotaLimitServiceClient(conn).Get(ctx, in, opts...)
}

// List implements quotamanager.QuotaLimitServiceClient
func (c *QuotaLimitServiceClient) List(ctx context.Context, in *quotamanager.ListQuotaLimitsRequest, opts ...grpc.CallOption) (*quotamanager.ListQuotaLimitsResponse, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return quotamanager.NewQuotaLimitServiceClient(conn).List(ctx, in, opts...)
}

type QuotaLimitIterator struct {
	ctx  context.Context
	opts []grpc.CallOption

	err           error
	started       bool
	requestedSize int64
	pageSize      int64

	client  *QuotaLimitServiceClient
	request *quotamanager.ListQuotaLimitsRequest

	items []*quotamanager.QuotaLimit
}

func (c *QuotaLimitServiceClient) QuotaLimitIterator(ctx context.Context, req *quotamanager.ListQuotaLimitsRequest, opts ...grpc.CallOption) *QuotaLimitIterator {
	var pageSize int64
	const defaultPageSize = 1000
	pageSize = req.PageSize
	if pageSize == 0 {
		pageSize = defaultPageSize
	}
	return &QuotaLimitIterator{
		ctx:      ctx,
		opts:     opts,
		client:   c,
		request:  req,
		pageSize: pageSize,
	}
}

func (it *QuotaLimitIterator) Next() bool {
	if it.err != nil {
		return false
	}
	if len(it.items) > 1 {
		it.items[0] = nil
		it.items = it.items[1:]
		return true
	}
	it.items = nil // consume last item, if any

	if it.started && it.request.PageToken == "" {
		return false
	}
	it.started = true

	if it.requestedSize == 0 || it.requestedSize > it.pageSize {
		it.request.PageSize = it.pageSize
	} else {
		it.request.PageSize = it.requestedSize
	}

	response, err := it.client.List(it.ctx, it.request, it.opts...)
	it.err = err
	if err != nil {
		return false
	}

	it.items = response.QuotaLimits
	it.request.PageToken = response.NextPageToken
	return len(it.items) > 0
}

func (it *QuotaLimitIterator) Take(size int64) ([]*quotamanager.QuotaLimit, error) {
	if it.err != nil {
		return nil, it.err
	}

	if size == 0 {
		size = 1 << 32 // something insanely large
	}
	it.requestedSize = size
	defer func() {
		// reset iterator for future calls.
		it.requestedSize = 0
	}()

	var result []*quotamanager.QuotaLimit

	for it.requestedSize > 0 && it.Next() {
		it.requestedSize--
		result = append(result, it.Value())
	}

	if it.err != nil {
		return nil, it.err
	}

	return result, nil
}

func (it *QuotaLimitIterator) TakeAll() ([]*quotamanager.QuotaLimit, error) {
	return it.Take(0)
}

func (it *QuotaLimitIterator) Value() *quotamanager.QuotaLimit {
	if len(it.items) == 0 {
		panic("calling Value on empty iterator")
	}
	return it.items[0]
}

func (it *QuotaLimitIterator) Error() error {
	return it.err
}

// ListServices implements quotamanager.QuotaLimitServiceClient
func (c *QuotaLimitServiceClient) ListServices(ctx context.Context, in *quotamanager.ListServicesRequest, opts ...grpc.CallOption) (*quotamanager.ListServicesResponse, error) {
	conn, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	return quotamanager.NewQuotaLimitServiceClient(conn).ListServices(ctx, in, opts...)
}

type QuotaLimitServicesIterator struct {
	ctx  context.Context
	opts []grpc.CallOption

	err           error
	started       bool
	requestedSize int64
	pageSize      int64

	client  *QuotaLimitServiceClient
	request *quotamanager.ListServicesRequest

	items []*quotamanager.Service
}

func (c *QuotaLimitServiceClient) QuotaLimitServicesIterator(ctx context.Context, req *quotamanager.ListServicesRequest, opts ...grpc.CallOption) *QuotaLimitServicesIterator {
	var pageSize int64
	const defaultPageSize = 1000
	pageSize = req.PageSize
	if pageSize == 0 {
		pageSize = defaultPageSize
	}
	return &QuotaLimitServicesIterator{
		ctx:      ctx,
		opts:     opts,
		client:   c,
		request:  req,
		pageSize: pageSize,
	}
}

func (it *QuotaLimitServicesIterator) Next() bool {
	if it.err != nil {
		return false
	}
	if len(it.items) > 1 {
		it.items[0] = nil
		it.items = it.items[1:]
		return true
	}
	it.items = nil // consume last item, if any

	if it.started && it.request.PageToken == "" {
		return false
	}
	it.started = true

	if it.requestedSize == 0 || it.requestedSize > it.pageSize {
		it.request.PageSize = it.pageSize
	} else {
		it.request.PageSize = it.requestedSize
	}

	response, err := it.client.ListServices(it.ctx, it.request, it.opts...)
	it.err = err
	if err != nil {
		return false
	}

	it.items = response.Services
	it.request.PageToken = response.NextPageToken
	return len(it.items) > 0
}

func (it *QuotaLimitServicesIterator) Take(size int64) ([]*quotamanager.Service, error) {
	if it.err != nil {
		return nil, it.err
	}

	if size == 0 {
		size = 1 << 32 // something insanely large
	}
	it.requestedSize = size
	defer func() {
		// reset iterator for future calls.
		it.requestedSize = 0
	}()

	var result []*quotamanager.Service

	for it.requestedSize > 0 && it.Next() {
		it.requestedSize--
		result = append(result, it.Value())
	}

	if it.err != nil {
		return nil, it.err
	}

	return result, nil
}

func (it *QuotaLimitServicesIterator) TakeAll() ([]*quotamanager.Service, error) {
	return it.Take(0)
}

func (it *QuotaLimitServicesIterator) Value() *quotamanager.Service {
	if len(it.items) == 0 {
		panic("calling Value on empty iterator")
	}
	return it.items[0]
}

func (it *QuotaLimitServicesIterator) Error() error {
	return it.err
}

type QuotaExporter struct {
	client *Client
}

func NewQuotaExporter(cfg *config.Config) (*QuotaExporter, error) {
	client, err := NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &QuotaExporter{client: client}, nil
}

func (qe *QuotaExporter) ExportMetrics() {
	servicesIter := qe.client.SDK.QuotaManager().QuotaLimit().QuotaLimitServicesIterator(context.Background(), &quotamanager.ListServicesRequest{
		ResourceType: defaultResourceType,
	})
	services, err := servicesIter.TakeAll()
	if err != nil {
		log.Fatal(err)
	}

	for _, service := range services {
		quotaLimits, err := qe.client.SDK.QuotaManager().QuotaLimit().List(context.Background(), &quotamanager.ListQuotaLimitsRequest{
			Resource: &quotamanager.Resource{
				Id:   qe.client.Config.CloudID,
				Type: defaultResourceType,
			},
			Service: service.Id,
		})
		if err != nil {
			log.Printf("Error getting quota limits for service %s: %v", service.Id, err)
			qe.exportQuotaMetric(service, nil)
			continue
		}

		for _, quotaLimit := range quotaLimits.QuotaLimits {
			qe.exportQuotaMetric(service, quotaLimit)
		}
	}
}

func (qe *QuotaExporter) exportQuotaMetric(service *quotamanager.Service, quota *quotamanager.QuotaLimit) {
	var limitValue, usageValue float64
	var quotaID string

	if quota == nil {
		limitValue = -1
		usageValue = -1
		quotaID = "unknown"
	} else {
		limitValue = float64(quota.Limit.GetValue())
		usageValue = float64(quota.Usage.GetValue())
		quotaID = quota.QuotaId
	}

	limitMetric := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "yandex_cloud_quota_limit",
		Help:        fmt.Sprintf("Quota limit"),
		ConstLabels: prometheus.Labels{"service": service.Id, "quota": quotaID},
	})
	limitMetric.Set(limitValue)
	prometheus.MustRegister(limitMetric)

	usageMetric := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "yandex_cloud_quota_usage",
		Help:        fmt.Sprintf("Quota usage"),
		ConstLabels: prometheus.Labels{"service": service.Id, "quota": quotaID},
	})
	usageMetric.Set(usageValue)
	prometheus.MustRegister(usageMetric)
}
