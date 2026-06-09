package configs

type Server struct {
	Mysql Mysql     `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis Redis     `mapstructure:"redis" json:"redis" yaml:"redis"`
	Kafka Kafka     `mapstructure:"kafka" json:"kafka" yaml:"kafka"`
	Iot   IotConfig `mapstructure:"iot" json:"iot" yaml:"iot"`
}
