package diagnostics

import (
	"context"
	"errors"
	"sync"

	"github.com/go-logr/logr"
	"github.com/kong/go-database-reconciler/pkg/file"
	"github.com/samber/mo"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/fallback"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/logging"
	managercfg "github.com/kong/kubernetes-ingress-controller/v3/pkg/manager/config"
)

const (
	// diagnosticConfigBufferDepth is the size of the channel buffer for receiving diagnostic
	// config dumps from the proxy sync loop. The chosen size is essentially arbitrary: we don't
	// expect that the receive end will get backlogged (it only assigns the value to a local
	// variable) but do want a small amount of leeway to account for goroutine scheduling, so it
	// is not zero.
	diagnosticConfigBufferDepth = 3

	// diffHistorySize is the number of diffs to keep in history.
	diffHistorySize = 5
)

// Collector collects diagnostic information from the proxy sync loop via Client it returns from Client()
// method. It can be queried for everything it collects via the Provider interface it implements.
type Collector struct {
	logger logr.Logger

	clientDiagnostic Client

	lastSuccessfulConfigDump mo.Option[file.Content]
	lastSuccessHash          string

	lastFailedConfigDump mo.Option[file.Content]
	lastFailedHash       string
	lastRawErrBody       []byte

	currentFallbackCacheMetadata mo.Option[fallback.GeneratedCacheMetadata]

	diffs diffMap

	configLock   sync.RWMutex
	fallbackLock sync.RWMutex
	diffLock     sync.RWMutex
}

func NewCollector(logger logr.Logger, cfg managercfg.Config) *Collector {
	return &Collector{
		logger: logger,
		diffs:  newDiffMap(diffHistorySize),
		clientDiagnostic: Client{
			DumpsIncludeSensitive: cfg.DumpSensitiveConfig,
			Configs:               make(chan ConfigDump, diagnosticConfigBufferDepth),
			FallbackCacheMetadata: make(chan fallback.GeneratedCacheMetadata, diagnosticConfigBufferDepth),
			Diffs:                 make(chan ConfigDiff, diagnosticConfigBufferDepth),
		},
	}
}

// Start starts the diagnostic collection loop. It will block until the context is done.
func (s *Collector) Start(ctx context.Context) error {
	return s.receiveDiagnostics(ctx)
}

// Client returns an object allowing dumping succeeded and failed configuration updates.
func (s *Collector) Client() Client {
	return s.clientDiagnostic
}

// LastSuccessfulConfigDump returns the last successful configuration dump.
func (s *Collector) LastSuccessfulConfigDump() (file.Content, string, bool) {
	s.configLock.RLock()
	defer s.configLock.RUnlock()

	if d, ok := s.lastSuccessfulConfigDump.Get(); ok {
		return *d.DeepCopy(), s.lastSuccessHash, true
	}
	return file.Content{}, "", false
}

// LastFailedConfigDump returns the last failed configuration dump.
func (s *Collector) LastFailedConfigDump() (file.Content, string, bool) {
	s.configLock.RLock()
	defer s.configLock.RUnlock()

	if d, ok := s.lastFailedConfigDump.Get(); ok {
		return *d.DeepCopy(), s.lastFailedHash, true
	}
	return file.Content{}, "", false
}

// LastErrorBody returns the raw error body of the last failed configuration push.
func (s *Collector) LastErrorBody() ([]byte, bool) {
	s.configLock.RLock()
	defer s.configLock.RUnlock()

	if s.lastRawErrBody != nil {
		return s.lastRawErrBody, true
	}
	return nil, false
}

// CurrentFallbackCacheMetadata returns the current fallback cache metadata.
func (s *Collector) CurrentFallbackCacheMetadata() mo.Option[fallback.GeneratedCacheMetadata] {
	s.fallbackLock.RLock()
	defer s.fallbackLock.RUnlock()

	return s.currentFallbackCacheMetadata
}

// LastConfigDiffHash returns the hash of the last config diff.
func (s *Collector) LastConfigDiffHash() string {
	s.diffLock.RLock()
	defer s.diffLock.RUnlock()

	return s.diffs.Latest()
}

// ConfigDiffByHash returns the config diff by hash.
func (s *Collector) ConfigDiffByHash(hash string) (ConfigDiff, bool) {
	s.diffLock.RLock()
	defer s.diffLock.RUnlock()

	return s.diffs.ByHash(hash)
}

// AvailableConfigDiffsHashes returns the hashes of available config diffs.
func (s *Collector) AvailableConfigDiffsHashes() []DiffIndex {
	s.diffLock.RLock()
	defer s.diffLock.RUnlock()

	return s.diffs.Available()
}

// receiveDiagnostics watches the diagnostic update channels.
func (s *Collector) receiveDiagnostics(ctx context.Context) error {
	for {
		select {
		case dump := <-s.clientDiagnostic.Configs:
			s.onConfigDump(dump)
		case meta := <-s.clientDiagnostic.FallbackCacheMetadata:
			s.onFallbackCacheMetadata(meta)
		case diff := <-s.clientDiagnostic.Diffs:
			s.onDiff(diff)
		case <-ctx.Done():
			if err := ctx.Err(); err != nil && !errors.Is(err, context.Canceled) {
				s.logger.Error(err, "Shutting down diagnostic collection: context completed with error")
				return err
			}
			s.logger.V(logging.InfoLevel).Info("Shutting down diagnostic collection: context completed")
			return nil
		}
	}
}

// onConfigDump handles a new configuration dump.
func (s *Collector) onConfigDump(dump ConfigDump) {
	s.configLock.Lock()
	defer s.configLock.Unlock()

	if dump.Meta.Failed {
		// If the config push failed, we need to keep the failed config dump and the raw error body.
		s.lastFailedConfigDump = mo.Some(dump.Config)
		s.lastFailedHash = dump.Meta.Hash
		s.lastRawErrBody = dump.RawResponseBody
	} else {
		// If the config push was successful, we need to keep successful config dump and the hash.
		s.lastSuccessfulConfigDump = mo.Some(dump.Config)
		s.lastSuccessHash = dump.Meta.Hash

		// If the regular config push was successful, we can drop the fallback cache metadata as it is
		// no longer relevant.
		if !dump.Meta.Fallback {
			s.fallbackLock.Lock()
			s.currentFallbackCacheMetadata = mo.None[fallback.GeneratedCacheMetadata]()
			s.fallbackLock.Unlock()
		}
	}
}

// onFallbackCacheMetadata handles a new fallback cache metadata.
func (s *Collector) onFallbackCacheMetadata(meta fallback.GeneratedCacheMetadata) {
	s.fallbackLock.Lock()
	defer s.fallbackLock.Unlock()

	s.currentFallbackCacheMetadata = mo.Some(meta)
}

// onDiff handles a new configuration diff.
func (s *Collector) onDiff(diff ConfigDiff) {
	s.diffLock.Lock()
	defer s.diffLock.Unlock()

	s.diffs.Update(diff)
}
