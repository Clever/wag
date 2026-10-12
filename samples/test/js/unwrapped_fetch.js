const assert = require("assert");
const http = require("http");

const { Client, unwrappedFetch } = require("./client");

// These tests cover the behavior of the client generated with x-unwrapped-fetch that differs
// from the default client, using a real HTTP server so the client sees real connection behavior.
(unwrappedFetch ? describe : describe.skip)("unwrapped fetch", function() {
  let server;
  let address;
  let handler;

  before(function(done) {
    server = http.createServer((req, res) => handler(req, res));
    server.listen(0, "127.0.0.1", () => {
      address = `http://127.0.0.1:${server.address().port}`;
      done();
    });
  });

  after(function(done) {
    // fetch keeps idle connections alive, so close them for server.close to finish.
    server.closeAllConnections();
    server.close(done);
  });

  let clients = [];
  afterEach(function() {
    clients.forEach((c) => c.close());
    clients = [];
  });

  function newClient(options) {
    const c = new Client(Object.assign({address, retryPolicy: Client.RetryPolicies.None}, options));
    clients.push(c);
    return c;
  }

  function sendJSON(res, status, body) {
    res.writeHead(status, {"Content-Type": "application/json"});
    res.end(JSON.stringify(body));
  }

  it("URI-encodes path parameters", async function() {
    let requestURL;
    handler = (req, res) => {
      requestURL = req.url;
      sendJSON(res, 200, {id: 1});
    };
    await newClient().getBookByID2("a/b c?");
    assert.equal(requestURL, "/v1/books2/a%2Fb%20c%3F");
  });

  it("sends array query parameters as repeated parameters and skips undefined parameters", async function() {
    let requestURL;
    handler = (req, res) => {
      requestURL = req.url;
      sendJSON(res, 200, []);
    };
    await newClient().getBooks({authors: ["Twain", "Austen"], available: false, state: undefined});
    assert.equal(requestURL, "/v1/books?authors=Twain&authors=Austen&available=false");
  });

  it("does not send undefined header parameters", async function() {
    let headers;
    handler = (req, res) => {
      headers = req.headers;
      sendJSON(res, 200, []);
    };
    await newClient().getBooks({});
    assert(!("authorization" in headers), "unexpected authorization header: " + headers.authorization);
  });

  it("passes the fetch Response and body to the retry policy", async function() {
    const calls = [];
    const retryPolicy = {
      backoffs() { return [10]; },
      retry(request, err, res, body) {
        calls.push({request, err, res, body});
        return res.status >= 500;
      },
    };
    let requestCount = 0;
    handler = (req, res) => {
      requestCount++;
      if (requestCount === 1) {
        sendJSON(res, 503, {message: "try again"});
      } else {
        sendJSON(res, 200, {id: 1});
      }
    };
    const book = await newClient({retryPolicy}).getBookByID2("abc");
    assert.deepEqual(book, {id: 1});
    assert.equal(calls.length, 1);
    assert.equal(calls[0].request.method, "GET");
    assert.equal(calls[0].request.url, `${address}/v1/books2/abc`);
    assert.equal(calls[0].err, undefined);
    assert(calls[0].res instanceof Response);
    assert.equal(calls[0].res.status, 503);
    assert.equal(calls[0].body, `{"message":"try again"}`);
  });

  it("times out while reading a slow response body", async function() {
    handler = (req, res) => {
      res.writeHead(200, {"Content-Type": "application/json"});
      res.write("[");
      // Never finish the response.
    };
    await assert.rejects(newClient({timeout: 100}).getBooks({}), (err) => {
      assert.equal(err.name, "TimeoutError");
      assert.equal(err._fromRequest, true);
      return true;
    });
  });

  it("rejects with fetch's error when the server cannot be reached", async function() {
    const closed = http.createServer();
    await new Promise((resolve) => closed.listen(0, "127.0.0.1", resolve));
    const closedAddress = `http://127.0.0.1:${closed.address().port}`;
    await new Promise((resolve) => closed.close(resolve));

    await assert.rejects(newClient({address: closedAddress}).getBooks({}), (err) => {
      assert(err instanceof TypeError);
      assert.equal(err.cause.code, "ECONNREFUSED");
      assert.equal(err._fromRequest, true);
      return true;
    });
  });

  it("follows redirects", async function() {
    handler = (req, res) => {
      if (req.url === "/v1/books2/old") {
        res.writeHead(302, {"Location": "/v1/books2/new"});
        res.end();
      } else {
        sendJSON(res, 200, {id: 2});
      }
    };
    assert.deepEqual(await newClient().getBookByID2("old"), {id: 2});
  });

  it("rejects with a SyntaxError when a response body is not valid JSON", async function() {
    handler = (req, res) => {
      res.writeHead(200, {"Content-Type": "text/html"});
      res.end("<html></html>");
    };
    await assert.rejects(newClient().getBookByID2("abc"), SyntaxError);
  });
});
