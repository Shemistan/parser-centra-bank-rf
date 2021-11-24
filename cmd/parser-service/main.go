package main

import (
	"fmt"
	"github.com/Shemistan/parser-centra-bank-rf/internal/utils"
	"log"
)

func main() {
	url := "https://www.cbr.ru/scripts/XML_daily_eng.asp?date_req=11/11/2020)"

	makeReq, errReq := utils.MakeRequest(url)
	if errReq != nil {
		log.Printf("error request: %s", errReq.Error())
	}

	editSlice, err := utils.EditingSlice(makeReq)
	if err != nil {
		log.Printf("error editing slice: %s", err.Error())
	}

	val, name, errSearch := utils.SearchMinMax(editSlice)
	if errSearch != nil {
		log.Printf("error searching min/max value: %s", errSearch.Error())
	}

	fmt.Printf("data: %v", name["day"])
	fmt.Printf("max value %v - %v", name["max"], val["max"])
	fmt.Printf("min value %v - %v", name["min"], val["min"])

}
