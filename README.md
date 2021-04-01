# universal-tip-script
Swiss army knife script for [Universal Tip](https://github.com/tanin47/tip#readme) tool

# Installation

## Dependencies

1. Have [Universal Tip](https://github.com/tanin47/tip#readme) installed
1. Have Go 1.14 or above installed
  `brew install go`

## Instructions

```
mv provider.script default.script
mkdir -p "~/Library/Application Scripts/tanin.tip" && cd "~/Library/Application Scripts/tanin.tip"
git clone https://github.com/JamesDunne/universal-tip-script .
go build
```

`go build` should produce a single binary named `provider.script` in the current folder. Universal Tip will invoke this
executable and pass it the selected text.

# Usage

Select some text on your screen in any application.

Press the shortcut key for Universal Tip. The default shortcut key is `Cmd + Shift + 7` aka `Cmd + &`. Follow instructions
on their site if you want to change this shortcut key.

Universal Tip will invoke the `provider.script` executable that we compiled using Go above and will pass it the selected text
as the first command-line argument.

The executable will attempt to interpret the text in multiple ways and provide several kinds of useful transformations on it.

Currently it can:
* base64 encode and decode
* generate a uuid
* parse RFC3339 timestamps and adjust timezones
* parse integer timestamps expressed in sec, msec, or nsec units and convert between them

Feel free to customize the code to add more transformations and utilities.
