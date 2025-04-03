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
		ID      int `json:"jobid"`
		AgentId int `json:"agent_id"`
	}
	var newJobs = new([]NewJobBuild)

	ticker := time.NewTicker(5 * time.Second)

	go func() {
		for {
			<-ticker.C

			err := dao.CicdJob.Ctx(ctx).
				Fields("id,agent_id").
				Where("job_type", "BUILD").
				WhereIn("job_status", g.Slice{"pending"}).
				Scan(newJobs)

			if err != nil {
				g.Log().Debug(ctx, "Failed to get BUILD jobs:", err)
			}

			for _, newJob := range *newJobs {
				agentId := strconv.Itoa(newJob.AgentId)
				jobId := strconv.Itoa(newJob.ID)

				userProfileKey := "ciagent:" + agentId
				userProfileData := map[string]interface{}{
					"name":  "John Doe",
					"age":   "32",
					"jobId": jobId,
				}
				_, err = g.Redis().HSet(ctx, userProfileKey, userProfileData)
				if err != nil {
					g.Log().Fatal(ctx, err)
				}

				_, err = g.Redis().Expire(ctx, userProfileKey, 10*60)
				if err != nil {
					g.Log().Fatal(ctx, err)
				}
			}

		}
	}()
}
