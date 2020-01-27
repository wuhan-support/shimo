package shimo

import (
	"io/ioutil"
	"testing"
)

func TestDocument_GetJSON(t *testing.T) {
	d := NewDocument("xxx", "xxx")
	d.EliminateSuffix = "（"
	data, err := d.GetJSON()
	if err != nil {
		t.Fatalf("failed to get json: %v", err)
	}
	//spew.Dump(data)
	data, err = d.GetJSON()
	if err != nil {
		t.Fatalf("failed to get json: %v", err)
	}
	//spew.Dump(data)
	marshalled, err := data.MarshalJSON()
	if err != nil {
		t.Fatalf("failed to marshal json: %v", err)
	}
	err = ioutil.WriteFile("out.json", marshalled, 0644)
	if err != nil {
		t.Fatalf("failed to write json file: %v", err)
	}
}
