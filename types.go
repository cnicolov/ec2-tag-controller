package main

type Tag struct {
	Key   string
	Value string
}

type TagMapping struct {
	Key      string
	Value    string
	TagKey   string `mapstructure:"tag_key"`
	TagValue string `mapstructure:"tag_value"`
}

type Config struct {
	Mappings []TagMapping `mapstructure:"mapping"`
}
