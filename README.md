# mps2opt

An simple .mps MIPLIB converter to other formats. 

Currently supports conversion to minizinc and LP format of MIPLIB mps files of type BP
and IP with integer coefficients. This converter is pre-alpha and works only under these assumptions. 

To install download the golang distribution and use *go get*. 

```bash
go get git@github.com:vale1410/mps2opt.git
```

```bash
mps2opt -f air04.mps -minizinc
```

Please report any errors you encounter. 

Converted Instances
-------------------

It has converted successfully the following MIPLIB instances into minizinc: 

* 30n20b8
* air04
* ash608gpia-3col
* d10200
* d20200
* enlight13
* enlight14
* enlight15
* enlight16
* enlight9
* ex10
* ex1010-pi
* ex9
* f2000
* go19
* hanoi5
* iis-100-0-cov
* iis-bupa-cov
* iis-pima-cov
* lectsched-1
* lectsched-1-obj
* lectsched-2
* lectsched-3
* lectsched-4-obj
* m100n500k4r1
* macrophage
* methanosarcina
* neos-1109824
* neos-1337307
* neos-1440225
* neos16
* neos-1616732
* neos-1620770
* neos18
* neos-631710
* neos-686190
* neos-777800
* neos-785912
* neos-807456
* neos-820146
* neos-820157
* neos-849702
* neos858960
* ns1854840
* ns1952667
* nu120-pr3
* nu60-pr9
* opm2-z10-s2
* opm2-z11-s8
* opm2-z12-s14
* opm2-z12-s7
* opm2-z7-s2
* p6b
* pw-myciel4
* queens-30
* ramos3
* rococoB10-011000
* rococoC10-001000
* rococoC11-011100
* rococoC12-111000
* seymour
* sts405
* sts729
* t1717
* t1722
* tanglegram1
* tanglegram2
* toll-like
* tw-myciel4
* vpphard
