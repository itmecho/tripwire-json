# Tripwire JSON

Tripwire JSON is a simple tool for parsing a tripwire report and outputting it as JSON.

## Usage

To use tripwire json, you can either pass it a (non-encrypted) tripwire report file or pipe directly into it. There is a flag available (`-pretty`) to pretty print the output using 2 spaces

### From a file
```
sudo tripwire --check > /tmp/report
tripwire-json -file /tmp/report
```

### Direct via a pipe
```
sudo tripwire --check | tripwire-json
```

## Development

### Building the binary
Use the included Makefile to build the binary. This will set the version and git commit hash at build time.

```
make
```