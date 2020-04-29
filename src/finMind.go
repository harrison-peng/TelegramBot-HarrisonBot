package harrisonbot

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func getTaiwanStockNews(stockID string) ([]News, error) {
	var newsList []News

	URL, err := url.Parse(FINMINDURL)
	if err != nil {
		return newsList, err
	}
	params := url.Values{}
	params.Add("dataset", "TaiwanStockNews")
	params.Add("stock_id", stockID)
	params.Add("date", time.Now().Format("2006-01-02"))
	URL.RawQuery = params.Encode()

	res, err := http.Get(URL.String())
	if err != nil {
		return newsList, err
	}

	news := &TaiwanStockNews{}
	if err := json.NewDecoder(res.Body).Decode(news); err != nil {
		return newsList, err
	}

	for i := 0; i < len(news.Data.Date); i++ {
		tinyURL := TINYURL + "?url=" + news.Data.Link[i]
		resp, err := http.Get(tinyURL)
		if err != nil {
			return newsList, err
		}
		newsURLBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return newsList, err
		}

		newsList = append(newsList, News{
			Title: news.Data.Title[i],
			URL:   string(newsURLBytes),
		})
	}
	return newsList, nil
}

// TaiwanStockNews is the mdoel of TaiwanStockNews
type TaiwanStockNews struct {
	Message string `json:"msg"`
	Status  int    `json:"status"`
	Data    struct {
		Date        []string `json:"date"`
		Description []string `json:"description"`
		Link        []string `json:"link"`
		Title       []string `json:"title"`
	} `json:"data"`
}

// News is the modle of the news
type News struct {
	Title string
	URL   string
}
