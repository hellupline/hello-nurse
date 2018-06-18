package main

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/hellupline/hello-nurse/nursetags"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
)

const uploadURL = "http://localhost:8080/v1/dataset/upload/gob"

type (
	PipelineCallback func(*TagPage, chan<- *TagPage)
	DomainURL        func(*TagPage) *url.URL

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
		URL     string `xml:"file_url,attr"`
		Key     string `xml:"id,attr"`
		RawTags string `xml:"tags,attr"`
	}
)

var (
	domainURLs = map[string]DomainURL{
		"konachan.net": func(tag *TagPage) *url.URL {
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
		},
	}
	basePath string
)

func init() {
	log.SetLevel(log.InfoLevel)
	log.SetOutput(os.Stderr)

	home, err := homedir.Dir()
	if err != nil {
		log.WithError(err).Fatal()
	}
	basePath = filepath.Join(home, ".booru")
}

func main() {
	stage1 := Pipeline(TagGenerator(), func(tag *TagPage, out chan<- *TagPage) {
		tag.Logger.Info("fetch started")
		if err := tag.Fetch(); err != nil {
			tag.Logger.WithField("url", tag.URL()).WithError(err).Error("Failed to fetch TagPage")
			return
		}

		pages := tag.Pages()
		tag.Logger.WithFields(log.Fields{
			"count": tag.Count,
			"pages": pages,
		}).Info("fetch done")

		for i := 0; i < pages; i++ {
			out <- tag.NextPage(i)
		}
	})
	stage2 := Pipeline(stage1, func(tag *TagPage, out chan<- *TagPage) {
		tag.Logger.Info("fetch started")
		if err := tag.Fetch(); err != nil {
			tag.Logger.WithField("url", tag.URL()).WithError(err).Error("Failed to fetch TagPage")
			return
		}
		tag.Logger.WithFields(log.Fields{
			"count": tag.Count,
			"pages": tag.Pages(),
		}).Info("fetch done")
		out <- tag
	})
	stage3 := Pipeline(stage2, func(tag *TagPage, out chan<- *TagPage) {
		for _, post := range tag.Posts {
			nursetags.DefaultDatabase.PostCreate(nursetags.PostData{
				PostKey: nursetags.PostKey{
					Namespace: tag.Domain,
					ID:        post.Key,
				},
				Tags:  post.Tags(),
				Value: "https:" + post.URL,
			})
		}
		out <- tag
	})

	// wait all jobs
	for range stage3 {
	}

	buffer := &bytes.Buffer{}
	if err := nursetags.DefaultDatabase.Write(buffer); err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", uploadURL, buffer)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := http.DefaultClient.Do(req); err != nil {
		log.Fatal(err)
	}
}

func TagGenerator() chan *TagPage {
	// tags := ReadTagsFile()
	tags := []*TagPage{
		NewTagPage("konachan.net", "landscape"),
		NewTagPage("konachan.net", "moon"),
		NewTagPage("konachan.net", "night"),
		NewTagPage("konachan.net", "scenic"),
		NewTagPage("konachan.net", "sky"),
		NewTagPage("konachan.net", "star"),
		NewTagPage("konachan.net", "sunset"),
	}

	out := make(chan *TagPage, len(tags))
	go func() {
		defer close(out)
		for _, tag := range tags {
			out <- tag
		}
	}()
	return out
}

func ReadTagsFile() []*TagPage {
	tags := make([]*TagPage, 0)

	f, err := os.Open("./tags.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close() // nolint: errcheck

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		tags = append(tags, NewTagPage("konachan.net", scanner.Text()))
	}
	return tags
}

func Pipeline(in <-chan *TagPage, cb PipelineCallback) <-chan *TagPage {
	out := make(chan *TagPage, 1)

	wg := sync.WaitGroup{}
	worker := func() {
		defer wg.Done()
		for tag := range in {
			cb(tag, out)
		}
	}
	wg.Add(4)
	for i := 0; i < 4; i++ {
		go worker()
	}
	go func() {
		defer close(out)
		wg.Wait()
	}()

	return out
}

// TagPage API

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
	if t.Page < 1 {
		t.CreateCacheDir()
	}
	filename := t.Filename()
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		u := t.URL().String()
		resp, err := http.Get(u)
		if err != nil {
			return err
		}
		defer resp.Body.Close() // nolint: errcheck

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(filename, body, 0660)
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
	prefix := filepath.Join(basePath, "cache", t.Domain, t.Name)
	return filepath.Join(prefix, fmt.Sprintf("%04d.xml", t.Page))
}

func (t *TagPage) CreateCacheDir() {
	prefix := filepath.Join(basePath, "cache", t.Domain, t.Name)
	os.MkdirAll(prefix, os.ModePerm)
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
