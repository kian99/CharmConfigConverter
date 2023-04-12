## Description

This project can be used to convert a charm's config.yaml file to a vars.tf suitable for use with a Terraform spec.

Build the project with
```
go build .
```

Run the project with
```
./config_converter <path-to-config.yaml> <path-to-output>
```