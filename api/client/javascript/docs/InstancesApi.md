# Cloudburst.InstancesApi

All URIs are relative to *http://localhost/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**updateInstances**](InstancesApi.md#updateInstances) | **PUT** /targets/{name}/instances | Update Instances for a ScrapeTarget



## updateInstances

> [Instance] updateInstances(name, instance)

Update Instances for a ScrapeTarget

### Example

```javascript
import Cloudburst from 'cloudburst';

let apiInstance = new Cloudburst.InstancesApi();
let name = "name_example"; // String | 
let instance = [new Cloudburst.Instance()]; // [Instance] | 
apiInstance.updateInstances(name, instance).then((data) => {
  console.log('API called successfully. Returned data: ' + data);
}, (error) => {
  console.error(error);
});

```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **String**|  | 
 **instance** | [**[Instance]**](Instance.md)|  | 

### Return type

[**[Instance]**](Instance.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

