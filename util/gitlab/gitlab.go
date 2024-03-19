package gitlab

import (
	"log"

	"github.com/xanzy/go-gitlab"
)

func CreateSngService() error {
	client, err := gitlab.NewBasicAuthClient(
		"zwf@k7.cn",
		"juexing20131009",
		gitlab.WithBaseURL("https://gitlab.kaiqitech.com/k7game/server"),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	projectOpts := &gitlab.CreateProjectOptions{
		Name:                 gitlab.String("test1-service"),
		Path:                 gitlab.String("test1-service"),
		Description:          gitlab.String("Just a test project to play with"),
		MergeRequestsEnabled: gitlab.Bool(true),
		SnippetsEnabled:      gitlab.Bool(true),
		Visibility:           gitlab.Visibility(gitlab.PublicVisibility),
		DefaultBranch:        gitlab.String("dev"),
	}
	if _, _, err := client.Projects.CreateProject(projectOpts); err != nil {
		log.Fatal(err)
	}
	return nil
}
