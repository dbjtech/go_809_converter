.PHONY: build_base build_master push build_and_push all
build_base:
	docker build -t registry.cn-qingdao.aliyuncs.com/dbjtech/go_809_converter:go_ent -f Dockerfile_base .

build_master:
	docker build -t registry.cn-qingdao.aliyuncs.com/dbjtech/go_809_converter:master -f Dockerfile .

push:
	docker push registry.cn-qingdao.aliyuncs.com/dbjtech/go_809_converter:master

build_and_push:
	build_master push

all:
	build_and_push