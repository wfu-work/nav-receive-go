package configs

type Mysql struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`
	Username string `mapstructure:"username" json:"username" yaml:"username"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	DBName   string `mapstructure:"dbname" json:"dbname" yaml:"dbname"`
	Charset  string `mapstructure:"charset" json:"charset" yaml:"charset"`
}
