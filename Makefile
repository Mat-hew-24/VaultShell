# Makefile
.PHONY: build clean install

# Name of the final executable
APP_NAME = vaultshell

build:
	pip install pyinstaller
	pyinstaller --onefile --name $(APP_NAME) vaultshell.py

clean:
	rm -rf build/ dist/ $(APP_NAME).spec
	rm -rf __pycache__

install: build
	sudo cp dist/$(APP_NAME) /usr/local/bin/$(APP_NAME)
