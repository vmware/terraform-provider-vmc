package vmc

import (
	"context"
	"fmt"
	"time"

	"gitlab.eng.vmware.com/het/vmc-go-sdk/vmc"
)

func waitForTask(client *vmc.APIClient, orgID string, taskID string) error {
	for {
		task, _, err := client.TaskApi.OrgsOrgTasksTaskGet(context.Background(), orgID, taskID)

		if err != nil {
			return fmt.Errorf("Error while getting task %s: %v", taskID, err)
		}

		if task.Status == "STARTED" || task.Status == "CANCELING" {
			waitInterval := 2 * time.Second
			fmt.Printf("Wait for task result for another %s", waitInterval)
			time.Sleep(waitInterval)
			continue
		}

		return nil
	}
}
