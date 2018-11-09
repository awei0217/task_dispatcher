package domain

import (
	"task_dispatcher/config"
	"log"
	"github.com/goinggo/mapstructure"
	"time"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"strconv"
	"task_dispatcher/common"
)
const (
	PERFORMANCE_MONITOR_SELECT_FILED= "ip,use_cpu as useCpu,use_memory as useMemory,update_time as updateTime"
)
/**
	性能监控结构体
 */
type PerformanceMonitor struct {
	Id int  `json:"id"`
	Ip string `json:"ip"`
	UseCpu string  `json:"useCpu"`
	UseMemory string `json:"useMemory"`
	UpdateTime string  `json:"updateTime"`
}

func InsertMonitorInfo()(){
	for{
		cpuLocal,_ := cpu.Percent(1 * time.Second,false);
		useCpu := strconv.FormatFloat(cpuLocal[0],'G',4,64)
		vm,_ := mem.VirtualMemory();
		useMemory := strconv.FormatFloat(vm.UsedPercent,'G',4,64)
		updateTime := time.Now().Format(common.YYYY_MM_DD_HH_MM_SS)
		ip := common.GetIp()
		_,err := config.Conn.Table("performance_monitor").Data(map[string]interface{}{
			"ip":ip,
			"use_cpu":useCpu,
			"use_memory":useMemory,
			"update_time":updateTime,
		}).Insert();
		if err != nil {
			log.Println("新增监控信息失败",err)
		}
		if cpuLocal[0] > 90 || vm.UsedPercent > 80{
			email := &common.Email{
				To:"sunpengwei1992@aliyun.com",
				Subject:"统一任务调度系统负载过高邮件推送",
				Message:ip+":CPU和内存负载过高，cpu百分比:"+useCpu+" 内存百分比:"+useMemory,
			}
			go common.SendMail(email)
		}
		time.Sleep(10 * time.Second)
	}
}

func FindPerformanceMonitor(ip string,limit int) (*[]PerformanceMonitor, error) {
	rcCon := config.Conn.Table("performance_monitor")
	//查询数据
	res, err1 := rcCon.Fields(PERFORMANCE_MONITOR_SELECT_FILED).
		Where("ip",ip).
		Order("id desc").
		Offset(0).
		Limit(limit).Get()
	if err1 != nil {
		log.Fatalln("查询性能监控信息异常", err1)
		return nil, err1
	}
	if res == nil {
		return nil,nil
	}
	pmResult := &[]PerformanceMonitor{}
	mapstructure.Decode(res, pmResult)
	return pmResult, nil
}

func FindPerformanceMonitorByIpAndTime(ip string,startTime,endTime string) (*[]PerformanceMonitor, error) {
	rcCon := config.Conn.Table("performance_monitor")
	//查询数据
	res, err1 := rcCon.Fields(PERFORMANCE_MONITOR_SELECT_FILED).
		Where("update_time",">=",startTime,).
		Where("update_time","<=",endTime).
		Order("id desc").
		Get()
	if err1 != nil {
		log.Fatalln("查询性能监控信息异常", err1)
		return nil, err1
	}
	if res == nil {
		return nil,nil
	}
	pmResult := &[]PerformanceMonitor{}
	mapstructure.Decode(res, pmResult)
	return pmResult, nil
}

func init()  {
	go InsertMonitorInfo()
}

