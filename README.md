<div align="center"">
    <img src="./assets/banner.png" style="width: 700px;"/>
    <h1>Get cozy with your logs 🔥🪵</h1>
</div>

I hated working with log files. They were obtuse, easy to search yet hard to filter, and just plain difficult to look at for long periods of time. But with how much logging we do (some might say too much), we should at least treat viewing them as a first class consideration.

This project came about as I was working with the (ever wonderful) [BubbleTea](https://github.com/charmbracelet/bubbletea) and was getting lost in the flow of messages. I wanted to use normal logging strategies, but obviously that won't work with an app taking the standard I/O. Using `tail` to monitor logs in realtime just felt so... inelegant, especially given the goal is to make the terminal pretty and fun!

Enter... `campfire`!

<div align="center">
    <h2>Why use <code>campfire</code> ❓</h2>
</div>

Make your log files helpful instead of scary and cumbersome.

- Monitor files as they update in real time

![realtime update example](./demo/realtime.gif)

- Filter your files by log type or keyword, hiding things you don't care about

![filtering example](./demo/filtering.gif)

- Continously monitor files by name, whether they exist or not

![file example](./demo/monitoring.gif)

- All of the above at once!

<div align="center">
    <h2>Usage ⚙️</h2>
</div>

- Just run `campfire [file]` with whatever file you want to monitor. That's it!

<div align="center">
    <h2>Installation ⬇️</h2>
</div>

### Github Releases 🐙

- Go to the `Releases` tab of the repo [here](https://github.com/DaltonSW/campfire/releases)
- Download the latest archive for your OS/architecture
- Extract it and place the resulting binary on your `$PATH` and ensure it is executable
```sh
cd ~/Downloads                          # ... or wherever else you downloaded it
tar -xvf campfire_[etc].tar.gz          # x: Extract; v: Verbose output; f: Give filename
chmod +x campfire                       # Make file executable
mv campfire [somewhere on your $PATH]   # Move the file to somewhere on your path
```

### Homebrew 🍺 

- Have `brew` installed ([brew.sh](https://brew.sh))
- Run the following:
```sh
brew install daltonsw/tap/campfire
```

### Go 🖥️ 

- Have `Go` 
- Have your `Go` install location on your `$PATH`
- Run the following: 
```sh
go install go.dalton.dog/campfire@latest
```

<div align="center">
    <h2>Credits 🗨️</h2>
</div>

- [Campfire Icon - Created by Freepik/Flaticon](https://www.flaticon.com/free-icons/campfire)

<div align="center">
    <h2>License ⚖️</h2>
</div>

Copyright 2025 - Dalton Williams  
Check [LICENSE](./LICENSE.md) in repo for full details
