package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"net"
	"net/http"
	"time"
)

var (
	host = "http://test.ziyi.com"
)

func main() {
	esClient := CreateEsClient()
	fmt.Println(esClient)
	//1. 操作索引
	//isExistIndex(esClient)
	//createIndex(esClient)
	//deleteIndex(esClient)

	//2. 操作doc文档（记录）
	//addDoc(esClient)
	//updateDoc(esClient)
	//deleteDoc(esClient)

	//3. 批处理请求
	//sendBulkRequest(esClient)

	//4. 查询
	//simpleSearch(esClient)
	//searchAfterSearch(esClient)

}

func searchAfterSearch(esClient *elastic.Client) {
	var lastHit *elastic.SearchHit
	for {
		q := elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("book_name", "士"))
		searchSource := elastic.NewSearchSource().Query(q).Size(2).Sort("_id", false)
		if lastHit != nil {
			fmt.Printf("search After %+v\n", lastHit.Sort)
			searchSource.SearchAfter(lastHit.Sort...)
		}
		dsl, err := searchSource.MarshalJSON()
		if err != nil {
			panic(err)
		}
		fmt.Printf("dsl %s\n", string(dsl))
		searchResult, err := esClient.Search().Index("test-ziyi-1-100004-100136").SearchSource(searchSource).Do(context.Background())
		if err != nil {
			panic(err)
		}
		if len(searchResult.Hits.Hits) == 0 {
			fmt.Println("no more data")
			break
		}
		for _, hit := range searchResult.Hits.Hits {
			res := make(map[string]interface{})
			if err = json.Unmarshal(hit.Source, &res); err != nil {
				panic(err)
			}
			fmt.Printf("search %s %s\n", hit.Id, res["author"])
		}
		lastHit = searchResult.Hits.Hits[len(searchResult.Hits.Hits)-1]
	}
}

func simpleSearch(esClient *elastic.Client) {
	response, err := esClient.Search([]string{"test-ziyi-1-100004-100136"}...).Query(elastic.NewTermQuery("author", "袁朗")).Size(100).Do(context.TODO())
	if err != nil {
		panic(err)
	}
	fmt.Println(response.Hits.Hits)
}

func sendBulkRequest(esClient *elastic.Client) {
	bulkRequest := esClient.Bulk()
	for i := 0; i < 10; i++ {
		docMappings := map[string]interface{}{
			"book_name":    fmt.Sprintf("士兵突击-%d", i),
			"author":       "aa",
			"on_sale_time": "2000-01-05",
			"book_desc":    "一个关于部队的...",
		}
		bulkRequest = bulkRequest.Add(elastic.NewBulkIndexRequest().Index("test-ziyi-1-100004-100136").Doc(docMappings))
	}
	bulkResponse, err := bulkRequest.Do(context.Background())
	if err != nil {
		panic(err)
	}
	if bulkResponse.Errors {
		for _, item := range bulkResponse.Items {
			for _, action := range item {
				if action.Error != nil {
					fmt.Printf("Error for item: %s: %s", action.Error.Index, action.Error.Reason)
				}
			}
		}
	} else {
		fmt.Println("All bulk requests executed successfully")
	}
}

func deleteDoc(esClient *elastic.Client) {
	//response返回删除的doc Id，如果要删除的doc不存在，则直接返回err not found
	response, err := esClient.Delete().Index("test-ziyi-1-100004-100136").Id("1").Do(context.Background())
	if err != nil {
		panic(err)
	}
	println(response.Id)
}

func updateDoc(esClient *elastic.Client) {
	type m map[string]interface{}
	docMappings := m{
		"book_name":    "士兵突击",
		"author":       "袁朗",
		"word_count":   100000,
		"on_sale_time": "2000-01-05",
		"book_desc":    "一个关于部队的...",
	}
	//覆盖式修改(response返回doc记录的id)
	response, err := esClient.Update().Index("test-ziyi-1-100004-100136").Id("1").Doc(docMappings).Do(context.Background())
	//指定字段修改
	//response, err := esClient.Update().Index("test-ziyi-1-100004-100136").Id("1").Doc(map[string]interface{}{
	//	"book_name": "我的团长我的团",
	//}).Do(context.Background())
	if err != nil {
		panic(err)
	}
	println(response.Id)
}

func addDoc(esClient *elastic.Client) {
	type m map[string]interface{}
	documentMappings := m{
		"book_name":    "士兵突击",
		"author":       "兰晓龙",
		"word_count":   100000,
		"on_sale_time": "2000-01-05",
		"book_desc":    "一个关于部队的...",
	}
	//如果不指定id，则es会自动生成一个id（杂乱无序不好维护），response为返回的文档id
	//response, err := esClient.Index().Index("test-ziyi-1-100004-100136").BodyJson(documentMappings).Do(context.Background())
	response, err := esClient.Index().Index("test-ziyi-1-100004-100136").Id("1").BodyJson(documentMappings).Do(context.Background()) //指定id
	if err != nil {
		panic(err)
	}
	println(response.Id)
}

func deleteIndex(esClient *elastic.Client) {
	response, err := esClient.DeleteIndex("test-ziyi-1-100004-100136").Do(context.Background())
	if err != nil {
		panic(err)
	}
	println(response.Acknowledged)
}

func createIndex(esClient *elastic.Client) {
	type m map[string]interface{}
	indexMapping := m{
		"settings": m{
			"number_of_shards":   5, //分片数
			"number_of_replicas": 1, //副本数
		},
		"mappings": m{
			"properties": m{ //索引属性值
				"book_name": m{ //索引属性名
					"type": "text", //filed类型
					//"analyzer": "ik_max_word", //使用ik分词器进行分词
					"index": true,  //当前field可以被用于查询条件
					"store": false, //是否额外存储
				},
				"author": m{
					"type": "keyword", //作为关键字不分词
				},
				"word_count": m{
					"type": "long",
				},
				"on_sale_time": m{
					"type":   "date",
					"format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis",
				},
				"book_desc": m{
					"type": "text",
					//"analyzer": "ik_max_word",
				},
			},
		},
	}
	result, err := esClient.CreateIndex("test-ziyi-1-100004-100136").BodyJson(indexMapping).Do(context.Background())
	if err != nil {
		panic(err)
	}
	if result.Acknowledged {
		println("create index success")
	} else {
		println("create index failed")
	}
}

func isExistIndex(esClient *elastic.Client) {
	isExist, err := esClient.IndexExists("test-ziyi-1-100004-100136").Do(context.TODO())
	if err != nil {
		panic(err)
	}
	if isExist {
		println("index exists")
	} else {
		println("index not exists")
	}
}

func CreateEsClient() *elastic.Client {
	esClient, err := elastic.NewClient(
		elastic.SetURL(host),
		elastic.SetSniff(false),
		elastic.SetBasicAuth("", ""),
		elastic.SetHttpClient(&http.Client{Transport: &DecoratedTransport{
			tp: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				ForceAttemptHTTP2:     true,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}}),
	)
	if err != nil {
		panic(err)
	}
	return esClient
}

type DecoratedTransport struct {
	tp http.RoundTripper
}

func (d *DecoratedTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	request.Host = "test.ziyi.com"
	return d.tp.RoundTrip(request)
}
