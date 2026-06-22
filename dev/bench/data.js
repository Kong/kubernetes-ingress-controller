window.BENCHMARK_DATA = {
  "lastUpdate": 1782114850983,
  "repoUrl": "https://github.com/Kong/kubernetes-ingress-controller",
  "entries": {
    "Go Benchmark": [
      {
        "commit": {
          "author": {
            "email": "zhangjintao9020@gmail.com",
            "name": "Jintao Zhang",
            "username": "tao12345666333"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": false,
          "id": "8bc40954943c87bf1c804a94f03894576e132446",
          "message": "sync CHANGELOG (#7800)\n\nupdate for KIC 3.4.9/3.4.10/3.5.2/3.5.3\n\nSigned-off-by: Jintao Zhang <zhangjintao9020@gmail.com>",
          "timestamp": "2025-11-26T19:35:30Z",
          "tree_id": "0d9ac78934d8a753879a99de978d9eed00f721a2",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/8bc40954943c87bf1c804a94f03894576e132446"
        },
        "date": 1764187560598,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 608445,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "1698 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 608445,
            "unit": "ns/op",
            "extra": "1698 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819453,
            "unit": "B/op",
            "extra": "1698 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1698 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7014,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "156936 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7014,
            "unit": "ns/op",
            "extra": "156936 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "156936 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "156936 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.81,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16434453 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.81,
            "unit": "ns/op",
            "extra": "16434453 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16434453 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16434453 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21811,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54483 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21811,
            "unit": "ns/op",
            "extra": "54483 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54483 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54483 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 228027,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5245 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 228027,
            "unit": "ns/op",
            "extra": "5245 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5245 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5245 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2417362,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "484 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2417362,
            "unit": "ns/op",
            "extra": "484 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "484 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "484 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 41070776,
            "unit": "ns/op\t24010754 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 41070776,
            "unit": "ns/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010754,
            "unit": "B/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 11839,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "92973 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 11839,
            "unit": "ns/op",
            "extra": "92973 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "92973 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "92973 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7555873,
            "unit": "ns/op\t 4527541 B/op\t   69230 allocs/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7555873,
            "unit": "ns/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527541,
            "unit": "B/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8154,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "145407 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8154,
            "unit": "ns/op",
            "extra": "145407 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "145407 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "145407 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 852830,
            "unit": "ns/op\t  392939 B/op\t    6173 allocs/op",
            "extra": "1370 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 852830,
            "unit": "ns/op",
            "extra": "1370 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392939,
            "unit": "B/op",
            "extra": "1370 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1370 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11156,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "109167 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11156,
            "unit": "ns/op",
            "extra": "109167 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "109167 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "109167 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8375422,
            "unit": "ns/op\t 4914206 B/op\t   75241 allocs/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8375422,
            "unit": "ns/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914206,
            "unit": "B/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 320.5,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3713754 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 320.5,
            "unit": "ns/op",
            "extra": "3713754 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3713754 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3713754 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.3,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 561.8,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2142327 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 561.8,
            "unit": "ns/op",
            "extra": "2142327 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2142327 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2142327 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 977.2,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1253302 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 977.2,
            "unit": "ns/op",
            "extra": "1253302 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1253302 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1253302 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "49699333+dependabot[bot]@users.noreply.github.com",
            "name": "dependabot[bot]",
            "username": "dependabot[bot]"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "6ec804e5f26d59623a06f00379bbb8c8b177c076",
          "message": "chore(deps): bump golang from 1.25.4 to 1.25.5 (#7806)\n\nBumps golang from 1.25.4 to 1.25.5.\n\n---\nupdated-dependencies:\n- dependency-name: golang\n  dependency-version: 1.25.5\n  dependency-type: direct:production\n  update-type: version-update:semver-patch\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2025-12-04T02:56:42Z",
          "tree_id": "0a99e82a74f2fdd1483070aea09ad2856ba496e2",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/6ec804e5f26d59623a06f00379bbb8c8b177c076"
        },
        "date": 1764818840961,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 684952,
            "unit": "ns/op\t  819454 B/op\t       5 allocs/op",
            "extra": "1599 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 684952,
            "unit": "ns/op",
            "extra": "1599 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819454,
            "unit": "B/op",
            "extra": "1599 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1599 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7501,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "161691 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7501,
            "unit": "ns/op",
            "extra": "161691 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "161691 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "161691 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.71,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16468134 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.71,
            "unit": "ns/op",
            "extra": "16468134 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16468134 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16468134 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21831,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53853 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21831,
            "unit": "ns/op",
            "extra": "53853 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53853 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53853 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 205722,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5012 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 205722,
            "unit": "ns/op",
            "extra": "5012 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5012 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5012 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2884773,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "412 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2884773,
            "unit": "ns/op",
            "extra": "412 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "412 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "412 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 30498926,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 30498926,
            "unit": "ns/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9610,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "126139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9610,
            "unit": "ns/op",
            "extra": "126139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "126139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "126139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7019245,
            "unit": "ns/op\t 4527544 B/op\t   69230 allocs/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7019245,
            "unit": "ns/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527544,
            "unit": "B/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7951,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "151461 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7951,
            "unit": "ns/op",
            "extra": "151461 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "151461 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "151461 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 841296,
            "unit": "ns/op\t  392952 B/op\t    6173 allocs/op",
            "extra": "1383 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 841296,
            "unit": "ns/op",
            "extra": "1383 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392952,
            "unit": "B/op",
            "extra": "1383 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1383 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11024,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "107455 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11024,
            "unit": "ns/op",
            "extra": "107455 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "107455 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "107455 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8044358,
            "unit": "ns/op\t 4914203 B/op\t   75241 allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8044358,
            "unit": "ns/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914203,
            "unit": "B/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 316.8,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3794570 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 316.8,
            "unit": "ns/op",
            "extra": "3794570 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3794570 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3794570 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.27,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.27,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 571.9,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2097601 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 571.9,
            "unit": "ns/op",
            "extra": "2097601 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2097601 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2097601 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 820.9,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1471269 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 820.9,
            "unit": "ns/op",
            "extra": "1471269 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1471269 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1471269 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "pmalek@users.noreply.github.com",
            "name": "Patryk Małek",
            "username": "pmalek"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "b9a88601090b70d2d1c4086a2312574e19734af0",
          "message": "ci: fix .github/workflows/validate_kong_image_trigger_via_label.yaml workflow (#7816)",
          "timestamp": "2025-12-16T15:58:46+01:00",
          "tree_id": "89e92124e00d9110ef9b4637bf75f7fa087dad30",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/b9a88601090b70d2d1c4086a2312574e19734af0"
        },
        "date": 1765897292712,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 866827,
            "unit": "ns/op\t  819450 B/op\t       5 allocs/op",
            "extra": "1407 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 866827,
            "unit": "ns/op",
            "extra": "1407 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819450,
            "unit": "B/op",
            "extra": "1407 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1407 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 9072,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "128734 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 9072,
            "unit": "ns/op",
            "extra": "128734 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "128734 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "128734 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 65.81,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18090970 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 65.81,
            "unit": "ns/op",
            "extra": "18090970 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "18090970 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "18090970 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21890,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "51277 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21890,
            "unit": "ns/op",
            "extra": "51277 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "51277 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "51277 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 242136,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4434 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 242136,
            "unit": "ns/op",
            "extra": "4434 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4434 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4434 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2750445,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "416 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2750445,
            "unit": "ns/op",
            "extra": "416 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "416 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "416 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 46511567,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "57 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 46511567,
            "unit": "ns/op",
            "extra": "57 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "57 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "57 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 10182,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "117421 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10182,
            "unit": "ns/op",
            "extra": "117421 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "117421 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "117421 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7685664,
            "unit": "ns/op\t 4527548 B/op\t   69230 allocs/op",
            "extra": "154 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7685664,
            "unit": "ns/op",
            "extra": "154 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527548,
            "unit": "B/op",
            "extra": "154 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "154 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8522,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "145436 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8522,
            "unit": "ns/op",
            "extra": "145436 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "145436 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "145436 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 960017,
            "unit": "ns/op\t  392947 B/op\t    6173 allocs/op",
            "extra": "1258 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 960017,
            "unit": "ns/op",
            "extra": "1258 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392947,
            "unit": "B/op",
            "extra": "1258 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1258 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 12025,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "96436 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 12025,
            "unit": "ns/op",
            "extra": "96436 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "96436 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "96436 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8926688,
            "unit": "ns/op\t 4914247 B/op\t   75241 allocs/op",
            "extra": "132 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8926688,
            "unit": "ns/op",
            "extra": "132 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914247,
            "unit": "B/op",
            "extra": "132 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "132 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 318.4,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3785394 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 318.4,
            "unit": "ns/op",
            "extra": "3785394 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3785394 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3785394 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 11.32,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 11.32,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 727.5,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "1639101 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 727.5,
            "unit": "ns/op",
            "extra": "1639101 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "1639101 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1639101 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 925,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1250644 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 925,
            "unit": "ns/op",
            "extra": "1250644 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1250644 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1250644 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "pmalek@users.noreply.github.com",
            "name": "Patryk Małek",
            "username": "pmalek"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "92b4ae93001cfec26df4b5a9a58008d30d568be8",
          "message": "ci: fix args passed to .github/workflows/validate_kong_image_trigger_via_label.yaml workflow (#7817)",
          "timestamp": "2025-12-16T16:11:10+01:00",
          "tree_id": "e7c277a6ac4f2b35a0c0acdf39b0107f5e490f3c",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/92b4ae93001cfec26df4b5a9a58008d30d568be8"
        },
        "date": 1765898041744,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 755547,
            "unit": "ns/op\t  819454 B/op\t       5 allocs/op",
            "extra": "1407 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 755547,
            "unit": "ns/op",
            "extra": "1407 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819454,
            "unit": "B/op",
            "extra": "1407 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1407 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7580,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "157725 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7580,
            "unit": "ns/op",
            "extra": "157725 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "157725 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "157725 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.64,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16472269 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.64,
            "unit": "ns/op",
            "extra": "16472269 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16472269 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16472269 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21940,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53631 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21940,
            "unit": "ns/op",
            "extra": "53631 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53631 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53631 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 212181,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4920 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 212181,
            "unit": "ns/op",
            "extra": "4920 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4920 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4920 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 3034528,
            "unit": "ns/op\t 2408450 B/op\t       2 allocs/op",
            "extra": "429 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 3034528,
            "unit": "ns/op",
            "extra": "429 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408450,
            "unit": "B/op",
            "extra": "429 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "429 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 35199134,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 35199134,
            "unit": "ns/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9670,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "123175 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9670,
            "unit": "ns/op",
            "extra": "123175 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "123175 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "123175 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7133474,
            "unit": "ns/op\t 4527541 B/op\t   69230 allocs/op",
            "extra": "168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7133474,
            "unit": "ns/op",
            "extra": "168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527541,
            "unit": "B/op",
            "extra": "168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8046,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "150622 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8046,
            "unit": "ns/op",
            "extra": "150622 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "150622 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "150622 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 853164,
            "unit": "ns/op\t  392935 B/op\t    6173 allocs/op",
            "extra": "1372 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 853164,
            "unit": "ns/op",
            "extra": "1372 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392935,
            "unit": "B/op",
            "extra": "1372 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1372 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11079,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "108241 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11079,
            "unit": "ns/op",
            "extra": "108241 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "108241 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "108241 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8121448,
            "unit": "ns/op\t 4914203 B/op\t   75241 allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8121448,
            "unit": "ns/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914203,
            "unit": "B/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 319.5,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3723254 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 319.5,
            "unit": "ns/op",
            "extra": "3723254 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3723254 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3723254 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.29,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.29,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 570.3,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2116950 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 570.3,
            "unit": "ns/op",
            "extra": "2116950 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2116950 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2116950 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 815.4,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1469856 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 815.4,
            "unit": "ns/op",
            "extra": "1469856 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1469856 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1469856 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "49699333+dependabot[bot]@users.noreply.github.com",
            "name": "dependabot[bot]",
            "username": "dependabot[bot]"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "e7f799f4a22a19d004a00a94dd75cc55fa387f97",
          "message": "chore(deps): bump golang from 1.25.5 to 1.25.6 (#7823)\n\nBumps golang from 1.25.5 to 1.25.6.\n\n---\nupdated-dependencies:\n- dependency-name: golang\n  dependency-version: 1.25.6\n  dependency-type: direct:production\n  update-type: version-update:semver-patch\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2026-02-05T10:24:42Z",
          "tree_id": "a9d002a892e44d3245df0088b9c5d7fd102c8bf9",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/e7f799f4a22a19d004a00a94dd75cc55fa387f97"
        },
        "date": 1770288929846,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 668132,
            "unit": "ns/op\t  819450 B/op\t       5 allocs/op",
            "extra": "1654 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 668132,
            "unit": "ns/op",
            "extra": "1654 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819450,
            "unit": "B/op",
            "extra": "1654 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1654 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 9183,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "127734 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 9183,
            "unit": "ns/op",
            "extra": "127734 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "127734 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "127734 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.94,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16357864 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.94,
            "unit": "ns/op",
            "extra": "16357864 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16357864 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16357864 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22020,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53866 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22020,
            "unit": "ns/op",
            "extra": "53866 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53866 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53866 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 223356,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4827 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 223356,
            "unit": "ns/op",
            "extra": "4827 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4827 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4827 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2377734,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "440 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2377734,
            "unit": "ns/op",
            "extra": "440 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "440 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "440 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 42329161,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "34 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 42329161,
            "unit": "ns/op",
            "extra": "34 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "34 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "34 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 11029,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "98932 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 11029,
            "unit": "ns/op",
            "extra": "98932 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "98932 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "98932 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7570336,
            "unit": "ns/op\t 4527544 B/op\t   69230 allocs/op",
            "extra": "156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7570336,
            "unit": "ns/op",
            "extra": "156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527544,
            "unit": "B/op",
            "extra": "156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8290,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "144706 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8290,
            "unit": "ns/op",
            "extra": "144706 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "144706 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "144706 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 874218,
            "unit": "ns/op\t  392941 B/op\t    6173 allocs/op",
            "extra": "1344 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 874218,
            "unit": "ns/op",
            "extra": "1344 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392941,
            "unit": "B/op",
            "extra": "1344 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1344 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11529,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "102535 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11529,
            "unit": "ns/op",
            "extra": "102535 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "102535 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "102535 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8720284,
            "unit": "ns/op\t 4914235 B/op\t   75241 allocs/op",
            "extra": "134 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8720284,
            "unit": "ns/op",
            "extra": "134 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914235,
            "unit": "B/op",
            "extra": "134 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "134 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 321.7,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3718606 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 321.7,
            "unit": "ns/op",
            "extra": "3718606 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3718606 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3718606 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.29,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.29,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 568.7,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2103313 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 568.7,
            "unit": "ns/op",
            "extra": "2103313 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2103313 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2103313 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 1062,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 1062,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "29139614+renovate[bot]@users.noreply.github.com",
            "name": "renovate[bot]",
            "username": "renovate[bot]"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "676716870aa41380ac19074b926e01eb943295ac",
          "message": "chore(deps): update dependency rhysd/actionlint to v1.7.10 (#7801)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2026-02-06T02:44:25Z",
          "tree_id": "da5e94043d91c6312130d5c712c63c87fb61d11b",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/676716870aa41380ac19074b926e01eb943295ac"
        },
        "date": 1770347679227,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 646830,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "1617 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 646830,
            "unit": "ns/op",
            "extra": "1617 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819453,
            "unit": "B/op",
            "extra": "1617 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1617 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7416,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "164223 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7416,
            "unit": "ns/op",
            "extra": "164223 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "164223 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "164223 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 74.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16419555 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 74.5,
            "unit": "ns/op",
            "extra": "16419555 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16419555 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16419555 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 25379,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "42949 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 25379,
            "unit": "ns/op",
            "extra": "42949 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "42949 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "42949 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 208353,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5319 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 208353,
            "unit": "ns/op",
            "extra": "5319 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5319 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5319 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2281025,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "549 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2281025,
            "unit": "ns/op",
            "extra": "549 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "549 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "549 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 36911351,
            "unit": "ns/op\t24010754 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 36911351,
            "unit": "ns/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010754,
            "unit": "B/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 11217,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "122655 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 11217,
            "unit": "ns/op",
            "extra": "122655 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "122655 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "122655 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7516054,
            "unit": "ns/op\t 4527541 B/op\t   69230 allocs/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7516054,
            "unit": "ns/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527541,
            "unit": "B/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7902,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "153914 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7902,
            "unit": "ns/op",
            "extra": "153914 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "153914 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "153914 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 844701,
            "unit": "ns/op\t  392942 B/op\t    6173 allocs/op",
            "extra": "1380 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 844701,
            "unit": "ns/op",
            "extra": "1380 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392942,
            "unit": "B/op",
            "extra": "1380 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1380 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11106,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "106718 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11106,
            "unit": "ns/op",
            "extra": "106718 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "106718 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "106718 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8072996,
            "unit": "ns/op\t 4914180 B/op\t   75241 allocs/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8072996,
            "unit": "ns/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914180,
            "unit": "B/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 326.8,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3682038 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 326.8,
            "unit": "ns/op",
            "extra": "3682038 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3682038 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3682038 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.49,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.49,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 570.9,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2100988 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 570.9,
            "unit": "ns/op",
            "extra": "2100988 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2100988 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2100988 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 873.2,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1458199 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 873.2,
            "unit": "ns/op",
            "extra": "1458199 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1458199 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1458199 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "29139614+renovate[bot]@users.noreply.github.com",
            "name": "renovate[bot]",
            "username": "renovate[bot]"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "76d5d18fab54477dc0d79f9ab10106052ff16de1",
          "message": "chore(deps): update kindest/node@only-patch docker tag to v1.31.14 (#7727)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2026-02-11T02:23:15Z",
          "tree_id": "8c8f3623e94d16462b8d8dbc93ce2432e4a8c5fd",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/76d5d18fab54477dc0d79f9ab10106052ff16de1"
        },
        "date": 1770778427671,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 705548,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "1545 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 705548,
            "unit": "ns/op",
            "extra": "1545 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819453,
            "unit": "B/op",
            "extra": "1545 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1545 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7691,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "161871 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7691,
            "unit": "ns/op",
            "extra": "161871 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "161871 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "161871 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 73.18,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16468917 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 73.18,
            "unit": "ns/op",
            "extra": "16468917 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16468917 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16468917 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21568,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54985 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21568,
            "unit": "ns/op",
            "extra": "54985 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54985 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54985 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 205488,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5112 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 205488,
            "unit": "ns/op",
            "extra": "5112 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5112 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5112 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2430920,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "498 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2430920,
            "unit": "ns/op",
            "extra": "498 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "498 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "498 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 41513596,
            "unit": "ns/op\t24010754 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 41513596,
            "unit": "ns/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010754,
            "unit": "B/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9317,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "127951 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9317,
            "unit": "ns/op",
            "extra": "127951 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "127951 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "127951 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6911186,
            "unit": "ns/op\t 4527542 B/op\t   69230 allocs/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6911186,
            "unit": "ns/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527542,
            "unit": "B/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7800,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "153025 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7800,
            "unit": "ns/op",
            "extra": "153025 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "153025 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "153025 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 832538,
            "unit": "ns/op\t  392934 B/op\t    6173 allocs/op",
            "extra": "1422 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 832538,
            "unit": "ns/op",
            "extra": "1422 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392934,
            "unit": "B/op",
            "extra": "1422 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1422 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10744,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "111036 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10744,
            "unit": "ns/op",
            "extra": "111036 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "111036 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "111036 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7894144,
            "unit": "ns/op\t 4914232 B/op\t   75241 allocs/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7894144,
            "unit": "ns/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914232,
            "unit": "B/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 319.3,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3771070 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 319.3,
            "unit": "ns/op",
            "extra": "3771070 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3771070 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3771070 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.39,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.39,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 659.2,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2073902 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 659.2,
            "unit": "ns/op",
            "extra": "2073902 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2073902 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2073902 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 844.9,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1395243 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 844.9,
            "unit": "ns/op",
            "extra": "1395243 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1395243 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1395243 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "zhangjintao9020@gmail.com",
            "name": "Jintao Zhang",
            "username": "tao12345666333"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": false,
          "id": "340fa9e920a7dece13e1341f93c6a7acf745207f",
          "message": "fix(konnect): separate config sync concurrency to reduce 429s (#7834)\n\n* feat(konnect): add config sync concurrency\n\nSigned-off-by: Jintao Zhang <zhangjintao9020@gmail.com>\n\n* add changelog entry for konnect concurrency\n\nSigned-off-by: Jintao Zhang <zhangjintao9020@gmail.com>\n\n* make generate\n\nSigned-off-by: Jintao Zhang <zhangjintao9020@gmail.com>\n\n* docs: reclassify konnect concurrency changelog\n\nSigned-off-by: Jintao Zhang <zhangjintao9020@gmail.com>\n\n---------\n\nSigned-off-by: Jintao Zhang <zhangjintao9020@gmail.com>",
          "timestamp": "2026-02-12T06:57:06Z",
          "tree_id": "b53949d0f3324f91869480f897357c1b8cbe2e6e",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/340fa9e920a7dece13e1341f93c6a7acf745207f"
        },
        "date": 1770881249173,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 676901,
            "unit": "ns/op\t  819454 B/op\t       5 allocs/op",
            "extra": "1528 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 676901,
            "unit": "ns/op",
            "extra": "1528 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819454,
            "unit": "B/op",
            "extra": "1528 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1528 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7542,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "151915 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7542,
            "unit": "ns/op",
            "extra": "151915 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "151915 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "151915 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.78,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16479596 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.78,
            "unit": "ns/op",
            "extra": "16479596 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16479596 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16479596 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21956,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53622 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21956,
            "unit": "ns/op",
            "extra": "53622 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53622 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53622 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 218224,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4917 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 218224,
            "unit": "ns/op",
            "extra": "4917 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4917 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4917 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2738722,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "518 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2738722,
            "unit": "ns/op",
            "extra": "518 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "518 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "518 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 32817930,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 32817930,
            "unit": "ns/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9719,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "123478 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9719,
            "unit": "ns/op",
            "extra": "123478 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "123478 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "123478 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7404507,
            "unit": "ns/op\t 4527543 B/op\t   69230 allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7404507,
            "unit": "ns/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527543,
            "unit": "B/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8055,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "150058 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8055,
            "unit": "ns/op",
            "extra": "150058 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "150058 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "150058 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 856439,
            "unit": "ns/op\t  392952 B/op\t    6173 allocs/op",
            "extra": "1384 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 856439,
            "unit": "ns/op",
            "extra": "1384 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392952,
            "unit": "B/op",
            "extra": "1384 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1384 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11281,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "105992 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11281,
            "unit": "ns/op",
            "extra": "105992 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "105992 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "105992 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8309693,
            "unit": "ns/op\t 4914309 B/op\t   75241 allocs/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8309693,
            "unit": "ns/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914309,
            "unit": "B/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 320.5,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3737593 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 320.5,
            "unit": "ns/op",
            "extra": "3737593 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3737593 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3737593 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.3,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 571.6,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2120774 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 571.6,
            "unit": "ns/op",
            "extra": "2120774 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2120774 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2120774 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 819,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1470456 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 819,
            "unit": "ns/op",
            "extra": "1470456 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1470456 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1470456 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "tao.yi@konghq.com",
            "name": "Tao Yi",
            "username": "randmonkey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "83286077deaef271ecf18434ccfef1f774d8a497",
          "message": "fix(konnect/tracing): Add values about Konnect sync round for tracing and logging in Konnect config syncer (#7838)\n\n* add values about Konnect sync round for tracing and logging\n\n* update changelog",
          "timestamp": "2026-02-16T14:36:10Z",
          "tree_id": "773119f5ee99b22e32310db5ef68ef14ef64000c",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/83286077deaef271ecf18434ccfef1f774d8a497"
        },
        "date": 1771254405778,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 685121,
            "unit": "ns/op\t  819449 B/op\t       5 allocs/op",
            "extra": "1606 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 685121,
            "unit": "ns/op",
            "extra": "1606 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819449,
            "unit": "B/op",
            "extra": "1606 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1606 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7123,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "148876 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7123,
            "unit": "ns/op",
            "extra": "148876 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "148876 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "148876 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.77,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16438549 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.77,
            "unit": "ns/op",
            "extra": "16438549 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16438549 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16438549 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21969,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54465 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21969,
            "unit": "ns/op",
            "extra": "54465 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54465 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54465 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 208354,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4993 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 208354,
            "unit": "ns/op",
            "extra": "4993 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4993 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4993 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2344683,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "501 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2344683,
            "unit": "ns/op",
            "extra": "501 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "501 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "501 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 29891387,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "75 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 29891387,
            "unit": "ns/op",
            "extra": "75 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "75 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "75 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9675,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "124024 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9675,
            "unit": "ns/op",
            "extra": "124024 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "124024 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "124024 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7018857,
            "unit": "ns/op\t 4527547 B/op\t   69230 allocs/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7018857,
            "unit": "ns/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527547,
            "unit": "B/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8026,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "150518 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8026,
            "unit": "ns/op",
            "extra": "150518 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "150518 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "150518 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 846969,
            "unit": "ns/op\t  392937 B/op\t    6173 allocs/op",
            "extra": "1375 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 846969,
            "unit": "ns/op",
            "extra": "1375 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392937,
            "unit": "B/op",
            "extra": "1375 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1375 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11055,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "108285 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11055,
            "unit": "ns/op",
            "extra": "108285 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "108285 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "108285 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8077306,
            "unit": "ns/op\t 4914208 B/op\t   75241 allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8077306,
            "unit": "ns/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914208,
            "unit": "B/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 321.9,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3724881 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 321.9,
            "unit": "ns/op",
            "extra": "3724881 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3724881 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3724881 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.32,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.32,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 583.4,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2059345 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 583.4,
            "unit": "ns/op",
            "extra": "2059345 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2059345 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2059345 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 846,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1414335 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 846,
            "unit": "ns/op",
            "extra": "1414335 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1414335 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1414335 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "49699333+dependabot[bot]@users.noreply.github.com",
            "name": "dependabot[bot]",
            "username": "dependabot[bot]"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "99b95e8dd31842cfd9df0265844ffff9aa77ff08",
          "message": "chore(deps): bump golang from 1.25.6 to 1.25.7 (#7830)\n\nBumps golang from 1.25.6 to 1.25.7.\n\n---\nupdated-dependencies:\n- dependency-name: golang\n  dependency-version: 1.25.7\n  dependency-type: direct:production\n  update-type: version-update:semver-patch\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2026-02-16T16:45:56+01:00",
          "tree_id": "95b1f40171ec27c3df419b7eeb7f41972785aa1f",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/99b95e8dd31842cfd9df0265844ffff9aa77ff08"
        },
        "date": 1771256918522,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 694724,
            "unit": "ns/op\t  819457 B/op\t       5 allocs/op",
            "extra": "1602 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 694724,
            "unit": "ns/op",
            "extra": "1602 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819457,
            "unit": "B/op",
            "extra": "1602 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1602 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7675,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "140557 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7675,
            "unit": "ns/op",
            "extra": "140557 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "140557 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "140557 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.77,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16450908 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.77,
            "unit": "ns/op",
            "extra": "16450908 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16450908 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16450908 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21559,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "55918 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21559,
            "unit": "ns/op",
            "extra": "55918 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "55918 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "55918 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 208302,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5046 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 208302,
            "unit": "ns/op",
            "extra": "5046 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5046 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5046 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2357208,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "444 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2357208,
            "unit": "ns/op",
            "extra": "444 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "444 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "444 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 31938015,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 31938015,
            "unit": "ns/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 12511,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "96631 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 12511,
            "unit": "ns/op",
            "extra": "96631 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "96631 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "96631 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7065297,
            "unit": "ns/op\t 4527544 B/op\t   69230 allocs/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7065297,
            "unit": "ns/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527544,
            "unit": "B/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7998,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "153226 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7998,
            "unit": "ns/op",
            "extra": "153226 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "153226 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "153226 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 850233,
            "unit": "ns/op\t  392934 B/op\t    6173 allocs/op",
            "extra": "1378 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 850233,
            "unit": "ns/op",
            "extra": "1378 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392934,
            "unit": "B/op",
            "extra": "1378 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1378 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11103,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "109274 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11103,
            "unit": "ns/op",
            "extra": "109274 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "109274 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "109274 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8182737,
            "unit": "ns/op\t 4914212 B/op\t   75241 allocs/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8182737,
            "unit": "ns/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914212,
            "unit": "B/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 325.5,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3744342 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 325.5,
            "unit": "ns/op",
            "extra": "3744342 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3744342 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3744342 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.5,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 580.2,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2067358 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 580.2,
            "unit": "ns/op",
            "extra": "2067358 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2067358 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2067358 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 995.5,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1223416 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 995.5,
            "unit": "ns/op",
            "extra": "1223416 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1223416 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1223416 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "pmalek@users.noreply.github.com",
            "name": "Patryk Małek",
            "username": "pmalek"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "38dbe0d3d5af3a40557be1feffcce14d398da725",
          "message": "chore: add changelog entries for 3.5.4 and 3.4.11 (#7847)",
          "timestamp": "2026-02-17T12:01:54+01:00",
          "tree_id": "3545802dafc65d689c759a98b229fe3e809c0db8",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/38dbe0d3d5af3a40557be1feffcce14d398da725"
        },
        "date": 1771326280450,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 559608,
            "unit": "ns/op\t  819450 B/op\t       5 allocs/op",
            "extra": "1864 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 559608,
            "unit": "ns/op",
            "extra": "1864 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819450,
            "unit": "B/op",
            "extra": "1864 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1864 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7120,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "170742 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7120,
            "unit": "ns/op",
            "extra": "170742 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "170742 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "170742 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 92.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12869629 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 92.9,
            "unit": "ns/op",
            "extra": "12869629 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12869629 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12869629 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21880,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54746 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21880,
            "unit": "ns/op",
            "extra": "54746 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54746 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54746 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 206867,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5091 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 206867,
            "unit": "ns/op",
            "extra": "5091 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5091 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5091 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2292352,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "505 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2292352,
            "unit": "ns/op",
            "extra": "505 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "505 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "505 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 30400647,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 30400647,
            "unit": "ns/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9706,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "122491 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9706,
            "unit": "ns/op",
            "extra": "122491 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "122491 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "122491 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7008599,
            "unit": "ns/op\t 4527536 B/op\t   69230 allocs/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7008599,
            "unit": "ns/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527536,
            "unit": "B/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8172,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "147553 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8172,
            "unit": "ns/op",
            "extra": "147553 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "147553 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "147553 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 849287,
            "unit": "ns/op\t  392931 B/op\t    6173 allocs/op",
            "extra": "1399 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 849287,
            "unit": "ns/op",
            "extra": "1399 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392931,
            "unit": "B/op",
            "extra": "1399 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1399 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11133,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "108588 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11133,
            "unit": "ns/op",
            "extra": "108588 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "108588 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "108588 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8134355,
            "unit": "ns/op\t 4914185 B/op\t   75241 allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8134355,
            "unit": "ns/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914185,
            "unit": "B/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 325.3,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3733461 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 325.3,
            "unit": "ns/op",
            "extra": "3733461 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3733461 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3733461 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.52,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.52,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 577.7,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2078300 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 577.7,
            "unit": "ns/op",
            "extra": "2078300 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2078300 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2078300 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 832.9,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1437765 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 832.9,
            "unit": "ns/op",
            "extra": "1437765 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1437765 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1437765 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "29139614+renovate[bot]@users.noreply.github.com",
            "name": "renovate[bot]",
            "username": "renovate[bot]"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "c0b43159af0b1d6c38a1cbaafec20f6184edc9b2",
          "message": "chore(deps): update dependency gke to v1.35.1 (#7832)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2026-03-03T02:34:04Z",
          "tree_id": "eafa12a780f6214c46e537f76835d1640e5bb710",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/c0b43159af0b1d6c38a1cbaafec20f6184edc9b2"
        },
        "date": 1772507032943,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 715650,
            "unit": "ns/op\t  819450 B/op\t       5 allocs/op",
            "extra": "1470 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 715650,
            "unit": "ns/op",
            "extra": "1470 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819450,
            "unit": "B/op",
            "extra": "1470 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1470 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7557,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "157688 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7557,
            "unit": "ns/op",
            "extra": "157688 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "157688 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "157688 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 65.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18221300 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 65.8,
            "unit": "ns/op",
            "extra": "18221300 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "18221300 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "18221300 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 20926,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "56574 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 20926,
            "unit": "ns/op",
            "extra": "56574 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "56574 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "56574 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 223792,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4603 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 223792,
            "unit": "ns/op",
            "extra": "4603 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4603 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4603 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 3531505,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "382 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 3531505,
            "unit": "ns/op",
            "extra": "382 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "382 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "382 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 44370282,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "50 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 44370282,
            "unit": "ns/op",
            "extra": "50 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "50 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "50 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9828,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "123210 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9828,
            "unit": "ns/op",
            "extra": "123210 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "123210 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "123210 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7178312,
            "unit": "ns/op\t 4527547 B/op\t   69230 allocs/op",
            "extra": "165 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7178312,
            "unit": "ns/op",
            "extra": "165 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527547,
            "unit": "B/op",
            "extra": "165 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "165 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8110,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "146768 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8110,
            "unit": "ns/op",
            "extra": "146768 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "146768 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "146768 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 886280,
            "unit": "ns/op\t  392928 B/op\t    6173 allocs/op",
            "extra": "1359 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 886280,
            "unit": "ns/op",
            "extra": "1359 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392928,
            "unit": "B/op",
            "extra": "1359 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1359 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11431,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "103180 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11431,
            "unit": "ns/op",
            "extra": "103180 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "103180 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "103180 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8260774,
            "unit": "ns/op\t 4914180 B/op\t   75241 allocs/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8260774,
            "unit": "ns/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914180,
            "unit": "B/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 314.4,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3809245 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 314.4,
            "unit": "ns/op",
            "extra": "3809245 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3809245 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3809245 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 11.06,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 11.06,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 566.7,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2106571 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 566.7,
            "unit": "ns/op",
            "extra": "2106571 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2106571 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2106571 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 787.1,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1521294 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 787.1,
            "unit": "ns/op",
            "extra": "1521294 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1521294 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1521294 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "justin.davies@konghq.com",
            "name": "Justin Davies",
            "username": "justindavies"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "2651cc1044a2e41127b3f201ad073aa4efb6d9e3",
          "message": "Document KIC and Kong Operator lifecycle (#7850)\n\nAdded a section about KIC and Kong Operator to the README.",
          "timestamp": "2026-03-11T15:01:38+01:00",
          "tree_id": "aaa85a4e4d3a65daec283cbbdcbfad29c7c981e9",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/2651cc1044a2e41127b3f201ad073aa4efb6d9e3"
        },
        "date": 1773237877286,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 689056,
            "unit": "ns/op\t  819449 B/op\t       5 allocs/op",
            "extra": "1629 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 689056,
            "unit": "ns/op",
            "extra": "1629 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819449,
            "unit": "B/op",
            "extra": "1629 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1629 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7133,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "167538 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7133,
            "unit": "ns/op",
            "extra": "167538 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "167538 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "167538 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.83,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16433013 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.83,
            "unit": "ns/op",
            "extra": "16433013 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16433013 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16433013 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22256,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54399 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22256,
            "unit": "ns/op",
            "extra": "54399 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54399 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54399 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 228321,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4724 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 228321,
            "unit": "ns/op",
            "extra": "4724 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4724 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4724 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2556065,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "502 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2556065,
            "unit": "ns/op",
            "extra": "502 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "502 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "502 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 41294626,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 41294626,
            "unit": "ns/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 10002,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "117840 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10002,
            "unit": "ns/op",
            "extra": "117840 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "117840 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "117840 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7541438,
            "unit": "ns/op\t 4527537 B/op\t   69230 allocs/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7541438,
            "unit": "ns/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527537,
            "unit": "B/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8356,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "144175 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8356,
            "unit": "ns/op",
            "extra": "144175 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "144175 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "144175 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 888004,
            "unit": "ns/op\t  392933 B/op\t    6173 allocs/op",
            "extra": "1344 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 888004,
            "unit": "ns/op",
            "extra": "1344 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392933,
            "unit": "B/op",
            "extra": "1344 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1344 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11706,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "105009 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11706,
            "unit": "ns/op",
            "extra": "105009 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "105009 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "105009 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 9021891,
            "unit": "ns/op\t 4914211 B/op\t   75241 allocs/op",
            "extra": "132 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 9021891,
            "unit": "ns/op",
            "extra": "132 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914211,
            "unit": "B/op",
            "extra": "132 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "132 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 326.7,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3666798 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 326.7,
            "unit": "ns/op",
            "extra": "3666798 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3666798 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3666798 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.53,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.53,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 592.3,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2025296 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 592.3,
            "unit": "ns/op",
            "extra": "2025296 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2025296 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2025296 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 840.7,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1428721 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 840.7,
            "unit": "ns/op",
            "extra": "1428721 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1428721 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1428721 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "pmalek@users.noreply.github.com",
            "name": "Patryk Małek",
            "username": "pmalek"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "18c142392544be72d8c9a8da4ec4f403dbace62c",
          "message": "fix: fix generating SNIs in dbless (#7853)\n\n* fix: fix generating SNIs in dbless\n\n* tests: fix golden tests files\n\n* tests: fix golden tests files\n\n* chore: add changelog",
          "timestamp": "2026-03-11T16:00:29+01:00",
          "tree_id": "4039dbc61e197b1bead3c53f2d7ec3f2cf642667",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/18c142392544be72d8c9a8da4ec4f403dbace62c"
        },
        "date": 1773241412473,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 663835,
            "unit": "ns/op\t  819454 B/op\t       5 allocs/op",
            "extra": "1566 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 663835,
            "unit": "ns/op",
            "extra": "1566 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819454,
            "unit": "B/op",
            "extra": "1566 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1566 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7589,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "160114 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7589,
            "unit": "ns/op",
            "extra": "160114 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "160114 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "160114 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.78,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16502133 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.78,
            "unit": "ns/op",
            "extra": "16502133 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16502133 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16502133 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21571,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "55599 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21571,
            "unit": "ns/op",
            "extra": "55599 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "55599 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "55599 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 204902,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5178 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 204902,
            "unit": "ns/op",
            "extra": "5178 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5178 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5178 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2464751,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "424 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2464751,
            "unit": "ns/op",
            "extra": "424 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "424 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "424 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 31152998,
            "unit": "ns/op\t24010754 B/op\t       2 allocs/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 31152998,
            "unit": "ns/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010754,
            "unit": "B/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9704,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "125233 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9704,
            "unit": "ns/op",
            "extra": "125233 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "125233 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "125233 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7060034,
            "unit": "ns/op\t 4527537 B/op\t   69230 allocs/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7060034,
            "unit": "ns/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527537,
            "unit": "B/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8061,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "147574 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8061,
            "unit": "ns/op",
            "extra": "147574 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "147574 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "147574 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 855001,
            "unit": "ns/op\t  392943 B/op\t    6173 allocs/op",
            "extra": "1392 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 855001,
            "unit": "ns/op",
            "extra": "1392 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392943,
            "unit": "B/op",
            "extra": "1392 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1392 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11140,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "108960 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11140,
            "unit": "ns/op",
            "extra": "108960 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "108960 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "108960 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8148659,
            "unit": "ns/op\t 4914218 B/op\t   75241 allocs/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8148659,
            "unit": "ns/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914218,
            "unit": "B/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 323.4,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3693282 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 323.4,
            "unit": "ns/op",
            "extra": "3693282 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3693282 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3693282 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.33,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.33,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 584.2,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2076608 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 584.2,
            "unit": "ns/op",
            "extra": "2076608 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2076608 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2076608 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 826.5,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1452193 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 826.5,
            "unit": "ns/op",
            "extra": "1452193 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1452193 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1452193 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "tao.yi@konghq.com",
            "name": "Tao Yi",
            "username": "randmonkey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "6d6b65eb2f839c6b0a5b39ece8ad9ff64995e448",
          "message": "Fix: set protocols in translated expression based routes (#7860)\n\n* set protocols in translated expression based routes\n\n* update golden tests",
          "timestamp": "2026-03-16T06:09:04Z",
          "tree_id": "6d0b36e5d5fda2ca24d5bad6aeac9b368b9af765",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/6d6b65eb2f839c6b0a5b39ece8ad9ff64995e448"
        },
        "date": 1773643120149,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 641359,
            "unit": "ns/op\t  819454 B/op\t       5 allocs/op",
            "extra": "1736 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 641359,
            "unit": "ns/op",
            "extra": "1736 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819454,
            "unit": "B/op",
            "extra": "1736 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1736 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7707,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "153850 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7707,
            "unit": "ns/op",
            "extra": "153850 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "153850 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "153850 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.87,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16479348 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.87,
            "unit": "ns/op",
            "extra": "16479348 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16479348 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16479348 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 26010,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "41828 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 26010,
            "unit": "ns/op",
            "extra": "41828 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "41828 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "41828 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 212358,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4975 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 212358,
            "unit": "ns/op",
            "extra": "4975 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4975 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4975 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2341871,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "451 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2341871,
            "unit": "ns/op",
            "extra": "451 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "451 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "451 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 31151402,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 31151402,
            "unit": "ns/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9571,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "124734 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9571,
            "unit": "ns/op",
            "extra": "124734 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "124734 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "124734 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 8680195,
            "unit": "ns/op\t 4527542 B/op\t   69230 allocs/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 8680195,
            "unit": "ns/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527542,
            "unit": "B/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8089,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "128264 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8089,
            "unit": "ns/op",
            "extra": "128264 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "128264 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "128264 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 850915,
            "unit": "ns/op\t  392958 B/op\t    6173 allocs/op",
            "extra": "1400 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 850915,
            "unit": "ns/op",
            "extra": "1400 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392958,
            "unit": "B/op",
            "extra": "1400 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1400 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11054,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "106012 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11054,
            "unit": "ns/op",
            "extra": "106012 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "106012 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "106012 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8286351,
            "unit": "ns/op\t 4914234 B/op\t   75241 allocs/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8286351,
            "unit": "ns/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914234,
            "unit": "B/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 322.3,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3731424 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 322.3,
            "unit": "ns/op",
            "extra": "3731424 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3731424 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3731424 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.4,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 591.7,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2035178 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 591.7,
            "unit": "ns/op",
            "extra": "2035178 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2035178 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2035178 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 856.1,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1405597 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 856.1,
            "unit": "ns/op",
            "extra": "1405597 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1405597 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1405597 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "pmalek@users.noreply.github.com",
            "name": "Patryk Małek",
            "username": "pmalek"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "a5878aead35832b692e9b71b3c906ff954bfeb70",
          "message": "Revert \"fix: fix generating SNIs in dbless (#7853)\" (#7872)\n\nThis reverts commit 18c142392544be72d8c9a8da4ec4f403dbace62c.",
          "timestamp": "2026-03-31T15:26:47+02:00",
          "tree_id": "c60591c69957a0c725f7e0421ea42d260d94c988",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/a5878aead35832b692e9b71b3c906ff954bfeb70"
        },
        "date": 1774963776413,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 641294,
            "unit": "ns/op\t  819454 B/op\t       5 allocs/op",
            "extra": "1640 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 641294,
            "unit": "ns/op",
            "extra": "1640 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819454,
            "unit": "B/op",
            "extra": "1640 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1640 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7658,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "159668 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7658,
            "unit": "ns/op",
            "extra": "159668 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "159668 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "159668 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.83,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16370313 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.83,
            "unit": "ns/op",
            "extra": "16370313 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16370313 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16370313 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22186,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53776 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22186,
            "unit": "ns/op",
            "extra": "53776 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53776 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53776 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 214113,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5142 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 214113,
            "unit": "ns/op",
            "extra": "5142 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5142 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5142 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2317416,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "510 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2317416,
            "unit": "ns/op",
            "extra": "510 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "510 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "510 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 30421159,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 30421159,
            "unit": "ns/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9825,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "118896 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9825,
            "unit": "ns/op",
            "extra": "118896 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "118896 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "118896 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7140625,
            "unit": "ns/op\t 4527542 B/op\t   69230 allocs/op",
            "extra": "166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7140625,
            "unit": "ns/op",
            "extra": "166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527542,
            "unit": "B/op",
            "extra": "166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8118,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "148830 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8118,
            "unit": "ns/op",
            "extra": "148830 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "148830 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "148830 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 856089,
            "unit": "ns/op\t  392948 B/op\t    6173 allocs/op",
            "extra": "1348 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 856089,
            "unit": "ns/op",
            "extra": "1348 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392948,
            "unit": "B/op",
            "extra": "1348 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1348 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11226,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "108432 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11226,
            "unit": "ns/op",
            "extra": "108432 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "108432 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "108432 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8269949,
            "unit": "ns/op\t 4914152 B/op\t   75241 allocs/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8269949,
            "unit": "ns/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914152,
            "unit": "B/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 323.3,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3695716 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 323.3,
            "unit": "ns/op",
            "extra": "3695716 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3695716 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3695716 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.55,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.55,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 590.8,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2041844 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 590.8,
            "unit": "ns/op",
            "extra": "2041844 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2041844 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2041844 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 858.4,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1395471 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 858.4,
            "unit": "ns/op",
            "extra": "1395471 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1395471 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1395471 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "pmalek@users.noreply.github.com",
            "name": "Patryk Małek",
            "username": "pmalek"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "bc5e417dc2926cc82cee11add4e2dd05ab5d78a2",
          "message": "chore: update changelog (#7881)",
          "timestamp": "2026-04-01T10:15:57+02:00",
          "tree_id": "71e0c0237845786516f6ee6641309d1e2b7f677d",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/bc5e417dc2926cc82cee11add4e2dd05ab5d78a2"
        },
        "date": 1775031519692,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 687952,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "1615 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 687952,
            "unit": "ns/op",
            "extra": "1615 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819453,
            "unit": "B/op",
            "extra": "1615 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1615 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7124,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "147076 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7124,
            "unit": "ns/op",
            "extra": "147076 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "147076 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "147076 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.89,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16466817 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.89,
            "unit": "ns/op",
            "extra": "16466817 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16466817 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16466817 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22092,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53730 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22092,
            "unit": "ns/op",
            "extra": "53730 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53730 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53730 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 211903,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5025 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 211903,
            "unit": "ns/op",
            "extra": "5025 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5025 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5025 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2348827,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "518 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2348827,
            "unit": "ns/op",
            "extra": "518 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "518 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "518 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 30653639,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 30653639,
            "unit": "ns/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 10392,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "98746 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10392,
            "unit": "ns/op",
            "extra": "98746 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "98746 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "98746 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7074478,
            "unit": "ns/op\t 4527539 B/op\t   69230 allocs/op",
            "extra": "168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7074478,
            "unit": "ns/op",
            "extra": "168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527539,
            "unit": "B/op",
            "extra": "168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8013,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "149704 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8013,
            "unit": "ns/op",
            "extra": "149704 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "149704 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "149704 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 845966,
            "unit": "ns/op\t  392931 B/op\t    6173 allocs/op",
            "extra": "1402 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 845966,
            "unit": "ns/op",
            "extra": "1402 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392931,
            "unit": "B/op",
            "extra": "1402 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1402 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11027,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "108726 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11027,
            "unit": "ns/op",
            "extra": "108726 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "108726 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "108726 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8135312,
            "unit": "ns/op\t 4914220 B/op\t   75241 allocs/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8135312,
            "unit": "ns/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914220,
            "unit": "B/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 330.6,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3460519 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 330.6,
            "unit": "ns/op",
            "extra": "3460519 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3460519 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3460519 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.4,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 603.9,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2031003 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 603.9,
            "unit": "ns/op",
            "extra": "2031003 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2031003 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2031003 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 1101,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 1101,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "sandro.modarelli@gmail.com",
            "name": "Sandro Modarelli",
            "username": "sandromodarelli"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "03482d653ee08a914a85fdbd3a6ce439422ad4f5",
          "message": "chore(deps): update GKE test dependency to v1.35.3 (#7884)",
          "timestamp": "2026-04-10T08:29:51Z",
          "tree_id": "6c56ebb3945dfc5682df14d978a47b11e6f482c6",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/03482d653ee08a914a85fdbd3a6ce439422ad4f5"
        },
        "date": 1775811646029,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 714878,
            "unit": "ns/op\t  819455 B/op\t       5 allocs/op",
            "extra": "1446 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 714878,
            "unit": "ns/op",
            "extra": "1446 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819455,
            "unit": "B/op",
            "extra": "1446 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1446 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7724,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "132022 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7724,
            "unit": "ns/op",
            "extra": "132022 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "132022 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "132022 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.83,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16499550 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.83,
            "unit": "ns/op",
            "extra": "16499550 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16499550 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16499550 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22271,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53734 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22271,
            "unit": "ns/op",
            "extra": "53734 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53734 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53734 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 211431,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5022 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 211431,
            "unit": "ns/op",
            "extra": "5022 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5022 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5022 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2712289,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "506 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2712289,
            "unit": "ns/op",
            "extra": "506 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "506 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "506 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 30966893,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 30966893,
            "unit": "ns/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9815,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "123693 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9815,
            "unit": "ns/op",
            "extra": "123693 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "123693 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "123693 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7224073,
            "unit": "ns/op\t 4527540 B/op\t   69230 allocs/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7224073,
            "unit": "ns/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527540,
            "unit": "B/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8227,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "150621 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8227,
            "unit": "ns/op",
            "extra": "150621 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "150621 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "150621 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 859713,
            "unit": "ns/op\t  392938 B/op\t    6173 allocs/op",
            "extra": "1392 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 859713,
            "unit": "ns/op",
            "extra": "1392 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392938,
            "unit": "B/op",
            "extra": "1392 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1392 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11295,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "107238 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11295,
            "unit": "ns/op",
            "extra": "107238 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "107238 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "107238 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8224663,
            "unit": "ns/op\t 4914233 B/op\t   75241 allocs/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8224663,
            "unit": "ns/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914233,
            "unit": "B/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 324.9,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3659769 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 324.9,
            "unit": "ns/op",
            "extra": "3659769 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3659769 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3659769 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.59,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.59,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 606.9,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "1985859 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 606.9,
            "unit": "ns/op",
            "extra": "1985859 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "1985859 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1985859 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 873.3,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1377117 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 873.3,
            "unit": "ns/op",
            "extra": "1377117 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1377117 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1377117 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "tao.yi@konghq.com",
            "name": "Tao Yi",
            "username": "randmonkey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "e25003e3d4b419fca13df5d33d9aeb869178ee84",
          "message": "bump source version of KIC upgrade test (#7885)",
          "timestamp": "2026-04-13T09:11:20Z",
          "tree_id": "371018d0db26a9f6867d723b57483bef144d6120",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/e25003e3d4b419fca13df5d33d9aeb869178ee84"
        },
        "date": 1776073283437,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 626308,
            "unit": "ns/op\t  819452 B/op\t       5 allocs/op",
            "extra": "1809 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 626308,
            "unit": "ns/op",
            "extra": "1809 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819452,
            "unit": "B/op",
            "extra": "1809 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1809 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7942,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "132638 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7942,
            "unit": "ns/op",
            "extra": "132638 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "132638 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "132638 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.95,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16469548 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.95,
            "unit": "ns/op",
            "extra": "16469548 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16469548 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16469548 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22390,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53743 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22390,
            "unit": "ns/op",
            "extra": "53743 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53743 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53743 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 220972,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4796 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 220972,
            "unit": "ns/op",
            "extra": "4796 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4796 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4796 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2310148,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "508 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2310148,
            "unit": "ns/op",
            "extra": "508 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "508 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "508 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 30505276,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "54 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 30505276,
            "unit": "ns/op",
            "extra": "54 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
            "unit": "B/op",
            "extra": "54 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9683,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "123668 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9683,
            "unit": "ns/op",
            "extra": "123668 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "123668 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "123668 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7191146,
            "unit": "ns/op\t 4527548 B/op\t   69230 allocs/op",
            "extra": "165 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7191146,
            "unit": "ns/op",
            "extra": "165 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527548,
            "unit": "B/op",
            "extra": "165 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "165 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8235,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "149859 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8235,
            "unit": "ns/op",
            "extra": "149859 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "149859 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "149859 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 861146,
            "unit": "ns/op\t  392938 B/op\t    6173 allocs/op",
            "extra": "1353 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 861146,
            "unit": "ns/op",
            "extra": "1353 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392938,
            "unit": "B/op",
            "extra": "1353 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1353 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11278,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "107329 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11278,
            "unit": "ns/op",
            "extra": "107329 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "107329 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "107329 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8531597,
            "unit": "ns/op\t 4914154 B/op\t   75241 allocs/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8531597,
            "unit": "ns/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914154,
            "unit": "B/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 328.3,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3679327 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 328.3,
            "unit": "ns/op",
            "extra": "3679327 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3679327 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3679327 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.3,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 617.8,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "1924778 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 617.8,
            "unit": "ns/op",
            "extra": "1924778 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "1924778 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1924778 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 886.5,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1351315 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 886.5,
            "unit": "ns/op",
            "extra": "1351315 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1351315 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1351315 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "jakub.warczarek@konghq.com",
            "name": "Jakub Warczarek",
            "username": "programmer04"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "c188b25990da2b3c7250126864e2ae70950411a6",
          "message": "chore(deps): bump `kubernetes-testing-framework` to `v0.47.2-patch.3` (#7895)\n\n* chore(deps): bump `kubernetes-testing-framework` to `v0.47.2-patch.1`\n\n* fix\n\n* chore(deps): bump `kubernetes-testing-framework` to `v0.47.2-patch.3`",
          "timestamp": "2026-04-17T11:07:05Z",
          "tree_id": "ea4b8dda6a2eb04b95370d2624c0e5a21ce85834",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/c188b25990da2b3c7250126864e2ae70950411a6"
        },
        "date": 1776425934518,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 643534,
            "unit": "ns/op\t  819454 B/op\t       5 allocs/op",
            "extra": "1644 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 643534,
            "unit": "ns/op",
            "extra": "1644 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819454,
            "unit": "B/op",
            "extra": "1644 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1644 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7238,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "159256 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7238,
            "unit": "ns/op",
            "extra": "159256 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "159256 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "159256 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.93,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16424997 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.93,
            "unit": "ns/op",
            "extra": "16424997 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16424997 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16424997 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22346,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52954 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22346,
            "unit": "ns/op",
            "extra": "52954 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52954 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52954 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 208336,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5623 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 208336,
            "unit": "ns/op",
            "extra": "5623 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5623 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5623 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2327360,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "504 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2327360,
            "unit": "ns/op",
            "extra": "504 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "504 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "504 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 30454267,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "75 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 30454267,
            "unit": "ns/op",
            "extra": "75 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "75 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "75 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9697,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "123543 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9697,
            "unit": "ns/op",
            "extra": "123543 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "123543 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "123543 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7069998,
            "unit": "ns/op\t 4527538 B/op\t   69230 allocs/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7069998,
            "unit": "ns/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527538,
            "unit": "B/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8115,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "147787 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8115,
            "unit": "ns/op",
            "extra": "147787 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "147787 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "147787 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 846489,
            "unit": "ns/op\t  392930 B/op\t    6173 allocs/op",
            "extra": "1377 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 846489,
            "unit": "ns/op",
            "extra": "1377 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392930,
            "unit": "B/op",
            "extra": "1377 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1377 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11085,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "110019 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11085,
            "unit": "ns/op",
            "extra": "110019 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "110019 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "110019 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8128061,
            "unit": "ns/op\t 4914225 B/op\t   75241 allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8128061,
            "unit": "ns/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914225,
            "unit": "B/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 332.4,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3490921 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 332.4,
            "unit": "ns/op",
            "extra": "3490921 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3490921 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3490921 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.54,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.54,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 611.7,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "1972453 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 611.7,
            "unit": "ns/op",
            "extra": "1972453 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "1972453 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1972453 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 878.1,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1371996 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 878.1,
            "unit": "ns/op",
            "extra": "1371996 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1371996 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1371996 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "jakub.warczarek@konghq.com",
            "name": "Jakub Warczarek",
            "username": "programmer04"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "effe78af48b8c5725bfe38602e63223a94adef1d",
          "message": "Revert \"Fix: set protocols in translated expression based routes (#7860)\" (#7897)\n\nThis reverts commit 6d6b65eb2f839c6b0a5b39ece8ad9ff64995e448.",
          "timestamp": "2026-04-17T11:55:17Z",
          "tree_id": "b7db1c265cbeec16ea61d8b72fc8be0f1c0fd611",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/effe78af48b8c5725bfe38602e63223a94adef1d"
        },
        "date": 1776428835822,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 651355,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "1789 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 651355,
            "unit": "ns/op",
            "extra": "1789 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819453,
            "unit": "B/op",
            "extra": "1789 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1789 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7802,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "155251 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7802,
            "unit": "ns/op",
            "extra": "155251 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "155251 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "155251 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.87,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16471267 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.87,
            "unit": "ns/op",
            "extra": "16471267 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16471267 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16471267 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22906,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52431 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22906,
            "unit": "ns/op",
            "extra": "52431 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52431 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52431 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 226197,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4779 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 226197,
            "unit": "ns/op",
            "extra": "4779 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4779 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4779 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2342975,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "501 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2342975,
            "unit": "ns/op",
            "extra": "501 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "501 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "501 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 42059704,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 42059704,
            "unit": "ns/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 10047,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "117048 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10047,
            "unit": "ns/op",
            "extra": "117048 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "117048 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "117048 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7520063,
            "unit": "ns/op\t 4527536 B/op\t   69230 allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7520063,
            "unit": "ns/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527536,
            "unit": "B/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8358,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "141514 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8358,
            "unit": "ns/op",
            "extra": "141514 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "141514 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "141514 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 881089,
            "unit": "ns/op\t  392929 B/op\t    6173 allocs/op",
            "extra": "1327 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 881089,
            "unit": "ns/op",
            "extra": "1327 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392929,
            "unit": "B/op",
            "extra": "1327 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1327 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11822,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "102087 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11822,
            "unit": "ns/op",
            "extra": "102087 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "102087 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "102087 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8965710,
            "unit": "ns/op\t 4914165 B/op\t   75241 allocs/op",
            "extra": "133 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8965710,
            "unit": "ns/op",
            "extra": "133 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914165,
            "unit": "B/op",
            "extra": "133 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "133 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 356.7,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3311301 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 356.7,
            "unit": "ns/op",
            "extra": "3311301 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3311301 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3311301 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.31,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.31,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 616.3,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "1957221 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 616.3,
            "unit": "ns/op",
            "extra": "1957221 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "1957221 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1957221 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 885.8,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1350160 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 885.8,
            "unit": "ns/op",
            "extra": "1350160 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1350160 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1350160 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "tao.yi@konghq.com",
            "name": "Tao Yi",
            "username": "randmonkey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "0c0013e764eb7f91adef63ef66c52ca3f6553c57",
          "message": "remove trigger image validation by label (#7887)",
          "timestamp": "2026-04-20T03:43:18Z",
          "tree_id": "8cc60383e0fcf6470f7919e08439ea57d14c3174",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/0c0013e764eb7f91adef63ef66c52ca3f6553c57"
        },
        "date": 1776658464455,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 773026,
            "unit": "ns/op\t  819452 B/op\t       5 allocs/op",
            "extra": "1350 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 773026,
            "unit": "ns/op",
            "extra": "1350 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819452,
            "unit": "B/op",
            "extra": "1350 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1350 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7498,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "158791 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7498,
            "unit": "ns/op",
            "extra": "158791 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "158791 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "158791 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 65.83,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18221466 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 65.83,
            "unit": "ns/op",
            "extra": "18221466 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "18221466 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "18221466 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21249,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "55944 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21249,
            "unit": "ns/op",
            "extra": "55944 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "55944 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "55944 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 223144,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5334 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 223144,
            "unit": "ns/op",
            "extra": "5334 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5334 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5334 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2553246,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "464 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2553246,
            "unit": "ns/op",
            "extra": "464 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "464 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "464 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 37876083,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "62 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 37876083,
            "unit": "ns/op",
            "extra": "62 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "62 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "62 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9701,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "125516 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9701,
            "unit": "ns/op",
            "extra": "125516 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "125516 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "125516 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7178421,
            "unit": "ns/op\t 4527550 B/op\t   69231 allocs/op",
            "extra": "168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7178421,
            "unit": "ns/op",
            "extra": "168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527550,
            "unit": "B/op",
            "extra": "168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69231,
            "unit": "allocs/op",
            "extra": "168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8146,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "145684 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8146,
            "unit": "ns/op",
            "extra": "145684 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "145684 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "145684 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 903553,
            "unit": "ns/op\t  392947 B/op\t    6173 allocs/op",
            "extra": "1318 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 903553,
            "unit": "ns/op",
            "extra": "1318 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392947,
            "unit": "B/op",
            "extra": "1318 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1318 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11332,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "104968 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11332,
            "unit": "ns/op",
            "extra": "104968 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "104968 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "104968 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8241610,
            "unit": "ns/op\t 4914184 B/op\t   75241 allocs/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8241610,
            "unit": "ns/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914184,
            "unit": "B/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 316.2,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3804182 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 316.2,
            "unit": "ns/op",
            "extra": "3804182 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3804182 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3804182 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 11.49,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 11.49,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 562,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2131996 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 562,
            "unit": "ns/op",
            "extra": "2131996 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2131996 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2131996 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 792.2,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1520121 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 792.2,
            "unit": "ns/op",
            "extra": "1520121 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1520121 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1520121 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "tao.yi@konghq.com",
            "name": "Tao Yi",
            "username": "randmonkey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "37230d26f552e3c0b25c220cee197cbacb90aacc",
          "message": "use eventually to check listener's condition become RefNotResolved (#7864)",
          "timestamp": "2026-04-29T07:18:22Z",
          "tree_id": "fd1b2d78eb160a0f7b1b44d27e23f52fe692ac23",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/37230d26f552e3c0b25c220cee197cbacb90aacc"
        },
        "date": 1777448984475,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 778474,
            "unit": "ns/op\t  819456 B/op\t       5 allocs/op",
            "extra": "1482 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 778474,
            "unit": "ns/op",
            "extra": "1482 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819456,
            "unit": "B/op",
            "extra": "1482 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1482 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7753,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "151328 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7753,
            "unit": "ns/op",
            "extra": "151328 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "151328 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "151328 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.83,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16450155 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.83,
            "unit": "ns/op",
            "extra": "16450155 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16450155 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16450155 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22631,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53271 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22631,
            "unit": "ns/op",
            "extra": "53271 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53271 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53271 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 229290,
            "unit": "ns/op\t  245762 B/op\t       2 allocs/op",
            "extra": "4510 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 229290,
            "unit": "ns/op",
            "extra": "4510 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245762,
            "unit": "B/op",
            "extra": "4510 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4510 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2438603,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "493 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2438603,
            "unit": "ns/op",
            "extra": "493 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "493 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "493 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 34880054,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 34880054,
            "unit": "ns/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 10292,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "118872 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10292,
            "unit": "ns/op",
            "extra": "118872 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "118872 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "118872 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7750714,
            "unit": "ns/op\t 4527539 B/op\t   69230 allocs/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7750714,
            "unit": "ns/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527539,
            "unit": "B/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8752,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "125257 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8752,
            "unit": "ns/op",
            "extra": "125257 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "125257 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "125257 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 892222,
            "unit": "ns/op\t  392941 B/op\t    6173 allocs/op",
            "extra": "1327 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 892222,
            "unit": "ns/op",
            "extra": "1327 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392941,
            "unit": "B/op",
            "extra": "1327 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1327 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11789,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "101236 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11789,
            "unit": "ns/op",
            "extra": "101236 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "101236 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "101236 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8991038,
            "unit": "ns/op\t 4914187 B/op\t   75241 allocs/op",
            "extra": "134 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8991038,
            "unit": "ns/op",
            "extra": "134 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914187,
            "unit": "B/op",
            "extra": "134 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "134 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 368.8,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3187120 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 368.8,
            "unit": "ns/op",
            "extra": "3187120 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3187120 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3187120 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.29,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.29,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 616.1,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "1965931 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 616.1,
            "unit": "ns/op",
            "extra": "1965931 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "1965931 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1965931 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 890.3,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1352161 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 890.3,
            "unit": "ns/op",
            "extra": "1352161 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1352161 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1352161 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "jakub.warczarek@konghq.com",
            "name": "Jakub Warczarek",
            "username": "programmer04"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "eae5a24e1d43ebd1c367f4ea82f2a473f8b2b96b",
          "message": "fix: secure by default (defaults from Gateway), working sectionName for `HTTPRoute` (#7901)",
          "timestamp": "2026-04-29T09:34:02Z",
          "tree_id": "9e4d536e28d80a5c16ffc6fc69a1e617b0e50b6e",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/eae5a24e1d43ebd1c367f4ea82f2a473f8b2b96b"
        },
        "date": 1777457162902,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 740097,
            "unit": "ns/op\t  819452 B/op\t       5 allocs/op",
            "extra": "1472 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 740097,
            "unit": "ns/op",
            "extra": "1472 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819452,
            "unit": "B/op",
            "extra": "1472 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1472 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7507,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "159213 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7507,
            "unit": "ns/op",
            "extra": "159213 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "159213 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "159213 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.87,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16443838 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.87,
            "unit": "ns/op",
            "extra": "16443838 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16443838 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16443838 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22803,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52329 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22803,
            "unit": "ns/op",
            "extra": "52329 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52329 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52329 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 228651,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4680 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 228651,
            "unit": "ns/op",
            "extra": "4680 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4680 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4680 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2448581,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "495 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2448581,
            "unit": "ns/op",
            "extra": "495 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "495 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "495 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 37591422,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 37591422,
            "unit": "ns/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 10063,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "120308 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10063,
            "unit": "ns/op",
            "extra": "120308 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "120308 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "120308 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7676067,
            "unit": "ns/op\t 4527547 B/op\t   69230 allocs/op",
            "extra": "156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7676067,
            "unit": "ns/op",
            "extra": "156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527547,
            "unit": "B/op",
            "extra": "156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8425,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "143383 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8425,
            "unit": "ns/op",
            "extra": "143383 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "143383 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "143383 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 869329,
            "unit": "ns/op\t  392951 B/op\t    6173 allocs/op",
            "extra": "1339 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 869329,
            "unit": "ns/op",
            "extra": "1339 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392951,
            "unit": "B/op",
            "extra": "1339 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1339 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11821,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "104611 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11821,
            "unit": "ns/op",
            "extra": "104611 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "104611 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "104611 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 9045475,
            "unit": "ns/op\t 4914192 B/op\t   75241 allocs/op",
            "extra": "130 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 9045475,
            "unit": "ns/op",
            "extra": "130 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914192,
            "unit": "B/op",
            "extra": "130 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "130 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 365.5,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3360934 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 365.5,
            "unit": "ns/op",
            "extra": "3360934 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3360934 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3360934 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.39,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.39,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 789.7,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "1534959 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 789.7,
            "unit": "ns/op",
            "extra": "1534959 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "1534959 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1534959 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 902.8,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1326866 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 902.8,
            "unit": "ns/op",
            "extra": "1326866 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1326866 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1326866 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "tao.yi@konghq.com",
            "name": "Tao Yi",
            "username": "randmonkey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "da0be9c8f5ada3fd46c4cb569082185924cdda7c",
          "message": "Use assert.Eventually to update HTTPRoute (#7911)",
          "timestamp": "2026-05-04T08:16:48Z",
          "tree_id": "cc6567b81dbdbe7f4aef17302bfb3d223f3a8864",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/da0be9c8f5ada3fd46c4cb569082185924cdda7c"
        },
        "date": 1777884514286,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 791429,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "1477 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 791429,
            "unit": "ns/op",
            "extra": "1477 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819453,
            "unit": "B/op",
            "extra": "1477 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1477 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7988,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "149198 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7988,
            "unit": "ns/op",
            "extra": "149198 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "149198 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "149198 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.86,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16185381 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.86,
            "unit": "ns/op",
            "extra": "16185381 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16185381 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16185381 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 24811,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52179 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 24811,
            "unit": "ns/op",
            "extra": "52179 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52179 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52179 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 306295,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "3662 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 306295,
            "unit": "ns/op",
            "extra": "3662 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "3662 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3662 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2471074,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "482 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2471074,
            "unit": "ns/op",
            "extra": "482 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "482 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "482 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 39304694,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "61 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 39304694,
            "unit": "ns/op",
            "extra": "61 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "61 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "61 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 10344,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "114855 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10344,
            "unit": "ns/op",
            "extra": "114855 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "114855 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "114855 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 8294766,
            "unit": "ns/op\t 4527544 B/op\t   69230 allocs/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 8294766,
            "unit": "ns/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527544,
            "unit": "B/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8432,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "137590 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8432,
            "unit": "ns/op",
            "extra": "137590 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "137590 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "137590 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 889388,
            "unit": "ns/op\t  392945 B/op\t    6173 allocs/op",
            "extra": "1320 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 889388,
            "unit": "ns/op",
            "extra": "1320 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392945,
            "unit": "B/op",
            "extra": "1320 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1320 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11960,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "100531 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11960,
            "unit": "ns/op",
            "extra": "100531 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "100531 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "100531 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 9198361,
            "unit": "ns/op\t 4914157 B/op\t   75241 allocs/op",
            "extra": "129 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 9198361,
            "unit": "ns/op",
            "extra": "129 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914157,
            "unit": "B/op",
            "extra": "129 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "129 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 326.5,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3688545 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 326.5,
            "unit": "ns/op",
            "extra": "3688545 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3688545 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3688545 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.36,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.36,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 596.8,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2005525 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 596.8,
            "unit": "ns/op",
            "extra": "2005525 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2005525 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2005525 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 924.6,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1361080 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 924.6,
            "unit": "ns/op",
            "extra": "1361080 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1361080 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1361080 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "pmalek@users.noreply.github.com",
            "name": "Patryk Małek",
            "username": "pmalek"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "38415bdb02014ac69417cf1027c6278ef88bf233",
          "message": "fix: fix unsanitized KongPlugin being included in sanitzed config dump (#7912)\n\n* fix: fix unsanitized KongPlugin being included in sanitzed config dump\n\n* tests: fix TestKongState_SanitizedCopy\n\n* chore: add changelog",
          "timestamp": "2026-05-04T16:06:29Z",
          "tree_id": "ed53c020100d16a9b42ef75daa0b944f2e5b049f",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/38415bdb02014ac69417cf1027c6278ef88bf233"
        },
        "date": 1777912744963,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 676994,
            "unit": "ns/op\t  819451 B/op\t       5 allocs/op",
            "extra": "1572 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 676994,
            "unit": "ns/op",
            "extra": "1572 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819451,
            "unit": "B/op",
            "extra": "1572 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1572 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 6644,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "178689 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 6644,
            "unit": "ns/op",
            "extra": "178689 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "178689 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "178689 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 94.29,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12737697 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 94.29,
            "unit": "ns/op",
            "extra": "12737697 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12737697 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12737697 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 19700,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "61312 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 19700,
            "unit": "ns/op",
            "extra": "61312 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "61312 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "61312 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 202634,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5163 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 202634,
            "unit": "ns/op",
            "extra": "5163 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5163 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5163 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2396684,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "428 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2396684,
            "unit": "ns/op",
            "extra": "428 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "428 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "428 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 32804709,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 32804709,
            "unit": "ns/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9294,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "126738 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9294,
            "unit": "ns/op",
            "extra": "126738 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "126738 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "126738 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6628437,
            "unit": "ns/op\t 4527547 B/op\t   69230 allocs/op",
            "extra": "180 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6628437,
            "unit": "ns/op",
            "extra": "180 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527547,
            "unit": "B/op",
            "extra": "180 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "180 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7807,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "156535 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7807,
            "unit": "ns/op",
            "extra": "156535 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "156535 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "156535 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 848495,
            "unit": "ns/op\t  392922 B/op\t    6173 allocs/op",
            "extra": "1389 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 848495,
            "unit": "ns/op",
            "extra": "1389 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392922,
            "unit": "B/op",
            "extra": "1389 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1389 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11090,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "108529 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11090,
            "unit": "ns/op",
            "extra": "108529 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "108529 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "108529 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7871634,
            "unit": "ns/op\t 4914200 B/op\t   75241 allocs/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7871634,
            "unit": "ns/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914200,
            "unit": "B/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 314.9,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3814089 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 314.9,
            "unit": "ns/op",
            "extra": "3814089 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3814089 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3814089 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.23,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.23,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 579.5,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2059454 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 579.5,
            "unit": "ns/op",
            "extra": "2059454 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2059454 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2059454 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 874.5,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1367485 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 874.5,
            "unit": "ns/op",
            "extra": "1367485 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1367485 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1367485 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "jakub.warczarek@konghq.com",
            "name": "Jakub Warczarek",
            "username": "programmer04"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "34bfcda970b9c04551931bc4e6b72a2ad9bd9b06",
          "message": "fix: more robust `HTTPRoute` validation (#7913)",
          "timestamp": "2026-05-06T08:32:44Z",
          "tree_id": "49cebd636b71126576db5ae2009825030cc19a54",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/34bfcda970b9c04551931bc4e6b72a2ad9bd9b06"
        },
        "date": 1778058369440,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 767587,
            "unit": "ns/op\t  819450 B/op\t       5 allocs/op",
            "extra": "1420 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 767587,
            "unit": "ns/op",
            "extra": "1420 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819450,
            "unit": "B/op",
            "extra": "1420 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1420 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7837,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "145304 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7837,
            "unit": "ns/op",
            "extra": "145304 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "145304 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "145304 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 73.04,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16438458 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 73.04,
            "unit": "ns/op",
            "extra": "16438458 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16438458 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16438458 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22028,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54781 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22028,
            "unit": "ns/op",
            "extra": "54781 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54781 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54781 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 215572,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5025 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 215572,
            "unit": "ns/op",
            "extra": "5025 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5025 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5025 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2355830,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "512 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2355830,
            "unit": "ns/op",
            "extra": "512 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "512 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "512 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 40151511,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 40151511,
            "unit": "ns/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9769,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "123060 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9769,
            "unit": "ns/op",
            "extra": "123060 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "123060 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "123060 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7423234,
            "unit": "ns/op\t 4527537 B/op\t   69230 allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7423234,
            "unit": "ns/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527537,
            "unit": "B/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8293,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "143252 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8293,
            "unit": "ns/op",
            "extra": "143252 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "143252 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "143252 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 876413,
            "unit": "ns/op\t  392950 B/op\t    6173 allocs/op",
            "extra": "1381 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 876413,
            "unit": "ns/op",
            "extra": "1381 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392950,
            "unit": "B/op",
            "extra": "1381 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1381 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11395,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "108187 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11395,
            "unit": "ns/op",
            "extra": "108187 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "108187 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "108187 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8092218,
            "unit": "ns/op\t 4914219 B/op\t   75241 allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8092218,
            "unit": "ns/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914219,
            "unit": "B/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 327.1,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3665766 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 327.1,
            "unit": "ns/op",
            "extra": "3665766 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3665766 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3665766 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.58,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.58,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 611.1,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "1967617 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 611.1,
            "unit": "ns/op",
            "extra": "1967617 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "1967617 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1967617 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 876.4,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1369053 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 876.4,
            "unit": "ns/op",
            "extra": "1369053 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1369053 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1369053 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "tao.yi@konghq.com",
            "name": "Tao Yi",
            "username": "randmonkey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": false,
          "id": "451dd707b002475d48af49541743a578372d5d7e",
          "message": "tests(conformance): Skip partial conformance tests related to `GRPCRoute` for bumping grpc package and use separate vars in update to reduce data race (#7916)\n\n* chore(deps): bump google.golang.org/grpc from 1.73.0 to 1.79.3\n\nBumps [google.golang.org/grpc](https://github.com/grpc/grpc-go) from 1.73.0 to 1.79.3.\n- [Release notes](https://github.com/grpc/grpc-go/releases)\n- [Commits](https://github.com/grpc/grpc-go/compare/v1.73.0...v1.79.3)\n\n---\nupdated-dependencies:\n- dependency-name: google.golang.org/grpc\n  dependency-version: 1.79.3\n  dependency-type: direct:production\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\n\n* skip grpcroute related tests\n\n* fix data race by deepcopy\n\n* define new var for updating HTTPRoute in tests\n\n* add comments\n\n---------\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2026-05-07T11:41:10Z",
          "tree_id": "e24f3825893ffebea41e4968fb8ca58b80091248",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/451dd707b002475d48af49541743a578372d5d7e"
        },
        "date": 1778157726076,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 831384,
            "unit": "ns/op\t  819454 B/op\t       5 allocs/op",
            "extra": "1338 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 831384,
            "unit": "ns/op",
            "extra": "1338 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819454,
            "unit": "B/op",
            "extra": "1338 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1338 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 8642,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "134506 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 8642,
            "unit": "ns/op",
            "extra": "134506 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "134506 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "134506 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 82.42,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14800434 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 82.42,
            "unit": "ns/op",
            "extra": "14800434 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14800434 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14800434 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 19730,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "60712 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 19730,
            "unit": "ns/op",
            "extra": "60712 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "60712 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "60712 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 213998,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5077 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 213998,
            "unit": "ns/op",
            "extra": "5077 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5077 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5077 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2485621,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "483 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2485621,
            "unit": "ns/op",
            "extra": "483 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "483 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "483 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 36673808,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 36673808,
            "unit": "ns/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9704,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "121384 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9704,
            "unit": "ns/op",
            "extra": "121384 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "121384 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "121384 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6898149,
            "unit": "ns/op\t 4527545 B/op\t   69230 allocs/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6898149,
            "unit": "ns/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527545,
            "unit": "B/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8090,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "148839 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8090,
            "unit": "ns/op",
            "extra": "148839 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "148839 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "148839 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 860120,
            "unit": "ns/op\t  392950 B/op\t    6173 allocs/op",
            "extra": "1383 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 860120,
            "unit": "ns/op",
            "extra": "1383 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392950,
            "unit": "B/op",
            "extra": "1383 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1383 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11365,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "105914 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11365,
            "unit": "ns/op",
            "extra": "105914 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "105914 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "105914 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8047674,
            "unit": "ns/op\t 4914200 B/op\t   75241 allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8047674,
            "unit": "ns/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914200,
            "unit": "B/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 317.3,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3806298 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 317.3,
            "unit": "ns/op",
            "extra": "3806298 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3806298 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3806298 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.59,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.59,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 583.6,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2033508 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 583.6,
            "unit": "ns/op",
            "extra": "2033508 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2033508 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2033508 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 847.3,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1412384 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 847.3,
            "unit": "ns/op",
            "extra": "1412384 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1412384 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1412384 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "sandro.modarelli@gmail.com",
            "name": "Sandro Modarelli",
            "username": "sandromodarelli"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "84fa1189cfcf84239d06c240043d972902845721",
          "message": "fix: implement ReferenceGrant checks for cross-namespace certificate references (#7920)\n\n* fix: implement ReferenceGrant checks for cross-namespace certificate references\n\n* chore(changelog): update changelog for ReferenceGrant checks on cross-namespace certificate references\n\n* test: improve readiness check assertions in manager tests",
          "timestamp": "2026-05-07T14:07:19Z",
          "tree_id": "7479fc6014286b05c9332c9b730c38b64ad820d7",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/84fa1189cfcf84239d06c240043d972902845721"
        },
        "date": 1778164851055,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 777031,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "1413 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 777031,
            "unit": "ns/op",
            "extra": "1413 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819453,
            "unit": "B/op",
            "extra": "1413 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1413 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7682,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "156958 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7682,
            "unit": "ns/op",
            "extra": "156958 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "156958 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "156958 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.76,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16469608 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.76,
            "unit": "ns/op",
            "extra": "16469608 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16469608 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16469608 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22078,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53552 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22078,
            "unit": "ns/op",
            "extra": "53552 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53552 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53552 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 250821,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5095 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 250821,
            "unit": "ns/op",
            "extra": "5095 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5095 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5095 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2397630,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "469 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2397630,
            "unit": "ns/op",
            "extra": "469 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "469 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "469 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 35013062,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "30 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 35013062,
            "unit": "ns/op",
            "extra": "30 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "30 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "30 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9685,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "124108 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9685,
            "unit": "ns/op",
            "extra": "124108 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "124108 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "124108 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7114069,
            "unit": "ns/op\t 4527540 B/op\t   69230 allocs/op",
            "extra": "168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7114069,
            "unit": "ns/op",
            "extra": "168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527540,
            "unit": "B/op",
            "extra": "168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8100,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "148316 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8100,
            "unit": "ns/op",
            "extra": "148316 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "148316 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "148316 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 857575,
            "unit": "ns/op\t  392929 B/op\t    6173 allocs/op",
            "extra": "1357 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 857575,
            "unit": "ns/op",
            "extra": "1357 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392929,
            "unit": "B/op",
            "extra": "1357 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1357 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11245,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "107728 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11245,
            "unit": "ns/op",
            "extra": "107728 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "107728 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "107728 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8058038,
            "unit": "ns/op\t 4914135 B/op\t   75241 allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8058038,
            "unit": "ns/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914135,
            "unit": "B/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 326.6,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3630164 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 326.6,
            "unit": "ns/op",
            "extra": "3630164 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3630164 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3630164 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.3,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 590.5,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2030194 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 590.5,
            "unit": "ns/op",
            "extra": "2030194 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2030194 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2030194 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 867.5,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1385792 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 867.5,
            "unit": "ns/op",
            "extra": "1385792 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1385792 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1385792 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "jakub.warczarek@konghq.com",
            "name": "Jakub Warczarek",
            "username": "programmer04"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "ef09bd70c9eea91f0398a5482ff75bb194496e7a",
          "message": "chore(deps): bump golang from 1.25.7 to 1.25.10 (#7924)",
          "timestamp": "2026-05-11T10:17:53Z",
          "tree_id": "9bc9ee62d8d9b2ccd030002ad22e34b6fe7e88db",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/ef09bd70c9eea91f0398a5482ff75bb194496e7a"
        },
        "date": 1778496626135,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 587971,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "1742 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 587971,
            "unit": "ns/op",
            "extra": "1742 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819453,
            "unit": "B/op",
            "extra": "1742 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1742 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7154,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "165558 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7154,
            "unit": "ns/op",
            "extra": "165558 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "165558 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "165558 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.78,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16452750 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.78,
            "unit": "ns/op",
            "extra": "16452750 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16452750 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16452750 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22237,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "51198 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22237,
            "unit": "ns/op",
            "extra": "51198 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "51198 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "51198 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 209863,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5017 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 209863,
            "unit": "ns/op",
            "extra": "5017 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5017 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5017 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2294410,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "511 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2294410,
            "unit": "ns/op",
            "extra": "511 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "511 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "511 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 37487286,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 37487286,
            "unit": "ns/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 11415,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "97184 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 11415,
            "unit": "ns/op",
            "extra": "97184 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "97184 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "97184 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6928215,
            "unit": "ns/op\t 4527540 B/op\t   69230 allocs/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6928215,
            "unit": "ns/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527540,
            "unit": "B/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8027,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "147170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8027,
            "unit": "ns/op",
            "extra": "147170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "147170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "147170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 843382,
            "unit": "ns/op\t  392932 B/op\t    6173 allocs/op",
            "extra": "1387 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 843382,
            "unit": "ns/op",
            "extra": "1387 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392932,
            "unit": "B/op",
            "extra": "1387 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1387 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11107,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "109196 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11107,
            "unit": "ns/op",
            "extra": "109196 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "109196 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "109196 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7968126,
            "unit": "ns/op\t 4914200 B/op\t   75241 allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7968126,
            "unit": "ns/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914200,
            "unit": "B/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 325.6,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3706826 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 325.6,
            "unit": "ns/op",
            "extra": "3706826 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3706826 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3706826 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.29,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.29,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 596.2,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2023396 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 596.2,
            "unit": "ns/op",
            "extra": "2023396 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2023396 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2023396 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 1126,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 1126,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "jakub.warczarek@konghq.com",
            "name": "Jakub Warczarek",
            "username": "programmer04"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "4ed9df170fbf015d806459ff3fa1f2db0a914ae2",
          "message": "chore: adjust CHANGELOG for `3.4.14` and `3.5.7` releases (#7929)",
          "timestamp": "2026-05-11T18:26:20+02:00",
          "tree_id": "552d3fac285b44ed8657be4bf3ac6eb49dd938e0",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/4ed9df170fbf015d806459ff3fa1f2db0a914ae2"
        },
        "date": 1778517005692,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 753875,
            "unit": "ns/op\t  819454 B/op\t       5 allocs/op",
            "extra": "1368 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 753875,
            "unit": "ns/op",
            "extra": "1368 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819454,
            "unit": "B/op",
            "extra": "1368 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1368 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 6879,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "147841 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 6879,
            "unit": "ns/op",
            "extra": "147841 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "147841 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "147841 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 81.21,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14735714 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 81.21,
            "unit": "ns/op",
            "extra": "14735714 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14735714 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14735714 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 20171,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "56605 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 20171,
            "unit": "ns/op",
            "extra": "56605 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "56605 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "56605 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 219138,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4938 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 219138,
            "unit": "ns/op",
            "extra": "4938 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4938 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4938 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2756531,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "504 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2756531,
            "unit": "ns/op",
            "extra": "504 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "504 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "504 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 33847818,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 33847818,
            "unit": "ns/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9574,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "123730 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9574,
            "unit": "ns/op",
            "extra": "123730 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "123730 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "123730 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6797281,
            "unit": "ns/op\t 4527542 B/op\t   69230 allocs/op",
            "extra": "176 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6797281,
            "unit": "ns/op",
            "extra": "176 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527542,
            "unit": "B/op",
            "extra": "176 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "176 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8014,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "148648 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8014,
            "unit": "ns/op",
            "extra": "148648 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "148648 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "148648 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 870142,
            "unit": "ns/op\t  392938 B/op\t    6173 allocs/op",
            "extra": "1370 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 870142,
            "unit": "ns/op",
            "extra": "1370 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392938,
            "unit": "B/op",
            "extra": "1370 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1370 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11125,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "109644 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11125,
            "unit": "ns/op",
            "extra": "109644 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "109644 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "109644 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7802976,
            "unit": "ns/op\t 4914173 B/op\t   75241 allocs/op",
            "extra": "154 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7802976,
            "unit": "ns/op",
            "extra": "154 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914173,
            "unit": "B/op",
            "extra": "154 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "154 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 320.5,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3757335 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 320.5,
            "unit": "ns/op",
            "extra": "3757335 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3757335 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3757335 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.55,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.55,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 588.1,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2061542 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 588.1,
            "unit": "ns/op",
            "extra": "2061542 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2061542 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2061542 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 843.1,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1415164 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 843.1,
            "unit": "ns/op",
            "extra": "1415164 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1415164 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1415164 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "49699333+dependabot[bot]@users.noreply.github.com",
            "name": "dependabot[bot]",
            "username": "dependabot[bot]"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": false,
          "id": "35e1368bda2061a8eb7491ae7755285442c7dbdc",
          "message": "chore(deps): bump github.com/moby/spdystream from 0.5.0 to 0.5.1 (#7896)\n\nBumps [github.com/moby/spdystream](https://github.com/moby/spdystream) from 0.5.0 to 0.5.1.\n- [Release notes](https://github.com/moby/spdystream/releases)\n- [Commits](https://github.com/moby/spdystream/compare/v0.5.0...v0.5.1)\n\n---\nupdated-dependencies:\n- dependency-name: github.com/moby/spdystream\n  dependency-version: 0.5.1\n  dependency-type: indirect\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2026-05-11T16:54:59Z",
          "tree_id": "333011ce105b24a89263e1a04f26781ab8ea30d4",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/35e1368bda2061a8eb7491ae7755285442c7dbdc"
        },
        "date": 1778520467995,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 852489,
            "unit": "ns/op\t  819451 B/op\t       5 allocs/op",
            "extra": "1267 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 852489,
            "unit": "ns/op",
            "extra": "1267 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819451,
            "unit": "B/op",
            "extra": "1267 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1267 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7286,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "161350 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7286,
            "unit": "ns/op",
            "extra": "161350 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "161350 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "161350 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 65.68,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18269846 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 65.68,
            "unit": "ns/op",
            "extra": "18269846 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "18269846 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "18269846 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21584,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "55141 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21584,
            "unit": "ns/op",
            "extra": "55141 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "55141 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "55141 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 238609,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5013 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 238609,
            "unit": "ns/op",
            "extra": "5013 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5013 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5013 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2721370,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "475 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2721370,
            "unit": "ns/op",
            "extra": "475 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "475 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "475 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 43329136,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "38 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 43329136,
            "unit": "ns/op",
            "extra": "38 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "38 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "38 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 10036,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "120122 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10036,
            "unit": "ns/op",
            "extra": "120122 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "120122 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "120122 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7501750,
            "unit": "ns/op\t 4527546 B/op\t   69230 allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7501750,
            "unit": "ns/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527546,
            "unit": "B/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8280,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "149670 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8280,
            "unit": "ns/op",
            "extra": "149670 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "149670 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "149670 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 900924,
            "unit": "ns/op\t  392935 B/op\t    6173 allocs/op",
            "extra": "1327 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 900924,
            "unit": "ns/op",
            "extra": "1327 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392935,
            "unit": "B/op",
            "extra": "1327 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1327 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11581,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "104966 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11581,
            "unit": "ns/op",
            "extra": "104966 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "104966 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "104966 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8356112,
            "unit": "ns/op\t 4914156 B/op\t   75241 allocs/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8356112,
            "unit": "ns/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914156,
            "unit": "B/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 316.5,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3744235 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 316.5,
            "unit": "ns/op",
            "extra": "3744235 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3744235 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3744235 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 11.44,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 11.44,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 565.2,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2112451 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 565.2,
            "unit": "ns/op",
            "extra": "2112451 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2112451 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2112451 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 794.9,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1512289 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 794.9,
            "unit": "ns/op",
            "extra": "1512289 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1512289 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1512289 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "zhangjintao9020@gmail.com",
            "name": "Jintao Zhang",
            "username": "tao12345666333"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "c55648d748c0455c9bf9598c5604e300c7f52e04",
          "message": "chore: bump kong gateway to 3.12 for testing (#7905)\n\nSigned-off-by: Jintao Zhang <zhangjintao9020@gmail.com>",
          "timestamp": "2026-05-13T21:59:58Z",
          "tree_id": "e86c223f8c5d3d890b24cf1a76f4bac73b089eb3",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/c55648d748c0455c9bf9598c5604e300c7f52e04"
        },
        "date": 1778711432005,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 711846,
            "unit": "ns/op\t  819452 B/op\t       5 allocs/op",
            "extra": "2008 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 711846,
            "unit": "ns/op",
            "extra": "2008 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819452,
            "unit": "B/op",
            "extra": "2008 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "2008 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7853,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "147796 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7853,
            "unit": "ns/op",
            "extra": "147796 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "147796 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "147796 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.97,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16444323 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.97,
            "unit": "ns/op",
            "extra": "16444323 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16444323 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16444323 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22606,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52594 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22606,
            "unit": "ns/op",
            "extra": "52594 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52594 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52594 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 223899,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4762 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 223899,
            "unit": "ns/op",
            "extra": "4762 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4762 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4762 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2391828,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "487 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2391828,
            "unit": "ns/op",
            "extra": "487 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "487 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "487 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 32966967,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 32966967,
            "unit": "ns/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 10074,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "118878 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10074,
            "unit": "ns/op",
            "extra": "118878 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "118878 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "118878 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7465620,
            "unit": "ns/op\t 4527537 B/op\t   69230 allocs/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7465620,
            "unit": "ns/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527537,
            "unit": "B/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8410,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "144864 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8410,
            "unit": "ns/op",
            "extra": "144864 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "144864 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "144864 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 868295,
            "unit": "ns/op\t  392926 B/op\t    6173 allocs/op",
            "extra": "1358 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 868295,
            "unit": "ns/op",
            "extra": "1358 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392926,
            "unit": "B/op",
            "extra": "1358 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1358 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11363,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "106900 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11363,
            "unit": "ns/op",
            "extra": "106900 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "106900 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "106900 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8294384,
            "unit": "ns/op\t 4914220 B/op\t   75241 allocs/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8294384,
            "unit": "ns/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914220,
            "unit": "B/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 333,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3539959 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 333,
            "unit": "ns/op",
            "extra": "3539959 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3539959 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3539959 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.29,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.29,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 729.4,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "1731934 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 729.4,
            "unit": "ns/op",
            "extra": "1731934 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "1731934 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1731934 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 937.4,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1264584 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 937.4,
            "unit": "ns/op",
            "extra": "1264584 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1264584 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1264584 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "jakub.warczarek@konghq.com",
            "name": "Jakub Warczarek",
            "username": "programmer04"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": false,
          "id": "bf2c0c7934b9f5fc11f497eda2f3ba8be186d640",
          "message": "chore: set protocol to `http` for Ingress in e2e tests (#7931)",
          "timestamp": "2026-05-19T02:16:42Z",
          "tree_id": "f1ce15b7ebc053bac70448607aaa7322d5e882e6",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/bf2c0c7934b9f5fc11f497eda2f3ba8be186d640"
        },
        "date": 1779158842505,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 654294,
            "unit": "ns/op\t  819454 B/op\t       5 allocs/op",
            "extra": "1537 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 654294,
            "unit": "ns/op",
            "extra": "1537 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819454,
            "unit": "B/op",
            "extra": "1537 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1537 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7167,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "167631 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7167,
            "unit": "ns/op",
            "extra": "167631 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "167631 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "167631 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.75,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16464480 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.75,
            "unit": "ns/op",
            "extra": "16464480 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16464480 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16464480 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22586,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53619 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22586,
            "unit": "ns/op",
            "extra": "53619 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53619 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53619 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 267309,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "3805 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 267309,
            "unit": "ns/op",
            "extra": "3805 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "3805 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3805 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2346554,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "516 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2346554,
            "unit": "ns/op",
            "extra": "516 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "516 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "516 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 37087329,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 37087329,
            "unit": "ns/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9632,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "126084 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9632,
            "unit": "ns/op",
            "extra": "126084 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "126084 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "126084 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7041110,
            "unit": "ns/op\t 4527573 B/op\t   69230 allocs/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7041110,
            "unit": "ns/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527573,
            "unit": "B/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8122,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "145050 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8122,
            "unit": "ns/op",
            "extra": "145050 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "145050 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "145050 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 853470,
            "unit": "ns/op\t  392951 B/op\t    6173 allocs/op",
            "extra": "1374 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 853470,
            "unit": "ns/op",
            "extra": "1374 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392951,
            "unit": "B/op",
            "extra": "1374 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1374 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11114,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "107708 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11114,
            "unit": "ns/op",
            "extra": "107708 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "107708 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "107708 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8026312,
            "unit": "ns/op\t 4914175 B/op\t   75241 allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8026312,
            "unit": "ns/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914175,
            "unit": "B/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 325.7,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3673111 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 325.7,
            "unit": "ns/op",
            "extra": "3673111 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3673111 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3673111 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.42,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.42,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 583.4,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2054895 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 583.4,
            "unit": "ns/op",
            "extra": "2054895 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2054895 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2054895 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 847.9,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1410975 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 847.9,
            "unit": "ns/op",
            "extra": "1410975 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1410975 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1410975 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "29139614+renovate[bot]@users.noreply.github.com",
            "name": "renovate[bot]",
            "username": "renovate[bot]"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "7daa0653531c7adf7295453b8d5e6ac5732fef7b",
          "message": "chore(deps): update kong/kong-gateway docker tag to v3.12.0.6 (#7930)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2026-05-26T08:24:07Z",
          "tree_id": "d9f15a3c96ba7819b62a4568d800c56260fc0f9e",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/7daa0653531c7adf7295453b8d5e6ac5732fef7b"
        },
        "date": 1779785694012,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 719117,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "1458 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 719117,
            "unit": "ns/op",
            "extra": "1458 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819453,
            "unit": "B/op",
            "extra": "1458 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1458 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 8320,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "144535 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 8320,
            "unit": "ns/op",
            "extra": "144535 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "144535 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "144535 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 82.49,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14737272 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 82.49,
            "unit": "ns/op",
            "extra": "14737272 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14737272 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14737272 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 19553,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "61179 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 19553,
            "unit": "ns/op",
            "extra": "61179 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "61179 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "61179 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 209219,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5146 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 209219,
            "unit": "ns/op",
            "extra": "5146 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5146 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5146 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2535536,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "415 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2535536,
            "unit": "ns/op",
            "extra": "415 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "415 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "415 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 35827041,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 35827041,
            "unit": "ns/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9261,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "123778 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9261,
            "unit": "ns/op",
            "extra": "123778 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "123778 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "123778 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6690073,
            "unit": "ns/op\t 4527549 B/op\t   69230 allocs/op",
            "extra": "177 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6690073,
            "unit": "ns/op",
            "extra": "177 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527549,
            "unit": "B/op",
            "extra": "177 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "177 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7729,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "158000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7729,
            "unit": "ns/op",
            "extra": "158000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "158000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "158000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 838614,
            "unit": "ns/op\t  392938 B/op\t    6173 allocs/op",
            "extra": "1388 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 838614,
            "unit": "ns/op",
            "extra": "1388 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392938,
            "unit": "B/op",
            "extra": "1388 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1388 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10782,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "113212 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10782,
            "unit": "ns/op",
            "extra": "113212 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "113212 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "113212 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7762082,
            "unit": "ns/op\t 4914248 B/op\t   75241 allocs/op",
            "extra": "154 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7762082,
            "unit": "ns/op",
            "extra": "154 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914248,
            "unit": "B/op",
            "extra": "154 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "154 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 310.5,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3854967 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 310.5,
            "unit": "ns/op",
            "extra": "3854967 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3854967 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3854967 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.22,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.22,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 655.7,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2081336 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 655.7,
            "unit": "ns/op",
            "extra": "2081336 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2081336 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2081336 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 849,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1378816 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 849,
            "unit": "ns/op",
            "extra": "1378816 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1378816 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1378816 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "29139614+renovate[bot]@users.noreply.github.com",
            "name": "renovate[bot]",
            "username": "renovate[bot]"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "cdafba58c1d560ace63cbee060379ba2a17cbd87",
          "message": "chore(deps): update dependency gke to v1.35.5 (#7936)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2026-05-28T06:36:40Z",
          "tree_id": "27f358cd0736d780c99bbd8eee76d8327bd59aba",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/cdafba58c1d560ace63cbee060379ba2a17cbd87"
        },
        "date": 1779952073338,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 725994,
            "unit": "ns/op\t  819454 B/op\t       5 allocs/op",
            "extra": "1482 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 725994,
            "unit": "ns/op",
            "extra": "1482 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819454,
            "unit": "B/op",
            "extra": "1482 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1482 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 6998,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "172250 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 6998,
            "unit": "ns/op",
            "extra": "172250 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "172250 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "172250 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 82.23,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14717870 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 82.23,
            "unit": "ns/op",
            "extra": "14717870 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14717870 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14717870 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 19485,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "59623 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 19485,
            "unit": "ns/op",
            "extra": "59623 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "59623 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "59623 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 201866,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5214 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 201866,
            "unit": "ns/op",
            "extra": "5214 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5214 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5214 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2376191,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "432 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2376191,
            "unit": "ns/op",
            "extra": "432 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "432 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "432 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 43260575,
            "unit": "ns/op\t24010754 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 43260575,
            "unit": "ns/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010754,
            "unit": "B/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9811,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "103226 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9811,
            "unit": "ns/op",
            "extra": "103226 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "103226 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "103226 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6719430,
            "unit": "ns/op\t 4527546 B/op\t   69230 allocs/op",
            "extra": "178 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6719430,
            "unit": "ns/op",
            "extra": "178 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527546,
            "unit": "B/op",
            "extra": "178 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "178 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7773,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "155012 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7773,
            "unit": "ns/op",
            "extra": "155012 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "155012 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "155012 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 851549,
            "unit": "ns/op\t  392937 B/op\t    6173 allocs/op",
            "extra": "1381 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 851549,
            "unit": "ns/op",
            "extra": "1381 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392937,
            "unit": "B/op",
            "extra": "1381 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1381 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10862,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "111926 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10862,
            "unit": "ns/op",
            "extra": "111926 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "111926 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "111926 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7770734,
            "unit": "ns/op\t 4914196 B/op\t   75241 allocs/op",
            "extra": "153 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7770734,
            "unit": "ns/op",
            "extra": "153 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914196,
            "unit": "B/op",
            "extra": "153 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "153 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 312.1,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3855974 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 312.1,
            "unit": "ns/op",
            "extra": "3855974 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3855974 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3855974 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.27,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.27,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 625.9,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2087714 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 625.9,
            "unit": "ns/op",
            "extra": "2087714 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2087714 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2087714 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 1018,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 1018,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "pmalek@users.noreply.github.com",
            "name": "Patryk Małek",
            "username": "pmalek"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "1f2a41b03c7e51f9fc61d67d84142441a9f7699c",
          "message": "Revert \"fix: fix unsanitized KongPlugin being included in sanitzed config dump (#7912)\" (#7937)\n\nThis reverts commit 38415bdb02014ac69417cf1027c6278ef88bf233.",
          "timestamp": "2026-05-28T12:15:41+02:00",
          "tree_id": "39f87482ac6e1d2ccfa8497e37a9cdb6991a2e1e",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/1f2a41b03c7e51f9fc61d67d84142441a9f7699c"
        },
        "date": 1779963566704,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 669452,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "1652 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 669452,
            "unit": "ns/op",
            "extra": "1652 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819453,
            "unit": "B/op",
            "extra": "1652 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1652 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7360,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "161817 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7360,
            "unit": "ns/op",
            "extra": "161817 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "161817 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "161817 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.79,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16473651 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.79,
            "unit": "ns/op",
            "extra": "16473651 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16473651 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16473651 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22143,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54026 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22143,
            "unit": "ns/op",
            "extra": "54026 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54026 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54026 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 268245,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5894 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 268245,
            "unit": "ns/op",
            "extra": "5894 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5894 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5894 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2385767,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "500 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2385767,
            "unit": "ns/op",
            "extra": "500 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "500 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "500 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 47164575,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "34 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 47164575,
            "unit": "ns/op",
            "extra": "34 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "34 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "34 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9497,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "125198 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9497,
            "unit": "ns/op",
            "extra": "125198 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "125198 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "125198 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6940059,
            "unit": "ns/op\t 4527540 B/op\t   69230 allocs/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6940059,
            "unit": "ns/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527540,
            "unit": "B/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7963,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "148555 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7963,
            "unit": "ns/op",
            "extra": "148555 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "148555 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "148555 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 838701,
            "unit": "ns/op\t  392945 B/op\t    6173 allocs/op",
            "extra": "1400 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 838701,
            "unit": "ns/op",
            "extra": "1400 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392945,
            "unit": "B/op",
            "extra": "1400 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1400 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10920,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "108783 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10920,
            "unit": "ns/op",
            "extra": "108783 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "108783 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "108783 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8012852,
            "unit": "ns/op\t 4914184 B/op\t   75241 allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8012852,
            "unit": "ns/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914184,
            "unit": "B/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 324.4,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3714332 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 324.4,
            "unit": "ns/op",
            "extra": "3714332 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3714332 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3714332 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.29,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.29,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 583.1,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2063209 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 583.1,
            "unit": "ns/op",
            "extra": "2063209 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2063209 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2063209 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 845.8,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1419082 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 845.8,
            "unit": "ns/op",
            "extra": "1419082 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1419082 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1419082 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "sandro.modarelli@gmail.com",
            "name": "Sandro Modarelli",
            "username": "sandromodarelli"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "e0d7b00a236242e3c5b096386077ce5886ca1cc9",
          "message": "fix(adminapi): preserve TLSServerName on admin client recreation (#7950)\n\n* fix(adminapi): preserve TLSServerName on admin client recreation\n\n* test(clients): cover TLSServerName preservation on readiness recovery\n\n* docs(changelog): note admin API SNI fix on proxy restart",
          "timestamp": "2026-06-04T12:23:37Z",
          "tree_id": "b1a89bb873d02db55c585235911f763c48c1a5bd",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/e0d7b00a236242e3c5b096386077ce5886ca1cc9"
        },
        "date": 1780577741723,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 680824,
            "unit": "ns/op\t  819450 B/op\t       5 allocs/op",
            "extra": "1548 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 680824,
            "unit": "ns/op",
            "extra": "1548 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819450,
            "unit": "B/op",
            "extra": "1548 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1548 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 6938,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "167421 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 6938,
            "unit": "ns/op",
            "extra": "167421 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "167421 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "167421 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 97.04,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12318478 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 97.04,
            "unit": "ns/op",
            "extra": "12318478 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12318478 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12318478 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 19398,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "62161 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 19398,
            "unit": "ns/op",
            "extra": "62161 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "62161 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "62161 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 209919,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5134 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 209919,
            "unit": "ns/op",
            "extra": "5134 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5134 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5134 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2612772,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "506 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2612772,
            "unit": "ns/op",
            "extra": "506 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "506 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "506 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 36175525,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 36175525,
            "unit": "ns/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9227,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "129076 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9227,
            "unit": "ns/op",
            "extra": "129076 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "129076 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "129076 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6848219,
            "unit": "ns/op\t 4527555 B/op\t   69231 allocs/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6848219,
            "unit": "ns/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527555,
            "unit": "B/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69231,
            "unit": "allocs/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7794,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "152271 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7794,
            "unit": "ns/op",
            "extra": "152271 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "152271 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "152271 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 886414,
            "unit": "ns/op\t  392939 B/op\t    6173 allocs/op",
            "extra": "1380 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 886414,
            "unit": "ns/op",
            "extra": "1380 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392939,
            "unit": "B/op",
            "extra": "1380 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1380 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10849,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "107456 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10849,
            "unit": "ns/op",
            "extra": "107456 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "107456 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "107456 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7727041,
            "unit": "ns/op\t 4914185 B/op\t   75241 allocs/op",
            "extra": "154 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7727041,
            "unit": "ns/op",
            "extra": "154 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914185,
            "unit": "B/op",
            "extra": "154 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "154 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 317,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3770247 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 317,
            "unit": "ns/op",
            "extra": "3770247 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3770247 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3770247 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.24,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.24,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 583.1,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2109153 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 583.1,
            "unit": "ns/op",
            "extra": "2109153 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2109153 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2109153 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 833.6,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1444092 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 833.6,
            "unit": "ns/op",
            "extra": "1444092 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1444092 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1444092 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "49699333+dependabot[bot]@users.noreply.github.com",
            "name": "dependabot[bot]",
            "username": "dependabot[bot]"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "4c94d3d30645210620bb036f477b4e33bd52693b",
          "message": "chore(deps): bump go.opentelemetry.io/otel from 1.39.0 to 1.41.0 (#7903)\n\nBumps [go.opentelemetry.io/otel](https://github.com/open-telemetry/opentelemetry-go) from 1.39.0 to 1.41.0.\n- [Release notes](https://github.com/open-telemetry/opentelemetry-go/releases)\n- [Changelog](https://github.com/open-telemetry/opentelemetry-go/blob/main/CHANGELOG.md)\n- [Commits](https://github.com/open-telemetry/opentelemetry-go/compare/v1.39.0...v1.41.0)\n\n---\nupdated-dependencies:\n- dependency-name: go.opentelemetry.io/otel\n  dependency-version: 1.41.0\n  dependency-type: indirect\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2026-06-09T08:00:58Z",
          "tree_id": "e51610c767c04844cb59cc9a56adc08442e7ea5a",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/4c94d3d30645210620bb036f477b4e33bd52693b"
        },
        "date": 1780994027062,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 576469,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "1814 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 576469,
            "unit": "ns/op",
            "extra": "1814 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819453,
            "unit": "B/op",
            "extra": "1814 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1814 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7125,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "165391 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7125,
            "unit": "ns/op",
            "extra": "165391 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "165391 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "165391 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 81.13,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14711348 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 81.13,
            "unit": "ns/op",
            "extra": "14711348 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14711348 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14711348 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 19531,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "61087 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 19531,
            "unit": "ns/op",
            "extra": "61087 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "61087 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "61087 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 221129,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4743 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 221129,
            "unit": "ns/op",
            "extra": "4743 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4743 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4743 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2461278,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "475 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2461278,
            "unit": "ns/op",
            "extra": "475 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "475 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "475 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 36038503,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "28 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 36038503,
            "unit": "ns/op",
            "extra": "28 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
            "unit": "B/op",
            "extra": "28 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "28 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 10249,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "128100 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10249,
            "unit": "ns/op",
            "extra": "128100 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "128100 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "128100 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7764487,
            "unit": "ns/op\t 4527544 B/op\t   69230 allocs/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7764487,
            "unit": "ns/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527544,
            "unit": "B/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7760,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "154952 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7760,
            "unit": "ns/op",
            "extra": "154952 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "154952 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "154952 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 846431,
            "unit": "ns/op\t  392930 B/op\t    6173 allocs/op",
            "extra": "1394 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 846431,
            "unit": "ns/op",
            "extra": "1394 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392930,
            "unit": "B/op",
            "extra": "1394 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1394 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11001,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "106185 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11001,
            "unit": "ns/op",
            "extra": "106185 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "106185 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "106185 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8000777,
            "unit": "ns/op\t 4914209 B/op\t   75241 allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8000777,
            "unit": "ns/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914209,
            "unit": "B/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 314.6,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3802402 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 314.6,
            "unit": "ns/op",
            "extra": "3802402 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3802402 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3802402 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.22,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.22,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 584.6,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2037120 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 584.6,
            "unit": "ns/op",
            "extra": "2037120 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2037120 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2037120 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 847.3,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1420125 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 847.3,
            "unit": "ns/op",
            "extra": "1420125 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1420125 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1420125 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "49699333+dependabot[bot]@users.noreply.github.com",
            "name": "dependabot[bot]",
            "username": "dependabot[bot]"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "58b32286e7a7eda638d14dc7e9eb5db3fbefc106",
          "message": "chore(deps): bump golang from 1.25.10 to 1.26.4 (#7949)\n\nBumps golang from 1.25.10 to 1.26.4.\n\n---\nupdated-dependencies:\n- dependency-name: golang\n  dependency-version: 1.26.4\n  dependency-type: direct:production\n  update-type: version-update:semver-minor\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2026-06-09T09:25:32Z",
          "tree_id": "e924148a6dacb67c20c495d9d13c918a4bd9e697",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/58b32286e7a7eda638d14dc7e9eb5db3fbefc106"
        },
        "date": 1780999095264,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 675343,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "1560 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 675343,
            "unit": "ns/op",
            "extra": "1560 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819453,
            "unit": "B/op",
            "extra": "1560 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1560 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7194,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "167176 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7194,
            "unit": "ns/op",
            "extra": "167176 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "167176 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "167176 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 81.24,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14790865 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 81.24,
            "unit": "ns/op",
            "extra": "14790865 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14790865 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14790865 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 19590,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "61527 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 19590,
            "unit": "ns/op",
            "extra": "61527 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "61527 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "61527 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 217640,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4974 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 217640,
            "unit": "ns/op",
            "extra": "4974 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4974 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4974 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2657491,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "388 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2657491,
            "unit": "ns/op",
            "extra": "388 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "388 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "388 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 35411603,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 35411603,
            "unit": "ns/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9245,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "129891 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9245,
            "unit": "ns/op",
            "extra": "129891 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "129891 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "129891 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6907020,
            "unit": "ns/op\t 4527548 B/op\t   69230 allocs/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6907020,
            "unit": "ns/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527548,
            "unit": "B/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7897,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "152476 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7897,
            "unit": "ns/op",
            "extra": "152476 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "152476 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "152476 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 850304,
            "unit": "ns/op\t  392943 B/op\t    6173 allocs/op",
            "extra": "1305 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 850304,
            "unit": "ns/op",
            "extra": "1305 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392943,
            "unit": "B/op",
            "extra": "1305 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1305 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11049,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "106390 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11049,
            "unit": "ns/op",
            "extra": "106390 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "106390 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "106390 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8193753,
            "unit": "ns/op\t 4914203 B/op\t   75241 allocs/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8193753,
            "unit": "ns/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914203,
            "unit": "B/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 313.1,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3810001 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 313.1,
            "unit": "ns/op",
            "extra": "3810001 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3810001 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3810001 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.23,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.23,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 576.5,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2070760 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 576.5,
            "unit": "ns/op",
            "extra": "2070760 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2070760 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2070760 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 814.5,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1473458 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 814.5,
            "unit": "ns/op",
            "extra": "1473458 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1473458 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1473458 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "pmalek@users.noreply.github.com",
            "name": "Patryk Małek",
            "username": "pmalek"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "eae8114a250c1290da0928c8dee2884dca23a059",
          "message": "chore: bump go to 1.25.11 (#7958)\n\n* chore: bump go to 1.25.11\n\n* Apply suggestions from code review\n\nCo-authored-by: Jakub Warczarek <jakub.warczarek@gmail.com>\n\n---------\n\nCo-authored-by: Jakub Warczarek <jakub.warczarek@gmail.com>",
          "timestamp": "2026-06-09T14:31:56+02:00",
          "tree_id": "56bb8a1ae8201af2fa5f6d05bd8f23378047d896",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/eae8114a250c1290da0928c8dee2884dca23a059"
        },
        "date": 1781008535345,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 662972,
            "unit": "ns/op\t  819430 B/op\t       4 allocs/op",
            "extra": "1530 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 662972,
            "unit": "ns/op",
            "extra": "1530 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819430,
            "unit": "B/op",
            "extra": "1530 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1530 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7484,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "157944 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7484,
            "unit": "ns/op",
            "extra": "157944 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "157944 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "157944 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.83,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16406882 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.83,
            "unit": "ns/op",
            "extra": "16406882 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16406882 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16406882 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22074,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54124 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22074,
            "unit": "ns/op",
            "extra": "54124 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54124 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54124 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 206975,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5073 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 206975,
            "unit": "ns/op",
            "extra": "5073 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5073 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5073 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2392896,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "465 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2392896,
            "unit": "ns/op",
            "extra": "465 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "465 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "465 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 42820066,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "34 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 42820066,
            "unit": "ns/op",
            "extra": "34 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
            "unit": "B/op",
            "extra": "34 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "34 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9470,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "126118 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9470,
            "unit": "ns/op",
            "extra": "126118 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "126118 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "126118 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7668076,
            "unit": "ns/op\t 4335499 B/op\t   63229 allocs/op",
            "extra": "166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7668076,
            "unit": "ns/op",
            "extra": "166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335499,
            "unit": "B/op",
            "extra": "166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8513,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "125676 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8513,
            "unit": "ns/op",
            "extra": "125676 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "125676 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "125676 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 843843,
            "unit": "ns/op\t  392927 B/op\t    6173 allocs/op",
            "extra": "1366 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 843843,
            "unit": "ns/op",
            "extra": "1366 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392927,
            "unit": "B/op",
            "extra": "1366 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1366 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10976,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "109347 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10976,
            "unit": "ns/op",
            "extra": "109347 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "109347 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "109347 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7553320,
            "unit": "ns/op\t 4722114 B/op\t   69239 allocs/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7553320,
            "unit": "ns/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722114,
            "unit": "B/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 215.5,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5633799 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 215.5,
            "unit": "ns/op",
            "extra": "5633799 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5633799 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5633799 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 9.983,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "120220407 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 9.983,
            "unit": "ns/op",
            "extra": "120220407 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "120220407 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "120220407 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 521.8,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2317873 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 521.8,
            "unit": "ns/op",
            "extra": "2317873 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2317873 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2317873 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 738.4,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1625126 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 738.4,
            "unit": "ns/op",
            "extra": "1625126 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1625126 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1625126 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "pmalek@users.noreply.github.com",
            "name": "Patryk Małek",
            "username": "pmalek"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "3f4d839afc2ff548076fc3a1a733cbf7fe1f3fec",
          "message": "ci: fix codecov (#7961)",
          "timestamp": "2026-06-09T15:25:59+02:00",
          "tree_id": "844a8a5482e5a73e57e887294928033a3cfca071",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/3f4d839afc2ff548076fc3a1a733cbf7fe1f3fec"
        },
        "date": 1781011786990,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 866966,
            "unit": "ns/op\t  819431 B/op\t       4 allocs/op",
            "extra": "1290 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 866966,
            "unit": "ns/op",
            "extra": "1290 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819431,
            "unit": "B/op",
            "extra": "1290 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1290 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7293,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "163413 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7293,
            "unit": "ns/op",
            "extra": "163413 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "163413 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "163413 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 81.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14737140 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 81.6,
            "unit": "ns/op",
            "extra": "14737140 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14737140 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14737140 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 20159,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "59776 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 20159,
            "unit": "ns/op",
            "extra": "59776 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "59776 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "59776 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 230965,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4680 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 230965,
            "unit": "ns/op",
            "extra": "4680 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4680 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4680 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2635873,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "453 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2635873,
            "unit": "ns/op",
            "extra": "453 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "453 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "453 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 37987281,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 37987281,
            "unit": "ns/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
            "unit": "B/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9778,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "122076 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9778,
            "unit": "ns/op",
            "extra": "122076 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "122076 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "122076 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7042360,
            "unit": "ns/op\t 4335504 B/op\t   63229 allocs/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7042360,
            "unit": "ns/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335504,
            "unit": "B/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8130,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "144184 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8130,
            "unit": "ns/op",
            "extra": "144184 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "144184 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "144184 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 860771,
            "unit": "ns/op\t  392920 B/op\t    6173 allocs/op",
            "extra": "1316 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 860771,
            "unit": "ns/op",
            "extra": "1316 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392920,
            "unit": "B/op",
            "extra": "1316 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1316 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11330,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "108793 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11330,
            "unit": "ns/op",
            "extra": "108793 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "108793 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "108793 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7887178,
            "unit": "ns/op\t 4722092 B/op\t   69239 allocs/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7887178,
            "unit": "ns/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722092,
            "unit": "B/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 205.5,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5857885 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 205.5,
            "unit": "ns/op",
            "extra": "5857885 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5857885 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5857885 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 9.502,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "126293478 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 9.502,
            "unit": "ns/op",
            "extra": "126293478 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "126293478 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "126293478 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 531.1,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2233749 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 531.1,
            "unit": "ns/op",
            "extra": "2233749 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2233749 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2233749 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 745.4,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1594026 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 745.4,
            "unit": "ns/op",
            "extra": "1594026 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1594026 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1594026 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "tao.yi@konghq.com",
            "name": "Tao Yi",
            "username": "randmonkey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "bf8ae01dfd69a14180a82ab50722dfe76238ab06",
          "message": "bump test dependencies to latest (#7935)",
          "timestamp": "2026-06-09T15:30:47+02:00",
          "tree_id": "672352d03e495acf27b490beecc9f633c785a028",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/bf8ae01dfd69a14180a82ab50722dfe76238ab06"
        },
        "date": 1781012080057,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 767170,
            "unit": "ns/op\t  819430 B/op\t       4 allocs/op",
            "extra": "1323 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 767170,
            "unit": "ns/op",
            "extra": "1323 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819430,
            "unit": "B/op",
            "extra": "1323 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1323 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7538,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "158617 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7538,
            "unit": "ns/op",
            "extra": "158617 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "158617 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "158617 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.37,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16572253 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.37,
            "unit": "ns/op",
            "extra": "16572253 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16572253 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16572253 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21708,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "55136 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21708,
            "unit": "ns/op",
            "extra": "55136 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "55136 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "55136 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 215195,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5241 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 215195,
            "unit": "ns/op",
            "extra": "5241 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5241 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5241 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2356499,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "508 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2356499,
            "unit": "ns/op",
            "extra": "508 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "508 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "508 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 38727058,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 38727058,
            "unit": "ns/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 11437,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "124621 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 11437,
            "unit": "ns/op",
            "extra": "124621 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "124621 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "124621 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6604091,
            "unit": "ns/op\t 4335506 B/op\t   63229 allocs/op",
            "extra": "181 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6604091,
            "unit": "ns/op",
            "extra": "181 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335506,
            "unit": "B/op",
            "extra": "181 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "181 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7954,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "150031 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7954,
            "unit": "ns/op",
            "extra": "150031 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "150031 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "150031 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 857431,
            "unit": "ns/op\t  392918 B/op\t    6173 allocs/op",
            "extra": "1333 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 857431,
            "unit": "ns/op",
            "extra": "1333 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392918,
            "unit": "B/op",
            "extra": "1333 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1333 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11227,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "108414 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11227,
            "unit": "ns/op",
            "extra": "108414 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "108414 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "108414 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7770320,
            "unit": "ns/op\t 4722122 B/op\t   69239 allocs/op",
            "extra": "153 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7770320,
            "unit": "ns/op",
            "extra": "153 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722122,
            "unit": "B/op",
            "extra": "153 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "153 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 215.3,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5555070 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 215.3,
            "unit": "ns/op",
            "extra": "5555070 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5555070 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5555070 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.04,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "119206760 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.04,
            "unit": "ns/op",
            "extra": "119206760 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "119206760 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "119206760 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 519.1,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2311071 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 519.1,
            "unit": "ns/op",
            "extra": "2311071 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2311071 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2311071 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 962.6,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1299528 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 962.6,
            "unit": "ns/op",
            "extra": "1299528 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1299528 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1299528 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "49699333+dependabot[bot]@users.noreply.github.com",
            "name": "dependabot[bot]",
            "username": "dependabot[bot]"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "969b6e2a2b3df27b20cf94f3271b0c48a5c7bc65",
          "message": "chore(deps): bump github/codeql-action from 3.29.7 to 4.31.2 (#7766)\n\nBumps [github/codeql-action](https://github.com/github/codeql-action) from 3.29.7 to 4.31.2.\n- [Release notes](https://github.com/github/codeql-action/releases)\n- [Changelog](https://github.com/github/codeql-action/blob/main/CHANGELOG.md)\n- [Commits](https://github.com/github/codeql-action/compare/51f77329afa6477de8c49fc9c7046c15b9a4e79d...8a06050a8c0348fb4738f28e0cfbb6727cf054ce)\n\n---\nupdated-dependencies:\n- dependency-name: github/codeql-action\n  dependency-version: 4.31.2\n  dependency-type: direct:production\n  update-type: version-update:semver-major\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2026-06-15T17:06:34+02:00",
          "tree_id": "8318a6a76c79d85ee42c6c3415417487cff80b0b",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/969b6e2a2b3df27b20cf94f3271b0c48a5c7bc65"
        },
        "date": 1781536229858,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 926442,
            "unit": "ns/op\t  819431 B/op\t       4 allocs/op",
            "extra": "1096 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 926442,
            "unit": "ns/op",
            "extra": "1096 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819431,
            "unit": "B/op",
            "extra": "1096 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1096 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 8457,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "132074 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 8457,
            "unit": "ns/op",
            "extra": "132074 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "132074 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "132074 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.47,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16562686 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.47,
            "unit": "ns/op",
            "extra": "16562686 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16562686 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16562686 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22445,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52542 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22445,
            "unit": "ns/op",
            "extra": "52542 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52542 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52542 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 237450,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4570 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 237450,
            "unit": "ns/op",
            "extra": "4570 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4570 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4570 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2456309,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "414 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2456309,
            "unit": "ns/op",
            "extra": "414 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "414 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "414 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 52439675,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "31 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 52439675,
            "unit": "ns/op",
            "extra": "31 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
            "unit": "B/op",
            "extra": "31 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "31 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 10790,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "119836 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10790,
            "unit": "ns/op",
            "extra": "119836 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "119836 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "119836 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 8019096,
            "unit": "ns/op\t 4335497 B/op\t   63229 allocs/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 8019096,
            "unit": "ns/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335497,
            "unit": "B/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8075,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "149132 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8075,
            "unit": "ns/op",
            "extra": "149132 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "149132 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "149132 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 864823,
            "unit": "ns/op\t  392921 B/op\t    6173 allocs/op",
            "extra": "1344 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 864823,
            "unit": "ns/op",
            "extra": "1344 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392921,
            "unit": "B/op",
            "extra": "1344 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1344 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11311,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "104228 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11311,
            "unit": "ns/op",
            "extra": "104228 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "104228 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "104228 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8030741,
            "unit": "ns/op\t 4722106 B/op\t   69239 allocs/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8030741,
            "unit": "ns/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722106,
            "unit": "B/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 214.8,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5612082 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 214.8,
            "unit": "ns/op",
            "extra": "5612082 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5612082 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5612082 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.02,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.02,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 511.8,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2331016 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 511.8,
            "unit": "ns/op",
            "extra": "2331016 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2331016 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2331016 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 737.9,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1602507 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 737.9,
            "unit": "ns/op",
            "extra": "1602507 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1602507 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1602507 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "49699333+dependabot[bot]@users.noreply.github.com",
            "name": "dependabot[bot]",
            "username": "dependabot[bot]"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "45363045750528505271e64c3394bf1afa4522be",
          "message": "chore(deps): bump golang from `68cb6d6` to `11fd8f7` (#7967)\n\nBumps golang from `68cb6d6` to `11fd8f7`.\n\n---\nupdated-dependencies:\n- dependency-name: golang\n  dependency-version: 1.26.4\n  dependency-type: direct:production\n  update-type: version-update:semver-patch\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2026-06-15T17:06:22+02:00",
          "tree_id": "7e842170cbf03866ebcd59f6e30098da270a1ff6",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/45363045750528505271e64c3394bf1afa4522be"
        },
        "date": 1781536233126,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 942437,
            "unit": "ns/op\t  819431 B/op\t       4 allocs/op",
            "extra": "1128 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 942437,
            "unit": "ns/op",
            "extra": "1128 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819431,
            "unit": "B/op",
            "extra": "1128 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1128 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 9536,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "124549 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 9536,
            "unit": "ns/op",
            "extra": "124549 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "124549 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "124549 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.78,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15563355 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.78,
            "unit": "ns/op",
            "extra": "15563355 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15563355 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15563355 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22353,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53192 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22353,
            "unit": "ns/op",
            "extra": "53192 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53192 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53192 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 230503,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4774 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 230503,
            "unit": "ns/op",
            "extra": "4774 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4774 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4774 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2424209,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "442 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2424209,
            "unit": "ns/op",
            "extra": "442 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "442 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "442 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 41266173,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "78 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 41266173,
            "unit": "ns/op",
            "extra": "78 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
            "unit": "B/op",
            "extra": "78 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "78 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9848,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "121828 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9848,
            "unit": "ns/op",
            "extra": "121828 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "121828 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "121828 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7021533,
            "unit": "ns/op\t 4335505 B/op\t   63229 allocs/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7021533,
            "unit": "ns/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335505,
            "unit": "B/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8110,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "147962 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8110,
            "unit": "ns/op",
            "extra": "147962 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "147962 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "147962 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 867985,
            "unit": "ns/op\t  392912 B/op\t    6173 allocs/op",
            "extra": "1341 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 867985,
            "unit": "ns/op",
            "extra": "1341 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392912,
            "unit": "B/op",
            "extra": "1341 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1341 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11536,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "105984 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11536,
            "unit": "ns/op",
            "extra": "105984 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "105984 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "105984 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7977494,
            "unit": "ns/op\t 4722130 B/op\t   69239 allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7977494,
            "unit": "ns/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722130,
            "unit": "B/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 217,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5524056 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 217,
            "unit": "ns/op",
            "extra": "5524056 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5524056 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5524056 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.03,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "120140834 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.03,
            "unit": "ns/op",
            "extra": "120140834 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "120140834 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "120140834 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 734.6,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "1659486 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 734.6,
            "unit": "ns/op",
            "extra": "1659486 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "1659486 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1659486 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 846,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1359644 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 846,
            "unit": "ns/op",
            "extra": "1359644 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1359644 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1359644 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "pmalek@users.noreply.github.com",
            "name": "Patryk Małek",
            "username": "pmalek"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "2af62b1bee9fcb1d9cf28619f940472490927491",
          "message": "fix: fix credential recreation on every Konnect sync (#7978)\n\n* fix: fix credential recreation on every Konnect sync\n\n* Update CHANGELOG.md\n\nCo-authored-by: Jintao Zhang <zhangjintao9020@gmail.com>\n\n* chore: rename credentialOps fields\n\n* chore: comment for credential interface\n\n---------\n\nCo-authored-by: Jintao Zhang <zhangjintao9020@gmail.com>",
          "timestamp": "2026-06-16T11:14:03+02:00",
          "tree_id": "7eec616d54b15c7c1af512295e53e8ddfbfdb3f9",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/2af62b1bee9fcb1d9cf28619f940472490927491"
        },
        "date": 1781601477605,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 841384,
            "unit": "ns/op\t  819431 B/op\t       4 allocs/op",
            "extra": "1200 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 841384,
            "unit": "ns/op",
            "extra": "1200 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819431,
            "unit": "B/op",
            "extra": "1200 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1200 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7780,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "152845 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7780,
            "unit": "ns/op",
            "extra": "152845 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "152845 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "152845 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.47,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16600698 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.47,
            "unit": "ns/op",
            "extra": "16600698 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16600698 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16600698 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21849,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54692 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21849,
            "unit": "ns/op",
            "extra": "54692 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54692 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54692 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 226350,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4648 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 226350,
            "unit": "ns/op",
            "extra": "4648 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4648 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4648 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2451168,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "469 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2451168,
            "unit": "ns/op",
            "extra": "469 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "469 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "469 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 43724776,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "76 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 43724776,
            "unit": "ns/op",
            "extra": "76 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
            "unit": "B/op",
            "extra": "76 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "76 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9333,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "129861 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9333,
            "unit": "ns/op",
            "extra": "129861 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "129861 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "129861 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6647092,
            "unit": "ns/op\t 4335500 B/op\t   63229 allocs/op",
            "extra": "180 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6647092,
            "unit": "ns/op",
            "extra": "180 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335500,
            "unit": "B/op",
            "extra": "180 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "180 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7663,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "156794 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7663,
            "unit": "ns/op",
            "extra": "156794 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "156794 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "156794 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 840017,
            "unit": "ns/op\t  392912 B/op\t    6173 allocs/op",
            "extra": "1353 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 840017,
            "unit": "ns/op",
            "extra": "1353 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392912,
            "unit": "B/op",
            "extra": "1353 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1353 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10658,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "108630 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10658,
            "unit": "ns/op",
            "extra": "108630 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "108630 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "108630 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7602286,
            "unit": "ns/op\t 4722093 B/op\t   69239 allocs/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7602286,
            "unit": "ns/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722093,
            "unit": "B/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 263.2,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5621702 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 263.2,
            "unit": "ns/op",
            "extra": "5621702 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5621702 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5621702 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.21,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.21,
            "unit": "ns/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 511.3,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2349367 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 511.3,
            "unit": "ns/op",
            "extra": "2349367 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2349367 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2349367 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 730.8,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1638734 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 730.8,
            "unit": "ns/op",
            "extra": "1638734 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1638734 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1638734 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "tao.yi@konghq.com",
            "name": "Tao Yi",
            "username": "randmonkey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "c88a0911a1eddb6f79ce68f3cd083b787f297cbb",
          "message": "bump go.opentelemetry.io/otel/sdk to 1.43.0 (#7976)",
          "timestamp": "2026-06-16T09:49:07Z",
          "tree_id": "55bac03906960560c02fe8a6c0003286ddf29030",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/c88a0911a1eddb6f79ce68f3cd083b787f297cbb"
        },
        "date": 1781606949827,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 822056,
            "unit": "ns/op\t  819430 B/op\t       4 allocs/op",
            "extra": "1250 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 822056,
            "unit": "ns/op",
            "extra": "1250 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819430,
            "unit": "B/op",
            "extra": "1250 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1250 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 8284,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "140712 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 8284,
            "unit": "ns/op",
            "extra": "140712 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "140712 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "140712 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 65.78,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "17960278 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 65.78,
            "unit": "ns/op",
            "extra": "17960278 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "17960278 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "17960278 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21326,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "55785 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21326,
            "unit": "ns/op",
            "extra": "55785 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "55785 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "55785 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 238918,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4406 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 238918,
            "unit": "ns/op",
            "extra": "4406 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4406 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4406 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2713621,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "440 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2713621,
            "unit": "ns/op",
            "extra": "440 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "440 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "440 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 43166156,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "63 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 43166156,
            "unit": "ns/op",
            "extra": "63 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
            "unit": "B/op",
            "extra": "63 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "63 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9938,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "121903 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9938,
            "unit": "ns/op",
            "extra": "121903 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "121903 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "121903 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7001110,
            "unit": "ns/op\t 4335507 B/op\t   63229 allocs/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7001110,
            "unit": "ns/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335507,
            "unit": "B/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8309,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "142594 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8309,
            "unit": "ns/op",
            "extra": "142594 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "142594 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "142594 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 934955,
            "unit": "ns/op\t  392905 B/op\t    6173 allocs/op",
            "extra": "1286 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 934955,
            "unit": "ns/op",
            "extra": "1286 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392905,
            "unit": "B/op",
            "extra": "1286 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1286 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11552,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "104397 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11552,
            "unit": "ns/op",
            "extra": "104397 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "104397 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "104397 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8081020,
            "unit": "ns/op\t 4722121 B/op\t   69239 allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8081020,
            "unit": "ns/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722121,
            "unit": "B/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 213,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5628187 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 213,
            "unit": "ns/op",
            "extra": "5628187 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5628187 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5628187 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 9.808,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "123161514 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 9.808,
            "unit": "ns/op",
            "extra": "123161514 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "123161514 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "123161514 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 716.6,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "1714477 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 716.6,
            "unit": "ns/op",
            "extra": "1714477 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "1714477 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1714477 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 806.7,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1422140 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 806.7,
            "unit": "ns/op",
            "extra": "1422140 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1422140 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1422140 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "tao.yi@konghq.com",
            "name": "Tao Yi",
            "username": "randmonkey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "e1f58a00c122ad47603f16d754856377f4ebc247",
          "message": "fix data race by using Eventually (#7979)",
          "timestamp": "2026-06-17T06:43:18Z",
          "tree_id": "657a17064574e6554600671a995dbabc8aa9dc94",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/e1f58a00c122ad47603f16d754856377f4ebc247"
        },
        "date": 1781680606065,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 868078,
            "unit": "ns/op\t  819431 B/op\t       4 allocs/op",
            "extra": "1214 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 868078,
            "unit": "ns/op",
            "extra": "1214 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819431,
            "unit": "B/op",
            "extra": "1214 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1214 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7081,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "142969 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7081,
            "unit": "ns/op",
            "extra": "142969 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "142969 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "142969 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.47,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16557373 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.47,
            "unit": "ns/op",
            "extra": "16557373 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16557373 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16557373 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22247,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53250 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22247,
            "unit": "ns/op",
            "extra": "53250 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53250 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53250 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 256348,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5000 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 256348,
            "unit": "ns/op",
            "extra": "5000 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5000 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5000 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2500322,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "418 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2500322,
            "unit": "ns/op",
            "extra": "418 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "418 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "418 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 29563035,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "44 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 29563035,
            "unit": "ns/op",
            "extra": "44 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
            "unit": "B/op",
            "extra": "44 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "44 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9103,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "131257 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9103,
            "unit": "ns/op",
            "extra": "131257 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "131257 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "131257 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6494480,
            "unit": "ns/op\t 4335500 B/op\t   63229 allocs/op",
            "extra": "184 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6494480,
            "unit": "ns/op",
            "extra": "184 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335500,
            "unit": "B/op",
            "extra": "184 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "184 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7579,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "158169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7579,
            "unit": "ns/op",
            "extra": "158169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "158169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "158169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 841393,
            "unit": "ns/op\t  392920 B/op\t    6173 allocs/op",
            "extra": "1384 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 841393,
            "unit": "ns/op",
            "extra": "1384 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392920,
            "unit": "B/op",
            "extra": "1384 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1384 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10545,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "112630 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10545,
            "unit": "ns/op",
            "extra": "112630 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "112630 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "112630 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7456656,
            "unit": "ns/op\t 4722093 B/op\t   69239 allocs/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7456656,
            "unit": "ns/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722093,
            "unit": "B/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 215.9,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5575180 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 215.9,
            "unit": "ns/op",
            "extra": "5575180 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5575180 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5575180 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 9.98,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "120222590 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 9.98,
            "unit": "ns/op",
            "extra": "120222590 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "120222590 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "120222590 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 525.7,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2288360 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 525.7,
            "unit": "ns/op",
            "extra": "2288360 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2288360 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2288360 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 743.9,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1609356 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 743.9,
            "unit": "ns/op",
            "extra": "1609356 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1609356 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1609356 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "pmalek@users.noreply.github.com",
            "name": "Patryk Małek",
            "username": "pmalek"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "287127f9f14196e2dffc6d9ebb75c8c4c8189b4d",
          "message": "fix: fix missing certificate tags (#7991)\n\n* fix: fix missing certificate tags\n\n* Apply suggestion from @programmer04\n\nCo-authored-by: Jakub Warczarek <jakub.warczarek@konghq.com>\n\n* Apply suggestion from @programmer04\n\nCo-authored-by: Jakub Warczarek <jakub.warczarek@konghq.com>\n\n* chore: fix linter issues\n\n---------\n\nCo-authored-by: Jakub Warczarek <jakub.warczarek@konghq.com>",
          "timestamp": "2026-06-18T12:36:13+02:00",
          "tree_id": "4211bf05c0e5cf04d955497cd7e8fd992c3567e5",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/287127f9f14196e2dffc6d9ebb75c8c4c8189b4d"
        },
        "date": 1781779188225,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 789512,
            "unit": "ns/op\t  819431 B/op\t       4 allocs/op",
            "extra": "1291 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 789512,
            "unit": "ns/op",
            "extra": "1291 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819431,
            "unit": "B/op",
            "extra": "1291 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1291 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7468,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "153979 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7468,
            "unit": "ns/op",
            "extra": "153979 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "153979 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "153979 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 82.52,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13733943 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 82.52,
            "unit": "ns/op",
            "extra": "13733943 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13733943 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13733943 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 20158,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "59455 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 20158,
            "unit": "ns/op",
            "extra": "59455 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "59455 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "59455 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 205201,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5130 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 205201,
            "unit": "ns/op",
            "extra": "5130 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5130 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5130 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2551682,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "477 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2551682,
            "unit": "ns/op",
            "extra": "477 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "477 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "477 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 49120182,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "33 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 49120182,
            "unit": "ns/op",
            "extra": "33 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
            "unit": "B/op",
            "extra": "33 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "33 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9927,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "132980 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9927,
            "unit": "ns/op",
            "extra": "132980 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "132980 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "132980 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7247740,
            "unit": "ns/op\t 4335497 B/op\t   63229 allocs/op",
            "extra": "153 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7247740,
            "unit": "ns/op",
            "extra": "153 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335497,
            "unit": "B/op",
            "extra": "153 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "153 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7396,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "158708 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7396,
            "unit": "ns/op",
            "extra": "158708 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "158708 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "158708 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 823105,
            "unit": "ns/op\t  392911 B/op\t    6173 allocs/op",
            "extra": "1401 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 823105,
            "unit": "ns/op",
            "extra": "1401 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392911,
            "unit": "B/op",
            "extra": "1401 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1401 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10437,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "114740 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10437,
            "unit": "ns/op",
            "extra": "114740 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "114740 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "114740 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7235549,
            "unit": "ns/op\t 4722109 B/op\t   69239 allocs/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7235549,
            "unit": "ns/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722109,
            "unit": "B/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 205.8,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5889162 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 205.8,
            "unit": "ns/op",
            "extra": "5889162 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5889162 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5889162 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 9.546,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "125594997 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 9.546,
            "unit": "ns/op",
            "extra": "125594997 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "125594997 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "125594997 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 528,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2270497 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 528,
            "unit": "ns/op",
            "extra": "2270497 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2270497 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2270497 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 736.6,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1649808 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 736.6,
            "unit": "ns/op",
            "extra": "1649808 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1649808 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1649808 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "tao.yi@konghq.com",
            "name": "Tao Yi",
            "username": "randmonkey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "53e10714663b31ddce75e7a5a3c14cfecc515104",
          "message": "force licenses if Kong version >= 3.15 (#7941)\n\n* force licenses if Kong version >= 3.15\n\n* support 4 digit version\n\n* add patch for license for all-in-one tests\n\n* add deployment name var for patch\n\n* add license check for all e2e\n\n* add konnect licensing patch and set workspace for plugin test\n\n* change tab to space in patch",
          "timestamp": "2026-06-22T07:18:36Z",
          "tree_id": "f0bd2de187dd18ebca59fda18f4f19d93e840a5b",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/53e10714663b31ddce75e7a5a3c14cfecc515104"
        },
        "date": 1782114685686,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 657809,
            "unit": "ns/op\t  819431 B/op\t       4 allocs/op",
            "extra": "1567 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 657809,
            "unit": "ns/op",
            "extra": "1567 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819431,
            "unit": "B/op",
            "extra": "1567 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1567 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7308,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "169989 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7308,
            "unit": "ns/op",
            "extra": "169989 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "169989 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "169989 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 73.99,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16468222 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 73.99,
            "unit": "ns/op",
            "extra": "16468222 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16468222 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16468222 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22057,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54284 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22057,
            "unit": "ns/op",
            "extra": "54284 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54284 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54284 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 214488,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5076 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 214488,
            "unit": "ns/op",
            "extra": "5076 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5076 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5076 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2310475,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "507 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2310475,
            "unit": "ns/op",
            "extra": "507 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "507 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "507 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 36919726,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "31 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 36919726,
            "unit": "ns/op",
            "extra": "31 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
            "unit": "B/op",
            "extra": "31 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "31 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 10766,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "124111 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10766,
            "unit": "ns/op",
            "extra": "124111 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "124111 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "124111 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7684655,
            "unit": "ns/op\t 4335503 B/op\t   63229 allocs/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7684655,
            "unit": "ns/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335503,
            "unit": "B/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7973,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "148770 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7973,
            "unit": "ns/op",
            "extra": "148770 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "148770 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "148770 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 843538,
            "unit": "ns/op\t  392919 B/op\t    6173 allocs/op",
            "extra": "1368 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 843538,
            "unit": "ns/op",
            "extra": "1368 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392919,
            "unit": "B/op",
            "extra": "1368 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1368 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11026,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "109218 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11026,
            "unit": "ns/op",
            "extra": "109218 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "109218 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "109218 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7543327,
            "unit": "ns/op\t 4722123 B/op\t   69239 allocs/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7543327,
            "unit": "ns/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722123,
            "unit": "B/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 216.4,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5561911 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 216.4,
            "unit": "ns/op",
            "extra": "5561911 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5561911 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5561911 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 9.971,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "120338893 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 9.971,
            "unit": "ns/op",
            "extra": "120338893 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "120338893 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "120338893 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 524.4,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2290292 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 524.4,
            "unit": "ns/op",
            "extra": "2290292 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2290292 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2290292 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 767.2,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1610287 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 767.2,
            "unit": "ns/op",
            "extra": "1610287 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1610287 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1610287 times\n4 procs"
          }
        ]
      }
    ]
  }
}