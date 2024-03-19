package gw_tpls

const ValueYeetownTpl = `replicaCount: 2

image:
  repository: "harbor-intranet.yeetown.cc/sng/{{.name}}"
  pullPolicy: IfNotPresent
  tag: ""

service:
  create: true
  headless: false
  type: ClusterIP
  port: 8000
  mgmtport: 8080

serviceMonitor:
  enabled: true

ingress:
  enabled: false

ingress_haproxy:
  enabled: false

resources: 
  limits:
    cpu: 2
    memory: 2Gi
  requests:
    cpu: 100m
    memory: 128Mi

nodeSelector: 

affinity: 
`
