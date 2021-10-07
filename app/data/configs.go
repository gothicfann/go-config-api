package data

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/imdario/mergo"
	"github.com/tidwall/gjson"
)

type Config struct {
	Name     string                 `json:"name"`
	Metadata map[string]interface{} `json:"metadata"`
}

type ConfigList []*Config

var ErrConfigNotFound = fmt.Errorf("Config not found")
var ErrConfigAlreadyExists = fmt.Errorf("Config already exists")

// configFound is helper function which searches for config with specified name parameter.
// it returns index of config element and true if configs was found, otherwise -1, false is returned
func configFound(n string) (int, bool) {
	for i, v := range db {
		if v.Name == n {
			return i, true
		}
	}
	return -1, false
}

// ToJSON method encodes ConfigList and writes to io.Writer interface
func (c *ConfigList) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(c)
}

// ToJSON method encodes Config and writes to io.Writer interface
func (c *Config) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(c)
}

// FromJSON method encodes Config and reads from io.Reader interface
func (c *Config) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(c)
}

// GetConfigs returns database
func GetConfigs() ConfigList {
	return db
}

// GetConfig requires name of config and returns *Config if found, otherwise returns ErrConfigNotFound error.
func GetConfig(n string) (*Config, error) {
	if i, ok := configFound(n); ok {
		return db[i], nil
	}
	return nil, ErrConfigNotFound
}

// AddConfig adds *Config to database. It returns error if config is already present.
func AddConfig(c *Config) error {
	if _, ok := configFound(c.Name); !ok {
		db = append(db, c)
		return nil
	}
	return ErrConfigAlreadyExists
}

// DeleteConfig deletes config from database by specified name parameter.
func DeleteConfig(n string) {
	if i, ok := configFound(n); ok {
		db[i] = db[len(db)-1]
		db = db[:len(db)-1]
	}
}

// PutConfig updates config in database if it exists, otherwise returns ErrConfigNotFound error.
func PutConfig(c *Config) error {
	if i, ok := configFound(c.Name); ok {
		db[i] = c
		return nil
	}
	return ErrConfigNotFound
}

// PatchConfig patches config in database if it exists, otherwise returns error.
func PatchConfig(c *Config) error {
	d, err := GetConfig(c.Name)
	if err != nil {
		return err
	}
	err = mergo.Merge(d, *c, mergo.WithOverride)
	if err != nil {
		return err
	}

	return nil
}

// QueryConfig receives key and value string type parameters to search database.
// Key should be json path (e.g: "metadata.monitoring.enabled")
// Value should be desired value of the key (e.g: "true")
// It returns empty ConfigList if no match.
func QueryConfig(key, value string) (ConfigList, error) {
	cl := ConfigList{}
	for i, c := range db {
		bs, err := json.Marshal(c)
		if err != nil {
			return nil, err
		}
		if gjson.Get(string(bs), key).String() == value {
			cl = append(cl, db[i])
		}
	}
	return cl, nil
}

// db is playground database
var db = ConfigList{
	{
		Name: "datacenter-1",
		Metadata: map[string]interface{}{
			"monitoring": map[string]interface{}{
				"enabled": "true",
			},
			"limits": map[string]interface{}{
				"cpu": map[string]interface{}{
					"enabled": "true",
					"value":   "250m",
				},
			},
		},
	},
	{
		Name: "datacenter-2",
		Metadata: map[string]interface{}{
			"monitoring": map[string]interface{}{
				"enabled": "false",
			},
			"limits": map[string]interface{}{
				"cpu": map[string]interface{}{
					"enabled": "true",
					"value":   "260m",
				},
			},
		},
	},
	{
		Name: "burger-nutrition",
		Metadata: map[string]interface{}{
			"calories": 230,
			"fats": map[string]interface{}{
				"enabled": "false",
			},
			"carbohydrates": map[string]interface{}{
				"dietary-fiber": "4g",
				"sugars":        "1g",
			},
			"allergens": map[string]interface{}{
				"nuts":     "false",
				"searfood": "false",
				"eggs":     "true",
			},
		},
	},
}
