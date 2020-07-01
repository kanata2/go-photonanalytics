package main

import (
	"log"
	"os"

	"github.com/kanata2/go-photonanalytics"
)

func main() {
	client, err := photonanalytics.New(os.Getenv("PHOTON_TOKEN"))
	if err != nil {
		panic(err)
	}
	v, err := client.GetAppValue(&photonanalytics.GetAppValueRequest{
		AppID:    os.Getenv("PHOTON_APP_ID"),
		Region:   "jp",
		Template: "Ccu",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(v)
}
