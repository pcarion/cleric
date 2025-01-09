<br/>
<p align="center">
<img src="assets/cleric-logo.svg" width="100" alt="Cleric logo">
</p>
<br/>

# Cleric

Cleric is a desktop application that helps you manage Model Context Protocol (MCP) servers for [Claude desktop](https://claude.ai/download). It provides an easy-to-use interface for adding, removing, and configuring your MCP servers.

## How It Works

Cleric manages two configuration files:

1. A personal server list at `$HOME/.cleric.json`
   - Stores all your available MCP servers
   - Acts as a backup of your server configurations

2. The Claude Desktop configuration at `<OS specific path>/claude_desktop_config.json`
   - Contains the active MCP servers currently used by Claude Desktop

## Key Features

- Seamlessly switch MCP servers in and out of Claude Desktop
- Automatically preserves server configurations when making changes
- Syncs with Claude Desktop: Any servers you manually add to Claude Desktop will be automatically added to your Cleric server list on next launch
- Real-time saving: Changes are saved automatically as you make them


## How to use

On the left side of the application, you can see the list of MCP servers that you have configured.

The servers with a checkmark are the ones that are currently configured in Claude desktop.

The right side of the application shows the information about the selected server:

![MCP server information](assets/cleric-screenshot-01.png)

The `Remove From Claude` button will remove the server from Claude desktop but the server will remain in your personal server list.

The `✄` button will remove the server from both your personal server list and the list of servers that are currently configured in Claude desktop.

The `🖊️` button will allow you to edit the server information:

![MCP server information](assets/cleric-screenshot-02.png)


## Important Note

Before using Cleric for the first time, please back up your `claude_desktop_config.json` file if you have manually configured MCP servers in Claude Desktop.


# Development

## Prerequisites

To build and run this project, you need:

- Go 1.23.4 or later
- Fyne toolkit and its dependencies - read the [Getting Started with Fyne](https://docs.fyne.io/started/) documentation for more information

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

### Running the Application

To run the application directly:
```bash
make run
```

## Project Structure

```
.
├── cmd/
│   └── cleric/          # Main application entry point
├── internal/
│   └── ui/             # UI components and logic
│   └── config/         # Code to manage the list of MCP servers
└── build/              # Build artifacts
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

