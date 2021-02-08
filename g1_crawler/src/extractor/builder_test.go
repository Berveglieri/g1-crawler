package extractor

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestBuilder(t *testing.T) {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS","/home/felipe/go/src/crawler/src/creds/auth.json")
	fmt.Println(Extract())
	log.Println("Finished")
}
