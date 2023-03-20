## ![jsubfinder logo](https://user-images.githubusercontent.com/17349277/146734055-8b836305-7a13-4c66-a02b-d92932322b42.png)




JSubFinder is a tool writtin in golang to search webpages & javascript for hidden subdomains and secrets in the given URL. Developed with BugBounty hunters in mind JSubFinder takes advantage of Go's amazing performance allowing it to utilize large data sets & be easily chained with other tools.


![z69D8q](https://user-images.githubusercontent.com/17349277/147615346-9c1471a6-a9a8-45cb-a429-f789b255950c.gif)

## Install
---
Install the application and download the signatures needed to find secrets

Using GO:

```bash
go install github.com/ThreatUnkown/jsubfinder@latest
wget https://raw.githubusercontent.com/ThreatUnkown/jsubfinder/master/.jsf_signatures.yaml && mv .jsf_signatures.yaml ~/.jsf_signatures.yaml
```

or

[Downloads Page](https://github.com/hiddengearz/jsubfinder/tags)


## Basic Usage
---

### Search

Search the given url's for subdomains and secrets

```text
$ jsubfinder search -h

Execute the command specified

Usage:
  JSubFinder search [flags]

Flags:
  -c, --crawl              Enable crawling
  -g, --greedy             Check all files for URL's not just Javascript
  -h, --help               help for search
  -f, --inputFile string   File containing domains
  -t, --threads int        Ammount of threads to be used (default 5)
  -u, --url strings        Url to check

Global Flags:
  -d, --debug               Enable debug mode. Logs are stored in log.info
  -K, --nossl               Skip SSL cert verification (default true)
  -o, --outputFile string   name/location to store the file
  -s, --secrets             Check results for secrets e.g api keys
      --sig string          Location of signatures for finding secrets
  -S, --silent              Disable printing to the console
```

Examples (results are the same in this case):

```bash
$ jsubfinder search -u www.google.com
$ jsubfinder search -f file.txt
$ echo www.google.com | jsubfinder search
$ echo www.google.com | httpx --silent | jsubfinder search$

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



#### With Secrets Enabled
*note `--secrets=""` will save the secret results in a secrets.txt file*
```bash

$ echo www.youtube.com | jsubfinder search --secrets=""
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

Google Cloud API Key <redacted> found in content of https://www.youtube.com
Google Cloud API Key <redacted> found in content of https://www.youtube.com
Google Cloud API Key <redacted> found in content of https://www.youtube.com
Google Cloud API Key <redacted> found in content of https://www.youtube.com
Google Cloud API Key <redacted> found in content of https://www.youtube.com
Google Cloud API Key <redacted> found in content of https://www.youtube.com
```


#### Advanced examples
```bash
$ echo www.google.com | jsubfinder search -crawl -s "google_secrets.txt" -S -o jsf_google.txt -t 10 -g
```

* `-crawl` use the default crawler to crawl pages for other URL's to analyze
* `-s` enables JSubFinder to search for secrets
* `-S` Silence output to console
* `-o <file>` save output to specified file
* `-t 10` use 10 threads
* `-g` search every URL for JS, even ones we don't think have any

### Proxy
Enables the upstream HTTP proxy with TLS MITM sypport. This allows you to:

1) Browse sites in realtime and have JSubFinder search for subdomains and secrets real time.
2) If needed run jsubfinder on another server to offload the workload

```text
$ JSubFinder proxy -h

Execute the command specified

Usage:
  JSubFinder proxy [flags]

Flags:
  -h, --help                    help for proxy
  -p, --port int                Port for the proxy to listen on (default 8444)
      --scope strings           Url's in scope seperated by commas. e.g www.google.com,www.netflix.com
  -u, --upstream-proxy string   Adress of upsteam proxy e.g http://127.0.0.1:8888 (default "http://127.0.0.1:8888")

Global Flags:
  -d, --debug               Enable debug mode. Logs are stored in log.info
  -K, --nossl               Skip SSL cert verification (default true)
  -o, --outputFile string   name/location to store the file
  -s, --secrets             Check results for secrets e.g api keys
      --sig string          Location of signatures for finding secrets
  -S, --silent              Disable printing to the console
```

```bash
$ jsubfinder proxy
Proxy started on :8444
Subdomain: out.reddit.com
Subdomain: www.reddit.com
Subdomain: 2Fwww.reddit.com
Subdomain: alb.reddit.com
Subdomain: about.reddit.com
```

#### With Burp Suite
1) Configure Burp Suite to forward traffic to an upstream proxy/ (User Options > Connections > Upsteam Proxy Servers > Add)
2) Run JSubFinder in proxy mode

Burp Suite will now forward all traffic proxied through it to JSubFinder. JSubFinder will retrieve the response, return it to burp and in another thread search for subdomains and secrets.

#### With Proxify
1) Launch [Proxify](https://github.com/projectdiscovery/proxify) & dump traffic to a folder `proxify -output logs`
2) Configure Burp Suite, a Browser or other tool to forward traffic to Proxify (see instructions on their [github page](https://github.com/projectdiscovery/proxify))
3) Launch JSubFinder in proxy mode & set the upstream proxy as Proxify `jsubfinder proxy -u http://127.0.0.1:8443`
4) Use Proxify's replay utility to replay the dumped traffic to jsubfinder `replay -output logs -burp-addr http://127.0.0.1:8444`


#### Run on another server
Simple, run JSubFinder in proxy mode on another server e.g 192.168.1.2. Follow the proxy steps above but set your applications upstream proxy as 192.168.1.2:8443

#### Advanced Examples

```bash
$ jsubfinder proxy --scope www.reddit.com -p 8081 -S -o jsf_reddit.txt
```

* `--scope` limits JSubFinder to only analyze responses from www.reddit.com
* `-p` port JSubFinders proxy server is running on
* `-S` silence output to the console/stdout
* `-o <file>` output examples to this file
