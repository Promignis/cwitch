<p align="center">
  <img alt="cwitch logo" src="https://raw.githubusercontent.com/promignis/cwitch/master/assets/logo.png" />
  <p align="center">
      <a href="https://circleci.com/gh/Promignis/cwitch"><img alt="CircleCi Build Status" src="https://circleci.com/gh/Promignis/cwitch.svg?style=shield"></a>
      <a href="/LICENSE"><img alt="Software License" src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=shield"></a>

  </p>
</p>

Conscious Swtich (cwitch, pronounced switch)
is a cross platform system tray app that allows you
create some modes via a json file(`data.json`) and
tracks the time you spend between these modes.

Great way to switch between tasks, modes etc consciously and mindfully.

It's aim is to help you reduce unconscious and mindless activities(browser tab switches, social media, youtube etc)

## Cwitch Json file format

```
{
  "modes": [
    {
      "mode": "Reading",
      "emoji": "📚"
    }
  ]
}
```

#### Example output
![Cwitch output](assets/output.png)

`--data` flag to specify the path to Cwitch Json file
Node: (optional, if not specificed will look for `data.json` in local folder)

eg: `cwitch --data ./data.main.json`


`--debug` flag for seeing debug logs

eg: `cwitch --debug`
