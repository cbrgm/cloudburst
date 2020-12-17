# InstanceEvent

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Type** | **string** |  | 
**Target** | **string** |  | 
**Data** | [**[]Instance**](Instance.md) |  | 

## Methods

### NewInstanceEvent

`func NewInstanceEvent(type_ string, target string, data []Instance, ) *InstanceEvent`

NewInstanceEvent instantiates a new InstanceEvent object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInstanceEventWithDefaults

`func NewInstanceEventWithDefaults() *InstanceEvent`

NewInstanceEventWithDefaults instantiates a new InstanceEvent object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetType

`func (o *InstanceEvent) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *InstanceEvent) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *InstanceEvent) SetType(v string)`

SetType sets Type field to given value.


### GetTarget

`func (o *InstanceEvent) GetTarget() string`

GetTarget returns the Target field if non-nil, zero value otherwise.

### GetTargetOk

`func (o *InstanceEvent) GetTargetOk() (*string, bool)`

GetTargetOk returns a tuple with the Target field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTarget

`func (o *InstanceEvent) SetTarget(v string)`

SetTarget sets Target field to given value.


### GetData

`func (o *InstanceEvent) GetData() []Instance`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *InstanceEvent) GetDataOk() (*[]Instance, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *InstanceEvent) SetData(v []Instance)`

SetData sets Data field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


