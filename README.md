# profilegen

A simple Go CLI tool to generate mock profiling data and send it to a Pyroscope server over OTLP gRPC.

## Usage

Build the tool:

```
go build -o profilegen
```

Run the tool with required arguments:

```
./profilegen --service_name=my-mock-service --pyroscope_url=http://localhost:4040
```

- `--service_name`: The service name tag to attach to the profile (for filtering in Pyroscope UI)
- `--pyroscope_url`: The Pyroscope server API endpoint URL

The tool will generate CPU and memory load for 30 seconds and send the profile data to the specified Pyroscope server.

## Requirements
- Go 1.18+
- A running Pyroscope server (OSS or Grafana Cloud)

## Notes
- The tool uses [pyroscope-go](https://github.com/grafana/pyroscope-go) for profiling.
- The service name is attached as a tag for easy filtering in the Pyroscope UI. 