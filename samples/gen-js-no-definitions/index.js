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

module.exports = class SwaggerTest {

  constructor(options) {
    options = options || {};

    if (options.discovery) {
      try {
        this.address = discovery("swagger-test", "http").url();
      } catch (e) {
        this.address = discovery("swagger-test", "default").url();
      };
    } else if (options.address) {
      this.address = options.address;
    } else {
      throw new Error("Cannot initialize swagger-test without discovery or address");
    }
    if (options.timeout) {
      this.timeout = options.timeout
    }
  }

  deleteBook(id, options, cb) {
    const params = {};
    params["id"] = id;

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
      span.logEvent("DELETE /v1/books/{id}");
    }

    const requestOptions = {
      method: "DELETE",
      uri: this.address + "/v1/books/" + params.id + "",
      json: true,
      timeout,
      headers,
      qs: query,
    };

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
