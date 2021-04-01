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

`go test ./compare/ -long`

To run all benchmarks (including long benchmarks but not tests):

`go test ./compare/ -run ! -bench . -long`

### Benchmarks

Some benchmarks require a `large` directory in the `test` directory (`./test/large/`) that contains many large files. The `large` directory should contain a nested hierarchy of directories also containing large files. Any file with a file size greater than 1 MB is sufficient. The `large` directory should exceed 1 GB to ensure reliable results. The `large` directory is intentionally ignored by `git`.
