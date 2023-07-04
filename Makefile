BINARY_NAME=MarkDown.app
APP_NAME=MarkDown
VERSION=1.0.0

## build: build binary and package app
build:
	rm -rf ${BINARY_NAME}
	rm -f fynemd
	fyne package -appVersion ${VERSION} -name ${APP_NAME} -release

## run: builds and runs the applicationrun:
run:
	go run .

## clean: runs go clean and deletes binariesclean:
clean:
	@echo "Cleaning..."
	@go clean@rm -rf ${BINARY_NAME}
	@echo "Cleaned!"

## test: runs all tests
test:
	go test -v ./...