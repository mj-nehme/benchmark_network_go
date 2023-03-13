# Go-Baseline

This simple code is just to measure the throughput between a client and a server using:

- gob package
- binary package
- net

To run the server: `go run *.go server`

To run the client: `go run *.go client`

My results on the localhost (without real network) were as follows:
- Throughput using `GOB` Package:           `54.38Gb/s`
- Throughput using `Binary` Package:        `53.89Gb/s`
- Throughput using `net` Package:           `117.84Gb/s`