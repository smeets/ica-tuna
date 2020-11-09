package main

import (
    "math"
    "regexp"
    "strconv"
    "sort"
    "time"
    "fmt"
    "errors"
)

type deal struct {
    Seen time.Time `json:"seen"`
    Unit float64 `json:"unit_price"`
    Total float64 `json:"total_price"`
}

// format is: X/Y?st
// X = total price
// Y = minimum but
var pricematch *regexp.Regexp = regexp.MustCompile(`(\d+)/(\d+)?(st|frp|kg|hg)`)
var errInvalidFormat error = errors.New("Invalid offer format")

func NewDeal(it *item) (deal, error) {
    dealwithit := pricematch.FindStringSubmatch(it.Offer)
    if len(dealwithit) < 2 {
        // match error?
        return deal{}, fmt.Errorf("%w: %s", errInvalidFormat, it.Offer)
    }

    total, err := strconv.ParseFloat(dealwithit[1], 64)
    if err != nil {
        return deal{}, fmt.Errorf("total: %w (%v)", err, dealwithit)
    }

    var units float64 = 1
    if len(dealwithit[2]) > 0 {
        units, err = strconv.ParseFloat(dealwithit[2], 64)
        if err != nil {
            return deal{}, fmt.Errorf("units: %w (%v)", err, dealwithit)
        }
    }

    // diy round to 0.01 ...
    unit := yolo(total / units)
    return deal{
        Seen: it.Seen,
        Unit: unit,
        Total: total,
    }, nil
}

type history map[string][]deal

func (h history) record(it *item) {
    deals := h[it.Label]
    decoded, err := NewDeal(it)
    if err != nil {
        return
    }
    deals = append(deals, decoded)
    h[it.Label] = deals
}

func (h history) compare(it *item) (float64, float64, float64, float64, bool) {
    deals := h[it.Label]
    if len(deals) == 0 {
        return 0,0,0,0,false
    }

    sort.Slice(deals[:], func(i, j int) bool {
        return deals[i].Unit < deals[j].Unit
    })
    medianUnitPrice := deals[len(deals)/2].Unit

    sort.Slice(deals[:], func(i, j int) bool {
        return deals[i].Total < deals[j].Total
    })
    medianTotalPrice := deals[len(deals)/2].Total

    var avgUnitPrice float64
    var avgTotalPrice float64
    for _, known := range deals {
        avgUnitPrice += known.Unit
        avgTotalPrice += known.Total
    }
    avgUnitPrice /= float64(len(deals))
    avgTotalPrice /= float64(len(deals))

    avgUnitPrice = yolo(avgUnitPrice)
    avgTotalPrice = yolo(avgTotalPrice)

    return medianUnitPrice, avgUnitPrice, medianTotalPrice, avgTotalPrice, true
}

func loadhistory() (history, error) {
    items, err := loaditems()
    if err != nil {
        return nil, err
    }

    offers := make(history)
    for _, item := range items {
        offers.record(&item)
    }

    return offers, nil
}

func yolo(f float64) float64 {
    return f - (f - math.Round(f / 0.01) * 0.01)
}
