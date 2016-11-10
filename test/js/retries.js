const assert = require("assert");
const nock = require("nock");
const util = require("util");

const Client = require("swagger-test");
const {RetryPolicies} = Client;

const mockAddress = "http://localhost:8000";

describe("retries", function() {
  it("performs exponential retries when set backoff by default", function(done) {
    const client = new Client({
      address: mockAddress,
      retryPolicy: RetryPolicies.Exponential
    });
    let requestCount = 0;
    let requestTimes = [];
    const scope = nock(mockAddress)
            .get("/v1/books")
            .times(3)
            .reply(function(uri, requestBody) {
              requestTimes.push(new Date());
              requestCount++;
              if (requestCount <= 2) {
                return [500, ""];
              } else {
                return [200, "[]"];
              }
            });
    client.getBooks({}, function(err, books) {
      assert.equal(requestCount, 3);
      assert(!err, "unexpected error: " + util.inspect(err));
      const backoffMS1 = requestTimes[1].getTime() - requestTimes[0].getTime();
      const backoffMS2 = requestTimes[2].getTime() - requestTimes[1].getTime();
      assert(backoffMS1 < 200 && backoffMS1 > 80, "expected first backoff to be roughly 100ms, got " + backoffMS1)
      assert(backoffMS2 < 240 && backoffMS2 > 160, "expected second backoff to be roughly 200ms, got " + backoffMS2);
      scope.done();
      done();
    });
  });

  it("lets the user configure a custom retry policy on the client", function(done) {
    const client = new Client({
      address: mockAddress,
      retryPolicy: RetryPolicies.None
    });
    let requestCount = 0;
    let requestTimes = [];
    const scope = nock(mockAddress)
            .get("/v1/books")
            .times(3)
            .reply(function(uri, requestBody) {
              requestTimes.push(new Date());
              requestCount++;
              if (requestCount <= 2) {
                return [500, ""];
              } else {
                return [200, "[]"];
              }
            });
    client.getBooks({}, function(err, books) {
      assert.equal(requestCount, 1);
      assert(err);
      assert(!scope.isDone());
      nock.cleanAll();
      done();
    });
  });

  it("lets the user configure a custom retry policy on a single request", function(done) {
    const client = new Client({address: mockAddress});
    let requestCount = 0;
    let requestTimes = [];
    const scope = nock(mockAddress)
            .get("/v1/books")
            .times(3)
            .reply(function(uri, requestBody) {
              requestTimes.push(new Date());
              requestCount++;
              if (requestCount <= 2) {
                return [500, ""];
              } else {
                return [200, "[]"];
              }
            });
    client.getBooks({}, {retryPolicy: RetryPolicies.None}, function(err, books) {
      assert.equal(requestCount, 1);
      assert(err);
      assert(!scope.isDone());
      nock.cleanAll();
      done();
    });
  });

  it("does not retry POSTs", function(done) {
    const client = new Client({address: mockAddress});
    let requestCount = 0;
    let requestTimes = [];
    const scope = nock(mockAddress)
            .post("/v1/books")
            .times(3)
            .reply(function(uri, requestBody) {
              requestTimes.push(new Date());
              requestCount++;
              if (requestCount <= 2) {
                return [500, ""];
              } else {
                return [200, "[]"];
              }
            });
    client.createBook({}, function(err, book) {
      assert.equal(requestCount, 1);
      assert(err);
      assert(!scope.isDone());
      nock.cleanAll();
      done();
    });
  });

  it("does not retry network errors", function(done) {
    let errorReceived = null;
    let requestCount = 0;
    const retryPolicy = {
      backoffs() { return [100, 100, 100]; },
      retry(requestOptions, err, respponse, body) {
        errorReceived = err;
        requestCount += 1;
        return false;
      }
    };
    const client = new Client({address: "https://thisshouldnotresolve1234567890.com", retryPolicy});
    client.getBooks({}, function(err, books) {
      assert(err);
      assert.equal(errorReceived, err);
      assert.equal(requestCount, 1);
      done();
    });
  });

});
