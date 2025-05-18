package main

import (
	"github.com/spf13/cobra"
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
	rootCmd.AddCommand(Version)
	rootCmd.Flags().StringVarP(&configFile, "config", "c", "./configs/config.yaml", "path to configuration file")
}

func main() {
	rootCmd.Execute()
}
