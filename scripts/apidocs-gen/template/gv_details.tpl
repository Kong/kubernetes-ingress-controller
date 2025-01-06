{{- define "gvDetails" -}}
{{- $gv := . -}}

## {{ $gv.GroupVersionString }}

{{ $gv.Doc }}

{{- if $gv.Kinds  }}
{{- range $gv.SortedKinds }}
{{- $typ := $gv.TypeForKind . }}
{{- /* Display only KIC supported kinds */ -}}
{{- if index $typ.Markers "apireference:kic:include" }}
- {{ $gv.TypeForKind . | markdownRenderTypeLink }}
{{- end }}
{{- end }}
{{ end }}

{{- /* Display exported Kinds first */ -}}
{{- range $gv.SortedKinds -}}
{{- $typ := $gv.TypeForKind . }}
{{- $isKind := true -}}
{{- /* Display only KIC supported kinds */ -}}
{{- if index $typ.Markers "apireference:kic:include" -}}
{{ template "type" (dict "type" $typ "isKind" $isKind) }}
{{- end }}
{{ end -}}

### Types

In this section you will find types that the CRDs rely on.

{{- /* Display Types that are not exported Kinds */ -}}
{{- range $typ := $gv.SortedTypes -}}
{{- $isKind := false -}}
{{- range $kind := $gv.SortedKinds -}}
{{- if eq $typ.Name $kind -}}
{{- $isKind = true -}}
{{- end -}}
{{- end -}}
{{- /* Display only KIC supported types */ -}}
{{- if and (not $isKind) (index $typ.Markers "apireference:kic:include") }}
{{ template "type" (dict "type" $typ "isKind" $isKind) }}
{{ end -}}
{{- end }}

{{- end -}}
