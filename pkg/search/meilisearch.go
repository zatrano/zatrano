package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type meilisearchDriver struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

func newMeilisearchDriver(baseURL, apiKey string, hc *http.Client) Driver {
	return &meilisearchDriver{
		baseURL: strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		apiKey:  strings.TrimSpace(apiKey),
		http:    hc,
	}
}

func (m *meilisearchDriver) IndexName(logical string) string { return logical }

func (m *meilisearchDriver) UpsertDocuments(ctx context.Context, logicalIndex string, docs []Document) error {
	if len(docs) == 0 {
		return nil
	}
	uid := m.IndexName(logicalIndex)
	payload := make([]map[string]any, 0, len(docs))
	for _, d := range docs {
		row := map[string]any{"id": d.ID}
		for k, v := range d.Fields {
			row[k] = v
		}
		payload = append(payload, row)
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, m.baseURL+"/indexes/"+url.PathEscape(uid)+"/documents", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if m.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+m.apiKey)
	}
	resp, err := m.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("meilisearch: %s: %s", resp.Status, string(b))
	}
	return nil
}

func (m *meilisearchDriver) DeleteDocument(ctx context.Context, logicalIndex, id string) error {
	uid := m.IndexName(logicalIndex)
	u := m.baseURL + "/indexes/" + url.PathEscape(uid) + "/documents/" + url.PathEscape(id)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return err
	}
	if m.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+m.apiKey)
	}
	resp, err := m.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil
	}
	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("meilisearch delete: %s: %s", resp.Status, string(b))
	}
	return nil
}
