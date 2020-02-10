.PHONY: davai comin bd

build-migrator:
	go build -o bin/migrator cmd/migrator/*

migrator: build-migrator
	bin/migrator $(dir) $(lv)

davai:
	go run main.go

comin:
	docker exec -it reporter_mongo_1 mongo -u who -p dat

bd:
	docker-compose -f docker-compose.yaml up -d --build

sd:
	docker-compose up -d