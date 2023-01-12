# Changelog

## [Unreleased](https://github.com/softwarekang/knetty/tree/HEAD)

[Full Changelog](https://github.com/softwarekang/knetty/compare/v0.1.0...HEAD)

**Implemented enhancements:**

- feat: knetty need provides SetLogger func for users to set their own logger [\#57](https://github.com/Softwarekang/knetty/issues/57)

## [v0.1.0](https://github.com/softwarekang/knetty/tree/v0.1.0) (2023-01-12)

[Full Changelog](https://github.com/softwarekang/knetty/compare/58765016485041135cadceb568df7c1640c53dd7...v0.1.0)

**Implemented enhancements:**

- feat:Implementing a lock-free ringbuffer for use by the connection  [\#54](https://github.com/Softwarekang/knetty/issues/54)
- \[Feature Request\]:An example of using knetty to implement an HTTP server [\#52](https://github.com/Softwarekang/knetty/issues/52)
- \[Feature Request\]:Optimize the connection bytebuffer cache mechanism and reduce GC by using a pool. [\#51](https://github.com/Softwarekang/knetty/issues/51)
- \[Feature Request\]: Implementing a non-blocking read function using chan in place of a for loop. [\#49](https://github.com/Softwarekang/knetty/issues/49)
- \[Feature Request\]: Session interface exposes APIs that are not needed by the user and need to be optimized. [\#45](https://github.com/Softwarekang/knetty/issues/45)
- \[Feature Request\]:Supports knetty configuration files and logging framework [\#40](https://github.com/Softwarekang/knetty/issues/40)
- \[Feature Request\]:Knetty requires support for running on LINUX systems. [\#39](https://github.com/Softwarekang/knetty/issues/39)
- Creating issue templates [\#35](https://github.com/Softwarekang/knetty/issues/35)
- chore\(ci\): ci for golint & mdlint [\#32](https://github.com/Softwarekang/knetty/issues/32)
- knet session api [\#17](https://github.com/Softwarekang/knetty/issues/17)

**Fixed bugs:**

- bug: Graceful shutdown of client and server failed [\#47](https://github.com/Softwarekang/knetty/issues/47)

**Merged pull requests:**

- feat\(\*\): provides SetLogger func for user set custom Logger [\#58](https://github.com/Softwarekang/knetty/pull/58) ([Softwarekang](https://github.com/Softwarekang))
- feat\(buffer\): Supporting a lock-free and thread-safe RingBuffer to address the issue of memory reuse [\#56](https://github.com/Softwarekang/knetty/pull/56) ([Softwarekang](https://github.com/Softwarekang))
- feat\(example\): add http server example [\#55](https://github.com/Softwarekang/knetty/pull/55) ([Softwarekang](https://github.com/Softwarekang))
- chore\(issues\_templates\): feature\_request [\#53](https://github.com/Softwarekang/knetty/pull/53) ([Chever-John](https://github.com/Chever-John))
- feat\(connection\): pref connection read without timeout [\#50](https://github.com/Softwarekang/knetty/pull/50) ([Softwarekang](https://github.com/Softwarekang))
- fix\(server \): fix the bug that causes graceful shutdown of the server and client to be ineffective. [\#48](https://github.com/Softwarekang/knetty/pull/48) ([Softwarekang](https://github.com/Softwarekang))
- feat\(session\): Optimizing and reducing the exposure of session and connection APIs to the user [\#46](https://github.com/Softwarekang/knetty/pull/46) ([Softwarekang](https://github.com/Softwarekang))
- Create LICENSE [\#44](https://github.com/Softwarekang/knetty/pull/44) ([Softwarekang](https://github.com/Softwarekang))
- feat\(poll\): Using the epoll  implement multiplexing. [\#43](https://github.com/Softwarekang/knetty/pull/43) ([Softwarekang](https://github.com/Softwarekang))
- chore\(issue\_template\): add two template [\#42](https://github.com/Softwarekang/knetty/pull/42) ([Chever-John](https://github.com/Chever-John))
- chore\(CI && doc\): support lint for markdown && go and update README [\#41](https://github.com/Softwarekang/knetty/pull/41) ([Chever-John](https://github.com/Chever-John))
- Fix typo of sever to server [\#36](https://github.com/Softwarekang/knetty/pull/36) ([sunbinnnnn](https://github.com/sunbinnnnn))
- chore: update golang workflow and support ci-lint for golang [\#34](https://github.com/Softwarekang/knetty/pull/34) ([Chever-John](https://github.com/Chever-John))
- feat\(client\):Added client impl and client usage case [\#31](https://github.com/Softwarekang/knetty/pull/31) ([Softwarekang](https://github.com/Softwarekang))
- chore\(\*\): rename project knet to knetty [\#30](https://github.com/Softwarekang/knetty/pull/30) ([Softwarekang](https://github.com/Softwarekang))
- docs\(README\): update readme [\#29](https://github.com/Softwarekang/knetty/pull/29) ([Softwarekang](https://github.com/Softwarekang))
- docs\(README\): update readme [\#28](https://github.com/Softwarekang/knetty/pull/28) ([Softwarekang](https://github.com/Softwarekang))
- feat\(\*\): Implement a knet server object that manages session connecti… [\#27](https://github.com/Softwarekang/knetty/pull/27) ([Softwarekang](https://github.com/Softwarekang))
- feat\(server\): define server options and main func [\#26](https://github.com/Softwarekang/knetty/pull/26) ([Softwarekang](https://github.com/Softwarekang))
- pref\(\*\): pref kNet project annotation and test case [\#25](https://github.com/Softwarekang/knetty/pull/25) ([Softwarekang](https://github.com/Softwarekang))
- fix\(session\): fix session handleTcpPackage buffer read bug [\#23](https://github.com/Softwarekang/knetty/pull/23) ([Softwarekang](https://github.com/Softwarekang))
- feat\(example\): add example for session [\#22](https://github.com/Softwarekang/knetty/pull/22) ([Softwarekang](https://github.com/Softwarekang))
- refactor\(kNet\): refactor project Knet dir [\#21](https://github.com/Softwarekang/knetty/pull/21) ([Softwarekang](https://github.com/Softwarekang))
- feat\(session\): impl knet session func [\#20](https://github.com/Softwarekang/knetty/pull/20) ([Softwarekang](https://github.com/Softwarekang))
- fix\(connection\): fix connection write buffer bug [\#19](https://github.com/Softwarekang/knetty/pull/19) ([Softwarekang](https://github.com/Softwarekang))
- pref\(connection\):pref connection flush buffer、OnWrite func  [\#18](https://github.com/Softwarekang/knetty/pull/18) ([Softwarekang](https://github.com/Softwarekang))
- perf\(connection\): impl connection write 、pref connection example [\#16](https://github.com/Softwarekang/knetty/pull/16) ([Softwarekang](https://github.com/Softwarekang))
- build\(workflow\): add workflow action step codecov [\#15](https://github.com/Softwarekang/knetty/pull/15) ([Softwarekang](https://github.com/Softwarekang))
- build\(workflow\): add workflow action step codecov [\#12](https://github.com/Softwarekang/knetty/pull/12) ([Softwarekang](https://github.com/Softwarekang))
- Create knet.yml [\#11](https://github.com/Softwarekang/knetty/pull/11) ([Softwarekang](https://github.com/Softwarekang))
- chore\(connection\_reader\): delete unuse code [\#10](https://github.com/Softwarekang/knetty/pull/10) ([Softwarekang](https://github.com/Softwarekang))
- Update README.md [\#9](https://github.com/Softwarekang/knetty/pull/9) ([Softwarekang](https://github.com/Softwarekang))
- refactor\(knet\): refactor project [\#8](https://github.com/Softwarekang/knetty/pull/8) ([Softwarekang](https://github.com/Softwarekang))
- perf\(connection\): Implement the connection reader interface [\#7](https://github.com/Softwarekang/knetty/pull/7) ([Softwarekang](https://github.com/Softwarekang))
- perf\(connection\): perf connection refactor onRead func [\#6](https://github.com/Softwarekang/knetty/pull/6) ([Softwarekang](https://github.com/Softwarekang))
- fix\(example\): fix kqueue example server connection read bug [\#5](https://github.com/Softwarekang/knetty/pull/5) ([Softwarekang](https://github.com/Softwarekang))
- pref\(example\):pref kqueue connection read data [\#4](https://github.com/Softwarekang/knetty/pull/4) ([Softwarekang](https://github.com/Softwarekang))
- feat\(knet\): add kqueue impl for poller [\#3](https://github.com/Softwarekang/knetty/pull/3) ([Softwarekang](https://github.com/Softwarekang))
- feat\(knet\): Defining the poll Interface [\#2](https://github.com/Softwarekang/knetty/pull/2) ([Softwarekang](https://github.com/Softwarekang))
- chore\(gitignore\): ignore .idea directory [\#1](https://github.com/Softwarekang/knetty/pull/1) ([Softwarekang](https://github.com/Softwarekang))



\* *This Changelog was automatically generated by [github_changelog_generator](https://github.com/github-changelog-generator/github-changelog-generator)*
