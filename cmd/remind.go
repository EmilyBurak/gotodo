package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/spf13/cobra"
)

var remindCmd = &cobra.Command{
	Use:   "remind",
	Short: "Invoke a lambda function to remind you about a task",
	Long:  `Invoke a lambda function to remind you about a task`,
	Run: func(cmd *cobra.Command, args []string) {

		task := cmd.Flag("task").Value.String()
		message := cmd.Flag("message").Value.String()
		// Create a new Lambda client
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			fmt.Println("Error loading configuration:", err)
			return
		}

		client := lambda.NewFromConfig(cfg)

		// Define the payload
		payload, err := json.Marshal(map[string]string{
			"task":    task,
			"message": message,
		})

		if err != nil {
			fmt.Println("Error marshaling payload:", err)
			return
		}

		// Invoke the Lambda function
		_, err = client.Invoke(context.TODO(), &lambda.InvokeInput{
			FunctionName: aws.String("goscheduler"),
			Payload:      payload,
		})
		if err != nil {
			fmt.Println("Error invoking lambda function:", err)
			return
		}

		fmt.Println("Message sent to lambda function:" + string(payload))
	},
}

func init() {
	rootCmd.AddCommand(remindCmd)
	remindCmd.Flags().StringP("task", "t", "", "Task to remind you about")
	remindCmd.MarkFlagRequired("task")
	remindCmd.Flags().StringP("message", "m", "", "Message to send to the lambda function")

}
