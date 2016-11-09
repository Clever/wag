module.exports.Errors = {};

/**
 * ExtendedError
 * @extends Error
 * @memberof module:swagger-test
 * @alias module:swagger-test.Errors.ExtendedError
 * @property {number} code
 * @property {string} message
 */
module.exports.Errors.ExtendedError = class extends Error {
  constructor(body) {
    super(body.message);
    for (const k of Object.keys(body)) {
      this[k] = body[k];
    }
  }
};

/**
 * NotFound
 * @extends Error
 * @memberof module:swagger-test
 * @alias module:swagger-test.Errors.NotFound
 * @property {string} message
 */
module.exports.Errors.NotFound = class extends Error {
  constructor(body) {
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
 * @property {number} code
 * @property {string} message
 */
module.exports.Errors.InternalError = class extends Error {
  constructor(body) {
    super(body.message);
    for (const k of Object.keys(body)) {
      this[k] = body[k];
    }
  }
};

