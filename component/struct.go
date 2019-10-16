package component

import (
    "io/ioutil"
    "encoding/json"
)

type Config struct {
    LogConfig *LogFileConfig `json:"log_config,omitempty"`
    Redis *Redis `json:"redis,omitempty"`
    Etcd  *EtcdConfig `json:"etcd,omitempty"`
}

type EtcdConfig struct {
    Enabled bool `json:"enabled,omitempty"`
    Addr []string `json:"addr,omitempty"`
    TimeOut int64 `json:"time_out,omitempty"`
    LeaseTime int64 `json:"lease_time,omitempty"`
}

type LogFileConfig struct {
    FilePath string `json:"file_path,omitempty"` // default /tmp/{YmdHis}.log
}

type Redis struct{
    Host string `json:"host,omitempty"`
    Db   string `json:"db,omitempty"`
    Pw   string `json:"pw,omitempty"`
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
