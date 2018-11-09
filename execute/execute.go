package execute

import (
	"task_dispatcher/domain"
	"io/ioutil"
	"time"
	"task_dispatcher/common"
	"net/http"
	"bytes"
	"task_dispatcher/enum"
	"strings"
	"strconv"
)
const (
	MAX_CONCURRENCY_NUM = 100
)
type Execute struct {

}
func (scheduler *Execute)ExecuteTask(tk *domain.Task)  {
	// 并发数最多100
	if tk.ConcurrencyNum > MAX_CONCURRENCY_NUM{
		tk.ConcurrencyNum = MAX_CONCURRENCY_NUM
	}
	for i:=0;i<tk.ConcurrencyNum ;i++  {
		go func() {
			// 调用 http 请求
			status,errorMsg:= invoker(tk.Url,tk.Param,i,tk.IsRegion,tk.RequestMethod)

			if tk.IsRecordLog == 0{ // 如果等于0 表示不记录日志
				return
			}
			// 根据返回结果新增执行记录
			if status == 1 {
				CreateTaskExecuteRecord(tk, 1, errorMsg, nil).AddTaskExecuteRecord()
			} else {
				CreateTaskExecuteRecord(tk, 0, errorMsg, nil).AddTaskExecuteRecord()
			}

		}()
	}
}

func CreateTaskExecuteRecord(tk *domain.Task,status int,errorMsg string,err error)(ter *domain.TaskExecuteRecord) {

	if err != nil{
		errorMsg = errorMsg+err.Error()
	}
	if errorMsg != ""{
		 common.GetLog().Errorln(tk.Url,"errorMsg",errorMsg)
	}
	if len(errorMsg) > 1000 {
		errorMsg = string([]rune(errorMsg)[0:999])
	}
	return &domain.TaskExecuteRecord{
		TaskNo:tk.TaskNo,
		GroupNo:tk.GroupNo,
		TaskName:tk.Name,
		Status:status,
		ErrorMsg:errorMsg,
		Result:"",
		CreateTime:time.Now().Format(common.YYYY_MM_DD_HH_MM_SS),
		IsDelete:0,
	}
}
/**
	0 代表http接口执行失败
    1 代表http接口执行成功
 */
func invoker(url,param string,region,isRegion,requestMethod int) (int,string){
	client := &http.Client{}

	// 如果分区
	if isRegion == enum.TASK_IS_REGION_YES{
		//对 param 做特殊处理
		paramArray := strings.Split(param,"]")
		r:= strconv.FormatInt(int64(region),10)
		param = paramArray[0]+","+ r +"]"
	}
	var req *http.Request
	if requestMethod == 1{
		req, _ = http.NewRequest("POST", url, bytes.NewBuffer([]byte(param)))
	}else{
		req, _ = http.NewRequest("GET", url, nil)
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("token", "123456")// jsf 的token
	res, err := client.Do(req)
	if err != nil {
		common.GetLog().Errorln("请求错误:",url)
		return 0,err.Error()
	}
	defer  res.Body.Close()
	bs,_ := ioutil.ReadAll(res.Body)

	/*result := &common.Result{}
	json.Unmarshal(bs,result)
	return result.Code,result.Message// 返回执行结果*/

	return 1,string(bs)
}




