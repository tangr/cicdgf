package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"cicdgf/internal/model"

	"github.com/gofrs/flock"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
)

type WsAgentSend struct {
	Items      []WsAgentSendMap `json:"items"`
	TimeoutSec int              `json:"timeoutSec"`
}

type WsAgentSendMap = model.WsAgentSendMap

type WsServerSend = model.WsServerSend

type WsServerSendMap = model.WsServerSendMap

var AgentCICD = agentCICD{}

type agentCICD struct{}

var (
	ctx                                   = context.Background()
	apiUrl                                = g.Cfg().MustGet(ctx, "agent.ApiUrl").String()
	syncInterval                          = g.Cfg().MustGet(ctx, "agent.SyncInterval").Int32()
	dataPathDir                           = g.Cfg().MustGet(ctx, "agent.DataPathDir").String()
	jobFlash                              = g.Cfg().MustGet(ctx, "agent.JobFlash").String()
	maxrunningjobs int                    = g.Cfg().MustGet(ctx, "agent.MaxRunningJobs").Int()
	runningJobs    map[int]*gproc.Process = make(map[int]*gproc.Process)
	envPrefix      string                 = g.Cfg().MustGet(ctx, "agent.EnvPrefix").String()
	agents         AgentsList             = make(AgentsList, 0)
	// agentInclude   string                 = g.Cfg().MustGet(ctx, "agent.Include").String()

)

type AgentsMap struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type JobMeta struct {
	ID        int    `json:"jobid"`
	JobStatus string `json:"status"`
}

type AgentsList []AgentsMap

var (
	client *gclient.Client
)

func init() {
	client = g.Client()
	client.SetTimeout(100 * time.Second)
	header := g.MapStrStr{
		"Content-Type": "application/json",
	}
	client.SetHeaderMap(header)
}

func main() {
	AgentCICD.AgentRun()
}

func (s *agentCICD) GetAgentsList(isreload bool) AgentsList {
	if len(agents) != 0 && !isreload {
		return agents
	}

	var newagents AgentsList = make(AgentsList, 0)
	var agentsList AgentsList
	var agentsStr string = g.Cfg().MustGet(ctx, "agents").String()

	if agentsStr != "" {
		if err := gjson.DecodeTo(agentsStr, &agentsList); err != nil {
			g.Log().Errorf(ctx, "decode failed. %s", err)
		}
		newagents = append(newagents, agentsList...)
	}

	if len(newagents) > 0 {
		agents = newagents
		jobFlashStatus, err := json.Marshal(agents)
		if err != nil {
			g.Log().Error(ctx, err)
		}
		g.Log().Info(ctx, "jobFlashStatus: ", jobFlashStatus)
		jobFlashPath := dataPathDir + jobFlash
		g.Log().Debug(ctx, "jobFlashPath: ", jobFlashPath)
		s.WriteFile(jobFlashPath, string(jobFlashStatus))
	}
	return agents
}

func (s *agentCICD) PrepareAgentStatusUpdate() WsAgentSend {
	var agentsList AgentsList
	var agentSent = WsAgentSend{
		Items:      make([]WsAgentSendMap, 0),
		TimeoutSec: 30,
	}
	// var agentSentMap = WsAgentSendMap{}

	agentsList = s.GetAgentsList(false)
	for _, agent := range agentsList {
		agentSentMap := WsAgentSendMap{
			AgentId:   agent.ID,
			AgentName: agent.Name,
			// JobId, JobStatus, JobOutput 可以根据需要设置默认值或保持为零值
		}
		agentSent.Items = append(agentSent.Items, agentSentMap)

		// agentSentMap.AgentId = agent.ID
		// agentSentMap.AgentName = agent.Name
		// agentSent = append(agentSent, agentSentMap)
	}
	return agentSent
}

