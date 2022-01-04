# Native Inserter

Fork of the Clipboard Inserter whose purpose it is to automatically insert the
data received by a native application into a broswer page.

It is useful for use with text extraction tools like Textractor for looking up
words in the browser.

This repository contains the browser plugin, native application, native
messaging manifest and an example client. Textractor plugin is not included.

Currently this addon is not available on any of the browsers addon stores, nor
are releases past the fork signed by our benevolent overlords at
Mozilla/Google.

Only tested on GNU/Linux with Firefox 91 ESR.

## Install

You will need:

- This code
- A go compiler
- A browser supporting native extensions (eg modern Firefox, Chrome)

1) Compile and install the native application
2) Edit `native_inserter.json` to find the native application
3) Put `native_inserter.json` into your browsers native messaging manifest directory
4) Build and install the addon in your browser

For this to be useful you will also need a sender, for example a [Textractor
plugin](https://github.com/45Tatami/Textractor-TCPSender) and a HTML page ([for
example](https://pastebin.com/raw/DRDE075L)) that the can be inserted into.

### Firefox

This addon is not signed. From the [official Firefox
documentation](https://extensionworkshop.com/documentation/publish/signing-and-distribution-overview/):

> Unsigned extensions can be installed in Developer Edition, Nightly, and ESR
> versions of Firefox, after toggling the xpinstall.signatures.required
> preference in about:config.

You can download the unsigned build from the releases or follow the official
instruction for building addons, for example via `web-ext`:

```
$ npm install -g web-ext
$ web-ext build
```

You can then install the resulting zip via the `Install Add-On From File`
dialog under `about:addons`.

The native messaging manifest directory under linux is
`~/.mozilla/native-messaging-hosts`.

## How does it work

The browser plugin starts a native application which creates a raw TCP listen
socket currently on all interfaces on port 30501 for incoming connections.

One or more applications (eg Textractor plugin) will connect to this socket and
send messages. Messages consist of a 4 Byte little-endian length header and
UTF-8 payload.

The native application will forward the UTF-8 data to the browser plugin which
in turn will insert the text into a webpage similar to the original clipboard
version.

## Troubleshooting

The native application is writing to stderr. Browsers usually forward this to
their default log.

The browser addon itself will log to the console. Problems with the native
messaging manifest might not necessarily be logged.

If the native application does not run in the background after loading the
add-on, double-check if the native messaging manifest is in the correct
directory and has the binary path set correctly.

## (Not so) FAQ

#### Why not the clipboard

Abusing the clipboard functionality to transfer data to the browser is ugly. It
does also not work under Wayland which does not allow unfocused applications to
access the clipboard.

(No opinion here whether the native application/messaging setup is more or less
ugly than the clipboard approach)

It also has the bonus of working with VMs without setting up clipboard
forwarding between host and guest.

#### Why a native application

The APIs including TCP listen sockets are no longer supported by modern
Browsers. The only alternative would be connecting to a WebSocket server, but
this has two drawbacks:

1) The API only supports the client side protocol and does not implement a
server. This would mean the sending side would need to be the server and the
native application the client, which would not work with multiple senders.

2) The sending side (eg C++ Textractor plugin) would need to implement a
Websocket server which is a lot more work than raw tcp.

## Bugs

- Native messaging expects the length prefix in native byte order. Go does
  neither have a non-unsafe way to convert an integer to a byte stream in
  native order for writing nor a way to detect endianness. The native
  application will write in little-endian, which will break the protocol on
  big-endian machines
- Messages over the native-messaging 1MB limit will be discarded instead of
  split up
