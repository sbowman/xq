/*
Copyright © 2020 Sean Bowman
*/
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/antchfx/xmlquery"
	"github.com/antchfx/xpath"
	"github.com/spf13/cobra"
)

var find, exec string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "xq",
	Short: "Run XPath queries against an XML document from stdin",

	Run: func(cmd *cobra.Command, args []string) {
		if find == "" && exec == "" {
			_, _ = fmt.Fprintln(os.Stderr, "Please indicate an XPath expression")
			os.Exit(1)
		}

		if find != "" && exec != "" {
			_, _ = fmt.Fprintln(os.Stderr, "Please indicate either --find or --exec, but not both")
			os.Exit(1)
		}

		fi, _ := os.Stdin.Stat()
		if (fi.Mode() & os.ModeCharDevice) == 0 {
			doc, err := xmlquery.Parse(os.Stdin)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Unable to parse XML: %s\n", err)
				os.Exit(1)
			}

			if find != "" {
				matches, err := xmlquery.QueryAll(doc, find)
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "Invalid query: %s\n", err)
					os.Exit(1)
				}

				for _, match := range matches {
					switch match.Type {
					case xmlquery.TextNode,
						xmlquery.CharDataNode,
						xmlquery.CommentNode,
						xmlquery.AttributeNode:
						fmt.Println(match.InnerText())
					default:
						fmt.Println(match.OutputXML(true))
					}
				}
				return
			}

			expr, err := xpath.Compile(exec)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Invalid xpath expression: %s\n", err)
				os.Exit(1)
			}

			results := expr.Evaluate(xmlquery.CreateXPathNavigator(doc))
			switch r := results.(type) {
			case bool:
				fmt.Printf("%t\n", r)
			case string:
				fmt.Println(r)
			case float64:
				fmt.Println(strconv.FormatFloat(r, 'f', -1, 64))
			case xmlquery.NodeNavigator:
				for {
					fmt.Println(r.Current().OutputXML(true))
					if !r.MoveToNext() {
						break
					}
				}
			}
		}

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&find, "find", "f", "", "Search the XML for matching elements or attributes using XPath")
	rootCmd.Flags().StringVarP(&exec, "exec", "e", "", "Compile and run an XPath expression over the XML document")
}
