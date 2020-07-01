package photonanalytics

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"path"
	"strconv"
)

type GetAppGraphRequest struct {
	AppID    string
	Region   string
	Template string
	Start    string
	End      string
	Width    int
	Height   int
}

func (c *APIClient) GetAppGraph(req *GetAppGraphRequest) (image.Image, error) {
	if req.AppID == "" || req.Region == "" || req.Template == "" {
		return nil, errors.New("appId, region, template must be set")
	}

	params := &Params{
		queries: make(map[string]string),
	}
	if req.Start != "" {
		params.queries["start"] = req.Start
	}
	if req.End != "" {
		params.queries["end"] = req.End
	}
	if req.Width > 0 {
		params.queries["width"] = strconv.Itoa(req.Width)
	}
	if req.Height > 0 {
		params.queries["height"] = strconv.Itoa(req.Height)
	}

	rpath := path.Join("graph", "app", req.AppID, req.Region, req.Template)
	ctx := context.Background()
	r, err := c.newRequest(ctx, http.MethodGet, rpath, params)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status: %s", resp.Status)
	}

	return png.Decode(resp.Body)
}

type GetAppValueRequest struct {
	AppID    string
	Region   string
	Template string
	Start    string
	End      string
}

func (c *APIClient) GetAppValue(req *GetAppValueRequest) (float64, error) {
	if req.AppID == "" || req.Region == "" || req.Template == "" {
		return 0, errors.New("appId, region, template must be set")
	}

	params := &Params{
		queries: make(map[string]string),
	}
	if req.Start != "" {
		params.queries["start"] = req.Start
	}
	if req.End != "" {
		params.queries["end"] = req.End
	}

	rpath := path.Join("data", "app", req.AppID, req.Region, req.Template)
	ctx := context.Background()
	r, err := c.newRequest(ctx, http.MethodGet, rpath, params)
	if err != nil {
		return 0, err
	}

	resp, err := c.HTTPClient.Do(r)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("invalid status: %s", resp.Status)
	}

	var b bytes.Buffer
	if _, err := io.Copy(&b, resp.Body); err != nil {
		return 0, err
	}
	return strconv.ParseFloat(b.String(), 64)
}

type GetMultipleAppsValueRequest struct {
	Queries []BulkdataAppQuery
}

type BulkdataAppQuery struct {
	ID       string `json:"id"`
	AppID    string `json:"appid"`
	Cloud    string `json:"cloud"`
	Cluster  string `json:"cluster"`
	Region   string `json:"region"`
	Template string `json:"template"`
}

func (c *APIClient) GetMultipleAppsValue(req *GetMultipleAppsValueRequest) (
	resp map[string]int64,
	_ error,
) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&req.Queries); err != nil {
		return nil, err
	}

	params := &Params{
		headers: map[string]string{
			"Content-Type": "application/json",
		},
		body: &buf,
	}
	ctx := context.Background()
	r, err := c.newRequest(ctx, http.MethodPost, "bulkdata/app", params)
	if err != nil {
		return nil, err
	}

	httpResp, err := c.HTTPClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()
	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status: %s", httpResp.Status)
	}

	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, err
	}
	return resp, nil
}

type GetMultipleAppsValueSeriesWithSpanRequest struct {
	Queries []BulkxportAppQuery `json:"queries"`
}

type BulkxportAppQuery struct {
	ID        string `json:"id"`
	AppID     string `json:"appid"`
	Cloud     string `json:"cloud"`
	Cluster   string `json:"cluster"`
	Region    string `json:"region"`
	Template  string `json:"template"`
	Start     string `json:"start,omitempty"`
	End       string `json:"end,omitempty"`
	Xporttime string `json:"xporttim,omitemptye"`
	Normalize string `json:"normalize,omitempty"`
}

type GetMultipleAppsValueSeriesWithSpanResponse struct {
}

func (c *APIClient) GetMultipleAppsValueSeriesWithSpan(
	req *GetMultipleAppsValueSeriesWithSpanRequest,
) (
	resp map[string]map[string][]int64,
	_ error,
) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&req.Queries); err != nil {
		return nil, err
	}

	params := &Params{
		headers: map[string]string{
			"Content-Type": "application/json",
		},
		body: &buf,
	}
	ctx := context.Background()
	r, err := c.newRequest(ctx, http.MethodPost, "bulkxport/app", params)
	if err != nil {
		return nil, err
	}

	httpResp, err := c.HTTPClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()
	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status: %s", httpResp.Status)
	}

	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, err
	}
	return resp, nil
}

type GetMultipleAppsValueAndValueSeriesWithSpanRequest struct {
	Queries []BulkAppQuery `json:"queries"`
}

type BulkAppQuery struct {
	ID        string `json:"id"`
	AppID     string `json:"appid"`
	Cloud     string `json:"cloud"`
	Cluster   string `json:"cluster"`
	Region    string `json:"region"`
	Template  string `json:"template"`
	Start     string `json:"start,omitempty"`
	End       string `json:"end,omitempty"`
	Xporttime string `json:"xporttim,omitemptye"`
	Normalize string `json:"normalize,omitempty"`
}

type BulkAllEntry struct {
	Data  int64
	Xport map[string][]int64
}

func (c *APIClient) GetMultipleAppsValueAndValueSeriesWithSpan(
	req *GetMultipleAppsValueAndValueSeriesWithSpanRequest,
) (
	resp map[string]BulkAllEntry,
	_ error,
) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&req.Queries); err != nil {
		return nil, err
	}

	params := &Params{
		headers: map[string]string{
			"Content-Type": "application/json",
		},
		body: &buf,
	}
	ctx := context.Background()
	r, err := c.newRequest(ctx, http.MethodPost, "bulk/app", params)
	if err != nil {
		return nil, err
	}

	httpResp, err := c.HTTPClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()
	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status: %s", httpResp.Status)
	}

	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, err
	}
	return resp, nil
}
