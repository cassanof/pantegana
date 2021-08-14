# Pantegana - A Botnet RAT Made With Go
### <center>FOR EDUCATIONAL AND RESEARCH USE ONLY</center>  

    ▄▀▀▄▀▀▀▄  ▄▀▀█▄   ▄▀▀▄ ▀▄  ▄▀▀▀█▀▀▄  ▄▀▀█▄▄▄▄  ▄▀▀▀▀▄    ▄▀▀█▄   ▄▀▀▄ ▀▄  ▄▀▀█▄
    █   █   █ ▐ ▄▀ ▀▄ █  █ █ █ █    █  ▐ ▐  ▄▀   ▐ █         ▐ ▄▀ ▀▄ █  █ █ █ ▐ ▄▀ ▀▄
    ▐  █▀▀▀▀    █▄▄▄█ ▐  █  ▀█ ▐   █       █▄▄▄▄▄  █    ▀▄▄    █▄▄▄█ ▐  █  ▀█   █▄▄▄█
       █       ▄▀   █   █   █     █        █    ▌  █     █ █  ▄▀   █   █   █   ▄▀   █
     ▄▀       █   ▄▀  ▄▀   █    ▄▀        ▄▀▄▄▄▄   ▐▀▄▄▄▄▀ ▐ █   ▄▀  ▄▀   █   █   ▄▀
    █         ▐   ▐   █    ▐   █          █    ▐   ▐         ▐   ▐   █    ▐   ▐   ▐ 
    ▐                 ▐        ▐          ▐                          ▐

## Features:
 - Pretty and clean interactive shell (using <a href="https://github.com/desertbit/grumble" target="_blank">grumble</a>)
 - HTTPS covert channel for communications
 - Undetected by AVs (behavioral AVs might detect it if its not running on port 443)
 - Direct command execution (not using bash or sh)
 - Multiple sessions handling
 - File Upload/Download
 - System fingerprinting

## TODO:
 - Full Windows and OSx integration (currently it's partial)
 - Gracefully closing sessions server-side on client-side crash
 - bash/cmd/psh shell dropping
 - TOR routing?
 - Implement Twitter-Transfer-Protocol (<a href="https://github.com/elleven11/twitter-transfer-protocol" target="_blank">ttp</a>)

## Building:
To build the program you will need `openssl` and `go-bindata`.  
Use: `go get -u github.com/go-bindata/go-bindata/...`  

When running `make` you will need to specify any external IP or domain to include in the SSL certificate.  
***This is done to prevent people stealing your binary and using it with malicious intent.***  
Example: `make IP=1.1.1.1 DOMAIN=example.com`.  
If you do not want to specify an IP or a domain, use `127.0.0.1` and `localhost` respectively.  
Example: `make IP=127.0.0.1 DOMAIN=localhost`    

Check Makefile for different build/running options
