package rsync

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path"
	"strconv"
	"sync"

	"github.com/rinetd/drone-rsync/utils"
)

type (
	// Config Rsync
	// Config struct {
	// 	Hosts     []string `json:"hosts"`
	// 	Port      int      `json:"port"`
	// 	User      string   `json:"user"`
	// 	Key       string   `json:"ssh-key"`
	// 	Password  string   `json:"password"`
	// 	Source    string   `json:"source"`
	// 	Target    string   `json:"target"`
	// 	Chmod     string   `json:"chmod"`
	// 	Chown     string   `json:"chown"`
	// 	Recursive bool     `json:"recursive"`
	// 	Delete    bool     `json:"delete"`
	// 	Include   []string `json:"include"`
	// 	Exclude   []string `json:"exclude"`
	// 	Filter    []string `json:"filter"`
	// 	Script    []string `json:"script"`
	// }

	// Config in Rsync
	Config struct {
		Hosts     []string
		Port      int
		User      string
		Key       string
		Password  string
		Source    string
		Target    string
		Chmod     string
		Chown     string
		Recursive bool
		Delete    bool
		Sync      bool
		Include   []string
		Exclude   []string
		Filter    []string
		Script    string
	}
	// Plugin -> Config
	Plugin struct {
		Config
	}
)

// Exec Plugin
func (p *Plugin) Exec() error {
	log.SetPrefix("[rsync]: ")
	// upstreamName := p.Config.Hosts

	// 1.rsync 同步

	// 2.执行ssh

	// write the rsa private key if provided
	if err := utils.WriteKey(p.Config.Key); err != nil {
		log.Panicln("write key error")
		return err
	}

	if len(p.Config.Hosts) == 0 {
		p.Config.Hosts = []string{"localhost", "127.0.0.1"}
	}

	// default values
	if p.Config.Port == 0 {
		p.Config.Port = 22
	}

	if len(p.Config.User) == 0 {
		p.Config.User = "root"
	}
	if len(p.Config.Source) == 0 {
		p.Config.Source = "./"
	}
	if len(p.Config.Target) == 0 {
		p.Config.Target = "/tmp/drone/drone"
	}

	// p.Config.Sync = true
	// p.Config.Hosts = []string{"d2", "d4"}
	// p.Config.Port = 3009
	// p.Config.User = "root"
	// p.Config.Target = "/tmp/drone/drone"
	// p.Config.Recursive = true

	// log.Println("[key:]", p.Config.Key)
	log.Println("[Hosts:]", p.Config.Hosts)
	log.Println("[Port:]", p.Config.Port)
	log.Println("[User:]", p.Config.User)
	log.Println("[Source:]", p.Config.Source)
	log.Println("[Target:]", p.Config.Target)
	log.Println("[Chmod:]", p.Config.Chmod)
	log.Println("[Chown:]", p.Config.Chown)
	log.Println("[Recursive:]", p.Config.Recursive)
	log.Println("[Sync:]", p.Config.Sync)
	log.Println("[Include:]", p.Config.Include)
	log.Println("[Exclude:]", p.Config.Exclude)
	log.Println("[Filter:]", p.Config.Filter)
	log.Println("[Script:]", p.Config.Script)
	p.genScript()
	p.asyncRun()
	// execute for each host
	// for _, host := range p.Config.Hosts {
	// sync the files on the remote machine
	// rs := p.commandRsync(host)
	// trace(rs)
	// err := p.Config.Run()
	// if err != nil {
	// 	return err
	// }

	// // continue if no commands
	// if len(v.Commands) == 0 {
	// 	continue
	// }

	// // execute commands on remote server (reboot instance, etc)
	// if err := v.run(w.Keys, host); err != nil {
	// 	return err
	// }
	// }

	return nil
}

func (p *Plugin) asyncRun() error {

	wg := sync.WaitGroup{}
	wg.Add(len(p.Config.Hosts))
	errChannel := make(chan error, 2)

	finished := make(chan bool, 1)
	go func() {
		wg.Wait()
		close(finished)
	}()
	for _, host := range p.Config.Hosts {

		if p.Config.Sync {
			log.Println("===Sync Runing===\n", host)
			p.commandRsync(host)
			p.commandSSH(host)
			wg.Done()
			// p.commandSSH(host).Run()
		} else {
			log.Println("===Async Runing===\n", host)
			go func(host string, wg *sync.WaitGroup, errChannel chan error) {
				p.commandRsync(host)
				p.commandSSH(host)

				wg.Done()
			}(host, &wg, errChannel)
		}
	}

	select {
	case <-finished:
	case err := <-errChannel:
		if err != nil {
			log.Println(err)
			return err
		}
	}
	return nil
}

var Filename string = ".drone-script.sh"

func (p *Plugin) genScript() error {
	// if len(p.Config.Script) > 0 {
	// 	log.Println(p.Config.Script)

	// }

	file := path.Join(p.Config.Source, Filename)
	// log.Println("[sourcePath:]", file)

	var buf = bytes.Buffer{}
	buf.WriteString("#! /bin/bash\n")
	s := utils.Split(p.Config.Script, "\\", ",")
	for _, v := range s {
		v += "\n"
		buf.WriteString(v)
	}

	// buf.WriteString(utils.Replace(p.Config.Script, "\\", ",", "\\n"))
	// buf.WriteString("rm -- \"$0\"")
	ioutil.WriteFile(file, buf.Bytes(), 0755)

	return nil
}

