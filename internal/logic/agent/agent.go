package agent

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
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type WsAgentSend struct {
	Items      []WsAgentSendMap `json:"items"`
	TimeoutSec int              `json:"timeoutSec"`
}

type WsAgentSendMap = model.WsAgentSendMap

type WsAgentSendLogMap = model.WsAgentSendLogMap

type WsServerSend = model.WsServerSend

type WsServerSendMap = model.WsServerSendMap

type TaskInfoMap = model.TaskInfoMap

type Script = model.Script

var AgentCICD = agentCICD{}

type agentCICD struct{}

var (
	ctx                                   = context.Background()
	apiUrl                                = g.Cfg().MustGet(ctx, "agent.ApiUrl").String()
	syncInterval                          = g.Cfg().MustGet(ctx, "agent.SyncInterval").Int32()
	dataPathDir                           = g.Cfg().MustGet(ctx, "agent.DataPathDir").String()
	jobFlash                              = g.Cfg().MustGet(ctx, "agent.JobFlash").String()
	maxRunningJobs int                    = g.Cfg().MustGet(ctx, "agent.MaxRunningJobs").Int()
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

// func main() {
// 	AgentCICD.AgentRun()
// }

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

func (s *agentCICD) HandleJob(ctx context.Context, jobv *WsServerSendMap) {
	var sendMap = &WsAgentSendLogMap{}
	taskId := jobv.TaskId

	taskInfo := s.GetTaskInfoById(taskId)

	sendMap.JobType = taskInfo.JobType
	sendMap.Ipaddr = "127.0.0.1"

	// jobStatus := jobv.JobStatus
	sendMap.AgentId = jobv.AgentId
	sendMap.JobId = taskInfo.JobId
	sendMap.PipelineId = taskInfo.PipelineId

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			g.Log().Debug(ctx, "File reading stopped")
			return
		case <-ticker.C:
			oldStatus := s.GetStatus(taskId)
			if oldStatus == "" {
				if err := s.SetStatus(taskId, "pending"); err != nil {
					g.Log().Error(ctx, taskId, err)
				}

				g.Log().Debug(ctx, "HandleJob GetScriptByTask")

				var script Script = s.GetScriptByTask(taskId)
				script_body := script.Body
				script_envs := script.Envs
				script_args := script.Args

				g.Log().Debugf(ctx, "HandleJob GetScriptByTask script_body: %s", script_body)
				g.Log().Debugf(ctx, "HandleJob GetScriptByTask script_args: %s", script_args)
				g.Log().Debugf(ctx, "HandleJob GetScriptByTask script_envs: %s", script_envs)

				jobPath := dataPathDir + strconv.Itoa(taskId)
				jobPathOutput := jobPath + ".output"

				if _, ok := runningJobs[taskId]; !ok {
					scriptBody := script_body + "\n"
					scriptBody = strings.Replace(scriptBody, "\r\n", "\n", -1)
					jobPathscriptBody := jobPath + ".scriptbody"
					s.WriteFile(jobPathscriptBody, scriptBody)
					scriptArgs := script_args + "\n"
					scriptArgs = strings.Replace(scriptArgs, "\r\n", "\n", -1)
					jobPathscriptArgs := jobPath + ".scriptargs"
					s.WriteFile(jobPathscriptArgs, scriptArgs)
					var scriptEnvs []string
					envAgentName := strings.Split(jobv.AgentName, ":")[0]
					scriptEnvs = append(scriptEnvs, envPrefix+"AGENTNAME"+"="+envAgentName)
					for k, v := range script_envs {
						scriptEnvs = append(scriptEnvs, envPrefix+k+"="+v)
					}
					execommand := s.GetExecutable(scriptBody)
					if execommand != "" {
						runcommand := execommand + " " + jobPathscriptBody + " " + jobPathscriptArgs + " >>" + jobPathOutput + " 2>&1"
						g.Log().Debugf(ctx, "Run taskId: %d with Command: %s and scriptEnvs: %s", taskId, runcommand, scriptEnvs)
						go s.RunCommand(taskId, runcommand, scriptEnvs)
					}
				}
			}

			jobPath := dataPathDir + strconv.Itoa(taskId)
			jobPathOutput := jobPath + ".output"
			output := s.ReadFile(jobPathOutput)
			g.Log().Debug(ctx, "File content: %s\n", string(output))
			sendMap.Output = output
			taskStatus := s.GetStatus(taskId)
			sendMap.TaskStatus = taskStatus
			currentTime := gtime.Timestamp()
			sendMap.UpdatedAt = currentTime

			g.Log().Debugf(ctx, "currentTime: %d-%d", taskId, currentTime)
			g.Log().Debugf(ctx, "sendMap: %s", gconv.String(sendMap))
			url := apiUrl + "/log/" + strconv.Itoa(taskId)
			response, err := client.Put(ctx, url, sendMap)
			if err != nil {
				g.Log().Errorf(ctx, "发送状态更新失败: %v", err)
			}
			res := response.ReadAll()
			g.Log().Debugf(ctx, "Receive Put response: %s", gconv.String(res))
			g.Log().Debugf(ctx, "Receive Put StatusCode: %s", gconv.String(response.StatusCode))

			g.Log().Debugf(ctx, "Send Put req: %s", gconv.String(taskStatus))

			if taskStatus == "success" || taskStatus == "failed" {
				g.Log().Debugf(ctx, "Send Put req2: %s", gconv.String(taskStatus))
				return
			}
		}
	}

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
		taskId := jobv.TaskId
		if taskId == 0 {
			continue
		}
		g.Log().Debugf(ctx, "len runningJobs: %d %d", len(runningJobs), maxRunningJobs)
		if len(runningJobs) >= maxRunningJobs {
			if _, ok := runningJobs[taskId]; !ok {
				continue
			}
		}
		g.Log().Debugf(ctx, "recvjson: %#v", jobv)
		// ctx, cancel := context.WithCancel(context.Background())
		// defer cancel()
		go s.HandleJob(ctx, &jobv)
		// var sendMap = s.HandleJob(&jobv)
		// g.Log().Debugf(ctx, "sendjson: %#v", sendMap)
		// sendJson.Items = append(sendJson.Items, *sendMap)
	}

}

