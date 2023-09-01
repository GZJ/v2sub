## v2sub

v2sub is a program that generates v2ray configuration files from subscription URL.

## Quick Start
```
go run v2sub.go -name="xxx" \
                -sub_url="https://***" \
                -proxy_path="http://[proxy address]:[proxy port]" \
                -tmpl_file="*.tmpl" \
                -config_output_path="xxx"
```
