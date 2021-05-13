REGISTRY = yametech
VERSION = 0.2.0

default:
	docker build -t ${REGISTRY}/global-ipam:${VERSION} .
	docker push ${REGISTRY}/global-ipam:${VERSION}