package config

type Config struct {
	Host          string `valid:"minstringlength(3)"`
	LogLevel      string `valid:"in(debug|info|warn|error|fatal)"`
	LogPath       string
	DBPath        string
	BuildVersion  string
	BuildDate     string
	ServerKeyPath string
}
