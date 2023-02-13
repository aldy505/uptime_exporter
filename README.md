# Uptime Exporter

Like [UptimeRobot](https://uptimerobot.com/), [Uptime Kuma](https://github.com/louislam/uptime-kuma)
or any other status page monitoring tool, but for Prometheus. So you can display it on your own Grafana.

## Installation

Simple!
1. Download the binary for your OS from the Release page. If you don't see one for your OS/Arch, I'm afraid that you'll need to build it from source.
2. Create a configuration file, it can be JSON, TOML, or YAML, depending on what you like:
   ```json
   {
     "endpoints": [
       {
         "name": "GitHub",
         "address": "https://github.com/healthz"
       },
       {
         "name": "Reinaldy",
         "address": "https://code.reinaldyrafli.com",
         "timeout": 60
       }
     ]
   }
   ```
   ```yaml
   endpoints:
       - name: "GitHub"
         address: https://github.com/healthz
       - name: "Reinaldy"
         address: https://code.reinaldyrafli.com
         timeout: 60
   ```
   ```toml
   [[endpoints]]
   name = "Github"
   address = "https://github.com/healthz"
   [[endpoints]]
   name = "Reinaldy"
   address = "htps://code.reinaldyrafli.com"
   timeout = 60
   ```
3. Run the binary!
   ```bash
   uptime_exporter --config-file ./path/to/config.json --web.listen-address :9428
   ```
   
If you want it to listen to a specific interface, you can use flags like so:  `--web.listen-address 127.0.0.1:9428`

## Usage

If you're using the default settings, which listens to `0.0.0.0:9428`, and you open: http://your-ip:9428/metrics,
you will get:

```
# HELP uptime_is_up Is it up?
# TYPE uptime_is_up gauge
uptime_is_up{endpoint_address="https://code.reinaldyrafli.com",endpoint_name="Reinaldy"} 1
uptime_is_up{endpoint_address="https://github.com/healthz",endpoint_name="GitHub"} 1
# HELP uptime_latency_seconds Measured latency on last scrape
# TYPE uptime_latency_seconds gauge
uptime_latency_seconds{endpoint_address="https://code.reinaldyrafli.com",endpoint_name="Reinaldy"} 0.1414042
uptime_latency_seconds{endpoint_address="https://github.com/healthz",endpoint_name="GitHub"} 0.3379525
```

`uptime_is_up` indicates whether the website (as stated by the `endpoint_address` and `endpoint_name` fields)
is reachable or not by the HTTP request the server just made when Prometheus scrape the page. If the value is 1,
then it's up and healthy. If the value is 0, then it's down.

`uptime_latency_seconds` indicates the time it took to do a HTTP request to that specified `endpoint_address`
field in second.

## Configuration file schema

If you want to do config.json.

```json
{
    "endpoints": [
        {
            "name": "string",
            "address": "https://url",
            "method": "GET",
            "timeout": 123,
            "successful_status_code": "2xx",
            "inverse_status": false,
            "tls_configuration": {
                "certificate_authority_path": "/path/to/certificate_authority.pem",
                "client_certificate_path": "/path/to/client_certificate.pem",
                "client_key_path": "/path/to/client_key.pem",
                "insecure_skip_verify": true
            }
        }
    ]
}
```

If you want to do config.yaml (or .yml if you prefer).

```yaml
endpoints:
    - name: "string"
      address: "https://url"
      method: "GET"
      timeout: 123
      successful_status_code: 2xx
      inverse_status: false
      tls_configuration:
        certificate_authority_path: "/path/to/certificate_authority.pem",
        client_certificate_path: "/path/to/client_certificate.pem",
        client_key_path: "/path/to/client_key.pem",
        insecure_skip_verify: true
```

## Build from source

Install [Go](https://go.dev/dl).

```sh
go build -o uptime_exporter .

# For windows
go build -o uptime_exporter.exe .
```

## License

```
MIT License

Copyright (c) 2023 Reinaldy Rafli <aldy505@proton.me>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

See [LICENSE](./LICENSE)