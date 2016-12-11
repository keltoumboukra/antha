// list.go: Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
//
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

package cmd

import (
	//"encoding/json"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"strings"
)

var reportTemplateCmd = &cobra.Command{
	Use:   "reporttemplate",
	Short: "produce a report template",
	RunE:  readme,
}

func readme(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	switch viper.GetString("output") {
	case jsonOutput:
		_, err := fmt.Println("Json not valid for report templates")
		return err

	default:
		var Urlstring string = "https://raw.githubusercontent.com/antha-lang/antha/AnthaAcademyVM2/antha/AnthaStandardLibrary/Packages/Templates/reporttemplate.md"
		file := "report" + fmt.Sprint(time.Now().Format("20060102150405")) + ".md"
		var err error

		anthacommit, err := GitCommit()
		if err != nil {
			return err
		}
		fmt.Println(anthacommit)
		if _, err = os.Stat(file); os.IsNotExist(err) {
			if err := os.MkdirAll(filepath.Dir(file), 0777); err != nil {
				return err
			}

			res, err := http.Get(Urlstring)
			if err != nil {
				return err
			}
			defer res.Body.Close()

			f, err := os.Create(file)
			if err != nil {
				return err
			}
			defer f.Close()

			var buf bytes.Buffer
			if _, err := io.Copy(&buf, res.Body); err != nil {
				return err
			}
			readme := string(buf.Bytes())

			newreadme := strings.Replace(readme, "***ANTHACOMMIT****", anthacommit, 1)
			if err := ioutil.WriteFile(file, []byte(newreadme), 0666); err != nil {
				return err
			}
			return nil
		}

		return err
	}
}

func init() {
	c := reportTemplateCmd
	flags := c.Flags()
	RootCmd.AddCommand(c)

	flags.String(
		"output",
		textOutput,
		fmt.Sprintf("Output format: one of {%s}", strings.Join([]string{textOutput, jsonOutput}, ",")))
}
