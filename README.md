# JumpCloud hashServer

This program accepts a single argument- the port to listen on.  

The following requests will be allowed:

<li>POST /hash - accepts formdata in the payload, looking for a form value of "password".  Hashes the provided value using the Sha512 algorithm.  Returns a unique Id for the request
<li>GET /hash/{hashId} - no payload.  Parameter is the id returned from a previous POST to /hash.  Value will not be availabe until 5 seconds after the original request.  Until 
  then, this call will return a 404.
<li>GET /statistics - will return a JSON-formatted string with the values of the current number of requests processed, as well as the average time to process each request.
<li>POST /shutdown - will shutdown the entire hashServer

There are several improvements I would make if I knew Go a bit better:
1. Improved http request handling.  There appear to be several public packages that help setup handlers with the methods allowed for that handler.  My 
handling of this is admittedly a bit simplistic, but it was not the main focus of my effort.
2. Improved managing of shared data.  A better knowledge of Go's conucrrency mechanisms would probably allow for a cleaner or more efficient implementation.

Build the Dockerfile from the root direcotry of the project - NOT the /build directory: "docker build -f build/Dockerfile ."  The image will export post 8080 in the container.  This can be mapped to any port on the host.

