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

// ignoring some fields here
type CloudflareDnsRecords struct {
	Result []*CloudflareDnsRecord `json:"result"`
}

// ignoring some fields here
type CloudflareDnsRecord struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Content    string   `json:"content"`
	Proxiable  bool     `json:"proxiable"`
	Proxied    bool     `json:"proxied"`
}

func CheckDnsExists(config *CloudflareConfig, dns string) bool {
	log.Infof("Checking if dns record %v exsits", dns)
	dnsRecords := getDnsRecords(config)
	for _, dnsRecord := range dnsRecords.Result {
		if dnsRecord.Name == dns {
			log.Infof("Dns record %v exists", dns)
			return true
		}
	}

	log.Infof("Dns record %v does not exist", dns)
	return false
}

func getDnsRecords(config *CloudflareConfig) *CloudflareDnsRecords {
	var cloudflareUrl = fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%v/dns_records", config.ZoneID)
	log.Infof("Cloudflare get dns records url: %v", cloudflareUrl)
	
	client := resty.New()
	dnsRecords := &CloudflareDnsRecords{}
	res, err := client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(config.APIKey).
		SetResult(dnsRecords).
		Get(cloudflareUrl)

	if err != nil {
		log.Errorf("Error calling cloudflare api: %v", err)
		panic(err)
	}

	if res.IsError() {
		log.Errorf("Error response from cloudflare api status=%v body=%v", res.StatusCode(), res.String())
		panic(res)
	}

	log.Info("Successfully got dns records")
	return dnsRecords
}

func UpdateCloudflareForDynamicIp(config *CloudflareConfig) {
	log.Info("Comparing ip of dns records to cluster nodes")
	ip := getCurrentIp()
	dnsRecords := getDnsRecords(config)

	for _, dnsRecord := range dnsRecords.Result {
		if dnsRecord.Type == "A" && dnsRecord.Content != ip {
			log.Infof("Mismatch dns ip %v with current nodes ip %v -- Patching record", dnsRecord.Content, ip)
			dnsRecord.Content = ip
			patchDnsRecord(config, dnsRecord)
		}
	}
}

func getCurrentIp() string {
	log.Infof("Getting current ip")
	const ipUrl = "https://api.ipify.org"

	client := resty.New()
	res, err := client.R().Get(ipUrl)

	if err != nil {
		log.Errorf("Error calling cloudflare api: %v", err)
		panic(err)
	}

	if res.IsError() {
		log.Errorf("Error response from cloudflare api status=%v body=%v", res.StatusCode(), res.String())
		panic(res)
	}

	ip := res.String()
	log.Infof("IP: %v", ip)
	return ip
}

func patchDnsRecord(config *CloudflareConfig, dnsRecord *CloudflareDnsRecord) {
	var cloudflareUrl = fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%v/dns_records/%v", config.ZoneID, dnsRecord.ID)
	log.Infof("Cloudflare patch dns record url: %v", cloudflareUrl)

	client := resty.New()
	res, err := client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(config.APIKey).
		SetBody(dnsRecord).
		Patch(cloudflareUrl)

	if err != nil {
		log.Errorf("Error calling cloudflare api: %v", err)
		panic(err)
	}

	if res.IsError() {
		log.Errorf("Error response from cloudflare api status=%v body=%v", res.StatusCode(), res.String())
		panic(res)
	}

	log.Infof("Successfully updated record name=%v", dnsRecord.Name)
}

// func CreateDnsRecordFromIngress() {
// 	dnsRecord := &CloudflareDnsRecord{
// 		Name: 
// 	}
// }

func createDnsRecord(config *CloudflareConfig, dnsRecord *CloudflareDnsRecord) {
	var cloudflareUrl = fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%v/dns_records/%v", config.ZoneID, dnsRecord.ID)
	log.Infof("Cloudflare patch dns record url: %v", cloudflareUrl)

	client := resty.New()
	res, err := client.R().
		SetHeader("Accept", "application/json").
		SetAuthToken(config.APIKey).
		SetBody(dnsRecord).
		Put(cloudflareUrl)

	if err != nil {
		log.Errorf("Error calling cloudflare api: %v", err)
		panic(err)
	}

	if res.IsError() {
		log.Errorf("Error response from cloudflare api status=%v body=%v", res.StatusCode(), res.String())
		panic(res)
	}

	log.Infof("Successfully updated record name=%v", dnsRecord.Name)
}