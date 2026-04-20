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

type typesenseDriver struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

func newTypesenseDriver(baseURL, apiKey string, hc *http.Client) Driver {
	return &typesenseDriver{
		baseURL: strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		apiKey:  strings.TrimSpace(apiKey),
		http:    hc,
	}
}

func (t *typesenseDriver) IndexName(logical string) string { return logical }

func (t *typesenseDriver) UpsertDocuments(ctx context.Context, logicalIndex string, docs []Document) error {
	if len(docs) == 0 {
		return nil
	}
	name := t.IndexName(logicalIndex)
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	for _, d := range docs {
		row := map[string]any{"id": d.ID}
		for k, v := range d.Fields {
			row[k] = v
		}
		if err := enc.Encode(row); err != nil {
			return err
		}
	}
	u := t.baseURL + "/collections/" + url.PathEscape(name) + "/documents/import?action=upsert"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(buf.Bytes()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-ndjson")
	req.Header.Set("X-TYPESENSE-API-KEY", t.apiKey)
	resp, err := t.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("typesense: %s: %s", resp.Status, string(b))
	}
	return nil
}

func (t *typesenseDriver) DeleteDocument(ctx context.Context, logicalIndex, id string) error {
	name := t.IndexName(logicalIndex)
	u := t.baseURL + "/collections/" + url.PathEscape(name) + "/documents/" + url.PathEscape(id)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-TYPESENSE-API-KEY", t.apiKey)
	resp, err := t.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil
	}
	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("typesense delete: %s: %s", resp.Status, string(b))
	}
	return nil
}
