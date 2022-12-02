{{- define "gvDetails" -}}
{{- $gv := . -}}

## {{ $gv.GroupVersionString }}

{{ $gv.Doc }}

{{- if $gv.Kinds  }}
{{- range $gv.SortedKinds }}
- {{ $gv.TypeForKind . | markdownRenderTypeLink }}
{{- end }}
{{ end }}

{{- /* Display exported Kinds first */ -}}
{{- range $gv.SortedKinds -}}
{{- $typ := $gv.TypeForKind . }}
{{ template "type" $typ }}
{{ end -}}

{{- /* Display Types that are not exported Kinds */ -}}
{{- range $typ := $gv.SortedTypes -}}
{{- $isKind := false -}}
{{- range $kind := $gv.SortedKinds -}}
{{- if eq $typ.Name $kind -}}
{{- $isKind = true -}}
{{- end -}}
{{- end -}}
{{- if not $isKind }}
{{ template "type" $typ }}
{{ end -}}
{{- end }}

{{- end -}}
