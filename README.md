# GoTypo

simple, fast, and kinda-smart spellchecker written in go. no more typos ruining your vibe.
## What's the deal?

GoTypo is a lightweight spellchecker that uses:

  - trie (for lightning-fast dictionary lookup)
  - levenshtein distance (fancy name for typo-catching math)
  - a basic suggestion engine so your typos don’t embarrass you

its built for speed and simplicity, zero bloated dependencies, just clean go code. and also i love to code in go.
## How to use it?

first, clone this repo:
```bash
git clone https://github.com/mvarr/GoTypo.git
cd GoTypo
```
make sure you have a file named words.txt in the project folder. each line = one word, easy stuff.

then, run it:
```bash
go run main.go
```
boom. enter a word, hit enter, and see if you messed up.

type exit to stop the madness.
## How it works?

  - trie: checks instantly if your word exists.
  - levenshtein math: finds words that are kinda close to what you typed. despite i hate math, its kinda fun to write.
  - suggestions: shows a tidy list of up to 10 corrections. (in other way it spams like crazy, so i have to limit to 10)

## Possible Next Steps (TODO)

  - maybe throw a basic gui in there.
  - improve suggestion algorithm (markov chain style? well see.)
  - option to add custom dictionaries so you can flex your weird vocab.

## Contribute?

sick. drop a PR, open an issue, or just fork it and do your own thing. keep it clean, though nobody likes messy code.

## License

MIT. basically, do whatever you want just dont blame me if things explode.