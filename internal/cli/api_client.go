// internal/cli/api_client.go
package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type APIClient struct {
	BaseURL string
}

func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{BaseURL: baseURL}
}

func (c *APIClient) ListItems() ([]map[string]any, error) {
	resp, err := http.Get(c.BaseURL + "/items")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var items []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, err
	}
	return items, nil
}

func (c *APIClient) SubmitItem(filename string, content []byte) error {
	url := c.BaseURL + "/items/upload"
	body := bytes.NewReader(content)

	// Simple upload (later switch to multipart for real files)
	resp, err := http.Post(url, "application/octet-stream", body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Upload response:", string(b))
	return nil
}

func (c *APIClient) Search(q string) ([]map[string]any, error) {
	resp, err := http.Get(c.BaseURL + "/search?q=" + q)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var results []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}
	return results, nil
}
