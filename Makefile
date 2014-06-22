# constants
#

-include Makefile.def

# targets
#

all: clean get test compile run
werker: get test compile
reload: compile run

get:
ifndef GOPATH
	@echo "$(ERROR_COLOR)GOPATH should be specified$(NO_COLOR)"
	@exit 2
else
	@echo "$(OK_COLOR)==> Get$(NO_COLOR)"
	@go get $(PACKAGES)
	-@go get
endif

test:
	@echo "$(OK_COLOR)==> Tests$(NO_COLOR)"
	@go test -v $(PACKAGES)

compile:
	@echo "$(OK_COLOR)==> Compile$(NO_COLOR)"
	@go build -ldflags "-X main.Build $(VERSION)" $(APP_MAIN)

run:
	@echo "$(OK_COLOR)==> Run$(NO_COLOR)"
	@./$(APP_NAME)

clean:
	@echo "$(OK_COLOR)==> Clean$(NO_COLOR)"
	@rm -f $(APP_NAME)
