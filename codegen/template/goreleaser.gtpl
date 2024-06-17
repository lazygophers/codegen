version: 2
project_name: {{ .PB.GoPackageName }}
env:
    - GO111MODULE=on
    - GOSUMDB=off
    - CGO_ENABLED=0

gomod:
    proxy: true
    mod: mod

builds:
    -   id: {{ .PB.GoPackageName }}
        binary: {{ .PB.GoPackageName }}
        main: ./cmd/
        tags:
            - release
        ldflags:
            - -s -w
            - --extldflags "-static -fpic"
            - -X github.com/lazygophers/utils/app.Name={{ Self ".ProjectName" }}}}
            - -X github.com/lazygophers/utils/app.Commit={{ Self ".FullCommit" }}
            - -X github.com/lazygophers/utils/app.ShortCommit={{ Self ".ShortCommit" }}
            - -X github.com/lazygophers/utils/app.Branch={{ Self ".Branch" }}
            - -X github.com/lazygophers/utils/app.Version={{ Self ".Version" }}
            - -X github.com/lazygophers/utils/app.Tag={{ Self ".Tag" }}
            - -X github.com/lazygophers/utils/app.BuildDate={{ Self ".Date" }}
            - -X github.com/lazygophers/utils/app.GoVersion={{ .EnvSelf ".GOVERSION" }}
            - -X github.com/lazygophers/utils/app.GoOS={{ Self ".Os" }}
            - -X github.com/lazygophers/utils/app.Goarch={{ Self ".Arch" }}
            - -X github.com/lazygophers/utils/app.Goarm={{ Self ".Arm" }}
            - -X github.com/lazygophers/utils/app.Goamd64={{ Self ".Amd64" }}
            - -X github.com/lazygophers/utils/app.Gomips={{ Self ".Mips" }}
        gcflags:
            - -N -l
        targets:
            - linux_amd64_v3

dockers:
    -   id: {{ .ProjectName }}
        skip_build: true
        skip_push: true
        push_flags:
            - --tls-verify=false
        extra_files:
            - config.yaml
        image_templates:
            - "{{ .ProjectName }}:latest"

archives:
    -   format: tar.gz
        name_template: '{{ .ProjectName }}_{{ .Os }}_{{ .Arch }}{{ with .Arm }}v{{ . }}{{ end }}{{ with .Mips }}_{{ . }}{{ end }}{{ if not (eq .Amd64 "v1") }}{{ .Amd64 }}{{ end }}'
        format_overrides:
            -   goos: windows
                format: zip
        files:
            - config.example.yaml

checksum:
    algorithm: sha512
    name_template: '{{ .ProjectName }}_{{ .Version }}.sha512'

source:
    enabled: false
    files:
        - LICENSE
        - README.md
        - config.example.yaml
        - '*.go'

report_sizes: true
