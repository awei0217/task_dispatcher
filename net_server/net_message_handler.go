package net_server

import (
	"net"
	"strconv"
	"strings"
	"task_dispatcher/common"
	"task_dispatcher/domain"
	"task_dispatcher/task_manager"
)

func HandlerNetMessage(message string, conn net.Conn) {

	if message == "NeedTasks" {

		// 判断当前机器的负载，是否需要剥离一些任务出去
		// 查询当前机器的最近半小时的 cpu 内存，进行统计求平均值
		performanceMonitor, _ := domain.FindPerformanceMonitor(common.GetIp(), 1000)
		var cpuCount, memCount float64
		for _, p := range *performanceMonitor {
			cpu, _ := strconv.ParseFloat(p.UseCpu, 64)
			mem, _ := strconv.ParseFloat(p.UseMemory, 64)
			cpuCount = cpuCount + cpu
			memCount = memCount + mem
		}
		if (int(cpuCount)/len(*performanceMonitor)) < 50 && (int(memCount)/len(*performanceMonitor)) < 50 {
			common.GetLog().Println(common.GetIp(), "当前负载过低,不给", conn.RemoteAddr().String(), "分配任务")
			return
		}
		//获取当前机器的任务分片
		rc, _ := domain.FindRegisterCenterByIp(common.GetIp())
		taskSlice := rc.TaskSliceCollection
		if taskSlice == "" {
			return
		}
		tss := strings.Split(taskSlice, ",")
		newSlice := make([]interface{}, 0)
		newSliceString := make([]string, 0)
		oldSlice := make([]string, 0)

		currentValueSlice := len(tss)
		for i := 0; i < currentValueSlice; i++ {
			if i < (currentValueSlice / 2) {
				newSlice = append(newSlice, tss[i])
				newSliceString = append(newSliceString, tss[i])
			} else {
				oldSlice = append(oldSlice, tss[i])
			}
		}
		rc.TaskSliceCollection = strings.Join(oldSlice, ",")
		rc.UpdateRegisterCenter()
		tasks, _ := domain.FindByTaskSlice(newSlice)
		if len(*tasks) != 0 {
			for _, t := range *tasks {
				task_manager.GetTaskManger().StopTask(t.Id)
			}
		}
		ip := strings.Split(conn.RemoteAddr().String(), ":")[0]
		rcDb, _ := domain.FindRegisterCenterByIp(ip)
		if rcDb.TaskSliceCollection == "" {
			rcDb.TaskSliceCollection = strings.Join(newSliceString, ",")
		} else {
			rcDb.TaskSliceCollection = rcDb.TaskSliceCollection + "," + strings.Join(newSliceString, ",")
		}
		_, err := rcDb.UpdateRegisterCenter()
		if err == nil {
			common.GetLog().Info("给", ip, "任务分配成功")
		} else {
			common.GetLog().Error("给", ip, "任务分配失败")
		}
	}
}
