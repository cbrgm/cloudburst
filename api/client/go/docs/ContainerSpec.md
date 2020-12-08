# ContainerSpec

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** |  | [optional] 
**Image** | Pointer to **string** |  | [optional] 

## Methods

### NewContainerSpec

`func NewContainerSpec() *ContainerSpec`

NewContainerSpec instantiates a new ContainerSpec object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewContainerSpecWithDefaults

`func NewContainerSpecWithDefaults() *ContainerSpec`

NewContainerSpecWithDefaults instantiates a new ContainerSpec object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ContainerSpec) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ContainerSpec) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ContainerSpec) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ContainerSpec) HasName() bool`

HasName returns a boolean if a field has been set.

### GetImage

`func (o *ContainerSpec) GetImage() string`

GetImage returns the Image field if non-nil, zero value otherwise.

### GetImageOk

`func (o *ContainerSpec) GetImageOk() (*string, bool)`

GetImageOk returns a tuple with the Image field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImage

`func (o *ContainerSpec) SetImage(v string)`

SetImage sets Image field to given value.

### HasImage

`func (o *ContainerSpec) HasImage() bool`

HasImage returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


