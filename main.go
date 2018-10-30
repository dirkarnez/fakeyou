package main

import (
	"encoding/json"
	"flag"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

type MethodProfile struct {
	MethodType string `json:"methodType"`
	RandomFailure string `json:"randomFailure"`
	FakeReturn json.RawMessage `json:"fakeReturn"`
}

var (
	profile string
	addr string
)

func main() {
	flag.StringVar(&profile, "profile", "", "help message for flagname")
	flag.StringVar(&addr, "addr", "", "help message for flagname")
	flag.Parse()

	parseProfile(profile)
}

func parseProfile(profile string) error {
	raw, err := ioutil.ReadFile(profile)
	if err != nil {
		return err
	}
	var dat []map[string]MethodProfile

	if err := json.Unmarshal(raw, &dat); err != nil {
		return err
	}

	r := gin.Default()
	for _, method := range dat {
		for methodName, methodProfile := range method {
			r.Handle(methodProfile.MethodType, methodName, func(c *gin.Context) {
				indicator := 0

				if methodProfile.RandomFailure == "true" {
					rand.Seed(time.Now().UnixNano())
					indicator = rand.Intn(2)
				}

				if indicator == 0 {
					c.JSON(200, methodProfile.FakeReturn)
				} else {
					c.Status(http.StatusNotFound)
				}
			})

		}
	}
	r.Run(":7890")

	return nil
}