func (s *agentCICD) AgentRun() {
	g.Log().Debug(ctx, dataPathDir)
	if err := gfile.Mkdir(dataPathDir); err != nil {
		g.Log().Error(ctx, err)
		panic(err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	reload := make(chan os.Signal, 1)
	signal.Notify(reload, syscall.SIGUSR1)

	// 移除 ticker，使用 done channel 来控制循环退出
	done := make(chan bool)

	// 启动主循环协程
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				// 准备要发送的Agent状态数据
				agentStatus := s.PrepareAgentStatusUpdate()

				// 发送Agent状态到服务器
				g.Log().Debugf(ctx, "******************************")
				g.Log().Debugf(ctx, "Send agentStatus: %s", gconv.String(agentStatus))
				response, err := client.Post(ctx, apiUrl+"/notifys/v1", agentStatus)
				if err != nil {
					g.Log().Errorf(ctx, "发送状态更新失败: %v", err)
					// 发生错误时等待一段时间再重试，避免频繁请求
					time.Sleep(time.Duration(syncInterval) * time.Second)
					continue
				}

				if response.StatusCode != 200 {
					g.Log().Debugf(ctx, "Get StatusCode: %d, %s", response.StatusCode, response.ReadAll())
					// 状态码不是200时等待一段时间再重试
					time.Sleep(time.Duration(syncInterval) * time.Second)
					continue
				}

				// 解析服务器响应
				var serverTasks WsServerSend

				res := response.ReadAll()

				g.Log().Debugf(ctx, "Receive res: %s", gconv.String(res))

				err = json.Unmarshal(res, &serverTasks)
				if err != nil {
					g.Log().Errorf(ctx, "解析服务器响应失败: %v", err)
					continue
				}

				// 处理服务器下发的任务
				g.Log().Debugf(ctx, "Receive serverTasks: %s", gconv.String(serverTasks))

				s.HandleRecvJson(&serverTasks)
				time.Sleep(time.Duration(1) * time.Second)
			}
		}
	}()

	// 主线程监听信号
	for {
		select {
		case <-interrupt:
			g.Log().Info(ctx, "程序被中断，正在退出...")
			done <- true
			return
		case <-reload:
			g.Log().Info(ctx, "正在重新加载配置...")
			s.GetAgentsList(true)
		}
	}
}
