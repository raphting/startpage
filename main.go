package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

var globalPage string
var APIKey string

func main() {
	APIKey = os.Getenv("STARTPAGE_API_KEY")
	if APIKey == "" {
		fmt.Println("The environment variable STARTPAGE_API_KEY (Unsplash API access key) is required.")
		return
	}

	generateFront()
	go runLoop()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, globalPage)
	})
	log.Fatal(http.ListenAndServe(":9691", nil))
}

func runLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	for {
		select {
		case <-ticker.C:
			generateFront()
		}
	}
}

func generateFront() {
	fmt.Println("Generate and publish new front page")
	uname, url, _ := callUnsplash(APIKey)
	parseFrontToGlobal(uname, url)
}

func parseFrontToGlobal(username, imageurl string) error {
	tpl, err := template.ParseFiles("./front.html")
	if err != nil {
		return err
	}

	var b bytes.Buffer
	foo := bufio.NewWriter(&b)

	type front struct {
		Username string
		ImageURL string
	}

	err = tpl.Execute(foo, front{username, imageurl})
	if err != nil {
		return err
	}
	err = foo.Flush()
	if err != nil {
		return err
	}
	globalPage = b.String()
	return nil
}

func callUnsplash(api_key string) (string, string, error) {
	// Prepare call to API
	uapi, _ := url.Parse("https://api.unsplash.com/photos/random?orientation=landscape")
	authHeader := http.Header{}
	authHeader.Add("Authorization", "Client-ID "+api_key)
	req := http.Request{
		Method: "GET",
		URL:    uapi,
		Header: authHeader,
	}

	// Call API
	c := http.Client{
		Timeout: time.Second * 30,
	}
	resp, err := c.Do(&req)
	if err != nil {
		return "", "", err
	}

	type randomPhoto struct {
		RNDPhotoURLs struct {
			Regular string `json:"regular"`
		} `json:"urls"`
		User struct {
			Username string `json:"username"`
		}
	}

	rp := randomPhoto{}

	ubody, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(ubody, &rp)
	return rp.User.Username, rp.RNDPhotoURLs.Regular, err
}
