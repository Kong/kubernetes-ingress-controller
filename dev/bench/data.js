window.BENCHMARK_DATA = {
  "lastUpdate": 1753431988603,
  "repoUrl": "https://github.com/Kong/kubernetes-ingress-controller",
  "entries": {
    "Go Benchmark": [
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
          "distinct": false,
          "id": "66781c956a47a27f37721b448e2012cda84cf6ff",
          "message": "chore(deps): update dependency go-delve/delve to v1.25.0 (#7507)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2025-06-30T09:04:21Z",
          "tree_id": "ee70da0ce65ac53b2a8906e033e3d7672cfde176",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/66781c956a47a27f37721b448e2012cda84cf6ff"
        },
        "date": 1751276352480,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 389405,
            "unit": "ns/op\t  819450 B/op\t       5 allocs/op",
            "extra": "2660 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 389405,
            "unit": "ns/op",
            "extra": "2660 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819450,
            "unit": "B/op",
            "extra": "2660 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "2660 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7642,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "154039 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7642,
            "unit": "ns/op",
            "extra": "154039 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "154039 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "154039 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 105.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11401461 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 105.2,
            "unit": "ns/op",
            "extra": "11401461 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11401461 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11401461 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21906,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53163 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21906,
            "unit": "ns/op",
            "extra": "53163 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53163 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53163 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 218986,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5486 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 218986,
            "unit": "ns/op",
            "extra": "5486 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5486 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5486 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2358445,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "502 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2358445,
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
            "value": 35670961,
            "unit": "ns/op\t24010754 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 35670961,
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
            "value": 9243,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "130767 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9243,
            "unit": "ns/op",
            "extra": "130767 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "130767 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "130767 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7093192,
            "unit": "ns/op\t 4527310 B/op\t   69224 allocs/op",
            "extra": "165 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7093192,
            "unit": "ns/op",
            "extra": "165 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527310,
            "unit": "B/op",
            "extra": "165 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "165 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7745,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "149880 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7745,
            "unit": "ns/op",
            "extra": "149880 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "149880 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "149880 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 876545,
            "unit": "ns/op\t  396619 B/op\t    6226 allocs/op",
            "extra": "1262 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 876545,
            "unit": "ns/op",
            "extra": "1262 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396619,
            "unit": "B/op",
            "extra": "1262 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6226,
            "unit": "allocs/op",
            "extra": "1262 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10847,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "110128 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10847,
            "unit": "ns/op",
            "extra": "110128 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "110128 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "110128 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8196470,
            "unit": "ns/op\t 4914013 B/op\t   75235 allocs/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8196470,
            "unit": "ns/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914013,
            "unit": "B/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 324.4,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3694686 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 324.4,
            "unit": "ns/op",
            "extra": "3694686 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3694686 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3694686 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6226,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6226,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 686.8,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2191153 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 686.8,
            "unit": "ns/op",
            "extra": "2191153 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2191153 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2191153 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 777.1,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1546605 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 777.1,
            "unit": "ns/op",
            "extra": "1546605 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1546605 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1546605 times\n4 procs"
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
          "distinct": false,
          "id": "d538f2efd5942330bc5e8d683f29fbd6efd03d83",
          "message": "chore(deps): update dependency kubernetes-sigs/kustomize to v5.7.0 (#7561)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>\nCo-authored-by: Jakub Warczarek <jakub.warczarek@konghq.com>",
          "timestamp": "2025-06-30T09:09:36Z",
          "tree_id": "ca78acb26093d09285e204123a1d497c3520c90b",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/d538f2efd5942330bc5e8d683f29fbd6efd03d83"
        },
        "date": 1751276425695,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 328973,
            "unit": "ns/op\t  819452 B/op\t       5 allocs/op",
            "extra": "4047 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 328973,
            "unit": "ns/op",
            "extra": "4047 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819452,
            "unit": "B/op",
            "extra": "4047 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "4047 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 8173,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "132068 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 8173,
            "unit": "ns/op",
            "extra": "132068 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "132068 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "132068 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 104.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11475487 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 104.7,
            "unit": "ns/op",
            "extra": "11475487 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11475487 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11475487 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22379,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53774 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22379,
            "unit": "ns/op",
            "extra": "53774 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53774 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53774 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 226876,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5160 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 226876,
            "unit": "ns/op",
            "extra": "5160 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5160 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5160 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2396947,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "427 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2396947,
            "unit": "ns/op",
            "extra": "427 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "427 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "427 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 43212211,
            "unit": "ns/op\t24010754 B/op\t       2 allocs/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 43212211,
            "unit": "ns/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010754,
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
            "value": 9543,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "127311 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9543,
            "unit": "ns/op",
            "extra": "127311 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "127311 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "127311 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7047958,
            "unit": "ns/op\t 4527309 B/op\t   69224 allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7047958,
            "unit": "ns/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527309,
            "unit": "B/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7772,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "151306 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7772,
            "unit": "ns/op",
            "extra": "151306 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "151306 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "151306 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 854351,
            "unit": "ns/op\t  396549 B/op\t    6226 allocs/op",
            "extra": "1275 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 854351,
            "unit": "ns/op",
            "extra": "1275 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396549,
            "unit": "B/op",
            "extra": "1275 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6226,
            "unit": "allocs/op",
            "extra": "1275 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10750,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "111046 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10750,
            "unit": "ns/op",
            "extra": "111046 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "111046 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "111046 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8035565,
            "unit": "ns/op\t 4913999 B/op\t   75235 allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8035565,
            "unit": "ns/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913999,
            "unit": "B/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 393.7,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3750075 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 393.7,
            "unit": "ns/op",
            "extra": "3750075 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3750075 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3750075 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6238,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6238,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 546.1,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2186395 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 546.1,
            "unit": "ns/op",
            "extra": "2186395 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2186395 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2186395 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 775.7,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1546731 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 775.7,
            "unit": "ns/op",
            "extra": "1546731 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1546731 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1546731 times\n4 procs"
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
          "id": "dd6538bf66f2e27805786fd6a21a7a93475f0323",
          "message": "feat(translator): Generate KongConfigurationTranslationFailed event when algorithm=sticky-sessions and kong version < 3.11 (#7538)\n\n* generate KongConfigurationTranslationFailed event when algorithm=sticky-sessions and kong version < 3.11\n\n* bump kubernetes-configuration, add envtest and minor adjustments\n\n---------\n\nCo-authored-by: Jakub Warczarek <jakub.warczarek@konghq.com>",
          "timestamp": "2025-06-30T09:11:00Z",
          "tree_id": "9f52604c18753794edb33bbe9234baeb94bc93a2",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/dd6538bf66f2e27805786fd6a21a7a93475f0323"
        },
        "date": 1751276494623,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 309442,
            "unit": "ns/op\t  819455 B/op\t       5 allocs/op",
            "extra": "3876 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 309442,
            "unit": "ns/op",
            "extra": "3876 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819455,
            "unit": "B/op",
            "extra": "3876 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "3876 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7514,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "161760 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7514,
            "unit": "ns/op",
            "extra": "161760 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "161760 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "161760 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 104.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11483790 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 104.7,
            "unit": "ns/op",
            "extra": "11483790 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11483790 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11483790 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21654,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "55200 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21654,
            "unit": "ns/op",
            "extra": "55200 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "55200 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "55200 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 213043,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5296 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 213043,
            "unit": "ns/op",
            "extra": "5296 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5296 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5296 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2375708,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "460 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2375708,
            "unit": "ns/op",
            "extra": "460 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "460 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "460 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 38610806,
            "unit": "ns/op\t24010754 B/op\t       2 allocs/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 38610806,
            "unit": "ns/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010754,
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
            "value": 9480,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "129244 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9480,
            "unit": "ns/op",
            "extra": "129244 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "129244 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "129244 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6991440,
            "unit": "ns/op\t 4527312 B/op\t   69224 allocs/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6991440,
            "unit": "ns/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527312,
            "unit": "B/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7835,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "154302 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7835,
            "unit": "ns/op",
            "extra": "154302 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "154302 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "154302 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 854229,
            "unit": "ns/op\t  396520 B/op\t    6225 allocs/op",
            "extra": "1288 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 854229,
            "unit": "ns/op",
            "extra": "1288 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396520,
            "unit": "B/op",
            "extra": "1288 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6225,
            "unit": "allocs/op",
            "extra": "1288 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10918,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "108798 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10918,
            "unit": "ns/op",
            "extra": "108798 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "108798 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "108798 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8025087,
            "unit": "ns/op\t 4913944 B/op\t   75235 allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8025087,
            "unit": "ns/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913944,
            "unit": "B/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 314.6,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3807343 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 314.6,
            "unit": "ns/op",
            "extra": "3807343 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3807343 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3807343 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6219,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6219,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 599.9,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2234940 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 599.9,
            "unit": "ns/op",
            "extra": "2234940 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2234940 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2234940 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 771.2,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1349367 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 771.2,
            "unit": "ns/op",
            "extra": "1349367 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1349367 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1349367 times\n4 procs"
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
          "id": "0702a04cad00cd7e86429919abf4ba9d4b595368",
          "message": "chore(deps): bump golang from `dc3de88` to `10c1318` (#7497)\n\nBumps golang from `dc3de88` to `10c1318`.\n\n---\nupdated-dependencies:\n- dependency-name: golang\n  dependency-version: 1.24.4\n  dependency-type: direct:production\n  update-type: version-update:semver-patch\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2025-06-30T10:04:51Z",
          "tree_id": "482e9144cc96cfae04a750910e27427fda9af98f",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/0702a04cad00cd7e86429919abf4ba9d4b595368"
        },
        "date": 1751279686583,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 291809,
            "unit": "ns/op\t  819452 B/op\t       5 allocs/op",
            "extra": "4132 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 291809,
            "unit": "ns/op",
            "extra": "4132 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819452,
            "unit": "B/op",
            "extra": "4132 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "4132 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 6986,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "171178 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 6986,
            "unit": "ns/op",
            "extra": "171178 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "171178 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "171178 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 104.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11458549 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 104.8,
            "unit": "ns/op",
            "extra": "11458549 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11458549 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11458549 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 24521,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54724 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 24521,
            "unit": "ns/op",
            "extra": "54724 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54724 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54724 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 209172,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5755 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 209172,
            "unit": "ns/op",
            "extra": "5755 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5755 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5755 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2447456,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "486 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2447456,
            "unit": "ns/op",
            "extra": "486 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "486 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "486 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 33503473,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 33503473,
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
            "value": 9353,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "127768 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9353,
            "unit": "ns/op",
            "extra": "127768 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "127768 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "127768 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7024197,
            "unit": "ns/op\t 4527312 B/op\t   69224 allocs/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7024197,
            "unit": "ns/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527312,
            "unit": "B/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7721,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "155068 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7721,
            "unit": "ns/op",
            "extra": "155068 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "155068 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "155068 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 859455,
            "unit": "ns/op\t  396236 B/op\t    6221 allocs/op",
            "extra": "1396 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 859455,
            "unit": "ns/op",
            "extra": "1396 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396236,
            "unit": "B/op",
            "extra": "1396 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6221,
            "unit": "allocs/op",
            "extra": "1396 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10833,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "110868 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10833,
            "unit": "ns/op",
            "extra": "110868 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "110868 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "110868 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8125235,
            "unit": "ns/op\t 4914106 B/op\t   75235 allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8125235,
            "unit": "ns/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914106,
            "unit": "B/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 320,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3795584 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 320,
            "unit": "ns/op",
            "extra": "3795584 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3795584 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3795584 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6263,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6263,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 543.3,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2219469 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 543.3,
            "unit": "ns/op",
            "extra": "2219469 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2219469 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2219469 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 915.7,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1545122 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 915.7,
            "unit": "ns/op",
            "extra": "1545122 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1545122 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1545122 times\n4 procs"
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
          "id": "d6f30920bd17e80a19cc9bd0a04f71e7e1c9998c",
          "message": "fix(konnect): Do not mark pod ready when no konnect licenses available (#7520)\n\n* Do not mark pod ready when no konnect licenses available\n\n* update changelog",
          "timestamp": "2025-07-01T02:20:22Z",
          "tree_id": "64052d5ca49150cb902f2eb9e8cfd8c7fedc3428",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/d6f30920bd17e80a19cc9bd0a04f71e7e1c9998c"
        },
        "date": 1751338216781,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 349933,
            "unit": "ns/op\t  819452 B/op\t       5 allocs/op",
            "extra": "3432 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 349933,
            "unit": "ns/op",
            "extra": "3432 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819452,
            "unit": "B/op",
            "extra": "3432 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "3432 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7527,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "157592 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7527,
            "unit": "ns/op",
            "extra": "157592 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "157592 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "157592 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 105.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11456581 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 105.2,
            "unit": "ns/op",
            "extra": "11456581 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11456581 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11456581 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22837,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54183 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22837,
            "unit": "ns/op",
            "extra": "54183 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54183 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54183 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 276686,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4347 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 276686,
            "unit": "ns/op",
            "extra": "4347 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4347 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4347 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2467936,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "428 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2467936,
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
            "value": 44923655,
            "unit": "ns/op\t24010754 B/op\t       2 allocs/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 44923655,
            "unit": "ns/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010754,
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
            "value": 9618,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "123438 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9618,
            "unit": "ns/op",
            "extra": "123438 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "123438 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "123438 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7864210,
            "unit": "ns/op\t 4527308 B/op\t   69224 allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7864210,
            "unit": "ns/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527308,
            "unit": "B/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7978,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "149053 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7978,
            "unit": "ns/op",
            "extra": "149053 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "149053 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "149053 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 871316,
            "unit": "ns/op\t  396795 B/op\t    6229 allocs/op",
            "extra": "1201 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 871316,
            "unit": "ns/op",
            "extra": "1201 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396795,
            "unit": "B/op",
            "extra": "1201 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6229,
            "unit": "allocs/op",
            "extra": "1201 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11244,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "103849 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11244,
            "unit": "ns/op",
            "extra": "103849 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "103849 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "103849 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8546951,
            "unit": "ns/op\t 4913975 B/op\t   75235 allocs/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8546951,
            "unit": "ns/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913975,
            "unit": "B/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 316.6,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3774387 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 316.6,
            "unit": "ns/op",
            "extra": "3774387 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3774387 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3774387 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6239,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6239,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 543.3,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2189498 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 543.3,
            "unit": "ns/op",
            "extra": "2189498 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2189498 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2189498 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 776.9,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1557122 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 776.9,
            "unit": "ns/op",
            "extra": "1557122 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1557122 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1557122 times\n4 procs"
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
          "id": "2f2a7c52e99a05b4c9732ed108dfb8427c42b454",
          "message": "chore(deps): bump golang from `10c1318` to `1bb140b` (#7570)\n\nBumps golang from `10c1318` to `1bb140b`.\n\n---\nupdated-dependencies:\n- dependency-name: golang\n  dependency-version: 1.24.4\n  dependency-type: direct:production\n  update-type: version-update:semver-patch\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2025-07-01T15:00:49Z",
          "tree_id": "5bdbfeeda5fb510896ebaf76b34018430cd719ef",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/2f2a7c52e99a05b4c9732ed108dfb8427c42b454"
        },
        "date": 1751383835848,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 290690,
            "unit": "ns/op\t  819452 B/op\t       5 allocs/op",
            "extra": "4119 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 290690,
            "unit": "ns/op",
            "extra": "4119 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819452,
            "unit": "B/op",
            "extra": "4119 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "4119 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 6994,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "169976 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 6994,
            "unit": "ns/op",
            "extra": "169976 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "169976 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "169976 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 104.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11475148 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 104.4,
            "unit": "ns/op",
            "extra": "11475148 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11475148 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11475148 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22112,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53827 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22112,
            "unit": "ns/op",
            "extra": "53827 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53827 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53827 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 263032,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5583 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 263032,
            "unit": "ns/op",
            "extra": "5583 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5583 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5583 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2428980,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "492 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2428980,
            "unit": "ns/op",
            "extra": "492 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "492 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "492 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 38180234,
            "unit": "ns/op\t24010754 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 38180234,
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
            "value": 9500,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "124848 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9500,
            "unit": "ns/op",
            "extra": "124848 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "124848 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "124848 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7241238,
            "unit": "ns/op\t 4527314 B/op\t   69224 allocs/op",
            "extra": "165 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7241238,
            "unit": "ns/op",
            "extra": "165 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527314,
            "unit": "B/op",
            "extra": "165 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "165 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7949,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "149800 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7949,
            "unit": "ns/op",
            "extra": "149800 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "149800 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "149800 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 860870,
            "unit": "ns/op\t  396586 B/op\t    6226 allocs/op",
            "extra": "1266 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 860870,
            "unit": "ns/op",
            "extra": "1266 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396586,
            "unit": "B/op",
            "extra": "1266 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6226,
            "unit": "allocs/op",
            "extra": "1266 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10942,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "107283 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10942,
            "unit": "ns/op",
            "extra": "107283 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "107283 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "107283 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8251013,
            "unit": "ns/op\t 4913943 B/op\t   75235 allocs/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8251013,
            "unit": "ns/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913943,
            "unit": "B/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 320.2,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3808490 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 320.2,
            "unit": "ns/op",
            "extra": "3808490 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3808490 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3808490 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6236,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6236,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 537.5,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2226801 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 537.5,
            "unit": "ns/op",
            "extra": "2226801 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2226801 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2226801 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 769.4,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1550541 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 769.4,
            "unit": "ns/op",
            "extra": "1550541 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1550541 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1550541 times\n4 procs"
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
          "id": "f3d04cb0b1507041dc0dde408a28f468ee1ebfff",
          "message": "feat: Remove feature gate `CombinedServicesFromDifferentHTTPRoutes`, introduce CLI flag (#7569)",
          "timestamp": "2025-07-02T05:12:41Z",
          "tree_id": "fc6cfa2566192ec8b5c2b234537a48f43128887f",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/f3d04cb0b1507041dc0dde408a28f468ee1ebfff"
        },
        "date": 1751434902275,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 370441,
            "unit": "ns/op\t  819451 B/op\t       5 allocs/op",
            "extra": "3981 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 370441,
            "unit": "ns/op",
            "extra": "3981 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819451,
            "unit": "B/op",
            "extra": "3981 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "3981 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7525,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "134989 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7525,
            "unit": "ns/op",
            "extra": "134989 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "134989 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "134989 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 105,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11470678 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 105,
            "unit": "ns/op",
            "extra": "11470678 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11470678 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11470678 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22266,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "55365 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22266,
            "unit": "ns/op",
            "extra": "55365 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "55365 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "55365 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 206924,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5049 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 206924,
            "unit": "ns/op",
            "extra": "5049 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5049 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5049 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2300672,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "508 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2300672,
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
            "value": 39196475,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 39196475,
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
            "value": 9592,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "128020 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9592,
            "unit": "ns/op",
            "extra": "128020 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "128020 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "128020 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7154080,
            "unit": "ns/op\t 4527308 B/op\t   69224 allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7154080,
            "unit": "ns/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527308,
            "unit": "B/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7942,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "151616 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7942,
            "unit": "ns/op",
            "extra": "151616 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "151616 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "151616 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 857216,
            "unit": "ns/op\t  396645 B/op\t    6227 allocs/op",
            "extra": "1243 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 857216,
            "unit": "ns/op",
            "extra": "1243 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396645,
            "unit": "B/op",
            "extra": "1243 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6227,
            "unit": "allocs/op",
            "extra": "1243 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10836,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "110827 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10836,
            "unit": "ns/op",
            "extra": "110827 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "110827 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "110827 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8089164,
            "unit": "ns/op\t 4913974 B/op\t   75235 allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8089164,
            "unit": "ns/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913974,
            "unit": "B/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 324.7,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3721645 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 324.7,
            "unit": "ns/op",
            "extra": "3721645 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3721645 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3721645 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6224,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6224,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 550.8,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2197575 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 550.8,
            "unit": "ns/op",
            "extra": "2197575 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2197575 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2197575 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 877.7,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1333466 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 877.7,
            "unit": "ns/op",
            "extra": "1333466 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1333466 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1333466 times\n4 procs"
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
          "id": "c743635361ed11daa51e95e167149f4cf049cdff",
          "message": "chore(dep): update kubernetes-configuration to v1.5.1 (#7573)\n\nSigned-off-by: Jintao Zhang <zhangjintao9020@gmail.com>",
          "timestamp": "2025-07-02T05:15:00Z",
          "tree_id": "98ae616031c3905110beabcc16f40910ae9c4c4c",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/c743635361ed11daa51e95e167149f4cf049cdff"
        },
        "date": 1751435119699,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 301318,
            "unit": "ns/op\t  819454 B/op\t       5 allocs/op",
            "extra": "3942 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 301318,
            "unit": "ns/op",
            "extra": "3942 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819454,
            "unit": "B/op",
            "extra": "3942 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "3942 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7433,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "159002 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7433,
            "unit": "ns/op",
            "extra": "159002 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "159002 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "159002 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 104.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11492581 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 104.7,
            "unit": "ns/op",
            "extra": "11492581 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11492581 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11492581 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 26729,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "55688 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 26729,
            "unit": "ns/op",
            "extra": "55688 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "55688 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "55688 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 209984,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4950 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 209984,
            "unit": "ns/op",
            "extra": "4950 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4950 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4950 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2424313,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "502 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2424313,
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
            "value": 37899499,
            "unit": "ns/op\t24010754 B/op\t       2 allocs/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 37899499,
            "unit": "ns/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010754,
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
            "value": 9467,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "113728 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9467,
            "unit": "ns/op",
            "extra": "113728 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "113728 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "113728 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7173624,
            "unit": "ns/op\t 4527313 B/op\t   69224 allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7173624,
            "unit": "ns/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527313,
            "unit": "B/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7689,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "153096 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7689,
            "unit": "ns/op",
            "extra": "153096 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "153096 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "153096 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 848181,
            "unit": "ns/op\t  396543 B/op\t    6225 allocs/op",
            "extra": "1280 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 848181,
            "unit": "ns/op",
            "extra": "1280 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396543,
            "unit": "B/op",
            "extra": "1280 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6225,
            "unit": "allocs/op",
            "extra": "1280 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10873,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "105370 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10873,
            "unit": "ns/op",
            "extra": "105370 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "105370 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "105370 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8029652,
            "unit": "ns/op\t 4913982 B/op\t   75235 allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8029652,
            "unit": "ns/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913982,
            "unit": "B/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 322.2,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3736948 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 322.2,
            "unit": "ns/op",
            "extra": "3736948 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3736948 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3736948 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6236,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6236,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 642.3,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2194030 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 642.3,
            "unit": "ns/op",
            "extra": "2194030 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2194030 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2194030 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 790.3,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1367019 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 790.3,
            "unit": "ns/op",
            "extra": "1367019 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1367019 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1367019 times\n4 procs"
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
          "id": "b75425e0c72c29f5c2fb17bede2ed42cd95d91ff",
          "message": "chore(deps): bump sigs.k8s.io/kustomize/api from 0.19.0 to 0.20.0 (#7563)\n\nBumps [sigs.k8s.io/kustomize/api](https://github.com/kubernetes-sigs/kustomize) from 0.19.0 to 0.20.0.\n- [Release notes](https://github.com/kubernetes-sigs/kustomize/releases)\n- [Commits](https://github.com/kubernetes-sigs/kustomize/compare/api/v0.19.0...api/v0.20.0)\n\n---\nupdated-dependencies:\n- dependency-name: sigs.k8s.io/kustomize/api\n  dependency-version: 0.20.0\n  dependency-type: direct:production\n  update-type: version-update:semver-minor\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2025-07-02T09:40:29Z",
          "tree_id": "769158d68aec3ec2331b7c8ee3094ff5517d8bd7",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/b75425e0c72c29f5c2fb17bede2ed42cd95d91ff"
        },
        "date": 1751451082918,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 359375,
            "unit": "ns/op\t  819452 B/op\t       5 allocs/op",
            "extra": "3392 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 359375,
            "unit": "ns/op",
            "extra": "3392 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819452,
            "unit": "B/op",
            "extra": "3392 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "3392 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7604,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "156250 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7604,
            "unit": "ns/op",
            "extra": "156250 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "156250 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "156250 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 105.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11448116 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 105.6,
            "unit": "ns/op",
            "extra": "11448116 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11448116 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11448116 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22255,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "51519 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22255,
            "unit": "ns/op",
            "extra": "51519 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "51519 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "51519 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 230028,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5292 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 230028,
            "unit": "ns/op",
            "extra": "5292 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5292 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5292 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2451288,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "487 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2451288,
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
            "value": 35048022,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 35048022,
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
            "value": 9441,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "127444 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9441,
            "unit": "ns/op",
            "extra": "127444 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "127444 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "127444 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7050912,
            "unit": "ns/op\t 4527304 B/op\t   69224 allocs/op",
            "extra": "168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7050912,
            "unit": "ns/op",
            "extra": "168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527304,
            "unit": "B/op",
            "extra": "168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7795,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "153022 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7795,
            "unit": "ns/op",
            "extra": "153022 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "153022 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "153022 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 863807,
            "unit": "ns/op\t  396591 B/op\t    6226 allocs/op",
            "extra": "1261 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 863807,
            "unit": "ns/op",
            "extra": "1261 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396591,
            "unit": "B/op",
            "extra": "1261 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6226,
            "unit": "allocs/op",
            "extra": "1261 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10966,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "108891 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10966,
            "unit": "ns/op",
            "extra": "108891 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "108891 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "108891 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8455904,
            "unit": "ns/op\t 4913989 B/op\t   75235 allocs/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8455904,
            "unit": "ns/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913989,
            "unit": "B/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 320.3,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3825674 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 320.3,
            "unit": "ns/op",
            "extra": "3825674 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3825674 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3825674 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6232,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6232,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 647.7,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "1939233 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 647.7,
            "unit": "ns/op",
            "extra": "1939233 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "1939233 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1939233 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 786.2,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1540826 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 786.2,
            "unit": "ns/op",
            "extra": "1540826 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1540826 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1540826 times\n4 procs"
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
          "id": "cbfeebd3fac65ef9922774bc8e496b5f0f595314",
          "message": "chore(deps): bump ncipollo/release-action from 1.16.0 to 1.18.0 (#7565)\n\nBumps [ncipollo/release-action](https://github.com/ncipollo/release-action) from 1.16.0 to 1.18.0.\n- [Release notes](https://github.com/ncipollo/release-action/releases)\n- [Commits](https://github.com/ncipollo/release-action/compare/440c8c1cb0ed28b9f43e4d1d670870f059653174...bcfe5470707e8832e12347755757cec0eb3c22af)\n\n---\nupdated-dependencies:\n- dependency-name: ncipollo/release-action\n  dependency-version: 1.18.0\n  dependency-type: direct:production\n  update-type: version-update:semver-minor\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2025-07-02T10:28:15Z",
          "tree_id": "28875a1b010dbd69cddc9e463d795b9c77fb503c",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/cbfeebd3fac65ef9922774bc8e496b5f0f595314"
        },
        "date": 1751453912853,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 289146,
            "unit": "ns/op\t  819451 B/op\t       5 allocs/op",
            "extra": "4119 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 289146,
            "unit": "ns/op",
            "extra": "4119 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819451,
            "unit": "B/op",
            "extra": "4119 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "4119 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 6911,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "171169 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 6911,
            "unit": "ns/op",
            "extra": "171169 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "171169 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "171169 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 104.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11484954 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 104.5,
            "unit": "ns/op",
            "extra": "11484954 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11484954 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11484954 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21593,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "55190 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21593,
            "unit": "ns/op",
            "extra": "55190 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "55190 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "55190 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 207414,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5707 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 207414,
            "unit": "ns/op",
            "extra": "5707 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5707 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5707 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2775199,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "510 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2775199,
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
            "value": 35039288,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 35039288,
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
            "value": 9271,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "127568 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9271,
            "unit": "ns/op",
            "extra": "127568 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "127568 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "127568 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6926236,
            "unit": "ns/op\t 4527308 B/op\t   69224 allocs/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6926236,
            "unit": "ns/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527308,
            "unit": "B/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7869,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "156217 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7869,
            "unit": "ns/op",
            "extra": "156217 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "156217 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "156217 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 842007,
            "unit": "ns/op\t  396513 B/op\t    6225 allocs/op",
            "extra": "1284 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 842007,
            "unit": "ns/op",
            "extra": "1284 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396513,
            "unit": "B/op",
            "extra": "1284 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6225,
            "unit": "allocs/op",
            "extra": "1284 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10752,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "110377 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10752,
            "unit": "ns/op",
            "extra": "110377 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "110377 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "110377 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7958754,
            "unit": "ns/op\t 4913953 B/op\t   75235 allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7958754,
            "unit": "ns/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913953,
            "unit": "B/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 359.3,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3668894 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 359.3,
            "unit": "ns/op",
            "extra": "3668894 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3668894 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3668894 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6252,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6252,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 547.4,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2096454 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 547.4,
            "unit": "ns/op",
            "extra": "2096454 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2096454 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2096454 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 785.6,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1560796 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 785.6,
            "unit": "ns/op",
            "extra": "1560796 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1560796 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1560796 times\n4 procs"
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
          "id": "542fc32f223e706bcd9619c327011295edaf404b",
          "message": "prepare changelog for 3.5.0 (#7575)\n\n* prepare changelog for 3.5.0\n\nSigned-off-by: Jintao Zhang <zhangjintao9020@gmail.com>\n\n* add entry for 3.5.0 compare\n\nSigned-off-by: Jintao Zhang <zhangjintao9020@gmail.com>\n\n---------\n\nSigned-off-by: Jintao Zhang <zhangjintao9020@gmail.com>",
          "timestamp": "2025-07-02T11:12:34Z",
          "tree_id": "5d6577ca86eaf8ca8e5b0bcc73305e493aeb9181",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/542fc32f223e706bcd9619c327011295edaf404b"
        },
        "date": 1751456600431,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 379891,
            "unit": "ns/op\t  819451 B/op\t       5 allocs/op",
            "extra": "3856 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 379891,
            "unit": "ns/op",
            "extra": "3856 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819451,
            "unit": "B/op",
            "extra": "3856 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "3856 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7379,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "157288 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7379,
            "unit": "ns/op",
            "extra": "157288 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "157288 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "157288 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 104.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11482436 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 104.7,
            "unit": "ns/op",
            "extra": "11482436 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11482436 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11482436 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22203,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53752 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22203,
            "unit": "ns/op",
            "extra": "53752 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53752 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53752 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 214199,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5587 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 214199,
            "unit": "ns/op",
            "extra": "5587 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5587 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5587 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2350769,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "447 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2350769,
            "unit": "ns/op",
            "extra": "447 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "447 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "447 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 40011238,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 40011238,
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
            "value": 9498,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "126277 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9498,
            "unit": "ns/op",
            "extra": "126277 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "126277 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "126277 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7282800,
            "unit": "ns/op\t 4527309 B/op\t   69224 allocs/op",
            "extra": "166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7282800,
            "unit": "ns/op",
            "extra": "166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527309,
            "unit": "B/op",
            "extra": "166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7868,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "151671 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7868,
            "unit": "ns/op",
            "extra": "151671 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "151671 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "151671 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 860300,
            "unit": "ns/op\t  396622 B/op\t    6226 allocs/op",
            "extra": "1261 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 860300,
            "unit": "ns/op",
            "extra": "1261 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396622,
            "unit": "B/op",
            "extra": "1261 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6226,
            "unit": "allocs/op",
            "extra": "1261 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10957,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "109444 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10957,
            "unit": "ns/op",
            "extra": "109444 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "109444 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "109444 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8427112,
            "unit": "ns/op\t 4913975 B/op\t   75235 allocs/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8427112,
            "unit": "ns/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913975,
            "unit": "B/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 316.6,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3763698 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 316.6,
            "unit": "ns/op",
            "extra": "3763698 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3763698 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3763698 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.626,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.626,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 593.7,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2215741 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 593.7,
            "unit": "ns/op",
            "extra": "2215741 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2215741 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2215741 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 784.4,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1318488 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 784.4,
            "unit": "ns/op",
            "extra": "1318488 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1318488 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1318488 times\n4 procs"
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
          "id": "21cbe2ab3d4f209c43863eb9f7f2ae561eb19942",
          "message": "add postStart hook to set open file limit (#7572)",
          "timestamp": "2025-07-02T13:12:30Z",
          "tree_id": "802d1c67bf0ea442d031497e0f7162a7b0cf9605",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/21cbe2ab3d4f209c43863eb9f7f2ae561eb19942"
        },
        "date": 1751468400991,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 350516,
            "unit": "ns/op\t  819452 B/op\t       5 allocs/op",
            "extra": "3412 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 350516,
            "unit": "ns/op",
            "extra": "3412 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819452,
            "unit": "B/op",
            "extra": "3412 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "3412 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7548,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "151980 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7548,
            "unit": "ns/op",
            "extra": "151980 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "151980 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "151980 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 104.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11485147 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 104.4,
            "unit": "ns/op",
            "extra": "11485147 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11485147 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11485147 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22245,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54403 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22245,
            "unit": "ns/op",
            "extra": "54403 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54403 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54403 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 267184,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4291 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 267184,
            "unit": "ns/op",
            "extra": "4291 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4291 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4291 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2395267,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "496 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2395267,
            "unit": "ns/op",
            "extra": "496 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "496 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "496 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 39944491,
            "unit": "ns/op\t24010754 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 39944491,
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
            "value": 9472,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "129080 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9472,
            "unit": "ns/op",
            "extra": "129080 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "129080 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "129080 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6931276,
            "unit": "ns/op\t 4527312 B/op\t   69224 allocs/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6931276,
            "unit": "ns/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527312,
            "unit": "B/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7703,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "153643 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7703,
            "unit": "ns/op",
            "extra": "153643 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "153643 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "153643 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 845319,
            "unit": "ns/op\t  396643 B/op\t    6227 allocs/op",
            "extra": "1244 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 845319,
            "unit": "ns/op",
            "extra": "1244 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396643,
            "unit": "B/op",
            "extra": "1244 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6227,
            "unit": "allocs/op",
            "extra": "1244 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10846,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "110286 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10846,
            "unit": "ns/op",
            "extra": "110286 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "110286 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "110286 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8085090,
            "unit": "ns/op\t 4914031 B/op\t   75235 allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8085090,
            "unit": "ns/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914031,
            "unit": "B/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 320.2,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3766218 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 320.2,
            "unit": "ns/op",
            "extra": "3766218 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3766218 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3766218 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6327,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6327,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 547.1,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2175730 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 547.1,
            "unit": "ns/op",
            "extra": "2175730 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2175730 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2175730 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 787.5,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1521711 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 787.5,
            "unit": "ns/op",
            "extra": "1521711 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1521711 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1521711 times\n4 procs"
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
          "id": "8636dd3d47e1d04f8a67d0389b0a1925983448f8",
          "message": "chore(deps): bump golang from `1bb140b` to `20a022e` (#7578)\n\nBumps golang from `1bb140b` to `20a022e`.\n\n---\nupdated-dependencies:\n- dependency-name: golang\n  dependency-version: 1.24.4\n  dependency-type: direct:production\n  update-type: version-update:semver-patch\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2025-07-02T15:07:45Z",
          "tree_id": "acbc8a9d05cbcfad24cead640d61ce2dc2a3672a",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/8636dd3d47e1d04f8a67d0389b0a1925983448f8"
        },
        "date": 1751470684861,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 291151,
            "unit": "ns/op\t  819451 B/op\t       5 allocs/op",
            "extra": "4102 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 291151,
            "unit": "ns/op",
            "extra": "4102 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819451,
            "unit": "B/op",
            "extra": "4102 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "4102 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 6961,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "173262 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 6961,
            "unit": "ns/op",
            "extra": "173262 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "173262 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "173262 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 105.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11382816 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 105.1,
            "unit": "ns/op",
            "extra": "11382816 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11382816 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11382816 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 24205,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "41658 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 24205,
            "unit": "ns/op",
            "extra": "41658 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "41658 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "41658 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 221588,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5454 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 221588,
            "unit": "ns/op",
            "extra": "5454 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5454 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5454 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2381072,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "506 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2381072,
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
            "value": 34695827,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 34695827,
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
            "value": 9349,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "127543 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9349,
            "unit": "ns/op",
            "extra": "127543 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "127543 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "127543 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6975816,
            "unit": "ns/op\t 4527309 B/op\t   69224 allocs/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6975816,
            "unit": "ns/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527309,
            "unit": "B/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7892,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "153116 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7892,
            "unit": "ns/op",
            "extra": "153116 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "153116 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "153116 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 865083,
            "unit": "ns/op\t  396786 B/op\t    6229 allocs/op",
            "extra": "1203 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 865083,
            "unit": "ns/op",
            "extra": "1203 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396786,
            "unit": "B/op",
            "extra": "1203 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6229,
            "unit": "allocs/op",
            "extra": "1203 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10862,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "108873 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10862,
            "unit": "ns/op",
            "extra": "108873 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "108873 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "108873 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8423760,
            "unit": "ns/op\t 4913959 B/op\t   75235 allocs/op",
            "extra": "136 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8423760,
            "unit": "ns/op",
            "extra": "136 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913959,
            "unit": "B/op",
            "extra": "136 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "136 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 319,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3792704 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 319,
            "unit": "ns/op",
            "extra": "3792704 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3792704 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3792704 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6239,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6239,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 539.5,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2229072 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 539.5,
            "unit": "ns/op",
            "extra": "2229072 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2229072 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2229072 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 772.9,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1555503 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 772.9,
            "unit": "ns/op",
            "extra": "1555503 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1555503 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1555503 times\n4 procs"
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
          "id": "d1708de26cc3363c7cc05c0aa9bd24ab7d8a1858",
          "message": "chore(deps): bump github/codeql-action from 3.29.1 to 3.29.2 (#7564)\n\nBumps [github/codeql-action](https://github.com/github/codeql-action) from 3.29.1 to 3.29.2.\n- [Release notes](https://github.com/github/codeql-action/releases)\n- [Changelog](https://github.com/github/codeql-action/blob/main/CHANGELOG.md)\n- [Commits](https://github.com/github/codeql-action/compare/39edc492dbe16b1465b0cafca41432d857bdb31a...181d5eefc20863364f96762470ba6f862bdef56b)\n\n---\nupdated-dependencies:\n- dependency-name: github/codeql-action\n  dependency-version: 3.29.2\n  dependency-type: direct:production\n  update-type: version-update:semver-patch\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2025-07-02T15:07:53Z",
          "tree_id": "0dd2a709504cb2d05b33579164ba3a7ae3769e59",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/d1708de26cc3363c7cc05c0aa9bd24ab7d8a1858"
        },
        "date": 1751470747739,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 294592,
            "unit": "ns/op\t  819452 B/op\t       5 allocs/op",
            "extra": "3440 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 294592,
            "unit": "ns/op",
            "extra": "3440 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819452,
            "unit": "B/op",
            "extra": "3440 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "3440 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7403,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "161907 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7403,
            "unit": "ns/op",
            "extra": "161907 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "161907 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "161907 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 104.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11441575 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 104.6,
            "unit": "ns/op",
            "extra": "11441575 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11441575 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11441575 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 23424,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53090 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 23424,
            "unit": "ns/op",
            "extra": "53090 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53090 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53090 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 206466,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4936 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 206466,
            "unit": "ns/op",
            "extra": "4936 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4936 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4936 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2348876,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "516 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2348876,
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
            "value": 36787137,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 36787137,
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
            "value": 9185,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "131053 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9185,
            "unit": "ns/op",
            "extra": "131053 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "131053 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "131053 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6771231,
            "unit": "ns/op\t 4527314 B/op\t   69224 allocs/op",
            "extra": "176 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6771231,
            "unit": "ns/op",
            "extra": "176 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527314,
            "unit": "B/op",
            "extra": "176 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "176 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7540,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "156076 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7540,
            "unit": "ns/op",
            "extra": "156076 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "156076 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "156076 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 835814,
            "unit": "ns/op\t  396475 B/op\t    6224 allocs/op",
            "extra": "1302 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 835814,
            "unit": "ns/op",
            "extra": "1302 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396475,
            "unit": "B/op",
            "extra": "1302 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6224,
            "unit": "allocs/op",
            "extra": "1302 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10653,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "113090 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10653,
            "unit": "ns/op",
            "extra": "113090 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "113090 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "113090 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7865438,
            "unit": "ns/op\t 4913923 B/op\t   75235 allocs/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7865438,
            "unit": "ns/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913923,
            "unit": "B/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 324.4,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3298510 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 324.4,
            "unit": "ns/op",
            "extra": "3298510 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3298510 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3298510 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6218,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6218,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 542.6,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2207100 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 542.6,
            "unit": "ns/op",
            "extra": "2207100 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2207100 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2207100 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 771.5,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1563615 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 771.5,
            "unit": "ns/op",
            "extra": "1563615 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1563615 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1563615 times\n4 procs"
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
          "id": "3487f809313ec5cfba9567226f737b3162ef033a",
          "message": "chore: update kic image to 3.5 (#7577)\n\nSigned-off-by: Jintao Zhang <zhangjintao9020@gmail.com>",
          "timestamp": "2025-07-02T21:48:57Z",
          "tree_id": "d69a0a61321f1400f9fdee1b985d4645eeb456db",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/3487f809313ec5cfba9567226f737b3162ef033a"
        },
        "date": 1751494731034,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 294836,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "4076 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 294836,
            "unit": "ns/op",
            "extra": "4076 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819453,
            "unit": "B/op",
            "extra": "4076 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "4076 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 8142,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "170714 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 8142,
            "unit": "ns/op",
            "extra": "170714 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "170714 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "170714 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 104.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11466962 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 104.9,
            "unit": "ns/op",
            "extra": "11466962 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11466962 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11466962 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22073,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54060 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22073,
            "unit": "ns/op",
            "extra": "54060 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54060 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54060 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 220873,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4611 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 220873,
            "unit": "ns/op",
            "extra": "4611 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4611 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4611 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2403224,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "493 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2403224,
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
            "value": 33119068,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 33119068,
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
            "value": 9347,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "126760 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9347,
            "unit": "ns/op",
            "extra": "126760 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "126760 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "126760 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7253170,
            "unit": "ns/op\t 4527311 B/op\t   69224 allocs/op",
            "extra": "166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7253170,
            "unit": "ns/op",
            "extra": "166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527311,
            "unit": "B/op",
            "extra": "166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8051,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "150162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8051,
            "unit": "ns/op",
            "extra": "150162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "150162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "150162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 882442,
            "unit": "ns/op\t  396683 B/op\t    6228 allocs/op",
            "extra": "1234 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 882442,
            "unit": "ns/op",
            "extra": "1234 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396683,
            "unit": "B/op",
            "extra": "1234 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6228,
            "unit": "allocs/op",
            "extra": "1234 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11047,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "110863 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11047,
            "unit": "ns/op",
            "extra": "110863 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "110863 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "110863 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8518201,
            "unit": "ns/op\t 4913939 B/op\t   75235 allocs/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8518201,
            "unit": "ns/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913939,
            "unit": "B/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 317.2,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3732052 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 317.2,
            "unit": "ns/op",
            "extra": "3732052 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3732052 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3732052 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6222,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6222,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 550,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2189379 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 550,
            "unit": "ns/op",
            "extra": "2189379 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2189379 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2189379 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 891.3,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1531778 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 891.3,
            "unit": "ns/op",
            "extra": "1531778 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1531778 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1531778 times\n4 procs"
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
          "id": "dcb3ad4ac6609895ca7007d0479000c66d11f326",
          "message": "use whitelist to check httproute filters in validator (#7582)",
          "timestamp": "2025-07-03T07:54:50Z",
          "tree_id": "e264ee8e3e81ca6e8deaf5453fa69b8e01e4554d",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/dcb3ad4ac6609895ca7007d0479000c66d11f326"
        },
        "date": 1751531085134,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 384965,
            "unit": "ns/op\t  819451 B/op\t       5 allocs/op",
            "extra": "4035 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 384965,
            "unit": "ns/op",
            "extra": "4035 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819451,
            "unit": "B/op",
            "extra": "4035 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "4035 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7418,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "162211 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7418,
            "unit": "ns/op",
            "extra": "162211 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "162211 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "162211 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 104.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11481099 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 104.5,
            "unit": "ns/op",
            "extra": "11481099 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11481099 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11481099 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21945,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "55131 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21945,
            "unit": "ns/op",
            "extra": "55131 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "55131 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "55131 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 205876,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5740 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 205876,
            "unit": "ns/op",
            "extra": "5740 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5740 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5740 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2454960,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "500 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2454960,
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
            "value": 39792392,
            "unit": "ns/op\t24010754 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 39792392,
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
            "value": 9162,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "127688 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9162,
            "unit": "ns/op",
            "extra": "127688 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "127688 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "127688 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6929195,
            "unit": "ns/op\t 4527310 B/op\t   69224 allocs/op",
            "extra": "174 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6929195,
            "unit": "ns/op",
            "extra": "174 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527310,
            "unit": "B/op",
            "extra": "174 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "174 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7588,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "154614 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7588,
            "unit": "ns/op",
            "extra": "154614 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "154614 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "154614 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 838711,
            "unit": "ns/op\t  396541 B/op\t    6225 allocs/op",
            "extra": "1285 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 838711,
            "unit": "ns/op",
            "extra": "1285 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396541,
            "unit": "B/op",
            "extra": "1285 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6225,
            "unit": "allocs/op",
            "extra": "1285 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10691,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "112617 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10691,
            "unit": "ns/op",
            "extra": "112617 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "112617 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "112617 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7897325,
            "unit": "ns/op\t 4913978 B/op\t   75235 allocs/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7897325,
            "unit": "ns/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913978,
            "unit": "B/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 314.1,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3817971 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 314.1,
            "unit": "ns/op",
            "extra": "3817971 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3817971 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3817971 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6223,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6223,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 638,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2222400 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 638,
            "unit": "ns/op",
            "extra": "2222400 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2222400 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2222400 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 770.7,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1496353 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 770.7,
            "unit": "ns/op",
            "extra": "1496353 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1496353 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1496353 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "patryk.malek@konghq.com",
            "name": "Patryk Małek",
            "username": "pmalek"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "09f6012507f145f3eda1e239446096f8db03442e",
          "message": "fix: fix KongUpstreamPolicy hash on translation (#7583)\n\n* fix: fix KongUpstreamPolicy hash on translation\n\n* tests: add a valid test case for hash_on uri_capture and hash_fallback with hreader",
          "timestamp": "2025-07-03T18:13:08Z",
          "tree_id": "33dcfbb51170755b5f55acb6b274c8d377bf5c8c",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/09f6012507f145f3eda1e239446096f8db03442e"
        },
        "date": 1751568233650,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 291723,
            "unit": "ns/op\t  819452 B/op\t       5 allocs/op",
            "extra": "4118 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 291723,
            "unit": "ns/op",
            "extra": "4118 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819452,
            "unit": "B/op",
            "extra": "4118 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "4118 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 8206,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "169231 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 8206,
            "unit": "ns/op",
            "extra": "169231 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "169231 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "169231 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 104.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10118294 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 104.7,
            "unit": "ns/op",
            "extra": "10118294 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "10118294 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "10118294 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21926,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "55296 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21926,
            "unit": "ns/op",
            "extra": "55296 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "55296 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "55296 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 209464,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5830 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 209464,
            "unit": "ns/op",
            "extra": "5830 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5830 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5830 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2328992,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "451 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2328992,
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
            "value": 39397832,
            "unit": "ns/op\t24010754 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 39397832,
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
            "value": 9272,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "127837 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9272,
            "unit": "ns/op",
            "extra": "127837 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "127837 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "127837 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6939447,
            "unit": "ns/op\t 4527310 B/op\t   69224 allocs/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6939447,
            "unit": "ns/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527310,
            "unit": "B/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7780,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "155989 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7780,
            "unit": "ns/op",
            "extra": "155989 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "155989 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "155989 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 855156,
            "unit": "ns/op\t  396515 B/op\t    6225 allocs/op",
            "extra": "1293 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 855156,
            "unit": "ns/op",
            "extra": "1293 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396515,
            "unit": "B/op",
            "extra": "1293 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6225,
            "unit": "allocs/op",
            "extra": "1293 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10748,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "111572 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10748,
            "unit": "ns/op",
            "extra": "111572 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "111572 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "111572 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8040958,
            "unit": "ns/op\t 4913940 B/op\t   75235 allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8040958,
            "unit": "ns/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913940,
            "unit": "B/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 326.7,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3703279 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 326.7,
            "unit": "ns/op",
            "extra": "3703279 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3703279 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3703279 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.625,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.625,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 546,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2190189 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 546,
            "unit": "ns/op",
            "extra": "2190189 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2190189 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2190189 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 939.1,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1542650 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 939.1,
            "unit": "ns/op",
            "extra": "1542650 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1542650 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1542650 times\n4 procs"
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
          "id": "063b122e350f6d73cf56e9be2f787c9beca25e20",
          "message": "chore(deps): update kong/kong-gateway@regenerate docker tag to v3.11 (#7588)\n\n* chore(deps): update kong/kong-gateway@regenerate docker tag to v3.11\n\n* chore: regenerate\n\n---------\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>\nCo-authored-by: github-actions <github-actions@users.noreply.github.com>",
          "timestamp": "2025-07-04T05:59:40Z",
          "tree_id": "7a10767b344b91f8305b763c0e408acbf982aabf",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/063b122e350f6d73cf56e9be2f787c9beca25e20"
        },
        "date": 1751610587666,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 360240,
            "unit": "ns/op\t  819452 B/op\t       5 allocs/op",
            "extra": "3459 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 360240,
            "unit": "ns/op",
            "extra": "3459 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819452,
            "unit": "B/op",
            "extra": "3459 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "3459 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 8924,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "160902 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 8924,
            "unit": "ns/op",
            "extra": "160902 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "160902 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "160902 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 104.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11073974 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 104.8,
            "unit": "ns/op",
            "extra": "11073974 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11073974 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11073974 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21517,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "55605 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21517,
            "unit": "ns/op",
            "extra": "55605 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "55605 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "55605 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 207464,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5936 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 207464,
            "unit": "ns/op",
            "extra": "5936 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5936 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5936 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2359711,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "458 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2359711,
            "unit": "ns/op",
            "extra": "458 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "458 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "458 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 43101645,
            "unit": "ns/op\t24010754 B/op\t       2 allocs/op",
            "extra": "44 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 43101645,
            "unit": "ns/op",
            "extra": "44 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010754,
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
            "value": 9437,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "127646 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9437,
            "unit": "ns/op",
            "extra": "127646 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "127646 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "127646 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6944574,
            "unit": "ns/op\t 4527308 B/op\t   69224 allocs/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6944574,
            "unit": "ns/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527308,
            "unit": "B/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7659,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "155979 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7659,
            "unit": "ns/op",
            "extra": "155979 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "155979 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "155979 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 846677,
            "unit": "ns/op\t  396558 B/op\t    6225 allocs/op",
            "extra": "1281 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 846677,
            "unit": "ns/op",
            "extra": "1281 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396558,
            "unit": "B/op",
            "extra": "1281 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6225,
            "unit": "allocs/op",
            "extra": "1281 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10717,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "97861 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10717,
            "unit": "ns/op",
            "extra": "97861 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "97861 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "97861 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7938819,
            "unit": "ns/op\t 4913907 B/op\t   75235 allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7938819,
            "unit": "ns/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913907,
            "unit": "B/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 312.6,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3827166 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 312.6,
            "unit": "ns/op",
            "extra": "3827166 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3827166 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3827166 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6229,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6229,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 540.7,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2219209 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 540.7,
            "unit": "ns/op",
            "extra": "2219209 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2219209 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2219209 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 771.2,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1540932 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 771.2,
            "unit": "ns/op",
            "extra": "1540932 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1540932 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1540932 times\n4 procs"
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
          "id": "db9d02434d9f3ef3bd8d16423e5a0f3983eb1fca",
          "message": "chore: update kubernetes-configuration to v1.5.2 (#7589)\n\nSigned-off-by: Jintao Zhang <zhangjintao9020@gmail.com>",
          "timestamp": "2025-07-04T06:48:15Z",
          "tree_id": "1d6338e181b32e8c30a1d847239389d8c2f1cb69",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/db9d02434d9f3ef3bd8d16423e5a0f3983eb1fca"
        },
        "date": 1751613529908,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 292066,
            "unit": "ns/op\t  819451 B/op\t       5 allocs/op",
            "extra": "4105 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 292066,
            "unit": "ns/op",
            "extra": "4105 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819451,
            "unit": "B/op",
            "extra": "4105 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "4105 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 8321,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "173895 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 8321,
            "unit": "ns/op",
            "extra": "173895 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "173895 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "173895 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 104.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11360212 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 104.8,
            "unit": "ns/op",
            "extra": "11360212 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11360212 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11360212 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21917,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54669 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21917,
            "unit": "ns/op",
            "extra": "54669 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54669 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54669 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 228317,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5122 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 228317,
            "unit": "ns/op",
            "extra": "5122 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5122 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5122 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2415920,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "446 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2415920,
            "unit": "ns/op",
            "extra": "446 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "446 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "446 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 45423769,
            "unit": "ns/op\t24010754 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 45423769,
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
            "value": 9706,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "123681 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9706,
            "unit": "ns/op",
            "extra": "123681 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "123681 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "123681 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7821819,
            "unit": "ns/op\t 4527307 B/op\t   69224 allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7821819,
            "unit": "ns/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527307,
            "unit": "B/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8121,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "147574 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8121,
            "unit": "ns/op",
            "extra": "147574 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "147574 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "147574 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 870190,
            "unit": "ns/op\t  396856 B/op\t    6230 allocs/op",
            "extra": "1189 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 870190,
            "unit": "ns/op",
            "extra": "1189 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396856,
            "unit": "B/op",
            "extra": "1189 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6230,
            "unit": "allocs/op",
            "extra": "1189 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10919,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "108639 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10919,
            "unit": "ns/op",
            "extra": "108639 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "108639 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "108639 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8180931,
            "unit": "ns/op\t 4914000 B/op\t   75235 allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8180931,
            "unit": "ns/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914000,
            "unit": "B/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 315.5,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3815667 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 315.5,
            "unit": "ns/op",
            "extra": "3815667 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3815667 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3815667 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6233,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6233,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 545,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2196267 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 545,
            "unit": "ns/op",
            "extra": "2196267 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2196267 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2196267 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 844,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1253580 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 844,
            "unit": "ns/op",
            "extra": "1253580 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1253580 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1253580 times\n4 procs"
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
          "id": "7b419af751a80fedf1c5cac049e587760adcb511",
          "message": "chore(deps): bump github.com/docker/docker (#7584)\n\nBumps [github.com/docker/docker](https://github.com/docker/docker) from 28.3.0+incompatible to 28.3.1+incompatible.\n- [Release notes](https://github.com/docker/docker/releases)\n- [Commits](https://github.com/docker/docker/compare/v28.3.0...v28.3.1)\n\n---\nupdated-dependencies:\n- dependency-name: github.com/docker/docker\n  dependency-version: 28.3.1+incompatible\n  dependency-type: direct:production\n  update-type: version-update:semver-patch\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2025-07-04T08:52:37Z",
          "tree_id": "d846ebbf424c2474307730e4f134df1ae3b4bfbb",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/7b419af751a80fedf1c5cac049e587760adcb511"
        },
        "date": 1751622568024,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 291269,
            "unit": "ns/op\t  819450 B/op\t       5 allocs/op",
            "extra": "4099 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 291269,
            "unit": "ns/op",
            "extra": "4099 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819450,
            "unit": "B/op",
            "extra": "4099 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "4099 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7833,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "167821 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7833,
            "unit": "ns/op",
            "extra": "167821 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "167821 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "167821 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 104.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11361898 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 104.7,
            "unit": "ns/op",
            "extra": "11361898 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11361898 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11361898 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21687,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "55171 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21687,
            "unit": "ns/op",
            "extra": "55171 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "55171 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "55171 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 211530,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5695 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 211530,
            "unit": "ns/op",
            "extra": "5695 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5695 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5695 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2347683,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "499 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2347683,
            "unit": "ns/op",
            "extra": "499 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "499 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "499 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 39783081,
            "unit": "ns/op\t24010754 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 39783081,
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
            "value": 9185,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "127344 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9185,
            "unit": "ns/op",
            "extra": "127344 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "127344 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "127344 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7196674,
            "unit": "ns/op\t 4527310 B/op\t   69224 allocs/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7196674,
            "unit": "ns/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527310,
            "unit": "B/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7576,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "152126 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7576,
            "unit": "ns/op",
            "extra": "152126 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "152126 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "152126 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 846910,
            "unit": "ns/op\t  396585 B/op\t    6226 allocs/op",
            "extra": "1269 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 846910,
            "unit": "ns/op",
            "extra": "1269 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396585,
            "unit": "B/op",
            "extra": "1269 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6226,
            "unit": "allocs/op",
            "extra": "1269 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10734,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "112617 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10734,
            "unit": "ns/op",
            "extra": "112617 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "112617 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "112617 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7977384,
            "unit": "ns/op\t 4913979 B/op\t   75235 allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7977384,
            "unit": "ns/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913979,
            "unit": "B/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 323.3,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3739444 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 323.3,
            "unit": "ns/op",
            "extra": "3739444 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3739444 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3739444 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6222,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6222,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 539.7,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2203522 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 539.7,
            "unit": "ns/op",
            "extra": "2203522 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2203522 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2203522 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 816.7,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1552946 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 816.7,
            "unit": "ns/op",
            "extra": "1552946 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1552946 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1552946 times\n4 procs"
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
          "id": "55c77a7abee5f8091362bcf27aca9d489211db53",
          "message": "chore(deps): update istio/istioctl docker tag to v1.26.2 (#7465)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2025-07-04T10:55:39Z",
          "tree_id": "9cbf222c8bd15a3b6b40b50d457056b00c4a9788",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/55c77a7abee5f8091362bcf27aca9d489211db53"
        },
        "date": 1751628364440,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 292040,
            "unit": "ns/op\t  819452 B/op\t       5 allocs/op",
            "extra": "4138 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 292040,
            "unit": "ns/op",
            "extra": "4138 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819452,
            "unit": "B/op",
            "extra": "4138 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "4138 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7137,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "168758 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7137,
            "unit": "ns/op",
            "extra": "168758 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "168758 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "168758 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 104.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11463464 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 104.9,
            "unit": "ns/op",
            "extra": "11463464 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11463464 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11463464 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21673,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "55428 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21673,
            "unit": "ns/op",
            "extra": "55428 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "55428 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "55428 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 208240,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5374 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 208240,
            "unit": "ns/op",
            "extra": "5374 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5374 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5374 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2356143,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "460 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2356143,
            "unit": "ns/op",
            "extra": "460 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "460 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "460 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 32201611,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 32201611,
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
            "value": 9101,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "127290 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9101,
            "unit": "ns/op",
            "extra": "127290 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "127290 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "127290 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6859891,
            "unit": "ns/op\t 4527313 B/op\t   69224 allocs/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6859891,
            "unit": "ns/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527313,
            "unit": "B/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7529,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "155103 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7529,
            "unit": "ns/op",
            "extra": "155103 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "155103 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "155103 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 851596,
            "unit": "ns/op\t  396525 B/op\t    6225 allocs/op",
            "extra": "1293 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 851596,
            "unit": "ns/op",
            "extra": "1293 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396525,
            "unit": "B/op",
            "extra": "1293 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6225,
            "unit": "allocs/op",
            "extra": "1293 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10539,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "112941 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10539,
            "unit": "ns/op",
            "extra": "112941 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "112941 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "112941 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7886031,
            "unit": "ns/op\t 4914009 B/op\t   75235 allocs/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7886031,
            "unit": "ns/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914009,
            "unit": "B/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 313.5,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3866924 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 313.5,
            "unit": "ns/op",
            "extra": "3866924 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3866924 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3866924 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6246,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6246,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 540.8,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2229253 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 540.8,
            "unit": "ns/op",
            "extra": "2229253 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2229253 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2229253 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 771.1,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1552290 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 771.1,
            "unit": "ns/op",
            "extra": "1552290 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1552290 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1552290 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "patryk.malek@konghq.com",
            "name": "Patryk Małek",
            "username": "pmalek"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "700d09cb3d14cc9317b9f8076f955ab923931adf",
          "message": "chore: changelog entry for 2.12.8 (#7591)\n\nCo-authored-by: Tao Yi <tao.yi@konghq.com>",
          "timestamp": "2025-07-04T13:55:35Z",
          "tree_id": "2f25b4ed736401e0968de7bd4f6305a9012097c2",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/700d09cb3d14cc9317b9f8076f955ab923931adf"
        },
        "date": 1751639171424,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 301740,
            "unit": "ns/op\t  819454 B/op\t       5 allocs/op",
            "extra": "3978 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 301740,
            "unit": "ns/op",
            "extra": "3978 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819454,
            "unit": "B/op",
            "extra": "3978 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "3978 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7534,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "159662 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7534,
            "unit": "ns/op",
            "extra": "159662 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "159662 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "159662 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 104.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11456140 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 104.5,
            "unit": "ns/op",
            "extra": "11456140 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11456140 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11456140 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21604,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54946 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21604,
            "unit": "ns/op",
            "extra": "54946 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54946 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54946 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 223214,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5572 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 223214,
            "unit": "ns/op",
            "extra": "5572 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5572 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5572 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2459418,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "510 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2459418,
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
            "value": 33434988,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 33434988,
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
            "value": 9307,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "129499 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9307,
            "unit": "ns/op",
            "extra": "129499 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "129499 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "129499 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7025427,
            "unit": "ns/op\t 4527311 B/op\t   69224 allocs/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7025427,
            "unit": "ns/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527311,
            "unit": "B/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7830,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "150210 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7830,
            "unit": "ns/op",
            "extra": "150210 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "150210 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "150210 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 844960,
            "unit": "ns/op\t  396574 B/op\t    6226 allocs/op",
            "extra": "1274 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 844960,
            "unit": "ns/op",
            "extra": "1274 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396574,
            "unit": "B/op",
            "extra": "1274 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6226,
            "unit": "allocs/op",
            "extra": "1274 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10886,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "110835 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10886,
            "unit": "ns/op",
            "extra": "110835 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "110835 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "110835 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8165721,
            "unit": "ns/op\t 4913974 B/op\t   75235 allocs/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8165721,
            "unit": "ns/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913974,
            "unit": "B/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 316.4,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3834270 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 316.4,
            "unit": "ns/op",
            "extra": "3834270 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3834270 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3834270 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6228,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6228,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 543.8,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2231794 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 543.8,
            "unit": "ns/op",
            "extra": "2231794 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2231794 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2231794 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 782,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1481671 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 782,
            "unit": "ns/op",
            "extra": "1481671 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1481671 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1481671 times\n4 procs"
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
          "id": "00e77761c785c28a2b7fe2d0182b948109447784",
          "message": "chore(deps): update helm release kong to v2.50.0 (#7592)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2025-07-07T07:50:53Z",
          "tree_id": "4c91754bc9c5a7dd4e5dde6b7ab779b55360d48b",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/00e77761c785c28a2b7fe2d0182b948109447784"
        },
        "date": 1751876441948,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 414107,
            "unit": "ns/op\t  819451 B/op\t       5 allocs/op",
            "extra": "2449 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 414107,
            "unit": "ns/op",
            "extra": "2449 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819451,
            "unit": "B/op",
            "extra": "2449 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "2449 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7411,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "162460 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7411,
            "unit": "ns/op",
            "extra": "162460 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "162460 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "162460 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 104.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10938927 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 104.5,
            "unit": "ns/op",
            "extra": "10938927 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "10938927 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "10938927 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21899,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54445 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21899,
            "unit": "ns/op",
            "extra": "54445 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54445 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54445 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 211314,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4934 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 211314,
            "unit": "ns/op",
            "extra": "4934 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4934 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4934 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2330957,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "498 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2330957,
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
            "value": 36288187,
            "unit": "ns/op\t24010754 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 36288187,
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
            "value": 9163,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "129357 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9163,
            "unit": "ns/op",
            "extra": "129357 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "129357 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "129357 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6882219,
            "unit": "ns/op\t 4527311 B/op\t   69224 allocs/op",
            "extra": "174 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6882219,
            "unit": "ns/op",
            "extra": "174 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527311,
            "unit": "B/op",
            "extra": "174 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "174 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7614,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "153979 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7614,
            "unit": "ns/op",
            "extra": "153979 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "153979 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "153979 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 841744,
            "unit": "ns/op\t  396492 B/op\t    6224 allocs/op",
            "extra": "1304 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 841744,
            "unit": "ns/op",
            "extra": "1304 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396492,
            "unit": "B/op",
            "extra": "1304 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6224,
            "unit": "allocs/op",
            "extra": "1304 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10555,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "115744 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10555,
            "unit": "ns/op",
            "extra": "115744 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "115744 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "115744 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8022662,
            "unit": "ns/op\t 4913975 B/op\t   75235 allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8022662,
            "unit": "ns/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913975,
            "unit": "B/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 313.3,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3845368 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 313.3,
            "unit": "ns/op",
            "extra": "3845368 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3845368 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3845368 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.625,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.625,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 540.5,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2219946 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 540.5,
            "unit": "ns/op",
            "extra": "2219946 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2219946 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2219946 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 1001,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 1001,
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
          "id": "191d5ec1ec166164b51d40062985d4f6f29a80b7",
          "message": "chore(deps): update dependency golangci/golangci-lint to v2.2.1 (#7560)\n\n* chore(deps): update dependency golangci/golangci-lint to v2.2.1\n\n* chore: make linter happy\n\n---------\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>\nCo-authored-by: Patryk Małek <patryk.malek@konghq.com>",
          "timestamp": "2025-07-07T14:06:54Z",
          "tree_id": "6ed1a1a4244260b372db5fde7b79fb6a0e0e8f73",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/191d5ec1ec166164b51d40062985d4f6f29a80b7"
        },
        "date": 1751899008251,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 299638,
            "unit": "ns/op\t  819454 B/op\t       5 allocs/op",
            "extra": "3973 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 299638,
            "unit": "ns/op",
            "extra": "3973 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819454,
            "unit": "B/op",
            "extra": "3973 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "3973 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7356,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "161865 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7356,
            "unit": "ns/op",
            "extra": "161865 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "161865 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "161865 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 104.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11473545 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 104.7,
            "unit": "ns/op",
            "extra": "11473545 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11473545 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11473545 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22082,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54997 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22082,
            "unit": "ns/op",
            "extra": "54997 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54997 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54997 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 211419,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5587 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 211419,
            "unit": "ns/op",
            "extra": "5587 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5587 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5587 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2373919,
            "unit": "ns/op\t 2408450 B/op\t       2 allocs/op",
            "extra": "466 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2373919,
            "unit": "ns/op",
            "extra": "466 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408450,
            "unit": "B/op",
            "extra": "466 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "466 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 33041006,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 33041006,
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
            "value": 9179,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "131445 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9179,
            "unit": "ns/op",
            "extra": "131445 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "131445 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "131445 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6919469,
            "unit": "ns/op\t 4527307 B/op\t   69224 allocs/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6919469,
            "unit": "ns/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527307,
            "unit": "B/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7557,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "158926 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7557,
            "unit": "ns/op",
            "extra": "158926 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "158926 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "158926 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 849928,
            "unit": "ns/op\t  396512 B/op\t    6225 allocs/op",
            "extra": "1298 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 849928,
            "unit": "ns/op",
            "extra": "1298 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396512,
            "unit": "B/op",
            "extra": "1298 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6225,
            "unit": "allocs/op",
            "extra": "1298 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10553,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "113244 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10553,
            "unit": "ns/op",
            "extra": "113244 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "113244 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "113244 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8003910,
            "unit": "ns/op\t 4913937 B/op\t   75235 allocs/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8003910,
            "unit": "ns/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913937,
            "unit": "B/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 316.7,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3827379 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 316.7,
            "unit": "ns/op",
            "extra": "3827379 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3827379 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3827379 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6228,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6228,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 544.6,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2189817 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 544.6,
            "unit": "ns/op",
            "extra": "2189817 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2189817 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2189817 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 772.6,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1548859 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 772.6,
            "unit": "ns/op",
            "extra": "1548859 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1548859 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1548859 times\n4 procs"
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
          "id": "c11582b8c9da2de61b639ee31fa0768f4f60e186",
          "message": "fix(ci): fix release docs workflow (#7596)\n\n* fix(ci): fix release docs workflow\n\nSigned-off-by: Jintao Zhang <zhangjintao9020@gmail.com>\n\n* add error handling for WriteString func\n\nSigned-off-by: Jintao Zhang <zhangjintao9020@gmail.com>\n\n* store generated docs in KIC repo\n\nSigned-off-by: Jintao Zhang <zhangjintao9020@gmail.com>\n\n---------\n\nSigned-off-by: Jintao Zhang <zhangjintao9020@gmail.com>",
          "timestamp": "2025-07-08T16:52:44Z",
          "tree_id": "26d10159a0a3e1bd59b455440e445adb168e44b3",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/c11582b8c9da2de61b639ee31fa0768f4f60e186"
        },
        "date": 1751995427275,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 291259,
            "unit": "ns/op\t  819452 B/op\t       5 allocs/op",
            "extra": "4130 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 291259,
            "unit": "ns/op",
            "extra": "4130 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819452,
            "unit": "B/op",
            "extra": "4130 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "4130 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 6896,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "170506 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 6896,
            "unit": "ns/op",
            "extra": "170506 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "170506 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "170506 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 104.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11495721 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 104.7,
            "unit": "ns/op",
            "extra": "11495721 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11495721 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11495721 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21612,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "55738 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21612,
            "unit": "ns/op",
            "extra": "55738 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "55738 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "55738 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 249371,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5763 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 249371,
            "unit": "ns/op",
            "extra": "5763 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5763 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5763 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2319668,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "442 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2319668,
            "unit": "ns/op",
            "extra": "442 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
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
            "value": 37499870,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 37499870,
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
            "value": 9245,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "127336 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9245,
            "unit": "ns/op",
            "extra": "127336 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "127336 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "127336 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6899247,
            "unit": "ns/op\t 4527312 B/op\t   69224 allocs/op",
            "extra": "174 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6899247,
            "unit": "ns/op",
            "extra": "174 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527312,
            "unit": "B/op",
            "extra": "174 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "174 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7695,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "155425 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7695,
            "unit": "ns/op",
            "extra": "155425 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "155425 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "155425 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 848798,
            "unit": "ns/op\t  396516 B/op\t    6225 allocs/op",
            "extra": "1285 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 848798,
            "unit": "ns/op",
            "extra": "1285 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396516,
            "unit": "B/op",
            "extra": "1285 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6225,
            "unit": "allocs/op",
            "extra": "1285 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10677,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "110692 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10677,
            "unit": "ns/op",
            "extra": "110692 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "110692 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "110692 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7907859,
            "unit": "ns/op\t 4913948 B/op\t   75235 allocs/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7907859,
            "unit": "ns/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913948,
            "unit": "B/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 315.2,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3816349 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 315.2,
            "unit": "ns/op",
            "extra": "3816349 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3816349 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3816349 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6231,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6231,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 540.9,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2237684 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 540.9,
            "unit": "ns/op",
            "extra": "2237684 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2237684 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2237684 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 768,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1554073 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 768,
            "unit": "ns/op",
            "extra": "1554073 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1554073 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1554073 times\n4 procs"
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
          "distinct": false,
          "id": "85aaf82bd62e3a4ed60e9869b69239ddcf94ab4f",
          "message": "chore(deps): update helm release kong to v2.51.0 (#7595)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2025-07-09T07:58:49Z",
          "tree_id": "6b5ed54a16e9c786275baf5594485ce2891c0b80",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/85aaf82bd62e3a4ed60e9869b69239ddcf94ab4f"
        },
        "date": 1752051321484,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 299369,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "3686 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 299369,
            "unit": "ns/op",
            "extra": "3686 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819453,
            "unit": "B/op",
            "extra": "3686 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "3686 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 8719,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "160218 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 8719,
            "unit": "ns/op",
            "extra": "160218 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "160218 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "160218 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 104.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11410057 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 104.6,
            "unit": "ns/op",
            "extra": "11410057 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11410057 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11410057 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21674,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "55141 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21674,
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
            "value": 212172,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5164 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 212172,
            "unit": "ns/op",
            "extra": "5164 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5164 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5164 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2315690,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "447 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2315690,
            "unit": "ns/op",
            "extra": "447 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "447 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "447 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 39369104,
            "unit": "ns/op\t24010754 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 39369104,
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
            "value": 9263,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "132356 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9263,
            "unit": "ns/op",
            "extra": "132356 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "132356 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "132356 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6866735,
            "unit": "ns/op\t 4527305 B/op\t   69224 allocs/op",
            "extra": "175 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6866735,
            "unit": "ns/op",
            "extra": "175 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527305,
            "unit": "B/op",
            "extra": "175 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "175 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7546,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "157351 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7546,
            "unit": "ns/op",
            "extra": "157351 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "157351 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "157351 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 851954,
            "unit": "ns/op\t  396576 B/op\t    6226 allocs/op",
            "extra": "1273 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 851954,
            "unit": "ns/op",
            "extra": "1273 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396576,
            "unit": "B/op",
            "extra": "1273 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6226,
            "unit": "allocs/op",
            "extra": "1273 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10579,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "113600 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10579,
            "unit": "ns/op",
            "extra": "113600 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "113600 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "113600 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7900167,
            "unit": "ns/op\t 4913976 B/op\t   75235 allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7900167,
            "unit": "ns/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913976,
            "unit": "B/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 323.8,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3659527 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 323.8,
            "unit": "ns/op",
            "extra": "3659527 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3659527 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3659527 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6239,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6239,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 540.5,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2179722 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 540.5,
            "unit": "ns/op",
            "extra": "2179722 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2179722 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2179722 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 852.3,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1535990 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 852.3,
            "unit": "ns/op",
            "extra": "1535990 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1535990 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1535990 times\n4 procs"
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
          "id": "8daefdeb761377168b0daaa304787922d94a2010",
          "message": "chore(docs): maintain consistency with the original value format (#7597)\n\nSigned-off-by: Jintao Zhang <zhangjintao9020@gmail.com>\nCo-authored-by: Patryk Małek <patryk.malek@konghq.com>",
          "timestamp": "2025-07-09T15:08:46Z",
          "tree_id": "645d0ef936323f12ddf9b95f8ea59a506b17d09c",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/8daefdeb761377168b0daaa304787922d94a2010"
        },
        "date": 1752075519445,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 293959,
            "unit": "ns/op\t  819450 B/op\t       5 allocs/op",
            "extra": "4075 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 293959,
            "unit": "ns/op",
            "extra": "4075 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819450,
            "unit": "B/op",
            "extra": "4075 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "4075 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7461,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "162442 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7461,
            "unit": "ns/op",
            "extra": "162442 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "162442 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "162442 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 104.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11490315 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 104.7,
            "unit": "ns/op",
            "extra": "11490315 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11490315 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11490315 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21656,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "55023 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21656,
            "unit": "ns/op",
            "extra": "55023 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "55023 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "55023 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 211210,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5928 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 211210,
            "unit": "ns/op",
            "extra": "5928 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5928 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5928 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2334273,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "505 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2334273,
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
            "value": 37291637,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 37291637,
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
            "value": 9203,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "131408 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9203,
            "unit": "ns/op",
            "extra": "131408 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "131408 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "131408 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6826293,
            "unit": "ns/op\t 4527309 B/op\t   69224 allocs/op",
            "extra": "175 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6826293,
            "unit": "ns/op",
            "extra": "175 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527309,
            "unit": "B/op",
            "extra": "175 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "175 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7637,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "156495 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7637,
            "unit": "ns/op",
            "extra": "156495 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "156495 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "156495 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 833817,
            "unit": "ns/op\t  396539 B/op\t    6225 allocs/op",
            "extra": "1282 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 833817,
            "unit": "ns/op",
            "extra": "1282 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396539,
            "unit": "B/op",
            "extra": "1282 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6225,
            "unit": "allocs/op",
            "extra": "1282 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10567,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "112999 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10567,
            "unit": "ns/op",
            "extra": "112999 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "112999 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "112999 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7901138,
            "unit": "ns/op\t 4913938 B/op\t   75235 allocs/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7901138,
            "unit": "ns/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913938,
            "unit": "B/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 316.7,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3824775 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 316.7,
            "unit": "ns/op",
            "extra": "3824775 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3824775 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3824775 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6224,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6224,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 545.1,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2210935 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 545.1,
            "unit": "ns/op",
            "extra": "2210935 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2210935 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2210935 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 765.9,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1554222 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 765.9,
            "unit": "ns/op",
            "extra": "1554222 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1554222 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1554222 times\n4 procs"
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
          "id": "da352b8fe7eee884c560cf2c72528b62dbd07356",
          "message": "chore(deps): bump golang from 1.24.4 to 1.24.5 (#7599)\n\nBumps golang from 1.24.4 to 1.24.5.\n\n---\nupdated-dependencies:\n- dependency-name: golang\n  dependency-version: 1.24.5\n  dependency-type: direct:production\n  update-type: version-update:semver-patch\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2025-07-09T15:24:33Z",
          "tree_id": "c5a0812b6814d2107c62f69d2f85b470098388ea",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/da352b8fe7eee884c560cf2c72528b62dbd07356"
        },
        "date": 1752076462306,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 417394,
            "unit": "ns/op\t  819451 B/op\t       5 allocs/op",
            "extra": "3469 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 417394,
            "unit": "ns/op",
            "extra": "3469 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819451,
            "unit": "B/op",
            "extra": "3469 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "3469 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7482,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "163567 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7482,
            "unit": "ns/op",
            "extra": "163567 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "163567 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "163567 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 104.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11451456 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 104.7,
            "unit": "ns/op",
            "extra": "11451456 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11451456 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11451456 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21744,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54343 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21744,
            "unit": "ns/op",
            "extra": "54343 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54343 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54343 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 206633,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5769 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 206633,
            "unit": "ns/op",
            "extra": "5769 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5769 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5769 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2310508,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "511 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2310508,
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
            "value": 31491693,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 31491693,
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
            "value": 9355,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "127758 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9355,
            "unit": "ns/op",
            "extra": "127758 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "127758 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "127758 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6888997,
            "unit": "ns/op\t 4527308 B/op\t   69224 allocs/op",
            "extra": "174 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6888997,
            "unit": "ns/op",
            "extra": "174 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527308,
            "unit": "B/op",
            "extra": "174 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "174 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7669,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "153213 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7669,
            "unit": "ns/op",
            "extra": "153213 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "153213 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "153213 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 853602,
            "unit": "ns/op\t  396632 B/op\t    6227 allocs/op",
            "extra": "1255 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 853602,
            "unit": "ns/op",
            "extra": "1255 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396632,
            "unit": "B/op",
            "extra": "1255 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6227,
            "unit": "allocs/op",
            "extra": "1255 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10760,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "109903 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10760,
            "unit": "ns/op",
            "extra": "109903 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "109903 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "109903 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7883815,
            "unit": "ns/op\t 4913939 B/op\t   75235 allocs/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7883815,
            "unit": "ns/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913939,
            "unit": "B/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 314.4,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3827983 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 314.4,
            "unit": "ns/op",
            "extra": "3827983 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3827983 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3827983 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.624,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.624,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 543.8,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2202532 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 543.8,
            "unit": "ns/op",
            "extra": "2202532 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2202532 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2202532 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 822.3,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1241133 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 822.3,
            "unit": "ns/op",
            "extra": "1241133 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1241133 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1241133 times\n4 procs"
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
          "distinct": false,
          "id": "50258ed4d8470c0af1bf50923b28f8221d428dd1",
          "message": "chore(deps): update dependency mikefarah/yq to v4.46.1 (#7601)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2025-07-10T08:45:48Z",
          "tree_id": "7422ccb55c008c3cd7b0974d1e3be89458203915",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/50258ed4d8470c0af1bf50923b28f8221d428dd1"
        },
        "date": 1752138958085,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 291612,
            "unit": "ns/op\t  819451 B/op\t       5 allocs/op",
            "extra": "3511 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 291612,
            "unit": "ns/op",
            "extra": "3511 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819451,
            "unit": "B/op",
            "extra": "3511 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "3511 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7309,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "163065 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7309,
            "unit": "ns/op",
            "extra": "163065 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "163065 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "163065 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 104.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11430339 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 104.9,
            "unit": "ns/op",
            "extra": "11430339 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11430339 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11430339 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21593,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "55077 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21593,
            "unit": "ns/op",
            "extra": "55077 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "55077 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "55077 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 217791,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5906 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 217791,
            "unit": "ns/op",
            "extra": "5906 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5906 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5906 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2448156,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "418 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2448156,
            "unit": "ns/op",
            "extra": "418 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
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
            "value": 36171908,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "30 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 36171908,
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
            "value": 9312,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "128732 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9312,
            "unit": "ns/op",
            "extra": "128732 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "128732 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "128732 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7004888,
            "unit": "ns/op\t 4527306 B/op\t   69224 allocs/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7004888,
            "unit": "ns/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527306,
            "unit": "B/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7688,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "156272 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7688,
            "unit": "ns/op",
            "extra": "156272 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "156272 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "156272 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 851381,
            "unit": "ns/op\t  396531 B/op\t    6225 allocs/op",
            "extra": "1286 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 851381,
            "unit": "ns/op",
            "extra": "1286 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396531,
            "unit": "B/op",
            "extra": "1286 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6225,
            "unit": "allocs/op",
            "extra": "1286 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10847,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "109975 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10847,
            "unit": "ns/op",
            "extra": "109975 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "109975 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "109975 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8076545,
            "unit": "ns/op\t 4913955 B/op\t   75235 allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8076545,
            "unit": "ns/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913955,
            "unit": "B/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 322,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3688240 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 322,
            "unit": "ns/op",
            "extra": "3688240 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3688240 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3688240 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6238,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6238,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 545.8,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2220447 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 545.8,
            "unit": "ns/op",
            "extra": "2220447 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2220447 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2220447 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 942.2,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1536418 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 942.2,
            "unit": "ns/op",
            "extra": "1536418 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1536418 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1536418 times\n4 procs"
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
          "distinct": false,
          "id": "7d68e77e042b8acb578b4d3735571621c95cd1ae",
          "message": "chore(deps): update helm release kuma to v2.11.2 (#7536)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2025-07-10T08:45:56Z",
          "tree_id": "4ac438e48b99d16f5cf2565302a3bf72062f8744",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/7d68e77e042b8acb578b4d3735571621c95cd1ae"
        },
        "date": 1752138961312,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 297771,
            "unit": "ns/op\t  819452 B/op\t       5 allocs/op",
            "extra": "4063 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 297771,
            "unit": "ns/op",
            "extra": "4063 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819452,
            "unit": "B/op",
            "extra": "4063 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "4063 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7461,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "159044 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7461,
            "unit": "ns/op",
            "extra": "159044 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "159044 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "159044 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 104.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11467942 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 104.7,
            "unit": "ns/op",
            "extra": "11467942 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11467942 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11467942 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21603,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "55305 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21603,
            "unit": "ns/op",
            "extra": "55305 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "55305 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "55305 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 214300,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5557 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 214300,
            "unit": "ns/op",
            "extra": "5557 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5557 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5557 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2316478,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "500 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2316478,
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
            "value": 33249000,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 33249000,
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
            "value": 9414,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "123852 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9414,
            "unit": "ns/op",
            "extra": "123852 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "123852 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "123852 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7192003,
            "unit": "ns/op\t 4527313 B/op\t   69224 allocs/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7192003,
            "unit": "ns/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527313,
            "unit": "B/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7741,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "153439 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7741,
            "unit": "ns/op",
            "extra": "153439 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "153439 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "153439 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 866761,
            "unit": "ns/op\t  396595 B/op\t    6226 allocs/op",
            "extra": "1266 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 866761,
            "unit": "ns/op",
            "extra": "1266 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396595,
            "unit": "B/op",
            "extra": "1266 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6226,
            "unit": "allocs/op",
            "extra": "1266 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10872,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "109676 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10872,
            "unit": "ns/op",
            "extra": "109676 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "109676 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "109676 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8174024,
            "unit": "ns/op\t 4914005 B/op\t   75235 allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8174024,
            "unit": "ns/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914005,
            "unit": "B/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 320.4,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3758000 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 320.4,
            "unit": "ns/op",
            "extra": "3758000 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3758000 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3758000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6233,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6233,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 538,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2223348 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 538,
            "unit": "ns/op",
            "extra": "2223348 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2223348 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2223348 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 774.2,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1544700 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 774.2,
            "unit": "ns/op",
            "extra": "1544700 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1544700 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1544700 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "patryk.malek@konghq.com",
            "name": "Patryk Małek",
            "username": "pmalek"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "3b3729665207191020b0cdb5adc2660ba40fe576",
          "message": "chore: 3.4.8 changelog (#7602)",
          "timestamp": "2025-07-10T13:08:54+02:00",
          "tree_id": "03e736be8e3159534762c688865ad5b12c068f24",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/3b3729665207191020b0cdb5adc2660ba40fe576"
        },
        "date": 1752145904564,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 293238,
            "unit": "ns/op\t  819452 B/op\t       5 allocs/op",
            "extra": "4102 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 293238,
            "unit": "ns/op",
            "extra": "4102 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819452,
            "unit": "B/op",
            "extra": "4102 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "4102 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7026,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "169860 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7026,
            "unit": "ns/op",
            "extra": "169860 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "169860 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "169860 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 106.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11333431 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 106.3,
            "unit": "ns/op",
            "extra": "11333431 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11333431 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11333431 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21634,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "55730 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21634,
            "unit": "ns/op",
            "extra": "55730 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "55730 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "55730 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 212872,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5640 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 212872,
            "unit": "ns/op",
            "extra": "5640 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5640 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5640 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2355017,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "496 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2355017,
            "unit": "ns/op",
            "extra": "496 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "496 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "496 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 33737773,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 33737773,
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
            "value": 9378,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "125382 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9378,
            "unit": "ns/op",
            "extra": "125382 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "125382 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "125382 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7051993,
            "unit": "ns/op\t 4527309 B/op\t   69224 allocs/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7051993,
            "unit": "ns/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527309,
            "unit": "B/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7790,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "155232 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7790,
            "unit": "ns/op",
            "extra": "155232 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "155232 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "155232 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 840037,
            "unit": "ns/op\t  396504 B/op\t    6225 allocs/op",
            "extra": "1291 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 840037,
            "unit": "ns/op",
            "extra": "1291 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396504,
            "unit": "B/op",
            "extra": "1291 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6225,
            "unit": "allocs/op",
            "extra": "1291 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10747,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "111199 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10747,
            "unit": "ns/op",
            "extra": "111199 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "111199 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "111199 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7910562,
            "unit": "ns/op\t 4913952 B/op\t   75235 allocs/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7910562,
            "unit": "ns/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913952,
            "unit": "B/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 314,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3710313 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 314,
            "unit": "ns/op",
            "extra": "3710313 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3710313 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3710313 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6243,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6243,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 544.4,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2213634 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 544.4,
            "unit": "ns/op",
            "extra": "2213634 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2213634 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2213634 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 938.7,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1530404 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 938.7,
            "unit": "ns/op",
            "extra": "1530404 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1530404 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1530404 times\n4 procs"
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
          "id": "a7bf8059c924e93241a973b96bac6c76b6a84d97",
          "message": "chore(ci): add modernize as one of the linter steps (#7603)\n\n* chore(ci): add modernize as one of the linter steps\n\n* Update internal/annotations/annotations.go\n\nCo-authored-by: Patryk Małek <patryk.malek@konghq.com>\n\n---------\n\nCo-authored-by: Patryk Małek <patryk.malek@konghq.com>",
          "timestamp": "2025-07-10T11:04:56Z",
          "tree_id": "773d362dc7a87b8d62dbfc35f171b93f9c29bbe8",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/a7bf8059c924e93241a973b96bac6c76b6a84d97"
        },
        "date": 1752147538952,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 694414,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "1485 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 694414,
            "unit": "ns/op",
            "extra": "1485 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819453,
            "unit": "B/op",
            "extra": "1485 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1485 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7360,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "162282 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7360,
            "unit": "ns/op",
            "extra": "162282 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "162282 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "162282 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 107.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11132899 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 107.6,
            "unit": "ns/op",
            "extra": "11132899 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11132899 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11132899 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22145,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54087 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22145,
            "unit": "ns/op",
            "extra": "54087 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54087 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54087 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 217824,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4710 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 217824,
            "unit": "ns/op",
            "extra": "4710 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4710 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4710 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2420255,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "423 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2420255,
            "unit": "ns/op",
            "extra": "423 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "423 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "423 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 41753330,
            "unit": "ns/op\t24010754 B/op\t       2 allocs/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 41753330,
            "unit": "ns/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010754,
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
            "value": 10756,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "103059 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10756,
            "unit": "ns/op",
            "extra": "103059 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "103059 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "103059 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7320721,
            "unit": "ns/op\t 4527316 B/op\t   69224 allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7320721,
            "unit": "ns/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527316,
            "unit": "B/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7795,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "154268 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7795,
            "unit": "ns/op",
            "extra": "154268 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "154268 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "154268 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 838119,
            "unit": "ns/op\t  392710 B/op\t    6167 allocs/op",
            "extra": "1417 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 838119,
            "unit": "ns/op",
            "extra": "1417 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392710,
            "unit": "B/op",
            "extra": "1417 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6167,
            "unit": "allocs/op",
            "extra": "1417 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10797,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "110905 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10797,
            "unit": "ns/op",
            "extra": "110905 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "110905 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "110905 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8010835,
            "unit": "ns/op\t 4913928 B/op\t   75235 allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8010835,
            "unit": "ns/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913928,
            "unit": "B/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 317.1,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3764042 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 317.1,
            "unit": "ns/op",
            "extra": "3764042 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3764042 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3764042 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.37,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.37,
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
            "value": 553.6,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2174001 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 553.6,
            "unit": "ns/op",
            "extra": "2174001 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2174001 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2174001 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 1012,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 1012,
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
          "distinct": false,
          "id": "b2f6af1598ec38069740b8c50433c6b0fe44efbd",
          "message": "chore(deps): update kindest/node@only-patch docker tag to v1.31.9 (#7438)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2025-07-10T12:02:21Z",
          "tree_id": "d50a612a4a5da75da941a75ab22666124e89287c",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/b2f6af1598ec38069740b8c50433c6b0fe44efbd"
        },
        "date": 1752150717606,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 699756,
            "unit": "ns/op\t  819452 B/op\t       5 allocs/op",
            "extra": "1555 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 699756,
            "unit": "ns/op",
            "extra": "1555 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819452,
            "unit": "B/op",
            "extra": "1555 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1555 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7245,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "165063 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7245,
            "unit": "ns/op",
            "extra": "165063 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "165063 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "165063 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 109.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11017520 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 109.1,
            "unit": "ns/op",
            "extra": "11017520 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11017520 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11017520 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 23295,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "50773 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 23295,
            "unit": "ns/op",
            "extra": "50773 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "50773 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "50773 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 219355,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5418 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 219355,
            "unit": "ns/op",
            "extra": "5418 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5418 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5418 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2517043,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "453 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2517043,
            "unit": "ns/op",
            "extra": "453 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
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
            "value": 38860876,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "62 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 38860876,
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
            "value": 11640,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "104089 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 11640,
            "unit": "ns/op",
            "extra": "104089 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "104089 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "104089 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7109383,
            "unit": "ns/op\t 4527313 B/op\t   69224 allocs/op",
            "extra": "168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7109383,
            "unit": "ns/op",
            "extra": "168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527313,
            "unit": "B/op",
            "extra": "168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8058,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "153178 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8058,
            "unit": "ns/op",
            "extra": "153178 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "153178 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "153178 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 875614,
            "unit": "ns/op\t  392704 B/op\t    6167 allocs/op",
            "extra": "1330 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 875614,
            "unit": "ns/op",
            "extra": "1330 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392704,
            "unit": "B/op",
            "extra": "1330 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6167,
            "unit": "allocs/op",
            "extra": "1330 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11210,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "107820 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11210,
            "unit": "ns/op",
            "extra": "107820 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "107820 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "107820 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8288901,
            "unit": "ns/op\t 4913941 B/op\t   75235 allocs/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8288901,
            "unit": "ns/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913941,
            "unit": "B/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 332.8,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3600283 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 332.8,
            "unit": "ns/op",
            "extra": "3600283 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3600283 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3600283 times\n4 procs"
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
            "value": 564.5,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2117877 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 564.5,
            "unit": "ns/op",
            "extra": "2117877 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2117877 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2117877 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 981.1,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1240852 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 981.1,
            "unit": "ns/op",
            "extra": "1240852 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1240852 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1240852 times\n4 procs"
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
          "distinct": false,
          "id": "b2d3657ea3395b79c6b25853476b6ad5937e135f",
          "message": "chore(deps): update kindest/node docker tag to v1.33.1 (#7436)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>\nCo-authored-by: Patryk Małek <patryk.malek@konghq.com>",
          "timestamp": "2025-07-10T12:04:57Z",
          "tree_id": "de91fd0153bbdce890dad025fce1af00f331080c",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/b2d3657ea3395b79c6b25853476b6ad5937e135f"
        },
        "date": 1752150879343,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 641848,
            "unit": "ns/op\t  819456 B/op\t       5 allocs/op",
            "extra": "1754 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 641848,
            "unit": "ns/op",
            "extra": "1754 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819456,
            "unit": "B/op",
            "extra": "1754 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1754 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7214,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "164913 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7214,
            "unit": "ns/op",
            "extra": "164913 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "164913 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "164913 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 107.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11094151 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 107.8,
            "unit": "ns/op",
            "extra": "11094151 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11094151 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11094151 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22356,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53602 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22356,
            "unit": "ns/op",
            "extra": "53602 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53602 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53602 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 231403,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4509 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 231403,
            "unit": "ns/op",
            "extra": "4509 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4509 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4509 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2381178,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "493 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2381178,
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
            "value": 45102260,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "51 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 45102260,
            "unit": "ns/op",
            "extra": "51 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "51 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "51 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9578,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "124774 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9578,
            "unit": "ns/op",
            "extra": "124774 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "124774 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "124774 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7615873,
            "unit": "ns/op\t 4527312 B/op\t   69224 allocs/op",
            "extra": "156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7615873,
            "unit": "ns/op",
            "extra": "156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527312,
            "unit": "B/op",
            "extra": "156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7962,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "149214 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7962,
            "unit": "ns/op",
            "extra": "149214 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "149214 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "149214 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 865798,
            "unit": "ns/op\t  392704 B/op\t    6167 allocs/op",
            "extra": "1370 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 865798,
            "unit": "ns/op",
            "extra": "1370 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392704,
            "unit": "B/op",
            "extra": "1370 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6167,
            "unit": "allocs/op",
            "extra": "1370 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11194,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "106692 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11194,
            "unit": "ns/op",
            "extra": "106692 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "106692 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "106692 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8623656,
            "unit": "ns/op\t 4913948 B/op\t   75235 allocs/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8623656,
            "unit": "ns/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913948,
            "unit": "B/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 319,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3751455 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 319,
            "unit": "ns/op",
            "extra": "3751455 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3751455 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3751455 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.26,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.26,
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
            "value": 554.2,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2159060 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 554.2,
            "unit": "ns/op",
            "extra": "2159060 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2159060 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2159060 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 802.1,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1489764 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 802.1,
            "unit": "ns/op",
            "extra": "1489764 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1489764 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1489764 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "patryk.malek@konghq.com",
            "name": "Patryk Małek",
            "username": "pmalek"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "cf9a5596f679717613695e522d1370e8229fe3f9",
          "message": "chore: fix 3.4.8 changelog (#7605)",
          "timestamp": "2025-07-10T12:09:09Z",
          "tree_id": "c7eb56b0ea863586b1c7299f8742daaa637b2b51",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/cf9a5596f679717613695e522d1370e8229fe3f9"
        },
        "date": 1752151127991,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 635170,
            "unit": "ns/op\t  819450 B/op\t       5 allocs/op",
            "extra": "1765 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 635170,
            "unit": "ns/op",
            "extra": "1765 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819450,
            "unit": "B/op",
            "extra": "1765 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1765 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7394,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "165874 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7394,
            "unit": "ns/op",
            "extra": "165874 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "165874 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "165874 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 109.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "9707235 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 109.4,
            "unit": "ns/op",
            "extra": "9707235 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "9707235 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "9707235 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22111,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53770 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22111,
            "unit": "ns/op",
            "extra": "53770 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53770 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53770 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 236616,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4657 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 236616,
            "unit": "ns/op",
            "extra": "4657 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4657 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4657 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2696648,
            "unit": "ns/op\t 2408450 B/op\t       2 allocs/op",
            "extra": "430 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2696648,
            "unit": "ns/op",
            "extra": "430 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408450,
            "unit": "B/op",
            "extra": "430 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "430 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 33721179,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 33721179,
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
            "value": 9425,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "128763 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9425,
            "unit": "ns/op",
            "extra": "128763 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "128763 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "128763 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7106868,
            "unit": "ns/op\t 4527313 B/op\t   69224 allocs/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7106868,
            "unit": "ns/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527313,
            "unit": "B/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7863,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "151368 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7863,
            "unit": "ns/op",
            "extra": "151368 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "151368 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "151368 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 845148,
            "unit": "ns/op\t  392711 B/op\t    6167 allocs/op",
            "extra": "1388 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 845148,
            "unit": "ns/op",
            "extra": "1388 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392711,
            "unit": "B/op",
            "extra": "1388 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6167,
            "unit": "allocs/op",
            "extra": "1388 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10972,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "107480 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10972,
            "unit": "ns/op",
            "extra": "107480 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "107480 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "107480 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8325279,
            "unit": "ns/op\t 4913971 B/op\t   75235 allocs/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8325279,
            "unit": "ns/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913971,
            "unit": "B/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 317.9,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3773286 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 317.9,
            "unit": "ns/op",
            "extra": "3773286 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3773286 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3773286 times\n4 procs"
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
            "value": 558.6,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2138328 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 558.6,
            "unit": "ns/op",
            "extra": "2138328 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2138328 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2138328 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 800.1,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1498048 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 800.1,
            "unit": "ns/op",
            "extra": "1498048 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1498048 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1498048 times\n4 procs"
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
          "id": "d8eb9a0bb3f673ea1d52008567368a3fdb6b44da",
          "message": "chore: move license controller to internal (#7606)\n\n* chore: move license controller to internal\n\n* chore: move license controller to internal",
          "timestamp": "2025-07-10T14:29:01Z",
          "tree_id": "1d114f6250f433c881ccacef9aff090cd2a7236c",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/d8eb9a0bb3f673ea1d52008567368a3fdb6b44da"
        },
        "date": 1752159493381,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 693946,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "1514 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 693946,
            "unit": "ns/op",
            "extra": "1514 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819453,
            "unit": "B/op",
            "extra": "1514 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1514 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7126,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "170436 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7126,
            "unit": "ns/op",
            "extra": "170436 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "170436 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "170436 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 107.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11164881 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 107.5,
            "unit": "ns/op",
            "extra": "11164881 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11164881 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11164881 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21942,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54558 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21942,
            "unit": "ns/op",
            "extra": "54558 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54558 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54558 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 214013,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4761 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 214013,
            "unit": "ns/op",
            "extra": "4761 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4761 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4761 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2923083,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "424 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2923083,
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
            "value": 32562438,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 32562438,
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
            "value": 9132,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "131100 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9132,
            "unit": "ns/op",
            "extra": "131100 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "131100 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "131100 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6887678,
            "unit": "ns/op\t 4527309 B/op\t   69224 allocs/op",
            "extra": "174 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6887678,
            "unit": "ns/op",
            "extra": "174 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527309,
            "unit": "B/op",
            "extra": "174 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "174 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7607,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "151226 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7607,
            "unit": "ns/op",
            "extra": "151226 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "151226 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "151226 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 835794,
            "unit": "ns/op\t  392697 B/op\t    6167 allocs/op",
            "extra": "1411 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 835794,
            "unit": "ns/op",
            "extra": "1411 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392697,
            "unit": "B/op",
            "extra": "1411 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6167,
            "unit": "allocs/op",
            "extra": "1411 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10634,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "114088 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10634,
            "unit": "ns/op",
            "extra": "114088 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "114088 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "114088 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7916325,
            "unit": "ns/op\t 4913938 B/op\t   75235 allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7916325,
            "unit": "ns/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913938,
            "unit": "B/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 316.3,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3779130 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 316.3,
            "unit": "ns/op",
            "extra": "3779130 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3779130 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3779130 times\n4 procs"
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
            "value": 550.7,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2199116 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 550.7,
            "unit": "ns/op",
            "extra": "2199116 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2199116 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2199116 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 787.3,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1528756 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 787.3,
            "unit": "ns/op",
            "extra": "1528756 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1528756 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1528756 times\n4 procs"
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
          "id": "59fadd5e714b70c7a96f06ae61168ed499f69ca2",
          "message": "chore(deps): bump kconf to v2 (#7614)",
          "timestamp": "2025-07-11T14:18:12Z",
          "tree_id": "9fbdef2f533aa4ef966cacb64b339ee88f505e9d",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/59fadd5e714b70c7a96f06ae61168ed499f69ca2"
        },
        "date": 1752245321275,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 625063,
            "unit": "ns/op\t  819454 B/op\t       5 allocs/op",
            "extra": "1746 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 625063,
            "unit": "ns/op",
            "extra": "1746 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819454,
            "unit": "B/op",
            "extra": "1746 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1746 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 8297,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "158156 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 8297,
            "unit": "ns/op",
            "extra": "158156 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "158156 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "158156 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 107.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11131551 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 107.6,
            "unit": "ns/op",
            "extra": "11131551 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11131551 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11131551 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22032,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54332 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22032,
            "unit": "ns/op",
            "extra": "54332 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54332 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54332 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 215921,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4702 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 215921,
            "unit": "ns/op",
            "extra": "4702 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4702 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4702 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2367878,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "484 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2367878,
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
            "value": 37221630,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "27 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 37221630,
            "unit": "ns/op",
            "extra": "27 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
            "unit": "B/op",
            "extra": "27 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "27 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 11110,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "129529 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 11110,
            "unit": "ns/op",
            "extra": "129529 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "129529 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "129529 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7012044,
            "unit": "ns/op\t 4527316 B/op\t   69224 allocs/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7012044,
            "unit": "ns/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527316,
            "unit": "B/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7717,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "156054 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7717,
            "unit": "ns/op",
            "extra": "156054 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "156054 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "156054 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 855356,
            "unit": "ns/op\t  392709 B/op\t    6167 allocs/op",
            "extra": "1370 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 855356,
            "unit": "ns/op",
            "extra": "1370 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392709,
            "unit": "B/op",
            "extra": "1370 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6167,
            "unit": "allocs/op",
            "extra": "1370 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10754,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "111747 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10754,
            "unit": "ns/op",
            "extra": "111747 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "111747 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "111747 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8130607,
            "unit": "ns/op\t 4913978 B/op\t   75235 allocs/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8130607,
            "unit": "ns/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913978,
            "unit": "B/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 321.2,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3678710 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 321.2,
            "unit": "ns/op",
            "extra": "3678710 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3678710 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3678710 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.28,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.28,
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
            "value": 550.4,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2174737 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 550.4,
            "unit": "ns/op",
            "extra": "2174737 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2174737 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2174737 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 920.3,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1382194 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 920.3,
            "unit": "ns/op",
            "extra": "1382194 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1382194 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1382194 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "aman.gupta@konghq.com",
            "name": "amankong",
            "username": "amankong"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "01c2c29ea22306b50ab5e30ebdec485de47de54a",
          "message": "added harden runner to all the jobs (#7574)",
          "timestamp": "2025-07-14T02:30:57Z",
          "tree_id": "ab46179599128be9025e19a05e1970fa4c654f8b",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/01c2c29ea22306b50ab5e30ebdec485de47de54a"
        },
        "date": 1752462139245,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 592193,
            "unit": "ns/op\t  819454 B/op\t       5 allocs/op",
            "extra": "1836 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 592193,
            "unit": "ns/op",
            "extra": "1836 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819454,
            "unit": "B/op",
            "extra": "1836 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1836 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7621,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "154465 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7621,
            "unit": "ns/op",
            "extra": "154465 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "154465 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "154465 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 107.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11179336 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 107.2,
            "unit": "ns/op",
            "extra": "11179336 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11179336 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11179336 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22622,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "51811 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22622,
            "unit": "ns/op",
            "extra": "51811 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "51811 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "51811 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 217793,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5593 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 217793,
            "unit": "ns/op",
            "extra": "5593 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5593 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5593 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2417965,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "513 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2417965,
            "unit": "ns/op",
            "extra": "513 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "513 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "513 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 41119707,
            "unit": "ns/op\t24010754 B/op\t       2 allocs/op",
            "extra": "63 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 41119707,
            "unit": "ns/op",
            "extra": "63 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010754,
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
            "value": 9650,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "123914 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9650,
            "unit": "ns/op",
            "extra": "123914 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "123914 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "123914 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7190860,
            "unit": "ns/op\t 4527310 B/op\t   69224 allocs/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7190860,
            "unit": "ns/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527310,
            "unit": "B/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7959,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "148387 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7959,
            "unit": "ns/op",
            "extra": "148387 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "148387 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "148387 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 890176,
            "unit": "ns/op\t  392702 B/op\t    6167 allocs/op",
            "extra": "1274 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 890176,
            "unit": "ns/op",
            "extra": "1274 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392702,
            "unit": "B/op",
            "extra": "1274 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6167,
            "unit": "allocs/op",
            "extra": "1274 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10884,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "111488 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10884,
            "unit": "ns/op",
            "extra": "111488 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "111488 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "111488 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8135920,
            "unit": "ns/op\t 4913919 B/op\t   75235 allocs/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8135920,
            "unit": "ns/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913919,
            "unit": "B/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 319.1,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3755640 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 319.1,
            "unit": "ns/op",
            "extra": "3755640 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3755640 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3755640 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.28,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.28,
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
            "value": 575.1,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2097073 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 575.1,
            "unit": "ns/op",
            "extra": "2097073 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2097073 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2097073 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 819.2,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1468036 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 819.2,
            "unit": "ns/op",
            "extra": "1468036 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1468036 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1468036 times\n4 procs"
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
          "id": "616c0fa5ca8654a85d796c28d1cdd136dcd50b70",
          "message": "chore(deps): update dependency golangci/golangci-lint to v2.2.2 (#7616)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2025-07-15T07:57:28Z",
          "tree_id": "8205b287634580fd0b2ecaea37ed3701664b6a81",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/616c0fa5ca8654a85d796c28d1cdd136dcd50b70"
        },
        "date": 1752568135049,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 621312,
            "unit": "ns/op\t  819450 B/op\t       5 allocs/op",
            "extra": "1784 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 621312,
            "unit": "ns/op",
            "extra": "1784 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819450,
            "unit": "B/op",
            "extra": "1784 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1784 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7219,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "166423 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7219,
            "unit": "ns/op",
            "extra": "166423 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "166423 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "166423 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 107.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11174017 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 107.2,
            "unit": "ns/op",
            "extra": "11174017 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11174017 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11174017 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 26937,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52578 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 26937,
            "unit": "ns/op",
            "extra": "52578 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52578 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52578 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 240909,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4816 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 240909,
            "unit": "ns/op",
            "extra": "4816 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4816 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4816 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2393783,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "492 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2393783,
            "unit": "ns/op",
            "extra": "492 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "492 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "492 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 34959899,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 34959899,
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
            "value": 9561,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "124686 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9561,
            "unit": "ns/op",
            "extra": "124686 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "124686 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "124686 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7399163,
            "unit": "ns/op\t 4527314 B/op\t   69224 allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7399163,
            "unit": "ns/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527314,
            "unit": "B/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7700,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "156836 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7700,
            "unit": "ns/op",
            "extra": "156836 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "156836 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "156836 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 861505,
            "unit": "ns/op\t  392705 B/op\t    6167 allocs/op",
            "extra": "1384 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 861505,
            "unit": "ns/op",
            "extra": "1384 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392705,
            "unit": "B/op",
            "extra": "1384 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6167,
            "unit": "allocs/op",
            "extra": "1384 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10966,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "107464 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10966,
            "unit": "ns/op",
            "extra": "107464 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "107464 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "107464 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8601426,
            "unit": "ns/op\t 4913986 B/op\t   75235 allocs/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8601426,
            "unit": "ns/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913986,
            "unit": "B/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 319.7,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3757359 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 319.7,
            "unit": "ns/op",
            "extra": "3757359 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3757359 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3757359 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.28,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.28,
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
            "value": 568.9,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2102905 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 568.9,
            "unit": "ns/op",
            "extra": "2102905 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2102905 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2102905 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 906.8,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1433894 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 906.8,
            "unit": "ns/op",
            "extra": "1433894 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1433894 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1433894 times\n4 procs"
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
          "id": "26a93b3b877f69d413774bf6220ce3b470d7207f",
          "message": "fix integration EE tests with Kong 3.11 (#7620)\n\n* fix integration EE tests with Kong 3.11\n\nSigned-off-by: Jintao Zhang <zhangjintao9020@gmail.com>\n\n* chore(deps): update kong/kong-gateway docker tag to v3.11.0.0\n\n---------\n\nSigned-off-by: Jintao Zhang <zhangjintao9020@gmail.com>\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2025-07-21T02:29:21Z",
          "tree_id": "e3ed56ca85e8fbcceed55db46f68d66d78723bcd",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/26a93b3b877f69d413774bf6220ce3b470d7207f"
        },
        "date": 1753066840439,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 710718,
            "unit": "ns/op\t  819454 B/op\t       5 allocs/op",
            "extra": "1456 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 710718,
            "unit": "ns/op",
            "extra": "1456 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819454,
            "unit": "B/op",
            "extra": "1456 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1456 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7629,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "160926 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7629,
            "unit": "ns/op",
            "extra": "160926 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "160926 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "160926 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 107.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11149105 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 107.7,
            "unit": "ns/op",
            "extra": "11149105 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11149105 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11149105 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21818,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "55032 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21818,
            "unit": "ns/op",
            "extra": "55032 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "55032 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "55032 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 206299,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5077 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 206299,
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
            "value": 2318433,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "513 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2318433,
            "unit": "ns/op",
            "extra": "513 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "513 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "513 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 42307277,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "37 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 42307277,
            "unit": "ns/op",
            "extra": "37 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "37 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "37 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9925,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "104608 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9925,
            "unit": "ns/op",
            "extra": "104608 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "104608 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "104608 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6936181,
            "unit": "ns/op\t 4527313 B/op\t   69224 allocs/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6936181,
            "unit": "ns/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527313,
            "unit": "B/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7709,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "152725 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7709,
            "unit": "ns/op",
            "extra": "152725 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "152725 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "152725 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 849043,
            "unit": "ns/op\t  392710 B/op\t    6167 allocs/op",
            "extra": "1390 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 849043,
            "unit": "ns/op",
            "extra": "1390 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392710,
            "unit": "B/op",
            "extra": "1390 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6167,
            "unit": "allocs/op",
            "extra": "1390 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10711,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "110420 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10711,
            "unit": "ns/op",
            "extra": "110420 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "110420 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "110420 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7951104,
            "unit": "ns/op\t 4914003 B/op\t   75235 allocs/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7951104,
            "unit": "ns/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914003,
            "unit": "B/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 317.5,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3792703 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 317.5,
            "unit": "ns/op",
            "extra": "3792703 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3792703 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3792703 times\n4 procs"
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
            "value": 589.4,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2133200 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 589.4,
            "unit": "ns/op",
            "extra": "2133200 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2133200 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2133200 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 947.3,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1211330 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 947.3,
            "unit": "ns/op",
            "extra": "1211330 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1211330 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1211330 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "patryk.malek@konghq.com",
            "name": "Patryk Małek",
            "username": "pmalek"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "f11890ef76693848914a0fd005fad378e5a38edb",
          "message": "chore: fix CLI args docs generation (#7628)",
          "timestamp": "2025-07-22T08:46:38Z",
          "tree_id": "fe08a127cffeb8cc0d5bc0d4ea8e8e4dba10b2c8",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/f11890ef76693848914a0fd005fad378e5a38edb"
        },
        "date": 1753175813828,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 608857,
            "unit": "ns/op\t  819452 B/op\t       5 allocs/op",
            "extra": "1798 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 608857,
            "unit": "ns/op",
            "extra": "1798 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819452,
            "unit": "B/op",
            "extra": "1798 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1798 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7330,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "163282 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7330,
            "unit": "ns/op",
            "extra": "163282 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "163282 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "163282 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 109.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10990280 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 109.1,
            "unit": "ns/op",
            "extra": "10990280 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "10990280 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "10990280 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21920,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54337 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21920,
            "unit": "ns/op",
            "extra": "54337 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54337 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54337 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 259406,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4672 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 259406,
            "unit": "ns/op",
            "extra": "4672 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4672 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4672 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2390765,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "504 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2390765,
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
            "value": 32613641,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 32613641,
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
            "value": 9551,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "115261 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9551,
            "unit": "ns/op",
            "extra": "115261 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "115261 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "115261 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7146933,
            "unit": "ns/op\t 4527307 B/op\t   69224 allocs/op",
            "extra": "165 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7146933,
            "unit": "ns/op",
            "extra": "165 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527307,
            "unit": "B/op",
            "extra": "165 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "165 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7812,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "154196 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7812,
            "unit": "ns/op",
            "extra": "154196 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "154196 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "154196 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 853887,
            "unit": "ns/op\t  392711 B/op\t    6167 allocs/op",
            "extra": "1364 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 853887,
            "unit": "ns/op",
            "extra": "1364 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392711,
            "unit": "B/op",
            "extra": "1364 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6167,
            "unit": "allocs/op",
            "extra": "1364 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10829,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "111657 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10829,
            "unit": "ns/op",
            "extra": "111657 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "111657 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "111657 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8313102,
            "unit": "ns/op\t 4914001 B/op\t   75235 allocs/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8313102,
            "unit": "ns/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914001,
            "unit": "B/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 321.5,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3725395 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 321.5,
            "unit": "ns/op",
            "extra": "3725395 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3725395 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3725395 times\n4 procs"
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
            "value": 568.3,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2123218 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 568.3,
            "unit": "ns/op",
            "extra": "2123218 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2123218 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2123218 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 806.4,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1494124 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 806.4,
            "unit": "ns/op",
            "extra": "1494124 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1494124 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1494124 times\n4 procs"
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
          "id": "e15688dc25cdb012a93c49dd8bc630aebfd44a4f",
          "message": "chore(deps): bump github/codeql-action from 3.29.2 to 3.29.3 (#7629)\n\n---\nupdated-dependencies:\n- dependency-name: github/codeql-action\n  dependency-version: 3.29.3\n  dependency-type: direct:production\n  update-type: version-update:semver-patch\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2025-07-23T07:00:01Z",
          "tree_id": "e694de15907d7e6394ee0e85b90963ba3d94e4a3",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/e15688dc25cdb012a93c49dd8bc630aebfd44a4f"
        },
        "date": 1753255809343,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 680344,
            "unit": "ns/op\t  819450 B/op\t       5 allocs/op",
            "extra": "1628 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 680344,
            "unit": "ns/op",
            "extra": "1628 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819450,
            "unit": "B/op",
            "extra": "1628 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1628 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7521,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "160117 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7521,
            "unit": "ns/op",
            "extra": "160117 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "160117 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "160117 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 107.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11172901 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 107.4,
            "unit": "ns/op",
            "extra": "11172901 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11172901 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11172901 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22118,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54225 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22118,
            "unit": "ns/op",
            "extra": "54225 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54225 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54225 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 219919,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4784 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 219919,
            "unit": "ns/op",
            "extra": "4784 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4784 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4784 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2382540,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "502 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2382540,
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
            "value": 49528496,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "31 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 49528496,
            "unit": "ns/op",
            "extra": "31 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
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
            "value": 10422,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "97417 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10422,
            "unit": "ns/op",
            "extra": "97417 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "97417 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "97417 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7196321,
            "unit": "ns/op\t 4527309 B/op\t   69224 allocs/op",
            "extra": "166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7196321,
            "unit": "ns/op",
            "extra": "166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527309,
            "unit": "B/op",
            "extra": "166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7708,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "155982 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7708,
            "unit": "ns/op",
            "extra": "155982 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "155982 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "155982 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 887272,
            "unit": "ns/op\t  392700 B/op\t    6167 allocs/op",
            "extra": "1398 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 887272,
            "unit": "ns/op",
            "extra": "1398 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392700,
            "unit": "B/op",
            "extra": "1398 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6167,
            "unit": "allocs/op",
            "extra": "1398 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11473,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "98834 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11473,
            "unit": "ns/op",
            "extra": "98834 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "98834 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "98834 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 9076822,
            "unit": "ns/op\t 4913925 B/op\t   75235 allocs/op",
            "extra": "129 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 9076822,
            "unit": "ns/op",
            "extra": "129 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913925,
            "unit": "B/op",
            "extra": "129 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "129 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 319.4,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3755641 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 319.4,
            "unit": "ns/op",
            "extra": "3755641 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3755641 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3755641 times\n4 procs"
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
            "value": 573,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2143694 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 573,
            "unit": "ns/op",
            "extra": "2143694 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2143694 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2143694 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 1059,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 1059,
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
          "id": "5e49174c57815d20a8ae787abeca2a6726b2bbf4",
          "message": "chore(deps): update dependency go-delve/delve to v1.25.1 (#7626)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2025-07-23T08:03:17Z",
          "tree_id": "fbb6f7f01deb4b6691d5e0b52c860272b5c1c317",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/5e49174c57815d20a8ae787abeca2a6726b2bbf4"
        },
        "date": 1753259615832,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 639025,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "1658 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 639025,
            "unit": "ns/op",
            "extra": "1658 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819453,
            "unit": "B/op",
            "extra": "1658 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1658 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7576,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "157675 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7576,
            "unit": "ns/op",
            "extra": "157675 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "157675 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "157675 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 107.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11163902 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 107.3,
            "unit": "ns/op",
            "extra": "11163902 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11163902 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11163902 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22020,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53887 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22020,
            "unit": "ns/op",
            "extra": "53887 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53887 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53887 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 224955,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5155 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 224955,
            "unit": "ns/op",
            "extra": "5155 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5155 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5155 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2651339,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "410 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2651339,
            "unit": "ns/op",
            "extra": "410 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "410 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "410 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 33232826,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 33232826,
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
            "value": 9442,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "127616 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9442,
            "unit": "ns/op",
            "extra": "127616 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "127616 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "127616 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7004843,
            "unit": "ns/op\t 4527308 B/op\t   69224 allocs/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7004843,
            "unit": "ns/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527308,
            "unit": "B/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7741,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "154542 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7741,
            "unit": "ns/op",
            "extra": "154542 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "154542 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "154542 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 848102,
            "unit": "ns/op\t  392702 B/op\t    6167 allocs/op",
            "extra": "1368 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 848102,
            "unit": "ns/op",
            "extra": "1368 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392702,
            "unit": "B/op",
            "extra": "1368 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6167,
            "unit": "allocs/op",
            "extra": "1368 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10707,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "109911 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10707,
            "unit": "ns/op",
            "extra": "109911 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "109911 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "109911 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8027246,
            "unit": "ns/op\t 4913941 B/op\t   75235 allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8027246,
            "unit": "ns/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913941,
            "unit": "B/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 316.7,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3781806 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 316.7,
            "unit": "ns/op",
            "extra": "3781806 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3781806 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3781806 times\n4 procs"
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
            "value": 556.5,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2159217 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 556.5,
            "unit": "ns/op",
            "extra": "2159217 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2159217 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2159217 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 806.2,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1496785 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 806.2,
            "unit": "ns/op",
            "extra": "1496785 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1496785 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1496785 times\n4 procs"
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
          "id": "b01b273bfaac24c84f4863c10f0df6039da4fb70",
          "message": "chore(deps): update dependency golangci/golangci-lint to v2.3.0 (#7631)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2025-07-23T09:07:17Z",
          "tree_id": "b1839bfd407a72b2cdb2cb192286c3d075f1cdcc",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/b01b273bfaac24c84f4863c10f0df6039da4fb70"
        },
        "date": 1753263474606,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 623213,
            "unit": "ns/op\t  819450 B/op\t       5 allocs/op",
            "extra": "1621 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 623213,
            "unit": "ns/op",
            "extra": "1621 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819450,
            "unit": "B/op",
            "extra": "1621 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1621 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7154,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "165550 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7154,
            "unit": "ns/op",
            "extra": "165550 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "165550 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "165550 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 111.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10473072 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 111.2,
            "unit": "ns/op",
            "extra": "10473072 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "10473072 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "10473072 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21899,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "55089 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21899,
            "unit": "ns/op",
            "extra": "55089 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "55089 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "55089 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 207789,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5073 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 207789,
            "unit": "ns/op",
            "extra": "5073 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
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
            "value": 2390686,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "478 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2390686,
            "unit": "ns/op",
            "extra": "478 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "478 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "478 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 32219377,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 32219377,
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
            "value": 9332,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "127921 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9332,
            "unit": "ns/op",
            "extra": "127921 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "127921 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "127921 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6967296,
            "unit": "ns/op\t 4527309 B/op\t   69224 allocs/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6967296,
            "unit": "ns/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527309,
            "unit": "B/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7628,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "158007 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7628,
            "unit": "ns/op",
            "extra": "158007 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "158007 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "158007 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 845491,
            "unit": "ns/op\t  392698 B/op\t    6167 allocs/op",
            "extra": "1380 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 845491,
            "unit": "ns/op",
            "extra": "1380 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392698,
            "unit": "B/op",
            "extra": "1380 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6167,
            "unit": "allocs/op",
            "extra": "1380 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10775,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "111092 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10775,
            "unit": "ns/op",
            "extra": "111092 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "111092 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "111092 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8037242,
            "unit": "ns/op\t 4913920 B/op\t   75235 allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8037242,
            "unit": "ns/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913920,
            "unit": "B/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 317.2,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3772585 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 317.2,
            "unit": "ns/op",
            "extra": "3772585 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3772585 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3772585 times\n4 procs"
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
            "value": 561.4,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2136722 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 561.4,
            "unit": "ns/op",
            "extra": "2136722 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2136722 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2136722 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 809.7,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1482103 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 809.7,
            "unit": "ns/op",
            "extra": "1482103 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1482103 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1482103 times\n4 procs"
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
          "id": "e2aeee555eae88ce4494ee567778423105c04e13",
          "message": "chore(deps): bump google.golang.org/api from 0.240.0 to 0.243.0 (#7635)\n\nBumps [google.golang.org/api](https://github.com/googleapis/google-api-go-client) from 0.240.0 to 0.243.0.\n- [Release notes](https://github.com/googleapis/google-api-go-client/releases)\n- [Changelog](https://github.com/googleapis/google-api-go-client/blob/main/CHANGES.md)\n- [Commits](https://github.com/googleapis/google-api-go-client/compare/v0.240.0...v0.243.0)\n\n---\nupdated-dependencies:\n- dependency-name: google.golang.org/api\n  dependency-version: 0.243.0\n  dependency-type: direct:production\n  update-type: version-update:semver-minor\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2025-07-23T14:59:33Z",
          "tree_id": "bff4a3fdbc5c5f331b2e129247f76a83b2116007",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/e2aeee555eae88ce4494ee567778423105c04e13"
        },
        "date": 1753284654446,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 628063,
            "unit": "ns/op\t  819455 B/op\t       5 allocs/op",
            "extra": "1762 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 628063,
            "unit": "ns/op",
            "extra": "1762 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819455,
            "unit": "B/op",
            "extra": "1762 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1762 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7013,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "171568 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7013,
            "unit": "ns/op",
            "extra": "171568 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "171568 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "171568 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 107.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11142456 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 107.7,
            "unit": "ns/op",
            "extra": "11142456 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11142456 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11142456 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22115,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52251 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22115,
            "unit": "ns/op",
            "extra": "52251 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52251 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52251 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 211691,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4934 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 211691,
            "unit": "ns/op",
            "extra": "4934 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4934 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4934 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2440686,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "410 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2440686,
            "unit": "ns/op",
            "extra": "410 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "410 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "410 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 33987572,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "68 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 33987572,
            "unit": "ns/op",
            "extra": "68 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "68 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "68 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9368,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "126214 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9368,
            "unit": "ns/op",
            "extra": "126214 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "126214 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "126214 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7013323,
            "unit": "ns/op\t 4527308 B/op\t   69224 allocs/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7013323,
            "unit": "ns/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527308,
            "unit": "B/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7832,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "154884 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7832,
            "unit": "ns/op",
            "extra": "154884 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "154884 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "154884 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 859185,
            "unit": "ns/op\t  392702 B/op\t    6167 allocs/op",
            "extra": "1368 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 859185,
            "unit": "ns/op",
            "extra": "1368 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392702,
            "unit": "B/op",
            "extra": "1368 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6167,
            "unit": "allocs/op",
            "extra": "1368 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10837,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "111079 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10837,
            "unit": "ns/op",
            "extra": "111079 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "111079 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "111079 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8170671,
            "unit": "ns/op\t 4913939 B/op\t   75235 allocs/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8170671,
            "unit": "ns/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913939,
            "unit": "B/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 325.1,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3672325 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 325.1,
            "unit": "ns/op",
            "extra": "3672325 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3672325 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3672325 times\n4 procs"
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
            "value": 557.1,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2136693 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 557.1,
            "unit": "ns/op",
            "extra": "2136693 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2136693 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2136693 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 806.2,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1490773 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 806.2,
            "unit": "ns/op",
            "extra": "1490773 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1490773 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1490773 times\n4 procs"
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
          "id": "930b723b49391da28778448d73f76fc130ba9a27",
          "message": "chore(deps): update dependency kubernetes-sigs/kustomize to v5.7.1 (#7637)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2025-07-24T07:04:03Z",
          "tree_id": "532b855c254caa924f0f92871a9edd8bbf7fea1c",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/930b723b49391da28778448d73f76fc130ba9a27"
        },
        "date": 1753342477034,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 708237,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "1516 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 708237,
            "unit": "ns/op",
            "extra": "1516 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819453,
            "unit": "B/op",
            "extra": "1516 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1516 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7705,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "154946 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7705,
            "unit": "ns/op",
            "extra": "154946 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "154946 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "154946 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 108.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11106499 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 108.2,
            "unit": "ns/op",
            "extra": "11106499 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11106499 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11106499 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21975,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54012 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21975,
            "unit": "ns/op",
            "extra": "54012 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54012 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54012 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 217570,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4696 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 217570,
            "unit": "ns/op",
            "extra": "4696 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4696 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4696 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 3044015,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "382 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 3044015,
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
            "value": 32528890,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 32528890,
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
            "value": 9255,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "126421 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9255,
            "unit": "ns/op",
            "extra": "126421 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "126421 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "126421 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7048569,
            "unit": "ns/op\t 4527314 B/op\t   69224 allocs/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7048569,
            "unit": "ns/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527314,
            "unit": "B/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7615,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "157045 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7615,
            "unit": "ns/op",
            "extra": "157045 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "157045 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "157045 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 854106,
            "unit": "ns/op\t  392713 B/op\t    6167 allocs/op",
            "extra": "1392 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 854106,
            "unit": "ns/op",
            "extra": "1392 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392713,
            "unit": "B/op",
            "extra": "1392 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6167,
            "unit": "allocs/op",
            "extra": "1392 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10651,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "112231 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10651,
            "unit": "ns/op",
            "extra": "112231 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "112231 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "112231 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8184867,
            "unit": "ns/op\t 4913954 B/op\t   75235 allocs/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8184867,
            "unit": "ns/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913954,
            "unit": "B/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 319.1,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3783482 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 319.1,
            "unit": "ns/op",
            "extra": "3783482 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3783482 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3783482 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.26,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.26,
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
            "value": 563.2,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2135820 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 563.2,
            "unit": "ns/op",
            "extra": "2135820 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2135820 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2135820 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 810.6,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1473831 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 810.6,
            "unit": "ns/op",
            "extra": "1473831 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1473831 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1473831 times\n4 procs"
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
          "id": "733e5a89a56239fe4d4aa9fa6bb902a40406c9f4",
          "message": "chore(deps): bump github/codeql-action from 3.29.3 to 3.29.4 (#7636)\n\nBumps [github/codeql-action](https://github.com/github/codeql-action) from 3.29.3 to 3.29.4.\n- [Release notes](https://github.com/github/codeql-action/releases)\n- [Changelog](https://github.com/github/codeql-action/blob/main/CHANGELOG.md)\n- [Commits](https://github.com/github/codeql-action/compare/d6bbdef45e766d081b84a2def353b0055f728d3e...4e828ff8d448a8a6e532957b1811f387a63867e8)\n\n---\nupdated-dependencies:\n- dependency-name: github/codeql-action\n  dependency-version: 3.29.4\n  dependency-type: direct:production\n  update-type: version-update:semver-patch\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2025-07-24T08:15:16Z",
          "tree_id": "ba52e563b651176720d06397367b51af36293b7d",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/733e5a89a56239fe4d4aa9fa6bb902a40406c9f4"
        },
        "date": 1753346796900,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 626774,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "1687 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 626774,
            "unit": "ns/op",
            "extra": "1687 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819453,
            "unit": "B/op",
            "extra": "1687 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1687 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7652,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "151182 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7652,
            "unit": "ns/op",
            "extra": "151182 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "151182 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "151182 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 107.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11152674 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 107.7,
            "unit": "ns/op",
            "extra": "11152674 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11152674 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11152674 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 26257,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "41584 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 26257,
            "unit": "ns/op",
            "extra": "41584 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "41584 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "41584 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 226988,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5228 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 226988,
            "unit": "ns/op",
            "extra": "5228 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5228 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5228 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2435739,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "489 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2435739,
            "unit": "ns/op",
            "extra": "489 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "489 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "489 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 35035004,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "68 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 35035004,
            "unit": "ns/op",
            "extra": "68 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "68 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "68 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 11281,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "124765 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 11281,
            "unit": "ns/op",
            "extra": "124765 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "124765 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "124765 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7850357,
            "unit": "ns/op\t 4527308 B/op\t   69224 allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7850357,
            "unit": "ns/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527308,
            "unit": "B/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7793,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "153452 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7793,
            "unit": "ns/op",
            "extra": "153452 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "153452 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "153452 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 872947,
            "unit": "ns/op\t  392706 B/op\t    6167 allocs/op",
            "extra": "1360 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 872947,
            "unit": "ns/op",
            "extra": "1360 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392706,
            "unit": "B/op",
            "extra": "1360 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6167,
            "unit": "allocs/op",
            "extra": "1360 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11258,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "107460 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11258,
            "unit": "ns/op",
            "extra": "107460 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "107460 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "107460 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8358225,
            "unit": "ns/op\t 4913955 B/op\t   75235 allocs/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8358225,
            "unit": "ns/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913955,
            "unit": "B/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 348.8,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3467030 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 348.8,
            "unit": "ns/op",
            "extra": "3467030 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3467030 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3467030 times\n4 procs"
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
            "value": 563.4,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2115331 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 563.4,
            "unit": "ns/op",
            "extra": "2115331 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2115331 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2115331 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 916.7,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1404088 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 916.7,
            "unit": "ns/op",
            "extra": "1404088 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1404088 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1404088 times\n4 procs"
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
          "id": "2137f2865bccabde6915249706962c3d34629fcf",
          "message": "chore(deps): bump github.com/spf13/pflag from 1.0.6 to 1.0.7 (#7640)\n\n---\nupdated-dependencies:\n- dependency-name: github.com/spf13/pflag\n  dependency-version: 1.0.7\n  dependency-type: direct:production\n  update-type: version-update:semver-patch\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2025-07-24T14:34:52Z",
          "tree_id": "fc55e5b596dae36445d2f93ef770220b290f43e7",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/2137f2865bccabde6915249706962c3d34629fcf"
        },
        "date": 1753369583383,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 754986,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "1455 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 754986,
            "unit": "ns/op",
            "extra": "1455 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819453,
            "unit": "B/op",
            "extra": "1455 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1455 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7384,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "159661 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7384,
            "unit": "ns/op",
            "extra": "159661 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "159661 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "159661 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 107.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11139654 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 107.6,
            "unit": "ns/op",
            "extra": "11139654 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11139654 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11139654 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22222,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53286 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22222,
            "unit": "ns/op",
            "extra": "53286 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53286 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53286 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 225614,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4714 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 225614,
            "unit": "ns/op",
            "extra": "4714 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4714 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4714 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2399943,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "498 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2399943,
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
            "value": 34130435,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "68 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 34130435,
            "unit": "ns/op",
            "extra": "68 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "68 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "68 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9328,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "122340 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9328,
            "unit": "ns/op",
            "extra": "122340 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "122340 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "122340 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6956844,
            "unit": "ns/op\t 4527312 B/op\t   69224 allocs/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6956844,
            "unit": "ns/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527312,
            "unit": "B/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7696,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "154927 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7696,
            "unit": "ns/op",
            "extra": "154927 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "154927 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "154927 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 849242,
            "unit": "ns/op\t  392714 B/op\t    6167 allocs/op",
            "extra": "1394 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 849242,
            "unit": "ns/op",
            "extra": "1394 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392714,
            "unit": "B/op",
            "extra": "1394 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6167,
            "unit": "allocs/op",
            "extra": "1394 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10642,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "114682 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10642,
            "unit": "ns/op",
            "extra": "114682 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "114682 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "114682 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7999938,
            "unit": "ns/op\t 4914017 B/op\t   75235 allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7999938,
            "unit": "ns/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914017,
            "unit": "B/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 322.3,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3770388 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 322.3,
            "unit": "ns/op",
            "extra": "3770388 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3770388 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3770388 times\n4 procs"
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
            "value": 558.9,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2122112 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 558.9,
            "unit": "ns/op",
            "extra": "2122112 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2122112 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2122112 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 801.7,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1500118 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 801.7,
            "unit": "ns/op",
            "extra": "1500118 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1500118 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1500118 times\n4 procs"
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
          "id": "1517e44ffa605d297449a8837f5da4beeee04b7d",
          "message": "chore(deps): update kindest/node docker tag to v1.33.2 (#7630)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2025-07-25T06:19:27Z",
          "tree_id": "b67c12e4535bfae004a7a9c19c9ce1af2b13db3a",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/1517e44ffa605d297449a8837f5da4beeee04b7d"
        },
        "date": 1753426211787,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 769220,
            "unit": "ns/op\t  819449 B/op\t       5 allocs/op",
            "extra": "1383 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 769220,
            "unit": "ns/op",
            "extra": "1383 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819449,
            "unit": "B/op",
            "extra": "1383 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1383 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7608,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "157094 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7608,
            "unit": "ns/op",
            "extra": "157094 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "157094 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "157094 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 107.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11197072 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 107.2,
            "unit": "ns/op",
            "extra": "11197072 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11197072 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11197072 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22088,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53920 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22088,
            "unit": "ns/op",
            "extra": "53920 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53920 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53920 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 212050,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4915 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 212050,
            "unit": "ns/op",
            "extra": "4915 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4915 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4915 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2370329,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "502 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2370329,
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
            "value": 31411366,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 31411366,
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
            "value": 11694,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "99831 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 11694,
            "unit": "ns/op",
            "extra": "99831 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "99831 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "99831 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7087652,
            "unit": "ns/op\t 4527314 B/op\t   69225 allocs/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7087652,
            "unit": "ns/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527314,
            "unit": "B/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69225,
            "unit": "allocs/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7813,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "148342 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7813,
            "unit": "ns/op",
            "extra": "148342 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "148342 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "148342 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 858353,
            "unit": "ns/op\t  392722 B/op\t    6167 allocs/op",
            "extra": "1398 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 858353,
            "unit": "ns/op",
            "extra": "1398 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392722,
            "unit": "B/op",
            "extra": "1398 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6167,
            "unit": "allocs/op",
            "extra": "1398 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10855,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "110522 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10855,
            "unit": "ns/op",
            "extra": "110522 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "110522 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "110522 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8244056,
            "unit": "ns/op\t 4913940 B/op\t   75235 allocs/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8244056,
            "unit": "ns/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913940,
            "unit": "B/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 315.9,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3803946 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 315.9,
            "unit": "ns/op",
            "extra": "3803946 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3803946 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3803946 times\n4 procs"
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
            "value": 570.7,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2088477 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 570.7,
            "unit": "ns/op",
            "extra": "2088477 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2088477 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2088477 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 947.1,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1313913 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 947.1,
            "unit": "ns/op",
            "extra": "1313913 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1313913 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1313913 times\n4 procs"
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
          "id": "a278f8c32c5cf0b29b2bc724968ecd3495341d42",
          "message": "chore(deps): update dependency gke to v1.33.2 (#7555)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>\nCo-authored-by: Tao Yi <tao.yi@konghq.com>",
          "timestamp": "2025-07-25T07:53:35Z",
          "tree_id": "fdc5cdd46ec66e70998d31fcc92f4e3203898fa8",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/a278f8c32c5cf0b29b2bc724968ecd3495341d42"
        },
        "date": 1753431838577,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 717081,
            "unit": "ns/op\t  819454 B/op\t       5 allocs/op",
            "extra": "1494 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 717081,
            "unit": "ns/op",
            "extra": "1494 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819454,
            "unit": "B/op",
            "extra": "1494 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1494 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7581,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "139002 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7581,
            "unit": "ns/op",
            "extra": "139002 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "139002 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "139002 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 106.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11201538 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 106.9,
            "unit": "ns/op",
            "extra": "11201538 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11201538 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11201538 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 27700,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54512 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 27700,
            "unit": "ns/op",
            "extra": "54512 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54512 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54512 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 210119,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5659 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 210119,
            "unit": "ns/op",
            "extra": "5659 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5659 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5659 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2305540,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "514 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2305540,
            "unit": "ns/op",
            "extra": "514 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "514 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "514 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 36428557,
            "unit": "ns/op\t24010754 B/op\t       2 allocs/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 36428557,
            "unit": "ns/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010754,
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
            "value": 9303,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "128916 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9303,
            "unit": "ns/op",
            "extra": "128916 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "128916 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "128916 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7078171,
            "unit": "ns/op\t 4527310 B/op\t   69224 allocs/op",
            "extra": "168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7078171,
            "unit": "ns/op",
            "extra": "168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527310,
            "unit": "B/op",
            "extra": "168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7762,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "156578 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7762,
            "unit": "ns/op",
            "extra": "156578 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "156578 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "156578 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 873179,
            "unit": "ns/op\t  392709 B/op\t    6167 allocs/op",
            "extra": "1365 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 873179,
            "unit": "ns/op",
            "extra": "1365 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392709,
            "unit": "B/op",
            "extra": "1365 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6167,
            "unit": "allocs/op",
            "extra": "1365 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10746,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "112383 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10746,
            "unit": "ns/op",
            "extra": "112383 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "112383 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "112383 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7992059,
            "unit": "ns/op\t 4913925 B/op\t   75235 allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7992059,
            "unit": "ns/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913925,
            "unit": "B/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75235,
            "unit": "allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 317.1,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3772186 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 317.1,
            "unit": "ns/op",
            "extra": "3772186 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3772186 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3772186 times\n4 procs"
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
            "value": 559.3,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2152591 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 559.3,
            "unit": "ns/op",
            "extra": "2152591 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2152591 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2152591 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 992.1,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1228261 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 992.1,
            "unit": "ns/op",
            "extra": "1228261 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1228261 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1228261 times\n4 procs"
          }
        ]
      }
    ]
  }
}