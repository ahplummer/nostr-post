package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/leekchan/accounting"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
	log "github.com/sirupsen/logrus"
	"io"
	"math"
	"net/http"
	"os"
	"strings"
	"time"
)

var nsec, cmc string

type Bitcoin struct {
	Data []struct {
		Tags []struct {
			Slug     string `json:"slug"`
			Name     string `json:"name"`
			Category string `json:"category"`
		} `json:"tags"`
		ID                            int         `json:"id"`
		Name                          string      `json:"name"`
		Symbol                        string      `json:"symbol"`
		Slug                          string      `json:"slug"`
		IsActive                      int         `json:"is_active"`
		InfiniteSupply                bool        `json:"infinite_supply"`
		IsFiat                        int         `json:"is_fiat"`
		CirculatingSupply             int         `json:"circulating_supply"`
		TotalSupply                   int         `json:"total_supply"`
		MaxSupply                     int         `json:"max_supply"`
		DateAdded                     time.Time   `json:"date_added"`
		NumMarketPairs                int         `json:"num_market_pairs"`
		CmcRank                       int         `json:"cmc_rank"`
		LastUpdated                   time.Time   `json:"last_updated"`
		TvlRatio                      interface{} `json:"tvl_ratio"`
		Platform                      interface{} `json:"platform"`
		SelfReportedCirculatingSupply interface{} `json:"self_reported_circulating_supply"`
		SelfReportedMarketCap         interface{} `json:"self_reported_market_cap"`
		Quote                         []struct {
			ID                    int         `json:"id"`
			Symbol                string      `json:"symbol"`
			Price                 float64     `json:"price"`
			Volume24H             float64     `json:"volume_24h"`
			VolumeChange24H       float64     `json:"volume_change_24h"`
			PercentChange1H       float64     `json:"percent_change_1h"`
			PercentChange24H      float64     `json:"percent_change_24h"`
			PercentChange7D       float64     `json:"percent_change_7d"`
			PercentChange30D      float64     `json:"percent_change_30d"`
			PercentChange60D      float64     `json:"percent_change_60d"`
			PercentChange90D      float64     `json:"percent_change_90d"`
			MarketCap             float64     `json:"market_cap"`
			MarketCapDominance    float64     `json:"market_cap_dominance"`
			FullyDilutedMarketCap float64     `json:"fully_diluted_market_cap"`
			Tvl                   interface{} `json:"tvl"`
			LastUpdated           time.Time   `json:"last_updated"`
		} `json:"quote"`
	} `json:"data"`
	Status struct {
		Timestamp    time.Time `json:"timestamp"`
		ErrorCode    string    `json:"error_code"`
		ErrorMessage string    `json:"error_message"`
		Elapsed      int       `json:"elapsed"`
		CreditCount  int       `json:"credit_count"`
	} `json:"status"`
}
type FearAndGreedScore struct {
	Data struct {
		Value               int       `json:"value"`
		UpdateTime          time.Time `json:"update_time"`
		ValueClassification string    `json:"value_classification"`
	} `json:"data"`
	Status struct {
		Timestamp    time.Time `json:"timestamp"`
		ErrorCode    string    `json:"error_code"`
		ErrorMessage string    `json:"error_message"`
		Elapsed      int       `json:"elapsed"`
		CreditCount  int       `json:"credit_count"`
	} `json:"status"`
}

