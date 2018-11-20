package ndb

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	baseURL               string = "api.nal.usda.gov/ndb"
	defaultSearchPageSize int    = 1500
)

// FoodIter is a type for iterating over a result set of foods from a food search.
type FoodIter struct {
	client     *HTTPClient
	err        error
	query      string
	totalItems int
	offset     int
	pageSize   int
	foods      []*Food
}

// Getter is an interface for making HTTP GET requests.
type Getter interface {
	Get(url string) (*http.Response, error)
}

// HTTPClient is a type for making HTTP requests to the NDB API.
type HTTPClient struct {
	getter Getter
	token  string
}

// NewHTTPClient returns an NDB client backed by the HTTP REST API.
func NewHTTPClient(getter Getter, token string) *HTTPClient {
	return &HTTPClient{
		getter: getter,
		token:  token,
	}
}

// FoodSearch returns a food iterator for lazily producing search results.
// The returned iterator is guaranteed to be non-nil.
func (c *HTTPClient) FoodSearch(query string) *FoodIter {
	return &FoodIter{
		client:     c,
		query:      query,
		totalItems: 0,
		offset:     0,
		pageSize:   defaultSearchPageSize,
		foods:      []*Food{},
	}
}

// Next returns the next food in the iterator.
// Callers should check it.Err() if Next returns nil.
func (it *FoodIter) Next() *Food {
	if it.offset%it.pageSize == 0 {
		fsResp, err := it.client.GetFoodSearchPage(it.query, it.offset, it.pageSize)
		if err != nil {
			it.err = errors.Wrap(err, "Getting food search page.")
			return nil
		}
		if fsResp.Results == nil {
			it.err = errors.Errorf("No results for query \"%s\".", it.query)
			return nil
		}
		it.totalItems = fsResp.Results.Total
		it.foods = append(it.foods, fsResp.Results.Foods...)

	}

	if it.offset >= it.totalItems {
		it.err = io.EOF
		return nil
	}

	defer func() {
		it.offset += 1
	}()

	return it.foods[it.offset]
}

// Err returns any error encountered during iteration.
// It returns io.EOF if the iterator is exhausted.
// Otherwise, it returns a descriptive error of the issue.
// It returns nil if there was no error.
func (it *FoodIter) Err() error {
	return it.err
}

// GetFoodSearchPage gets a page of foods matching the given query.
// It returns an error if any encoding/decoding errors are encountered,
// or there was an issue making the request to the NDB API.
func (c *HTTPClient) GetFoodSearchPage(query string, offset int, limit int) (*FoodSearchResponse, error) {
	url := fmt.Sprintf(
		"http://%s/search?api_key=%s&q=%s&offset=%d&max=%d",
		baseURL,
		c.token,
		url.QueryEscape(query),
		offset,
		limit,
	)

	resp, err := c.getter.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "Making search request to NDB.")
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.Errorf("Status code %d from NDB Search", resp.StatusCode)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Reading response body from NDB Search.")
	}

	searchResp := new(FoodSearchResponse)
	if err := json.Unmarshal(respBytes, searchResp); err != nil {
		return nil, errors.Wrap(err, "Unmarshaling NDB search response.")
	}

	return searchResp, nil
}

// FoodReport returns an NDB food report for the food with the given ndbno.
// It returns an error if the request was unsuccessful.
func (c *HTTPClient) FoodReport(ndbno string) (*Food, error) {
	url := fmt.Sprintf(
		"http://%s/reports?api_key=%s&ndbno=%s",
		baseURL,
		c.token,
		ndbno,
	)

	resp, err := c.getter.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "Getting NDB food report.")
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.Errorf("Status code %d from NDB food report.", resp.StatusCode)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Reading response body from NDB food report.")
	}

	foodReportResp := new(FoodReportResponse)
	if err := json.Unmarshal(respBytes, foodReportResp); err != nil {
		return nil, errors.Wrap(err, "Unmarshaling food report response body.")
	}

	return foodReportResp.Report.Food, nil
}
