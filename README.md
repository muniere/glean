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

# Register new scraping task
$ glean scrape "https://example.com/"

# Cancel scraping task
$ glean cancel 1234
```

## License

This software is licensed under MIT.
