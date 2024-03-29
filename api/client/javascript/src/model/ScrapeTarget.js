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

import ApiClient from '../ApiClient';
import InstanceSpec from './InstanceSpec';
import ProviderSpec from './ProviderSpec';
import StaticSpec from './StaticSpec';

/**
 * The ScrapeTarget model module.
 * @module model/ScrapeTarget
 * @version 0.0.0
 */
class ScrapeTarget {
    /**
     * Constructs a new <code>ScrapeTarget</code>.
     * @alias module:model/ScrapeTarget
     * @param name {String} 
     * @param path {String} 
     * @param description {String} 
     * @param query {String} 
     * @param providerSpec {module:model/ProviderSpec} 
     * @param instanceSpec {module:model/InstanceSpec} 
     * @param staticSpec {module:model/StaticSpec} 
     */
    constructor(name, path, description, query, providerSpec, instanceSpec, staticSpec) { 
        
        ScrapeTarget.initialize(this, name, path, description, query, providerSpec, instanceSpec, staticSpec);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj, name, path, description, query, providerSpec, instanceSpec, staticSpec) { 
        obj['name'] = name;
        obj['path'] = path;
        obj['description'] = description;
        obj['query'] = query;
        obj['providerSpec'] = providerSpec;
        obj['instanceSpec'] = instanceSpec;
        obj['staticSpec'] = staticSpec;
    }

    /**
     * Constructs a <code>ScrapeTarget</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/ScrapeTarget} obj Optional instance to populate.
     * @return {module:model/ScrapeTarget} The populated <code>ScrapeTarget</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new ScrapeTarget();

            if (data.hasOwnProperty('name')) {
                obj['name'] = ApiClient.convertToType(data['name'], 'String');
            }
            if (data.hasOwnProperty('path')) {
                obj['path'] = ApiClient.convertToType(data['path'], 'String');
            }
            if (data.hasOwnProperty('description')) {
                obj['description'] = ApiClient.convertToType(data['description'], 'String');
            }
            if (data.hasOwnProperty('query')) {
                obj['query'] = ApiClient.convertToType(data['query'], 'String');
            }
            if (data.hasOwnProperty('providerSpec')) {
                obj['providerSpec'] = ProviderSpec.constructFromObject(data['providerSpec']);
            }
            if (data.hasOwnProperty('instanceSpec')) {
                obj['instanceSpec'] = InstanceSpec.constructFromObject(data['instanceSpec']);
            }
            if (data.hasOwnProperty('staticSpec')) {
                obj['staticSpec'] = StaticSpec.constructFromObject(data['staticSpec']);
            }
        }
        return obj;
    }


}

/**
 * @member {String} name
 */
ScrapeTarget.prototype['name'] = undefined;

/**
 * @member {String} path
 */
ScrapeTarget.prototype['path'] = undefined;

/**
 * @member {String} description
 */
ScrapeTarget.prototype['description'] = undefined;

/**
 * @member {String} query
 */
ScrapeTarget.prototype['query'] = undefined;

/**
 * @member {module:model/ProviderSpec} providerSpec
 */
ScrapeTarget.prototype['providerSpec'] = undefined;

/**
 * @member {module:model/InstanceSpec} instanceSpec
 */
ScrapeTarget.prototype['instanceSpec'] = undefined;

/**
 * @member {module:model/StaticSpec} staticSpec
 */
ScrapeTarget.prototype['staticSpec'] = undefined;






export default ScrapeTarget;

