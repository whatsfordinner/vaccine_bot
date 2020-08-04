package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type diseaseList struct {
	Diseases []string `json:"diseases"`
}

type twitterConfig struct {
	consumerKey    string
	consumerSecret string
	accessToken    string
	accessSecret   string
}

func getTwitterConfig() (*twitterConfig, error) {
	newTwitterConfig := new(twitterConfig)
	errorString := "Required envrionment variable %s is not defined"

	newConsumerKey, exists := os.LookupEnv("TWITTER_CONSUMER_KEY")

	if !exists {
		return nil, fmt.Errorf(errorString, "TWITTER_CONSUMER_KEY")
	}

	newConsumerSecret, exists := os.LookupEnv("TWITTER_CONSUMER_SECRET")

	if !exists {
		return nil, fmt.Errorf(errorString, "TWITTER_CONSUMER_SECRET")
	}

	newAccessToken, exists := os.LookupEnv("TWITTER_ACCESS_TOKEN")

	if !exists {
		return nil, fmt.Errorf(errorString, "TWITTER_ACCESS_TOKEN")
	}

	newAccessSecret, exists := os.LookupEnv("TWITTER_ACCESS_SECRET")

	if !exists {
		return nil, fmt.Errorf(errorString, "TWITTER_ACCESS_SECRET")
	}

	newTwitterConfig.consumerKey = newConsumerKey
	newTwitterConfig.consumerSecret = newConsumerSecret
	newTwitterConfig.accessToken = newAccessToken
	newTwitterConfig.accessSecret = newAccessSecret

	return newTwitterConfig, nil
}

func getDiseaseListFromFile(diseaseFilename string) (*diseaseList, error) {
	log.Printf("Creating new disease list from file: %s", diseaseFilename)
	diseases := new(diseaseList)
	diseaseFile, err := os.Open(diseaseFilename)

	if err != nil {
		return nil, err
	}

	defer diseaseFile.Close()
	byteValue, err := ioutil.ReadAll(diseaseFile)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(byteValue), diseases)

	if err != nil {
		return nil, err
	}

	if len(diseases.Diseases) == 0 {
		errorString := fmt.Sprintf("Source JSON %s is malformed or has no diseases", diseaseFilename)
		return nil, errors.New(errorString)
	}

	return diseases, nil
}

func (dl *diseaseList) getDisease() (string, error) {
	log.Print("Getting random disease from list of diseases")
	rand.Seed(time.Now().UnixNano())

	if len(dl.Diseases) == 0 {
		return "", errors.New("List of diseases has no entries")
	}

	return dl.Diseases[rand.Intn(len(dl.Diseases))], nil
}

func buildTweet(disease string) string {
	log.Printf("Generating new tweet with disease: %s", disease)
	return fmt.Sprintf("💉 Get vaccinated against %s", disease)
}

func sendTweet(newTweet string) error {
	twitterConfig, err := getTwitterConfig()

	if err != nil {
		return err
	}

	log.Print("Creating new Twitter client")
	config := oauth1.NewConfig(twitterConfig.consumerKey, twitterConfig.consumerSecret)
	token := oauth1.NewToken(twitterConfig.accessToken, twitterConfig.accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	twitterClient := twitter.NewClient(httpClient)

	log.Print("Sending tweet")
	tweet, resp, err := twitterClient.Statuses.Update(newTweet, nil)

	if err != nil {
		return err
	}

	log.Printf("Received response:\n%s\n%s", prettyPrint(tweet), prettyPrint(resp))

	return nil
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "    ")
	return string(s)
}

func handleRequest(ctx context.Context, req events.CloudWatchEvent) {
	log.Printf("Handling new request")
	dl, err := getDiseaseListFromFile("diseases.json")

	if err != nil {
		log.Panic(err.Error())
	}

	disease, err := dl.getDisease()

	if err != nil {
		log.Panic(err.Error())
	}

	tweet := buildTweet(disease)

	log.Printf("Created new tweet: %s", tweet)

	err = sendTweet(tweet)

	if err != nil {
		panic(err.Error())
	}
}

func main() {
	log.Printf("New execution context created")
	lambda.Start(handleRequest)
}
