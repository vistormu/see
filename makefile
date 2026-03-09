PROJECT_NAME := see
DIST_DIR := dist
VERSION := 0.0.5
ROOT_DIR := $(CURDIR)

# Force module mode so local/global go.work is ignored.
GO_WORK := off
# Keep cache inside the repo for reproducible/local CI runs.
GO_CACHE_DIR ?= $(ROOT_DIR)/.gocache
GO_ENV := GOWORK=$(GO_WORK) GOCACHE=$(GO_CACHE_DIR)


.PHONY: build release upload install clean clean-cache

build: clean
	# linux
	$(GO_ENV) GOOS=linux GOARCH=arm64 go build -o $(DIST_DIR)/$(PROJECT_NAME)
	tar -czf $(DIST_DIR)/$(PROJECT_NAME)_$(VERSION)_linux_arm64.tar.gz -C $(DIST_DIR) $(PROJECT_NAME)
	mv $(DIST_DIR)/$(PROJECT_NAME) $(DIST_DIR)/$(PROJECT_NAME)_linux_arm64

	$(GO_ENV) GOOS=linux GOARCH=amd64 go build -o $(DIST_DIR)/$(PROJECT_NAME)
	tar -czf $(DIST_DIR)/$(PROJECT_NAME)_$(VERSION)_linux_amd64.tar.gz -C $(DIST_DIR) $(PROJECT_NAME)
	mv $(DIST_DIR)/$(PROJECT_NAME) $(DIST_DIR)/$(PROJECT_NAME)_linux_amd64

	# darwin
	$(GO_ENV) GOOS=darwin GOARCH=arm64 go build -o $(DIST_DIR)/$(PROJECT_NAME)
	tar -czf $(DIST_DIR)/$(PROJECT_NAME)_$(VERSION)_darwin_arm64.tar.gz -C $(DIST_DIR) $(PROJECT_NAME)
	mv $(DIST_DIR)/$(PROJECT_NAME) $(DIST_DIR)/$(PROJECT_NAME)_darwin_arm64

	$(GO_ENV) GOOS=darwin GOARCH=amd64 go build -o $(DIST_DIR)/$(PROJECT_NAME)
	tar -czf $(DIST_DIR)/$(PROJECT_NAME)_$(VERSION)_darwin_amd64.tar.gz -C $(DIST_DIR) $(PROJECT_NAME)
	mv $(DIST_DIR)/$(PROJECT_NAME) $(DIST_DIR)/$(PROJECT_NAME)_darwin_amd64


release:
	git tag -a v$(VERSION) -m "release v$(VERSION)"
	git push origin v$(VERSION)

upload: release

install: build
	sudo cp $(DIST_DIR)/$(PROJECT_NAME)_darwin_arm64 /usr/local/bin/$(PROJECT_NAME)
	sudo chmod +x /usr/local/bin/$(PROJECT_NAME)

clean:
	rm -rf $(DIST_DIR)

clean-cache:
	rm -rf $(GO_CACHE_DIR)
