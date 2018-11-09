package domain

import (
	"task_dispatcher/config"
	"time"
)

/**
任务分组
*/
type TaskGroup struct {
	Id         int
	GroupNo    string
	GroupName  string
	CreateTime time.Time
	UpdateTime time.Time
	CreateUser time.Time
	UpdateUser time.Time
	IsDelete   int //是否删除 0 否 1 是
}

/**
	添加一个任务分组
*/
func (tg *TaskGroup) AddTaskGroup() (int64, error) {
	tg.Id = time.Now().Nanosecond()
	tg.CreateTime = time.Now()
	tg.UpdateTime = time.Now()
	tg.IsDelete = 0
	res, err := config.Conn.Table("task_group").Data(map[string]interface{}{
		"id":          tg.Id,
		"group_no":    tg.GroupNo,
		"group_name":  tg.GroupName,
		"create_time": tg.CreateTime,
		"update_time": tg.UpdateTime,
		"create_user": tg.CreateUser,
		"update_user": tg.UpdateUser,
		"is_delete":   tg.IsDelete,
	}).Insert()
	return res, err
}

/**
	根据ID更新任务分组
*/
func (tg *TaskGroup) UpdateTaskGroupById() (int64, error) {
	tg.UpdateTime = time.Now()
	res, err := config.Conn.Table("task_group").Where("id", tg.Id).Data(map[string]interface{}{
		"group_no":    tg.GroupNo,
		"group_name":  tg.GroupName,
		"update_time": tg.UpdateTime,
		"update_user": tg.UpdateUser,
		"is_delete":   tg.IsDelete,
	}).Update()
	return res, err
}
