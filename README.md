# plancks-cli | plancksktl


## plancksktl deploy

Info needed
- docker hub / other repo
- docker TEAM/PROJECT
- git commit 
- Dockerfile to use
- endpoint pc
- service json
- route json

Sources of info
- DCR: `project.json`
- TEAM/PROJECT: `project.json`
- Git commit: terminal command.
- Dockerfile: `project.json`
- endpoint: `project.json`
- service.json: `project.json`
- route.json: `project.json`


Process
- Get **git** commit
- `Docker build -t REPO/TEAM/PROJECT:git .` or that with other dockerfile.
- `Docker push ...`
- apply service.json @endpoint
- apply route.json  @endpoint
- 

## Future 
- for Docker Hub builds to be used
