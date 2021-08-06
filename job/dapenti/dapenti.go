package dapenti

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

func Run(out string) error {
	// 判断输出目录是否存在
	if _, err := os.Stat(out); err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		if err := os.Mkdir(out, os.ModePerm); err != nil {
			return err
		}
	}

	rssPath := path.Join(out, "dapenti.xml")

	newRss, err := downloadRss()
	if err != nil {
		return err
	}

	newRss.Channel.Title = "喷嚏图卦"

	// 把“喷嚏图卦”类型的文章摘取出来
	rssItem := []DapentiRssItem{}
	for _, item := range newRss.Channel.Items {
		if strings.Contains(item.Title, "喷嚏图卦") {
			rssItem = append(rssItem, item)
		}
	}
	newRss.Channel.Items = rssItem

	oldRss, err := readRss(rssPath)
	if err != nil {
		return err
	}
	if oldRss == nil {
		return writeRss(newRss, rssPath)
	}

	for _, oldItem := range oldRss.Channel.Items {
		isSame := false

		for _, newItem := range newRss.Channel.Items {
			if newItem.Title == oldItem.Title {
				isSame = true
				break
			}
		}

		if !isSame {
			rssItem = append(rssItem, oldItem)
		}
	}
	newRss.Channel.Items = rssItem

	return writeRss(newRss, rssPath)
}

type DapentiRss struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Content string   `xml:"content,attr"`
	Dc      string   `xml:"dc,attr"`
	Channel struct {
		Text        string           `xml:",chardata"`
		Language    string           `xml:"language"`
		Copyright   string           `xml:"copyright"`
		Generator   string           `xml:"generator"`
		WebMaster   string           `xml:"webMaster"`
		Title       string           `xml:"title"`
		Link        string           `xml:"link"`
		Description string           `xml:"description"`
		Items       []DapentiRssItem `xml:"item"`
	} `xml:"channel"`
}

type DapentiRssItem struct {
	Text        string `xml:",chardata"`
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Author      string `xml:"author"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
}

func toRss(txt []byte) (*DapentiRss, error) {
	rss := &DapentiRss{}
	if err := xml.Unmarshal(txt, rss); err != nil {
		return nil, err
	}

	return rss, nil
}

func downloadRss() (*DapentiRss, error) {
	rep, err := http.Get("https://www.dapenti.com/blog/rss2.asp")
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		return nil, err
	}

	if rep.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("下载喷嚏网 RSS 失败，错误：%v", string(body))
	}

	return toRss(body)
}

func readRss(rssPath string) (*DapentiRss, error) {
	if _, err := os.Stat(rssPath); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	txt, err := ioutil.ReadFile(rssPath)
	if err != nil {
		return nil, err
	}

	return toRss(txt)
}

func writeRss(rss *DapentiRss, rsspath string) error {
	// 不保存超过 10 条，以免文件太大
	if len(rss.Channel.Items) > 10 {
		rss.Channel.Items = rss.Channel.Items[:10]
	}

	txt, err := xml.Marshal(rss)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(rsspath, txt, os.ModePerm); err != nil {
		return err
	}

	return nil
}
