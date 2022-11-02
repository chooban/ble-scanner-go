debbuild: build
	cp ./ble-scanner-go ble-scanner/usr/local/sbin/
	dpkg-deb --build ble-scanner

build:
	GOOS=linux GOARCH=arm GOARM=5 go build
