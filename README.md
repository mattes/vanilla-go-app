# Vanilla Go Server

## Usage

```
$ ./vanilla-go-app -listen ':8080' -http-keep-alive true

# Returns 200
curl http://localhost:8080

# Returns 200 and fixed binary body size
curl http://localhost:8080/bin/0KB
curl http://localhost:8080/bin/1KB
curl http://localhost:8080/bin/10KB
curl http://localhost:8080/bin/100KB
curl http://localhost:8080/bin/1000KB

# Returns 204 and just reads the input body
curl http://localhost:8080/readall -X POST -d 'foobar'

# Returns 201, reads and returns the same body back
curl http://localhost:8080/echo -X POST -d 'foobar'

# Returns 202 and sleeps for some time
curl http://localhost:8080/sleep?ms=500 
```

