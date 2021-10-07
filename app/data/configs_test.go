package data

import (
	"testing"

	"github.com/imdario/mergo"
	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	db = ConfigList{}
	input := "config-1"
	expected := &Config{
		Name: "config-1",
		Metadata: map[string]interface{}{
			"monitoring": map[string]interface{}{
				"enabled": "true",
			},
		},
	}

	AddConfig(expected)
	d, _ := GetConfig(input)

	assert.Equal(t, expected, d)
}

func TestGetConfigErr(t *testing.T) {
	db = ConfigList{}
	input := "config-1"

	_, err := GetConfig(input)

	assert.Equal(t, ErrConfigNotFound, err)
}

func TestAddConfig(t *testing.T) {
	db = ConfigList{}
	input := &Config{
		Name: "config-1",
		Metadata: map[string]interface{}{
			"monitoring": map[string]interface{}{
				"enabled": "true",
			},
		},
	}

	err := AddConfig(input)

	assert.Equal(t, nil, err)
}

func TestAddConfigErr(t *testing.T) {
	db = ConfigList{}
	input := &Config{
		Name: "config-1",
		Metadata: map[string]interface{}{
			"monitoring": map[string]interface{}{
				"enabled": "true",
			},
		},
	}

	AddConfig(input)
	err := AddConfig(input)

	assert.Equal(t, ErrConfigAlreadyExists, err)
}

func TestDeleteConfig(t *testing.T) {
	db = ConfigList{}
	input := "config-1"

	DeleteConfig(input)
	_, err := GetConfig(input)

	assert.Equal(t, ErrConfigNotFound, err)
}

func TestPutConfig(t *testing.T) {
	db = ConfigList{{
		Name: "config-1",
		Metadata: map[string]interface{}{
			"monitoring": map[string]interface{}{
				"enabled": "true",
				"grafana": "false",
			},
		},
	}}
	input := &Config{
		Name: "config-1",
		Metadata: map[string]interface{}{
			"monitoring": map[string]interface{}{
				"enabled": "true",
			},
		},
	}

	PutConfig(input)
	expected, _ := GetConfig(input.Name)

	assert.Equal(t, expected, input)
}

func TestPutConfigErr(t *testing.T) {
	db = ConfigList{}
	input := &Config{
		Name: "config-1",
		Metadata: map[string]interface{}{
			"monitoring": map[string]interface{}{
				"enabled": "true",
			},
		},
	}

	err := PutConfig(input)

	assert.Equal(t, ErrConfigNotFound, err)
}

func TestPatchConfig(t *testing.T) {
	db = ConfigList{
		{
			Name: "config-1",
			Metadata: map[string]interface{}{
				"monitoring": map[string]interface{}{
					"enabled": "true",
				},
			},
		},
	}
	input := &Config{
		Name: "config-1",
		Metadata: map[string]interface{}{
			"monitoring": map[string]interface{}{
				"enabled": "true",
				"grafana": "true",
			},
		},
	}

	expected, _ := GetConfig(input.Name)
	mergo.Merge(expected, *input, mergo.WithOverride)

	PatchConfig(input)
	d, _ := GetConfig(input.Name)

	assert.Equal(t, expected, d)
}

func TestPatchConfigErr(t *testing.T) {
	db = ConfigList{}
	input := &Config{
		Name: "config-1",
		Metadata: map[string]interface{}{
			"monitoring": map[string]interface{}{
				"enabled": "true",
			},
		},
	}

	err := PatchConfig(input)

	assert.Equal(t, ErrConfigNotFound, err)
}

func TestQueryConfig(t *testing.T) {
	db = ConfigList{
		{
			Name: "config-1",
			Metadata: map[string]interface{}{
				"monitoring": map[string]interface{}{
					"enabled": "true",
				},
			},
		},
		{
			Name: "config-2",
			Metadata: map[string]interface{}{
				"monitoring": map[string]interface{}{
					"enabled": "false",
				},
			},
		},
	}
	key := "metadata.monitoring.enabled"
	value := "true"
	expected := ConfigList{
		{
			Name: "config-1",
			Metadata: map[string]interface{}{
				"monitoring": map[string]interface{}{
					"enabled": "true",
				},
			},
		},
	}

	d, _ := QueryConfig(key, value)

	assert.Equal(t, expected, d)
}

func TestQueryConfigKeyNotFound(t *testing.T) {
	db = ConfigList{
		{
			Name: "config-1",
			Metadata: map[string]interface{}{
				"monitoring": map[string]interface{}{
					"enabled": "true",
				},
			},
		},
	}
	key := "metadata.monitoring.grafana"
	value := "true"
	expected := ConfigList{}

	d, _ := QueryConfig(key, value)

	assert.Equal(t, expected, d)
}

func TestQueryConfigValueNotFound(t *testing.T) {
	db = ConfigList{
		{
			Name: "config-1",
			Metadata: map[string]interface{}{
				"monitoring": map[string]interface{}{
					"enabled": "true",
				},
			},
		},
	}
	key := "metadata.monitoring.enabled"
	value := "false"
	expected := ConfigList{}

	d, _ := QueryConfig(key, value)

	assert.Equal(t, expected, d)
}
