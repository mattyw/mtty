package mtty

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"
)

// Simple getting Started:
// run mtty
// run a command
// run :f to save the output of command to a file
// run :o to list options for opening the file at line
// vim should have been run using a command such as vim --servername vim
// TODO Unit Test for Open function to get it working.
var fileReg = regexp.MustCompile(`\S+:[0-9]+:?[0-9]?`)

type filelinecol struct {
	filename string
	line     string
	col      string
}

func split(flc string) filelinecol {
	s := strings.Split(flc, ":")
	if len(s) == 2 {
		return filelinecol{s[0], s[1], ""}
	}
	if len(s) == 3 {
		return filelinecol{s[0], s[1], s[2]}
	}
	if len(s) == 4 {
		return filelinecol{s[0], s[1], s[2]}
	}
	panic(string(flc))
}

type Mtty struct {
	lastOut []byte //Make buffer we can read and write to.

	Stdout io.Writer
	Stderr io.Writer
	Stdin  io.Reader
}

func (m *Mtty) Save(filename string) error {
	return ioutil.WriteFile(filename, m.lastOut, 0777)
}

func (m *Mtty) Open() {
	options := fileReg.FindAllString(string(m.lastOut), -1)
	if len(options) == 0 {
		return
	}
	opts := make([]filelinecol, len(options))
	for i, option := range options {
		fmt.Fprintf(m.Stdout, "%d) %s\n", i, string(option))
		opts[i] = split(string(option))
	}
	fmt.Println(opts)
	exec.Command(
		"vim",
		"--remote",
		fmt.Sprintf("+%s", opts[0].line),
		opts[0].filename,
	).Run()
}

func (m *Mtty) SetLastOut(b []byte) {
	m.lastOut = b
}

func (m *Mtty) runCommand(cmdS string) {
	sp := strings.Split(cmdS, " ")
	cmd := exec.Command(sp[0], sp[1:]...)
	cout, _ := cmd.CombinedOutput()
	fmt.Fprint(m.Stdout, string(cout))
	m.SetLastOut(cout)
}

func Loop(in io.Reader, out, errOut io.Writer) {
	tty := Mtty{
		Stdout: out,
		Stderr: errOut,
		Stdin:  in,
	}
	fmt.Fprintf(errOut, "$ ")
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		text := scanner.Text()
		switch text {
		case ":q":
			return
		case ":f":
			tty.Save("/tmp/mtty.txt") // TODO Accept an arg to take the name.
		case ":o":
			tty.Open()
		case ":h":
			fmt.Fprintln(out, "help not implemented yet") // TODO help
		default:
			tty.runCommand(text)
		}
		fmt.Fprintf(errOut, "$ ")
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(errOut, "reading standard input:", err)
	}
}
