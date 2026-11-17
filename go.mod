module github.com/nattokin/go-backlog

go 1.26

require github.com/stretchr/testify v1.12.1

require go.yaml.in/yaml/v3 v3.0.5 // indirect

retract v0.25.6 // Contains breaking changes; use v0.26.0 instead.
