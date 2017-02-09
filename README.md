# shadow

**shadow** is a simple and fast http proxy written by Go.

## Usage

Start a transparent http proxy

```
./shadow -addr :7778
```

With parent proxy

```
./shadow -addr :7777 -proxy :7778
```

Using configuration file instead of command line args

```
./shadow -c shadow.json
```

## Args

* config, c: configuration file, all other args arg ignored if config is specified.
* addr, a: Listen address (default ":1080").
* proxy, p: Parent proxy address.
* timeout, t: Timeout for waiting response headers, by milliseconds.
* verbose, v: Verbose output.

## Configuration

Here is a config sample

```json
{
  "addr": ":7778",
  "timeout": 30000,
  "proxy": "",
  "verbose": true
}
```