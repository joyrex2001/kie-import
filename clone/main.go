package clone

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	ssh2 "golang.org/x/crypto/ssh"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/client"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

// repositoryUrl will return the ssh url for the repo, using the viper
// interface to get the corresponding Drools git settings. It will return
// an error if the settings are invalid.
func repositoryUrl() (string, error) {
	host := viper.GetString("drools.host")
	if host == "" {
		return "", errors.New("drools host not provided")
	}
	port := viper.GetString("drools.git_ssh_port")
	if port == "" {
		return "", errors.New("drools git server port not provided")
	}
	repo := viper.GetString("drools.git_repo")
	if repo == "" {
		return "", errors.New("drools repository not provided")
	}
	return fmt.Sprintf("ssh://%s:%s/%s", host, port, repo), nil
}

// gitCredentials will return the ssh url for the repo, using the viper
// interface to get the corresponding Drools git settings. It will return
// an error if the settings are invalid.
func gitCredentials() (string, string, error) {
	user := viper.GetString("git.user")
	if user == "" {
		return "", "", errors.New("git username not provided")
	}
	passwd := viper.GetString("git.password")
	if user == "" {
		return "", "", errors.New("git password not provided")
	}
	return user, passwd, nil
}

// getDestination will return the destination folder for cloning this repo,
// using the viper interface to get the corresponding Drools git settings. It
// will return an error if the settings are invalid.
func getDestination() (string, error) {
	dest := viper.GetString("git.destination")
	if dest == "" {
		return "", errors.New("git destination not provided")
	}
	return dest, nil
}

// Main is the main entry point, based the settings initiated by cmd.
func Main(cmd *cobra.Command, args []string) {
	repo, err := repositoryUrl()
	CheckIfError(err)

	user, passwd, err := gitCredentials()
	CheckIfError(err)

	dest, err := getDestination()
	CheckIfError(err)

	// This is a specific client config that will allow connecting to the
	// drools internal git repo, without any user interaction. It will add
	// the (legacy) ssh-dss key algorithm, and it will ignore the known_hosts
	// file.
	ssh_config := &ssh2.ClientConfig{
		HostKeyAlgorithms: []string{
			ssh2.KeyAlgoRSA,
			ssh2.KeyAlgoDSA, // this is the extra required keyalgo
			ssh2.KeyAlgoECDSA256,
			ssh2.KeyAlgoECDSA384,
			ssh2.KeyAlgoECDSA521,
			ssh2.KeyAlgoED25519,
		},
		HostKeyCallback: ssh2.InsecureIgnoreHostKey(), // ftw, accept everything
		User:            user,
		Auth:            []ssh2.AuthMethod{ssh2.Password(passwd)},
	}
	client.InstallProtocol("ssh", ssh.NewClient(ssh_config))

	// This is the regular go-git auth object; note that it replicated the
	// user and auth method from the client config.
	auth := &ssh.Password{
		User:     user,
		Password: passwd,
	}

	// Check if we should do some monkey patching...
	if viper.GetBool("git.monkey-dsa2048") {
		Info("Monkey patches applied for ignoring dsa signature checking")
		PatchDsa2048()
	}

	fmt.Printf("Cloning %s to %s\n", repo, dest)

	r, err := git.PlainClone(dest, false, &git.CloneOptions{
		URL:  repo,
		Auth: auth,
	})
	CheckIfError(err)

	ref, err := r.Head()
	CheckIfError(err)

	commit, err := r.CommitObject(ref.Hash())
	CheckIfError(err)

	fmt.Printf("Imported %s\n", repo)
	fmt.Println(commit)
}
