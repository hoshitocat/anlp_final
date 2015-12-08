package main

import (
	"fmt"
	"log"
	"errors"
	"os"
	"io/ioutil"
	"bufio"
	"strings"
	"strconv"
	"sort"
)

const TRAIN_DATA = "train_datas/neko.num"
const TRAIN_DATA_DIC = "train_datas/neko.dic.txt"
const BIGRAM_MODEL_PATH = "./bigram.model"
const TRIGRAM_MODEL_PATH = "./trigram.model"

var data []int
var dataDic map[string]int

type Model struct {
	Word string
	Index int
	Count int
	Prob float64
}

type Models []Model

type Estimater struct {
	TypeNgram string
	TargetWords []string
	Models Models
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

func calcAbsDiscount(nr int, models []Model) float64 {
	numOfNr := 0
	numOfNextNr := 0

	for _, model := range models {
		if model.Count == nr { numOfNr++ }
		if model.Count == (nr + 1) { numOfNextNr++ }
	}

	return (float64(numOfNr) / float64(numOfNr + 2 * numOfNextNr))
}

func calcBigram(w string) []Model {
	models := make([]Model, 0)
	sum, indexes := appearCount(dataDic[w])
	absDiscount := make(map[int]float64)

	for key, val := range dataDic {
		count := appearCountAfterWords(val, indexes)
		model := Model{ Word: key, Index: val, Count: count, Prob: 0 }
		models = append(models, model)
	}

	for i, model := range models {
		if model.Count == 0 {
			models[i].Prob = float64(0)
			continue
		}
		if _, ok := absDiscount[model.Count]; !ok { absDiscount[model.Count] = calcAbsDiscount(model.Count, models) }
		models[i].Prob = (float64(model.Count) - absDiscount[model.Count]) / float64(sum)
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

// NOTE: ここからソート
func (m Models) Len() int {
	return len(m)
}

func (m Models) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m Models) Less(i, j int) bool {
	return m[i].Index < m[j].Index
}

func modelWrite(models Models) {
	sumPer := 0.0
	for _, model := range models { sumPer += model.Prob }

	content := fmt.Sprintf("%20.17e\n", 1 - sumPer)
	for _, model := range models {
		content += fmt.Sprintf("%20.17e\n", model.Prob)
	}

	ioutil.WriteFile(BIGRAM_MODEL_PATH, []byte(content), os.ModePerm)
}
// NOTE: ここまで

func backOff(bigramModels Models, unigramModels Models) {
	sumPer := 0.0
	var remainPer float64
	var zeroPerIndexes []int
	zeroPerIndexes = make([]int, 0)

	for i, model := range bigramModels {
		sumPer += model.Prob
		if model.Prob == 0 { zeroPerIndexes = append(zeroPerIndexes, i) }
	}
	remainPer = 1 - sumPer

	for _, index := range zeroPerIndexes {
		unigramPer := 0.0
		for _, model := range unigramModels {
			if model.Index == bigramModels[index].Index { unigramPer = model.Prob }
		}
		bigramModels[index].Prob = remainPer * unigramPer
	}
}

func main() {
	unigramEstimater, err := NewEstimater()
	if err != nil { log.Fatal(err) }

	bigramEstimater, err := NewEstimater("て")
	if err != nil { log.Fatal(err) }

	backOff(bigramEstimater.Models, unigramEstimater.Models)

	sort.Sort(bigramEstimater.Models)
	modelWrite(bigramEstimater.Models)

	// log.Println(bigramEstimater.Models)
	// log.Println(bigramEstimater.TypeNgram)
	// sum := 0.0
	// for _, model := range bigramEstimater.Models {
	// 	sum += model.Prob
	// }
	// log.Println(sum)
}

