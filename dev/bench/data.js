window.BENCHMARK_DATA = {
  "lastUpdate": 1734435270996,
  "repoUrl": "https://github.com/Kong/kubernetes-ingress-controller",
  "entries": {
    "Go Benchmark": [
      {
        "commit": {
          "author": {
            "email": "lavacca.mattia@gmail.com",
            "name": "Mattia Lavacca",
            "username": "mlavacca"
          },
          "committer": {
            "email": "lavacca.mattia@gmail.com",
            "name": "Mattia Lavacca",
            "username": "mlavacca"
          },
          "distinct": true,
          "id": "ff047c96e916dfc952fbc8baae7b6b24b0a93841",
          "message": "configmap for CA certificates complete\n\nSigned-off-by: Mattia Lavacca <lavacca.mattia@gmail.com>",
          "timestamp": "2024-12-06T13:56:27+01:00",
          "tree_id": "bc9a069879a8f7c4d5ac467333ac5a8357af33a9",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/ff047c96e916dfc952fbc8baae7b6b24b0a93841"
        },
        "date": 1733489973822,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1459536,
            "unit": "ns/op\t  819441 B/op\t       5 allocs/op",
            "extra": "955 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1459536,
            "unit": "ns/op",
            "extra": "955 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819441,
            "unit": "B/op",
            "extra": "955 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "955 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7308,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "163861 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7308,
            "unit": "ns/op",
            "extra": "163861 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "163861 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "163861 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13804606 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.7,
            "unit": "ns/op",
            "extra": "13804606 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13804606 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13804606 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22305,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52594 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22305,
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
            "value": 220837,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4838 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 220837,
            "unit": "ns/op",
            "extra": "4838 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4838 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4838 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2560779,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "459 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2560779,
            "unit": "ns/op",
            "extra": "459 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "459 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "459 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 39535018,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "68 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 39535018,
            "unit": "ns/op",
            "extra": "68 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 10052,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "118285 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10052,
            "unit": "ns/op",
            "extra": "118285 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "118285 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "118285 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7489097,
            "unit": "ns/op\t 4594801 B/op\t   75254 allocs/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7489097,
            "unit": "ns/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594801,
            "unit": "B/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75254,
            "unit": "allocs/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8378,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "140980 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8378,
            "unit": "ns/op",
            "extra": "140980 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "140980 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "140980 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 878269,
            "unit": "ns/op\t  396800 B/op\t    6234 allocs/op",
            "extra": "1212 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 878269,
            "unit": "ns/op",
            "extra": "1212 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396800,
            "unit": "B/op",
            "extra": "1212 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6234,
            "unit": "allocs/op",
            "extra": "1212 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11420,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "103514 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11420,
            "unit": "ns/op",
            "extra": "103514 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "103514 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "103514 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8470681,
            "unit": "ns/op\t 4981266 B/op\t   81263 allocs/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8470681,
            "unit": "ns/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981266,
            "unit": "B/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81263,
            "unit": "allocs/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 227.3,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5236648 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 227.3,
            "unit": "ns/op",
            "extra": "5236648 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5236648 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5236648 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6194,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6194,
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
            "value": 294.9,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4099933 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 294.9,
            "unit": "ns/op",
            "extra": "4099933 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4099933 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4099933 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 458.9,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2686168 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 458.9,
            "unit": "ns/op",
            "extra": "2686168 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2686168 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2686168 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "lavacca.mattia@gmail.com",
            "name": "Mattia Lavacca",
            "username": "mlavacca"
          },
          "committer": {
            "email": "lavacca.mattia@gmail.com",
            "name": "Mattia Lavacca",
            "username": "mlavacca"
          },
          "distinct": true,
          "id": "a0be3e4da552ea9893db4031c128789411fb1ca5",
          "message": "configmap for CA certificates complete\n\nSigned-off-by: Mattia Lavacca <lavacca.mattia@gmail.com>",
          "timestamp": "2024-12-06T14:01:54+01:00",
          "tree_id": "0a77bcedde3e95ff1a1cf547efc83418c73335b2",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/a0be3e4da552ea9893db4031c128789411fb1ca5"
        },
        "date": 1733490300162,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1138885,
            "unit": "ns/op\t  819440 B/op\t       5 allocs/op",
            "extra": "1008 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1138885,
            "unit": "ns/op",
            "extra": "1008 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819440,
            "unit": "B/op",
            "extra": "1008 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1008 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 9232,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "114627 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 9232,
            "unit": "ns/op",
            "extra": "114627 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "114627 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "114627 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 78.95,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15215827 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 78.95,
            "unit": "ns/op",
            "extra": "15215827 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15215827 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15215827 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22874,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52789 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22874,
            "unit": "ns/op",
            "extra": "52789 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52789 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52789 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 320517,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "3133 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 320517,
            "unit": "ns/op",
            "extra": "3133 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "3133 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3133 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2586343,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "466 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2586343,
            "unit": "ns/op",
            "extra": "466 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
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
            "value": 43478322,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 43478322,
            "unit": "ns/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 10136,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "115623 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10136,
            "unit": "ns/op",
            "extra": "115623 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "115623 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "115623 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7921011,
            "unit": "ns/op\t 4594314 B/op\t   75252 allocs/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7921011,
            "unit": "ns/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594314,
            "unit": "B/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75252,
            "unit": "allocs/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8486,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "141631 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8486,
            "unit": "ns/op",
            "extra": "141631 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "141631 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "141631 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 885822,
            "unit": "ns/op\t  396744 B/op\t    6233 allocs/op",
            "extra": "1224 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 885822,
            "unit": "ns/op",
            "extra": "1224 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396744,
            "unit": "B/op",
            "extra": "1224 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6233,
            "unit": "allocs/op",
            "extra": "1224 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11930,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "101541 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11930,
            "unit": "ns/op",
            "extra": "101541 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "101541 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "101541 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8943254,
            "unit": "ns/op\t 4981529 B/op\t   81264 allocs/op",
            "extra": "132 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8943254,
            "unit": "ns/op",
            "extra": "132 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981529,
            "unit": "B/op",
            "extra": "132 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81264,
            "unit": "allocs/op",
            "extra": "132 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 230.3,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5246916 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 230.3,
            "unit": "ns/op",
            "extra": "5246916 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5246916 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5246916 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6203,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6203,
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
            "value": 306.2,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4014692 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 306.2,
            "unit": "ns/op",
            "extra": "4014692 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4014692 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4014692 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 453.4,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2669840 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 453.4,
            "unit": "ns/op",
            "extra": "2669840 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2669840 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2669840 times\n4 procs"
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
          "id": "9031b7b77e7f5c5a307b11f1882278e4da73e67a",
          "message": "chore(ci): renovate update kubernetes-configuration to unstable versions (#6789)",
          "timestamp": "2024-12-06T16:36:18+01:00",
          "tree_id": "90740db9d4300c3cda67edde040015bf1bbcd059",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/9031b7b77e7f5c5a307b11f1882278e4da73e67a"
        },
        "date": 1733499586089,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1449977,
            "unit": "ns/op\t  819447 B/op\t       5 allocs/op",
            "extra": "909 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1449977,
            "unit": "ns/op",
            "extra": "909 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819447,
            "unit": "B/op",
            "extra": "909 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "909 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7424,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "162100 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7424,
            "unit": "ns/op",
            "extra": "162100 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "162100 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "162100 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 78.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15104344 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 78.8,
            "unit": "ns/op",
            "extra": "15104344 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15104344 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15104344 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22366,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "51957 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22366,
            "unit": "ns/op",
            "extra": "51957 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "51957 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "51957 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 224549,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5144 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 224549,
            "unit": "ns/op",
            "extra": "5144 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5144 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5144 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2632300,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "397 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2632300,
            "unit": "ns/op",
            "extra": "397 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
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
            "value": 36910658,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "28 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 36910658,
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
            "value": 9504,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "124910 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9504,
            "unit": "ns/op",
            "extra": "124910 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "124910 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "124910 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7328766,
            "unit": "ns/op\t 4594468 B/op\t   75247 allocs/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7328766,
            "unit": "ns/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594468,
            "unit": "B/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75247,
            "unit": "allocs/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7913,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "149848 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7913,
            "unit": "ns/op",
            "extra": "149848 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "149848 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "149848 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 850627,
            "unit": "ns/op\t  396432 B/op\t    6226 allocs/op",
            "extra": "1248 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 850627,
            "unit": "ns/op",
            "extra": "1248 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396432,
            "unit": "B/op",
            "extra": "1248 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6226,
            "unit": "allocs/op",
            "extra": "1248 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10911,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "109500 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10911,
            "unit": "ns/op",
            "extra": "109500 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "109500 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "109500 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8420384,
            "unit": "ns/op\t 4981379 B/op\t   81259 allocs/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8420384,
            "unit": "ns/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981379,
            "unit": "B/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81259,
            "unit": "allocs/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 230.5,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5294964 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 230.5,
            "unit": "ns/op",
            "extra": "5294964 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5294964 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5294964 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6197,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6197,
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
            "value": 337.8,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4063923 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 337.8,
            "unit": "ns/op",
            "extra": "4063923 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4063923 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4063923 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 449.7,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2676600 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 449.7,
            "unit": "ns/op",
            "extra": "2676600 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2676600 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2676600 times\n4 procs"
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
          "id": "6704a8b4155fdcd2ff8366f3ec1ac6067289cd6c",
          "message": "chore(deps): bump github.com/kong/go-kong from 0.60.0 to 0.61.0 (#6788)\n\nBumps [github.com/kong/go-kong](https://github.com/kong/go-kong) from 0.60.0 to 0.61.0.\r\n- [Release notes](https://github.com/kong/go-kong/releases)\r\n- [Changelog](https://github.com/Kong/go-kong/blob/main/CHANGELOG.md)\r\n- [Commits](https://github.com/kong/go-kong/compare/v0.60.0...v0.61.0)\r\n\r\n---\r\nupdated-dependencies:\r\n- dependency-name: github.com/kong/go-kong\r\n  dependency-type: direct:production\r\n  update-type: version-update:semver-minor\r\n...\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2024-12-06T15:58:41Z",
          "tree_id": "910e86507a8df0b672833de4401405513c5fcf74",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/6704a8b4155fdcd2ff8366f3ec1ac6067289cd6c"
        },
        "date": 1733500907871,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1146610,
            "unit": "ns/op\t  819442 B/op\t       5 allocs/op",
            "extra": "946 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1146610,
            "unit": "ns/op",
            "extra": "946 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819442,
            "unit": "B/op",
            "extra": "946 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "946 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 9379,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "108258 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 9379,
            "unit": "ns/op",
            "extra": "108258 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "108258 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "108258 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 80.05,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14924242 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 80.05,
            "unit": "ns/op",
            "extra": "14924242 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14924242 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14924242 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22490,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53061 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22490,
            "unit": "ns/op",
            "extra": "53061 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53061 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53061 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 225012,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5148 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 225012,
            "unit": "ns/op",
            "extra": "5148 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5148 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5148 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2695614,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "445 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2695614,
            "unit": "ns/op",
            "extra": "445 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "445 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "445 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 38316326,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 38316326,
            "unit": "ns/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 9707,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "123596 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9707,
            "unit": "ns/op",
            "extra": "123596 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "123596 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "123596 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7445203,
            "unit": "ns/op\t 4594296 B/op\t   75247 allocs/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7445203,
            "unit": "ns/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594296,
            "unit": "B/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75247,
            "unit": "allocs/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8048,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "147744 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8048,
            "unit": "ns/op",
            "extra": "147744 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "147744 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "147744 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 864667,
            "unit": "ns/op\t  396632 B/op\t    6229 allocs/op",
            "extra": "1191 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 864667,
            "unit": "ns/op",
            "extra": "1191 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396632,
            "unit": "B/op",
            "extra": "1191 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6229,
            "unit": "allocs/op",
            "extra": "1191 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11092,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "108094 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11092,
            "unit": "ns/op",
            "extra": "108094 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "108094 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "108094 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8639186,
            "unit": "ns/op\t 4981149 B/op\t   81258 allocs/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8639186,
            "unit": "ns/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981149,
            "unit": "B/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81258,
            "unit": "allocs/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 228.9,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5275154 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 228.9,
            "unit": "ns/op",
            "extra": "5275154 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5275154 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5275154 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6217,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6217,
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
            "value": 295.8,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4080502 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 295.8,
            "unit": "ns/op",
            "extra": "4080502 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4080502 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4080502 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 449,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2679910 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 449,
            "unit": "ns/op",
            "extra": "2679910 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2679910 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2679910 times\n4 procs"
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
          "id": "9fa371ac1053b7ad513fc3f37fc10c9297cd6454",
          "message": "chore(deps): update dependency mikefarah/yq to v4.44.6 (#6796)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2024-12-09T09:34:07+01:00",
          "tree_id": "1e6b154b4af5ce95bd21887a8393a6071585443c",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/9fa371ac1053b7ad513fc3f37fc10c9297cd6454"
        },
        "date": 1733733441484,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1165550,
            "unit": "ns/op\t  819442 B/op\t       5 allocs/op",
            "extra": "999 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1165550,
            "unit": "ns/op",
            "extra": "999 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819442,
            "unit": "B/op",
            "extra": "999 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "999 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7403,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "158524 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7403,
            "unit": "ns/op",
            "extra": "158524 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "158524 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "158524 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.62,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14993121 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.62,
            "unit": "ns/op",
            "extra": "14993121 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14993121 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14993121 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22694,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52353 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22694,
            "unit": "ns/op",
            "extra": "52353 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52353 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52353 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 229524,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5186 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 229524,
            "unit": "ns/op",
            "extra": "5186 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5186 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5186 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2553906,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "462 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2553906,
            "unit": "ns/op",
            "extra": "462 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "462 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "462 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 34918021,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 34918021,
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
            "value": 10023,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "118675 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10023,
            "unit": "ns/op",
            "extra": "118675 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "118675 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "118675 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7645411,
            "unit": "ns/op\t 4594480 B/op\t   75247 allocs/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7645411,
            "unit": "ns/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594480,
            "unit": "B/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75247,
            "unit": "allocs/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8215,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "145006 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8215,
            "unit": "ns/op",
            "extra": "145006 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "145006 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "145006 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 874183,
            "unit": "ns/op\t  396471 B/op\t    6226 allocs/op",
            "extra": "1237 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 874183,
            "unit": "ns/op",
            "extra": "1237 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396471,
            "unit": "B/op",
            "extra": "1237 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6226,
            "unit": "allocs/op",
            "extra": "1237 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11206,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "105024 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11206,
            "unit": "ns/op",
            "extra": "105024 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "105024 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "105024 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8747125,
            "unit": "ns/op\t 4981141 B/op\t   81258 allocs/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8747125,
            "unit": "ns/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981141,
            "unit": "B/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81258,
            "unit": "allocs/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 240.8,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5201721 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 240.8,
            "unit": "ns/op",
            "extra": "5201721 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5201721 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5201721 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6526,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6526,
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
            "value": 300.7,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "3589444 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 300.7,
            "unit": "ns/op",
            "extra": "3589444 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "3589444 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "3589444 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 452.1,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2690300 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 452.1,
            "unit": "ns/op",
            "extra": "2690300 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2690300 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2690300 times\n4 procs"
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
          "id": "7bf240cffd5e4a0cd45fad1a6b52a19c81105772",
          "message": "fix: verify if the referenced KongPlugin/KongClusterPlugin exists and set Programmed condition accordingly (#6791)\n\n* fix: verify if the referenced KongPlugin/KongClusterPlugin exists and set Programmed condition accordingly\r\n\r\n* Update CHANGELOG.md\r\n\r\n---------\r\n\r\nCo-authored-by: Tao Yi <tao.yi@konghq.com>",
          "timestamp": "2024-12-09T09:14:57Z",
          "tree_id": "091802deed728cfa8e234e30dc814ce179095e74",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/7bf240cffd5e4a0cd45fad1a6b52a19c81105772"
        },
        "date": 1733735896288,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1187048,
            "unit": "ns/op\t  819443 B/op\t       5 allocs/op",
            "extra": "979 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1187048,
            "unit": "ns/op",
            "extra": "979 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819443,
            "unit": "B/op",
            "extra": "979 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "979 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 9399,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "107607 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 9399,
            "unit": "ns/op",
            "extra": "107607 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "107607 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "107607 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 78.87,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14960872 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 78.87,
            "unit": "ns/op",
            "extra": "14960872 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14960872 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14960872 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22657,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52824 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22657,
            "unit": "ns/op",
            "extra": "52824 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52824 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52824 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 225128,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5199 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 225128,
            "unit": "ns/op",
            "extra": "5199 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5199 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5199 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2548148,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "474 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2548148,
            "unit": "ns/op",
            "extra": "474 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "474 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "474 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 36390883,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "67 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 36390883,
            "unit": "ns/op",
            "extra": "67 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 9773,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "123674 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9773,
            "unit": "ns/op",
            "extra": "123674 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "123674 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "123674 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7480507,
            "unit": "ns/op\t 4594484 B/op\t   75248 allocs/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7480507,
            "unit": "ns/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594484,
            "unit": "B/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75248,
            "unit": "allocs/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8134,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "147009 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8134,
            "unit": "ns/op",
            "extra": "147009 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "147009 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "147009 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 872436,
            "unit": "ns/op\t  396454 B/op\t    6226 allocs/op",
            "extra": "1242 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 872436,
            "unit": "ns/op",
            "extra": "1242 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396454,
            "unit": "B/op",
            "extra": "1242 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6226,
            "unit": "allocs/op",
            "extra": "1242 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11076,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "108318 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11076,
            "unit": "ns/op",
            "extra": "108318 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "108318 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "108318 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8497331,
            "unit": "ns/op\t 4981269 B/op\t   81258 allocs/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8497331,
            "unit": "ns/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981269,
            "unit": "B/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81258,
            "unit": "allocs/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 226.6,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5287880 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 226.6,
            "unit": "ns/op",
            "extra": "5287880 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5287880 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5287880 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.62,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.62,
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
            "value": 294.2,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4038254 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 294.2,
            "unit": "ns/op",
            "extra": "4038254 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4038254 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4038254 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 446.8,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2660138 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 446.8,
            "unit": "ns/op",
            "extra": "2660138 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2660138 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2660138 times\n4 procs"
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
          "id": "67990d91b24bb25fff016deecb827cf34915ff5b",
          "message": "chore(deps): update dependency gke to v1.31.3 (#6792)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2024-12-09T09:57:52Z",
          "tree_id": "d594dc9bb9ead11d88c83d24d79d9eef4f881368",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/67990d91b24bb25fff016deecb827cf34915ff5b"
        },
        "date": 1733738464800,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1174529,
            "unit": "ns/op\t  819443 B/op\t       5 allocs/op",
            "extra": "986 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1174529,
            "unit": "ns/op",
            "extra": "986 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819443,
            "unit": "B/op",
            "extra": "986 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "986 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7216,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "166713 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7216,
            "unit": "ns/op",
            "extra": "166713 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "166713 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "166713 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 81.66,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15151464 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 81.66,
            "unit": "ns/op",
            "extra": "15151464 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15151464 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15151464 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22367,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52782 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22367,
            "unit": "ns/op",
            "extra": "52782 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52782 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52782 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 270522,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4947 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 270522,
            "unit": "ns/op",
            "extra": "4947 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4947 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4947 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2536341,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "463 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2536341,
            "unit": "ns/op",
            "extra": "463 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "463 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "463 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 40656630,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "68 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 40656630,
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
            "value": 9517,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "125560 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9517,
            "unit": "ns/op",
            "extra": "125560 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "125560 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "125560 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7329443,
            "unit": "ns/op\t 4594638 B/op\t   75248 allocs/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7329443,
            "unit": "ns/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594638,
            "unit": "B/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75248,
            "unit": "allocs/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7896,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "153711 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7896,
            "unit": "ns/op",
            "extra": "153711 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "153711 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "153711 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 848090,
            "unit": "ns/op\t  396381 B/op\t    6225 allocs/op",
            "extra": "1267 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 848090,
            "unit": "ns/op",
            "extra": "1267 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396381,
            "unit": "B/op",
            "extra": "1267 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6225,
            "unit": "allocs/op",
            "extra": "1267 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10844,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "109468 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10844,
            "unit": "ns/op",
            "extra": "109468 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "109468 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "109468 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8410484,
            "unit": "ns/op\t 4981222 B/op\t   81258 allocs/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8410484,
            "unit": "ns/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981222,
            "unit": "B/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81258,
            "unit": "allocs/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 227.1,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5276564 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 227.1,
            "unit": "ns/op",
            "extra": "5276564 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5276564 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5276564 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6247,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6247,
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
            "value": 298.4,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4076534 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 298.4,
            "unit": "ns/op",
            "extra": "4076534 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4076534 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4076534 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 493.3,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2604814 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 493.3,
            "unit": "ns/op",
            "extra": "2604814 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2604814 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2604814 times\n4 procs"
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
          "id": "d6f7d65951edbead7f989b3f509f743dcd95f1ed",
          "message": "feat: Add `--secret-label-selector` to set the label selector for `Secrets` to ingest (#6795)\n\nBy setting this flag, the secrets that are ingested will be limited to those having this label set to \"true\".\r\nThis can reduce the memory usage in scenarios with a large number of giant secrets.\r\n\r\nSigned-off-by: Jintao Zhang <zhangjintao9020@gmail.com>\r\nCo-authored-by: Patryk Małek <patryk.malek@konghq.com>\r\nCo-authored-by: Tao Yi <tao.yi@konghq.com>",
          "timestamp": "2024-12-09T10:32:46Z",
          "tree_id": "9ab7b67b728d96d7f607a31ca64dc860b0d2ecd2",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/d6f7d65951edbead7f989b3f509f743dcd95f1ed"
        },
        "date": 1733740559172,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1442232,
            "unit": "ns/op\t  819446 B/op\t       5 allocs/op",
            "extra": "1119 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1442232,
            "unit": "ns/op",
            "extra": "1119 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819446,
            "unit": "B/op",
            "extra": "1119 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1119 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7286,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "164348 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7286,
            "unit": "ns/op",
            "extra": "164348 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "164348 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "164348 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 78.94,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15232450 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 78.94,
            "unit": "ns/op",
            "extra": "15232450 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15232450 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15232450 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22794,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52828 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22794,
            "unit": "ns/op",
            "extra": "52828 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52828 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52828 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 278007,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5658 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 278007,
            "unit": "ns/op",
            "extra": "5658 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5658 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5658 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2835260,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "450 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2835260,
            "unit": "ns/op",
            "extra": "450 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "450 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "450 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 34329888,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 34329888,
            "unit": "ns/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 9875,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "122522 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9875,
            "unit": "ns/op",
            "extra": "122522 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "122522 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "122522 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7468642,
            "unit": "ns/op\t 4594599 B/op\t   75248 allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7468642,
            "unit": "ns/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594599,
            "unit": "B/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75248,
            "unit": "allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8069,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "144451 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8069,
            "unit": "ns/op",
            "extra": "144451 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "144451 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "144451 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 876037,
            "unit": "ns/op\t  396540 B/op\t    6227 allocs/op",
            "extra": "1221 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 876037,
            "unit": "ns/op",
            "extra": "1221 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396540,
            "unit": "B/op",
            "extra": "1221 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6227,
            "unit": "allocs/op",
            "extra": "1221 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11096,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "106843 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11096,
            "unit": "ns/op",
            "extra": "106843 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "106843 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "106843 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8747079,
            "unit": "ns/op\t 4981156 B/op\t   81258 allocs/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8747079,
            "unit": "ns/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981156,
            "unit": "B/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81258,
            "unit": "allocs/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 231.4,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5214784 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 231.4,
            "unit": "ns/op",
            "extra": "5214784 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5214784 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5214784 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6196,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6196,
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
            "value": 299.5,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "3925424 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 299.5,
            "unit": "ns/op",
            "extra": "3925424 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "3925424 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "3925424 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 454.4,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2650081 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 454.4,
            "unit": "ns/op",
            "extra": "2650081 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2650081 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2650081 times\n4 procs"
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
          "id": "1523c85306b4074feb5404674ba5cc3ac5034ae1",
          "message": "fix: resolve panic in resourceErrorFromEntityAction (#6800)",
          "timestamp": "2024-12-09T15:56:32+01:00",
          "tree_id": "56173b40c5932a5a2560b0b60ad51605e6ad4197",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/1523c85306b4074feb5404674ba5cc3ac5034ae1"
        },
        "date": 1733756383429,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1153426,
            "unit": "ns/op\t  819447 B/op\t       5 allocs/op",
            "extra": "967 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1153426,
            "unit": "ns/op",
            "extra": "967 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819447,
            "unit": "B/op",
            "extra": "967 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "967 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 6824,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "175572 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 6824,
            "unit": "ns/op",
            "extra": "175572 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "175572 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "175572 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 78.92,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15186282 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 78.92,
            "unit": "ns/op",
            "extra": "15186282 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15186282 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15186282 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 24295,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52666 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 24295,
            "unit": "ns/op",
            "extra": "52666 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52666 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52666 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 255428,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5104 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 255428,
            "unit": "ns/op",
            "extra": "5104 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5104 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5104 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2552519,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "457 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2552519,
            "unit": "ns/op",
            "extra": "457 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "457 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "457 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 35133221,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "66 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 35133221,
            "unit": "ns/op",
            "extra": "66 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
            "unit": "B/op",
            "extra": "66 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "66 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9697,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "119125 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9697,
            "unit": "ns/op",
            "extra": "119125 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "119125 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "119125 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7515137,
            "unit": "ns/op\t 4594617 B/op\t   75248 allocs/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7515137,
            "unit": "ns/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594617,
            "unit": "B/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75248,
            "unit": "allocs/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8067,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "147837 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8067,
            "unit": "ns/op",
            "extra": "147837 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "147837 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "147837 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 876286,
            "unit": "ns/op\t  396639 B/op\t    6229 allocs/op",
            "extra": "1190 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 876286,
            "unit": "ns/op",
            "extra": "1190 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396639,
            "unit": "B/op",
            "extra": "1190 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6229,
            "unit": "allocs/op",
            "extra": "1190 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11492,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "102883 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11492,
            "unit": "ns/op",
            "extra": "102883 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "102883 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "102883 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8672770,
            "unit": "ns/op\t 4981059 B/op\t   81257 allocs/op",
            "extra": "135 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8672770,
            "unit": "ns/op",
            "extra": "135 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981059,
            "unit": "B/op",
            "extra": "135 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81257,
            "unit": "allocs/op",
            "extra": "135 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 227.7,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5207916 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 227.7,
            "unit": "ns/op",
            "extra": "5207916 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5207916 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5207916 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6198,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6198,
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
            "value": 295.7,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4099551 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 295.7,
            "unit": "ns/op",
            "extra": "4099551 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4099551 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4099551 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 548.8,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2681275 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 548.8,
            "unit": "ns/op",
            "extra": "2681275 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2681275 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2681275 times\n4 procs"
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
          "id": "5ac94d2468e9ff66b3142fd20c955b8903044968",
          "message": "fix(debug): fix /raw-error endpoint not serving raw response body (#6801)",
          "timestamp": "2024-12-09T16:36:28+01:00",
          "tree_id": "36a341e240b2e41011307973965b5de7966329aa",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/5ac94d2468e9ff66b3142fd20c955b8903044968"
        },
        "date": 1733758805620,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1236190,
            "unit": "ns/op\t  819442 B/op\t       5 allocs/op",
            "extra": "1208 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1236190,
            "unit": "ns/op",
            "extra": "1208 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819442,
            "unit": "B/op",
            "extra": "1208 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1208 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 6904,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "174867 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 6904,
            "unit": "ns/op",
            "extra": "174867 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "174867 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "174867 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 80.25,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14960245 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 80.25,
            "unit": "ns/op",
            "extra": "14960245 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14960245 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14960245 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22689,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53023 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22689,
            "unit": "ns/op",
            "extra": "53023 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53023 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53023 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 228560,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4579 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 228560,
            "unit": "ns/op",
            "extra": "4579 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4579 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4579 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2611333,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "415 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2611333,
            "unit": "ns/op",
            "extra": "415 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
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
            "value": 38470202,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "26 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 38470202,
            "unit": "ns/op",
            "extra": "26 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
            "unit": "B/op",
            "extra": "26 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "26 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9482,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "123952 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9482,
            "unit": "ns/op",
            "extra": "123952 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "123952 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "123952 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7340940,
            "unit": "ns/op\t 4594619 B/op\t   75248 allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7340940,
            "unit": "ns/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594619,
            "unit": "B/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75248,
            "unit": "allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7860,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "150594 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7860,
            "unit": "ns/op",
            "extra": "150594 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "150594 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "150594 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 852190,
            "unit": "ns/op\t  396444 B/op\t    6226 allocs/op",
            "extra": "1244 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 852190,
            "unit": "ns/op",
            "extra": "1244 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396444,
            "unit": "B/op",
            "extra": "1244 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6226,
            "unit": "allocs/op",
            "extra": "1244 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10913,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "107982 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10913,
            "unit": "ns/op",
            "extra": "107982 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "107982 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "107982 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8372611,
            "unit": "ns/op\t 4981144 B/op\t   81258 allocs/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8372611,
            "unit": "ns/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981144,
            "unit": "B/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81258,
            "unit": "allocs/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 233.7,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5194551 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 233.7,
            "unit": "ns/op",
            "extra": "5194551 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5194551 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5194551 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6201,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6201,
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
            "value": 295.1,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "3989481 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 295.1,
            "unit": "ns/op",
            "extra": "3989481 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "3989481 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "3989481 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 514.5,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2062425 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 514.5,
            "unit": "ns/op",
            "extra": "2062425 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2062425 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2062425 times\n4 procs"
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
          "id": "f343327f2379f0485492168c86caecefe8c725be",
          "message": "chore(logs): DEBUG instead of ERROR for K8s Event push (#6790)",
          "timestamp": "2024-12-10T16:14:33+08:00",
          "tree_id": "08229dcebd1a20400810fc15ff762c2d036a7b30",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/f343327f2379f0485492168c86caecefe8c725be"
        },
        "date": 1733818660898,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1122150,
            "unit": "ns/op\t  819447 B/op\t       5 allocs/op",
            "extra": "1002 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1122150,
            "unit": "ns/op",
            "extra": "1002 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819447,
            "unit": "B/op",
            "extra": "1002 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1002 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 9794,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "105543 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 9794,
            "unit": "ns/op",
            "extra": "105543 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "105543 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "105543 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 81.63,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14705521 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 81.63,
            "unit": "ns/op",
            "extra": "14705521 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14705521 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14705521 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22603,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52388 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22603,
            "unit": "ns/op",
            "extra": "52388 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52388 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52388 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 225813,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5204 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 225813,
            "unit": "ns/op",
            "extra": "5204 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5204 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5204 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2558056,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "458 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2558056,
            "unit": "ns/op",
            "extra": "458 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
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
            "value": 37687856,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 37687856,
            "unit": "ns/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 9759,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "121694 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9759,
            "unit": "ns/op",
            "extra": "121694 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "121694 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "121694 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7549783,
            "unit": "ns/op\t 4594958 B/op\t   75249 allocs/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7549783,
            "unit": "ns/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594958,
            "unit": "B/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75249,
            "unit": "allocs/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8140,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "147787 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8140,
            "unit": "ns/op",
            "extra": "147787 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "147787 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "147787 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 873899,
            "unit": "ns/op\t  396578 B/op\t    6228 allocs/op",
            "extra": "1207 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 873899,
            "unit": "ns/op",
            "extra": "1207 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396578,
            "unit": "B/op",
            "extra": "1207 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6228,
            "unit": "allocs/op",
            "extra": "1207 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11172,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "107694 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11172,
            "unit": "ns/op",
            "extra": "107694 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "107694 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "107694 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8687636,
            "unit": "ns/op\t 4980891 B/op\t   81257 allocs/op",
            "extra": "134 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8687636,
            "unit": "ns/op",
            "extra": "134 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4980891,
            "unit": "B/op",
            "extra": "134 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81257,
            "unit": "allocs/op",
            "extra": "134 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 228.5,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5205571 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 228.5,
            "unit": "ns/op",
            "extra": "5205571 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5205571 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5205571 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6193,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6193,
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
            "value": 294.2,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4080684 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 294.2,
            "unit": "ns/op",
            "extra": "4080684 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4080684 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4080684 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 446.5,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2684512 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 446.5,
            "unit": "ns/op",
            "extra": "2684512 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2684512 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2684512 times\n4 procs"
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
          "id": "d9111427d0614904a304fd844135fdbd43423b53",
          "message": "chore(deps): update helm release kuma to v2.9.2 (#6803)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2024-12-10T12:05:57+01:00",
          "tree_id": "ef5267bb44a3ea768f4e3bbfa71a3c312a91870b",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/d9111427d0614904a304fd844135fdbd43423b53"
        },
        "date": 1733828970249,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1068368,
            "unit": "ns/op\t  819445 B/op\t       5 allocs/op",
            "extra": "1138 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1068368,
            "unit": "ns/op",
            "extra": "1138 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819445,
            "unit": "B/op",
            "extra": "1138 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1138 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 8965,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "113254 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 8965,
            "unit": "ns/op",
            "extra": "113254 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "113254 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "113254 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 78.97,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15122325 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 78.97,
            "unit": "ns/op",
            "extra": "15122325 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15122325 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15122325 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22479,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53494 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22479,
            "unit": "ns/op",
            "extra": "53494 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53494 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53494 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 220544,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4666 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 220544,
            "unit": "ns/op",
            "extra": "4666 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4666 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4666 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2606553,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "452 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2606553,
            "unit": "ns/op",
            "extra": "452 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "452 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "452 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 35418349,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "30 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 35418349,
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
            "value": 9673,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "123411 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9673,
            "unit": "ns/op",
            "extra": "123411 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "123411 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "123411 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7473803,
            "unit": "ns/op\t 4594529 B/op\t   75248 allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7473803,
            "unit": "ns/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594529,
            "unit": "B/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75248,
            "unit": "allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7994,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "146923 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7994,
            "unit": "ns/op",
            "extra": "146923 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "146923 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "146923 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 858791,
            "unit": "ns/op\t  396526 B/op\t    6227 allocs/op",
            "extra": "1220 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 858791,
            "unit": "ns/op",
            "extra": "1220 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396526,
            "unit": "B/op",
            "extra": "1220 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6227,
            "unit": "allocs/op",
            "extra": "1220 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11054,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "109192 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11054,
            "unit": "ns/op",
            "extra": "109192 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "109192 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "109192 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8708720,
            "unit": "ns/op\t 4980999 B/op\t   81257 allocs/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8708720,
            "unit": "ns/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4980999,
            "unit": "B/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81257,
            "unit": "allocs/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 226,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5284652 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 226,
            "unit": "ns/op",
            "extra": "5284652 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5284652 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5284652 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6188,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6188,
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
            "value": 293.8,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4086757 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 293.8,
            "unit": "ns/op",
            "extra": "4086757 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4086757 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4086757 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 443.7,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2657769 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 443.7,
            "unit": "ns/op",
            "extra": "2657769 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2657769 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2657769 times\n4 procs"
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
          "id": "5c93b46327642f2171b52fee314740b63ae90f17",
          "message": "chore(deps): bump github.com/docker/docker (#6805)\n\nBumps [github.com/docker/docker](https://github.com/docker/docker) from 27.3.1+incompatible to 27.4.0+incompatible.\r\n- [Release notes](https://github.com/docker/docker/releases)\r\n- [Commits](https://github.com/docker/docker/compare/v27.3.1...v27.4.0)\r\n\r\n---\r\nupdated-dependencies:\r\n- dependency-name: github.com/docker/docker\r\n  dependency-type: direct:production\r\n  update-type: version-update:semver-minor\r\n...\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2024-12-10T19:02:00Z",
          "tree_id": "f575d334501c2ee2508903915cadb02c64052298",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/5c93b46327642f2171b52fee314740b63ae90f17"
        },
        "date": 1733857507921,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1165427,
            "unit": "ns/op\t  819441 B/op\t       5 allocs/op",
            "extra": "1082 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1165427,
            "unit": "ns/op",
            "extra": "1082 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819441,
            "unit": "B/op",
            "extra": "1082 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1082 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7189,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "165153 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7189,
            "unit": "ns/op",
            "extra": "165153 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "165153 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "165153 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.51,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15215703 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.51,
            "unit": "ns/op",
            "extra": "15215703 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15215703 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15215703 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22744,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52928 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22744,
            "unit": "ns/op",
            "extra": "52928 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52928 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52928 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 251355,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4821 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 251355,
            "unit": "ns/op",
            "extra": "4821 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4821 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4821 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2588302,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "416 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2588302,
            "unit": "ns/op",
            "extra": "416 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
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
            "value": 38998282,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 38998282,
            "unit": "ns/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 9774,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "122324 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9774,
            "unit": "ns/op",
            "extra": "122324 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "122324 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "122324 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7443564,
            "unit": "ns/op\t 4594502 B/op\t   75248 allocs/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7443564,
            "unit": "ns/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594502,
            "unit": "B/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75248,
            "unit": "allocs/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8020,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "147868 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8020,
            "unit": "ns/op",
            "extra": "147868 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "147868 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "147868 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 870367,
            "unit": "ns/op\t  396510 B/op\t    6227 allocs/op",
            "extra": "1226 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 870367,
            "unit": "ns/op",
            "extra": "1226 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396510,
            "unit": "B/op",
            "extra": "1226 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6227,
            "unit": "allocs/op",
            "extra": "1226 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11089,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "107610 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11089,
            "unit": "ns/op",
            "extra": "107610 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "107610 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "107610 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8570333,
            "unit": "ns/op\t 4981287 B/op\t   81258 allocs/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8570333,
            "unit": "ns/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981287,
            "unit": "B/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81258,
            "unit": "allocs/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 230.3,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5281480 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 230.3,
            "unit": "ns/op",
            "extra": "5281480 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5281480 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5281480 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6197,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6197,
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
            "value": 298.2,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4038774 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 298.2,
            "unit": "ns/op",
            "extra": "4038774 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4038774 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4038774 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 452.5,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2661175 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 452.5,
            "unit": "ns/op",
            "extra": "2661175 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2661175 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2661175 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "571832+rainest@users.noreply.github.com",
            "name": "Travis Raines",
            "username": "rainest"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "4f6373f853da5d8a91c9dd78357136ec26f0f163",
          "message": "feat(diagnostics) add diff endpoints (#6131)\n\n* wip: diff diagnostic endpoint\r\n\r\n* chore: rename ConfigDump to ClientDiagnostic\r\n\r\nRename the diagnostics.ConfigDump struct to ClientDiagnostic to reflect\r\nthat it now handles multiple types of diagnostic data (config dumps,\r\nconfig diffs, and fallback status) from the Kong client.\r\n\r\n* chore: plumb diagnostic to DB mode strategies\r\n\r\n* feat: add diff collector and send to it\r\n\r\nAdd a struct to hold diff history in the diagnostic server. It holds a\r\nfixed number of diffs (5, arbitrarily) and provides functions to update\r\nand view history by config hash.\r\n\r\nAdd a diagnostics function to generate diagnostic diff structs from GDR\r\nevent data.\r\n\r\nAdd a diagnostics update loop case to handle inbound diffs.\r\n\r\nCall the diagnostic diff generator when processing an event and append\r\nthe diagnostic diff to the diff list to the DB mode event handler.\r\n\r\n* fix(tests) add missing CLI flag prefix\r\n\r\n* feat: add diff endpoint\r\n\r\n* chore: dump sensitive config during integration\r\n\r\nDump sensitive config during integration, mostly to support debugging\r\nthe diff endpoint. While there's no diff diagnostic integration test,\r\nhaving it available allows for easier interactive debugging of the\r\nendpoint, or viewing diff states when running DB-backed tests.\r\n\r\n* feat: wrap up diff endpoint\r\n\r\nAdd timestamps and a list of available diffs to the diff endpoint.\r\n\r\nAdd unit tests for the diff endpoint.\r\n\r\n* chore: clean up comments and lints\r\n\r\nRemove scaffolding comments and unfinished wishlist items. Add missing\r\ndoc comments.\r\n\r\nClean up some garbage the linter caught.\r\n\r\n* chore: add missing JSON key\r\n\r\n* chore: add changelog entry\r\n\r\n* chore: fix missed envtest type rename\r\n\r\n* fix lockups in UpdateStrategyDBMode.HandleEvent\r\n\r\n* tests: add tests for diagnostics/diff.go\r\n\r\n* Update internal/dataplane/sendconfig/dbmode.go\r\n\r\nCo-authored-by: Jakub Warczarek <jakub.warczarek@konghq.com>\r\n\r\n---------\r\n\r\nCo-authored-by: Yi Tao <tao.yi@konghq.com>\r\nCo-authored-by: Patryk Małek <patryk.malek@konghq.com>\r\nCo-authored-by: Jakub Warczarek <jakub.warczarek@konghq.com>",
          "timestamp": "2024-12-11T07:57:09+08:00",
          "tree_id": "45db63bd36097640be7de087e1429721db29b7f8",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/4f6373f853da5d8a91c9dd78357136ec26f0f163"
        },
        "date": 1733875186161,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1295724,
            "unit": "ns/op\t  819442 B/op\t       5 allocs/op",
            "extra": "841 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1295724,
            "unit": "ns/op",
            "extra": "841 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819442,
            "unit": "B/op",
            "extra": "841 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "841 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 8496,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "163633 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 8496,
            "unit": "ns/op",
            "extra": "163633 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "163633 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "163633 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.94,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15156915 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.94,
            "unit": "ns/op",
            "extra": "15156915 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15156915 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15156915 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22929,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53460 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22929,
            "unit": "ns/op",
            "extra": "53460 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53460 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53460 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 234232,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5265 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 234232,
            "unit": "ns/op",
            "extra": "5265 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5265 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5265 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2616033,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "462 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2616033,
            "unit": "ns/op",
            "extra": "462 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "462 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "462 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 32764883,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 32764883,
            "unit": "ns/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 9931,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "120919 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9931,
            "unit": "ns/op",
            "extra": "120919 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "120919 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "120919 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7924412,
            "unit": "ns/op\t 4594641 B/op\t   75248 allocs/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7924412,
            "unit": "ns/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594641,
            "unit": "B/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75248,
            "unit": "allocs/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8225,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "144411 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8225,
            "unit": "ns/op",
            "extra": "144411 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "144411 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "144411 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 878221,
            "unit": "ns/op\t  396638 B/op\t    6229 allocs/op",
            "extra": "1190 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 878221,
            "unit": "ns/op",
            "extra": "1190 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396638,
            "unit": "B/op",
            "extra": "1190 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6229,
            "unit": "allocs/op",
            "extra": "1190 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11513,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "102900 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11513,
            "unit": "ns/op",
            "extra": "102900 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "102900 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "102900 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8978615,
            "unit": "ns/op\t 4981008 B/op\t   81257 allocs/op",
            "extra": "133 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8978615,
            "unit": "ns/op",
            "extra": "133 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981008,
            "unit": "B/op",
            "extra": "133 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81257,
            "unit": "allocs/op",
            "extra": "133 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 228.5,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5238081 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 228.5,
            "unit": "ns/op",
            "extra": "5238081 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5238081 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5238081 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6199,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6199,
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
            "value": 296.7,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4048740 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 296.7,
            "unit": "ns/op",
            "extra": "4048740 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4048740 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4048740 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 525.6,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2653063 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 525.6,
            "unit": "ns/op",
            "extra": "2653063 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2653063 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2653063 times\n4 procs"
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
          "id": "79f97b9503f4f08a73c098798162d84ff7e640bf",
          "message": "feat(custom_entity): introduce schema validation (#6802)\n\n* feat(custom_entity): introduce schema validation\r\n\r\n* chore: make linter happy\r\n\r\n* tests: adjust tests\r\n\r\n* chore: rename schemaGetter -> schemaService",
          "timestamp": "2024-12-11T10:19:33+01:00",
          "tree_id": "5a4674a3449eb7124e1c38b9cd7a9f68854c5fe4",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/79f97b9503f4f08a73c098798162d84ff7e640bf"
        },
        "date": 1733908919981,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1266403,
            "unit": "ns/op\t  819442 B/op\t       5 allocs/op",
            "extra": "1008 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1266403,
            "unit": "ns/op",
            "extra": "1008 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819442,
            "unit": "B/op",
            "extra": "1008 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1008 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 9120,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "111662 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 9120,
            "unit": "ns/op",
            "extra": "111662 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "111662 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "111662 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.05,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15196514 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.05,
            "unit": "ns/op",
            "extra": "15196514 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15196514 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15196514 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22595,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53139 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22595,
            "unit": "ns/op",
            "extra": "53139 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53139 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53139 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 222992,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4906 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 222992,
            "unit": "ns/op",
            "extra": "4906 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4906 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4906 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2578252,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "456 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2578252,
            "unit": "ns/op",
            "extra": "456 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "456 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "456 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 32139350,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 32139350,
            "unit": "ns/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 9593,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "124968 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9593,
            "unit": "ns/op",
            "extra": "124968 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "124968 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "124968 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7515457,
            "unit": "ns/op\t 4594574 B/op\t   75248 allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7515457,
            "unit": "ns/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594574,
            "unit": "B/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75248,
            "unit": "allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7889,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "151645 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7889,
            "unit": "ns/op",
            "extra": "151645 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "151645 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "151645 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 853971,
            "unit": "ns/op\t  396376 B/op\t    6225 allocs/op",
            "extra": "1268 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 853971,
            "unit": "ns/op",
            "extra": "1268 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396376,
            "unit": "B/op",
            "extra": "1268 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6225,
            "unit": "allocs/op",
            "extra": "1268 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10943,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "111280 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10943,
            "unit": "ns/op",
            "extra": "111280 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "111280 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "111280 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8479508,
            "unit": "ns/op\t 4981209 B/op\t   81258 allocs/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8479508,
            "unit": "ns/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981209,
            "unit": "B/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81258,
            "unit": "allocs/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 229.9,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5309185 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 229.9,
            "unit": "ns/op",
            "extra": "5309185 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5309185 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5309185 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6217,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6217,
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
            "value": 359.7,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4043563 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 359.7,
            "unit": "ns/op",
            "extra": "4043563 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4043563 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4043563 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 447.2,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2676244 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 447.2,
            "unit": "ns/op",
            "extra": "2676244 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2676244 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2676244 times\n4 procs"
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
          "id": "63de0149aaf6e4a8314bf02896502d7f8d6f1cd0",
          "message": "feat(translator): combine services across httproutes expression router (#6766)\n\n* combine Kong services for HTTPRoutes when using expression routes\r\n\r\n* update tests and changelog",
          "timestamp": "2024-12-11T10:42:46+01:00",
          "tree_id": "d1d72408226f80ce3fa89d74149914780d0f832f",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/63de0149aaf6e4a8314bf02896502d7f8d6f1cd0"
        },
        "date": 1733910312849,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1162030,
            "unit": "ns/op\t  819441 B/op\t       5 allocs/op",
            "extra": "1006 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1162030,
            "unit": "ns/op",
            "extra": "1006 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819441,
            "unit": "B/op",
            "extra": "1006 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1006 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 8740,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "124777 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 8740,
            "unit": "ns/op",
            "extra": "124777 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "124777 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "124777 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 84.51,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15140793 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 84.51,
            "unit": "ns/op",
            "extra": "15140793 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15140793 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15140793 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22412,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53260 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22412,
            "unit": "ns/op",
            "extra": "53260 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53260 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53260 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 221245,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5280 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 221245,
            "unit": "ns/op",
            "extra": "5280 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5280 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5280 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2588790,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "424 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2588790,
            "unit": "ns/op",
            "extra": "424 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
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
            "value": 31065153,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 31065153,
            "unit": "ns/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 9738,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "125481 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9738,
            "unit": "ns/op",
            "extra": "125481 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "125481 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "125481 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7350988,
            "unit": "ns/op\t 4594527 B/op\t   75248 allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7350988,
            "unit": "ns/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594527,
            "unit": "B/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75248,
            "unit": "allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8057,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "146264 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8057,
            "unit": "ns/op",
            "extra": "146264 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "146264 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "146264 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 867386,
            "unit": "ns/op\t  396436 B/op\t    6226 allocs/op",
            "extra": "1249 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 867386,
            "unit": "ns/op",
            "extra": "1249 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396436,
            "unit": "B/op",
            "extra": "1249 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6226,
            "unit": "allocs/op",
            "extra": "1249 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11021,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "108883 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11021,
            "unit": "ns/op",
            "extra": "108883 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "108883 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "108883 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8662885,
            "unit": "ns/op\t 4981234 B/op\t   81258 allocs/op",
            "extra": "135 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8662885,
            "unit": "ns/op",
            "extra": "135 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981234,
            "unit": "B/op",
            "extra": "135 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81258,
            "unit": "allocs/op",
            "extra": "135 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 230.4,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5240973 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 230.4,
            "unit": "ns/op",
            "extra": "5240973 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5240973 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5240973 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6244,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6244,
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
            "value": 298.4,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4029788 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 298.4,
            "unit": "ns/op",
            "extra": "4029788 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4029788 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4029788 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 446.6,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2672124 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 446.6,
            "unit": "ns/op",
            "extra": "2672124 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2672124 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2672124 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "lavacca.mattia@gmail.com",
            "name": "Mattia Lavacca",
            "username": "mlavacca"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "4d1771041f5698cf0b16b54d379bd9cc7fc19868",
          "message": "feat: convert `BackendTLSPolicies` into service annotations (#6753)\n\n* feat: convert BackendTLSPolicies into service anns\r\n\r\nAll the BackendTLSPolicies are converted into a set of annotations that\r\nare already supported by KIC.\r\n\r\n* chore: CHANGELOG entry added\r\n\r\n* address review comments\r\n\r\n* Apply suggestions from code review\r\n\r\n* Update internal/controllers/reference/reference.go\r\n\r\n* address reviews comments\r\n\r\n* add secret-label-selector\r\n\r\n* refactor getCACerts\r\n\r\n* refactor UpdateReferencesToSecret and removeOutdatedReferencesToSecret\r\n\r\n* add default for --configmap-label-selector\r\n\r\n* Apply suggestions from code review\r\n\r\n---------\r\n\r\nSigned-off-by: Mattia Lavacca <lavacca.mattia@gmail.com>\r\nCo-authored-by: Jakub Warczarek <jakub.warczarek@konghq.com>\r\nCo-authored-by: Patryk Małek <patryk.malek@konghq.com>",
          "timestamp": "2024-12-11T11:25:06+01:00",
          "tree_id": "78bd167f89ba853a2d3bc8c4fe05ad6a3ea80e32",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/4d1771041f5698cf0b16b54d379bd9cc7fc19868"
        },
        "date": 1733912850802,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1179521,
            "unit": "ns/op\t  819442 B/op\t       5 allocs/op",
            "extra": "957 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1179521,
            "unit": "ns/op",
            "extra": "957 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819442,
            "unit": "B/op",
            "extra": "957 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "957 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7229,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "167233 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7229,
            "unit": "ns/op",
            "extra": "167233 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "167233 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "167233 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15047498 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.1,
            "unit": "ns/op",
            "extra": "15047498 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15047498 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15047498 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22185,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53944 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22185,
            "unit": "ns/op",
            "extra": "53944 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53944 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53944 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 221242,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5560 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 221242,
            "unit": "ns/op",
            "extra": "5560 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5560 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5560 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2544632,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "463 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2544632,
            "unit": "ns/op",
            "extra": "463 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "463 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "463 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 33319678,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 33319678,
            "unit": "ns/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 9997,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "119085 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9997,
            "unit": "ns/op",
            "extra": "119085 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "119085 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "119085 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7438057,
            "unit": "ns/op\t 4595117 B/op\t   75255 allocs/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7438057,
            "unit": "ns/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4595117,
            "unit": "B/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75255,
            "unit": "allocs/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8323,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "142743 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8323,
            "unit": "ns/op",
            "extra": "142743 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "142743 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "142743 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 887772,
            "unit": "ns/op\t  396724 B/op\t    6232 allocs/op",
            "extra": "1234 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 887772,
            "unit": "ns/op",
            "extra": "1234 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396724,
            "unit": "B/op",
            "extra": "1234 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6232,
            "unit": "allocs/op",
            "extra": "1234 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11331,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "105679 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11331,
            "unit": "ns/op",
            "extra": "105679 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "105679 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "105679 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8503971,
            "unit": "ns/op\t 4980971 B/op\t   81262 allocs/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8503971,
            "unit": "ns/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4980971,
            "unit": "B/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81262,
            "unit": "allocs/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 229.3,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5273456 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 229.3,
            "unit": "ns/op",
            "extra": "5273456 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5273456 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5273456 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6194,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6194,
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
            "value": 294.8,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4045495 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 294.8,
            "unit": "ns/op",
            "extra": "4045495 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4045495 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4045495 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 475.6,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2713351 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 475.6,
            "unit": "ns/op",
            "extra": "2713351 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2713351 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2713351 times\n4 procs"
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
          "id": "6d50108f520437abd4e6427ea1aa9fb94329e1cf",
          "message": "tests(kongintegration): retry when creating kong's container (#6808)",
          "timestamp": "2024-12-11T12:40:11+01:00",
          "tree_id": "138d133c63e47d0666dcb889c3d32cf1b37c4376",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/6d50108f520437abd4e6427ea1aa9fb94329e1cf"
        },
        "date": 1733917359531,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1232569,
            "unit": "ns/op\t  819442 B/op\t       5 allocs/op",
            "extra": "991 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1232569,
            "unit": "ns/op",
            "extra": "991 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819442,
            "unit": "B/op",
            "extra": "991 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "991 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 9072,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "119500 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 9072,
            "unit": "ns/op",
            "extra": "119500 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "119500 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "119500 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 80.89,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14860485 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 80.89,
            "unit": "ns/op",
            "extra": "14860485 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14860485 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14860485 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22724,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53192 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22724,
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
            "value": 227602,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4740 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 227602,
            "unit": "ns/op",
            "extra": "4740 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4740 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4740 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2530844,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "442 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2530844,
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
            "value": 31777539,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "75 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 31777539,
            "unit": "ns/op",
            "extra": "75 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 10292,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "117334 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10292,
            "unit": "ns/op",
            "extra": "117334 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "117334 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "117334 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7774395,
            "unit": "ns/op\t 4594475 B/op\t   75253 allocs/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7774395,
            "unit": "ns/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594475,
            "unit": "B/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75253,
            "unit": "allocs/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8445,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "139890 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8445,
            "unit": "ns/op",
            "extra": "139890 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "139890 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "139890 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 870375,
            "unit": "ns/op\t  396757 B/op\t    6233 allocs/op",
            "extra": "1224 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 870375,
            "unit": "ns/op",
            "extra": "1224 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396757,
            "unit": "B/op",
            "extra": "1224 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6233,
            "unit": "allocs/op",
            "extra": "1224 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11502,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "102849 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11502,
            "unit": "ns/op",
            "extra": "102849 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "102849 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "102849 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8840359,
            "unit": "ns/op\t 4981196 B/op\t   81263 allocs/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8840359,
            "unit": "ns/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981196,
            "unit": "B/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81263,
            "unit": "allocs/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 229.1,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5260126 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 229.1,
            "unit": "ns/op",
            "extra": "5260126 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5260126 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5260126 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6208,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6208,
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
            "value": 294.2,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4057128 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 294.2,
            "unit": "ns/op",
            "extra": "4057128 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4057128 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4057128 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 450.3,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2676130 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 450.3,
            "unit": "ns/op",
            "extra": "2676130 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2676130 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2676130 times\n4 procs"
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
          "id": "93856b1bb353de1664f1aee0f21a36b41523179a",
          "message": "chore(deps): update dependency kubernetes/code-generator to v0.31.4 (#6806)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2024-12-11T12:40:24+01:00",
          "tree_id": "9ff35de8f8722447c36ead2cb15ee84b9d1741d8",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/93856b1bb353de1664f1aee0f21a36b41523179a"
        },
        "date": 1733917378552,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1185797,
            "unit": "ns/op\t  819442 B/op\t       5 allocs/op",
            "extra": "996 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1185797,
            "unit": "ns/op",
            "extra": "996 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819442,
            "unit": "B/op",
            "extra": "996 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "996 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 6973,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "171277 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 6973,
            "unit": "ns/op",
            "extra": "171277 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "171277 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "171277 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 78.88,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15053511 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 78.88,
            "unit": "ns/op",
            "extra": "15053511 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15053511 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15053511 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 23006,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "50726 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 23006,
            "unit": "ns/op",
            "extra": "50726 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "50726 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "50726 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 240573,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4482 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 240573,
            "unit": "ns/op",
            "extra": "4482 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4482 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4482 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2639258,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "423 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2639258,
            "unit": "ns/op",
            "extra": "423 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
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
            "value": 34575130,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 34575130,
            "unit": "ns/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 10175,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "118977 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10175,
            "unit": "ns/op",
            "extra": "118977 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "118977 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "118977 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7878014,
            "unit": "ns/op\t 4594867 B/op\t   75254 allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7878014,
            "unit": "ns/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594867,
            "unit": "B/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75254,
            "unit": "allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8456,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "142095 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8456,
            "unit": "ns/op",
            "extra": "142095 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "142095 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "142095 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 888535,
            "unit": "ns/op\t  396953 B/op\t    6236 allocs/op",
            "extra": "1164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 888535,
            "unit": "ns/op",
            "extra": "1164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396953,
            "unit": "B/op",
            "extra": "1164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6236,
            "unit": "allocs/op",
            "extra": "1164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11687,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "102061 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11687,
            "unit": "ns/op",
            "extra": "102061 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "102061 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "102061 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 9364698,
            "unit": "ns/op\t 4981210 B/op\t   81263 allocs/op",
            "extra": "129 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 9364698,
            "unit": "ns/op",
            "extra": "129 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981210,
            "unit": "B/op",
            "extra": "129 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81263,
            "unit": "allocs/op",
            "extra": "129 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 229.5,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5185796 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 229.5,
            "unit": "ns/op",
            "extra": "5185796 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5185796 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5185796 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6489,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6489,
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
            "value": 330,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "3193851 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 330,
            "unit": "ns/op",
            "extra": "3193851 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "3193851 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "3193851 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 452.4,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2639718 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 452.4,
            "unit": "ns/op",
            "extra": "2639718 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2639718 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2639718 times\n4 procs"
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
          "id": "33e0f8f8b524224acee90993b780caf76fb7e10a",
          "message": "chore(deps): bump the k8s-io group with 7 updates (#6810)\n\nBumps the k8s-io group with 7 updates:\r\n\r\n| Package | From | To |\r\n| --- | --- | --- |\r\n| [k8s.io/api](https://github.com/kubernetes/api) | `0.31.3` | `0.31.4` |\r\n| [k8s.io/apiextensions-apiserver](https://github.com/kubernetes/apiextensions-apiserver) | `0.31.3` | `0.31.4` |\r\n| [k8s.io/apimachinery](https://github.com/kubernetes/apimachinery) | `0.31.3` | `0.31.4` |\r\n| [k8s.io/client-go](https://github.com/kubernetes/client-go) | `0.31.3` | `0.31.4` |\r\n| [k8s.io/component-base](https://github.com/kubernetes/component-base) | `0.31.3` | `0.31.4` |\r\n| [k8s.io/cli-runtime](https://github.com/kubernetes/cli-runtime) | `0.31.3` | `0.31.4` |\r\n| [k8s.io/kubectl](https://github.com/kubernetes/kubectl) | `0.31.3` | `0.31.4` |\r\n\r\n\r\nUpdates `k8s.io/api` from 0.31.3 to 0.31.4\r\n- [Commits](https://github.com/kubernetes/api/compare/v0.31.3...v0.31.4)\r\n\r\nUpdates `k8s.io/apiextensions-apiserver` from 0.31.3 to 0.31.4\r\n- [Release notes](https://github.com/kubernetes/apiextensions-apiserver/releases)\r\n- [Commits](https://github.com/kubernetes/apiextensions-apiserver/compare/v0.31.3...v0.31.4)\r\n\r\nUpdates `k8s.io/apimachinery` from 0.31.3 to 0.31.4\r\n- [Commits](https://github.com/kubernetes/apimachinery/compare/v0.31.3...v0.31.4)\r\n\r\nUpdates `k8s.io/client-go` from 0.31.3 to 0.31.4\r\n- [Changelog](https://github.com/kubernetes/client-go/blob/master/CHANGELOG.md)\r\n- [Commits](https://github.com/kubernetes/client-go/compare/v0.31.3...v0.31.4)\r\n\r\nUpdates `k8s.io/component-base` from 0.31.3 to 0.31.4\r\n- [Commits](https://github.com/kubernetes/component-base/compare/v0.31.3...v0.31.4)\r\n\r\nUpdates `k8s.io/cli-runtime` from 0.31.3 to 0.31.4\r\n- [Commits](https://github.com/kubernetes/cli-runtime/compare/v0.31.3...v0.31.4)\r\n\r\nUpdates `k8s.io/kubectl` from 0.31.3 to 0.31.4\r\n- [Commits](https://github.com/kubernetes/kubectl/compare/v0.31.3...v0.31.4)\r\n\r\n---\r\nupdated-dependencies:\r\n- dependency-name: k8s.io/api\r\n  dependency-type: direct:production\r\n  update-type: version-update:semver-patch\r\n  dependency-group: k8s-io\r\n- dependency-name: k8s.io/apiextensions-apiserver\r\n  dependency-type: direct:production\r\n  update-type: version-update:semver-patch\r\n  dependency-group: k8s-io\r\n- dependency-name: k8s.io/apimachinery\r\n  dependency-type: direct:production\r\n  update-type: version-update:semver-patch\r\n  dependency-group: k8s-io\r\n- dependency-name: k8s.io/client-go\r\n  dependency-type: direct:production\r\n  update-type: version-update:semver-patch\r\n  dependency-group: k8s-io\r\n- dependency-name: k8s.io/component-base\r\n  dependency-type: direct:production\r\n  update-type: version-update:semver-patch\r\n  dependency-group: k8s-io\r\n- dependency-name: k8s.io/cli-runtime\r\n  dependency-type: direct:production\r\n  update-type: version-update:semver-patch\r\n  dependency-group: k8s-io\r\n- dependency-name: k8s.io/kubectl\r\n  dependency-type: direct:production\r\n  update-type: version-update:semver-patch\r\n  dependency-group: k8s-io\r\n...\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2024-12-11T15:54:00+01:00",
          "tree_id": "1c025ae88e064e4d7b62fec8327bea4ef912d21e",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/33e0f8f8b524224acee90993b780caf76fb7e10a"
        },
        "date": 1733929035922,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1776225,
            "unit": "ns/op\t  819441 B/op\t       5 allocs/op",
            "extra": "769 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1776225,
            "unit": "ns/op",
            "extra": "769 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819441,
            "unit": "B/op",
            "extra": "769 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "769 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 6953,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "147472 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 6953,
            "unit": "ns/op",
            "extra": "147472 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "147472 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "147472 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.29,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15072190 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.29,
            "unit": "ns/op",
            "extra": "15072190 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15072190 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15072190 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 23203,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52434 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 23203,
            "unit": "ns/op",
            "extra": "52434 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52434 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52434 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 237833,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4441 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 237833,
            "unit": "ns/op",
            "extra": "4441 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4441 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4441 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2674205,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "448 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2674205,
            "unit": "ns/op",
            "extra": "448 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "448 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "448 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 33931116,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 33931116,
            "unit": "ns/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 10328,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "114213 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10328,
            "unit": "ns/op",
            "extra": "114213 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "114213 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "114213 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7966943,
            "unit": "ns/op\t 4595004 B/op\t   75254 allocs/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7966943,
            "unit": "ns/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4595004,
            "unit": "B/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75254,
            "unit": "allocs/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8574,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "138462 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8574,
            "unit": "ns/op",
            "extra": "138462 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "138462 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "138462 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 892756,
            "unit": "ns/op\t  396743 B/op\t    6233 allocs/op",
            "extra": "1228 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 892756,
            "unit": "ns/op",
            "extra": "1228 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396743,
            "unit": "B/op",
            "extra": "1228 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6233,
            "unit": "allocs/op",
            "extra": "1228 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 12123,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "97624 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 12123,
            "unit": "ns/op",
            "extra": "97624 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "97624 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "97624 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 9593192,
            "unit": "ns/op\t 4981670 B/op\t   81265 allocs/op",
            "extra": "130 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 9593192,
            "unit": "ns/op",
            "extra": "130 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981670,
            "unit": "B/op",
            "extra": "130 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81265,
            "unit": "allocs/op",
            "extra": "130 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 234.9,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5194628 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 234.9,
            "unit": "ns/op",
            "extra": "5194628 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5194628 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5194628 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6573,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6573,
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
            "value": 298.5,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "3444804 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 298.5,
            "unit": "ns/op",
            "extra": "3444804 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "3444804 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "3444804 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 454.2,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2628952 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 454.2,
            "unit": "ns/op",
            "extra": "2628952 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2628952 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2628952 times\n4 procs"
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
          "id": "f3e7618fa2e2cc2d965278a442c1b60bdf0c2d51",
          "message": "chore(deps): bump google.golang.org/api from 0.210.0 to 0.211.0 (#6811)\n\nBumps [google.golang.org/api](https://github.com/googleapis/google-api-go-client) from 0.210.0 to 0.211.0.\r\n- [Release notes](https://github.com/googleapis/google-api-go-client/releases)\r\n- [Changelog](https://github.com/googleapis/google-api-go-client/blob/main/CHANGES.md)\r\n- [Commits](https://github.com/googleapis/google-api-go-client/compare/v0.210.0...v0.211.0)\r\n\r\n---\r\nupdated-dependencies:\r\n- dependency-name: google.golang.org/api\r\n  dependency-type: direct:production\r\n  update-type: version-update:semver-minor\r\n...\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2024-12-11T15:15:34Z",
          "tree_id": "caf12429ca03f7b85b8363824cc4715481de2ff4",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/f3e7618fa2e2cc2d965278a442c1b60bdf0c2d51"
        },
        "date": 1733930319828,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1047451,
            "unit": "ns/op\t  819442 B/op\t       5 allocs/op",
            "extra": "980 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1047451,
            "unit": "ns/op",
            "extra": "980 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819442,
            "unit": "B/op",
            "extra": "980 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "980 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 8311,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "165970 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 8311,
            "unit": "ns/op",
            "extra": "165970 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "165970 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "165970 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 80.01,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12815504 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 80.01,
            "unit": "ns/op",
            "extra": "12815504 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12815504 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12815504 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22473,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "51729 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22473,
            "unit": "ns/op",
            "extra": "51729 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "51729 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "51729 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 222109,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4700 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 222109,
            "unit": "ns/op",
            "extra": "4700 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4700 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4700 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2669950,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "406 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2669950,
            "unit": "ns/op",
            "extra": "406 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "406 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "406 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 36544350,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "75 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 36544350,
            "unit": "ns/op",
            "extra": "75 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 10263,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "120696 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10263,
            "unit": "ns/op",
            "extra": "120696 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "120696 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "120696 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7299112,
            "unit": "ns/op\t 4594833 B/op\t   75254 allocs/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7299112,
            "unit": "ns/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594833,
            "unit": "B/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75254,
            "unit": "allocs/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8243,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "147060 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8243,
            "unit": "ns/op",
            "extra": "147060 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "147060 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "147060 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 853613,
            "unit": "ns/op\t  396697 B/op\t    6232 allocs/op",
            "extra": "1234 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 853613,
            "unit": "ns/op",
            "extra": "1234 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396697,
            "unit": "B/op",
            "extra": "1234 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6232,
            "unit": "allocs/op",
            "extra": "1234 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11232,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "105865 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11232,
            "unit": "ns/op",
            "extra": "105865 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "105865 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "105865 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8451022,
            "unit": "ns/op\t 4981909 B/op\t   81266 allocs/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8451022,
            "unit": "ns/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981909,
            "unit": "B/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81266,
            "unit": "allocs/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 228.3,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5271499 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 228.3,
            "unit": "ns/op",
            "extra": "5271499 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5271499 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5271499 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6193,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6193,
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
            "value": 294.3,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4110666 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 294.3,
            "unit": "ns/op",
            "extra": "4110666 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4110666 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4110666 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 450.4,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2681743 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 450.4,
            "unit": "ns/op",
            "extra": "2681743 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2681743 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2681743 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "lavacca.mattia@gmail.com",
            "name": "Mattia Lavacca",
            "username": "mlavacca"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "9f0ded8121426d8b924d025fa24d4f33fb1b57f7",
          "message": "test: `BackendTLSPolicy` integration test (#6799)\n\n* feat: convert BackendTLSPolicies into service anns\r\n\r\nAll the BackendTLSPolicies are converted into a set of annotations that\r\nare already supported by KIC.\r\n\r\nSigned-off-by: Mattia Lavacca <lavacca.mattia@gmail.com>\r\n\r\n* chore: CHANGELOG entry added\r\n\r\nSigned-off-by: Mattia Lavacca <lavacca.mattia@gmail.com>\r\n\r\n* feat: convert BackendTLSPolicies into service anns\r\n\r\nAll the BackendTLSPolicies are converted into a set of annotations that\r\nare already supported by KIC.\r\n\r\nSigned-off-by: Mattia Lavacca <lavacca.mattia@gmail.com>\r\n\r\n* test: integration test for BackendTLSPolicies\r\n\r\nSigned-off-by: Mattia Lavacca <lavacca.mattia@gmail.com>\r\n\r\n* chore: deep copy k8s service in translator\r\n\r\nSigned-off-by: Mattia Lavacca <lavacca.mattia@gmail.com>\r\n\r\n* chore: remove unused code\r\n\r\n* chore(tests): fix using invalid tls-verify-depth key\r\n\r\n* chore: extract TLSVerifyDepthKey to a named const\r\n\r\n* chore: remove unnecessary /status RBAC policy rules\r\n\r\n---------\r\n\r\nSigned-off-by: Mattia Lavacca <lavacca.mattia@gmail.com>\r\nCo-authored-by: Patryk Małek <patryk.malek@konghq.com>",
          "timestamp": "2024-12-11T16:05:17Z",
          "tree_id": "f1b34218480cdbcc67742e5b2e0c9c9ffde86016",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/9f0ded8121426d8b924d025fa24d4f33fb1b57f7"
        },
        "date": 1733933305500,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1165488,
            "unit": "ns/op\t  819442 B/op\t       5 allocs/op",
            "extra": "1202 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1165488,
            "unit": "ns/op",
            "extra": "1202 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819442,
            "unit": "B/op",
            "extra": "1202 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1202 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7238,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "140440 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7238,
            "unit": "ns/op",
            "extra": "140440 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "140440 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "140440 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.03,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14634019 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.03,
            "unit": "ns/op",
            "extra": "14634019 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14634019 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14634019 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22896,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53452 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22896,
            "unit": "ns/op",
            "extra": "53452 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53452 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53452 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 273310,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5722 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 273310,
            "unit": "ns/op",
            "extra": "5722 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5722 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5722 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2579631,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "448 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2579631,
            "unit": "ns/op",
            "extra": "448 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "448 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "448 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 30954948,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 30954948,
            "unit": "ns/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 9893,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "119982 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9893,
            "unit": "ns/op",
            "extra": "119982 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "119982 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "119982 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7413309,
            "unit": "ns/op\t 4594717 B/op\t   75253 allocs/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7413309,
            "unit": "ns/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594717,
            "unit": "B/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75253,
            "unit": "allocs/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8262,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "142964 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8262,
            "unit": "ns/op",
            "extra": "142964 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "142964 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "142964 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 874083,
            "unit": "ns/op\t  396714 B/op\t    6232 allocs/op",
            "extra": "1239 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 874083,
            "unit": "ns/op",
            "extra": "1239 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396714,
            "unit": "B/op",
            "extra": "1239 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6232,
            "unit": "allocs/op",
            "extra": "1239 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11270,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "105016 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11270,
            "unit": "ns/op",
            "extra": "105016 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "105016 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "105016 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8670880,
            "unit": "ns/op\t 4981870 B/op\t   81265 allocs/op",
            "extra": "136 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8670880,
            "unit": "ns/op",
            "extra": "136 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981870,
            "unit": "B/op",
            "extra": "136 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81265,
            "unit": "allocs/op",
            "extra": "136 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 230.3,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5232985 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 230.3,
            "unit": "ns/op",
            "extra": "5232985 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5232985 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5232985 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6204,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6204,
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
            "value": 296.6,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4018754 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 296.6,
            "unit": "ns/op",
            "extra": "4018754 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4018754 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4018754 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 449.6,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2652102 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 449.6,
            "unit": "ns/op",
            "extra": "2652102 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2652102 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2652102 times\n4 procs"
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
          "id": "488c1c9938ec6e78127dd438117476bde9772815",
          "message": "chore(deps): bump github.com/kong/go-database-reconciler (#6812)\n\nBumps [github.com/kong/go-database-reconciler](https://github.com/kong/go-database-reconciler) from 1.16.1 to 1.17.0.\r\n- [Release notes](https://github.com/kong/go-database-reconciler/releases)\r\n- [Commits](https://github.com/kong/go-database-reconciler/compare/v1.16.1...v1.17.0)\r\n\r\n---\r\nupdated-dependencies:\r\n- dependency-name: github.com/kong/go-database-reconciler\r\n  dependency-type: direct:production\r\n  update-type: version-update:semver-minor\r\n...\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2024-12-11T17:01:12Z",
          "tree_id": "9fe5d880d6a281400d7643cde94c6d58a67b358c",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/488c1c9938ec6e78127dd438117476bde9772815"
        },
        "date": 1733936671595,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1297238,
            "unit": "ns/op\t  819441 B/op\t       5 allocs/op",
            "extra": "994 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1297238,
            "unit": "ns/op",
            "extra": "994 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819441,
            "unit": "B/op",
            "extra": "994 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "994 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 9111,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "112674 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 9111,
            "unit": "ns/op",
            "extra": "112674 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "112674 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "112674 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.98,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15137959 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.98,
            "unit": "ns/op",
            "extra": "15137959 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15137959 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15137959 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 23085,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "51462 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 23085,
            "unit": "ns/op",
            "extra": "51462 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "51462 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "51462 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 238006,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4666 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 238006,
            "unit": "ns/op",
            "extra": "4666 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4666 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4666 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2570328,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "432 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2570328,
            "unit": "ns/op",
            "extra": "432 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
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
            "value": 36413629,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 36413629,
            "unit": "ns/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 10434,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "117182 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10434,
            "unit": "ns/op",
            "extra": "117182 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "117182 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "117182 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7845999,
            "unit": "ns/op\t 4594811 B/op\t   75254 allocs/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7845999,
            "unit": "ns/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594811,
            "unit": "B/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75254,
            "unit": "allocs/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8432,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "137630 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8432,
            "unit": "ns/op",
            "extra": "137630 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "137630 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "137630 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 904503,
            "unit": "ns/op\t  396787 B/op\t    6234 allocs/op",
            "extra": "1212 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 904503,
            "unit": "ns/op",
            "extra": "1212 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396787,
            "unit": "B/op",
            "extra": "1212 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6234,
            "unit": "allocs/op",
            "extra": "1212 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11581,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "103010 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11581,
            "unit": "ns/op",
            "extra": "103010 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "103010 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "103010 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 9297539,
            "unit": "ns/op\t 4981478 B/op\t   81264 allocs/op",
            "extra": "129 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 9297539,
            "unit": "ns/op",
            "extra": "129 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981478,
            "unit": "B/op",
            "extra": "129 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81264,
            "unit": "allocs/op",
            "extra": "129 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 275.8,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5128304 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 275.8,
            "unit": "ns/op",
            "extra": "5128304 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5128304 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5128304 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.631,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.631,
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
            "value": 297.8,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4057440 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 297.8,
            "unit": "ns/op",
            "extra": "4057440 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4057440 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4057440 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 459.8,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2630516 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 459.8,
            "unit": "ns/op",
            "extra": "2630516 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2630516 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2630516 times\n4 procs"
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
          "id": "c68ea01e360a4c0ae98ec1ed1f6d8a6c462f99b6",
          "message": "chore(deps): bump google.golang.org/api from 0.209.0 to 0.210.0 (#6784)\n\nBumps [google.golang.org/api](https://github.com/googleapis/google-api-go-client) from 0.209.0 to 0.210.0.\r\n- [Release notes](https://github.com/googleapis/google-api-go-client/releases)\r\n- [Changelog](https://github.com/googleapis/google-api-go-client/blob/main/CHANGES.md)\r\n- [Commits](https://github.com/googleapis/google-api-go-client/compare/v0.209.0...v0.210.0)\r\n\r\n---\r\nupdated-dependencies:\r\n- dependency-name: google.golang.org/api\r\n  dependency-type: direct:production\r\n  update-type: version-update:semver-minor\r\n...\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2024-12-05T15:12:34Z",
          "tree_id": "c8a91a0c746ad0d6a375a4ab5922f9fde429cd86",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/c68ea01e360a4c0ae98ec1ed1f6d8a6c462f99b6"
        },
        "date": 1733991235729,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1331970,
            "unit": "ns/op\t  819442 B/op\t       5 allocs/op",
            "extra": "910 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1331970,
            "unit": "ns/op",
            "extra": "910 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819442,
            "unit": "B/op",
            "extra": "910 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "910 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 8689,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "157539 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 8689,
            "unit": "ns/op",
            "extra": "157539 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "157539 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "157539 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 81.25,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14241560 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 81.25,
            "unit": "ns/op",
            "extra": "14241560 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14241560 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14241560 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22601,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52221 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22601,
            "unit": "ns/op",
            "extra": "52221 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52221 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52221 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 224205,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5461 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 224205,
            "unit": "ns/op",
            "extra": "5461 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5461 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5461 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2645829,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "445 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2645829,
            "unit": "ns/op",
            "extra": "445 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "445 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "445 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 35588439,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 35588439,
            "unit": "ns/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 9733,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "123232 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9733,
            "unit": "ns/op",
            "extra": "123232 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "123232 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "123232 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7544898,
            "unit": "ns/op\t 4594514 B/op\t   75248 allocs/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7544898,
            "unit": "ns/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594514,
            "unit": "B/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75248,
            "unit": "allocs/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8110,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "146283 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8110,
            "unit": "ns/op",
            "extra": "146283 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "146283 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "146283 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 883405,
            "unit": "ns/op\t  396526 B/op\t    6227 allocs/op",
            "extra": "1221 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 883405,
            "unit": "ns/op",
            "extra": "1221 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396526,
            "unit": "B/op",
            "extra": "1221 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6227,
            "unit": "allocs/op",
            "extra": "1221 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11349,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "107589 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11349,
            "unit": "ns/op",
            "extra": "107589 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "107589 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "107589 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8851382,
            "unit": "ns/op\t 4981052 B/op\t   81257 allocs/op",
            "extra": "136 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8851382,
            "unit": "ns/op",
            "extra": "136 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981052,
            "unit": "B/op",
            "extra": "136 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81257,
            "unit": "allocs/op",
            "extra": "136 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 240.8,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5254743 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 240.8,
            "unit": "ns/op",
            "extra": "5254743 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5254743 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5254743 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.63,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.63,
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
            "value": 296.8,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "3800536 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 296.8,
            "unit": "ns/op",
            "extra": "3800536 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "3800536 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "3800536 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 451.2,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2654971 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 451.2,
            "unit": "ns/op",
            "extra": "2654971 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2654971 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2654971 times\n4 procs"
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
          "id": "81574837b35fb5deadabe6e5859fd07310d1ecdd",
          "message": "chore(deps): bump golang.org/x/crypto from 0.30.0 to 0.31.0 (#6820)\n\nBumps [golang.org/x/crypto](https://github.com/golang/crypto) from 0.30.0 to 0.31.0.\r\n- [Commits](https://github.com/golang/crypto/compare/v0.30.0...v0.31.0)\r\n\r\n---\r\nupdated-dependencies:\r\n- dependency-name: golang.org/x/crypto\r\n  dependency-type: indirect\r\n...\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2024-12-12T11:58:17+01:00",
          "tree_id": "c227e35c208a53d956270304fab52809dfcc5e92",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/81574837b35fb5deadabe6e5859fd07310d1ecdd"
        },
        "date": 1734001292101,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1139713,
            "unit": "ns/op\t  819447 B/op\t       5 allocs/op",
            "extra": "968 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1139713,
            "unit": "ns/op",
            "extra": "968 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819447,
            "unit": "B/op",
            "extra": "968 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "968 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 6803,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "166198 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 6803,
            "unit": "ns/op",
            "extra": "166198 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "166198 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "166198 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 78.81,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15241983 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 78.81,
            "unit": "ns/op",
            "extra": "15241983 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15241983 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15241983 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22456,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53146 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22456,
            "unit": "ns/op",
            "extra": "53146 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53146 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53146 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 224612,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5076 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 224612,
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
            "value": 2538686,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "465 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2538686,
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
            "value": 38270475,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 38270475,
            "unit": "ns/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 9703,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "122719 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9703,
            "unit": "ns/op",
            "extra": "122719 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "122719 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "122719 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7348802,
            "unit": "ns/op\t 4594623 B/op\t   75253 allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7348802,
            "unit": "ns/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594623,
            "unit": "B/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75253,
            "unit": "allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8153,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "143542 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8153,
            "unit": "ns/op",
            "extra": "143542 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "143542 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "143542 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 852574,
            "unit": "ns/op\t  396682 B/op\t    6232 allocs/op",
            "extra": "1245 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 852574,
            "unit": "ns/op",
            "extra": "1245 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396682,
            "unit": "B/op",
            "extra": "1245 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6232,
            "unit": "allocs/op",
            "extra": "1245 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11107,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "105614 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11107,
            "unit": "ns/op",
            "extra": "105614 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "105614 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "105614 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8390860,
            "unit": "ns/op\t 4981385 B/op\t   81264 allocs/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8390860,
            "unit": "ns/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981385,
            "unit": "B/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81264,
            "unit": "allocs/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 240.1,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5051894 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 240.1,
            "unit": "ns/op",
            "extra": "5051894 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5051894 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5051894 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6377,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6377,
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
            "value": 301.4,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "3551864 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 301.4,
            "unit": "ns/op",
            "extra": "3551864 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "3551864 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "3551864 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 450.5,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2673975 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 450.5,
            "unit": "ns/op",
            "extra": "2673975 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2673975 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2673975 times\n4 procs"
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
          "id": "9247df3f03a42d134ca2d4585a596bf2105dfe9a",
          "message": "revert GDR to 1.16.1 and go-kong to 0.60.0 (#6822)\n\nCo-authored-by: Yi Tao <tao.yi@konghq.com>",
          "timestamp": "2024-12-12T13:32:32+01:00",
          "tree_id": "7f4d507f8d3db2c0820eb7b79857ef3cd79105be",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/9247df3f03a42d134ca2d4585a596bf2105dfe9a"
        },
        "date": 1734006941217,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1307250,
            "unit": "ns/op\t  819449 B/op\t       5 allocs/op",
            "extra": "787 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1307250,
            "unit": "ns/op",
            "extra": "787 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819449,
            "unit": "B/op",
            "extra": "787 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "787 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 18563,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "55662 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 18563,
            "unit": "ns/op",
            "extra": "55662 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "55662 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "55662 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 88.73,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11353833 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 88.73,
            "unit": "ns/op",
            "extra": "11353833 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11353833 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11353833 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22927,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52890 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22927,
            "unit": "ns/op",
            "extra": "52890 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52890 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52890 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 239498,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5301 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 239498,
            "unit": "ns/op",
            "extra": "5301 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5301 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5301 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2647947,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "433 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2647947,
            "unit": "ns/op",
            "extra": "433 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "433 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "433 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 34978517,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 34978517,
            "unit": "ns/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 10450,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "112726 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10450,
            "unit": "ns/op",
            "extra": "112726 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "112726 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "112726 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7886632,
            "unit": "ns/op\t 4595050 B/op\t   75255 allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7886632,
            "unit": "ns/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4595050,
            "unit": "B/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75255,
            "unit": "allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8779,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "139166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8779,
            "unit": "ns/op",
            "extra": "139166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "139166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "139166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 902744,
            "unit": "ns/op\t  396899 B/op\t    6235 allocs/op",
            "extra": "1184 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 902744,
            "unit": "ns/op",
            "extra": "1184 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396899,
            "unit": "B/op",
            "extra": "1184 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6235,
            "unit": "allocs/op",
            "extra": "1184 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11786,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "102217 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11786,
            "unit": "ns/op",
            "extra": "102217 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "102217 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "102217 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 9258178,
            "unit": "ns/op\t 4981518 B/op\t   81264 allocs/op",
            "extra": "132 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 9258178,
            "unit": "ns/op",
            "extra": "132 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981518,
            "unit": "B/op",
            "extra": "132 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81264,
            "unit": "allocs/op",
            "extra": "132 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 230.9,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5115807 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 230.9,
            "unit": "ns/op",
            "extra": "5115807 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5115807 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5115807 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6207,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6207,
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
            "value": 364.5,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "3634981 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 364.5,
            "unit": "ns/op",
            "extra": "3634981 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "3634981 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "3634981 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 458.7,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2611429 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 458.7,
            "unit": "ns/op",
            "extra": "2611429 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2611429 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2611429 times\n4 procs"
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
          "id": "5fac68875bd9d1286dcc1185679f14e71d01b8d8",
          "message": "chore(deps): update dependency kubernetes/code-generator to v0.32.0 (#6818)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2024-12-12T12:58:03Z",
          "tree_id": "76049319e8fe67caf4ca18348d16a6a36a73c400",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/5fac68875bd9d1286dcc1185679f14e71d01b8d8"
        },
        "date": 1734008471855,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1331717,
            "unit": "ns/op\t  819446 B/op\t       5 allocs/op",
            "extra": "1268 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1331717,
            "unit": "ns/op",
            "extra": "1268 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819446,
            "unit": "B/op",
            "extra": "1268 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1268 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7171,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "146583 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7171,
            "unit": "ns/op",
            "extra": "146583 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "146583 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "146583 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.48,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15149308 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.48,
            "unit": "ns/op",
            "extra": "15149308 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15149308 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15149308 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22493,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52546 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22493,
            "unit": "ns/op",
            "extra": "52546 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52546 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52546 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 224305,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5130 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 224305,
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
            "value": 2569938,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "466 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2569938,
            "unit": "ns/op",
            "extra": "466 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
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
            "value": 30612651,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 30612651,
            "unit": "ns/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 9707,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "124946 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9707,
            "unit": "ns/op",
            "extra": "124946 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "124946 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "124946 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7354027,
            "unit": "ns/op\t 4594520 B/op\t   75253 allocs/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7354027,
            "unit": "ns/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594520,
            "unit": "B/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75253,
            "unit": "allocs/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8127,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "147794 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8127,
            "unit": "ns/op",
            "extra": "147794 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "147794 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "147794 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 868136,
            "unit": "ns/op\t  396661 B/op\t    6231 allocs/op",
            "extra": "1254 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 868136,
            "unit": "ns/op",
            "extra": "1254 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396661,
            "unit": "B/op",
            "extra": "1254 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6231,
            "unit": "allocs/op",
            "extra": "1254 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11127,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "106869 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11127,
            "unit": "ns/op",
            "extra": "106869 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "106869 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "106869 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8455671,
            "unit": "ns/op\t 4981206 B/op\t   81263 allocs/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8455671,
            "unit": "ns/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981206,
            "unit": "B/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81263,
            "unit": "allocs/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 229.2,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5201364 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 229.2,
            "unit": "ns/op",
            "extra": "5201364 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5201364 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5201364 times\n4 procs"
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
            "value": 330.8,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "3244161 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 330.8,
            "unit": "ns/op",
            "extra": "3244161 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "3244161 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "3244161 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 446,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2711824 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 446,
            "unit": "ns/op",
            "extra": "2711824 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2711824 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2711824 times\n4 procs"
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
          "id": "59759c81e80b5fbde1f2232ba8ef52f4ca6c2560",
          "message": "chore(deps): update kindest/node docker tag to v1.32.0 (#6819)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2024-12-12T15:18:45+01:00",
          "tree_id": "3e5cfa58c1f85a991793190d8547fa6f313cefb9",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/59759c81e80b5fbde1f2232ba8ef52f4ca6c2560"
        },
        "date": 1734013328625,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1329787,
            "unit": "ns/op\t  819442 B/op\t       5 allocs/op",
            "extra": "981 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1329787,
            "unit": "ns/op",
            "extra": "981 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819442,
            "unit": "B/op",
            "extra": "981 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "981 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7481,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "163381 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7481,
            "unit": "ns/op",
            "extra": "163381 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "163381 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "163381 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.27,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15105888 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.27,
            "unit": "ns/op",
            "extra": "15105888 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15105888 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15105888 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 26026,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "49488 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 26026,
            "unit": "ns/op",
            "extra": "49488 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "49488 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "49488 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 271168,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4161 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 271168,
            "unit": "ns/op",
            "extra": "4161 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4161 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4161 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2714926,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "436 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2714926,
            "unit": "ns/op",
            "extra": "436 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "436 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "436 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 36177885,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 36177885,
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
            "value": 10458,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "113162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10458,
            "unit": "ns/op",
            "extra": "113162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "113162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "113162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 8337301,
            "unit": "ns/op\t 4594845 B/op\t   75254 allocs/op",
            "extra": "136 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 8337301,
            "unit": "ns/op",
            "extra": "136 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594845,
            "unit": "B/op",
            "extra": "136 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75254,
            "unit": "allocs/op",
            "extra": "136 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8931,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "135780 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8931,
            "unit": "ns/op",
            "extra": "135780 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "135780 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "135780 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 936119,
            "unit": "ns/op\t  396944 B/op\t    6236 allocs/op",
            "extra": "1167 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 936119,
            "unit": "ns/op",
            "extra": "1167 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396944,
            "unit": "B/op",
            "extra": "1167 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6236,
            "unit": "allocs/op",
            "extra": "1167 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 12619,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "97476 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 12619,
            "unit": "ns/op",
            "extra": "97476 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "97476 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "97476 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 10290828,
            "unit": "ns/op\t 4981144 B/op\t   81263 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 10290828,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981144,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81263,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 238.1,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "4907938 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 238.1,
            "unit": "ns/op",
            "extra": "4907938 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "4907938 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4907938 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6221,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6221,
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
            "value": 303.6,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "3921228 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 303.6,
            "unit": "ns/op",
            "extra": "3921228 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "3921228 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "3921228 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 557.9,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2490196 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 557.9,
            "unit": "ns/op",
            "extra": "2490196 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2490196 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2490196 times\n4 procs"
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
          "id": "06f120792c1fec982b2027ffb8978d46bf25c185",
          "message": "chore(deps): bump github.com/goccy/go-json from 0.10.3 to 0.10.4 (#6825)\n\nBumps [github.com/goccy/go-json](https://github.com/goccy/go-json) from 0.10.3 to 0.10.4.\r\n- [Release notes](https://github.com/goccy/go-json/releases)\r\n- [Changelog](https://github.com/goccy/go-json/blob/master/CHANGELOG.md)\r\n- [Commits](https://github.com/goccy/go-json/compare/v0.10.3...v0.10.4)\r\n\r\n---\r\nupdated-dependencies:\r\n- dependency-name: github.com/goccy/go-json\r\n  dependency-type: direct:production\r\n  update-type: version-update:semver-patch\r\n...\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2024-12-12T15:58:30+01:00",
          "tree_id": "118a59640dd1fc561ced456946a1e1884f39d46c",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/06f120792c1fec982b2027ffb8978d46bf25c185"
        },
        "date": 1734015696081,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1227172,
            "unit": "ns/op\t  819441 B/op\t       5 allocs/op",
            "extra": "931 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1227172,
            "unit": "ns/op",
            "extra": "931 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819441,
            "unit": "B/op",
            "extra": "931 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "931 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 10133,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "113107 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 10133,
            "unit": "ns/op",
            "extra": "113107 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "113107 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "113107 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15041661 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.9,
            "unit": "ns/op",
            "extra": "15041661 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15041661 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15041661 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22763,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53254 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22763,
            "unit": "ns/op",
            "extra": "53254 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53254 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53254 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 235576,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4839 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 235576,
            "unit": "ns/op",
            "extra": "4839 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4839 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4839 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2583542,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "459 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2583542,
            "unit": "ns/op",
            "extra": "459 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "459 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "459 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 32434929,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 32434929,
            "unit": "ns/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 9835,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "119586 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9835,
            "unit": "ns/op",
            "extra": "119586 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "119586 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "119586 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7418590,
            "unit": "ns/op\t 4595124 B/op\t   75255 allocs/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7418590,
            "unit": "ns/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4595124,
            "unit": "B/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75255,
            "unit": "allocs/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8298,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "145347 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8298,
            "unit": "ns/op",
            "extra": "145347 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "145347 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "145347 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 860064,
            "unit": "ns/op\t  396695 B/op\t    6232 allocs/op",
            "extra": "1240 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 860064,
            "unit": "ns/op",
            "extra": "1240 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396695,
            "unit": "B/op",
            "extra": "1240 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6232,
            "unit": "allocs/op",
            "extra": "1240 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11661,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "104185 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11661,
            "unit": "ns/op",
            "extra": "104185 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "104185 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "104185 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 9060427,
            "unit": "ns/op\t 4981195 B/op\t   81263 allocs/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 9060427,
            "unit": "ns/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981195,
            "unit": "B/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81263,
            "unit": "allocs/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 228.3,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5168052 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 228.3,
            "unit": "ns/op",
            "extra": "5168052 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5168052 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5168052 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.619,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.619,
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
            "value": 297.1,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4039045 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 297.1,
            "unit": "ns/op",
            "extra": "4039045 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4039045 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4039045 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 478,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2696830 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 478,
            "unit": "ns/op",
            "extra": "2696830 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2696830 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2696830 times\n4 procs"
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
          "id": "537590971f66036bfd7f9d0c7191a3d6350c73f8",
          "message": "chore(examples): add BackendTLSPolicy example manifest (#6813)",
          "timestamp": "2024-12-12T16:00:55+01:00",
          "tree_id": "34605aafb22713f2f5e3aa1340eab0ef7fe3af42",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/537590971f66036bfd7f9d0c7191a3d6350c73f8"
        },
        "date": 1734015843372,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1272752,
            "unit": "ns/op\t  819441 B/op\t       5 allocs/op",
            "extra": "889 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1272752,
            "unit": "ns/op",
            "extra": "889 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819441,
            "unit": "B/op",
            "extra": "889 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "889 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 18568,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "56478 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 18568,
            "unit": "ns/op",
            "extra": "56478 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "56478 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "56478 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.32,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13964728 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.32,
            "unit": "ns/op",
            "extra": "13964728 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13964728 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13964728 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22644,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53176 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22644,
            "unit": "ns/op",
            "extra": "53176 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53176 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53176 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 280039,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5203 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 280039,
            "unit": "ns/op",
            "extra": "5203 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5203 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5203 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2571987,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "427 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2571987,
            "unit": "ns/op",
            "extra": "427 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
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
            "value": 39386716,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 39386716,
            "unit": "ns/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 9937,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "118430 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9937,
            "unit": "ns/op",
            "extra": "118430 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "118430 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "118430 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7575402,
            "unit": "ns/op\t 4594774 B/op\t   75254 allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7575402,
            "unit": "ns/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594774,
            "unit": "B/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75254,
            "unit": "allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8489,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "143464 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8489,
            "unit": "ns/op",
            "extra": "143464 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "143464 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "143464 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 881328,
            "unit": "ns/op\t  396806 B/op\t    6234 allocs/op",
            "extra": "1207 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 881328,
            "unit": "ns/op",
            "extra": "1207 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396806,
            "unit": "B/op",
            "extra": "1207 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6234,
            "unit": "allocs/op",
            "extra": "1207 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11461,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "102888 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11461,
            "unit": "ns/op",
            "extra": "102888 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "102888 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "102888 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8900931,
            "unit": "ns/op\t 4981686 B/op\t   81265 allocs/op",
            "extra": "136 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8900931,
            "unit": "ns/op",
            "extra": "136 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981686,
            "unit": "B/op",
            "extra": "136 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81265,
            "unit": "allocs/op",
            "extra": "136 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 229.7,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5206748 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 229.7,
            "unit": "ns/op",
            "extra": "5206748 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5206748 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5206748 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.62,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.62,
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
            "value": 295,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4028067 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 295,
            "unit": "ns/op",
            "extra": "4028067 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4028067 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4028067 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 447.6,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2697304 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 447.6,
            "unit": "ns/op",
            "extra": "2697304 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2697304 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2697304 times\n4 procs"
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
          "id": "7bfa78b20747d84bb0e2ce477601f2e61ac22b1d",
          "message": "chore(deps): bump the k8s-io group with 7 updates (#6824)\n\nBumps the k8s-io group with 7 updates:\r\n\r\n| Package | From | To |\r\n| --- | --- | --- |\r\n| [k8s.io/api](https://github.com/kubernetes/api) | `0.31.4` | `0.32.0` |\r\n| [k8s.io/apiextensions-apiserver](https://github.com/kubernetes/apiextensions-apiserver) | `0.31.4` | `0.32.0` |\r\n| [k8s.io/apimachinery](https://github.com/kubernetes/apimachinery) | `0.31.4` | `0.32.0` |\r\n| [k8s.io/client-go](https://github.com/kubernetes/client-go) | `0.31.4` | `0.32.0` |\r\n| [k8s.io/component-base](https://github.com/kubernetes/component-base) | `0.31.4` | `0.32.0` |\r\n| [k8s.io/cli-runtime](https://github.com/kubernetes/cli-runtime) | `0.31.4` | `0.32.0` |\r\n| [k8s.io/kubectl](https://github.com/kubernetes/kubectl) | `0.31.4` | `0.32.0` |\r\n\r\n\r\nUpdates `k8s.io/api` from 0.31.4 to 0.32.0\r\n- [Commits](https://github.com/kubernetes/api/compare/v0.31.4...v0.32.0)\r\n\r\nUpdates `k8s.io/apiextensions-apiserver` from 0.31.4 to 0.32.0\r\n- [Release notes](https://github.com/kubernetes/apiextensions-apiserver/releases)\r\n- [Commits](https://github.com/kubernetes/apiextensions-apiserver/compare/v0.31.4...v0.32.0)\r\n\r\nUpdates `k8s.io/apimachinery` from 0.31.4 to 0.32.0\r\n- [Commits](https://github.com/kubernetes/apimachinery/compare/v0.31.4...v0.32.0)\r\n\r\nUpdates `k8s.io/client-go` from 0.31.4 to 0.32.0\r\n- [Changelog](https://github.com/kubernetes/client-go/blob/master/CHANGELOG.md)\r\n- [Commits](https://github.com/kubernetes/client-go/compare/v0.31.4...v0.32.0)\r\n\r\nUpdates `k8s.io/component-base` from 0.31.4 to 0.32.0\r\n- [Commits](https://github.com/kubernetes/component-base/compare/v0.31.4...v0.32.0)\r\n\r\nUpdates `k8s.io/cli-runtime` from 0.31.4 to 0.32.0\r\n- [Commits](https://github.com/kubernetes/cli-runtime/compare/v0.31.4...v0.32.0)\r\n\r\nUpdates `k8s.io/kubectl` from 0.31.4 to 0.32.0\r\n- [Commits](https://github.com/kubernetes/kubectl/compare/v0.31.4...v0.32.0)\r\n\r\n---\r\nupdated-dependencies:\r\n- dependency-name: k8s.io/api\r\n  dependency-type: direct:production\r\n  update-type: version-update:semver-minor\r\n  dependency-group: k8s-io\r\n- dependency-name: k8s.io/apiextensions-apiserver\r\n  dependency-type: direct:production\r\n  update-type: version-update:semver-minor\r\n  dependency-group: k8s-io\r\n- dependency-name: k8s.io/apimachinery\r\n  dependency-type: direct:production\r\n  update-type: version-update:semver-minor\r\n  dependency-group: k8s-io\r\n- dependency-name: k8s.io/client-go\r\n  dependency-type: direct:production\r\n  update-type: version-update:semver-minor\r\n  dependency-group: k8s-io\r\n- dependency-name: k8s.io/component-base\r\n  dependency-type: direct:production\r\n  update-type: version-update:semver-minor\r\n  dependency-group: k8s-io\r\n- dependency-name: k8s.io/cli-runtime\r\n  dependency-type: direct:production\r\n  update-type: version-update:semver-minor\r\n  dependency-group: k8s-io\r\n- dependency-name: k8s.io/kubectl\r\n  dependency-type: direct:production\r\n  update-type: version-update:semver-minor\r\n  dependency-group: k8s-io\r\n...\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2024-12-12T17:08:06+01:00",
          "tree_id": "0ab6f27a6068ac286adcfa5d0b09fcdde40d0ce0",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/7bfa78b20747d84bb0e2ce477601f2e61ac22b1d"
        },
        "date": 1734019891535,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1170678,
            "unit": "ns/op\t  819440 B/op\t       5 allocs/op",
            "extra": "1024 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1170678,
            "unit": "ns/op",
            "extra": "1024 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819440,
            "unit": "B/op",
            "extra": "1024 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1024 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 8606,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "162130 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 8606,
            "unit": "ns/op",
            "extra": "162130 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "162130 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "162130 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.37,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15043028 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.37,
            "unit": "ns/op",
            "extra": "15043028 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15043028 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15043028 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 23199,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "51452 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 23199,
            "unit": "ns/op",
            "extra": "51452 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "51452 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "51452 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 243322,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4476 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 243322,
            "unit": "ns/op",
            "extra": "4476 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4476 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4476 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2646810,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "427 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2646810,
            "unit": "ns/op",
            "extra": "427 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
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
            "value": 34670736,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 34670736,
            "unit": "ns/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 10206,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "118491 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10206,
            "unit": "ns/op",
            "extra": "118491 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "118491 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "118491 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 8365419,
            "unit": "ns/op\t 4594992 B/op\t   75254 allocs/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 8365419,
            "unit": "ns/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594992,
            "unit": "B/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75254,
            "unit": "allocs/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8557,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "138096 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8557,
            "unit": "ns/op",
            "extra": "138096 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "138096 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "138096 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 908074,
            "unit": "ns/op\t  396894 B/op\t    6235 allocs/op",
            "extra": "1186 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 908074,
            "unit": "ns/op",
            "extra": "1186 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396894,
            "unit": "B/op",
            "extra": "1186 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6235,
            "unit": "allocs/op",
            "extra": "1186 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 12129,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "97834 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 12129,
            "unit": "ns/op",
            "extra": "97834 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "97834 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "97834 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 9412762,
            "unit": "ns/op\t 4981494 B/op\t   81264 allocs/op",
            "extra": "126 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 9412762,
            "unit": "ns/op",
            "extra": "126 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981494,
            "unit": "B/op",
            "extra": "126 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81264,
            "unit": "allocs/op",
            "extra": "126 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 232.2,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5140596 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 232.2,
            "unit": "ns/op",
            "extra": "5140596 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5140596 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5140596 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6202,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6202,
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
            "value": 300,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4044790 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 300,
            "unit": "ns/op",
            "extra": "4044790 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4044790 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4044790 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 574.9,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2644879 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 574.9,
            "unit": "ns/op",
            "extra": "2644879 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2644879 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2644879 times\n4 procs"
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
          "id": "d6f15bc1fea1a0fd91f3ba2b53f92fdb3e441edb",
          "message": "chore(deps): update kong/kong-gateway docker tag to v3.9.0.0 (#6830)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2024-12-13T09:46:26+01:00",
          "tree_id": "13996d042b3ce8485865a93238fac000b130b9ce",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/d6f15bc1fea1a0fd91f3ba2b53f92fdb3e441edb"
        },
        "date": 1734079783722,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1276489,
            "unit": "ns/op\t  819441 B/op\t       5 allocs/op",
            "extra": "958 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1276489,
            "unit": "ns/op",
            "extra": "958 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819441,
            "unit": "B/op",
            "extra": "958 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "958 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7274,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "166892 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7274,
            "unit": "ns/op",
            "extra": "166892 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "166892 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "166892 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.65,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15011764 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.65,
            "unit": "ns/op",
            "extra": "15011764 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15011764 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15011764 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22579,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52915 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22579,
            "unit": "ns/op",
            "extra": "52915 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52915 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52915 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 224833,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4812 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 224833,
            "unit": "ns/op",
            "extra": "4812 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4812 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4812 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2562230,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "463 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2562230,
            "unit": "ns/op",
            "extra": "463 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "463 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "463 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 41440186,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 41440186,
            "unit": "ns/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 10222,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "118326 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10222,
            "unit": "ns/op",
            "extra": "118326 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "118326 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "118326 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7479837,
            "unit": "ns/op\t 4594810 B/op\t   75254 allocs/op",
            "extra": "156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7479837,
            "unit": "ns/op",
            "extra": "156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594810,
            "unit": "B/op",
            "extra": "156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75254,
            "unit": "allocs/op",
            "extra": "156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8440,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "144430 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8440,
            "unit": "ns/op",
            "extra": "144430 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "144430 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "144430 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 918777,
            "unit": "ns/op\t  396786 B/op\t    6234 allocs/op",
            "extra": "1212 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 918777,
            "unit": "ns/op",
            "extra": "1212 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396786,
            "unit": "B/op",
            "extra": "1212 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6234,
            "unit": "allocs/op",
            "extra": "1212 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11737,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "96072 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11737,
            "unit": "ns/op",
            "extra": "96072 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "96072 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "96072 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8937095,
            "unit": "ns/op\t 4981615 B/op\t   81265 allocs/op",
            "extra": "132 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8937095,
            "unit": "ns/op",
            "extra": "132 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981615,
            "unit": "B/op",
            "extra": "132 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81265,
            "unit": "allocs/op",
            "extra": "132 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 234.4,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5162510 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 234.4,
            "unit": "ns/op",
            "extra": "5162510 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5162510 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5162510 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6205,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6205,
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
            "value": 298.7,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4023843 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 298.7,
            "unit": "ns/op",
            "extra": "4023843 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4023843 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4023843 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 451.3,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2665285 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 451.3,
            "unit": "ns/op",
            "extra": "2665285 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2665285 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2665285 times\n4 procs"
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
          "id": "84f4fa0f5b639aedee132d0a30d809983d78e40d",
          "message": "chore(deps): update kong/kong-gateway@regenerate docker tag to v3.9 (#6831)\n\n* chore(deps): update kong/kong-gateway@regenerate docker tag to v3.9\n\n* chore: regenerate\n\n---------\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>\nCo-authored-by: github-actions <github-actions@users.noreply.github.com>",
          "timestamp": "2024-12-13T10:08:32+01:00",
          "tree_id": "2ce03bafe5dd00fe9dba37a42d32521288bac591",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/84f4fa0f5b639aedee132d0a30d809983d78e40d"
        },
        "date": 1734081106084,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1268344,
            "unit": "ns/op\t  819441 B/op\t       5 allocs/op",
            "extra": "1076 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1268344,
            "unit": "ns/op",
            "extra": "1076 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819441,
            "unit": "B/op",
            "extra": "1076 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1076 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 15181,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "66334 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 15181,
            "unit": "ns/op",
            "extra": "66334 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "66334 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "66334 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 78.83,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13919998 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 78.83,
            "unit": "ns/op",
            "extra": "13919998 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13919998 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13919998 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22500,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53738 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22500,
            "unit": "ns/op",
            "extra": "53738 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53738 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53738 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 224155,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5334 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 224155,
            "unit": "ns/op",
            "extra": "5334 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
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
            "value": 2501309,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "406 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2501309,
            "unit": "ns/op",
            "extra": "406 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "406 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "406 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 31824197,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 31824197,
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
            "value": 10083,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "120456 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10083,
            "unit": "ns/op",
            "extra": "120456 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "120456 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "120456 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7645910,
            "unit": "ns/op\t 4594796 B/op\t   75254 allocs/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7645910,
            "unit": "ns/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594796,
            "unit": "B/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75254,
            "unit": "allocs/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8419,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "143932 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8419,
            "unit": "ns/op",
            "extra": "143932 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "143932 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "143932 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 881274,
            "unit": "ns/op\t  396813 B/op\t    6234 allocs/op",
            "extra": "1204 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 881274,
            "unit": "ns/op",
            "extra": "1204 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396813,
            "unit": "B/op",
            "extra": "1204 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6234,
            "unit": "allocs/op",
            "extra": "1204 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11558,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "103114 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11558,
            "unit": "ns/op",
            "extra": "103114 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "103114 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "103114 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8521592,
            "unit": "ns/op\t 4981449 B/op\t   81264 allocs/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8521592,
            "unit": "ns/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981449,
            "unit": "B/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81264,
            "unit": "allocs/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 238.2,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5184646 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 238.2,
            "unit": "ns/op",
            "extra": "5184646 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5184646 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5184646 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6195,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6195,
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
            "value": 297.9,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4035518 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 297.9,
            "unit": "ns/op",
            "extra": "4035518 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4035518 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4035518 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 451.1,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2594722 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 451.1,
            "unit": "ns/op",
            "extra": "2594722 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2594722 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2594722 times\n4 procs"
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
          "id": "d881a383ec2eba6b98bbe4ce5f8a1217292c8991",
          "message": "chore(deps): update helm release kong to v2.44.1 (#6829)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2024-12-13T09:40:55Z",
          "tree_id": "9af5b57d4991f2a4d9dff5d1a73397f54a777d0a",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/d881a383ec2eba6b98bbe4ce5f8a1217292c8991"
        },
        "date": 1734083055940,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1449811,
            "unit": "ns/op\t  819442 B/op\t       5 allocs/op",
            "extra": "1270 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1449811,
            "unit": "ns/op",
            "extra": "1270 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819442,
            "unit": "B/op",
            "extra": "1270 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1270 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 6825,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "175732 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 6825,
            "unit": "ns/op",
            "extra": "175732 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "175732 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "175732 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.39,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15162104 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.39,
            "unit": "ns/op",
            "extra": "15162104 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15162104 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15162104 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22806,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52281 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22806,
            "unit": "ns/op",
            "extra": "52281 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52281 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52281 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 230925,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4551 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 230925,
            "unit": "ns/op",
            "extra": "4551 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4551 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4551 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2736841,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "448 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2736841,
            "unit": "ns/op",
            "extra": "448 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "448 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "448 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 45381798,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 45381798,
            "unit": "ns/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 10186,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "118268 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10186,
            "unit": "ns/op",
            "extra": "118268 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "118268 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "118268 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 8005180,
            "unit": "ns/op\t 4594655 B/op\t   75253 allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 8005180,
            "unit": "ns/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594655,
            "unit": "B/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75253,
            "unit": "allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8434,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "138210 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8434,
            "unit": "ns/op",
            "extra": "138210 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "138210 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "138210 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 891856,
            "unit": "ns/op\t  396779 B/op\t    6233 allocs/op",
            "extra": "1219 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 891856,
            "unit": "ns/op",
            "extra": "1219 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396779,
            "unit": "B/op",
            "extra": "1219 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6233,
            "unit": "allocs/op",
            "extra": "1219 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11774,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "102472 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11774,
            "unit": "ns/op",
            "extra": "102472 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "102472 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "102472 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 9188439,
            "unit": "ns/op\t 4981521 B/op\t   81264 allocs/op",
            "extra": "130 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 9188439,
            "unit": "ns/op",
            "extra": "130 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981521,
            "unit": "B/op",
            "extra": "130 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81264,
            "unit": "allocs/op",
            "extra": "130 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 259.7,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5156712 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 259.7,
            "unit": "ns/op",
            "extra": "5156712 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5156712 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5156712 times\n4 procs"
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
            "value": 298.8,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "3688434 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 298.8,
            "unit": "ns/op",
            "extra": "3688434 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "3688434 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "3688434 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 451.9,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2617196 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 451.9,
            "unit": "ns/op",
            "extra": "2617196 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2617196 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2617196 times\n4 procs"
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
          "id": "d85a6a04e0575401d7767b451e0c09587520a908",
          "message": "chore: predicate.And -> predicate.Or in secret and configmap controllers (#6828)",
          "timestamp": "2024-12-13T13:23:33+01:00",
          "tree_id": "7e1c9d0f2a542bf4cc29baf2aacbec54796b43b5",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/d85a6a04e0575401d7767b451e0c09587520a908"
        },
        "date": 1734092818881,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1189230,
            "unit": "ns/op\t  819447 B/op\t       5 allocs/op",
            "extra": "1015 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1189230,
            "unit": "ns/op",
            "extra": "1015 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819447,
            "unit": "B/op",
            "extra": "1015 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1015 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 9923,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "102742 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 9923,
            "unit": "ns/op",
            "extra": "102742 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "102742 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "102742 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.04,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15123627 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.04,
            "unit": "ns/op",
            "extra": "15123627 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15123627 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15123627 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 23341,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "50343 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 23341,
            "unit": "ns/op",
            "extra": "50343 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "50343 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "50343 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 238474,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4448 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 238474,
            "unit": "ns/op",
            "extra": "4448 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4448 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4448 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2571712,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "452 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2571712,
            "unit": "ns/op",
            "extra": "452 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "452 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "452 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 43893301,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "68 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 43893301,
            "unit": "ns/op",
            "extra": "68 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 10381,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "117258 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10381,
            "unit": "ns/op",
            "extra": "117258 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "117258 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "117258 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7628220,
            "unit": "ns/op\t 4594607 B/op\t   75253 allocs/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7628220,
            "unit": "ns/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594607,
            "unit": "B/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75253,
            "unit": "allocs/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8399,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "140168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8399,
            "unit": "ns/op",
            "extra": "140168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "140168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "140168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 885619,
            "unit": "ns/op\t  396801 B/op\t    6234 allocs/op",
            "extra": "1210 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 885619,
            "unit": "ns/op",
            "extra": "1210 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396801,
            "unit": "B/op",
            "extra": "1210 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6234,
            "unit": "allocs/op",
            "extra": "1210 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 12023,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "88250 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 12023,
            "unit": "ns/op",
            "extra": "88250 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "88250 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "88250 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 9008930,
            "unit": "ns/op\t 4981418 B/op\t   81264 allocs/op",
            "extra": "133 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 9008930,
            "unit": "ns/op",
            "extra": "133 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981418,
            "unit": "B/op",
            "extra": "133 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81264,
            "unit": "allocs/op",
            "extra": "133 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 276.8,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5085627 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 276.8,
            "unit": "ns/op",
            "extra": "5085627 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5085627 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5085627 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6376,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6376,
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
            "value": 300.7,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4007497 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 300.7,
            "unit": "ns/op",
            "extra": "4007497 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4007497 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4007497 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 459.6,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2600970 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 459.6,
            "unit": "ns/op",
            "extra": "2600970 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2600970 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2600970 times\n4 procs"
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
          "id": "aeb6c54ec92317bd8830a2dbe405c2e5ee7543a0",
          "message": "feat(KongPlugin): generate K8s Events for misconfiguration (#6814)",
          "timestamp": "2024-12-13T13:24:50+01:00",
          "tree_id": "9480030f79e6b4f6a991df63471395273e962bbb",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/aeb6c54ec92317bd8830a2dbe405c2e5ee7543a0"
        },
        "date": 1734092888816,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1197149,
            "unit": "ns/op\t  819442 B/op\t       5 allocs/op",
            "extra": "943 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1197149,
            "unit": "ns/op",
            "extra": "943 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819442,
            "unit": "B/op",
            "extra": "943 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "943 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 9846,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "104688 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 9846,
            "unit": "ns/op",
            "extra": "104688 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "104688 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "104688 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 80.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15150013 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 80.4,
            "unit": "ns/op",
            "extra": "15150013 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15150013 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15150013 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22509,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52831 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22509,
            "unit": "ns/op",
            "extra": "52831 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52831 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52831 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 222905,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5076 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 222905,
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
            "value": 3023599,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "441 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 3023599,
            "unit": "ns/op",
            "extra": "441 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "441 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "441 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 38210638,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 38210638,
            "unit": "ns/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 9754,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "120951 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9754,
            "unit": "ns/op",
            "extra": "120951 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "120951 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "120951 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7363908,
            "unit": "ns/op\t 4594849 B/op\t   75254 allocs/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7363908,
            "unit": "ns/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594849,
            "unit": "B/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75254,
            "unit": "allocs/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8162,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "146449 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8162,
            "unit": "ns/op",
            "extra": "146449 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "146449 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "146449 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 858113,
            "unit": "ns/op\t  396850 B/op\t    6234 allocs/op",
            "extra": "1198 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 858113,
            "unit": "ns/op",
            "extra": "1198 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396850,
            "unit": "B/op",
            "extra": "1198 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6234,
            "unit": "allocs/op",
            "extra": "1198 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11456,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "106660 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11456,
            "unit": "ns/op",
            "extra": "106660 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "106660 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "106660 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8646519,
            "unit": "ns/op\t 4981214 B/op\t   81263 allocs/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8646519,
            "unit": "ns/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981214,
            "unit": "B/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81263,
            "unit": "allocs/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 232.6,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5196302 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 232.6,
            "unit": "ns/op",
            "extra": "5196302 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5196302 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5196302 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6198,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6198,
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
            "value": 513.8,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2316318 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 513.8,
            "unit": "ns/op",
            "extra": "2316318 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2316318 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2316318 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 744.3,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1612310 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 744.3,
            "unit": "ns/op",
            "extra": "1612310 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1612310 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1612310 times\n4 procs"
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
          "id": "cd2336345ff9d1d8b755b3769bd0b4751f54b051",
          "message": "chore(deps): bump google.golang.org/grpc from 1.68.1 to 1.69.0 (#6835)\n\nBumps [google.golang.org/grpc](https://github.com/grpc/grpc-go) from 1.68.1 to 1.69.0.\r\n- [Release notes](https://github.com/grpc/grpc-go/releases)\r\n- [Commits](https://github.com/grpc/grpc-go/compare/v1.68.1...v1.69.0)\r\n\r\n---\r\nupdated-dependencies:\r\n- dependency-name: google.golang.org/grpc\r\n  dependency-type: direct:production\r\n  update-type: version-update:semver-minor\r\n...\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2024-12-13T15:03:03Z",
          "tree_id": "fae7aa5bb6eaf86c105ae509a695e9464c874ff9",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/cd2336345ff9d1d8b755b3769bd0b4751f54b051"
        },
        "date": 1734102386586,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1326712,
            "unit": "ns/op\t  819449 B/op\t       5 allocs/op",
            "extra": "759 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1326712,
            "unit": "ns/op",
            "extra": "759 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819449,
            "unit": "B/op",
            "extra": "759 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "759 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 21177,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "55413 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 21177,
            "unit": "ns/op",
            "extra": "55413 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "55413 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "55413 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 101,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10872864 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 101,
            "unit": "ns/op",
            "extra": "10872864 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "10872864 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "10872864 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22998,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52303 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22998,
            "unit": "ns/op",
            "extra": "52303 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52303 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52303 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 240533,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4700 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 240533,
            "unit": "ns/op",
            "extra": "4700 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4700 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4700 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2597313,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "398 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2597313,
            "unit": "ns/op",
            "extra": "398 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "398 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "398 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 46125134,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "63 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 46125134,
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
            "value": 10133,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "120008 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10133,
            "unit": "ns/op",
            "extra": "120008 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "120008 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "120008 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7787739,
            "unit": "ns/op\t 4594956 B/op\t   75254 allocs/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7787739,
            "unit": "ns/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594956,
            "unit": "B/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75254,
            "unit": "allocs/op",
            "extra": "151 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8372,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "137590 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8372,
            "unit": "ns/op",
            "extra": "137590 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "137590 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "137590 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 879747,
            "unit": "ns/op\t  396805 B/op\t    6234 allocs/op",
            "extra": "1207 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 879747,
            "unit": "ns/op",
            "extra": "1207 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396805,
            "unit": "B/op",
            "extra": "1207 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6234,
            "unit": "allocs/op",
            "extra": "1207 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11734,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "103791 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11734,
            "unit": "ns/op",
            "extra": "103791 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "103791 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "103791 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 9075682,
            "unit": "ns/op\t 4981820 B/op\t   81265 allocs/op",
            "extra": "129 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 9075682,
            "unit": "ns/op",
            "extra": "129 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981820,
            "unit": "B/op",
            "extra": "129 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81265,
            "unit": "allocs/op",
            "extra": "129 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 240.6,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5162575 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 240.6,
            "unit": "ns/op",
            "extra": "5162575 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5162575 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5162575 times\n4 procs"
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
            "value": 545.5,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2287167 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 545.5,
            "unit": "ns/op",
            "extra": "2287167 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2287167 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2287167 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 824,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1581690 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 824,
            "unit": "ns/op",
            "extra": "1581690 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1581690 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1581690 times\n4 procs"
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
          "id": "ae372cd321d36ae64b20aa4a901f88894d3b5ed0",
          "message": "feat(backendtlspolicy): enqueue for ConfigMaps (#6837)",
          "timestamp": "2024-12-13T17:40:32+01:00",
          "tree_id": "397725ebde376daca5baeca5e10cc008ceac1610",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/ae372cd321d36ae64b20aa4a901f88894d3b5ed0"
        },
        "date": 1734108192464,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1408967,
            "unit": "ns/op\t  819447 B/op\t       5 allocs/op",
            "extra": "963 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1408967,
            "unit": "ns/op",
            "extra": "963 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819447,
            "unit": "B/op",
            "extra": "963 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "963 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7434,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "161890 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7434,
            "unit": "ns/op",
            "extra": "161890 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "161890 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "161890 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.72,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14910828 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.72,
            "unit": "ns/op",
            "extra": "14910828 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14910828 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14910828 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22583,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52340 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22583,
            "unit": "ns/op",
            "extra": "52340 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52340 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52340 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 222870,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5134 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 222870,
            "unit": "ns/op",
            "extra": "5134 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
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
            "value": 2967847,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "405 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2967847,
            "unit": "ns/op",
            "extra": "405 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "405 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "405 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 33373094,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 33373094,
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
            "value": 9957,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "119473 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9957,
            "unit": "ns/op",
            "extra": "119473 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "119473 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "119473 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7482423,
            "unit": "ns/op\t 4594686 B/op\t   75253 allocs/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7482423,
            "unit": "ns/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594686,
            "unit": "B/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75253,
            "unit": "allocs/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8455,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "142405 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8455,
            "unit": "ns/op",
            "extra": "142405 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "142405 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "142405 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 873601,
            "unit": "ns/op\t  396806 B/op\t    6234 allocs/op",
            "extra": "1208 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 873601,
            "unit": "ns/op",
            "extra": "1208 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396806,
            "unit": "B/op",
            "extra": "1208 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6234,
            "unit": "allocs/op",
            "extra": "1208 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11404,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "104823 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11404,
            "unit": "ns/op",
            "extra": "104823 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "104823 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "104823 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8492192,
            "unit": "ns/op\t 4981088 B/op\t   81263 allocs/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8492192,
            "unit": "ns/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981088,
            "unit": "B/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81263,
            "unit": "allocs/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 231.5,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5138775 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 231.5,
            "unit": "ns/op",
            "extra": "5138775 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5138775 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5138775 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.619,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.619,
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
            "value": 514.4,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2331774 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 514.4,
            "unit": "ns/op",
            "extra": "2331774 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2331774 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2331774 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 741.9,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1615675 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 741.9,
            "unit": "ns/op",
            "extra": "1615675 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1615675 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1615675 times\n4 procs"
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
          "id": "1976f88961197ddca5203e0545c7c8857e9179ca",
          "message": "refactor: update conflict error is not an error, don't report it but requeue (#6833)",
          "timestamp": "2024-12-13T20:34:31+01:00",
          "tree_id": "64eea583680670a0c6e2dd99bfe7d0c7b49777df",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/1976f88961197ddca5203e0545c7c8857e9179ca"
        },
        "date": 1734118633257,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1275184,
            "unit": "ns/op\t  819448 B/op\t       5 allocs/op",
            "extra": "904 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1275184,
            "unit": "ns/op",
            "extra": "904 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819448,
            "unit": "B/op",
            "extra": "904 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "904 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 10082,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "99636 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 10082,
            "unit": "ns/op",
            "extra": "99636 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "99636 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "99636 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.08,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15190530 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.08,
            "unit": "ns/op",
            "extra": "15190530 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15190530 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15190530 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22557,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53125 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22557,
            "unit": "ns/op",
            "extra": "53125 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53125 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53125 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 270933,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5484 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 270933,
            "unit": "ns/op",
            "extra": "5484 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5484 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5484 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2578137,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "456 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2578137,
            "unit": "ns/op",
            "extra": "456 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "456 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "456 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 36804416,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "75 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 36804416,
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
            "value": 9743,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "116718 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9743,
            "unit": "ns/op",
            "extra": "116718 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "116718 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "116718 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7465061,
            "unit": "ns/op\t 4594707 B/op\t   75253 allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7465061,
            "unit": "ns/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594707,
            "unit": "B/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75253,
            "unit": "allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8210,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "145540 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8210,
            "unit": "ns/op",
            "extra": "145540 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "145540 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "145540 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 864914,
            "unit": "ns/op\t  396688 B/op\t    6232 allocs/op",
            "extra": "1245 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 864914,
            "unit": "ns/op",
            "extra": "1245 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396688,
            "unit": "B/op",
            "extra": "1245 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6232,
            "unit": "allocs/op",
            "extra": "1245 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11226,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "106562 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11226,
            "unit": "ns/op",
            "extra": "106562 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "106562 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "106562 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8708078,
            "unit": "ns/op\t 4981151 B/op\t   81263 allocs/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8708078,
            "unit": "ns/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981151,
            "unit": "B/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81263,
            "unit": "allocs/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 233.1,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5168005 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 233.1,
            "unit": "ns/op",
            "extra": "5168005 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5168005 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5168005 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6217,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6217,
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
            "value": 515.7,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2324166 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 515.7,
            "unit": "ns/op",
            "extra": "2324166 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2324166 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2324166 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 749.7,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1600682 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 749.7,
            "unit": "ns/op",
            "extra": "1600682 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1600682 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1600682 times\n4 procs"
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
          "id": "65896ed0a6b0f2ad6913616d627f96e4273009a2",
          "message": "chore(deps): update helm release kong to v2.45.0 (#6839)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2024-12-16T16:23:53+08:00",
          "tree_id": "e38e4155ce7b3b084be34155c7c780b79252a520",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/65896ed0a6b0f2ad6913616d627f96e4273009a2"
        },
        "date": 1734337589418,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1281843,
            "unit": "ns/op\t  819443 B/op\t       5 allocs/op",
            "extra": "1028 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1281843,
            "unit": "ns/op",
            "extra": "1028 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819443,
            "unit": "B/op",
            "extra": "1028 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1028 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 6998,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "171093 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 6998,
            "unit": "ns/op",
            "extra": "171093 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "171093 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "171093 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.88,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15159490 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.88,
            "unit": "ns/op",
            "extra": "15159490 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15159490 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15159490 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22531,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52921 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22531,
            "unit": "ns/op",
            "extra": "52921 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52921 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52921 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 224617,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5324 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 224617,
            "unit": "ns/op",
            "extra": "5324 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5324 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5324 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2601760,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "452 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2601760,
            "unit": "ns/op",
            "extra": "452 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "452 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "452 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 31025398,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 31025398,
            "unit": "ns/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 9952,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "117691 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9952,
            "unit": "ns/op",
            "extra": "117691 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "117691 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "117691 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7491760,
            "unit": "ns/op\t 4594913 B/op\t   75254 allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7491760,
            "unit": "ns/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594913,
            "unit": "B/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75254,
            "unit": "allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8341,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "143578 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8341,
            "unit": "ns/op",
            "extra": "143578 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "143578 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "143578 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 877465,
            "unit": "ns/op\t  396684 B/op\t    6232 allocs/op",
            "extra": "1242 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 877465,
            "unit": "ns/op",
            "extra": "1242 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396684,
            "unit": "B/op",
            "extra": "1242 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6232,
            "unit": "allocs/op",
            "extra": "1242 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11472,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "104802 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11472,
            "unit": "ns/op",
            "extra": "104802 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "104802 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "104802 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8563728,
            "unit": "ns/op\t 4981700 B/op\t   81265 allocs/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8563728,
            "unit": "ns/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981700,
            "unit": "B/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81265,
            "unit": "allocs/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 230.4,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5144526 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 230.4,
            "unit": "ns/op",
            "extra": "5144526 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5144526 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5144526 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6193,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6193,
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
            "value": 515.4,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2340772 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 515.4,
            "unit": "ns/op",
            "extra": "2340772 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2340772 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2340772 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 820.4,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1616986 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 820.4,
            "unit": "ns/op",
            "extra": "1616986 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1616986 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1616986 times\n4 procs"
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
          "id": "425d3ef39a411fea73bd9c3b61a726d7254578a5",
          "message": "fix/prevent_panic_in_cachestores (#6840)",
          "timestamp": "2024-12-16T17:19:44+08:00",
          "tree_id": "697eebb923016da20889f572fcd865b07575b560",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/425d3ef39a411fea73bd9c3b61a726d7254578a5"
        },
        "date": 1734340948017,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1283207,
            "unit": "ns/op\t  819442 B/op\t       5 allocs/op",
            "extra": "903 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1283207,
            "unit": "ns/op",
            "extra": "903 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819442,
            "unit": "B/op",
            "extra": "903 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "903 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7568,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "160021 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7568,
            "unit": "ns/op",
            "extra": "160021 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "160021 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "160021 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.22,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15193149 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.22,
            "unit": "ns/op",
            "extra": "15193149 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15193149 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15193149 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22562,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53348 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22562,
            "unit": "ns/op",
            "extra": "53348 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53348 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53348 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 233083,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4735 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 233083,
            "unit": "ns/op",
            "extra": "4735 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4735 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4735 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2542397,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "457 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2542397,
            "unit": "ns/op",
            "extra": "457 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "457 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "457 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 34076641,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 34076641,
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
            "value": 9931,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "118815 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9931,
            "unit": "ns/op",
            "extra": "118815 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "118815 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "118815 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 8000781,
            "unit": "ns/op\t 4594814 B/op\t   75254 allocs/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 8000781,
            "unit": "ns/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594814,
            "unit": "B/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75254,
            "unit": "allocs/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8288,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "141638 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8288,
            "unit": "ns/op",
            "extra": "141638 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "141638 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "141638 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 872480,
            "unit": "ns/op\t  396663 B/op\t    6232 allocs/op",
            "extra": "1249 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 872480,
            "unit": "ns/op",
            "extra": "1249 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396663,
            "unit": "B/op",
            "extra": "1249 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6232,
            "unit": "allocs/op",
            "extra": "1249 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11383,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "103591 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11383,
            "unit": "ns/op",
            "extra": "103591 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "103591 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "103591 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8611128,
            "unit": "ns/op\t 4981395 B/op\t   81264 allocs/op",
            "extra": "135 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8611128,
            "unit": "ns/op",
            "extra": "135 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981395,
            "unit": "B/op",
            "extra": "135 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81264,
            "unit": "allocs/op",
            "extra": "135 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 230.8,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "4828424 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 230.8,
            "unit": "ns/op",
            "extra": "4828424 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "4828424 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4828424 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6203,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6203,
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
            "value": 516.9,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2305443 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 516.9,
            "unit": "ns/op",
            "extra": "2305443 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2305443 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2305443 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 752.9,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1612218 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 752.9,
            "unit": "ns/op",
            "extra": "1612218 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1612218 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1612218 times\n4 procs"
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
          "id": "31e6fc59f8bd562d4062adede718e85b0db27d9d",
          "message": "tests(envtest): disable metrics by default (#6843)",
          "timestamp": "2024-12-16T13:26:39+01:00",
          "tree_id": "27c96100e96ee7a4b88f8882b3e1d1de776743cd",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/31e6fc59f8bd562d4062adede718e85b0db27d9d"
        },
        "date": 1734352167202,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1316946,
            "unit": "ns/op\t  819449 B/op\t       5 allocs/op",
            "extra": "937 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1316946,
            "unit": "ns/op",
            "extra": "937 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819449,
            "unit": "B/op",
            "extra": "937 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "937 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 18248,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "56155 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 18248,
            "unit": "ns/op",
            "extra": "56155 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "56155 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "56155 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 93.34,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11146592 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 93.34,
            "unit": "ns/op",
            "extra": "11146592 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11146592 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11146592 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22811,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52598 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22811,
            "unit": "ns/op",
            "extra": "52598 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52598 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52598 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 238981,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4503 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 238981,
            "unit": "ns/op",
            "extra": "4503 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4503 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4503 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 3088301,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "374 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 3088301,
            "unit": "ns/op",
            "extra": "374 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "374 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "374 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 38519395,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "67 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 38519395,
            "unit": "ns/op",
            "extra": "67 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 10141,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "112700 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10141,
            "unit": "ns/op",
            "extra": "112700 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "112700 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "112700 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7883750,
            "unit": "ns/op\t 4594654 B/op\t   75253 allocs/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7883750,
            "unit": "ns/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594654,
            "unit": "B/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75253,
            "unit": "allocs/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8426,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "140470 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8426,
            "unit": "ns/op",
            "extra": "140470 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "140470 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "140470 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 881643,
            "unit": "ns/op\t  396805 B/op\t    6234 allocs/op",
            "extra": "1206 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 881643,
            "unit": "ns/op",
            "extra": "1206 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396805,
            "unit": "B/op",
            "extra": "1206 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6234,
            "unit": "allocs/op",
            "extra": "1206 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11727,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "104623 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11727,
            "unit": "ns/op",
            "extra": "104623 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "104623 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "104623 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 9058651,
            "unit": "ns/op\t 4981732 B/op\t   81265 allocs/op",
            "extra": "132 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 9058651,
            "unit": "ns/op",
            "extra": "132 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981732,
            "unit": "B/op",
            "extra": "132 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81265,
            "unit": "allocs/op",
            "extra": "132 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 233.7,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5123323 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 233.7,
            "unit": "ns/op",
            "extra": "5123323 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5123323 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5123323 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6591,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6591,
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
            "value": 570.5,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "1764705 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 570.5,
            "unit": "ns/op",
            "extra": "1764705 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "1764705 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1764705 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 741.8,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1593528 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 741.8,
            "unit": "ns/op",
            "extra": "1593528 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1593528 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1593528 times\n4 procs"
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
          "id": "3cf39d23af2a91663ae5aae85596ddae29ffd3c6",
          "message": "feat(KongPlugin): generate K8s events for missing grant (#6841)",
          "timestamp": "2024-12-16T13:17:39Z",
          "tree_id": "213a0cbc428dde8cbd7bc44949040a877af8daac",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/3cf39d23af2a91663ae5aae85596ddae29ffd3c6"
        },
        "date": 1734355218360,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1305154,
            "unit": "ns/op\t  819450 B/op\t       5 allocs/op",
            "extra": "858 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1305154,
            "unit": "ns/op",
            "extra": "858 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819450,
            "unit": "B/op",
            "extra": "858 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "858 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 21013,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "57404 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 21013,
            "unit": "ns/op",
            "extra": "57404 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "57404 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "57404 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 98.44,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10678996 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 98.44,
            "unit": "ns/op",
            "extra": "10678996 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "10678996 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "10678996 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22953,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52597 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22953,
            "unit": "ns/op",
            "extra": "52597 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52597 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52597 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 238580,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4405 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 238580,
            "unit": "ns/op",
            "extra": "4405 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4405 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4405 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2589453,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "446 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2589453,
            "unit": "ns/op",
            "extra": "446 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
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
            "value": 43344404,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "28 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 43344404,
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
            "value": 10139,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "116671 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10139,
            "unit": "ns/op",
            "extra": "116671 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "116671 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "116671 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 8118556,
            "unit": "ns/op\t 4594990 B/op\t   75254 allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 8118556,
            "unit": "ns/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594990,
            "unit": "B/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75254,
            "unit": "allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8468,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "141082 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8468,
            "unit": "ns/op",
            "extra": "141082 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "141082 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "141082 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 880802,
            "unit": "ns/op\t  396687 B/op\t    6232 allocs/op",
            "extra": "1240 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 880802,
            "unit": "ns/op",
            "extra": "1240 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396687,
            "unit": "B/op",
            "extra": "1240 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6232,
            "unit": "allocs/op",
            "extra": "1240 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11650,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "104340 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11650,
            "unit": "ns/op",
            "extra": "104340 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "104340 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "104340 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 9388221,
            "unit": "ns/op\t 4980986 B/op\t   81262 allocs/op",
            "extra": "133 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 9388221,
            "unit": "ns/op",
            "extra": "133 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4980986,
            "unit": "B/op",
            "extra": "133 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81262,
            "unit": "allocs/op",
            "extra": "133 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 231.9,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5124832 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 231.9,
            "unit": "ns/op",
            "extra": "5124832 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5124832 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5124832 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6212,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6212,
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
            "value": 675.5,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "1977818 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 675.5,
            "unit": "ns/op",
            "extra": "1977818 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "1977818 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1977818 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 750,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1605576 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 750,
            "unit": "ns/op",
            "extra": "1605576 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1605576 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1605576 times\n4 procs"
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
          "id": "15200c712bcb77c57a69a852bd1689caaa884d96",
          "message": "chore(deps): update module github.com/kong/kubernetes-configuration to v1 (#6832)\n\n* chore(deps): update module github.com/kong/kubernetes-configuration to v1\r\n\r\n* chore: regenerate\r\n\r\n---------\r\n\r\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>\r\nCo-authored-by: github-actions <github-actions@users.noreply.github.com>",
          "timestamp": "2024-12-16T13:22:24Z",
          "tree_id": "ec3275f649b2d1475f0a30a862875eaceaa0b9d9",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/15200c712bcb77c57a69a852bd1689caaa884d96"
        },
        "date": 1734355555639,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1359297,
            "unit": "ns/op\t  819442 B/op\t       5 allocs/op",
            "extra": "1086 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1359297,
            "unit": "ns/op",
            "extra": "1086 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819442,
            "unit": "B/op",
            "extra": "1086 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1086 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 8180,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "162372 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 8180,
            "unit": "ns/op",
            "extra": "162372 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "162372 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "162372 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 80.04,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12687148 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 80.04,
            "unit": "ns/op",
            "extra": "12687148 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12687148 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12687148 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 23868,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "51169 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 23868,
            "unit": "ns/op",
            "extra": "51169 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "51169 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "51169 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 241023,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5128 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 241023,
            "unit": "ns/op",
            "extra": "5128 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
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
            "value": 2565980,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "408 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2565980,
            "unit": "ns/op",
            "extra": "408 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "408 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "408 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 43364795,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "55 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 43364795,
            "unit": "ns/op",
            "extra": "55 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
            "unit": "B/op",
            "extra": "55 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "55 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 10760,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "111254 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10760,
            "unit": "ns/op",
            "extra": "111254 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "111254 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "111254 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 8795406,
            "unit": "ns/op\t 4594920 B/op\t   75254 allocs/op",
            "extra": "132 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 8795406,
            "unit": "ns/op",
            "extra": "132 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594920,
            "unit": "B/op",
            "extra": "132 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75254,
            "unit": "allocs/op",
            "extra": "132 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8697,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "135682 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8697,
            "unit": "ns/op",
            "extra": "135682 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "135682 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "135682 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 901156,
            "unit": "ns/op\t  396960 B/op\t    6236 allocs/op",
            "extra": "1164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 901156,
            "unit": "ns/op",
            "extra": "1164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396960,
            "unit": "B/op",
            "extra": "1164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6236,
            "unit": "allocs/op",
            "extra": "1164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 12189,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "99154 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 12189,
            "unit": "ns/op",
            "extra": "99154 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "99154 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "99154 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 9476860,
            "unit": "ns/op\t 4981420 B/op\t   81264 allocs/op",
            "extra": "122 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 9476860,
            "unit": "ns/op",
            "extra": "122 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981420,
            "unit": "B/op",
            "extra": "122 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81264,
            "unit": "allocs/op",
            "extra": "122 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 234.1,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5132185 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 234.1,
            "unit": "ns/op",
            "extra": "5132185 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5132185 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5132185 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6194,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6194,
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
            "value": 522.1,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2309656 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 522.1,
            "unit": "ns/op",
            "extra": "2309656 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2309656 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2309656 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 749.1,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1606149 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 749.1,
            "unit": "ns/op",
            "extra": "1606149 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1606149 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1606149 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "czeslavo@gmail.com",
            "name": "Grzegorz Burzyński",
            "username": "czeslavo"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "7d370af358d2cad97bbda43375ccd8707276df30",
          "message": "chore: make CA certificate annotations use plural form (#6845)",
          "timestamp": "2024-12-16T13:43:06Z",
          "tree_id": "5e9c35ab650e2e39c7536e1ab9bb839bc1b45966",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/7d370af358d2cad97bbda43375ccd8707276df30"
        },
        "date": 1734356779366,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1387687,
            "unit": "ns/op\t  819441 B/op\t       5 allocs/op",
            "extra": "727 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1387687,
            "unit": "ns/op",
            "extra": "727 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819441,
            "unit": "B/op",
            "extra": "727 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "727 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 22743,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "52818 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 22743,
            "unit": "ns/op",
            "extra": "52818 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "52818 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "52818 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 78.95,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15131665 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 78.95,
            "unit": "ns/op",
            "extra": "15131665 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15131665 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15131665 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 23082,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52180 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 23082,
            "unit": "ns/op",
            "extra": "52180 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52180 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52180 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 241258,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4515 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 241258,
            "unit": "ns/op",
            "extra": "4515 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4515 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4515 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2740172,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "446 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2740172,
            "unit": "ns/op",
            "extra": "446 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
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
            "value": 45599775,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 45599775,
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
            "value": 10231,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "117782 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10231,
            "unit": "ns/op",
            "extra": "117782 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "117782 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "117782 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7935013,
            "unit": "ns/op\t 4594679 B/op\t   75253 allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7935013,
            "unit": "ns/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594679,
            "unit": "B/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75253,
            "unit": "allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8510,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "140018 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8510,
            "unit": "ns/op",
            "extra": "140018 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "140018 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "140018 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 901394,
            "unit": "ns/op\t  396930 B/op\t    6236 allocs/op",
            "extra": "1171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 901394,
            "unit": "ns/op",
            "extra": "1171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396930,
            "unit": "B/op",
            "extra": "1171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6236,
            "unit": "allocs/op",
            "extra": "1171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11956,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "101559 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11956,
            "unit": "ns/op",
            "extra": "101559 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "101559 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "101559 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 9095585,
            "unit": "ns/op\t 4981366 B/op\t   81264 allocs/op",
            "extra": "127 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 9095585,
            "unit": "ns/op",
            "extra": "127 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981366,
            "unit": "B/op",
            "extra": "127 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81264,
            "unit": "allocs/op",
            "extra": "127 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 273.8,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "4530897 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 273.8,
            "unit": "ns/op",
            "extra": "4530897 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "4530897 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4530897 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6227,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6227,
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
            "value": 527.8,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2239450 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 527.8,
            "unit": "ns/op",
            "extra": "2239450 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2239450 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2239450 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 752.7,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1552827 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 752.7,
            "unit": "ns/op",
            "extra": "1552827 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1552827 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1552827 times\n4 procs"
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
          "id": "965a7f4d1f8c7e04a44aa8d91742cf6880142e5a",
          "message": "chore(deps): bump github.com/Kong/sdk-konnect-go from 0.1.11 to 0.1.12 (#6846)\n\nBumps [github.com/Kong/sdk-konnect-go](https://github.com/Kong/sdk-konnect-go) from 0.1.11 to 0.1.12.\n- [Release notes](https://github.com/Kong/sdk-konnect-go/releases)\n- [Changelog](https://github.com/Kong/sdk-konnect-go/blob/main/.goreleaser.yml)\n- [Commits](https://github.com/Kong/sdk-konnect-go/compare/v0.1.11...v0.1.12)\n\n---\nupdated-dependencies:\n- dependency-name: github.com/Kong/sdk-konnect-go\n  dependency-type: direct:production\n  update-type: version-update:semver-patch\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2024-12-16T15:16:16Z",
          "tree_id": "aa6d450166c881689574f0ea87dce9d3c5ccd0d4",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/965a7f4d1f8c7e04a44aa8d91742cf6880142e5a"
        },
        "date": 1734362378253,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1204040,
            "unit": "ns/op\t  819441 B/op\t       5 allocs/op",
            "extra": "901 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1204040,
            "unit": "ns/op",
            "extra": "901 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819441,
            "unit": "B/op",
            "extra": "901 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "901 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7522,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "143128 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7522,
            "unit": "ns/op",
            "extra": "143128 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "143128 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "143128 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.17,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14887690 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.17,
            "unit": "ns/op",
            "extra": "14887690 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14887690 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14887690 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 27214,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52588 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 27214,
            "unit": "ns/op",
            "extra": "52588 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52588 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52588 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 229186,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5252 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 229186,
            "unit": "ns/op",
            "extra": "5252 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5252 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5252 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2559141,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "404 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2559141,
            "unit": "ns/op",
            "extra": "404 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "404 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "404 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 38635927,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 38635927,
            "unit": "ns/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 10042,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "118041 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10042,
            "unit": "ns/op",
            "extra": "118041 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "118041 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "118041 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7569928,
            "unit": "ns/op\t 4594692 B/op\t   75253 allocs/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7569928,
            "unit": "ns/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594692,
            "unit": "B/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75253,
            "unit": "allocs/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8569,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "141897 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8569,
            "unit": "ns/op",
            "extra": "141897 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "141897 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "141897 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 891799,
            "unit": "ns/op\t  396889 B/op\t    6235 allocs/op",
            "extra": "1183 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 891799,
            "unit": "ns/op",
            "extra": "1183 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396889,
            "unit": "B/op",
            "extra": "1183 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6235,
            "unit": "allocs/op",
            "extra": "1183 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11579,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "102874 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11579,
            "unit": "ns/op",
            "extra": "102874 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "102874 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "102874 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8692588,
            "unit": "ns/op\t 4981700 B/op\t   81265 allocs/op",
            "extra": "135 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8692588,
            "unit": "ns/op",
            "extra": "135 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981700,
            "unit": "B/op",
            "extra": "135 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81265,
            "unit": "allocs/op",
            "extra": "135 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 231.1,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5235456 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 231.1,
            "unit": "ns/op",
            "extra": "5235456 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5235456 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5235456 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6187,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6187,
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
            "value": 516.7,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2311052 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 516.7,
            "unit": "ns/op",
            "extra": "2311052 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2311052 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2311052 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 748.6,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1576574 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 748.6,
            "unit": "ns/op",
            "extra": "1576574 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1576574 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1576574 times\n4 procs"
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
          "id": "8d56485b72e5210eaa873b5ee96184391fdc9174",
          "message": "chore(deps): update kindest/node@only-patch docker tag to v1.29.12 (#6849)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2024-12-17T10:52:52+08:00",
          "tree_id": "f37da3d8d2e92e1952c4d7790886b8aaa89590c0",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/8d56485b72e5210eaa873b5ee96184391fdc9174"
        },
        "date": 1734404172721,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1323652,
            "unit": "ns/op\t  819447 B/op\t       5 allocs/op",
            "extra": "898 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1323652,
            "unit": "ns/op",
            "extra": "898 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819447,
            "unit": "B/op",
            "extra": "898 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "898 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7555,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "158866 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7555,
            "unit": "ns/op",
            "extra": "158866 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "158866 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "158866 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 82.02,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15177154 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 82.02,
            "unit": "ns/op",
            "extra": "15177154 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15177154 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15177154 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22520,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53034 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22520,
            "unit": "ns/op",
            "extra": "53034 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53034 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53034 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 228550,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5403 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 228550,
            "unit": "ns/op",
            "extra": "5403 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5403 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5403 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2535604,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "430 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2535604,
            "unit": "ns/op",
            "extra": "430 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
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
            "value": 40130688,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 40130688,
            "unit": "ns/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 9865,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "119654 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9865,
            "unit": "ns/op",
            "extra": "119654 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "119654 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "119654 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7524102,
            "unit": "ns/op\t 4594909 B/op\t   75254 allocs/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7524102,
            "unit": "ns/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594909,
            "unit": "B/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75254,
            "unit": "allocs/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8459,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "143214 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8459,
            "unit": "ns/op",
            "extra": "143214 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "143214 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "143214 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 859416,
            "unit": "ns/op\t  396660 B/op\t    6232 allocs/op",
            "extra": "1250 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 859416,
            "unit": "ns/op",
            "extra": "1250 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396660,
            "unit": "B/op",
            "extra": "1250 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6232,
            "unit": "allocs/op",
            "extra": "1250 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11316,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "106581 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11316,
            "unit": "ns/op",
            "extra": "106581 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "106581 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "106581 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8609125,
            "unit": "ns/op\t 4981067 B/op\t   81263 allocs/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8609125,
            "unit": "ns/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981067,
            "unit": "B/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81263,
            "unit": "allocs/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 236.7,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5164173 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 236.7,
            "unit": "ns/op",
            "extra": "5164173 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5164173 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5164173 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6201,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6201,
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
            "value": 518,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2319902 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 518,
            "unit": "ns/op",
            "extra": "2319902 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2319902 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2319902 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 740.6,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1615134 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 740.6,
            "unit": "ns/op",
            "extra": "1615134 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1615134 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1615134 times\n4 procs"
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
          "id": "700afd6f582ef8c236aa14841444f0c2b38f3bd8",
          "message": "chore(deps): update kindest/node@only-patch docker tag to v1.30.8 (#6850)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2024-12-17T09:48:20Z",
          "tree_id": "a31c18c7b8ad37d99d36c6c3459878870cd13205",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/700afd6f582ef8c236aa14841444f0c2b38f3bd8"
        },
        "date": 1734429095710,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1235377,
            "unit": "ns/op\t  819448 B/op\t       5 allocs/op",
            "extra": "912 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1235377,
            "unit": "ns/op",
            "extra": "912 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819448,
            "unit": "B/op",
            "extra": "912 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "912 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 13177,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "77526 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 13177,
            "unit": "ns/op",
            "extra": "77526 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "77526 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "77526 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 81.95,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15068947 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 81.95,
            "unit": "ns/op",
            "extra": "15068947 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15068947 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15068947 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22290,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53426 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22290,
            "unit": "ns/op",
            "extra": "53426 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53426 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53426 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 219208,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4928 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 219208,
            "unit": "ns/op",
            "extra": "4928 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4928 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4928 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2550751,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "469 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2550751,
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
            "value": 34087665,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 34087665,
            "unit": "ns/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 9878,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "118406 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9878,
            "unit": "ns/op",
            "extra": "118406 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "118406 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "118406 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7388557,
            "unit": "ns/op\t 4594465 B/op\t   75253 allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7388557,
            "unit": "ns/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594465,
            "unit": "B/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75253,
            "unit": "allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8332,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "140958 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8332,
            "unit": "ns/op",
            "extra": "140958 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "140958 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "140958 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 876192,
            "unit": "ns/op\t  396672 B/op\t    6232 allocs/op",
            "extra": "1251 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 876192,
            "unit": "ns/op",
            "extra": "1251 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396672,
            "unit": "B/op",
            "extra": "1251 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6232,
            "unit": "allocs/op",
            "extra": "1251 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11340,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "106020 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11340,
            "unit": "ns/op",
            "extra": "106020 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "106020 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "106020 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8513190,
            "unit": "ns/op\t 4981791 B/op\t   81265 allocs/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8513190,
            "unit": "ns/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981791,
            "unit": "B/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81265,
            "unit": "allocs/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 234,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5208817 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 234,
            "unit": "ns/op",
            "extra": "5208817 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5208817 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5208817 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6197,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6197,
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
            "value": 525.4,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2325352 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 525.4,
            "unit": "ns/op",
            "extra": "2325352 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2325352 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2325352 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 891.9,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1603174 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 891.9,
            "unit": "ns/op",
            "extra": "1603174 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1603174 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1603174 times\n4 procs"
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
          "id": "e3fd99e068bbd13c21d4cc5fb8fb8dbbd2dd8812",
          "message": "chore: add acl plugin examples for acl groups and consumer groups (#6847)",
          "timestamp": "2024-12-17T12:05:03+01:00",
          "tree_id": "ce73f0765742e3b8ddd1e927a66bbacea7d11c7b",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/e3fd99e068bbd13c21d4cc5fb8fb8dbbd2dd8812"
        },
        "date": 1734433702036,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1283790,
            "unit": "ns/op\t  819448 B/op\t       5 allocs/op",
            "extra": "817 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1283790,
            "unit": "ns/op",
            "extra": "817 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819448,
            "unit": "B/op",
            "extra": "817 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "817 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 19116,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "57120 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 19116,
            "unit": "ns/op",
            "extra": "57120 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "57120 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "57120 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 89.78,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11981442 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 89.78,
            "unit": "ns/op",
            "extra": "11981442 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11981442 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11981442 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 23183,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "51823 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 23183,
            "unit": "ns/op",
            "extra": "51823 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "51823 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "51823 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 235227,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4482 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 235227,
            "unit": "ns/op",
            "extra": "4482 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4482 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4482 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2918169,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "458 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2918169,
            "unit": "ns/op",
            "extra": "458 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
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
            "value": 39102728,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 39102728,
            "unit": "ns/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 10121,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "119618 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10121,
            "unit": "ns/op",
            "extra": "119618 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "119618 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "119618 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7699099,
            "unit": "ns/op\t 4594734 B/op\t   75254 allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7699099,
            "unit": "ns/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594734,
            "unit": "B/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75254,
            "unit": "allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8392,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "140823 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8392,
            "unit": "ns/op",
            "extra": "140823 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "140823 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "140823 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 877719,
            "unit": "ns/op\t  396697 B/op\t    6232 allocs/op",
            "extra": "1238 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 877719,
            "unit": "ns/op",
            "extra": "1238 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396697,
            "unit": "B/op",
            "extra": "1238 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6232,
            "unit": "allocs/op",
            "extra": "1238 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11770,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "102775 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11770,
            "unit": "ns/op",
            "extra": "102775 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "102775 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "102775 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8793255,
            "unit": "ns/op\t 4981664 B/op\t   81265 allocs/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8793255,
            "unit": "ns/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981664,
            "unit": "B/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81265,
            "unit": "allocs/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 291.4,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5115670 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 291.4,
            "unit": "ns/op",
            "extra": "5115670 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5115670 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5115670 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6339,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6339,
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
            "value": 520.4,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2314236 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 520.4,
            "unit": "ns/op",
            "extra": "2314236 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2314236 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2314236 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 751,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1601931 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 751,
            "unit": "ns/op",
            "extra": "1601931 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1601931 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1601931 times\n4 procs"
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
          "id": "3794e98d14d748fa03403575befdaa3c27de26dd",
          "message": "fix(credentials): ignore missing secret field on secret validation for JWT credentials with HMAC algo (#6848)",
          "timestamp": "2024-12-17T11:28:37Z",
          "tree_id": "52f8b04e434946cba22747062e6bcafe683ee0a8",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/3794e98d14d748fa03403575befdaa3c27de26dd"
        },
        "date": 1734435123634,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1160970,
            "unit": "ns/op\t  819442 B/op\t       5 allocs/op",
            "extra": "992 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1160970,
            "unit": "ns/op",
            "extra": "992 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819442,
            "unit": "B/op",
            "extra": "992 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "992 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 8214,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "135259 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 8214,
            "unit": "ns/op",
            "extra": "135259 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "135259 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "135259 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 86.29,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15146059 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 86.29,
            "unit": "ns/op",
            "extra": "15146059 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15146059 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15146059 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22491,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52327 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22491,
            "unit": "ns/op",
            "extra": "52327 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52327 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52327 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 226825,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5502 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 226825,
            "unit": "ns/op",
            "extra": "5502 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5502 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5502 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2534417,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "430 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2534417,
            "unit": "ns/op",
            "extra": "430 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
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
            "value": 38590018,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "67 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 38590018,
            "unit": "ns/op",
            "extra": "67 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 9973,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "120014 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9973,
            "unit": "ns/op",
            "extra": "120014 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "120014 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "120014 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7583568,
            "unit": "ns/op\t 4595049 B/op\t   75255 allocs/op",
            "extra": "154 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7583568,
            "unit": "ns/op",
            "extra": "154 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4595049,
            "unit": "B/op",
            "extra": "154 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75255,
            "unit": "allocs/op",
            "extra": "154 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8412,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "142328 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8412,
            "unit": "ns/op",
            "extra": "142328 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "142328 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "142328 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 872971,
            "unit": "ns/op\t  396688 B/op\t    6232 allocs/op",
            "extra": "1244 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 872971,
            "unit": "ns/op",
            "extra": "1244 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396688,
            "unit": "B/op",
            "extra": "1244 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6232,
            "unit": "allocs/op",
            "extra": "1244 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11476,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "100899 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11476,
            "unit": "ns/op",
            "extra": "100899 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "100899 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "100899 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8637487,
            "unit": "ns/op\t 4981237 B/op\t   81263 allocs/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8637487,
            "unit": "ns/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981237,
            "unit": "B/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81263,
            "unit": "allocs/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 231.9,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5191309 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 231.9,
            "unit": "ns/op",
            "extra": "5191309 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5191309 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5191309 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6206,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6206,
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
            "value": 515.2,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2325567 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 515.2,
            "unit": "ns/op",
            "extra": "2325567 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2325567 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2325567 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 745.3,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1621293 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 745.3,
            "unit": "ns/op",
            "extra": "1621293 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1621293 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1621293 times\n4 procs"
          }
        ]
      }
    ]
  }
}