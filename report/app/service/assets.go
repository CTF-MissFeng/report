package service

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"strconv"
	"strings"

	"assets/app/dao"
	"assets/app/model"
	"assets/library/logger"

	_ "github.com/mattn/go-sqlite3"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
)

var Assets = new(serviceAssets)

type serviceAssets struct{}

// ManagerAdd 添加安全管理员
func (s *serviceAssets) ManagerAdd(r *model.RequestAssetsManagerAdd)error{
	count,err := dao.AssetsManager.Where("manager_name=?",r.ManagerName).Count()
	if err != nil{
		logger.WebLog.Warningf("资产管理-添加安全管理员失败:%s", err.Error())
		return errors.New("添加安全管理员失败,数据库错误")
	}
	if count > 0{
		return errors.New("添加安全管理员失败,已存在该安全管理员")
	}
	_,err = dao.AssetsManager.Insert(r)
	if err != nil{
		logger.WebLog.Warningf("资产管理-添加安全管理员失败:%s", err.Error())
		return errors.New("添加安全管理员失败,数据库错误")
	}
	return nil
}

// ManagerDelete 删除安全管理员-将渗透测试报告和业务系统资产的安全管理员一并更新
func (s *serviceAssets) ManagerDelete(r *model.RequestAssetsManagerAdd)error{
	count,err := dao.AssetsManager.Where("manager_name=?",r.ManagerName).Count()
	if err != nil{
		logger.WebLog.Warningf("资产管理-删除安全管理员失败:%s", err.Error())
		return errors.New("删除安全管理员失败,数据库错误")
	}
	if count == 0{
		return errors.New("删除安全管理员失败,该管理员不存在")
	}
	if _,err = dao.AssetsManager.Where("manager_name=?",r.ManagerName).Delete(); err != nil{
		return errors.New(fmt.Sprintf("删除安全管理员失败:%s", err.Error()))
	}
	dao.AssetsReports.Update(g.Map{"manager_name":"已删除"},"manager_name", r.ManagerName)
	dao.AssetsWeb.Update(g.Map{"manager_name":"已删除"},"manager_name", r.ManagerName)
	return nil
}

// SearchManager 安全管理员模糊分页查询
func (s *serviceAssets) SearchManager(page, limit int, search interface{})*model.RsponseManager{
	var result []*model.AssetsManager
	SearchModel := dao.AssetsManager.Clone()
	searchStr := gconv.String(search)
	if search != ""{
		j := gjson.New(searchStr)
		if gconv.String(j.Get("managerName")) != ""{
			SearchModel = SearchModel.Where("manager_name like ?", "%"+gconv.String(j.Get("managerName"))+"%")
		}
	}
	count,_ := SearchModel.Count()
	if page > 0 && limit > 0 {
		err := SearchModel.Order("id desc").Limit((page-1)*limit,limit).Scan(&result)
		if err != nil {
			logger.WebLog.Warningf("安全管理员分页查询 数据库错误:%s", err.Error())
			return &model.RsponseManager{Code:201, Msg:"查询失败,数据库错误", Count:0, Data:nil}
		}
	}else{
		return &model.RsponseManager{Code:201, Msg:"查询失败,分页参数有误", Count:0, Data:nil}
	}
	index := (page-1)*limit
	results := make([]model.RsponseManagerInfo,0)
	for i,_:=range result{
		index++
		result[i].Id = uint(index)
		webCount,_ := dao.AssetsWeb.Where("manager_name=?", result[i].ManagerName).Count() // web资产数
		levelNoCount,_ := dao.AssetsReports.Where("manager_name=?",result[i].ManagerName).Where("level_status=?",2).Count() // 已整改
		levelYesCount,_ := dao.AssetsReports.Where("manager_name=?",result[i].ManagerName).Where("level_status=?",1).Count() // 未整改
		levelCount,_ := dao.AssetsReports.Where("manager_name=?",result[i].ManagerName).Count() //漏洞总数
		results = append(results, model.RsponseManagerInfo{
			Id: index,
			ManagerName: result[i].ManagerName,
			CusTime: result[i].CreateAt,
			WebCount: webCount,
			LevelNoCount: levelNoCount,
			LevelYesCount: levelYesCount,
			LevelCount: levelCount,
		})
	}
	return &model.RsponseManager{Code:0, Msg:"ok", Count:int64(count), Data:results}
}


// GroupAssetsType 返回主机资产厂商组
func (s *serviceAssets) GroupAssetsType(page, limit int, search interface{})*model.ResponseAssetsTypeGroup{
	var result []model.ResponseAssetsTypeGroupInfo
	SearchModel := dao.AssetsType.Clone()
	searchStr := gconv.String(search)
	if searchStr != ""{
		SearchModel = SearchModel.Where("type_name like ?", "%"+searchStr+"%")
	}
	count,_ := SearchModel.Fields("type_name").Group("type_name").Count()
	if page > 0 && limit > 0 {
		err := SearchModel.Limit((page-1)*limit,limit).Fields("type_name").Group("type_name").Scan(&result)
		if err != nil {
			return &model.ResponseAssetsTypeGroup{Code:201, Msg:"查询失败,数据库错误", Count:0, Data:nil}
		}
	}else{
		return &model.ResponseAssetsTypeGroup{Code:201, Msg:"查询失败,分页参数有误", Count:0, Data:nil}
	}
	index := (page-1)*limit
	for i,_:=range result{
		index++
		result[i].ID = index
	}
	return &model.ResponseAssetsTypeGroup{Code:0, Msg:"ok", Count:int64(count), Data:result}
}

// GroupAssetsAttribution 返回主机资产应用系统组
func (s *serviceAssets) GroupAssetsAttribution(page, limit int, search interface{})*model.ResponseAssetsTypeAttributionGroup{
	var result []model.ResponseAssetsTypeAttributionGroupInfo
	SearchModel := dao.AssetsType.Clone()
	searchStr := gconv.String(search)
	if searchStr != ""{
		SearchModel = SearchModel.Where("attribution like ?", "%"+searchStr+"%")
	}
	count,_ := SearchModel.Fields("attribution").Group("attribution").Count()
	if page > 0 && limit > 0 {
		err := SearchModel.Limit((page-1)*limit,limit).Fields("attribution").Group("attribution").Scan(&result)
		if err != nil {
			return &model.ResponseAssetsTypeAttributionGroup{Code:201, Msg:"查询失败,数据库错误", Count:0, Data:nil}
		}
	}else{
		return &model.ResponseAssetsTypeAttributionGroup{Code:201, Msg:"查询失败,分页参数有误", Count:0, Data:nil}
	}
	index := (page-1)*limit
	for i,_:=range result{
		index++
		result[i].ID = index
	}
	return &model.ResponseAssetsTypeAttributionGroup{Code:0, Msg:"ok", Count:int64(count), Data:result}
}

// GroupAssetsDepartment 返回主机资产部门组
func (s *serviceAssets) GroupAssetsDepartment(page, limit int, search interface{})*model.ResponseAssetsTypeDepartmentGroup{
	var result []model.ResponseAssetsTypeDepartmentGroupInfo
	SearchModel := dao.AssetsType.Clone()
	searchStr := gconv.String(search)
	if searchStr != ""{
		SearchModel = SearchModel.Where("department like ?", "%"+searchStr+"%")
	}
	count,_ := SearchModel.Fields("department").Group("department").Count()
	if page > 0 && limit > 0 {
		err := SearchModel.Limit((page-1)*limit,limit).Fields("department").Group("department").Scan(&result)
		if err != nil {
			return &model.ResponseAssetsTypeDepartmentGroup{Code:201, Msg:"查询失败,数据库错误", Count:0, Data:nil}
		}
	}else{
		return &model.ResponseAssetsTypeDepartmentGroup{Code:201, Msg:"查询失败,分页参数有误", Count:0, Data:nil}
	}
	index := (page-1)*limit
	for i,_:=range result{
		index++
		result[i].ID = index
	}
	return &model.ResponseAssetsTypeDepartmentGroup{Code:0, Msg:"ok", Count:int64(count), Data:result}
}

