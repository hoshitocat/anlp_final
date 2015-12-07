package main

import (
	"log"
	"errors"
)

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

func (self *Estimater) Initialize(words ...string) *Estimater {
	typeNgram, err := ngram(len(words))
	if err != nil { log.Errorf(err) }
	estimater := &Estimater{ TypeNgram: typeNgram, , nil }
}

func ngram(n int) string, error {
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

func main() {
	unigramEstimater := &Estimater{"unigram", nil, nil}
	log.Println(unigramEstimater.Type)
}

