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
)

type diseaseList struct {
	Diseases []string `json:"diseases"`
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
	return fmt.Sprintf("Vaccinate your kids against %s", disease)
}

func sendTweet() {

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
}

func main() {
	log.Printf("New execution context created")
	lambda.Start(handleRequest)
}
