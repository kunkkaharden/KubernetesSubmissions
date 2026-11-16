{{/*
Expand the name of the chart.
*/}}
{{- define "docker-registry-ui.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "docker-registry-ui.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "docker-registry-ui.labels" -}}
app.kubernetes.io/name: {{ include "docker-registry-ui.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Labels por componente (ui | registry)
*/}}
{{- define "docker-registry-ui.componentLabels" -}}
{{ include "docker-registry-ui.labels" (index . 0) }}
app.kubernetes.io/component: {{ index . 1 }}
{{- end }}
