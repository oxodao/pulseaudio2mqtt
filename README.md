# mqtt2pulseaudio

Pass mosquito MQTT messages to playerctl.

Tried with a [Ikea SYMFONISK remote control](https://www.ikea.com/fr/fr/p/symfonisk-telecommande-volume-blanc-60370480/) and zigbee2mqtt.

Requires playerctl to be installed and mpris-based player (Spotify, web browser, ...)

Usage:
Download the binary release and put it in `/opt/mqtt2pulseaudio/mqtt2pulseaudio` and chmod +x it.

```sh
$ sudo ln -s /opt/mqtt2pulseaudio/mqtt2pulseaudio /usr/bin/mqtt2pulseaudio
```

Copy `config.yaml.dist` to `/etc/mqtt2pulseaudio.yaml` and fill the required stuff.

Copy `contrib/mqtt2pulseaudio.service` to `/etc/systemd/user/mqtt2pulseaudio.service` (shamelessly stolen from [spotifyd](https://github.com/Spotifyd/spotifyd/blob/master/contrib/spotifyd.service))
```sh
$ systemctl enable --now --user mqtt2pulseaudio
```

## Todo

Next step in this "project" is to handle messages not as "symfonisk" messages but rather custom formatted messages, this will require to add a step in node-red to convert the message in the correct format, but this will support other buttons / rotary encoders than the Ikea's Symfonisk.

Payload should be:
```json
{
        "action": "VOL_SET|VOL_DOWN|VOL_UP|VOL_STOP|PLAY|PAUSE|PLAY-PAUSE",
        "mode": "CONTINUOUS|VALUE",
        "value": 100, # Percentage
        "delay": 50, # currently called delta_duration server-side
        "delta_percentage": 1 # currently called delta_percentage server-side
}
```

delay and delta_percentage will be removed from the server-side as they will be contained on the message's body. This will let multiple-manufacturers buttons be able to work together by tuning those value on a per-button basis.

Other thing to do is enable the use of passwordless broker as I think most people don't bother with it since they are running it locally.

Also this should register as an Home Assistant entity exposing each sink / source volumes

## License
```
           DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE
                   Version 2, December 2004
 
Copyright (C) 2004 Sam Hocevar <sam@hocevar.net>

Everyone is permitted to copy and distribute verbatim or modified
copies of this license document, and changing it is allowed as long
as the name is changed.
 
           DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE
  TERMS AND CONDITIONS FOR COPYING, DISTRIBUTION AND MODIFICATION

 0. You just DO WHAT THE FUCK YOU WANT TO.
 ```
