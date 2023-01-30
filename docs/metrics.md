## Metrics

The following metrics are available to be scraped by Prometheus on traffic port and `/metrics` path

| Metric Name | Meaning |
| ----------- | ----------- |
go_* | Standard Go instrumentation metrics
process_* | Standard Go instrumentation metrics
promhttp_metric_handler_requests_in_flight | Amount of parallel requests in flight
promhttp_metric_handler_requests_total | Total count of metric requests (scrapes)
kubewebhook_mutating_webhook_review_duration_seconds | Webhook duration, in buckets for percentiles
kubewebhook_webhook_review_warnings_total	| Webhook warnings
k8s_image_swapper_ecr_errors | Number of errors related to ecr provider
k8s_image_swapper_main_errors | Number of errors
k8s_image_swapper_cache_hits | Number of registry cache hits
k8s_image_swapper_cache_misses | Number of registry cache misses
k8s_image_swapper_cache_filtered | Number of registry cache filtered out
k8s_image_swapper_cache_images_copied | Number of images copied
k8s_image_swapper_cache_repos_created | Number of repositories created