func Publish(message string) {

	prefix, privKey, err := nip19.Decode(nsec)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("nsec key %s shows to be prefix: %s, value: %s", nsec, prefix, privKey)
	pub, _ := nostr.GetPublicKey(fmt.Sprintf("%s", privKey))
	log.Infof("Public key %s", pub)

	tag := []nostr.Tag{nostr.Tag{"bitcoin"}}
	ev := nostr.Event{
		PubKey:    pub,
		CreatedAt: nostr.Now(),
		Kind:      nostr.KindTextNote,
		Tags:      tag,
		Content:   message,
	}

	// calling Sign sets the event ID field and the event Sig field
	err = ev.Sign(privKey.(string))
	if err != nil {
		log.Fatal(err)
	}
	// publish the event to two relays
	ctx := context.Background()
	relays := []string{"nos.lol", "nostr-pub.wellorder.net", "relay.damus.io", "http://umbrel.local:4848"} // Creates a slice with initial values
	//[]string{"wss://relay.stoner.com", "wss://nostr-pub.wellorder.net"
	for _, url := range relays {
		relay, err := nostr.RelayConnect(ctx, url)
		if err != nil {
			log.Errorf("Error: %s", err)
			continue
		}
		if err := relay.Publish(ctx, ev); err != nil {
			log.Errorf("Error: %s", err)
			continue
		}
		log.Infof("published to %s\n", url)
	}

}
func FormatPrice(price float64) string {
	ac := accounting.Accounting{Symbol: "$", Precision: 2}
	return ac.FormatMoney(price)
}
func getFearAndGreed(api_key string) (index int, classification string) {
	url := "https://pro-api.coinmarketcap.com/v3/fear-and-greed/latest"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("X-CMC_PRO_API_KEY", api_key)
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	var f FearAndGreedScore
	err = json.Unmarshal(body, &f)
	if err != nil {
		log.Fatal(err)
	}
	index = f.Data.Value
	classification = f.Data.ValueClassification
	return index, classification
}
func getBlockHeight() string {
	url := "https://mempool.space/api/blocks/tip/height"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}
func getUSDPrice(api_key string) (price string, sats int) {
	//number=`curl -s --header "X-CMC_PRO_API_KEY: "$CMC_API"" "https://pro-api.coinmarketcap.com/v3/cryptocurrency/quotes/latest?slug=bitcoin"  | jq -r '.data[0].quote[0].price'`
	url := "https://pro-api.coinmarketcap.com/v3/cryptocurrency/quotes/latest?slug=bitcoin"
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("X-CMC_PRO_API_KEY", api_key)
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	var b Bitcoin
	err = json.Unmarshal(body, &b)
	if err != nil {
		log.Fatal(err)
	}
	p := b.Data[0].Quote[0].Price
	price = FormatPrice(p)
	log.Infof("Price: %s", price)
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	//get sats per dollar
	exactsats := p / 100000000
	exactsatsperdollar := 1 / exactsats
	sats = int(math.Round(exactsatsperdollar))
	return price, sats
}
func initlogging() {
	//Trace, Debug, Info, Warn, Error, Fatal, and Panic are valid
	logfmt := os.Getenv("LOG_FORMAT")
	if strings.ToUpper(logfmt) == "JSON" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{
			DisableColors: false,
			FullTimestamp: true,
		})
	}
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	// LOG_LEVEL not set, let's default to debug
	if !ok {
		lvl = "debug"
	}
	// parse string, this is built-in feature of logrus
	ll, err := log.ParseLevel(lvl)
	if err != nil {
		ll = log.DebugLevel
	}
	// set global log level
	log.SetLevel(ll)
	log.Debug("Finished configuring logger.")
}
func initConfig() {
	initlogging()
	cmc = os.Getenv("CMC_API")
	nsec = os.Getenv("NSEC")
	log.Debugf("Finished initConfig, will use nsec %s, cmc key %s.", nsec, cmc)
}

func main() {
	initConfig()
	log.Info("nostr-post API started")
	price, sats := getUSDPrice(cmc)
	index, classification := getFearAndGreed(cmc)
	blockHeight := getBlockHeight()
	message := fmt.Sprintf("The block height is %s, and the current price for bitcoin is %s USD. This means that you can get %d sats for one dollar. The fear/greed index is %d: %s. Buy some ₿, it may catch on.", blockHeight, price, sats, index, classification)
	log.Infof("%s", message)
	Publish(message)
}
