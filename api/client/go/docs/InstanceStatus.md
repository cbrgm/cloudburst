# InstanceStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Agent** | **string** |  | 
**Status** | **string** |  | 
**Started** | [**time.Time**](time.Time.md) |  | 

## Methods

### NewInstanceStatus

`func NewInstanceStatus(agent string, status string, started time.Time, ) *InstanceStatus`

NewInstanceStatus instantiates a new InstanceStatus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInstanceStatusWithDefaults

`func NewInstanceStatusWithDefaults() *InstanceStatus`

NewInstanceStatusWithDefaults instantiates a new InstanceStatus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAgent

`func (o *InstanceStatus) GetAgent() string`

GetAgent returns the Agent field if non-nil, zero value otherwise.

### GetAgentOk

`func (o *InstanceStatus) GetAgentOk() (*string, bool)`

GetAgentOk returns a tuple with the Agent field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAgent

`func (o *InstanceStatus) SetAgent(v string)`

SetAgent sets Agent field to given value.


### GetStatus

`func (o *InstanceStatus) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *InstanceStatus) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *InstanceStatus) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetStarted

`func (o *InstanceStatus) GetStarted() time.Time`

GetStarted returns the Started field if non-nil, zero value otherwise.

### GetStartedOk

`func (o *InstanceStatus) GetStartedOk() (*time.Time, bool)`

GetStartedOk returns a tuple with the Started field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStarted

`func (o *InstanceStatus) SetStarted(v time.Time)`

SetStarted sets Started field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


