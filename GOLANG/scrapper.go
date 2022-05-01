package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func WriteFile(text string) {
	f, err := os.OpenFile("text.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(text + "\n"); err != nil {
		log.Println(err)
	}
}

func MakeRequest(gid int) (resp string) {
	u := "https://api.csgorun.pro/games/" + strconv.Itoa(gid)
	opt := hhttp.NewOptions().Retry(10)
	opt.Timeout(5 * time.Second)
	opt.DNS("1.1.1.1:53")
	r, err := hhttp.NewClient().SetOptions(opt).Get(u).Do()
	fmt.Println("popitki",r.Attempts)
	if err != nil {
		log.Fatal(err)
	}
	return r.Body.String()

}

func Unparse(resp string) (stroka string, err error) {
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
	}
	var dat Get

	if err := json.Unmarshal([]byte(resp), &dat); err != nil {
		return "", err
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
	stroka = fmt.Sprint(dat.Data.Id, ";", len(dat.Data.Bets), ";", allitems, ";", summ, ";", dat.Data.Crash)
	fmt.Println(stroka)
	return stroka, nil
}

func main() {
	gid := 2146677
	count := 0
	for {
		if gid >= 2169483 {
			break
		}
		resp := MakeRequest(gid)
		fmt.Println(count)
		text, err := Unparse(resp)
		count++
		gid++
		if err != nil {
			continue
		}
		WriteFile(text)
	}

}