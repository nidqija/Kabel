package command

// This file contains the code for the "root" command, which allows users to interact with the application. It serves as the entry point for all other commands and provides a common interface for users to access the various features of the application.
// pass for now, will implement later.

import (
	"fmt"
	"os"
	cobra "github.com/spf13/cobra"
)


var (
	VerboseMode bool
)

// root command is the main command that serves as the entry point for all other commands
// provides a common interface for users to access the various feature of the application
var rootCMD = &cobra.Command{
	
	Use : "Kabel", 
	Short: "Kabel is a high performance local private Database as a Service ",
    Long: "An automated local platform infrastructure management tool that simplifies the deployment and management of local applications and services, providing a seamless experience for developers and users alike.",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Kabel - Your Local Database as a Service!")
		fmt.Println("Use 'kabel --help' to see available commands and options.")
	},

}

func Execute(){

	if err := rootCMD.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}


func init(){
	rootCMD.PersistentFlags().BoolVarP(&VerboseMode, "verbose" , "v"  , false , "Enable verbose output for debugging purposes")
}


