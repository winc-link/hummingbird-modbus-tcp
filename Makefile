.PHONY: build

build:
	docker buildx build --platform linux/amd64 -t 'registry.cn-shanghai.aliyuncs.com/winc-driver/modbus-tcp-driver:2.7.1' -f docker/Dockerfile . --push

