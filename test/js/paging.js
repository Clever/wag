const assert = require("assert");
const nock = require("nock");
const util = require("util");

const Client = require("swagger-test");

const mockAddress = "http://localhost:8000";

describe("paging", function() {
  it("iterators support map", async () => {
    const c = new Client({address: mockAddress});
    const scopeFirst = nock(mockAddress)
      .get("/v1/books")
      .reply(
        200,
        [{id: 1, name: "first"}, {id: 2, name: "second"}],
        {"X-Next-Page-Path": "/v1/books?startingAfter=2"}
      );
    const scopeSecond = nock(mockAddress)
      .get("/v1/books")
      .query({startingAfter: "2"})
      .reply(200, [{id: 3, name: "third"}]);
    const bookNames = await c.getBooksIter({}, {}).map(b => b.name);
    assert.deepEqual(bookNames, ["first", "second", "third"]);
    assert(scopeFirst.isDone());
    assert(scopeSecond.isDone());
  });

  it("iterators support forEach", async () => {
    const c = new Client({address: mockAddress});
    const scopeFirst = nock(mockAddress)
      .get("/v1/books")
      .reply(
        200,
        [{id: 1, name: "first"}, {id: 2, name: "second"}],
        {"X-Next-Page-Path": "/v1/books?startingAfter=2"},
      );
    const scopeSecond = nock(mockAddress)
      .get("/v1/books")
      .query({startingAfter: "2"})
      .reply(200, [{id: 3, name: "third"}]);

    // do a manual map
    const results = [];
    const bookNames = await c.getBooksIter({}, {}).forEach(b => results.push(b.name));
    assert.deepEqual(results, ["first", "second", "third"]);
    assert(scopeFirst.isDone());
    assert(scopeSecond.isDone());
  });

  it("iterators support resource path", async () => {
    const c = new Client({address: mockAddress});
    const scope = nock(mockAddress)
      .get("/v1/authors")
      .reply(
        200,
        {authorSet: {
          results: [{id: "abc", name: "John"}, {id: "def", name: "Julie"}],
        }},
      );
    const authorNames = await c.getAuthorsIter({}, {}).map(a => a.name);
    assert.deepEqual(authorNames, ["John", "Julie"]);
    assert(scope.isDone());
  });

  it("iterators pass through errors", async () => {
    const c = new Client({address: mockAddress});
    const scopeFirst = nock(mockAddress)
      .get("/v1/books")
      .reply(
        200,
        [{id: 1, name: "first"}, {id: 2, name: "second"}],
        {"X-Next-Page-Path": "/v1/books?startingAfter=2"}
      );
    const scopeSecond = nock(mockAddress)
      .get("/v1/books")
      .times(2)
      .query({startingAfter: "2"})
      .reply(500, {message: "fail"});
    try {
      await c.getBooksIter({}, {}).map(b => b.name);
    } catch (e) {
      assert.equal(e.message, "fail");
      assert(scopeFirst.isDone());
      assert(scopeSecond.isDone());
      return;
    }
    assert.fail(null, null, "expected error");
  });

  it("iterators handle callbacks", done => {
    const c = new Client({address: mockAddress});
    const scopeFirst = nock(mockAddress)
      .get("/v1/books")
      .reply(
        200,
        [{id: 1, name: "first"}, {id: 2, name: "second"}],
        {"X-Next-Page-Path": "/v1/books?startingAfter=2"}
      );
    const scopeSecond = nock(mockAddress)
      .get("/v1/books")
      .query({startingAfter: "2"})
      .reply(200, [{id: 3, name: "third"}]);
    c.getBooksIter({}, {}).toArray((err, books) => {
      assert(!err, "unexpected error: " + util.inspect(err));
      assert.deepEqual(
        books,
        [{id: 1, name: "first"}, {id: 2, name: "second"}, {id: 3, name: "third"}],
      );
      assert(scopeFirst.isDone());
      assert(scopeSecond.isDone());
      done();
    });
  });
});
