//go:build generate_gateway_api_urls
// +build generate_gateway_api_urls

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

//go:generate go run --tags generate_gateway_api_urls . -crds-standard-url $CRDS_STANDARD_URL -crds-experimental-url $CRDS_EXPERIMENTAL_URL -raw-repo-url $RAW_REPO_URL -in $INPUT -out $OUTPUT

var (
	crdsStandardURLFlag     = flag.String("crds-standard-url", "", "The URL of standard Gateway API CRDs to be consumed by kustomize")
	crdsExperimentalURLFlag = flag.String("crds-experimental-url", "", "The URL of experimental Gateway API CRDs to be consumed by kustomize")
	rawRepoURLFlag          = flag.String("raw-repo-url", "", "The raw URL of Gateway API repository")
	inFlag                  = flag.String("in", "", "Template file path")
	outFlag                 = flag.String("out", "", "Output file path where the generate file will be placed")
)

type Data struct {
	CRDsStandardKustomizeURL     string
	CRDsExperimentalKustomizeURL string
	RawRepoURL                   string
}

func main() {
	flagParse()

	data := Data{
		CRDsStandardKustomizeURL:     *crdsStandardURLFlag,
		CRDsExperimentalKustomizeURL: *crdsExperimentalURLFlag,
		RawRepoURL:                   *rawRepoURLFlag,
	}
	processTemplate(*inFlag, *outFlag, data)
}

func must(err error, errMsg string) {
	if err != nil {
		log.Fatalf("%s: %v", errMsg, err)
	}
}

func flagParse() {
	flag.Parse()
	if *crdsStandardURLFlag == "" {
		log.Print("Please provide the 'crds-standard-url' flag")
		os.Exit(0)
	}
	if *crdsExperimentalURLFlag == "" {
		log.Print("Please provide the 'crds-experimental-url' flag")
		os.Exit(0)
	}
	if *rawRepoURLFlag == "" {
		log.Print("Please provide the 'raw-repo-url' flag")
		os.Exit(0)
	}
	if *inFlag == "" {
		log.Print("Please provide the 'in' flag")
		os.Exit(0)
	}
	if *outFlag == "" {
		log.Print("Please provide the 'out' flag")
		os.Exit(0)
	}
}

func processTemplate(fileName string, outputFile string, data Data) {
	base := filepath.Base(fileName)
	tmpl, err := template.New(base).ParseFiles(fileName)
	must(err, "Unable to parse template file")

	var processed bytes.Buffer
	err = tmpl.Execute(&processed, data)
	must(err, "Unable to parse data into template")

	formatted, err := format.Source(processed.Bytes())
	must(err, "Unable to format resulting file")

	outputPath := outputFile

	f, err := os.Create(outputPath)
	must(err, fmt.Sprintf("Unable to create file: %v", outputPath))

	w := bufio.NewWriter(f)
	_, err = w.Write(formatted)
	must(err, "Unable to underlying buffered writer")

	err = w.Flush()
	must(err, "Unable to flush")
}
