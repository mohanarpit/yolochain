build: ## Build the program
	go build

protogen: ## Generate the protobuf code
	/Users/arpit/code/protobuf/bin/protoc -I ./protofiles service.proto --go_out=plugins=grpc:blockchainGrpc

run: ## Run the program in default mode (POS)
	go run main.go -mode=pos

help: ## This help dialog.
	@IFS=$$'\n' ; \
	help_lines=(`fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//'`); \
	for help_line in $${help_lines[@]}; do \
		IFS=$$'#' ; \
		help_split=($$help_line) ; \
		help_command=`echo $${help_split[0]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
		help_info=`echo $${help_split[2]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
		printf "%-30s %s\n" $$help_command $$help_info ; \
	done
