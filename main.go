package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

type MethodProfile struct {
	MethodType string `json:"methodType"`
	RandomFailure string `json:"randomFailure"`
	FakeReturn *json.RawMessage `json:"fakeReturn"`
}

type FakeProfileNameSpace struct {
	Namespace string `json:"namespace"`
	Routes []map[string]MethodProfile `json:"routes"`
}

var (
	profile string
	addr string
)

func main() {
	flag.StringVar(&profile, "profile", "", "help message for flagname")
	flag.StringVar(&addr, "addr", "", "help message for flagname")
	flag.Parse()

	err := parseProfile(profile)
	if err != nil {
		fmt.Println(err)
	}
}

func parseProfile(profile string) error {
	raw, err := ioutil.ReadFile(profile)
	if err != nil {
		return err
	}

	var fakeProfileNameSpaces []FakeProfileNameSpace
	if err := json.Unmarshal(raw, &fakeProfileNameSpaces); err != nil {
		return err
	}
	
	r := gin.Default()

	for _, fakeProfileNameSpace := range fakeProfileNameSpaces {
		group := r.Group(fakeProfileNameSpace.Namespace)
		{
			for _, method := range fakeProfileNameSpace.Routes {
				for methodName, methodProfile := range method {
					group.Handle(methodProfile.MethodType, methodName, func(c *gin.Context) {
						indicator := 0

						if methodProfile.RandomFailure == "true" {
							rand.Seed(time.Now().UnixNano())
							indicator = rand.Intn(2)
						}

						if indicator == 0 {
							if methodProfile.FakeReturn != nil {
								c.JSON(200, *methodProfile.FakeReturn)
							} else {
								c.Status(http.StatusOK)
							}
						} else {
							c.Status(http.StatusNotFound)
						}
					})
				}
			}
		}
	}
	
	r.Run(addr)

	return nil
}

//package main
//
//import (
//"fmt"
//"github.com/xeipuuv/gojsonschema"
//)
//
//func main() {
//
//	schemaLoader := gojsonschema.NewReferenceLoader("file:///home/me/schema.json")
//	documentLoader := gojsonschema.NewReferenceLoader("file:///home/me/document.json")
//
//	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
//	if err != nil {
//		panic(err.Error())
//	}
//
//	if result.Valid() {
//		fmt.Printf("The document is valid\n")
//	} else {
//		fmt.Printf("The document is not valid. see errors :\n")
//		for _, desc := range result.Errors() {
//			fmt.Printf("- %s\n", desc)
//		}
//	}
//}