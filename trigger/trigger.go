package trigger

import (
	"bytes"
	"log"
	"os"
	"os/exec"
)

func SetupPipeline(output string) {
	_, err := os.Stat(output)
	if err != nil {
		os.Mkdir(output, 0755)
	}
}

func TriggerBuildPipeline(cmds []string) {
	var stdout, stderr bytes.Buffer
	log.Println("Starting commands pipeline.")
	for _, command := range cmds {
		log.Println("Running command:", command)
		cmd := exec.Command("sh", "-c", command)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err := cmd.Run()
		cmd.Wait()

		if err != nil {
			log.Printf("Error: %s", stderr.String())
		}
		if stdout.String() != "" {
			log.Println("Output:", stdout.String())
		}
	}	
}
