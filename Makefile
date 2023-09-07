install-swagger:
	go install github.com/go-swagger/go-swagger/cmd/swagger@latest

swagger: install-swagger
	swagger generate spec -o ./swagger.yaml --scan-models

serve-swagger: swagger
	swagger serve -F=swagger swagger.yaml