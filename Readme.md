# Directory Compare

Command line utility to compare one directory to another.

## Testing

To run tests:

`go test ./compare/`

To run tests and benchmarks:

`go test ./compare/ -bench .`

To only run benchmarks (without tests):

`go test ./compare/ -run ! -bench .`

To run all tests (including long tests):

`go test ./compare -long`

To run all benchmarks (including long benchmarks but not tests):

`go test ./compare/ -run ! -bench . -long`
