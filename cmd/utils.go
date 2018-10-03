// Copyright © 2018 Rodrigo Navarro <rodrigonavarro23@gmail.com>
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
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	. "github.com/logrusorgru/aurora"
	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
)

func checkErr(err error) {
	if err != nil {
		color.Set(color.FgMagenta)
		defer color.Unset()
		fmt.Print(err)
		os.Exit(0)
	}
}

func parseTemplate(template string) string {
	for _, v := range variables {
		template = strings.Replace(template, "{{"+v.Key+"}}", v.Value, -1)
	}

	return template
}

func promptList() {
	settings := viper.AllSettings()
	prompts := settings["prompt"].([]interface{})
	for _, v := range prompts {
		result := ""
		pr := v.(map[interface{}]interface{})
		if options := pr["OPTIONS"]; options == nil {
			validate := func(input string) error {
				if input == "" {
					return errors.New("Empty value")
				}
				return nil
			}
			p := promptui.Prompt{
				Label: pr["LABEL"],
				// Templates: templates,
				Validate: validate,
			}
			r, err := p.Run()
			checkErr(err)
			result = r
		} else {

			templates := &promptui.SelectTemplates{
				Label:    "{{ . | bold }}",
				Active:   "\U0001F449 {{ .Value | cyan }}  {{ .Desc | faint }}",
				Inactive: "   {{ .Value }}  {{ .Desc | faint }}",
				Selected: "\U0001F44D {{ .Value |  bold }}",
			}
			optList := pr["OPTIONS"].([]interface{})
			var opts []*option
			for _, o := range optList {
				op := o.(map[interface{}]interface{})
				opts = append(opts, &option{Value: op["VALUE"].(string), Desc: op["DESC"].(string)})
			}
			p := promptui.Select{
				Label:     pr["LABEL"],
				Items:     opts,
				Templates: templates,
			}
			i, _, err := p.Run()
			checkErr(err)
			result = opts[i].Value
		}

		variables = append(variables, keyValue{
			Key:   pr["KEY"].(string),
			Value: result,
		})
	}
}

func commit(message string) (err error) {
	cmdGit := exec.Command("git", "commit", "-m", message)
	lastCommit := "Last commit: " + message
	_, err = cmdGit.Output()
	checkErr(err)
	fmt.Println(Gray(lastCommit))
	fmt.Println(Green("Done"))

	return
}

func init() {
	cmdGit := exec.Command("git", "diff", "--cached", "--exit-code")
	_, err := cmdGit.Output()
	if err == nil {
		checkErr(errors.New("No changes added to commit"))
	}
}
