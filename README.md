# Native Inserter

Fork of the Clipboard Inserter whose purpose it is to automatically insert the
data received by an in this repository included native application into a
broswer page.

It is useful for use with text extraction tools like Textractor for looking up
words in the browser.

This repository contains the browser plugin, native application and native
manifest. Textractor plugin is not included.

Currently this addon is not available on any of the browsers Addon Store.

Only tested on GNU/Linux with Firefox.

## Install

You will need:

 - This code
 - A go compiler
 - A browser supporting native extensions (eg modern Firefox, Chrome)

1) Compile and install the native application
2) Put the native application manifest into your browsers expected directory
3) Install the addon in your browser

For this to be useful you will also need a sender, for example a Textractor
plugin (not included) and a browser page that the can be inserted to.

## How does it work

The browser plugin starts a native application which creates a raw TCP listen
socket for incoming connections.

One or more applications (eg Textractor plugin) will connect to this socket and
send messages. Messages consist of a 4 Byte little-endian length header and
UTF-8 payload.

The native application will forward the UTF-8 data to the browser plugin which
in turn will insert the text into a webpage similar to the original clipboard
version.

## Troubleshooting

The native application is writing to stderr. Browsers usually forward this to
their default log.

The browser addon itself 

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

2) The sending side (eg C++ Textractor plugin) will need to implement a
Websocket server which is a lot more work than raw tcp.

## Bugs

- Native messaging requires native host order prefixed messages. Golang does
  not have an easy way to retrieve this information. The native application
  will use little-endian which is the default on most common architectures
- Messages over the native-messaging 1MB limit will be discarded instead of
  split up
