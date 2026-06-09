package configs

type Kafka struct {
	Name           string   `mapstructure:"name" json:"name" yaml:"name"`
	Servers        []string `mapstructure:"servers" json:"servers" yaml:"servers"`
	GroupID        string   `mapstructure:"group-id" json:"group-id" yaml:"group-id"`
	RtloggingTopic string   `mapstructure:"rtlogging-topic" json:"rtlogging-topic" yaml:"rtlogging-topic"`
	SensorTopic    string   `mapstructure:"sensor-topic" json:"sensor-topic" yaml:"sensor-topic"`
}
