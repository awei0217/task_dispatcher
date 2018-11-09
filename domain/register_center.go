package domain

import (
	"task_dispatcher/common"
	"task_dispatcher/config"
	"task_dispatcher/enum"
	"task_dispatcher/net_client"
	"github.com/gohouse/gorose"
	"github.com/goinggo/mapstructure"
	"github.com/kataras/iris/core/errors"
	"strconv"
	"strings"
	"time"
)

const (
	REGISTE_CENTER_SELECT_FILED = "id,ip,register_time as registerTime,update_time as updateTime,is_death as isDeath,task_slice_collection as taskSliceCollection"
)

type RegisterCenter struct {
	Id                  int    `json:"id"`                   // 主键
	Ip                  string `json:"ip"`                   // Ip地址
	RegisterTime        string `json:"registerTime"`         // 注册时间
	UpdateTime          string `json:"updateTime"`           //更新时间
	IsDeath             int    `json:"isDeath"`              //是否死亡 0 否 1是
	TaskSliceCollection string `json:"taskSliceCollection"'` // 本机扫描这些分片集合的任务
	IsDelete            int    `json:"isDelete"`

	IsDeathName string `json:"isDeathName"` // 是否死亡的名称
}

/**
注册
*/
func (rc *RegisterCenter) Register() (int64, error) {
	dbRc, registerErr := FindRegisterCenterByIp(rc.Ip)
	if registerErr != nil {
		return 0, registerErr
	}
	s := strings.Replace(rc.Ip, ".", "_", -1)
	rc.TaskSliceCollection, _ = common.GetYmlFile().Get(s)
	if dbRc != nil {
		res, updateErr := config.Conn.Table("register_center").Where("ip", rc.Ip).Data(map[string]interface{}{
			"update_time": time.Now().Format(common.YYYY_MM_DD_HH_MM_SS),
			"is_death":    0,
		}).Update()
		// 任务分片为空,向全网发送消息，给我点任务
		if dbRc.TaskSliceCollection == "" {
			common.GetLog().Println(common.GetIp(), "开始向全网发送索要任务的消息")
			go AllNetSendNeedTask()
		}
		return res, updateErr
	} else {
		rc.RegisterTime = time.Now().Format(common.YYYY_MM_DD_HH_MM_SS)
		rc.UpdateTime = time.Now().Format(common.YYYY_MM_DD_HH_MM_SS)
		rc.IsDeath = 0
		rc.IsDelete = 0
		res, insertErr := config.Conn.Table("register_center").Data(map[string]interface{}{
			"ip":                    rc.Ip,
			"register_time":         rc.RegisterTime,
			"update_time":           rc.UpdateTime,
			"task_slice_collection": rc.TaskSliceCollection,
			"is_death":              rc.IsDeath,
			"is_delete":             rc.IsDelete,
		}).Insert()
		// 任务分片为空,向全网发送消息，给我点任务
		if rc.TaskSliceCollection == "" {
			go AllNetSendNeedTask()
		}
		return res, insertErr
	}
}

func (rc *RegisterCenter) UpdateRegisterCenter() (int64, error) {
	return config.Conn.Table("register_center").Where("ip", rc.Ip).Data(map[string]interface{}{
		"update_time":           time.Now().Format(common.YYYY_MM_DD_HH_MM_SS),
		"is_death":              rc.IsDeath,
		"task_slice_collection": rc.TaskSliceCollection,
	}).Update()
}

