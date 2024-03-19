package gw_tpls

const TestConnectionTpl = `apiVersion: v1
kind: Pod
metadata:
  name: "{{ include "sng.fullname" . }}-test-connection"
  labels:
    {{- include "sng.labels" . | nindent 4 }}
  annotations:
    "helm.sh/hook": test
spec:
  containers:
    - name: wget
      image: busybox
      command: ['wget']
      args: ['{{ include "sng.fullname" . }}:{{ .Values.service.port }}']
  restartPolicy: Never
`
