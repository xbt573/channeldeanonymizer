package cmd

type Config struct {
	Token  string `mapstructure:"token"`
	ChatId int64  `mapstructure:"chat_id"`
}
