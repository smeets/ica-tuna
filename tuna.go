package main

import "net/http"
import "log"
import "time"
import "io/ioutil"

var datastore []byte
var lastfetch time.Time

func update() error {
	lastfetch = time.Now()

	resp, err := http.Get("https://www.ica.se/butiker/kvantum/lund/ica-kvantum-malmborgs-tuna-2780/erbjudanden/")
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	datastore = body
	return nil
}

func readHtml(w http.ResponseWriter, r *http.Request) {
	if time.Since(lastfetch) > time.Hour {
		if err := update(); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	w.Write(datastore)
}

func main() {
	update()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.HandleFunc("/html", readHtml)
	log.Println("Listening @ 5000...")
	log.Fatal(http.ListenAndServe(":5000", nil))
}