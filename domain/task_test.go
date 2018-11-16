package domain

import (
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"
)

var wg sync.WaitGroup

func BenchmarkTask_AddTask(b *testing.B) {
	start := time.Now().Unix()
	//maxProcs := runtime.NumCPU()   //获取cpu个数
	//runtime.GOMAXPROCS(maxProcs)  //限制同时运行的goroutines数量
	for m := 0; m < 10; m++ {
		wg.Add(1)
		go insert()
	}
	wg.Wait()
	end := time.Now().Unix()
	log.Println("执行完毕,耗时：", (end - start), " 秒")
}

func insert() {
	for i := 0; i < 1000; i++ {
		task := &Task{}
		task.TaskNo = "100"
		task.GroupNo = "100"
		task.Name = "test"
		task.TaskSlice = 1
		task.CreateTime = "2018-05-01 12:12:12"
		task.UpdateTime = "2018-05-01 12:12:12"
		task.Cron = "0,0,0,1,0,?"
		task.Description = "测试任务"
		task.IsActivity = 1
		task.Mail = "sunpengwei1992@aliyun.com"
		task.Param = ""

		task.AddTask()
	}
	wg.Done()
}

func TestTask_AddTask(t *testing.T) {

	task := &Task{}
	task.TaskNo = "100"
	task.GroupNo = "100"
	task.Name = "test"
	task.TaskSlice = 1
	task.CreateTime = "2018-05-01 12:12:12"
	task.UpdateTime = "2018-05-01 12:12:12"
	task.Cron = "0,0,0,1,0,?"
	task.Description = "测试任务"
	task.IsActivity = 1
	task.Mail = "sunpengwei1992@aliyun.com"
	task.Param = ""
	s := rand.Intn(127)

	log.Println("分片:", s)
}

/*func TestTask_FindAll(t *testing.T) {
	task.FindAll()
}*/

func TestTask_FindByTaskSlice(t *testing.T) {
	FindByTaskSlice([]interface{}{1, 2})
}
