## Gotodo

A to-do list, how novel! Practicing some Golang based on [Dreams of Code's list of golang projects](https://github.com/dreamsofcode-io/goprojects/tree/main/01-todo-list).

The **Technical Considerations** and **Extra Features** there haven't been addressed yet.

### Key takeaways:

- How golang opens, closes, writes and reads from files, including permissions.
- CSVs are handled weirdly in `encoding/csv` -- why no update record function?
- Lots of string manipulation and printing
- First CLI from scratch in Golang, whoo
- Cobra and `cobra-cli` are very nice to work with from the ground-up.
- Scopes and shadowing are tough, golang is nice and typed like that
- Managing concurrency is made easier by Go, but it still can be hard to get ahold of all the moving parts.

### What's next?:

- Implement breaks between pomodoros
- Break code up into functions more
- Write tests
- Deployment and DevOps stuff
- Convert to a TUI with [bubbletea](https://github.com/charmbracelet/bubbletea)
