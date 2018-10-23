package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Team struct {
	id        string
	name      string
	lastScore int
	bestScore int
	time      string
}

var idToNameMap = map[string]string{
	"1":  "しっこくのすぐちむ",
	"2":  "ふぇっちでこーどえぐぜくらいと",
	"3":  "おかねかぞく",
	"4":  "ているど",
	"5":  "どっとだっと",
	"6":  "ここでいっく",
	"7":  "しとうのはてに",
	"8":  "たけだし",
	"9":  "なるせじゅん",
	"10": "だいあもんどぷりんせす",
	"11": "ちーむにんげんせい",
	"12": "えむえぬしす",
	"13": "でぃめんじょなるはいそさいえてぃぬれねずみ",
	"14": "しゃからん",
	"15": "やまがたぐみ",
	"16": "けーぜろに",
	"17": "さいだいのてきはじさ",
	"18": "のんだくれおんけんは",
	"19": "りーくれい",
	"20": "ひゃくまんえんどりぶん",
	"21": "あいきゅーいち",
	"22": "こくさいこうとうきょういくいんとう",
	"23": "けーえすてぃーえむ",
	"24": "じんごにっく",
	"25": "せあぶら",
	"26": "じゅけんせいのかたき",
	"27": "うるとらふぁすとごーふぁー",
	"28": "かんがえちゅう",
	"29": "えすきゅーえるいんじぇくしょん",
	"30": "しょぼーん",
}

var (
	name = os.Getenv("NAME")
	pass = os.Getenv("PASSWORD")
	data = make(map[string]Team)
)

func main() {

	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: cookieJar}

	if err := login(client); err != nil {
		log.Println(err)
		return
	}

	log.Println("init")
	doc, err := scores(client)
	if err != nil {
		log.Println(err)
		return
	}

	doc.Find(".content .table tbody tr").Each(func(i int, s *goquery.Selection) {
		lastScore, _ := strconv.Atoi(strings.Replace(s.Find("td").Eq(4).Text(), ",", "", -1))
		bestScore, _ := strconv.Atoi(strings.Replace(s.Find("td").Eq(5).Text(), ",", "", -1))
		team := Team{
			id:        s.Find("td").Eq(1).Text(),
			name:      s.Find("td").Eq(2).Text(),
			lastScore: lastScore,
			bestScore: bestScore,
			time:      s.Find("td").Eq(6).Text(),
		}
		data[team.id] = team
	})

	//log.Printf("%#v", data)

	for {
		doc, err := scores(client)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second * 5)
			continue
		}
		log.Println("check")

		doc.Find(".content .table tbody tr").Each(func(i int, s *goquery.Selection) {
			lastScore, _ := strconv.Atoi(strings.Replace(s.Find("td").Eq(4).Text(), ",", "", -1))
			bestScore, _ := strconv.Atoi(strings.Replace(s.Find("td").Eq(5).Text(), ",", "", -1))
			curr := Team{
				id:        s.Find("td").Eq(1).Text(),
				name:      s.Find("td").Eq(2).Text(),
				lastScore: lastScore,
				bestScore: bestScore,
				time:      s.Find("td").Eq(6).Text(),
			}
			if prev, _ := data[curr.id]; prev.time != curr.time && prev.bestScore < curr.lastScore {
				diff := curr.lastScore - prev.bestScore
				fmt.Println(idToNameMap[curr.id], " すこあ ", curr.lastScore, " ぷらす ", diff)
				data[curr.id] = curr
			}
		})

		time.Sleep(time.Second * 30)
	}
}

func login(client *http.Client) error {
	resp, err := client.PostForm(
		"https://portal.isucon8.flying-chair.net/admin/login",
		url.Values{"name": {name}, "password": {pass}},
	)
	if err != nil {
		return err
	}
	return resp.Body.Close()
}

func scores(client *http.Client) (*goquery.Document, error) {
	resp, err := client.Get("https://portal.isucon8.flying-chair.net/admin/scores")
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromResponse(resp)
}
