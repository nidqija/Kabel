package Command

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"kabel/Engine"
	"kabel/Strategy"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var (
	dbUsername string;
	dbPassword string;
	engineType string;
)


var InstanceCMD = &cobra.Command{
	Use:   "instance",
	Short: "Manage local database workspace instances",

}

var InsertInstanceCMD = &cobra.Command{
	Use:   "insert",
	Short: "Insert a new database instance",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		runtimePort , err := Engine.FindAvailablePort()

		if err != nil {
			fmt.Println("Error finding available port: ", err);
			return err;
		}


		var targetStrategy Strategy.DatabaseStrategyInterface;

		switch engineType {
		 case "postgres":
			targetStrategy = &Strategy.PostgresStrategy{}

		 default:
			return fmt.Errorf("unsupported engine type: %s", engineType)
		}

		absDataDir, err := filepath.Abs(filepath.Join("data", args[0]))

		if err != nil {
			fmt.Println("Error getting absolute path: ", err)
			return err;
		}

		_ = os.MkdirAll(absDataDir, 0755)

		cfg := Strategy.DatabaseStrategy{
			InstanceName: args[0],
			HostPort: runtimePort,
			LocalDatabaseName: args[0],
			LocalDatabaseDir: absDataDir,
			Username: dbUsername,
			Password: dbPassword,
		}

		cli , err := client.NewClientWithOpts(client.FromEnv , client.WithAPIVersionNegotiation())

		if err != nil {
			fmt.Println("Error creating Docker client: ", err)
			return err;
		}

		defer cli.Close();

		fmt.Printf("🚀 Kabel: Provisioning your isolated '%s' deployment...\n", args[0])

		containerID , err := targetStrategy.Deploy(ctx, cli , cfg)

		if err != nil {
			fmt.Println("Error deploying database instance: ", err)
			return err;
		}

		err = targetStrategy.Connect(ctx, cli , cfg)

		if err != nil {
			fmt.Println("Error connecting to database instance: ", err)
			return err;
		}

		defer targetStrategy.Close(ctx)


		fmt.Printf("\n✅ Infrastructure instance fully initialized!\n")
		fmt.Printf("📌 Local Endpoint:  127.0.0.1:%s\n", runtimePort)
		fmt.Printf("👤 Root Username:   %s\n", dbUsername)
		fmt.Printf("🔑 Master Password: %s\n", dbPassword)
		fmt.Printf("🆔 Docker Ref ID:   %s\n", containerID[:12])


        return nil
		

	},
}


func init() {
	InstanceCMD.AddCommand(InsertInstanceCMD)

	InsertInstanceCMD.Flags().StringVarP(&dbUsername, "username", "u", "admin", "Database root username")
	InsertInstanceCMD.Flags().StringVarP(&dbPassword, "password", "p", "admin", "Database root password")
	InsertInstanceCMD.Flags().StringVarP(&engineType, "engine", "e", "postgres", "Database engine type (e.g., postgres)")
}


