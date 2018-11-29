package csdn

import (
	"github.com/PuerkitoBio/goquery"
	"ibenben.org/blog/entity"
	"ibenben.org/util"
	"ibenben.org/blog"
	"strings"
	"io/ioutil"
	"os"
	"net/http"
)

//
// 获取器
//

type Acquirer struct {
	doc     *goquery.Document
	article *entity.Article
}

func NewAcquirer(url string) *Acquirer {
	obj := new(Acquirer)
	obj.init(url)
	return obj
}

func (this *Acquirer) init(url string) {
	//fmt.Println(url)
	this.doc = Doc(url)
	this.article = new(entity.Article)

	this.article.Id = strings.Split(url, "article/details/")[1]
}

func (this *Acquirer) Parse() {
	this.parseHeader()

	content := this.doc.Find("#content_views")
	this.parseContent(content)

	blog.MarkDown().CreateMDFile(this.article)
}


//解析文章内容
func (this *Acquirer) parseContent(content *goquery.Selection) {
	content.Children().Each(func(i int, s *goquery.Selection) {

		if s.Is("h1") {
			this.article.Context += blog.MarkDown().H2(s.Text())
		} else if s.Is("h2") {
			this.article.Context += blog.MarkDown().H2(s.Text())
		} else if s.Is("h3") {
			this.article.Context += blog.MarkDown().Strong(s.Text())
		} else if s.Is("p") {
			this.article.Context += this.parsePhase(s)
		} else if s.Is("pre") {
			if len(s.Find("code").Nodes) > 0 {
				this.article.Context += blog.MarkDown().Code(s.Text())
			} else {
				this.article.Context += blog.MarkDown().Common(s.Text())
			}
		} else if s.Is("table") {
			this.article.Context += this.parseTable(s)
		} else if s.Is("ol") {
			this.article.Context += this.parseOl(s)
		} else if s.Is("ul") {
			this.article.Context += this.parseUl(s)
		} else if s.Is("blockquote") {
			p := s.Find("p")
			this.article.Context += blog.MarkDown().Quote(p.Text())
		} else if s.Is("div") {
			this.parseContent(s)
		} else {
			this.article.Context += blog.MarkDown().Common(s.Text())
		}
	})
}

//解析有序列表
func (this *Acquirer) parseOl(parent *goquery.Selection) string {
	if parent.Children().Length() == 1 {
		c := parent.Find("li").First()
		start := "1"
		if v, ok := parent.Attr("start"); ok {
			start = v
		}
		msg := start + ". " + strings.TrimSpace(c.Text())
		return blog.MarkDown().Strong(msg)

	}

	arr := make([]string, 0)

	parent.Children().Each(func(i int, c *goquery.Selection) {

		msg := strings.TrimSpace(c.Text())

		arr = append(arr, msg)
	})
	return blog.MarkDown().Ol(arr)
}

//解析无序列表
func (this *Acquirer) parseUl(parent *goquery.Selection) string {
	arr := make([]string, 0)
	isTodo := false
	parent.Children().Each(func(i int, c *goquery.Selection) {
		isTodo = c.HasClass("task-list-item")
		msg := strings.TrimSpace(c.Text())
		arr = append(arr, msg)
	})
	if isTodo {
		return blog.MarkDown().Todo(arr)
	}
	return blog.MarkDown().Ul(arr)
}

//解析段落
func (this *Acquirer) parsePhase(parent *goquery.Selection) string {
	context := ""
	parent.Children().Each(func(i int, c *goquery.Selection) {
		if c.Is("strong") {
			context += blog.MarkDown().Strong(c.Text())
		} else if c.Is("em") {
			context += blog.MarkDown().Italic(c.Text())
		} else if c.Is("s") {
			context += blog.MarkDown().Strikeout(c.Text())
		} else if c.Is("img") {
			url, _ := c.Attr("src")
			alt, _ := c.Attr("alt")
			path := this.saveImage(url)
			context += blog.MarkDown().Common(parent.Text())
			context += blog.MarkDown().Image(path, alt)

		}
	})

	if context == "" {
		context = blog.MarkDown().Common(parent.Text())
	}

	return context
}
func (this *Acquirer) parseTable(table *goquery.Selection) string {
	content := ""
	isNotHeader := true
	arr := make([]string, 0)
	table.Find("thead").Find("tr").Children().Each(func(i int, s *goquery.Selection) {
		arr = append(arr, strings.TrimSpace(s.Text()))
	})
	if len(arr) > 0 {
		isNotHeader = false
		content += blog.MarkDown().TableHead(arr)
	}

	table.Find("tbody").Children().Each(func(i int, s *goquery.Selection) {
		arr = make([]string, 0)

		s.Find("td").Each(func(i int, t *goquery.Selection) {
			msg := strings.TrimSpace(t.Text())
			if msg == "" {
				msg = " - "
			}
			arr = append(arr, msg)
		})

		if isNotHeader && i == 0 {
			content += blog.MarkDown().TableHead(arr)
			return
		}

		content += blog.MarkDown().TableRow(arr)
	})

	return content
}

//解析头部信息
func (this *Acquirer) parseHeader() {
	header := this.doc.Find(".article-header-box")

	s := header.Find(".title-article")
	this.article.Title = s.Text()

	t := header.Find(".time")
	this.article.Time = util.Chinese2Base(t.Text())

	header.Find(".tags-box").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			//tags
			this.article.Tags = this.parseAs(s)
		} else {
			this.article.Category = this.parseAs(s)
		}
	})
}

//遍历节点下的a子节点，返回子节点的文本
func (this *Acquirer) parseAs(parent *goquery.Selection) []string {
	arr := make([]string, 1)
	parent.Find("a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		v := strings.TrimSpace(s.Text())
		arr = append(arr, v)
	})

	return arr
}

//下载图片
func (this *Acquirer) saveImage(url string) (imgPath string) {

	url = strings.Split(url, "?")[0]
	//log.Println(url)

	//去掉最左边的'/'
	imgName := strings.Replace(url, "https://img-blog.csdn.net/", "", -1) + ".png"
	filename := entity.ArticleImageDst + imgName

	response, _ := http.Get(url)

	defer response.Body.Close()

	data, _ := ioutil.ReadAll(response.Body)

	image, _ := os.Create(filename)

	defer image.Close()
	image.Write(data)

	imgPath = entity.ArticleImageUrl + imgName
	return
}
