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

// StaticSpec struct for StaticSpec
type StaticSpec struct {
	Endpoints []string `json:"endpoints"`
}

// NewStaticSpec instantiates a new StaticSpec object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewStaticSpec(endpoints []string, ) *StaticSpec {
	this := StaticSpec{}
	this.Endpoints = endpoints
	return &this
}

// NewStaticSpecWithDefaults instantiates a new StaticSpec object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewStaticSpecWithDefaults() *StaticSpec {
	this := StaticSpec{}
	return &this
}

// GetEndpoints returns the Endpoints field value
func (o *StaticSpec) GetEndpoints() []string {
	if o == nil  {
		var ret []string
		return ret
	}

	return o.Endpoints
}

// GetEndpointsOk returns a tuple with the Endpoints field value
// and a boolean to check if the value has been set.
func (o *StaticSpec) GetEndpointsOk() (*[]string, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.Endpoints, true
}

// SetEndpoints sets field value
func (o *StaticSpec) SetEndpoints(v []string) {
	o.Endpoints = v
}

func (o StaticSpec) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["endpoints"] = o.Endpoints
	}
	return json.Marshal(toSerialize)
}

type NullableStaticSpec struct {
	value *StaticSpec
	isSet bool
}

func (v NullableStaticSpec) Get() *StaticSpec {
	return v.value
}

func (v *NullableStaticSpec) Set(val *StaticSpec) {
	v.value = val
	v.isSet = true
}

func (v NullableStaticSpec) IsSet() bool {
	return v.isSet
}

func (v *NullableStaticSpec) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableStaticSpec(val *StaticSpec) *NullableStaticSpec {
	return &NullableStaticSpec{value: val, isSet: true}
}

func (v NullableStaticSpec) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableStaticSpec) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