// SearchType 主机资产所属模糊分页查询
func (s *serviceAssets) SearchType(page, limit int, search interface{})*model.ResponseType{
	var result []*model.ResponseTypeInfo
	SearchModel := dao.AssetsType.Clone()
	searchStr := gconv.String(search)
	if search != ""{
		j := gjson.New(searchStr)
		if gconv.String(j.Get("TypeName")) != ""{
			SearchModel = SearchModel.Where("type_name like ?", "%"+gconv.String(j.Get("TypeName"))+"%")
		}
		if gconv.String(j.Get("AttriBution")) != ""{
			SearchModel = SearchModel.Where("attribution like ?", "%"+gconv.String(j.Get("AttriBution"))+"%")
		}
		if gconv.String(j.Get("SubDomain")) != ""{
			SearchModel = SearchModel.Where("subdomain like ?", "%"+gconv.String(j.Get("SubDomain"))+"%")
		}
		if gconv.String(j.Get("PublicIp")) != ""{
			SearchModel = SearchModel.Where("public_ip like ?", "%"+gconv.String(j.Get("PublicIp"))+"%")
		}
		if gconv.String(j.Get("IntranetIp")) != ""{
			SearchModel = SearchModel.Where("intranet_ip like ?", "%"+gconv.String(j.Get("IntranetIp"))+"%")
		}
		if gconv.String(j.Get("Department")) != ""{
			SearchModel = SearchModel.Where("department like ?", "%"+gconv.String(j.Get("Department"))+"%")
		}
	}
	count,_ := SearchModel.Count()
	if page > 0 && limit > 0 {
		err := SearchModel.Order("id desc").Limit((page-1)*limit,limit).Scan(&result)
		if err != nil {
			logger.WebLog.Warningf("主机资产分页查询 数据库错误:%s", err.Error())
			return &model.ResponseType{Code:201, Msg:"查询失败,数据库错误", Count:0, Data:nil}
		}
	}else{
		return &model.ResponseType{Code:201, Msg:"查询失败,分页参数有误", Count:0, Data:nil}
	}
	index := (page-1)*limit
	for i,_:=range result{
		index++
		result[i].ID1 = result[i].Id
		result[i].Id = uint(index)
		result[i].Subdomain = strings.ReplaceAll(result[i].Subdomain,"\n","<br/>")
		result[i].PublicIp = strings.ReplaceAll(result[i].PublicIp,"\n","<br/>")
		result[i].IntranetIp = strings.ReplaceAll(result[i].IntranetIp,"\n","<br/>")
	}
	return &model.ResponseType{Code:0, Msg:"ok", Count:int64(count), Data:result}
}

// TypeAdd 添加修改主机资产-修改主机资产，将关联的渗透测试和业务系统一并更新
func (s *serviceAssets) TypeAdd(r *model.RequestAssetsTypeAdd)error{
	if len(r.TypeID) != 0 {
		result,err := dao.AssetsType.Where("id=?", r.TypeID).FindOne()
		if err != nil {
			logger.WebLog.Warningf("修改主机资产,查询对应应用系统名错误:%s", err.Error())
			return errors.New("修改主机资产失败,该主机ID不存在")
		}
		if _,err := dao.AssetsType.Update(r,"id",r.TypeID); err != nil{
			logger.WebLog.Warningf("资产管理-修改主机资产失败:%s", err.Error())
			return errors.New("修改主机资产失败,数据库错误")
		}
		if result.Attribution != r.AttriBution{ // 用户修改了主机资产中的应用系统名
			count,err:=dao.AssetsWeb.Where("attribution=?", result.Attribution).Count()
			if err != nil {
				logger.WebLog.Warningf("修改主机资产,查询业务系统对应的应用名错误:%s", err.Error())
				return nil
			}else if count > 0{ // 若业务系统中有此应用名，则更新
				if _,err = dao.AssetsWeb.Data(g.Map{"attribution":r.AttriBution}).Where("attribution", result.Attribution).Update();err != nil{
					logger.WebLog.Warningf("修改主机资产,更新业务系统应用名错误:%s", err.Error())
					return nil
				}
			}
			count,err=dao.AssetsReports.Where("attribution=?", result.Attribution).Count()
			if err != nil {
				logger.WebLog.Warningf("修改主机资产,查询渗透测试报告对应的应用名错误:%s", err.Error())
				return nil
			}else if count > 0{ // 若渗透测试报告中有此应用名，则更新
				if _,err = dao.AssetsReports.Data(g.Map{"attribution":r.AttriBution}).Where("attribution", result.Attribution).Update();err != nil{
					logger.WebLog.Warningf("修改主机资产,更新渗透测试报告应用名错误:%s", err.Error())
					return nil
				}
			}
		}
	}else{
		if count,err := dao.AssetsType.Where("attribution=?",r.AttriBution).Count(); err != nil{
			logger.WebLog.Warningf("资产管理-添加主机资产失败:%s", err.Error())
			return errors.New("添加主机资产失败,数据库错误")
		}else if count>0{
			return errors.New("添加主机资产失败,已存在该应用名")
		}
		if _,err := dao.AssetsType.Insert(r); err != nil{
			logger.WebLog.Warningf("资产管理-添加主机资产失败:%s", err.Error())
			return errors.New("添加主机资产失败,数据库错误")
		}
	}
	return nil
}

// TypeDelete 删除主机资产
func (s *serviceAssets) TypeDelete(r *model.RequestAssetsTypeDelete)error{
	count,err := dao.AssetsType.Where("id=?",r.ID).Count()
	if err != nil{
		logger.WebLog.Warningf("资产管理-删除主机资产失败:%s", err.Error())
		return errors.New("删除主机资产失败,数据库错误")
	}
	if count == 0{
		return errors.New("删除主机资产失败,该资产不存在")
	}
	if _,err = dao.AssetsType.Where("id=?",r.ID).Delete(); err != nil{
		return errors.New(fmt.Sprintf("删除主机资产失败:%s", err.Error()))
	}
	return nil
}

// TypeShow 查看指定ID的主机资产信息
func (s *serviceAssets) TypeShow(r *model.RequestAssetsTypeDelete)(*model.AssetsType,error){
	result,err := dao.AssetsType.Where("id=?",r.ID).FindOne()
	if err != nil{
		logger.WebLog.Warningf("资产管理-查看主机资产失败:%s", err.Error())
		return nil,errors.New("查看主机资产失败,数据库错误")
	}
	return result,nil
}

// ExportType 导出主机资产
func (s *serviceAssets) ExportType()(*bytes.Buffer, error){
	var result []model.AssetsType
	err := dao.AssetsType.Scan(&result)
	if err != nil {
		return &bytes.Buffer{},errors.New("导出主机资产数据失败,数据库错误")
	}
	if len(result) == 0{
		return &bytes.Buffer{},errors.New("导出主机资产数据失败,无资产数据")
	}
	xlsx := excelize.NewFile()
	index := xlsx.NewSheet("主机资产表")
	xlsx.DeleteSheet("Sheet1")
	table_name := "主机资产表"
	xlsx.SetCellValue(table_name,"A1","厂商")
	xlsx.SetCellValue(table_name,"B1","应用系统")
	xlsx.SetCellValue(table_name,"C1","所属部门")
	xlsx.SetCellValue(table_name,"D1","管理员")
	xlsx.SetCellValue(table_name,"E1","子域名")
	xlsx.SetCellValue(table_name,"F1","内网IP")
	xlsx.SetCellValue(table_name,"G1","公网IP")
	for i, info := range result {
		xlsx.SetCellValue(table_name, "A" + strconv.Itoa(i+2), info.TypeName)
		xlsx.SetCellValue(table_name, "B" + strconv.Itoa(i+2), info.Attribution)
		xlsx.SetCellValue(table_name, "C" + strconv.Itoa(i+2), info.Department)
		xlsx.SetCellValue(table_name, "D" + strconv.Itoa(i+2), info.AssetsUsername)
		xlsx.SetCellValue(table_name, "E" + strconv.Itoa(i+2), info.Subdomain)
		xlsx.SetCellValue(table_name, "F" + strconv.Itoa(i+2), info.IntranetIp)
		xlsx.SetCellValue(table_name, "G" + strconv.Itoa(i+2), info.PublicIp)
	}
	xlsx.SetColWidth(table_name, "A", "A", 18)
	xlsx.SetColWidth(table_name, "B", "B", 25)
	xlsx.SetColWidth(table_name, "C", "C", 20)
	xlsx.SetColWidth(table_name, "D", "D", 15)
	xlsx.SetColWidth(table_name, "E", "E", 20)
	xlsx.SetColWidth(table_name, "F", "F", 20)
	xlsx.SetColWidth(table_name, "G", "G", 20)
	xlsx.SetActiveSheet(index) // 设置工作簿的默认工作表
	var buf bytes.Buffer
	err = xlsx.Write(&buf)
	if err != nil {
		return &bytes.Buffer{},errors.New("导出主机资产数据失败,xlsx流写入失败")
	}
	return &buf,nil
}

// DownloadTypeTemplate 下载主机资产模板
func (s *serviceAssets) DownloadTypeTemplate()(*bytes.Buffer, error){
	xlsx := excelize.NewFile()
	index := xlsx.NewSheet("主机资产表")
	xlsx.DeleteSheet("Sheet1")
	table_name := "主机资产表"
	xlsx.SetCellValue(table_name,"A1","厂商")
	xlsx.SetCellValue(table_name,"B1","应用系统")
	xlsx.SetCellValue(table_name,"C1","所属部门")
	xlsx.SetCellValue(table_name,"D1","管理员")
	xlsx.SetCellValue(table_name,"E1","子域名")
	xlsx.SetCellValue(table_name,"F1","内网IP")
	xlsx.SetCellValue(table_name,"G1","公网IP")
	xlsx.SetCellValue(table_name, "A" + strconv.Itoa(2), "厂商")
	xlsx.SetCellValue(table_name, "B" + strconv.Itoa(2), "xx系统")
	xlsx.SetCellValue(table_name, "C" + strconv.Itoa(2), "xx部门")
	xlsx.SetCellValue(table_name, "D" + strconv.Itoa(2), "张三")
	xlsx.SetCellValue(table_name, "E" + strconv.Itoa(2), "")
	xlsx.SetCellValue(table_name, "F" + strconv.Itoa(2), "")
	xlsx.SetCellValue(table_name, "G" + strconv.Itoa(2), "")
	xlsx.SetColWidth(table_name, "A", "A", 18)
	xlsx.SetColWidth(table_name, "B", "B", 25)
	xlsx.SetColWidth(table_name, "C", "C", 20)
	xlsx.SetColWidth(table_name, "D", "D", 15)
	xlsx.SetColWidth(table_name, "E", "E", 20)
	xlsx.SetColWidth(table_name, "F", "F", 20)
	xlsx.SetColWidth(table_name, "G", "G", 20)
	xlsx.SetActiveSheet(index)
	var buf bytes.Buffer
	err := xlsx.Write(&buf)
	if err != nil {
		return &bytes.Buffer{},errors.New("下载主机资产表模板错误,xlsx流写入失败")
	}
	return &buf,nil
}

