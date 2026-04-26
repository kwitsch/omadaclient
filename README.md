# Golang client for Omada Rest Api

Golang library to interface the [Omada](https://www.tp-link.com/de/omada-sdn/) controller API.

---
This library is work in progress and only contains a limited API subset.  
New methods will be added as required for other projects.

## Omada Open API (generated client)

The `openapi/` package contains a fully generated Go client for the official
[Omada Open API](https://community.tp-link.com/en/business/forum/topic/590430).  
The client is generated from the OpenAPI 3.0 specification at `api/openapi.yaml`
using [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen).

### Regenerate the client

```bash
make generate
```

This runs `oapi-codegen` against `api/openapi.yaml` and writes the result to
`openapi/openapi.gen.go`.

### Updating the API spec

Edit `api/openapi.yaml` and run `make generate` to regenerate the client.

---
