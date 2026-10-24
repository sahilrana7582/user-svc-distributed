package config

type Config struct {
	App *Application
	Db  *Database
}

type Application struct {
	Name    string
	Port    string
	NatsUrl string
}

type Database struct {
	Host     string
	Port     string
	User     string
	Name     string
	Password string
	SSLMode  string
}
