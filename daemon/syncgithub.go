package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// ghAPIBase is the GitHub REST API root. Tests override it with an httptest URL.
var ghAPIBase = "https://api.github.com"

type ghContentEntry struct {
	Name string `json:"name"`
	SHA  string `json:"sha"`
	Type string `json:"type"`
}

type ghContentFile struct {
	Content string `json:"content"`
	SHA     string `json:"sha"`
	Name    string `json:"name"`
}

type ghPutResponse struct {
	Content struct {
		SHA string `json:"sha"`
	} `json:"content"`
}

func (e *syncEngine) ghClient() *http.Client {
	return &http.Client{Timeout: 30 * time.Second}
}

func (e *syncEngine) ghHeaders() http.Header {
	h := http.Header{}
	h.Set("Authorization", "Bearer "+e.getToken())
	h.Set("Accept", "application/vnd.github.v3+json")
	h.Set("Content-Type", "application/json")
	return h
}

func (e *syncEngine) syncRepo() string {
	settingsMu.Lock()
	repo := curSettings.SyncRepo
	settingsMu.Unlock()
	return repo
}

func (e *syncEngine) ghContentsURL(filename string) string {
	repo := e.syncRepo()
	if filename == "" {
		return ghAPIBase + "/repos/" + repo + "/contents/"
	}
	return ghAPIBase + "/repos/" + repo + "/contents/" + url.PathEscape(filename)
}

func (e *syncEngine) verifyRepo(repo, token string) (int, error) {
	req, err := http.NewRequest(http.MethodGet, ghAPIBase+"/repos/"+repo, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	resp, err := e.ghClient().Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body) //nolint:errcheck
	return resp.StatusCode, nil
}

func (e *syncEngine) ghListNotes() ([]ghContentEntry, int, error) {
	req, err := http.NewRequest(http.MethodGet, e.ghContentsURL(""), nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header = e.ghHeaders()
	resp, err := e.ghClient().Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusNotFound {
		return []ghContentEntry{}, resp.StatusCode, nil
	}
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, resp.StatusCode, fmt.Errorf("auth rejected")
	}
	if resp.StatusCode != http.StatusOK {
		return []ghContentEntry{}, resp.StatusCode, nil
	}
	var entries []ghContentEntry
	if err := json.Unmarshal(body, &entries); err != nil {
		return []ghContentEntry{}, resp.StatusCode, nil
	}
	return entries, resp.StatusCode, nil
}

func (e *syncEngine) ghGetFile(filename string) (*ghContentFile, int, error) {
	req, err := http.NewRequest(http.MethodGet, e.ghContentsURL(filename), nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header = e.ghHeaders()
	resp, err := e.ghClient().Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusNotFound {
		return nil, resp.StatusCode, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, fmt.Errorf("github GET %d", resp.StatusCode)
	}
	var f ghContentFile
	if err := json.Unmarshal(body, &f); err != nil {
		return nil, resp.StatusCode, err
	}
	return &f, resp.StatusCode, nil
}

func ghDecodeBytes(b64 string) ([]byte, error) {
	b64 = strings.ReplaceAll(b64, "\n", "")
	return base64.StdEncoding.DecodeString(b64)
}

func ghDecodeContent(b64 string) (string, error) {
	raw, err := ghDecodeBytes(b64)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func (e *syncEngine) ghPutBytes(filename string, data []byte, sha string) (string, int, error) {
	payload := map[string]string{
		"message": "Writerdeck: " + filename,
		"content": base64.StdEncoding.EncodeToString(data),
	}
	if sha != "" {
		payload["sha"] = sha
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPut, e.ghContentsURL(filename), bytes.NewReader(body))
	if err != nil {
		return "", 0, err
	}
	req.Header = e.ghHeaders()
	resp, err := e.ghClient().Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusConflict || resp.StatusCode == 422 {
		return "", resp.StatusCode, fmt.Errorf("clash")
	}
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		e.setLastError("GitHub token rejected")
		return "", resp.StatusCode, fmt.Errorf("auth rejected")
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", resp.StatusCode, fmt.Errorf("github PUT %d", resp.StatusCode)
	}
	var out ghPutResponse
	if err := json.Unmarshal(respBody, &out); err != nil {
		return "", resp.StatusCode, err
	}
	return out.Content.SHA, resp.StatusCode, nil
}

func (e *syncEngine) ghPutFile(filename, content, sha string) (string, int, error) {
	payload := map[string]string{
		"message": "Writerdeck: " + filename,
		"content": base64.StdEncoding.EncodeToString([]byte(content)),
	}
	if sha != "" {
		payload["sha"] = sha
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPut, e.ghContentsURL(filename), bytes.NewReader(body))
	if err != nil {
		return "", 0, err
	}
	req.Header = e.ghHeaders()
	resp, err := e.ghClient().Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusConflict || resp.StatusCode == 422 {
		return "", resp.StatusCode, fmt.Errorf("clash")
	}
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		e.setLastError("GitHub token rejected")
		return "", resp.StatusCode, fmt.Errorf("auth rejected")
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", resp.StatusCode, fmt.Errorf("github PUT %d", resp.StatusCode)
	}
	var out ghPutResponse
	if err := json.Unmarshal(respBody, &out); err != nil {
		return "", resp.StatusCode, err
	}
	return out.Content.SHA, resp.StatusCode, nil
}

func (e *syncEngine) ghDeleteFile(filename, sha string) error {
	if sha == "" {
		return nil
	}
	payload, _ := json.Marshal(map[string]string{
		"message": "Writerdeck: delete " + filename,
		"sha":     sha,
	})
	req, err := http.NewRequest(http.MethodDelete, e.ghContentsURL(filename), bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header = e.ghHeaders()
	resp, err := e.ghClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body) //nolint:errcheck
	return nil
}
