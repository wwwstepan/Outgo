package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func loadClassifier() {
	clsf := loadClassifierFromJSON()
	if len(clsf) == 0 {
		// default classifier
		clsf = append(clsf, outgoType{'1', "еда", 1})
		clsf = append(clsf, outgoType{'2', "рестораны", 2})
		clsf = append(clsf, outgoType{'3', "здоровье", 3})
		clsf = append(clsf, outgoType{'4', "одежда", 4})
		clsf = append(clsf, outgoType{'5', "техника", 5})
		clsf = append(clsf, outgoType{'6', "путешествия", 6})
		clsf = append(clsf, outgoType{'7', "для дома", 7})
		clsf = append(clsf, outgoType{'8', "помощь, подарки", 8})
		clsf = append(clsf, outgoType{'9', "оплата жилья", 9})
		clsf = append(clsf, outgoType{'0', "прочее", 10})
	}
	ogTypeClassifier = clsf
}

func loadClassifierFromJSON() []outgoType {
	clsf := make([]outgoType, 0, 20)

	f, err := os.Open("outgoitems.json")
	if err != nil {
		return clsf
	}
	defer f.Close()

	byteValue, err := ioutil.ReadAll(f)
	if err != nil {
		return clsf
	}

	var clsfJ outgoClassifierJSON
	json.Unmarshal(byteValue, &clsfJ)

	uniCodes := make(map[rune]bool)
	uniNames := make(map[string]bool)
	uniInd := make(map[int32]bool)

	for _, v := range clsfJ.Items {

		acode := []rune(v.Key)[:1] // Only 1st symbol

		if len(acode) == 0 {
			continue
		}
		code := acode[0]

		if _, ok := uniCodes[code]; ok {
			continue // code not unical
		}

		rname := []rune(v.Description)
		if len(rname) == 0 {
			continue
		}
		if len(rname) > 20 { // max length of description is 20
			rname = rname[:20]
		}

		name := string(rname)

		if _, ok := uniNames[name]; ok {
			continue // name not unical
		}

		ind := v.Code

		if _, ok := uniInd[ind]; ok {
			continue // index not unical
		}

		if ind < 1 || ind > 1023 { // index range: 1..1023
			continue
		}

		// Correct item
		clsf = append(clsf, outgoType{code: code, name: name, ind: ind})
		uniCodes[code] = true
		uniInd[ind] = true
	}

	return clsf
}
