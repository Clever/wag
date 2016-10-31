const discovery = require("@clever/discovery");
const request = require("request");
const url = require("url");
const opentracing = require("opentracing");

// go-swagger treats handles/expects arrays in the query string to be a string of comma joined values
// so...do that thing. It's worth noting that this has lots of issues ("what if my values have commas in them?")
// but that's an issue with go-swagger
function serializeQueryString(data) {
  if (Array.isArray(data)) {
    return data.join(",");
  }
  return data;
}

module.exports = class WagPatch {

  constructor(options) {
    options = options || {};

    if (options.discovery) {
      try {
        this.address = discovery("wag-patch", "http").url();
      } catch (e) {
        this.address = discovery("wag-patch", "default").url();
      };
    } else if (options.address) {
      this.address = options.address;
    } else {
      throw new Error("Cannot initialize wag-patch without discovery or address");
    }
    if (options.timeout) {
      this.timeout = options.timeout
    }
  }

  wagpatch(specialPatch, options, cb) {
    const params = {};
    params["specialPatch"] = specialPatch;

    if (!cb && typeof options === "function") {
      cb = options;
      options = undefined;
    }

    if (!options) {
      options = {};
    }

    const timeout = options.timeout || this.timeout;
    const span = options.span;

    const headers = {};

    const query = {};

    if (span) {
      opentracing.inject(span, opentracing.FORMAT_TEXT_MAP, headers);
      span.logEvent("PATCH /wagpatch");
    }

    const requestOptions = {
      method: "PATCH",
      uri: this.address + "/wagpatch",
      json: true,
      timeout,
      headers,
      qs: query,
    };

    requestOptions.body = params.specialPatch;

    return new Promise((resolve, reject) => {
      const rejecter = (err) => {
        reject(err);
        if (cb) {
          cb(err);
        }
      }
      const resolver = (data) => {
        resolve(data);
        if (cb) {
          cb(null, data);
        }
      }

      request(requestOptions, (err, response, body) => {
        if (err) {
          return rejecter(err);
        }
        if (response.statusCode >= 400) {
          return rejecter(new Error(body));
        }
        resolver(body);
      });
    });
  }

  wagpatch2(params, options, cb) {
    if (!cb && typeof options === "function") {
      cb = options;
      options = undefined;
    }

    if (!options) {
      options = {};
    }

    const timeout = options.timeout || this.timeout;
    const span = options.span;

    const headers = {};

    const query = {};
    query["other"] = serializeQueryString(params.other);

    if (span) {
      opentracing.inject(span, opentracing.FORMAT_TEXT_MAP, headers);
      span.logEvent("PATCH /wagpatch2");
    }

    const requestOptions = {
      method: "PATCH",
      uri: this.address + "/wagpatch2",
      json: true,
      timeout,
      headers,
      qs: query,
    };

    requestOptions.body = params.specialPatch;

    return new Promise((resolve, reject) => {
      const rejecter = (err) => {
        reject(err);
        if (cb) {
          cb(err);
        }
      }
      const resolver = (data) => {
        resolve(data);
        if (cb) {
          cb(null, data);
        }
      }

      request(requestOptions, (err, response, body) => {
        if (err) {
          return rejecter(err);
        }
        if (response.statusCode >= 400) {
          return rejecter(new Error(body));
        }
        resolver(body);
      });
    });
  }
}
