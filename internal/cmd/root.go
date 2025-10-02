/*
Copyright © 2025 Peter Shaan <petershaan12@gmail.com>
*/
package cmd

import (
	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-auth-clean-arch",
	Short: "An application for user authentication using clean architecture principles by Peter Shaan",
	Long:  "A longer description that spans multiple lines and likely contains examples and usage of using your application.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
