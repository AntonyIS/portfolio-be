package config

type config struct {
	Env  string
	Port string
}

func Config(env string) *config {
	return &config{
		Env:  env,
		Port: "8080",
	}
}
