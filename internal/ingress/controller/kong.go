/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/kong/deck/file"
	"github.com/kong/kubernetes-ingress-controller/pkg/deckgen"
	"github.com/kong/kubernetes-ingress-controller/pkg/kongstate"
	"github.com/kong/kubernetes-ingress-controller/pkg/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/pkg/util"
)

// OnUpdate is called periodically by syncQueue to keep the configuration in sync.
// returning nil implies the synchronization finished correctly.
// Returning an error means requeue the update.
func (n *KongController) OnUpdate(ctx context.Context, state *kongstate.KongState) error {
	var customEntities []byte
	var err error
	// process any custom entities
	if n.cfg.InMemory && n.cfg.KongCustomEntitiesSecret != "" {
		customEntities, err = n.fetchCustomEntities()
		if err != nil {
			// failure to fetch custom entities shouldn't block updates
			n.Logger.Errorf("failed to fetch custom entities: %v", err)
		}
	}

	filterTags := sendconfig.GetIngressControllerTags(n.cfg.Kong)

	targetContent := deckgen.ToDeckContent(ctx, n.Logger, state, &n.PluginSchemaStore, filterTags)

	newSHA, err := sendconfig.PerformUpdate(ctx,
		n.Logger,
		&n.cfg.Kong,
		n.cfg.InMemory,
		n.cfg.EnableReverseSync,
		targetContent,
		filterTags,
		customEntities,
		n.runningConfigHash,
	)

	if n.cfg.DumpConfig != util.ConfigDumpModeOff {
		if n.cfg.DumpConfig == util.ConfigDumpModeEnabled {
			targetContent = deckgen.ToDeckContent(ctx, n.Logger, state.SanitizedCopy(), &n.PluginSchemaStore,
				filterTags)
		}
		dumpErr := dumpConfig(err != nil, n.cfg.DumpDir, targetContent)
		if dumpErr != nil {
			n.Logger.WithError(err).Warn("failed to dump configuration")
		}
	}

	n.runningConfigHash = newSHA
	return err
}

func dumpConfig(failed bool, dumpDir string, targetContent *file.Content) error {
	target, err := json.Marshal(targetContent)
	if err != nil {
		return err
	}
	filename := "last_good.json"
	if failed {
		filename = "last_bad.json"
	}
	return ioutil.WriteFile(filepath.Join(dumpDir, filename), target, 0600)
}

func (n *KongController) fetchCustomEntities() ([]byte, error) {
	ns, name, err := util.ParseNameNS(n.cfg.KongCustomEntitiesSecret)
	if err != nil {
		return nil, fmt.Errorf("parsing kong custom entities secret: %w", err)
	}
	secret, err := n.store.GetSecret(ns, name)
	if err != nil {
		return nil, fmt.Errorf("fetching secret: %w", err)
	}
	config, ok := secret.Data["config"]
	if !ok {
		return nil, fmt.Errorf("'config' key not found in "+
			"custom entities secret '%v'", n.cfg.KongCustomEntitiesSecret)
	}
	return config, nil
}
