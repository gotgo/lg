package main

import (
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/hcl"
)

type LogConfig struct {
	TODO string `hcl:"to_do"`
}

func LoadConfig(path string) (*LogConfig, error) {

	d, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf(
			"Error reading %s: %s", path, err)
	}

	obj, err := hcl.Parse(string(d))
	if err != nil {
		return nil, fmt.Errorf(
			"Error parsing %s: %s", path, err)
	}

	// Build up the result
	var result LogConfig
	if err := hcl.DecodeObject(&result, obj); err != nil {
		return nil, err
	}

	return &result, nil
}
