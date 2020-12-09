# Cloudburst.InstancesApi

All URIs are relative to *http://localhost/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**getInstances**](InstancesApi.md#getInstances) | **GET** /targets/{target}/instances | Get Instances for a ScrapeTarget
[**saveInstances**](InstancesApi.md#saveInstances) | **PUT** /targets/{target}/instances | Update Instances for a ScrapeTarget



## getInstances

> [Instance] getInstances(target)

Get Instances for a ScrapeTarget

### Example

```javascript
import Cloudburst from 'cloudburst';

let apiInstance = new Cloudburst.InstancesApi();
let target = "target_example"; // String | 
apiInstance.getInstances(target).then((data) => {
  console.log('API called successfully. Returned data: ' + data);
}, (error) => {
  console.error(error);
});

```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **target** | **String**|  | 

### Return type

[**[Instance]**](Instance.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## saveInstances

> [Instance] saveInstances(target, instance)

Update Instances for a ScrapeTarget

### Example

```javascript
import Cloudburst from 'cloudburst';

let apiInstance = new Cloudburst.InstancesApi();
let target = "target_example"; // String | 
let instance = [new Cloudburst.Instance()]; // [Instance] | 
apiInstance.saveInstances(target, instance).then((data) => {
  console.log('API called successfully. Returned data: ' + data);
}, (error) => {
  console.error(error);
});

```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **target** | **String**|  | 
 **instance** | [**[Instance]**](Instance.md)|  | 

### Return type

[**[Instance]**](Instance.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

