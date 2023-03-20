mock_store:
	mockgen -source=internal/repository/repo_interface.go -destination=internal/repository/mocks/repo_mock.go -package=mocks

mock_cash:
	mockgen -source=internal/repository/cash/cash_interface.go -destination=internal/repository/mocks/cash_mock.go -package=mocks Cash

mock_usecase:
	mockgen -source=internal/usecase/usecase_interface.go -destination=internal/usecase/mocks/usecase_mock.go -package=mocks

mock_filestorage:
	mockgen -source=internal/filestorage/diskFileStorage.go -destination=internal/filestorage/mocks/filestorage_mock.go -package=mocks FileStorager

up:
	docker-compose up -d

down:
	docker-compose down

run:
	go run ./cmd/onlineShopBackend/main.go

swag:
	swag init -d ./internal/delivery -g delivery.go -o ./internal/delivery/swagger/docs

swag_fmt:
	swag fmt -d ./internal/delivery -g delivery.go

up-win:
	git clone https://github.com/ZavNatalia/gb-store.git&&cd gb-store&&git checkout feature/new-api&&cd ..&&copy front.Dockerfile .\gb-store&&cd ./gb-store&&docker build -f ./front.Dockerfile -t front:latest .&&cd ..&&RD /s/q .\gb-store&&docker-compose up -d

up-lin:
	git clone https://github.com/ZavNatalia/gb-store.git&&cd gb-store&&git checkout feature/new-api&&cd ..&&cp ./front.Dockerfile ./gb-store&&cd ./gb-store&&docker build -f ./front.Dockerfile -t front:latest .&&cd ..&&rm -R -f ./gb-store&&docker-compose up -d

test:
	go test ./... -v