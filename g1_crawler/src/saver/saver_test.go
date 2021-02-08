package saver

import (
	"log"
	"os"
	"testing"
)

func TestSaver(t *testing.T)  {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS","/home/felipe/go/src/crawler/src/creds/auth.json")
	UploadScrapeToBucket("globo_g1")
	log.Println("Finished")
}
