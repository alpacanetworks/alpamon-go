package io

import (
	"context"
	"time"

	"github.com/alpacanetworks/alpamon-go/pkg/collector/check/base"
	"github.com/alpacanetworks/alpamon-go/pkg/db/ent"
	"github.com/alpacanetworks/alpamon-go/pkg/db/ent/diskioperhour"
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

	queryset, err := c.getDiskIOPerHour(ctx)
	if err != nil {
		checkError.GetQueryError = err
	}

	metric := base.MetricData{
		Type: base.DISK_IO_PER_DAY,
		Data: []base.CheckResult{},
	}
	if checkError.GetQueryError == nil {
		for _, row := range queryset {
			data := base.CheckResult{
				Timestamp:      time.Now(),
				Device:         row.Device,
				PeakWriteBytes: uint64(row.PeakWriteBytes),
				PeakReadBytes:  uint64(row.PeakReadBytes),
				AvgWriteBytes:  uint64(row.AvgWriteBytes),
				AvgReadBytes:   uint64(row.AvgReadBytes),
			}
			metric.Data = append(metric.Data, data)

			if err := c.deleteDiskIOPerHour(ctx); err != nil {
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

func (c *Check) getDiskIOPerHour(ctx context.Context) ([]base.DiskIOQuerySet, error) {
	client := c.GetClient()
	now := time.Now()
	from := now.Add(-24 * time.Hour)

	var queryset []base.DiskIOQuerySet
	err := client.DiskIOPerHour.Query().
		Where(diskioperhour.TimestampGTE(from), diskioperhour.TimestampLTE(now)).
		GroupBy(diskioperhour.FieldDevice).
		Aggregate(
			ent.As(ent.Max(diskioperhour.FieldPeakReadBytes), "peak_read_bytes"),
			ent.As(ent.Max(diskioperhour.FieldPeakWriteBytes), "peak_write_bytes"),
			ent.As(ent.Mean(diskioperhour.FieldAvgReadBytes), "avg_read_bytes"),
			ent.As(ent.Mean(diskioperhour.FieldAvgWriteBytes), "avg_write_bytes"),
		).
		Scan(ctx, &queryset)
	if err != nil {
		log.Debug().Msg(err.Error())
		return queryset, err
	}

	return queryset, nil
}

func (c *Check) deleteDiskIOPerHour(ctx context.Context) error {
	client := c.GetClient()
	now := time.Now()
	from := now.Add(-24 * time.Hour)

	_, err := client.DiskIOPerHour.Delete().
		Where(diskioperhour.TimestampGTE(from), diskioperhour.TimestampLTE(now)).
		Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
