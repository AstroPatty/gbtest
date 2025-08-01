/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/astropatty/gbtest/auth"
	"github.com/astropatty/gbtest/stack"
)


// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the initial infrastructure",
	Long: `Deploy the initial infrastructure to AWS. This command only needs to be run once, and will error if it is run a second time.`,

	Run: func(cmd *cobra.Command, args []string) {
		err := auth.CheckCredentials()
		if err != nil {
			msg := fmt.Sprintf("Error while trying to authenticate: %s", err)
			fmt.Println(msg)
			return
		}
		fmt.Println("Permissions valid!")
		stack.DeployStack()
		fmt.Println("Synthesized stack")
	},
}

var destroyCommand = &cobra.Command {
	Use: "destroy",
	Short: "Remove all infrastructure in AWS",
	Long: "Remove all infrastructure in AWS. This will delete ALL DATA",
	Run: func(cmd *cobra.Command, args[]string) {
		err := stack.DestroyStack()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Sucessfully removed resources")
		}
	},
	
}

func init() {
	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(destroyCommand)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
