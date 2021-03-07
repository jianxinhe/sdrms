package controllers

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jianxinhe/sdrms/enums"
	"github.com/jianxinhe/sdrms/models"

	"github.com/astaxie/beego/orm"
)

type CourseWareController struct {
	BaseController
	tree []*models.Resource
}

func (c *CourseWareController) Prepare() {
	//先执行
	c.BaseController.Prepare()
	//如果一个Controller的少数Action需要权限控制，则将验证放到需要控制的Action里
	//c.checkAuthor("TreeGrid", "UserMenuTree", "ParentTreeGrid", "Select")
	//如果一个Controller的所有Action都需要登录验证，则将验证放到Prepare
	//这里注释了权限控制，因此这里需要登录验证
	//c.checkLogin()
	tree := models.GetTreeGrid()
	//转换UrlFor 2 LinkUrl
	c.UrlFor2Link(tree)
	c.tree = tree
}
func (c *CourseWareController) Index() {
	//需要权限控制
	c.checkAuthor()
	//将页面左边菜单的某项激活
	c.Data["activeSidebarUrl"] = c.URLFor(c.controllerName + "." + c.actionName)
	c.setTpl()
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "courseware/index_headcssjs.html"
	c.LayoutSections["footerjs"] = "courseware/index_footerjs.html"
	//页面里按钮权限控制
	c.Data["canEdit"] = c.checkActionAuthor("CourseWareController", "Edit")
	c.Data["canDelete"] = c.checkActionAuthor("CourseWareController", "Delete")
}

//TreeGrid 获取文件夹中所有资源的列表
func (c *CourseWareController) TreeGrid() {
	tree := models.GetTreeGrid()
	//转换UrlFor 2 LinkUrl
	c.UrlFor2Link(tree)
	c.tree = tree
	c.jsonResult(enums.JRCodeSucc, "", tree)
}

//UserMenuTree 获取用户有权管理的菜单、区域列表
func (c *CourseWareController) UserMenuTree() {
	userid := c.curUser.Id
	//获取用户有权管理的菜单列表（包括区域）
	tree := models.ResourceTreeGridByUserId(userid, 1)
	//转换UrlFor 2 LinkUrl
	c.UrlFor2Link(tree)
	c.jsonResult(enums.JRCodeSucc, "", tree)
}

//ParentTreeGrid 获取可以成为某节点的父节点列表
func (c *CourseWareController) ParentTreeGrid() {
	Id, _ := c.GetInt("id", 0)
	tree := models.ResourceTreeGrid5Parent(Id)
	//转换UrlFor 2 LinkUrl
	c.UrlFor2Link(tree)
	c.jsonResult(enums.JRCodeSucc, "", tree)
}

// UrlFor2LinkOne 使用URLFor方法，将资源表里的UrlFor值转成LinkUrl
func (c *CourseWareController) UrlFor2LinkOne(urlfor string) string {
	return urlfor
}

//UrlFor2Link 使用URLFor方法，批量将资源表里的UrlFor值转成LinkUrl
func (c *CourseWareController) UrlFor2Link(src []*models.Resource) {
	for _, item := range src {
		item.LinkUrl = c.UrlFor2LinkOne(item.UrlFor)
	}
}

