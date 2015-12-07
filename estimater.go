package main

import (
	"log"
	"errors"
	"os"
	"bufio"
	"strings"
)

const TRAIN_DATA = "train_datas/neko.num"
const TRAIN_DATA_DIC = "train_datas/neko.dic.txt"

var data []string
var dataDic map[string]int

type Model struct {
	Word string
	Index int
	Count int
	Prob float64
}

type Estimater struct {
	TypeNgram string
	TargetWords []string
	Models []Model
}

func NewEstimater(words ...string) *Estimater {
	typeNgram, err := ngram(len(words))
	if err != nil { log.Fatal(err) }
	estimater := &Estimater{ TypeNgram: typeNgram, TargetWords: words, Models: nil }
	return estimater
}

func ngram(n int) (string, error) {
	switch n {
	case 0:
		return "unigram", nil
	case 1:
		return "bigram", nil
	case 2:
		return "trigram", nil
	default:
		return "", errors.New("ngram: This program use until trigram.(argument out of range)")
	}
}

func loadDic() {
	f, err := os.Open(TRAIN_DATA_DIC)
	if err != nil { log.Fatal(err) }

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "	")
		dataDic[line[1]] = line[0]
	}

	log.Println(dataDic)

	if err = scanner.Err(); err != nil { log.Fatal(err) }
}

func loadData() {
}

func init() {
	loadDic()
	loadData()
}

func main() {
	unigramEstimater := NewEstimater()
	log.Println(unigramEstimater)
}

