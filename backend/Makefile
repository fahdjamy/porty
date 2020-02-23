# note: call scripts from /scripts

BUILDPATH=&{CURDIR}
GOINSTALL=${GO} install
GOCLEAN=${GO} clean
GOGET=${GO} get

default:
	@echo "=============building Local API============="
	docker build -f cmd/app/main.go

up: default
	@echo "=============starting api locally============="
	docker-compose up -d

setdevenv:
	@echo "============= setting MODE to DEV_MODE  ============="
	export MODE=DEV_MODE
	@echo MODE is ==> $$MODE

start:
	CompileDaemon -build="go build cmd/app/main.go" -command="./main" -include=Makefile -include="*.less" -include="*.tmpl"

rundocker:
	docker run -p 8080:8008 -it raedar/backend

runtest:
	go test -v -cover ./...

settestenv:
	@echo "============= setting MODE to TESTING_MODE  ============="
	export MODE=TESTING_MODE

testapp: settestenv runtest

compile:
	GOOS=freebsd GOARCH=386 go build -o bin/main-freebsd-386 cmd/app/main.go
	GOOS=linux GOARCH=386 go build -o bin/main-linux-386 cmd/app/main.go
	GOOS=windows GOARCH=386 go build -o bin/main-windows-386 cmd/app/main.go

clean: down
	@echo "=============cleaning up============="
	rm -f api
	docker system prune -f
	docker volume prune -f
