# Wachat

## Install

Download the latest release from the [releases page](https://github.com/threeproto/wachat/releases)

### Linux

1. download the binary `wachat-linux`
2. make it executable with command, `chmod +x wachat-linux`
3. now you can open and use the `wachat-linux` app.

### Mac

1. download and unzip the binary `wachat-darwin.app.zip`, and you should get a `wachat` app.
2. move `wachat` app to the Applications folder.
3. open Applications folder, find `wachat` app, right click and select `Open` to open it.

### Windows

_Note_: build is broken for windows.

1. download the binary `wachat-windows.exe`.
2. double click the `wachat-windows.exe` to open it and click `Yes` if asking for permission.


## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.
