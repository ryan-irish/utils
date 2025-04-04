.PHONY: build install clean

build:
	go build -o standup

install: build
	@echo "Installing standup to /usr/local/bin/"
	@mv standup /usr/local/bin/ || \
		(echo "Failed to move binary. Try using sudo: sudo make install" && exit 1)

clean:
	@rm -f standup

# Default target
all: build 