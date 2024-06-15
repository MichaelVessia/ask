# Ask

Ask ai stuff, and then easily execute it in command line.

## Usage

Need an api key, set that up with GROQ:

```bash
export GROQ_API_KEY="your_groq_api_key_here"
```

Ensure ~/ask-output.txt exists:

The runit.sh script expects ~/ask-output.txt to contain a valid shell command.

```bash
touch ~/ask-output.txt
```

Move the binary for `ask` and `runit` onto your path. For example:

```bash
sudo cp ask /usr/local/bin/
sudo cp runit /usr/local/bin/
```

Then you can just do something like this:

```bash
$ ask "how do i count the number of lines in main.go?"
 wc main.go
$ runit
 101  289 2256 main.go
```
