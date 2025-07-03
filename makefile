project_name := see
dist_dir := dist
version := 0.0.4


.PHONY: build upload install clean

build: clean
	# linux
	GOOS=linux GOARCH=arm64 go build -o $(dist_dir)/$(project_name)_linux_arm64 main.go
	tar -czf $(dist_dir)/$(project_name)_$(version)_linux_arm64.tar.gz -C $(dist_dir) $(project_name)_linux_arm64

	GOOS=linux GOARCH=amd64 go build -o $(dist_dir)/$(project_name)_linux_amd64 main.go
	tar -czf $(dist_dir)/$(project_name)_$(version)_linux_amd64.tar.gz -C $(dist_dir) $(project_name)_linux_amd64

	# darwin
	GOOS=darwin GOARCH=arm64 go build -o $(dist_dir)/$(project_name)_darwin_arm64 main.go
	tar -czf $(dist_dir)/$(project_name)_$(version)_darwin_arm64.tar.gz -C $(dist_dir) $(project_name)_darwin_arm64

	GOOS=darwin GOARCH=amd64 go build -o $(dist_dir)/$(project_name)_darwin_amd64 main.go
	tar -czf $(dist_dir)/$(project_name)_$(version)_darwin_amd64.tar.gz -C $(dist_dir) $(project_name)_darwin_amd64


upload: build
	git add .
	git commit -m "release v$(version)"
	git tag -a v$(version) -m "release v$(version)"
	git push origin main
	git push origin --tags
	gh release create v$(version) $(dist_dir)/* --title "release v$(version)" --notes "release v$(version)"

install: build
	sudo cp $(dist_dir)/$(project_name)_darwin_arm64 /usr/local/bin/$(project_name)
	sudo chmod +x /usr/local/bin/$(project_name)

clean:
	rm -rf $(dist_dir)
