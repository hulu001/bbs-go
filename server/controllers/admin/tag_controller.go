package admin

import (
	"github.com/kataras/iris"
	"github.com/mlogclub/bbs-go/model"
	"github.com/mlogclub/bbs-go/services"
	"github.com/mlogclub/simple"
	"strconv"
)

type TagController struct {
	Ctx iris.Context
}

func (this *TagController) GetBy(id int64) *simple.JsonResult {
	t := services.TagService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return simple.JsonData(t)
}

func (this *TagController) AnyList() *simple.JsonResult {
	list, paging := services.TagService.FindPageByParams(simple.NewQueryParams(this.Ctx).
		LikeByReq("name").
		EqByReq("status").
		PageByReq().Desc("id"))
	return simple.JsonData(&simple.PageResult{Results: list, Page: paging})
}

func (this *TagController) PostCreate() *simple.JsonResult {
	t := &model.Tag{}
	err := this.Ctx.ReadForm(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	if len(t.Name) == 0 {
		return simple.JsonErrorMsg("name is required")
	}
	if services.TagService.GetByName(t.Name) != nil {
		return simple.JsonErrorMsg("标签「" + t.Name + "」已存在")
	}

	t.Status = model.TagStatusOk
	t.CreateTime = simple.NowTimestamp()
	t.UpdateTime = simple.NowTimestamp()

	err = services.TagService.Create(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (this *TagController) PostUpdate() *simple.JsonResult {
	id, err := simple.FormValueInt64(this.Ctx, "id")
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	t := services.TagService.Get(id)
	if t == nil {
		return simple.JsonErrorMsg("entity not found")
	}

	err = this.Ctx.ReadForm(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}

	if len(t.Name) == 0 {
		return simple.JsonErrorMsg("name is required")
	}
	if tmp := services.TagService.GetByName(t.Name); tmp != nil && tmp.Id != id {
		return simple.JsonErrorMsg("标签「" + t.Name + "」已存在")
	}

	t.UpdateTime = simple.NowTimestamp()
	err = services.TagService.Update(t)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	return simple.JsonData(t)
}

func (this *TagController) AnyListAll() *simple.JsonResult {
	categoryId, err := strconv.ParseInt(this.Ctx.FormValue("categoryId"), 10, 64)
	if err != nil {
		return simple.JsonErrorMsg(err.Error())
	}
	if categoryId < 0 {
		return simple.JsonErrorMsg("请指定categoryId")
	}
	list, err := services.TagService.ListAll(categoryId)
	if err != nil {
		return simple.JsonData([]interface{}{})
	}
	return simple.JsonData(list)
}

// 标签数据级联选择器
func (this *TagController) GetCascader() *simple.JsonResult {
	categories, err := services.CategoryService.GetCategories()
	if err != nil {
		return simple.JsonErrorMsg("数据加载失败")
	}

	var results []map[string]interface{}

	for _, cat := range categories {
		tags, err := services.TagService.ListAll(cat.Id)
		if err != nil || len(tags) == 0 {
			continue
		}

		var tagOptions []map[string]interface{}
		for _, tag := range tags {
			tagOption := make(map[string]interface{})
			tagOption["value"] = tag.Id
			tagOption["label"] = tag.Name
			tagOptions = append(tagOptions, tagOption)
		}

		option := make(map[string]interface{})
		option["value"] = cat.Id
		option["label"] = cat.Name
		option["children"] = tagOptions

		results = append(results, option)
	}

	return simple.JsonData(results)

}
