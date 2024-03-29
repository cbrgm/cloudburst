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

/**
 * The StaticSpec model module.
 * @module model/StaticSpec
 * @version 0.0.0
 */
class StaticSpec {
    /**
     * Constructs a new <code>StaticSpec</code>.
     * @alias module:model/StaticSpec
     * @param endpoints {Array.<String>} 
     */
    constructor(endpoints) { 
        
        StaticSpec.initialize(this, endpoints);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj, endpoints) { 
        obj['endpoints'] = endpoints;
    }

    /**
     * Constructs a <code>StaticSpec</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/StaticSpec} obj Optional instance to populate.
     * @return {module:model/StaticSpec} The populated <code>StaticSpec</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new StaticSpec();

            if (data.hasOwnProperty('endpoints')) {
                obj['endpoints'] = ApiClient.convertToType(data['endpoints'], ['String']);
            }
        }
        return obj;
    }


}

/**
 * @member {Array.<String>} endpoints
 */
StaticSpec.prototype['endpoints'] = undefined;






export default StaticSpec;

