# Experiments

```
ceil(((((rate(example_sorting_request_duration_seconds_sum[1m]) / rate(example_sorting_request_duration_seconds_count[1m])) / 0.040) -1) * 10) + 0.5)


workload-normal-constant.yaml

  phases:
    - duration: 600s
      arrivalRate: 500
```