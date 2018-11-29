package blog

import (
	"ibenben.org/blog/entity"
	"io/ioutil"
	"fmt"
	"strings"
)

const (
	MDMsg           = "---"     //相关信息（包括title、date、tags等）
	MDStrong        = "**"      //加粗
	MDItalic        = "*"       //斜体
	MDStrikeout     = "~~"      //删除的样式
	MDHeaderFirst   = " # "     //一级标题
	MDHeaderSecond  = " ## "    //二级标题
	MDList          = ". "      //有序列表（数字加.）
	MDListUnordered = " - "     //无序列表
	MDTodo          = " - [ ] " //待办事件
	MDQuote         = "> "      //引用
	MDCode          = "```"     //代码块
	MDTableRow      = "|"       //表格行分隔
	MDTableSeparate = " -- "    //表格标题与内容分割
	MDNewline       = "\r\n"    //换行
)

type markdown struct {
}

var m *markdown

func MarkDown() *markdown {
	if m == nil {
		m = new(markdown)
	}

	return m
}

func (this *markdown) CreateMDFile(article *entity.Article) {
	f := entity.ArticleDst + article.Id + ".md"
	data := []byte(this.addHeadMsg(article) + article.Context)
	ioutil.WriteFile(f, data, 0644)
}

func (this *markdown) addHeadMsg(article *entity.Article) string {
	head := MDMsg + MDNewline
	head += "layout: post" + MDNewline
	head += "title: " + article.Title + MDNewline
	head += "date: " + article.Time + MDNewline
	head += "comments: true" + MDNewline
	head += "categories: " + MDNewline
	for _, m := range article.Category {
		head += MDListUnordered + m + MDNewline
	}

	head += "tags: " + MDNewline
	for _, m := range article.Tags {
		head += MDListUnordered + m + MDNewline
	}

	head += MDMsg + MDNewline

	return head
}

func (this *markdown) Common(msg string) string {
	return MDNewline + strings.TrimSpace(msg) + MDNewline
}

func (this *markdown) Strikeout(msg string) string {
	return this.betweenMD(msg, MDStrikeout)
}

func (this *markdown) Italic(msg string) string {
	return this.betweenMD(msg, MDItalic)
}

func (this *markdown) Strong(msg string) string {
	return this.betweenMD(msg, MDStrong)
}

func (this *markdown) H1(msg string) string {
	if msg == "" {
		return msg
	}
	return this.betweenMD(msg, MDHeaderFirst)
}

func (this *markdown) H2(msg string) string {
	if msg == "" {
		return msg
	}
	return this.betweenMD(msg, MDHeaderSecond)
}

func (this *markdown) Image(url, alt string) string {
	return "![" + alt + "](" + url + ")" + MDNewline
}

func (this *markdown) Code(msg string) string {
	return this.betweenMD(msg, MDCode+MDNewline)
}
func (this *markdown) Quote(msg string) string {
	return this.afterMD(msg, MDQuote)
}

func (this *markdown) TableHead(msg []string) string {
	head := MDTableRow
	separate := MDTableRow
	for _, m := range msg {
		head += m + MDTableRow
		separate += MDTableSeparate + MDTableRow
	}
	return head + MDNewline + separate + MDNewline
}

func (this *markdown) TableRow(msg []string) string {
	row := MDTableRow
	for _, m := range msg {
		row += m + MDTableRow
	}

	return row + MDNewline
}

func (this *markdown) Ul(arr []string) string {
	ul := MDNewline
	for _, m := range arr {
		ul += this.afterMD(m, MDListUnordered)
	}

	return ul + MDNewline
}

func (this *markdown) Todo(arr []string) string {
	ul := MDNewline
	for _, m := range arr {
		ul += this.afterMD(m, MDTodo)
	}

	return ul + MDNewline
}

func (this *markdown) Ol(arr []string) string {
	ul := ""
	for i, m := range arr {
		md := fmt.Sprintf("%d%s", 1+i, MDList)
		ul += this.afterMD(m, md)
	}

	return ul + MDNewline + MDNewline + MDNewline
}

func (this *markdown) betweenMD(msg, md string) string {
	return md + msg + md + MDNewline + MDNewline + MDNewline
}

func (this *markdown) afterMD(msg, md string) string {

	return md + msg + MDNewline
}
