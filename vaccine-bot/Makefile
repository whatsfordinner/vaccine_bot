build-VaccineTweeter:
	export GOOS=linux
	go get ./...
	go build -o ${ARTIFACTS_DIR}/vaccine-bot
	cp diseases.json ${ARTIFACTS_DIR}