package main

import "net/http"
import "io"
import "encoding/json"
import "log"
import "time"
import "fmt"
import "io/ioutil"
import "flag"

var datastore []byte
var lastfetch time.Time
var pricefilepath string
var runanalysis bool

func init() {
	flag.StringVar(&pricefilepath, "pricedb", "./pricedb.tsv", "path to price database")
	flag.BoolVar(&runanalysis, "report", false, "run analysis")
}

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

	if time.Since(lastfetch) > time.Hour * 24 * 3 ||
		(int(lastfetch.Weekday()) > int(time.Now().Weekday())) ||
		(lastfetch.Weekday() == time.Sunday && time.Now().Weekday() == time.Monday) {
		if err := update(); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	w.Write(datastore)
}

func recordPrices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	items, err := getitems(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := dumpitems(items); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`OK`))
}

func comparePrices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	items, err := getitems(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hist, err := loadhistory()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c := struct {
		Items []item `json:"items"`
	} { }
	for _, it := range items {
		thisdeal, err := NewDeal(&it)
		if err != nil {
			log.Println(it.Label, err)
			continue
		}

		unitMed, unitAvg, totalMed, totalAvg, seenBefore := hist.compare(&it)

		var decision string
		if !seenBefore {
			decision = "item not yet tracked"
		} else if thisdeal.Unit < unitMed && thisdeal.Total < totalMed &&
			thisdeal.Unit < unitAvg && thisdeal.Total < totalAvg {
			decision = "outstanding deal"
		} else if thisdeal.Unit < unitMed && thisdeal.Total < totalMed {
			decision = "better than median"
		} else if thisdeal.Unit < unitMed {
			decision = "unit price better than median"
		} else if thisdeal.Unit < unitAvg && thisdeal.Total < totalAvg {
			decision = "better than average"
		} else if thisdeal.Unit < unitAvg {
			decision = "unit price better than average"
		} else if thisdeal.Unit == unitAvg || thisdeal.Total == totalAvg ||
			thisdeal.Unit == unitMed || thisdeal.Total == totalMed {
			decision = "normal deal"
		} else {
			decision = "probably not a good deal"
		}

		// fmt.Println(it.Label, decision)

		c.Items = append(c.Items, item{
			Seen: it.Seen,
			Label: it.Label,
			Offer: decision,
		})
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(c)
}

func proxyGet(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	url := qs.Get("q")

	res, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	w.Header().Set("Content-Type", res.Header.Get("Content-Type"))
	w.WriteHeader(res.StatusCode)
    io.Copy(w, res.Body)
}

func main() {
	flag.Parse()
	if runanalysis {
		hist, err := loadhistory()
		if err != nil {
			log.Fatalln(err)
		}
		for k, v := range hist {
			fmt.Println(k, v)
		}
		return
	}

	update()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.HandleFunc("/html", readHtml)
	http.HandleFunc("/record", recordPrices)
	http.HandleFunc("/compare", comparePrices)
	http.HandleFunc("/proxy", proxyGet)
	log.Println("Listening @ 5000...")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
