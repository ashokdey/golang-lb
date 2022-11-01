# Golang Load Balancer

A simple load balancer implementation in Golang with support of config files and list of microservices.

## Sample YML file 

```yaml

microservices:
  - service:
      properties:
        name: user
        servers: [
          "http://127.0.0.1:4000",
          "http://127.0.0.1:4001",
          "http://127.0.0.1:4002",
          "http://127.0.0.1:4003",
        ]
        version: 1.0
        apiPrefix: users-api
        healthRoute: "/"

  - service:
      properties:
        name: auth
        servers: [
          "http://127.0.0.1:6000",
          "http://127.0.0.1:6001",
        ]
        version: 1.0
        apiPrefix: auth-api
        healthRoute: "/"

```

### Usages 

1. Clone the repo
2. Run `go run ./cmd`

### Notes 

This was inspired from the video by [Arjun Mahishi](https://github.com/arjunmahishi) on creating a LB in Golang