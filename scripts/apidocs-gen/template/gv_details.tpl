{{- define "gvDetails" -}}
{{- $gv := . -}}

## {{ $gv.GroupVersionString }}

{{ $gv.Doc }}

{{- if $gv.Kinds  }}
{{- range $gv.SortedKinds }}
- {{ $gv.TypeForKind . | markdownRenderTypeLink }}
{{- end }}
{{ end }}

{{ range $gv.SortedTypes }}
{{ template "type" . }}
{{ end }}

{{- end -}}
