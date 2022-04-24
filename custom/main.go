package main

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic("error create client: " + err.Error())
	}

	/*reader, err := cli.ImagePull(ctx, "testapp1:latest", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)*/

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "testapp1:latest",
		//Cmd:   []string{"cat", "/etc/hosts"},
		Tty: true,
	}, nil, nil, &specs.Platform{
		OS:           "linux",
		Architecture: "arm64",
	}, "testapp1")
	if err != nil {
		panic("error create container: " + err.Error())
	}

	if err := cli.ContainerStart(ctx, resp.ID,
		types.ContainerStartOptions{
			CheckpointID:  "",
			CheckpointDir: "",
		}); err != nil {
		panic("error start container: " + err.Error())
	}

	fmt.Println("container started")
	time.Sleep(3 * time.Second)

	c := types.ExecConfig{
		User:         "",
		Privileged:   false,
		Tty:          true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Detach:       false,
		DetachKeys:   "",
		Env:          nil,
		WorkingDir:   "",
		Cmd:          []string{"./main"},
	}
	execID, err := cli.ContainerExecCreate(ctx, resp.ID, c)
	if err != nil {
		panic("error container exec create: " + err.Error())
	}
	fmt.Println("exec id", execID, err, resp.ID)

	config := types.ExecStartCheck{
		Detach: false,
		Tty:    true,
	}
	conn, err := cli.ContainerExecAttach(ctx, execID.ID, config)
	if err != nil {
		panic("error container exec attach: " + err.Error())
	}
	defer conn.Close()

	err = cli.ContainerExecStart(ctx, execID.ID, config)
	if err != nil {
		panic("error container exec start: " + err.Error())
	}
	content, _, _ := conn.Reader.ReadLine()
	fmt.Println("content1", string(content))

	fmt.Println("writing...")
	_, err = conn.Conn.Write([]byte("Ivan\n"))
	if err != nil {
		fmt.Println("err write conn:", err.Error())
	}
	content, _, _ = conn.Reader.ReadLine()
	fmt.Println("content2", string(content))

	content, _, _ = conn.Reader.ReadLine()
	fmt.Println("content2", string(content))
}
