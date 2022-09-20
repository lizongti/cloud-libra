package hierarchy

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var ErrModeAssetsArgsNotEnough = errors.New("mode assets args not enough")

type PlainAssets struct{}

func (*PlainAssets) Plain(s string) (map[string][]byte, error) {
	assets := make(map[string][]byte)
	assets[""] = []byte(s)

	return assets, nil
}

type ModeAssets struct {
	hierarchy *Hierarchy
}

func (*ModeAssets) parse(s string) []string {
	return strings.Split(s, ":")
}

func (*ModeAssets) Args(s string) (map[string][]byte, error) {
	assets := make(map[string][]byte)
	assets[""] = []byte(s)

	return assets, nil
}

func (*ModeAssets) Env(s string) (map[string][]byte, error) {
	assets := make(map[string][]byte)
	assets[""] = []byte(s)

	return assets, nil
}

func (*ModeAssets) Flags(s string) (map[string][]byte, error) {
	assets := make(map[string][]byte)
	assets[""] = []byte(s)

	return assets, nil
}

func (*ModeAssets) Stdin(s string) (map[string][]byte, error) {
	data, err := io.ReadAll(bufio.NewReader(os.Stdin))
	if err != nil {
		return nil, err
	}

	assets := make(map[string][]byte)
	assets[""] = []byte(data)

	return assets, nil
}

func (ma *ModeAssets) Hierarchy(s string) (map[string][]byte, error) {
	strs := ma.parse(s)
	if len(strs) < 2 {
		return nil, fmt.Errorf("%w: %s", ErrModeAssetsArgsNotEnough, s)
	}
	key := strs[1]
	val := ma.hierarchy.GetString(key)

	assets := make(map[string][]byte)
	assets[""] = []byte(val)

	return assets, nil
}

type URLAssets struct{}

func (*URLAssets) File(s string) (map[string][]byte, error) {
	panic("implement me")
}

func (*URLAssets) Http(s string) (map[string][]byte, error) {
	panic("implement me")
}

func (*URLAssets) Https(s string) (map[string][]byte, error) {
	panic("implement me")
}
