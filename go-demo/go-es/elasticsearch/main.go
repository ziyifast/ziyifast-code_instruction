package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/estransport"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	url      = []string{"http://test.ziyi.com"}
	username = ""
	password = ""
	sort     = json.RawMessage(`[{"_id":{"order":"desc"}}]`)
	aggs     = json.RawMessage(`{"size": {"sum": {"field": "size"}},"count":{"value_count": {"field": "_id"}}}`)
	size     = 2
	indices  = []string{"test-ziyi-1-100004-100136"}
)

func main() {
	esClient, err := CreateClient(url, username, password)
	if err != nil {
		panic(err)
	}
	var searchAfter []interface{}

	for {
		dsl := Dsl{
			Sort:        sort,
			Size:        size,
			SearchAfter: searchAfter,
			Query: map[string]interface{}{
				"bool": map[string]interface{}{
					"must": map[string]interface{}{
						"wildcard": map[string]interface{}{
							"book_name": "士",
						},
						//"match_all": map[string]interface{}{},
					},
				},
			},
		}
		queryJson, err := json.MarshalIndent(dsl, "", "\t")
		if err != nil {
			panic(err)
		}
		fmt.Printf("queryJson:%s\n", queryJson)
		res, err := esClient.Search(
			esClient.Search.WithContext(context.Background()),
			esClient.Search.WithIndex(indices...),
			esClient.Search.WithBody(strings.NewReader(string(queryJson))),
			esClient.Search.WithTrackTotalHits(false),
		)
		if err != nil {
			panic(err)
		}
		var result struct {
			Hits struct {
				Hits []struct {
					Index  string                 `json:"_index"`
					ID     string                 `json:"_id"`
					Sort   []interface{}          `json:"sort"`
					Source map[string]interface{} `json:"_source"`
				} `json:"hits"`
			} `json:"hits"`
		}

		if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
			panic(err)
		}
		err = res.Body.Close()
		if err != nil {
			panic(err)
		}
		if len(result.Hits.Hits) > 0 {
			lastHit := result.Hits.Hits[len(result.Hits.Hits)-1]
			searchAfter = lastHit.Sort
		} else {
			break
		}
		for _, h := range result.Hits.Hits {
			fmt.Printf("=====id:%s book_name:%s\n", h.ID, h.Source["book_name"])
		}
	}
}

type Dsl struct {
	Sort        json.RawMessage        `json:"sort"`
	Size        int                    `json:"size"`
	SearchAfter []interface{}          `json:"search_after,omitempty"`
	Query       map[string]interface{} `json:"query"`
}

func CreateClient(url []string, username, password string) (*elasticsearch.Client, error) {
	es, err := elasticsearch.NewClient(genConfig(url, username, password))
	if err != nil {
		panic(err)
		return nil, err
	}
	res, err := es.Info()
	if err != nil {
		panic(err)
		return nil, err
	}
	defer res.Body.Close()
	return es, nil

}

type DecoratedTransport struct {
	tp http.RoundTripper
}

func (d *DecoratedTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	request.Host = "test.ziyi.com"
	return d.tp.RoundTrip(request)
}

func genConfig(url []string, username, password string) elasticsearch.Config {
	retryBackoff := backoff.NewExponentialBackOff()
	cfg := elasticsearch.Config{
		Addresses:     url,
		Logger:        &estransport.ColorLogger{Output: os.Stdout},
		Username:      username,
		Password:      password,
		RetryOnStatus: []int{502, 503, 504, 429},
		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				retryBackoff.Reset()
			}
			return retryBackoff.NextBackOff()
		},
		MaxRetries: 5,
		Transport: &DecoratedTransport{
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
		},
	}
	return cfg
}
