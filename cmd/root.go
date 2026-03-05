package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xbt573/channeldeanonymizer/internal/bot"
)

var (
	config Config
)

func init() {
	rootCmd.PersistentFlags().StringP("token", "t", "", "Telegram Bot API token (also available as TOKEN env variable)")
	rootCmd.PersistentFlags().Int64P("chatid", "i", 0, "Channel ID (also available as CHAT_ID env variable)")

	rootCmd.MarkFlagRequired("token")
	rootCmd.MarkFlagRequired("chatid")

	viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("chat_id", rootCmd.PersistentFlags().Lookup("chatid"))

	viper.BindEnv("token", "TOKEN")
	viper.BindEnv("chat_id", "CHAT_ID")

	cobra.OnInitialize(func() {
		if err := viper.Unmarshal(&config); err != nil {
			slog.Error("failed to unmarshal config", "err", err)
			os.Exit(1)
		}
	})
}

var rootCmd = &cobra.Command{
	Use:   "channeldeanonymizer",
	Short: "Telegram bot to delete anonymous messages in channel",
	Args: func(cmd *cobra.Command, args []string) error {
		// ну пиздец, мне ещё руками это валидировать надо?
		if config.Token == "" {
			return fmt.Errorf("invalid token provided: %v", config.Token)
		}

		if config.ChatId == 0 {
			return fmt.Errorf("invalid chat id provided: %v", config.ChatId)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		botInstance := bot.New(bot.Options{
			Token:  config.Token,
			ChatId: config.ChatId,
		})

		return botInstance.Start()
	},
}

func Execute() {
	rootCmd.Execute()
}
