# Text Editor in Go

> Simple text editor written in Go. It provides a basic set of features that you would expect from a text editor.

## Features

- Undo and Redo
- Select, Cut, Copy, Paste
- Reading and Writing from file
- Cursor movements

## Key Bindings

| Key        | Action                           |
| ---------- | -------------------------------- |
| Ctrl+S     | Save                             |
| ArrowRight | Move Cursor Right                |
| ArrowLeft  | Move Cursor Left                 |
| ArrowUp    | Move Cursor Up                   |
| ArrowDown  | Move Cursor Down                 |
| Ctrl+L     | Move Cursor to End of Line       |
| Ctrl+H     | Move Cursor to Beginning of Line |
| Ctrl+Q     | Toggle Selection Mode            |
| Ctrl+C     | Copy Selection                   |
| Ctrl+X     | Cut Selection                    |
| Ctrl+V     | Paste                            |
| Ctrl+Z     | Undo                             |
| Ctrl+R     | Redo                             |
| Space      | Insert Space                     |
| Enter      | Add New Line                     |
| Backspace  | Remove Character                 |

## Development

To develop this project, you need to have Go installed on your machine. If you don't have Go installed, you can download it from the [official website](https://golang.org/dl/).

Once you have Go installed, clone this repository to your local machine:

```bash
git clone https://github.com/VandanRogheliya/gotext.git
```

### Build

To build the project, navigate to the project directory and run the following command:

```bash
go build -o bin/gotext
```

### Run

To run the project, use the following command:

```bash
./bin/gotext path_to_your_file
```

Replace `path_to_your_file` with the path to the file you want to edit.

## Author

[Vandan Rogheliya](https://github.com/VandanRogheliya)
