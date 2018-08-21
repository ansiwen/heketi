package glusterd

import (
	"github.com/gluster/glusterd2/pkg/restclient"
	"github.com/heketi/heketi/executors"
	"github.com/heketi/heketi/executors/cmdexec"
	"github.com/heketi/heketi/pkg/utils"
)

type executor struct {
	client *restclient.Client
	config Config
	cmdexec.CmdExecutor
}

type Config struct {
	Schema     string `json:"url_schema"`
	ClientPort string `json:"client_port"`
	CertPath   string `json:"cert_path"`
	Insecure   bool   `json:"insecure"`
}

var (
	logger = utils.NewLogger("[glusterd]", utils.LEVEL_DEBUG)
)

func NewExecutor(config *Config) (executors.Executor, error) {
	g := executor{}
	//TODO add code read certfile and pass it
	g.config = *config
	return &g, nil
}

func (g *executor) createClient(host string) {
	//add default ip if not present
	url := g.config.Schema + "://" + host + ":" + g.config.ClientPort
	g.client = restclient.New(url, "", "", g.config.CertPath, g.config.Insecure)
}
