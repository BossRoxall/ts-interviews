{{/*
Expand the name of the chart.
*/}}
{{- define "rate-limiter.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name. If the release name contains the
chart name it is used as the full name.
*/}}
{{- define "rate-limiter.fullname" -}}
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

{{- define "rate-limiter.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "rate-limiter.labels" -}}
helm.sh/chart: {{ include "rate-limiter.chart" . }}
{{ include "rate-limiter.selectorLabels" . }}
app.kubernetes.io/version: {{ .Values.image.tag | default .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Stable labels only: the Service and Deployment selector must not change between
releases.
*/}}
{{- define "rate-limiter.selectorLabels" -}}
app.kubernetes.io/name: {{ include "rate-limiter.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "rate-limiter.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "rate-limiter.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}
