#!/bin/sh

benchmark_results () {
  outdir=$(basename -- "$1")
  outdir="${outdir%.*}"
  outdir="results/$outdir"
  mkdir -p $outdir
  get_benchmark_data $outdir
}

get_benchmark_data () {
  echo "creating benchmark graphs..."
  styx matplotlib --prometheus http://localhost:9090 --duration 10m 'rate(cloudburst_proxy_custom_total{key="incoming.HTTP/1.1"}[30s])' > $1/http_requests_total.py
  styx matplotlib --prometheus http://localhost:9090 --duration 10m '(rate(cloudburst_proxy_custom_total{key="incoming.HTTP/1.1"}[30s]) / 5)' > $1/http_requests_slo.py
  styx matplotlib --prometheus http://localhost:9090 --duration 10m 'cloudburst_calculator_instances_total{status="running"} + 1' > $1/instances_total_running.py
  styx matplotlib --prometheus http://localhost:9090 --duration 10m 'cloudburst_calculator_instances_total{provider="cloud-provider-A",status="running"}' > $1/instances_total_provider_a.py
  styx matplotlib --prometheus http://localhost:9090 --duration 10m 'cloudburst_calculator_instances_total{provider="cloud-provider-B",status="running"}' > $1/instances_total_provider_b.py
  styx matplotlib --prometheus http://localhost:9090 --duration 10m 'cloudburst_calculator_instances_total{provider="cloud-provider-C",status="running"}' > $1/instances_total_provider_c.py
  styx matplotlib --prometheus http://localhost:9090 --duration 10m 'sum(rate(example_sorting_requests_total[30s])) * 10' > $1/http_requests_single_service.py
  echo "done!"
  sleep 120
}

for filename in workloads/*.yaml; do
  echo "Running benchmark file $filename"
  docker run --rm -it --network="host" -v $(pwd):/workload cbrgm/artillery:latest /bin/sh -c "artillery run ${filename}"
  benchmark_results $filename
done