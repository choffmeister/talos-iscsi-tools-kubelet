```yaml
- op: add
  path: /machine/install/extensions
  value:
    - image: ghcr.io/siderolabs/iscsi-tools:v0.1.1
    - image: ghcr.io/diztortion/talos-iscsi-tools-kubelet:v0.1.0
- op: add
  path: /machine/kubelet/extraMounts
  value:
    - destination: /usr/local/sbin
      type: bind
      source: /usr/local/sbin
      options:
        - bind
        - rshared
        - rw
```
