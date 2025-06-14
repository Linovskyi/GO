# GoQueensResolver

GoQueensResolver is a Go application that detects colored regions on a chessboard-like grid in an image and solves a variant of the N-Queens problem, placing one queen per color region such that no two queens threaten each other. The solution is visualized by overlaying queen emojis on the original image.

## Features

- Automatically detects grid regions and their colors from an input image.
- Solves the N-Queens problem with the constraint of one queen per color region.
- Outputs the solution as a PNG image with queen emojis drawn on the board.

## Usage

1. Place your input image (e.g., `photo.jpg`) on your Desktop or specify the path as a command-line argument.
2. Run the program:

   ```sh
   go run Main.go [optional-path-to-image]
   ```

3. The solution will be printed in the console and saved as `result_with_queens.png` in the current directory.

## Example

**Before:**
![photo](https://github.com/user-attachments/assets/91b4a260-f68a-4481-85bb-9b7d4b88fab2)

**After:**

![result_with_queens](https://github.com/user-attachments/assets/3bbb0a4d-815f-4a9d-bebe-9b653108e877)

## Dependencies

- [github.com/fogleman/gg](https://github.com/fogleman/gg)
- [github.com/golang/freetype](https://github.com/golang/freetype)
- [golang.org/x/image](https://pkg.go.dev/golang.org/x/image)

Install dependencies with:

```sh
go mod tidy
```

## License

MIT License

---

*Developed by Moonlight.*
