package models

type Config struct {
	Database struct {
		Connection string `yaml:"connection"`
	} `yaml:"database"`
}
