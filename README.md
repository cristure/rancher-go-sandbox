# rancher-go-sandbox

Login scaffold test for rancher on single-node Docker install.

The test expects a docker volume created pointing at a rancher backup. If you don't have one you can download one [here](https://drive.usercontent.google.com/download?id=1UxD91cgdjqRJ6p5PRHvPaoe4KYbLxpKd&export=download&authuser=0&confirm=t&uuid=cf9b5265-6777-4683-a366-528de22d5750&at=APZUnTVg3SwptgpBiq4jIIAWpU4t:1721565539878).

## Pre-requisites
1. Docker
2. Go, version >= 1.21

## Setup

1. The test expects a user with a password already defined. You can download my backup from [here](https://drive.usercontent.google.com/download?id=1UxD91cgdjqRJ6p5PRHvPaoe4KYbLxpKd&export=download&authuser=0&confirm=t&uuid=cf9b5265-6777-4683-a366-528de22d5750&at=APZUnTVg3SwptgpBiq4jIIAWpU4t:1721565539878).
2. Extract the `.tar.gz`
```azure
tar xvf <DOWNLOAD_FOLDER>.tar.gz
```
3. Create a volume from the respective backup
```
docker volume create --name my_test_volume --opt type=none --opt device=/home/cristu/rancher/var/lib/rancher --opt o=bind
```

NOTE: You can replace the name arg. I left it `my_test_volume` because it is hardcoded as such in the test. If you want another name, please change the volume name in the test too.

4. Go to the source repository and run the tests.
```
export SCHEME=https
go test -v ./...
```