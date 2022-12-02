{{- define "type_members_html" -}}
{{- $type := . -}}
{{- "<table><tbody>" -}}
{{- range $type.Members -}}
<tr><td>`{{ .Name  }}` _{{ markdownRenderType .Type }}_</td><td>{{- template "type_member" . -}}</td></tr>
{{- end -}}
</tbody></table>

{{- end -}}
