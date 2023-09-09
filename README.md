# goanner
A Cool Port Scanner Written In Go

## Overview
A basic port scanner written entirely in Go it uses the `net` package for connecting and interacting with the target. The script utilizes goroutines and channels to optimize performance.

### Features:
- The ability to provide a port range (example: 100-8000)
- Validating inputs and error handling
- Verbose mode for more informative outputs (Coming Soon)

### Supported Scan Types:
- TCP Scan
- UDP Scan (Coming Soon)

## Usage
```
git clone https://github.com/mohabgabber/goanner.git && cd goanner
```
To run the script directly:

```
go run main.go -t [TARGET ADDRESS (DOMAIN/IP)] -p [PORT RANGE (100-2300)]
```
You can also compile the script:
```
go build main.go
```
