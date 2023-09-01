package bindings

import (
	"fmt"
	"log"

	js "github.com/onflow/flixkit-go/bindings/js"
	"github.com/onflow/flixkit-go/common"
)

func Generate(lang string, flix *common.FlowInteractionTemplate, templateLocation string) (string, error) {
	var contents string
	var err error
	switch lang {
		case "javascript", "js":
			contents, err = js.GenerateJavaScript(flix, templateLocation)
		default:
			return "", fmt.Errorf("language %s not supported", lang)
	}

    if err != nil {
        log.Fatalf("Error generating JavaScript: %v", err)
    }

	return contents, err
}
