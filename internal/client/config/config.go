package config

import "time"

type Config struct {
	Host                string `valid:"minstringlength(3)"`
	LogLevel            string `valid:"in(debug|info|warn|error|fatal)"`
	LogPath             string
	DBPath              string
	BuildVersion        string
	BuildDate           string
	ServerKeyPath       string
	FileStoragePath     string
	QueueWorkersCnt     int
	QueueTimeout        time.Duration
	ShowMsgTimeout      time.Duration
	ShutdownTimeout     time.Duration
	SyncInterval        time.Duration
	ViewRefreshInterval time.Duration
}
