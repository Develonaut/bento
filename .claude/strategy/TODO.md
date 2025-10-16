A list of TODOs to look into for improving the codebase.

1. Bubbletea Tabs. It might benefit us to organize the views using tabs.
    - https://github.com/charmbracelet/bubbletea/blob/main/examples/tabs/main.go
2. Bubbles Help and and Key components
    - We should leverage the Bubbles help and key components more for our footer help menu.
    - https://github.com/charmbracelet/bubbletea/blob/main/examples/help/main.go
    - https://github.com/charmbracelet/bubbles?tab=readme-ov-file#key
3. Simplify Editor Have Two Section
   1. Top Section is a Table View of the Nodes in the Bento
   2. Bottom Section is the interactive huh form for editing/creating nodes
    - As you add nodes they get added to the table view above. This should work well with recursive nodes like loops and groups. since we should be able to render another table inside a table cell.
    - A future improvements would be allowing arrow keys to move around the table and select nodes to edit. and have the bottom section be huh form driven for editing or to make changes.
    - https://github.com/charmbracelet/bubbletea/blob/main/examples/table/main.go
    - https://github.com/charmbracelet/huh
    - This would simplify the editor view and make it easier to understand.
    - We could use a table to display the nodes in order, with nested tables for loops or groups.
    - Any Node Editing/Creation would use Huh forms.
4. Make Header a bit nicer with some styling.
    -  🍱   Bento v0.0.1
            View: Browser
    - Two line header with same nice background styling we have just make the emoji bigger to fill the two row header

5. Make Executor Progress Feedback Nicer with Sequence Bubbletea component

[Already Existing Bento Info Section]


🍱 Preparing Bento...
[]string{"🍣", "🍙", "🥢", "🍥"} [Tasting|Sampling|Savoring|Finished] [Get Hello Message]… [Execution time]
[]string{"🍣", "🍙", "🥢", "🍥"} [Tasting|Sampling|Savoring|Finished] [Extract Slideshow Title]… [Execution time]
[]string{"🍣", "🍙", "🥢", "🍥"} [Tasting|Sampling|Savoring|Finished] [Echo Result]… [Execution time]
🍱 Bento Finished. Delicious✨!

[Output] (We need a nice way to handle large outputs here without them taking over the whole screen)

[Already Existing Progress Bar/Status Section]

6. Add VHS Charm functionality for capturing gifs of the TUI in action for Github README
    - https://github.com/charmbracelet/vhs
7. Review https://leg100.github.io/en/posts/building-bubbletea-programs/ for best practices and see if we can improve our codebase.
