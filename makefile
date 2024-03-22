# Makefile

## variables
build_dir := build# define the build directory

build:
	
# ensure go is present
	@echo "Checking if go is installed"
	@if ! command -v go > /dev/null; then \
		echo "Go is not installed. Please install go."; \
		exit 1; \
	fi

# create build directory
	@echo "Creating build directory"
	@if [ ! -d "$(build_dir)" ]; then \
        mkdir -p "$(build_dir)"; \
    fi

# build the binary
	@echo "Building proxmox-cli"
	@go build -o "$(build_dir)/proxmox-cli" . >> "$(build_dir)/build.log" 2>&1

# copying environment file
	@echo "Copying environment file"
	@cp environments.json $(build_dir)/environments.json

# # check if the binary was built successfully
	@echo "Checking if the binary was built successfully"
# 	@if  ! -f $(build_dir)/proxmox-cli ; then \
# 		echo "Binary was not built successfully. Please check the logs."; \
# 		exit 1; \
# 	fi
	@echo "Binary was built successfully. You can find it at $(build_dir)/proxmox-cli"

createConfig:
# check if the user is running as sudo
	@make checkSudo

# create the config directory
	@echo "Creating config directory"
	@mkdir -p /etc/proxmox-cli

# copy the environments.json file
	@echo "Copying environments.json file"
	@cp environments.json /etc/proxmox-cli/environments.json

	@echo "Config file was created successfully"

checkSudo:
# confirm the user is running as sudo
	@echo "Checking for sudo"
	@if [ `id -u` -ne 0 ]; then \
        echo "Are you running sudo?"; \
        exit 1; \
    fi

install: 
	@echo "Installing proxmox-cli"
	@if [ ! -f "$(build_dir)/proxmox-cli" ] ; then \
        echo "You need to run make build first"; \
        exit 1; \
    fi

# check if the user is running as sudo
	@make checkSudo

# create the config directory
	@make createConfig

# copy the binary to /usr/local/bin
	@cp "$(build_dir)/proxmox-cli" /usr/local/bin/proxmox-cli
# make the binary executable
	@chmod +x /usr/local/bin/proxmox-cli

	@echo "proxmox-cli was installed successfully"

clean: 
	@echo "Cleaning up build directory"
	@rm -rf $(build_dir)

uninstall:

# check if the user is running as sudo
	@make checkSudo

# remove the binary from /usr/local/bin
	@echo "Uninstalling proxmox-cli"
	@rm -f /usr/local/bin/proxmox-cli

.PHONY: install build clean

# End of makefile