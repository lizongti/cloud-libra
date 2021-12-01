package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// Nil is a wrapper for redis.Nil
const Nil = redis.Nil

type (
	// Client is a wrapper for redis.Client
	Client = redis.Client

	// Options is a wrapper for redis.Options
	Options = redis.Options

	// Pipeliner is a wrapper for redis.Pipeliner
	Pipeliner = redis.Pipeliner

	// Pipeline is a wrapper for redis.Pipeline
	Pipeline = redis.Pipeline

	// Cmder is a wrapper for redis.Cmder
	Cmder = redis.Cmder

	// Cmd is a wrapper for redis.Cmd
	Cmd = redis.Cmd

	// SliceCmd is a wrapper for redis.SliceCmd
	SliceCmd = redis.SliceCmd

	// StatusCmd is a wrapper for redis.StatusCmd
	StatusCmd = redis.StatusCmd

	// IntCmd is a wrapper for redis.IntCmd
	IntCmd = redis.IntCmd

	// DurationCmd is a wrapper for redis.DurationCmd
	DurationCmd = redis.DurationCmd

	// TimeCmd is a wrapper for redis.TimeCmd
	TimeCmd = redis.TimeCmd

	// BoolCmd is a wrapper for redis.BoolCmd
	BoolCmd = redis.BoolCmd

	// StringCmd is a wrapper for redis.StringCmd
	StringCmd = redis.StringCmd

	// FloatCmd is a wrapper for redis.FloatCmd
	FloatCmd = redis.FloatCmd

	// StringSliceCmd is a wrapper for redis.StringSliceCmd
	StringSliceCmd = redis.StringSliceCmd

	// BoolSliceCmd is a wrapper for redis.BoolSliceCmd
	BoolSliceCmd = redis.BoolSliceCmd

	// StringStringMapCmd is a wrapper for redis.StringStringMapCmd
	StringStringMapCmd = redis.StringStringMapCmd

	// StringIntMapCmd is a wrapper for redis.StringIntMapCmd
	StringIntMapCmd = redis.StringIntMapCmd

	// StringStructMapCmd is a wrapper for redis.StringStructMapCmd
	StringStructMapCmd = redis.StringStructMapCmd

	// XMessageSliceCmd is a wrapper for redis.XMessageSliceCmd
	XMessageSliceCmd = redis.XMessageSliceCmd

	// XStreamSliceCmd is a wrapper for redis.XStreamSliceCmd
	XStreamSliceCmd = redis.XStreamSliceCmd

	// XPendingCmd is a wrapper for redis.XPendingCmd
	XPendingCmd = redis.XPendingCmd

	// XPendingExtCmd is a wrapper for redis.XPendingExtCmd
	XPendingExtCmd = redis.XPendingExtCmd

	// ZSliceCmd is a wrapper for redis.ZSliceCmd
	ZSliceCmd = redis.ZSliceCmd

	// ZWithKeyCmd is a wrapper for redis.ZWithKeyCmd
	ZWithKeyCmd = redis.ZWithKeyCmd

	// ScanCmd is a wrapper for redis.ScanCmd
	ScanCmd = redis.ScanCmd

	// ClusterSlotsCmd is a wrapper for redis.ClusterSlotsCmd
	ClusterSlotsCmd = redis.ClusterSlotsCmd

	// GeoLocationCmd is a wrapper for redis.GeoLocationCmd
	GeoLocationCmd = redis.GeoLocationCmd

	// GeoPosCmd is a wrapper for redis.GeoPosCmd
	GeoPosCmd = redis.GeoPosCmd

	// CommandsInfoCmd is a wrapper for redis.CommandsInfoCmd
	CommandsInfoCmd = redis.CommandsInfoCmd

	// Z is a wrapper for redis.Z
	Z = redis.Z

	// ZRangeBy is a wrapper for redis.ZRangeBy
	ZRangeBy = redis.ZRangeBy

	// ZStore is a wrapper for redis.ZStore
	ZStore = redis.ZStore
)

// NewCmd is a wrapper for redis.NewCmd
func NewCmd(ctx context.Context, args ...interface{}) *redis.Cmd {
	return redis.NewCmd(ctx, args...)
}

// NewSliceCmd is a wrapper for redis.NewSliceCmd
func NewSliceCmd(ctx context.Context, args ...interface{}) *redis.SliceCmd {
	return redis.NewSliceCmd(ctx, args...)
}

// NewStatusCmd is a wrapper for redis.NewStatusCmd
func NewStatusCmd(ctx context.Context, args ...interface{}) *redis.StatusCmd {
	return redis.NewStatusCmd(ctx, args...)
}

// NewIntCmd is a wrapper for redis.NewIntCmd
func NewIntCmd(ctx context.Context, args ...interface{}) *redis.IntCmd {
	return redis.NewIntCmd(ctx, args...)
}

// NewDurationCmd is a wrapper for redis.NewDurationCmd
func NewDurationCmd(ctx context.Context, precision time.Duration, args ...interface{}) *redis.DurationCmd {
	return redis.NewDurationCmd(ctx, precision, args...)
}

