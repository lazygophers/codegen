message {{ .RequestType }} {
	enum ListOption {
		ListOptionNil = {{ IncrKey $.RequestType 0 }};{{ range $key, $value := .ListOptions }}{{ with $value }}
        // @type: {{ $value }}{{ end }}
		ListOption{{ ToCamel $key }} = {{ IncrKey $.RequestType 0 }};{{ end }}
	}

	// @validate: required
	core.ListOption list_option = 1;
}
