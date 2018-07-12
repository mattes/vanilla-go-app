# Vanilla Go App [![Build Status](https://travis-ci.org/templarbit/vanilla-go-app.svg?branch=master)](https://travis-ci.org/templarbit/vanilla-go-app)

Potentially useful for establishing baseline benchmarks.
Supports all the HTTP, TCP timeouts and keep alives.

## Usage

```
./vanilla-go-app -listen :8080

# Returns 200
curl http://localhost:8080

# Returns 200 and fixed size binary body
curl http://localhost:8080/bin/0KB
curl http://localhost:8080/bin/1KB
curl http://localhost:8080/bin/10KB
curl http://localhost:8080/bin/100KB
curl http://localhost:8080/bin/1000KB

# Read and discard the body, returns 204
curl http://localhost:8080/readall -X POST -d 'foobar'

# Reads the body and returns it with status 201
curl http://localhost:8080/echo -X POST -d 'foobar'

# Sleeps for some specified time and returns 202
curl http://localhost:8080/sleep?ms=500 
```

