package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/onflow/flixkit-go"
)

type TemplateInfo struct {
	flixQuery string
	tmplReferencePath string
	filename string
}

func main() {
	action := flag.String("action", "", "Specify the action to take.")
	flag.Parse()
	
	switch *action {
	case "generate":
		fmt.Println("Generating ...")
		
		templates := [2]TemplateInfo{TemplateInfo{
			"file://templates/flow-transfer-tokens.template.json",
			"../templates/flow-transfer-tokens.template.json",
			"flow-transfer-tokens.template.js",
		}, TemplateInfo{
			"file://templates/multiply.template.json",
			"../templates/multiply.template.json",
			"multiply.template.js",
		}}
		bindingDirectory := "./bindings"

		for _, element := range templates {
			GenerateTemplateBinding(element.flixQuery, element.tmplReferencePath, element.filename, bindingDirectory)
		}
	
	default:
		fmt.Println("Invalid or missing action.")
	}
}

func GenerateTemplateBinding(flixQuery string, tmplReferencePath string, filename string, bindingDirectory string) {
	flixService := flixkit.NewFlixService(&flixkit.Config{})
	ctx := context.Background()

	code, err := flixService.GenFlixBinding(ctx, flixQuery, "javascript", tmplReferencePath)
	if err != nil {
		fmt.Errorf("could not process flix with path of %s: %w", flixQuery, err)
	}
	fmt.Println(code)
	SaveBindingFile(bindingDirectory, filename, code)
}

func SaveBindingFile(dirPath string, filename string, code string) error {
	// Create the directory if it doesn't exist
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			fmt.Printf("Error creating directory: %v\n", err)
			return err
		}
	}

	// Write the string content to the file
	err := os.WriteFile(dirPath + "/" + filename, []byte(code), 0644)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return err
	}
	return nil
}