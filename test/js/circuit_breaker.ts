const assert = require("assert");
const nock = require("nock");

const Client = require("swagger-test");

const mockAddress = "http://localhost:8000";

const {commandFactory,metricsFactory,circuitFactory} = require("hystrixjs");

async function sleep(ms) { return new Promise((resolve) => {
  setTimeout(resolve, ms);
});};

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
        debug: false,
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

  it("doesn't do anything if debug is set to true", async () => {
    const c = new Client({
      address: mockAddress,
      circuit: {
        debug: true,
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
      circuit: {
        debug: true,
        requestVolumeThreshold: 20,
        logIntervalMs: 1000,
        logger: (data) => {
          loggerCalls++;
          if (loggerCalls == 1) {
            assert.equal(data.errorCount, 20, "expected log to show 20 errors");
            assert.equal(data.errorPercentage, 100, "expected error percent to be 100");
          }
        },
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
});

