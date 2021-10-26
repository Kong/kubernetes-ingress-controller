---
status: declined
---

# NOTES

The maintainers created this KEP originally to respond to a performance issue, but in the conversation about this KEP so far we've decided there simply isn't enough supporting evidence that this problem is actively affecting anyone and so we can't justify starting work on it yet.

This KEP exists in a `declined` state for posterity, and so we can pick back up where we left off if new reports surface that help inform its priority. If you're reading this and looking to get it started again, please feel free to check in on the [discussions](https://github.com/Kong/kubernetes-ingress-controller/discussions) or make an update to this KEP and provide the context about why you feel this needs attention.

# Unreferenced Kubernetes Object Filtration

<!-- toc -->
- [Summary](#summary)
- [Motivation](#motivation)
  - [Goals](#goals)
  - [Non-Goals](#non-goals)
- [Proposal](#proposal)
- [Drawbacks](#drawbacks)
- [Alternatives](#alternatives)
<!-- /toc -->

## Summary

Historically the [Kong Kubernetes Ingress Controller (KIC)][kic] has gathered up _all_ [Service][svc], [Endpoints][endp] and [Secrets][secrets] on the cluster in order to make sure they are available in the cache when building the backend Kong Gateway configuration. We want to improve performance for the KIC by filtering out `Services`, `Endpoints` and other resources which aren't referenced by `Ingress`, `TCPIngress`, `UDPIngress` and related resources we're managing.

[kic]:https://github.com/kong/kubernetes-ingress-controller
[svc]:https://kubernetes.io/docs/concepts/services-networking/service/
[endp]:https://kubernetes.io/docs/concepts/services-networking/endpoint-slices/
[secrets]:https://kubernetes.io/docs/concepts/configuration/secret/

## Motivation

- reduce excess CPU workload for the KIC controller manager
- reduce excess memory utilization for the KIC controller manager
- add logging for ingest of managed API resources

### Goals

- add filtration in our controller watches and reconcilation code paths to avoid processing and logging `Service`, `Endpoint` and `Secret` resources that are irrelevant
- avoid using memory and CPU resources to cache resources that aren't going to be used due to a lack of reference by other `Ingress` type objects we're managing
- add logging to indicate when objects are being processed by the backend proxy

## Proposal

The main problems we need to solve when writing logic which will make decisions about whether an object will be included in the cache are:

- identifying object references in a performant and threadsafe manner
- the trigger mechanism to populate the object into the cache
- dealing with multiple references to the same object

In order to solve these problems we can keep track of any referenced object using a reference counted index and implement that for filtration in places like the [watch predicates][watch] so that these unusable objects will never make it to reconcilation (and consequently will never make it into our cache) until it becomes referenced, and it wont be removed from the cache until all references to it are gone.

In order to populate the reference counted index we will need to hook into places like the `Proxy` implementations `UpdateObject()` and `DeleteObject()` methods, OR the reconcilation loops of the parent objects themselves. Inbound Kubernetes objects will need to be processed to derive lists of referenced objects and populate the index with them.

The index will be threadsafe and READ optimized (e.g. `sync.RWMutex` with `RLock()`) so that we can safely access it from any watch predicate function for objects we manage.

[watch]:https://cloud.redhat.com/blog/kubernetes-operators-best-practices#Creating-Watches

### WIP

The above proposal is high level and we still need to iron out some of the details.

There are a few different variations of the above concept which will be enumerated here with pros and cons for considerations. This is all a WIP, we maybe derive further variations after thinking on and discussing these further.

Ultimately if we decide to go with an entirely different proposal, move this into the `Alternatives` section for posterity.

#### Variation 1 - synchronous time of access indexing and processing

In this variation of the above proposal all indexing is done synchronously as objects are processed, meaning that the `UpdateObject()` and `DeleteObject()` methods are fully responsible for index population (and cleanup) and operationally the index is managed synchronously alongside changes to the cache.

When the index is changed, the cache is changed meaning that an `Ingress` object is not considered properly "cached" with `UpdateObject()` until the referenced objects are also populated in the cache and the index has been updated to reference them appropriately.

##### Pros

- fastest complete solution - this is likely the most straightforward implementation to maintain from a coding perspective which eliminates _all_ extraneous objects from being cached
- lowest complexity - less moving parts, less code, less things that can go wrong

##### Cons

- performance overhead - needing to make additional Kubernetes API requests while processing cache objects will make caching more expensive in CPU and Network utilization (though reducing memory utilization) and creates a bottleneck at the cache layer
- reduced separation of concerns (monolith) - this adds a new dynamic and responsibility to `Proxy` cache server implementations which are already slightly overloaded by their backend responsibilities

#### Variation 2 - asynchronous index processing

In this variation the `UpdateObject()` and `DeleteObject()` methods are responsible only for modifying the reference counted index and a discreet goroutine routine will be responsible for populating the cache with referenced objects. We will call this goroutine the `ReferencedObjectInformer`.

The `ReferencedObjectInformer` will receive updates from the index for new objects and referenced objects can be updated in the Kubernetes API by being given new status `metav1.Conditions` or metadata labels/annotations which will serve to both indicate that the object became referenced by a Kong managed API resource and will also serve to trigger reconcilation of the object which will then pass the watch predicates so that the object becomes consequently cached. Thereafter we can rely on the retry mechanisms for failure conditions already granted to us by `controller-runtime`.

A garbage collection routine will be responsible for identifying objects with a `0` reference count and pruning them from the cache and the index on a tunable interval.

##### Pros

- robust - this implementation is intended to be efficient and reduce performance overhead
- code hygiene - this implementation produces more discreet components which can be maintained and tested independently
- separation of concerns - doesn't add more responsibilities to `Proxy` implementations

##### Cons

The biggest drawback with this variation are complexities:

- more moving parts at runtime: more things that may need to be tuned
- the reference counted index will need to provide an update channel
- longest for initial implementation and more overall code to maintain long term

#### Variation 3 - namespace indexing only

This variation trades off overall efficacy of the solution in favor of short term gains and low maintenance costs. It can be thought of as a "quick win" solution.

In this variation only namespace objects are referenced in the index and watch predicates simply filter out unreferenced namespaces.

This is most effective when using the default "watch all namespaces" functionality of the controller manager.

##### Pros

- minimal - quick and easy to develop, very small amount of code, low maintenance costs
- lightweight - very little operational resource utilization overhead

##### Cons

- incomplete solution - reduces waste but leaves the overall problem partially intact

## Design Details

Here is an example implementation of the threadsafe reference counted index which can be used to reference count any Kubernetes object by its unique `<namespace>/<name>`:

```go
// ReferenceCountedStringIndex is a reference counting index of strings where the
// implementations must be threadsafe.
type ReferenceCountedStringIndex interface {
	// Insert increases the reference count for each provided entry by 1
	Insert(entries ...string)

	// Delete decreases the reference count for each provided entry by 1
	Delete(entries ...string)

	// ReferenceCount provides the current reference count for a provided entry
	ReferenceCount(entry string) (referenceCount int)

	// Len provides the total number of references currently being tracked
	Len() int
}

// NewReferenceCountedStringIndex provides a new ReferenceCountedStringIndex
func NewReferenceCountedStringIndex() ReferenceCountedStringIndex {
	return &index{
		index: make(map[string]int),
		lock:  sync.RWMutex{},
	}
}

type index struct {
	index map[string]int
	lock  sync.RWMutex
}

func (s *index) Insert(entries ...string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, entry := range entries {
		s.index[entry]++
	}
}

func (s *index) Delete(entries ...string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, entry := range entries {
		if s.index[entry] == 1 || s.index[entry] == 0 {
			delete(s.index, entry)
		} else {
			s.index[entry]--
		}
	}
}

func (s *index) ReferenceCount(entry string) int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.index[entry]
}

func (s *index) Len() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.index)
}
```

TODO: more design details once the proposal is further along.

## Drawbacks

The main drawback so far is that we haven't seen any production issues reported yet indicating that this has a significant consequence in real world applications. The problem domain as reported thus far has more theoretical than practical standing at the time of writing. As such we should be considering the largest gain we can get for the smallest cost, and/or perhaps even considering doing some additional performance testing and consequently perhaps even declining this KEP and using it only for posterity until such a time that there are more real reports of the issue.

## Notes - Prior Art

Our neighbors have some prior art which tries at some level to deal with this problem domain, and we should keep them in mind for our own implementation considerations.

### NGinx Ingress

The [NGinx Ingress Controller][ing-nginx] has some [filtration for secrets][ing-nginx-sec] that aren't relevant because they're deployed by [Helm][helm].

[ing-nginx]:https://github.com/kubernetes/ingress-nginx
[ing-nginx-sec]:https://github.com/kubernetes/ingress-nginx/blob/402f21bcb7402942f91258c8971ff472d81f5322/internal/ingress/controller/store/store.go#L252-L277
[helm]:https://github.com/helm/helm

### Emissary

The [Emissary Ingress Controller][emissary] builds [essentially unfiltered queries for relevant types][emissary-queries] which you can specify your own labels or field selectors for (empty by default) and [watches them][emissary-watch].

[emissary]:https://github.com/emissary-ingress/emissary
[emissary-queries]:https://github.com/emissary-ingress/emissary/blob/v2.0.1-rc.2/cmd/entrypoint/watcher.go#L36-L37
[emissary-watch]:https://github.com/emissary-ingress/emissary/blob/v2.0.1-rc.2/cmd/entrypoint/watcher.go#L168

### Tyk Operator

Tyk side-steps the problem by [using service hostnames always][tyk-impl] so it doesn't need a `Service` or `Endpoint` controller.

It performs filters on [well-defined fields][tyk-filters] that aren't a feasible solution for us because we care about non-TLS Secrets as well.

[tyk-impl]:https://github.com/TykTechnologies/tyk-operator/blob/v0.7.1/controllers/ingress_controller.go#L157
[tyk-filters]:https://github.com/TykTechnologies/tyk-operator/blob/v0.7.1/controllers/secret_cert_controller.go#L99-L103
