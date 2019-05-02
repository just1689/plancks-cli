# plancks-cli | plancksktl


## plancksktl deploy

Info needed
- docker TEAM/PROJECT
- git commit 
- Dockerfile to use
- endpoint pc
- service json

Sources of info
- TEAM/PROJECT: `project.json`
- Git commit: terminal command.
- Dockerfile: `project.json`
- endpoint: `project.json`
- service.json: `project.json`


Process
- Get **git** commit
- `Docker build -t TEAM/PROJECT:git .` or that with other dockerfile.
- `Docker push ...`
- mangle, apply service.json @endpoint
- 

## Future 
- for Docker Hub builds to be used
