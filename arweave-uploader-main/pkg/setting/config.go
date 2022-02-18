package setting

type LogFile struct {
	LogPath    string `mapstructure:"log_path"`
	NamePrefix string `mapstructure:"name_prefix"`
}

type LoggerConfig struct {
	EnableContext bool    `mapstructure:"enable_context"`
	Level         string  `mapstructure:"level"`
	Mode          string  `mapstructure:"mode"` // "stdout" or "file"
	File          LogFile `mapstructure:"file"`
}

type ImageConfig struct {
	MaxFileSize int64 `mapstructure:"maxfilesize"`
}

type AppConfig struct {
	Node    string      `mapstructure:"node"`
	Keyfile string      `mapstructure:"keyfile"`
	Image   ImageConfig `mapstructure:"image"`
}

type ServerConfig struct {
	RunMode string `mapstructure:"mode"`
	Port    int    `mapstructure:"port"`
	Timeout struct {
		Read  int `mapstructure:"read"`
		Write int `mapstructure:"write"`
	} `mapstructure:"timeout"`
}

type Config struct {
	App    AppConfig    `mapstructure:"app"`
	Server ServerConfig `mapstructure:"server"`
	Logger LoggerConfig `mapstructure:"logger"`
}
