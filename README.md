# fpush

fpush is a small tool to send push notifications when photos
in a given folder change.

The detected photo will be attached to the notification.

It uses
[Pushover](https://pushover.net)
to send push notifications to devices,
so you will need an API key and a device key.

## Usage

Set these environment variables:

* `PUSHOVER_KEY` Your API key
* `PUSHOVER_RECIPIENT_KEY` Your device key

Run the program like this:

`fpush -path=/foo/bar`

If you leave out the `-path` flag
it will fall back to the current directory.

## Background

I use this for my security camera
which is a Raspberry Pi with a camera module
that uses
[motion](https://motion-project.github.io)
to detect changes in photos.
I run fpush every five minutes as a cron job
to see if there are new photos in motion's photos folder
and notify me about them.
