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
const defaultFileName = ".snykctl.yaml"
const defaultUrl = "https://snyk.io/api/v1"

type ConfigProperties struct {
	token      string
	id         string
	timeout    int
	workerSize int
	url        string
	sync       bool
	filename   string
}

func (c *ConfigProperties) SetSync(b bool) {
	c.sync = b
}

func (c *ConfigProperties) Init(sync bool) error {
	c.filename = defaultFileName
	c.timeout = defaultTimeout
	c.workerSize = defaultWorkerSize
	c.url = defaultUrl
	if sync {
		return c.Sync()
	}
	return nil
}

func (c *ConfigProperties) Sync() error {
	if !c.sync {
		return c.ReadConf()
	}
	return nil
}

func (c *ConfigProperties) SetFilename(file string) {
	c.filename = file
}

func (c ConfigProperties) Filename() string {
	return c.filename
}

func (c ConfigProperties) Url() string {
	return c.url
}

func (c ConfigProperties) Token() string {
	return c.token
}

func (c ConfigProperties) Id() string {
	return c.id
}

func (c ConfigProperties) ObfuscatedToken() string {
	if len(c.token) > 6 {
		return c.token[len(c.token)-6:]
	}

	return ""
}

func (c ConfigProperties) ObfuscatedId() string {
	if len(c.id) > 6 {
		return c.id[len(c.id)-6:]
	}

	return ""
}

func (c ConfigProperties) Timeout() int {
	return c.timeout
}

func (c ConfigProperties) WorkerSize() int {
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

func (c *ConfigProperties) WriteConf() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("writeConf failed: %s", err)
	}
	filename := home + "/.snykctl.yaml"
	confStr := fmt.Sprintf("token: %s\nid: %s\ntimeout: %d\nworkerSize: %d\n", c.token, c.id, c.timeout, c.workerSize)
	d1 := []byte(confStr)
	err = ioutil.WriteFile(filename, d1, 0644)
	if err != nil {
		return fmt.Errorf("writeConf failed: %s", err)
	}
	c.sync = true
	return nil
}

func (c *ConfigProperties) ReadConf() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("readConf failed: %s", err)
	}

	filepath := home + "/" + c.filename
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("readConf failed: %s", err)
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
						return fmt.Errorf("readConf failed: %s", err)
					}
					c.timeout = t
				}

				if key == "workerSize" {
					w, err := strconv.Atoi(value)
					if err != nil {
						return fmt.Errorf("readConf failed: %s", err)
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
		return fmt.Errorf("readConf failed: %s", err)
	}

	c.sync = true
	return nil
}

var Instance ConfigProperties
