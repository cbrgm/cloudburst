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
import ScrapeTarget from '../model/ScrapeTarget';

/**
* Targets service.
* @module api/TargetsApi
* @version 0.0.0
*/
export default class TargetsApi {

    /**
    * Constructs a new TargetsApi. 
    * @alias module:api/TargetsApi
    * @class
    * @param {module:ApiClient} [apiClient] Optional API client implementation to use,
    * default to {@link module:ApiClient#instance} if unspecified.
    */
    constructor(apiClient) {
        this.apiClient = apiClient || ApiClient.instance;
    }



    /**
     * List ScrapeTargets
     * @return {Promise} a {@link https://www.promisejs.org/|Promise}, with an object containing data of type {@link Array.<module:model/ScrapeTarget>} and HTTP response
     */
    listScrapeTargetsWithHttpInfo() {
      let postBody = null;

      let pathParams = {
      };
      let queryParams = {
      };
      let headerParams = {
      };
      let formParams = {
      };

      let authNames = [];
      let contentTypes = [];
      let accepts = ['application/json'];
      let returnType = [ScrapeTarget];
      return this.apiClient.callApi(
        '/targets', 'GET',
        pathParams, queryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, null
      );
    }

    /**
     * List ScrapeTargets
     * @return {Promise} a {@link https://www.promisejs.org/|Promise}, with data of type {@link Array.<module:model/ScrapeTarget>}
     */
    listScrapeTargets() {
      return this.listScrapeTargetsWithHttpInfo()
        .then(function(response_and_data) {
          return response_and_data.data;
        });
    }


}
