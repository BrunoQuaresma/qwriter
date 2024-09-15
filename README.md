# QWriter

QWriter is a command-line tool designed to enhance your text effortlessly. Whether you need to polish your writing or find better wording, QWriter is here to help.

## Installation

To install QWriter, simply run the following command in your terminal:

```bash
curl -fsSL https://raw.githubusercontent.com/BrunoQuaresma/qwriter/main/scripts/install.sh | sh
```

## Setup

Before you can start using QWriter, you'll need to set your OpenAI API key. This key allows QWriter to connect to OpenAI's powerful language models.

1. **Generate an API Key**: Follow the instructions in the OpenAI [quickstart guide](https://platform.openai.com/docs/quickstart/create-and-export-an-api-key) to create your API key.
2. **Export the API Key**: Set the `QWRITER_OPENAI_KEY` environment variable with your API key:
```bash
export QWRITER_OPENAI_KEY=your-key-here
```

## Usage

Once your API key is set, you can start using QWriter to enhance your text. The syntax is straightforward:

```bash
qwriter improve "Please improve this text"
```

QWriter will analyze the provided text and return a refined version, helping you to communicate more effectively.

## Examples

Discover how QWriter CLI can enhance your workflow:

### Generate Better Commit Messages

QWriter can help you craft more meaningful and impactful commit messages. For example:

```bash
git commit -m "$(qwriter improve 'Add commit message example')"
```

### Revise and Improve Documentation

QWriter makes it easy to refine your documentation. You can iterate over all your Markdown files and update them with improved content in one go:

```bash
for file in /path/to/folder/*.md; do qwriter improve "$(cat "$file")" > "$file"; done
```