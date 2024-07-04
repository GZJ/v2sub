## v2sub

v2sub is a program that generates v2ray configuration files from subscription URL.

## Quick Start
```
go install github.com/gzj/v2sub@latest

v2sub -name="xxx" \
      -sub-url="https://***" \
      -proxy-path="http://[proxy address]:[proxy port]" \
      -tmpl-file="*.tmpl" \
      -config-output-path="xxx"
```
