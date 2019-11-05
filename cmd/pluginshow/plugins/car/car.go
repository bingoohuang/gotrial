package main

type car string

func (c *car) Parse([]byte) (meta map[string]string, data map[string]float64, err error) {
	meta = map[string]string{"key1": "carcard"}
	data = map[string]float64{"key1": 1}
	return meta, data, nil
}

var Car car
