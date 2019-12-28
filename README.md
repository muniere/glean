# Glean

Glean is a web crawler with Server-Client style.

## Get Started

```bash
# Clone
$ git clone git@github.com:muniere/glean.git
$ cd glean

# Install
$ make && make install

# Kick server
$ gleand
```

## Control via CLI

```bash
# Check current status
$ glean status

# Register new crawling
$ glean launch "https://example.com/"

# Unregister or abort crawling
$ glean cancel "https://example.com/"
```

## License

This software is licensed under MIT.
