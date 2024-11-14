{{- define "type" -}}
{{- $type := $.type -}}
{{- $isKind := $.isKind -}}
{{- if markdownShouldRenderType $type -}}
{{- if not (index $type.Markers "apireference:kic:exclude") -}}

{{- if $isKind -}}
### {{ $type.Name }}
{{ else -}}
#### {{ $type.Name }}
{{ end -}}

{{ if $type.IsAlias }}_Underlying type:_ `{{ markdownRenderTypeLink $type.UnderlyingType  }}`{{ end }}

{{ $type.Doc | replace "\n\n" "<br /><br />" }}

{{ if $type.GVK -}}
<!-- {{ snakecase $type.Name }} description placeholder -->
{{- end }}

{{ if $type.Members -}}
| Field | Description |
| --- | --- |
{{ if $type.GVK -}}
| `apiVersion` _string_ | `{{ $type.GVK.Group }}/{{ $type.GVK.Version }}`
| `kind` _string_ | `{{ $type.GVK.Kind }}`
{{ end -}}

{{ range $type.Members -}}
| `{{ .Name  }}` _{{ markdownRenderType .Type }}_ | {{ template "type_members" . }} |
{{ end -}}

{{ end }}

{{ if $type.References -}}
_Appears in:_
{{- range $type.SortedReferences }}
- {{ markdownRenderTypeLink . }}
{{- end }}
{{- end }}

{{- end -}}
{{- end -}}
{{- end -}}
