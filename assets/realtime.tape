Output realtime.gif

Require tmux
Require campfire

Set Width 1920
Set Height 1080
Set FontFamily "MartianMono NFM"
Set FontSize 14

Hide
  # Start new tmux session with default configuration
  Type "tmux -f /dev/null -L test new-session -- bash" Enter

  # Disable tmux status line
  Type "tmux set status" Enter

  # Set border colors to black
  Type 'tmux setw pane-border-style "fg=white" && \' Enter
  Type 'tmux setw pane-active-border-style "fg=white"' Enter

  # Clear any existing files
  Type "rm vhs.log" Enter

  # Clear Screen
  Ctrl+L

  Sleep 0.5
Show

# Create a horizontal split
Type "tmux split-window -h -d -p 70 -- bash" Enter
Ctrl+L
Sleep 0.5

Ctrl+B
Type o
Sleep 0.5


Type "campfire vhs.log" Enter
Sleep 1

# Swap back to main pane
Ctrl+B
Type o
Sleep 0.5

Type "go run ../demo/demo.go vhs.log" Enter
Sleep 10
