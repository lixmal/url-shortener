# url-shortener

## Description
This is a simple url shortener web service written in Go.

It provides an API to shorten URLs and also to look up the original URLs based on
the given identifier.
Invalid requests will return HTTP `400`, requests unknown identifiers will return HTTP `404`. If there is a body returned it will contain a JSON document.

See [API](#API) for a more detailed description.

Requests to shorten URLs that are already in the database will receive a new identifier regardless.
This will allow us to separate expiration for distinct users in a future release.

Users can enter the shortened link in their browsers to be redirected to the orignal URL.


## Usage

### Build
    go build -o url-shortener github.com/lixmal/url-shortener/cmd


### Test
    go test -v ./...


### Run
    GIN_MODE=release ./url-shortener

## API

### Shortening a URL

Request: `POST` to `∕api∕v1∕shorten∕<url to shorten>`

Response success:
* Code: `200`
* Body: JSON `{"shortened_url":"http://host:port/<identifier>"}`

Response on malformed request:
* Code: `400`
* Body: JSON `{"error":"<error message>"}`

Response server error:
* Code: `500`
* Body: JSON `{"error":"<error message>"}`

### Looking up a URL

Request: `GET` to `∕api∕v1∕lookup∕<identifier>`

Response success:
* Code: `200`
* Body: JSON `{"url":"<original url"}`

Response on malformed request:
* Code: `400`
* Body: JSON `{"error":"<error message>"}`

Response on unknown identifier:
* Code: `404`
* Body: JSON `{"error":"given identifier was not found"}`

### Accessing the original url through the shortened url

Request: `GET` to `/<identifier>`

Response success:
* Code: `307`
* Header `Location: <original url>`

Response on malformed request:
* Code: `400`
* Body: JSON `{"error":"<error message>"}`

Response on unknown identifier:
* Code: `404`


## Limitations

* The service currently stores all url <-> identifier mappings in memory,
therefore all data is lost once the service is shut down.

* Created shortened urls are valid forever. A time or consumption based deletion facility
could be added in a future release.

* No form of authz, authn or transport encryption is currently supported.

* Identfiers are not cryptographically secure (= can be brute forced), don't create shortened urls with the mindset that they are secret.

