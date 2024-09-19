This is a public pastebin server. You may upload files up to 10MB in size and
have them shared publicly. FTP is used to upload files to the server.

## Usage

Uploading is quite simple. Using your favourite FTP client you can `ftp
paste.nilsu.org`, `cd incoming` (the public upload directory), and `put
file.whatever` then run `sha1sum file.whatever` on your file and the public
address will be `paste.cf/yoursha1.extension`

To make uploading a bit quicker I wrote a [tiny
client](https://git.sr.ht/~kota/pcf) that FTPs the file to the server,
calculates the hash, and prints what the resulting url should be. You can
combine `pcf` with other programs to do cool stuff. Like take a screenshot,
upload it, and put the url in your clipboard. 

`scrot -q 85 /tmp/scrot.png && pcf /tmp/scrot.png | xclip -in -selection c`

## How

An FTP server configured to allow public uploads into the `incoming` directory.
That public FTP directory is watched with `incrond`. Incrond is configured to
run a script when new files are placed in the public FTP directory. The script
checks that the file is under 10Mb. Then it caculates the sha1sum of the file
and renames + moves the file into the web server's public directory.

Additionally the server uses cron to run another script which checks if the disk
usage is above a certain percentage. If so it deletes uploaded files based on
age * size.

With basic server knowledge you should be able to duplicate this for yourself
(and even have your own private version if you'd like). Then you can set your
`pcf` environment variable to your ip/hostname and you're golden!
