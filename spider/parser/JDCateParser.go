package parser

import (
	"JD-price-tracker/model"
	"regexp"
)

func ParseJDCategory(bodyText []byte) []model.Category {
	var cateList []model.Category
	re := regexp.MustCompile(`"list.jd.com/list.html\?cat=(.*?)"`)
	matches := re.FindAllSubmatch(bodyText, -1)
	for _, match := range matches {
		cate := model.Category{
			Site: "www.jd.com",
			Url: "https://" + string(match[0]),
			CategoryId: string(match[1]),
		}
		cateList = append(cateList, cate)
	}
	return cateList
}
