package Command

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"kabel/Engine"
	"kabel/Strategy"

	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var (
	dbUsername string
	dbPassword string
	engineType string
)

var InstanceCMD = &cobra.Command{
	Use:   "instance",
	Short: "Manage local database workspace instances",
}

var InsertInstanceCMD = &cobra.Command{
	Use:   "insert [instance_name]",
	Short: "Insert a new database instance",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		instanceName := args[0]

		runtimePort, err := Engine.FindAvailablePort()
		if err != nil {
			fmt.Println("Error finding available port: ", err)
			return err
		}

		// 1. Resolve strategy based on engineType flag
		var targetStrategy Strategy.DatabaseStrategyInterface

		switch engineType {
		case "postgres":
			targetStrategy = &Strategy.PostgresStrategy{}
		case "pocketbase":
			targetStrategy = &Strategy.PocketBaseStrategy{}
		default:
			return fmt.Errorf("unsupported engine type: %s", engineType)
		}

		absDataDir, err := filepath.Abs(filepath.Join("data", instanceName))
		if err != nil {
			fmt.Println("Error getting absolute path: ", err)
			return err
		}

		_ = os.MkdirAll(absDataDir, 0755)

		cfg := Strategy.DatabaseStrategy{
			InstanceName:      instanceName,
			HostPort:          runtimePort,
			LocalDatabaseName: instanceName,
			LocalDatabaseDir:  absDataDir,
			Username:          dbUsername,
			Password:          dbPassword,
		}

		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			fmt.Println("Error creating Docker client: ", err)
			return err
		}
		defer cli.Close()

		fmt.Printf("🚀 Kabel: Provisioning your isolated '%s' (%s) deployment...\n", instanceName, engineType)

		containerID, err := targetStrategy.Deploy(ctx, cli, cfg)
		if err != nil {
			fmt.Println("Error deploying database instance: ", err)
			return err
		}

		err = targetStrategy.Connect(ctx, cli, cfg)
		if err != nil {
			fmt.Println("Error connecting to database instance: ", err)
			return err
		}
		defer targetStrategy.Close(ctx)

		fmt.Printf("\n✅ Infrastructure instance fully initialized!\n")
		fmt.Printf("📌 Local Endpoint:  http://127.0.0.1:%s\n", runtimePort)
		if engineType == "postgres" {
			fmt.Printf("👤 Root Username:   %s\n", dbUsername)
			fmt.Printf("🔑 Master Password: %s\n", dbPassword)
		} else {
			fmt.Printf("🔧 Admin Dashboard: http://127.0.0.1:%s/_/\n", runtimePort)
		}
		fmt.Printf("🆔 Docker Ref ID:   %s\n", containerID[:12])

		// 2. Generate engine-specific environment configuration file content
		var envContent string
		currentTime := time.Now().Format("2006-01-02 15:04:05")

		if engineType == "postgres" {
			envContent = fmt.Sprintf(`# Generated automatically by Kabel on %s
				# Core connection parameters 
				DATABASE_HOST=127.0.0.1
				DATABASE_PORT=%s
				DATABASE_USER=%s
				DATABASE_PASSWORD=%s
				DATABASE_NAME=%s

				# Connection String (for convenience)
				DATABASE_URL=postgresql://%s:%s@127.0.0.1:%s/%s

				# Docker container reference
				DOCKER_CONTAINER_ID=%s
				`, currentTime, runtimePort, dbUsername, dbPassword, instanceName, dbUsername, dbPassword, runtimePort, instanceName, containerID[:12])
						} else if engineType == "pocketbase" {
							envContent = fmt.Sprintf(`# Generated automatically by Kabel on %s
				# PocketBase parameters
				POCKETBASE_URL=http://127.0.0.1:%s
				POCKETBASE_ADMIN_UI=http://127.0.0.1:%s/_/

				# Docker container reference
				DOCKER_CONTAINER_ID=%s
			`, currentTime, runtimePort, runtimePort, containerID[:12])

			// if the engine type is pocketbase , generate environment variables specific to pocketbase configuration
					} else if engineType == "pocketbase"{
			envContent = fmt.Sprintf(`# Generated automatically by Kabel on %s
				# PocketBase parameters
				POCKETBASE_URL=http://127.0.0.1:%s
				POCKETBASE_ADMIN_UI=http://127.0.0.1:%s/_/

				# Docker container reference
				DOCKER_CONTAINER_ID=%s
			`, currentTime, runtimePort, runtimePort, containerID[:12])
					}

		envFilePath := filepath.Join("data", instanceName, ".env")
		err = os.WriteFile(envFilePath, []byte(envContent), 0644)
		if err != nil {
			fmt.Println("Error writing .env file: ", err)
			return err
		}

		return nil
	},
}

func init() {
	InstanceCMD.AddCommand(InsertInstanceCMD)

	InsertInstanceCMD.Flags().StringVarP(&dbUsername, "username", "u", "admin", "Database root username")
	InsertInstanceCMD.Flags().StringVarP(&dbPassword, "password", "p", "admin", "Database root password")
	InsertInstanceCMD.Flags().StringVarP(&engineType, "engine", "e", "postgres", "Database engine type (postgres, pocketbase)")
}