{{- define "type" -}}
{{- $type := . -}}
{{- if markdownShouldRenderType $type -}}

### {{ $type.Name }}

{{ if $type.IsAlias }}_Underlying type:_ `{{ markdownRenderTypeLink $type.UnderlyingType  }}`{{ end }}

{{ $type.Doc }}

{{ if $type.Members -}}
<table>
<thead><tr>
    <td>Field</td><td>Description</td>
</tr></thead>
<tbody>
{{ if $type.GVK -}}
<tr>
    <td>`apiVersion` _string_</td>
    <td>`{{ $type.GVK.Group }}/{{ $type.GVK.Version }}`</td>
</tr>
<tr>
    <td>`kind` _string_</td>
    <td>`{{ $type.GVK.Kind }}`</td>
</tr>
{{ end -}}
{{ range $type.Members -}}
<tr>
    <td>`{{ .Name  }}` _{{ markdownRenderType .Type }}_</td>
    <td>{{ template "type_member" . }} </td>
</tr>
{{ end -}}

</tbody></table>
{{ end }}

{{ if $type.References -}}
_Appears in:_
{{- range $type.SortedReferences }}
- {{ markdownRenderTypeLink . }}
{{- end }}
{{- end }}

{{- end -}}
{{- end -}}
