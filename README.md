# Pantegana - A Botnet RAT Made With Go
### <center>FOR EDUCATIONAL AND RESEARCH USE ONLY</center>  

```
   ___            _                               
  / _ \__ _ _ __ | |_ ___  __ _  __ _ _ __   __ _ 
 / /_)/ _` | '_ \| __/ _ \/ _` |/ _` | '_ \ / _` |
/ ___/ (_| | | | | ||  __/ (_| | (_| | | | | (_| |
\/    \__,_|_| |_|\__\___|\__, |\__,_|_| |_|\__,_|
                          |___/                   
```

## Features:
 - Pretty and clean interactive shell (using <a href="https://github.com/desertbit/grumble" target="_blank">grumble</a>)
 - Cross-platform payload client (Windows, Linux, OSX)
 - HTTPS covert channel for communications
 - Undetected by AVs (behavioral AVs might detect it if its not running on port 443)
 - Direct command execution (not using bash or sh)
 - Multiple sessions handling
 - File Upload/Download
 - System fingerprinting
 - Gracefully closing sessions server-side on client-side crash

## TODO:
 - Full Windows and OSX integration (currently it's partial)
 - bash/cmd/psh shell dropping
 - TOR routing?
 - Implement Twitter-Transfer-Protocol (<a href="https://github.com/elleven11/twitter-transfer-protocol" target="_blank">ttp</a>)

## Building:
**Requires Go 1.16 and up**  
To build the program you will need `openssl` and `go-bindata`.  
Use: `go get -u github.com/go-bindata/go-bindata/...`  

By default the client is set to listen on `127.0.0.1:1337`.  
To change that, you can edit the config object in to your liking `./main/client/main.go`  

When running `make` you will need to specify any external IP or domain to include in the SSL certificate.  
***This is done to prevent people stealing your binary and using it with malicious intent.***  
Example: `make IP=1.1.1.1 DOMAIN=example.com`.  
By default the Makefile sets `IP=127.0.0.1` and `DOMAIN=localhost`. If you want to keep that you can just ommit the variables in the make command.  
Example: `make`    
You will find your client and server builds in the `out` directory.  

Check Makefile for additional build/running options  
