# InstanceSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Container** | [**ContainerSpec**](ContainerSpec.md) |  | 

## Methods

### NewInstanceSpec

`func NewInstanceSpec(container ContainerSpec, ) *InstanceSpec`

NewInstanceSpec instantiates a new InstanceSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInstanceSpecWithDefaults

`func NewInstanceSpecWithDefaults() *InstanceSpec`

NewInstanceSpecWithDefaults instantiates a new InstanceSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetContainer

`func (o *InstanceSpec) GetContainer() ContainerSpec`

GetContainer returns the Container field if non-nil, zero value otherwise.

### GetContainerOk

`func (o *InstanceSpec) GetContainerOk() (*ContainerSpec, bool)`

GetContainerOk returns a tuple with the Container field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContainer

`func (o *InstanceSpec) SetContainer(v ContainerSpec)`

SetContainer sets Container field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


