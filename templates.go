package main

const variable_template = `
{{range $key, $value := .}}
variable "{{ $key }}" {
	type = {{ $value.Type}}
	{{ with $value.Default -}}
		{{ if eq $value.Type "string" -}}
	default={{printf "%q" $value.Default }} 
		{{ else }}
	default={{ $value.Default }} 
		{{- end }}
	{{- end }}
}
{{ end }}
`
