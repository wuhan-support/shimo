package shimo

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/dchest/uniuri"
	"github.com/tealeg/xlsx"
	"github.com/wuhan-support/shimo/x2j"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

var (
	client = http.Client{
		Timeout: time.Minute,
	}
	sheetIndex int
)

type Document struct {
	GUID           string        `json:"guid"`
	SheetIndex     int           `json:"sheet_index"`

	Cookie string `json:"-"`
	Suffix string `json:"suffix"`

	CacheTTL       time.Duration `json:"cache_ttl"`
	CacheBytes     []byte        `json:"cache_bytes"`
	CacheUpdatedAt time.Time     `json:"cache_updated_at"`
	JSONCache *json.RawMessage
	CSVCache []byte
}

type Response struct {
	DownloadURL string `json:"downloadUrl"`
}

// NewDocument returns a new Document
func NewDocument(guid string, cookie string) *Document {
	return &Document{
		GUID:     guid,
		Cookie: cookie,
		CacheTTL: time.Minute * 30,
	}
}

// update updates a document and writing new data to cache fields
func (d *Document) update() error {
	log.Printf("fetching download url from api...")

	val := url.Values{}
	val.Add("type", "xlsx")
	val.Add("file", d.GUID)
	val.Add("returnJson", "1")
	val.Add("name", uniuri.NewLen(32))
	val.Add("isAsync", "0")
	query := val.Encode()

	u := fmt.Sprintf("https://shimo.im/lizard-api/files/%s/export?%s", d.GUID, query)
	request, err := http.NewRequest("GET", u, nil)
	if err != nil {
		log.Printf("failed to fetch http resource: %v", err)
		return err
	}
	request.Header.Set("Cookie", d.Cookie)
	request.Header.Set("Referer", fmt.Sprintf("https://shimo.im/sheets/%s/MODOC", d.GUID))

	var response *http.Response
	response, err = client.Do(request)
	if err != nil {
		response, err = client.Do(request)
		if err != nil {
			return err
		}
	}
	if response.StatusCode != 200 {
		log.Printf("failed to fetch http resource: status code %v", response.StatusCode)
	}
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	var resp Response
	err = json.Unmarshal(responseBytes, &resp)
	if err != nil {
		return err
	}

	log.Printf("got download url: %s", resp.DownloadURL)
	log.Printf("downloading file...")

	response, err = client.Get(resp.DownloadURL)
	if err != nil {
		return err
	}

	log.Printf("got file. converting...")
	fileBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	d.CacheBytes = fileBytes
	d.CacheUpdatedAt = time.Now()

	d.CSVCache = nil
	d.JSONCache = nil

	return nil
}

// touchCache checks against the timeout and update the document accordingly
func (d *Document) touchCache() error {
	// check cache, and update if it has expired
	if d.CacheBytes == nil || d.CacheUpdatedAt.Add(d.CacheTTL).Before(time.Now()) {
		err := d.update()
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Document) GetJSON() (*json.RawMessage, error) {
	err := d.touchCache()
	if err != nil {
		return nil, err
	}

	if d.JSONCache != nil {
		return d.JSONCache, nil
	}

	xlsxFile, err := xlsx.OpenBinary(d.CacheBytes)
	if err != nil {
		return nil, err
	}

	x := x2j.New()
	x.EliminateSuffix = d.Suffix
	message, err := x.Convert(xlsxFile)
	if err != nil {
		return nil, err
	}
	d.JSONCache = &message
	return d.JSONCache, nil
}

func (d *Document) GetCSV() ([]byte, error) {
	err := d.touchCache()
	if err != nil {
		return nil, err
	}

	if d.CSVCache != nil {
		return d.CSVCache, nil
	}

	xlsxFile, err := xlsx.OpenBinary(d.CacheBytes)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBufferString("")

	cw := csv.NewWriter(buf)
	sheet := xlsxFile.Sheets[sheetIndex]
	var vals []string
	for _, row := range sheet.Rows {
		if row != nil {
			vals = vals[:0]
			for _, cell := range row.Cells {
				str, err := cell.FormattedValue()
				if err != nil {
					vals = append(vals, err.Error())
				}
				vals = append(vals, str)
			}
		}
		cw.Write(vals)
	}
	cw.Flush()
	if err := cw.Error(); err != nil {
		return nil, err
	}

	d.CSVCache = buf.Bytes()

	return buf.Bytes(), nil
}
