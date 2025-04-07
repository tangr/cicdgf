package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gproc"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gorilla/websocket"
)

var AgentCICD = agentCICD{}

type agentCICD struct{}

var (
	ctx                                   = context.Background()
	wsUrl                                 = g.Cfg().MustGet(ctx, "agent.WsUrl").String()
	ApiUrl                                = g.Cfg().MustGet(ctx, "agent.ApiUrl").String()
	syncInterval                          = g.Cfg().MustGet(ctx, "agent.SyncInterval").Int32()
	dataPathDir                           = g.Cfg().MustGet(ctx, "agent.DataPathDir").String()
	jobFlash                              = g.Cfg().MustGet(ctx, "agent.JobFlash").String()
	maxrunningjobs int                    = g.Cfg().MustGet(ctx, "agent.MaxRunningJobs").Int()
	runningJobs    map[int]*gproc.Process = make(map[int]*gproc.Process)
	envPrefix      string                 = g.Cfg().MustGet(ctx, "agent.EnvPrefix").String()
	agentInclude   string                 = g.Cfg().MustGet(ctx, "agent.Include").String()
	// wsAgentSend    chan model.WsAgentSend = make(chan model.WsAgentSend)
)

func main() {
	AgentCICD.AgentRun()
}

func (s *agentCICD) AgentRun() {
	if err := gfile.Mkdir(dataPathDir); err != nil {
		g.Log().Error(ctx, err)
		panic(err)
		// os.Exit(1)
	}

	interrupt := make(chan os.Signal, 1)
	// signal.Notify(interrupt, os.Interrupt, syscall.SIGUSR1)
	signal.Notify(interrupt, os.Interrupt)
	// signal.Notify(interrupt, syscall.SIGTERM)
	reload := make(chan os.Signal, 1)
	signal.Notify(reload, syscall.SIGUSR1)

	var recvJson = new(model.WsServerSend)

	client := ghttp.NewWebSocketClient()
	client.HandshakeTimeout = time.Second    // 设置超时时间
	client.Proxy = http.ProxyFromEnvironment // 设置代理

	// for i := 0; i < 10; i++ {
	for {
		// time.Sleep(time.Second)
		select {
		case <-interrupt:
			g.Log().Info(ctx, "interrupt2")
			os.Exit(1)
		case <-time.After(time.Second):
		}

		// conn, _, err := client.Dial("ws://127.0.0.1:8070/wsv1/wsci", nil)
		conn, _, err := client.Dial(wsUrl, nil)
		if err != nil {
			// panic(err)
			g.Log().Error(ctx, "dial:", err)
			continue
		}
		defer conn.Close()

		done := make(chan struct{})

		go func() {
			defer close(done)
			for {
				err := conn.ReadJSON(&recvJson)
				if err != nil {
					time.Sleep(time.Second)
					g.Log().Error(ctx, "read:", err)
					g.Log().Infof(ctx, "recv+v: %+v", recvJson)
					// continue
					break
					// return
				}
				// g.Log().Infof("recv+v: %+v", recvJson)

				newjobs, _ := json.Marshal(recvJson)
				g.Log().Infof(ctx, "recvjson: %s", string(newjobs))
				s.HandleRecvJson(recvJson)
			}
		}()

		ticker := time.NewTicker(time.Duration(1000000000 * syncInterval))
		defer ticker.Stop()

	L:
		for {
			// T:
			select {
			case <-done:
				break L
			case <-ticker.C:
				// g.Log().Info("*********************************")
				sendJson := s.SendJson()
				err := conn.WriteJSON(sendJson)
				if err != nil {
					time.Sleep(time.Second)
					g.Log().Error(ctx, "write:", err)
					g.Log().Infof(ctx, "send+v: %+v", sendJson)
					// continue
					break
					// return
				}
				// g.Log().Infof("send+v: %+v", sendJson)
				// g.Log().Infof("send#v: %#v", sendJson)
				newjobs, _ := json.Marshal(sendJson)
				g.Log().Infof(ctx, "sendjson: %s", string(newjobs))
				// g.Log().Info("###################################")
			case <-interrupt:
				g.Log().Info(ctx, "interrupt1")
				err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					g.Log().Warningf(ctx, "write close:", err)
					return
				}
				select {
				case <-done:
				case <-time.After(time.Second):
				}
				return
			case <-reload:
				g.Log().Info(ctx, "reload")
				s.GetAgentsList(true)
			}
		}
	}
}
