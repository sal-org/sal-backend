package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	Model "salbackend/model"
	"strings"
	"time"
)

func GetWordPressPageBySlug(baseURL, slug string) (*Model.WordPressPage, error) {

	apiURL := fmt.Sprintf(
		"%s/wp-json/wp/v2/pages?slug=%s",
		baseURL,
		url.QueryEscape(slug),
	)

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "MentalHealthCMS/1.0")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("wordpress error: %s", string(body))
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected content type %q: %s", contentType, string(body))
	}

	var pages []Model.WordPressPage

	if err := json.NewDecoder(resp.Body).Decode(&pages); err != nil {
		return nil, err
	}

	if len(pages) == 0 {
		return nil, fmt.Errorf("page not found")
	}

	return &pages[0], nil
}

func GetWordPressPageByID(baseURL string, pageID int) (*Model.WordPressPage, error) {

	apiURL := fmt.Sprintf(
		"%s/wp-json/wp/v2/pages/%d",
		baseURL,
		pageID,
	)

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "MentalHealthCMS/1.0")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("wordpress error: %s", string(body))
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected content type %q: %s", contentType, string(body))
	}

	var page Model.WordPressPage

	if err := json.NewDecoder(resp.Body).Decode(&page); err != nil {
		return nil, err
	}

	return &page, nil
}
