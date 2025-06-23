package main

import (
	"github.com/spf13/cobra"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/cmd"
)

var configFile string

var rootCmd = &cobra.Command{
	Use:   "calendar",
	Short: `Сервис "Календарь"`,
	Long: `Сервис "Календарь" представляет собой сервис для хранения календарных событий и отправки уведомлений.
Сервис предоставляет возможность:
	- добавить/обновить событие;
	- получить список событий на день/неделю/месяц;
	- получить уведомление за несколько дней до события.`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	rootCmd.AddCommand(cmd.Version)
	rootCmd.Flags().StringVarP(&configFile, "config", "c", "./configs/calendar_config.yaml", "path to configuration file")
}

func main() {
	rootCmd.Execute()
}
