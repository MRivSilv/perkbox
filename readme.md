# PERKBOX

Perkbox is a console based password manager, at the moment only works locally storaging your passwords, but soon there's gonna be a wa to host it by yourself

## Features
* Developed in golang to allow max performance for the user and comodity for the contributors
* "SECURE" Passwords live for 10 seconds in the clipboard, then the clipboard gets emptied
* CLI Interface, ngl I just like terminals, so I don't care about having interfaces, add one if you want to

## Tutorial
* Clone the repo
* Run "go build"
* Syntax goes: `./perkbox <command> <service>`

## Command list
* **"add"**: Use it to add new services to the manager
  * Example: `./perkbox add furryporn.com`
* **"get"**: Get the password to your specific service.
  * Example: `./perkbox get xvideos.com`
* **"delete"**: Same shit as add but to delete, like cmon bro is not that hard

## Master Password
In order to keep every password safe you will need to set a master passwords after adding, this password can be the same or different from the other ones

## Suggestions
I'm accepting suggestions since I'm new to this programming language, like bro, every company is asking if I know about this language and IDK why.
