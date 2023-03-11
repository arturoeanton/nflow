package playbook

import (
	"fmt"
	"strings"

	"github.com/arturoeanton/nFlow/pkg/process"
	"github.com/google/shlex"
)

func command(cmd string) {
	/*if cmd[0] == '!' {
		cmd = cmd[1:]
		prg, err := shlex.Split(cmd)
		if err != nil {
			SendConsoleWS("\n" + err.Error())
			return
		}
		cmd := exec.Command(prg[0], prg[1:]...)
		stdout, err := cmd.Output()
		if err != nil {
			SendConsoleWS("\n" + err.Error())
			return
		}

		SendConsoleWS("\n" + string(stdout))
		return
	}

	*/
	terms, err := shlex.Split(cmd)
	if err != nil {
		SendConsoleWS("\n" + err.Error())
		return
	}

	if terms[0] == "nflow" || terms[0] == "d" {
		if terms[1] == "js" && len(terms) == 2 {
			var b strings.Builder
			for _, p := range Plugins {
				for key, fx := range p.AddFeatureJS() {
					fmt.Fprintf(&b, "\n%-30s \t\t %30s", key, fx)
				}
			}
			SendConsoleWS(b.String())
			return
		}

		if terms[1] == "vars" && len(terms) == 2 {
			var b strings.Builder
			for key, comment := range jsVars {
				fmt.Fprintf(&b, "\n%-30s \t\t %30s", key, comment)
			}

			SendConsoleWS(b.String())
			return
		}
		if (terms[1] == "ver" || terms[1] == "v") && len(terms) == 2 {
			msg := `nFlow v0.1`

			SendConsoleWS(msg)
			return
		}

		SendConsoleWS("error| dromedary param")
		return

	}

	if terms[0] == "help" {
		SendConsoleWS(`
ps
kill [wid]
killall
nflwo js || d js
nflow vars || d vars
nflow ver || d ver
help`)

	}
	if terms[0] == "ps" {
		SendConsoleWS(process.Ps())
		return
	}
	if terms[0] == "killall" {
		process.WKillAll()
		SendConsoleWS("killed all processes")
		return
	}

	if terms[0] == "kill" {
		if len(terms) == 2 {
			process.WKill(terms[1])
			SendConsoleWS("ok")
		}
		SendConsoleWS("error| kill [wid]")
		return
	}

}
