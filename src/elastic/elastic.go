package elastic

import (
	"context"
	"crypto/tls"
	"elasticbulk/common"
	"elasticbulk/settings"
	"fmt"
	"github.com/olivere/elastic"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

var client *elastic.Client

type Network struct {
	CreateTime string `json:"createTime"`
	Host string `json:"host"`
	Name string `json:"name"`
	NodeId string `json:"nodeId"`
	Offset int `json:"offset"`
	ReportTime string `json:"reportTime"`
	StrictDateTime string `json:"strictDateTime"`
	Territory string `json:"territory"`

	Data Data
	Fields Fields

}

type Data struct {
	Count int `json:"count"`
	DhcpUsed int `json:"dhcpused"`
	IpamUsed int `json:"ipamused"`
	Total int `json:"total"`
	Used int `json:"used"`
}

type Fields struct {
	Tag string `json:"tag"`
}

func InitElastic() error {
	errorLog := log.New(os.Stdout, "elasticbulk", log.LstdFlags)
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	var err error
	client, err = elastic.NewClient(
		elastic.SetErrorLog(errorLog),
		elastic.SetHttpClient(httpClient),
		elastic.SetURL(settings.App.Elastic.Url),
		elastic.SetScheme("https"),
		elastic.SetBasicAuth(settings.App.Elastic.UserName, settings.App.Elastic.Password))
	if err != nil {
		return err
	}
	info, code, err := client.Ping(settings.App.Elastic.Url).Do(context.Background())
	if err != nil {
		return err
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s \n", code, info.Version.Number)

	version, err := client.ElasticsearchVersion(settings.App.Elastic.Url)
	if err != nil {
		return err
	}
	fmt.Printf("Elasticsearch version %s \n", version)

	return nil
}

func BulkInsert () {
	//data := Data{1, 2, 2, 2, 2}
	//fields := Fields{"network"}
	//network := Network{"2021-12-08 22:55:51", "localhost.localdomain", "192.168.100.0/24", "61adba5dab40b94f4432b8cb", 751044, "1638888651000", "2021-12-08T12:55:51.700+0800", "河南--周口", data, fields}

	networkList := CreateNetworkList()

	bulkRequest := client.Bulk()
	for _, network := range networkList {
		bulkRequest.Add(elastic.NewBulkIndexRequest().Index(settings.App.Elastic.IndexName).Type(settings.App.Elastic.Index).Doc(network))
	}
	response, err := bulkRequest.Do(context.Background())
	//result, err := client.Index().Index("network-2021.12.09").Type("network").BodyJson(networkList).Do(context.Background())
	if err != nil {
		fmt.Printf(err.Error())
	}
	failed := response.Failed()
	l := len(failed)
	if l > 0 {
		fmt.Printf("Error(%d)", l, response.Errors)
	}
	// fmt.Printf("Indexed tweet %s to index s%s, type %s\n", result.Id, result.Index, result.Type)
}

func CreateNetworkList () []Network  {
	config := settings.App.Elastic
	// 计算2个时间相差多少分钟
	differ := common.GetMinutesDiffer(config.From, config.To)
	// 计算次方
	count := math.Pow(2, float64(32 - config.Mask))
	fromTime, _ := time.Parse("2006-01-02 15:04:05", config.From)
	networks := make([]Network, differ)
	for i := 0; i < differ; i++ {
		// ip段使用数
		total := common.RandInt64(1, int64(count))
		// dhcp使用数
		used :=common.RandInt64(1, total)

		ipamused := math.Ceil(float64(total) / count * 100)
		dhcpused := math.Ceil((float64(used) / float64(total)) * 100)

		data := Data{int(count), int(dhcpused), int(ipamused), int(total), int(used)}
		fields := Fields{config.Index }

		// 每次加一分钟
		m, _ := time.ParseDuration(strconv.Itoa(i + 1) + "m")
		m1 := fromTime.Add(m)

		networks[i] = Network{m1.Format("2006-01-02 03:04:05"), "localhost.localdomain", config.Ip + "/" + strconv.Itoa(config.Mask), config.NodeId, config.Offset + i, strconv.Itoa(int(m1.UnixNano() / 1e6 - 28800000)), m1.Format("2006-01-02T15:04:05.000+0800"), config.Territory, data, fields}
	}
	return networks
}


