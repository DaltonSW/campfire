Output filtering.gif

Require tmux
Require campfire

Set Width 1920
Set Height 1080
Set FontFamily "MartianMono NFM"
Set FontSize 14
Set TypingSpeed 0.5

Hide
  # Clear any existing files
  Type "rm vhs.log" Enter

  Type "go run . vhs.log" Enter
  Sleep 40
  Ctrl+C
  Sleep 0.5

  # Clear Screen
  Ctrl+L
Show

Type@200ms "campfire vhs.log" Enter
Sleep 2

Type "1"
Sleep 1

Type "3"
Sleep 1

Type "13f"
Sleep 2
Type "user"
Sleep 3
Backspace 4
Sleep 2
Type "auth"
Sleep 3
Backspace 4
Sleep 2
Type "oogabooga"
Sleep 3
Backspace 9
Sleep 3
