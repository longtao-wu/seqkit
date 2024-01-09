// Copyright © 2016-2019 Wei Shen <shenwei356@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/shenwei356/bio/seqio/fastx"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "seqkit",
	Short: "a cross-platform and ultrafast toolkit for FASTA/Q file manipulation",
	Long: fmt.Sprintf(`SeqKit -- a cross-platform and ultrafast toolkit for FASTA/Q file manipulation

Version: %s

Author: Wei Shen <shenwei356@gmail.com>

Documents  : http://bioinf.shenwei.me/seqkit
Source code: https://github.com/shenwei356/seqkit
Please cite: https://doi.org/10.1371/journal.pone.0163962


Seqkit utlizies the pgzip (https://github.com/klauspost/pgzip) package to
read and write gzip file, and the outputted gzip file would be slighty
larger than files generated by GNU gzip.

Seqkit writes gzip files very fast, much faster than the multi-threaded pigz,
therefore there's no need to pipe the result to gzip/pigz.

Seqkit also supports reading and writing xz (.xz) and zstd (.zst) formats since v2.2.0.
Bzip2 format is supported since v2.4.0.

Compression level:
  format   range   default  comment
  gzip     1-9     5        https://github.com/klauspost/pgzip sets 5 as the default value.
  xz       NA      NA       https://github.com/ulikunitz/xz does not support.
  zstd     1-4     2        roughly equals to zstd 1, 3, 7, 11, respectively.
  bzip     1-9     6        https://github.com/dsnet/compress

`, VERSION),
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.AddGroup(
		&cobra.Group{
			ID:    "basic",
			Title: "Commands for Basic Operation:",
		},
		&cobra.Group{
			ID:    "format",
			Title: "Commands for Format Conversion:",
		},
		&cobra.Group{
			ID:    "search",
			Title: "Commands for Searching:",
		},
		&cobra.Group{
			ID:    "set",
			Title: "Commands for Set Operation:",
		},
		&cobra.Group{
			ID:    "edit",
			Title: "Commands for Edit:",
		},
		&cobra.Group{
			ID:    "order",
			Title: "Commands for Ordering:",
		},
		&cobra.Group{
			ID:    "bam",
			Title: "Commands for BAM Processing:",
		},
		&cobra.Group{
			ID:    "misc",
			Title: "Commands for Miscellaneous:",
		},
	)

	defaultThreads := runtime.NumCPU()
	if defaultThreads > 4 {
		defaultThreads = 4
	}
	envThreads := os.Getenv("SEQKIT_THREADS")
	if envThreads != "" {
		t, err := strconv.Atoi(envThreads)
		if err == nil {
			defaultThreads = t
		}
	}
	if defaultThreads < 1 {
		defaultThreads = runtime.NumCPU()
	}
	RootCmd.PersistentFlags().StringP("seq-type", "t", "auto", "sequence type (dna|rna|protein|unlimit|auto) (for auto, it automatically detect by the first sequence)")
	RootCmd.PersistentFlags().IntP("threads", "j", defaultThreads, "number of CPUs. can also set with environment variable SEQKIT_THREADS)")
	RootCmd.PersistentFlags().IntP("line-width", "w", 60, "line width when outputting FASTA format (0 for no wrap)")
	RootCmd.PersistentFlags().StringP("id-regexp", "", fastx.DefaultIDRegexp, "regular expression for parsing ID")
	RootCmd.PersistentFlags().BoolP("id-ncbi", "", false, "FASTA head is NCBI-style, e.g. >gi|110645304|ref|NC_002516.2| Pseud...")
	RootCmd.PersistentFlags().StringP("out-file", "o", "-", `out file ("-" for stdout, suffix .gz for gzipped out)`)
	RootCmd.PersistentFlags().BoolP("quiet", "", false, "be quiet and do not show extra information")
	RootCmd.PersistentFlags().IntP("alphabet-guess-seq-length", "", 10000, "length of sequence prefix of the first FASTA record based on which seqkit guesses the sequence type (0 for whole seq)")
	RootCmd.PersistentFlags().StringP("infile-list", "X", "", "file of input files list (one file per line), if given, they are appended to files from cli arguments")
	RootCmd.PersistentFlags().IntP("compress-level", "", -1, `compression level for gzip, zstd, xz and bzip2. type "seqkit -h" for the range and default value for each format`)

	RootCmd.CompletionOptions.DisableDefaultCmd = true
	RootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	RootCmd.SetUsageTemplate(usageTemplate(""))

}

func usageTemplate(s string) string {
	return fmt.Sprintf(`Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}} %s{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}

Available Commands:{{range $cmds}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{else}}{{range $group := .Groups}}

{{.Title}}{{range $cmds}}{{if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if not .AllChildCommandsHaveGroup}}

Additional Commands:{{range $cmds}}{{if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsagesWrapped 110 | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsagesWrapped 110 | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`, s)
}
