package hierarchy

import (
	"errors"

	"github.com/spf13/viper"
)

// examples:
//  - (project:aries) (runtime:default)
//     project = aries, runtime = default
//  - [lua]{a=1,b=2} [json]{"a":1,"b":2}
//     a = 1, b = 2
//  - {args}
//     Get Hierarchy from args
//  - {flags}
//     Get Hierarchy from flags
//  - {stdio:yaml}
//     Get hierarchy from stdio.
//  - {env:Aries}
// 	   Get hierarchy from env with prefix Aries.
//  - {cluster:default}
//     Get Config from cluster. Ger hierarchy from hierarchy.
//  - {hierarchy:recursive}
//     Get hierarchy from hierarchy key
//  - <http://filestone.com/file.json>
// 	   Get hierarchy from http.
//  - <file:///E:/Filename/file.ini> ...
//	   Get hierarchy from local file.
//  - <etcd://192.168.1.2:2379@usr:passwd/aries/hierarchy>
//     Get hierarchy from etcd.
//  - <redis://192.168.1.2:6379@usr:passwd/0/aries/hierarchy>
//     Get hierarchy from redis.

var (
	ErrArgNotMatch   = errors.New("arg not match")
	ErrInvalidString = errors.New("invalid string")
)

func IsArgNotMatch(err error) bool {
	return errors.Is(err, ErrArgNotMatch)
}

type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(config string) (*viper.Viper, error) {
	panic("implement me")
}
