install:
	brew tap go-swagger/go-swagger
	brew install go-swagger

swagger:
	GO111MODULE=off swagger generate spec -o ./swagger.yaml --scan-models