## jsubfinder

![jsubfinder logo](https://user-images.githubusercontent.com/17349277/146628329-af844760-2278-47b8-9ec1-022254774af7.png)


jsubfinder searches webpages for javascript & analyzes them for hidden subdomains and secrets (wip). From it's inception jsubfinder has been designed with performance in mind, to utilize large data sets & to be chained with other tools. It utilizes the [fasthttp go library](https://github.com/valyala/fasthttp) & golang's amazing concurency for blazing fast results.

This tool is still in active development thus hasn't been refactored and has alot of room for optomizations.

## Install

```
▶ go install github.com/hiddengearz/jsubfinder@latest
wget https://raw.githubusercontent.com/hiddengearz/jsubfinder/master/.jsf_signatures.yaml && mv .jsf_signatures.yaml ~/.jsf_signatures.yaml
```

## Basic Usage

jsubfinder accepts line-delimited domains on `stdin` & file input:

Examples (results are the same in this case):
```
▶ jsubfinder -u www.google.com
▶ jsubfinder -f google.txt
▶ echo www.google.com | jsubfinder
▶ echo www.google.com | jsubfinder -crawl
▶ echo www.google.com | httpx --silent | jsubfinder

apis.google.com
ogs.google.com
store.google.com
mail.google.com
accounts.google.com
www.google.com
policies.google.com
support.google.com
adservice.google.com
play.google.com
```

### With Secrets Enabled

```
▶ echo www.youtube.com | ./jsubfinder -s
www.youtube.com
youtubei.youtube.com
payments.youtube.com
2Fwww.youtube.com
252Fwww.youtube.com
m.youtube.com
tv.youtube.com
music.youtube.com
creatoracademy.youtube.com
artists.youtube.com
Google Cloud API Key AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8 found in content of https://www.youtube.com
Google Cloud API Key AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8 found in content of https://www.youtube.com
Google Cloud API Key AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8 found in content of https://www.youtube.com
Google Cloud API Key AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8 found in content of https://www.youtube.com
Google Cloud API Key AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8 found in content of https://www.youtube.com
Google Cloud API Key AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8 found in content of https://www.youtube.com
```

flag          | Description
------------- | -------------
-c            | Set the concurrency level (Default 10)
-f            | file with urls on each line
-u            | single url address to scan
-o            | file to output subdomains to. If secrets is enabled it's output file will be abreviated with secret_
-crawl        | Enable the basic crawler
-sig          | (optional) Location of signature file, by default is ~/.jsf_signatures.yaml
-d            | Enable debug mode
-g            | Enables greedy regex which scans all urls and not just JS files. This is disabled by default but will likely be enabled by default in the future
-s            | Enable secrets (beta), can result in alot of false positives.

## Credits

* The secrets (beta) funtion of this tool is heavily based off of [eth0izzle's](https://github.com/eth0izzle) [shhgit](https://github.com/eth0izzle/shhgit)
* jsubfinder is inspired by [nsonaniya2010's](https://github.com/nsonaniya2010) [SubDomainizer](https://github.com/nsonaniya2010/SubDomainizer)
