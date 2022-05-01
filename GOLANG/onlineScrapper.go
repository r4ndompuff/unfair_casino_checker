package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"github.com/keymetrics/pm2-io-apm-go/structures"
	"github.com/keymetrics/pm2-io-apm-go"
)

func WriteFile(m chan string) {
	f, err := os.OpenFile("textOnline.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	for text := range m {
		if _, err := f.WriteString(text); err != nil {
			log.Println(err)
		}
	}
}
func ParseTick(tick float32) int {
	if tick < 1.2 {
		return 0
	} else if tick < 2.0 {
		return 1
	} else if tick < 4.0 {
		return 2

	} else if tick < 8.0 {
		return 3

	} else if tick < 25.0 {
		return 4
	}
	return 5
}

func MakeRequest(gid int) (resp string) {
	u := "https://api.csgorun.pro/games/" + strconv.Itoa(gid)
	opt := hhttp.NewOptions().Retry(10, 5*time.Second)
	opt.Timeout(5 * time.Second)
	opt.DNS("1.1.1.1:53")
	opt.KeepAlive(false)
	r, err := hhttp.NewClient().SetOptions(opt).Get(u).Do()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(r.Attempts)
	return r.Body.String()

}

func Unparse(resp string) (stroka string, success bool, err error) {
	type Get struct {
		Data struct {
			Crash float32 `json:"crash"`
			Id    int     `json:"id"`
			Bets  []struct {
				Deposit struct {
					Amount float32       `json:"amount"`
					Items  []interface{} `json:"items"`
				} `json:"deposit"`
			} `json:"bets"`
		} `json:"data"`
		Success bool `json:"success"`
	}
	var dat Get

	if err := json.Unmarshal([]byte(resp), &dat); err != nil {
		return "", false, err
	}
	var summ float32
	summ = 0
	for _, s := range dat.Data.Bets {
		//fmt.Println(s.Deposit.Amount)
		summ = summ + s.Deposit.Amount
	}
	allitems := 0
	for _, s := range dat.Data.Bets {
		//fmt.Println(s.Deposit.Amount)
		allitems = allitems + len(s.Deposit.Items)
	}
	//fmt.Println(dat.Data.Id)
	//fmt.Println(len(dat.Data.Bets)) //kolvo igrokov
	//fmt.Println(allitems)
	//fmt.Println(summ) //summa babla
	//fmt.Println(dat.Data.Crash)
	//fmt.Println("\n\n\n")
	fmt.Println(dat.Success)
	stroka = fmt.Sprint(dat.Data.Id, ";", len(dat.Data.Bets), ";", allitems, ";", summ, ";", dat.Data.Crash)
	fmt.Println(stroka)
	parsedtick := ParseTick(dat.Data.Crash)
	strokaa := fmt.Sprint(parsedtick, " ")
	return strokaa, dat.Success, nil
}

func main() {
	pm2 := pm2io.Pm2Io{
	    Config: &structures.Config{
	      PublicKey: "publick_key",    // define the public key given in the dashboard
	      PrivateKey: "private_key",  // define the private key given in the dashboard
	      Name: "onlineScrapper",       // define an application name
	    },
	}
	pm2.Start()
	gid := 2179249
	count := 0
	m := make(chan string)
	go WriteFile(m)
	for {
		resp := MakeRequest(gid)
		fmt.Println(count)
		text, success, err := Unparse(resp)
		if err != nil {
			continue
		}
		if success == true {
			count++
			gid++
			m <- text
		} else {
			fmt.Println("game now ", gid)
		}
		fmt.Println("sleep 1 second")
		time.Sleep(1 * time.Second)

	}

}