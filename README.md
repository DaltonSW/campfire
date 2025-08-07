<div align="center"">
    <img src="./assets/banner.png" style="width: 700px;"/>
    <h2>Cozy up to your logs 🔥🪵</h2>
</div>

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
cd ~/Downloads # Assuming you downloaded it here
tar -xvf campfire_[whatever].tar.gz # x: Extract; v: Verbose output; f: Specify filename
chmod +x campfire # Make file executable
mv campfire [somewhere on your $PATH] # Move the file to somewhere on your path for easy execution
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
