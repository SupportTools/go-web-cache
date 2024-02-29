# GO-Web-Cache

GO-Web-Cache is a simple caching server written in Go that acts as a reverse proxy to cache responses from a backend server. It includes features such as caching based on Cache-Control directives, health checks, metrics, and an admin API to manage cached items.

## Project Structure

The project is organized into the following packages:

- `pkg/logging`: Contains logging utilities for request and general logging.
- `pkg/cache`: Provides in-memory caching functionality and utilities.
- `pkg/metrics`: Handles Prometheus metrics collection and server.
- `pkg/config`: Loads and manages configuration settings.
- `pkg/proxy`: Implements the reverse proxy functionality with caching logic.
- `pkg/admin`: Includes handlers for the admin API endpoints.
- `pkg/health`: Implements health check and version information handlers.
- `pkg/security`: Contains security checks, such as WordPress login cookie detection.
- `main.go`: Main entry point for starting the caching server.

## Usage

To start the GO-Web-Cache server, run the following command:

```bash
go run main.go
```

This will start the caching server based on the configurations provided in `config.json` or through environment variables.

## Installation

### Kubernetes

```bash
helm repo supporttools https://charts.support.tools
helm repo update
helm install go-web-cache supporttools/go-web-cache
```

### Local Development

#### Prerequisites

- Go 1.16 or later

```bash
go get github.com/SupportTools/go-web-cache
cd $GOPATH/src/github.com/SupportTools/go-web-cache
export BACKEND_SERVER="http://example.com"
export PORT=8080
export DEBUG=true
go run main.go
```

### Docker

To run Go-Web-Cache in a Docker container, use the following command:

```bash
docker build -t supporttools/go-web-cache .
docker run -p 8080:8080 -e BACKEND_SERVER="http://example.com" -e PORT=8080 -e DEBUG=true supporttools/go-web-cache
```

## Configuration

The server can be configured using a `config.json` file or environment variables. The following settings are available:

## Enhancements and Customization

- Extend functionalities for different caching strategies or storage backends.
- Implement additional logging or security checks as needed.
- Customize cache eviction policies or caching rules based on specific requirements.

Feel free to explore and enhance the project to suit your caching needs and server requirements.

## Contributing

We welcome contributions from the community. Please refer to the [CONTRIBUTING.md](CONTRIBUTING.md) file for more details.

## License

This project is licensed under the Apache License - see the [LICENSE](LICENSE) file for details.

## Support

If you have any questions or encounter issues, please file an issue on the GitHub repository or contact us at [mmattox@support.tools].
