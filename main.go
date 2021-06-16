package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type ResponseBody struct {
	LastUpdateId int64      `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

// func hello(w http.ResponseWriter, r *http.Request) {
// 	io.WriteString(w, "hello world")
// 	send("Deployed Successfully, Happy trading", "Deployment Status")
// 	i := true
// 	for i {
// 		plot(6.0)
// 		time.Sleep(30 * time.Second)
// 		//i = false
// 	}
// }
func main() {
	//fmt.Println("happy trading")
	// log.Fatal(http.ListenAndServe(":"+port, nil))
	send("Deployed Successfully, Happy trading", "Deployment Status")
	ii := 0
	i := true
	for i {
		plot(6.0)
		time.Sleep(60 * time.Second)

		//i = false
		ii = ii + 1
		sendSMS("testing Deployment " + strconv.Itoa(ii))

	}

}

func plot(threshold float64) {

	var sym, lim []string
	bla := false

	sym = []string{"BTC", "ETH", "BNB", "ADA", "XRP", "DOGE", "DOT", "UNI", "ICP", "LINK", "BCH", "LTC", "MATIC", "SOL", "XLM", "VET", "ETC", "FIL", "TRX", "EOS", "XMR", "KSM", "AAVE", "NEO", "MKR", "CAKE", "TFUEL", "ALGO", "XTZ", "ATOM", "SHIB", "AVAX", "LUNA", "BTT", "RUNE", "COMP", "HBAR", "DCR", "DASH", "ZEC", "EGLD", "WAVES", "YFI", "HOT", "CHZ", "BTG", "SUSHI", "ZIL", "NEAR", "SNX", "MANA", "ENJ"}

	lim = []string{"1000", "500", "500", "100", "100", "500", "100", "100", "100", "100", "100", "100", "100", "100"}

	mailBody := ""
	smsBody := ""
	for i, v := range sym {
		//fmt.Println(v) //1

		li := "100"

		//var orders ResponseBody
		if i > 40 {
			li = "50"
		} else if i > 13 {
			li = "100"
		} else {
			li = lim[i]
		}

		orders := getOrders(v, li)
		//fmt.Println("size ", len(orders.Asks)) //2
		askVol, perChange1 := findTotalVolbyPrice(orders.Asks)
		askVol = askVol * (-1)
		bidVol, perChange2 := findTotalVolbyPrice(orders.Bids)

		perChange := (perChange2 - perChange1) / 2
		//fmt.Println("Potential ", perChange)

		buyfactor := bidVol / askVol
		//fmt.Println("buy factor for ", v, buyfactor) //
		if buyfactor > threshold {
			//fmt.Println("Buy", "buy "+v+" at "+fmt.Sprintf("%f", buyfactor), "nil")
			bla = true
			mssgBody := `
			Buy ` + v + `
			Buy Factor =` + fmt.Sprintf("%f", buyfactor) + `
			Price =` + orders.Bids[0][0] + `
			Delta =` + fmt.Sprintf("%f", perChange) + `
			Upper Price =` + orders.Asks[len(orders.Asks)-1][0] + `
			Lower Price =` + orders.Bids[len(orders.Bids)-1][0] + `
			
			`
			smsBody = smsBody + `
			Buy ` + v + ` potential ` + fmt.Sprintf("%f", buyfactor)
			mailBody = mailBody + mssgBody

			//send(mssgBody,"Buy "+v)
		}

	}
	//fmt.Println(mailBody)
	if bla {
		send(mailBody, "Buy at threshold "+fmt.Sprintf("%f", threshold))
		sendSMS(smsBody)
	}
}

func getOrders(sym, lim string) ResponseBody {
	url := "https://api.binance.com/api/v3/depth?symbol=" + sym + "USDT&limit=" + lim
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ResponseBody{}
	}

	resp, responseErr := client.Do(req)
	if responseErr != nil {
		fmt.Println("Error : " + err.Error())
		return ResponseBody{}
	}
	defer resp.Body.Close()
	var orders ResponseBody
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	json.Unmarshal([]byte(string(body)), &orders)

	return orders
}

func findTotalVolbyPrice(pricevsVol [][]string) (float64, float64) {
	vol := 0.0
	for i := range pricevsVol {
		ii, _ := strconv.ParseFloat(pricevsVol[i][1], 64)
		vol = vol + ii

	}

	h, _ := strconv.ParseFloat(pricevsVol[0][0], 64)
	l, _ := strconv.ParseFloat(pricevsVol[len(pricevsVol)-1][0], 64)
	pricediff := h - l
	//fmt.Println(h, l)
	per := ((h - l) / h) * 100
	return (vol / pricediff), per
	//return (vol)
}

func send(body, sub string) {
	from := "apigeetest23@gmail.com"
	pass := ".mnbvcxz.1"
	to := "shahulhameed28p@gmail.com"

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + sub + "\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}

	log.Print("sent")
}

func sendSMS(body string) {
	accountSid := "AC0b9e944cbf5a09e163c975cb558c2efb"
	authToken := "1a6b0683872a56b2ff5d56e023de20ec"
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"
	msgData := url.Values{}
	msgData.Set("To", "+918137803478")
	msgData.Set("From", "+14158775523")
	msgData.Set("Body", body)
	msgDataReader := *strings.NewReader(msgData.Encode())
	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			fmt.Println(data["sid"])
		}
	} else {
		fmt.Println(resp.Status)
	}
	//+14158775523
}
