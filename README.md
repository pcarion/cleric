<br/>
<p align="center">
<img src="assets/cleric-logo-color.svg" width="225" alt="Cleric logo">
</p>
<br/>

# Cleric
Cleric is a desktop application that provides a graphical user interface to configure the list of Model Context Proptocol (MCP) servers supported by [Claude desktop](https://claude.ai/download).


## Prerequisites
To build and run this project, you need:

* Go 1.23 or later
* [Fyne toolkit](https://docs.fyne.io/) and its dependencies

## Features

- Server management capabilities
- Cross-platform compatibility (thanks to Fyne and Go)
- Dark theme support

## Prerequisites

To build and run this project, you need:

- Go 1.23.4 or later
- Fyne toolkit and its dependencies
- golangci-lint (for development)

## Installation

- Clone the repository:

```bash
git clone https://github.com/pcarion/cleric.git
cd cleric
```

2. Build the project:

```bash
make build
```

The built binary will be available in the `build` directory.

## Development

### Running the Application

To run the application directly:
```bash
make run
```

### Available Make Commands

- `make build`: Builds the application
- `make run`: Runs the application directly
- `make bundle`: Bundles assets into the application
- `make package`: Creates a distributable package

### Project Structure

```
.
├── cmd/
│   └── cleric/          # Main application entry point
├── internal/
│   └── ui/             # UI components and logic
│   └── config/         # Code to manage the list of MCP servers
├── build/              # Build artifacts
└── cleric.app/         # Application bundle
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the terms found in the `LICENSE` file in the root directory.

## Version

The application version information is embedded at build time from `cmd/cleric/version.txt`.

# Documentation
* [Fyne documentation](https://docs.fyne.io/)
* [Fyne demo source code](https://github.com/fyne-io/fyne/blob/master/cmd/fyne_demo/main.go)

# Acknowledgments

This project was my introduction to the Fyne toolkit. Special thanks to Alex Ballas's [go2tv](https://github.com/alexballas/go2tv) project, which served as an excellent reference. His well-structured codebase demonstrates Fyne's capabilities beautifully and provided valuable insights for building this application.

