/*
 * Cloudburst
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 0.0.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
)

// ProviderSpec struct for ProviderSpec
type ProviderSpec struct {
	Weights map[string]float32 `json:"weights"`
}

// NewProviderSpec instantiates a new ProviderSpec object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewProviderSpec(weights map[string]float32, ) *ProviderSpec {
	this := ProviderSpec{}
	this.Weights = weights
	return &this
}

// NewProviderSpecWithDefaults instantiates a new ProviderSpec object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewProviderSpecWithDefaults() *ProviderSpec {
	this := ProviderSpec{}
	return &this
}

// GetWeights returns the Weights field value
func (o *ProviderSpec) GetWeights() map[string]float32 {
	if o == nil  {
		var ret map[string]float32
		return ret
	}

	return o.Weights
}

// GetWeightsOk returns a tuple with the Weights field value
// and a boolean to check if the value has been set.
func (o *ProviderSpec) GetWeightsOk() (*map[string]float32, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.Weights, true
}

// SetWeights sets field value
func (o *ProviderSpec) SetWeights(v map[string]float32) {
	o.Weights = v
}

func (o ProviderSpec) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["weights"] = o.Weights
	}
	return json.Marshal(toSerialize)
}

type NullableProviderSpec struct {
	value *ProviderSpec
	isSet bool
}

func (v NullableProviderSpec) Get() *ProviderSpec {
	return v.value
}

func (v *NullableProviderSpec) Set(val *ProviderSpec) {
	v.value = val
	v.isSet = true
}

func (v NullableProviderSpec) IsSet() bool {
	return v.isSet
}

func (v *NullableProviderSpec) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableProviderSpec(val *ProviderSpec) *NullableProviderSpec {
	return &NullableProviderSpec{value: val, isSet: true}
}

func (v NullableProviderSpec) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableProviderSpec) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


