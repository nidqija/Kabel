package Command

import (
	"fmt"
	cobra "github.com/spf13/cobra"

)



// used to describe the marketplace command and its functionality, 
// providing users with an overview of what they can expect when using this command. 
// It serves as a guide for users to understand the purpose and benefits of the marketplace command, and 
// how it can enhance their experience with the application.

var marketplaceCMD = &cobra.Command{
	Use : "marketplace",
	Short : "Access the Kabel marketplace to discover and manage applications and services",
	Long : "The marketplace command provides users with access to the Kabel marketplace, where they can discover and manage applications and services that are compatible with the Kabel platform. Users can browse through a wide range of applications, view details about each application, and easily install or manage them directly from the command line interface.",


	Run : func (cmd *cobra.Command , args[] string){
	fmt.Println("Welcome to the Kabel Marketplace! Here you can discover and manage applications and services compatible with the Kabel platform.")
	},


}



// this function is used to initialize the marketplace command and add it to the root command.
func init(){
	rootCMD.AddCommand(marketplaceCMD)
}


