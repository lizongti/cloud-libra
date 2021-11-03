package encoding_test

import (
	"reflect"
	"testing"

	"github.com/aceaura/libra/encoding"
	"github.com/aceaura/libra/magic"
)

func TestEncoding(t *testing.T) {
	type TestStruct struct {
		Integer int
		String  string
		Bool    bool
		Slice   []byte
	}

	ts1 := &TestStruct{
		Integer: 1,
		String:  "this is test text",
		Bool:    false,
		Slice:   []byte("this is slice"),
	}
	e1 := encoding.NewEncoding().WithEncoder(
		"json.base64.lazy", magic.SeparatorPeriod, magic.SeparatorUnderscore,
	).WithDecoder(
		"lazy.base64.xml", magic.SeparatorPeriod, magic.SeparatorUnderscore,
	)
	e2 := e1.Reverse()

	data1, err := encoding.Encode(e1, ts1)
	if err != nil {
		t.Fatal(err)
	}
	ts2 := &TestStruct{}
	err = encoding.Decode(e2, data1, ts2)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(ts1, ts2) {
		t.Fatalf("expecting ts1 %v equals ts2 %v", ts1, ts2)
	}

	data2Bytes, err := encoding.Marshal(e2, ts2)
	if err != nil {
		t.Fatal(err)
	}
	ts3 := &TestStruct{}
	err = encoding.Unmarshal(e1, encoding.MakeBytes(data2Bytes), ts3)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(ts2, ts3) {
		t.Fatalf("expecting ts2 %v equals ts3 %v", ts2, ts3)
	}
	t.Logf("data1: %s", string(data1))
	t.Logf("data2: %s", string(data2Bytes.Data))
	t.Logf("ts1: %+v", ts1)
	t.Logf("ts2: %+v", ts2)
	t.Logf("ts3: %+v", ts3)
}
