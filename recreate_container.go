package main

import (
	"errors"
	recreate "github.com/falafeljan/docker-recreate"
	"github.com/fsouza/go-dockerclient"
	"github.com/graphql-go/graphql"
)

var recreationResponseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ContainerRecreationResponse",
	Fields: graphql.Fields{
		"previousContainerID": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The ID of the previous container",
		},
		"newContainerID": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "The ID of the newly created container.",
		},
	},
})

func recreateContainer(
	tokenContext TokenContext,
	registries []recreate.RegistryConf,
	accessToken string,
) (*recreate.Recreation, error) {
	tokenConf, err := tokenContext.GetGrant(accessToken)

	if tokenConf.ContainerID == "" || tokenConf.ImageTag == "" {
		return nil, errors.New("internal error")
	}

	client, err := docker.NewClient("unix:///var/run/docker.sock")
	if err != nil {
		return nil, err
	}

	dockerOptions := recreate.DockerOptions{
		PullImage:       true,
		DeleteContainer: true,
		Registries:      registries}
	context := recreate.NewContextWithClient(dockerOptions, client)

	recreation, err := context.Recreate(
		tokenConf.ContainerID,
		tokenConf.ImageTag,
		recreate.ContainerOptions{})

	if err != nil {
		return nil, err
	}

	return recreation, nil
}

func createContainerRecreationMutation(tokenContext TokenContext, registries []recreate.RegistryConf) *graphql.Field {
	return &graphql.Field{
		Type: graphql.NewNonNull(recreationResponseType),
		Args: graphql.FieldConfigArgument{
			"accessToken": &graphql.ArgumentConfig{
				Description: "A token with respective grants",
				Type:        graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			accessToken, ok := params.Args["accessToken"].(string)
			if !ok {
				return nil, errors.New("`accessToken` is expected to be a string")
			}

			return recreateContainer(tokenContext, registries, accessToken)
		},
	}
}
