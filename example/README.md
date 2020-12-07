
## Workload Simulation


## Manual Workloads

### Request
```
curl -X POST -i http://localhost:8080/bubblesort --data '{
  "numbers": [
    252346,
    323,
    3422,
    55523
  ]
}'
```

### Result

```
{
  "numbers": [
    252346,
    323,
    3422,
    55523
  ],
  "sorted": [
    323,
    3422,
    55523,
    252346
  ]
}
```

## Metrics

```
curl http://localhost:8080/metrics
```

## Benchmarking

```
docker run --rm -it --network="host" -v $(pwd):/workload cbrgm/artillery:latest /bin/sh
```