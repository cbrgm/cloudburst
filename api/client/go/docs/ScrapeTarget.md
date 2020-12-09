# ScrapeTarget

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** |  | [optional] 
**Description** | Pointer to **string** |  | [optional] 
**Query** | Pointer to **string** |  | [optional] 
**InstanceSpec** | Pointer to [**InstanceSpec**](InstanceSpec.md) |  | [optional] 

## Methods

### NewScrapeTarget

`func NewScrapeTarget() *ScrapeTarget`

NewScrapeTarget instantiates a new ScrapeTarget object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewScrapeTargetWithDefaults

`func NewScrapeTargetWithDefaults() *ScrapeTarget`

NewScrapeTargetWithDefaults instantiates a new ScrapeTarget object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ScrapeTarget) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ScrapeTarget) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ScrapeTarget) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ScrapeTarget) HasName() bool`

HasName returns a boolean if a field has been set.

### GetDescription

`func (o *ScrapeTarget) GetDescription() string`

GetDescription returns the Description field if non-nil, zero value otherwise.

### GetDescriptionOk

`func (o *ScrapeTarget) GetDescriptionOk() (*string, bool)`

GetDescriptionOk returns a tuple with the Description field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDescription

`func (o *ScrapeTarget) SetDescription(v string)`

SetDescription sets Description field to given value.

### HasDescription

`func (o *ScrapeTarget) HasDescription() bool`

HasDescription returns a boolean if a field has been set.

### GetQuery

`func (o *ScrapeTarget) GetQuery() string`

GetQuery returns the Query field if non-nil, zero value otherwise.

### GetQueryOk

`func (o *ScrapeTarget) GetQueryOk() (*string, bool)`

GetQueryOk returns a tuple with the Query field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetQuery

`func (o *ScrapeTarget) SetQuery(v string)`

SetQuery sets Query field to given value.

### HasQuery

`func (o *ScrapeTarget) HasQuery() bool`

HasQuery returns a boolean if a field has been set.

### GetInstanceSpec

`func (o *ScrapeTarget) GetInstanceSpec() InstanceSpec`

GetInstanceSpec returns the InstanceSpec field if non-nil, zero value otherwise.

### GetInstanceSpecOk

`func (o *ScrapeTarget) GetInstanceSpecOk() (*InstanceSpec, bool)`

GetInstanceSpecOk returns a tuple with the InstanceSpec field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInstanceSpec

`func (o *ScrapeTarget) SetInstanceSpec(v InstanceSpec)`

SetInstanceSpec sets InstanceSpec field to given value.

### HasInstanceSpec

`func (o *ScrapeTarget) HasInstanceSpec() bool`

HasInstanceSpec returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


