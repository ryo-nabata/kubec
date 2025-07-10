package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/ryo-nabata/kubec/utils"
)

var showCurrent bool

var rootCmd = &cobra.Command{
	Use:   "kubec",
	Short: "A tool to easily switch Kubernetes contexts",
	Long:  `kubec is a command-line tool for easily switching Kubernetes current-context.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Show current context
		if showCurrent {
			currentContext := utils.GetCurrentContext()
			if currentContext != "" {
				fmt.Printf("Current context: %s\n", color.GreenString(currentContext))
			} else {
				fmt.Println("No current context is set")
			}
			return
		}

		// Direct context name specification
		if len(args) > 0 && shouldRunDirectContextSwitch(args[0]) {
			contextName := args[0]
			contexts := utils.GetContexts()
			
			// Check if the specified context exists
			found := false
			for _, c := range contexts {
				if c == contextName {
					found = true
					break
				}
			}
			
			if !found {
				fmt.Printf("Context '%s' not found\n", color.RedString(contextName))
				return
			}
			
			// Switch context
			err := utils.SetCurrentContext(contextName)
			if err != nil {
				log.Fatalf("Failed to switch context: %v", err)
			}
			
			fmt.Printf("Switched to context '%s'\n", color.GreenString(contextName))
			return
		}

		// Interactive mode
		contexts := utils.GetContexts()
		if len(contexts) == 0 {
			fmt.Println("No available contexts found")
			return
		}

		currentContext := utils.GetCurrentContext()
		
		// Create prompt template
		templates := &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   "→ {{ . | cyan }}",
			Inactive: "  {{ . | white }}",
			Selected: "✓ {{ . | green }}",
		}

		prompt := promptui.Select{
			Label:     "Select a context",
			Items:     contexts,
			Templates: templates,
		}

		// Set current context as initial selection
		if currentContext != "" {
			for i, context := range contexts {
				if context == currentContext {
					prompt.CursorPos = i
					break
				}
			}
		}

		_, selectedContext, err := prompt.Run()
		if err != nil {
			fmt.Printf("Selection cancelled: %v\n", err)
			return
		}

		// Switch context
		err = utils.SetCurrentContext(selectedContext)
		if err != nil {
			log.Fatalf("Failed to switch context: %v", err)
		}

		fmt.Printf("Switched to context '%s'\n", color.GreenString(selectedContext))
	},
}

func shouldRunDirectContextSwitch(arg string) bool {
	// Exclude special commands like help, version
	excludedArgs := []string{"help", "version", "--help", "-h", "--version", "-v"}
	
	for _, excluded := range excludedArgs {
		if strings.EqualFold(arg, excluded) {
			return false
		}
	}
	
	return true
}

func init() {
	rootCmd.Flags().BoolVarP(&showCurrent, "current", "c", false, "Show current context")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}