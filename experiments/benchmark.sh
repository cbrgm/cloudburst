#!/bin/sh

benchmark_results () {
  outdir=$(basename -- "$1")
  outdir="${outdir%.*}"
  outdir="results/$outdir"
  mkdir -p $outdir
  get_benchmark_data $outdir
}

get_benchmark_data () {
  styx matplotlib --prometheus http://localhost:9090 --duration 5m 'rate(cloudburst_proxy_custom_total{key="incoming.HTTP/1.1"}[30s])' > $1/http_requests_total.py
  styx matplotlib --prometheus http://localhost:9090 --duration 5m 'rate(cloudburst_proxy_custom_total{key="incoming.HTTP/1.1"}[30s]) / 15' > $1/http_requests_slo.py
  styx matplotlib --prometheus http://localhost:9090 --duration 5m 'cloudburst_api_instances_total{status="running"} + 1' > $1/instances_total.py
  styx matplotlib --prometheus http://localhost:9090 --duration 5m 'sum(rate(example_sorting_requests_total[30s]))' > $1/http_requests_single_service.py
}

for filename in workloads/*.yaml; do
  echo "Running benchmark file $filename"
  docker run --rm -it --network="host" -v $(pwd):/workload cbrgm/artillery:latest /bin/sh -c "artillery run ${filename}"
  benchmark_results $filename
done