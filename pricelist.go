package main

import "encoding/json"
import "encoding/csv"
import "net/http"
import "os"
import "io"
import "log"
import "time"

type item struct {
    Seen time.Time `json:"seen"`
    Label string `json:"label"`
    Offer string `json:"offer"`
}

// getitems decodes json-encoded items from request body
func getitems(r *http.Request) ([]item, error) {
    var items []item
    if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
        return nil, err
    }
    return items, nil
}

// dumpitems appends items to file as tab-separated values
func dumpitems(items []item) error {
    f, err := os.OpenFile(pricefilepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }

    // weekly offers
    t := time.Now().Truncate(time.Hour * 24 * 7)
    tsbytes, _ := t.MarshalText()
    ts := string(tsbytes)

    tsv := csv.NewWriter(f)
    tsv.Comma = '\t'
    tsv.UseCRLF = false

    for _, item := range items {
        if err := tsv.Write([]string{ts, item.Label, item.Offer}); err != nil {
            f.Close()
            return err
        }
    }

    tsv.Flush()

    if err := f.Close(); err != nil {
        return err
    }

    return nil
}

// loaditems deserializes tab-separated items from the db.
func loaditems() ([]item, error) {
    f, err := os.Open(pricefilepath)
    if err != nil {
        return nil, err
    }

    tsv := csv.NewReader(f)
    tsv.Comma = '\t'
    tsv.ReuseRecord = true

    items := make([]item, 0, 32)

    for {
        vals, err := tsv.Read()
        if err != nil && err != io.EOF {
            f.Close()
            return nil, err
        }

        if len(vals) == 0 && err == io.EOF {
            break
        }

        var t time.Time
        if err := t.UnmarshalText([]byte(vals[0])); err != nil {
            log.Println(vals[0], err)
            continue
        }
        items = append(items, item{
            Seen: t,
            Label: vals[1],
            Offer: vals[2],
        })
    }

    if err := f.Close(); err != nil {
        return nil, err
    }

    return items, nil
}
