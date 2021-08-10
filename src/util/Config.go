package util

import (
	"io/ioutil"
	"strings"
)

type Config struct {
	file   string
	fields map[string]string
}

func LoadConfig(fileName string) (*Config, error) {
	var c Config
	c.file = fileName
	DebugMsg("Config", "Loading:"+fileName)
	c.fields = make(map[string]string)
	fileByte, err := ioutil.ReadFile(c.file)

	if err != nil {
		return nil, err
	}
	fileStr := strings.ReplaceAll(string(fileByte), "\r", "")
	lines := strings.Split(fileStr, "\n")
	for _, line := range lines {
		field := strings.Split(line, "=")
		if len(field) >= 2 {
			c.fields[field[0]] = field[1]
		}
	}
	return &c, nil
}
func (c *Config) Set(key string, value string) {
	c.fields[key] = value
}
func (c *Config) Get(key string) (string, bool) {
	value, ok := c.fields[key]
	return value, ok
}

func (c *Config) GetAnyway(key string, prefer string) string {
	value, ok := c.fields[key]
	if !ok {
		return prefer
	} else {
		return value
	}
}

func (c *Config) Remove(key string) bool {
	_, ok := c.fields[key]
	if !ok {
		return false
	}
	delete(c.fields, key)
	return true
}
func (c *Config) Sync() error {
	text := ""
	for key, value := range c.fields {
		text += key + "=" + value + "\n"
	}
	return ioutil.WriteFile(c.file, []byte(text), 0777)
}
