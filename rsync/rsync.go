package rsync

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/rinetd/drone-rsync/utils"
)

type (

	// Config in Rsync
	Config struct {
		Hosts     []string
		Port      string
		User      string
		Key       string
		Password  string
		Source    string
		Target    string
		Args      []string
		Chmod     string
		Chown     string
		Verbose   string
		Recursive bool
		Delete    bool
		Sync      bool
		Include   []string
		Exclude   []string
		Filter    []string
		Export    []string
		Script    string
	}
	// Plugin -> Config
	Plugin struct {
		Config
	}
)

// Exec Plugin
func (p *Plugin) Exec() error {

	// write the rsa private key if provided
	utils.WriteKey(p.Config.Key)

	if len(p.User) == 0 {
		p.Config.User = "root"
	}
	// if len(p.Config.Hosts) == 0 {
	// 	p.Config.Hosts = []string{"d2", "d4"}
	// }

	// if len(p.Config.Source) == 0 {
	// 	p.Config.Source = "./"
	// }
	// if len(p.Config.Target) == 0 {
	// 	p.Config.Target = "/tmp/drone/drone"
	// }
	// p.Config.Source = "/home/ubuntu/test/test/"
	// p.Config.Sync = true
	// p.Config.Hosts = []string{"d2", "d4"}
	// p.Config.Port = "3009"
	// p.Config.User = "root"
	// p.Config.Target = "~/test/"
	// p.Config.Chmod = "0755"
	// p.Config.Chown = "33:33"
	// p.Config.Recursive = true
	// p.Config.Delete = true
	// log.Println("Verbose", p.Config.Verbose)

	// log.Println("[key:]", p.Config.Key)
	// log.Println("[Hosts:]", p.Config.Hosts)
	// log.Println("[Port:]", p.Config.Port)
	// log.Println("[User:]", p.Config.User)
	// log.Println("[Source:]", p.Config.Source)
	// log.Println("[Target:]", p.Config.Target)
	// log.Println("[Chmod:]", p.Config.Chmod)
	// log.Println("[Chown:]", p.Config.Chown)
	// log.Println("[Recursive:]", p.Config.Recursive)
	// log.Println("[Sync:]", p.Config.Sync)
	// log.Println("[Include:]", p.Config.Include)
	// log.Println("[Exclude:]", p.Config.Exclude)
	// log.Println("[Filter:]", p.Config.Filter)
	// log.Println("[Script:]", p.Config.Script)
	// p.genScript()

	return p.Run()
}

func (p *Plugin) Run() error {
	if p.Config.Sync {
		for i, host := range p.Config.Hosts {
			log.Println("===Sync host [", i, host, "]Runing===")
			p.commandRsync(host)
			p.commandSSH(host)
		}
	} else {
		//parallel
		wg := sync.WaitGroup{}
		wg.Add(len(p.Config.Hosts))
		errChannel := make(chan error, 2)

		finished := make(chan bool, 1)

		for i, host := range p.Config.Hosts {

			log.Println("===Async host [", i, host, "] Runing===")
			go func(host string, wg *sync.WaitGroup, errChannel chan error) {
				p.commandRsync(host)
				p.commandSSH(host)
				wg.Done()
			}(host, &wg, errChannel)
		}

		go func() {
			wg.Wait()
			close(finished)
		}()

		select {
		case <-finished:
			return nil
		case err := <-errChannel:
			if err != nil {
				log.Println(err)
				return err
			}
		}
	}
	return nil
}

