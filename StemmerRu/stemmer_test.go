package StemmerRu

import (
	"testing"
	"io/ioutil"
	"encoding/json"
	"path"
)

var testFile = path.Join(`..`, `tests.json`)

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
		result := StemWord(source);
		if expected != result {
			t.Errorf(`Expected "%s" (source: %s) but got "%s"`, result, source, result)
		}
	}
}
