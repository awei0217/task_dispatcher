package domain

import (
	"task_dispatcher/common"
	"task_dispatcher/config"
	"github.com/gohouse/gorose"
	"github.com/goinggo/mapstructure"
	"github.com/kataras/iris/core/errors"
	"math/rand"
	"strconv"
	"time"
	"task_dispatcher/enum"
)

const (
	//数据库需要查询的字段
	TASK_SELECT_FIELD = "id,task_no as taskNo,group_no as groupNo,name,task_type as taskType,url,is_activity as isActivity,is_record_log as isRecordLog," +
		"param,cron,mail,description,concurrency_num as concurrencyNum,task_slice as taskSlice,create_time as createTime,status,is_region as isRegion,request_method as requestMethod"
)

/**
任务实体
*/
type Task struct {
	Id             int    `json:"id"`             //任务ID
	IsDelete       int    `json:"isDelete"`       //是否删除 0否 1是
	TaskNo         string `json:"taskNo"`         //任务编号
	GroupNo        string `json:"groupNo"`        //分组编号
	Name           string `json:"name"`           //任务名称
	TaskSlice      int    `json:"taskSlice"`      //任务分片 (0~127) 也就是说调度中心的机器数量最多为128台
	Description    string `json:"description"`    // 任务描述
	Url            string `json:"url"`            //目标地址
	ConcurrencyNum int    `json:"concurrencyNum"` // 并发数量（一次启多少协成执行）
	Status 		  int `json:"status"`             // 任务状态
	IsActivity     int    `json:"isActivity"`     //是否激活 0 否 1 是
	IsRecordLog  int    `json:"isRecordLog"`     //是否记录日志 0 否 1 是
	Param          string `json:"param"`          //传递给任务的参数 json串
	TaskType       int    `json:"taskType"`       //任务分类 1 定时任务 2 一次性任务
	Cron           string `json:"cron"`           // 定时任务表达式
	CreateTime     string `json:"createTime"`
	UpdateTime     string `json:"updateTime"`
	CreateUser     string `json:"createUser"`
	UpdateUser     string `json:"updateUser"`
	Mail           string `json:"mail"` //邮箱
	IsRegion	   int    `json:"isRegion"` //是否分区 1 否 2 是
	RequestMethod  int 	  `json:"requestMethod"` // 请求方式

	GroupNameUI    string `json:"groupNameUI"`    //页面回显时使用  分组名称
	TaskTypeNameUI string `json:"taskTypeNameUI"` //页面回显时使用 任务类型名称
	TaskIsActivityNameUI string `json:"taskIsActivityNameUI"` //页面回显时使用 任务是否运行
	TaskIsRecordLogNameUI string `json:"taskIsRecordLogNameUI"` //页面回显时使用 任务是否运行
	StatusNameUI string `json:"statusNameUI"` //任务启动 停止状态
	IsRegionNameUI string `json:"isRegionNameUI"`
}

/**
添加一条任务
*/
func (task *Task) AddTask() (int64, error) {
	dbRes, _ := config.Conn.Table("task").Fields(TASK_SELECT_FIELD).Where("is_delete", 0).Where("url", task.Url).Get()
	if len(dbRes) != 0 {
		return 0, errors.New("此任务已存在")
	}
	id := time.Now().Nanosecond()
	task.TaskNo = strconv.Itoa(id)
	task.TaskSlice = rand.Intn(127) // 返回 0 到 127 的随机数 ，包含0和127
	task.CreateTime = time.Now().Format(common.YYYY_MM_DD_HH_MM_SS)
	task.UpdateTime = time.Now().Format(common.YYYY_MM_DD_HH_MM_SS)
	res, err := config.Conn.Table("task").Data(map[string]interface{}{
		"task_no":         task.TaskNo,
		"group_no":        task.GroupNo,
		"name":            task.Name,
		"task_slice":      task.TaskSlice,
		"description":     task.Description,
		"url":             task.Url,
		"concurrency_num": task.ConcurrencyNum,
		"is_activity":     task.IsActivity,
		"is_record_log":     task.IsRecordLog,
		"param":           task.Param,
		"task_type":       task.TaskType,
		"cron":            task.Cron,
		"create_time":     task.CreateTime,
		"update_time":     task.UpdateTime,
		"create_user":     task.CreateUser,
		"update_user":     task.UpdateUser,
		"mail":            task.Mail,
		"is_region":	   task.IsRegion,
		"request_method":   task.RequestMethod,
	}).Insert()
	return res, err
}

/**
根据ID修改任务
*/
func (task *Task) UpdateTask() (int64, error) {

	task.UpdateTime = time.Now().Format(common.YYYY_MM_DD_HH_MM_SS)
	res, err := config.Conn.Table("task").Where("id", task.Id).Data(map[string]interface{}{
		"group_no":        task.GroupNo,
		"name":            task.Name,
		"task_slice":      task.TaskSlice,
		"description":     task.Description,
		"url":             task.Url,
		"concurrency_num": task.ConcurrencyNum,
		"is_activity":     task.IsActivity,
		"status":     task.Status,
		"is_record_log":     task.IsRecordLog,
		"param":           task.Param,
		"cron":            task.Cron,
		"update_time":     task.UpdateTime,
		"update_user":     task.UpdateUser,
		"mail":            task.Mail,
		"task_type":            task.TaskType,
		"is_region":		task.IsRegion,
		"request_method":   task.RequestMethod,
	}).Update()
	return res, err
}

