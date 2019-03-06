# go-http-monitor

> Network monitoring is the use of a system that constantly monitors a computer network for slow or failing components and that notifies the network administrator (via email, SMS or other alarms) in case of outages or other trouble. -- [Wikipedia](https://en.wikipedia.org/wiki/Network_monitoring)

This program monitors an actively updated `access.log` file, and notifies users when thresholds are either exceeded or recover. The thresholds are configurable via command-line flags.

This project was designed as part of a interview take-home code test, and so the program itself generates the `access.log` and populates it with data in order to simulate real traffic patterns (that logic is, as you can imagine, very basic).

## Flags

```
-e int
      monitoring evaluation period in minutes (shorthand) (default 2)
-evaluation int
      monitoring evaluation period in minutes (default 2)
-help
      show available command flags
-l string
      location of access.log file to monitor (shorthand) (default "./access.log")
-location string
      location of access.log file to monitor (default "./access.log")
-s int
      statistic output interval in seconds (shorthand) (default 10)
-stats int
      statistic output interval in seconds (default 10)
-t int
      average alarm threshold time period (shorthand) (default 10)
-threshold int
      average alarm threshold time period (default 10)
-u string
      unit of time of the alarm threshold (shorthand) (default "second")
-unit string
      unit of time of the alarm threshold (default "second")
```

## Access Log Format

We're monitoring access logs based upon the W3C httpd [common log file format](https://www.w3.org/Daemon/User/Config/Logging.html#common-logfile-format).

- remotehost: Remote hostname (or IP number if DNS hostname is not available, or if DNSLookup is Off.
- rfc931: The remote logname of the user.
- authuser: The username as which the user has authenticated himself.
- [date]: Date and time of the request.
- "request": The request line exactly as it came from the client.
- status: The HTTP status code returned to the client.
- bytes: The content-length of the document transferred.

### Example Access Logs

```
127.0.0.1 - james [09/May/2018:16:00:39 +0000] "GET /report HTTP/1.0" 200 123
```

Which (as described by the common log file format) breaks down to:

- `127.0.0.1`: remote hostname
- `-`: remote logname
- `james`: username
- `[09/May/2018:16:00:39 +0000]`: request date/time
- `"GET /report HTTP/1.0"`: http request
- `200`: http response status code
- `123`: content-length of response

## Requirements

The following requirements were what I was provided with:

- Consume an actively written-to w3c-formatted HTTP access log (https://www.w3.org/Daemon/User/Config/Logging.html). It should default to reading `/tmp/access.log` and be overrideable.
- Display stats every 10s about the traffic during those 10s: the sections of the web site with the most hits, as well as interesting summary statistics on the traffic as a whole. A section is defined as being what's before the second '/' in the resource section of the log line. For example, the section for "/pages/create" is "/pages".
- Make sure a user can keep the app running and monitor the log file continuously.
- Whenever total traffic for the past 2 minutes exceeds a certain number on average, add a message saying that “High traffic generated an alert - hits = {value}, triggered at {time}”. The default threshold should be 10 requests per second, and should be overridable.
- Whenever the total traffic drops again below that value on average for the past 2 minutes, add another message detailing when the alert recovered.
- Make sure all messages showing when alerting thresholds are crossed remain visible on the page for historical reasons.
- Write a test for the alerting logic.
- Explain how you'd improve on this application design.
- Have something else write to the log file.
- Dockerize if possible.

## Generated Access Log

Because I don't have an actively written to `access.log` file by default, I'm generating one as part of this program. If you happen to have one you should be able to provide it to the program using the `-location` flag (the default value is `"./access.log"` which is a file comitted as part of this repository and is _reset_ on each call to `make run` -- see the [Makefile](./Makefile) for more details).

## Running the program

I use a [Makefile](./Makefile) to generate shortened commands, to make running the program easier:

- `make run`
- `make build`
- `make test`
- `make clean`

If you want the program to generate access logs for you, then ensure you use the `-populate` flag, like so:

```bash
make run flags="-populate"
```

## Statistical Output

I don't do anything fancy with the log data other than pretty print specific users and their requests over the past N seconds, as defined by the `-stats` flag (see the [flags](#flags) section above for more details).

```
--------------------------------------
Stats for last 10 seconds of requests:

{
  "Bob": {
    "/bar/f": 1,
    "/bar/n": 1,
    "/bar/p": 2,
    "/baz/u": 1,
    "/baz/v": 1,
    "/foo/g": 1,
    "/foo/n": 3,
    "/foo/s": 1,
    "/qiz/m": 3,
    "/qiz/o": 1,
    "/qiz/w": 1,
    "/qux/c": 2,
    "/qux/d": 2,
    "/qux/t": 1
  },
  "Jane": {
    "/bar/f": 1,
    "/bar/n": 1,
    "/bar/p": 2,
    "/baz/u": 1,
    "/baz/v": 1,
    "/foo/g": 1,
    "/foo/n": 3,
    "/foo/s": 1,
    "/qiz/m": 3,
    "/qiz/o": 1,
    "/qiz/w": 1,
    "/qux/c": 2,
    "/qux/d": 2,
    "/qux/t": 1
  },
  "Lisa": {
    "/bar/f": 1,
    "/bar/n": 1,
    "/bar/p": 2,
    "/baz/u": 1,
    "/baz/v": 1,
    "/foo/g": 1,
    "/foo/n": 3,
    "/foo/s": 1,
    "/qiz/m": 3,
    "/qiz/o": 1,
    "/qiz/w": 1,
    "/qux/c": 2,
    "/qux/d": 2,
    "/qux/t": 1
  },
  "Mark": {
    "/bar/m": 1,
    "/bar/s": 1,
    "/qiz/y": 1,
    "/qux/m": 1,
    "/qux/n": 1
  }
}
--------------------------------------
```

## TODO

- more idiomatic package names
- don't generate a `stats.Stat` struct until stats analysis is needed
- write test(s)
- dockerize
