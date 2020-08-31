package logging

import (
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"reflect"
	"testing"
)

func TestBasic(t *testing.T) {
	formatter := LogTagFormatter{}

	b, err := formatter.Format(logrus.WithFields(logrus.Fields{
		"stringValue": "value",
		"intValue":    123,
		"floatValue":  123.09,
		"boolValue":   true,
		"arrayValue":  [2]int{1, 2},
		"sliceValue":  []int{1, 2, 3},
	}))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["stringValue"] != "value" {
		t.Fatal("Error string field not set")
	}
	if (entry["intValue"]).(float64) != 123 {
		t.Fatal("Error int field not set")
	}
	if entry["floatValue"].(float64) != 123.09 {
		t.Fatal("Error float field not set")
	}
	if entry["boolValue"] != true {
		t.Fatal("Error bool field not set")
	}
	if reflect.DeepEqual(entry["arrayValue"], [2]int{1, 2}) {
		t.Fatal("Error array field not set")
	}
	if reflect.DeepEqual(entry["sliceValue"], []int{1, 2, 3}) {
		t.Fatal("Error slice field not set")
	}
}

func TestBasicPointers(t *testing.T) {
	formatter := LogTagFormatter{}
	stringValue := "value"
	intValue := 123
	floatValue := 123.09
	boolValue := true

	b, err := formatter.Format(logrus.WithFields(logrus.Fields{
		"stringPtr": &stringValue,
		"intPtr":    &intValue,
		"floatPtr":  &floatValue,
		"boolPtr":   &boolValue,
	}))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["stringPtr"] != "value" {
		t.Fatal("Error string field not set")
	}
	if (entry["intPtr"]).(float64) != 123 {
		t.Fatal("Error int field not set")
	}
	if entry["floatPtr"].(float64) != 123.09 {
		t.Fatal("Error float field not set")
	}
	if entry["boolPtr"] != true {
		t.Fatal("Error bool field not set")
	}
}

func TestStructEmpty(t *testing.T) {
	formatter := LogTagFormatter{}
	s := struct{}{}

	b, err := formatter.Format(logrus.WithField("s", s))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	emptyStruct := entry["s"].(map[string]interface{})

	if len(emptyStruct) != 0 {
		t.Fatal("Error empty struct not empty")
	}
}

func TestStructWithBasicAndPointers(t *testing.T) {
	formatter := LogTagFormatter{}
	stringValue := "value"
	intValue := 123
	floatValue := 123.09
	boolValue := true
	s := struct {
		StringValue   string   `log:"stringValue"`
		IntValue      int      `log:"intValue"`
		FloatValue    float64  `log:"floatValue"`
		BoolValue     bool     `log:"boolValue"`
		ArrayValue    [2]int   `log:"arrayValue"`
		SliceValue    []int    `log:"sliceValue"`
		StringPtr     *string  `log:"stringPtr"`
		IntPtr        *int     `log:"intPtr"`
		FloatPtr      *float64 `log:"floatPtr"`
		BoolPtr       *bool    `log:"boolPtr"`
		SkippedField  int      `log:"-"`
		EmptyLogField int      `log:""`
		NoTagField    int
	}{
		StringValue:   "value",
		IntValue:      123,
		FloatValue:    123.09,
		BoolValue:     true,
		ArrayValue:    [2]int{1, 2},
		SliceValue:    []int{1, 2, 3},
		StringPtr:     &stringValue,
		IntPtr:        &intValue,
		FloatPtr:      &floatValue,
		BoolPtr:       &boolValue,
		SkippedField:  1,
		EmptyLogField: 2,
		NoTagField:    3,
	}

	b, err := formatter.Format(logrus.WithFields(logrus.Fields{
		"basicStruct": s,
	}))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	basicStructEntry := entry["basicStruct"].(map[string]interface{})

	if basicStructEntry["stringValue"] != "value" {
		t.Fatal("Error string field not set")
	}
	if (basicStructEntry["intValue"]).(float64) != 123 {
		t.Fatal("Error int field not set")
	}
	if basicStructEntry["floatValue"].(float64) != 123.09 {
		t.Fatal("Error float field not set")
	}
	if basicStructEntry["boolValue"] != true {
		t.Fatal("Error bool field not set")
	}
	if reflect.DeepEqual(basicStructEntry["arrayValue"], [2]int{1, 2}) {
		t.Fatal("Error array field not set")
	}
	if reflect.DeepEqual(basicStructEntry["sliceValue"], []int{1, 2, 3}) {
		t.Fatal("Error slice field not set")
	}
	if basicStructEntry["stringPtr"] != "value" {
		t.Fatal("Error string field not set")
	}
	if (basicStructEntry["intPtr"]).(float64) != 123 {
		t.Fatal("Error int field not set")
	}
	if basicStructEntry["floatPtr"].(float64) != 123.09 {
		t.Fatal("Error float field not set")
	}
	if basicStructEntry["boolPtr"] != true {
		t.Fatal("Error bool field not set")
	}
	if _, ok := basicStructEntry["SkippedField"]; ok {
		t.Fatal("Error skipped field set")
	}
	if _, ok := basicStructEntry["EmptyLogField"]; ok {
		t.Fatal("Error empty log field set")
	}
	if _, ok := basicStructEntry["NoTagField"]; ok {
		t.Fatal("Error no tag field set")
	}
}

