const assert = require("assert");

const Client = require("swagger-test");

const mockAddress = "http://localhost:8000";
const alternateMockAddress = "http://localhost:8001";

describe("initialization", function() {
  it("fails if not given address or discovery", function() {
    assert.throws(function() {
      new Client();
    }, "Cannot initialize swagger-test without discovery or address");
  });

  it("succeeds given `address`", function() {
    const c = new Client({ address: mockAddress });
    assert.equal(c.address, mockAddress);
  });

  describe("given `discovery`", function() {
    beforeEach(function() {
      delete process.env.SERVICE_SWAGGER_TEST_HTTP_PROTO;
      delete process.env.SERVICE_SWAGGER_TEST_HTTP_HOST;
      delete process.env.SERVICE_SWAGGER_TEST_HTTP_PORT;

      delete process.env.SERVICE_SWAGGER_TEST_DEFAULT_PROTO;
      delete process.env.SERVICE_SWAGGER_TEST_DEFAULT_HOST;
      delete process.env.SERVICE_SWAGGER_TEST_DEFAULT_PORT;

      delete process.env.SERVICE_ALTERNATE_SWAGGER_TEST_DEFAULT_PROTO;
      delete process.env.SERVICE_ALTERNATE_SWAGGER_TEST_DEFAULT_HOST;
      delete process.env.SERVICE_ALTERNATE_SWAGGER_TEST_DEFAULT_PORT;
    });

    it("fails with no env vars", function() {
      assert.throws(function() {
        new Client({discovery: true});
      }, "Missing env var SERVICE_SWAGGER_TEST_DEFAULT_PROTO");
    });

    it("succeeds with default expose", function() {
      process.env.SERVICE_SWAGGER_TEST_DEFAULT_PROTO = "http"
      process.env.SERVICE_SWAGGER_TEST_DEFAULT_HOST = "localhost"
      process.env.SERVICE_SWAGGER_TEST_DEFAULT_PORT = "8000"
      const c = new Client({ discovery: true });
      assert.equal(c.address, mockAddress);
    });

    it("succeeds with legacy `http` expose", function() {
      process.env.SERVICE_SWAGGER_TEST_HTTP_PROTO = "http"
      process.env.SERVICE_SWAGGER_TEST_HTTP_HOST = "localhost"
      process.env.SERVICE_SWAGGER_TEST_HTTP_PORT = "8000"
      const c = new Client({ discovery: true });
      assert.equal(c.address, mockAddress);
    });

    it("uses serviceName when one is passed in", function() {
      process.env.SERVICE_ALTERNATE_SWAGGER_TEST_HTTP_PROTO = "http"
      process.env.SERVICE_ALTERNATE_SWAGGER_TEST_HTTP_HOST = "localhost"
      process.env.SERVICE_ALTERNATE_SWAGGER_TEST_HTTP_PORT = "8001"
      const c = new Client({ discovery: true, serviceName: "alternate-swagger-test" });
      assert.equal(c.address, alternateMockAddress);
    })
  });
});
