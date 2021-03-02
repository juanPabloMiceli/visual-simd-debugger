install: microjail
	sudo cp microjail /usr/bin/

microjail: microjail.c
	gcc microjail.c -o microjail

clean:
	rm microjail