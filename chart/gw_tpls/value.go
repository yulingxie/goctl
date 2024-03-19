package gw_tpls

const ValueTpl = `# Default values for sng.
# This is a YAML-formatted file.
# Declare variables to be passed into your templates.

serviceName: {{.name}}

replicaCount: 2

image:
  repository: "harbor-10001.k7.cn/sng/{{.name}}"
  pullPolicy: IfNotPresent
  # Overrides the image tag whose default is the chart appVersion.
  tag: ""

imagePullSecrets: 
   - name: harbor-10001-pullonly
nameOverride: "{{.name}}"
fullnameOverride: "sng-{{.name}}"

serviceAccount:
  # Specifies whether a service account should be created
  create: true
  # Annotations to add to the service account
  annotations: {}
  # The name of the service account to use.
  # If not set and create is true, a name is generated using the fullname template
  name: ""

env: []

envFrom: 
  - configMapRef:
      name: sng-config

podAnnotations: {}

podSecurityContext: {}
  # fsGroup: 2000

securityContext: {}
  # capabilities:
  #   drop:
  #   - ALL
  # readOnlyRootFilesystem: true
  # runAsNonRoot: true
  # runAsUser: 1000

service:
  create: true
  headless: true
  type: ClusterIP
  port: 8000
  mgmtport: 8080

serviceMonitor:
  enabled: true
  
ingress:
  enabled: true
  annotations: 
    # kubernetes.io/ingress.class: nginx
    # kubernetes.io/tls-acme: "true"
    nginx.ingress.kubernetes.io/configuration-snippet: |
      proxy_ignore_client_abort  on;
  hosts:
    - host: {{.name}}.svc.qipai007cs.com
      paths: 
        - path: /
          backend:
            service: sng-{{.name}}
            port:
              number: 8000
  tls: []
  #  - secretName: chart-example-tls
  #    hosts:
  #      - chart-example.local

ingress_haproxy:
  enabled: true
  annotations: 
    kubernetes.io/ingress.class: haproxy
    haproxy-ingress.github.io/allowlist-source-range: 10.0.0.0/8, 192.168.1.0/24
  hosts:
    - host: {{.name}}.svc.qipai007cs.com
      paths: 
        - path: /
          pathType: Prefix
  tls: []
  #  - secretName: chart-example-tls
  #    hosts:
  #      - chart-example.local

resources: 
  limits:
    cpu: 3
    memory: 2Gi
  requests:
    cpu: 300m
    memory: 128Mi

autoscaling:
  enabled: false
  minReplicas: 1
  maxReplicas: 100
  targetCPUUtilizationPercentage: 80
  # targetMemoryUtilizationPercentage: 80

nodeSelector: 
  zone: zw

tolerations: []

affinity: {}
`
