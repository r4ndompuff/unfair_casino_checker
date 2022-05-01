package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

var g float32

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
func WriteFile(m chan string) {
	f, err := os.OpenFile("text.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	for text := range m {
		if _, err := f.WriteString(text + "\n"); err != nil {
			log.Println(err)
		}
	}
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

func Unparse(resp string) (stroka string, err error) {
	badids := []string{"76561199085192846"}
	type Get struct {
		Date string `json:"date"`
		Data struct {
			Crash float32 `json:"crash"`
			Id    int     `json:"id"`
			Bets  []struct {
				CreatedAt string `json:"createdAt"`
				Deposit   struct {
					Amount float32       `json:"amount"`
					Items  []interface{} `json:"items"`
				} `json:"deposit"`
				Withdraw struct {
					Amount float32 `json:"amount"`
				} `json:"withdraw"`
				User struct {
					Steamid string `json:"steamId"`
				} `json:"user"`
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

	var withdrawsumm float32
	withdrawsumm = 0
	for _, s := range dat.Data.Bets {
		withdrawsumm = withdrawsumm + s.Withdraw.Amount
	}

	allitems := 0
	for _, s := range dat.Data.Bets {
		//fmt.Println(s.Deposit.Amount)
		allitems = allitems + len(s.Deposit.Items)
	}

	timemap := make(map[string]int)
	for _, s := range dat.Data.Bets {
		layout := "2006-01-02T15:04:05.000Z"
		t, err := time.Parse(layout, s.CreatedAt)

		if err != nil {
			fmt.Println(err)
			continue
		}

		parsedtime := t.Format("2006-01-02 15:04")
		timemap[parsedtime] += 1
	}
	fmt.Println(timemap)
	var maxk string
	var maxv = 0
	for k, v := range timemap {
		if v > maxv {
			maxv = v
			maxk = k
		}

	}
	fmt.Println(summ) //summa babla

	//stroka = fmt.Sprint(dat.Data.Id, ";", summ-withdrawsumm, ";", summ, ";", withdrawsumm, ";", maxk)
	//	if (summ-withdrawsumm) > 3000.0 || (summ-withdrawsumm) < -3000.0 {
	//		fmt.Println("BADPROFIT")
	//	} else {
	//		g = g + (summ - withdrawsumm)
	//	}
	badgame := false
	for _, s := range dat.Data.Bets {
		if contains(badids, s.User.Steamid) {
			badgame = true
			fmt.Println("BADPROFIT")
			break
		}
	}
	if badgame == false {
		g = g + (summ - withdrawsumm)
	}

	//ID;SUMMA_VSEH-SUMMA_VIPLAT;SUMMA_VSEH;SUMMA_VIPLAT;GAMETIME
	stroka = fmt.Sprint(dat.Data.Id, ";", summ-withdrawsumm, ";", summ, ";", withdrawsumm, ";", maxk, ";", g, ";", badgame)
	fmt.Println(stroka)
	fmt.Println(g)
	return stroka, nil
}

func main() {
	gid := 2146676
	count := 0
	m := make(chan string)
	go WriteFile(m)
	for {
		if gid >= 2146678 {
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
		m <- text
		time.Sleep(250 * time.Millisecond)
	}
	fmt.Println(g)

}
