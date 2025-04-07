# ECG Monitoring Tool

A WebSocket-based real-time ECG data processing tool that monitors heart rate and detects irregularities.

## Features

- Real-time ECG data transmission via WebSockets
- Detection of heart irregularities:
  - Tachycardia (high heart rate)
  - Bradycardia (low heart rate)
  - Arrhythmia (irregular heartbeat patterns)
- Notifications for detected irregularities:
  - Logging to file
  - Console beep alerts
  - Visual indicators in terminal
- Configurable simulation parameters

## Project Structure

```
├── cmd/
│   ├── server/        # WebSocket server (ECG data sender)
│   ├── client/        # WebSocket client (ECG data receiver)
│   └── sound_test/    # Utility to test sound functionality
├── pkg/
│   ├── model/         # Data structures
│   ├── detection/     # Irregularity detection
│   └── notification/  # Logging and beeping
```

## Requirements

- Go 1.16 or higher
- Gorilla WebSocket package

## Installation

1. Clone the repository
2. Install dependencies:

```bash
go get github.com/gorilla/websocket
```

3. Build the project (optional):

```bash
make build
```

## Usage

### Using the Makefile

The simplest way to use the application is with the provided Makefile:

```bash
# Start the server
make run-server

# In another terminal, start the client
make run-client

# Test sound functionality
make test-sound

# Run all tests
make test
```

### Running Manually

#### Running the Server (ECG Data Sender)

```bash
go run cmd/server/main.go [options]
```

Options:
- `-addr`: Server address and port (default: ":8080")
- `-interval`: Interval between sending data (default: "1s", accepts formats like "500ms")
- `-irregular`: Simulate irregular heartbeats (default: true)

#### Running the Client (ECG Data Receiver)

```bash
go run cmd/client/main.go [options]
```

Options:
- `-addr`: Server address and port to connect to (default: "localhost:8080")
- `-logdir`: Directory to store log files (default: "logs")
- `-beep`: Enable beep sounds for alerts (default: true)
- `-lognormal`: Log normal readings in addition to irregularities (default: false)
- `-debug`: Enable debug mode for extra logging (default: false)

#### Sound Test Utility

A special utility is provided to test the sound functionality:

```bash
go run cmd/sound_test/main.go [options]
```

Options:
- `-mode`: Sound test mode: "tachycardia", "bradycardia", "arrhythmia", "system", or "all" (default: "all")
- `-severity`: Severity level: "mild", "moderate", or "severe" (default: "severe")

## Example Workflow

1. Start the server:
   ```bash
   make run-server
   ```

2. In another terminal, start the client:
   ```bash
   make run-client
   ```

3. The client will connect to the server and start receiving ECG data.
4. When irregularities are detected, the client will:
   - Display a "!" character in the terminal
   - Log the condition to the console and log file
   - Generate a beep alert if enabled

5. Check the log files in the specified log directory for detailed information.


## Running Tests

```bash
make test
```

## Configuration

You can adjust the thresholds for detecting irregularities in `pkg/model/ecg.go`. 
