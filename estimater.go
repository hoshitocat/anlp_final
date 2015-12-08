package main

import (
	"log"
	"errors"
	"os"
	"bufio"
	"strings"
	"strconv"
)

const TRAIN_DATA = "train_datas/neko.num"
const TRAIN_DATA_DIC = "train_datas/neko.dic.txt"

var data []int
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

// NOTE: UnigramのときはIndexesは不要
func appearCount(wordIndex int) (int, []int) {
	sum := 0
	indexes := make([]int, 0)
	for k, v := range data {
		if v == wordIndex {
			sum++
			indexes = append(indexes, k)
		}
	}
	return sum, indexes
}

func appearCountAfterWords(wordIndex int, indexes []int) int {
	sum := 0
	for _, index := range indexes {
		if len(data) < (index + 1) { continue }
		if wordIndex == data[index + 1] { sum++ }
	}

	return sum
}

func calcUnigram() []Model {
	models := make([]Model, 0)
	sum := len(data)

	for key, val := range dataDic {
		count, _ := appearCount(val)
		model := Model{ Word: key, Index: val, Count: count, Prob: (float64(count) / float64(sum)) }
		models = append(models, model)
	}

	return models
}

func calcBigram(w string) []Model {
	models := make([]Model, 0)
	sum, indexes := appearCount(dataDic[w])

	for key, val := range dataDic {
		count := appearCountAfterWords(val, indexes)
		model := Model{ Word: key, Index: val, Count: count, Prob: (float64(count) / float64(sum)) }
		models = append(models, model)
	}
	return models
}

func calcTrigram(w1 string, w2 string) []Model {
	models := make([]Model, len(dataDic))
	return models
}

func NewEstimater(words ...string) (*Estimater, error) {
	var typeNgram string
	var models []Model

	switch len(words) {
	case 0:
		typeNgram = "unigram"
		models = calcUnigram()
	case 1:
		typeNgram = "bigram"
		models = calcBigram(words[0])
	case 2:
		typeNgram = "trigram"
		models = calcTrigram(words[0], words[1])
	default:
		return nil, errors.New("ngram: This program use until trigram.(argument out of range)")
	}

	estimater := &Estimater{ TypeNgram: typeNgram, TargetWords: words, Models: models }
	return estimater, nil
}

func loadDic() {
	dataDic = make(map[string]int)
	f, err := os.Open(TRAIN_DATA_DIC)
	if err != nil { log.Fatal(err) }

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "\t")
		key := line[1]
		val, err := strconv.Atoi(line[0])
		if err != nil { log.Fatal(err) }
		dataDic[key] = val
	}

	if err = scanner.Err(); err != nil { log.Fatal(err) }
}

func loadData() {
	data = make([]int, len(dataDic))
	f, err := os.Open(TRAIN_DATA)
	if err != nil { log.Fatal(err) }

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		for _, w := range line {
			w, err := strconv.Atoi(w)
			if err != nil { log.Fatal(err) }
			data = append(data, w)
		}
	}

	if err = scanner.Err(); err != nil { log.Fatal(err) }
}

func init() {
	loadDic()
	loadData()
}

func main() {
	// unigramEstimater, err := NewEstimater()
	// if err != nil { log.Fatal(err) }

	bigramEstimater, err := NewEstimater("て")
	if err != nil { log.Fatal(err) }
}