// ImportType 导入主机资产表
func (s *serviceAssets) ImportType(filename string)(string,error){
	xlsx, err := excelize.OpenFile(filename)
	if err != nil {
		logger.WebLog.Warningf("导入主机资产表格失败:%s", err.Error())
		return "",errors.New("导入主机资产失败,不是有效的xlsx文档")
	}
	rows,err := xlsx.GetRows("主机资产表")
	if err != nil {
		logger.WebLog.Warningf("导入主机资产表格失败:%s", err.Error())
		return "",errors.New("导入主机资产失败,不是有效的资产模板")
	}
	var Inserts []*model.AssetsType
	errorCount := 0 // 记录不合格的行数
	for i,row := range rows {
		if i == 0{ // 过滤首行
			continue
		}
		if len(row[0]) == 0{
			errorCount++
			continue
		}
		insert := model.AssetsType{}
		for index,v := range row{
			switch index{
			case 0:
				insert.TypeName = gstr.Trim(v)
			case 1:
				insert.Attribution = gstr.Trim(v)
			case 2:
				insert.Department = gstr.Trim(v)
			case 3:
				insert.AssetsUsername = gstr.Trim(v)
			case 4:
				insert.Subdomain = gstr.Trim(strings.ReplaceAll(v,"/","\n"))
			case 5:
				insert.IntranetIp = gstr.Trim(strings.ReplaceAll(v,"/","\n"))
			case 6:
				insert.PublicIp = gstr.Trim(strings.ReplaceAll(v,"/","\n"))
			}
		}
		Inserts = append(Inserts, &insert)
	}
	if len(Inserts) == 0{
		return "",errors.New("导入主机资产失败,没有有效的数据")
	}
	_,err = dao.AssetsType.Data(Inserts).Insert()
	if err != nil {
		logger.WebLog.Warningf("导入主机资产表格失败，数据库错误:%s", err.Error())
		return "",errors.New("导入主机资产失败,插入数据库错误")
	}
	if errorCount != 0{
		msg := fmt.Sprintf("成功导入%d条记录,%d条导入失败,因为不是有效数据", len(Inserts), errorCount)
		return msg,nil
	}else{
		msg := fmt.Sprintf("成功导入%d条记录", len(Inserts))
		return msg,nil
	}
}

// EmptyType 清空主机资产表
func (s *serviceAssets) EmptyType()error{
	_,err := dao.AssetsType.Where("1=",1).Delete()
	if err != nil {
		logger.WebLog.Warningf("清空主机资产表,数据库错误:%s", err.Error())
		return errors.New("清空主机资产失败,数据库错误")
	}
	dao.AssetsReports.Update(g.Map{"attribution":""},1,1)
	dao.AssetsWeb.Update(g.Map{"attribution":""},1, 1)
	return nil
}


// GroupAssetsPcDepartment 返回终端资产部门组
func (s *serviceAssets) GroupAssetsPcDepartment(page, limit int, search interface{})*model.ResponseAssetsPcDepartmentGroup{
	var result []model.ResponseAssetsPcDepartmentGroupInfo
	SearchModel := dao.AssetsComputer.Clone()
	searchStr := gconv.String(search)
	if searchStr != ""{
		SearchModel = SearchModel.Where("department like ?", "%"+searchStr+"%")
	}
	count,_ := SearchModel.Fields("department").Group("department").Count()
	if page > 0 && limit > 0 {
		err := SearchModel.Limit((page-1)*limit,limit).Fields("department").Group("department").Scan(&result)
		if err != nil {
			return &model.ResponseAssetsPcDepartmentGroup{Code:201, Msg:"查询失败,数据库错误", Count:0, Data:nil}
		}
	}else{
		return &model.ResponseAssetsPcDepartmentGroup{Code:201, Msg:"查询失败,分页参数有误", Count:0, Data:nil}
	}
	index := (page-1)*limit
	for i,_:=range result{
		index++
		result[i].ID = index
	}
	return &model.ResponseAssetsPcDepartmentGroup{Code:0, Msg:"ok", Count:int64(count), Data:result}
}

// GroupAssetsPcDepartmentSub 返回终端资产二级部门组
func (s *serviceAssets) GroupAssetsPcDepartmentSub(page, limit int, search interface{})*model.ResponseAssetsPcDepartmentSubGroup{
	var result []model.ResponseAssetsPcDepartmentSubGroupInfo
	SearchModel := dao.AssetsComputer.Clone()
	searchStr := gconv.String(search)
	if searchStr != ""{
		SearchModel = SearchModel.Where("department_sub like ?", "%"+searchStr+"%")
	}
	count,_ := SearchModel.Fields("department_sub").Group("department_sub").Count()
	if page > 0 && limit > 0 {
		err := SearchModel.Limit((page-1)*limit,limit).Fields("department_sub").Group("department_sub").Scan(&result)
		if err != nil {
			return &model.ResponseAssetsPcDepartmentSubGroup{Code:201, Msg:"查询失败,数据库错误", Count:0, Data:nil}
		}
	}else{
		return &model.ResponseAssetsPcDepartmentSubGroup{Code:201, Msg:"查询失败,分页参数有误", Count:0, Data:nil}
	}
	index := (page-1)*limit
	for i,_:=range result{
		index++
		result[i].ID = index
	}
	return &model.ResponseAssetsPcDepartmentSubGroup{Code:0, Msg:"ok", Count:int64(count), Data:result}
}

// SearchPc 终端资产所属模糊分页查询
func (s *serviceAssets) SearchPc(page, limit int, search interface{})*model.ResponseAssetsPc{
	var result []*model.ResponseAssetsPcInfo
	SearchModel := dao.AssetsComputer.Clone()
	searchStr := gconv.String(search)
	if search != ""{
		j := gjson.New(searchStr)
		if gconv.String(j.Get("Department")) != ""{
			SearchModel = SearchModel.Where("department like ?", "%"+gconv.String(j.Get("Department"))+"%")
		}
		if gconv.String(j.Get("DepartmentSub")) != ""{
			SearchModel = SearchModel.Where("department_sub like ?", "%"+gconv.String(j.Get("DepartmentSub"))+"%")
		}
		if gconv.String(j.Get("WorkNumber")) != ""{
			SearchModel = SearchModel.Where("work_number like ?", "%"+gconv.String(j.Get("WorkNumber"))+"%")
		}
		if gconv.String(j.Get("PersonName")) != ""{
			SearchModel = SearchModel.Where("person_name like ?", "%"+gconv.String(j.Get("PersonName"))+"%")
		}
		if gconv.String(j.Get("Address")) != ""{
			SearchModel = SearchModel.Where("address like ?", "%"+gconv.String(j.Get("Address"))+"%")
		}
		if gconv.String(j.Get("InternetFlag")) != ""{
			SearchModel = SearchModel.Where("internet_flag like ?", "%"+gconv.String(j.Get("InternetFlag"))+"%")
		}
		if gconv.String(j.Get("VpnFlag")) != ""{
			SearchModel = SearchModel.Where("vpn_flag like ?", "%"+gconv.String(j.Get("VpnFlag"))+"%")
		}
	}
	count,_ := SearchModel.Count()
	if page > 0 && limit > 0 {
		err := SearchModel.Order("id desc").Limit((page-1)*limit,limit).Scan(&result)
		if err != nil {
			logger.WebLog.Warningf("终端资产分页查询 数据库错误:%s", err.Error())
			return &model.ResponseAssetsPc{Code:201, Msg:"查询失败,数据库错误", Count:0, Data:nil}
		}
	}else{
		return &model.ResponseAssetsPc{Code:201, Msg:"查询失败,分页参数有误", Count:0, Data:nil}
	}
	index := (page-1)*limit
	for i,_:=range result{
		index++
		result[i].ID1 = result[i].Id
		result[i].Id = uint(index)
	}
	return &model.ResponseAssetsPc{Code:0, Msg:"ok", Count:int64(count), Data:result}
}

// PcAdd 添加修改终端资产
func (s *serviceAssets) PcAdd(r *model.RequestAssetsPcAdd)error{
	if len(r.FlagId) != 0 {
		if _,err := dao.AssetsComputer.Update(r,"id",r.FlagId); err != nil{
			logger.WebLog.Warningf("资产管理-修改终端资产失败:%s", err.Error())
			return errors.New("修改终端资产失败,数据库错误")
		}
	}else{
		if _,err := dao.AssetsComputer.Insert(r); err != nil{
			logger.WebLog.Warningf("资产管理-添加终端资产失败:%s", err.Error())
			return errors.New("添加终端资产失败,数据库错误")
		}
	}
	return nil
}