func (s *agentCICD) GetExecutable(scriptbody string) string {
	if len(scriptbody) < 3 {
		g.Log().Error(ctx, "scriptbody is empty")
		return ""
	}
	if scriptbody[0:2] == "#!" {
		return scriptbody[2:strings.Index(scriptbody, "\n")]
	} else {
		return "/usr/bin/env bash"
	}
}

func (s *agentCICD) WriteFile(path string, content string) error {
	g.Log().Debug(ctx, "Write file: ", path)
	if err := gfile.PutContents(path, content); err != nil {
		g.Log().Error(ctx, err)
		return err
	}
	return nil
}

func FileExists(name string) bool {
	g.Log().Debug(ctx, "FileExists: ", name)
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func (s *agentCICD) ReadFile(path string) string {
	g.Log().Debug(ctx, "Read file: ", path)
	if !FileExists(path) {
		g.Log().Debug(ctx, "file not exist: ", path)
		return ""
	}
	content := gfile.GetContents(path)
	return content
}

func (s *agentCICD) SetStatus(jobId int, jobStatus string) error {
	jobPathscriptJson := dataPathDir + strconv.Itoa(jobId) + ".json"
	oldJobStatus := s.GetStatus(jobId)
	if jobStatus == oldJobStatus {
		return nil
	}
	var jobMeta = JobMeta{}
	jobMeta.ID = jobId
	jobMeta.JobStatus = jobStatus
	jobJson, _ := json.Marshal(&jobMeta)
	fileLock := flock.New(jobPathscriptJson)
	err := fileLock.Lock()
	if err != nil {
		g.Log().Error(ctx, err)
	}
	if err := s.WriteFile(jobPathscriptJson, string(jobJson)); err != nil {
		g.Log().Error(ctx, err)
		return err
	}
	fileLock.Unlock()
	return nil
}

func (s *agentCICD) GetStatus(jobId int) string {
	jobPathscriptJson := dataPathDir + strconv.Itoa(jobId) + ".json"
	fileLock := flock.New(jobPathscriptJson)
	err := fileLock.RLock()
	if err != nil {
		g.Log().Error(ctx, err)
	}
	jobJson := s.ReadFile(jobPathscriptJson)
	fileLock.Unlock()
	if jobJson == "" {
		g.Log().Debugf(ctx, "fileName %s with content empty!", jobPathscriptJson)
		return ""
	}
	var jobMeta = JobMeta{}
	if err := json.Unmarshal([]byte(jobJson), &jobMeta); err != nil {
		g.Log().Error(ctx, err)
		g.Log().Debugf(ctx, "fileName %s with content: %s !", jobPathscriptJson, jobJson)
		return ""
	}
	return jobMeta.JobStatus
}

func (s *agentCICD) KillJob(jobId int) {
	if runningProcess, ok := runningJobs[jobId]; ok {
		g.Log().Warningf(ctx, "kill jobid: %d, pid: %d ", jobId, runningProcess.Cmd.Process.Pid)
		syscall.Kill(-runningProcess.Cmd.Process.Pid, syscall.SIGKILL)
		delete(runningJobs, jobId)

		if err := s.SetStatus(jobId, "failed"); err != nil {
			g.Log().Error(ctx, runningProcess.Cmd.Process.Pid, err)
		}
	}
}

func (s *agentCICD) RunCommand(jobId int, runCommand string, scriptEnvs []string) {
	// defer delete(runningJobs, jobId)
	g.Log().Debugf(ctx, "recvScriptEnvs: %+v", scriptEnvs)
	g.Log().Debugf(ctx, "recvScriptEnvs: %#v", scriptEnvs)
	newprocess := gproc.NewProcessCmd(runCommand, scriptEnvs)
	newprocess.Cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	newpid, err := newprocess.Start(ctx)
	if err != nil {
		g.Log().Error(ctx, newpid, err)
	}
	g.Log().Debugf(ctx, "Run newjob: %d pid: %d", jobId, newpid)
	if err := s.SetStatus(jobId, "running"); err != nil {
		g.Log().Error(ctx, newpid, err)
	}
	runningJobs[jobId] = newprocess
	if err = newprocess.Wait(); err != nil {
		g.Log().Warningf(ctx, "Command finished with error: %v", err)
	}
	g.Log().Debugf(ctx, "Finished Run newjob: %d pid: %d", jobId, newpid)

	if newprocess.ProcessState.Exited() {
		exitCode := newprocess.ProcessState.ExitCode()
		g.Log().Debugf(ctx, "Exit newjob: %d pid: %d exitcode: %d", jobId, newpid, exitCode)
		if exitCode == 0 {
			if err := s.SetStatus(jobId, "success"); err != nil {
				g.Log().Error(ctx, newpid, err)
			}
		} else {
			if err := s.SetStatus(jobId, "failed"); err != nil {
				g.Log().Error(ctx, newpid, err)
			}
		}
		delete(runningJobs, jobId)
	}
}

func (s *agentCICD) HandleJob2(jobv *WsServerSendMap) *WsAgentSendMap {
	var sendMap = &WsAgentSendMap{}
	jobId := jobv.JobId
	jobStatus := jobv.JobStatus
	sendMap.AgentId = jobv.AgentId
	sendMap.AgentName = jobv.AgentName
	sendMap.JobId = jobId
	g.Log().Debug(ctx, "HandleJob:jobStatus: ", jobStatus)
	if jobStatus == "success" || jobStatus == "failed" {
		sendMap.JobStatus = jobStatus
		jobPath := dataPathDir + strconv.Itoa(jobId)
		jobPathOutput := jobPath + ".output"
		output := s.ReadFile(jobPathOutput)
		sendMap.JobOutput = output
		return sendMap
	}
	if jobStatus == "running" {
		localJobStatus := s.GetStatus(jobId)
		if localJobStatus == "running" {
			if _, ok := runningJobs[jobId]; !ok {
				sendMap.JobStatus = "pending"
				return sendMap
			}
		}
		sendMap.JobStatus = localJobStatus
		jobPath := dataPathDir + strconv.Itoa(jobId)
		jobPathOutput := jobPath + ".output"
		output := s.ReadFile(jobPathOutput)
		sendMap.JobOutput = output
		return sendMap
	}
	if jobStatus == "aborted" {
		s.KillJob(jobId)
		sendMap.JobStatus = s.GetStatus(jobId)
		jobPath := dataPathDir + strconv.Itoa(jobId)
		jobPathOutput := jobPath + ".output"
		output := s.ReadFile(jobPathOutput)
		sendMap.JobOutput = output
		return sendMap
	}
	if jobStatus == "rerun" {
		oldJobStatus := s.GetStatus(jobId)
		if oldJobStatus == "success" || oldJobStatus == "failed" {
			sendMap.JobStatus = oldJobStatus
			return sendMap
		}
	}
	g.Log().Debug(ctx, 2222)
	g.Log().Debug(ctx, "HandleJob:jobId: %d", jobId)
	oldJobStatus := s.GetStatus(jobId)
	g.Log().Debug(ctx, 3333)
	g.Log().Debugf(ctx, "oldJobStatus: %s", oldJobStatus)
	if oldJobStatus == "success" || oldJobStatus == "failed" {
		sendMap.JobStatus = oldJobStatus
		jobPath := dataPathDir + strconv.Itoa(jobId)
		jobPathOutput := jobPath + ".output"
		output := s.ReadFile(jobPathOutput)
		sendMap.JobOutput = output
		return sendMap
	}
	if oldJobStatus == "" {
		if err := s.SetStatus(jobId, "pending"); err != nil {
			g.Log().Error(ctx, jobId, err)
		}
	}
	g.Log().Debugf(ctx, "oldJobStatus2: %s", oldJobStatus)
	g.Log().Debug(ctx, "jobStatus: ", jobStatus)
	jobPath := dataPathDir + strconv.Itoa(jobId)
	jobPathOutput := jobPath + ".output"
	if jobv.Body != "" {
		if _, ok := runningJobs[jobId]; !ok {
			scriptBody := jobv.Body + "\n"
			scriptBody = strings.Replace(scriptBody, "\r\n", "\n", -1)
			jobPathscriptBody := jobPath + ".scriptbody"
			s.WriteFile(jobPathscriptBody, scriptBody)
			scriptArgs := jobv.Args + "\n"
			scriptArgs = strings.Replace(scriptArgs, "\r\n", "\n", -1)
			jobPathscriptArgs := jobPath + ".scriptargs"
			s.WriteFile(jobPathscriptArgs, scriptArgs)
			var scriptEnvs []string
			envAgentName := strings.Split(jobv.AgentName, ":")[0]
			scriptEnvs = append(scriptEnvs, envPrefix+"AGENTNAME"+"="+envAgentName)
			for k, v := range jobv.Envs {
				scriptEnvs = append(scriptEnvs, envPrefix+k+"="+v)
			}
			execommand := s.GetExecutable(scriptBody)
			if execommand != "" {
				runcommand := execommand + " " + jobPathscriptBody + " " + jobPathscriptArgs + " >>" + jobPathOutput + " 2>&1"
				g.Log().Debugf(ctx, "Run jobId: %d with Command: %s and scriptEnvs: %s", jobId, runcommand, scriptEnvs)
				go s.RunCommand(jobId, runcommand, scriptEnvs)
			}
		}
	}
	sendMap.JobStatus = s.GetStatus(jobId)
	output := s.ReadFile(jobPathOutput)
	sendMap.JobOutput = output
	return sendMap
}

func (s *agentCICD) HandleJob(ctx context.Context, jobv *WsServerSendMap) {
	var sendMap = &WsAgentSendMap{}
	jobId := jobv.JobId
	jobStatus := jobv.JobStatus
	sendMap.AgentId = jobv.AgentId
	sendMap.AgentName = jobv.AgentName
	sendMap.JobId = jobId

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			g.Log().Debug(ctx, "File reading stopped")
			return
		case <-ticker.C:
			oldJobStatus := s.GetStatus(jobId)
			if oldJobStatus == "" {
				if err := s.SetStatus(jobId, "pending"); err != nil {
					g.Log().Error(ctx, jobId, err)
				}
				// jobPath := dataPathDir + strconv.Itoa(jobId)
				// jobPathOutput := jobPath + ".output"
				var taskRes model.TaskRes

				url := apiUrl + "/log/" + strconv.Itoa(jobv.TaskId)
				response, err := client.Get(ctx, url)
				if err != nil {
					g.Log().Error(ctx, jobId, err)
				}
				res := response.ReadAll()
				g.Log().Debug(ctx, res)
				err = json.Unmarshal(res, &taskRes)
				if err != nil {
					g.Log().Errorf(ctx, "解析服务器响应失败: %v", err)
					continue
				}
				jobId = taskRes.Data.JobId

				var jobRes model.JobRes

				url = apiUrl + "/job/" + strconv.Itoa(jobId)
				response, err = client.Get(ctx, url)
				if err != nil {
					g.Log().Error(ctx, jobId, err)
				}
				res = response.ReadAll()
				g.Log().Debug(ctx, res)
				err = json.Unmarshal(res, &jobRes)
				if err != nil {
					g.Log().Errorf(ctx, "解析服务器响应失败: %v", err)
					continue
				}
				jobPath := dataPathDir + strconv.Itoa(jobId)
				jobPathOutput := jobPath + ".output"
				script_body := jobRes.Data.Script.Body

				if _, ok := runningJobs[jobId]; !ok {
					scriptBody := script_body + "\n"
					scriptBody = strings.Replace(scriptBody, "\r\n", "\n", -1)
					jobPathscriptBody := jobPath + ".scriptbody"
					s.WriteFile(jobPathscriptBody, scriptBody)
					scriptArgs := jobv.Args + "\n"
					scriptArgs = strings.Replace(scriptArgs, "\r\n", "\n", -1)
					jobPathscriptArgs := jobPath + ".scriptargs"
					s.WriteFile(jobPathscriptArgs, scriptArgs)
					var scriptEnvs []string
					envAgentName := strings.Split(jobv.AgentName, ":")[0]
					scriptEnvs = append(scriptEnvs, envPrefix+"AGENTNAME"+"="+envAgentName)
					for k, v := range jobv.Envs {
						scriptEnvs = append(scriptEnvs, envPrefix+k+"="+v)
					}
					execommand := s.GetExecutable(scriptBody)
					if execommand != "" {
						runcommand := execommand + " " + jobPathscriptBody + " " + jobPathscriptArgs + " >>" + jobPathOutput + " 2>&1"
						g.Log().Debugf(ctx, "Run jobId: %d with Command: %s and scriptEnvs: %s", jobId, runcommand, scriptEnvs)
						go s.RunCommand(jobId, runcommand, scriptEnvs)
					}
				}
			}

			jobPath := dataPathDir + strconv.Itoa(jobId)
			jobPathOutput := jobPath + ".output"
			output := s.ReadFile(jobPathOutput)
			g.Log().Debug(ctx, "File content: %s\n", string(output))
			sendMap.JobOutput = output
			jobStatus := s.GetStatus(jobId)
			sendMap.JobStatus = jobStatus
			url := apiUrl + "/log/" + strconv.Itoa(jobv.TaskId)
			client.Put(ctx, url, sendMap)
		}
	}

	// var sendMap = &WsAgentSendMap{}
	// jobId := jobv.JobId
	// jobStatus := jobv.JobStatus
	// sendMap.AgentId = jobv.AgentId
	// sendMap.AgentName = jobv.AgentName
	// sendMap.JobId = jobId
	g.Log().Debug(ctx, "HandleJob:jobStatus: ", jobStatus)
	if jobStatus == "success" || jobStatus == "failed" {
		sendMap.JobStatus = jobStatus
		jobPath := dataPathDir + strconv.Itoa(jobId)
		jobPathOutput := jobPath + ".output"
		output := s.ReadFile(jobPathOutput)
		sendMap.JobOutput = output
		// return sendMap
	}
	if jobStatus == "running" {
		localJobStatus := s.GetStatus(jobId)
		if localJobStatus == "running" {
			if _, ok := runningJobs[jobId]; !ok {
				sendMap.JobStatus = "pending"
				// return sendMap
			}
		}
		sendMap.JobStatus = localJobStatus
		jobPath := dataPathDir + strconv.Itoa(jobId)
		jobPathOutput := jobPath + ".output"
		output := s.ReadFile(jobPathOutput)
		sendMap.JobOutput = output
		// return sendMap
	}
	if jobStatus == "aborted" {
		s.KillJob(jobId)
		sendMap.JobStatus = s.GetStatus(jobId)
		jobPath := dataPathDir + strconv.Itoa(jobId)
		jobPathOutput := jobPath + ".output"
		output := s.ReadFile(jobPathOutput)
		sendMap.JobOutput = output
		// return sendMap
	}
	if jobStatus == "rerun" {
		oldJobStatus := s.GetStatus(jobId)
		if oldJobStatus == "success" || oldJobStatus == "failed" {
			sendMap.JobStatus = oldJobStatus
			// return sendMap
		}
	}
	g.Log().Debug(ctx, 2222)
	g.Log().Debug(ctx, "HandleJob:jobId: %d", jobId)
	oldJobStatus := s.GetStatus(jobId)
	g.Log().Debug(ctx, 3333)
	g.Log().Debugf(ctx, "oldJobStatus: %s", oldJobStatus)
	if oldJobStatus == "success" || oldJobStatus == "failed" {
		sendMap.JobStatus = oldJobStatus
		jobPath := dataPathDir + strconv.Itoa(jobId)
		jobPathOutput := jobPath + ".output"
		output := s.ReadFile(jobPathOutput)
		sendMap.JobOutput = output
		// return sendMap
	}
	if oldJobStatus == "" {
		if err := s.SetStatus(jobId, "pending"); err != nil {
			g.Log().Error(ctx, jobId, err)
		}
	}
	g.Log().Debugf(ctx, "oldJobStatus2: %s", oldJobStatus)
	g.Log().Debug(ctx, "jobStatus: ", jobStatus)
	jobPath := dataPathDir + strconv.Itoa(jobId)
	jobPathOutput := jobPath + ".output"
	if jobv.Body != "" {
		if _, ok := runningJobs[jobId]; !ok {
			scriptBody := jobv.Body + "\n"
			scriptBody = strings.Replace(scriptBody, "\r\n", "\n", -1)
			jobPathscriptBody := jobPath + ".scriptbody"
			s.WriteFile(jobPathscriptBody, scriptBody)
			scriptArgs := jobv.Args + "\n"
			scriptArgs = strings.Replace(scriptArgs, "\r\n", "\n", -1)
			jobPathscriptArgs := jobPath + ".scriptargs"
			s.WriteFile(jobPathscriptArgs, scriptArgs)
			var scriptEnvs []string
			envAgentName := strings.Split(jobv.AgentName, ":")[0]
			scriptEnvs = append(scriptEnvs, envPrefix+"AGENTNAME"+"="+envAgentName)
			for k, v := range jobv.Envs {
				scriptEnvs = append(scriptEnvs, envPrefix+k+"="+v)
			}
			execommand := s.GetExecutable(scriptBody)
			if execommand != "" {
				runcommand := execommand + " " + jobPathscriptBody + " " + jobPathscriptArgs + " >>" + jobPathOutput + " 2>&1"
				g.Log().Debugf(ctx, "Run jobId: %d with Command: %s and scriptEnvs: %s", jobId, runcommand, scriptEnvs)
				go s.RunCommand(jobId, runcommand, scriptEnvs)
			}
		}
	}
	sendMap.JobStatus = s.GetStatus(jobId)
	output := s.ReadFile(jobPathOutput)
	sendMap.JobOutput = output
	// return sendMap
}

