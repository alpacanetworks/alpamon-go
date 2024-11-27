package usage

import (
	"context"
	"time"

	"github.com/alpacanetworks/alpamon-go/pkg/collector/check/base"
	"github.com/alpacanetworks/alpamon-go/pkg/db/ent"
	"github.com/alpacanetworks/alpamon-go/pkg/db/ent/diskusageperhour"
	"github.com/rs/zerolog/log"
)

type Check struct {
	base.BaseCheck
}

func NewCheck(args *base.CheckArgs) base.CheckStrategy {
	return &Check{
		BaseCheck: base.NewBaseCheck(args),
	}
}

func (c *Check) Execute(ctx context.Context) {
	var checkError base.CheckError

	queryset, err := c.getDiskUsagePerHour(ctx)
	if err != nil {
		checkError.GetQueryError = err
	}

	metric := base.MetricData{
		Type: base.DISK_USAGE_PER_DAY,
		Data: []base.CheckResult{},
	}
	if checkError.GetQueryError == nil {
		for _, row := range queryset {
			data := base.CheckResult{
				Timestamp:  time.Now(),
				Device:     row.Device,
				MountPoint: row.MountPoint,
				PeakUsage:  row.Max,
				AvgUsage:   row.AVG,
			}
			metric.Data = append(metric.Data, data)

			if err := c.deleteDiskUsagePerHour(ctx); err != nil {
				checkError.DeleteQueryError = err
			}
		}
	}

	if ctx.Err() != nil {
		return
	}

	buffer := c.GetBuffer()
	if checkError.GetQueryError != nil || checkError.DeleteQueryError != nil {
		buffer.FailureQueue <- metric
	} else {
		buffer.SuccessQueue <- metric
	}
}

func (c *Check) getDiskUsagePerHour(ctx context.Context) ([]base.DiskUsageQuerySet, error) {
	client := c.GetClient()
	now := time.Now()
	from := now.Add(-24 * time.Hour)

	var queryset []base.DiskUsageQuerySet
	err := client.DiskUsagePerHour.Query().
		Where(diskusageperhour.TimestampGTE(from), diskusageperhour.TimestampLTE(now)).
		GroupBy(diskusageperhour.FieldDevice, diskusageperhour.FieldMountPoint).
		Aggregate(
			ent.Max(diskusageperhour.FieldPeakUsage),
			ent.Mean(diskusageperhour.FieldAvgUsage),
		).
		Scan(ctx, &queryset)
	if err != nil {
		log.Debug().Msg(err.Error())
		return queryset, err
	}

	return queryset, nil
}

func (c *Check) deleteDiskUsagePerHour(ctx context.Context) error {
	client := c.GetClient()
	now := time.Now()
	from := now.Add(-24 * time.Hour)

	_, err := client.DiskUsagePerHour.Delete().
		Where(diskusageperhour.TimestampGTE(from), diskusageperhour.TimestampLTE(now)).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
