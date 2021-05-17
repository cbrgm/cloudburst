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
  styx matplotlib --prometheus http://localhost:9090 --duration 10m 'rate(go_memstats_alloc_bytes_total{instance="example-app:9997", job="example-app"}[30s])' > $1/memory_single_service.py
  styx matplotlib --prometheus http://localhost:9090 --duration 10m '((sum(rate(go_memstats_alloc_bytes_total{instance="example-app:9997", job="example-app"}[30s])) / 3000000) * ( 1 + sum(cloudburst_calculator_instances_total{status="running"})) + 0.5)' > $1/memory_requests_slo.py
  styx matplotlib --prometheus http://localhost:9090 --duration 10m 'sum(cloudburst_calculator_instances_total{status="running"}) + 1' > $1/instances_total_running.py
  styx matplotlib --prometheus http://localhost:9090 --duration 10m 'cloudburst_calculator_instances_total{provider="cloud-provider-alpha",status="running"}' > $1/instances_total_provider_alpha.py
  styx matplotlib --prometheus http://localhost:9090 --duration 10m 'cloudburst_calculator_instances_total{provider="cloud-provider-beta",status="running"}' > $1/instances_total_provider_beta.py
  styx matplotlib --prometheus http://localhost:9090 --duration 10m 'cloudburst_calculator_instances_total{provider="cloud-provider-gamma",status="running"}' > $1/instances_total_provider_gamma.py
  echo "done!"
  sleep 120
}

for filename in workloads/*.yaml; do
  echo "Running benchmark file $filename"
  docker run --rm -it --network="host" -v $(pwd):/workload cbrgm/artillery:latest /bin/sh -c "artillery run ${filename}"
  benchmark_results $filename
done