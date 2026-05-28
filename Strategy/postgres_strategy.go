package Strategy

// implement later

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq"
)


type PostgresStrategy struct {
	containerID string;
	db *sql.DB;
	cfg DatabaseStrategy;


}


func (p *PostgresStrategy) Deploy(ctx context.Context, cli *client.Client , cfg DatabaseStrategy) (string , error) {
        
	// pull the postgres image
	p.cfg = cfg;
	targetImage := "postgres:16-alpine";

	fmt.Println("Pulling Postgres image..." , targetImage);

	reader, err := cli.ImagePull(ctx, targetImage, image.PullOptions{})

	if err != nil {
		return "", err;
	}

	defer reader.Close()

	_, _ = io.Copy(io.Discard, reader)

	fmt.Println("Postgres image pulled successfully")

	internalPort := nat.Port("5432/tcp")

	portMap := nat.PortMap{
		internalPort: []nat.PortBinding{
			{
				HostIP: "127.0.0.1",
				HostPort: p.cfg.HostPort,
			},
		},
	}


	volumeBinding := fmt.Sprintf("%s:/var/lib/postgresql/data", p.cfg.LocalDatabaseDir)

	postgresEnv := []string{
		fmt.Sprintf("POSTGRES_DB=%s", p.cfg.LocalDatabaseName),
		fmt.Sprintf("POSTGRES_USER=%s", p.cfg.Username),
		fmt.Sprintf("POSTGRES_PASSWORD=%s", p.cfg.Password),
	}

	resp , err := cli.ContainerCreate(ctx , 
		&container.Config{
		Image : targetImage,
		Env: postgresEnv,
		ExposedPorts: nat.PortSet{internalPort: struct{}{}},
	},

	&container.HostConfig{
		PortBindings: portMap,
		Binds: []string{volumeBinding},
	},

	nil , nil , p.cfg.InstanceName,)

	if err != nil { return "", err;}

	err = cli.ContainerStart(ctx, resp.ID , container.StartOptions{})

	if err != nil { return "", err;}

	p.containerID = resp.ID;

	fmt.Println("Postgres container started successfully with ID: " , p.containerID);

	return resp.ID, nil

}

func (p *PostgresStrategy) Connect(ctx context.Context, cli *client.Client , cfg DatabaseStrategy) (error) {
	
	connectionString := fmt.Sprintf("host=localhost port=%s user=%s password=%s dbname=%s sslmode=disable", p.cfg.HostPort, p.cfg.Username, p.cfg.Password, p.cfg.LocalDatabaseName)

	db , err := sql.Open("postgres", connectionString)

	if err != nil {
		return err;
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(5 * time.Minute)

	fmt.Println("Probing Database socket inside docker environment...")

	var pingErr error

	for i := 0; i < 10; i++ {
		pingErr = db.PingContext(ctx)

		if pingErr == nil {
			break;
		}

		time.Sleep(2 * time.Second)
	}

	if pingErr != nil {
		db.Close();
		return fmt.Errorf("Failed to connect to the database after multiple attempts: %v", pingErr);
	}

	p.db = db;

	fmt.Println("Connected to Postgres database successfully!");

	return nil;
	
	
}


func (p *PostgresStrategy) Close(ctx context.Context) error {
	if p.db != nil {
		fmt.Print("Closing Database Pool...");
		err := p.db.Close();
		if err != nil {
			return err;
		}
	}
return nil;

}


