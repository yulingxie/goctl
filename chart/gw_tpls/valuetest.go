package gw_tpls

const ValueTestTpl = `image:
repository: "dev-harbor.k7.cn/sng/{{.name}}"

ingress_haproxy:
enabled: true
annotations: 
  kubernetes.io/ingress.class: haproxy
  haproxy-ingress.github.io/allowlist-source-range: ""
hosts:
  - host: test-{{.name}}.svc.qipai007cs.com
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
