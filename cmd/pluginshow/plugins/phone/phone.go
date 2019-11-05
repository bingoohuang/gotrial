package main

type phone string

func (p *phone) Parse([]byte) (meta map[string]string, data map[string]float64, err error) {
	meta = map[string]string{"key1": "phonephone"}
	data = map[string]float64{"key1": 2}
	return meta, data, nil
}

var Phone phone
