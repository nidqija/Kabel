package Strategy

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)


// create a new struct that implements the databasestrategy interface for pocketbase
type PocketBaseStrategy struct {
	containerID string
	cfg         DatabaseStrategy
	httpClient  *http.Client
}


// function to deploy a new pocketbase instance using the docker client
func (p *PocketBaseStrategy) Deploy(ctx context.Context, cli *client.Client, cfg DatabaseStrategy) (string, error) {

	// pull the pocketbase community image
	p.cfg = cfg
	targetImage := "ghcr.io/muchobien/pocketbase:latest"

	fmt.Println("Pulling PocketBase image...", targetImage)

	// pull the image
	reader, err := cli.ImagePull(ctx, targetImage, image.PullOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to pull image: %w", err)
	}

	// ensure the reader is closed after we're done
	// defer is used to guarantee that the reader is closed 
	// even if an error occurs during the pull process
	defer reader.Close()

	// read the output to ensure the pull completes before proceeding
	_, _ = io.Copy(io.Discard, reader)
	fmt.Println("PocketBase image pulled successfully")

	// serve the container on 8090 port and bind the data directory
	internalPort := nat.Port("8090/tcp")

	// define the port mapping for the container
	portMap := nat.PortMap{
		internalPort: []nat.PortBinding{
			{
				HostIP:   "127.0.0.1",
				HostPort: p.cfg.HostPort,
			},
		},
	}

	// bind the local data directory to the container's expected data directory
	volumeBinding := fmt.Sprintf("%s:/pb_data", p.cfg.LocalDatabaseDir)

	// create the container for pocketbase with the specified configuration
	resp, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image:        targetImage,
			ExposedPorts: nat.PortSet{internalPort: struct{}{}},
		},
		&container.HostConfig{
			PortBindings: portMap,
			Binds:        []string{volumeBinding},
		},
		nil, nil, p.cfg.InstanceName,
	)

	// log the error if container creation fails and return it to the caller
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}
    
	// start the container and log any errors that occur during startup
	err = cli.ContainerStart(ctx, resp.ID, container.StartOptions{})

	
	if err != nil {
		return "", fmt.Errorf("failed to start container: %w", err)
	}

	// store the container ID in struct for later reference
	p.containerID = resp.ID

	// log the successful start of the container with its ID
	fmt.Println("PocketBase container started successfully with ID:", p.containerID)

	// return the container ID to the caller for reference
	return p.containerID, nil
}


// function to connect to the pocketbase instance and check if it's healthy
func (p *PocketBaseStrategy) Connect(ctx context.Context, cli *client.Client, cfg DatabaseStrategy) error {
	p.httpClient = &http.Client{Timeout: 5 * time.Second}

	// construct the health check url using the host port from the configuration
	healthURL := fmt.Sprintf("http://127.0.0.1:%s/api/health", p.cfg.HostPort)
	fmt.Println("Checking PocketBase health at:", healthURL)

	var healthErr error
    // retry the health check multiple times with a delay to allow the container to start up and become healthy
	for i := 0; i < 10; i++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, healthURL, nil)
		if err != nil {
			healthErr = err
			fmt.Println("Error creating health check request:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		// if the request is successful, check the status code 
		// to determine if the service is healthy

		resp, err := p.httpClient.Do(req)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				fmt.Println("PocketBase is healthy and ready to accept connections.")
				return nil
			}

			// if the status code is not 200, consider it an error 
			// and log the unexpected status
			healthErr = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		} else {

			// if there was an error on the request
			// log the error and store it for potential return after retries are exhausted
			healthErr = err
			fmt.Println("Error performing health check retry:", err)
		}

		// add time buffer between retries to allow the container to finish starting up
		time.Sleep(2 * time.Second)
	}

	// if we exhaust all retries and still have an error, return it to the caller
	if healthErr != nil {
		return fmt.Errorf("failed to connect to PocketBase after multiple attempts: %v", healthErr)
	}

	// log the successful connection to the PocketBase instance
	fmt.Println("Connected to PocketBase successfully!")
	return nil
}


// function to clean up resources when closing the connection
func (p *PocketBaseStrategy) Close(ctx context.Context) error {
	if p.httpClient != nil {
		p.httpClient.CloseIdleConnections()
		fmt.Println("PocketBase client wrapper connections idle-closed.")
	}
	return nil
}