// buildRsync rsync command
func (p *Plugin) commandRsync(host string) ([]byte, error) {
	log.Println("[Rsync：] Runing")
	// p.genScript()
	if len(p.Config.Target) == 0 {
		return nil, nil
	}
	args := []string{
		"-az",
	}
	// append recursive flag
	if p.Config.Recursive {
		args = append(args, "-r")
		// append delete flag
		if p.Config.Delete {
			args = append(args, "--del")
		}
	}
	// append custom ssh parameters
	args = append(args, "-e", fmt.Sprintf("ssh -p %d -o UserKnownHostsFile=/dev/null -o LogLevel=quiet -o StrictHostKeyChecking=no", p.Config.Port))
	args = append(args, "--rsync-path", fmt.Sprintf("mkdir -p %s && rsync", p.Config.Target))

	if len(p.Config.Chown) > 0 {
		args = append(args, fmt.Sprintf("--owner --group --chown=%s", p.Config.Chown))
	}
	if len(p.Config.Chmod) > 0 {
		args = append(args, fmt.Sprintf("--perms --chmod=%s", p.Config.Chmod))
	}

	// append filtering rules
	for _, pattern := range p.Config.Include {
		args = append(args, fmt.Sprintf("--include=%s", pattern))
	}

	args = append(args, "--exclude", ".git")
	for _, pattern := range p.Config.Exclude {
		args = append(args, fmt.Sprintf("--exclude=%s", pattern))
	}

	for _, pattern := range p.Config.Filter {
		args = append(args, fmt.Sprintf("--filter=%s", pattern))
	}
	//
	if len(p.Config.Script) > 0 {
		log.Println(p.Config.Script)

	}

	// args = append(args, p.globSource(root)...)
	if len(p.Config.Source) > 0 {
		args = append(args, p.Config.Source)
	} else {
		args = append(args, ".")
	}
	args = append(args, fmt.Sprintf("%s@%s:%s", p.Config.User, host, p.Config.Target))

	cmd := exec.Command("rsync", args...)
	// cmd.Dir = workspacePath
	// if p.Config.Sync {
	// 	cmd.Stdout = os.Stdout
	// 	cmd.Stderr = os.Stderr
	// } else {
	// 	cmd.Stdout = ioutil.Discard
	// 	cmd.Stderr = ioutil.Discard
	// }
	log.Println("[Rsync:]", args)

	// fmt.Println("%v", cmd)
	// return exec.Command("rsync", args...)
	b, err := cmd.Output()
	log.Println("[Rsync:]", string(b))
	return b, err
}

func (p *Plugin) commandSSH(host string) ([]byte, error) {
	log.Println("[SSH：] Runing")

	args := []string{
		"-A",
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "LogLevel=quiet",
	}
	if p.Config.Port != 22 {
		args = append(args, "-p", strconv.Itoa(p.Config.Port))
	}
	if len(p.Config.User) != 0 {
		args = append(args, "-l", string(p.Config.User))
	}
	args = append(args, host)
	args = append(args, "bash -c", path.Join(p.Config.Target, Filename))
	cmd := exec.Command("ssh", args...)
	log.Println(args)
	b, err := cmd.Output()
	log.Println("[SSH:]", string(b))
	return b, err
}

// func (p *Plugin) globSource(root string) []string {
// 	src := p.Config.Source
// 	if !path.IsAbs(p.Config.Source) {
// 		src = path.Join(root, p.Config.Source)
// 	}
// 	srcs, err := filepath.Glob(src)
// 	if err != nil || len(srcs) == 0 {
// 		return []string{p.Config.Source}
// 	}
// 	sep := fmt.Sprintf("%c", os.PathSeparator)
// 	if strings.HasSuffix(p.Config.Source, sep) {
// 		// Add back the trailing slash removed by path.Join()
// 		for i := range srcs {
// 			srcs[i] += sep
// 		}
// 	}
// 	return srcs
// }

// Run commands on the remote host

// func (p Plugin) run(keys *drone.Key, host string) error {

// 	// join the host and port if necessary
// 	addr := net.JoinHostPort(host, strconv.Itoa(p.Config.Port))

// 	// trace command used for debugging in the build logs
// 	fmt.Printf("$ ssh %s@%s -p %d\n", p.Config.User, addr, p.Config.Port)

// 	signer, err := ssh.ParsePrivateKey([]byte(keys.Private))
// 	if err != nil {
// 		return fmt.Errorf("Error parsing private key. %s.", err)
// 	}

// 	config := &ssh.ClientConfig{
// 		User: p.Config.User,
// 		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
// 	}

// 	client, err := ssh.Dial("tcp", addr, config)
// 	if err != nil {
// 		return fmt.Errorf("Error dialing server. %s.", err)
// 	}

// 	session, err := client.NewSession()
// 	if err != nil {
// 		return fmt.Errorf("Error starting ssh session. %s.", err)
// 	}
// 	defer session.Close()

// 	session.Stdout = os.Stdout
// 	session.Stderr = os.Stderr
// 	return session.Run(strings.Join(p.Config.Commands, "\n"))
// }

// globSource returns the names of all files matching the source pattern.
// If there are no matches or an error occurs, the original source string is
// returned.
//
// If the source path is not absolute the root path will be prepended to the
// source path prior to matching.

// Trace writes each command to standard error (preceded by a ‘$ ’) before it
// is executed. Used for debugging your build.
// func trace(cmd *exec.Cmd) {
// 	fmt.Println("$", strings.Join(cmd.Args, " "))
// }
