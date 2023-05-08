package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/chamzzzzzz/supersimplesoup"
)

func main() {
	for {
		archive()
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location())
		log.Printf("next archive at %s\n", next.Format("2006-01-02 15:04:05"))
		time.Sleep(next.Sub(now))
	}
}

func archive() {
	log.Printf("start archive at %s\n", time.Now().Format("2006-01-02 15:04:05"))

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.kuaishou.com/?isHome=1", nil)
	if err != nil {
		log.Printf("new request failed, err:%v\n", err)
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.5005.61 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("get http reponse failed, err:%v\n", err)
		return
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("read body failed, err:%v\n", err)
		return
	}

	dom, err := supersimplesoup.Parse(bytes.NewReader(b))
	if err != nil {
		log.Printf("parse html failed, err:%v\n", err)
		return
	}

	os.Mkdir("archives", 0755)
	name := fmt.Sprintf("archives/%s.txt", time.Now().Format("2006-01-02"))
	b, err = os.ReadFile(name)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("read archive file failed, err:%v\n", err)
			return
		}
	}

	var words []string
	if len(b) > 0 {
		words = strings.Split(string(b), "\r\n")
	}

	n := 0
	for _, p := range dom.QueryAll("p", "class", "rank-name") {
		word := p.Text()
		word = strings.TrimSpace(word)
		word = strings.ReplaceAll(word, "\r\n", "")
		has := false
		for _, w := range words {
			if w == word {
				has = true
				break
			}
		}
		if !has {
			words = append(words, word)
			n++
		}
	}

	err = os.WriteFile(name, []byte(strings.Join(words, "\r\n")), 0755)
	if err != nil {
		log.Printf("write archive file failed, err:%v\n", err)
		return
	}

	log.Printf("archived %d new words\n", n)
	log.Printf("finish archive at %s\n", time.Now().Format("2006-01-02 15:04:05"))
}