func (s *agentCICD) HandleRecvJson(recvJson *WsServerSend) {
	// var sendJson WsAgentSend
	// var sendJson = WsAgentSend{
	// 	Items:      make([]WsAgentSendMap, 0),
	// 	TimeoutSec: 30, // 设置默认超时时间，可根据需要调整
	// }

	recvData := *recvJson
	if len(recvData.Data) < 1 {
		return
	}

	for _, jobv := range recvData.Data {
		// if jobv.ErrMsg != "" {
		// 	g.Log().Errorf(ctx, "jobId: %d errmsg: %s", jobv.JobId, jobv.ErrMsg)
		// 	continue
		// }
		if jobv.JobId == 0 || jobv.JobStatus == "" {
			continue
		}
		g.Log().Debugf(ctx, "len runningJobs: %d %d", len(runningJobs), maxrunningjobs)
		if len(runningJobs) >= maxrunningjobs {
			jobId := jobv.JobId
			if _, ok := runningJobs[jobId]; !ok {
				continue
			}
		}
		g.Log().Debugf(ctx, "recvjson: %#v", jobv)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go s.HandleJob(ctx, &jobv)
		// var sendMap = s.HandleJob(&jobv)
		// g.Log().Debugf(ctx, "sendjson: %#v", sendMap)
		// sendJson.Items = append(sendJson.Items, *sendMap)
	}

	// jobv := recvData.Data

	// if jobv.ErrMsg != "" {
	// 	g.Log().Errorf(ctx, "jobId: %d errmsg: %s", jobv.JobId, jobv.ErrMsg)
	// }

	// g.Log().Debug(ctx, "HandleRecvJson jobv: %+v", jobv)
	// if jobv.JobId == 0 || jobv.JobStatus == "" {
	// }
	// g.Log().Debugf(ctx, "len runningJobs: %d %d", len(runningJobs), maxrunningjobs)
	// if len(runningJobs) >= maxrunningjobs {
	// 	jobId := jobv.JobId
	// 	if _, ok := runningJobs[jobId]; !ok {
	// 		g.Log().Error(ctx, ok)
	// 	}
	// }
	// g.Log().Debugf(ctx, "recvjson: %#v", jobv)
	// var sendMap = s.HandleJob(&jobv)
	// g.Log().Debugf(ctx, "sendjson: %#v", sendMap)
	// sendJson.Items = append(sendJson.Items, *sendMap)

}

