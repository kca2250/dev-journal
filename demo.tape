# djou Demo
# Run: vhs demo.tape

Output demo.gif

Set Shell "bash"
Set FontSize 16
Set Width 800
Set Height 500
Set Framerate 30
Set Theme "Dracula"
Set Padding 20

# Set language to English
Type "djou config set language en"
Sleep 500ms
Enter
Sleep 1s

# Title
Type "# djou - Dev Journal CLI"
Enter
Sleep 1s

# Quick Record
Type "djou record -q 'Implemented user authentication' -e 2 -a 2.5 -ai 45 -t 'feature,backend'"
Sleep 500ms
Enter
Sleep 2s

# Another record
Type "djou record -q 'Fixed CORS issue' -e 1 -a 0.5 -ai 15 -t 'bugfix,api'"
Sleep 500ms
Enter
Sleep 2s

# List entries
Type "djou list"
Sleep 500ms
Enter
Sleep 3s

# Stats
Type "djou stats"
Sleep 500ms
Enter
Sleep 3s

# End
Type "# Track your development journey!"
Enter
Sleep 2s
