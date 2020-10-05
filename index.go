package cloudant

import (
	"bufio"
	"bytes"
	"errors"
	"io/ioutil"
	"strings"
)

// IndexRow contains one row from Cloudant search index
type IndexRow struct {
	ID     string      `json:"id"`
	Fields interface{} `json:"fields"`
}

func (d *Database) indexRequest(pathStr string, q *IndexQuery) (*Job, error) {
	verb := "GET"
	var body []byte
	var err error

	urlStr, err := Endpoint(*d.URL, pathStr, q.URLValues)
	if err != nil {
		return nil, err
	}

	return d.client.request(verb, urlStr, bytes.NewReader(body))
}

// indexChannel returns a channel for a given index path in which any row interface can be received
func (d *Database) indexChannel(pathStr string, q *IndexQuery) (<-chan []byte, error) {
	job, err := d.indexRequest(pathStr, q)
	if err != nil {
		if job != nil {
			job.done() // close the body reader to avoid leakage
		}
		return nil, err
	}

	err = expectedReturnCodes(job, 200)
	if err != nil {
		job.done() // close the body reader to avoid leakage
		return nil, err
	}

	results := make(chan []byte, 1000)

	go func(job *Job, results chan<- []byte) {
		defer job.Close()

		reader := bufio.NewReader(job.response.Body)

		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				close(results)
				return
			}
			lineStr := string(line)
			lineStr = strings.TrimSpace(lineStr)      // remove whitespace
			lineStr = strings.TrimRight(lineStr, ",") // remove trailing comma

			if len(lineStr) > 7 && lineStr[0:7] == "{\"id\":\"" {
				results <- []byte(lineStr)
			}
		}
	}(job, results)

	return results, nil
}

// Index returns a channel of search index documents in which matching row types can be received.
func (d *Database) Index(designName, indexName string, q *IndexQuery) (<-chan []byte, error) {
	pathStr := "/_design/" + designName + "/_search/" + indexName
	return d.indexChannel(pathStr, q)
}

// IndexRaw allows querying search indexes with arbitrary output
func (d *Database) IndexRaw(designName, indexName string, q *IndexQuery) ([]byte, error) {
	pathStr := "/_design/" + designName + "/_search/" + indexName
	job, err := d.indexRequest(pathStr, q)
	defer job.Close()

	err = expectedReturnCodes(job, 200)
	if err != nil {
		return nil, err
	}

	if job.response == nil {
		return nil, errors.New("Empty response")
	}

	return ioutil.ReadAll(job.response.Body)
}
