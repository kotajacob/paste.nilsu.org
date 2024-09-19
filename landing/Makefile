all:
	hugo --cleanDestinationDir --minify
	rsync -rdvP public/ paste.nilsu.org:/var/www/html/paste/
