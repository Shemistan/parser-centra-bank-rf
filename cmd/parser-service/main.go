package main

import (
	"fmt"
	"github.com/Shemistan/parser-centra-bank-rf/internal/parser"
	"log"
	"os"
)

func main() {
	var date string
	var amount int

	fmt.Print("Input the start of the analysis period in day/month/year format(\"23/08/2018\"):\n--->")
	_, err := fmt.Fscan(os.Stdin, &date)
	if err != nil {
		return
	}

	fmt.Print("Input the required number of days (90):\n--->")
	_, err = fmt.Fscan(os.Stdin, &amount)
	if err != nil {
		return
	}

	exmPars := parser.NewParser()
	err = exmPars.Init(date, amount)
	if err != nil {
		log.Printf("error init: %s", err.Error())
	}

	err = exmPars.Run()
	if err != nil {
		log.Printf("error run: %s", err.Error())
	}
	fmt.Print("\n \n \n \n \n")
	err = exmPars.Show()
	if err != nil {
		log.Printf("error show: %s", err.Error())
	}
}
