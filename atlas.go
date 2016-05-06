package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/codegangsta/cli"
	"github.com/hashicorp/atlas-go/v1"
)

var (
	Version string
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func ToJson(obj interface{}) (interface{}, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	return string(data), nil
}

func search(c *cli.Context) {

	if c.String("user") == "" || c.String("artifact") == "" || c.String("type") == "" {
		cli.ShowAppHelp(c)
		os.Exit(1)
	}

	searchOpts := &atlas.ArtifactSearchOpts{
		User: c.String("user"),
		Name: c.String("artifact"),
		Type: c.String("type"),
	}
	if len(c.StringSlice("meta")) > 0 {
		filter := map[string]string{}
		for _, m := range c.StringSlice("meta") {
			pair := strings.Split(m, "=")
			filter[pair[0]] = pair[1]
		}
		fmt.Fprintf(os.Stderr, "FILTER: %#v\n", filter)
		searchOpts.Metadata = filter
	}
	client := atlas.DefaultClient()

	versions, err := client.ArtifactSearch(searchOpts)
	if err != nil {
		fmt.Errorf("search error: %#v", err)
		os.Exit(1)
	}

	fnMap := template.FuncMap{}
	fnMap["json"] = ToJson
	tmpl, err := template.New("artifact").Funcs(fnMap).Parse(c.String("format"))
	if err != nil {
		panic(err)
	}

	for _, v := range versions {
		//fmt.Println("ver: %#v", v)
		err = tmpl.Execute(os.Stdout, v)
		if err != nil {
			panic(err)
		}
	}

}

func main() {
	fmt.Fprintln(os.Stderr, "Search atlas.hashicorp artifacts ...")

	app := cli.NewApp()
	app.Name = "atlifacts"
	app.Usage = "query atlas.hashicorp.com artifacts"
	app.Version = Version

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "user, u",
			Usage:  "atlas user",
			EnvVar: "ATLAS_USER",
		},
		cli.StringFlag{
			Name:   "artifact, a",
			Usage:  "atlas artifact",
			EnvVar: "ATLAS_ARTIFACT_NAME",
		},
		cli.StringFlag{
			Name:   "type, t",
			Usage:  "atlas artifact type",
			EnvVar: "ATLAS_ARTIFACT_TYPE",
		},
		cli.StringFlag{
			Name:  "format, f",
			Value: "{{.Slug}}\n",
			Usage: "output format in golang template",
		},
		cli.StringSliceFlag{
			Name:  "meta, m",
			Usage: "meta field as fielter",
		},
	}
	app.Action = search

	app.Run(os.Args)

}
