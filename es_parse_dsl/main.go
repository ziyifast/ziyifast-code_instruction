package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/olivere/elastic/v7"
	"github.com/ziyifast/log"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var EsClient *elasticsearch.Client

func init() {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
		Username:  "",
		Password:  "",
	})
	if err != nil {
		log.Fatalf("Error creating the Elasticsearch client: %s", err)
	}
	EsClient = client
	log.Infof("%v", client)
}

func main() {
	file, err := os.Open("demo_home/es_parse_dsl/bbb.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	header, err := reader.Read()
	if err != nil {
		panic(err)
	}
	var size, count int64
	for {
		dsl, finished, err := GetDsl(header, reader)
		if err != nil {
			log.Errorf("%v", err)
			return
		}
		newAggDsl, err := GetAggregationDsl(dsl)
		if err != nil {
			log.Errorf("%v", err)
			return
		}
		c, s, err := GetCountByDsl(newAggDsl, []string{"book-index-xxxxx"})
		if err != nil {
			log.Errorf("%v", err)
			return
		}
		count = count + c
		size = size + s
		log.Infof("size: %d, count: %d", size, count)
		if finished {
			break
		}
	}
	fmt.Println("done....", "size=", size, "count=", count)

}

func GetDsl(header []string, reader *csv.Reader) (string, bool, error) {
	//es has a length limit on the dsl
	queries := make([]elastic.Query, 0, 1024)
	finished := false
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				finished = true
				break
			}
			return "", finished, err
		}
		terms := make([]elastic.Query, 0, len(header))
		for i, term := range record {
			query := elastic.NewTermQuery(header[i], term)
			terms = append(terms, query)
		}
		query := elastic.NewBoolQuery().Must(terms...)
		queries = append(queries, query)
		if len(queries) > 1000 {
			break
		}
	}
	should := elastic.NewBoolQuery().Should(queries...)
	source, err := should.Source()
	if err != nil {
		log.Errorf("%v", err)
		return "", false, err
	}
	marshal, err := json.MarshalIndent(source, "", "\t")
	if err != nil {
		log.Errorf("%v", err)
		return "", false, err
	}
	//log.Infof("\n%s\n", string(marshal))
	return string(marshal), finished, nil
}

type aggDsl struct {
	Aggs  json.RawMessage        `json:"aggs"`
	Query map[string]interface{} `json:"query"`
}

var aggs = json.RawMessage(`{"size": {"sum": {"field": "size"}},"count":{"value_count": {"field": "_id"}}}`)

func GetAggregationDsl(dsl string) (string, error) {
	agg := new(aggDsl)
	agg.Aggs = aggs
	if len(dsl) > 0 {
		query := make(map[string]interface{})
		err := json.Unmarshal([]byte(dsl), &query)
		if err != nil {
			return "", err
		}
		agg.Query = query
	}
	indent, err := json.MarshalIndent(agg, "", "\t")
	if err != nil {
		return "", err
	}
	log.Infof("========================")
	log.Infof("\n%s\n", string(indent))
	return string(indent), nil
}

func GetCountByDsl(newDsl string, indices []string) (int64, int64, error) {
	search, err := EsClient.Search(
		EsClient.Search.WithContext(context.Background()),
		EsClient.Search.WithBody(strings.NewReader(newDsl)),
		EsClient.Search.WithIndex(indices...),
		EsClient.Search.WithTrackTotalHits(false),
	)
	if err != nil {
		log.Errorf("%v", err)
		return 0, 0, err
	}
	if search.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("response unknown code:%d", search.StatusCode)
	}
	all, err := ioutil.ReadAll(search.Body)
	if err != nil {
		log.Errorf("%v", err)
		return 0, 0, err
	}
	result := &elastic.SearchResult{}
	err = json.Unmarshal(all, result)
	if err != nil {
		err := fmt.Errorf("unExpect return:%s", string(all))
		log.Errorf("%v", err)
		return 0, 0, err

	}
	//count the files size
	//这里的size是我们开始在aggs变量中定义的聚合结果名
	size, b := result.Aggregations.Sum("size")
	if !b {
		return 0, 0, fmt.Errorf("size agg not found")
	}
	count, b := result.Aggregations.Sum("count")
	if !b {
		return 0, 0, fmt.Errorf("count agg not found")
	}
	return int64(aws.ToFloat64(count.Value)), int64(aws.ToFloat64(size.Value)), nil
}
