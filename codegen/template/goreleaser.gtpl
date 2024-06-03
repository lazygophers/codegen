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
            - -X github.com/lazygophers/utils/app.Name={{ .ProjectName }}
            - -X github.com/lazygophers/utils/app.Commit={{ .FullCommit }}
            - -X github.com/lazygophers/utils/app.ShortCommit={{ .ShortCommit }}
            - -X github.com/lazygophers/utils/app.Branch={{ .Branch }}
            - -X github.com/lazygophers/utils/app.Version={{ .Version }}
            - -X github.com/lazygophers/utils/app.Tag={{ .Tag }}
            - -X github.com/lazygophers/utils/app.BuildDate={{ .Date }}
            - -X github.com/lazygophers/utils/app.GoVersion={{ .Env.GOVERSION }}
            - -X github.com/lazygophers/utils/app.GoOS={{ .Os }}
            - -X github.com/lazygophers/utils/app.Goarch={{ .Arch }}
            - -X github.com/lazygophers/utils/app.Goarm={{ .Arm }}
            - -X github.com/lazygophers/utils/app.Goamd64={{ .Amd64 }}
            - -X github.com/lazygophers/utils/app.Gomips={{ .Mips }}
        gcflags:
            - -N -l
        targets:
            - linux_amd64_v3

dockers:
	-   id: {{ .PB.GoPackageName }}
		image_templates:
			- "{{ .ProjectName }}:{{ .Tag }}"
			- "{{ .ProjectName }}:latest"
		skip_build: true
		skip_push: true
		push_flags:
			- --tls-verify=false
		extra_files:
			- config.yaml

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
