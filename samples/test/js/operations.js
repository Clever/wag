const assert = require("assert");
const nock = require("nock");

const Client = require("swagger-test");

const mockAddress = "http://localhost:8000";
const alternateMockAddress = "http://localhost:8001";

describe("operations", function() {
  it("properly combines store and options baggage", function(done) {
    const context = new Map();
    context.set("ClientIP", "123.123.123.123");
    context.set("session_id", "1234567890");
    store = new Map([["context", context]])
    const c = new Client({address: mockAddress, asynclocalstore: store});
    const params = {
      bookID: "the-mediocre-gatsby",
    }
    const mockResponse = {hello: "world", headers: {}};
    let baggageHeader = ""
    const scope = nock(mockAddress)
    .get(`/v1/books/${params.bookID}`)
    .reply(200, function(uri, requestBody) {
      baggageHeader = this.req.headers["baggage"];
      return mockResponse;
    });


    baggage = new Map()
    baggage.set("foo", "bar")
    options = { baggage }
    c.getBookByID({bookID: params.bookID}, options, function(err, resp) {
      assert.ifError(err);
      assert.strictEqual(baggageHeader, 'ClientIP=123.123.123.123,session_id=1234567890,foo=bar');
      scope.done();
      done();
    });
  });

  it("properly sets baggage headers", function(done) {

    const c = new Client({address: mockAddress});
    const params = {
      bookID: "the-mediocre-gatsby",
    }
    const mockResponse = {hello: "world", headers: {}};
    const req = {
      ClientIP: "123.123.123.123",
      SessionID: "1234567890"
    }
    let baggageHeader = ""
    const scope = nock(mockAddress)
    .get(`/v1/books/${params.bookID}`)
    .reply(200, function(uri, requestBody) {
      baggageHeader = this.req.headers["baggage"];
      return mockResponse;
    });
    const baggage = new Map();
    baggage.set("ClientIP", req.ClientIP);
    baggage.set("session_id", req.SessionID);

    const options = {
      baggage
    }
    c.getBookByID({bookID: params.bookID}, options, function(err, resp) {
      assert.ifError(err);
      assert.strictEqual(baggageHeader, 'ClientIP=123.123.123.123,session_id=1234567890');
      scope.done();
      done();
    });
  });

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

  it("works with multiple clients as long as a different serviceName is used", function(done) {
    let doneCount = 0;

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
      doneCount++;
      if (doneCount == 2) {
        done();
      }
    });

    const alternateClient = new Client({address: alternateMockAddress, serviceName: "alternate-swagger-test"});
    const alternateScope = nock(alternateMockAddress, {reqheaders: {'authorization': params.authorization}})
      .get(`/v1/books/${params.bookID}`)
      .query({randomBytes: params.randomBytes})
      .reply(200, mockResponse);
    alternateClient.getBookByID(params, function(err, resp) {
      assert(alternateScope.isDone());
      assert.ifError(err);
      assert.deepEqual(resp, mockResponse)
      doneCount++;
      if (doneCount == 2) {
        done();
      }
    });
  });

  it("error on empty path param", function(done) {
    const c = new Client({address: mockAddress});
    const params = {
      bookID: "",
    };
    c.getBookByID(params, function(err, resp) {
      assert.equal(err.toString(), "Error: bookID must be non-empty because it's a path parameter")
      done()
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

  it("sends a version header", function(done) {
    const c = new Client({address: mockAddress});
    const scope = nock(mockAddress)
      .get(`/v1/books`)
      .reply(function(uri, requestBody) {
        assert.equal(Client.Version, this.req.headers[Client.VersionHeader.toLowerCase()]);
        return [200, {}];
        
      });
    c.getBooks({}, {}).then(function() {
      assert(scope.isDone());
      done();
    });
  });
});


