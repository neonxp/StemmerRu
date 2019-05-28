package StemmerRu

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

var testFile = `tests.json`

func TestStemWord(t *testing.T) {
	file, err := ioutil.ReadFile(testFile)
	if err != nil {
		t.Error("Can't open file", testFile)
	}
	tests := &map[string]string{}
	err = json.Unmarshal(file, tests)
	if err != nil {
		t.Error("Can't parse json", err)
	}
	for source, expected := range *tests {
		result := Stem(source)
		if expected != result {
			t.Errorf(`Expected "%s" (source: %s) but got "%s"`, result, source, result)
		}
	}
}
