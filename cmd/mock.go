package main

import (
	"awesomeProject/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	forceMock bool
)

var mockCmd = &cobra.Command{
	Use:   "mock",
	Short: "Generate mock data",
	Long:  `Generate mock data for development and testing`,
	Run: func(cmd *cobra.Command, args []string) {
		//if forceMock {
		logrus.Info("force generating mock data")
		utils.MockCourseData()
		//} else {
		//	logrus.Info("mock data generation skipped (use --force to enable)")
		//}
	},
}

func init() {
	//mockCmd.Flags().BoolVarP(&forceMock, "force", "f", false, "Force generate mock data")
	rootCmd.AddCommand(mockCmd)
}
