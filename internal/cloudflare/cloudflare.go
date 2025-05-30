package cloudflare

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

type CloudflareConfig struct {
	ZoneID string
	APIKey string
}

type CloudflareDnsRecords struct {
	Result     []*CloudflareDnsRecord `json:"result"`
}

type CloudflareDnsRecord struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Content    string   `json:"content"`
	Proxiable  bool     `json:"proxiable"`
	Proxied    bool     `json:"proxied"`
}

func CheckDNSExists(config *CloudflareConfig, dns string) bool {
	log.Infof("Checking if dns record %v exsits", dns)
	var url = fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%v/dns_records", config.ZoneID)
	log.Infof("Cloudflare api url: %v", url)
	
	client := resty.New()
	dnsRecords := &CloudflareDnsRecords{}
	res, err := client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(config.APIKey).
		SetResult(dnsRecords).
		Get(url)

	if err != nil {
		log.Error("Error calling cloudflare api", err)
		panic(err)
	}

	if res.IsError() {
		log.Errorf("Error response from cloudflare api status=%v body=%v", res.StatusCode(), res.String())
		panic(res)
	}

	log.Info("Successfully got dns records")
	for _, dnsRecord := range dnsRecords.Result {
		if dnsRecord.Name == dns {
			log.Infof("Dns record %v exists", dns)
			return true
		}
	}

	log.Infof("Dns record %v does not exist", dns)
	return false
}