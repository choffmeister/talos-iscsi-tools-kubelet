REGISTRY=ghcr.io
IMAGE=diztortion/talos-iscsi-tools-kubelet
REGISTRY_IMAGE=${REGISTRY}/${IMAGE}

.PHONY: all build push clean clean_all

all: build push clean

build:
	@echo Building ${REGISTRY_IMAGE}:${TAG}
	@cd iscsiadm-injector && CGO_ENABLED=0 go build -trimpath -o ../rootfs/usr/local/lib/containers/iscsi-tools-kubelet/iscsiadm-injector main.go
	@docker build --label "org.opencontainers.image.created=$(shell date --rfc-3339=seconds)" --label "org.opencontainers.image.version=$(TAG)" --build-arg TAG=${TAG} --tag ${REGISTRY_IMAGE}:${TAG} --pull .
	@docker images --filter label=name=${PACKAGE_NAME} --filter label=stage=builder --quiet | xargs --no-run-if-empty docker rmi

push:
	@echo Pushing ${REGISTRY_IMAGE}:${TAG}
	@docker push ${REGISTRY_IMAGE}:${TAG}

clean:
	@docker rmi ${REGISTRY_IMAGE}:${TAG}

clean_all:
	@docker images --filter "reference=${REGISTRY_IMAGE}" --quiet | xargs --no-run-if-empty docker rmi
	@docker images --filter label=name=${PACKAGE_NAME} --filter label=stage=builder --quiet | xargs --no-run-if-empty docker rmi
