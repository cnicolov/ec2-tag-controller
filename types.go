package main

type Mapping struct {
	Key      string
	Value    string
	TagKey   string `mapstructure:"tag_key"`
	TagValue string `mapstructure:"tag_value"`
}

type Config struct {
	Mappings []Mapping `mapstructure:"mapping"`
}
