package main

import "net/http"
import "io"
import "log"

func readHtml(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://www.ica.se/butiker/kvantum/lund/ica-kvantum-malmborgs-tuna-2780/erbjudanden/")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	defer resp.Body.Close()
	io.Copy(w, resp.Body)
	//w.Write(resp.Body)
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.HandleFunc("/html", readHtml)
	log.Println("Listening @ 5000...")
	log.Fatal(http.ListenAndServe(":5000", nil))
}