# GoTrics Node

The **GoTrics Node** is a lightweight agent that runs on multiple machines to collect system metrics (CPU, memory) and sends the data to the **GoTrics Server**. This component helps gather real-time metrics to be displayed on the frontend dashboard.

---

### Features

- **Metrics Collection**: Monitors system performance, such as CPU usage and memory consumption.
- **Periodic Data Sending**: Sends metrics to the GoTrics server every 5 seconds.
- **Written in Go**: Lightweight agent built using Go.

### Installation

1. Clone the repository:

```bash
git clone https://github.com/MathDesigns/gotrics-node.git
cd gotrics-node
```

2. Configure the node by editing the `config.yaml` to set the correct server URL and node ID.

3. Build and run the node:
```bash
go build -o gotrics-node ./gotrics-node
```

The node will start sending metrics to the configured server.

---

### License

This project is licensed under the MIT License.