/**
根据IP往注册中心发送心跳，更新时间
*/
func HeartBeat() {
	for {
		// 每隔 90 秒心跳一次
		time.Sleep(90 * time.Second)
		tempAds := common.GetIp()
		_, err := config.Conn.Table("register_center").Where("ip", tempAds).Data(map[string]interface{}{
			"update_time": time.Now().Format(common.YYYY_MM_DD_HH_MM_SS),
			"is_death":    0,
		}).Update()
		if err != nil {
			common.GetLog().Errorln(tempAds, "心跳失败", err)
		}
		//每隔2分钟心跳一次
		//同时检测是否有死亡的机器
		rc, findErr := config.Conn.Table("register_center").Fields(REGISTE_CENTER_SELECT_FILED).Get()
		if findErr != nil {
			common.GetLog().Errorln(tempAds, "查询注册中心检测死亡机器失败", err)
		}
		// 循环检测是否有死亡的机器
		if len(rc) == 0 {
			time.Sleep(120 * time.Second)
			continue
		}
		for _, n := range rc {
			ut, _ := time.Parse(common.YYYY_MM_DD_HH_MM_SS, n["updateTime"].(string))
			registerCenter := &RegisterCenter{}
			mapstructure.Decode(n, registerCenter)
			// 检测本机器任务分片是否为空，为空的话，向全网发送需要任务的消息
			if registerCenter.Ip == common.GetIp() {
				if registerCenter.TaskSliceCollection == "" {
					common.GetLog().Println(common.GetIp(), "开始向全网发送索要任务的消息")
					go AllNetSendNeedTask()
				}
			}
			t, _ := time.Parse(common.YYYY_MM_DD_HH_MM_SS, time.Now().Format(common.YYYY_MM_DD_HH_MM_SS))
			if t.Unix()-ut.Unix() < 60*3 { // 3分钟 说明机器活着，进入下一个循环
				continue
			}
			// 将长时间没有心跳的机器更新为死亡
			registerCenter.SetDeath()
			email := &common.Email{
				To:      "sunpengwei1992@aliyun.com",
				Subject: "统一任务调度系统节点死亡邮件推送",
				Message: registerCenter.Ip + ":节点死亡",
			}
			go common.SendMail(email)
			if registerCenter.TaskSliceCollection == "" { //说明有应用程序死亡了，但死亡应用程序上的任务已经被分配，进入下一个循环
				continue
			}
			// 向死亡机器发送网络通信，检测是否真的死亡，还是假死亡
			count, err := net_client.SendNetInfo(registerCenter.Ip, "HeartBeat")
			if err == nil && count != 0 { // 说明发送成功,应用程序假死亡（可能只是连不上数据库）
				continue
			}
			// 网络通信中断,说明应用程序确实死亡
			// 抢死亡占机器上的任务
			db := config.Conn.Table("register_center")

			// 开启事务，自动提交事务
			db.Transaction(func() error {
				// 查询死亡机器的注册中心
				deathRcMap, err := db.Query("select ip,task_slice_collection as taskSliceCollection,is_death as isDeath from register_center where ip = " + "'" + registerCenter.Ip + "'" + " for update")
				deathRcStructs := &[]RegisterCenter{}
				mapstructure.Decode(deathRcMap, deathRcStructs)
				dr1 := (*deathRcStructs)[0]
				if dr1.TaskSliceCollection == "" {
					return err
				}
				// 查询当前机器的最近半小时的 cpu 内存，进行统计求平均值
				performanceMonitor, _ := FindPerformanceMonitor(common.GetIp(), 1000)
				var cpuCount, memCount float64
				for _, p := range *performanceMonitor {
					cpu, _ := strconv.ParseFloat(p.UseCpu, 64)
					mem, _ := strconv.ParseFloat(p.UseMemory, 64)
					cpuCount = cpuCount + cpu
					memCount = memCount + mem
				}
				if ((int(cpuCount) / len(*performanceMonitor)) > 50) || ((int(memCount) / len(*performanceMonitor)) > 50) {
					common.GetLog().Println(common.GetIp(), "当前负载过高,不索要死亡机器的任务")
					return nil
				}
				// 如果平均值小于50，把死亡机器的任务更新到当前机器上
				//更新任务死亡机器上执行的任务的 status 为停止 isActivity 为 否
				deathTaskSlice := strings.Split(dr1.TaskSliceCollection, ",")
				newSls := make([]interface{}, len(deathTaskSlice))
				for i, v := range deathTaskSlice {
					newSls[i] = v
				}
				UpdateTaskByTaskSlice(newSls)
				// 把死亡机器上的任务切边集合更新到当前机器
				currentRegisterCenter, _ := FindRegisterCenterByIp(common.GetIp())
				if currentRegisterCenter.TaskSliceCollection == "" {
					currentRegisterCenter.TaskSliceCollection = dr1.TaskSliceCollection

				} else {
					currentRegisterCenter.TaskSliceCollection = currentRegisterCenter.TaskSliceCollection + "," + dr1.TaskSliceCollection
				}
				currentRegisterCenter.UpdateRegisterCenter()
				// 更新死亡机器的任务切片为空
				db.Table("register_center").Data(map[string]interface{}{"task_slice_collection": ""}).Where("ip", dr1.Ip).Update()
				return err
			})
		}
	}
}

