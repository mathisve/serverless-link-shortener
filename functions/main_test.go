package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"testing"
)

var (
	// example url
	url = "https://3hljbavp0k.execute-api.eu-central-1.amazonaws.com"
	getURL = url + "/get-link/%s"
	newURL = url + "/new-link"
)

type request struct {
	Link string `json:"link"`
}

func getRandomString() string {
	return uuid.NewString()[:8]
}

func getRandomBool() bool {
	if rand.Intn(1) == 0 {
		return true
	}

	return false
}

func generateLink() string {
	return fmt.Sprintf("https://www.%s.com", getRandomString())
}

func makeRequest(url string, payload io.Reader) (*http.Response, string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return nil, "", err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return res, "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return res, "", err
	}

	return res, string(body), err
}

func hash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func TestLinks(t *testing.T) {

	for i:=0; i <= 10; i+=1 {
		link := generateLink()
		h := hash(link)[:5]

		makeFirst := getRandomBool()
		if makeFirst {
			payload := strings.NewReader(fmt.Sprintf(`{"link": "%s"}`, link))

			res, body, err := makeRequest(newURL, payload)
			if err != nil {
				t.Error(err)
			}

			if res.StatusCode != 200 {
				t.Error("status code not 200\n")
			}

			if body != h {
				t.Error(fmt.Printf("returned hash not the same as self produced hash, expected: %s, recieved: %s\n", h, body))
			}
		}

		res, body, err := makeRequest(fmt.Sprintf(getURL, h), nil)
		if err != nil {
			t.Error(err)
		}

		// made the link, expect 200 status code
		if makeFirst && res.StatusCode != 200 {
			t.Error(fmt.Sprintf("expected status code: 200, received: %d\n", res.StatusCode))
		} else if !makeFirst && res.StatusCode != 404 {
			t.Error(fmt.Sprintf("expected status code: 404, received: %d\n", res.StatusCode))
		}


		if makeFirst && !strings.Contains(body, link) {
			t.Error(fmt.Sprintf("body does not contain link. body: %s, link: %s\n", body, link))
		}
	}
}
