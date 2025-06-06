package agent

import (
	"cicdgf/internal/model"
	"encoding/json"
	"strconv"

	"github.com/gogf/gf/v2/frame/g"
)

type JobInfoMap model.JobInfoMap

func (j *JobInfoMap) GetScript() Script {
	var script Script
	err := json.Unmarshal([]byte(j.Script), &script)
	if err != nil {
		g.Log().Debugf(ctx, "GetScript %s", err)
	}
	return script
}

func (s *agentCICD) GetScriptByTask(taskid int) Script {
	taskInfo := s.GetTaskInfoById(taskid)
	jobId := taskInfo.JobId

	var jobGetRes model.JobGetRes

	url := apiUrl + "/job/" + strconv.Itoa(jobId)
	response, err := client.Get(ctx, url)
	if err != nil {
		g.Log().Error(ctx, jobId, err)
	}

	res := response.ReadAll()
	g.Log().Debugf(ctx, "GetScriptByTask res %s", res)

	err = json.Unmarshal(res, &jobGetRes)
	if err != nil {
		g.Log().Errorf(ctx, "GetScriptByTask failed: %v", err)
	}

	// 将 model.JobInfoMap 转换为本地的 JobInfoMap
	localJobInfo := JobInfoMap(jobGetRes.Data)
	script := localJobInfo.GetScript()
	if err != nil {
		g.Log().Errorf(ctx, "Parse script failed: %v", err)
	}

	return script
}

func (s *agentCICD) GetTaskInfoById(taskid int) TaskInfoMap {
	var taskGetRes model.TaskGetRes

	url := apiUrl + "/log/" + strconv.Itoa(taskid)
	response, err := client.Get(ctx, url)
	if err != nil {
		g.Log().Error(ctx, taskid, err)
	}
	res := response.ReadAll()
	g.Log().Debug(ctx, res)
	err = json.Unmarshal(res, &taskGetRes)
	if err != nil {
		g.Log().Errorf(ctx, "GetTaskInfoById failed: %v", err)
	}
	taskInfo := taskGetRes.Data
	return taskInfo
}
