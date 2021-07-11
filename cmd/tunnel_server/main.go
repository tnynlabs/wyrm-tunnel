package main

import (
	"log"
	"os"
	"sync"

	"github.com/tnynlabs/wyrm-tunnel/pkg/manager"
	"github.com/tnynlabs/wyrm-tunnel/pkg/transport/tungrpc"
	"github.com/tnynlabs/wyrm-tunnel/pkg/tunnels"
	"github.com/tnynlabs/wyrm/pkg/devices"
	"github.com/tnynlabs/wyrm/pkg/storage/postgres"

	"github.com/joho/godotenv"
)

func main() {
	if devFlag := os.Getenv("WYRM_DEV"); devFlag == "1" {
		// Load environment variables from .env file
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Error loading .env file (error: %v)", err)
		}
	}

	db, err := postgres.GetFromEnv()
	if err != nil {
		log.Fatalln(err)
	}

	deviceRepo := postgres.CreateDeviceRepository(db)
	deviceService := devices.CreateDeviceService(deviceRepo)

	registry := tunnels.CreateRegistry()

	// Wait group to allow server to notify main when they exit
	// so we can cleanup before exiting
	wg := new(sync.WaitGroup)
	wg.Add(1)

	managerServer := manager.NewServer(registry)
	go func() {
		log.Println("Running tunnel manager server...")
		err := manager.RunServer(":5050", managerServer)
		if err != nil {
			log.Printf("Failed to run tunnel manager server (%v)\n", err)
		}
		wg.Done()
	}()

	grpcTunnelServer := tungrpc.NewServer(registry, deviceService)
	go func() {
		log.Println("Running grpc device tunnel transport server...")
		err := tungrpc.RunServer(":5051", grpcTunnelServer)
		if err != nil {
			log.Printf("Failed to run grpc device tunnel transport server (%v)\n", err)
		}
		wg.Done()
	}()

	// If any of the running servers cleanup and exit
	wg.Wait()
	// TODO: Cleanup
}
