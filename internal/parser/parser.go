package parser

import (
	"errors"
	"fmt"
	"github.com/Shemistan/parser-centra-bank-rf/internal/utils"
	"strings"
)

type Parser interface {
	Init(string, int) error
	Run() error
	Show() error
}

type parser struct {
	maxDate string
	minDate string

	maxVal     float64
	maxNameVal string

	minVal     float64
	minNameVal string

	meanRub   map[string]float64
	urls      map[string]string
	isInit    bool
	beginDate string
	endData   string
	amountDay int
}

func NewParser() Parser {
	return &parser{
		isInit: false,
	}
}

func (p *parser) Init(beginDate string, amount int) error {
	p.beginDate = beginDate
	p.amountDay = amount
	if p.isInit {
		return errors.New("the parser has already been initialized")
	}

	var err error
	p.isInit = true

	dataUrls, err := utils.GenerationUrlSlice(beginDate, amount)
	if err != nil {
		return err
	}

	// Map key is date and value is url for this date
	p.urls = dataUrls.ResMap

	p.endData = dataUrls.Date

	// The map key is the name of the currency, and the value is the ruble rate
	p.meanRub, err = utils.Parsing(p.urls[beginDate])
	if err != nil {
		return err
	}

	// Structure with maps, 1st map - min / max values, 2nd map - names of these currencies
	minMax, err := utils.SearchMinMax(p.meanRub)
	if err != nil {
		return err
	}

	name := minMax.Name
	val := minMax.Value

	p.maxDate = p.beginDate
	p.minDate = p.beginDate

	p.maxNameVal = name["max"]
	p.minNameVal = name["min"]

	p.maxVal = val["max"]
	p.minVal = val["min"]

	return nil
}

func (p *parser) Run() error {
	if !p.isInit {
		return errors.New("error parser is not initialized")
	}

	for date, url := range p.urls {
		if date == p.beginDate {
			continue
		}

		ratePerDay, err := utils.Parsing(url)
		if err != nil {
			return err
		}

		for k, v := range ratePerDay {
			if _, ok := p.meanRub[k]; ok == true {
				p.meanRub[k] += v
			} else {
				p.meanRub[k] = v
			}

			if v > p.maxVal {
				p.maxVal = v
				p.maxNameVal = k
				p.maxDate = date
			}

			if v < p.minVal {
				p.minVal = v
				p.minNameVal = k
				p.minDate = date
			}
		}
	}

	for k, v := range p.meanRub {
		p.meanRub[k] = v / float64(p.amountDay)
	}

	return nil
}

func (p *parser) Show() error {
	begin := strings.ReplaceAll(p.beginDate, "/", ".")
	end := strings.ReplaceAll(p.endData, "/", ".")
	fmt.Printf("Analysis result for the period from %s to %s:\n", begin, end)

	p.maxDate = strings.ReplaceAll(p.maxDate, "/", ".")
	fmt.Printf("The maximum value of the ruble:\ndate: %s\ncurrency name: %s\nvalue: %v\n\n",
		p.maxDate, p.maxNameVal, p.maxVal)

	p.minDate = strings.ReplaceAll(p.minDate, "/", ".")
	fmt.Printf("The minimum value of the ruble:\ndate: %s\ncurrency name: %s\nvalue: %v\n\n",
		p.minDate, p.minNameVal, p.minVal)

	fmt.Println("Average value of the ruble exchange rate for the entire period:")
	for k, v := range p.meanRub {
		fmt.Printf("%s  :  %v\n", k, v)
	}

	return nil
}
