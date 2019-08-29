# malpedia_cli

Malpedia_cli is a tool to interact with the malpedia service located [here](https://malpedia.caad.fkie.fraunhofer.de). Some of the endpoints commands require an api key due to restrictions with the service itself. It simplifies some of the endpoints and exposes the features that I beleive are the most important. 

## Configuration of the tool
The application requires an API for some of the endpoints, which can be passed by arugment or a YAML file at `$HOME/.malpedia_cli.yaml`. Currently it only allows for an apikey, so an example would look like the following 

```
apikey: <apikey>
```

## Currently supported commands
<<<<<<< HEAD
- get samples via hash 
- get a list of all tracked actors 
- get information about a specific actor 
- get a list of all tracked malware families 
- get information about a specific malware family 
- get yara rules by TLP level 
=======
- [X] get samples via hash 
- [X] get a list of all tracked actors 
- [X] get information about a specific actor 
- [X] get a list of all tracked malware families 
- [X] get information about a specific malware family 
- [X] get yara rules by TLP level 
- [X] get yara rule by family 
- [X] get the malpedia version
- [X] get all hashes for a family 
- 
>>>>>>> 09f1726... rough implementation of post request and scanBinary against yara rules

## Commands to be added 
- download all samples from a family 
- upload a file to be checked against yara rules 
- generic search (will return a family or actor)
- get all samples from an actor 

## Examples 
```
- malpedia_cli version
- malpedia_cli getYaraRules white
- malpedia_cli getSample 12f38f9be4df1909a1370d77588b74c60b25f65a098a08cf81389c97d3352f82 -p infected123 -o samples1234.zip
- malpedia_cli getSample 12f38f9be4df1909a1370d77588b74c60b25f65a098a08cf81389c97d3352f82 -r 
- malpedia_cli getActors --json
- malpedia_cli getActor apt28
- malpedia_cli getFamilies
- malpedia_cli getFamily ursnif
- malpedia_cli getYara ursnif 
- malpedia_cli getYara njrat -o njrat.zip
```