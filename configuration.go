package configuration

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/vipinmakode/configuration/hocon"
)

func ParseString(text, baseDir string, includeCallback ...hocon.IncludeCallback) *Config {
	var callback hocon.IncludeCallback
	if len(includeCallback) > 0 {
		callback = includeCallback[0]
	} else {
		callback = defaultIncludeCallback
	}
	root := hocon.Parse(text, baseDir, callback)
	return NewConfigFromRoot(root)
}

func LoadConfig(filename string) *Config {
	data, baseDir, err := readFile(filename)
	if err != nil {
		panic(err)
	}

	return ParseString(string(data), baseDir, defaultIncludeCallback)
}

func FromObject(obj interface{}) *Config {
	data, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}

	return ParseString(string(data), "", defaultIncludeCallback)
}

func defaultIncludeCallback(filename string) *hocon.HoconRoot {
	data, baseDir, err := readFile(filename)
	if err != nil {
		// Ignore the missing included file
		return hocon.Parse("", "", defaultIncludeCallback)
	}

	return hocon.Parse(string(data), baseDir, defaultIncludeCallback)
}

func readFile(filename string) (string, string, error) {
	if strings.HasPrefix(filename, "http://") || strings.HasPrefix(filename, "https://") || strings.HasPrefix(filename, "file://") {
		panic("url is not yet supported.")
	} else {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return "", "", err
		}
		return string(data), filepath.Dir(filename), err
	}
}
