package main

import (
	"fmt"
	"log"

	"github.com/suzuki-shunsuke/gen-go-jsonschema/jsonschema"
	"github.com/suzuki-shunsuke/tfprovidercheck/pkg/controller"
)

func main() {
	if err := core(); err != nil {
		log.Fatal(err)
	}
}

func core() error {
	if err := jsonschema.Write(&controller.Config{}, "json-schema/tfprovidercheck.json"); err != nil {
		return fmt.Errorf("create or update a JSON Schema: %w", err)
	}
	return nil
}
