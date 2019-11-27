package goapiclient

import "encoding/json"

// Go representation
type Contact struct {
	Name          string                 `json:"name"`
	UnknownFields map[string]interface{} `json:"-"`
}

type contact struct {
	Name string `json:"first_name"`
}

func (c *Contact) CompatibleUnmarshal(data []byte) error {
	if err := json.Unmarshal(data, c); err != nil {
		return err
	}

	c.UnknownFields = make(map[string]interface{})

	if err := json.Unmarshal(data, &c.UnknownFields); err != nil {
		return err
	} else {
		delete(c.UnknownFields, "name")
		delete(c.UnknownFields, "first_name")
	}

	if c.Name == "" {
		var legacy contact
		if err := json.Unmarshal(data, &legacy); err != nil {
			return err
		}

		c.Name = legacy.Name
	}

	return nil
}
