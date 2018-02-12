default: dev

dev: test
	go build
	cd ${GOPATH}/src/github.com/xtracdev/automated-perf-test/ui-src && \
		ng build
	./automated-perf-test -ui

e2e:
	go build
	cd ${GOPATH}/src/github.com/xtracdev/automated-perf-test/ui-src && \
		ng build
	../automated-perf-test -ui&
	cd ${GOPATH}/src/github.com/xtracdev/automated-perf-test/ui-src && npm run e2e

prod: test
	go build
	cd ${GOPATH}/src/github.com/xtracdev/automated-perf-test/ui-src && \
		ng build --prod
	./automated-perf-test -ui

test:
	cd ${GOPATH}/src/github.com/xtracdev/automated-perf-test/ && \
		go fmt ./... && \
		go test ./... && \
		go vet ./... && \
		cd uiServices/test/ && \
		godog *.feature