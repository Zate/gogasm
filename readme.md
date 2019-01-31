# GOGASM - Go Gaming Automated Server Manager

Automated server management for game servers.

Well that is what it might end up as one day, right now simply running ./gogasm
will query the live servers and display their pop.

Next up is a web interface to display all the info it grabs with A2S_Rules and A2S_Info

## Docker

```bash
docker build -t gogasm .
docker run gogasm -live napve
```

Output should look like:

```bash
2019/01/28 11:00:31.265895 main.go:420: A1 | Atlas_A1 - (v16.21) | 37.10.126.130:57555 Pop: 2
2019/01/28 11:00:31.858132 main.go:420: A2 | Atlas_A2 - (v16.21) | 37.10.126.130:57557 Pop: 2
2019/01/28 11:00:32.444956 main.go:420: A3 | Atlas_A3 - (v16.21) | 37.10.126.130:57559 Pop: 1
2019/01/28 11:00:33.036845 main.go:420: A4 | Atlas_A4 - (v16.21) | 37.10.126.130:57561 Pop: 8
2019/01/28 11:00:33.646473 main.go:420: A5 | Lawless Region - (v16.21) | 37.10.126.131:57555 Pop: 20
2019/01/28 11:00:34.243117 main.go:420: A6 | Northwest Tropical Freeport - (v16.21) | 37.10.126.131:57557 Pop: 17
2019/01/28 11:00:34.811976 main.go:420: A7 | Lawless Region - (v16.21) | 37.10.126.131:57559 Pop: 12
2019/01/28 11:00:35.387430 main.go:420: A8 | Atlas_A8 - (v16.21) | 37.10.126.131:57561 Pop: 14
2019/01/28 11:00:35.988003 main.go:420: A9 | Lawless Region - (v16.21) | 37.10.126.132:57555 Pop: 7
2019/01/28 11:00:36.595844 main.go:420: A10 | Southwest Tropical Freeport - (v16.21) | 37.10.126.132:57557 Pop: 13
2019/01/28 11:00:37.192902 main.go:420: A11 | Lawless Region - (v16.21) | 37.10.126.132:57559 Pop: 20
2019/01/28 11:00:37.797957 main.go:420: A12 | Atlas_A12 - (v16.21) | 37.10.126.132:57561 Pop: 5
2019/01/28 11:00:38.377579 main.go:420: A13 | Atlas_A13 - (v16.21) | 37.10.126.133:57555 Pop: 0
2019/01/28 11:00:38.971894 main.go:420: A14 | Atlas_A14 - (v16.21) | 37.10.126.133:57557 Pop: 1
2019/01/28 11:00:39.563928 main.go:420: A15 | Atlas_A15 - (v16.21) | 37.10.126.133:57559 Pop: 1
2019/01/28 11:00:40.167091 main.go:420: B1 | Atlas_B1 - (v16.21) | 37.10.126.133:57561 Pop: 0
2019/01/28 11:00:40.754423 main.go:420: B2 | Atlas_B2 - (v16.21) | 37.10.126.134:57555 Pop: 0
2019/01/28 11:00:41.321914 main.go:420: B3 | Atlas_B3 - (v16.21) | 37.10.126.134:57557 Pop: 2
2019/01/28 11:00:41.889598 main.go:420: B4 | Atlas_B4 - (v16.21) | 37.10.126.134:57559 Pop: 3
2019/01/28 11:00:42.472781 main.go:420: B5 | Atlas_B5 - (v16.21) | 37.10.126.134:57561 Pop: 6
```
