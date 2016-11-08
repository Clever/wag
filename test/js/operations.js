const assert = require("assert");
const nock = require("nock");

const Client = require("swagger-test");

const mockAddress = "http://localhost:8000";

describe("operations", function() {
  it("correctly fill out path, header and query params", function(done) {
    const c = new Client({address: mockAddress});
    const params = {
      authorization: "Bearer YOGI",
      bookID: "the-mediocre-gatsby",
      randomBytes: "4" // Very random
    };
    const mockResponse = {hello: "world"};
    const scope = nock(mockAddress, {reqheaders: {'authorization': params.authorization}})
      .get(`/v1/books/${params.bookID}`)
      .query({randomBytes: params.randomBytes})
      .reply(200, mockResponse);
    c.getBookByID(params, function(err, resp) {
      assert(scope.isDone());
      assert.ifError(err);
      assert.deepEqual(resp, mockResponse)
      done();
    });
  });

  it("return a error in failure cases", function(done) {
    const c = new Client({address: mockAddress});
    const scope = nock(mockAddress)
      .get(`/v1/books`)
      .reply(400, {problems: 99, message: "hit me"});
    c.getBooks({}, {}, function(err, resp) {
      assert.equal(resp, undefined);
      assert.equal(err.message, "hit me");
      assert.equal(err.problems, 99);
      done();
    });
  });

  it("return a promise", function(done) {
    const c = new Client({address: mockAddress});
    const scope = nock(mockAddress)
      .get(`/v1/books`)
      .reply(200);
    c.getBooks({}, {}).then(function() {
      assert(scope.isDone());
      done();
    });
  });

  it("return a error to the promise", function(done) {
    const c = new Client({address: mockAddress});
    const scope = nock(mockAddress)
      .get(`/v1/books`)
      .reply(400, {problems: 99, message: "hit me"});
    c.getBooks({}, {}).then(function() {
      assert(false, "then callback should not have been called");
    }).catch(function(err) {
      assert.equal(err.message, "hit me");
      assert.equal(err.problems, 99);
      done();
    });
  });
});