// PcShow 查看指定ID的终端资产信息
func (s *serviceAssets) PcShow(r *model.RequestAssetsTypeDelete)(*model.AssetsComputer,error){
	result,err := dao.AssetsComputer.Where("id=?",r.ID).FindOne()
	if err != nil{
		logger.WebLog.Warningf("资产管理-查看终端资产失败:%s", err.Error())
		return nil,errors.New("查看终端资产失败,数据库错误")
	}
	return result,nil
}

// PcDelete 删除终端资产
func (s *serviceAssets) PcDelete(r *model.RequestAssetsTypeDelete)error{
	count,err := dao.AssetsComputer.Where("id=?",r.ID).Count()
	if err != nil{
		logger.WebLog.Warningf("资产管理-删除终端资产失败:%s", err.Error())
		return errors.New("删除终端资产失败,数据库错误")
	}
	if count == 0{
		return errors.New("删除终端资产失败,该资产不存在")
	}
	if _,err = dao.AssetsComputer.Where("id=?",r.ID).Delete(); err != nil{
		return errors.New(fmt.Sprintf("删除终端资产失败:%s", err.Error()))
	}
	return nil
}

// PcEmpty 清空终端资产
func (s *serviceAssets) PcEmpty()error{
	_,err := dao.AssetsComputer.Where("1=",1).Delete()
	if err != nil {
		logger.WebLog.Warningf("清空终端资产表,数据库错误:%s", err.Error())
		return errors.New("清空终端资产失败,数据库错误")
	}
	return nil
}

// DownloadPcTemplate 下载终端资产模板
func (s *serviceAssets) DownloadPcTemplate()(*bytes.Buffer, error){
	xlsx := excelize.NewFile()
	index := xlsx.NewSheet("终端资产表")
	xlsx.DeleteSheet("Sheet1")
	table_name := "终端资产表"
	xlsx.SetCellValue(table_name,"A1","部门")
	xlsx.SetCellValue(table_name,"B1","二级部门")
	xlsx.SetCellValue(table_name,"C1","员工姓名")
	xlsx.SetCellValue(table_name,"D1","工号")
	xlsx.SetCellValue(table_name,"E1","计算机类型")
	xlsx.SetCellValue(table_name,"F1","计算机名")
	xlsx.SetCellValue(table_name,"G1","涉密等级")
	xlsx.SetCellValue(table_name,"H1","IP地址")
	xlsx.SetCellValue(table_name,"I1","外网权限")
	xlsx.SetCellValue(table_name,"J1","拷贝权限")
	xlsx.SetCellValue(table_name,"K1","邮箱权限")
	xlsx.SetCellValue(table_name,"L1","VPN权限")
	xlsx.SetCellValue(table_name,"M1","PDM权限")
	xlsx.SetCellValue(table_name,"N1","备注")
	xlsx.SetCellValue(table_name, "A" + strconv.Itoa(2), "xx部门")
	xlsx.SetCellValue(table_name, "B" + strconv.Itoa(2), "xx处")
	xlsx.SetCellValue(table_name, "C" + strconv.Itoa(2), "张三")
	xlsx.SetCellValue(table_name, "D" + strconv.Itoa(2), "0012345")
	xlsx.SetCellValue(table_name, "E" + strconv.Itoa(2), "台式机")
	xlsx.SetCellValue(table_name, "F" + strconv.Itoa(2), "cq-zhangsan")
	xlsx.SetCellValue(table_name, "G" + strconv.Itoa(2), "普通商密")
	xlsx.SetCellValue(table_name, "H" + strconv.Itoa(2), "192.168.1.1")
	xlsx.SetCellValue(table_name, "I" + strconv.Itoa(2), "是")
	xlsx.SetCellValue(table_name, "J" + strconv.Itoa(2), "是")
	xlsx.SetCellValue(table_name, "K" + strconv.Itoa(2), "")
	xlsx.SetCellValue(table_name, "L" + strconv.Itoa(2), "")
	xlsx.SetCellValue(table_name, "M" + strconv.Itoa(2), "是")
	xlsx.SetCellValue(table_name, "N" + strconv.Itoa(2), "xx楼层xx工位")
	xlsx.SetColWidth(table_name, "A", "A", 12)
	xlsx.SetColWidth(table_name, "B", "B", 12)
	xlsx.SetColWidth(table_name, "C", "C", 15)
	xlsx.SetColWidth(table_name, "D", "D", 14)
	xlsx.SetColWidth(table_name, "E", "E", 13)
	xlsx.SetColWidth(table_name, "F", "F", 15)
	xlsx.SetColWidth(table_name, "G", "G", 15)
	xlsx.SetColWidth(table_name, "H", "H", 15)
	xlsx.SetActiveSheet(index)
	var buf bytes.Buffer
	err := xlsx.Write(&buf)
	if err != nil {
		return &bytes.Buffer{},errors.New("下载终端资产表模板错误,xlsx流写入失败")
	}
	return &buf,nil
}

// ExportPc 导出终端资产
func (s *serviceAssets) ExportPc()(*bytes.Buffer, error){
	var result []model.AssetsComputer
	err := dao.AssetsComputer.Scan(&result)
	if err != nil {
		return &bytes.Buffer{},errors.New("导出终端资产数据失败,数据库错误")
	}
	if len(result) == 0{
		return &bytes.Buffer{},errors.New("导出终端资产数据失败,无资产数据")
	}
	xlsx := excelize.NewFile()
	index := xlsx.NewSheet("终端资产表")
	xlsx.DeleteSheet("Sheet1")
	table_name := "终端资产表"
	xlsx.SetCellValue(table_name,"A1","部门")
	xlsx.SetCellValue(table_name,"B1","二级部门")
	xlsx.SetCellValue(table_name,"C1","员工姓名")
	xlsx.SetCellValue(table_name,"D1","员工编号")
	xlsx.SetCellValue(table_name,"E1","计算机类型")
	xlsx.SetCellValue(table_name,"F1","计算机名")
	xlsx.SetCellValue(table_name,"G1","涉密界别")
	xlsx.SetCellValue(table_name,"H1","IP地址")
	xlsx.SetCellValue(table_name,"I1","外网权限")
	xlsx.SetCellValue(table_name,"J1","拷贝权限")
	xlsx.SetCellValue(table_name,"K1","邮箱权限")
	xlsx.SetCellValue(table_name,"L1","VPN权限")
	xlsx.SetCellValue(table_name,"M1","PDM权限")
	xlsx.SetCellValue(table_name,"N1","备注")
	for i, info := range result {
		xlsx.SetCellValue(table_name, "A" + strconv.Itoa(i+2), info.Department)
		xlsx.SetCellValue(table_name, "B" + strconv.Itoa(i+2), info.DepartmentSub)
		xlsx.SetCellValue(table_name, "C" + strconv.Itoa(i+2), info.PersonName)
		xlsx.SetCellValue(table_name, "D" + strconv.Itoa(i+2), info.WorkNumber)
		xlsx.SetCellValue(table_name, "E" + strconv.Itoa(i+2), info.ComputerType)
		xlsx.SetCellValue(table_name, "F" + strconv.Itoa(i+2), info.ComputerName)
		xlsx.SetCellValue(table_name, "H" + strconv.Itoa(i+2), info.Address)
		xlsx.SetCellValue(table_name, "G" + strconv.Itoa(i+2), info.SecretLevel)
		xlsx.SetCellValue(table_name, "I" + strconv.Itoa(i+2), info.InternetFlag)
		xlsx.SetCellValue(table_name, "J" + strconv.Itoa(i+2), info.FileCopyFlag)
		xlsx.SetCellValue(table_name, "K" + strconv.Itoa(i+2), info.EmailFlag)
		xlsx.SetCellValue(table_name, "L" + strconv.Itoa(i+2), info.VpnFlag)
		xlsx.SetCellValue(table_name, "M" + strconv.Itoa(i+2), info.PdmFlag)
		xlsx.SetCellValue(table_name, "N" + strconv.Itoa(i+2), info.Remarks)
	}
	xlsx.SetColWidth(table_name, "A", "A", 18)
	xlsx.SetColWidth(table_name, "B", "B", 18)
	xlsx.SetColWidth(table_name, "C", "C", 15)
	xlsx.SetColWidth(table_name, "D", "D", 18)
	xlsx.SetColWidth(table_name, "E", "E", 18)
	xlsx.SetColWidth(table_name, "F", "F", 20)
	xlsx.SetColWidth(table_name, "G", "G", 15)
	xlsx.SetColWidth(table_name, "H", "H", 15)
	xlsx.SetColWidth(table_name, "J", "J", 15)
	xlsx.SetColWidth(table_name, "K", "K", 15)
	xlsx.SetColWidth(table_name, "L", "L", 15)
	xlsx.SetColWidth(table_name, "M", "M", 15)
	xlsx.SetColWidth(table_name, "N", "M", 20)
	xlsx.SetActiveSheet(index) // 设置工作簿的默认工作表
	var buf bytes.Buffer
	err = xlsx.Write(&buf)
	if err != nil {
		return &bytes.Buffer{},errors.New("导出终端资产数据失败,xlsx流写入失败")
	}
	return &buf,nil
}

