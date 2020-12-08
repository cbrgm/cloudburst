# Cloudburst.TargetsApi

All URIs are relative to *http://localhost/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**listScrapeTargets**](TargetsApi.md#listScrapeTargets) | **GET** /targets | List ScrapeTargets



## listScrapeTargets

> ScrapeTarget listScrapeTargets()

List ScrapeTargets

### Example

```javascript
import Cloudburst from 'cloudburst';

let apiInstance = new Cloudburst.TargetsApi();
apiInstance.listScrapeTargets().then((data) => {
  console.log('API called successfully. Returned data: ' + data);
}, (error) => {
  console.error(error);
});

```

### Parameters

This endpoint does not need any parameter.

### Return type

[**ScrapeTarget**](ScrapeTarget.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

