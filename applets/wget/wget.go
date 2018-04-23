package wget

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

const binaryName = "wget"

type Option struct {
	helpFlag   bool
	spiderFlag bool
	showFlag   bool
	outputFile string
	timeout    int
	ua         string
}

func NewFlagSet() (*flag.FlagSet, *Option) {
	ret := flag.NewFlagSet(binaryName, flag.ExitOnError)

	ret.Usage = func() {
		fmt.Println("wget [-STOU] [url...]")
		ret.PrintDefaults()
	}

	var opt Option

	ret.BoolVar(&opt.helpFlag, "help", false, "show this message")
	ret.BoolVar(&opt.spiderFlag, "spider", false, "don't download anything")
	ret.BoolVar(&opt.showFlag, "S", false, "print server response")
	ret.IntVar(&opt.timeout, "T", 15, "set all timeout to seconds")
	ret.StringVar(&opt.outputFile, "O", "", "write documents to FILE")
	ret.StringVar(&opt.ua, "U", "Wget/0.0.1", " identify as AGENT instead of Wget/VERSION")
	ret.StringVar(&opt.ua, "-user-agent", "Wget/0.0.1", " identify as AGENT instead of Wget/VERSION")

	return ret, &opt
}

func Main(stdout io.Writer, args []string) error {
	flagSet, opt := NewFlagSet()
	flagSet.Parse(args)

	if flagSet.NArg() < 1 || opt.helpFlag {
		flagSet.Usage()
		return nil
	}

	return wget(stdout, flagSet.Args(), opt)
}

func wget(w io.Writer, urls []string, opt *Option) error {
	for _, url := range urls {
		if err := wgetOne(w, url, opt); err != nil {
			return err
		}
	}

	return nil
}

func wgetOne(w io.Writer, url string, opt *Option) error {
	timeout := time.Duration(opt.timeout) * time.Second

	output, err := getOutputFile(url, opt)
	if err != nil {
		return err
	}
	defer output.Close()

	client := http.Client{
		Timeout: timeout,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("User-Agent", opt.ua)
	r, err := client.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if opt.spiderFlag {
		return nil
	}
	if opt.showFlag {
		showServerResponse(w, r)
	}

	if _, err = io.Copy(output, r.Body); err != nil {
		return err
	}
	return nil
}

func getFilenameFromURL(urlStr string, opt *Option) (string, error) {
	url, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	fname := path.Base(url.Path)
	if len(fname) == 0 || fname == "." {
		fname = "index.html"
	}
	return fname, nil
}

func getOutputFile(urlStr string, opt *Option) (io.WriteCloser, error) {
	var filename string
	var err error

	if opt.outputFile == "-" {
		return os.Stdout, nil
	} else if opt.outputFile == "" {
		filename, err = getFilenameFromURL(urlStr, opt)
		if err != nil {
			return nil, err
		}
	} else {
		filename = opt.outputFile
	}
	return os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
}

func showServerResponse(w io.Writer, resp *http.Response) {
	fmt.Fprintf(w, "%s %s\n", resp.Proto, resp.StatusCode)

	for key, v := range resp.Header {
		fmt.Fprintf(w, "%s: %s\n", key, v)
	}

}
