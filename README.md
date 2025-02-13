# Resource Simulator

This repository simulates CPU and memory usage.

## Build

```bash
go build -o resource-simulator cmd/server/main.go
```

## Run

```bash
./resource-simulator
```

## API Specification

| Endpoint | Parameters                                                                                             | Response                                                                                                                   | Description                                |
| -------- | ------------------------------------------------------------------------------------------------------ | -------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------ |
| /ping    | None                                                                                                   | JSON:<br> - `result`(bool)                                                                                                 | Health check endpoint.                     |
| /init    | - `cpu` (int): CPU load <br>- `mem` (int): Memory load in MB<br>- `duration`(int): Simulation duration | JSON:<br> - `cpuLoad(int)`: CPU load <br> - `memLoadMB(int)`: Memory load in MB<br> - `duration(int)`: Simulation duration | Simulates CPU, Memory load for a duration. |
| /abort   | None                                                                                                   | JSON:<br> - `result`(bool)                                                                                                 | Aborts the simulation.                     |
