# malpedia_cli

[![Go Report Card](https://goreportcard.com/badge/github.com/PimmyTrousers/malpedia_cli)](https://goreportcard.com/report/github.com/PimmyTrousers/malpedia_cli)
[![Build Status](https://travis-ci.org/PimmyTrousers/malpedia_cli.svg?branch=master)](https://travis-ci.org/PimmyTrousers/malpedia_cli)

Malpedia_cli is a tool to interact with the malpedia [service](https://malpedia.caad.fkie.fraunhofer.de). Some of the endpoints commands require an api key due to restrictions with the service itself but the tool will tell you if you need one or not for the request. Its goal is to simplify usage and allows users to seamlessly work with the resources contained with the malpedia service. 

Malpedia_cli can be used for getting information about a actor, getting information about a malware family, acquiring samples, uploading yara rules, downloading yara rules, and uploading samples to be scanned against their malware corpus.

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
![Ursnif output](images/Screen&#32;Shot&#32;2019-11-02&#32;at&#32;1.41.53&#32;PM.png)

![FIN7 output](images/Screen&#32;Shot&#32;2019-11-02&#32;at&#32;1.42.31&#32;PM.png)

![Yara scan results](images/Screen&#32;Shot&#32;2019-11-02&#32;at&#32;1.40.19&#32;PM.png)
## TODO
- [X] Command to download all samples from a family 
- [X] Scan malpedia's malware catalog against a yara rule
- [ ] Remove apikey argument from functions that don't need it 
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