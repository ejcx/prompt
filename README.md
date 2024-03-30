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

I create new prompts that I store as a files:
```
$ echo "example prompt, return the word 'foo'" > example-prompt.md
$ prompt add example-prompt.md
```

and list them with the prompt subcommands
```
$ prompt list
- example-prompt.md
- meal-plan.md
- exercise-plan.md
```

I use the `prompt` command to generate todo lists that I write
to my daily notes.
```
$ prompt meal-plan.md | today
```
