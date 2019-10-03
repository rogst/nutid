VERSION := "v0.1.0"

clean:
	rm -rf releases/

build:
	go build

release: clean
	mkdir -p releases/nutid-${VERSION}-linux-amd64
	mkdir -p releases/nutid-${VERSION}-darwin-amd64
	cp LICENSE README.md releases/nutid-${VERSION}-linux-amd64/
	cp LICENSE README.md releases/nutid-${VERSION}-darwin-amd64/
	GOOS=linux go build -o releases/nutid-${VERSION}-linux-amd64/nutid
	GOOS=darwin go build -o releases/nutid-${VERSION}-darwin-amd64/nutid
	cd releases && tar -czf nutid-${VERSION}-linux-amd64.tar.gz nutid-${VERSION}-linux-amd64
	cd releases && tar -czf nutid-${VERSION}-darwin-amd64.tar.gz nutid-${VERSION}-darwin-amd64