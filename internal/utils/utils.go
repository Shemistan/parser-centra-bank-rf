package utils

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func MakeRequest(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	spliting := strings.Split(string(body), "<Valute ID=")

	return spliting, nil
}

func EditingSlice(data []string) (map[string]float64, error) {
	resOneDay := make(map[string]float64)
	begin := strings.Index(data[0], "Date=") + 6
	end := strings.Index(data[0], " name") - 1
	dateToday := data[0][begin:end]

	for _, v := range data[1:] {
		beginKey := strings.Index(v, "<Name>") + 6
		endKey := strings.Index(v, "</Name>")
		key := v[beginKey:endKey]

		beginVal := strings.Index(v, "<Value>") + 7
		endVal := strings.Index(v, "</Value>")

		val, errFloat := strconv.ParseFloat(strings.ReplaceAll(v[beginVal:endVal], ",", "."), 64)
		if errFloat != nil {
			return nil, errFloat
		}

		resOneDay[key] = val
	}

	resOneDay[dateToday] = -1

	return resOneDay, nil
}

func SearchMinMax(data map[string]float64) (map[string]float64, map[string]string, error) {

	maxVal := float64(0)
	maxName := ""

	minVal := data["British Pound Sterling"]
	minName := ""
	day := ""

	for k, v := range data {
		if v == -1 {
			day = k
			continue
		}

		if v > maxVal {
			maxVal = v
			maxName = k
		}

		if v < minVal {
			minVal = v
			minName = k
		}
	}
	val := map[string]float64{
		"max": maxVal,
		"min": minVal,
	}
	name := map[string]string{
		"max": maxName,
		"min": minName,
		"day": day,
	}

	return val, name, nil

}
