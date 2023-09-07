package flixkit

import (
	"fmt"
	"log"
)

func Generate(lang string, flix *FlowInteractionTemplate, templateLocation string, isLocal bool) (string, error) {
	var contents string
	var err error
	switch lang {
		case "javascript", "js":
			contents, err = GenerateJavaScript(flix, templateLocation, isLocal)
		default:
			return "", fmt.Errorf("language %s not supported", lang)
	}

    if err != nil {
        log.Fatalf("Error generating JavaScript: %v", err)
    }

	return contents, err
}
