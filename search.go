package goflags

import (
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
)

// searchUsageFunc is a usage function that searches for a flag
func (flagSet *FlagSet) searchUsageFunc() error {
	promptui.SearchPrompt = "Search for a flag: "
	prompt := promptui.Prompt{
		Label: "Search",
	}
	searchTerm, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("search prompt failed %v", err)
	}
	return flagSet.search(searchTerm)
}

// search searches for a string in the flagset
func (flagSet *FlagSet) search(term string) error {
	cliOutput := flagSet.CommandLine.Output()
	selectedFlags := flagSet.searchTerm(term)
	if len(selectedFlags) > 0 {
		prompt := promptui.Select{
			Label: "Select a flag",
			Items: selectedFlags,
		}

		_, result, err := prompt.Run()

		if err != nil {
			return fmt.Errorf("prompt failed %v", err)
		}

		fmt.Fprintf(cliOutput, "displaying help menu for: %q\n", result)
		os.Args[2] = result
		return nil
	}
	return fmt.Errorf("no match found for: %s", term)
}

// searchTerm returns a list of flags that contains the search term
func (flagSet *FlagSet) searchTerm(term string) []string {
	var res []string
	allKeys := flagSet.flagKeys.all()
	for _, flag := range allKeys {
		if strings.Contains(flag, term) {
			res = append(res, flag)
		}
	}
	return res
}
