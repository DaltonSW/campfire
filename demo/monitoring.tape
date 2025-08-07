Output monitoring.gif

Source shared.tape
Require tmux

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
Sleep 1

Type "echo INFO Hello world! >> vhs.log" Enter
Sleep 2
Type "echo WARN Uh oh! Incoming... >> vhs.log" Enter
Sleep 2
Type "echo ERROR OH NO FILE GETTING DELETED >> vhs.log" Enter
Sleep 2
Type "rm vhs.log" Enter
Sleep 2
Type "echo INFO Oh awesome! No crash! >> vhs.log" Enter
Sleep 2
Type "echo INFO Nope! Continuous monitoring! >> vhs.log" Enter
Sleep 5

