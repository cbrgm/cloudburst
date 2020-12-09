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

// ScrapeTarget struct for ScrapeTarget
type ScrapeTarget struct {
	Name *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Query *string `json:"query,omitempty"`
	InstanceSpec *InstanceSpec `json:"instanceSpec,omitempty"`
}

// NewScrapeTarget instantiates a new ScrapeTarget object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewScrapeTarget() *ScrapeTarget {
	this := ScrapeTarget{}
	return &this
}

// NewScrapeTargetWithDefaults instantiates a new ScrapeTarget object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewScrapeTargetWithDefaults() *ScrapeTarget {
	this := ScrapeTarget{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *ScrapeTarget) GetName() string {
	if o == nil || o.Name == nil {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ScrapeTarget) GetNameOk() (*string, bool) {
	if o == nil || o.Name == nil {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *ScrapeTarget) HasName() bool {
	if o != nil && o.Name != nil {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *ScrapeTarget) SetName(v string) {
	o.Name = &v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *ScrapeTarget) GetDescription() string {
	if o == nil || o.Description == nil {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ScrapeTarget) GetDescriptionOk() (*string, bool) {
	if o == nil || o.Description == nil {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *ScrapeTarget) HasDescription() bool {
	if o != nil && o.Description != nil {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *ScrapeTarget) SetDescription(v string) {
	o.Description = &v
}

// GetQuery returns the Query field value if set, zero value otherwise.
func (o *ScrapeTarget) GetQuery() string {
	if o == nil || o.Query == nil {
		var ret string
		return ret
	}
	return *o.Query
}

// GetQueryOk returns a tuple with the Query field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ScrapeTarget) GetQueryOk() (*string, bool) {
	if o == nil || o.Query == nil {
		return nil, false
	}
	return o.Query, true
}

// HasQuery returns a boolean if a field has been set.
func (o *ScrapeTarget) HasQuery() bool {
	if o != nil && o.Query != nil {
		return true
	}

	return false
}

// SetQuery gets a reference to the given string and assigns it to the Query field.
func (o *ScrapeTarget) SetQuery(v string) {
	o.Query = &v
}

// GetInstanceSpec returns the InstanceSpec field value if set, zero value otherwise.
func (o *ScrapeTarget) GetInstanceSpec() InstanceSpec {
	if o == nil || o.InstanceSpec == nil {
		var ret InstanceSpec
		return ret
	}
	return *o.InstanceSpec
}

// GetInstanceSpecOk returns a tuple with the InstanceSpec field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ScrapeTarget) GetInstanceSpecOk() (*InstanceSpec, bool) {
	if o == nil || o.InstanceSpec == nil {
		return nil, false
	}
	return o.InstanceSpec, true
}

// HasInstanceSpec returns a boolean if a field has been set.
func (o *ScrapeTarget) HasInstanceSpec() bool {
	if o != nil && o.InstanceSpec != nil {
		return true
	}

	return false
}

// SetInstanceSpec gets a reference to the given InstanceSpec and assigns it to the InstanceSpec field.
func (o *ScrapeTarget) SetInstanceSpec(v InstanceSpec) {
	o.InstanceSpec = &v
}

func (o ScrapeTarget) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Name != nil {
		toSerialize["name"] = o.Name
	}
	if o.Description != nil {
		toSerialize["description"] = o.Description
	}
	if o.Query != nil {
		toSerialize["query"] = o.Query
	}
	if o.InstanceSpec != nil {
		toSerialize["instanceSpec"] = o.InstanceSpec
	}
	return json.Marshal(toSerialize)
}

type NullableScrapeTarget struct {
	value *ScrapeTarget
	isSet bool
}

func (v NullableScrapeTarget) Get() *ScrapeTarget {
	return v.value
}

func (v *NullableScrapeTarget) Set(val *ScrapeTarget) {
	v.value = val
	v.isSet = true
}

func (v NullableScrapeTarget) IsSet() bool {
	return v.isSet
}

func (v *NullableScrapeTarget) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableScrapeTarget(val *ScrapeTarget) *NullableScrapeTarget {
	return &NullableScrapeTarget{value: val, isSet: true}
}

func (v NullableScrapeTarget) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableScrapeTarget) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


