## Gotodo

A to-do list, how novel! Practicing some Golang based on [Dreams of Code's list of golang projects](https://github.com/dreamsofcode-io/goprojects/tree/main/01-todo-list).

![A demonstration of the Gotodo CLI: adding tasks, deleting, completing, and listing uncompleted and all tasks](https://github.com/EmilyBurak/gotodo/blob/main/render1732725257933.gif)

The **Technical Considerations** and **Extra Features** there haven't been addressed yet.

### Key takeaways:

- How golang opens, closes, writes and reads from files, including permissions.
- CSVs are handled weirdly in `encoding/csv` -- why no update record function?
- Lots of string manipulation and printing
- First CLI from scratch in Golang, whoo
- Cobra and `cobra-cli` are very nice to work with from the ground-up.

### What's next?:

- Convert to a TUI with [bubbletea](https://github.com/charmbracelet/bubbletea)