func TestStructNested(t *testing.T) {
	formatter := LogTagFormatter{}
	type A struct {
		A int `log:"aa"`
		B int
	}
	type S struct {
		A A `log:"sa"`
		B A
		C *A `log:"sc"`
		D *A
	}
	s := S{
		A: A{
			A: 1,
			B: 2,
		},
		B: A{
			A: 3,
			B: 4,
		},
		C: &A{
			A: 5,
			B: 6,
		},
		D: &A{
			A: 7,
			B: 8,
		},
	}

	b, err := formatter.Format(logrus.WithFields(logrus.Fields{"s": s}))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	nestedStructEntry := entry["s"].(map[string]interface{})

	if sa, ok := nestedStructEntry["sa"].(map[string]interface{}); true {
		if !ok {
			t.Fatal("Error nested struct not set")
		}
		if sa["aa"].(float64) != 1 {
			t.Fatal("Error nested struct field not set")
		}
	}

	if sc, ok := nestedStructEntry["sc"].(map[string]interface{}); true {
		if !ok {
			t.Fatal("Error nested struct pointer not set")
		}
		if sc["aa"].(float64) != 5 {
			t.Fatal("Error nested struct pointer field not set")
		}
	}
}

func TestStructInline(t *testing.T) {
	formatter := LogTagFormatter{}
	type A struct {
		A int `log:"aa"`
		B int
	}
	type B struct {
		A int `log:"ba"`
		B int
	}
	type S struct {
		A  `log:"sa"`
		*B `log:"sb"`
	}
	s := S{
		A: A{
			A: 1,
			B: 2,
		},
		B: &B{
			A: 3,
			B: 4,
		},
	}

	b, err := formatter.Format(logrus.WithField("s", s))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	inlineStruct := entry["s"].(map[string]interface{})

	if sa, ok := inlineStruct["sa"].(map[string]interface{}); true {
		if !ok {
			t.Fatal("Error inline struct not set")
		}
		if sa["aa"].(float64) != 1 {
			t.Fatal("Error inline struct field not set")
		}
	}

	if sb, ok := inlineStruct["sb"].(map[string]interface{}); true {
		if !ok {
			t.Fatal("Error inline struct pointer not set")
		}
		if sb["ba"].(float64) != 3 {
			t.Fatal("Error inline struct pointer field not set")
		}
	}
}

func TestStructMultipleFields(t *testing.T) {
	formatter := LogTagFormatter{}
	type A struct {
		A int `log:"aa"`
		B int
	}
	type B struct {
		A int `log:"ba"`
		B int
	}
	sa := A{
		A: 1,
		B: 2,
	}
	sb := &B{
		A: 3,
		B: 4,
	}

	b, err := formatter.Format(logrus.
		WithField("sa", sa).
		WithField("sb", sb))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	saEntry := entry["sa"].(map[string]interface{})
	sbEntry := entry["sb"].(map[string]interface{})

	if saEntry["aa"].(float64) != 1 {
		t.Fatal("Error multiple struct with field aa not set")
	}
	if _, ok := saEntry["b"]; ok {
		t.Fatal("Error multiple struct with field b set")
	}
	if sbEntry["ba"].(float64) != 3 {
		t.Fatal("Error multiple struct with field ba not set")
	}
	if _, ok := sbEntry["b"]; ok {
		t.Fatal("Error multiple struct with field b set")
	}
}

func TestLogTagFormatter_Format_DefaultFormatter(t *testing.T) {
	formatter := LogTagFormatter{defaultFormatter: &logrus.TextFormatter{}}
	s := struct {
		A int `log:"a"`
	}{
		A: 1,
	}

	b, err := formatter.Format(logrus.WithField("s", s))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	if !bytes.Contains(b, []byte(`a:1`)) {
		t.Fatal("Error default formatter not set")
	}
}
