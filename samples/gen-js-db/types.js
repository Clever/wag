module.exports.Errors = {};

/**
 * BadRequest
 * @extends Error
 * @memberof module:swagger-test
 * @alias module:swagger-test.Errors.BadRequest
 * @property {string} message
 */
module.exports.Errors.BadRequest = class extends Error {
  constructor(body = {}) {
    super(body.message);
    for (const k of Object.keys(body)) {
      this[k] = body[k];
    }
  }
};

/**
 * InternalError
 * @extends Error
 * @memberof module:swagger-test
 * @alias module:swagger-test.Errors.InternalError
 * @property {string} message
 */
module.exports.Errors.InternalError = class extends Error {
  constructor(body = {}) {
    super(body.message);
    for (const k of Object.keys(body)) {
      this[k] = body[k];
    }
  }
};

