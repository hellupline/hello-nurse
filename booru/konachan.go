package booru

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
)

func init() {
	RegisterClient(NewKonachanClient, "konachan.net")
}

type (
	KonachanClient struct { // nolint: golint
		CacheDir string
		Limit    int
	}

	KonachanTag struct { // nolint: golint
		Count    int            `xml:"count,attr"`
		RawPosts []KonachanPost `xml:"posts"`

		TagName string `xml:"-"`
		Page    int    `xml:"-"`
		Limit   int    `xml:"-"`

		CacheDir string `xml:"-"`
	}

	KonachanPost struct { // nolint: golint
		PreviewURL string `xml:"preview_url,attr"`
		FileURL    string `xml:"file_url,attr"`
		SampleURL  string `xml:"sample_url,attr"`

		ID      string `xml:"id,attr"`
		RawTags string `xml:"tags,attr"`
	}
)

func NewKonachanClient(cacheDir string, limit int) Client { // nolint: golint
	return &KonachanClient{
		CacheDir: cacheDir,
		Limit:    limit,
	}
}

func (c *KonachanClient) NewTag(tagName string, page int) Tag { // nolint: golint
	return &KonachanTag{
		TagName: tagName,
		Page:    page,
		Limit:   c.Limit,
	}
}

func (c *KonachanClient) Name() string { // nolint: golint
	return "konachan.net"
}

func (t *KonachanTag) Fetch() error { // nolint: golint
	filename := t.filename()
	body, err := ioutil.ReadFile(filename)

	if err != nil {
		u := t.url().String()
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

func (t *KonachanTag) filename() string { // nolint: golint
	base := filepath.Join(t.CacheDir, "cache", "konachan.net", t.TagName)
	_ = os.MkdirAll(base, 0750)

	return filepath.Join(base, fmt.Sprintf("%04d.xml", t.Page))
}

func (t *KonachanTag) url() *url.URL { // nolint: golint
	return &url.URL{
		RawQuery: url.Values{
			"limit": []string{strconv.Itoa(t.Limit)},
			"page":  []string{strconv.Itoa(t.Page)},
			"tags":  []string{t.TagName},
		}.Encode(),
		Scheme: "https",
		Host:   "konachan.net",
		Path:   "/post.xml",
	}
}

func (t *KonachanTag) Posts() []Post { // nolint: golint
	posts := make([]Post, len(t.RawPosts))
	for i := range t.RawPosts {
		posts[i] = &t.RawPosts[i]
	}
	return posts
}

func (t *KonachanTag) Pages() int { // nolint: golint
	return (t.Count / t.Limit) + 1
}

func (p *KonachanPost) Tags() []string { // nolint: golint
	return strings.Split(strings.TrimSpace(p.RawTags), " ")
}

func (p *KonachanPost) Key() string { // nolint: golint
	return p.ID
}

func (p *KonachanPost) Body() map[string]string { // nolint: golint
	return map[string]string{
		"preview_url": p.PreviewURL,
		"file_url":    p.FileURL,
		"sample_url":  p.SampleURL,
	}
}
