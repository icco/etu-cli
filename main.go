package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/machinebox/graphql"
	"github.com/peterh/liner"
	"github.com/urfave/cli/v2"
)

var (
	history_fn = filepath.Join(os.TempDir(), ".etu.history")
)

type Config struct {
	Env string
	Key string
}

// Etu is the personifcation of time according to the Lakota.
func main() {
	cfg := &Config{}
	app := &cli.App{
		Name:  "etu",
		Usage: "log a project to etu.natwelch.com",
		Commands: []*cli.Command{
			{
				Name:    "print",
				Aliases: []string{"p"},
				Usage:   "print recent entries",
				Action:  cfg.Print,
			},
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "add a log",
				Action:  cfg.Add,
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "api_key",
				Usage:       "authorize your user",
				EnvVars:     []string{"GQL_TOKEN"},
				Destination: &cfg.Key,
			},
			&cli.StringFlag{
				Name:        "env",
				Usage:       "set which graphql server to talk to",
				Value:       "production",
				EnvVars:     []string{"NAT_ENV"},
				Destination: &cfg.Env,
			},
		},
	}

	err := app.RunContext(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

type AddHeaderTransport struct {
	T   http.RoundTripper
	Key string
}

func (adt *AddHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("X-API-AUTH", adt.Key)

	return adt.T.RoundTrip(req)
}

func (cfg *Config) Client() (*graphql.Client, error) {
	url := ""
	switch cfg.Env {
	case "production":
		url = "https://graphql.natwelch.com/graphql"
	case "development":
		url = "http://localhost:9393/graphql"
	default:
		return nil, fmt.Errorf("unknown environment %q", cfg.Env)
	}

	httpclient := &http.Client{Transport: &AddHeaderTransport{T: http.DefaultTransport, Key: cfg.Key}}
	return graphql.NewClient(url, graphql.WithHTTPClient(httpclient)), nil
}

func (cfg *Config) Add(c *cli.Context) error {
	line := liner.NewLiner()
	defer line.Close()

	if f, err := os.Open(history_fn); err == nil {
		line.ReadHistory(f)
		f.Close()
	}

	project, err := line.Prompt("What Project? ")
	if err != nil {
		return err
	}

	fmt.Println("Categories:")
	fmt.Println(" 1. Educating ")
	fmt.Println(" 2. Building")
	fmt.Println(" 3. Living")

	typeStr, err := line.Prompt("Category [1-3]? ")
	if err != nil {
		return err
	}

	focusStr, err := line.Prompt("Focus [1-9]? ")
	if err != nil {
		return err
	}

	introversionStr, err := line.Prompt("Introversion [1-9]? ")
	if err != nil {
		return err
	}
	code := fmt.Sprintf("%s%s%s", typeStr, focusStr, introversionStr)

	line.SetMultiLineMode(true)
	comment, err := line.Prompt("Comment? ")
	if err != nil {
		return err
	}

	if f, err := os.Create(history_fn); err != nil {
		log.Print("Error writing history file: ", err)
	} else {
		line.WriteHistory(f)
		f.Close()
	}

	client, err := cfg.Client()
	if err != nil {
		return err
	}

	gql := `
  mutation SaveLog($content: String!, $project: String!, $code: String!) {
    insertLog(
      input: { code: $code, description: $content, project: $project }
    ) {
      id
      datetime
    }
  }`

	req := graphql.NewRequest(gql)
	req.Var("content", comment)
	req.Var("code", code)
	req.Var("project", project)

	return client.Run(c.Context, req, nil)
}

func (cfg *Config) Print(c *cli.Context) error {
	return nil
}
