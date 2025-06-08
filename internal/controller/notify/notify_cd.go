package notify

import (
	"cicdgf/internal/dao"
	"context"
	"strconv"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

func SyncNewCDJob(ctx context.Context) {
	type NewJobBuild struct {
		Id      uint `json:"id"`
		AgentId uint `json:"agent_id"`
	}
	// var newJobs = new([]NewJobBuild)
	var newJobs []NewJobBuild

	ticker := time.NewTicker(15 * time.Second)

	go func() {
		for {
			<-ticker.C

			// now := time.Now()

			g.Log().Debug(ctx, "SyncNewCDJob CicdJob")

			newJobs = newJobs[:0]
			err := dao.CicdLog.Ctx(ctx).
				Fields("id,agent_id").
				Where("job_type", "DEPLOY").
				WhereIn("task_status", g.Slice{"pending"}).
				Scan(&newJobs)

			if err != nil {
				g.Log().Debug(ctx, "Failed to get DEPLOY jobs:", err)
			}

			g.Log().Debug(ctx, "SyncNewCDJob newJobs:", newJobs)

			expireSecs := int64(10 * 60)

			for _, newJob := range newJobs {
				agentId := newJob.AgentId
				taskId := newJob.Id

				cdAgentKey := "cdagent:" + strconv.FormatUint(uint64(agentId), 10)
				count, err := g.Redis().Exists(ctx, cdAgentKey)
				if err != nil {
					g.Log().Fatal(ctx, err)

				}
				AddNotification(agentId, taskId)

				if count == 0 {
					err = g.Redis().SetEX(ctx, cdAgentKey, taskId, expireSecs)
					if err != nil {
						g.Log().Fatal(ctx, err)
					}
				}
			}
		}
	}()
}
