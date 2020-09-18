## jsubfinder

jsubfinder searches webpages for javascript & analyzes them for hidden subdomains and secrets (wip). From it's inception jsubfinder has been designed with performance in mind, to utilize large data sets & to be chained with other tools. It utilizes the [fasthttp go library](https://github.com/valyala/fasthttp) & golang's amazing concurency for blazing fast results.

This tool is still in active development thus hasn't been refactored and has alot of room for optomizations.

## Install

```
▶ go get -u github.com/hiddengearz/jsubfinder
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
▶ echo www.youtube.com | ./jsubfinder -s -o youtube.com
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

▶ cat secrets_youtube.com
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
-d            | Enable debug mode
-g            | Enables greedy regex which scans all urls and not just JS files. This is disabled by default but will likely be enabled by default in the future
-s            | Enable secrets (beta), can result in alot of false positives.

## Credits

* The secrets (beta) funtion of this tool is heavily based off of [eth0izzle's](https://github.com/eth0izzle) [shhgit](https://github.com/eth0izzle/shhgit)
* jsubfinder is inspired by [nsonaniya2010's](https://github.com/nsonaniya2010) [SubDomainizer](https://github.com/nsonaniya2010/SubDomainizer)
