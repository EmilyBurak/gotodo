package cmd

import (
	"encoding/csv"
	"fmt"
	"os"

	tui "github.com/EmilyBurak/gotodo/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var listtuiCmd = &cobra.Command{
	Use:   "listtui",
	Short: "Display the list as a bubbletea tui",
	Long:  `Display the list of tasks as a bubbletea tui.`,
	Run: func(cmd *cobra.Command, args []string) {

		All, _ := cmd.Flags().GetBool("all")

		file, err := os.Open("tasks.csv")
		if err != nil {
			fmt.Println("No tasks to list")
			return
		}
		defer file.Close()

		csvReader := csv.NewReader(file)
		// Skip the headers
		_, err = csvReader.Read()
		if err != nil {
			panic(err)
		}

		// Loop through the records and add them to the tasks slice
		tasks := []tui.Task{}
		for {
			record, err := csvReader.Read()
			if err != nil {
				break
			}
			if !All && (record[2] == "done" || record[3] == "true") {
				continue
			}
			if record[5] == "" {
				record[5] = "N/A"
			}
			// Add the task to the tasks slice
			tasks = append(tasks, tui.Task{
				ID:              record[0],
				Name:            record[1],
				Status:          record[2],
				Deleted:         record[3],
				Pomodoros:       record[4],
				PomodorosNeeded: record[5],
			})
		}
		// Run the bubbletea program
		tea.NewProgram(tui.Model.InitList(tui.Model{}, 20, 20, tasks)).Run()
	},
}

func init() {
	rootCmd.AddCommand(listtuiCmd)
	listtuiCmd.Flags().BoolP("all", "a", false, "List all tasks")
}
