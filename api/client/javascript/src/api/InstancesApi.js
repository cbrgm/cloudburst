/**
 * Cloudburst
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 0.0.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 *
 */


import ApiClient from "../ApiClient";
import Instance from '../model/Instance';

/**
* Instances service.
* @module api/InstancesApi
* @version 0.0.0
*/
export default class InstancesApi {

    /**
    * Constructs a new InstancesApi. 
    * @alias module:api/InstancesApi
    * @class
    * @param {module:ApiClient} [apiClient] Optional API client implementation to use,
    * default to {@link module:ApiClient#instance} if unspecified.
    */
    constructor(apiClient) {
        this.apiClient = apiClient || ApiClient.instance;
    }



    /**
     * Update Instances for a ScrapeTarget
     * @param {String} name 
     * @param {Array.<module:model/Instance>} instance 
     * @return {Promise} a {@link https://www.promisejs.org/|Promise}, with an object containing data of type {@link Array.<module:model/Instance>} and HTTP response
     */
    updateInstancesWithHttpInfo(name, instance) {
      let postBody = instance;
      // verify the required parameter 'name' is set
      if (name === undefined || name === null) {
        throw new Error("Missing the required parameter 'name' when calling updateInstances");
      }
      // verify the required parameter 'instance' is set
      if (instance === undefined || instance === null) {
        throw new Error("Missing the required parameter 'instance' when calling updateInstances");
      }

      let pathParams = {
        'name': name
      };
      let queryParams = {
      };
      let headerParams = {
      };
      let formParams = {
      };

      let authNames = [];
      let contentTypes = ['application/json'];
      let accepts = ['application/json'];
      let returnType = [Instance];
      return this.apiClient.callApi(
        '/targets/{name}/instances', 'PUT',
        pathParams, queryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, null
      );
    }

    /**
     * Update Instances for a ScrapeTarget
     * @param {String} name 
     * @param {Array.<module:model/Instance>} instance 
     * @return {Promise} a {@link https://www.promisejs.org/|Promise}, with data of type {@link Array.<module:model/Instance>}
     */
    updateInstances(name, instance) {
      return this.updateInstancesWithHttpInfo(name, instance)
        .then(function(response_and_data) {
          return response_and_data.data;
        });
    }


}
