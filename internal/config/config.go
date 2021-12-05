package config

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const defaultTimeout = 10
const defaultWorkerSize = 10
const fileName = ".snykctl.yaml"
const defaultUrl = "https://snyk.io/api/v1"

type ConfigProperties struct {
	token      string
	id         string
	timeout    int
	workerSize int
	url        string
	sync       bool
}

func (c *ConfigProperties) SetSync(b bool) {
	c.sync = b
}

func (c *ConfigProperties) Sync() {
	if !c.sync {
		c.readConf()
	}
}

func (c ConfigProperties) Url() string {
	if !c.sync {
		c.readConf()
	}

	return c.url
}

func (c ConfigProperties) Token() string {
	c.Sync()
	return c.token
}

func (c ConfigProperties) Id() string {
	c.Sync()
	return c.id
}

func (c ConfigProperties) ObfuscatedToken() string {
	c.Sync()
	if len(c.token) > 6 {
		return c.token[len(c.token)-6:]
	}

	return ""
}

func (c ConfigProperties) ObfuscatedId() string {
	c.Sync()
	if len(c.id) > 6 {
		return c.id[len(c.id)-6:]
	}

	return ""
}

func (c ConfigProperties) Timeout() int {
	c.Sync()
	return c.timeout
}

func (c ConfigProperties) WorkerSize() int {
	c.Sync()
	return c.workerSize
}

func (c *ConfigProperties) SetToken(t string) {
	c.token = t
}

func (c *ConfigProperties) SetId(i string) {
	c.id = i
}

func (c *ConfigProperties) SetTimeout(t int) {
	c.timeout = t
}

func (c *ConfigProperties) SetTimeoutStr(t string) {
	tt, err := strconv.Atoi(t)
	if err != nil {
		panic(err)
	}
	c.timeout = tt
}

func (c *ConfigProperties) SetWorkerSize(w int) {
	c.workerSize = w
}

func (c *ConfigProperties) SetWorkerSizeStr(t string) {
	tt, err := strconv.Atoi(t)
	if err != nil {
		panic(err)
	}
	c.workerSize = tt
}

func (c *ConfigProperties) SetUrl(u string) {
	c.url = u
}

func (c *ConfigProperties) WriteConf() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	filename := home + "/.snykctl.yaml"
	confStr := fmt.Sprintf("token: %s\nid: %s\ntimeout: %d\nworkerSize: %d\n", c.token, c.id, c.timeout, c.workerSize)
	d1 := []byte(confStr)
	err = ioutil.WriteFile(filename, d1, 0644)
	if err != nil {
		panic(err)
	}
	c.sync = true
}

func (c *ConfigProperties) readConf() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	filename := home + "/.snykctl.yaml"

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if equal := strings.Index(line, ":"); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				if key == "token" {
					c.token = value
				}

				if key == "id" {
					c.id = value
				}

				if key == "timeout" {
					t, err := strconv.Atoi(value)
					if err != nil {
						panic(err)
					}
					c.timeout = t
				}

				if key == "workerSize" {
					w, err := strconv.Atoi(value)
					if err != nil {
						panic(err)
					}
					c.workerSize = w
				}

				if key == "url" {
					c.url = value
				}
			}
		}
	}

	if c.url == "" {
		c.url = defaultUrl
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	c.sync = true
}

var Instance ConfigProperties
