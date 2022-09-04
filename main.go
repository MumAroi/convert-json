package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/iancoleman/strcase"
	"github.com/urfave/cli"
)

var app = cli.NewApp()

func main() {
	info()
	commands()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func info() {
	app.Name = "Convert Json"
	app.Usage = "An basic CLI for convert josn"
	app.Author = "paramas"
	app.Version = "1.0.0"
}

func commands() {

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "file, f",
			Value: "sample.json",
			Usage: "Load json from `FILE`",
		},
		cli.StringFlag{
			Name:  "output, o",
			Value: "converted.json",
			Usage: "Load json from `FILE`",
		},
	}

	app.Action = func(c *cli.Context) error {

		inputFile := c.String("file")
		outputFile := c.String("output")

		fileContent, err := os.Open(inputFile)

		if err != nil {
			log.Fatal(err)
			return nil
		}

		defer fileContent.Close()

		byteResult, _ := ioutil.ReadAll(fileContent)

		var res map[string]interface{}
		json.Unmarshal([]byte(byteResult), &res)
		if err != nil {
			panic(err)
		}

		for k, _ := range res["inputs"].([]interface{}) {

			input := res["inputs"].([]interface{})[k]
			caseType := input.(map[string]interface{})["caseType"]
			sensitive := input.(map[string]interface{})["sensitive"]
			text := input.(map[string]interface{})["text"]
			valStr := fmt.Sprintf("%v", text)
			if caseType == "camelCase" {
				// log.Println(strcase.ToLowerCamel(valStr))
				res["inputs"].([]interface{})[k].(map[string]interface{})["text"] = strcase.ToLowerCamel(valStr)
			} else if caseType == "snakeCase" {
				// log.Println(strcase.ToSnake(valStr))
				res["inputs"].([]interface{})[k].(map[string]interface{})["text"] = strcase.ToSnake(valStr)
			} else if caseType == "kebabCase" {
				// log.Println(strcase.ToKebab(valStr))
				res["inputs"].([]interface{})[k].(map[string]interface{})["text"] = strcase.ToKebab(valStr)
			} else {
				if res["defaultCaseType"] == "camelCase" {
					// log.Println(strcase.ToLowerCamel(valStr))
					res["inputs"].([]interface{})[k].(map[string]interface{})["text"] = strcase.ToLowerCamel(valStr)
					res["inputs"].([]interface{})[k].(map[string]interface{})["caseType"] = "camelCase"
				} else if res["defaultCaseType"] == "snakeCase" {
					// log.Println(strcase.ToSnake(valStr))
					res["inputs"].([]interface{})[k].(map[string]interface{})["text"] = strcase.ToSnake(valStr)
					res["inputs"].([]interface{})[k].(map[string]interface{})["caseType"] = "snakeCase"
				} else if caseType == "kebabCase" {
					// log.Println(strcase.ToKebab(valStr))
					res["inputs"].([]interface{})[k].(map[string]interface{})["text"] = strcase.ToKebab(valStr)
					res["inputs"].([]interface{})[k].(map[string]interface{})["caseType"] = "kebabCase"
				}
			}

			if sensitive == "true" {
				res["inputs"].([]interface{})[k].(map[string]interface{})["text"] = strings.Repeat("*", utf8.RuneCountInString(valStr))
				delete(res["inputs"].([]interface{})[k].(map[string]interface{}), "sensitive")
			} else if res["defaultSensitive"] == "true" {
				res["inputs"].([]interface{})[k].(map[string]interface{})["text"] = strings.Repeat("*", utf8.RuneCountInString(valStr))
				delete(res["inputs"].([]interface{})[k].(map[string]interface{}), "sensitive")
			}

		}

		file, _ := json.MarshalIndent(res["inputs"], "", " ")
		_ = ioutil.WriteFile(outputFile, file, 0644)

		return nil
	}

	sort.Sort(cli.FlagsByName(app.Flags))
}