// NewTimeCmd is a wrapper for redis.NewTimeCmd
func NewTimeCmd(ctx context.Context, args ...interface{}) *redis.TimeCmd {
	return redis.NewTimeCmd(ctx, args...)
}

// NewBoolCmd is a wrapper for redis.NewBoolCmd
func NewBoolCmd(ctx context.Context, args ...interface{}) *redis.BoolCmd {
	return redis.NewBoolCmd(ctx, args...)
}

// NewStringCmd is a wrapper for redis.NewStringCmd
func NewStringCmd(ctx context.Context, args ...interface{}) *redis.StringCmd {
	return redis.NewStringCmd(ctx, args...)
}

// NewFloatCmd is a wrapper for redis.NewFloatCmd
func NewFloatCmd(ctx context.Context, args ...interface{}) *redis.FloatCmd {
	return redis.NewFloatCmd(ctx, args...)
}

// NewStringSliceCmd is a wrapper for redis.NewStringSliceCmd
func NewStringSliceCmd(ctx context.Context, args ...interface{}) *redis.StringSliceCmd {
	return redis.NewStringSliceCmd(ctx, args...)
}

// NewBoolSliceCmd is a wrapper for redis.NewBoolSliceCmd
func NewBoolSliceCmd(ctx context.Context, args ...interface{}) *redis.BoolSliceCmd {
	return redis.NewBoolSliceCmd(ctx, args...)
}

// NewStringStringMapCmd is a wrapper for redis.NewStringStringMapCmd
func NewStringStringMapCmd(ctx context.Context, args ...interface{}) *redis.StringStringMapCmd {
	return redis.NewStringStringMapCmd(ctx, args...)
}

// NewStringIntMapCmd is a wrapper for redis.NewStringIntMapCmd
func NewStringIntMapCmd(ctx context.Context, args ...interface{}) *redis.StringIntMapCmd {
	return redis.NewStringIntMapCmd(ctx, args...)
}

// NewStringStructMapCmd is a wrapper for redis.NewStringStructMapCmd
func NewStringStructMapCmd(ctx context.Context, args ...interface{}) *redis.StringStructMapCmd {
	return redis.NewStringStructMapCmd(ctx, args...)
}

// NewXMessageSliceCmd is a wrapper for redis.NewXMessageSliceCmd
func NewXMessageSliceCmd(ctx context.Context, args ...interface{}) *redis.XMessageSliceCmd {
	return redis.NewXMessageSliceCmd(ctx, args...)
}

// NewXStreamSliceCmd is a wrapper for redis.NewXStreamSliceCmd
func NewXStreamSliceCmd(ctx context.Context, args ...interface{}) *redis.XStreamSliceCmd {
	return redis.NewXStreamSliceCmd(ctx, args...)
}

// NewXPendingCmd is a wrapper for redis.NewXPendingCmd
func NewXPendingCmd(ctx context.Context, args ...interface{}) *redis.XPendingCmd {
	return redis.NewXPendingCmd(ctx, args...)
}

// NewXPendingExtCmd is a wrapper for redis.NewXPendingExtCmd
func NewXPendingExtCmd(ctx context.Context, args ...interface{}) *redis.XPendingExtCmd {
	return redis.NewXPendingExtCmd(ctx, args...)
}

// NewZSliceCmd is a wrapper for redis.NewZSliceCmd
func NewZSliceCmd(ctx context.Context, args ...interface{}) *redis.ZSliceCmd {
	return redis.NewZSliceCmd(ctx, args...)
}

// NewZWithKeyCmd is a wrapper for redis.NewZWithKeyCmd
func NewZWithKeyCmd(ctx context.Context, args ...interface{}) *redis.ZWithKeyCmd {
	return redis.NewZWithKeyCmd(ctx, args...)
}

// NewClusterSlotsCmd is a wrapper for redis.NewClusterSlotsCmd
func NewClusterSlotsCmd(ctx context.Context, args ...interface{}) *redis.ClusterSlotsCmd {
	return redis.NewClusterSlotsCmd(ctx, args...)
}

// NewGeoLocationCmd is a wrapper for redis.NewGeoLocationCmd
func NewGeoLocationCmd(ctx context.Context, q *redis.GeoRadiusQuery, args ...interface{}) *redis.GeoLocationCmd {
	return redis.NewGeoLocationCmd(ctx, q, args...)
}

// NewGeoPosCmd is a wrapper for redis.NewGeoPosCmd
func NewGeoPosCmd(ctx context.Context, args ...interface{}) *redis.GeoPosCmd {
	return redis.NewGeoPosCmd(ctx, args...)
}

// NewCommandsInfoCmd is a wrapper for redis.NewCommandsInfoCmd
func NewCommandsInfoCmd(ctx context.Context, args ...interface{}) *redis.CommandsInfoCmd {
	return redis.NewCommandsInfoCmd(ctx, args...)
}
