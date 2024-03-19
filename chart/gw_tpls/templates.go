package gw_tpls

const HelpersTpl = `{{/*
Expand the name of the chart.
*/}}
{{- define "sng.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "sng.fullname" -}}
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
Create chart name and version as used by the chart label.
*/}}
{{- define "sng.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "sng.labels" -}}
helm.sh/chart: {{ include "sng.chart" . }}
{{ include "sng.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "sng.selectorLabels" -}}
app.kubernetes.io/name: {{ include "sng.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "sng.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "sng.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}
`

const DeploymentTpl = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "sng.fullname" . }}
  labels:
    {{- include "sng.labels" . | nindent 4 }}
spec:
  {{- if not .Values.autoscaling.enabled }}
  replicas: {{ .Values.replicaCount }}
  {{- end }}
  selector:
    matchLabels:
      {{- include "sng.selectorLabels" . | nindent 6 }}
  template:
    metadata:
      {{- with .Values.podAnnotations }}
      annotations:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      labels:
        {{- include "sng.selectorLabels" . | nindent 8 }}
    spec:
      {{- with .Values.imagePullSecrets }}
      imagePullSecrets:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      serviceAccountName: {{ include "sng.serviceAccountName" . }}
      securityContext:
        {{- toYaml .Values.podSecurityContext | nindent 8 }}
      volumes:
        []
      containers:
        - name: {{ .Chart.Name }}
          securityContext:
            {{- toYaml .Values.securityContext | nindent 12 }}
          image: "{{ .Values.image.repository }}:{{ .Values.image.tag | default .Chart.AppVersion }}"
          imagePullPolicy: {{ .Values.image.pullPolicy }}
          args:
            - {{ .Values.serviceName }}
            - --nacos-address=$(NACOS_ADDRESS)
            - --nacos-namespace=$(NACOS_NAMESPACE)
            - --password=$(SNG_PASSWORD)
            - --stage=$(SNG_STAGE)
            - --log-level=$(SNG_LOG_LEVEL)
            - --port=:$(SNG_HTTP_PORT)
            - --mgmt-port=:$(SNG_MGMT_PORT)
            - --pprof-port=:$(SNG_PPROF_PORT)
            - --trace=$(SNG_TRACE)
            - --lang=$(SNG_LANG)
          {{- if .Values.env }}
          env:
            {{- toYaml .Values.env | nindent 12 }}
          {{- end }}
          {{- if .Values.envFrom }}
          envFrom:
            {{- toYaml .Values.envFrom | nindent 12 }}
          {{- end }}
          ports:
            - name: http
              containerPort: {{ .Values.service.port | default 8000 }}
              protocol: TCP
            - name: mgmt-port
              containerPort: {{ .Values.service.mgmtport | default 8080 }}
              protocol: TCP              
          livenessProbe:
            httpGet:
              path: /health
              port: mgmt-port
            failureThreshold: 10
            periodSeconds: 10
            successThreshold: 1
          startupProbe:  # 启动探针，启动后禁用其他探针
            httpGet:
              path: /health
              port: mgmt-port
            initialDelaySeconds: 10   # 延迟探测时间(秒)【 在k8s第一次探测前等待秒 】
            periodSeconds: 10         # 执行探测频率(秒) 【 每隔秒执行一次 】
            timeoutSeconds: 1         # 超时时间
            successThreshold: 1       # 健康阀值 
            failureThreshold: 10      # 不健康阀值 
          lifecycle:
            preStop:
              exec:
                command: ["/bin/sh","-c","kill -15 1;while kill -0 1; do sleep 5; done; sleep 5"]
          resources:
            {{- toYaml .Values.resources | nindent 12 }}
      terminationGracePeriodSeconds: 60
      {{- with .Values.nodeSelector }}
      nodeSelector:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      {{- with .Values.affinity }}
      affinity:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      {{- with .Values.tolerations }}
      tolerations:
        {{- toYaml . | nindent 8 }}
      {{- end }}
`

const HpaYamlTpl = `{{- if .Values.autoscaling.enabled }}
apiVersion: autoscaling/v2beta1
kind: HorizontalPodAutoscaler
metadata:
  name: {{ include "sng.fullname" . }}
  labels:
    {{- include "sng.labels" . | nindent 4 }}
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: {{ include "sng.fullname" . }}
  minReplicas: {{ .Values.autoscaling.minReplicas }}
  maxReplicas: {{ .Values.autoscaling.maxReplicas }}
  metrics:
    {{- if .Values.autoscaling.targetCPUUtilizationPercentage }}
    - type: Resource
      resource:
        name: cpu
        targetAverageUtilization: {{ .Values.autoscaling.targetCPUUtilizationPercentage }}
    {{- end }}
    {{- if .Values.autoscaling.targetMemoryUtilizationPercentage }}
    - type: Resource
      resource:
        name: memory
        targetAverageUtilization: {{ .Values.autoscaling.targetMemoryUtilizationPercentage }}
    {{- end }}
{{- end }}
`

const IngressTpl = `{{- if .Values.ingress.enabled -}}
{{- $fullName := include "sng.fullname" . -}}
{{- $svcPort := .Values.service.port -}}
{{- if semverCompare ">=1.14-0" .Capabilities.KubeVersion.GitVersion -}}
apiVersion: networking.k8s.io/v1beta1
{{- else -}}
apiVersion: extensions/v1beta1
{{- end }}
kind: Ingress
metadata:
  name: {{ $fullName }}
  labels:
    {{- include "sng.labels" . | nindent 4 }}
  {{- with .Values.ingress.annotations }}
  annotations:
    {{- toYaml . | nindent 4 }}
  {{- end }}
