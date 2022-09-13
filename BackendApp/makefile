JIKONI_DOCKER_IMAGE_NAME_PREFIX ?= rodneydav
BUILD_DIR = build
SERVICES = orders
DOCKERS = $(addprefix docker_,$(SERVICES))
DOCKERS_DEV = $(addprefix docker_dev_,$(SERVICES))
CGO_ENABLED ?= 0
GOARCH ?= amd64

define compile_service
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) \
	go build -mod=vendor -ldflags "-s -w" -o ${BUILD_DIR}/jikoni-$(1) cmd/$(1)/main.go
endef

define make_docker
	$(eval svc=$(subst docker_,,$(1)))

	docker build \
		--no-cache \
		--build-arg SVC=$(svc) \
		--build-arg GOARCH=$(GOARCH) \
		--build-arg GOARM=$(GOARM) \
		--tag=$(JIKONI_DOCKER_IMAGE_NAME_PREFIX)/jikoni-$(svc) \
		-f docker/Dockerfile .
endef

define make_docker_dev
	$(eval svc=$(subst docker_dev_,,$(1)))

	docker build \
		--no-cache \
		--build-arg SVC=$(svc) \
		--tag=$(JIKONI_DOCKER_IMAGE_NAME_PREFIX)/jikoni-$(svc) \
		-f docker/Dockerfile.dev ./build
endef

proto:
	protoc --gofast_out=plugins=grpc:. auth/*.proto

all: $(SERVICES)

.PHONY: all $(SERVICES) dockers dockers_dev latest release

start: all $(SERVICES) dockers_dev run

clean:
	rm -rf ${BUILD_DIR}

cleandocker:
	# Stops containers and removes containers, networks, volumes, and images created by up
	docker-compose -f docker/docker-compose.yml down -v --remove-orphans
	# Remove unused volumes
	docker volume ls -f name=$(JIKONI_DOCKER_IMAGE_NAME_PREFIX) -f dangling=true -q | xargs -r docker volume rm

install:
	cp ${BUILD_DIR}/* $(GOBIN)

$(SERVICES):
	$(call compile_service,$(@))

$(DOCKERS):
	$(call make_docker,$(@),$(GOARCH))

$(DOCKERS_DEV):
	$(call make_docker_dev,$(@))

dockers: $(DOCKERS)
dockers_dev: $(DOCKERS_DEV)

run:
	# docker-compose -f docker/docker-compose.yml build --no-cache
	docker-compose -f docker/docker-compose.yml up
	# docker-compose exec jikoni-roach-0 /cockroach/cockroach sql --certs-dir=/certs --host=jikoni-roach-0 --execute="CREATE USER IF NOT EXISTS roach WITH PASSWORD 'roach';"
	# docker-compose exec jikoni-roach-0 /cockroach/cockroach sql --certs-dir=/certs --host=jikoni-roach-0 --execute="GRANT ADMIN TO roach;"