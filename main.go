package main

import (
	"log"

	"github.com/fx2y/reflexion-go/reflexion"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// Set up the Temporal client
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	// Create the Reflexion Agent
	searchService := &MockSearchService{}
	agent := reflexion.NewReflexionAgent(searchService)

	// Create a Temporal worker
	w := worker.New(c, "reflexion-task-queue", worker.Options{})

	// Register the Reflexion Agent workflow and activities
	agent.RegisterWorkflowAndActivities(w)

	// Start the worker
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
