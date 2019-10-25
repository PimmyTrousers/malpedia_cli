# malpedia_cli

[![Go Report Card](https://goreportcard.com/badge/github.com/PimmyTrousers/malpedia_cli)](https://goreportcard.com/report/github.com/PimmyTrousers/malpedia_cli)

Malpedia_cli is a tool to interact with the malpedia service located [here](https://malpedia.caad.fkie.fraunhofer.de). Some of the endpoints commands require an api key due to restrictions with the service itself. It simplifies some of the endpoints and exposes the features that I beleive are the most important. 

## Configuration of the tool
The application requires an API for some of the endpoints, which can be passed by arugment or a YAML file at `$HOME/.malpedia_cli.yaml`. Currently it only allows for an apikey, so an example would look like the following 

```
apikey: <apikey>
```

## Currently supported features
- [X] download samples via hash 
- [X] get a list of all tracked actors 
- [X] get information about a specific actor 
- [X] get a list of all tracked malware families 
- [X] get information about a specific malware family 
- [X] download yara rules by TLP level 
- [X] download yara rules by family 
- [X] scan malpedia's malware catalog against a yara rule
- [X] validate API keys 
- [X] get the malpedia version
- [X] get all hashes for a family 
- [X] download all samples from a family
## Images 
![Ursnif output](images/Screen&#32;Shot&#32;2019-08-25&#32;at&#32;11.42.33&#32;AM.png)

![FIN7 output](images/Screen&#32;Shot&#32;2019-09-13&#32;at&#32;7.14.04&#32;PM.png)
## TODO
- [X] Command to download all samples from a family 
- [X] Scan malpedia's malware catalog against a yara rule
- [ ] Upload a file to be checked against yara rules (in the works)
- [ ] Generic search (will return a family or actor)
- [ ] Download all samples from an actor
- [ ] Verbose logging 
- [ ] Enable user choice if multiple results are returned for fuzzy search
- [ ] Support contexts
- [X] Reject commands that require an API key when one isnt applied

## Examples 
```
- malpedia_cli version
- malpedia_cli getYaraRules white
- malpedia_cli getYaraRules amber -z -o yara_rules.zip
- malpedia_cli getSample 12f38f9be4df1909a1370d77588b74c60b25f65a098a08cf81389c97d3352f82 -p infected123 -o samples.zip
- malpedia_cli getSample 12f38f9be4df1909a1370d77588b74c60b25f65a098a08cf81389c97d3352f82 -r 
- malpedia_cli actors --json
- malpedia_cli actor apt28
- malpedia_cli scanYara RAT_Nanocore.yar
- malpedia_cli families
- malpedia_cli downloadFamily ursnif
- malpedia_cli downloadYara ursnif 
- malpedia_cli downloadYara njrat -o njrat.zip
- malpedia_cli scanYaraAgainstFamily carbanak myRule.yar
```

## Build Instructions
Create a binary file at your current directory
```
go build -o ./malpedia_cli
```
Create a binary file and install it in your path
```
go install
```