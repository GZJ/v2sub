## v2sub

v2sub is a program that generates v2ray configuration files from subscription URL.

## Quick Start
```
go install github.com/GZJ/v2sub@latest

v2sub -name="xxx" \
      -sub_url="https://***" \
      -proxy_path="http://[proxy address]:[proxy port]" \
      -tmpl_file="*.tmpl" \
      -config_output_path="xxx"
```
