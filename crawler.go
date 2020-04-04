package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func getStockMap() []StockInfo {
	// stockMap := make(map[string]StockInfo)
	var stockList []StockInfo

	c := colly.NewCollector(
		colly.DetectCharset(),
	)

	c.OnRequest(func(r *colly.Request) {
		r.ResponseCharacterEncoding = "big5"
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	// On every a element which has href attribute call callback
	c.OnHTML("table.h4", func(table *colly.HTMLElement) {
		table.DOM.Find("tbody>tr").Each(func(i int, tr *goquery.Selection) {
			if i > 1 && tr.Find("td").Length() == 7 {
				info := &StockInfo{}
				var stockID string
				tr.Find("td").Each(func(i int, td *goquery.Selection) {
					text := td.Text()
					switch i {
					case 0:
						textArr := strings.Split(text, "　")
						stockID = textArr[0]
						info.ID = stockID
						info.Name = textArr[1]
					case 4:
						info.Type = text
					}
				})
				// if err := insertStockInfoToDB(*info); err != nil {
				// 	panic(err)
				// }
				stockList = append(stockList, *info)
				// stockMap[stockID] = *info
			}
		})
	})

	// Start scraping on https://hackerspaces.org
	c.Visit("http://isin.twse.com.tw/isin/C_public.jsp?strMode=2")

	fmt.Println(stockList)
	return stockList
}
