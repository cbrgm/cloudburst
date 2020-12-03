//
// my-functions.js
//
module.exports = {
    beforeRequest: beforeRequest
}

function beforeRequest(requestParams, context, ee, next) {
    var body = {
        numbers: randomArray(5000)
    };
    requestParams.json = body;
    return next();
}

function randomArray(size) {
    return Array.from({length: size}, () => randomNumber(1, 100000));
}

function randomNumber(min, max) {
    return Math.floor(Math.random() * (max - min) + min);
}
