project_name := see
dist_dir := dist
version := 0.0.2


.SILENT:

.PHONY: build clean

build:
	# linux
	GOOS=linux GOARCH=arm64 go build -o $(dist_dir)/$(project_name)_linux_arm64 main.go
	GOOS=linux GOARCH=amd64 go build -o $(dist_dir)/$(project_name)_linux_amd64 main.go

	# darwin
	GOOS=darwin GOARCH=arm64 go build -o $(dist_dir)/$(project_name)_darwin_arm64 main.go
	GOOS=darwin GOARCH=amd64 go build -o $(dist_dir)/$(project_name)_darwin_amd64 main.go

upload: build
	git add .
	git commit -m "release $(version)"
	git tag -a $(version) -m "release $(version)"
	git push origin main
	git push origin --tags

install: build
	sudo cp $(dist_dir)/$(project_name)_darwin_arm64 /usr/local/bin/$(project_name)
	sudo chmod +x /usr/local/bin/$(project_name)

clean:
	rm -rf $(dist_dir)