// ImportPc 导入终端资产表
func (s *serviceAssets) ImportPc(filename string)(string,error){
	xlsx, err := excelize.OpenFile(filename)
	if err != nil {
		logger.WebLog.Warningf("导入终端资产失败:%s", err.Error())
		return "",errors.New("导入终端资产失败,不是有效的xlsx文档")
	}
	rows,err := xlsx.GetRows("终端资产表")
	if err != nil {
		logger.WebLog.Warningf("导入终端资产表格失败:%s", err.Error())
		return "",errors.New("导入终端资产失败,不是有效的资产模板")
	}
	var Inserts []*model.AssetsComputer
	errorCount := 0 // 记录不合格的行数
	for i,row := range rows {
		if i == 0{ // 过滤首行
			continue
		}
		if len(row) < 12{
			errorCount++
			continue
		}
		if len(row[2]) == 0 || len(row[0]) == 0{ //姓名和部门不能为空
			errorCount++
			continue
		}
		insert := model.AssetsComputer{}
		for index,v := range row{
			switch index{
			case 0:
				insert.Department = gstr.Trim(v)
			case 1:
				insert.DepartmentSub = gstr.Trim(v)
			case 2:
				insert.PersonName = gstr.Trim(v)
			case 3:
				insert.WorkNumber = gstr.Trim(v)
			case 4:
				insert.ComputerType = gstr.Trim(v)
			case 5:
				insert.ComputerName = gstr.Trim(v)
			case 7:
				insert.Address = gstr.Trim(v)
			case 6:
				insert.SecretLevel = gstr.Trim(v)
			case 8:
				insert.InternetFlag = gstr.Trim(v)
			case 9:
				insert.FileCopyFlag = gstr.Trim(v)
			case 10:
				insert.EmailFlag = gstr.Trim(v)
			case 11:
				insert.VpnFlag = gstr.Trim(v)
			case 12:
				insert.PdmFlag = gstr.Trim(v)
			case 13:
				insert.Remarks = gstr.Trim(v)
			}
		}
		Inserts = append(Inserts, &insert)
	}
	if len(Inserts) == 0{
		return "",errors.New("导入终端资产失败,没有有效的数据")
	}
	_,err = dao.AssetsComputer.Data(Inserts).Insert()
	if err != nil {
		logger.WebLog.Warningf("导入终端资产失败，数据库错误:%s", err.Error())
		return "",errors.New("导入终端资产,插入数据库错误")
	}
	if errorCount != 0{
		msg := fmt.Sprintf("成功导入%d条记录,%d条导入失败,因为不是有效数据", len(Inserts), errorCount)
		return msg,nil
	}else{
		msg := fmt.Sprintf("成功导入%d条记录", len(Inserts))
		return msg,nil
	}
}

// IchanganUsers 加载员工数据
func (s *serviceAssets) IchanganUsers()error{
	count,err := dao.AssetsUsers.Count()
	if err != nil {
		return errors.New("查询员工数量失败，加载失败")
	}
	if count > 1{
		return errors.New("已加载员工信息，不能重复加载")
	}
	dbFile := "./public/ip/user.db"
	if !gfile.Exists(dbFile){
		return errors.New("员工数据文件不存在,加载失败")
	}
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil{
		return errors.New("员工数据打开失败,加载失败")
	}
	defer func() {
		if db != nil{
			db.Close()
		}
	}()
	var sql = "select * from userinfo"
	stmt, err := db.Prepare(sql)
	if err != nil{
		return errors.New("sql语句错误,加载失败")
	}
	rows, err := stmt.Query()
	if err != nil {
		return errors.New("sql查询失败,加载失败")
	}
	var results = make([]*model.AssetsUsers, 0)
	for rows.Next(){
		var (
			uid int
			UserLoginID string
			IDCard string
			ResMobile string
			DepartmentName string
			WorkerType string
			Sex string
			EntryReason string
			EntryDate string
			PostName string
			UserID string
		)
		err := rows.Scan(&uid, &UserLoginID, &IDCard, &ResMobile, &DepartmentName, &WorkerType, &Sex, &EntryReason, &EntryDate, &PostName, &UserID)
		if err != nil{
			continue
		}
		tmpUser := model.AssetsUsers{
			WorkNumber: UserLoginID,
			IdCard: IDCard,
			Phone: ResMobile,
			Department: DepartmentName,
			WorkerType: WorkerType,
			Sex: Sex,
			EntryReason: EntryReason,
			EntryDate: EntryDate,
			PortName: PostName,
			UserId: UserID,
		}
		results = append(results, &tmpUser)
	}
	_, err = dao.AssetsUsers.Insert(results)
	if err != nil {
		logger.WebLog.Warningf("加载员工数据错误:%s", err.Error())
	}
	return nil
}

// PcShowUser 查看指定ID的终端和员工信息
func (s *serviceAssets) PcShowUser(r *model.RequestAssetsTypeDelete)(*model.ResponseAssetsPcUserInfo,error){
	result,err := dao.AssetsComputer.Where("id=?",r.ID).FindOne()
	if err != nil{
		logger.WebLog.Warningf("资产管理-查看终端资产失败:%s", err.Error())
		return nil,errors.New("查看终端资产失败,数据库错误")
	}
	var results model.ResponseAssetsPcUserInfo
	results.Computer = result
	user,err := dao.AssetsUsers.Where("work_number=?", result.WorkNumber).FindOne()
	if err == nil{
		results.Users = user
	}
	return &results,nil
}


// GroupAssetsManagerName 返回安全管理员Group分组
func (s *serviceAssets) GroupAssetsManagerName(page, limit int, search interface{})*model.ResponseAssetsManagerNameGroup{
	var result []model.ResponseAssetsManagerNameGroupInfo
	SearchModel := dao.AssetsManager.Clone()
	searchStr := gconv.String(search)
	if searchStr != ""{
		SearchModel = SearchModel.Where("manager_name like ?", "%"+searchStr+"%")
	}
	count,_ := SearchModel.Fields("manager_name").Group("manager_name").Count()
	if page > 0 && limit > 0 {
		err := SearchModel.Limit((page-1)*limit,limit).Fields("manager_name").Group("manager_name").Scan(&result)
		if err != nil {
			return &model.ResponseAssetsManagerNameGroup{Code:201, Msg:"查询失败,数据库错误", Count:0, Data:nil}
		}
	}else{
		return &model.ResponseAssetsManagerNameGroup{Code:201, Msg:"查询失败,分页参数有误", Count:0, Data:nil}
	}
	index := (page-1)*limit
	for i,_:=range result{
		index++
		result[i].ID = index
	}
	return &model.ResponseAssetsManagerNameGroup{Code:0, Msg:"ok", Count:int64(count), Data:result}
}

// GroupAssetsWebManagerName 返回Web资产安全管理员Group分组
func (s *serviceAssets) GroupAssetsWebManagerName(page, limit int, search interface{})*model.ResponseAssetsManagerNameGroup{
	var result []model.ResponseAssetsManagerNameGroupInfo
	SearchModel := dao.AssetsWeb.Clone()
	searchStr := gconv.String(search)
	if searchStr != ""{
		SearchModel = SearchModel.Where("manager_name like ?", "%"+searchStr+"%")
	}
	count,_ := SearchModel.Fields("manager_name").Group("manager_name").Count()
	if page > 0 && limit > 0 {
		err := SearchModel.Limit((page-1)*limit,limit).Fields("manager_name").Group("manager_name").Scan(&result)
		if err != nil {
			return &model.ResponseAssetsManagerNameGroup{Code:201, Msg:"查询失败,数据库错误", Count:0, Data:nil}
		}
	}else{
		return &model.ResponseAssetsManagerNameGroup{Code:201, Msg:"查询失败,分页参数有误", Count:0, Data:nil}
	}
	index := (page-1)*limit
	for i,_:=range result{
		index++
		result[i].ID = index
	}
	return &model.ResponseAssetsManagerNameGroup{Code:0, Msg:"ok", Count:int64(count), Data:result}
}

// GroupAssetsWebAssetsName 返回web应用业务系统名Group分组
func (s *serviceAssets) GroupAssetsWebAssetsName(page, limit int, search interface{})*model.ResponseAssetsWebAssetsNameGroup{
	var result []model.ResponseAssetsWebAssetsNameGroupInfo
	SearchModel := dao.AssetsWeb.Clone()
	searchStr := gconv.String(search)
	if searchStr != ""{
		SearchModel = SearchModel.Where("assets_name like ?", "%"+searchStr+"%")
	}
	count,_ := SearchModel.Fields("assets_name").Group("assets_name").Count()
	if page > 0 && limit > 0 {
		err := SearchModel.Limit((page-1)*limit,limit).Fields("assets_name").Group("assets_name").Scan(&result)
		if err != nil {
			return &model.ResponseAssetsWebAssetsNameGroup{Code:201, Msg:"查询失败,数据库错误", Count:0, Data:nil}
		}
	}else{
		return &model.ResponseAssetsWebAssetsNameGroup{Code:201, Msg:"查询失败,分页参数有误", Count:0, Data:nil}
	}
	index := (page-1)*limit
	for i,_:=range result{
		index++
		result[i].ID = index
	}
	return &model.ResponseAssetsWebAssetsNameGroup{Code:0, Msg:"ok", Count:int64(count), Data:result}
}

