/*
 * Cloudburst
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 0.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// A InstancesApiController binds http requests to an api service and writes the service results to the http response
type InstancesApiController struct {
	service InstancesApiServicer
}

// NewInstancesApiController creates a default api controller
func NewInstancesApiController(s InstancesApiServicer) Router {
	return &InstancesApiController{ service: s }
}

// Routes returns all of the api route for the InstancesApiController
func (c *InstancesApiController) Routes() Routes {
	return Routes{ 
		{
			"GetInstances",
			strings.ToUpper("Get"),
			"/api/v1/targets/{target}/instances",
			c.GetInstances,
		},
		{
			"UpdateInstances",
			strings.ToUpper("Put"),
			"/api/v1/targets/{target}/instances",
			c.UpdateInstances,
		},
	}
}

// GetInstances - Get Instances for a ScrapeTarget
func (c *InstancesApiController) GetInstances(w http.ResponseWriter, r *http.Request) { 
	params := mux.Vars(r)
	target := params["target"]
	result, err := c.service.GetInstances(r.Context(), target)
	//If an error occured, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	//If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
	
}

// UpdateInstances - Update Instances for a ScrapeTarget
func (c *InstancesApiController) UpdateInstances(w http.ResponseWriter, r *http.Request) { 
	params := mux.Vars(r)
	target := params["target"]
	instance := &[]Instance{}
	if err := json.NewDecoder(r.Body).Decode(&instance); err != nil {
		w.WriteHeader(500)
		return
	}
	
	result, err := c.service.UpdateInstances(r.Context(), target, *instance)
	//If an error occured, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	//If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
	
}
