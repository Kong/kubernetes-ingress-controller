{{- define "type_members_html" -}}
{{- $type := . -}}
<table><tbody>
{{- range $type.Members -}}
<tr><td>{{ .Name  }} ({{ markdownRenderType .Type }})</td><td>{{- template "type_member" . -}}</td></tr>
{{- end -}}
</tbody></table>

{{- end -}}
