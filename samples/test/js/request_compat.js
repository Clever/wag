const assert = require("assert");
const http = require("http");

const Client = require("swagger-test");

// These tests use a real HTTP server instead of nock, so that the client sees the same
// connection-level behavior (resets, slow responses, redirects) it would in production.
// They check that the fetch-based client matches the behavior of the request library
// it replaced.
describe("request compatibility", function() {
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

  function expectError(promise) {
    return promise.then(() => {
      throw new Error("expected the request to fail");
    }, (err) => err);
  }

  it("does not follow redirects for POST requests", function() {
    handler = (req, res) => {
      if (req.method === "POST") {
        res.writeHead(302, {location: "/elsewhere"});
        res.end();
        return;
      }
      res.writeHead(200, {"content-type": "application/json"});
      res.end("{}");
    };
    const c = newClient();
    return expectError(c.createBook({name: "x"})).then((err) => {
      assert.equal(err.message, "Received unexpected statusCode 302");
    });
  });

  it("follows redirects for GET requests", function() {
    handler = (req, res) => {
      if (req.url.startsWith("/v1/books")) {
        res.writeHead(302, {location: "/elsewhere"});
        res.end();
        return;
      }
      res.writeHead(200, {"content-type": "application/json"});
      res.end(JSON.stringify([{name: req.url}]));
    };
    const c = newClient();
    return c.getBooks({}).then((books) => {
      assert.deepEqual(books, [{name: "/elsewhere"}]);
    });
  });

  it("sends null header params as the string \"null\"", function() {
    let authorization;
    handler = (req, res) => {
      authorization = req.headers.authorization;
      res.writeHead(200, {"content-type": "application/json"});
      res.end("{}");
    };
    const c = newClient();
    return c.getBookByID({bookID: 1, authorization: null}).then(() => {
      assert.equal(authorization, "null");
    });
  });

  it("returns ECONNRESET \"socket hang up\" when the server closes the connection before responding", function() {
    handler = (req) => req.socket.destroy();
    const c = newClient();
    return expectError(c.getBooks({})).then((err) => {
      assert.equal(err.code, "ECONNRESET");
      assert.equal(err.message, "socket hang up");
    });
  });

  it("returns ECONNRESET \"aborted\" when the server closes the connection mid-body", function() {
    handler = (req, res) => {
      res.writeHead(200, {"content-type": "application/json", "content-length": "100"});
      res.write("[1,2,");
      setTimeout(() => req.socket.destroy(), 20);
    };
    const c = newClient();
    return expectError(c.getBooks({})).then((err) => {
      assert.equal(err.code, "ECONNRESET");
      assert.equal(err.message, "aborted");
    });
  });

  it("returns ESOCKETTIMEDOUT when the server is slow to respond", function() {
    handler = (req, res) => {
      setTimeout(() => {
        res.writeHead(200, {"content-type": "application/json"});
        res.end("[]");
      }, 300);
    };
    const c = newClient({timeout: 50});
    return expectError(c.getBooks({})).then((err) => {
      assert.equal(err.code, "ESOCKETTIMEDOUT");
      assert.equal(err.connect, false);
    });
  });

  it("treats the timeout as an idle timeout while reading the body", function() {
    // The body takes ~300ms in total, but data arrives every 30ms, so a 100ms
    // timeout should not fire.
    handler = (req, res) => {
      res.writeHead(200, {"content-type": "application/json"});
      res.write("[");
      let n = 0;
      const interval = setInterval(() => {
        n++;
        if (n < 10) {
          res.write(n > 1 ? ",{}" : "{}");
        } else {
          clearInterval(interval);
          res.end("]");
        }
      }, 30);
    };
    const c = newClient({timeout: 100});
    return c.getBooks({}).then((books) => {
      assert.equal(books.length, 9);
    });
  });

  it("returns ESOCKETTIMEDOUT when the body stalls for longer than the timeout", function() {
    handler = (req, res) => {
      res.writeHead(200, {"content-type": "application/json", "content-length": "100"});
      res.write("[");
    };
    const c = newClient({timeout: 50});
    return expectError(c.getBooks({})).then((err) => {
      assert.equal(err.code, "ESOCKETTIMEDOUT");
    });
  });
});
