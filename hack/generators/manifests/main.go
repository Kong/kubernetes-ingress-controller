package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/yaml"
)

// -----------------------------------------------------------------------------
// This script exists as a workaround for some problems which came up between
// v1 and v2 of the KIC regarding upgrading and generating CRDs:
//
// 1. upgrades from v1beta1 CRD to v1 needed a fix for a backwards incompatible
//    change made in upstream that makes "preserveUnknownFields" no longer usable
//    in the same form that it previously was
// 2. controller-gen automatic generation of CRDs did not work out of the box for
//    some of our complex "Kong<TYPE>" resources, and so this script has a list
//    of types to filter out as we're still manually managing those CRDs
//
// -----------------------------------------------------------------------------

// -----------------------------------------------------------------------------
// Details on CRD v1's "PreserveUnknownFields"
//
// In the v1beta1 version of the API a field called "preserveUnknownFields" was
// used to trigger pruning for API fields that are sent in the request payload
// but not a valid part of the API. In v1 this field is effectively deprecated
// and MUST be set to false:
//
//  https://github.com/kubernetes/kubernetes/pull/93078
//
// Problematically: even though validation enforces the "false" setting of this
// deprecated attribute for v1 CRDs, the previous default was "true" in v1beta1
// and if you try to apply a v1 CRD that does NOT have this field set where a
// v1beta1 version of the CRD already exists on the cluster this can result in
// and error and fail to upgrade the CRD.
//
// This program accounts for this by applying "preserveUnknownFields: false" to
// any v1 CustomResourceDefinition to work around the upstream backwards
// incompatibility issue automatically.
//
// NOTE: When this was first written using controller-gen configuration options
//       and kubebuilder flags on the API types was attempted but none of these
//       things could produce the right result.
//
//       The controller-gen CLI can't really do it (easily) because the default
//       behavior of the marshaller in sigs.k8s.io/yaml is to not emit optional
//       boolean fields to YAML if they're set to false. 😑
//
//       It's quite possible there's a better way to do this and if you find one
//       we encourage you to put in a PR and replace this.
//
// -----------------------------------------------------------------------------

// -----------------------------------------------------------------------------
// Vars
// -----------------------------------------------------------------------------

var (
	// inputDir that will be used for walking the filesystem to find CRDs.
	inputDir string

	// outputDir is where this script will emit completed CRDs
	outputDir string

	// patchfile is the patch that will be used for each CRD file
	patchfile = "config/crd/patches/upgrade_compat.yaml"

	// excludedTypes is the list of API types which we DON'T want to automatically
	// generate CRDs for, as those CRDs are currently manually maintained.
	excludedTypes = [4]string{
		"KongClusterPlugin",
		"KongPlugin",
		"KongConsumer",
		"KongIngress",
	}
)

// -----------------------------------------------------------------------------
// Main
// -----------------------------------------------------------------------------

func main() {
	// handle arguments
	flag.StringVar(&inputDir, "input-directory", ".", "which directory to find CRDs in")
	flag.StringVar(&outputDir, "output-directory", ".", "which directory to emit completed CRDs to")
	flag.Parse()

	// find all the YAML CRDs in the provided directory and patch them
	if err := filepath.Walk(inputDir, processFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error processing files in %s: %s\n", inputDir, err)
		os.Exit(1)
	}
}

func processFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	// filter out any files that don't identify themselves as YAML files via file extensions
	if !strings.HasSuffix(info.Name(), ".yaml") && !strings.HasSuffix(info.Name(), ".yml") {
		return nil
	}

	// read in the file contents
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// filter out any YAML files that don't contain a v1.CustomResourceDefinition
	if !bytes.Contains(b, []byte(`kind: CustomResourceDefinition`)) || !bytes.Contains(b, []byte(`apiVersion: apiextensions.k8s.io/v1`)) {
		return nil
	}

	// don't allow files that contain multiple objects (controller-gen wont be configured to emit these anyway)
	yamlCRDs := bytes.Split(b, []byte(`\n---\n`))
	if len(yamlCRDs) > 1 {
		return fmt.Errorf("file %s contained multiple objects (%d) which this script doesn't yet support", path, len(yamlCRDs))
	}

	// filter out YAML files for CRDs we specifically want to exclude from automatic generation
	for _, excludedType := range excludedTypes {
		// check whether the excluded type is even present
		if bytes.Contains(b, []byte(excludedType)) {
			// the excluded type is present in this file, marshal it into YAML
			crd := v1.CustomResourceDefinition{}
			if err := yaml.Unmarshal(b, &crd); err != nil {
				return err
			}
			if crd.Spec.Names.Kind == excludedType {
				return nil
			}
		}
	}

	// generate a patch command to ensure "preserveUnknownFields" is set to false
	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)
	patchCMD := exec.Command("kubectl", "patch", "--type", "merge", "--local", "-o", "yaml", "-f", path, "--patch-file", patchfile)
	patchCMD.Stdout = stdout
	patchCMD.Stderr = stderr

	// run the patcher and make sure there are no errors
	if err := patchCMD.Run(); err != nil {
		return fmt.Errorf("failed to patch crd in %s: %w; STDOUT=(%s) STDERR=(%s)",
			path, err, strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String()),
		)
	}

	// controller-gen will add a YAML delimiter in the header, so we match that for consistency
	patchedCRD := append([]byte("\n---\n"), stdout.Bytes()...)

	// write the patched CRD back to the original file
	outputPath := fmt.Sprintf("%s/%s", outputDir, info.Name())
	return os.WriteFile(outputPath, patchedCRD, info.Mode().Perm())
}
