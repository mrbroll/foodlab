package ndb

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	baseURL            string = "api.nal.usda.gov/ndb"
	defaultSearchLimit int    = 50
	searchDataSource   string = "Standard Reference"
)

type Getter interface {
	Get(url string) (*http.Response, error)
}

type HTTPClient struct {
	getter Getter
	token  string
}

func NewHTTPClient(getter Getter, token string) *HTTPClient {
	return &HTTPClient{
		getter: getter,
		token:  token,
	}
}

// FoodSearch searches for foods matching the given query.
// It returns an error if any encoding/decoding errors are encountered,
// or there was an issue making the request to the NDB API.
func (c *HTTPClient) FoodSearch(query string) ([]*Food, error) {
	url := fmt.Sprintf(
		"http://%s/search?api_key=%s&q=%s&ds=%s&max=%d",
		baseURL,
		c.token,
		url.QueryEscape(query),
		// url.QueryEscape(searchDataSource),
		defaultSearchLimit,
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

	if searchResp.Results == nil {
		return nil, nil
	}

	return searchResp.Results.Foods, nil
}

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
