package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

const (
	API_KEY             = ""
	API_KEY_SECRET      = ""
	ACCESS_TOKEN        = ""
	ACCESS_TOKEN_SECRET = ""
	URL                 = "https://us-central1-faketweet-3818e.cloudfunctions.net/generateTweetCallable"
	DATA                = "{\"data\":{\"username\":\"%s\"}}"
	USERNAME            = "blobnom"
)

var api *anaconda.TwitterApi

type TweetResult struct {
	Result struct {
		Tweet     string `json:"tweet"`
		UserImage string `json:"userImage"`
		UserAlias string `json:"userAlias"`
	} `json:"result"`
}

func main() {
	log.Println("Starting...")

	anaconda.SetConsumerKey(API_KEY)
	anaconda.SetConsumerSecret(API_KEY_SECRET)
	api = anaconda.NewTwitterApi(
		ACCESS_TOKEN,
		ACCESS_TOKEN_SECRET,
	)

	go func() {
		for {
			generateTweet()
			time.Sleep(1 * time.Hour)
		}
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(
		sc,
		syscall.SIGINT,
		syscall.SIGTERM,
		os.Interrupt,
	)
	<-sc

	log.Println("Stopping...")
	return
}

func generateTweet() {
	data := []byte(fmt.Sprintf(DATA, USERNAME))
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if strings.HasPrefix(string(body), "<") {
		return
	}

	var result TweetResult
	json.Unmarshal(body, &result)

	api.PostTweet(result.Result.Tweet, nil)
	log.Println("Posted:", result.Result.Tweet)
}
