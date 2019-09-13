.PHONY: davai comin bd

davai:
	go run main.go

comin:
	docker exec -it reporter_mongo_1 mongo -u who -p dat

bd:
	docker-compose -f docker-compose.yaml up -d --build

sd:
	docker-compose up -d