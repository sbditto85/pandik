package main

import (
	"time"
)

type MonitorConf struct {
	Type string
	Url  string
	Freq string
	Data map[string]string
}

type MonitorLog struct {
	Up      bool
	Time    time.Time
	Message string
	Monitor *Monitor
}

type Monitor struct {
	Conf    *MonitorConf
	Checker Checker
	Up      bool
	Logs    []*MonitorLog
}

func NewMonitorLog(up bool, message string) *MonitorLog {
	return &MonitorLog{up, time.Now(), message, nil}
}

func NewMonitor(conf *MonitorConf) (*Monitor, error) {
	checker, err := GetChecker(conf.Type)
	if err != nil {
		return nil, err
	}

	return &Monitor{conf, checker, false, nil}, nil
}

func (m *Monitor) Watch(logChan chan *MonitorLog) {
	for {
		monitorLog := m.Checker(m.Conf)
		monitorLog.Monitor = m

		logChan <- monitorLog

		m.Logs = append(m.Logs, monitorLog)
		m.Up = monitorLog.Up

		nextCheck, _ := time.ParseDuration(m.Conf.Freq)
		time.Sleep(nextCheck)
	}
}
