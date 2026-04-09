.PHONY: build clean install

APP_NAME = vaultshell

build:
	pip install pyinstaller
	pyinstaller --onefile --name $(APP_NAME) vaultshell.py

clean:
	rm -rf build/ dist/ $(APP_NAME).spec
	rm -rf __pycache__

install: build
	sudo cp dist/$(APP_NAME) /usr/local/bin/$(APP_NAME)
	sudo chown root:docker /usr/local/bin/$(APP_NAME)
	sudo chmod 2755 /usr/local/bin/$(APP_NAME)
	sudo mkdir -p /var/lib/vaultshell
	sudo chown root:docker /var/lib/vaultshell
	sudo chmod 770 /var/lib/vaultshell
	sudo touch /var/lib/vaultshell/vaultshell.db
	sudo chown root:docker /var/lib/vaultshell/vaultshell.db
	sudo chmod 660 /var/lib/vaultshell/vaultshell.db
	sudo mkdir -p /var/lib/vaultshell/logs
	sudo chown root:docker /var/lib/vaultshell/logs
	sudo chmod 770 /var/lib/vaultshell/logs
	sudo mkdir -p /var/lib/vaultshell/logstxt
	sudo chown root:docker /var/lib/vaultshell/logstxt
	sudo chmod 770 /var/lib/vaultshell/logstxt
	sudo cp docker-wrapper.sh /usr/local/bin/docker
	sudo chmod +x /usr/local/bin/docker
	sudo cp watch.sh /usr/local/bin/watch.sh
	sudo cp logconverter.sh /usr/local/bin/logconverter.sh
	sudo chmod +x /usr/local/bin/watch.sh
	sudo chmod +x /usr/local/bin/logconverter.sh
