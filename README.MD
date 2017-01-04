# GBackup

Backup utility build in Golang

## Docker test

* docker pull golang:1.6
* docker run -it --rm -v "$PWD":/usr/src/mygoapp -w /usr/src/mygoapp golang:1.6
* go run main.go