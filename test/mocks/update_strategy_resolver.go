package mocks

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/stretchr/testify/require"

	"github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/diagnostics"
)

const mockUpdateReturnedConfigSize = 22

// UpdateStrategyResolver is a mock implementation of sendconfig.UpdateStrategyResolver.
type UpdateStrategyResolver struct {
	updateCalledForURLs       []string
	lastUpdatedContentForURLs map[string]sendconfig.ContentWithHash
	errorsToReturnOnUpdate    map[string][]error
	lock                      sync.RWMutex
}

func NewUpdateStrategyResolver() *UpdateStrategyResolver {
	return &UpdateStrategyResolver{
		lastUpdatedContentForURLs: map[string]sendconfig.ContentWithHash{},
		errorsToReturnOnUpdate:    map[string][]error{},
	}
}

// ResolveUpdateStrategy returns a mocked UpdateStrategy that will track which URLs were called.
func (f *UpdateStrategyResolver) ResolveUpdateStrategy(c sendconfig.UpdateClient, _ *diagnostics.ClientDiagnostic) sendconfig.UpdateStrategy {
	f.lock.Lock()
	defer f.lock.Unlock()

	url := c.AdminAPIClient().BaseRootURL()
	return &UpdateStrategy{onUpdate: f.updateCalledForURLCallback(url)}
}

// ReturnErrorOnUpdate will cause the mockUpdateStrategy with a given Admin API URL to return an error on Update().
// Errors will be returned following FIFO order. Each call to this function adds a new error to the queue.
func (f *UpdateStrategyResolver) ReturnErrorOnUpdate(url string) {
	f.lock.Lock()
	defer f.lock.Unlock()

	f.errorsToReturnOnUpdate[url] = append(f.errorsToReturnOnUpdate[url], errors.New("error on update"))
}

// ReturnSpecificErrorOnUpdate will cause the mockUpdateStrategy with a given Admin API URL to return a specific error
// on Update() call. Errors will be returned following FIFO order. Each call to this function adds a new error to the queue.
func (f *UpdateStrategyResolver) ReturnSpecificErrorOnUpdate(url string, err error) {
	f.lock.Lock()
	defer f.lock.Unlock()

	f.errorsToReturnOnUpdate[url] = append(f.errorsToReturnOnUpdate[url], err)
}

// GetUpdateCalledForURLs returns the called URLs.
func (f *UpdateStrategyResolver) GetUpdateCalledForURLs() []string {
	f.lock.RLock()
	defer f.lock.RUnlock()

	urls := make([]string, 0, len(f.updateCalledForURLs))
	urls = append(urls, f.updateCalledForURLs...)
	return urls
}

// LastUpdatedContentForURL returns the last updated content for the given URL.
func (f *UpdateStrategyResolver) LastUpdatedContentForURL(url string) (sendconfig.ContentWithHash, bool) {
	f.lock.RLock()
	defer f.lock.RUnlock()
	c, ok := f.lastUpdatedContentForURLs[url]
	return c, ok
}

// AssertUpdateCalledForURLs asserts that the mockUpdateStrategy was called for the given URLs.
func (f *UpdateStrategyResolver) AssertUpdateCalledForURLs(t *testing.T, urls []string, msgAndArgs ...any) {
	t.Helper()

	f.lock.RLock()
	defer f.lock.RUnlock()

	if len(msgAndArgs) == 0 {
		msgAndArgs = []any{"update was not called for all URLs"}
	}
	require.ElementsMatch(t, urls, f.updateCalledForURLs, msgAndArgs...)
}

// AssertUpdateCalledForURLsWithGivenCount asserts that Update was called for the given URLs with the given count.
func (f *UpdateStrategyResolver) AssertUpdateCalledForURLsWithGivenCount(t *testing.T, urlToCount map[string]int, msgAndArgs ...any) {
	t.Helper()

	f.lock.RLock()
	defer f.lock.RUnlock()
	actualURLToCount := lo.CountValues(f.updateCalledForURLs)
	for url, callCount := range urlToCount {
		m := []any{
			fmt.Sprintf("URL %s should receive %d update calls", url, callCount),
		}
		m = append(m, msgAndArgs...)
		require.Equal(t, callCount, actualURLToCount[url], m...)
	}
}

// AssertNoUpdateCalled asserts that no Update was not called.
func (f *UpdateStrategyResolver) AssertNoUpdateCalled(t *testing.T) {
	t.Helper()

	f.lock.RLock()
	defer f.lock.RUnlock()

	require.Empty(t, f.updateCalledForURLs, "update was called")
}

// EventuallyGetLastUpdatedContentForURL waits for the given URL to be called and returns the last updated content.
func (f *UpdateStrategyResolver) EventuallyGetLastUpdatedContentForURL(
	t *testing.T, url string, waitTime, waitTick time.Duration, msgAndArgs ...any,
) sendconfig.ContentWithHash {
	t.Helper()

	var content sendconfig.ContentWithHash
	if len(msgAndArgs) == 0 {
		msgAndArgs = []any{"update was not called for URL " + url}
	}
	require.Eventually(t, func() bool {
		c, ok := f.LastUpdatedContentForURL(url)
		if ok {
			content = c
			return true
		}
		return false
	}, waitTime, waitTick, msgAndArgs...)
	return content
}

// updateCalledForURLCallback returns a function that will be called when the mockUpdateStrategy is called.
// That enables us to track which URLs were called.
func (f *UpdateStrategyResolver) updateCalledForURLCallback(url string) func(sendconfig.ContentWithHash) (mo.Option[int], error) {
	return func(content sendconfig.ContentWithHash) (mo.Option[int], error) {
		f.lock.Lock()
		defer f.lock.Unlock()

		f.updateCalledForURLs = append(f.updateCalledForURLs, url)
		f.lastUpdatedContentForURLs[url] = content
		if errsToReturn, ok := f.errorsToReturnOnUpdate[url]; ok {
			if len(errsToReturn) > 0 {
				err := errsToReturn[0]
				f.errorsToReturnOnUpdate[url] = errsToReturn[1:]
				return mo.None[int](), err
			}
			return mo.Some(mockUpdateReturnedConfigSize), nil
		}
		return mo.Some(mockUpdateReturnedConfigSize), nil
	}
}
