package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"time"
)

func main() {
	type Person struct {
		Name  string
		Age   int
		Ratio float32
	}
	var (
		a []string
		b []interface{}
		c []io.Writer
		d []Person
	)

	fmt.Println(fill(&a), a) // pass
	fmt.Println(fill(&b), b) // pass
	fmt.Println(fill(&d), d) // pass
	fmt.Println(fill(&c), c) // fail

	jso := `{"results":[{"statement_id":0,"series":[{"name":"mem","tags":{"host":"beta-bq"},"columns":["time","available","available_percent","used","used_percent"],"values":[["2019-04-10T00:30:00Z",6780.0203125,86.67115414384232,898.2842447916667,11.483052949234583]]},{"name":"mem","tags":{"host":"beta-ca"},"columns":["time","available","available_percent","used","used_percent"],"values":[["2019-04-10T00:30:00Z",1954.0084635416667,24.979122733690968,5708.663671875,72.97686431033614]]},{"name":"mem","tags":{"host":"beta-hetong"},"columns":["time","available","available_percent","used","used_percent"],"values":[["2019-04-10T00:30:00Z",8579.579557291667,54.0070803835492,7180.895052083333,45.202585244878705]]},{"name":"mem","tags":{"host":"beta-qianming"},"columns":["time","available","available_percent","used","used_percent"],"values":[["2019-04-10T00:30:00Z",1510.322265625,39.84152728988894,2131.5911458333335,56.230281934207206]]},{"name":"mem","tags":{"host":"beta-qianzhang"},"columns":["time","available","available_percent","used","used_percent"],"values":[["2019-04-10T00:30:00Z",4243.741015625,54.24998389575031,3420.365234375,43.72433619281555]]},{"name":"mem","tags":{"host":"tencent-beta02"},"columns":["time","available","available_percent","used","used_percent"],"values":[["2019-04-10T00:30:00Z",5941.044270833333,37.3938036913039,10135.829947916667,63.7963997645268]]},{"name":"mem","tags":{"host":"tencent-beta03"},"columns":["time","available","available_percent","used","used_percent"],"values":[["2019-04-10T00:30:00Z",8451.113802083333,53.190260751669065,7655.479947916667,48.182640080941866]]},{"name":"mem","tags":{"host":"tencent-beta04"},"columns":["time","available","available_percent","used","used_percent"],"values":[["2019-04-10T00:30:00Z",870.2723958333333,22.947172035845103,2759.5984375,72.76455095938637]]},{"name":"mem","tags":{"host":"tencent-beta16"},"columns":["time","available","available_percent","used","used_percent"],"values":[["2019-04-10T00:30:00Z",2468.683984375,31.55795608832669,5218.279427083334,66.70689082070024]]},{"name":"mem","tags":{"host":"tencent-beta17"},"columns":["time","available","available_percent","used","used_percent"],"values":[["2019-04-10T00:30:00Z",11059.716666666667,69.61914660821645,4796.33125,30.192137695690224]]},{"name":"mem","tags":{"host":"tencent-beta18"},"columns":["time","available","available_percent","used","used_percent"],"values":[["2019-04-10T00:30:00Z",2899.252864583333,18.250332835565583,12820.065234375,80.70025914577543]]},{"name":"mem","tags":{"host":"tencent-beta19"},"columns":["time","available","available_percent","used","used_percent"],"values":[["2019-04-10T00:30:00Z",7888.863411458334,49.65913277596132,8042.551302083333,50.62657345385657]]}]}]}`

	type MemSeries struct {
		Host             string    `iql:"host/ip"`
		Time             time.Time `iql:"time"`
		Available        float32   `iql:"available"`
		AvailablePercent float32   `iql:"available_percent"`
		Used             float32   `iql:"used"`
		UsedPercent      float32   `iql:"used_percent"`
	}

	var memSeries []MemSeries

	ParseIqlResults(jso, func(m map[string]interface{}) {
		ms := MemSeries{}
		Map2Struct(m, &ms)
		fmt.Printf("%+v\n", ms)
	})

	fmt.Println(memSeries)

}

// UnmarshalError return error from errMsg
func UnmarshalError(errMsg *json.RawMessage) error {
	var s string
	_ = json.Unmarshal(*errMsg, &s)
	return errors.New(s)
}

func fill(ptr interface{}) error {
	slice, err := GetSliceByPtr(ptr)
	if err != nil {
		return fmt.Errorf("can't fill non-slice value")
	}
	slice.Set(reflect.MakeSlice(slice.Type(), 0, 10))
	elemType := slice.Type().Elem()

	if elemType.Kind() == reflect.Struct {
		ii := RandomInt(10) + 1
		EnsureSliceLen(slice, ii)
		for i := 0; i < ii; i++ {
			fillSlice(slice, i)
		}

		return nil
	}

	// validate the type of the slice. see below.
	if !CanAssign(elemType, reflect.String) {
		return fmt.Errorf("can't assign string to slice elements")
	}

	ii := RandomInt(10) + 1
	EnsureSliceLen(slice, ii)
	for i := 0; i < ii; i++ {
		slice.Index(i).Set(reflect.ValueOf(RandomString(10)))
	}
	return nil
}

func fillSlice(slice reflect.Value, i int) {
	structFields := CachedStructFields(slice.Index(0).Type(), "iql")

	for _, sf := range structFields {
		field := slice.Index(i).Field(sf.Index)
		switch sf.Kind {
		case reflect.String:
			field.SetString(RandomString(10))
		case reflect.Int:
			field.SetInt(RandomInt64())
		case reflect.Float32:
			field.SetFloat(RandomFloat64())
		}
	}
}
