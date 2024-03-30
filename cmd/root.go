/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"
	"os"

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
	println(resp)

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
