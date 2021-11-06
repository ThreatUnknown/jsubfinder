package main

import "github.com/hiddengearz/jsubfinder/cmd"

//"strconv"

func main() {
	cmd.Execute()
}

/*
func main() {
	if core.Debug {
		defer core.TimeTrack(time.Now(), "JSubfinder")
	}

	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	var urls []string
	var sig string

	concurrencyFlag := flag.Int("c", 10, "Concurrency")
	fileFlag := flag.String("f", "", "File with urls")
	flag.BoolVar(&core.Greedy, "g", false, "Use Greedy regex to find subdomains")
	urlFlag := flag.String("u", "", "Url address to scan")
	flag.BoolVar(&core.Debug, "d", false, "Enable Debug mode")
	outputFlag := flag.String("o", "", "Output results to this file")
	flag.BoolVar(&core.Crawl, "crawl", false, "Crawl")
	flag.BoolVar(&core.FindSecrets, "s", false, "Search for secrets")
	flag.StringVar(&sig, "sig", home+"/.jsf_signatures.yaml", "Location of signatures")

	flag.Parse()

	if core.FindSecrets {
		core.ConfigSigs.ParseConfig(sig) //https://github.com/eth0izzle/shhgit/blob/090e3586ee089f013659e02be16fd0472b629bc7/core/signatures.go
		core.Signatures = core.ConfigSigs.GetSignatures()
		core.Blacklisted_extensions = []string{".exe", ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".tif", ".psd", ".xcf", ".zip", ".tar.gz", ".ttf", ".lock"}
		if *outputFlag == "" {
			core.PrintSecrets = true
		}
	}
	//fmt.Println("compiled signatures: " + strconv.Itoa(len(core.Signatures)))
	//return

	if core.IsFlagPassed("f") && core.IsFlagPassed("u") {
		fmt.Println("Provide either -f or -u, you can't provide both")
	} else if core.IsFlagPassed("f") {
		urls = core.ReadFile(*fileFlag)
	} else if core.IsFlagPassed("u") {
		urls = append(urls, *urlFlag)
		//} else if Debug && fileExists("test.txt") {
		//fmt.Println("using test.txt")
		//urls = ReadFile("test.txt")
	} else {
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			urls = append(urls, strings.ToLower(sc.Text()))
		}
	}

	core.Exec(urls, *concurrencyFlag, *outputFlag)

}
*/