func FindRegisterCenterByIp(ip string) (*RegisterCenter, error) {
	if ip == "" {
		return nil, errors.New("ip is not null")
	}
	rc, err := config.Conn.Table("register_center").Fields(REGISTE_CENTER_SELECT_FILED).Where("ip", ip).Get()
	if err != nil {
		common.GetLog().Errorln("根据IP获取注册中心记录失败", err)
		return nil, err
	}
	if len(rc) != 0 {
		registerCenter := &[]RegisterCenter{}
		mapstructure.Decode(rc, registerCenter)
		return &((*registerCenter)[0]), err
	}
	return nil, err
}

/**
分页查询
*/
func FindRegisterCenterByPage(page *common.Page, rc *RegisterCenter) (*common.Page, error) {
	rcCon := config.Conn.Table("register_center")
	//查询数据
	res, err1 := rcCon.Fields(REGISTE_CENTER_SELECT_FILED).Where(func() { rcWhereCondition(rcCon, rc) }).Offset((page.Page - 1) * page.Limit).Limit(page.Limit).Get()
	//查询总数
	rcConCount := config.Conn.Table("register_center")
	count, err2 := rcConCount.Where(func() { rcWhereCondition(rcConCount, rc) }).Count()
	if err1 != nil {
		common.GetLog().Errorln("根据条件查询注册中心记录异常", err1)
		return page, err1
	}
	if err2 != nil {
		common.GetLog().Errorln("根据条件查询注册中心总数异常", err1)
		return page, err2
	}
	rcResult := &[]RegisterCenter{}
	mapstructure.Decode(res, rcResult)

	for index, tr := range *rcResult {
		(*rcResult)[index].IsDeathName = enum.FindRegisterCenterDeathNameByCode(tr.IsDeath)
	}
	page.Data = rcResult
	page.Total = count
	return page, nil
}

func rcWhereCondition(tableName *gorose.Session, rc *RegisterCenter) {
	tableName.Where("is_delete", 0)
	if rc.Ip != "" {
		tableName.Where("ip", rc.Ip)
	}
}

/**
根据IP设置机器死亡
*/
func (rc *RegisterCenter) SetDeath() (int64, error) {
	res, err := config.Conn.Table("register_center").Where("ip", rc.Ip).Data(map[string]interface{}{
		"is_death": 1,
	}).Update()
	return res, err
}

/**
根据IP设置机器死亡
*/
func (rc *RegisterCenter) SetLive() (int64, error) {
	res, err := config.Conn.Table("register_center").Where("ip", rc.Ip).Data(map[string]interface{}{
		"is_death": 0,
	}).Update()
	return res, err
}

/**
查询所有注册中心
*/
func FindAll() (*[]RegisterCenter, error) {
	res, err := config.Conn.Table("register_center").Fields("id,ip,register_time as registerTime,update_time updateTime,is_death as isDeath").Get()
	registerCenters := &[]RegisterCenter{}
	err1 := mapstructure.Decode(res, registerCenters)
	if err1 != nil {
		return registerCenters, err1
	}
	return registerCenters, err
}

func init() {
	tempAds := common.GetIp()
	registerCenter := &RegisterCenter{
		Id:           time.Now().Nanosecond(),
		Ip:           tempAds,
		RegisterTime: time.Now().Format(common.YYYY_MM_DD_HH_MM_SS),
		UpdateTime:   time.Now().Format(common.YYYY_MM_DD_HH_MM_SS),
		IsDeath:      0,
	}
	_, err1 := registerCenter.Register()
	if err1 != nil {
		common.GetLog().Infoln("ip:", tempAds, ":往注册中心注册失败:", err1)
	} else {
		common.GetLog().Infoln("ip:", tempAds, ":往注册中心注册成功")
		go HeartBeat()
	}
}

/**
向全网发送需要任务的消息
*/
func AllNetSendNeedTask() {
	rcs, _ := FindAll()
	for _, rc := range *rcs {
		if rc.Ip == common.GetIp() || rc.IsDeath == 1 {
			continue
		}
		count, err := net_client.SendNetInfo(rc.Ip, "NeedTasks")
		if err != nil || count == 0 {
			common.GetLog().Errorln("往目标IP:", rc.Ip, "发送失败", err)
		}
	}
}
