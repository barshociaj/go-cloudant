package cloudant

import (
	"encoding/json"
	"net/url"
	"strconv"
)

// IndexQuery is a helper utility to build Cloudant request parameters for views (including _all_docs)
//
// Example:
// 	q := NewIndexQuery().
//     Query("title:abc")
//
//	docs, err := db.All(q)

// IndexQuery object helps build Cloudant IndexQuery parameters
type IndexQuery struct {
	URLValues url.Values
}

// NewIndexQuery is a shortcut to create new Cloudant IndexQuery object with no parameters
func NewIndexQuery() *IndexQuery {
	return &IndexQuery{URLValues: url.Values{}}
}

// Query applies q=(query) parameter to Cloudant IndexQuery
func (q *IndexQuery) Query(query string) *IndexQuery {
	if query != "" {
		q.URLValues.Set("q", query)
	}
	return q
}

// IncludeDocs applies include_docs=true parameter to Cloudant IndexQuery
func (q *IndexQuery) IncludeDocs() *IndexQuery {
	q.URLValues.Set("include_docs", "true")
	return q
}

// Bookmark applies bookmark=(bookmark) parameter to Cloudant IndexQuery
func (q *IndexQuery) Bookmark(bookmark string) *IndexQuery {
	if bookmark != "" {
		q.URLValues.Set("bookmark", bookmark)
	}
	return q
}

// Limit applies limit parameter to Cloudant IndexQuery
func (q *IndexQuery) Limit(lim int) *IndexQuery {
	if lim > 0 {
		q.URLValues.Set("limit", strconv.Itoa(lim))
	}
	return q
}

// Stale applies stale=ok parameter to Cloudant IndexQuery
func (q *IndexQuery) Stale() *IndexQuery {
	q.URLValues.Set("stale", "ok")
	return q
}

// IncludeFields applies include_fields=(fields) parameter to Cloudant IndexQuery
func (q *IndexQuery) IncludeFields(fields []string) *IndexQuery {
	if len(fields) > 0 {
		data, err := json.Marshal(fields)
		if err == nil {
			q.URLValues.Set("include_fields", string(data[:]))
		}
	}
	return q
}
