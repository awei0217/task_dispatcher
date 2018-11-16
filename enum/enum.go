package enum

// 任务类型枚举
const (
	TASK_TYPE_CRON     = 1
	TASK_TYPE_ONE_TIME = 2
)

// 任务分区枚举
const (
	TASK_IS_REGION_NO  = 1 // 否
	TASK_IS_REGION_YES = 2 // 是
)

// 任务分组枚举
const (
	TASK_TGROUP_CANGCHU = "1"
	TASK_TGROUP_TC      = "2"
	TASK_TGROUP_DAJIAN  = "3"
)

// 任务运行停止枚举
const (
	TASK_IS_ACTIVITY_YES = 1 // 运行中
	TASK_IS_ACTIVITY_NO  = 2 // 停止运行
	TASK_IS_ACTIVITY_NOT = 3 // 默认(手动任务)
)

// 任务启动停止状态枚举
const (
	TASK_STATUS_START = 1 // 启动
	TASK_STATUS_STOP  = 2 // 停止
	TASK_STATUS_NOT   = 3 // 默认 (手动任务)
)

// 任务执行记录状态枚举
const (
	TASK_EXECUTE_RECORD_FAIL    = 0
	TASK_EXECUTE_RECORD_SUCCESS = 1
)

// 是否记录日志枚举
const (
	TASK_RECORD_LOG_NO  = 0
	TASK_RECORD_LOG_YES = 1
)

// 注册中心是否死亡枚举
const (
	REGISTER_CENTER_DEATH_NO  = 0
	REGISTER_CENTER_DEATH_YES = 1
)

var taskType = map[int]string{
	TASK_TYPE_CRON:     "定时任务",
	TASK_TYPE_ONE_TIME: "手动任务",
}

var isRegion = map[int]string{
	TASK_IS_REGION_NO:  "否",
	TASK_IS_REGION_YES: "是",
}
var taskGroup = map[string]string{
	TASK_TGROUP_CANGCHU: "仓储计费任务",
	TASK_TGROUP_TC:      "TC计费任务",
	TASK_TGROUP_DAJIAN:  "大件计费任务",
}

var taskIsActivity = map[int]string{
	TASK_IS_ACTIVITY_YES: "是",
	TASK_IS_ACTIVITY_NO:  "否",
	TASK_IS_ACTIVITY_NOT: "默认",
}

var taskIsRecordLog = map[int]string{
	TASK_RECORD_LOG_NO:  "否",
	TASK_RECORD_LOG_YES: "是",
}

var taskExecuteRecordStatus = map[int]string{
	TASK_EXECUTE_RECORD_FAIL:    "失败",
	TASK_EXECUTE_RECORD_SUCCESS: "成功",
}

var registerCenterDeath = map[int]string{
	REGISTER_CENTER_DEATH_NO:  "否",
	REGISTER_CENTER_DEATH_YES: "是",
}

var taskStatus = map[int]string{
	TASK_STATUS_START: "启动",
	TASK_STATUS_STOP:  "停止",
	TASK_STATUS_NOT:   "默认",
}

func FindTaskTypeNameByCode(code int) string {
	str, ok := taskType[code]
	if ok {
		return str
	}
	return ""
}

func FindTaskGroupNameByCode(code string) string {
	str, ok := taskGroup[code]
	if ok {
		return str
	}
	return ""
}
func FindTaskIsActivityNameByCode(code int) string {
	str, ok := taskIsActivity[code]
	if ok {
		return str
	}
	return ""
}

func FindTaskExecuteRecordStatusNameByCode(code int) string {
	str, ok := taskExecuteRecordStatus[code]
	if ok {
		return str
	}
	return ""
}

func FindTaskIsRecordLogNameByCode(code int) string {
	str, ok := taskIsRecordLog[code]
	if ok {
		return str
	}
	return ""
}

func FindRegisterCenterDeathNameByCode(code int) string {
	str, ok := registerCenterDeath[code]
	if ok {
		return str
	}
	return ""
}

func FindTaskStatusNameByCode(code int) string {
	str, ok := taskStatus[code]
	if ok {
		return str
	}
	return ""
}

func FindTaskIsRegionByCode(code int) string {
	str, ok := isRegion[code]
	if ok {
		return str
	}
	return ""
}
