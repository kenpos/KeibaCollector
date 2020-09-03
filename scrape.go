package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/saintfish/chardet"
	"golang.org/x/net/html/charset"
)

type URL struct {
	race_date       string
	race_course_num string
	race_info       string
	race_count      string
	race_no         string
}

//output the Err
func PrintlnErr(err error) {
	if err != nil {
		panic(err)
	}
}

// charactor replacer
func ReplaceCharactor(str string) string {
	tmp := strings.Replace(str, "\n", ",", -1)
	rep := regexp.MustCompile(`\,{2,}`)
	str = rep.ReplaceAllString(tmp, ",")
	str = removeCheckMark(str)
	str = rep.ReplaceAllString(str, ",")
	return str
}

func removeCheckMark(strtmp string) string {
	rep := regexp.MustCompile(`--◎◯▲△☆&#10003消`)
	str := rep.ReplaceAllString(strtmp, "")
	return str
}

// get the paged url
func GetPage(url string) {
	// Getリクエスト
	res, err := http.Get(url)
	PrintlnErr(err)

	defer res.Body.Close()

	// 読み取り
	buf, err := ioutil.ReadAll(res.Body)
	PrintlnErr(err)

	// 文字コード判定
	det := chardet.NewTextDetector()
	detRslt, err := det.DetectBest(buf)
	PrintlnErr(err)

	fmt.Println(detRslt.Charset)
	// => EUC-JP

	// 文字コード変換
	bReader := bytes.NewReader(buf)
	reader, err := charset.NewReaderLabel(detRslt.Charset, bReader)
	PrintlnErr(err)

	// HTMLパース
	doc, err := goquery.NewDocumentFromReader(reader)
	PrintlnErr(err)

	RaceListNameBox := doc.Find(".RaceNum")
	RaceListNameBox.Each(func(i int, s *goquery.Selection) {
		RaceNum := ReplaceCharactor(s.Text())
		fmt.Println(RaceNum)
	})

	RaceListHorseList := doc.Find(".HorseList")
	RaceListHorseList.Each(func(i int, s *goquery.Selection) {
		HorseList := ReplaceCharactor(s.Text())
		fmt.Println(HorseList)
	})
}

func createID() {
}

func main() {
	file, err := os.Create("./RaceList.csv")
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	url := "http://oldrace.netkeiba.com/?pid=race_old&id=c202008010601"
	GetPage(url)
}
