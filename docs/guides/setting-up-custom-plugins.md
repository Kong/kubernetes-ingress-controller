# Setting up custom plugin in Kubernetes environment

This guide goes through steps on installing a custom plugin
in Kong without using a Docker build.

## Prepare a directory with plugin code

First, we need to create either a ConfigMap or a Secret with
the plugin code inside it.
If you would like to install a plugin which is available as
a rock from Luarocks, then you need to download it, unzip it and create a
ConfigMap from all the Lua files of the plugin.

We are going to setup a dummy plugin next.
If you already have a real plugin, you can skip this step.

```shell
$ mkdir myheader && cd myheader
$ echo 'local MyHeader = {}

MyHeader.PRIORITY = 1000

function MyHeader:header_filter(conf)
  -- do custom logic here
  kong.response.set_header("myheader", conf.header_value)
end

return MyHeader
' > handler.lua

$ echo 'return {
  name = "myheader",
  fields = {
    { config = {
        type = "record",
        fields = {
          { header_value = { type = "string", default = "roar", }, },
        },
    }, },
  }
}
' > schema.lua
```

Once we have our plugin code available in a directory,
the directory should look something like this:

```shell
$ tree myheader 
myheader
├── handler.lua
└── schema.lua

0 directories, 2 files
```

You might have more files inside the directory as well.

## Create a ConfigMap or Secret with the plugin code

Next, we are going to create a ConfigMap or Secret based on the plugin
code.

Please ensure that this is created in the same namespace as the one
in which Kong is going to be installed.

```shell
# using ConfigMap; replace `myheader` with the name of your plugin
$ kubectl create configmap kong-plugin-myheader --from-file=myheader -n kong
configmap/kong-plugin-myheader created

# OR using Secret
$ kubectl create secret generic -n kong kong-plugin-myheader --from-file=myheader
secret/kong-plugin-myheader created
```

## Modify configuration

Next, we need to update Kong's Deployment to load our custom plugin.

Based on your installation method, this step will differ slightly.
The next section explains what changes are necessary.

### YAML

The following patch is necessary to load the plugin.
Notable changes:
- The plugin code is mounted into the pod via `volumeMounts` and `volumes`
  configuration property.
- `KONG_PLUGINS` environment variable is set to include the custom plugin
  alongwith all the plugins that come in Kong by default.
- `KONG_LUA_PACKAGE_PATH` environment variable directs Kong to look
  for plugins in the directory where we are mounting them.

If you have multiple plugins, simply mount multiple
ConfigMaps and include the plugin name in the `KONG_PLUGINS`
environment variable.
  
> Please note that if your plugin code involves database
  migration then you need to include the below patch to pod definition of your
  migration Job as well.

Please note that the below is not a complete definition of
the Deployment but merely a strategic patch which can be applied to
an existing Deployment.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ingress-kong
  namespace: kong
spec:
  template:
    spec:
      containers:
      - name: proxy
        env:
        - name: KONG_PLUGINS
          value: bundled,myheader
        - name: KONG_LUA_PACKAGE_PATH
          value: "/opt/?.lua;;"
        volumeMounts:
        - name: kong-plugin-myheader
          mountPath: /opt/kong/plugins/myheader
      volumes:
      - name: kong-plugin-myheader
        configMap:
          name: kong-plugin-myheader
```

This is also available as a Kustomization:

```shell
$ kustomize build github.com/hbagdi/yaml/kong/kong-custom-plugin
```

### Helm chart

With Helm, this is as simple as adding the following values to
your `values.yaml` file:

```yaml
# values.yaml
plugins:
  configMaps:                # change this to 'secrets' if you created a secret
  - name: kong-plugin-myheader
    pluginName: myheader
```

The chart automatically configures all the environment variables based on the
plugins you inject.

Please ensure that you add in other configuration values
you might need for your installation to work.

### Deploy

Once, you have all the pieces in place, you are ready
to deploy Kong Ingress Controller:

```shell
# using YAML or kustomize
kustomize build github.com/hbagdi/yaml/kong/kong-custom-plugin | kubectl apply -f -

# or helm
$ helm repo add kong https://charts.konghq.com
$ helm repo update

# Helm 2
$ helm install kong/kong --values values.yaml

# Helm 3
$ helm install kong/kong --generate-name --set ingressController.installCRDs=false --values values.yaml
```

Once you have got Kong up and running, configure your
custom plugin via [KongPlugin resource](using-kongplugin-resource.md).
