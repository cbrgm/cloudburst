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
	"context"
	"net/http"
)



// InstancesApiRouter defines the required methods for binding the api requests to a responses for the InstancesApi
// The InstancesApiRouter implementation should parse necessary information from the http request, 
// pass the data to a InstancesApiServicer to perform the required actions, then write the service results to the http response.
type InstancesApiRouter interface { 
	GetInstances(http.ResponseWriter, *http.Request)
	SaveInstances(http.ResponseWriter, *http.Request)
}
// TargetsApiRouter defines the required methods for binding the api requests to a responses for the TargetsApi
// The TargetsApiRouter implementation should parse necessary information from the http request, 
// pass the data to a TargetsApiServicer to perform the required actions, then write the service results to the http response.
type TargetsApiRouter interface { 
	ListScrapeTargets(http.ResponseWriter, *http.Request)
}


// InstancesApiServicer defines the api actions for the InstancesApi service
// This interface intended to stay up to date with the openapi yaml used to generate it, 
// while the service implementation can ignored with the .openapi-generator-ignore file 
// and updated with the logic required for the API.
type InstancesApiServicer interface { 
	GetInstances(context.Context, string) (ImplResponse, error)
	SaveInstances(context.Context, string, []Instance) (ImplResponse, error)
}


// TargetsApiServicer defines the api actions for the TargetsApi service
// This interface intended to stay up to date with the openapi yaml used to generate it, 
// while the service implementation can ignored with the .openapi-generator-ignore file 
// and updated with the logic required for the API.
type TargetsApiServicer interface { 
	ListScrapeTargets(context.Context) (ImplResponse, error)
}
