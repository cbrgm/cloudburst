# \TargetsApi

All URIs are relative to *http://localhost/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ListScrapeTargets**](TargetsApi.md#ListScrapeTargets) | **Get** /targets | List ScrapeTargets



## ListScrapeTargets

> ScrapeTarget ListScrapeTargets(ctx).Execute()

List ScrapeTargets

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.TargetsApi.ListScrapeTargets(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `TargetsApi.ListScrapeTargets``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListScrapeTargets`: ScrapeTarget
    fmt.Fprintf(os.Stdout, "Response from `TargetsApi.ListScrapeTargets`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiListScrapeTargetsRequest struct via the builder pattern


### Return type

[**ScrapeTarget**](ScrapeTarget.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

