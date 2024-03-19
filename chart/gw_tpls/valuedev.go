package gw_tpls

const ValueDevTpl = `image:
  repository: "dev-harbor.k7.cn/sng/{{.name}}"

ingress:
  enabled: false

ingress_haproxy:
  enabled: true
  annotations: 
    kubernetes.io/ingress.class: haproxy
  hosts:
    - host: dev-{{.name}}.svc.qipai007cs.com
      paths: 
        - path: /
          pathType: Prefix
  tls: []

resources:
  limits:
    cpu: 3
  requests:
    cpu: 1
    memory: 1G

nodeSelector: {}
`
