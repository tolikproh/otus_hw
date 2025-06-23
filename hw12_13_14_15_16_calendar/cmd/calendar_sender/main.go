package main

import (
	"github.com/spf13/cobra"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/cmd"
)

var configFile string

var rootCmd = &cobra.Command{
	Use:   "calendar_sender",
	Short: `Сервис "Рассыльщик"`,
	Long:  `Сервис "Рассыльщик" читает сообщения из очереди и шлёт уведомления.`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	rootCmd.AddCommand(cmd.Version)
	rootCmd.Flags().StringVarP(&configFile, "config", "c", "./configs/sender_config.yaml", "path to configuration file")
}

func main() {
	rootCmd.Execute()
}
