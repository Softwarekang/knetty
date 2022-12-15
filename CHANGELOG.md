## Content
- [Content](#content)
- [[2022-12-15] [Version v1.0.1]](#2022-12-15-version-v1-0-1)
  - [Feature](#feature)
  - [BugFix](#bugfix)

## [2022-12-15] [Version v1.0.1]

### Feature
- feat(client):Added client impl and client usage case(#28)
- feat(*): Implement a knet server object that manages session connections and user interfaces. (#27)
- feat(server): define server options and main func (#26)
- feat(example): add example for session (#22)
- feat(session): impl knet session func (#20)
- feat(buffer):add kNet buffer interface and byteBuffer
- feat(knet): add kqueue impl for poller (#3)
- feat(knet): Defining the poll Interface (#2)

### BugFix
- fix(session): fix session handleTcpPackage buffer read bug (#23)
- fix(connection): fix connection write buffer bug (#19)
- fix(example): fix kqueue example server connection read bug (#5)