// GroupAssetsWebAttribution 返回Web资产应用系统组
func (s *serviceAssets) GroupAssetsWebAttribution(page, limit int, search interface{})*model.ResponseAssetsTypeAttributionGroup{
	var result []model.ResponseAssetsTypeAttributionGroupInfo
	SearchModel := dao.AssetsWeb.Clone()
	searchStr := gconv.String(search)
	if searchStr != ""{
		SearchModel = SearchModel.Where("attribution like ?", "%"+searchStr+"%")
	}
	count,_ := SearchModel.Fields("attribution").Group("attribution").Count()
	if page > 0 && limit > 0 {
		err := SearchModel.Limit((page-1)*limit,limit).Fields("attribution").Group("attribution").Scan(&result)
		if err != nil {
			return &model.ResponseAssetsTypeAttributionGroup{Code:201, Msg:"查询失败,数据库错误", Count:0, Data:nil}
		}
	}else{
		return &model.ResponseAssetsTypeAttributionGroup{Code:201, Msg:"查询失败,分页参数有误", Count:0, Data:nil}
	}
	index := (page-1)*limit
	for i,_:=range result{
		index++
		result[i].ID = index
	}
	return &model.ResponseAssetsTypeAttributionGroup{Code:0, Msg:"ok", Count:int64(count), Data:result}
}

// SearchWeb 应用系统模糊分页查询
func (s *serviceAssets) SearchWeb(page, limit int, search interface{})*model.ResponseAssetsWeb{
	var result []*model.ResponseAssetsWebInfo
	SearchModel := dao.AssetsWeb.Clone()
	searchStr := gconv.String(search)
	if search != ""{
		j := gjson.New(searchStr)
		if gconv.String(j.Get("AttriBution")) != ""{
			SearchModel = SearchModel.Where("attribution like ?", "%"+gconv.String(j.Get("AttriBution"))+"%")
		}
		if gconv.String(j.Get("ManagerName")) != ""{
			SearchModel = SearchModel.Where("manager_name like ?", "%"+gconv.String(j.Get("ManagerName"))+"%")
		}
		if gconv.String(j.Get("AssetsName")) != ""{
			SearchModel = SearchModel.Where("assets_name like ?", "%"+gconv.String(j.Get("AssetsName"))+"%")
		}
		if gconv.String(j.Get("Urls")) != ""{
			SearchModel = SearchModel.Where("urls like ?", "%"+gconv.String(j.Get("Urls"))+"%")
		}
		if gconv.String(j.Get("FingerPrint")) != ""{
			SearchModel = SearchModel.Where("fingerprint like ?", "%"+gconv.String(j.Get("FingerPrint"))+"%")
		}
	}
	count,_ := SearchModel.Count()
	if page > 0 && limit > 0 {
		err := SearchModel.Order("id desc").Limit((page-1)*limit,limit).Scan(&result)
		if err != nil {
			logger.WebLog.Warningf("应用系统分页查询 数据库错误:%s", err.Error())
			return &model.ResponseAssetsWeb{Code:201, Msg:"查询失败,数据库错误", Count:0, Data:nil}
		}
	}else{
		return &model.ResponseAssetsWeb{Code:201, Msg:"查询失败,分页参数有误", Count:0, Data:nil}
	}
	index := (page-1)*limit
	for i,_:=range result{
		index++
		result[i].ID1 = result[i].Id
		result[i].Id = uint(index)
		levelCount,err := dao.AssetsReports.Where("attribution=?", result[i].Attribution).Count()
		if err == nil{
			result[i].LevelCount = levelCount
		}
		levelNoCount,err := dao.AssetsReports.Where("attribution=?", result[i].Attribution).Where("level_status=?", 2).Count()
		if err == nil{
			result[i].LevelNoCount = levelNoCount
		}
		result[i].Urls = strings.ReplaceAll(result[i].Urls,"\n","<br/>")
		if len(result[i].ScreenshotsPath) != 0{
			var imgsrc []model.ResponseAssetsWebInfoImgsrcInfo
			imgsrc = append(imgsrc, model.ResponseAssetsWebInfoImgsrcInfo{Src: result[i].ScreenshotsPath})
			result[i].Data = imgsrc
		}
	}
	return &model.ResponseAssetsWeb{Code:0, Msg:"ok", Count:int64(count), Data:result}
}

// WebAdd 添加修改业务系统资产
func (s *serviceAssets) WebAdd(r *model.RequestAssetsWebAdd)error{
	count,err := dao.AssetsType.Where("attribution=?", r.Attribution).Count()
	if err != nil {
		return errors.New("添加失败,数据库错误,查询应用系统失败")
	}
	if count == 0{
		return errors.New("添加失败,没有此应用系统名,请在主机资产里添加")
	}
	count,err = dao.AssetsManager.Where("manager_name=?",r.ManagerName).Count()
	if err != nil {
		return errors.New("添加失败,数据库错误,查询安全管理员失败")
	}
	if count == 0{
		return errors.New("添加失败,没有此安全管理员,请在安全管理员资产里添加")
	}
	if len(r.Fid) != 0 { // 修改业务系统
		res1,err := dao.AssetsWeb.Where("id=?",r.Fid).FindOne()
		if err != nil {
			logger.WebLog.Warningf("业务系统管理-修改查询资产失败:%s", err.Error())
		}
		if _,err = dao.AssetsWeb.Update(r, "id", r.Fid);err != nil{
			logger.WebLog.Warningf("业务系统管理-修改资产失败:%s", err.Error())
			return errors.New("修改失败,数据库错误")
		}
		if res1.Attribution != r.Attribution || res1.ManagerName != r.ManagerName || res1.AssetsName != r.AssetsName{ // 业务系统资产变化了则修改渗透测试报告中的字段
			_,err = dao.AssetsReports.Data(g.Map{"attribution": r.Attribution,"manager_name":r.ManagerName,"assets_name":r.AssetsName}).
				Where("attribution", res1.Attribution).Where("manager_name", res1.ManagerName).
				Where("assets_name",res1.AssetsName).Update()
			if err != nil {
				logger.WebLog.Warningf("更改渗透测试报告信息错误:%s", err.Error())
			}
		}

	}else{
		if _,err = dao.AssetsWeb.Insert(r); err != nil{
			logger.WebLog.Warningf("业务系统管理-添加资产失败:%s", err.Error())
			return errors.New("添加失败,数据库错误")
		}
	}
	return nil
}

// WebDelete 删除业务系统
func (s *serviceAssets) WebDelete(r *model.RequestAssetsTypeDelete)error{
	count,err := dao.AssetsWeb.Where("id=?",r.ID).Count()
	if err != nil{
		logger.WebLog.Warningf("业务系统-删除失败:%s", err.Error())
		return errors.New("删除业务系统失败,数据库错误")
	}
	if count == 0{
		return errors.New("删除业务系统失败,该资产不存在")
	}
	if _,err = dao.AssetsWeb.Where("id=?",r.ID).Delete(); err != nil{
		return errors.New(fmt.Sprintf("删除业务系统失败:%s", err.Error()))
	}
	return nil
}

// ExportWeb 导出业务系统
func (s *serviceAssets) ExportWeb()(*bytes.Buffer, error){
	var result []model.AssetsWeb
	err := dao.AssetsWeb.Scan(&result)
	if err != nil {
		return &bytes.Buffer{},errors.New("导出业务系统数据失败,数据库错误")
	}
	if len(result) == 0{
		return &bytes.Buffer{},errors.New("导出业务系统数据失败,无资产数据")
	}
	xlsx := excelize.NewFile()
	index := xlsx.NewSheet("业务系统资产表")
	xlsx.DeleteSheet("Sheet1")
	table_name := "业务系统资产表"
	xlsx.SetCellValue(table_name,"A1","应用系统")
	xlsx.SetCellValue(table_name,"B1","安全管理员")
	xlsx.SetCellValue(table_name,"C1","业务系统")
	xlsx.SetCellValue(table_name,"D1","Url")
	xlsx.SetCellValue(table_name,"E1","业务系统指纹")
	xlsx.SetCellValue(table_name,"F1","Web服务类型")
	xlsx.SetCellValue(table_name,"G1","备注")
	for i, info := range result {
		xlsx.SetCellValue(table_name, "A" + strconv.Itoa(i+2), info.Attribution)
		xlsx.SetCellValue(table_name, "B" + strconv.Itoa(i+2), info.ManagerName)
		xlsx.SetCellValue(table_name, "C" + strconv.Itoa(i+2), info.AssetsName)
		xlsx.SetCellValue(table_name, "D" + strconv.Itoa(i+2), info.Urls)
		xlsx.SetCellValue(table_name, "E" + strconv.Itoa(i+2), info.Fingerprint)
		xlsx.SetCellValue(table_name, "F" + strconv.Itoa(i+2), info.Webserver)
		xlsx.SetCellValue(table_name, "G" + strconv.Itoa(i+2), info.Remarks)
	}
	xlsx.SetColWidth(table_name, "A", "A", 20)
	xlsx.SetColWidth(table_name, "B", "B", 16)
	xlsx.SetColWidth(table_name, "C", "C", 20)
	xlsx.SetColWidth(table_name, "D", "D", 25)
	xlsx.SetColWidth(table_name, "E", "E", 20)
	xlsx.SetColWidth(table_name, "F", "F", 14)
	xlsx.SetColWidth(table_name, "G", "G", 20)
	xlsx.SetActiveSheet(index) // 设置工作簿的默认工作表
	var buf bytes.Buffer
	err = xlsx.Write(&buf)
	if err != nil {
		return &bytes.Buffer{},errors.New("导出业务系统数据失败,xlsx流写入失败")
	}
	return &buf,nil
}

