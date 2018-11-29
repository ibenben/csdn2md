package util

import "time"

const (
	ChineseFormat = "2006年01月02日 15:04:05"
	BaseFormat = "2006-01-02 15:04:05"
	OnlyNumberFormate  = "20060102150405"
)

func Chinese2Base(v string) string{
	t, _ := time.Parse(ChineseFormat, v)

	return t.Format(BaseFormat)
}