//Edit 资源编辑页面
func (c *CourseWareController) Edit() {
	//需要权限控制
	c.checkAuthor()
	//如果是POST请求，则由Save处理
	if c.Ctx.Request.Method == "POST" {
		c.Save()
	}
	Id, _ := c.GetInt(":id", 0)
	m := &models.Resource{}
	var err error
	if Id == 0 {
		m.Seq = 100
	} else {
		m, err = models.ResourceOne(Id)
		if err != nil {
			c.pageError("数据无效，请刷新后重试")
		}
	}
	if m.Parent != nil {
		c.Data["parent"] = m.Parent.Id
	} else {
		c.Data["parent"] = 0
	}
	//获取可以成为当前节点的父节点的列表
	c.Data["parents"] = models.ResourceTreeGrid5Parent(Id)
	//转换地址
	m.LinkUrl = c.UrlFor2LinkOne(m.UrlFor)
	c.Data["m"] = m
	if m.Parent != nil {
		c.Data["parent"] = m.Parent.Id
	} else {
		c.Data["parent"] = 0
	}

	c.setTpl("courseware/edit.html", "shared/layout_pullbox.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["footerjs"] = "courseware/edit_footerjs.html"
}

//Save 资源添加编辑 保存
func (c *CourseWareController) Save() {
	var err error
	o := orm.NewOrm()
	parent := &models.Resource{}
	m := models.Resource{}
	parentId, _ := c.GetInt("Parent", 0)
	//获取form里的值
	if err = c.ParseForm(&m); err != nil {
		c.jsonResult(enums.JRCodeFailed, "获取数据失败", m.Id)
	}
	//获取父节点
	if parentId > 0 {
		parent, err = models.ResourceOne(parentId)
		if err == nil && parent != nil {
			m.Parent = parent
		} else {
			c.jsonResult(enums.JRCodeFailed, "父节点无效", "")
		}
	}
	if m.Id == 0 {
		if _, err = o.Insert(&m); err == nil {
			c.jsonResult(enums.JRCodeSucc, "添加成功", m.Id)
		} else {
			c.jsonResult(enums.JRCodeFailed, "添加失败", m.Id)
		}

	} else {
		if _, err = o.Update(&m); err == nil {
			c.jsonResult(enums.JRCodeSucc, "编辑成功", m.Id)
		} else {
			c.jsonResult(enums.JRCodeFailed, "编辑失败", m.Id)
		}

	}
}

// Delete 删除
func (c *CourseWareController) Delete() {
	//需要权限控制
	c.checkAuthor()
	Id, _ := c.GetInt("Id", 0)
	if Id == 0 {
		c.jsonResult(enums.JRCodeFailed, "选择的数据无效", 0)
	}

	fileName := ""
	for _, v := range c.tree {
		if v.Id == Id {
			fileName = v.LinkUrl
			break
		}
	}
	if fileName == "" {
		// 删除的文件或文件夹不存在
		c.jsonResult(enums.JRCodeFailed, "删除的文件或文件夹不存在", 0)
	}

	err := os.Remove(fileName)
	if err != nil {
		// 删除失败
		c.jsonResult(enums.JRCodeFailed, "删除失败", 0)
	} else {
		// 删除成功
		c.jsonResult(enums.JRCodeSucc, fmt.Sprintf("删除成功"), 0)
	}
}

// Select 通用选择面板
func (c *CourseWareController) Select() {
	//获取调用者的类别 1表示 角色
	desttype, _ := c.GetInt("desttype", 0)
	//获取调用者的值
	destval, _ := c.GetInt("destval", 0)
	//返回的资源列表
	var selectedIds []string
	o := orm.NewOrm()
	if desttype > 0 && destval > 0 {
		//如果都大于0,则获取已选择的值，例如：角色，就是获取某个角色已关联的资源列表
		switch desttype {
		case 1:
			{
				role := models.Role{Id: destval}
				o.LoadRelated(&role, "RoleResourceRel")
				for _, item := range role.RoleResourceRel {
					selectedIds = append(selectedIds, strconv.Itoa(item.Resource.Id))
				}
			}
		}
	}
	c.Data["selectedIds"] = strings.Join(selectedIds, ",")
	c.setTpl("courseware/select.html", "shared/layout_pullbox.html")
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["headcssjs"] = "courseware/select_headcssjs.html"
	c.LayoutSections["footerjs"] = "courseware/select_footerjs.html"
}

//CheckUrlFor 填写UrlFor时进行验证
func (c *CourseWareController) CheckUrlFor() {
	urlfor := c.GetString("urlfor")
	link := c.UrlFor2LinkOne(urlfor)
	if len(link) > 0 {
		c.jsonResult(enums.JRCodeSucc, "解析成功", link)
	} else {
		c.jsonResult(enums.JRCodeFailed, "解析失败", link)
	}
}
func (c *CourseWareController) UpdateSeq() {

	Id, _ := c.GetInt("pk", 0)
	oM, err := models.ResourceOne(Id)
	if err != nil || oM == nil {
		c.jsonResult(enums.JRCodeFailed, "选择的数据无效", 0)
	}
	value, _ := c.GetInt("value", 0)
	oM.Seq = value
	if _, err := orm.NewOrm().Update(oM); err == nil {
		c.jsonResult(enums.JRCodeSucc, "修改成功", oM.Id)
	} else {
		c.jsonResult(enums.JRCodeFailed, "修改失败", oM.Id)
	}
}
