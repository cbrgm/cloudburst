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
	"errors"
)

// InstancesApiService is a service that implents the logic for the InstancesApiServicer
// This service should implement the business logic for every endpoint for the InstancesApi API. 
// Include any external packages or services that will be required by this service.
type InstancesApiService struct {
}

// NewInstancesApiService creates a default api service
func NewInstancesApiService() InstancesApiServicer {
	return &InstancesApiService{}
}

// GetInstances - Get InstanceDemand for a ScrapeTarget
func (s *InstancesApiService) GetInstances(ctx context.Context, target string) (ImplResponse, error) {
	// TODO - update GetInstances with the required logic for this service method.
	// Add api_instances_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []Instance{}) or use other options such as http.Ok ...
	//return Response(200, []Instance{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("GetInstances method not implemented")
}

// SaveInstances - Update InstanceDemand for a ScrapeTarget
func (s *InstancesApiService) SaveInstances(ctx context.Context, target string, instance []Instance) (ImplResponse, error) {
	// TODO - update SaveInstances with the required logic for this service method.
	// Add api_instances_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []Instance{}) or use other options such as http.Ok ...
	//return Response(200, []Instance{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("SaveInstances method not implemented")
}

