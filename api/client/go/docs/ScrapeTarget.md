# ScrapeTarget

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | 
**Path** | **string** |  | 
**Description** | **string** |  | 
**Query** | **string** |  | 
**InstanceSpec** | [**InstanceSpec**](InstanceSpec.md) |  | 
**StaticSpec** | [**StaticSpec**](StaticSpec.md) |  | 

## Methods

### NewScrapeTarget

`func NewScrapeTarget(name string, path string, description string, query string, instanceSpec InstanceSpec, staticSpec StaticSpec, ) *ScrapeTarget`

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


### GetPath

`func (o *ScrapeTarget) GetPath() string`

GetPath returns the Path field if non-nil, zero value otherwise.

### GetPathOk

`func (o *ScrapeTarget) GetPathOk() (*string, bool)`

GetPathOk returns a tuple with the Path field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPath

`func (o *ScrapeTarget) SetPath(v string)`

SetPath sets Path field to given value.


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


### GetStaticSpec

`func (o *ScrapeTarget) GetStaticSpec() StaticSpec`

GetStaticSpec returns the StaticSpec field if non-nil, zero value otherwise.

### GetStaticSpecOk

`func (o *ScrapeTarget) GetStaticSpecOk() (*StaticSpec, bool)`

GetStaticSpecOk returns a tuple with the StaticSpec field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStaticSpec

`func (o *ScrapeTarget) SetStaticSpec(v StaticSpec)`

SetStaticSpec sets StaticSpec field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


