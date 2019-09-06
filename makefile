.PHONY: davai comin

davai:
	go run main.go

comin:
	docker exec -it reporter_mongo_1 mongo -u who -p dat