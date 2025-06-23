package main

import (
	"github.com/spf13/cobra"
	"github.com/tolikproh/otus_hw/hw12_13_14_15_16_calendar/cmd"
)

var configFile string

var rootCmd = &cobra.Command{
	Use:   "calendar_scheduler",
	Short: `Сервис "Планировщик календаря"`,
	Long: `Сервис "Планировщик календаря" периодически сканирует основную базу данных, 
выбирая события о которых нужно напомнить.
Сервис выполняет следующий функционал:
	- выбырает события для которых следует отправка уведомлений;
	- складывает события в очередь для отправки уведомлений;
	- очистку старых собитий`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	rootCmd.AddCommand(cmd.Version)
	rootCmd.Flags().StringVarP(&configFile, "config", "c", "./configs/scheduler_config.yaml", "path to configuration file")
}

func main() {
	rootCmd.Execute()
}