/**
根据条件查询任务
*/
func FindAllTask(task *Task) (*[]Task, error) {
	tk := config.Conn.Table("task")
	res, err := tk.Fields(TASK_SELECT_FIELD).Where(func() { taskWhereCondition(tk, task) }).Get()
	tasks := &[]Task{}
	err1 := mapstructure.Decode(res, tasks)
	if err1 != nil {
		common.GetLog().Errorln("map转task结构体异常", err1)
		return tasks, err1
	}
	return tasks, err
}

/**
分页查询
*/
func FindTaskByPage(page *common.Page, task *Task) (*common.Page, error) {
	tk := config.Conn.Table("task")
	//查询数据
	res, err1 := tk.Fields(TASK_SELECT_FIELD).Where(func() { taskWhereCondition(tk, task) }).Offset((page.Page - 1) * page.Limit).Limit(page.Limit).Get()
	//查询总数
	tkCount := config.Conn.Table("task")
	count, err2 := tkCount.Where(func() { taskWhereCondition(tkCount, task) }).Count()
	if err1 != nil {
		common.GetLog().Errorln("根据条查询任务记录异常", err1)
		return page, err1
	}
	if err2 != nil {
		common.GetLog().Errorln("根据条查询任务总数异常", err1)
		return page, err2
	}
	tasksResult := &[]Task{}
	mapstructure.Decode(res, tasksResult)
	for index,tr:= range *tasksResult{
		(*tasksResult)[index].GroupNameUI = enum.FindTaskGroupNameByCode(tr.GroupNo)
		(*tasksResult)[index].TaskTypeNameUI = enum.FindTaskTypeNameByCode(tr.TaskType)
		(*tasksResult)[index].TaskIsActivityNameUI = enum.FindTaskIsActivityNameByCode(tr.IsActivity)
		(*tasksResult)[index].TaskIsRecordLogNameUI = enum.FindTaskIsRecordLogNameByCode(tr.IsActivity)
		(*tasksResult)[index].StatusNameUI = enum.FindTaskStatusNameByCode(tr.Status)
		(*tasksResult)[index].IsRegionNameUI = enum.FindTaskIsRegionByCode(tr.IsRegion)
	}

	page.Data = tasksResult
	page.Total = count
	return page, nil
}

/**
根据任务切片来查寻任务
*/
func FindByTaskSlice(taskSlice []interface{}) (*[]Task, error) {
	res, err := config.Conn.Table("task").Fields(TASK_SELECT_FIELD).Where("task_slice", "in", taskSlice).Get()
	tasks := &[]Task{}
	mapstructure.Decode(res, tasks)
	return tasks, err
}
func taskWhereCondition(tableName *gorose.Session, task *Task) {
	tableName.Where("is_delete", 0)
	if task.Name != "" {
		tableName.Where("name", task.Name)
	}
	if task.TaskNo != "" {
		tableName.Where("task_no", task.TaskNo)
	}
	if task.GroupNo != "" {
		tableName.Where("group_no", task.GroupNo)
	}
	if task.Id != 0 {
		tableName.Where("id", task.Id)
	}
}

func FindTaskById(id int) (*Task, error) {
	res, err := config.Conn.Table("task").Fields(TASK_SELECT_FIELD).Where("id", id).Get()
	if err != nil {
		common.GetLog().Errorln("根据ID查询任务异常", err)
		return nil, err
	}
	task := &[]Task{}
	// 将map转化成struct
	if res == nil{
		return nil,nil
	}
	mapstructure.Decode(res, task)
	return &((*task)[0]), err
}

func FindTaskByUrl(url string) (*Task, error) {
	res, err := config.Conn.Table("task").Fields(TASK_SELECT_FIELD).Where("url", url).Get()
	if err != nil {
		common.GetLog().Errorln("根据URL查询任务异常", err)
		return nil, err
	}

	task := &[]Task{}
	// 将map转化成struct
	mapstructure.Decode(res, task)

	return &((*task)[0]), err
}
/**
	删除任务
 */
func DeleteTask(id int)(int64,error){
	res, err :=config.Conn.Table("task").Where("id",id).Delete()
	if err != nil{
		return 0,err
	}
	return res,nil
}
/**
	统计总任务数量
 */
func CountTask()(int64, error){
	tkCount := config.Conn.Table("task")
	count, err := tkCount.Where().Count()
	return count,err
}

/**
	统计执行中务数量
 */
func CountExecutingTask()(int64, error){
	tkCount := config.Conn.Table("task")
	count, err := tkCount.Where("is_activity",1).Count()
	return count,err
}

/**
	统计定时任务数量
 */
func CountCronTask()(int64, error){
	tkCount := config.Conn.Table("task")
	count, err := tkCount.Where("task_type",1).Count()
	return count,err
}
/**
	统计普通任务数量
 */
func CountOneTask()(int64, error){
	tkCount := config.Conn.Table("task")
	count, err := tkCount.Where("task_type",2).Count()
	return count,err
}
/**
	跟新任务为启动未运行的状态
 */
func UpdateTaskByTaskSlice( taskSlice []interface{})  {
	config.Conn.Table("task").Where("task_slice","in",taskSlice).Where("task_type",1).Data(map[string]interface{}{
		"is_activity":2,
	}).Update()
}




