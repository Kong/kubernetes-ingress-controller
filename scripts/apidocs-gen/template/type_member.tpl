{{- define "type_member" -}}
{{- $field := . -}}
{{- $isNotBasic := not $field.Type.IsBasic -}}
{{- $isNotImported := not $field.Type.Imported -}}
{{- $isNotJSON := not (eq $field.Type.Name "JSON") -}}

{{- if eq $field.Name "metadata" -}}
Refer to Kubernetes API documentation for fields of `metadata`.
{{- else if and $isNotBasic $isNotImported $isNotJSON -}}
{{ $field.Doc }}<br/>{{ template "type_members_html" $field.Type }}
{{- else -}}
{{ $field.Doc }}
{{- end -}}
{{- end -}}
