package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// ReqBody from reqeust
type ReqBody struct {
	Repo    string `json:"repo"`
	Env     map[string]string
	Timeout int    `json:"timeout" default:"600"`
	Mode    string `json:"mode" default:"plan"`
	Varfile string `json:"varfile"`
	Extra   string `json:"extra"`
}

// ReqToCommand create command structure to run container
// from POST request
func ReqToCommand(req *http.Request) (*Command, error) {
	var d ReqBody
	if err := json.NewDecoder(req.Body).Decode(&d); err != nil {
		req.Body.Close()
		return nil, err
	}

	c := new(Command)
	c.Image = TerraformImage

	//make env string
	var env []string
	for k, v := range d.Env {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	c.Env = env

	// make command
	var cmdlist []string
	cmdlist = append(cmdlist, fmt.Sprintf("git clone %s&&", d.Repo))

	// get folder name
	s := strings.Split(d.Repo, "/")
	f := s[len(s)-1]
	f = f[:len(f)-4] // remove .git

	cmdlist = append(cmdlist, fmt.Sprintf("cd %s&&", f))
	cmdlist = append(cmdlist, "terraform init&&")
	if d.Mode == "apply" {
		log.Println("entering apply mode ...")
		cmdlist = append(cmdlist, fmt.Sprintf("%s %s&&%s", "terraform apply -auto-approve -var-file", d.Varfile, d.Extra))
	} else if d.Mode == "destroy" {
		log.Println("entering destroy mode ...")
		cmdlist = append(cmdlist, fmt.Sprintf("%s %s&&%s", "terraform destroy -auto-approve -var-file", d.Varfile, d.Extra))
	} else if d.Mode == "pull" {
		log.Println("show state info ...")
		cmdlist = append(cmdlist, fmt.Sprintf("%s&&%s", "terraform state pull", d.Extra))
	} else {
		cmdlist = append(cmdlist, fmt.Sprintf("%s %s&&%s", "terraform plan -var-file", d.Varfile, d.Extra))
	}

	cmdstr := ""
	for _, v := range cmdlist {
		cmdstr = cmdstr + v
	}
	var t []string
	t = append(t, "sh")
	t = append(t, "-c")
	t = append(t, cmdstr)
	c.Commands = t

	// set timeout
	c.Timeout = d.Timeout

	c.ContainerInstance = new(Container)
	c.ContainerInstance.Context = context.Background()
	log.Printf("%#v", c)
	return c, nil
}
