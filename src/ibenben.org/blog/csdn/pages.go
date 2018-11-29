package csdn

import (
	"github.com/PuerkitoBio/goquery"
	"ibenben.org/blog/entity"
	"time"
	"ibenben.org/util"
	"fmt"
)

type pages struct {
	doc      *goquery.Document
	isNoDate bool
}

func Pages() *pages {
	return new(pages)
}

func (this *pages) QueryAndParse() {
	page := 1
	for {
		go this.parseList(this.listUrl(page))
		time.Sleep(5 * time.Second)
		if this.isNoDate {
			break
		}
		page++
	}
}

func (this *pages) parseList(url string) {
	defer util.CatchErr()
	fmt.Println(url)
	ids := this.getIds(url)
	if len(ids) == 0 {
		this.isNoDate = true
	}

	for _, id := range ids {
		u := this.articleUrl(id)

		go this.parseDetail(u)
	}

}

func (this *pages) parseDetail(url string) {
	defer util.CatchErr()
	NewAcquirer(url).Parse()
}

func (this *pages) listUrl(page int) string {
	return fmt.Sprintf("%s%s%d", entity.Root, entity.ArticleList, page)
}

func (this *pages) articleUrl(id string) string {
	return entity.Root + entity.ArticleDetails + id
}

func (this *pages) getIds(url string) []string {
	doc := Doc(url)

	ids := make([]string, 0)

	doc.Find(".article-list").Children().Each(func(i int, s *goquery.Selection) {

		style, _ := s.Attr("style")
		if style == "display: none;" {
			return
		}
		id, _ := s.Attr("data-articleid")
		if len(id) == 0 {
			return
		}

		ids = append(ids, id)
	})
	//article-list

	return ids
}