func (s *agentCICD) AgentRun() {
	if err := gfile.Mkdir(dataPathDir); err != nil {
		g.Log().Error(ctx, err)
		panic(err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	reload := make(chan os.Signal, 1)
	signal.Notify(reload, syscall.SIGUSR1)

	// 创建HTTP客户端
	// client := g.Client()
	// client.SetTimeout(100 * time.Second)
	// header := g.MapStrStr{
	// 	"Content-Type": "application/json",
	// }
	// client.SetHeaderMap(header)

	ticker := time.NewTicker(time.Duration(syncInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-interrupt:
			g.Log().Info(ctx, "程序被中断，正在退出...")
			return
		case <-reload:
			g.Log().Info(ctx, "正在重新加载配置...")
			s.GetAgentsList(true)
		case <-ticker.C:
			// 准备要发送的Agent状态数据
			agentStatus := s.PrepareAgentStatusUpdate()

			// 发送Agent状态到服务器
			g.Log().Infof(ctx, "发送Agent状态更新：%+v", agentStatus)
			response, err := client.Post(ctx, apiUrl+"/notifys/v1", agentStatus)
			if err != nil {
				g.Log().Errorf(ctx, "发送状态更新失败: %v", err)
				continue
			}

			if response.StatusCode != 200 {
				continue
			}

			// 解析服务器响应
			var serverTasks WsServerSend

			g.Log().Debug(ctx, apiUrl+"/notifys/v1")
			g.Log().Debug(ctx, response.StatusCode)

			res := response.ReadAll()
			g.Log().Debug(ctx, res)
			err = json.Unmarshal(res, &serverTasks)
			if err != nil {
				g.Log().Errorf(ctx, "解析服务器响应失败: %v", err)
				continue
			}

			// 处理服务器下发的任务
			g.Log().Infof(ctx, "接收到服务器任务：%v", serverTasks)
			s.HandleRecvJson(&serverTasks)

			// // 处理服务器下发的任务
			// if len(serverTasks) > 0 {
			// 	g.Log().Infof(ctx, "接收到服务器任务：%v", serverTasks)
			// 	result := s.HandleRecvJson(&serverTasks)

			// 	// 上报任务处理结果
			// 	if len(result.Items) > 0 {
			// 		_, err := client.Post(ctx, apiUrl+"/agent/job/result", result)
			// 		if err != nil {
			// 			g.Log().Errorf(ctx, "上报任务结果失败: %v", err)
			// 		}
			// 	}
			// }

			// // 处理服务器下发的任务
			// g.Log().Infof(ctx, "接收到服务器任务：%v", serverTasks)
			// result := s.HandleRecvJson(&serverTasks)

			// // 上报任务处理结果
			// if len(result.Items) > 0 {
			// 	_, err := client.Post(ctx, apiUrl+"/agent/job/result", result)
			// 	if err != nil {
			// 		g.Log().Errorf(ctx, "上报任务结果失败: %v", err)
			// 	}
			// }
		}
	}
}
