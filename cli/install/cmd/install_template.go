package cmd

const postgresTemplate = `
`

const template = `
apiVersion: v1
kind: Namespace
metadata:
  name: {{ .Namespace }}

---

apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: kongplugins.configuration.konghq.com
spec:
  group: configuration.konghq.com
  version: v1
  scope: Namespaced
  names:
    kind: KongPlugin
    plural: kongplugins

---

apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: kongconsumers.configuration.konghq.com
spec:
  group: configuration.konghq.com
  version: v1
  scope: Namespaced
  names:
    kind: KongConsumer
    plural: kongconsumers

---

apiVersion: v1
kind: ServiceAccount
metadata:
  name: kong-serviceaccount
  namespace: {{ .Namespace }}

---

apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  name: kong-ingress-clusterrole
rules:
  - apiGroups:
      - ""
    resources:
      - endpoints
      - nodes
      - pods
      - secrets
    verbs:
      - list
      - watch
  - apiGroups:
      - ""
    resources:
      - nodes
    verbs:
      - get
  - apiGroups:
      - ""
    resources:
      - services
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - "extensions"
    resources:
      - ingresses
    verbs:
      - get
      - list
      - watch
  - apiGroups:
      - ""
    resources:
        - events
    verbs:
        - create
        - patch
  - apiGroups:
      - "extensions"
    resources:
      - ingresses/status
    verbs:
      - update

---

apiVersion: rbac.authorization.k8s.io/v1beta1
kind: Role
metadata:
  name: kong-ingress-role
  namespace: {{ .Namespace }}
rules:
  - apiGroups:
      - ""
    resources:
      - configmaps
      - pods
      - secrets
      - namespaces
    verbs:
      - get
  - apiGroups:
      - ""
    resources:
      - configmaps
    resourceNames:
      # Defaults to "<election-id>-<ingress-class>"
      # Here: "<ingress-controller-leader>-<nginx>"
      # This has to be adapted if you change either parameter
      # when launching the nginx-ingress-controller.
      - "ingress-controller-leader-nginx"
    verbs:
      - get
      - update
  - apiGroups:
      - ""
    resources:
      - configmaps
    verbs:
      - create
  - apiGroups:
      - ""
    resources:
      - endpoints
    verbs:
      - get

---

apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  name: kong-ingress-role-nisa-binding
  namespace: {{ .Namespace }}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: kong-ingress-role
subjects:
  - kind: ServiceAccount
    name: kong-serviceaccount
    namespace: {{ .Namespace }}

---

apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: kong-ingress-clusterrole-nisa-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: kong-ingress-clusterrole
subjects:
  - kind: ServiceAccount
    name: kong-serviceaccount
    namespace: {{ .Namespace }}

---

apiVersion: v1
kind: Service
metadata:
  name: kong-ingress-controller
  namespace: {{ .Namespace }}
spec:
  type: NodePort
  ports:
  - name: kong-admin
    port: 8001
    targetPort: 8001
    protocol: TCP
  selector:
    app: ingress-kong

---

apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  labels:
    app: ingress-kong
  name: kong-ingress-controller
  namespace: {{ .Namespace }}
spec:
  serviceAccountName: kong-serviceaccount
  selector:
    matchLabels:
      app: ingress-kong
  strategy:
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
    type: RollingUpdate
  template:
    metadata:
      annotations:
        # the returned metrics are related to the kong ingress controller not kong itself
        prometheus.io/port: "10254"
        prometheus.io/scrape: "true"
      labels:
        app: ingress-kong
    spec:
      initContainers:
      - name: kong-migration
        image: {{ .KongImage }}
        env:
          - name: KONG_PG_PASSWORD
            value: kong
          - name: KONG_PG_HOST
            value: postgres
        command: [ "/bin/sh", "-c", "kong migrations up" ]
      containers:
      - name: admin-api
        image: {{ .KongImage }}
        env:
          - name: KONG_VITALS
            value: 'off'
          - name: KONG_PG_PASSWORD
            value: kong
          - name: KONG_PG_HOST
            value: postgres
          - name: KONG_PROXY_ACCESS_LOG
            value: /dev/stdout
          - name: KONG_ADMIN_ACCESS_LOG
            value: /dev/stdout
          - name: KONG_PROXY_ERROR_LOG
            value: /dev/stderr
          - name: KONG_ADMIN_ERROR_LOG
            value: /dev/stderr
          - name: KONG_ADMIN_LISTEN
            value: 0.0.0.0:8001
          - name: KONG_PROXY_LISTEN
            value: 'off'
        ports:
        - name: kong-admin
          containerPort: 8001
        livenessProbe:
          failureThreshold: 3
          httpGet:
            path: /status
            port: 8001
            scheme: HTTP
          initialDelaySeconds: 60
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1
        readinessProbe:
          failureThreshold: 3
          httpGet:
            path: /status
            port: 8001
            scheme: HTTP
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1
      - name: ingress-controller
        args:
        - /kong-ingress-controller
        # the kong URL points to the kong admin api server
        - --kong-url=http://localhost:8001
        # the default service is the kong proxy service
        - --default-backend-service=kong/kong-proxy
        env:
        - name: POD_NAME
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: metadata.name
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: metadata.namespace
        image: {{ .KongIngressControllerImage }}
        imagePullPolicy: {{ .ImagePullPolicy }}
        livenessProbe:
          failureThreshold: 3
          httpGet:
            path: /healthz
            port: 10254
            scheme: HTTP
          initialDelaySeconds: 60
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1
        readinessProbe:
          failureThreshold: 3
          httpGet:
            path: /healthz
            port: 10254
            scheme: HTTP
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1

---

apiVersion: v1
kind: Service
metadata:
  name: kong-proxy
  namespace: {{ .Namespace }}
  annotations:
  {{ range $k, $v := .ServiceAnnotations }}
    {{ $k }}: {{ $v }}
  {{ end }}
  spec:
  externalTrafficPolicy: Local
  type: {{ .ServiceType }}
  ports:
  - name: kong-proxy
    port: 80
    targetPort: 8000
    protocol: TCP
  - name: kong-proxy-ssl
    port: 443
    targetPort: 8443
    protocol: TCP
  selector:
    app: kong

---

apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: kong
  namespace: {{ .Namespace }}
spec:
  template:
    metadata:
      labels:
        name: kong
        app: kong
    spec:
      containers:
      - name: kong-proxy
        image: {{ .KongImage }}
        env:
          - name: KONG_PG_PASSWORD
            value: kong
          - name: KONG_PG_HOST
            value: postgres
          - name: KONG_PROXY_ACCESS_LOG
            value: "/dev/stdout"
          - name: KONG_ADMIN_ACCESS_LOG
            value: "/dev/stdout"
          - name: KONG_PROXY_ERROR_LOG
            value: "/dev/stderr"
          - name: KONG_ADMIN_ERROR_LOG
            value: "/dev/stderr"
        ports:
        - name: proxy
          containerPort: 8000
          protocol: TCP
        - name: proxy-ssl
          containerPort: 8443
          protocol: TCP
        livenessProbe:
          failureThreshold: 3
          httpGet:
            path: /status
            port: 8001
            scheme: HTTP
          initialDelaySeconds: 60
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1
        readinessProbe:
          failureThreshold: 3
          httpGet:
            path: /status
            port: 8001
            scheme: HTTP
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1

---

`
