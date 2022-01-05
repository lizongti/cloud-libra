package redis

import "github.com/gomodule/redigo/redis"

type (
	SlowLog = redis.SlowLog
)

var (
	ErrNil = redis.ErrNil
)

func Int(reply interface{}, err error) (int, error) {
	return redis.Int(reply, err)
}

func Int64(reply interface{}, err error) (int64, error) {
	return redis.Int64(reply, err)
}

func Uint64(reply interface{}, err error) (uint64, error) {
	return redis.Uint64(reply, err)
}

func Float64(reply interface{}, err error) (float64, error) {
	return redis.Float64(reply, err)
}

func String(reply interface{}, err error) (string, error) {
	return redis.String(reply, err)
}

func Bytes(reply interface{}, err error) ([]byte, error) {
	return redis.Bytes(reply, err)
}

func Bool(reply interface{}, err error) (bool, error) {
	return redis.Bool(reply, err)
}

func Values(reply interface{}, err error) ([]interface{}, error) {
	return redis.Values(reply, err)
}

func Float64s(reply interface{}, err error) ([]float64, error) {
	return redis.Float64s(reply, err)
}

func Strings(reply interface{}, err error) ([]string, error) {
	return redis.Strings(reply, err)
}

func ByteSlices(reply interface{}, err error) ([][]byte, error) {
	return redis.ByteSlices(reply, err)
}

func Int64s(reply interface{}, err error) ([]int64, error) {
	return redis.Int64s(reply, err)
}

func Ints(reply interface{}, err error) ([]int, error) {
	return redis.Ints(reply, err)
}

func IntMap(reply interface{}, err error) (map[string]int, error) {
	return redis.IntMap(reply, err)
}

func Int64Map(reply interface{}, err error) (map[string]int64, error) {
	return redis.Int64Map(reply, err)
}

func Positions(reply interface{}, err error) ([]*[2]float64, error) {
	return redis.Positions(reply, err)
}

func Uint64s(reply interface{}, err error) ([]uint64, error) {
	return redis.Uint64s(reply, err)
}
func Uint64Map(reply interface{}, err error) (map[string]uint64, error) {
	return redis.Uint64Map(reply, err)
}

func SlowLogs(reply interface{}, err error) ([]SlowLog, error) {
	return redis.SlowLogs(reply, err)
}
