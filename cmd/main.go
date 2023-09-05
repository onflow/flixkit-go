package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/onflow/flixkit-go"
)

func main() {
	action := flag.String("action", "", "Specify the action to take.")
	flag.Parse()
	
	switch *action {
	case "generate":
		fmt.Println("Generating ...")

		templates := []string{
			"./templates/flow-transfer-tokens.template.json",
			"./templates/multiply.template.json",
			"https://flix.flow.com/v1/templates?name=transfer-flow",
		}

		for _, value := range templates {
			GenerateTemplateBinding(value)
		}
	
	default:
		fmt.Println("Invalid or missing action.")
	}
}

func GenerateTemplateBinding(flixQuery string) {
	flixService := flixkit.NewFlixService(&flixkit.Config{})
	ctx := context.Background()

	code, err := flixService.GenFlixBinding(ctx, flixQuery, "javascript")
	if err != nil {
		fmt.Errorf("could not process flix with path of %s: %w", flixQuery, err)
	}
	fmt.Println(code)
	parsedURL, err := url.Parse(flixQuery)
	filename := path.Base(parsedURL.Path)
	outputDirectory := path.Dir(parsedURL.Path)
	filenameNoExt := filename[0 : len(filename)-len("json")] + "js"
	SaveBindingFile(outputDirectory, filenameNoExt, code)
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