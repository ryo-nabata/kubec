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

		// Build selectable items: the "unset" entry first, then each context.
		// The current context is flagged so it can be marked in the list.
		items := []contextItem{{Name: unsetContextLabel, IsUnset: true}}
		for _, context := range contexts {
			items = append(items, contextItem{
				Name:      context,
				IsCurrent: context == currentContext,
			})
		}

		// Create prompt template. The current context is annotated with
		// "(current)" so it stays identifiable even though the cursor
		// defaults to the "unset" entry at the top.
		templates := &promptui.SelectTemplates{
			Label:    "{{ .Name }}",
			Active:   "→ {{ .Name | cyan }}{{ if .IsCurrent }} (current){{ end }}",
			Inactive: "  {{ .Name }}{{ if .IsCurrent }} (current){{ end }}",
			Selected: "✓ {{ .Name | green }}",
		}

		prompt := promptui.Select{
			Label:     "Select a context",
			Items:     items,
			Templates: templates,
			CursorPos: 0,                  // Default to the "unset" entry (no context selected)
			Stdout:    utils.NoBellStdout, // Suppress the terminal bell on navigation
		}

		selectedIndex, _, err := prompt.Run()
		if err != nil {
			fmt.Printf("Selection cancelled: %v\n", err)
			return
		}
		selected := items[selectedIndex]

		// Unset the current context
		if selected.IsUnset {
			if err := utils.UnsetCurrentContext(); err != nil {
				log.Fatalf("Failed to unset context: %v", err)
			}
			fmt.Println("Current context unset (no context selected)")
			return
		}

		// Switch context
		err = utils.SetCurrentContext(selected.Name)
		if err != nil {
			log.Fatalf("Failed to switch context: %v", err)
		}

		fmt.Printf("Switched to context '%s'\n", color.GreenString(selected.Name))
	},
}

// unsetContextLabel is the menu entry shown in interactive mode that clears
// the current-context. The angle brackets keep it from colliding with a real
// Kubernetes context name.
const unsetContextLabel = "<未選択 / unset current-context>"

// contextItem is a single entry in the interactive context selector.
type contextItem struct {
	Name      string
	IsCurrent bool
	IsUnset   bool
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