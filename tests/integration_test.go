package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestIntegration(t *testing.T) {
	// dummyLogin as moderator
	type Tok struct{ Token string }
	buf := bytes.NewBufferString(`{"role":"moderator"}`)
	res, _ := http.Post("http://localhost:8080/dummyLogin", "application/json", buf)
	var tok Tok
	json.NewDecoder(res.Body).Decode(&tok)
	// create PVZ
	buf = bytes.NewBufferString(`{"city":"Москва"}`)
	req, _ := http.NewRequest("POST", "http://localhost:8080/pvz", buf)
	req.Header.Set("Authorization", "Bearer "+tok.Token)
	res, _ = http.DefaultClient.Do(req)
	var pvz struct{ Id int64 }
	json.NewDecoder(res.Body).Decode(&pvz)
	// start session
	req, _ = http.NewRequest("POST", fmt.Sprintf("http://localhost:8080/pvz/%d/sessions", pvz.Id), nil)
	req.Header.Set("Authorization", "Bearer "+tok.Token)
	res, _ = http.DefaultClient.Do(req)
	var sess struct{ Id int64 }
	json.NewDecoder(res.Body).Decode(&sess)
	// add 50 items
	for i := 0; i < 50; i++ {
		buf = bytes.NewBufferString(`{"type":"electronics"}`)
		req, _ = http.NewRequest("POST", fmt.Sprintf("http://localhost:8080/pvz/%d/sessions/%d/items", pvz.Id, sess.Id), buf)
		req.Header.Set("Authorization", "Bearer "+tok.Token)
		http.DefaultClient.Do(req)
	}
	// close session
	req, _ = http.NewRequest("POST", fmt.Sprintf("http://localhost:8080/pvz/%d/sessions/%d/close", pvz.Id, sess.Id), nil)
	req.Header.Set("Authorization", "Bearer "+tok.Token)
	res, _ = http.DefaultClient.Do(req)
	if res.StatusCode != http.StatusOK {
		t.Fatal("expected 200")
	}
}
