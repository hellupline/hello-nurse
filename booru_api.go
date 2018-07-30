package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type (
	BooruDomainURL func(*BooruTag) *url.URL // nolint

	BooruTag struct { // nolint
		Domain string `xml:"-"`
		Name   string `xml:"-"`
		Page   int    `xml:"-"`
		Limit  int    `xml:"-"`

		Count int          `xml:"count,attr"`
		Posts []*BooruPost `xml:"post"`

		Logger *log.Entry `xml:"-"`
	}
	BooruPost struct { // nolint
		PreviewURL string `xml:"preview_url,attr"`
		SampleURL  string `xml:"sample_url,attr"`
		FileURL    string `xml:"file_url,attr"`

		Key     string `xml:"id,attr"`
		RawTags string `xml:"tags,attr"`
	}
)

var (
	domainURLs = map[string]BooruDomainURL{
		"konachan.net": KonachanURLBuilder,
	}
)

func KonachanURLBuilder(tag *BooruTag) *url.URL { // nolint
	return &url.URL{
		RawQuery: url.Values{
			"limit": []string{strconv.Itoa(tag.Limit)},
			"page":  []string{strconv.Itoa(tag.Page)},
			"tags":  []string{tag.Name},
		}.Encode(),
		Scheme: "https",
		Host:   "konachan.net",
		Path:   "/post.xml",
	}
}

func NewBooruTag(domain, name string, page int) *BooruTag { // nolint
	return &BooruTag{
		Domain: domain,
		Name:   name,
		Page:   page,
		Limit:  100,

		Logger: log.WithFields(log.Fields{
			"domain": domain,
			"name":   name,
			"page":   page,
			"limit":  100,
		}),
	}

}

func (t *BooruTag) Fetch() error { // nolint
	filename := t.Filename()
	body, err := ioutil.ReadFile(filename)

	if err != nil {
		u := t.URL().String()
		resp, err := http.Get(u)
		if err != nil {
			return err
		}
		defer resp.Body.Close() // nolint: errcheck

		// XXX: use io copy
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(filename, body, 0644)
		if err != nil {
			return err
		}
	}

	return xml.Unmarshal(body, t)
}

func (t *BooruTag) Filename() string { // nolint
	base := filepath.Join(baseDir, "cache", t.Domain, t.Name)
	_ = os.MkdirAll(base, 755)

	return filepath.Join(base, fmt.Sprintf("%04d.xml", t.Page))
}

func (t *BooruTag) URL() *url.URL { // nolint
	return domainURLs[t.Domain](t)
}

func (t *BooruTag) Pages() int { // nolint
	return (t.Count / t.Limit) + 1
}

func (p *BooruPost) Tags() []string { // nolint
	return strings.Split(strings.TrimSpace(p.RawTags), " ")
}
