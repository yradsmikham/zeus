package util

import (
	"context"

	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/build"
	"github.com/microsoft/azure-devops-go-api/azuredevops/pipelines"
	log "github.com/sirupsen/logrus"
)

// Execute pipeline
func ExecPipeline(ctx context.Context, connection *azuredevops.Connection, app string, id int, project string) (err error) {
	log.Info("Queuing run for ", app, "...")

	client, err := build.NewClient(ctx, connection)
	if err != nil {
		log.Error("Unable to connect to Azure DevOps.")
		log.Fatal(err)
	}
	
	/*
		branch := "refs/heads/master"
		repoName := "OHDISonAzure"
		repoUrl := "https://github.com/yradsmikham/OHDSIonAzure"
		buildRepository := build.BuildRepository{
			DefaultBranch: &branch,
			Name:          &repoName,
			Url:           &repoUrl,
		}

		//*buildRepository.DefaultBranch = "main"
	*/

	definition := build.DefinitionReference{
		Id: new(int),
	}
	*definition.Id = id

	buildParameters := build.Build{
		Definition: &definition,
		//Repository: &buildRepository,
	}

	queueBuildArgs := build.QueueBuildArgs{
		Build:   &buildParameters,
		Project: &project,
	}

	responseValue, err := client.QueueBuild(ctx, queueBuildArgs)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Info("For more information on the build, please visit: ", *responseValue.Url)
	}

	return err
}

func ReturnBuildId(app string, project string, ctx context.Context, connection *azuredevops.Connection) (id []int, err error) {
	log.Info("Retrieving Build ID for: ", app)
	client := pipelines.NewClient(ctx, connection)
	if err != nil {
		log.Error("Unable to connect to Azure DevOps.")
		log.Fatal(err)
	}

	listPipelinesArgs := pipelines.ListPipelinesArgs{
		Project: &project,
	}

	responseValue, err := client.ListPipelines(ctx, listPipelinesArgs)
	if err != nil {
		log.Fatal(err)
	}

	index := 0
	ids := []int{}
	for responseValue != nil {
		for _, pipes := range (*responseValue).Value {
			if *pipes.Name == app {
				log.Info("Build ID for ", app, ": ", *pipes.Id)
				ids = append(ids, *pipes.Id)
				return ids, err
			} else {
				index++
			}
		}
		log.Error("Build ID for ", app, " was not found. Please try again.")
	}
	return ids, err
}
