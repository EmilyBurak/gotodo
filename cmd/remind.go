package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	"github.com/aws/aws-sdk-go-v2/service/scheduler/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var remindCmd = &cobra.Command{
	Use:   "remind",
	Short: "Invoke a lambda function to remind you about a task",
	Long:  `Invoke a lambda function to remind you about a task`,
	Run: func(cmd *cobra.Command, args []string) {

		viper.SetConfigType("env")
		viper.SetConfigName(".env")
		viper.AddConfigPath(".")
		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}

		task := cmd.Flag("task").Value.String()
		message := cmd.Flag("message").Value.String()

		hour := cmd.Flag("hour").Value.String()
		minute := cmd.Flag("minute").Value.String()

		today := time.Now()
		formattedToday := today.Format("2006-01-02T")

		log.Printf("Task: %s, Message: %s, Hour: %s, Minute: %s", task, message, hour, minute)

		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			fmt.Println("Error loading configuration:", err)
			return
		}

		client := scheduler.NewFromConfig(cfg)

		_, err = client.CreateSchedule(context.TODO(), &scheduler.CreateScheduleInput{
			FlexibleTimeWindow: &types.FlexibleTimeWindow{Mode: types.FlexibleTimeWindowModeOff},
			Name:               aws.String("reminderfor" + task),
			// yyyy-mm-ddThh:mm:ss
			ScheduleExpression: aws.String("at(" + formattedToday + hour + ":" + minute + ":00)"),
			Target: &types.Target{
				Arn:     aws.String(viper.Get("SNS_TOPIC_ARN").(string)),
				RoleArn: aws.String(viper.Get("ROLE_ARN").(string)),
				Input:   aws.String(fmt.Sprintf(`Task: %s, Message: %s`, task, message)),
			},
			ScheduleExpressionTimezone: aws.String("America/Denver"),
			State:                      types.ScheduleStateEnabled,
		})
		if err != nil {
			fmt.Println("Error creating schedule:", err)
			return
		}
		fmt.Println("Reminder set for task: " + task)
	},
}

func init() {
	rootCmd.AddCommand(remindCmd)
	remindCmd.Flags().StringP("task", "t", "", "Task to remind you about")
	remindCmd.MarkFlagRequired("task")
	remindCmd.Flags().StringP("hour", "H", "12", "Hour to remind you about the task in 24-hour format")
	remindCmd.MarkFlagRequired("hour")
	remindCmd.Flags().StringP("minute", "M", "0", "Minute of hour to remind you about the task")
	remindCmd.MarkFlagRequired("minute")
	remindCmd.Flags().StringP("message", "m", "", "Message to send to AWS SNS")
}
