# go-http-monitor

> Network monitoring is the use of a system that constantly monitors a computer network for slow or failing components and that notifies the network administrator (via email, SMS or other alarms) in case of outages or other trouble. -- [Wikipedia](https://en.wikipedia.org/wiki/Network_monitoring)

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

The following requirements were provided by the business:

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

## TODO

- figure out how best to access.log and track changes over time period
- figure out how to populate access.log at runtime (simulating real-time requests)
- figure out how to calculate when requests per section go over/under a threshold
- figure out what 'interesting' stats to show every 10s
- write test(s)
- dockerize

- name some packages more appropriately (e.g. verbs not nouns)
