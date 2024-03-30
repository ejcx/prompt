# prompt
I'm using this to juggle the chatgpt prompts that I use to write
to my Obsidian notes and generate todo lists.

# How I use it
I have all of this in my zshrc file.
```
export PROMPT_DIR=~/Obsidian/prompts/
export DAILY_DIR=~/Obsidian/Daily/
export OPENAI_TOKEN="....."
alias today="cat >> $DAILY_DIR$(date '+%Y-%m-%d').md"
```

A create new prompts and list them with the prompt subcommands
```
$ prompt list
- example-prompt
- meal-plan
- exercise-plan
```

I use the `prompt` command to generate todo lists that I write
to my daily notes.
```
$ prompt meal-plan | today
```
