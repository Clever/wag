const assert = require("assert");
const nock = require("nock");
const opentracing = require('opentracing');
const Client = require("swagger-test");
const sinon = require("sinon");

const mockAddress = "http://localhost:8000";

describe("tracing", function() {
  it("injects span metadata into request headers", function(done) {
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
      .reply(function(uri, requestBody) {
        // ideally we inspect request headers here and assert that the opentracing
        // headers are set, but both the noop tracer and the mock tracer in
        // the opentracing lib have no-op inject()s. For now we assert that
        // inject() was called.
        return [200, mockResponse, {}]
      });
    const tracer = new opentracing.Tracer();
    opentracing.initGlobalTracer(tracer);
    // from reading the source, _inject appears to be the method to override for
    // custom behavior/spying in test (same goes for span methods below)
    tracer._inject = sinon.spy();
    const span = tracer.startSpan('foo_span');
    span._log = sinon.spy();
    span._addTags = sinon.spy();
    c.getBookByID(params, {span}, function(err, resp) {
      assert(scope.isDone());
      assert.ifError(err);
      assert.deepEqual(resp, mockResponse)
      assert(tracer._inject.calledOnce, "expected inject() to be called once");
      assert.deepEqual(tracer._inject.args[0][0], span.context(), "expected inject() to be called on span");
      assert(span._log.calledOnce, "expected to log one event");
      assert.deepEqual(span._log.args[0][0], {event: 'GET /v1/books/{book_id}', payload: undefined});
      assert(span._addTags.calledOnce, "expected to setTag to be called once");
      assert.deepEqual(span._addTags.args[0][0], {'span.kind': 'client'});
      done();
    });
  });
});
