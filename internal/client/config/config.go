package config

type Config struct {
	Host         string `valid:"minstringlength(3)"`
	LogLevel     string `valid:"in(debug|info|warn|error|fatal)"`
	LogPath      string
	BuildVersion string
	BuildDate    string
}
