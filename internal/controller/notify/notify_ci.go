package notify

import (
	"cicdgf/internal/dao"
	"context"
	"strconv"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

func SyncNewCIJob(ctx context.Context) {
	type NewJobBuild struct {
		Id      uint `json:"task_id"`
		AgentId uint `json:"agent_id"`
	}
	var newJobs = new([]NewJobBuild)

	ticker := time.NewTicker(15 * time.Second)

	go func() {
		for {
			<-ticker.C

			// now := time.Now()

			g.Log().Debug(ctx, "SyncNewCIJob CicdJob")

			err := dao.CicdLog.Ctx(ctx).
				Fields("id,agent_id").
				Where("job_type", "BUILD").
				WhereIn("task_status", g.Slice{"pending"}).
				Scan(newJobs)

			if err != nil {
				g.Log().Debug(ctx, "Failed to get BUILD jobs:", err)
			}

			expireSecs := int64(10 * 60)

			for _, newJob := range *newJobs {
				agentId := newJob.AgentId
				taskId := newJob.Id

				ciAgentKey := "ciagent:" + strconv.FormatUint(uint64(agentId), 10)
				count, err := g.Redis().Exists(ctx, ciAgentKey)
				if err != nil {
					g.Log().Fatal(ctx, err)

				}
				AddNotification(agentId, taskId)

				if count == 0 {
					err = g.Redis().SetEX(ctx, ciAgentKey, taskId, expireSecs)
					if err != nil {
						g.Log().Fatal(ctx, err)
					}
				}
			}
		}
	}()
}
