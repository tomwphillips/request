package requester

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type instruction struct {
	URL string
}

func decodeInstruction(m []byte) instruction {
	var i instruction
	err := json.Unmarshal(m, &i)
	if err != nil {
		log.Fatalf("Unmarshal failed: %v", err)
	}
	return i
}

func getURL(url *string) string {
	resp, err := http.Get(*url)
	if err != nil {
		log.Fatalf("Get request failed: %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Read response body failed: %v", err)
	}
	return string(body)
}

func executeInstruction(i instruction) string {
	return getURL(&i.URL)
}
