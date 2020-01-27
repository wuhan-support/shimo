package shimo

import (
	"io/ioutil"
	"testing"
)

func TestDocument_GetJSON(t *testing.T) {
	d := NewDocument("6c6GKvX83hRCVdG8", "shimo_gatedlaunch=5; shimo_kong=1; shimo_svc_edit=7958; shimo_sid=s%3Az53puugT5iYN0OBxa0e2cvIuvberbzBg.Er0zR1skqDhb37OEwqgNIDgxR7qfDvkyzfDLLd300dI; sensorsdata2015session=%7B%7D; deviceId=2b663e32-b56d-4079-8194-a0675166f9d1; deviceIdGenerateTime=1580009106800; _csrf=a8gC7RFXooqS5g_s1hnu1vVW; sensorsdata2015jssdkcross=%7B%22distinct_id%22%3A%229998849%22%2C%22%24device_id%22%3A%2216fdd6cf4eeb44-06c60c78fd6141-39607b0f-2073600-16fdd6cf4efcd9%22%2C%22props%22%3A%7B%22%24latest_traffic_source_type%22%3A%22%E7%9B%B4%E6%8E%A5%E6%B5%81%E9%87%8F%22%2C%22%24latest_referrer%22%3A%22%22%2C%22%24latest_referrer_host%22%3A%22%22%2C%22%24latest_search_keyword%22%3A%22%E6%9C%AA%E5%8F%96%E5%88%B0%E5%80%BC_%E7%9B%B4%E6%8E%A5%E6%89%93%E5%BC%80%22%7D%2C%22first_id%22%3A%2216fdd6cf4eeb44-06c60c78fd6141-39607b0f-2073600-16fdd6cf4efcd9%22%7D")

	d.Suffix = "（"
	data, err := d.GetJSON()
	if err != nil {
		t.Fatalf("failed to get json: %v", err)
	}
	bytes, err := d.GetCSV()
	if err != nil {
		t.Fatalf("failed to get json: %v", err)
	}
	marshalled, err := data.MarshalJSON()
	if err != nil {
		t.Fatalf("failed to marshal json: %v", err)
	}
	err = ioutil.WriteFile("out.json", marshalled, 0644)
	if err != nil {
		t.Fatalf("failed to write json file: %v", err)
	}
	err = ioutil.WriteFile("out.csv", bytes, 0644)
	if err != nil {
		t.Fatalf("failed to write csv file: %v", err)
	}
}
