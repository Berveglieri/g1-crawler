package crawler

import (
	"fmt"
	"github.com/crawler/src/saver"
	"net/http"
	"time"
)

func Trigger(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	startFormat := start.Format("2 Jan 2006 15:04:05")
	fmt.Printf("triggering scraper script\n")
	saver.UploadScrapeToBucket("globo_g1")
	finished := time.Now()
	finishedFormat := finished.Format("2 Jan 2006 15:04:05")
	fmt.Printf("started at %s\n",startFormat)
	fmt.Printf("scraper finished at %s",finishedFormat)
}
