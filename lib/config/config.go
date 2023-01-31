package config

type Config struct {
	Meteo MeteoConf `envPrefix:"API_METEO_CONCEPT_"`
	Slack SlackConf `envPrefix:"SLACK_"`
}

type MeteoConf struct {
	Token    string `env:"TOKEN"`
	Tries    int    `env:"MAX_TRIES" envDefault:"3"`
	Location string `env:"INSEE_LOCATION" envDefault:"75056"`
}

type SlackConf struct {
	APIURL  string `env:"API_URL"`
	Token   string `env:"BOT_TOKEN"`
	Channel string `env:"CHANNEL"`
}
