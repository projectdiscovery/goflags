package main

import (
	"fmt"

	"github.com/projectdiscovery/goflags"
)

// Options contains the configuration options for nuclei scanner.
type Options struct {
	// RandomAgent generates random User-Agent
	RandomAgent bool
	// Metrics enables display of metrics via an http endpoint
	Metrics bool
	// Debug mode allows debugging request/responses for the engine
	Debug bool
	// DebugRequests mode allows debugging request for the engine
	DebugRequests bool
	// DebugResponse mode allows debugging response for the engine
	DebugResponse bool
	// Silent suppresses any extra text and only writes found URLs on screen.
	Silent bool
	// Version specifies if we should just show version and exit
	Version bool
	// Verbose flag indicates whether to show verbose output or not
	Verbose bool
	// No-Color disables the colored output.
	NoColor bool
	// UpdateTemplates updates the templates installed at startup
	UpdateTemplates bool
	// JSON writes json output to files
	JSON bool
	// JSONRequests writes requests/responses for matches in JSON output
	JSONRequests bool
	// EnableProgressBar enables progress bar
	EnableProgressBar bool
	// TemplatesVersion shows the templates installed version
	TemplatesVersion bool
	// TemplateList lists available templates
	TemplateList bool
	// Stdin specifies whether stdin input was given to the process
	Stdin bool
	// StopAtFirstMatch stops processing template at first full match (this may break chained requests)
	StopAtFirstMatch bool
	// NoMeta disables display of metadata for the matches
	NoMeta bool
	// Project is used to avoid sending same HTTP request multiple times
	Project bool
	// MetricsPort is the port to show metrics on
	MetricsPort int
	// BulkSize is the of targets analyzed in parallel for each template
	BulkSize int
	// TemplateThreads is the number of templates executed in parallel
	TemplateThreads int
	// Timeout is the seconds to wait for a response from the server.
	Timeout int
	// Retries is the number of times to retry the request
	Retries int
	// Rate-Limit is the maximum number of requests per specified target
	RateLimit int
	// BurpCollaboratorBiid is the Burp Collaborator BIID for polling interactions.
	BurpCollaboratorBiid string
	// ProjectPath allows nuclei to use a user defined project folder
	ProjectPath string
	// Severity filters templates based on their severity and only run the matching ones.
	Severity goflags.StringSlice
	// Target is a single URL/Domain to scan using a template
	Target string
	// Targets specifies the targets to scan using templates.
	Targets string
	// Output is the file to write found results to.
	Output string
	// ProxyURL is the URL for the proxy server
	ProxyURL string
	// ProxySocksURL is the URL for the proxy socks server
	ProxySocksURL string
	// TemplatesDirectory is the directory to use for storing templates
	TemplatesDirectory string
	// TraceLogFile specifies a file to write with the trace of all requests
	TraceLogFile string
	// Templates specifies the template/templates to use
	Templates goflags.StringSlice
	// 	ExcludedTemplates  specifies the template/templates to exclude
	ExcludedTemplates goflags.StringSlice
	// CustomHeaders is the list of custom global headers to send with each request.
	CustomHeaders goflags.StringSlice
	// Normalized contains the list of normalized input formats for nuclei
	Normalized string
	// NormalizedOutput writes the internal normalized format representation to a file.
	NormalizedOutput string
}

