Output filtering.gif

Require tmux
Require campfire

Set Width 1920
Set Height 1080
Set FontFamily "MartianMono NFM"
Set FontSize 14

Hide
  # Clear any existing files
  Type "rm vhs.log" Enter

  Type "go run ../demo/demo.go vhs.log" Enter
  Sleep 20
  Ctrl+C
  Sleep 0.5

  # Clear Screen
  Ctrl+L
Show

Type "campfire vhs.log" Enter
Sleep 2

Type "1"
Sleep 1

Type "3"
Sleep 1

Type "13f"
Sleep 0.5
Type "user"
Sleep 1
Backspace 4
Type "auth"
Sleep 1
Backspace 4
Type "oogabooga"
Sleep 1
Backspace 9
Sleep 3
