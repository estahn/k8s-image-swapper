/*
Copyright © 2020 Enrico Stahn <enrico.stahn@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var (
	commonRepoLabels = []string{"resource_namespace", "registry", "repo"}

	errors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "k8s_image_swapper",
			Subsystem: "main",
			Name:      "errors",
			Help:      "Number of errors",
		},
		[]string{"error_type"},
	)
	ecrErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "k8s_image_swapper",
			Subsystem: "ecr",
			Name:      "errors",
			Help:      "Number of ecr errors",
		},
		append(commonRepoLabels, "error_type"),
	)

	cacheHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "k8s_image_swapper",
			Subsystem: "cache",
			Name:      "hits",
			Help:      "Number of registry cache hits",
		},
		commonRepoLabels,
	)

	cacheMisses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "k8s_image_swapper",
			Subsystem: "cache",
			Name:      "misses",
			Help:      "Number of registry cache misses",
		},
		commonRepoLabels,
	)
	cacheFiltered = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "k8s_image_swapper",
			Subsystem: "cache",
			Name:      "filtered",
			Help:      "Number of registry cache filtered out",
		},
		commonRepoLabels,
	)

	imageCopyDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "k8s_image_swapper",
			Subsystem: "cache",
			Name:      "image_copy_duration_seconds",
			Help:      "Image copy duration distribution in seconds",
			Buckets:   prometheus.ExponentialBuckets(4, 2, 10),
		},
		commonRepoLabels,
	)

	reposCreateRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "k8s_image_swapper",
			Subsystem: "cache",
			Name:      "repos_create_requests",
			Help:      "Number of repository create requests",
		},
		commonRepoLabels,
	)
	reposCreated = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: "cache",
			Name:      "repos_created",
			Help:      "Number of repositories created",
		},
		[]string{"registry", "repo"},
	)
)

var PromReg *prometheus.Registry

func init() {
	PromReg = prometheus.NewRegistry()
	PromReg.MustRegister(collectors.NewGoCollector())
	PromReg.MustRegister(errors)
	PromReg.MustRegister(ecrErrors)
	PromReg.MustRegister(cacheHits)
	PromReg.MustRegister(cacheMisses)
	PromReg.MustRegister(cacheFiltered)
	PromReg.MustRegister(imageCopyDuration)
	PromReg.MustRegister(reposCreateRequests)
	PromReg.MustRegister(reposCreated)
}

// Increments the counter of errors
func IncrementError(errType string) {
	errors.With(
		prometheus.Labels{
			"error_type": errType,
		},
	).Inc()
}

// Increments the counter of ecr errors
func IncrementEcrError(resource_namespace string, registry string, repo string, errType string) {
	ecrErrors.With(
		prometheus.Labels{
			"resource_namespace": resource_namespace,
			"registry":           registry,
			"repo":               repo,
			"error_type":         errType,
		},
	).Inc()
}

// Increments the counter of registry cache hits
func IncrementCacheHits(resource_namespace string, registry string, repo string) {
	cacheHits.With(
		prometheus.Labels{
			"resource_namespace": resource_namespace,
			"registry":           registry,
			"repo":               repo,
		},
	).Inc()
}

// Increments the counter of registry cache misses
func IncrementCacheMisses(resource_namespace string, registry string, repo string) {
	cacheMisses.With(
		prometheus.Labels{
			"resource_namespace": resource_namespace,
			"registry":           registry,
			"repo":               repo,
		},
	).Inc()
}

// Increments the counter of registry cache ignored/filtered out
func IncrementCacheFiltered(resource_namespace string, registry string, repo string) {
	cacheFiltered.With(
		prometheus.Labels{
			"resource_namespace": resource_namespace,
			"registry":           registry,
			"repo":               repo,
		},
	).Inc()
}

// Sets the duration of image copy operation
func SetImageCopyDuration(resource_namespace string, registry string, repo string, duration float64) {
	imageCopyDuration.With(
		prometheus.Labels{
			"resource_namespace": resource_namespace,
			"registry":           registry,
			"repo":               repo,
		},
	).Observe(duration)
}

// Increments the counter of repo create requests
func IncrementReposCreateRequests(resource_namespace string, registry string, repo string) {
	reposCreateRequests.With(
		prometheus.Labels{
			"resource_namespace": resource_namespace,
			"registry":           registry,
			"repo":               repo,
		},
	).Inc()
}

// Increments the counter of repos created
func IncrementReposCreated(registry string, repo string) {
	reposCreated.With(
		prometheus.Labels{
			"registry": registry,
			"repo":     repo,
		},
	).Inc()
}
