## v2sub

v2sub is a program that generates v2ray configuration files from subscription URL.

v2meu is an add-on tool for selecting and starting v2ray.

## Quick Start
```
v2sub -name="xxx" \
    -sub-url="https://***" \
    -proxy-path="http://[proxy address]:[proxy port]" \
    -tmpl-file="*.tmpl" \
    -config-output-path="xxx"

v2menu -dir="[v2ray config path]" \
    -stop-cmd="[kill command]" \
    -stop-args="[process name]" \
    -run-cmd="[v2ray path]" \
    -run-args="--config=%s"
```
