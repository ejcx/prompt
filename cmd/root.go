/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

var (
	promptDir   string
	openAIToken string
)

var rootCmd = &cobra.Command{
	Use:   "prompt",
	Short: "ai prompts done well",
	Long:  `Managed your prompts in a directory, invoke them by name.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		run(cmd, args)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all prompts",
	Long:  `List all prompts in the prompt directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		listPrompts()
	},
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add file to prompt list",
	Long:  `add file to prompt list`,
	Run: func(cmd *cobra.Command, args []string) {
		addPrompt(cmd, args)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	dir, ok := os.LookupEnv("PROMPT_DIR")
	if !ok {
		log.Fatal("PROMPT_DIR not set")
	}
	promptDir = dir

	token, ok := os.LookupEnv("OPENAI_TOKEN")
	if !ok {
		log.Fatal("OPENAI_TOKEN not set")
	}
	openAIToken = token

	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(addCmd)
}

func listPrompts() {
	finfos, err := os.ReadDir(promptDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, finfo := range finfos {
		println("-", finfo.Name())
	}
}

func run(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Help()
		return
	}
	promptName := args[0]
	var f = promptDir + "/" + promptName
	buf, err := os.ReadFile(f)
	if err != nil {
		log.Fatalf("could not read prompt %s: %s\n", f, err)
	}

	// send contents to openai
	resp, err := Prompt(string(buf))
	if err != nil {
		log.Fatalf("could not prompt: %s\n", err)
	}
	fmt.Println(resp)

}

func Prompt(query string) (string, error) {
	client := openai.NewClient(openAIToken)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4TurboPreview,
			TopP:  .3,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: query,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil

}

func addPrompt(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Help()
		return
	}
	promptName := args[0]

	// Get file name, strip off path
	fileParts := strings.Split(promptName, "/")
	fname := fileParts[len(fileParts)-1]

	// Check if this prompt already exists
	var f = promptDir + "/" + fname
	_, err := os.Stat(f)
	if err == nil {
		log.Fatalf("prompt %s already exists\n", f)
	}

	// Read the contents
	buf, err := os.ReadFile(promptName)
	if err != nil {
		log.Fatalf("could not read new prompt %s: %s\n", f, err)
	}

	// Write to the prompt dir
	err = os.WriteFile(f, buf, 0644)
	if err != nil {
		log.Fatalf("could not write new prompt %s: %s\n", f, err)
	}
	println("added prompt", f)
}
