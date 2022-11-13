# gitgc

A small utility for easily cloning github repositories to a pre-configured directory.

## Get

```
go install github.com/mrahhal/gitgc@latest
```

## Usage

```
gitgc [user]/[repo]
```

A config file called ".gitgc" in your home directory ("~/.gitgc" on unix, "%USERPROFILE%/.gitgc" on windows) will be used to determine the base directory where repos should be cloned. The ".gitgc" file should contain a single line that represents that directory.

Example ".gitgc" file:

```
C:\dev\git
```

It will be automatically created on the first run, and will default to "[user home]/git".
