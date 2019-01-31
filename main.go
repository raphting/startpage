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

	// Generate front page once and then start ticker to generate every n minutes
	generateFront()
	go runLoop()

	// Handle and serve front page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, globalPage)
		if err != nil {
			fmt.Println(err.Error())
		}
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
	resp, err := callUnsplash(APIKey)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = parseFrontToGlobal(resp.Name, resp.Username, resp.PhotoURLRegular)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func parseFrontToGlobal(name, username, imageurl string) error {
	tpl := template.New("front")
	tpl, err := tpl.Parse(getFrontTemplate())
	if err != nil {
		return err
	}

	var b bytes.Buffer
	foo := bufio.NewWriter(&b)

	type front struct {
		Name     string
		Username string
		ImageURL string
	}

	err = tpl.Execute(foo, front{name, username, imageurl})
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

type apiResponse struct {
	Name            string
	Username        string
	PhotoURLRegular string
}

func callUnsplash(api_key string) (apiResponse, error) {
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
		return apiResponse{}, err
	}

	// Prepare to pick what's needed from returned JSON
	type randomPhoto struct {
		RNDPhotoURLs struct {
			Regular string `json:"regular"`
		} `json:"urls"`
		User struct {
			Name     string `json:"name"`
			Username string `json:"username"`
		}
	}

	rp := randomPhoto{}
	ubody, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(ubody, &rp)

	maxNameLen := 30
	if len(rp.User.Name) > maxNameLen {
		rp.User.Name = rp.User.Name[0:maxNameLen]
		rp.User.Name += "..."
	}

	return apiResponse{rp.User.Name,
		rp.User.Username,
		rp.RNDPhotoURLs.Regular}, err
}
