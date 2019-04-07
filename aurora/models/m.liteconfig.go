package models

import (
	"errors"
	"io/ioutil"
	"regexp"
	"strings"
)

type LiteConfig struct {
	Config map[string]map[string]string
}

func NewConfig(path string) (*LiteConfig, error) {
	lc := &LiteConfig{}
	err := lc.EasyParse(path)
	return lc, err
}

func (config *LiteConfig) EasyParse(path string) (err error) {
	config.Config = make(map[string]map[string]string)

	if path == "" {
		err = errors.New("LiteConfig - Parser - Invalid filepath received")
		return
	}
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	lines := regexp.MustCompile("\n").Split(string(raw), -1)
	curSection := ""
	for _, l := range lines {
		if len(l) > 0 {
			if l[0] == 91 && l[len(l)-1] == 93 {
				curSection = l
				curSection = strings.Replace(curSection, "[", "", 1)
				curSection = strings.Replace(curSection, "]", "", 1)
				config.Config[curSection] = make(map[string]string)
				continue
			}
			kv := strings.Split(l, "=")
			if len(kv) > 1 {
				config.Config[string(curSection)][kv[0]] = kv[1]
			} else {
				config.Config[string(curSection)][kv[0]] = ""
			}
		}
	}
	return
}
