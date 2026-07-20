package main

import "testing"

func TestIdeBrowserUA(t *testing.T) {
	cursor := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Cursor/3.12.17 Chrome/144.0.7559.236 Electron/40.10.3 Safari/537.36"
	safari := "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1"
	chrome := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"
	if !ideBrowserUA(cursor) {
		t.Fatal("Cursor UA should be ide browser")
	}
	if !ideBrowserUA("Electron/40.0.0") {
		t.Fatal("Electron UA should be ide browser")
	}
	if ideBrowserUA(safari) {
		t.Fatal("Safari should not be ide browser")
	}
	if ideBrowserUA(chrome) {
		t.Fatal("Chrome should not be ide browser")
	}
	if ideBrowserUA("") {
		t.Fatal("empty UA should not be ide browser")
	}
}
