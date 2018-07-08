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
	TagPage struct {
		Domain string `xml:"-"`
		Name   string `xml:"-"`
		Page   int    `xml:"-"`
		Limit  int    `xml:"-"`

		Count int         `xml:"count,attr"`
		Posts []*PostPage `xml:"post"`

		Logger *log.Entry `xml:"-"`
	}
	PostPage struct {
		PreviewURL string `xml:"preview_url,attr"`
		SampleURL  string `xml:"sample_url,attr"`
		FileURL    string `xml:"file_url,attr"`

		Key     string `xml:"id,attr"`
		RawTags string `xml:"tags,attr"`
	}

	DomainURL func(*TagPage) *url.URL
)

var (
	domainURLs = map[string]DomainURL{
		"konachan.net": KonachanURLBuilder,
	}
)

func NewTagPage(domain, name string) *TagPage {
	return &TagPage{
		Domain: domain,
		Name:   name,
		Page:   0,
		Limit:  100,

		Logger: log.WithFields(log.Fields{
			"domain": domain,
			"name":   name,
			"page":   0,
			"limit":  100,
		}),
	}

}

func (t *TagPage) Fetch() error {
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

func (t *TagPage) NextPage(page int) *TagPage {
	return &TagPage{
		Domain: t.Domain,
		Name:   t.Name,
		Page:   page,
		Limit:  t.Limit,

		Logger: log.WithFields(log.Fields{
			"domain": t.Domain,
			"name":   t.Name,
			"page":   page,
			"limit":  t.Limit,
		}),
	}

}

func (t *TagPage) Filename() string {
	base := filepath.Join(baseDir, "cache", t.Domain, t.Name)
	os.MkdirAll(base, 755)

	return filepath.Join(base, fmt.Sprintf("%04d.xml", t.Page))
}

func (t *TagPage) URL() *url.URL {
	return domainURLs[t.Domain](t)
}

func (t *TagPage) Pages() int {
	return (t.Count / t.Limit) + 1
}

func (p *PostPage) Tags() []string {
	return strings.Split(strings.TrimSpace(p.RawTags), " ")
}

func KonachanURLBuilder(tag *TagPage) *url.URL {
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
