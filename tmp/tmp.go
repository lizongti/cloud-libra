package main

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

func main() {
	v := viper.New()

	data1 := map[string]interface{}{
		"a": []interface{}{
			map[string]interface{}{
				"a": "a",
			},
			map[string]interface{}{
				"a": "a",
			},
		},
	}

	v.MergeConfigMap(data1)

	fmt.Println(v.GetString("a.3.a"))
	// v.Set("a.0.a", "b")
	// v.Set("a.1.a", "b")
	data, err := json.Marshal(v.AllSettings())
	if err != nil {
		panic(err)
	}

	gjson.Get(string(data), "a").IsArray()

	ForEach(func(key, value gjson.Result) bool {
		fmt.Println(key, value)
		return true
	})
}