// WebExportAdd 增加渗透测试漏洞
func (s *serviceAssets) WebExportAdd(r *gjson.Json)error{
	if len(r.GetString("ReportDataTime")) == 0{
		return errors.New("渗透测试报告时间不能为空")
	}
	if len(r.GetString("report_file_name")) == 0{
		return errors.New("请上传渗透测试报告")
	}
	if len(r.GetString("ReportId")) == 0{
		return errors.New("业务系统ID不能空")
	}
	result,err := dao.AssetsWeb.Where("id=?", r.GetString("ReportId")).FindOne()
	if err != nil{
		return errors.New("业务系统ID错误")
	}
	if result == nil{
		return errors.New("该业务系统ID不存在")
	}
	reportTime := gtime.New(r.GetString("ReportDataTime"))
	reportPath := r.GetString("report_file_name")
	index := 0
	var inserts []*model.AssetsReports
	for i := 0;i<1000;i++{ // 单次提交最多不超过1000个漏洞
		index ++
		tmp := "level_name" +  strconv.Itoa(index)
		if len(r.GetString(tmp)) == 0{
			break
		}
		tmp1 := model.AssetsReports{}
		tmp1.Level = r.GetInt("level" +  strconv.Itoa(index))
		tmp1.LevelName = r.GetString(tmp)
		tmp1.LevelStatus = r.GetInt("level_status" +  strconv.Itoa(index))
		tmp1.FileDate = reportTime
		tmp1.FilePath = reportPath
		tmp1.AssetsName = result.AssetsName
		tmp1.ManagerName = result.ManagerName
		tmp1.Attribution = result.Attribution
		inserts = append(inserts, &tmp1)
	}
	if len(inserts) == 0{
		return errors.New("未发现漏洞,请添加")
	}
	if _,err = dao.AssetsReports.Insert(inserts); err != nil{
		logger.WebLog.Warningf("添加漏洞数据库错误:%s", err.Error())
		return errors.New("添加漏洞失败,插入数据库错误")
	}
	return nil
}

// SearchWebReport 渗透测试报告模糊分页查询
func (s *serviceAssets) SearchWebReport(page, limit int, assetsWebId int)*model.ResponseAssetsWebReport{
	res,err := dao.AssetsWeb.Where("id=?", assetsWebId).FindOne()
	if err != nil {
		return &model.ResponseAssetsWebReport{Code:201, Msg:"业务系统ID错误", Count:0, Data:nil}
	}
	var result []*model.AssetsReports
	SearchModel := dao.AssetsReports.Clone()
	SearchModel = SearchModel.Where("assets_name=?", res.AssetsName)
	count,_ := SearchModel.Count()
	if count == 0{
		return &model.ResponseAssetsWebReport{Code:201, Msg:"该业务系统无漏洞,请添加", Count:0, Data:nil}
	}
	if page > 0 && limit > 0 {
		err := SearchModel.Order("id desc").Limit((page-1)*limit,limit).Scan(&result)
		if err != nil {
			logger.WebLog.Warningf("渗透测试报告分页查询 数据库错误:%s", err.Error())
			return &model.ResponseAssetsWebReport{Code:201, Msg:"查询失败,数据库错误", Count:0, Data:nil}
		}
	}else{
		return &model.ResponseAssetsWebReport{Code:201, Msg:"查询失败,分页参数有误", Count:0, Data:nil}
	}
	return &model.ResponseAssetsWebReport{Code:0, Msg:"ok", Count:int64(count), Data:result}
}

// WebShow 查看指定ID的业务系统信息
func (s *serviceAssets) WebShow(r *model.RequestAssetsTypeDelete)(*model.AssetsWeb,error){
	result,err := dao.AssetsWeb.Where("id=?",r.ID).FindOne()
	if err != nil{
		logger.WebLog.Warningf("业务系统-查看信息失败:%s", err.Error())
		return nil,errors.New("查看业务系统失败,数据库错误")
	}
	return result,nil
}



// GroupAssetsReportLevelName 返回渗透测试报告漏洞名Group分组
func (s *serviceAssets) GroupAssetsReportLevelName(page, limit int, search interface{})*model.ResponseAssetsReportLevelNameGroup{
	var result []model.ResponseAssetsReportLevelNameGroupInfo
	SearchModel := dao.AssetsReports.Clone()
	searchStr := gconv.String(search)
	if searchStr != ""{
		SearchModel = SearchModel.Where("level_name like ?", "%"+searchStr+"%")
	}
	count,_ := SearchModel.Fields("level_name").Group("level_name").Count()
	if page > 0 && limit > 0 {
		err := SearchModel.Limit((page-1)*limit,limit).Fields("level_name").Group("level_name").Scan(&result)
		if err != nil {
			return &model.ResponseAssetsReportLevelNameGroup{Code:201, Msg:"查询失败,数据库错误", Count:0, Data:nil}
		}
	}else{
		return &model.ResponseAssetsReportLevelNameGroup{Code:201, Msg:"查询失败,分页参数有误", Count:0, Data:nil}
	}
	index := (page-1)*limit
	for i,_:=range result{
		index++
		result[i].ID = index
	}
	return &model.ResponseAssetsReportLevelNameGroup{Code:0, Msg:"ok", Count:int64(count), Data:result}
}

// GroupAssetsReportAttribution 返回渗透测试报告应用系统组
func (s *serviceAssets) GroupAssetsReportAttribution(page, limit int, search interface{})*model.ResponseAssetsTypeAttributionGroup{
	var result []model.ResponseAssetsTypeAttributionGroupInfo
	SearchModel := dao.AssetsReports.Clone()
	searchStr := gconv.String(search)
	if searchStr != ""{
		SearchModel = SearchModel.Where("attribution like ?", "%"+searchStr+"%")
	}
	count,_ := SearchModel.Fields("attribution").Group("attribution").Count()
	if page > 0 && limit > 0 {
		err := SearchModel.Limit((page-1)*limit,limit).Fields("attribution").Group("attribution").Scan(&result)
		if err != nil {
			return &model.ResponseAssetsTypeAttributionGroup{Code:201, Msg:"查询失败,数据库错误", Count:0, Data:nil}
		}
	}else{
		return &model.ResponseAssetsTypeAttributionGroup{Code:201, Msg:"查询失败,分页参数有误", Count:0, Data:nil}
	}
	index := (page-1)*limit
	for i,_:=range result{
		index++
		result[i].ID = index
	}
	return &model.ResponseAssetsTypeAttributionGroup{Code:0, Msg:"ok", Count:int64(count), Data:result}
}

// GroupAssetsReportManagerName 返回安全管理员Group分组
func (s *serviceAssets) GroupAssetsReportManagerName(page, limit int, search interface{})*model.ResponseAssetsManagerNameGroup{
	var result []model.ResponseAssetsManagerNameGroupInfo
	SearchModel := dao.AssetsReports.Clone()
	searchStr := gconv.String(search)
	if searchStr != ""{
		SearchModel = SearchModel.Where("manager_name like ?", "%"+searchStr+"%")
	}
	count,_ := SearchModel.Fields("manager_name").Group("manager_name").Count()
	if page > 0 && limit > 0 {
		err := SearchModel.Limit((page-1)*limit,limit).Fields("manager_name").Group("manager_name").Scan(&result)
		if err != nil {
			return &model.ResponseAssetsManagerNameGroup{Code:201, Msg:"查询失败,数据库错误", Count:0, Data:nil}
		}
	}else{
		return &model.ResponseAssetsManagerNameGroup{Code:201, Msg:"查询失败,分页参数有误", Count:0, Data:nil}
	}
	index := (page-1)*limit
	for i,_:=range result{
		index++
		result[i].ID = index
	}
	return &model.ResponseAssetsManagerNameGroup{Code:0, Msg:"ok", Count:int64(count), Data:result}
}

// GroupAssetsReportAssetsName 返回web应用业务系统名Group分组
func (s *serviceAssets) GroupAssetsReportAssetsName(page, limit int, search interface{})*model.ResponseAssetsWebAssetsNameGroup{
	var result []model.ResponseAssetsWebAssetsNameGroupInfo
	SearchModel := dao.AssetsReports.Clone()
	searchStr := gconv.String(search)
	if searchStr != ""{
		SearchModel = SearchModel.Where("assets_name like ?", "%"+searchStr+"%")
	}
	count,_ := SearchModel.Fields("assets_name").Group("assets_name").Count()
	if page > 0 && limit > 0 {
		err := SearchModel.Limit((page-1)*limit,limit).Fields("assets_name").Group("assets_name").Scan(&result)
		if err != nil {
			return &model.ResponseAssetsWebAssetsNameGroup{Code:201, Msg:"查询失败,数据库错误", Count:0, Data:nil}
		}
	}else{
		return &model.ResponseAssetsWebAssetsNameGroup{Code:201, Msg:"查询失败,分页参数有误", Count:0, Data:nil}
	}
	index := (page-1)*limit
	for i,_:=range result{
		index++
		result[i].ID = index
	}
	return &model.ResponseAssetsWebAssetsNameGroup{Code:0, Msg:"ok", Count:int64(count), Data:result}
}

