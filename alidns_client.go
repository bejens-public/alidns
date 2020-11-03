package alidns

import (
	"strings"
	"sync"
	"time"

	aliclouddns "github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/libdns/libdns"
)

type alidnsClient struct {
	RegionID        string `json:"region_id"`
	AccessKeyID     string `json:"access_key"`
	AccessKeySecret string `json:"access_secret"`

	alidnsClient *aliclouddns.Client

	records domainRecords

	clientMutex sync.Mutex
}

func (adc *alidnsClient) getClient() error {

	if adc.alidnsClient == nil {

		adc.clientMutex.Lock()
		defer adc.clientMutex.Unlock()

		if adc.alidnsClient == nil {
			client, err := aliclouddns.NewClientWithAccessKey(adc.RegionID, adc.AccessKeyID, adc.AccessKeySecret)
			if err != nil {
				return err
			}
			adc.alidnsClient = client
		}

	}

	return nil
}

func (adc *alidnsClient) getRecords(zone string) ([]libdns.Record, error) {

	domainName := strings.Trim(zone, ".")

	if domainRecords, ok := adc.records.load(domainName); ok {
		return domainRecords.([]libdns.Record), nil
	}

	if err := adc.getClient(); err != nil {
		return nil, err
	}

	request := aliclouddns.CreateDescribeDomainRecordsRequest()
	request.Scheme = "https"
	request.Domain = domainName

	response, err := adc.alidnsClient.DescribeDomainRecords(request)
	if err != nil {
		return nil, err
	}

	var records []libdns.Record
	for _, domainRecord := range response.DomainRecords.Record {
		record := libdns.Record{
			ID:    domainRecord.RecordId,
			Type:  domainRecord.Type,
			Name:  domainRecord.DomainName,
			Value: domainRecord.Value,
			TTL:   time.Duration(domainRecord.TTL),
		}
		records = append(records, record)
	}

	adc.domains.Store(domainName, records)
	return records, nil
}
