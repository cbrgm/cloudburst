# \InstancesApi

All URIs are relative to *http://localhost/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetInstances**](InstancesApi.md#GetInstances) | **Get** /targets/{target}/instances | Get Instances for a ScrapeTarget
[**SaveInstances**](InstancesApi.md#SaveInstances) | **Put** /targets/{target}/instances | Update Instances for a ScrapeTarget



## GetInstances

> []Instance GetInstances(ctx, target).Execute()

Get Instances for a ScrapeTarget

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
    target := "target_example" // string | 

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.InstancesApi.GetInstances(context.Background(), target).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `InstancesApi.GetInstances``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GetInstances`: []Instance
    fmt.Fprintf(os.Stdout, "Response from `InstancesApi.GetInstances`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**target** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiGetInstancesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**[]Instance**](Instance.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## SaveInstances

> []Instance SaveInstances(ctx, target).Instance(instance).Execute()

Update Instances for a ScrapeTarget

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
    target := "target_example" // string | 
    instance := []openapiclient.Instance{*openapiclient.NewInstance("Name_example", "Endpoint_example", "Provider_example", false, *openapiclient.NewContainerSpec("Name_example", "Image_example"), *openapiclient.NewInstanceStatus("Agent_example", "Status_example", "TODO"))} // []Instance | 

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.InstancesApi.SaveInstances(context.Background(), target).Instance(instance).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `InstancesApi.SaveInstances``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `SaveInstances`: []Instance
    fmt.Fprintf(os.Stdout, "Response from `InstancesApi.SaveInstances`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**target** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiSaveInstancesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **instance** | [**[]Instance**](Instance.md) |  | 

### Return type

[**[]Instance**](Instance.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

