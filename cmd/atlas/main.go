package main

import (
	"fmt"
	"log"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/bndrmrtn/interleaf/internal/db"
)

func main() {
	schema, err := gormschema.New("postgres").Load(db.GetModels()...)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(schema)
}
