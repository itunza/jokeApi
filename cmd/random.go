/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

// randomCmd represents the random command
var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "get a random joke",
	Long:  `get a random joke to chear you up`,
	Run: func(cmd *cobra.Command, args []string) {
		jokeTerm, _ := cmd.Flags().GetString("term")
		if jokeTerm != "" {
			getJokeDataWithTerm(jokeTerm)
		} else {
			getRandomJoke()
		}

	},
}

func init() {
	rootCmd.AddCommand(randomCmd)
	rootCmd.PersistentFlags().String("term", "", "A search term for a joke")

}

type Joke struct {
	ID     string `json:"id"`
	Joke   string `json:"joke"`
	Status int    `json:"status"`
}
type SerchResults struct {
	Results    json.RawMessage `json:"results"`
	SearchTerm string          `json:"search_term"`
	Status     int             `json:"status"`
	TotalJokes int             `json:"total_jokes"`
}

func getRandomJoke() {
	url := "https://icanhazdadjoke.com/"
	responseBytes := getJokeData(url)
	joke := Joke{}
	if err := json.Unmarshal(responseBytes, &joke); err != nil {
		log.Printf(" unmarshalling error %v", err)
	}
	fmt.Println(joke.Joke)
}

func getRandamJokeTerm(jokeTerm string) (totalJokes int, jokeList []Joke) {
	url := fmt.Sprintf("https://icanhazdadjoke.com/search?term=%s", jokeTerm)
	responseBytes := getJokeData(url)
	jokeListRaw := SerchResults{}

	if err := json.Unmarshal(responseBytes, &jokeListRaw); err != nil {
		log.Printf(" unmarshalling error %v", err)
	}
	jokes := []Joke{}
	if err := json.Unmarshal(jokeListRaw.Results, &jokes); err != nil {
		log.Printf(" unmarshalling error %v", err)
	}

	return jokeListRaw.TotalJokes, jokes

}

func randomiseJokeList(length int, jokeList []Joke) {
	rand.Seed(time.Now().Unix())
	min := 0
	max := length - 1

	if length <= 0 {
		err := fmt.Errorf("no jokes with the search term ")
		fmt.Println(err.Error())
	} else {
		randomNum := min + rand.Intn(max-min)
		fmt.Println(jokeList[randomNum].Joke)

	}
}

func getJokeData(baseURL string) []byte {
	request, err := http.NewRequest(
		http.MethodGet,
		baseURL,
		nil,
	)
	if err != nil {
		log.Printf("Error %v", err)
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("User-Agent", "Jokeapi for slackBoss group (github.com/itunza/jokeapi)")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Printf("Error %v", err)
	}
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error %v", err)
	}
	return responseBytes

}

func getJokeDataWithTerm(jokeTerm string) {
	total, results := getRandamJokeTerm(jokeTerm)
	randomiseJokeList(total, results)
}
