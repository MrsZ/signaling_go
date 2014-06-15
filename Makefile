# constants
#

-include Makefile.def

# targets
#

all: clean get test compile run
werker: get test compile

get:
ifndef GOPATH
	@echo "GOPATH should be specified"
	@exit 2
else
	@echo Get ...
	@go get $(PACKAGES)
endif

test:
	@echo Tests ...
	@go test -v $(PACKAGES)

compile:
	@echo Compile ...
	@go build $(APP_MAIN)

run:
	@echo Run ...
	@./$(APP_NAME)

clean:
	@echo Clean ...
	@rm -f $(APP_NAME)