spec:
  {{- if .Values.ingress.tls }}
  tls:
    {{- range .Values.ingress.tls }}
    - hosts:
        {{- range .hosts }}
        - {{ . | quote }}
        {{- end }}
      secretName: {{ .secretName }}
    {{- end }}
  {{- end }}
  rules:
    {{- range .Values.ingress.hosts }}
    - host: {{ .host | quote }}
      http:
        paths:
          {{- range .paths }}
          - path: {{ .path }}
            backend:
              serviceName: {{ $fullName }}
              servicePort: {{ $svcPort }}
          {{- end }}
    {{- end }}
  {{- end }}
`

const IngeressHaproxyTpl = `{{- if .Values.ingress_haproxy.enabled -}}
{{- $fullName := include "sng.fullname" . -}}
{{- $svcPort := .Values.service.port -}}
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: {{ $fullName }}-haproxy
  labels:
    {{- include "sng.labels" . | nindent 4 }}
  {{- with .Values.ingress_haproxy.annotations }}
  annotations:
    {{- toYaml . | nindent 4 }}
  {{- end }}
spec:
  {{- if .Values.ingress_haproxy.tls }}
  tls:
    {{- range .Values.ingress_haproxy.tls }}
    - hosts:
        {{- range .hosts }}
        - {{ . | quote }}
        {{- end }}
      secretName: {{ .secretName }}
    {{- end }}
  {{- end }}
  rules:
    {{- range .Values.ingress_haproxy.hosts }}
    - host: {{ .host | quote }}
      http:
        paths:
          {{- range .paths }}
          - path: {{ .path }}
            pathType: {{ .pathType}}
            backend:
              service:
                name: {{ $fullName }}
                port: 
                  number: {{ $svcPort }}
          {{- end }}
    {{- end }}
{{- end }}
`

const NotesTpl = `1. Get the application URL by running these commands:
{{- if .Values.ingress.enabled }}
{{- range $host := .Values.ingress.hosts }}
  {{- range .paths }}
  http{{ if $.Values.ingress.tls }}s{{ end }}://{{ $host.host }}{{ .path }}
  {{- end }}
{{- end }}
{{- else if contains "NodePort" .Values.service.type }}
  export NODE_PORT=$(kubectl get --namespace {{ .Release.Namespace }} -o jsonpath="{.spec.ports[0].nodePort}" services {{ include "sng.fullname" . }})
  export NODE_IP=$(kubectl get nodes --namespace {{ .Release.Namespace }} -o jsonpath="{.items[0].status.addresses[0].address}")
  echo http://$NODE_IP:$NODE_PORT
{{- else if contains "LoadBalancer" .Values.service.type }}
     NOTE: It may take a few minutes for the LoadBalancer IP to be available.
           You can watch the status of by running 'kubectl get --namespace {{ .Release.Namespace }} svc -w {{ include "sng.fullname" . }}'
  export SERVICE_IP=$(kubectl get svc --namespace {{ .Release.Namespace }} {{ include "sng.fullname" . }} --template "{{"{{ range (index .status.loadBalancer.ingress 0) }}{{.}}{{ end }}"}}")
  echo http://$SERVICE_IP:{{ .Values.service.port }}
{{- else if contains "ClusterIP" .Values.service.type }}
  export POD_NAME=$(kubectl get pods --namespace {{ .Release.Namespace }} -l "app.kubernetes.io/name={{ include "sng.name" . }},app.kubernetes.io/instance={{ .Release.Name }}" -o jsonpath="{.items[0].metadata.name}")
  export CONTAINER_PORT=$(kubectl get pod --namespace {{ .Release.Namespace }} $POD_NAME -o jsonpath="{.spec.containers[0].ports[0].containerPort}")
  echo "Visit http://127.0.0.1:8080 to use your application"
  kubectl --namespace {{ .Release.Namespace }} port-forward $POD_NAME 8080:$CONTAINER_PORT
{{- end }}
`

const ServiceTpl = `{{- if .Values.service.create -}}
apiVersion: v1
kind: Service
metadata:
  name: {{ include "sng.fullname" . }}
  labels:
    app: {{ include "sng.fullname" . }}
    {{- include "sng.labels" . | nindent 4 }}
spec:
  type: {{ .Values.service.type }}
  {{- with .Values.service.headless }}
  clusterIP: None
  {{- end }}
  ports:
    - port: {{ .Values.service.port }}
      targetPort: http
      protocol: TCP
      name: http
    - port: {{ .Values.service.mgmtport }}
      targetPort: mgmt-port
      protocol: TCP
      name: mgmt-port
  selector:
    {{- include "sng.selectorLabels" . | nindent 4 }}
{{- end  }}
`
const ServiceaccountTpl = `{{- if .Values.serviceAccount.create -}}
apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ include "sng.serviceAccountName" . }}
  labels:
    {{- include "sng.labels" . | nindent 4 }}
  {{- with .Values.serviceAccount.annotations }}
  annotations:
    {{- toYaml . | nindent 4 }}
  {{- end }}
{{- end }}
`
const ServiceMonitorTpl = `{{ if .Values.serviceMonitor.enabled }}
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: {{ include "sng.fullname" . }}
  #namespace: {{ .Release.Namespace | quote }}
  labels:
    release: kube-prometheus
    app: {{ include "sng.fullname" . }}
spec:
  endpoints:
    - port: mgmt-port
      path: /metrics/prometheus
  namespaceSelector:
    matchNames: 
      - {{ .Release.Namespace | quote }}
  selector:
    matchLabels:
      app: {{ include "sng.fullname" . }}
{{- end }}
`
