# Version Compatibility

Kong's Ingress Controller is compatible with different flavors of Kong.
The following sections detail on compatibility between versions.

## Kong

By Kong, we are here referring to the official distribution of the Open-Source
Kong Gateway.

| Kong Ingress Controller  | <= 0.0.4           | 0.0.5              | 0.1.x              | 0.2.x              | 0.3.x              | 0.4.x              | 0.5.x              | 0.6.x              | 0.7.x              | 0.8.x              |
|--------------------------|:------------------:|:------------------:|:------------------:|:------------------:|:------------------:|:------------------:|:------------------:|:------------------:|:------------------:|:------------------:|
| Kong 0.13.x              | :white_check_mark: | :white_check_mark: | :white_check_mark: | :x:                | :x:                | :x:                | :x:                | :x:                | :x:                | :x:                |
| Kong 0.14.x              | :x:                | :x:                | :x:                | :white_check_mark: | :x:                | :x:                | :x:                | :x:                | :x:                | :x:                |
| Kong 1.0.x               | :x:                | :x:                | :x:                | :x:                | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | :x:                |
| Kong 1.1.x               | :x:                | :x:                | :x:                | :x:                | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | :x:                |
| Kong 1.2.x               | :x:                | :x:                | :x:                | :x:                | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Kong 1.3.x               | :x:                | :x:                | :x:                | :x:                | :x:                | :x:                | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Kong 1.4.x               | :x:                | :x:                | :x:                | :x:                | :x:                | :x:                | :x:                | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Kong 2.0.x               | :x:                | :x:                | :x:                | :x:                | :x:                | :x:                | :x:                | :x:                | :white_check_mark: | :white_check_mark: |

## Kong-enterprise-k8s

Kong-enterprise-k8s is an official distribution by Kong, Inc. which bundles
all enterprise plugins into Open-Source Kong Gateway.

The compatibility for this distribution will largely follow that of the
Open-Source Kong Gateway compatibility (the previous section).

| Kong Ingress Controller     | 0.6.2+             | 0.7.x              | 0.8.x              |
|-----------------------------|:------------------:|:------------------:|:------------------:|
| Kong-enterprise-k8s 1.3.x.y | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Kong-enterprise-k8s 1.4.x.y | :white_check_mark: | :white_check_mark: | :white_check_mark: |

## Kong Enterprise

Kong Enterprise is the official enterprise distribution, which includes all
other enterprise functionality, built on top of the Open-Source Kong Gateway.

| Kong Ingress Controller  | 0.0.5              | 0.1.x              | 0.2.x              | 0.3.x              | 0.4.x              | 0.5.x              | 0.6.x              | 0.7.x              | 0.8.x              |
|--------------------------|:------------------:|:------------------:|:------------------:|:------------------:|:------------------:|:------------------:|:------------------:|:------------------:|:------------------:|
| Kong Enterprise 0.32-x   | :white_check_mark: | :white_check_mark: | :x:                | :x:                | :x:                | :x:                | :x:                | :x:                | :x:                |
| Kong Enterprise 0.33-x   | :white_check_mark: | :white_check_mark: | :x:                | :x:                | :x:                | :x:                | :x:                | :x:                | :x:                |
| Kong Enterprise 0.34-x   | :white_check_mark: | :white_check_mark: | :x:                | :x:                | :x:                | :x:                | :x:                | :x:                | :x:                |
| Kong Enterprise 0.35-x   | :x:                | :x:                | :x:                | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | :x:                |
| Kong Enterprise 0.36-x   | :x:                | :x:                | :x:                | :x:                | :x:                | :x:                | :white_check_mark: | :white_check_mark: | :white_check_mark: |
| Kong Enterprise 1.3.x    | :x:                | :x:                | :x:                | :x:                | :x:                | :x:                | :white_check_mark: | :white_check_mark: | :white_check_mark: |
