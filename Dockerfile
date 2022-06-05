FROM scratch

COPY assets/manifest.yaml .
COPY assets/iscsiadm-nsenter.sh /rootfs/usr/local/sbin/iscsiadm
