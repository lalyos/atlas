package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
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
	data, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return nil, err
	}
	return string(data), nil
}

// add returns the sum of a and b.
func add(b, a interface{}) (interface{}, error) {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	switch av.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return av.Int() + bv.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Int() + int64(bv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return float64(av.Int()) + bv.Float(), nil
		default:
			return nil, fmt.Errorf("add: unknown type for %q (%T)", bv, b)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return int64(av.Uint()) + bv.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Uint() + bv.Uint(), nil
		case reflect.Float32, reflect.Float64:
			return float64(av.Uint()) + bv.Float(), nil
		default:
			return nil, fmt.Errorf("add: unknown type for %q (%T)", bv, b)
		}
	case reflect.Float32, reflect.Float64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return av.Float() + float64(bv.Int()), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Float() + float64(bv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return av.Float() + bv.Float(), nil
		default:
			return nil, fmt.Errorf("add: unknown type for %q (%T)", bv, b)
		}
	default:
		return nil, fmt.Errorf("add: unknown type for %q (%T)", av, a)
	}
}

// subtract returns the difference of b from a.
func subtract(b, a interface{}) (interface{}, error) {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	switch av.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return av.Int() - bv.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Int() - int64(bv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return float64(av.Int()) - bv.Float(), nil
		default:
			return nil, fmt.Errorf("subtract: unknown type for %q (%T)", bv, b)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return int64(av.Uint()) - bv.Int(), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Uint() - bv.Uint(), nil
		case reflect.Float32, reflect.Float64:
			return float64(av.Uint()) - bv.Float(), nil
		default:
			return nil, fmt.Errorf("subtract: unknown type for %q (%T)", bv, b)
		}
	case reflect.Float32, reflect.Float64:
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return av.Float() - float64(bv.Int()), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return av.Float() - float64(bv.Uint()), nil
		case reflect.Float32, reflect.Float64:
			return av.Float() - bv.Float(), nil
		default:
			return nil, fmt.Errorf("subtract: unknown type for %q (%T)", bv, b)
		}
	default:
		return nil, fmt.Errorf("subtract: unknown type for %q (%T)", av, a)
	}
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
	fnMap["add"] = add
	fnMap["subtract"] = subtract
	tmpl, err := template.New("artifact").Funcs(fnMap).Parse(c.String("format") + "\n")
	if err != nil {
		panic(err)
	}

	if c.Bool("last") {
		versions = versions[:1]
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
			Value: "{{json .}}",
			Usage: "output format in golang template",
		},
		cli.StringSliceFlag{
			Name:  "meta, m",
			Usage: "meta field as filter",
		},
		cli.BoolFlag{
			Name:  "last, l",
			Usage: "only show the latest version",
		},
	}
	app.Action = search

	app.Run(os.Args)

}