// SearchReport 渗透测试报告模糊分页查询
func (s *serviceAssets) SearchReport(page, limit int, search interface{})*model.ResponseAssetsReport{
	var result []*model.ResponseAssetsReportInfo
	SearchModel := dao.AssetsReports.Clone()
	searchStr := gconv.String(search)
	if search != ""{
		j := gjson.New(searchStr)
		if gconv.String(j.Get("AttriBution")) != ""{
			SearchModel = SearchModel.Where("attribution like ?", "%"+gconv.String(j.Get("AttriBution"))+"%")
		}
		if gconv.String(j.Get("ManagerName")) != ""{
			SearchModel = SearchModel.Where("manager_name like ?", "%"+gconv.String(j.Get("ManagerName"))+"%")
		}
		if gconv.String(j.Get("AssetsName")) != ""{
			SearchModel = SearchModel.Where("assets_name like ?", "%"+gconv.String(j.Get("AssetsName"))+"%")
		}
		if gconv.String(j.Get("level_name")) != ""{
			SearchModel = SearchModel.Where("level_name like ?", "%"+gconv.String(j.Get("level_name"))+"%")
		}
		if gconv.String(j.Get("level")) != ""{
			SearchModel = SearchModel.Where("level = ?", gconv.String(j.Get("level")))
		}
		if gconv.String(j.Get("level_status")) != ""{
			SearchModel = SearchModel.Where("level_status = ?", gconv.String(j.Get("level_status")))
		}
	}
	count,_ := SearchModel.Count()
	if page > 0 && limit > 0 {
		err := SearchModel.Order("id desc").Limit((page-1)*limit,limit).Scan(&result)
		if err != nil {
			logger.WebLog.Warningf("渗透测试报告查询 数据库错误:%s", err.Error())
			return &model.ResponseAssetsReport{Code:201, Msg:"查询失败,数据库错误", Count:0, Data:nil}
		}
	}else{
		return &model.ResponseAssetsReport{Code:201, Msg:"查询失败,分页参数有误", Count:0, Data:nil}
	}
	index := (page-1)*limit
	for i,_:=range result{
		index++
		result[i].ID1 = result[i].Id
		result[i].Id = uint(index)
	}
	return &model.ResponseAssetsReport{Code:0, Msg:"ok", Count:int64(count), Data:result}
}

// ReportDelete 删除漏洞
func (s *serviceAssets) ReportDelete(r *model.RequestAssetsTypeDelete)error{
	count,err := dao.AssetsReports.Where("id=?",r.ID).Count()
	if err != nil{
		logger.WebLog.Warningf("渗透测试报告-漏洞-删除失败:%s", err.Error())
		return errors.New("删除该漏洞失败,数据库错误")
	}
	if count == 0{
		return errors.New("删除漏洞失败,该漏洞不存在")
	}
	if _,err = dao.AssetsReports.Where("id=?",r.ID).Delete(); err != nil{
		return errors.New(fmt.Sprintf("删除漏洞失败:%s", err.Error()))
	}
	return nil
}

// ExportReport 导出漏洞
func (s *serviceAssets) ExportReport()(*bytes.Buffer, error){
	var result []model.AssetsReports
	err := dao.AssetsReports.Scan(&result)
	if err != nil {
		return &bytes.Buffer{},errors.New("导出渗透测试报告漏洞数据失败,数据库错误")
	}
	if len(result) == 0{
		return &bytes.Buffer{},errors.New("导出渗透测试报告漏洞数据失败,无数据")
	}
	xlsx := excelize.NewFile()
	index := xlsx.NewSheet("漏洞统计表")
	xlsx.DeleteSheet("Sheet1")
	table_name := "漏洞统计表"
	xlsx.SetCellValue(table_name,"A1","应用系统")
	xlsx.SetCellValue(table_name,"B1","安全管理员")
	xlsx.SetCellValue(table_name,"C1","业务系统")
	xlsx.SetCellValue(table_name,"D1","漏洞名称")
	xlsx.SetCellValue(table_name,"E1","漏洞等级")
	xlsx.SetCellValue(table_name,"F1","整改情况")
	xlsx.SetCellValue(table_name,"G1","发现时间")
	for i, info := range result {
		xlsx.SetCellValue(table_name, "A" + strconv.Itoa(i+2), info.Attribution)
		xlsx.SetCellValue(table_name, "B" + strconv.Itoa(i+2), info.ManagerName)
		xlsx.SetCellValue(table_name, "C" + strconv.Itoa(i+2), info.AssetsName)
		xlsx.SetCellValue(table_name, "D" + strconv.Itoa(i+2), info.LevelName)
		if info.Level == 1{
			xlsx.SetCellValue(table_name, "E" + strconv.Itoa(i+2), "高危")
		} else if info.Level == 1{
			xlsx.SetCellValue(table_name, "E" + strconv.Itoa(i+2), "中危")
		} else{
			xlsx.SetCellValue(table_name, "E" + strconv.Itoa(i+2), "低危")
		}
		if info.LevelStatus == 1{
			xlsx.SetCellValue(table_name, "F" + strconv.Itoa(i+2), "已整改")
		} else if info.LevelStatus == 2{
			xlsx.SetCellValue(table_name, "F" + strconv.Itoa(i+2), "未整改")
		} else{
			xlsx.SetCellValue(table_name, "F" + strconv.Itoa(i+2), "已关闭")
		}
		xlsx.SetCellValue(table_name, "G" + strconv.Itoa(i+2), info.FileDate)
	}
	xlsx.SetColWidth(table_name, "A", "A", 20)
	xlsx.SetColWidth(table_name, "B", "B", 16)
	xlsx.SetColWidth(table_name, "C", "C", 20)
	xlsx.SetColWidth(table_name, "D", "D", 20)
	xlsx.SetColWidth(table_name, "E", "E", 15)
	xlsx.SetColWidth(table_name, "F", "F", 15)
	xlsx.SetColWidth(table_name, "G", "G", 15)
	xlsx.SetActiveSheet(index) // 设置工作簿的默认工作表
	var buf bytes.Buffer
	err = xlsx.Write(&buf)
	if err != nil {
		return &bytes.Buffer{},errors.New("导出漏洞数据失败,xlsx流写入失败")
	}
	return &buf,nil
}

// ReportShow 查看指定ID的漏洞报告
func (s *serviceAssets) ReportShow(r *model.RequestAssetsTypeDelete)(*model.AssetsReports,error){
	result,err := dao.AssetsReports.Where("id=?",r.ID).FindOne()
	if err != nil{
		logger.WebLog.Warningf("渗透测试报告-漏洞管理查看资产失败:%s", err.Error())
		return nil,errors.New("查看渗透测试报告漏洞详情失败,数据库错误")
	}
	return result,nil
}

// ReportEdit 修改漏洞报告
func (s *serviceAssets) ReportEdit(r *model.RequestReport)error{
	if len(r.ReportId) != 0 {
		if _,err := dao.AssetsReports.Data(g.Map{"level_name":r.LevelName,"level":r.Level,"level_status":r.LevelStatus}).Where("id",r.ReportId).Update(); err != nil{
			logger.WebLog.Warningf("渗透测试报告-漏洞管理-修改失败:%s", err.Error())
			return errors.New("渗透测试报告-漏洞修改失败,数据库错误")
		}
	}else{
		return errors.New("修改漏洞失败,漏洞ID不存在")
	}
	return nil
}


// TongJiCount 资产数据统计
func (s *serviceAssets) TongJiCount()model.ResponseTongjiInfo{
	var result model.ResponseTongjiInfo
	WebCount,_ := dao.AssetsWeb.Where("1=",1).Count()
	LevelCount,_ := dao.AssetsReports.Where("1=",1).Count()
	LevelYesCount,_ := dao.AssetsReports.Where("level_status=?",1).Count()
	LevelNoCount,_ := dao.AssetsReports.Where("level_status=?",2).Count()
	result.WebCount = WebCount
	result.LevelCount = LevelCount
	result.LevelYesCount = LevelYesCount
	result.LevelNoCount = LevelNoCount
	return result
}

// EchartsInfo Echarts图标统计信息
func (s *serviceAssets) EchartsInfo()*model.ResponseEchartsInfo{
	var result1 []model.UResponseEchartsInfoManagerLevel
	var result2 []model.UResponseEchartsInfoLevel
	var result3 []model.UResponseEchartsInfoAssetsName
	err := dao.AssetsReports.Fields("COUNT(level_name) Number, level_name").Group("level_name").Limit(10).Scan(&result2)
	if err != nil {
		return nil
	}
	err = dao.AssetsReports.Fields("COUNT(manager_name) Number, manager_name").Group("manager_name").Limit(10).Scan(&result1)
	if err != nil {
		return nil
	}
	err = dao.AssetsReports.Fields("COUNT(assets_name) Number, assets_name").Group("assets_name").Limit(10).Scan(&result3)
	if err != nil {
		return nil
	}
	return &model.ResponseEchartsInfo{Code:200,Msg:"ok",Data:result1, Data1:result2, Data2:result3}
}