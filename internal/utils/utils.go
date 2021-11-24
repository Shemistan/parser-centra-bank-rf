package utils

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	beginName = "<Name>"
	endName   = "</Name>"
	beginVal  = "<Value>"
	endVal    = "</Value>"
)

type MinMax struct {
	Value map[string]float64
	Name  map[string]string
}

type DateUrl struct {
	Date   string
	ResMap map[string]string
}

func makeRequest(url string) ([]string, error) {
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

func Parsing(url string) (map[string]float64, error) {
	ratePerDay := make(map[string]float64)

	resp, err := makeRequest(url)
	if err != nil {
		return nil, err
	}

	for _, v := range resp[1:] {
		begin := strings.Index(v, beginName) + len(beginName)
		end := strings.Index(v, endName)
		currencyName := v[begin:end]

		begin = strings.Index(v, beginVal) + len(beginVal)
		end = strings.Index(v, endVal)

		rubRate, errFloat := strconv.ParseFloat(strings.ReplaceAll(v[begin:end], ",", "."), 64)
		if errFloat != nil {
			return nil, errFloat
		}

		ratePerDay[currencyName] = rubRate
	}

	return ratePerDay, nil
}

func SearchMinMax(data map[string]float64) (*MinMax, error) {

	maxVal := float64(0)
	maxName := ""

	minVal := data["British Pound Sterling"]
	minName := "British Pound Sterling"

	for k, v := range data {

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
	}

	return &MinMax{
		Value: val,
		Name:  name,
	}, nil

}

func GenerationUrlSlice(beginDate string, amountDay int) (*DateUrl, error) {
	var dayMonthYear []int

	dateSlice := strings.Split(beginDate, "/")

	for _, v := range dateSlice {
		val, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}

		dayMonthYear = append(dayMonthYear, int(val))
	}

	urlSlice, err := generationUrlSlice(dayMonthYear[0], dayMonthYear[1], dayMonthYear[2], amountDay)
	if err != nil {
		return nil, err
	}

	return urlSlice, nil
}

func generationUrlSlice(d, m, y, amountDay int) (*DateUrl, error) {
	date := ""
	resMap := make(map[string]string)

	url := "https://www.cbr.ru/scripts/XML_daily_eng.asp?date_req="
	leapYear := definitionLeapYear(y)

	for i := 0; i <= amountDay; {
		numbersDayInMonth := definitionMonth(m, leapYear)

		strM := strconv.Itoa(m)
		if m < 10 {
			strM = "0" + strM
		}

		strY := strconv.Itoa(y)

		for d <= numbersDayInMonth {

			controlDate := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)

			if controlDate.After(time.Now()) {
				return &DateUrl{
					Date:   date,
					ResMap: resMap,
				}, nil
			}

			strD := strconv.Itoa(d)
			if d < 10 {
				strD = "0" + strD
			}
			date = strD + "/" + strM + "/" + strY

			if _, ok := resMap[date]; ok == false {
				resMap[date] = url + date
			}

			i++
			d++
			if i >= amountDay {
				return &DateUrl{
					Date:   date,
					ResMap: resMap,
				}, nil
			}
		}
		d = 1
		m++

		if m > 12 {
			m = 1
			y++
			leapYear = definitionLeapYear(y)
		}
	}
	return &DateUrl{
		Date:   date,
		ResMap: resMap,
	}, nil
}

func definitionMonth(m int, isLeap bool) int {
	switch {
	case m == 4 || m == 6 || m == 9 || m == 11:
		return 30
	case m == 2:
		if isLeap {
			return 29
		}
		return 28
	default:
		return 31
	}

}

func definitionLeapYear(y int) bool {
	if y%4 == 0 {
		if y%100 == 0 {
			if y%400 == 0 {
				return true
			}

			return false
		}

		return true
	}

	return true
}
