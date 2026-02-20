my-app/
├── main.go              # Entry point only — initializes and runs the program
├── model/
│   ├── model.go         # Root model (top-level Model struct, Init, Update, View)
│   └── keys.go          # Key bindings (keymap structs)
├── ui/
│   ├── views/
│   │   ├── home.go      # Each "screen" or "page" as its own model
│   │   ├── detail.go
│   │   └── form.go
│   └── components/
│       ├── list.go      # Reusable components (wrap bubbles components)
│       ├── spinner.go
│       └── table.go
├── styles/
│   └── styles.go        # Lipgloss styles, colors, theme
└── cmd/                 # (optional) if you have multiple binaries
