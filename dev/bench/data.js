window.BENCHMARK_DATA = {
  "lastUpdate": 1771257070423,
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
          "id": "e37c1cdc087dc0d385718d8ddb0bf14e03ae9c50",
          "message": "fix(konnect): Use correct label value of secret to store license (#7648)\n\n* fix label value of secret to store license\n\n* Update internal/labels/labels.go\n\nCo-authored-by: Patryk Małek <patryk.malek@konghq.com>\n\n* check secret for stored license\n\n---------\n\nCo-authored-by: Patryk Małek <patryk.malek@konghq.com>",
          "timestamp": "2025-07-30T08:10:20Z",
          "tree_id": "85490c64410f845966f14815dcad724b9dfd0d71",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/e37c1cdc087dc0d385718d8ddb0bf14e03ae9c50"
        },
        "date": 1753864883069,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 696356,
            "unit": "ns/op\t  819455 B/op\t       5 allocs/op",
            "extra": "1448 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 696356,
            "unit": "ns/op",
            "extra": "1448 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819455,
            "unit": "B/op",
            "extra": "1448 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1448 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7513,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "159117 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7513,
            "unit": "ns/op",
            "extra": "159117 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "159117 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "159117 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 108.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11141874 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 108.6,
            "unit": "ns/op",
            "extra": "11141874 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11141874 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11141874 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 26173,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "41896 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 26173,
            "unit": "ns/op",
            "extra": "41896 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "41896 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "41896 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 210019,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5109 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 210019,
            "unit": "ns/op",
            "extra": "5109 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5109 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5109 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2411905,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "476 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2411905,
            "unit": "ns/op",
            "extra": "476 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "476 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "476 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 38865881,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "37 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 38865881,
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
            "value": 9291,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "128754 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9291,
            "unit": "ns/op",
            "extra": "128754 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "128754 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "128754 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7810930,
            "unit": "ns/op\t 4527307 B/op\t   69224 allocs/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7810930,
            "unit": "ns/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527307,
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
            "value": 8830,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "123416 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8830,
            "unit": "ns/op",
            "extra": "123416 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "123416 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "123416 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 849286,
            "unit": "ns/op\t  392686 B/op\t    6167 allocs/op",
            "extra": "1424 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 849286,
            "unit": "ns/op",
            "extra": "1424 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392686,
            "unit": "B/op",
            "extra": "1424 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6167,
            "unit": "allocs/op",
            "extra": "1424 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10728,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "111650 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10728,
            "unit": "ns/op",
            "extra": "111650 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "111650 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "111650 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8034624,
            "unit": "ns/op\t 4913912 B/op\t   75235 allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8034624,
            "unit": "ns/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913912,
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
            "value": 317.2,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3796894 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 317.2,
            "unit": "ns/op",
            "extra": "3796894 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3796894 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3796894 times\n4 procs"
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
            "value": 555.5,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2165150 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 555.5,
            "unit": "ns/op",
            "extra": "2165150 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2165150 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2165150 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 807,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1489507 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 807,
            "unit": "ns/op",
            "extra": "1489507 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1489507 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1489507 times\n4 procs"
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
          "id": "5eea70b746aa44a696e799f78da0afda97c4894c",
          "message": "chore: fix kubernetes-configuration CRD installation in integration tests (#7659)",
          "timestamp": "2025-08-01T09:45:27Z",
          "tree_id": "59e5065ec9d61774950c0628a9f5e6c972a22585",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/5eea70b746aa44a696e799f78da0afda97c4894c"
        },
        "date": 1754043374027,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 694719,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "1516 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 694719,
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
            "value": 7541,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "149641 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7541,
            "unit": "ns/op",
            "extra": "149641 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "149641 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "149641 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 107.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11159410 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 107.5,
            "unit": "ns/op",
            "extra": "11159410 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11159410 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11159410 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21914,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54302 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21914,
            "unit": "ns/op",
            "extra": "54302 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54302 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54302 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 211727,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4754 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 211727,
            "unit": "ns/op",
            "extra": "4754 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4754 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4754 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2492708,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "475 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2492708,
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
            "value": 37061027,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "28 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 37061027,
            "unit": "ns/op",
            "extra": "28 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
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
            "value": 9300,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "129043 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9300,
            "unit": "ns/op",
            "extra": "129043 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "129043 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "129043 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 8344070,
            "unit": "ns/op\t 4527306 B/op\t   69224 allocs/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 8344070,
            "unit": "ns/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527306,
            "unit": "B/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8002,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "147844 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8002,
            "unit": "ns/op",
            "extra": "147844 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "147844 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "147844 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 849858,
            "unit": "ns/op\t  392691 B/op\t    6167 allocs/op",
            "extra": "1383 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 849858,
            "unit": "ns/op",
            "extra": "1383 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392691,
            "unit": "B/op",
            "extra": "1383 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6167,
            "unit": "allocs/op",
            "extra": "1383 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10654,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "113116 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10654,
            "unit": "ns/op",
            "extra": "113116 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "113116 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "113116 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8043953,
            "unit": "ns/op\t 4913976 B/op\t   75235 allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8043953,
            "unit": "ns/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913976,
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
            "value": 319.8,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3753456 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 319.8,
            "unit": "ns/op",
            "extra": "3753456 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3753456 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3753456 times\n4 procs"
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
            "value": 561,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2136082 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 561,
            "unit": "ns/op",
            "extra": "2136082 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2136082 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2136082 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 807,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1488049 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 807,
            "unit": "ns/op",
            "extra": "1488049 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1488049 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1488049 times\n4 procs"
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
          "id": "53c0947238ae1dc338b5b09ec69dbeea1fca8b2f",
          "message": "chore(docs): Update changelog of 3.5.1 and 3.4.8 in main branch (#7667)\n\n* update changelog for 3.5.1 (#7657)\n\n* update changelog for 3.4.8",
          "timestamp": "2025-08-04T13:42:40Z",
          "tree_id": "07060948494fa234662a4efe61e955323d438a62",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/53c0947238ae1dc338b5b09ec69dbeea1fca8b2f"
        },
        "date": 1754316772151,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 684928,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "1530 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 684928,
            "unit": "ns/op",
            "extra": "1530 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819453,
            "unit": "B/op",
            "extra": "1530 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1530 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7621,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "136198 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7621,
            "unit": "ns/op",
            "extra": "136198 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "136198 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "136198 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 107.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11162584 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 107.5,
            "unit": "ns/op",
            "extra": "11162584 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11162584 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11162584 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21985,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53226 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21985,
            "unit": "ns/op",
            "extra": "53226 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53226 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53226 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 216846,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4911 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 216846,
            "unit": "ns/op",
            "extra": "4911 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4911 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4911 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2375267,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "480 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2375267,
            "unit": "ns/op",
            "extra": "480 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "480 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "480 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 36868197,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 36868197,
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
            "value": 9127,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "130417 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9127,
            "unit": "ns/op",
            "extra": "130417 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "130417 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "130417 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7009740,
            "unit": "ns/op\t 4527310 B/op\t   69224 allocs/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7009740,
            "unit": "ns/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527310,
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
            "value": 7506,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "156704 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7506,
            "unit": "ns/op",
            "extra": "156704 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "156704 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "156704 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 843675,
            "unit": "ns/op\t  392703 B/op\t    6167 allocs/op",
            "extra": "1398 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 843675,
            "unit": "ns/op",
            "extra": "1398 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392703,
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
            "value": 10521,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "113938 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10521,
            "unit": "ns/op",
            "extra": "113938 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "113938 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "113938 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8245026,
            "unit": "ns/op\t 4913907 B/op\t   75235 allocs/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8245026,
            "unit": "ns/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913907,
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
            "value": 318.5,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3754712 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 318.5,
            "unit": "ns/op",
            "extra": "3754712 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3754712 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3754712 times\n4 procs"
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
            "value": 556.3,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2152322 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 556.3,
            "unit": "ns/op",
            "extra": "2152322 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2152322 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2152322 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 811.3,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1474071 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 811.3,
            "unit": "ns/op",
            "extra": "1474071 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1474071 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1474071 times\n4 procs"
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
          "id": "73ddb33115c2390db47e05d0e6d3018dd32b16f8",
          "message": "chore(deps): update dependency golangci/golangci-lint to v2.3.1 (#7669)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2025-08-05T03:01:08Z",
          "tree_id": "cfb62c4fc399a72f32828d4450a530544edd503e",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/73ddb33115c2390db47e05d0e6d3018dd32b16f8"
        },
        "date": 1754364673963,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 643425,
            "unit": "ns/op\t  819454 B/op\t       5 allocs/op",
            "extra": "1568 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 643425,
            "unit": "ns/op",
            "extra": "1568 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819454,
            "unit": "B/op",
            "extra": "1568 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1568 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7444,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "138627 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7444,
            "unit": "ns/op",
            "extra": "138627 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "138627 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "138627 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 107.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11169178 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 107.3,
            "unit": "ns/op",
            "extra": "11169178 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11169178 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11169178 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22157,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52971 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22157,
            "unit": "ns/op",
            "extra": "52971 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52971 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52971 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 213151,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5128 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 213151,
            "unit": "ns/op",
            "extra": "5128 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5128 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5128 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2834933,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "440 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2834933,
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
            "value": 34635940,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "30 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 34635940,
            "unit": "ns/op",
            "extra": "30 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 9411,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "127622 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9411,
            "unit": "ns/op",
            "extra": "127622 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "127622 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "127622 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7121289,
            "unit": "ns/op\t 4527306 B/op\t   69224 allocs/op",
            "extra": "168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7121289,
            "unit": "ns/op",
            "extra": "168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527306,
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
            "value": 7770,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "155367 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7770,
            "unit": "ns/op",
            "extra": "155367 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "155367 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "155367 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 852975,
            "unit": "ns/op\t  392719 B/op\t    6167 allocs/op",
            "extra": "1366 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 852975,
            "unit": "ns/op",
            "extra": "1366 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392719,
            "unit": "B/op",
            "extra": "1366 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6167,
            "unit": "allocs/op",
            "extra": "1366 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10861,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "111230 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10861,
            "unit": "ns/op",
            "extra": "111230 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "111230 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "111230 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8173368,
            "unit": "ns/op\t 4914017 B/op\t   75235 allocs/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8173368,
            "unit": "ns/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914017,
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
            "value": 316.9,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3788043 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 316.9,
            "unit": "ns/op",
            "extra": "3788043 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3788043 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3788043 times\n4 procs"
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
            "value": 561.3,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2156896 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 561.3,
            "unit": "ns/op",
            "extra": "2156896 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2156896 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2156896 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 809.8,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1476607 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 809.8,
            "unit": "ns/op",
            "extra": "1476607 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1476607 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1476607 times\n4 procs"
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
          "id": "110a07310a933ecf830f080ba1c47919d7ceb623",
          "message": "chore(deps): bump github.com/docker/docker (#7649)\n\nBumps [github.com/docker/docker](https://github.com/docker/docker) from 28.3.1+incompatible to 28.3.3+incompatible.\n- [Release notes](https://github.com/docker/docker/releases)\n- [Commits](https://github.com/docker/docker/compare/v28.3.1...v28.3.3)\n\n---\nupdated-dependencies:\n- dependency-name: github.com/docker/docker\n  dependency-version: 28.3.3+incompatible\n  dependency-type: direct:production\n  update-type: version-update:semver-patch\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2025-08-05T08:40:40Z",
          "tree_id": "a1c8cc031cb1ef770ed995f20dcc8b310f1981a3",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/110a07310a933ecf830f080ba1c47919d7ceb623"
        },
        "date": 1754385113381,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 647720,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "1586 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 647720,
            "unit": "ns/op",
            "extra": "1586 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819453,
            "unit": "B/op",
            "extra": "1586 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1586 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7082,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "156190 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7082,
            "unit": "ns/op",
            "extra": "156190 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "156190 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "156190 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 107.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11150566 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 107.6,
            "unit": "ns/op",
            "extra": "11150566 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11150566 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11150566 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22268,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53461 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22268,
            "unit": "ns/op",
            "extra": "53461 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53461 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53461 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 226974,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4809 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 226974,
            "unit": "ns/op",
            "extra": "4809 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4809 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4809 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2336980,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "518 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2336980,
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
            "value": 36759413,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "68 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 36759413,
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
            "value": 9565,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "121976 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9565,
            "unit": "ns/op",
            "extra": "121976 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "121976 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "121976 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7327236,
            "unit": "ns/op\t 4527309 B/op\t   69224 allocs/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7327236,
            "unit": "ns/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527309,
            "unit": "B/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7872,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "150181 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7872,
            "unit": "ns/op",
            "extra": "150181 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "150181 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "150181 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 868433,
            "unit": "ns/op\t  392693 B/op\t    6167 allocs/op",
            "extra": "1368 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 868433,
            "unit": "ns/op",
            "extra": "1368 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392693,
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
            "value": 11178,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "108210 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11178,
            "unit": "ns/op",
            "extra": "108210 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "108210 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "108210 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8368964,
            "unit": "ns/op\t 4913956 B/op\t   75235 allocs/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8368964,
            "unit": "ns/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913956,
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
            "value": 319.4,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3700482 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 319.4,
            "unit": "ns/op",
            "extra": "3700482 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3700482 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3700482 times\n4 procs"
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
            "value": 562.3,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2132424 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 562.3,
            "unit": "ns/op",
            "extra": "2132424 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2132424 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2132424 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 811.1,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1475484 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 811.1,
            "unit": "ns/op",
            "extra": "1475484 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1475484 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1475484 times\n4 procs"
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
          "id": "019c46f48b01a91611a7952581aa713223077206",
          "message": "chore(deps): bump ruby/setup-ruby from 1.245.0 to 1.254.0 (#7655)\n\nBumps [ruby/setup-ruby](https://github.com/ruby/setup-ruby) from 1.245.0 to 1.254.0.\n- [Release notes](https://github.com/ruby/setup-ruby/releases)\n- [Changelog](https://github.com/ruby/setup-ruby/blob/master/release.rb)\n- [Commits](https://github.com/ruby/setup-ruby/compare/a4effe49ee8ee5b8b5091268c473a4628afb5651...2a7b30092b0caf9c046252510f9273b4875f3db9)\n\n---\nupdated-dependencies:\n- dependency-name: ruby/setup-ruby\n  dependency-version: 1.254.0\n  dependency-type: direct:production\n  update-type: version-update:semver-minor\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2025-08-05T08:40:58Z",
          "tree_id": "127b82ff91fbe53f7a4d32adf6593c1bbce874f6",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/019c46f48b01a91611a7952581aa713223077206"
        },
        "date": 1754385443249,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 614378,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "1647 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 614378,
            "unit": "ns/op",
            "extra": "1647 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819453,
            "unit": "B/op",
            "extra": "1647 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1647 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7167,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "169462 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7167,
            "unit": "ns/op",
            "extra": "169462 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "169462 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "169462 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 107.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11150754 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 107.6,
            "unit": "ns/op",
            "extra": "11150754 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11150754 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11150754 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22058,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53866 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22058,
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
            "value": 210111,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5114 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 210111,
            "unit": "ns/op",
            "extra": "5114 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5114 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5114 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2446597,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "435 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2446597,
            "unit": "ns/op",
            "extra": "435 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "435 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "435 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 32192345,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 32192345,
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
            "value": 9071,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "130774 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9071,
            "unit": "ns/op",
            "extra": "130774 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "130774 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "130774 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7111687,
            "unit": "ns/op\t 4527311 B/op\t   69224 allocs/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7111687,
            "unit": "ns/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527311,
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
            "value": 7588,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "161587 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7588,
            "unit": "ns/op",
            "extra": "161587 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "161587 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "161587 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 842878,
            "unit": "ns/op\t  392715 B/op\t    6167 allocs/op",
            "extra": "1387 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 842878,
            "unit": "ns/op",
            "extra": "1387 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392715,
            "unit": "B/op",
            "extra": "1387 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6167,
            "unit": "allocs/op",
            "extra": "1387 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10468,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "115765 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10468,
            "unit": "ns/op",
            "extra": "115765 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "115765 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "115765 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8062949,
            "unit": "ns/op\t 4913937 B/op\t   75235 allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8062949,
            "unit": "ns/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913937,
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
            "value": 388.4,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3602799 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 388.4,
            "unit": "ns/op",
            "extra": "3602799 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3602799 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3602799 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.62,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "112663648 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.62,
            "unit": "ns/op",
            "extra": "112663648 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "112663648 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "112663648 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 555,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2163621 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 555,
            "unit": "ns/op",
            "extra": "2163621 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2163621 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2163621 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 808,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1472799 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 808,
            "unit": "ns/op",
            "extra": "1472799 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1472799 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1472799 times\n4 procs"
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
          "id": "36078b6de4466cdefe4f373d6aecdab64fc72788",
          "message": "chore(deps): bump distroless/static from `627d6c5` to `cdf4daa` (#7668)\n\nBumps distroless/static from `627d6c5` to `cdf4daa`.\n\n---\nupdated-dependencies:\n- dependency-name: distroless/static\n  dependency-version: nonroot\n  dependency-type: direct:production\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2025-08-05T08:41:27Z",
          "tree_id": "6a6254d2a5c4504d4975b7c336d77ac73cfa4cda",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/36078b6de4466cdefe4f373d6aecdab64fc72788"
        },
        "date": 1754386759158,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 865834,
            "unit": "ns/op\t  819454 B/op\t       5 allocs/op",
            "extra": "1188 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 865834,
            "unit": "ns/op",
            "extra": "1188 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819454,
            "unit": "B/op",
            "extra": "1188 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1188 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 9909,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "116785 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 9909,
            "unit": "ns/op",
            "extra": "116785 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "116785 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "116785 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 107.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11189950 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 107.1,
            "unit": "ns/op",
            "extra": "11189950 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11189950 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11189950 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21910,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54397 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21910,
            "unit": "ns/op",
            "extra": "54397 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54397 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54397 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 222729,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5704 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 222729,
            "unit": "ns/op",
            "extra": "5704 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5704 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5704 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2453422,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "477 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2453422,
            "unit": "ns/op",
            "extra": "477 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
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
            "value": 39087705,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "67 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 39087705,
            "unit": "ns/op",
            "extra": "67 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "67 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "67 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9496,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "124897 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9496,
            "unit": "ns/op",
            "extra": "124897 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "124897 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "124897 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7324918,
            "unit": "ns/op\t 4527311 B/op\t   69224 allocs/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7324918,
            "unit": "ns/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527311,
            "unit": "B/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69224,
            "unit": "allocs/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7812,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "152000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7812,
            "unit": "ns/op",
            "extra": "152000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "152000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "152000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 853229,
            "unit": "ns/op\t  392710 B/op\t    6167 allocs/op",
            "extra": "1389 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 853229,
            "unit": "ns/op",
            "extra": "1389 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392710,
            "unit": "B/op",
            "extra": "1389 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6167,
            "unit": "allocs/op",
            "extra": "1389 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10837,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "111655 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10837,
            "unit": "ns/op",
            "extra": "111655 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "111655 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "111655 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8470199,
            "unit": "ns/op\t 4913932 B/op\t   75235 allocs/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8470199,
            "unit": "ns/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913932,
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
            "value": 320.5,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3780265 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 320.5,
            "unit": "ns/op",
            "extra": "3780265 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3780265 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3780265 times\n4 procs"
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
            "value": 730.4,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "1641943 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 730.4,
            "unit": "ns/op",
            "extra": "1641943 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "1641943 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1641943 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 840.2,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1398384 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 840.2,
            "unit": "ns/op",
            "extra": "1398384 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1398384 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1398384 times\n4 procs"
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
          "id": "6b33d44fd79432beacc4ba654726a3d679840e0c",
          "message": "chore(deps): bump cloud.google.com/go/container from 1.43.0 to 1.44.0 (#7671)\n\nBumps [cloud.google.com/go/container](https://github.com/googleapis/google-cloud-go) from 1.43.0 to 1.44.0.\n- [Release notes](https://github.com/googleapis/google-cloud-go/releases)\n- [Changelog](https://github.com/googleapis/google-cloud-go/blob/main/CHANGES.md)\n- [Commits](https://github.com/googleapis/google-cloud-go/compare/pubsub/v1.43.0...pubsub/v1.44.0)\n\n---\nupdated-dependencies:\n- dependency-name: cloud.google.com/go/container\n  dependency-version: 1.44.0\n  dependency-type: direct:production\n  update-type: version-update:semver-minor\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2025-08-06T13:17:58Z",
          "tree_id": "b32e0ea3497280023e32a5cb69111e59fc54e696",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/6b33d44fd79432beacc4ba654726a3d679840e0c"
        },
        "date": 1754488181764,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 642800,
            "unit": "ns/op\t  819454 B/op\t       5 allocs/op",
            "extra": "1632 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 642800,
            "unit": "ns/op",
            "extra": "1632 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819454,
            "unit": "B/op",
            "extra": "1632 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1632 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7092,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "161986 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7092,
            "unit": "ns/op",
            "extra": "161986 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "161986 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "161986 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 107.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11200071 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 107.2,
            "unit": "ns/op",
            "extra": "11200071 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11200071 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11200071 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22351,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52549 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22351,
            "unit": "ns/op",
            "extra": "52549 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52549 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52549 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 218340,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4654 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 218340,
            "unit": "ns/op",
            "extra": "4654 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4654 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4654 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2407675,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "496 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2407675,
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
            "value": 35000750,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "68 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 35000750,
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
            "value": 9268,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "130582 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9268,
            "unit": "ns/op",
            "extra": "130582 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "130582 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "130582 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7172400,
            "unit": "ns/op\t 4527313 B/op\t   69224 allocs/op",
            "extra": "165 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7172400,
            "unit": "ns/op",
            "extra": "165 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527313,
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
            "value": 7693,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "156189 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7693,
            "unit": "ns/op",
            "extra": "156189 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "156189 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "156189 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 861168,
            "unit": "ns/op\t  392706 B/op\t    6167 allocs/op",
            "extra": "1371 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 861168,
            "unit": "ns/op",
            "extra": "1371 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392706,
            "unit": "B/op",
            "extra": "1371 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6167,
            "unit": "allocs/op",
            "extra": "1371 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10755,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "113430 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10755,
            "unit": "ns/op",
            "extra": "113430 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "113430 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "113430 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8301858,
            "unit": "ns/op\t 4913954 B/op\t   75235 allocs/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8301858,
            "unit": "ns/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913954,
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
            "value": 364.3,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3763196 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 364.3,
            "unit": "ns/op",
            "extra": "3763196 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3763196 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3763196 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.62,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "112832023 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.62,
            "unit": "ns/op",
            "extra": "112832023 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "112832023 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "112832023 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 559.9,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2129930 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 559.9,
            "unit": "ns/op",
            "extra": "2129930 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2129930 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2129930 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 817.7,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1461655 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 817.7,
            "unit": "ns/op",
            "extra": "1461655 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1461655 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1461655 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "Guaris@users.noreply.github.com",
            "name": "Angel",
            "username": "Guaris"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "2462b7e98eed8f9511b6adc9f8a568675de7d66f",
          "message": "Fix: update links in readme (#7690)",
          "timestamp": "2025-08-18T21:04:11+02:00",
          "tree_id": "21d19a29eec41f390bf539abe2face4eb14eda34",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/2462b7e98eed8f9511b6adc9f8a568675de7d66f"
        },
        "date": 1755544074055,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 589332,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "1845 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 589332,
            "unit": "ns/op",
            "extra": "1845 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819453,
            "unit": "B/op",
            "extra": "1845 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1845 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7539,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "156552 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7539,
            "unit": "ns/op",
            "extra": "156552 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "156552 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "156552 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 107.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11146158 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 107.7,
            "unit": "ns/op",
            "extra": "11146158 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11146158 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11146158 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22071,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53810 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22071,
            "unit": "ns/op",
            "extra": "53810 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53810 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53810 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 216192,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5185 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 216192,
            "unit": "ns/op",
            "extra": "5185 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5185 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5185 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2392717,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "478 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2392717,
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
            "value": 36170639,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "36 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 36170639,
            "unit": "ns/op",
            "extra": "36 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "36 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "36 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9232,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "128796 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9232,
            "unit": "ns/op",
            "extra": "128796 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "128796 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "128796 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7049207,
            "unit": "ns/op\t 4527305 B/op\t   69224 allocs/op",
            "extra": "168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7049207,
            "unit": "ns/op",
            "extra": "168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527305,
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
            "value": 7729,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "154578 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7729,
            "unit": "ns/op",
            "extra": "154578 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "154578 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "154578 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 846505,
            "unit": "ns/op\t  392701 B/op\t    6167 allocs/op",
            "extra": "1396 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 846505,
            "unit": "ns/op",
            "extra": "1396 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392701,
            "unit": "B/op",
            "extra": "1396 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6167,
            "unit": "allocs/op",
            "extra": "1396 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10749,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "113148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10749,
            "unit": "ns/op",
            "extra": "113148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "113148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "113148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8111956,
            "unit": "ns/op\t 4913948 B/op\t   75235 allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8111956,
            "unit": "ns/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913948,
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
            "value": 342.6,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3515420 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 342.6,
            "unit": "ns/op",
            "extra": "3515420 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3515420 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3515420 times\n4 procs"
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
            "value": 563,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2115327 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 563,
            "unit": "ns/op",
            "extra": "2115327 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2115327 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2115327 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 814.3,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1476292 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 814.3,
            "unit": "ns/op",
            "extra": "1476292 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1476292 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1476292 times\n4 procs"
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
          "id": "c29a8c84afdb1d6a04b93789b6a0faf8f9cf24e3",
          "message": "chore(deps): update dependency gke to v1.33.4 (#7682)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2025-08-21T13:23:12Z",
          "tree_id": "78172b4305aacae73fe3637da5cb75ee3cd21f9d",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/c29a8c84afdb1d6a04b93789b6a0faf8f9cf24e3"
        },
        "date": 1755784503818,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 699717,
            "unit": "ns/op\t  819450 B/op\t       5 allocs/op",
            "extra": "1437 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 699717,
            "unit": "ns/op",
            "extra": "1437 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819450,
            "unit": "B/op",
            "extra": "1437 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1437 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7617,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "131632 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7617,
            "unit": "ns/op",
            "extra": "131632 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "131632 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "131632 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 109.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10901620 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 109.9,
            "unit": "ns/op",
            "extra": "10901620 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "10901620 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "10901620 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21806,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54466 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21806,
            "unit": "ns/op",
            "extra": "54466 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54466 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54466 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 207440,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4935 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 207440,
            "unit": "ns/op",
            "extra": "4935 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4935 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4935 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2384770,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "465 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2384770,
            "unit": "ns/op",
            "extra": "465 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
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
            "value": 38062254,
            "unit": "ns/op\t24010754 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 38062254,
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
            "value": 9330,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "129662 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9330,
            "unit": "ns/op",
            "extra": "129662 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "129662 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "129662 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6987913,
            "unit": "ns/op\t 4527310 B/op\t   69224 allocs/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6987913,
            "unit": "ns/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527310,
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
            "value": 7706,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "157488 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7706,
            "unit": "ns/op",
            "extra": "157488 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "157488 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "157488 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 850508,
            "unit": "ns/op\t  392694 B/op\t    6167 allocs/op",
            "extra": "1370 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 850508,
            "unit": "ns/op",
            "extra": "1370 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392694,
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
            "value": 10779,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "111991 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10779,
            "unit": "ns/op",
            "extra": "111991 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "111991 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "111991 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8018313,
            "unit": "ns/op\t 4913936 B/op\t   75235 allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8018313,
            "unit": "ns/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913936,
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
            "value": 370.6,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3791442 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 370.6,
            "unit": "ns/op",
            "extra": "3791442 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3791442 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3791442 times\n4 procs"
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
            "value": 564.2,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2133927 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 564.2,
            "unit": "ns/op",
            "extra": "2133927 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2133927 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2133927 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 804.5,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1489683 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 804.5,
            "unit": "ns/op",
            "extra": "1489683 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1489683 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1489683 times\n4 procs"
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
          "id": "7342711be211b52fe7bef04c67b0a8247f748d4f",
          "message": "set postgresql image repo to bitnamilegacy/postgresql (#7710)",
          "timestamp": "2025-09-18T07:48:45Z",
          "tree_id": "6aeb4aa53104e19d9efadc79f366fecd106e3da6",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/7342711be211b52fe7bef04c67b0a8247f748d4f"
        },
        "date": 1758183635937,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 652358,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "1544 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 652358,
            "unit": "ns/op",
            "extra": "1544 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819453,
            "unit": "B/op",
            "extra": "1544 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1544 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7600,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "142314 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7600,
            "unit": "ns/op",
            "extra": "142314 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "142314 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "142314 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 107,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11197525 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 107,
            "unit": "ns/op",
            "extra": "11197525 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11197525 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11197525 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21916,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54144 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21916,
            "unit": "ns/op",
            "extra": "54144 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54144 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54144 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 206586,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5604 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 206586,
            "unit": "ns/op",
            "extra": "5604 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5604 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5604 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2350654,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "513 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2350654,
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
            "value": 31547232,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 31547232,
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
            "value": 9308,
            "unit": "ns/op\t    7720 B/op\t     179 allocs/op",
            "extra": "128149 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9308,
            "unit": "ns/op",
            "extra": "128149 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7720,
            "unit": "B/op",
            "extra": "128149 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 179,
            "unit": "allocs/op",
            "extra": "128149 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6996284,
            "unit": "ns/op\t 4527310 B/op\t   69224 allocs/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6996284,
            "unit": "ns/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527310,
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
            "value": 7730,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "155731 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7730,
            "unit": "ns/op",
            "extra": "155731 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "155731 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "155731 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 845822,
            "unit": "ns/op\t  392719 B/op\t    6167 allocs/op",
            "extra": "1388 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 845822,
            "unit": "ns/op",
            "extra": "1388 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392719,
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
            "value": 10722,
            "unit": "ns/op\t    8136 B/op\t     187 allocs/op",
            "extra": "111110 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10722,
            "unit": "ns/op",
            "extra": "111110 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8136,
            "unit": "B/op",
            "extra": "111110 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 187,
            "unit": "allocs/op",
            "extra": "111110 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7970456,
            "unit": "ns/op\t 4913969 B/op\t   75235 allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7970456,
            "unit": "ns/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4913969,
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
            "value": 318.9,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3802120 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 318.9,
            "unit": "ns/op",
            "extra": "3802120 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3802120 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3802120 times\n4 procs"
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
            "value": 558.7,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2147779 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 558.7,
            "unit": "ns/op",
            "extra": "2147779 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2147779 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2147779 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 803.8,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1497967 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 803.8,
            "unit": "ns/op",
            "extra": "1497967 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1497967 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1497967 times\n4 procs"
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
          "id": "25626627e3c762ac3cc12a0637b5802126b913be",
          "message": "fix(translator): check gateway class and do not skip listeners when gateway class is not managed by KIC (#7666)\n\n* check gateway class and do not skip listeners when gateway class is not managed by KIC\n\n* update changelog\n\n* store gatewayclass into cache in gwc controller\n\n* fix comments and move gwc controller name check after getting gwc",
          "timestamp": "2025-09-18T08:18:27Z",
          "tree_id": "671f0416cdb35ce1f5285a527b60da52bae3b853",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/25626627e3c762ac3cc12a0637b5802126b913be"
        },
        "date": 1758185364624,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 610905,
            "unit": "ns/op\t  819449 B/op\t       5 allocs/op",
            "extra": "1714 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 610905,
            "unit": "ns/op",
            "extra": "1714 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819449,
            "unit": "B/op",
            "extra": "1714 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1714 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7364,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "160399 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7364,
            "unit": "ns/op",
            "extra": "160399 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "160399 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "160399 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 107.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11179220 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 107.2,
            "unit": "ns/op",
            "extra": "11179220 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11179220 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11179220 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21539,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54912 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21539,
            "unit": "ns/op",
            "extra": "54912 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54912 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54912 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 209453,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5295 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 209453,
            "unit": "ns/op",
            "extra": "5295 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5295 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5295 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2645362,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "397 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2645362,
            "unit": "ns/op",
            "extra": "397 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "397 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "397 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 29375789,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "46 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 29375789,
            "unit": "ns/op",
            "extra": "46 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "46 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "46 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9302,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "127971 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9302,
            "unit": "ns/op",
            "extra": "127971 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "127971 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "127971 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6917921,
            "unit": "ns/op\t 4527542 B/op\t   69230 allocs/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6917921,
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
            "value": 7717,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "155433 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7717,
            "unit": "ns/op",
            "extra": "155433 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "155433 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "155433 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 832269,
            "unit": "ns/op\t  392946 B/op\t    6173 allocs/op",
            "extra": "1406 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 832269,
            "unit": "ns/op",
            "extra": "1406 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392946,
            "unit": "B/op",
            "extra": "1406 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1406 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10786,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "107012 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10786,
            "unit": "ns/op",
            "extra": "107012 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "107012 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "107012 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7993088,
            "unit": "ns/op\t 4914221 B/op\t   75241 allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7993088,
            "unit": "ns/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914221,
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
            "value": 341.2,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3523699 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 341.2,
            "unit": "ns/op",
            "extra": "3523699 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3523699 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3523699 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.48,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.48,
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
            "value": 555.7,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2148822 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 555.7,
            "unit": "ns/op",
            "extra": "2148822 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2148822 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2148822 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 811.7,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1487764 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 811.7,
            "unit": "ns/op",
            "extra": "1487764 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1487764 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1487764 times\n4 procs"
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
            "email": "tao.yi@konghq.com",
            "name": "Yi Tao",
            "username": "randmonkey"
          },
          "distinct": true,
          "id": "7a3a28ac0cde7a54fbdb762002c20e09f795c4e2",
          "message": "set postgresql image repo to bitnamilegacy/postgresql (#7710)\n\n(cherry picked from commit 7342711be211b52fe7bef04c67b0a8247f748d4f)",
          "timestamp": "2025-09-18T17:33:38+08:00",
          "tree_id": "661a7b12423edbf3345daa714f4948967180027c",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/7a3a28ac0cde7a54fbdb762002c20e09f795c4e2"
        },
        "date": 1758188165367,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkDeckgenGenerateSHA",
            "value": 71393,
            "unit": "ns/op\t   11109 B/op\t      12 allocs/op",
            "extra": "17677 times\n4 procs"
          },
          {
            "name": "BenchmarkDeckgenGenerateSHA - ns/op",
            "value": 71393,
            "unit": "ns/op",
            "extra": "17677 times\n4 procs"
          },
          {
            "name": "BenchmarkDeckgenGenerateSHA - B/op",
            "value": 11109,
            "unit": "B/op",
            "extra": "17677 times\n4 procs"
          },
          {
            "name": "BenchmarkDeckgenGenerateSHA - allocs/op",
            "value": 12,
            "unit": "allocs/op",
            "extra": "17677 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 142.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8958936 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 142.2,
            "unit": "ns/op",
            "extra": "8958936 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8958936 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8958936 times\n4 procs"
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
          "id": "360fbc44a9149bed2d94ed7185e8e30e937a8974",
          "message": "fix changelog for 7666 (#7713)",
          "timestamp": "2025-09-18T09:52:54Z",
          "tree_id": "218028c2cbae599d8bac69d3a2110f1d8a45018f",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/360fbc44a9149bed2d94ed7185e8e30e937a8974"
        },
        "date": 1758191061515,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 640211,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "1658 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 640211,
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
            "value": 7524,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "154009 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7524,
            "unit": "ns/op",
            "extra": "154009 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "154009 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "154009 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 110.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10805136 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 110.9,
            "unit": "ns/op",
            "extra": "10805136 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "10805136 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "10805136 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21893,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54690 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21893,
            "unit": "ns/op",
            "extra": "54690 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54690 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54690 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 262827,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5499 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 262827,
            "unit": "ns/op",
            "extra": "5499 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5499 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5499 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2424513,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "501 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2424513,
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
            "value": 35350735,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 35350735,
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
            "value": 9988,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "120643 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9988,
            "unit": "ns/op",
            "extra": "120643 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "120643 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "120643 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7712042,
            "unit": "ns/op\t 4527576 B/op\t   69230 allocs/op",
            "extra": "156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7712042,
            "unit": "ns/op",
            "extra": "156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527576,
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
            "value": 8209,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "146998 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8209,
            "unit": "ns/op",
            "extra": "146998 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "146998 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "146998 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 882674,
            "unit": "ns/op\t  392930 B/op\t    6173 allocs/op",
            "extra": "1346 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 882674,
            "unit": "ns/op",
            "extra": "1346 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392930,
            "unit": "B/op",
            "extra": "1346 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1346 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11715,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "102049 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11715,
            "unit": "ns/op",
            "extra": "102049 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "102049 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "102049 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 9017411,
            "unit": "ns/op\t 4914237 B/op\t   75241 allocs/op",
            "extra": "135 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 9017411,
            "unit": "ns/op",
            "extra": "135 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914237,
            "unit": "B/op",
            "extra": "135 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "135 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 321.9,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3736196 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 321.9,
            "unit": "ns/op",
            "extra": "3736196 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3736196 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3736196 times\n4 procs"
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
            "value": 569.2,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2128200 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 569.2,
            "unit": "ns/op",
            "extra": "2128200 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2128200 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2128200 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 812.2,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1471686 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 812.2,
            "unit": "ns/op",
            "extra": "1471686 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1471686 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1471686 times\n4 procs"
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
          "id": "8b3ba5181a8b3eb69561507c8ca0b2cdb5ef6266",
          "message": "set postgresql image repo to bitnamilegacy/postgresql (#7710) (#7714)\n\n(cherry picked from commit 7342711be211b52fe7bef04c67b0a8247f748d4f)",
          "timestamp": "2025-09-19T10:54:26+02:00",
          "tree_id": "661a7b12423edbf3345daa714f4948967180027c",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/8b3ba5181a8b3eb69561507c8ca0b2cdb5ef6266"
        },
        "date": 1758272185849,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkDeckgenGenerateSHA",
            "value": 70215,
            "unit": "ns/op\t   11102 B/op\t      12 allocs/op",
            "extra": "16730 times\n4 procs"
          },
          {
            "name": "BenchmarkDeckgenGenerateSHA - ns/op",
            "value": 70215,
            "unit": "ns/op",
            "extra": "16730 times\n4 procs"
          },
          {
            "name": "BenchmarkDeckgenGenerateSHA - B/op",
            "value": 11102,
            "unit": "B/op",
            "extra": "16730 times\n4 procs"
          },
          {
            "name": "BenchmarkDeckgenGenerateSHA - allocs/op",
            "value": 12,
            "unit": "allocs/op",
            "extra": "16730 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 149.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7404535 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 149.7,
            "unit": "ns/op",
            "extra": "7404535 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7404535 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7404535 times\n4 procs"
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
          "id": "532acdb6dcb55bead6c99f49b8b359eeeaf085fb",
          "message": "chore(deps): update dependency go-delve/delve to v1.25.2 (#7696)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2025-09-19T08:54:59Z",
          "tree_id": "9f3ac1c8da067f4d4c64b93e9d2c2e66c2e328b9",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/532acdb6dcb55bead6c99f49b8b359eeeaf085fb"
        },
        "date": 1758273918745,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 769829,
            "unit": "ns/op\t  819454 B/op\t       5 allocs/op",
            "extra": "1360 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 769829,
            "unit": "ns/op",
            "extra": "1360 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819454,
            "unit": "B/op",
            "extra": "1360 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1360 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7635,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "156288 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7635,
            "unit": "ns/op",
            "extra": "156288 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "156288 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "156288 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 111.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10762221 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 111.4,
            "unit": "ns/op",
            "extra": "10762221 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "10762221 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "10762221 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21660,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53818 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21660,
            "unit": "ns/op",
            "extra": "53818 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53818 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53818 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 219877,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5022 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 219877,
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
            "value": 2837076,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "390 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2837076,
            "unit": "ns/op",
            "extra": "390 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "390 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "390 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 31373935,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 31373935,
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
            "value": 9541,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "125947 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9541,
            "unit": "ns/op",
            "extra": "125947 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "125947 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "125947 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7004324,
            "unit": "ns/op\t 4527542 B/op\t   69230 allocs/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7004324,
            "unit": "ns/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527542,
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
            "value": 7918,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "154015 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7918,
            "unit": "ns/op",
            "extra": "154015 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "154015 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "154015 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 842255,
            "unit": "ns/op\t  392937 B/op\t    6173 allocs/op",
            "extra": "1381 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 842255,
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
            "value": 10962,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "111498 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10962,
            "unit": "ns/op",
            "extra": "111498 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "111498 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "111498 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8005289,
            "unit": "ns/op\t 4914201 B/op\t   75241 allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8005289,
            "unit": "ns/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914201,
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
            "value": 345.4,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3493134 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 345.4,
            "unit": "ns/op",
            "extra": "3493134 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3493134 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3493134 times\n4 procs"
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
            "value": 562.8,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2149840 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 562.8,
            "unit": "ns/op",
            "extra": "2149840 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2149840 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2149840 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 800.2,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1494897 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 800.2,
            "unit": "ns/op",
            "extra": "1494897 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1494897 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1494897 times\n4 procs"
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
          "id": "7ea22fafedb7ebc2bf6396d72aea665a2e808757",
          "message": "chore(deps): bump jdx/mise-action from 2.3.1 to 3.2.0 (#7694)\n\nBumps [jdx/mise-action](https://github.com/jdx/mise-action) from 2.3.1 to 3.2.0.\n- [Release notes](https://github.com/jdx/mise-action/releases)\n- [Changelog](https://github.com/jdx/mise-action/blob/main/CHANGELOG.md)\n- [Commits](https://github.com/jdx/mise-action/compare/5cb1df66ed5e1fb3c670ea0b62fd17a76979826a...5ac50f778e26fac95da98d50503682459e86d566)\n\n---\nupdated-dependencies:\n- dependency-name: jdx/mise-action\n  dependency-version: 3.2.0\n  dependency-type: direct:production\n  update-type: version-update:semver-major\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2025-09-19T08:55:06Z",
          "tree_id": "d3bf3f920352cdf5e536fd619144613200337c81",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/7ea22fafedb7ebc2bf6396d72aea665a2e808757"
        },
        "date": 1758273965229,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 697530,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "1623 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 697530,
            "unit": "ns/op",
            "extra": "1623 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819453,
            "unit": "B/op",
            "extra": "1623 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1623 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7329,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "163939 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7329,
            "unit": "ns/op",
            "extra": "163939 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "163939 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "163939 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 107.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11145846 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 107.5,
            "unit": "ns/op",
            "extra": "11145846 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11145846 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11145846 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21756,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54944 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21756,
            "unit": "ns/op",
            "extra": "54944 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54944 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54944 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 206272,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5156 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 206272,
            "unit": "ns/op",
            "extra": "5156 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5156 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5156 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2392111,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "490 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2392111,
            "unit": "ns/op",
            "extra": "490 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "490 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "490 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 40201918,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "31 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 40201918,
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
            "value": 9727,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "123966 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9727,
            "unit": "ns/op",
            "extra": "123966 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "123966 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "123966 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 8821703,
            "unit": "ns/op\t 4527535 B/op\t   69230 allocs/op",
            "extra": "135 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 8821703,
            "unit": "ns/op",
            "extra": "135 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527535,
            "unit": "B/op",
            "extra": "135 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "135 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7897,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "154094 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7897,
            "unit": "ns/op",
            "extra": "154094 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "154094 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "154094 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 848638,
            "unit": "ns/op\t  392943 B/op\t    6173 allocs/op",
            "extra": "1388 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 848638,
            "unit": "ns/op",
            "extra": "1388 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392943,
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
            "value": 11007,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "110422 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11007,
            "unit": "ns/op",
            "extra": "110422 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "110422 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "110422 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8112352,
            "unit": "ns/op\t 4914185 B/op\t   75241 allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8112352,
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
            "value": 346.6,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3402193 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 346.6,
            "unit": "ns/op",
            "extra": "3402193 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3402193 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3402193 times\n4 procs"
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
            "value": 562.7,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2130619 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 562.7,
            "unit": "ns/op",
            "extra": "2130619 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2130619 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2130619 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 809,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1482015 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 809,
            "unit": "ns/op",
            "extra": "1482015 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1482015 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1482015 times\n4 procs"
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
          "id": "1a574d1d636cb9daf38aefa361927616c8c40ca4",
          "message": "chore(deps): bump sigs.k8s.io/kustomize/kyaml from 0.20.0 to 0.20.1 (#7673)\n\nBumps [sigs.k8s.io/kustomize/kyaml](https://github.com/kubernetes-sigs/kustomize) from 0.20.0 to 0.20.1.\n- [Release notes](https://github.com/kubernetes-sigs/kustomize/releases)\n- [Commits](https://github.com/kubernetes-sigs/kustomize/compare/api/v0.20.0...api/v0.20.1)\n\n---\nupdated-dependencies:\n- dependency-name: sigs.k8s.io/kustomize/kyaml\n  dependency-version: 0.20.1\n  dependency-type: direct:production\n  update-type: version-update:semver-patch\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2025-09-19T08:55:51Z",
          "tree_id": "7db731edcc3b4e44653f3c92c0c848d6af49ec34",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/1a574d1d636cb9daf38aefa361927616c8c40ca4"
        },
        "date": 1758275675528,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 697358,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "1473 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 697358,
            "unit": "ns/op",
            "extra": "1473 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819453,
            "unit": "B/op",
            "extra": "1473 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1473 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7040,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "157044 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7040,
            "unit": "ns/op",
            "extra": "157044 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "157044 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "157044 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 107.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11184864 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 107.1,
            "unit": "ns/op",
            "extra": "11184864 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11184864 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11184864 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21923,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53641 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21923,
            "unit": "ns/op",
            "extra": "53641 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53641 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53641 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 219586,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5217 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 219586,
            "unit": "ns/op",
            "extra": "5217 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5217 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5217 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2336648,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "504 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2336648,
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
            "value": 43628330,
            "unit": "ns/op\t24010754 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 43628330,
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
            "value": 9700,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "123276 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9700,
            "unit": "ns/op",
            "extra": "123276 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "123276 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "123276 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7369342,
            "unit": "ns/op\t 4527547 B/op\t   69230 allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7369342,
            "unit": "ns/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527547,
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
            "value": 8172,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "150608 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8172,
            "unit": "ns/op",
            "extra": "150608 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "150608 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "150608 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 852553,
            "unit": "ns/op\t  392932 B/op\t    6173 allocs/op",
            "extra": "1375 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 852553,
            "unit": "ns/op",
            "extra": "1375 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392932,
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
            "value": 11419,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "106821 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11419,
            "unit": "ns/op",
            "extra": "106821 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "106821 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "106821 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8549982,
            "unit": "ns/op\t 4914242 B/op\t   75241 allocs/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8549982,
            "unit": "ns/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914242,
            "unit": "B/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 319.5,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3736623 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 319.5,
            "unit": "ns/op",
            "extra": "3736623 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3736623 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3736623 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "95050851 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.6,
            "unit": "ns/op",
            "extra": "95050851 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "95050851 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "95050851 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 564.2,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2126581 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 564.2,
            "unit": "ns/op",
            "extra": "2126581 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2126581 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2126581 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 820.3,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1465852 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 820.3,
            "unit": "ns/op",
            "extra": "1465852 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1465852 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1465852 times\n4 procs"
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
          "id": "cc01d31a1918731c63dc5280a855971a2ef76fbc",
          "message": "chore(deps): bump docker/metadata-action from 5.7.0 to 5.8.0 (#7664)\n\nBumps [docker/metadata-action](https://github.com/docker/metadata-action) from 5.7.0 to 5.8.0.\n- [Release notes](https://github.com/docker/metadata-action/releases)\n- [Commits](https://github.com/docker/metadata-action/compare/902fa8ec7d6ecbf8d84d538b9b233a880e428804...c1e51972afc2121e065aed6d45c65596fe445f3f)\n\n---\nupdated-dependencies:\n- dependency-name: docker/metadata-action\n  dependency-version: 5.8.0\n  dependency-type: direct:production\n  update-type: version-update:semver-minor\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2025-09-19T09:27:59Z",
          "tree_id": "2d3c385f371f44e314d38dbd2a8b9708aac0e3ad",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/cc01d31a1918731c63dc5280a855971a2ef76fbc"
        },
        "date": 1758275960980,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 711478,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "1540 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 711478,
            "unit": "ns/op",
            "extra": "1540 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819453,
            "unit": "B/op",
            "extra": "1540 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1540 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 8671,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "148371 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 8671,
            "unit": "ns/op",
            "extra": "148371 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "148371 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "148371 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 111.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10745558 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 111.2,
            "unit": "ns/op",
            "extra": "10745558 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "10745558 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "10745558 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21552,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "55358 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21552,
            "unit": "ns/op",
            "extra": "55358 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "55358 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "55358 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 205025,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5167 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 205025,
            "unit": "ns/op",
            "extra": "5167 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5167 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5167 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2399752,
            "unit": "ns/op\t 2408450 B/op\t       2 allocs/op",
            "extra": "500 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2399752,
            "unit": "ns/op",
            "extra": "500 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408450,
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
            "value": 32367569,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 32367569,
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
            "value": 9433,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "124898 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9433,
            "unit": "ns/op",
            "extra": "124898 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "124898 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "124898 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6933757,
            "unit": "ns/op\t 4527545 B/op\t   69230 allocs/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6933757,
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
            "value": 7866,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "152172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7866,
            "unit": "ns/op",
            "extra": "152172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "152172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "152172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 837434,
            "unit": "ns/op\t  392945 B/op\t    6173 allocs/op",
            "extra": "1419 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 837434,
            "unit": "ns/op",
            "extra": "1419 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392945,
            "unit": "B/op",
            "extra": "1419 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1419 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10889,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "110892 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10889,
            "unit": "ns/op",
            "extra": "110892 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "110892 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "110892 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7990663,
            "unit": "ns/op\t 4914263 B/op\t   75241 allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7990663,
            "unit": "ns/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914263,
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
            "value": 318.3,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3787234 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 318.3,
            "unit": "ns/op",
            "extra": "3787234 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3787234 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3787234 times\n4 procs"
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
            "value": 662.2,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "1997860 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 662.2,
            "unit": "ns/op",
            "extra": "1997860 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "1997860 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1997860 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 856.9,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1391112 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 856.9,
            "unit": "ns/op",
            "extra": "1391112 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1391112 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1391112 times\n4 procs"
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
          "id": "6108a8a3ea6c2cb7c32e35b8fc8118b57b135a85",
          "message": "chore(deps): bump docker/login-action from 3.4.0 to 3.5.0 (#7674)\n\nBumps [docker/login-action](https://github.com/docker/login-action) from 3.4.0 to 3.5.0.\n- [Release notes](https://github.com/docker/login-action/releases)\n- [Commits](https://github.com/docker/login-action/compare/74a5d142397b4f367a81961eba4e8cd7edddf772...184bdaa0721073962dff0199f1fb9940f07167d1)\n\n---\nupdated-dependencies:\n- dependency-name: docker/login-action\n  dependency-version: 3.5.0\n  dependency-type: direct:production\n  update-type: version-update:semver-minor\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2025-09-19T09:28:57Z",
          "tree_id": "c748862c5e2ac0bbef8ba600527fd91a0ba52d08",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/6108a8a3ea6c2cb7c32e35b8fc8118b57b135a85"
        },
        "date": 1758276019428,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 661025,
            "unit": "ns/op\t  819450 B/op\t       5 allocs/op",
            "extra": "1522 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 661025,
            "unit": "ns/op",
            "extra": "1522 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819450,
            "unit": "B/op",
            "extra": "1522 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1522 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7219,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "165500 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7219,
            "unit": "ns/op",
            "extra": "165500 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "165500 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "165500 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 108.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11184882 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 108.2,
            "unit": "ns/op",
            "extra": "11184882 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11184882 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11184882 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21438,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "55144 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21438,
            "unit": "ns/op",
            "extra": "55144 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "55144 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "55144 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 205879,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5084 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 205879,
            "unit": "ns/op",
            "extra": "5084 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5084 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5084 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2282378,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "513 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2282378,
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
            "value": 40510278,
            "unit": "ns/op\t24010754 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 40510278,
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
            "value": 9926,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "103996 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9926,
            "unit": "ns/op",
            "extra": "103996 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "103996 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "103996 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6982327,
            "unit": "ns/op\t 4527544 B/op\t   69230 allocs/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6982327,
            "unit": "ns/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527544,
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
            "value": 7954,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "150637 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7954,
            "unit": "ns/op",
            "extra": "150637 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "150637 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "150637 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 842525,
            "unit": "ns/op\t  392940 B/op\t    6173 allocs/op",
            "extra": "1388 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 842525,
            "unit": "ns/op",
            "extra": "1388 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392940,
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
            "value": 11011,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "110212 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11011,
            "unit": "ns/op",
            "extra": "110212 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "110212 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "110212 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7957340,
            "unit": "ns/op\t 4914172 B/op\t   75241 allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7957340,
            "unit": "ns/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914172,
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
            "value": 320.2,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3660111 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 320.2,
            "unit": "ns/op",
            "extra": "3660111 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3660111 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3660111 times\n4 procs"
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
            "value": 597.3,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2133547 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 597.3,
            "unit": "ns/op",
            "extra": "2133547 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2133547 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2133547 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 962.5,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1206325 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 962.5,
            "unit": "ns/op",
            "extra": "1206325 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1206325 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1206325 times\n4 procs"
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
          "id": "d6bfebc9821595ab5faad77012a3db4c802b1c93",
          "message": "chore(deps): bump distroless/static from `cdf4daa` to `a9f88e0` (#7695)\n\nBumps distroless/static from `cdf4daa` to `a9f88e0`.\n\n---\nupdated-dependencies:\n- dependency-name: distroless/static\n  dependency-version: nonroot\n  dependency-type: direct:production\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2025-09-19T09:50:47Z",
          "tree_id": "96a0a0fda9af313f69e63b2be2248cf7e8b6c03c",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/d6bfebc9821595ab5faad77012a3db4c802b1c93"
        },
        "date": 1758277814542,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 703453,
            "unit": "ns/op\t  819454 B/op\t       5 allocs/op",
            "extra": "1474 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 703453,
            "unit": "ns/op",
            "extra": "1474 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819454,
            "unit": "B/op",
            "extra": "1474 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1474 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7479,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "140455 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7479,
            "unit": "ns/op",
            "extra": "140455 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "140455 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "140455 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 107.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11111125 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 107.5,
            "unit": "ns/op",
            "extra": "11111125 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11111125 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11111125 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21987,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "55435 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21987,
            "unit": "ns/op",
            "extra": "55435 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "55435 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "55435 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 209519,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5080 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 209519,
            "unit": "ns/op",
            "extra": "5080 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5080 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5080 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2394698,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "500 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2394698,
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
            "value": 38504663,
            "unit": "ns/op\t24010754 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 38504663,
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
            "value": 9653,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "124845 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9653,
            "unit": "ns/op",
            "extra": "124845 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "124845 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "124845 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7028231,
            "unit": "ns/op\t 4527542 B/op\t   69230 allocs/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7028231,
            "unit": "ns/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527542,
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
            "value": 7976,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "151938 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7976,
            "unit": "ns/op",
            "extra": "151938 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "151938 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "151938 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 846305,
            "unit": "ns/op\t  392939 B/op\t    6173 allocs/op",
            "extra": "1407 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 846305,
            "unit": "ns/op",
            "extra": "1407 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392939,
            "unit": "B/op",
            "extra": "1407 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1407 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11093,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "107787 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11093,
            "unit": "ns/op",
            "extra": "107787 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "107787 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "107787 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8337608,
            "unit": "ns/op\t 4914184 B/op\t   75241 allocs/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8337608,
            "unit": "ns/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914184,
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
            "value": 353.2,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3384022 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 353.2,
            "unit": "ns/op",
            "extra": "3384022 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3384022 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3384022 times\n4 procs"
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
            "value": 571,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2137622 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 571,
            "unit": "ns/op",
            "extra": "2137622 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2137622 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2137622 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 819,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1466587 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 819,
            "unit": "ns/op",
            "extra": "1466587 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1466587 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1466587 times\n4 procs"
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
          "id": "cc67509b1befe32e4f1d89513935d43104a013a3",
          "message": "[Cherry-pick] fix: Do not add request-termination on TLSRoute / TCPRoute (#7700) (#7720)\n\n* fix: Do not add request-termination on TLSRoute / TCPRoute\n\n* fix: add import of strings and check nil\n\n* fix unit tests and golden test cases\n\n* add changelog\n\n---------\n\nCo-authored-by: Victor Durand <contact@victordurand.fr>",
          "timestamp": "2025-09-22T03:01:48Z",
          "tree_id": "53ad723d54c93df599dcb919c126fd3a524e08b6",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/cc67509b1befe32e4f1d89513935d43104a013a3"
        },
        "date": 1758511999313,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 646283,
            "unit": "ns/op\t  819450 B/op\t       5 allocs/op",
            "extra": "1726 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 646283,
            "unit": "ns/op",
            "extra": "1726 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819450,
            "unit": "B/op",
            "extra": "1726 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1726 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 8392,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "163468 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 8392,
            "unit": "ns/op",
            "extra": "163468 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "163468 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "163468 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 107.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11121952 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 107.7,
            "unit": "ns/op",
            "extra": "11121952 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11121952 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11121952 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22016,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "51738 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22016,
            "unit": "ns/op",
            "extra": "51738 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "51738 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "51738 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 213383,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5071 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 213383,
            "unit": "ns/op",
            "extra": "5071 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5071 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5071 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2391134,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "495 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2391134,
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
            "value": 34785372,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 34785372,
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
            "value": 9860,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "112544 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9860,
            "unit": "ns/op",
            "extra": "112544 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "112544 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "112544 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7312229,
            "unit": "ns/op\t 4527537 B/op\t   69230 allocs/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7312229,
            "unit": "ns/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527537,
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
            "value": 8172,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "146590 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8172,
            "unit": "ns/op",
            "extra": "146590 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "146590 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "146590 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 863137,
            "unit": "ns/op\t  392932 B/op\t    6173 allocs/op",
            "extra": "1332 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 863137,
            "unit": "ns/op",
            "extra": "1332 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392932,
            "unit": "B/op",
            "extra": "1332 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1332 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11363,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "106729 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11363,
            "unit": "ns/op",
            "extra": "106729 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "106729 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "106729 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8327046,
            "unit": "ns/op\t 4914178 B/op\t   75241 allocs/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8327046,
            "unit": "ns/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914178,
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
            "value": 322.3,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3739077 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 322.3,
            "unit": "ns/op",
            "extra": "3739077 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3739077 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3739077 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.41,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.41,
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
            "value": 650.3,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2089446 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 650.3,
            "unit": "ns/op",
            "extra": "2089446 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2089446 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2089446 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 913.4,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1299085 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 913.4,
            "unit": "ns/op",
            "extra": "1299085 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1299085 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1299085 times\n4 procs"
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
          "id": "a90e89edbf13d0ab6a7e4a023787f5ee79c516cd",
          "message": "chore(deps): bump golang.org/x/mod from 0.25.0 to 0.28.0 (#7703)\n\nBumps [golang.org/x/mod](https://github.com/golang/mod) from 0.25.0 to 0.28.0.\n- [Commits](https://github.com/golang/mod/compare/v0.25.0...v0.28.0)\n\n---\nupdated-dependencies:\n- dependency-name: golang.org/x/mod\n  dependency-version: 0.28.0\n  dependency-type: direct:production\n  update-type: version-update:semver-minor\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2025-09-22T09:35:32Z",
          "tree_id": "f62694c7c9bd68629bc7881af0ace23a96bf401f",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/a90e89edbf13d0ab6a7e4a023787f5ee79c516cd"
        },
        "date": 1758535635364,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 581287,
            "unit": "ns/op\t  819450 B/op\t       5 allocs/op",
            "extra": "1824 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 581287,
            "unit": "ns/op",
            "extra": "1824 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819450,
            "unit": "B/op",
            "extra": "1824 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1824 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7279,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "163628 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7279,
            "unit": "ns/op",
            "extra": "163628 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "163628 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "163628 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 107,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11214135 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 107,
            "unit": "ns/op",
            "extra": "11214135 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11214135 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11214135 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21772,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54522 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21772,
            "unit": "ns/op",
            "extra": "54522 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54522 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54522 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 204884,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4974 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 204884,
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
            "value": 2855702,
            "unit": "ns/op\t 2408450 B/op\t       2 allocs/op",
            "extra": "414 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2855702,
            "unit": "ns/op",
            "extra": "414 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408450,
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
            "value": 31456460,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 31456460,
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
            "value": 9370,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "125040 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9370,
            "unit": "ns/op",
            "extra": "125040 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "125040 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "125040 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6911049,
            "unit": "ns/op\t 4527540 B/op\t   69230 allocs/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6911049,
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
            "value": 7839,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "154036 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7839,
            "unit": "ns/op",
            "extra": "154036 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "154036 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "154036 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 835281,
            "unit": "ns/op\t  392946 B/op\t    6173 allocs/op",
            "extra": "1407 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 835281,
            "unit": "ns/op",
            "extra": "1407 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392946,
            "unit": "B/op",
            "extra": "1407 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1407 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10793,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "112261 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10793,
            "unit": "ns/op",
            "extra": "112261 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "112261 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "112261 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7864951,
            "unit": "ns/op\t 4914224 B/op\t   75241 allocs/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7864951,
            "unit": "ns/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914224,
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
            "value": 321.7,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3748248 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 321.7,
            "unit": "ns/op",
            "extra": "3748248 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3748248 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3748248 times\n4 procs"
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
            "value": 563.5,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2167262 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 563.5,
            "unit": "ns/op",
            "extra": "2167262 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2167262 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2167262 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 811.4,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1492867 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 811.4,
            "unit": "ns/op",
            "extra": "1492867 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1492867 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1492867 times\n4 procs"
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
          "id": "5e139f51a2508fa01327f655b08353831d2e07af",
          "message": "bump GKE to 1.33.5 (#7746)",
          "timestamp": "2025-10-13T06:06:17Z",
          "tree_id": "c10dd0059ce8dd7d2c3433152dc8adc5e19fd6ea",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/5e139f51a2508fa01327f655b08353831d2e07af"
        },
        "date": 1760337416704,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 685572,
            "unit": "ns/op\t  819452 B/op\t       5 allocs/op",
            "extra": "1569 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 685572,
            "unit": "ns/op",
            "extra": "1569 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819452,
            "unit": "B/op",
            "extra": "1569 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1569 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7370,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "163468 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7370,
            "unit": "ns/op",
            "extra": "163468 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "163468 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "163468 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 107.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11195516 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 107.1,
            "unit": "ns/op",
            "extra": "11195516 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11195516 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11195516 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21553,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54798 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21553,
            "unit": "ns/op",
            "extra": "54798 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54798 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54798 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 208535,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4969 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 208535,
            "unit": "ns/op",
            "extra": "4969 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4969 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4969 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2936027,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "363 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2936027,
            "unit": "ns/op",
            "extra": "363 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "363 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "363 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 31776920,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "68 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 31776920,
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
            "value": 9506,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "125500 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9506,
            "unit": "ns/op",
            "extra": "125500 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "125500 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "125500 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7013523,
            "unit": "ns/op\t 4527547 B/op\t   69230 allocs/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7013523,
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
            "value": 7970,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "152125 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7970,
            "unit": "ns/op",
            "extra": "152125 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "152125 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "152125 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 843525,
            "unit": "ns/op\t  392926 B/op\t    6173 allocs/op",
            "extra": "1402 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 843525,
            "unit": "ns/op",
            "extra": "1402 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392926,
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
            "value": 11003,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "109232 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11003,
            "unit": "ns/op",
            "extra": "109232 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "109232 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "109232 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8095159,
            "unit": "ns/op\t 4914166 B/op\t   75241 allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8095159,
            "unit": "ns/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914166,
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
            "value": 319.8,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3763354 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 319.8,
            "unit": "ns/op",
            "extra": "3763354 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3763354 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3763354 times\n4 procs"
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
            "value": 563.2,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2120475 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 563.2,
            "unit": "ns/op",
            "extra": "2120475 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2120475 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2120475 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 805.2,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1490474 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 805.2,
            "unit": "ns/op",
            "extra": "1490474 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1490474 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1490474 times\n4 procs"
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
          "id": "64e50933f0688324f3954145d4a2ef4705ed1c43",
          "message": "fix(translator): reject CA Secrets with multiple PEM certs (#7763)\n\n* docs: update API references and bump crd-ref-docs to 0.2.0\n\n* fix(translator): reject CA Secrets with multiple PEM certs\n\n* update CHANGELOG\n\nSigned-off-by: Jintao Zhang <zhangjintao9020@gmail.com>\n\n---------\n\nSigned-off-by: Jintao Zhang <zhangjintao9020@gmail.com>",
          "timestamp": "2025-10-30T08:40:22Z",
          "tree_id": "64aac333e31017a4b5b3c3134711d1e5bc728a7c",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/64e50933f0688324f3954145d4a2ef4705ed1c43"
        },
        "date": 1761815477910,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 663995,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "1549 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 663995,
            "unit": "ns/op",
            "extra": "1549 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819453,
            "unit": "B/op",
            "extra": "1549 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1549 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 9127,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "132025 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 9127,
            "unit": "ns/op",
            "extra": "132025 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "132025 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "132025 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 107.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11134431 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 107.7,
            "unit": "ns/op",
            "extra": "11134431 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11134431 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11134431 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21865,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54415 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21865,
            "unit": "ns/op",
            "extra": "54415 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54415 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54415 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 225572,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4833 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 225572,
            "unit": "ns/op",
            "extra": "4833 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4833 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4833 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2398143,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "499 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2398143,
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
            "value": 50818665,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "33 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 50818665,
            "unit": "ns/op",
            "extra": "33 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
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
            "value": 12450,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "94711 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 12450,
            "unit": "ns/op",
            "extra": "94711 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "94711 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "94711 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7245699,
            "unit": "ns/op\t 4527544 B/op\t   69230 allocs/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7245699,
            "unit": "ns/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527544,
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
            "value": 7981,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "150796 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7981,
            "unit": "ns/op",
            "extra": "150796 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "150796 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "150796 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 846790,
            "unit": "ns/op\t  392946 B/op\t    6173 allocs/op",
            "extra": "1365 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 846790,
            "unit": "ns/op",
            "extra": "1365 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392946,
            "unit": "B/op",
            "extra": "1365 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1365 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11211,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "106965 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11211,
            "unit": "ns/op",
            "extra": "106965 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "106965 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "106965 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8043836,
            "unit": "ns/op\t 4914254 B/op\t   75241 allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8043836,
            "unit": "ns/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914254,
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
            "value": 349.5,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3386828 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 349.5,
            "unit": "ns/op",
            "extra": "3386828 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3386828 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3386828 times\n4 procs"
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
            "value": 575.8,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2082396 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 575.8,
            "unit": "ns/op",
            "extra": "2082396 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2082396 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2082396 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 950.8,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1321680 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 950.8,
            "unit": "ns/op",
            "extra": "1321680 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1321680 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1321680 times\n4 procs"
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
          "distinct": false,
          "id": "1a4afb2725ce12e1bbd2f7a2f5d3e55fe288fa69",
          "message": "ci: fix envtest suite CI failures (#7771)",
          "timestamp": "2025-11-03T15:01:16Z",
          "tree_id": "6ecc44257c030053b316008ff40992c822ff625e",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/1a4afb2725ce12e1bbd2f7a2f5d3e55fe288fa69"
        },
        "date": 1762183914559,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 636367,
            "unit": "ns/op\t  819449 B/op\t       5 allocs/op",
            "extra": "1642 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 636367,
            "unit": "ns/op",
            "extra": "1642 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819449,
            "unit": "B/op",
            "extra": "1642 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1642 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7320,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "162775 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7320,
            "unit": "ns/op",
            "extra": "162775 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "162775 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "162775 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 107.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11175531 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 107.3,
            "unit": "ns/op",
            "extra": "11175531 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11175531 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11175531 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21704,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53559 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21704,
            "unit": "ns/op",
            "extra": "53559 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53559 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53559 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 206107,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "6000 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 206107,
            "unit": "ns/op",
            "extra": "6000 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "6000 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "6000 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2755235,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "500 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2755235,
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
            "value": 28986573,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "46 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 28986573,
            "unit": "ns/op",
            "extra": "46 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "46 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "46 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9340,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "125654 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9340,
            "unit": "ns/op",
            "extra": "125654 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "125654 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "125654 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7047232,
            "unit": "ns/op\t 4527541 B/op\t   69230 allocs/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7047232,
            "unit": "ns/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527541,
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
            "value": 7816,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "155492 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7816,
            "unit": "ns/op",
            "extra": "155492 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "155492 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "155492 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 835141,
            "unit": "ns/op\t  392928 B/op\t    6173 allocs/op",
            "extra": "1411 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 835141,
            "unit": "ns/op",
            "extra": "1411 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392928,
            "unit": "B/op",
            "extra": "1411 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1411 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10777,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "111087 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10777,
            "unit": "ns/op",
            "extra": "111087 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "111087 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "111087 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7941313,
            "unit": "ns/op\t 4914177 B/op\t   75241 allocs/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7941313,
            "unit": "ns/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914177,
            "unit": "B/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 347.8,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3430635 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 347.8,
            "unit": "ns/op",
            "extra": "3430635 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3430635 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3430635 times\n4 procs"
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
            "value": 561.1,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2131575 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 561.1,
            "unit": "ns/op",
            "extra": "2131575 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2131575 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2131575 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 815,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1483760 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 815,
            "unit": "ns/op",
            "extra": "1483760 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1483760 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1483760 times\n4 procs"
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
          "id": "5b8404b648e89bed96decf65754e5a83a5830e54",
          "message": "chore(deps): bump golang from 1.25.1 to 1.25.3 (#7749)\n\nBumps golang from 1.25.1 to 1.25.3.\n\n---\nupdated-dependencies:\n- dependency-name: golang\n  dependency-version: 1.25.3\n  dependency-type: direct:production\n  update-type: version-update:semver-patch\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2025-11-03T16:32:39+01:00",
          "tree_id": "5ad52e1e80239ab84fd7df2f4a26ca75683e0939",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/5b8404b648e89bed96decf65754e5a83a5830e54"
        },
        "date": 1762184133000,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 747392,
            "unit": "ns/op\t  819455 B/op\t       5 allocs/op",
            "extra": "1447 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 747392,
            "unit": "ns/op",
            "extra": "1447 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819455,
            "unit": "B/op",
            "extra": "1447 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1447 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7465,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "149394 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7465,
            "unit": "ns/op",
            "extra": "149394 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "149394 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "149394 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 107.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11178122 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 107.4,
            "unit": "ns/op",
            "extra": "11178122 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11178122 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11178122 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21748,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54255 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21748,
            "unit": "ns/op",
            "extra": "54255 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54255 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54255 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 254648,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "5090 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 254648,
            "unit": "ns/op",
            "extra": "5090 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "5090 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5090 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2492979,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "442 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2492979,
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
            "value": 37780693,
            "unit": "ns/op\t24010754 B/op\t       2 allocs/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 37780693,
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
            "value": 9521,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "126693 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9521,
            "unit": "ns/op",
            "extra": "126693 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "126693 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "126693 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7106043,
            "unit": "ns/op\t 4527542 B/op\t   69230 allocs/op",
            "extra": "166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7106043,
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
            "value": 7776,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "154642 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7776,
            "unit": "ns/op",
            "extra": "154642 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "154642 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "154642 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 845239,
            "unit": "ns/op\t  392935 B/op\t    6173 allocs/op",
            "extra": "1405 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 845239,
            "unit": "ns/op",
            "extra": "1405 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392935,
            "unit": "B/op",
            "extra": "1405 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1405 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10833,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "110683 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10833,
            "unit": "ns/op",
            "extra": "110683 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "110683 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "110683 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8021990,
            "unit": "ns/op\t 4914231 B/op\t   75241 allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8021990,
            "unit": "ns/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914231,
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
            "value": 318.3,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3765768 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 318.3,
            "unit": "ns/op",
            "extra": "3765768 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3765768 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3765768 times\n4 procs"
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
            "value": 566.9,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2114528 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 566.9,
            "unit": "ns/op",
            "extra": "2114528 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2114528 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2114528 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 819.9,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1463703 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 819.9,
            "unit": "ns/op",
            "extra": "1463703 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1463703 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1463703 times\n4 procs"
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
          "id": "2bb2029577de2e30d3b0dea01e36baa7c82c3157",
          "message": "fix(sendconfig): Do not cleanup `null`s in plugin configuration   (#7751)\n\n* remove cleanup of nulls in plugin config\n\n* update tests\n\n* update changelog\n\n* dump plugin config in tests\n\n* find plugin in multiple plugins\n\n* use different paths for plugin tests and fix workspace",
          "timestamp": "2025-11-04T06:07:24Z",
          "tree_id": "a262f0bc1cf7e850a34d5b387f29ed3b959ddf0b",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/2bb2029577de2e30d3b0dea01e36baa7c82c3157"
        },
        "date": 1762238295356,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 705002,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "1455 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 705002,
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
            "value": 7911,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "149803 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7911,
            "unit": "ns/op",
            "extra": "149803 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "149803 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "149803 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.79,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16447561 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.79,
            "unit": "ns/op",
            "extra": "16447561 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16447561 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16447561 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22395,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53107 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22395,
            "unit": "ns/op",
            "extra": "53107 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53107 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53107 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 232404,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4797 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 232404,
            "unit": "ns/op",
            "extra": "4797 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4797 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4797 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2372465,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "482 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2372465,
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
            "value": 36610883,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "68 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 36610883,
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
            "value": 9821,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "121671 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9821,
            "unit": "ns/op",
            "extra": "121671 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "121671 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "121671 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7900429,
            "unit": "ns/op\t 4527545 B/op\t   69230 allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7900429,
            "unit": "ns/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527545,
            "unit": "B/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69230,
            "unit": "allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8121,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "150070 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8121,
            "unit": "ns/op",
            "extra": "150070 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "150070 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "150070 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 870344,
            "unit": "ns/op\t  392946 B/op\t    6173 allocs/op",
            "extra": "1365 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 870344,
            "unit": "ns/op",
            "extra": "1365 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392946,
            "unit": "B/op",
            "extra": "1365 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1365 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11696,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "103420 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11696,
            "unit": "ns/op",
            "extra": "103420 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "103420 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "103420 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 9881194,
            "unit": "ns/op\t 4914213 B/op\t   75241 allocs/op",
            "extra": "120 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 9881194,
            "unit": "ns/op",
            "extra": "120 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914213,
            "unit": "B/op",
            "extra": "120 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "120 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 323.8,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3691597 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 323.8,
            "unit": "ns/op",
            "extra": "3691597 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3691597 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3691597 times\n4 procs"
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
            "value": 572.3,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2092558 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 572.3,
            "unit": "ns/op",
            "extra": "2092558 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2092558 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2092558 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 825.9,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1458398 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 825.9,
            "unit": "ns/op",
            "extra": "1458398 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1458398 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1458398 times\n4 procs"
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
            "email": "tao.yi@konghq.com",
            "name": "Yi Tao",
            "username": "randmonkey"
          },
          "distinct": true,
          "id": "7ad3c28caacd90266d742f9de50dc0ba348bfe22",
          "message": "fix(translator): reject CA Secrets with multiple PEM certs (#7763)\n\n* docs: update API references and bump crd-ref-docs to 0.2.0\n\n* fix(translator): reject CA Secrets with multiple PEM certs\n\n* update CHANGELOG\n\nSigned-off-by: Jintao Zhang <zhangjintao9020@gmail.com>\n\n---------\n\nSigned-off-by: Jintao Zhang <zhangjintao9020@gmail.com>\n(cherry picked from commit 64e50933f0688324f3954145d4a2ef4705ed1c43)",
          "timestamp": "2025-11-05T16:23:06+08:00",
          "tree_id": "708e9069b1491ba3c86fd5aecaefe7ad41b96bf0",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/7ad3c28caacd90266d742f9de50dc0ba348bfe22"
        },
        "date": 1762331150526,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkDeckgenGenerateSHA",
            "value": 71143,
            "unit": "ns/op\t   11111 B/op\t      12 allocs/op",
            "extra": "16454 times\n4 procs"
          },
          {
            "name": "BenchmarkDeckgenGenerateSHA - ns/op",
            "value": 71143,
            "unit": "ns/op",
            "extra": "16454 times\n4 procs"
          },
          {
            "name": "BenchmarkDeckgenGenerateSHA - B/op",
            "value": 11111,
            "unit": "B/op",
            "extra": "16454 times\n4 procs"
          },
          {
            "name": "BenchmarkDeckgenGenerateSHA - allocs/op",
            "value": 12,
            "unit": "allocs/op",
            "extra": "16454 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 185.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6578305 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 185.4,
            "unit": "ns/op",
            "extra": "6578305 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6578305 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6578305 times\n4 procs"
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
            "email": "tao.yi@konghq.com",
            "name": "Yi Tao",
            "username": "randmonkey"
          },
          "distinct": true,
          "id": "befe280cddce298822076aa2d97829381e05ea57",
          "message": "fix(translator): reject CA Secrets with multiple PEM certs (#7763)\n\n* docs: update API references and bump crd-ref-docs to 0.2.0\n\n* fix(translator): reject CA Secrets with multiple PEM certs\n\n* update CHANGELOG\n\nSigned-off-by: Jintao Zhang <zhangjintao9020@gmail.com>\n\n---------\n\nSigned-off-by: Jintao Zhang <zhangjintao9020@gmail.com>\n(cherry picked from commit 64e50933f0688324f3954145d4a2ef4705ed1c43)",
          "timestamp": "2025-11-05T16:46:11+08:00",
          "tree_id": "f829940a377313370cf4891eacc5d6c777fa27a4",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/befe280cddce298822076aa2d97829381e05ea57"
        },
        "date": 1762332514587,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkDeckgenGenerateSHA",
            "value": 70932,
            "unit": "ns/op\t   11109 B/op\t      12 allocs/op",
            "extra": "16838 times\n4 procs"
          },
          {
            "name": "BenchmarkDeckgenGenerateSHA - ns/op",
            "value": 70932,
            "unit": "ns/op",
            "extra": "16838 times\n4 procs"
          },
          {
            "name": "BenchmarkDeckgenGenerateSHA - B/op",
            "value": 11109,
            "unit": "B/op",
            "extra": "16838 times\n4 procs"
          },
          {
            "name": "BenchmarkDeckgenGenerateSHA - allocs/op",
            "value": 12,
            "unit": "allocs/op",
            "extra": "16838 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 152.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7795472 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 152.5,
            "unit": "ns/op",
            "extra": "7795472 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7795472 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7795472 times\n4 procs"
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
          "id": "8e3d9d016c535fc88b3ac2d141afab9bdf497a3b",
          "message": "fix(translator): reject CA Secrets with multiple PEM certs (#7763) (#7779)\n\n* docs: update API references and bump crd-ref-docs to 0.2.0\n\n* fix(translator): reject CA Secrets with multiple PEM certs\n\n* update CHANGELOG\n\n\n\n---------\n\n\n(cherry picked from commit 64e50933f0688324f3954145d4a2ef4705ed1c43)\n\nSigned-off-by: Jintao Zhang <zhangjintao9020@gmail.com>\nCo-authored-by: Jintao Zhang <zhangjintao9020@gmail.com>",
          "timestamp": "2025-11-05T17:42:48+08:00",
          "tree_id": "f829940a377313370cf4891eacc5d6c777fa27a4",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/8e3d9d016c535fc88b3ac2d141afab9bdf497a3b"
        },
        "date": 1762335889625,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkDeckgenGenerateSHA",
            "value": 69390,
            "unit": "ns/op\t   11109 B/op\t      12 allocs/op",
            "extra": "16788 times\n4 procs"
          },
          {
            "name": "BenchmarkDeckgenGenerateSHA - ns/op",
            "value": 69390,
            "unit": "ns/op",
            "extra": "16788 times\n4 procs"
          },
          {
            "name": "BenchmarkDeckgenGenerateSHA - B/op",
            "value": 11109,
            "unit": "B/op",
            "extra": "16788 times\n4 procs"
          },
          {
            "name": "BenchmarkDeckgenGenerateSHA - allocs/op",
            "value": 12,
            "unit": "allocs/op",
            "extra": "16788 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 140.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8910008 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 140.3,
            "unit": "ns/op",
            "extra": "8910008 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8910008 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8910008 times\n4 procs"
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
            "email": "tao.yi@konghq.com",
            "name": "Yi Tao",
            "username": "randmonkey"
          },
          "distinct": true,
          "id": "5f5f2435c6bfed208719ef68464c07d4970c8b59",
          "message": "fix(sendconfig): Do not cleanup `null`s in plugin configuration   (#7751)\n\n* remove cleanup of nulls in plugin config\n\n* update tests\n\n* update changelog\n\n* dump plugin config in tests\n\n* find plugin in multiple plugins\n\n* use different paths for plugin tests and fix workspace\n\n(cherry picked from commit 2bb2029577de2e30d3b0dea01e36baa7c82c3157)",
          "timestamp": "2025-11-05T17:50:43+08:00",
          "tree_id": "8b2ef89d682580b92c723373aea730567d0a440c",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/5f5f2435c6bfed208719ef68464c07d4970c8b59"
        },
        "date": 1762336406223,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkDeckgenGenerateSHA",
            "value": 74080,
            "unit": "ns/op\t   11117 B/op\t      12 allocs/op",
            "extra": "17101 times\n4 procs"
          },
          {
            "name": "BenchmarkDeckgenGenerateSHA - ns/op",
            "value": 74080,
            "unit": "ns/op",
            "extra": "17101 times\n4 procs"
          },
          {
            "name": "BenchmarkDeckgenGenerateSHA - B/op",
            "value": 11117,
            "unit": "B/op",
            "extra": "17101 times\n4 procs"
          },
          {
            "name": "BenchmarkDeckgenGenerateSHA - allocs/op",
            "value": 12,
            "unit": "allocs/op",
            "extra": "17101 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 64.11,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "20103526 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 64.11,
            "unit": "ns/op",
            "extra": "20103526 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "20103526 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "20103526 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "tao.yi@konghq.com",
            "name": "Yi Tao",
            "username": "randmonkey"
          },
          "committer": {
            "email": "tao.yi@konghq.com",
            "name": "Yi Tao",
            "username": "randmonkey"
          },
          "distinct": true,
          "id": "9a2c912f28f7bb5b8448cb0bc49d8836ee5bd612",
          "message": "bump go-kong to 0.60.0",
          "timestamp": "2025-11-05T18:33:07+08:00",
          "tree_id": "c10c6a407fdb08e7f833a5295ef8aa74f482ec47",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/9a2c912f28f7bb5b8448cb0bc49d8836ee5bd612"
        },
        "date": 1762338928341,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkDeckgenGenerateSHA",
            "value": 68782,
            "unit": "ns/op\t   11103 B/op\t      12 allocs/op",
            "extra": "17073 times\n4 procs"
          },
          {
            "name": "BenchmarkDeckgenGenerateSHA - ns/op",
            "value": 68782,
            "unit": "ns/op",
            "extra": "17073 times\n4 procs"
          },
          {
            "name": "BenchmarkDeckgenGenerateSHA - B/op",
            "value": 11103,
            "unit": "B/op",
            "extra": "17073 times\n4 procs"
          },
          {
            "name": "BenchmarkDeckgenGenerateSHA - allocs/op",
            "value": 12,
            "unit": "allocs/op",
            "extra": "17073 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 60,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "21464690 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 60,
            "unit": "ns/op",
            "extra": "21464690 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "21464690 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "21464690 times\n4 procs"
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
          "id": "d3cdba13d211eae1315adec2bd48d17e3d471af3",
          "message": "fix(translator): Add translate threshold in upstream policy translation (#7784)\n\n* add translate threshold in upstream policy translation\n\n* update changelog\n\n* Thanks to the contribution from elbrogan-vizio(@elbrogan-vizio) by PR #7756",
          "timestamp": "2025-11-20T03:37:46Z",
          "tree_id": "daf2f60bb9f8f505578110065fd7b0f97f0bc4b4",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/d3cdba13d211eae1315adec2bd48d17e3d471af3"
        },
        "date": 1763611694882,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 760736,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "1447 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 760736,
            "unit": "ns/op",
            "extra": "1447 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819453,
            "unit": "B/op",
            "extra": "1447 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1447 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 9123,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "113078 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 9123,
            "unit": "ns/op",
            "extra": "113078 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "113078 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "113078 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 80.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14383834 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 80.1,
            "unit": "ns/op",
            "extra": "14383834 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14383834 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14383834 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21817,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54625 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21817,
            "unit": "ns/op",
            "extra": "54625 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54625 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54625 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 205441,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4978 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 205441,
            "unit": "ns/op",
            "extra": "4978 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4978 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4978 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2264685,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "514 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2264685,
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
            "value": 40500736,
            "unit": "ns/op\t24010754 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 40500736,
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
            "value": 9957,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "101678 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9957,
            "unit": "ns/op",
            "extra": "101678 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "101678 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "101678 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7008004,
            "unit": "ns/op\t 4527543 B/op\t   69230 allocs/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7008004,
            "unit": "ns/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527543,
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
            "value": 7858,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "153405 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7858,
            "unit": "ns/op",
            "extra": "153405 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "153405 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "153405 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 840217,
            "unit": "ns/op\t  392940 B/op\t    6173 allocs/op",
            "extra": "1418 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 840217,
            "unit": "ns/op",
            "extra": "1418 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392940,
            "unit": "B/op",
            "extra": "1418 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1418 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10873,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "111178 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10873,
            "unit": "ns/op",
            "extra": "111178 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "111178 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "111178 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8003107,
            "unit": "ns/op\t 4914204 B/op\t   75241 allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8003107,
            "unit": "ns/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914204,
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
            "value": 318.9,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3792345 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 318.9,
            "unit": "ns/op",
            "extra": "3792345 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3792345 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3792345 times\n4 procs"
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
            "value": 622.6,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2111088 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 622.6,
            "unit": "ns/op",
            "extra": "2111088 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2111088 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2111088 times\n4 procs"
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
          "id": "7f165cd2cade53cb7c257b34fa17d0b0eb554327",
          "message": "chore(deps): bump golang.org/x/crypto from 0.41.0 to 0.45.0 (#7789)\n\nBumps [golang.org/x/crypto](https://github.com/golang/crypto) from 0.41.0 to 0.45.0.\n- [Commits](https://github.com/golang/crypto/compare/v0.41.0...v0.45.0)\n\n---\nupdated-dependencies:\n- dependency-name: golang.org/x/crypto\n  dependency-version: 0.45.0\n  dependency-type: indirect\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2025-11-20T10:14:58Z",
          "tree_id": "41f4709eba45b0e553b3ce75832c20137009c1be",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/7f165cd2cade53cb7c257b34fa17d0b0eb554327"
        },
        "date": 1763635599566,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 706838,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "1506 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 706838,
            "unit": "ns/op",
            "extra": "1506 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819453,
            "unit": "B/op",
            "extra": "1506 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1506 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7344,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "157170 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7344,
            "unit": "ns/op",
            "extra": "157170 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "157170 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "157170 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.62,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16402900 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.62,
            "unit": "ns/op",
            "extra": "16402900 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16402900 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16402900 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21953,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54790 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21953,
            "unit": "ns/op",
            "extra": "54790 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54790 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54790 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 212114,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4844 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 212114,
            "unit": "ns/op",
            "extra": "4844 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4844 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4844 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2297178,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "511 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2297178,
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
            "value": 37109919,
            "unit": "ns/op\t24010754 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 37109919,
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
            "value": 9616,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "123483 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9616,
            "unit": "ns/op",
            "extra": "123483 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "123483 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "123483 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7303584,
            "unit": "ns/op\t 4527544 B/op\t   69230 allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7303584,
            "unit": "ns/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527544,
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
            "value": 7974,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "147312 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7974,
            "unit": "ns/op",
            "extra": "147312 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "147312 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "147312 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 853634,
            "unit": "ns/op\t  392938 B/op\t    6173 allocs/op",
            "extra": "1387 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 853634,
            "unit": "ns/op",
            "extra": "1387 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392938,
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
            "value": 11423,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "107826 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11423,
            "unit": "ns/op",
            "extra": "107826 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "107826 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "107826 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8048944,
            "unit": "ns/op\t 4914187 B/op\t   75241 allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8048944,
            "unit": "ns/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914187,
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
            "value": 380.7,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3715480 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 380.7,
            "unit": "ns/op",
            "extra": "3715480 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3715480 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3715480 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.35,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.35,
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
            "extra": "2042050 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 577.7,
            "unit": "ns/op",
            "extra": "2042050 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2042050 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2042050 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 814.4,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1471688 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 814.4,
            "unit": "ns/op",
            "extra": "1471688 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1471688 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1471688 times\n4 procs"
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
          "id": "b49efbbab63f4f29775b1eff4ec030509e999072",
          "message": "chore(deps): bump golang from 1.25.3 to 1.25.4 (#7782)\n\nBumps golang from 1.25.3 to 1.25.4.\n\n---\nupdated-dependencies:\n- dependency-name: golang\n  dependency-version: 1.25.4\n  dependency-type: direct:production\n  update-type: version-update:semver-patch\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2025-11-20T10:16:46Z",
          "tree_id": "ada4959f43c02bce470edfeaf24d97194da11133",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/b49efbbab63f4f29775b1eff4ec030509e999072"
        },
        "date": 1763635683872,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 682564,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "1684 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 682564,
            "unit": "ns/op",
            "extra": "1684 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819453,
            "unit": "B/op",
            "extra": "1684 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1684 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7543,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "150760 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7543,
            "unit": "ns/op",
            "extra": "150760 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "150760 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "150760 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.67,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16534461 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.67,
            "unit": "ns/op",
            "extra": "16534461 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16534461 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16534461 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21876,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54048 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21876,
            "unit": "ns/op",
            "extra": "54048 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54048 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54048 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 217779,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4995 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 217779,
            "unit": "ns/op",
            "extra": "4995 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4995 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4995 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2471512,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "472 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2471512,
            "unit": "ns/op",
            "extra": "472 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "472 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "472 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 34210268,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "68 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 34210268,
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
            "value": 9654,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "125083 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9654,
            "unit": "ns/op",
            "extra": "125083 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "125083 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "125083 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7409444,
            "unit": "ns/op\t 4527547 B/op\t   69231 allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7409444,
            "unit": "ns/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527547,
            "unit": "B/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69231,
            "unit": "allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7963,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "147951 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7963,
            "unit": "ns/op",
            "extra": "147951 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "147951 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "147951 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 859009,
            "unit": "ns/op\t  392942 B/op\t    6173 allocs/op",
            "extra": "1374 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 859009,
            "unit": "ns/op",
            "extra": "1374 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392942,
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
            "value": 11356,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "104676 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11356,
            "unit": "ns/op",
            "extra": "104676 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "104676 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "104676 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8360738,
            "unit": "ns/op\t 4914250 B/op\t   75241 allocs/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8360738,
            "unit": "ns/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914250,
            "unit": "B/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 349.6,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3453697 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 349.6,
            "unit": "ns/op",
            "extra": "3453697 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3453697 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3453697 times\n4 procs"
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
            "value": 565.7,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2123264 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 565.7,
            "unit": "ns/op",
            "extra": "2123264 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2123264 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2123264 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 816.4,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1464795 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 816.4,
            "unit": "ns/op",
            "extra": "1464795 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1464795 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1464795 times\n4 procs"
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
          "id": "2fa8441a66783bcf6e8e999ba8e167dd814dddd1",
          "message": "chore(deps): update dependency rhysd/actionlint to v1.7.8 (#7769)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2025-11-20T10:41:57Z",
          "tree_id": "ab88c1025bee76dcef1477f295bc8689633a819b",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/2fa8441a66783bcf6e8e999ba8e167dd814dddd1"
        },
        "date": 1763637224826,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 804537,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "1346 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 804537,
            "unit": "ns/op",
            "extra": "1346 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819453,
            "unit": "B/op",
            "extra": "1346 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1346 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7819,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "149780 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7819,
            "unit": "ns/op",
            "extra": "149780 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "149780 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "149780 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 66.41,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18179055 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 66.41,
            "unit": "ns/op",
            "extra": "18179055 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "18179055 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "18179055 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21907,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54241 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21907,
            "unit": "ns/op",
            "extra": "54241 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54241 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54241 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 271690,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "3914 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 271690,
            "unit": "ns/op",
            "extra": "3914 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "3914 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3914 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 3111603,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "384 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 3111603,
            "unit": "ns/op",
            "extra": "384 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408449,
            "unit": "B/op",
            "extra": "384 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "384 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 55471776,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "56 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 55471776,
            "unit": "ns/op",
            "extra": "56 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010753,
            "unit": "B/op",
            "extra": "56 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "56 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 11006,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "109179 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 11006,
            "unit": "ns/op",
            "extra": "109179 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "109179 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "109179 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 8213814,
            "unit": "ns/op\t 4527553 B/op\t   69231 allocs/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 8213814,
            "unit": "ns/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4527553,
            "unit": "B/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 69231,
            "unit": "allocs/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8788,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "136180 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8788,
            "unit": "ns/op",
            "extra": "136180 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "136180 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "136180 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 984617,
            "unit": "ns/op\t  392926 B/op\t    6173 allocs/op",
            "extra": "1198 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 984617,
            "unit": "ns/op",
            "extra": "1198 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392926,
            "unit": "B/op",
            "extra": "1198 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1198 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 12000,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "99884 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 12000,
            "unit": "ns/op",
            "extra": "99884 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "99884 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "99884 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 9680707,
            "unit": "ns/op\t 4914216 B/op\t   75241 allocs/op",
            "extra": "124 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 9680707,
            "unit": "ns/op",
            "extra": "124 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914216,
            "unit": "B/op",
            "extra": "124 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 75241,
            "unit": "allocs/op",
            "extra": "124 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 323.7,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3689427 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 323.7,
            "unit": "ns/op",
            "extra": "3689427 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3689427 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3689427 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 11.26,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 11.26,
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
            "value": 573.6,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2065282 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 573.6,
            "unit": "ns/op",
            "extra": "2065282 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2065282 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2065282 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 804,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1507196 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 804,
            "unit": "ns/op",
            "extra": "1507196 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1507196 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1507196 times\n4 procs"
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
          "id": "4d15e6305213ef892fe4626ded471285c7f13bbd",
          "message": "chore(deps): update dependency gke to v1.34.1 (#7797)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2025-11-21T06:38:31Z",
          "tree_id": "eb0d17512b334a4a84743ce386cc4c4fe541c311",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/4d15e6305213ef892fe4626ded471285c7f13bbd"
        },
        "date": 1763708968045,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 657307,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "1563 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 657307,
            "unit": "ns/op",
            "extra": "1563 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819453,
            "unit": "B/op",
            "extra": "1563 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1563 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7536,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "158678 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7536,
            "unit": "ns/op",
            "extra": "158678 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "158678 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "158678 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.61,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16539066 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.61,
            "unit": "ns/op",
            "extra": "16539066 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16539066 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16539066 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21553,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "55383 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21553,
            "unit": "ns/op",
            "extra": "55383 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "55383 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "55383 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 210721,
            "unit": "ns/op\t  245761 B/op\t       2 allocs/op",
            "extra": "4910 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 210721,
            "unit": "ns/op",
            "extra": "4910 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245761,
            "unit": "B/op",
            "extra": "4910 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4910 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2413633,
            "unit": "ns/op\t 2408449 B/op\t       2 allocs/op",
            "extra": "444 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2413633,
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
            "value": 36084408,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "28 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 36084408,
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
            "value": 9830,
            "unit": "ns/op\t    7952 B/op\t     185 allocs/op",
            "extra": "124550 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9830,
            "unit": "ns/op",
            "extra": "124550 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7952,
            "unit": "B/op",
            "extra": "124550 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 185,
            "unit": "allocs/op",
            "extra": "124550 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7001501,
            "unit": "ns/op\t 4527547 B/op\t   69230 allocs/op",
            "extra": "170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7001501,
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
            "value": 7922,
            "unit": "ns/op\t    6704 B/op\t     171 allocs/op",
            "extra": "152608 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7922,
            "unit": "ns/op",
            "extra": "152608 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6704,
            "unit": "B/op",
            "extra": "152608 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 171,
            "unit": "allocs/op",
            "extra": "152608 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 845025,
            "unit": "ns/op\t  392933 B/op\t    6173 allocs/op",
            "extra": "1388 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 845025,
            "unit": "ns/op",
            "extra": "1388 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392933,
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
            "value": 11061,
            "unit": "ns/op\t    8368 B/op\t     193 allocs/op",
            "extra": "107298 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11061,
            "unit": "ns/op",
            "extra": "107298 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8368,
            "unit": "B/op",
            "extra": "107298 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 193,
            "unit": "allocs/op",
            "extra": "107298 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8137274,
            "unit": "ns/op\t 4914219 B/op\t   75241 allocs/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8137274,
            "unit": "ns/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4914219,
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
            "value": 320.4,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "3743313 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 320.4,
            "unit": "ns/op",
            "extra": "3743313 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "3743313 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3743313 times\n4 procs"
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
            "value": 561.8,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2138028 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 561.8,
            "unit": "ns/op",
            "extra": "2138028 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2138028 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2138028 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 818.7,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1462424 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 818.7,
            "unit": "ns/op",
            "extra": "1462424 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1462424 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1462424 times\n4 procs"
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
      }
    ]
  }
}