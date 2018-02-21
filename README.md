# screen_share_remote_go
This program is meant to be used together with [screen_share](https://github.com/rootkiwi/screen_share/).

This program can be started on for example your server somewhere, then you share your screen to this program
and this program can forward to connected web clients.

## Why rewrite from Java?
This program is a rewrite of the Java program [screen_share_remote](https://github.com/rootkiwi/screen_share_remote/).

I recently started looking into Go (I love it) and felt that this program was a great fit for Go.

## How to migrate from Java version
* You need to generate a new config file. They are almost the same except bcrypt is used instead of Argon2 for
password hashing.

* And if you're using a reverse proxy the Host header need to be forwarded as well. Because Gorilla WebSocket defaults
to allow only if Host and Origin headers matches.

## Usage
```
Usage: screen_share_remote_go-<VERSION>-<PLATFORM> [/path/to/conf | genconf | noconf]

Example: ./screen_share_remote_go-<VERSION>-<PLATFORM> genconf

Available command-line options:
1. /path/to/screen_share_remote.conf (start screen_share_remote)
2. genconf                           (generate new config file)
3. noconf                            (start without saving config)
```

## Example generate config
```
$ ./screen_share_remote_go-<VERSION>-<PLATFORM> genconf
Generate screen_share_remote.conf file in working directory which is:
<working_directory>

Will overwrite if already exists

The config file will contain these attributes:
1. port number                 (port number to enter in screen_share)
2. web server port number      (port the web server will serve on)
3. password                    (password to enter in screen_share)
4. self-signed TLS certificate (whose fingerprint to enter in screen_share)
5. RSA private key             (corresponding to certificate)

Do note that the RSA private key is stored in cleartext, so make sure
to make the config file inaccessible for unauthorized parties.
Or you could run in the 'noconf' mode, which means a new private key will be
generated each time. Without saving to disk.

Leave empty and press ENTER for the [default] value
1. enter port number (0-65535) [50000]: 
2. enter web server port number (0-65535) [8081]: 
3. enter password [random]: 

password:
B9p2aPreTRh6ya8fHYtUr5JbMDVCW6veRgzJSbiz

Generating a 4096-bit RSA key pair and a self-signed certificate... done.

certificate fingerprint:
13BEE4D234B31D8EE09542639FA98EEFB3DADCBE39AED4DEAF21C447971F3167

Config file created:
<working_directory>/screen_share_remote.conf

The settings 'port' and 'webPort' is changeable, the rest is not
If you need to change the password/certificate run genconf again
```

## Binaries
[screen_share_remote_go/releases/latest](https://github.com/rootkiwi/screen_share_remote_go/releases/latest)

## Example NGINX reverse proxy
You may want to set up `screen_share_remote_go` behind a reverse proxy like NGINX.
WebSockets is used so the reverse proxy need to be set up for that. Host header should be forwarded as well.
Following is an example location config block in NGINX.
```
location / {
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_pass http://localhost:8081;
}
```

## How to build
Docker is used for building ([Dockerfile](https://github.com/rootkiwi/screen_share_remote_go/blob/master/Dockerfile))
so make sure you have docker installed.

For dependency management [dep](https://github.com/golang/dep) is used, but code is checked out in `vendor/` so dep is
not required for building.

Run following to output in `bin/release/<VERSION>`
```
./build.sh release
```

Run following to output in `bin/dev`
```
./build.sh
```

## License
[GNU General Public License 3 or later](https://www.gnu.org/licenses/gpl-3.0.html)

See LICENSE for more details.

## 3RD Party Dependencies

### Gorilla WebSocket

[https://github.com/gorilla/websocket](https://github.com/gorilla/websocket)

[2-clause BSD License](https://github.com/gorilla/websocket/blob/master/LICENSE)


### Packr

[https://github.com/gobuffalo/packr](https://github.com/gobuffalo/packr)

[MIT License](https://github.com/gobuffalo/packr/blob/master/LICENSE.txt)


### Broadway

[https://github.com/mbebenita/Broadway/](https://github.com/mbebenita/Broadway/)

[3-clause BSD License](https://github.com/mbebenita/Broadway/blob/master/LICENSE)
