package domain

import (
	"github.com/gohouse/gorose"
	"github.com/goinggo/mapstructure"
	"task_dispatcher/common"
	"task_dispatcher/config"
	"task_dispatcher/enum"
	"time"
)

const (
	EXECUTE_RECORD_SELECT_FIELD = "id,task_no as taskNo,task_name as taskName,group_no as groupNo,status,error_msg as errorMsg,create_time as createTime"
)

/**
任务执行记录
*/
type TaskExecuteRecord struct {
	Id         int
	TaskNo     string `json:"taskNo"`
	TaskName   string `json:"taskName"`
	GroupNo    string `json:"groupNo"`
	Status     int    `json:"status"`
	ErrorMsg   string `json:"errorMsg"`
	Result     string `json:"result"`
	CreateTime string `json:"createTime"`
	IsDelete   int    `json:"isDelete"`

	StatusNameUI string `json:"statusNameUI"` //页面回显时使用  状态名称 0 失败  1 成功
}

/**
添加任务执行记录
*/
func (ter *TaskExecuteRecord) AddTaskExecuteRecord() (int64, error) {

	ter.CreateTime = time.Now().Format(common.YYYY_MM_DD_HH_MM_SS)
	ter.IsDelete = 0
	res, err := config.Conn.Table("task_execute_record").Data(map[string]interface{}{
		"task_no":     ter.TaskNo,
		"task_name":   ter.TaskName,
		"group_no":    ter.GroupNo,
		"status":      ter.Status,
		"error_msg":   ter.ErrorMsg,
		"result":      ter.Result,
		"create_time": ter.CreateTime,
		"is_delete":   ter.IsDelete,
	}).Insert()
	if err != nil {
		common.GetLog().Errorln("插入日志记录错误", err)
	}
	return res, err
}

/**
根据任务查询任务执行记录
*/
func FindTaskExecuteRecordByTaskNo(taskNo interface{}) (*[]TaskExecuteRecord, error) {
	res, err := config.Conn.Table("task_execute_record").Fields(EXECUTE_RECORD_SELECT_FIELD).Where("task_no", taskNo).Get()
	taskExecuteRecords := &[]TaskExecuteRecord{}
	err1 := mapstructure.Decode(res, taskExecuteRecords)
	if err1 == nil {
		return nil, err1
	}
	return taskExecuteRecords, err
}

/**
根据任务边和状态查询任务执行记录
*/
func FindTaskExecuteRecordByTaskNoAndStatus(taskNo, status interface{}) (*[]TaskExecuteRecord, error) {
	res, err := config.Conn.Table("task_execute_record").Fields(EXECUTE_RECORD_SELECT_FIELD).Where("task_no", taskNo).Where("status", status).Get()
	taskExecuteRecords := &[]TaskExecuteRecord{}
	err1 := mapstructure.Decode(res, taskExecuteRecords)
	if err1 == nil {
		return nil, err1
	}
	return taskExecuteRecords, err
}

/**
分页查询
*/
func FindTaskExecuteRecordByPage(page *common.Page, taskExecuteRecord *TaskExecuteRecord) (*common.Page, error) {
	ter := config.Conn.Table("task_execute_record")
	//查询数据
	taskExecuteRecordWhereCondition(ter, taskExecuteRecord)
	res, err1 := ter.Fields(EXECUTE_RECORD_SELECT_FIELD).Order("id desc").Offset((page.Page - 1) * page.Limit).Limit(page.Limit).Get()
	//查询总数
	terCount := config.Conn.Table("task_execute_record")
	count, err2 := terCount.Where("is_delete", 0).Count()
	if err1 != nil {

		return page, err1
	}
	if err2 != nil {

		return page, err2
	}
	taskExecuteRecordsResult := &[]TaskExecuteRecord{}
	mapstructure.Decode(res, taskExecuteRecordsResult)
	for index, ter := range *taskExecuteRecordsResult {
		(*taskExecuteRecordsResult)[index].StatusNameUI = enum.FindTaskExecuteRecordStatusNameByCode(ter.Status)
	}
	page.Data = taskExecuteRecordsResult
	page.Total = count
	return page, nil
}
func taskExecuteRecordWhereCondition(tableName *gorose.Session, record *TaskExecuteRecord) {
	tableName.Where("is_delete", 0)
	if record.Status != 0 {
		tableName.Where("status", record.Status)
	}
	if record.TaskNo != "" {
		tableName.Where("task_no", record.TaskNo)
	}
	if record.GroupNo != "" {
		tableName.Where("group_no", record.GroupNo)
	}
	if record.Id != 0 {
		tableName.Where("id", record.Id)
	}
}
