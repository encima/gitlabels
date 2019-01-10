package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

type User struct {
	AccessToken string
	Owner       string
}

type Config struct {
	User   User
	Labels []github.Label
}

//loadConfig loads specified file from the config path and sets it in the viper instance
func loadConfig(path string) Config {
	var config Config
	viper.SetConfigName("config")
	viper.AddConfigPath(path)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
	err = viper.Unmarshal(&config)
	handleErr(err)
	return config
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	config := loadConfig(".")
	fmt.Println(config)
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.User.AccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// list all repositories for the authenticated user
	repos, _, err := client.Repositories.List(ctx, config.User.Owner, nil)
	handleErr(err)

	for _, v := range repos {
		repo := v.GetName()
		fmt.Println(v.GetFullName())
		labels, _, err := client.Issues.ListLabels(ctx, config.User.Owner, v.GetName(), nil)
		handleErr(err)

		for _, v := range labels {
			lblName := v.GetName()
			fmt.Println(v.GetName())
			_, err := client.Issues.DeleteLabel(ctx, config.User.Owner, repo, lblName)
			handleErr(err)
			fmt.Println("Deleted")
		}

		for _, lbl := range config.Labels {
			fmt.Println(lbl)
			res, _, err := client.Issues.CreateLabel(ctx, config.User.Owner, repo, &lbl)
			handleErr(err)
			fmt.Println("Created Label", res)
		}
	}
}