// commandRsync rsync command
func (p *Plugin) commandRsync(host string) ([]byte, error) {
	if len(p.Config.Source) == 0 {
		log.Println("---Source is null , Skip Rsync---")
		return nil, nil
	}
	log.SetPrefix("[rsync]: ")

	if len(p.Config.Target) == 0 {
		log.Println("---Target is null , Skip Rsync---")
		return nil, nil
	}
	args := []string{
		"-az",
		// "-r", //--recursive          recurse into directories
		// "-l", //--links 				copy symlinks as symlinks
		// "-t", //--times              preserve modification times
		// "-p", //--perms 				preserve permissions
		// "-g", //--group              preserve group
		// "-o", //--owner              preserve owner (super-user only)
		// "-D", // same as --devices --specials
		// "-P", // same as --partial --progress
		"--partial", //                 keep partially transferred files
		// "-z", //--compress           compress file data during the transfer
		// "-h", //--human-readable     output numbers in a human-readable format
		// "-H", //--hard-links			preserve hard links
		// "-A", //--acls               preserve ACLs (implies --perms)
		// "-X", //--xattrs             preserve extended attributes
	}
	switch p.Config.Verbose {
	case "v", "-v":
		args = append(args, "-v")
	case "vv", "-vv":
		args = append(args, "-vv")
	case "vvv", "-vvv":
		args = append(args, "-vvv")
	default:
	}

	// append args rules
	for _, arg := range p.Config.Args {
		args = append(args, arg)
	}

	// append recursive flag
	// if p.Config.Recursive {
	// 	args = append(args, "-r")
	// }
	// append delete flag
	if p.Config.Delete {
		args = append(args, "-r")
		args = append(args, "--del")
	}
	if len(p.Config.Chown) > 0 {
		args = append(args, "--owner", "--group", "--chown", p.Config.Chown)
	}
	if len(p.Config.Chmod) > 0 {
		args = append(args, "--perms", "--chmod", p.Config.Chmod)
	}
	// append custom ssh parameters
	args = append(args, "-e", fmt.Sprintf("ssh -p %v -o ControlPersist=5m -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no", p.Config.Port))
	args = append(args, "--rsync-path", fmt.Sprintf("mkdir -p %s && rsync", path.Dir(p.Config.Target)))

	// append filtering rules
	for _, pattern := range p.Config.Include {
		args = append(args, fmt.Sprintf("--include=%s", pattern))
	}

	// args = append(args, "--exclude", ".git")
	for _, pattern := range p.Config.Exclude {
		args = append(args, fmt.Sprintf("--exclude=%s", pattern))
	}

	for _, pattern := range p.Config.Filter {
		args = append(args, fmt.Sprintf("--filter=%s", pattern))
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
	log.Println("[Rsync]", host, args)

	// fmt.Println("%v", cmd)
	// return exec.Command("rsync", args...)
	cmd.Stderr = os.Stderr
	b, err := cmd.Output()
	if err != nil {
		log.Println(err)
	}
	log.Println("[Rsync:", host, "output]\n", string(b))
	return b, err
}

func (p *Plugin) commandSSH(host string) ([]byte, error) {
	log.SetPrefix("[SSH]: ")

	if len(p.Config.Script) == 0 {
		log.Println("SSH Host:", host, "Skip Excute on SSH Remote with NULL Script !")
		return nil, nil

	}
	// log.Println("[SSH：]", host, " is Runing")

	args := []string{
		"-A",
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "LogLevel=quiet",
		"-o", "ControlPersist=5m",
	}
	if p.Config.Port != "22" {
		args = append(args, "-p", p.Config.Port)
	}
	if len(p.Config.User) != 0 {
		args = append(args, "-l", string(p.Config.User))
	}
	args = append(args, host)
	// args = append(args, "bash", path.Join(p.Config.Target, Filename))
	exp := p.genExport()
	sc := utils.Replace(p.Config.Script, "\\", ",", "&&")
	args = append(args, exp+sc)
	cmd := exec.Command("ssh", args...)

	log.Println("[SSH]", host, args)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	b, err := cmd.Output()
	if err != nil {
		log.Println("[SSH:", host, " ERROR output]\n", stderr.String())
		log.Println(err)
		return stderr.Bytes(), err
	}
	log.Println("[SSH:", host, " output]\n", string(b))
	return b, nil
}

// func (c *Cmd) Output() ([]byte, error) {

// 	var out bytes.Buffer

// 	cmd.Stdout = &out
// 	cmd.Stderr = &out
// 	err := cmd.Run()
// 	if err != nil {
// 		return out.Bytes(), err
// 	}
// 	return out.Bytes(), nil
// }

// var envdef = []string{"MYSQL_ROOT_PASSWORD"}

// var re_init = regexp.MustCompile(`[db_init (.*)]`)
var re_update = regexp.MustCompile(`\[db_update (.*?)\]`)
var re_delete = regexp.MustCompile(`\[db_delete (.*?)\]`)
var re_drop = regexp.MustCompile(`\[db_drop (.*?)\]`)
var re_backup = regexp.MustCompile(`\[db_backup (.*?)\]`)

func (p *Plugin) genExport() (ex string) {

	if ev := os.Getenv("DRONE_COMMIT_MESSAGE"); ev != "" {
		ex += "export DRONE_COMMIT_MESSAGE='" + strings.Replace(ev, "\n", " ", -1) + "';"
		st := strings.ToLower(ev)
		if strings.Contains(st, "[db_init") {
			ex += "export DB_INIT=true;"
			// 导出MYSQL
			if ev := os.Getenv("MYSQL_ROOT_PASSWORD"); ev != "" {
				ex += "export MYSQL_ROOT_PASSWORD='" + base64.StdEncoding.EncodeToString([]byte(ev)) + "';"
			}
		}
		if strings.Contains(st, "[db_update") {
			sql := re_update.FindStringSubmatch(st)[1]
			ex += "export DB_UPDATE=(" + sql + ");"
		}
		if strings.Contains(st, "[db_drop") {
			sql := re_drop.FindStringSubmatch(st)[1]
			ex += "export DB_DROP=(" + sql + ");"
		}
		if strings.Contains(st, "[db_delete") {
			sql := re_delete.FindStringSubmatch(st)[1]
			ex += "export DB_DELETE=(" + sql + ");"
		}
		if strings.Contains(st, "[db_backup") {
			sql := re_backup.FindStringSubmatch(st)[1]
			ex += "export DB_BACKUP=(" + sql + ");"
		}
	}
	// for _, v := range envdef {
	// 	ek := strings.ToUpper(v)
	// 	if ev := os.Getenv(ek); ev != "" {
	// 		er := strings.Replace(ev, "\n", " ", -1)
	// 		ex += "export " + ek + "='" + er + "';"
	// 		log.Println(ek, er)
	// 	}
	// }

	for _, v := range p.Config.Export {
		// log.Println(i, v)
		ek := strings.ToUpper(v)
		if ev := os.Getenv(ek); ev != "" {
			ex += "export " + ek + "=" + ev + ";"
		}
	}
	// log.Println("export：", ex)
	return
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

// var Filename string = ".drone-script.sh"

// func (p *Plugin) genScript() error {
// 	if len(p.Config.Script) == 0 {
// 		log.Println("p.Config.Script is Null")
// 		return nil

// 	}

// 	file := path.Join(p.Config.Source, Filename)
// 	// log.Println("[sourcePath:]", file)

// 	var buf = bytes.Buffer{}
// 	buf.WriteString("#! /bin/bash\n")
// 	s := utils.Split(p.Config.Script, "\\", ",")
// 	for _, v := range s {
// 		v += "\n"
// 		buf.WriteString(v)
// 	}

// 	// buf.WriteString(utils.Replace(p.Config.Script, "\\", ",", "\\n"))
// 	buf.WriteString("rm -- \"$0\"")
// 	ioutil.WriteFile(file, buf.Bytes(), 0755)

// 	return nil
// }
