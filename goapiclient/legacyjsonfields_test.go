package goapiclient_test

import (
	"github.com/bingoohuang/golang-trial/goapiclient"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContact_UnmarshalJSON(t *testing.T) {
	// JSON payload
	jsonPayload1 := `{ "name": "anthony" }`
	jsonPayload2 := `{ "first_name": "anthony" }`

	var c goapiclient.Contact

	assert.Nil(t, c.CompatibleUnmarshal([]byte(jsonPayload1)))
	assert.Equal(t, "anthony", c.Name)

	var c2 goapiclient.Contact

	assert.Nil(t, c2.CompatibleUnmarshal([]byte(jsonPayload2)))
	assert.Equal(t, c, c2)

	jsonPayload3 := `{ "name": "anthony", "age": 26 }`

	var c3 goapiclient.Contact
	assert.Nil(t, c3.CompatibleUnmarshal([]byte(jsonPayload3)))
	assert.Equal(t, goapiclient.Contact{
		Name:          "anthony",
		UnknownFields: map[string]interface{}{"age": float64(26)},
	}, c3)
}
