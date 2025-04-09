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
		ID      uint `json:"jobid"`
		AgentId uint `json:"agent_id"`
	}
	var newJobs = new([]NewJobBuild)

	ticker := time.NewTicker(5 * time.Second)

	go func() {
		for {
			<-ticker.C

			// now := time.Now()

			g.Log().Debug(ctx, "SyncNewCIJob CicdJob")
			err := dao.CicdJob.Ctx(ctx).
				Fields("id,agent_id").
				Where("job_type", "BUILD").
				WhereIn("job_status", g.Slice{"pending"}).
				Scan(newJobs)

			if err != nil {
				g.Log().Debug(ctx, "Failed to get BUILD jobs:", err)
			}

			expireSecs := int64(10 * 60)

			for _, newJob := range *newJobs {
				agentId := newJob.AgentId
				jobId := newJob.ID

				ciAgentKey := "ciagent:" + strconv.FormatUint(uint64(agentId), 10)
				count, err := g.Redis().Exists(ctx, ciAgentKey)
				if err != nil {
					g.Log().Fatal(ctx, err)

				}
				AddNotification(agentId, jobId)

				if count == 0 {
					err = g.Redis().SetEX(ctx, ciAgentKey, jobId, expireSecs)
					if err != nil {
						g.Log().Fatal(ctx, err)
					}

				}

			}

		}
	}()
}
