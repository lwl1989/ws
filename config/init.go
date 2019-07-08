package config

import (
    "encoding/json"
    "io/ioutil"
)

type Config struct {
    LogConfig *LogFileConfig `json:"log_config"`
    Redis *Redis `json:"redis"`
}

type LogFileConfig struct {
    FilePath string `json:"file_path,omitempty"` // default /tmp/{YmdHis}.log
}

type Redis struct{
    Host string `json:"host,omitempty"`
    Db   string `json:"db,omitempty"`
    Pw   string `json:"pw,omitempty"`
}

var Cf *Config

func init()   {
    Cf = &Config{
        LogConfig: &LogFileConfig{
            FilePath: "/tmp/ws_log",
        },
        Redis: &Redis{
            Host:"127.0.0.1:6379",
            Db:"1",
            Pw:"",
        },
    }
}

func (cf *Config) LoadConfig(path string) {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        data, err = ioutil.ReadFile("config.json")
        if err != nil {
            panic("load config err " + path)
        }
    }
    b := []byte(data)
    err = json.Unmarshal(b, Cf)
    if err != nil {
        panic("load config err " + path)
    }
}


type LogConfigInterface interface {
    GetLogConfig() (string,interface{})
}

func (file *LogFileConfig) GetLogConfig() (string,interface{}) {
    return "file", file
}

func (file *LogFileConfig) GetFilePath() interface{} {
    return file.FilePath
}
