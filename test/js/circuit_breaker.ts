const assert = require("assert");
const nock = require("nock");

const Client = require("swagger-test");
const { RetryPolicies } = Client;

const mockAddress = "http://localhost:8000";

const { commandFactory, metricsFactory, circuitFactory } = require("hystrixjs");

async function sleep(ms) {
  return new Promise((resolve) => {
    setTimeout(resolve, ms);
  });
};

describe("circuit", function() {
  beforeEach(() => {
    metricsFactory.resetCache();
    circuitFactory.resetCache();
    commandFactory.resetCache();
    nock.cleanAll();
  });

  it("opens when the number of failures exceeds requestVolumeThreshold, and closes if sees success after sleep window", async () => {
    const c = new Client({
      address: mockAddress,
      circuit: {
        forceClosed: false,
        requestVolumeThreshold: 20,
        sleepWindow: 500,
      },
    });

    // trip the circuit
    for (let i = 0; i < 20; i++) {
      const scope = nock(mockAddress).
        get("/v1/authors").
        replyWithError('something awful happened');
      try {
        const resp = await c.getAuthors({});
      } catch (err) {
        assert.equal(err.message, "something awful happened");
      }
      assert(scope.isDone(), "nock scope should be done");
    }

    // see circuit open
    let scope = nock(mockAddress).
      get("/v1/authors").
      replyWithError('something awful happened');
    try {
      const resp = await c.getAuthors({});
      throw new Error("got resp, should have gotten circuit open error");
    } catch (err) {
      assert.equal(err.message, "OpenCircuitError");
      assert.equal(scope.isDone(), false, "nock should not have gotten hit");
    }
    nock.cleanAll();

    // circuit should let something through after the sleep window
    await sleep(750);
    scope = nock(mockAddress).
      get("/v1/authors").
      reply(200);
    try {
      const resp = await c.getAuthors({});
    } catch (err) {
      assert.fail(`Circuit should have let a test request through, instead got error: ${err}`);
    }
    assert.equal(scope.isDone(), true, "nock should have gotten hit");

    // circuit should now be closed and letting requests through
    for (let i = 0; i < 20; i++) {
      const scope = nock(mockAddress).
        get("/v1/authors").
        reply(200);
      try {
        const resp = await c.getAuthors({});
      } catch (err) {
        assert.fail(`Circuit should be closed and letting requests through, but instead we got error: ${err}`);
      }
      assert(scope.isDone(), "nock scope should be done");
    }
  });

  it("doesn't do anything if forceClosed is set to true", async () => {
    const c = new Client({
      address: mockAddress,
      circuit: {
        forceClosed: true,
        requestVolumeThreshold: 20
      },
    });
    for (let i = 0; i < 40; i++) {
      const scope = nock(mockAddress).
        get("/v1/authors").
        replyWithError('something awful happened');
      try {
        const resp = await c.getAuthors({});
        throw new Error("got resp, should have gotten error");
      } catch (err) {
        assert.equal(err.message, "something awful happened");
      }
      assert(scope.isDone(), "nock scope should be done");
    }
  });

  it("logs state of the circuit", async () => {
    let loggerCalls = 0;
    const c = new Client({
      address: mockAddress,
      logger: {
        errorD: (title, data) => { },
        infoD: (title, data) => {
          loggerCalls++;
          if (loggerCalls == 1) {
            assert.equal(data.errorCount, 20, "expected log to show 20 errors");
            assert.equal(data.errorPercentage, 100, "expected error percent to be 100");
          }
        },
      },
      circuit: {
        forceClosed: true,
        requestVolumeThreshold: 20,
        logIntervalMs: 1000,
      },
    });
    for (let i = 0; i < 20; i++) {
      const scope = nock(mockAddress).
        get("/v1/authors").
        replyWithError('something awful happened');
      try {
        const resp = await c.getAuthors({});
        throw new Error("got resp, should have gotten error");
      } catch (err) {
        assert.equal(err.message, "something awful happened");
      }
      assert(scope.isDone(), "nock scope should be done");
    }
    await sleep(1500);
    assert.equal(loggerCalls, 1);
  });

  it("does not consider 4XXs errors", async () => {
    let loggerCalls = 0;
    const c = new Client({
      address: mockAddress,
      retryPolicy: RetryPolicies.None,
      logger: {
        errorD: (title, data) => { },
        infoD: (title, data) => {
          loggerCalls++;
          if (loggerCalls == 1) {
            assert.equal(data.errorCount, 0, "expected log to show 0 errors");
            assert.equal(data.errorPercentage, 0, "expected error percent to be 0");
          }
        },
      },
      circuit: {
        forceClosed: true,
        requestVolumeThreshold: 20,
        logIntervalMs: 1000,
      },
    });
    for (let i = 0; i < 20; i++) {
      const scope = nock(mockAddress).
        get("/v1/books/12345").
        reply(404, `{"message":"Not found"}`);
      try {
        const resp = await c.getBookByID({ bookID: 12345 });
        throw new Error("got resp, should have gotten error");
      } catch (err) {
        assert.equal(err.message, "Not found");
      }
      assert(scope.isDone(), "nock scope should be done");
    }
    await sleep(1500);
    assert.equal(loggerCalls, 1);
  });

  it("applies callback, if provided, to circuit-breaker error", (done) => {
    const c = new Client({
      address: mockAddress,
      circuit: {
        forceClosed: false,
        maxConcurrentRequests: 5,
        requestVolumeThreshold: 5,
      },
      logger: {
        errorD: () => { },
        infoD: () => { },
      },
      retryPolicy: RetryPolicies.None,
    });

    const promises = [];
    for (let i = 0; i < 10; i++) {
      const mockApi = nock(mockAddress).get("/v1/authors").reply(200, { authors: [] });
      promises.push(c.getAuthors({}, (err, data) => {
        if (i < 5) {
          assert.equal(null, err, "Expected no error for first 5 client calls");
          assert.deepEqual([], data.authors, "Expected data for first 5 client calls");
          assert.equal(true, mockApi.isDone(), "Expected API call for first 5 client calls");
        } else {
          assert.equal(
            "CommandRejected", err.message,
            "Expected circuit breaker to short circuit last 5 client calls");
          assert.equal(false, mockApi.isDone(), "Expected no API call for last 5 client calls");
        }
      }));
    }

    Promise.all(promises).then(() => {
      done();
    }).catch((err) => {
      done(err);
    });
  });
});
