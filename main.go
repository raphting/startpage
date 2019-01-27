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

func main() {
	API_KEY := os.Getenv("STARTPAGE_API_KEY")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		uname, url, _ := callUnsplash(API_KEY)
		front, _ := parseFront(uname, url)
		fmt.Fprintf(w, front)
	})
	log.Fatal(http.ListenAndServe(":9691", nil))
}

func parseFront(username, imageurl string) (string, error) {
	tpl, err := template.ParseFiles("./front.html")
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	foo := bufio.NewWriter(&b)

	type front struct {
		Username string
		ImageURL string
	}

	err = tpl.Execute(foo, front{username, imageurl})
	if err != nil {
		return "", err
	}
	err = foo.Flush()
	if err != nil {
		return "", err
	}

	return b.String(), err
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
