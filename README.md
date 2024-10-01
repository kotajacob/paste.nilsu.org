# [paste.nilsu.org](https://paste.nilsu.org/)

This is the source code for a pastebin server. The code is very short and
readable, but lacks options. It should be relatively easy to tweak and use as
you'd like.

There are three separate projects in this repo. Each repo has it's own Makefile.

# cli
A basic cli tool which uploads a file to the pastebin server and prints the URL.
If a folder is given it will be zipped and then uploaded. The cli binary is
named `pcf` for historical reasons.

# landing
The landing page is a simple static site generated with hugo. Running make will
generate the html and css files into a directory called `public/`.

# upload
A simple http request handler. It recieves post requests on the /upload route
and stores files in a folder (configured with the -storage flag).

By default the upload server will bind to post `2016`, but this can be
configured with the `-addr` flag.

To use host your own pastebin server you'll need a reverse proxy / webserver
such as [caddy](https://caddyserver.com/). Here's an example config with caddy
that will host the static landing page and the upload handler:
```
paste.nilsu.org {
	reverse_proxy /upload localhost:2016
	root * /var/www/html/paste.nilsu.org
	file_server
}
```

I would also recommend creating an init script to run the upload server on boot.
Here's an example with openrc:
```
#!/sbin/openrc-run
supervisor=supervise-daemon

name="paste-upload"
description="paste-upload"

command=${command:-/usr/bin/paste-upload}
command_background=true
command_user="paste:paste"
command_args="-storage /var/www/html/paste.nilsu.org"

pidfile="/run/${RC_SVCNAME}.pid"
output_log="/var/log/paste/paste.log"
error_log="/var/log/paste/paste.err"

depend() {
	need net
	use dns logger netmount
}
```