func main() {
	set := goflags.New()
	set.SetDescription(`Nuclei is a fast tool for configurable targeted scanning
based on templates offering massive extensibility and ease of use.`)

	options := &Options{}
	set.BoolVar(&options.Metrics, "metrics", false, "Expose nuclei metrics on a port")
	set.IntVar(&options.MetricsPort, "metrics-port", 9092, "Port to expose nuclei metrics on")
	set.StringVar(&options.Target, "target", "", "Target is a single target to scan using template")
	set.StringSliceVarP(&options.Templates, "templates", "t", []string{}, "Template input dir/file/files to run on host. Can be used multiple times. Supports globbing.")
	set.StringSliceVar(&options.ExcludedTemplates, "exclude", []string{}, "Template input dir/file/files to exclude. Can be used multiple times. Supports globbing.")
	set.StringVarP(&options.Normalized, "normalized", "n", "", "Normalized requests input dir/file/files.")
	set.StringVar(&options.NormalizedOutput, "normalized-output", "", "Optional File to write internal normalized format representation to")
	set.StringSliceVar(&options.Severity, "severity", []string{}, "Filter templates based on their severity and only run the matching ones. Comma-separated values can be used to specify multiple severities.")
	set.StringVarP(&options.Targets, "list", "l", "", "List of URLs to run templates on")
	set.StringVarP(&options.Output, "output", "o", "", "File to write output to (optional)")
	set.StringVar(&options.ProxyURL, "proxy-url", "", "URL of the proxy server")
	set.StringVar(&options.ProxySocksURL, "proxy-socks-url", "", "URL of the proxy socks server")
	set.BoolVar(&options.Silent, "silent", false, "Show only results in output")
	set.BoolVar(&options.Version, "version", false, "Show version of nuclei")
	set.BoolVarP(&options.Verbose, "verbose", "v", false, "Show Verbose output")
	set.BoolVar(&options.NoColor, "no-color", false, "Disable colors in output")
	set.IntVar(&options.Timeout, "timeout", 5, "Time to wait in seconds before timeout")
	set.IntVar(&options.Retries, "retries", 1, "Number of times to retry a failed request")
	set.BoolVar(&options.RandomAgent, "random-agent", false, "Use randomly selected HTTP User-Agent header value")
	set.StringSliceVarP(&options.CustomHeaders, "header", "H", []string{}, "Custom Header.")
	set.BoolVar(&options.Debug, "debug", false, "Allow debugging of request/responses")
	set.BoolVar(&options.DebugRequests, "debug-req", false, "Allow debugging of request")
	set.BoolVar(&options.DebugResponse, "debug-resp", false, "Allow debugging of response")
	set.BoolVar(&options.UpdateTemplates, "update-templates", false, "Update Templates updates the installed templates (optional)")
	set.StringVar(&options.TraceLogFile, "trace-log", "", "File to write sent requests trace log")
	set.StringVar(&options.TemplatesDirectory, "update-directory", "", "Directory to use for storing nuclei-templates")
	set.BoolVar(&options.JSON, "json", false, "Write json output to files")
	set.BoolVar(&options.JSONRequests, "include-rr", false, "Write requests/responses for matches in JSON output")
	set.BoolVar(&options.EnableProgressBar, "stats", false, "Display stats of the running scan")
	set.BoolVar(&options.TemplateList, "tl", false, "List available templates")
	set.IntVar(&options.RateLimit, "rate-limit", 150, "Rate-Limit (maximum requests/second")
	set.BoolVar(&options.StopAtFirstMatch, "stop-at-first-match", false, "Stop processing http requests at first match (this may break template/workflow logic)")
	set.IntVar(&options.BulkSize, "bulk-size", 25, "Maximum Number of hosts analyzed in parallel per template")
	set.IntVarP(&options.TemplateThreads, "concurrency", "c", 10, "Maximum Number of templates executed in parallel")
	set.BoolVar(&options.Project, "project", false, "Use a project folder to avoid sending same request multiple times")
	set.StringVar(&options.ProjectPath, "project-path", "", "Use a user defined project folder, temporary folder is used if not specified but enabled")
	set.BoolVar(&options.NoMeta, "no-meta", false, "Don't display metadata for the matches")
	set.BoolVar(&options.TemplatesVersion, "templates-version", false, "Shows the installed nuclei-templates version")
	set.StringVar(&options.BurpCollaboratorBiid, "burp-collaborator-biid", "", "Burp Collaborator BIID")

	set.SetConfigFile("config.yaml")
	set.Parse()
	fmt.Printf("%+v\n", options)
}
