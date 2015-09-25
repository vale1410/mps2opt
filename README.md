# mps2opt

An simple .mps MIPLIB converter to other formats. 

Currently supports conversion to minizinc format of MIPLIB mps files of type BP
and IP with integer coefficients. This converter is pre-alpha and works only under certain assumptions. 

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

* air04.mps
* ash608gpia-3col.mps
* eil33-2.mps
* eilA101-2.mps
* eilB101.mps
* ex1010-pi.mps
* ex10.mps
* ex9.mps
* f2000.mps
* go19.mps
* hanoi5.mps
* iis-100-0-cov.mps
* iis-bupa-cov.mps
* iis-pima-cov.mps
* m100n500k4r1.mps
* macrophage.mps
* methanosarcina.mps
* mine-166-5.mps
* mine-90-10.mps
* neos-1109824.mps
* neos-1337307.mps
* neos-1440225.mps
* neos-1616732.mps
* neos-1620770.mps
* neos18.mps
* neos-631710.mps
* neos-777800.mps
* neos-785912.mps
* neos-807456.mps
* neos-820146.mps
* neos-820157.mps
* neos-849702.mps
* neos858960.mps
* neos-952987.mps
* neos-957389.mps
* opm2-z10-s2.mps
* opm2-z11-s8.mps
* opm2-z12-s14.mps
* opm2-z12-s7.mps
* opm2-z7-s2.mps
* p6b.mps
* queens-30.mps
* rail01.mps
* ramos3.mps
* reblock166.mps
* reblock354.mps
* reblock420.mps
* reblock67.mps
* seymour.mps
* stp3d.mps
* sts405.mps
* sts729.mps
* t1717.mps
* t1722.mps
* tanglegram1.mps
* tanglegram2.mps
* toll-like.mps
* vpphard.mps
* 30n20b8.mps
* d10200.mps
* d20200.mps
* lectsched-1.mps
* lectsched-1-obj.mps
* lectsched-2.mps
* lectsched-3.mps
* lectsched-4-obj.mps
* neos16.mps
* neos-686190.mps
* ns1854840.mps
* rococoB10-011000.mps
* rococoC10-001000.mps
* rococoC11-011100.mps
* rococoC12-111000.mps

