// Code generated by OpenAPI Generator (https://openapi-generator.tech), manual changes will be lost - read more on https://github.com/algolia/api-clients-automation. DO NOT EDIT.
package search

import (
	"encoding/json"
	"fmt"
)

// Source Source.
type Source struct {
	// IP address range of the source.
	Source string `json:"source"`
	// Source description.
	Description *string `json:"description,omitempty"`
}

type SourceOption func(f *Source)

func WithSourceDescription(val string) SourceOption {
	return func(f *Source) {
		f.Description = &val
	}
}

// NewSource instantiates a new Source object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed.
func NewSource(source string, opts ...SourceOption) *Source {
	this := &Source{}
	this.Source = source
	for _, opt := range opts {
		opt(this)
	}
	return this
}

// NewEmptySource return a pointer to an empty Source object.
func NewEmptySource() *Source {
	return &Source{}
}

// GetSource returns the Source field value.
func (o *Source) GetSource() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Source
}

// GetSourceOk returns a tuple with the Source field value
// and a boolean to check if the value has been set.
func (o *Source) GetSourceOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Source, true
}

// SetSource sets field value.
func (o *Source) SetSource(v string) *Source {
	o.Source = v
	return o
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *Source) GetDescription() string {
	if o == nil || o.Description == nil {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Source) GetDescriptionOk() (*string, bool) {
	if o == nil || o.Description == nil {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *Source) HasDescription() bool {
	if o != nil && o.Description != nil {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *Source) SetDescription(v string) *Source {
	o.Description = &v
	return o
}

func (o Source) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]any{}
	toSerialize["source"] = o.Source
	if o.Description != nil {
		toSerialize["description"] = o.Description
	}
	serialized, err := json.Marshal(toSerialize)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Source: %w", err)
	}

	return serialized, nil
}

func (o Source) String() string {
	out := ""
	out += fmt.Sprintf("  source=%v\n", o.Source)
	out += fmt.Sprintf("  description=%v\n", o.Description)
	return fmt.Sprintf("Source {\n%s}", out)
}
