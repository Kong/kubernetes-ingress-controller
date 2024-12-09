window.BENCHMARK_DATA = {
  "lastUpdate": 1733733588921,
  "repoUrl": "https://github.com/Kong/kubernetes-ingress-controller",
  "entries": {
    "Go Benchmark": [
      {
        "commit": {
          "author": {
            "email": "patryk.malek@konghq.com",
            "name": "Patryk Małek",
            "username": "pmalek"
          },
          "committer": {
            "email": "patryk.malek@konghq.com",
            "name": "Patryk Małek",
            "username": "pmalek"
          },
          "distinct": true,
          "id": "2511996e7fceab64a174ba9381671c6335139220",
          "message": "chore(ci): fix performance tests artifacts",
          "timestamp": "2024-12-04T16:58:03+01:00",
          "tree_id": "7d8725ab8254b6fb6e820fc586163f291412a9b3",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/2511996e7fceab64a174ba9381671c6335139220"
        },
        "date": 1733328069148,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1391146,
            "unit": "ns/op\t  819441 B/op\t       5 allocs/op",
            "extra": "897 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1391146,
            "unit": "ns/op",
            "extra": "897 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819441,
            "unit": "B/op",
            "extra": "897 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "897 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 8401,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "156120 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 8401,
            "unit": "ns/op",
            "extra": "156120 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "156120 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "156120 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.64,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13329775 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.64,
            "unit": "ns/op",
            "extra": "13329775 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13329775 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13329775 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 24127,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "48074 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 24127,
            "unit": "ns/op",
            "extra": "48074 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "48074 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "48074 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 243903,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5107 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 243903,
            "unit": "ns/op",
            "extra": "5107 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5107 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5107 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2809402,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "423 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2809402,
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
            "value": 30050551,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "51 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 30050551,
            "unit": "ns/op",
            "extra": "51 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 10538,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "114996 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10538,
            "unit": "ns/op",
            "extra": "114996 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "114996 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "114996 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 8454248,
            "unit": "ns/op\t 4594500 B/op\t   75248 allocs/op",
            "extra": "135 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 8454248,
            "unit": "ns/op",
            "extra": "135 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594500,
            "unit": "B/op",
            "extra": "135 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75248,
            "unit": "allocs/op",
            "extra": "135 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8437,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "148257 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8437,
            "unit": "ns/op",
            "extra": "148257 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "148257 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "148257 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 901122,
            "unit": "ns/op\t  396662 B/op\t    6229 allocs/op",
            "extra": "1182 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 901122,
            "unit": "ns/op",
            "extra": "1182 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396662,
            "unit": "B/op",
            "extra": "1182 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6229,
            "unit": "allocs/op",
            "extra": "1182 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11717,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "105745 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11717,
            "unit": "ns/op",
            "extra": "105745 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "105745 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "105745 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 9522015,
            "unit": "ns/op\t 4981276 B/op\t   81258 allocs/op",
            "extra": "120 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 9522015,
            "unit": "ns/op",
            "extra": "120 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981276,
            "unit": "B/op",
            "extra": "120 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81258,
            "unit": "allocs/op",
            "extra": "120 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 230.2,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5127057 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 230.2,
            "unit": "ns/op",
            "extra": "5127057 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5127057 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5127057 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6214,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6214,
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
            "extra": "3974058 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 300,
            "unit": "ns/op",
            "extra": "3974058 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "3974058 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "3974058 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 456.7,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2644182 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 456.7,
            "unit": "ns/op",
            "extra": "2644182 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2644182 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2644182 times\n4 procs"
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
          "id": "024f89ee2cfee12f38c8538cae29ebb384a6cf5c",
          "message": "chore(tests): skip e2e TestDeployAllInOneDBLESS scaling proxy to 0 (#6777)",
          "timestamp": "2024-12-04T16:44:30Z",
          "tree_id": "f5af547bcbfd0bf825d922c77e40ecdfdd9558ac",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/024f89ee2cfee12f38c8538cae29ebb384a6cf5c"
        },
        "date": 1733330845812,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1152285,
            "unit": "ns/op\t  819441 B/op\t       5 allocs/op",
            "extra": "1058 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1152285,
            "unit": "ns/op",
            "extra": "1058 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819441,
            "unit": "B/op",
            "extra": "1058 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1058 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7532,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "159897 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7532,
            "unit": "ns/op",
            "extra": "159897 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "159897 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "159897 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.25,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15166698 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.25,
            "unit": "ns/op",
            "extra": "15166698 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15166698 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15166698 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22312,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53536 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22312,
            "unit": "ns/op",
            "extra": "53536 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53536 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53536 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 219496,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5463 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 219496,
            "unit": "ns/op",
            "extra": "5463 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5463 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5463 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2494272,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "470 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2494272,
            "unit": "ns/op",
            "extra": "470 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "470 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "470 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 29431410,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "75 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 29431410,
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
            "value": 9543,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "126982 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9543,
            "unit": "ns/op",
            "extra": "126982 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "126982 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "126982 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7262604,
            "unit": "ns/op\t 4594746 B/op\t   75248 allocs/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7262604,
            "unit": "ns/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594746,
            "unit": "B/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75248,
            "unit": "allocs/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8002,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "150793 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8002,
            "unit": "ns/op",
            "extra": "150793 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "150793 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "150793 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 843027,
            "unit": "ns/op\t  396435 B/op\t    6226 allocs/op",
            "extra": "1249 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 843027,
            "unit": "ns/op",
            "extra": "1249 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396435,
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
            "value": 10901,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "110227 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10901,
            "unit": "ns/op",
            "extra": "110227 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "110227 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "110227 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8382971,
            "unit": "ns/op\t 4981317 B/op\t   81258 allocs/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8382971,
            "unit": "ns/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981317,
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
            "value": 225.8,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5285564 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 225.8,
            "unit": "ns/op",
            "extra": "5285564 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5285564 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5285564 times\n4 procs"
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
            "value": 294.1,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4082334 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 294.1,
            "unit": "ns/op",
            "extra": "4082334 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4082334 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4082334 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 461.6,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2695912 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 461.6,
            "unit": "ns/op",
            "extra": "2695912 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2695912 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2695912 times\n4 procs"
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
            "email": "jakub.warczarek@konghq.com",
            "name": "Jakub Warczarek",
            "username": "programmer04"
          },
          "distinct": true,
          "id": "ba13166f47aa4736c5b7a7860656860238294dbb",
          "message": "fix(tests): keep all golden tests in one place",
          "timestamp": "2024-12-04T17:36:29+01:00",
          "tree_id": "7285cf201f0762f80176158e80bb5fca85ebba16",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/ba13166f47aa4736c5b7a7860656860238294dbb"
        },
        "date": 1733330875985,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1176373,
            "unit": "ns/op\t  819447 B/op\t       5 allocs/op",
            "extra": "1113 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1176373,
            "unit": "ns/op",
            "extra": "1113 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819447,
            "unit": "B/op",
            "extra": "1113 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1113 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7343,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "162458 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7343,
            "unit": "ns/op",
            "extra": "162458 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "162458 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "162458 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.41,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15076178 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.41,
            "unit": "ns/op",
            "extra": "15076178 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15076178 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15076178 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22471,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53467 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22471,
            "unit": "ns/op",
            "extra": "53467 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53467 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53467 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 250638,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4862 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 250638,
            "unit": "ns/op",
            "extra": "4862 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4862 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4862 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2529007,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "434 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2529007,
            "unit": "ns/op",
            "extra": "434 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "434 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "434 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 42903919,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "31 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 42903919,
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
            "value": 9724,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "122614 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9724,
            "unit": "ns/op",
            "extra": "122614 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "122614 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "122614 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7459805,
            "unit": "ns/op\t 4594524 B/op\t   75248 allocs/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7459805,
            "unit": "ns/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594524,
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
            "value": 8097,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "147050 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8097,
            "unit": "ns/op",
            "extra": "147050 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "147050 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "147050 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 879505,
            "unit": "ns/op\t  396524 B/op\t    6227 allocs/op",
            "extra": "1224 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 879505,
            "unit": "ns/op",
            "extra": "1224 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396524,
            "unit": "B/op",
            "extra": "1224 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6227,
            "unit": "allocs/op",
            "extra": "1224 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11130,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "105405 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11130,
            "unit": "ns/op",
            "extra": "105405 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "105405 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "105405 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8646830,
            "unit": "ns/op\t 4981011 B/op\t   81257 allocs/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8646830,
            "unit": "ns/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981011,
            "unit": "B/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81257,
            "unit": "allocs/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 226.6,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5235998 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 226.6,
            "unit": "ns/op",
            "extra": "5235998 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5235998 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5235998 times\n4 procs"
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
            "value": 296.4,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4081327 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 296.4,
            "unit": "ns/op",
            "extra": "4081327 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4081327 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4081327 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 445.9,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2689104 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 445.9,
            "unit": "ns/op",
            "extra": "2689104 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2689104 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2689104 times\n4 procs"
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
          "id": "48d335e211a1b197119020c50520ce5aec6b5de8",
          "message": "chore(ci): fix performance tests artifacts (#6778)",
          "timestamp": "2024-12-04T17:46:22+01:00",
          "tree_id": "b265187baa2d80f7a0b36667c8a759e3e00d9e60",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/48d335e211a1b197119020c50520ce5aec6b5de8"
        },
        "date": 1733330979033,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1185632,
            "unit": "ns/op\t  819448 B/op\t       5 allocs/op",
            "extra": "1006 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1185632,
            "unit": "ns/op",
            "extra": "1006 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819448,
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
            "value": 6907,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "172879 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 6907,
            "unit": "ns/op",
            "extra": "172879 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "172879 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "172879 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 78.93,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14037115 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 78.93,
            "unit": "ns/op",
            "extra": "14037115 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14037115 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14037115 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 24469,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "51853 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 24469,
            "unit": "ns/op",
            "extra": "51853 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "51853 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "51853 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 270391,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "3930 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 270391,
            "unit": "ns/op",
            "extra": "3930 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "3930 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3930 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2572276,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "460 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2572276,
            "unit": "ns/op",
            "extra": "460 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
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
            "value": 46097927,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "36 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 46097927,
            "unit": "ns/op",
            "extra": "36 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 9636,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "122222 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9636,
            "unit": "ns/op",
            "extra": "122222 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "122222 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "122222 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7373891,
            "unit": "ns/op\t 4594788 B/op\t   75249 allocs/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7373891,
            "unit": "ns/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594788,
            "unit": "B/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75249,
            "unit": "allocs/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8034,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "149119 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8034,
            "unit": "ns/op",
            "extra": "149119 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "149119 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "149119 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 857404,
            "unit": "ns/op\t  396422 B/op\t    6226 allocs/op",
            "extra": "1250 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 857404,
            "unit": "ns/op",
            "extra": "1250 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396422,
            "unit": "B/op",
            "extra": "1250 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6226,
            "unit": "allocs/op",
            "extra": "1250 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11039,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "107654 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11039,
            "unit": "ns/op",
            "extra": "107654 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "107654 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "107654 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8591838,
            "unit": "ns/op\t 4981201 B/op\t   81258 allocs/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8591838,
            "unit": "ns/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981201,
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
            "value": 225.9,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5277570 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 225.9,
            "unit": "ns/op",
            "extra": "5277570 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5277570 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5277570 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6192,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6192,
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
            "value": 351.8,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4093003 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 351.8,
            "unit": "ns/op",
            "extra": "4093003 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4093003 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4093003 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 445.2,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2587520 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 445.2,
            "unit": "ns/op",
            "extra": "2587520 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2587520 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2587520 times\n4 procs"
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
            "email": "jakub.warczarek@konghq.com",
            "name": "Jakub Warczarek",
            "username": "programmer04"
          },
          "distinct": true,
          "id": "3bfabc012ccafa080b31e147e4af24f0f5b1f773",
          "message": "fix(tests): keep all golden tests in one place and adjust",
          "timestamp": "2024-12-04T18:06:14+01:00",
          "tree_id": "91cae44335123505f16e05067d851c4274788a24",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/3bfabc012ccafa080b31e147e4af24f0f5b1f773"
        },
        "date": 1733332162775,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1304908,
            "unit": "ns/op\t  819442 B/op\t       5 allocs/op",
            "extra": "1010 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1304908,
            "unit": "ns/op",
            "extra": "1010 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819442,
            "unit": "B/op",
            "extra": "1010 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1010 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7489,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "161602 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7489,
            "unit": "ns/op",
            "extra": "161602 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "161602 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "161602 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.25,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15144075 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.25,
            "unit": "ns/op",
            "extra": "15144075 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15144075 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15144075 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22991,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52459 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22991,
            "unit": "ns/op",
            "extra": "52459 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52459 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52459 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 238808,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4221 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 238808,
            "unit": "ns/op",
            "extra": "4221 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4221 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4221 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2644107,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "388 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2644107,
            "unit": "ns/op",
            "extra": "388 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
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
            "value": 36537810,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "68 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 36537810,
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
            "value": 9942,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "117787 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9942,
            "unit": "ns/op",
            "extra": "117787 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "117787 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "117787 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7911779,
            "unit": "ns/op\t 4594299 B/op\t   75247 allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7911779,
            "unit": "ns/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594299,
            "unit": "B/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75247,
            "unit": "allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8226,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "146701 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8226,
            "unit": "ns/op",
            "extra": "146701 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "146701 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "146701 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 890515,
            "unit": "ns/op\t  396497 B/op\t    6227 allocs/op",
            "extra": "1230 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 890515,
            "unit": "ns/op",
            "extra": "1230 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396497,
            "unit": "B/op",
            "extra": "1230 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6227,
            "unit": "allocs/op",
            "extra": "1230 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11629,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "104360 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11629,
            "unit": "ns/op",
            "extra": "104360 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "104360 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "104360 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 9254370,
            "unit": "ns/op\t 4981065 B/op\t   81257 allocs/op",
            "extra": "128 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 9254370,
            "unit": "ns/op",
            "extra": "128 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981065,
            "unit": "B/op",
            "extra": "128 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81257,
            "unit": "allocs/op",
            "extra": "128 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 249.9,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5255187 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 249.9,
            "unit": "ns/op",
            "extra": "5255187 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5255187 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5255187 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.643,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.643,
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
            "value": 298.1,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4013037 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 298.1,
            "unit": "ns/op",
            "extra": "4013037 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4013037 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4013037 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 448.2,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2644342 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 448.2,
            "unit": "ns/op",
            "extra": "2644342 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2644342 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2644342 times\n4 procs"
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
            "email": "jakub.warczarek@konghq.com",
            "name": "Jakub Warczarek",
            "username": "programmer04"
          },
          "distinct": true,
          "id": "ead7d57b9f5438655904c9ee0afdbe731fa919ff",
          "message": "fix(tests): keep all golden tests in one place and adjust",
          "timestamp": "2024-12-04T18:13:11+01:00",
          "tree_id": "a5d155e36d8196405482f32226b092434647c411",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/ead7d57b9f5438655904c9ee0afdbe731fa919ff"
        },
        "date": 1733332581994,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1219295,
            "unit": "ns/op\t  819441 B/op\t       5 allocs/op",
            "extra": "1119 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1219295,
            "unit": "ns/op",
            "extra": "1119 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819441,
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
            "value": 6901,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "174109 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 6901,
            "unit": "ns/op",
            "extra": "174109 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "174109 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "174109 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 78.98,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15185272 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 78.98,
            "unit": "ns/op",
            "extra": "15185272 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15185272 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15185272 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22517,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52803 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22517,
            "unit": "ns/op",
            "extra": "52803 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52803 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52803 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 225016,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5611 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 225016,
            "unit": "ns/op",
            "extra": "5611 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5611 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5611 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2583155,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "412 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2583155,
            "unit": "ns/op",
            "extra": "412 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
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
            "value": 37087658,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "68 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 37087658,
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
            "value": 9644,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "118430 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9644,
            "unit": "ns/op",
            "extra": "118430 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "118430 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "118430 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7361235,
            "unit": "ns/op\t 4594398 B/op\t   75247 allocs/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7361235,
            "unit": "ns/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594398,
            "unit": "B/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75247,
            "unit": "allocs/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8009,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "148628 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8009,
            "unit": "ns/op",
            "extra": "148628 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "148628 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "148628 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 852232,
            "unit": "ns/op\t  396459 B/op\t    6226 allocs/op",
            "extra": "1239 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 852232,
            "unit": "ns/op",
            "extra": "1239 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396459,
            "unit": "B/op",
            "extra": "1239 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6226,
            "unit": "allocs/op",
            "extra": "1239 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11028,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "108085 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11028,
            "unit": "ns/op",
            "extra": "108085 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "108085 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "108085 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8548654,
            "unit": "ns/op\t 4981595 B/op\t   81259 allocs/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8548654,
            "unit": "ns/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981595,
            "unit": "B/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81259,
            "unit": "allocs/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 227.5,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5243743 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 227.5,
            "unit": "ns/op",
            "extra": "5243743 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5243743 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5243743 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6216,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6216,
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
            "value": 340.8,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "3273746 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 340.8,
            "unit": "ns/op",
            "extra": "3273746 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "3273746 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "3273746 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 449,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2657755 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 449,
            "unit": "ns/op",
            "extra": "2657755 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2657755 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2657755 times\n4 procs"
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
          "id": "98cf7aa31dbe167e29e069377c489c0ed7fd894f",
          "message": "fix(tests): keep all golden tests in one place and adjust (#6773)",
          "timestamp": "2024-12-04T18:35:41+01:00",
          "tree_id": "a5d155e36d8196405482f32226b092434647c411",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/98cf7aa31dbe167e29e069377c489c0ed7fd894f"
        },
        "date": 1733333939022,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1213652,
            "unit": "ns/op\t  819441 B/op\t       5 allocs/op",
            "extra": "985 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1213652,
            "unit": "ns/op",
            "extra": "985 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819441,
            "unit": "B/op",
            "extra": "985 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "985 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 9496,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "120685 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 9496,
            "unit": "ns/op",
            "extra": "120685 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "120685 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "120685 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.38,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15167247 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.38,
            "unit": "ns/op",
            "extra": "15167247 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15167247 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15167247 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22595,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53595 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22595,
            "unit": "ns/op",
            "extra": "53595 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53595 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53595 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 224064,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4510 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 224064,
            "unit": "ns/op",
            "extra": "4510 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
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
            "value": 2543374,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "422 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2543374,
            "unit": "ns/op",
            "extra": "422 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "422 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "422 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 32522174,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "33 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 32522174,
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
            "value": 9635,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "123879 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9635,
            "unit": "ns/op",
            "extra": "123879 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "123879 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "123879 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7391531,
            "unit": "ns/op\t 4594620 B/op\t   75248 allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7391531,
            "unit": "ns/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594620,
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
            "value": 8051,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "150174 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8051,
            "unit": "ns/op",
            "extra": "150174 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "150174 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "150174 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 863054,
            "unit": "ns/op\t  396562 B/op\t    6228 allocs/op",
            "extra": "1210 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 863054,
            "unit": "ns/op",
            "extra": "1210 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396562,
            "unit": "B/op",
            "extra": "1210 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6228,
            "unit": "allocs/op",
            "extra": "1210 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11127,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "109470 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11127,
            "unit": "ns/op",
            "extra": "109470 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "109470 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "109470 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8498558,
            "unit": "ns/op\t 4980898 B/op\t   81257 allocs/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8498558,
            "unit": "ns/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4980898,
            "unit": "B/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81257,
            "unit": "allocs/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 227.1,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5251978 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 227.1,
            "unit": "ns/op",
            "extra": "5251978 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5251978 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5251978 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6191,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6191,
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
            "value": 297.3,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4014541 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 297.3,
            "unit": "ns/op",
            "extra": "4014541 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4014541 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4014541 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 452.3,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2642479 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 452.3,
            "unit": "ns/op",
            "extra": "2642479 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2642479 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2642479 times\n4 procs"
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
          "id": "b7a39c36bda054872031b798e89e14340db4f647",
          "message": "update tests and changelog",
          "timestamp": "2024-12-05T10:43:23+08:00",
          "tree_id": "bb60dfaa3d8f80f609ebf9552d62b869dc28d863",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/b7a39c36bda054872031b798e89e14340db4f647"
        },
        "date": 1733366803592,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1469637,
            "unit": "ns/op\t  819445 B/op\t       5 allocs/op",
            "extra": "1232 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1469637,
            "unit": "ns/op",
            "extra": "1232 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819445,
            "unit": "B/op",
            "extra": "1232 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1232 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7358,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "162288 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7358,
            "unit": "ns/op",
            "extra": "162288 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "162288 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "162288 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.31,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14988336 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.31,
            "unit": "ns/op",
            "extra": "14988336 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14988336 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14988336 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 24400,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "47902 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 24400,
            "unit": "ns/op",
            "extra": "47902 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "47902 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "47902 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 256669,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4454 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 256669,
            "unit": "ns/op",
            "extra": "4454 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4454 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4454 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2889062,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "398 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2889062,
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
            "value": 41997321,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 41997321,
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
            "value": 10639,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "111856 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10639,
            "unit": "ns/op",
            "extra": "111856 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "111856 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "111856 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 9236621,
            "unit": "ns/op\t 4594789 B/op\t   75249 allocs/op",
            "extra": "128 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 9236621,
            "unit": "ns/op",
            "extra": "128 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594789,
            "unit": "B/op",
            "extra": "128 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75249,
            "unit": "allocs/op",
            "extra": "128 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8753,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "138079 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8753,
            "unit": "ns/op",
            "extra": "138079 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "138079 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "138079 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 945619,
            "unit": "ns/op\t  396909 B/op\t    6233 allocs/op",
            "extra": "1114 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 945619,
            "unit": "ns/op",
            "extra": "1114 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396909,
            "unit": "B/op",
            "extra": "1114 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6233,
            "unit": "allocs/op",
            "extra": "1114 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 12642,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "94929 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 12642,
            "unit": "ns/op",
            "extra": "94929 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "94929 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "94929 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 10652813,
            "unit": "ns/op\t 4981277 B/op\t   81258 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 10652813,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981277,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81258,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 247.9,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5088070 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 247.9,
            "unit": "ns/op",
            "extra": "5088070 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5088070 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5088070 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.7524,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.7524,
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
            "value": 313.1,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "3656401 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 313.1,
            "unit": "ns/op",
            "extra": "3656401 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "3656401 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "3656401 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 487.3,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2502640 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 487.3,
            "unit": "ns/op",
            "extra": "2502640 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2502640 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2502640 times\n4 procs"
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
          "id": "9df705d1a08bdf24d5366aeac27bb7e3b50cbbcf",
          "message": "update dbmodeUpdateStrategy",
          "timestamp": "2024-12-05T14:30:23+08:00",
          "tree_id": "a58f038fa450c35afc2a9f4876911308a9d3b91a",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/9df705d1a08bdf24d5366aeac27bb7e3b50cbbcf"
        },
        "date": 1733385220190,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1125421,
            "unit": "ns/op\t  819441 B/op\t       5 allocs/op",
            "extra": "1154 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1125421,
            "unit": "ns/op",
            "extra": "1154 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819441,
            "unit": "B/op",
            "extra": "1154 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1154 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 9069,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "113998 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 9069,
            "unit": "ns/op",
            "extra": "113998 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "113998 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "113998 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 80.21,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15137677 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 80.21,
            "unit": "ns/op",
            "extra": "15137677 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15137677 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15137677 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22300,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "55786 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22300,
            "unit": "ns/op",
            "extra": "55786 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "55786 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "55786 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 242688,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4849 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 242688,
            "unit": "ns/op",
            "extra": "4849 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4849 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4849 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2542669,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "399 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2542669,
            "unit": "ns/op",
            "extra": "399 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "399 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "399 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 34011358,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "31 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 34011358,
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
            "value": 9539,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "124432 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9539,
            "unit": "ns/op",
            "extra": "124432 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "124432 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "124432 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7337398,
            "unit": "ns/op\t 4594549 B/op\t   75248 allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7337398,
            "unit": "ns/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594549,
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
            "value": 7994,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "148576 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7994,
            "unit": "ns/op",
            "extra": "148576 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "148576 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "148576 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 832937,
            "unit": "ns/op\t  396382 B/op\t    6225 allocs/op",
            "extra": "1263 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 832937,
            "unit": "ns/op",
            "extra": "1263 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396382,
            "unit": "B/op",
            "extra": "1263 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6225,
            "unit": "allocs/op",
            "extra": "1263 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10546,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "110850 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10546,
            "unit": "ns/op",
            "extra": "110850 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "110850 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "110850 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8208885,
            "unit": "ns/op\t 4981221 B/op\t   81258 allocs/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8208885,
            "unit": "ns/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981221,
            "unit": "B/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81258,
            "unit": "allocs/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 222.2,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5295608 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 222.2,
            "unit": "ns/op",
            "extra": "5295608 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5295608 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5295608 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.5939,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.5939,
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
            "value": 290.6,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4222322 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 290.6,
            "unit": "ns/op",
            "extra": "4222322 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4222322 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4222322 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 440.2,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2713488 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 440.2,
            "unit": "ns/op",
            "extra": "2713488 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2713488 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2713488 times\n4 procs"
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
          "id": "3fb8cf31fba49aeec1b9fa59b77feb2b2474e7e1",
          "message": "fix lockups in UpdateStrategyDBMode.HandleEvent",
          "timestamp": "2024-12-05T17:04:12+08:00",
          "tree_id": "3e14750490cb08ceb31a4b6cf69a04d2ff1e44d7",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/3fb8cf31fba49aeec1b9fa59b77feb2b2474e7e1"
        },
        "date": 1733390113472,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 363247,
            "unit": "ns/op\t  819444 B/op\t       5 allocs/op",
            "extra": "3418 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 363247,
            "unit": "ns/op",
            "extra": "3418 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819444,
            "unit": "B/op",
            "extra": "3418 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "3418 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 9225,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "114302 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 9225,
            "unit": "ns/op",
            "extra": "114302 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "114302 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "114302 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.48,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15178873 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.48,
            "unit": "ns/op",
            "extra": "15178873 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15178873 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15178873 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22446,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53304 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22446,
            "unit": "ns/op",
            "extra": "53304 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53304 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53304 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 224497,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4915 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 224497,
            "unit": "ns/op",
            "extra": "4915 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
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
            "value": 2540536,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "462 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2540536,
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
            "value": 37975314,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 37975314,
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
            "value": 9791,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "122714 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9791,
            "unit": "ns/op",
            "extra": "122714 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "122714 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "122714 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7887377,
            "unit": "ns/op\t 4594460 B/op\t   75247 allocs/op",
            "extra": "156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7887377,
            "unit": "ns/op",
            "extra": "156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594460,
            "unit": "B/op",
            "extra": "156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75247,
            "unit": "allocs/op",
            "extra": "156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8286,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "136936 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8286,
            "unit": "ns/op",
            "extra": "136936 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "136936 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "136936 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 869078,
            "unit": "ns/op\t  396500 B/op\t    6227 allocs/op",
            "extra": "1230 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 869078,
            "unit": "ns/op",
            "extra": "1230 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396500,
            "unit": "B/op",
            "extra": "1230 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6227,
            "unit": "allocs/op",
            "extra": "1230 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11160,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "106942 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11160,
            "unit": "ns/op",
            "extra": "106942 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "106942 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "106942 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8722920,
            "unit": "ns/op\t 4980988 B/op\t   81257 allocs/op",
            "extra": "136 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8722920,
            "unit": "ns/op",
            "extra": "136 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4980988,
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
            "value": 226.2,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5317315 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 226.2,
            "unit": "ns/op",
            "extra": "5317315 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5317315 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5317315 times\n4 procs"
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
            "value": 294.3,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4020388 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 294.3,
            "unit": "ns/op",
            "extra": "4020388 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4020388 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4020388 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 526.5,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2684917 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 526.5,
            "unit": "ns/op",
            "extra": "2684917 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2684917 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2684917 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "charly.molter@konghq.com",
            "name": "Charly Molter",
            "username": "lahabana"
          },
          "committer": {
            "email": "charly.molter@konghq.com",
            "name": "Charly Molter",
            "username": "lahabana"
          },
          "distinct": true,
          "id": "8b68a7c32e6c6908931ae264389e7f126df7992b",
          "message": "docs(ISSUE_TEMPLATE): add epic template\n\nWe've been in the past few releases working out of umbrella issues\nfor complex tasks (.i.e: epic).\n\nWe've also in retros mentionned the need to do some things\nsystematically (docs, user acceptance testing...).\n\nThis issue template is an attempt at codifying this\n\nSigned-off-by: Charly Molter <charly.molter@konghq.com>",
          "timestamp": "2024-12-05T10:15:53+01:00",
          "tree_id": "36e1abb41234185cbdd9dd1ae7e6a9ab69088205",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/8b68a7c32e6c6908931ae264389e7f126df7992b"
        },
        "date": 1733390354563,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1262894,
            "unit": "ns/op\t  819446 B/op\t       5 allocs/op",
            "extra": "1243 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1262894,
            "unit": "ns/op",
            "extra": "1243 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819446,
            "unit": "B/op",
            "extra": "1243 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1243 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7046,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "171644 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7046,
            "unit": "ns/op",
            "extra": "171644 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "171644 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "171644 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 83.78,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15168975 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 83.78,
            "unit": "ns/op",
            "extra": "15168975 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15168975 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15168975 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22462,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52933 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22462,
            "unit": "ns/op",
            "extra": "52933 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52933 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52933 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 220289,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5612 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 220289,
            "unit": "ns/op",
            "extra": "5612 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5612 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5612 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2568669,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "423 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2568669,
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
            "value": 31161447,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "39 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 31161447,
            "unit": "ns/op",
            "extra": "39 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
            "unit": "B/op",
            "extra": "39 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "39 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9624,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "119960 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9624,
            "unit": "ns/op",
            "extra": "119960 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "119960 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "119960 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7599511,
            "unit": "ns/op\t 4594284 B/op\t   75247 allocs/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7599511,
            "unit": "ns/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594284,
            "unit": "B/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75247,
            "unit": "allocs/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8012,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "149709 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8012,
            "unit": "ns/op",
            "extra": "149709 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "149709 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "149709 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 857667,
            "unit": "ns/op\t  396457 B/op\t    6226 allocs/op",
            "extra": "1243 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 857667,
            "unit": "ns/op",
            "extra": "1243 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396457,
            "unit": "B/op",
            "extra": "1243 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6226,
            "unit": "allocs/op",
            "extra": "1243 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11038,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "106839 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11038,
            "unit": "ns/op",
            "extra": "106839 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "106839 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "106839 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8459628,
            "unit": "ns/op\t 4980805 B/op\t   81256 allocs/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8459628,
            "unit": "ns/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4980805,
            "unit": "B/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81256,
            "unit": "allocs/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 226,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5226253 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 226,
            "unit": "ns/op",
            "extra": "5226253 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5226253 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5226253 times\n4 procs"
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
            "extra": "4012114 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 297.9,
            "unit": "ns/op",
            "extra": "4012114 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4012114 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4012114 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 445.9,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2714276 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 445.9,
            "unit": "ns/op",
            "extra": "2714276 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2714276 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2714276 times\n4 procs"
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
          "id": "0a51c00ac115264ac11e51068c06ff72f9d86aa3",
          "message": "fix lockups in UpdateStrategyDBMode.HandleEvent",
          "timestamp": "2024-12-05T17:20:46+08:00",
          "tree_id": "3e14750490cb08ceb31a4b6cf69a04d2ff1e44d7",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/0a51c00ac115264ac11e51068c06ff72f9d86aa3"
        },
        "date": 1733390549158,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 401456,
            "unit": "ns/op\t  819443 B/op\t       5 allocs/op",
            "extra": "2928 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 401456,
            "unit": "ns/op",
            "extra": "2928 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819443,
            "unit": "B/op",
            "extra": "2928 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "2928 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7468,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "157616 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7468,
            "unit": "ns/op",
            "extra": "157616 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "157616 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "157616 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 78.92,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15122404 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 78.92,
            "unit": "ns/op",
            "extra": "15122404 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15122404 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15122404 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22473,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53469 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22473,
            "unit": "ns/op",
            "extra": "53469 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53469 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53469 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 228870,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5109 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 228870,
            "unit": "ns/op",
            "extra": "5109 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
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
            "value": 3057473,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "452 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 3057473,
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
            "value": 39637906,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 39637906,
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
            "value": 9565,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "118221 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9565,
            "unit": "ns/op",
            "extra": "118221 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "118221 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "118221 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7480155,
            "unit": "ns/op\t 4594667 B/op\t   75248 allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7480155,
            "unit": "ns/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594667,
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
            "value": 8028,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "146497 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8028,
            "unit": "ns/op",
            "extra": "146497 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "146497 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "146497 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 865140,
            "unit": "ns/op\t  396749 B/op\t    6231 allocs/op",
            "extra": "1156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 865140,
            "unit": "ns/op",
            "extra": "1156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396749,
            "unit": "B/op",
            "extra": "1156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6231,
            "unit": "allocs/op",
            "extra": "1156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11061,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "108928 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11061,
            "unit": "ns/op",
            "extra": "108928 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "108928 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "108928 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8817868,
            "unit": "ns/op\t 4980994 B/op\t   81257 allocs/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8817868,
            "unit": "ns/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4980994,
            "unit": "B/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81257,
            "unit": "allocs/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 229.5,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5255781 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 229.5,
            "unit": "ns/op",
            "extra": "5255781 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5255781 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5255781 times\n4 procs"
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
            "value": 347,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4043534 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 347,
            "unit": "ns/op",
            "extra": "4043534 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4043534 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4043534 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 442.9,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2474263 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 442.9,
            "unit": "ns/op",
            "extra": "2474263 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2474263 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2474263 times\n4 procs"
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
          "id": "c67aa264df92d3c2b1c50f1ce0c9b37c57d0075e",
          "message": "fix lockups in UpdateStrategyDBMode.HandleEvent",
          "timestamp": "2024-12-05T17:27:15+08:00",
          "tree_id": "bca22b8ef8d962f59e030bbe87c924c22ab28929",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/c67aa264df92d3c2b1c50f1ce0c9b37c57d0075e"
        },
        "date": 1733390930794,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 413399,
            "unit": "ns/op\t  819443 B/op\t       5 allocs/op",
            "extra": "2942 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 413399,
            "unit": "ns/op",
            "extra": "2942 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819443,
            "unit": "B/op",
            "extra": "2942 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "2942 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7554,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "148809 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7554,
            "unit": "ns/op",
            "extra": "148809 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "148809 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "148809 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.34,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15095619 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.34,
            "unit": "ns/op",
            "extra": "15095619 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15095619 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15095619 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22669,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52603 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22669,
            "unit": "ns/op",
            "extra": "52603 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52603 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52603 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 285491,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5088 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 285491,
            "unit": "ns/op",
            "extra": "5088 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5088 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5088 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2675923,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "394 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2675923,
            "unit": "ns/op",
            "extra": "394 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "394 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "394 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 35700005,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 35700005,
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
            "value": 10230,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "117046 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10230,
            "unit": "ns/op",
            "extra": "117046 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "117046 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "117046 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 8207762,
            "unit": "ns/op\t 4594439 B/op\t   75247 allocs/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 8207762,
            "unit": "ns/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594439,
            "unit": "B/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75247,
            "unit": "allocs/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8446,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "145411 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8446,
            "unit": "ns/op",
            "extra": "145411 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "145411 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "145411 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 894267,
            "unit": "ns/op\t  396517 B/op\t    6227 allocs/op",
            "extra": "1220 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 894267,
            "unit": "ns/op",
            "extra": "1220 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396517,
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
            "value": 11934,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "100617 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11934,
            "unit": "ns/op",
            "extra": "100617 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "100617 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "100617 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8911823,
            "unit": "ns/op\t 4981543 B/op\t   81259 allocs/op",
            "extra": "130 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8911823,
            "unit": "ns/op",
            "extra": "130 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981543,
            "unit": "B/op",
            "extra": "130 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81259,
            "unit": "allocs/op",
            "extra": "130 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 227.6,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5282218 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 227.6,
            "unit": "ns/op",
            "extra": "5282218 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5282218 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5282218 times\n4 procs"
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
            "value": 295,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4038400 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 295,
            "unit": "ns/op",
            "extra": "4038400 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4038400 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4038400 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 450.9,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2714155 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 450.9,
            "unit": "ns/op",
            "extra": "2714155 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2714155 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2714155 times\n4 procs"
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
            "email": "jakub.warczarek@konghq.com",
            "name": "Jakub Warczarek",
            "username": "programmer04"
          },
          "distinct": true,
          "id": "6094a27623adc697edab68990550c1cbab87f575",
          "message": "feat((TLS|TCP|UDP|GRPC)Route): propagate to Kong tag k8s-named-route-rule",
          "timestamp": "2024-12-05T10:36:38+01:00",
          "tree_id": "6ecd96822609e1bc391e889c08d62aca1461a941",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/6094a27623adc697edab68990550c1cbab87f575"
        },
        "date": 1733391601335,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1237004,
            "unit": "ns/op\t  819446 B/op\t       5 allocs/op",
            "extra": "1081 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1237004,
            "unit": "ns/op",
            "extra": "1081 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819446,
            "unit": "B/op",
            "extra": "1081 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1081 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 11295,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "99531 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 11295,
            "unit": "ns/op",
            "extra": "99531 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "99531 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "99531 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.25,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15223482 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.25,
            "unit": "ns/op",
            "extra": "15223482 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15223482 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15223482 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22405,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52971 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22405,
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
            "value": 219929,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5377 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 219929,
            "unit": "ns/op",
            "extra": "5377 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5377 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5377 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2545416,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "434 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2545416,
            "unit": "ns/op",
            "extra": "434 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "434 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "434 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 38015960,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 38015960,
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
            "value": 9603,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "123174 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9603,
            "unit": "ns/op",
            "extra": "123174 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "123174 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "123174 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7589837,
            "unit": "ns/op\t 4594656 B/op\t   75248 allocs/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7589837,
            "unit": "ns/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594656,
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
            "value": 7983,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "148694 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7983,
            "unit": "ns/op",
            "extra": "148694 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "148694 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "148694 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 856177,
            "unit": "ns/op\t  396417 B/op\t    6225 allocs/op",
            "extra": "1254 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 856177,
            "unit": "ns/op",
            "extra": "1254 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396417,
            "unit": "B/op",
            "extra": "1254 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6225,
            "unit": "allocs/op",
            "extra": "1254 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11041,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "108328 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11041,
            "unit": "ns/op",
            "extra": "108328 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "108328 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "108328 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8488426,
            "unit": "ns/op\t 4981247 B/op\t   81258 allocs/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8488426,
            "unit": "ns/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981247,
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
            "value": 255.2,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5201696 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 255.2,
            "unit": "ns/op",
            "extra": "5201696 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5201696 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5201696 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6343,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6343,
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
            "extra": "4065721 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 296.6,
            "unit": "ns/op",
            "extra": "4065721 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4065721 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4065721 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 448.1,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2652525 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 448.1,
            "unit": "ns/op",
            "extra": "2652525 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2652525 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2652525 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "charly.molter@konghq.com",
            "name": "Charly Molter",
            "username": "lahabana"
          },
          "committer": {
            "email": "charly.molter@konghq.com",
            "name": "Charly Molter",
            "username": "lahabana"
          },
          "distinct": true,
          "id": "8b68a7c32e6c6908931ae264389e7f126df7992b",
          "message": "docs(ISSUE_TEMPLATE): add epic template\n\nWe've been in the past few releases working out of umbrella issues\nfor complex tasks (.i.e: epic).\n\nWe've also in retros mentionned the need to do some things\nsystematically (docs, user acceptance testing...).\n\nThis issue template is an attempt at codifying this\n\nSigned-off-by: Charly Molter <charly.molter@konghq.com>",
          "timestamp": "2024-12-05T10:15:53+01:00",
          "tree_id": "36e1abb41234185cbdd9dd1ae7e6a9ab69088205",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/8b68a7c32e6c6908931ae264389e7f126df7992b"
        },
        "date": 1733391733638,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1015005,
            "unit": "ns/op\t  819442 B/op\t       5 allocs/op",
            "extra": "998 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1015005,
            "unit": "ns/op",
            "extra": "998 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819442,
            "unit": "B/op",
            "extra": "998 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "998 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7287,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "169471 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7287,
            "unit": "ns/op",
            "extra": "169471 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "169471 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "169471 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.73,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15131203 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.73,
            "unit": "ns/op",
            "extra": "15131203 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15131203 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15131203 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22885,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53068 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22885,
            "unit": "ns/op",
            "extra": "53068 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53068 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53068 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 220576,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5460 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 220576,
            "unit": "ns/op",
            "extra": "5460 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5460 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5460 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2618104,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "415 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2618104,
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
            "value": 39117487,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 39117487,
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
            "value": 9536,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "124875 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9536,
            "unit": "ns/op",
            "extra": "124875 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "124875 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "124875 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7410603,
            "unit": "ns/op\t 4594565 B/op\t   75248 allocs/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7410603,
            "unit": "ns/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594565,
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
            "value": 7935,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "148042 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7935,
            "unit": "ns/op",
            "extra": "148042 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "148042 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "148042 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 847561,
            "unit": "ns/op\t  396378 B/op\t    6225 allocs/op",
            "extra": "1269 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 847561,
            "unit": "ns/op",
            "extra": "1269 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396378,
            "unit": "B/op",
            "extra": "1269 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6225,
            "unit": "allocs/op",
            "extra": "1269 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10944,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "107904 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10944,
            "unit": "ns/op",
            "extra": "107904 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "107904 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "107904 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8480063,
            "unit": "ns/op\t 4981174 B/op\t   81258 allocs/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8480063,
            "unit": "ns/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981174,
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
            "value": 227.1,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5339337 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 227.1,
            "unit": "ns/op",
            "extra": "5339337 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5339337 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5339337 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6521,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6521,
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
            "value": 295.2,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "3919466 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 295.2,
            "unit": "ns/op",
            "extra": "3919466 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "3919466 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "3919466 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 448.2,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2684508 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 448.2,
            "unit": "ns/op",
            "extra": "2684508 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2684508 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2684508 times\n4 procs"
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
            "email": "jakub.warczarek@konghq.com",
            "name": "Jakub Warczarek",
            "username": "programmer04"
          },
          "distinct": true,
          "id": "872db9219ecd718478d66c7ff116481108963b12",
          "message": "feat((TLS|TCP|UDP|GRPC)Route): propagate to Kong tag k8s-named-route-rule",
          "timestamp": "2024-12-05T10:42:36+01:00",
          "tree_id": "b2935311cef45dcbddb3ffd4fd5577124a7cbc70",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/872db9219ecd718478d66c7ff116481108963b12"
        },
        "date": 1733391939840,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1196878,
            "unit": "ns/op\t  819442 B/op\t       5 allocs/op",
            "extra": "963 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1196878,
            "unit": "ns/op",
            "extra": "963 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819442,
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
            "value": 9825,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "103368 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 9825,
            "unit": "ns/op",
            "extra": "103368 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "103368 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "103368 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 78.88,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15147386 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 78.88,
            "unit": "ns/op",
            "extra": "15147386 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15147386 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15147386 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22296,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53222 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22296,
            "unit": "ns/op",
            "extra": "53222 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53222 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53222 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 224435,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4800 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 224435,
            "unit": "ns/op",
            "extra": "4800 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4800 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4800 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2778152,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "381 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2778152,
            "unit": "ns/op",
            "extra": "381 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "381 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "381 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 34463708,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 34463708,
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
            "value": 9698,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "121274 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9698,
            "unit": "ns/op",
            "extra": "121274 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "121274 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "121274 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7520093,
            "unit": "ns/op\t 4594452 B/op\t   75247 allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7520093,
            "unit": "ns/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594452,
            "unit": "B/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75247,
            "unit": "allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8066,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "146145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8066,
            "unit": "ns/op",
            "extra": "146145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "146145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "146145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 865012,
            "unit": "ns/op\t  396434 B/op\t    6226 allocs/op",
            "extra": "1248 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 865012,
            "unit": "ns/op",
            "extra": "1248 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396434,
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
            "value": 11084,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "105024 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11084,
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
            "value": 8510480,
            "unit": "ns/op\t 4980888 B/op\t   81257 allocs/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8510480,
            "unit": "ns/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4980888,
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
            "value": 226.9,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5298373 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 226.9,
            "unit": "ns/op",
            "extra": "5298373 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5298373 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5298373 times\n4 procs"
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
            "value": 293.9,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4072116 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 293.9,
            "unit": "ns/op",
            "extra": "4072116 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4072116 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4072116 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 446,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2688490 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 446,
            "unit": "ns/op",
            "extra": "2688490 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2688490 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2688490 times\n4 procs"
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
            "email": "jakub.warczarek@konghq.com",
            "name": "Jakub Warczarek",
            "username": "programmer04"
          },
          "distinct": true,
          "id": "caf26a4a5dda17d121321917322a2ff9506a5cb5",
          "message": "feat((TLS|TCP|UDP|GRPC)Route): propagate to Kong tag k8s-named-route-rule",
          "timestamp": "2024-12-05T10:54:37+01:00",
          "tree_id": "0d37dabf38ecde0b6fa6ca42a000e418e72a8d30",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/caf26a4a5dda17d121321917322a2ff9506a5cb5"
        },
        "date": 1733392660484,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1233264,
            "unit": "ns/op\t  819448 B/op\t       5 allocs/op",
            "extra": "885 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1233264,
            "unit": "ns/op",
            "extra": "885 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819448,
            "unit": "B/op",
            "extra": "885 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "885 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 16984,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "58971 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 16984,
            "unit": "ns/op",
            "extra": "58971 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "58971 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "58971 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 82.03,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13585507 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 82.03,
            "unit": "ns/op",
            "extra": "13585507 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13585507 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13585507 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22504,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53318 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22504,
            "unit": "ns/op",
            "extra": "53318 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53318 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53318 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 234613,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5299 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 234613,
            "unit": "ns/op",
            "extra": "5299 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5299 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5299 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2535900,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "438 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2535900,
            "unit": "ns/op",
            "extra": "438 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "438 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "438 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 38250244,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 38250244,
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
            "value": 9577,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "125788 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9577,
            "unit": "ns/op",
            "extra": "125788 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "125788 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "125788 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7595568,
            "unit": "ns/op\t 4594736 B/op\t   75248 allocs/op",
            "extra": "154 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7595568,
            "unit": "ns/op",
            "extra": "154 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594736,
            "unit": "B/op",
            "extra": "154 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75248,
            "unit": "allocs/op",
            "extra": "154 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8003,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "144543 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8003,
            "unit": "ns/op",
            "extra": "144543 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "144543 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "144543 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 856504,
            "unit": "ns/op\t  396451 B/op\t    6226 allocs/op",
            "extra": "1246 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 856504,
            "unit": "ns/op",
            "extra": "1246 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396451,
            "unit": "B/op",
            "extra": "1246 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6226,
            "unit": "allocs/op",
            "extra": "1246 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11472,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "109032 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11472,
            "unit": "ns/op",
            "extra": "109032 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "109032 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "109032 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8431520,
            "unit": "ns/op\t 4981124 B/op\t   81258 allocs/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8431520,
            "unit": "ns/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981124,
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
            "value": 227.2,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5261131 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 227.2,
            "unit": "ns/op",
            "extra": "5261131 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5261131 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5261131 times\n4 procs"
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
            "value": 295.9,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4071847 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 295.9,
            "unit": "ns/op",
            "extra": "4071847 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4071847 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4071847 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 449.2,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2682326 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 449.2,
            "unit": "ns/op",
            "extra": "2682326 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2682326 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2682326 times\n4 procs"
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
          "id": "d683567f89625fde90a78c2d453d15e063bcc2c0",
          "message": "feat((TLS|TCP|UDP|GRPC)Route): propagate to Kong tag k8s-named-route-rule (#6780)",
          "timestamp": "2024-12-05T18:21:25+08:00",
          "tree_id": "0d37dabf38ecde0b6fa6ca42a000e418e72a8d30",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/d683567f89625fde90a78c2d453d15e063bcc2c0"
        },
        "date": 1733394271561,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1240110,
            "unit": "ns/op\t  819441 B/op\t       5 allocs/op",
            "extra": "964 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1240110,
            "unit": "ns/op",
            "extra": "964 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819441,
            "unit": "B/op",
            "extra": "964 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "964 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 6895,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "171660 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 6895,
            "unit": "ns/op",
            "extra": "171660 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "171660 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "171660 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 78.97,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15155229 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 78.97,
            "unit": "ns/op",
            "extra": "15155229 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15155229 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15155229 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22500,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53359 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22500,
            "unit": "ns/op",
            "extra": "53359 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53359 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53359 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 227262,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4875 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 227262,
            "unit": "ns/op",
            "extra": "4875 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4875 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4875 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2891841,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "505 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2891841,
            "unit": "ns/op",
            "extra": "505 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
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
            "value": 38221914,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "27 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 38221914,
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
            "value": 9592,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "121545 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9592,
            "unit": "ns/op",
            "extra": "121545 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "121545 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "121545 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7613964,
            "unit": "ns/op\t 4594365 B/op\t   75247 allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7613964,
            "unit": "ns/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594365,
            "unit": "B/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75247,
            "unit": "allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8146,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "148741 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8146,
            "unit": "ns/op",
            "extra": "148741 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "148741 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "148741 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 879299,
            "unit": "ns/op\t  396528 B/op\t    6227 allocs/op",
            "extra": "1222 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 879299,
            "unit": "ns/op",
            "extra": "1222 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396528,
            "unit": "B/op",
            "extra": "1222 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6227,
            "unit": "allocs/op",
            "extra": "1222 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11457,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "103556 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11457,
            "unit": "ns/op",
            "extra": "103556 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "103556 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "103556 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 9261712,
            "unit": "ns/op\t 4980950 B/op\t   81257 allocs/op",
            "extra": "128 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 9261712,
            "unit": "ns/op",
            "extra": "128 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4980950,
            "unit": "B/op",
            "extra": "128 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81257,
            "unit": "allocs/op",
            "extra": "128 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 229.1,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5113872 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 229.1,
            "unit": "ns/op",
            "extra": "5113872 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5113872 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5113872 times\n4 procs"
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
            "value": 298.8,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4051878 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 298.8,
            "unit": "ns/op",
            "extra": "4051878 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4051878 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4051878 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 455.5,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2652808 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 455.5,
            "unit": "ns/op",
            "extra": "2652808 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2652808 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2652808 times\n4 procs"
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
          "id": "0412408d211134f2c52c8bdeb591aad6bf148df8",
          "message": "update tests and changelog",
          "timestamp": "2024-12-05T18:23:23+08:00",
          "tree_id": "62edb462f5aa02c926045d2b5dc849a7f3043a31",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/0412408d211134f2c52c8bdeb591aad6bf148df8"
        },
        "date": 1733394386790,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1176346,
            "unit": "ns/op\t  819453 B/op\t       5 allocs/op",
            "extra": "962 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1176346,
            "unit": "ns/op",
            "extra": "962 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819453,
            "unit": "B/op",
            "extra": "962 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "962 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 6892,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "173910 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 6892,
            "unit": "ns/op",
            "extra": "173910 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "173910 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "173910 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.14,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15124708 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.14,
            "unit": "ns/op",
            "extra": "15124708 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15124708 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15124708 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22599,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53745 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22599,
            "unit": "ns/op",
            "extra": "53745 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53745 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53745 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 221269,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5043 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 221269,
            "unit": "ns/op",
            "extra": "5043 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5043 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5043 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2981303,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "429 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2981303,
            "unit": "ns/op",
            "extra": "429 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
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
            "value": 36566180,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 36566180,
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
            "value": 9606,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "119787 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9606,
            "unit": "ns/op",
            "extra": "119787 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "119787 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "119787 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7363438,
            "unit": "ns/op\t 4594409 B/op\t   75247 allocs/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7363438,
            "unit": "ns/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594409,
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
            "value": 7899,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "152952 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7899,
            "unit": "ns/op",
            "extra": "152952 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "152952 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "152952 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 863560,
            "unit": "ns/op\t  396451 B/op\t    6226 allocs/op",
            "extra": "1246 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 863560,
            "unit": "ns/op",
            "extra": "1246 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396451,
            "unit": "B/op",
            "extra": "1246 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6226,
            "unit": "allocs/op",
            "extra": "1246 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10944,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "109693 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10944,
            "unit": "ns/op",
            "extra": "109693 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "109693 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "109693 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8468453,
            "unit": "ns/op\t 4981305 B/op\t   81258 allocs/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8468453,
            "unit": "ns/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981305,
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
            "value": 284.4,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5298576 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 284.4,
            "unit": "ns/op",
            "extra": "5298576 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5298576 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5298576 times\n4 procs"
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
            "value": 297.6,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4048652 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 297.6,
            "unit": "ns/op",
            "extra": "4048652 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4048652 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4048652 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 451.4,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2651180 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 451.4,
            "unit": "ns/op",
            "extra": "2651180 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2651180 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2651180 times\n4 procs"
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
          "id": "78fff3e89a489cdbbfc6637f6affcff24d730ada",
          "message": "fix lockups in UpdateStrategyDBMode.HandleEvent",
          "timestamp": "2024-12-05T18:24:49+08:00",
          "tree_id": "729f22ba23b8bea267c1e9733339295b356a80fb",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/78fff3e89a489cdbbfc6637f6affcff24d730ada"
        },
        "date": 1733394388283,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 360561,
            "unit": "ns/op\t  819444 B/op\t       5 allocs/op",
            "extra": "2880 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 360561,
            "unit": "ns/op",
            "extra": "2880 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819444,
            "unit": "B/op",
            "extra": "2880 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "2880 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 6907,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "170607 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 6907,
            "unit": "ns/op",
            "extra": "170607 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "170607 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "170607 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.27,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15070378 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.27,
            "unit": "ns/op",
            "extra": "15070378 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15070378 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15070378 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22434,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52773 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22434,
            "unit": "ns/op",
            "extra": "52773 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52773 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52773 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 221970,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4550 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 221970,
            "unit": "ns/op",
            "extra": "4550 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4550 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4550 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2476966,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "442 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2476966,
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
            "value": 35882933,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 35882933,
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
            "value": 9680,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "122200 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9680,
            "unit": "ns/op",
            "extra": "122200 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "122200 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "122200 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7610880,
            "unit": "ns/op\t 4594645 B/op\t   75248 allocs/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7610880,
            "unit": "ns/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594645,
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
            "value": 8119,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "145538 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8119,
            "unit": "ns/op",
            "extra": "145538 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "145538 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "145538 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 869757,
            "unit": "ns/op\t  396663 B/op\t    6229 allocs/op",
            "extra": "1182 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 869757,
            "unit": "ns/op",
            "extra": "1182 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396663,
            "unit": "B/op",
            "extra": "1182 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6229,
            "unit": "allocs/op",
            "extra": "1182 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11229,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "104862 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11229,
            "unit": "ns/op",
            "extra": "104862 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "104862 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "104862 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8679560,
            "unit": "ns/op\t 4981033 B/op\t   81257 allocs/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8679560,
            "unit": "ns/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981033,
            "unit": "B/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81257,
            "unit": "allocs/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 227.4,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5317770 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 227.4,
            "unit": "ns/op",
            "extra": "5317770 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5317770 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5317770 times\n4 procs"
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
            "value": 296.6,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4035624 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 296.6,
            "unit": "ns/op",
            "extra": "4035624 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4035624 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4035624 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 447.2,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2662814 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 447.2,
            "unit": "ns/op",
            "extra": "2662814 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2662814 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2662814 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "charly.molter@konghq.com",
            "name": "Charly Molter",
            "username": "lahabana"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "94ed805528eea668333d709cd6781e8370734c14",
          "message": "docs(ISSUE_TEMPLATE): add epic template (#6779)\n\nWe've been in the past few releases working out of umbrella issues\r\nfor complex tasks (.i.e: epic).\r\n\r\nWe've also in retros mentionned the need to do some things\r\nsystematically (docs, user acceptance testing...).\r\n\r\nThis issue template is an attempt at codifying this\r\n\r\nSigned-off-by: Charly Molter <charly.molter@konghq.com>",
          "timestamp": "2024-12-05T18:23:49+08:00",
          "tree_id": "9e7092004256135297ed29234268fd2f60c6894a",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/94ed805528eea668333d709cd6781e8370734c14"
        },
        "date": 1733394405979,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1284038,
            "unit": "ns/op\t  819441 B/op\t       5 allocs/op",
            "extra": "789 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1284038,
            "unit": "ns/op",
            "extra": "789 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819441,
            "unit": "B/op",
            "extra": "789 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "789 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7612,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "162728 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7612,
            "unit": "ns/op",
            "extra": "162728 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "162728 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "162728 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.35,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15138938 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.35,
            "unit": "ns/op",
            "extra": "15138938 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15138938 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15138938 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22720,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53278 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22720,
            "unit": "ns/op",
            "extra": "53278 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53278 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53278 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 234779,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5109 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 234779,
            "unit": "ns/op",
            "extra": "5109 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
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
            "value": 2886218,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "369 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2886218,
            "unit": "ns/op",
            "extra": "369 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "369 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "369 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 43762118,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 43762118,
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
            "value": 10001,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "120895 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10001,
            "unit": "ns/op",
            "extra": "120895 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "120895 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "120895 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 8006757,
            "unit": "ns/op\t 4594407 B/op\t   75247 allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 8006757,
            "unit": "ns/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594407,
            "unit": "B/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75247,
            "unit": "allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8309,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "142860 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8309,
            "unit": "ns/op",
            "extra": "142860 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "142860 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "142860 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 886831,
            "unit": "ns/op\t  396565 B/op\t    6228 allocs/op",
            "extra": "1210 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 886831,
            "unit": "ns/op",
            "extra": "1210 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396565,
            "unit": "B/op",
            "extra": "1210 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6228,
            "unit": "allocs/op",
            "extra": "1210 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11421,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "104078 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11421,
            "unit": "ns/op",
            "extra": "104078 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "104078 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "104078 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 9533002,
            "unit": "ns/op\t 4981282 B/op\t   81258 allocs/op",
            "extra": "129 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 9533002,
            "unit": "ns/op",
            "extra": "129 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981282,
            "unit": "B/op",
            "extra": "129 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81258,
            "unit": "allocs/op",
            "extra": "129 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 227.9,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5281881 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 227.9,
            "unit": "ns/op",
            "extra": "5281881 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5281881 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5281881 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6274,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6274,
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
            "value": 297.3,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4044298 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 297.3,
            "unit": "ns/op",
            "extra": "4044298 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4044298 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4044298 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 452.1,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2653352 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 452.1,
            "unit": "ns/op",
            "extra": "2653352 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2653352 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2653352 times\n4 procs"
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
            "name": "Tao Yi",
            "username": "randmonkey"
          },
          "distinct": true,
          "id": "1503cdedf0b9e9dd16e8cad253450115d93f8ed7",
          "message": "fix lockups in UpdateStrategyDBMode.HandleEvent",
          "timestamp": "2024-12-05T18:28:37+08:00",
          "tree_id": "6eae0360454009ca9711476d90aaf86bd9c340ee",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/1503cdedf0b9e9dd16e8cad253450115d93f8ed7"
        },
        "date": 1733394608132,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 360403,
            "unit": "ns/op\t  819444 B/op\t       5 allocs/op",
            "extra": "3367 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 360403,
            "unit": "ns/op",
            "extra": "3367 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819444,
            "unit": "B/op",
            "extra": "3367 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "3367 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7381,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "162376 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7381,
            "unit": "ns/op",
            "extra": "162376 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "162376 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "162376 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.19,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15125484 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.19,
            "unit": "ns/op",
            "extra": "15125484 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15125484 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15125484 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22682,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52507 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22682,
            "unit": "ns/op",
            "extra": "52507 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52507 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52507 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 237821,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4452 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 237821,
            "unit": "ns/op",
            "extra": "4452 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4452 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4452 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2601935,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "452 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2601935,
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
            "value": 34938199,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 34938199,
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
            "value": 9626,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "124819 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9626,
            "unit": "ns/op",
            "extra": "124819 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "124819 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "124819 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7654667,
            "unit": "ns/op\t 4594519 B/op\t   75248 allocs/op",
            "extra": "156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7654667,
            "unit": "ns/op",
            "extra": "156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594519,
            "unit": "B/op",
            "extra": "156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75248,
            "unit": "allocs/op",
            "extra": "156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8087,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "145346 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8087,
            "unit": "ns/op",
            "extra": "145346 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "145346 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "145346 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 868736,
            "unit": "ns/op\t  396462 B/op\t    6226 allocs/op",
            "extra": "1240 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 868736,
            "unit": "ns/op",
            "extra": "1240 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396462,
            "unit": "B/op",
            "extra": "1240 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6226,
            "unit": "allocs/op",
            "extra": "1240 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11138,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "107088 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11138,
            "unit": "ns/op",
            "extra": "107088 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "107088 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "107088 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 9803245,
            "unit": "ns/op\t 4981387 B/op\t   81258 allocs/op",
            "extra": "130 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 9803245,
            "unit": "ns/op",
            "extra": "130 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981387,
            "unit": "B/op",
            "extra": "130 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81258,
            "unit": "allocs/op",
            "extra": "130 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 228.3,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5253964 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 228.3,
            "unit": "ns/op",
            "extra": "5253964 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5253964 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5253964 times\n4 procs"
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
            "value": 295.9,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4018924 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 295.9,
            "unit": "ns/op",
            "extra": "4018924 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4018924 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4018924 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 447.4,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2677482 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 447.4,
            "unit": "ns/op",
            "extra": "2677482 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2677482 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2677482 times\n4 procs"
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
          "id": "250d107180e424976002f6b52fb1014c09ce4fd1",
          "message": "continue\n\nSigned-off-by: Mattia Lavacca <lavacca.mattia@gmail.com>",
          "timestamp": "2024-12-05T11:58:34+01:00",
          "tree_id": "e7d8cd53d07c5d2f39b047acf41a884aee775b6e",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/250d107180e424976002f6b52fb1014c09ce4fd1"
        },
        "date": 1733396518920,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1344496,
            "unit": "ns/op\t  819441 B/op\t       5 allocs/op",
            "extra": "907 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1344496,
            "unit": "ns/op",
            "extra": "907 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819441,
            "unit": "B/op",
            "extra": "907 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "907 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7200,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "162436 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7200,
            "unit": "ns/op",
            "extra": "162436 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "162436 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "162436 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15120429 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.2,
            "unit": "ns/op",
            "extra": "15120429 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15120429 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15120429 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22462,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53570 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22462,
            "unit": "ns/op",
            "extra": "53570 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53570 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53570 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 218588,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5449 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 218588,
            "unit": "ns/op",
            "extra": "5449 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5449 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5449 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2791945,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "393 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2791945,
            "unit": "ns/op",
            "extra": "393 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "393 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "393 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 36121349,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 36121349,
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
            "value": 9618,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "120928 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9618,
            "unit": "ns/op",
            "extra": "120928 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "120928 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "120928 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7417167,
            "unit": "ns/op\t 4594313 B/op\t   75247 allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7417167,
            "unit": "ns/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594313,
            "unit": "B/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75247,
            "unit": "allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8026,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "147728 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8026,
            "unit": "ns/op",
            "extra": "147728 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "147728 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "147728 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 871451,
            "unit": "ns/op\t  396501 B/op\t    6227 allocs/op",
            "extra": "1227 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 871451,
            "unit": "ns/op",
            "extra": "1227 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396501,
            "unit": "B/op",
            "extra": "1227 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6227,
            "unit": "allocs/op",
            "extra": "1227 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11270,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "105260 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11270,
            "unit": "ns/op",
            "extra": "105260 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "105260 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "105260 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8637232,
            "unit": "ns/op\t 4981165 B/op\t   81258 allocs/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8637232,
            "unit": "ns/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981165,
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
            "value": 39.14,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "30700274 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 39.14,
            "unit": "ns/op",
            "extra": "30700274 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "30700274 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "30700274 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6216,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6216,
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
            "value": 296.3,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4037397 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 296.3,
            "unit": "ns/op",
            "extra": "4037397 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4037397 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4037397 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 449.8,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2658582 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 449.8,
            "unit": "ns/op",
            "extra": "2658582 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2658582 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2658582 times\n4 procs"
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
          "id": "343d922d159741870bb51f96b25a64a0ea740576",
          "message": "continue\n\nSigned-off-by: Mattia Lavacca <lavacca.mattia@gmail.com>",
          "timestamp": "2024-12-05T11:59:00+01:00",
          "tree_id": "a6dac5137c78750f413844b43f936429ad716980",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/343d922d159741870bb51f96b25a64a0ea740576"
        },
        "date": 1733396522112,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1357096,
            "unit": "ns/op\t  819443 B/op\t       5 allocs/op",
            "extra": "860 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1357096,
            "unit": "ns/op",
            "extra": "860 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819443,
            "unit": "B/op",
            "extra": "860 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "860 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 17398,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "61934 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 17398,
            "unit": "ns/op",
            "extra": "61934 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "61934 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "61934 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 84.62,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13384503 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 84.62,
            "unit": "ns/op",
            "extra": "13384503 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13384503 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13384503 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22922,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52342 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22922,
            "unit": "ns/op",
            "extra": "52342 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52342 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52342 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 252673,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5301 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 252673,
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
            "value": 2682210,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "446 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2682210,
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
            "value": 44344527,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 44344527,
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
            "value": 9937,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "119400 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9937,
            "unit": "ns/op",
            "extra": "119400 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "119400 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "119400 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 8004948,
            "unit": "ns/op\t 4594781 B/op\t   75249 allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 8004948,
            "unit": "ns/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594781,
            "unit": "B/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75249,
            "unit": "allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8096,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "143271 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8096,
            "unit": "ns/op",
            "extra": "143271 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "143271 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "143271 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 889274,
            "unit": "ns/op\t  396569 B/op\t    6228 allocs/op",
            "extra": "1208 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 889274,
            "unit": "ns/op",
            "extra": "1208 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396569,
            "unit": "B/op",
            "extra": "1208 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6228,
            "unit": "allocs/op",
            "extra": "1208 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11298,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "103303 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11298,
            "unit": "ns/op",
            "extra": "103303 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "103303 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "103303 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 9058269,
            "unit": "ns/op\t 4981034 B/op\t   81257 allocs/op",
            "extra": "130 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 9058269,
            "unit": "ns/op",
            "extra": "130 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981034,
            "unit": "B/op",
            "extra": "130 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81257,
            "unit": "allocs/op",
            "extra": "130 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 229.6,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5209726 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 229.6,
            "unit": "ns/op",
            "extra": "5209726 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5209726 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5209726 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6191,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6191,
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
            "value": 301,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4016812 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 301,
            "unit": "ns/op",
            "extra": "4016812 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4016812 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4016812 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 454.5,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2662821 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 454.5,
            "unit": "ns/op",
            "extra": "2662821 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2662821 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2662821 times\n4 procs"
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
            "email": "jakub.warczarek@konghq.com",
            "name": "Jakub Warczarek",
            "username": "programmer04"
          },
          "distinct": true,
          "id": "c925f9e5d0f6f94befd3d334934dfa450a6de162",
          "message": "fix(log): get rid of `No targets found to create upstream ...`",
          "timestamp": "2024-12-05T12:14:18+01:00",
          "tree_id": "5d49d47024ac6f3d30449914d3213c3034d49a6c",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/c925f9e5d0f6f94befd3d334934dfa450a6de162"
        },
        "date": 1733397449575,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1216577,
            "unit": "ns/op\t  819447 B/op\t       5 allocs/op",
            "extra": "1034 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1216577,
            "unit": "ns/op",
            "extra": "1034 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819447,
            "unit": "B/op",
            "extra": "1034 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1034 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7837,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "166867 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7837,
            "unit": "ns/op",
            "extra": "166867 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "166867 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "166867 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 78.63,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15518935 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 78.63,
            "unit": "ns/op",
            "extra": "15518935 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15518935 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15518935 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22380,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53956 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22380,
            "unit": "ns/op",
            "extra": "53956 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53956 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53956 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 212613,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5949 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 212613,
            "unit": "ns/op",
            "extra": "5949 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5949 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5949 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2578955,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "412 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2578955,
            "unit": "ns/op",
            "extra": "412 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
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
            "value": 35235763,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "76 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 35235763,
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
            "value": 9423,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "117242 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9423,
            "unit": "ns/op",
            "extra": "117242 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "117242 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "117242 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7299332,
            "unit": "ns/op\t 4594646 B/op\t   75248 allocs/op",
            "extra": "166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7299332,
            "unit": "ns/op",
            "extra": "166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594646,
            "unit": "B/op",
            "extra": "166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75248,
            "unit": "allocs/op",
            "extra": "166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7746,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "151850 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7746,
            "unit": "ns/op",
            "extra": "151850 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "151850 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "151850 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 831649,
            "unit": "ns/op\t  396397 B/op\t    6225 allocs/op",
            "extra": "1261 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 831649,
            "unit": "ns/op",
            "extra": "1261 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396397,
            "unit": "B/op",
            "extra": "1261 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6225,
            "unit": "allocs/op",
            "extra": "1261 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10679,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "111690 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10679,
            "unit": "ns/op",
            "extra": "111690 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "111690 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "111690 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8345756,
            "unit": "ns/op\t 4981138 B/op\t   81258 allocs/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8345756,
            "unit": "ns/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981138,
            "unit": "B/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81258,
            "unit": "allocs/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 223.8,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5425897 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 223.8,
            "unit": "ns/op",
            "extra": "5425897 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5425897 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5425897 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.606,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.606,
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
            "value": 328.1,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "3308644 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 328.1,
            "unit": "ns/op",
            "extra": "3308644 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "3308644 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "3308644 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 447.2,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2708640 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 447.2,
            "unit": "ns/op",
            "extra": "2708640 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2708640 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2708640 times\n4 procs"
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
          "id": "c74c408637ae7a8a7557cbe07cf963780bd17bcf",
          "message": "ca from cm implemented\n\nSigned-off-by: Mattia Lavacca <lavacca.mattia@gmail.com>",
          "timestamp": "2024-12-05T12:16:00+01:00",
          "tree_id": "64858462992edb0a19df11a5a90d9399e7a0bb29",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/c74c408637ae7a8a7557cbe07cf963780bd17bcf"
        },
        "date": 1733397551606,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1265073,
            "unit": "ns/op\t  819448 B/op\t       5 allocs/op",
            "extra": "807 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1265073,
            "unit": "ns/op",
            "extra": "807 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819448,
            "unit": "B/op",
            "extra": "807 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "807 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 19858,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "51051 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 19858,
            "unit": "ns/op",
            "extra": "51051 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "51051 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "51051 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.27,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15038338 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.27,
            "unit": "ns/op",
            "extra": "15038338 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15038338 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15038338 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22393,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53248 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22393,
            "unit": "ns/op",
            "extra": "53248 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53248 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53248 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 222978,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4748 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 222978,
            "unit": "ns/op",
            "extra": "4748 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4748 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4748 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2619039,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "474 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2619039,
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
            "value": 30161155,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 30161155,
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
            "value": 9900,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "121329 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9900,
            "unit": "ns/op",
            "extra": "121329 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "121329 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "121329 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7418228,
            "unit": "ns/op\t 4595142 B/op\t   75255 allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7418228,
            "unit": "ns/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4595142,
            "unit": "B/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75255,
            "unit": "allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8231,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "143013 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8231,
            "unit": "ns/op",
            "extra": "143013 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "143013 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "143013 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 858711,
            "unit": "ns/op\t  396603 B/op\t    6231 allocs/op",
            "extra": "1267 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 858711,
            "unit": "ns/op",
            "extra": "1267 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396603,
            "unit": "B/op",
            "extra": "1267 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6231,
            "unit": "allocs/op",
            "extra": "1267 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11259,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "106710 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11259,
            "unit": "ns/op",
            "extra": "106710 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "106710 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "106710 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8489310,
            "unit": "ns/op\t 4981203 B/op\t   81263 allocs/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8489310,
            "unit": "ns/op",
            "extra": "139 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981203,
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
            "value": 227.2,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5278514 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 227.2,
            "unit": "ns/op",
            "extra": "5278514 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5278514 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5278514 times\n4 procs"
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
            "value": 302.3,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4061917 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 302.3,
            "unit": "ns/op",
            "extra": "4061917 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4061917 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4061917 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 502.4,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2072539 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 502.4,
            "unit": "ns/op",
            "extra": "2072539 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2072539 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2072539 times\n4 procs"
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
          "id": "7800ee434bfccaf3a37687c617db5203de5c6935",
          "message": "ca from cm implemented\n\nSigned-off-by: Mattia Lavacca <lavacca.mattia@gmail.com>",
          "timestamp": "2024-12-05T12:23:08+01:00",
          "tree_id": "1b24cbfcaee8bf8cb0cc0d73df08097a8174f1f1",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/7800ee434bfccaf3a37687c617db5203de5c6935"
        },
        "date": 1733397977527,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1349201,
            "unit": "ns/op\t  819442 B/op\t       5 allocs/op",
            "extra": "1080 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1349201,
            "unit": "ns/op",
            "extra": "1080 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819442,
            "unit": "B/op",
            "extra": "1080 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1080 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7323,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "141616 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7323,
            "unit": "ns/op",
            "extra": "141616 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "141616 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "141616 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.28,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15155307 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.28,
            "unit": "ns/op",
            "extra": "15155307 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15155307 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15155307 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22372,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53528 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22372,
            "unit": "ns/op",
            "extra": "53528 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53528 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53528 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 222540,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5614 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 222540,
            "unit": "ns/op",
            "extra": "5614 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5614 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5614 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2606764,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "404 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2606764,
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
            "value": 32136964,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "75 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 32136964,
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
            "value": 9940,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "118302 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9940,
            "unit": "ns/op",
            "extra": "118302 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "118302 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "118302 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7502908,
            "unit": "ns/op\t 4594997 B/op\t   75254 allocs/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7502908,
            "unit": "ns/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594997,
            "unit": "B/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75254,
            "unit": "allocs/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8335,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "140553 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8335,
            "unit": "ns/op",
            "extra": "140553 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "140553 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "140553 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 869092,
            "unit": "ns/op\t  396738 B/op\t    6233 allocs/op",
            "extra": "1227 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 869092,
            "unit": "ns/op",
            "extra": "1227 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396738,
            "unit": "B/op",
            "extra": "1227 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6233,
            "unit": "allocs/op",
            "extra": "1227 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11263,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "105242 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11263,
            "unit": "ns/op",
            "extra": "105242 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "105242 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "105242 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8424903,
            "unit": "ns/op\t 4981252 B/op\t   81263 allocs/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8424903,
            "unit": "ns/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981252,
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
            "value": 233.5,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5307123 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 233.5,
            "unit": "ns/op",
            "extra": "5307123 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5307123 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5307123 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6619,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6619,
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
            "value": 297.2,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4103559 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 297.2,
            "unit": "ns/op",
            "extra": "4103559 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4103559 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4103559 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 449.8,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2645884 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 449.8,
            "unit": "ns/op",
            "extra": "2645884 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2645884 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2645884 times\n4 procs"
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
            "email": "jakub.warczarek@konghq.com",
            "name": "Jakub Warczarek",
            "username": "programmer04"
          },
          "distinct": true,
          "id": "f2078b6d7658b8899658794fcb6e726f354306fd",
          "message": "fix(log): get rid of `No targets found to create upstream ...`",
          "timestamp": "2024-12-05T12:17:54+01:00",
          "tree_id": "d7c45210e94402386c9b25eec1caa843d1e61e4b",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/f2078b6d7658b8899658794fcb6e726f354306fd"
        },
        "date": 1733398456247,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1192889,
            "unit": "ns/op\t  819447 B/op\t       5 allocs/op",
            "extra": "1023 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1192889,
            "unit": "ns/op",
            "extra": "1023 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819447,
            "unit": "B/op",
            "extra": "1023 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1023 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 8516,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "169197 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 8516,
            "unit": "ns/op",
            "extra": "169197 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "169197 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "169197 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 80.82,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15168314 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 80.82,
            "unit": "ns/op",
            "extra": "15168314 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15168314 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15168314 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22524,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53017 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22524,
            "unit": "ns/op",
            "extra": "53017 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53017 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53017 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 224605,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5656 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 224605,
            "unit": "ns/op",
            "extra": "5656 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5656 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5656 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2626136,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "451 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2626136,
            "unit": "ns/op",
            "extra": "451 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
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
            "value": 45925164,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "40 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 45925164,
            "unit": "ns/op",
            "extra": "40 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
            "unit": "B/op",
            "extra": "40 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "40 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9744,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "123799 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9744,
            "unit": "ns/op",
            "extra": "123799 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "123799 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "123799 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7377847,
            "unit": "ns/op\t 4594751 B/op\t   75248 allocs/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7377847,
            "unit": "ns/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594751,
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
            "value": 7958,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "148723 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7958,
            "unit": "ns/op",
            "extra": "148723 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "148723 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "148723 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 850741,
            "unit": "ns/op\t  396473 B/op\t    6226 allocs/op",
            "extra": "1238 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 850741,
            "unit": "ns/op",
            "extra": "1238 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396473,
            "unit": "B/op",
            "extra": "1238 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6226,
            "unit": "allocs/op",
            "extra": "1238 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11044,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "107802 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11044,
            "unit": "ns/op",
            "extra": "107802 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "107802 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "107802 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8829643,
            "unit": "ns/op\t 4981249 B/op\t   81258 allocs/op",
            "extra": "132 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8829643,
            "unit": "ns/op",
            "extra": "132 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981249,
            "unit": "B/op",
            "extra": "132 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81258,
            "unit": "allocs/op",
            "extra": "132 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 226.7,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5211560 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 226.7,
            "unit": "ns/op",
            "extra": "5211560 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5211560 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5211560 times\n4 procs"
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
            "value": 299,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4039306 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 299,
            "unit": "ns/op",
            "extra": "4039306 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4039306 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4039306 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 450.8,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2643376 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 450.8,
            "unit": "ns/op",
            "extra": "2643376 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2643376 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2643376 times\n4 procs"
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
          "id": "a0e31747eb503715f405e60675ecb752792914b9",
          "message": "fix(log): get rid of `No targets found to create upstream ...` (#6781)",
          "timestamp": "2024-12-05T14:51:22+01:00",
          "tree_id": "d7c45210e94402386c9b25eec1caa843d1e61e4b",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/a0e31747eb503715f405e60675ecb752792914b9"
        },
        "date": 1733406869443,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1238774,
            "unit": "ns/op\t  819447 B/op\t       5 allocs/op",
            "extra": "972 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1238774,
            "unit": "ns/op",
            "extra": "972 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819447,
            "unit": "B/op",
            "extra": "972 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "972 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 6971,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "172322 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 6971,
            "unit": "ns/op",
            "extra": "172322 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "172322 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "172322 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.23,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15036654 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.23,
            "unit": "ns/op",
            "extra": "15036654 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15036654 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15036654 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 23169,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52046 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 23169,
            "unit": "ns/op",
            "extra": "52046 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52046 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52046 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 291425,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5344 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 291425,
            "unit": "ns/op",
            "extra": "5344 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5344 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5344 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2624617,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "445 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2624617,
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
            "value": 36223140,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "66 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 36223140,
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
            "value": 9960,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "117487 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9960,
            "unit": "ns/op",
            "extra": "117487 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "117487 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "117487 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7854595,
            "unit": "ns/op\t 4594594 B/op\t   75248 allocs/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7854595,
            "unit": "ns/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594594,
            "unit": "B/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75248,
            "unit": "allocs/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8416,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "143384 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8416,
            "unit": "ns/op",
            "extra": "143384 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "143384 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "143384 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 881301,
            "unit": "ns/op\t  396713 B/op\t    6230 allocs/op",
            "extra": "1167 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 881301,
            "unit": "ns/op",
            "extra": "1167 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396713,
            "unit": "B/op",
            "extra": "1167 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6230,
            "unit": "allocs/op",
            "extra": "1167 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11492,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "105190 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11492,
            "unit": "ns/op",
            "extra": "105190 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "105190 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "105190 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 9226456,
            "unit": "ns/op\t 4981387 B/op\t   81258 allocs/op",
            "extra": "130 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 9226456,
            "unit": "ns/op",
            "extra": "130 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981387,
            "unit": "B/op",
            "extra": "130 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81258,
            "unit": "allocs/op",
            "extra": "130 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 229.1,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5264378 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 229.1,
            "unit": "ns/op",
            "extra": "5264378 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5264378 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5264378 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6191,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6191,
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
            "value": 296.4,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4027392 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 296.4,
            "unit": "ns/op",
            "extra": "4027392 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4027392 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4027392 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 467,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2630683 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 467,
            "unit": "ns/op",
            "extra": "2630683 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2630683 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2630683 times\n4 procs"
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
          "id": "d27f67078bf768270effe7d032f2a5cfb66c992a",
          "message": "chore(deps): bump google.golang.org/grpc from 1.68.0 to 1.68.1\n\nBumps [google.golang.org/grpc](https://github.com/grpc/grpc-go) from 1.68.0 to 1.68.1.\n- [Release notes](https://github.com/grpc/grpc-go/releases)\n- [Commits](https://github.com/grpc/grpc-go/compare/v1.68.0...v1.68.1)\n\n---\nupdated-dependencies:\n- dependency-name: google.golang.org/grpc\n  dependency-type: direct:production\n  update-type: version-update:semver-patch\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>",
          "timestamp": "2024-12-05T14:13:43Z",
          "tree_id": "840e27d50ead5c2f3b0703e94b6544bdbaec6316",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/d27f67078bf768270effe7d032f2a5cfb66c992a"
        },
        "date": 1733408211615,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1159134,
            "unit": "ns/op\t  819441 B/op\t       5 allocs/op",
            "extra": "1134 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1159134,
            "unit": "ns/op",
            "extra": "1134 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819441,
            "unit": "B/op",
            "extra": "1134 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1134 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7473,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "163342 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7473,
            "unit": "ns/op",
            "extra": "163342 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "163342 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "163342 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.03,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14858499 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.03,
            "unit": "ns/op",
            "extra": "14858499 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14858499 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14858499 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22639,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52628 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22639,
            "unit": "ns/op",
            "extra": "52628 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52628 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52628 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 228209,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5036 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 228209,
            "unit": "ns/op",
            "extra": "5036 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5036 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5036 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2630032,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "447 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2630032,
            "unit": "ns/op",
            "extra": "447 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
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
            "value": 40117478,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 40117478,
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
            "value": 9811,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "122221 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9811,
            "unit": "ns/op",
            "extra": "122221 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "122221 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "122221 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7630694,
            "unit": "ns/op\t 4594705 B/op\t   75248 allocs/op",
            "extra": "154 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7630694,
            "unit": "ns/op",
            "extra": "154 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594705,
            "unit": "B/op",
            "extra": "154 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75248,
            "unit": "allocs/op",
            "extra": "154 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8144,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "145364 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8144,
            "unit": "ns/op",
            "extra": "145364 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "145364 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "145364 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 872218,
            "unit": "ns/op\t  396424 B/op\t    6225 allocs/op",
            "extra": "1255 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 872218,
            "unit": "ns/op",
            "extra": "1255 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396424,
            "unit": "B/op",
            "extra": "1255 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6225,
            "unit": "allocs/op",
            "extra": "1255 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11150,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "107391 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11150,
            "unit": "ns/op",
            "extra": "107391 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "107391 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "107391 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8695841,
            "unit": "ns/op\t 4980920 B/op\t   81257 allocs/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8695841,
            "unit": "ns/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4980920,
            "unit": "B/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81257,
            "unit": "allocs/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 226.3,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5268424 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 226.3,
            "unit": "ns/op",
            "extra": "5268424 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5268424 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5268424 times\n4 procs"
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
            "value": 357.2,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "3449750 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 357.2,
            "unit": "ns/op",
            "extra": "3449750 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "3449750 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "3449750 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 447.9,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2669815 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 447.9,
            "unit": "ns/op",
            "extra": "2669815 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2669815 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2669815 times\n4 procs"
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
          "id": "763876a4a837d8392d008e0b4074fedca248dddc",
          "message": "chore(deps): bump github.com/prometheus/common from 0.60.1 to 0.61.0\n\nBumps [github.com/prometheus/common](https://github.com/prometheus/common) from 0.60.1 to 0.61.0.\n- [Release notes](https://github.com/prometheus/common/releases)\n- [Changelog](https://github.com/prometheus/common/blob/main/RELEASE.md)\n- [Commits](https://github.com/prometheus/common/compare/v0.60.1...v0.61.0)\n\n---\nupdated-dependencies:\n- dependency-name: github.com/prometheus/common\n  dependency-type: direct:production\n  update-type: version-update:semver-minor\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>",
          "timestamp": "2024-12-05T14:13:58Z",
          "tree_id": "d9729f67b23f44ac4d508b19bcd6774b0e3a9ee3",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/763876a4a837d8392d008e0b4074fedca248dddc"
        },
        "date": 1733408233879,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1132601,
            "unit": "ns/op\t  819441 B/op\t       5 allocs/op",
            "extra": "916 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1132601,
            "unit": "ns/op",
            "extra": "916 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819441,
            "unit": "B/op",
            "extra": "916 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "916 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 17555,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "57670 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 17555,
            "unit": "ns/op",
            "extra": "57670 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "57670 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "57670 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 80.79,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13106280 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 80.79,
            "unit": "ns/op",
            "extra": "13106280 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13106280 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13106280 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22656,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52975 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22656,
            "unit": "ns/op",
            "extra": "52975 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52975 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52975 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 237507,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5480 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 237507,
            "unit": "ns/op",
            "extra": "5480 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5480 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5480 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2589375,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "486 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2589375,
            "unit": "ns/op",
            "extra": "486 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
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
            "value": 40180003,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "57 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 40180003,
            "unit": "ns/op",
            "extra": "57 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 9758,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "119763 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9758,
            "unit": "ns/op",
            "extra": "119763 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "119763 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "119763 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 8935652,
            "unit": "ns/op\t 4594315 B/op\t   75247 allocs/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 8935652,
            "unit": "ns/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594315,
            "unit": "B/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75247,
            "unit": "allocs/op",
            "extra": "145 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8716,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "137730 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8716,
            "unit": "ns/op",
            "extra": "137730 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "137730 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "137730 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 901932,
            "unit": "ns/op\t  396727 B/op\t    6230 allocs/op",
            "extra": "1162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 901932,
            "unit": "ns/op",
            "extra": "1162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396727,
            "unit": "B/op",
            "extra": "1162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6230,
            "unit": "allocs/op",
            "extra": "1162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11902,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "96964 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11902,
            "unit": "ns/op",
            "extra": "96964 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "96964 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "96964 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 10019471,
            "unit": "ns/op\t 4980945 B/op\t   81257 allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 10019471,
            "unit": "ns/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4980945,
            "unit": "B/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81257,
            "unit": "allocs/op",
            "extra": "100 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 234.5,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5161443 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 234.5,
            "unit": "ns/op",
            "extra": "5161443 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5161443 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5161443 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6867,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6867,
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
            "value": 302.9,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "3956737 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 302.9,
            "unit": "ns/op",
            "extra": "3956737 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "3956737 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "3956737 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 458.3,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2566771 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 458.3,
            "unit": "ns/op",
            "extra": "2566771 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2566771 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2566771 times\n4 procs"
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
          "id": "491211aff1c7048b621b865a258b95665c1d9f69",
          "message": "chore(deps): bump google.golang.org/api from 0.209.0 to 0.210.0\n\nBumps [google.golang.org/api](https://github.com/googleapis/google-api-go-client) from 0.209.0 to 0.210.0.\n- [Release notes](https://github.com/googleapis/google-api-go-client/releases)\n- [Changelog](https://github.com/googleapis/google-api-go-client/blob/main/CHANGES.md)\n- [Commits](https://github.com/googleapis/google-api-go-client/compare/v0.209.0...v0.210.0)\n\n---\nupdated-dependencies:\n- dependency-name: google.golang.org/api\n  dependency-type: direct:production\n  update-type: version-update:semver-minor\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>",
          "timestamp": "2024-12-05T14:14:18Z",
          "tree_id": "b0e330c844bf1b6935ba5ab073814aa85e3e6422",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/491211aff1c7048b621b865a258b95665c1d9f69"
        },
        "date": 1733408254985,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1523241,
            "unit": "ns/op\t  819449 B/op\t       5 allocs/op",
            "extra": "958 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1523241,
            "unit": "ns/op",
            "extra": "958 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819449,
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
            "value": 7436,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "160418 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7436,
            "unit": "ns/op",
            "extra": "160418 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "160418 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "160418 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 80.76,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15250387 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 80.76,
            "unit": "ns/op",
            "extra": "15250387 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15250387 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15250387 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22558,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52606 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22558,
            "unit": "ns/op",
            "extra": "52606 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52606 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52606 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 230042,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4708 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 230042,
            "unit": "ns/op",
            "extra": "4708 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4708 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4708 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2569027,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "452 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2569027,
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
            "value": 43207954,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 43207954,
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
            "value": 9687,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "123625 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9687,
            "unit": "ns/op",
            "extra": "123625 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "123625 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "123625 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7769220,
            "unit": "ns/op\t 4594711 B/op\t   75248 allocs/op",
            "extra": "156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7769220,
            "unit": "ns/op",
            "extra": "156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594711,
            "unit": "B/op",
            "extra": "156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75248,
            "unit": "allocs/op",
            "extra": "156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8119,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "148376 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8119,
            "unit": "ns/op",
            "extra": "148376 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "148376 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "148376 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 861105,
            "unit": "ns/op\t  396400 B/op\t    6225 allocs/op",
            "extra": "1261 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 861105,
            "unit": "ns/op",
            "extra": "1261 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396400,
            "unit": "B/op",
            "extra": "1261 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6225,
            "unit": "allocs/op",
            "extra": "1261 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11351,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "103474 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11351,
            "unit": "ns/op",
            "extra": "103474 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "103474 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "103474 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 9186884,
            "unit": "ns/op\t 4980884 B/op\t   81257 allocs/op",
            "extra": "135 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 9186884,
            "unit": "ns/op",
            "extra": "135 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4980884,
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
            "value": 236.5,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5305045 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 236.5,
            "unit": "ns/op",
            "extra": "5305045 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5305045 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5305045 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6693,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6693,
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
            "value": 297.7,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "3692080 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 297.7,
            "unit": "ns/op",
            "extra": "3692080 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "3692080 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "3692080 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 448,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2664826 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 448,
            "unit": "ns/op",
            "extra": "2664826 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2664826 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2664826 times\n4 procs"
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
          "id": "e0f7970103a1e4f6e6d5f7772f5bbde2988943c3",
          "message": "chore(deps): bump google.golang.org/grpc from 1.68.0 to 1.68.1 (#6782)\n\nBumps [google.golang.org/grpc](https://github.com/grpc/grpc-go) from 1.68.0 to 1.68.1.\r\n- [Release notes](https://github.com/grpc/grpc-go/releases)\r\n- [Commits](https://github.com/grpc/grpc-go/compare/v1.68.0...v1.68.1)\r\n\r\n---\r\nupdated-dependencies:\r\n- dependency-name: google.golang.org/grpc\r\n  dependency-type: direct:production\r\n  update-type: version-update:semver-patch\r\n...\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2024-12-05T14:33:28Z",
          "tree_id": "840e27d50ead5c2f3b0703e94b6544bdbaec6316",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/e0f7970103a1e4f6e6d5f7772f5bbde2988943c3"
        },
        "date": 1733409407433,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1349608,
            "unit": "ns/op\t  819441 B/op\t       5 allocs/op",
            "extra": "819 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1349608,
            "unit": "ns/op",
            "extra": "819 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819441,
            "unit": "B/op",
            "extra": "819 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "819 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 6953,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "171787 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 6953,
            "unit": "ns/op",
            "extra": "171787 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "171787 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "171787 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.45,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15185362 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.45,
            "unit": "ns/op",
            "extra": "15185362 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15185362 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15185362 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22641,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53000 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22641,
            "unit": "ns/op",
            "extra": "53000 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53000 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53000 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 233792,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4646 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 233792,
            "unit": "ns/op",
            "extra": "4646 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4646 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4646 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2651463,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "445 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2651463,
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
            "value": 34794824,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "69 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 34794824,
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
            "value": 9973,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "119432 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9973,
            "unit": "ns/op",
            "extra": "119432 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "119432 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "119432 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7849585,
            "unit": "ns/op\t 4594828 B/op\t   75249 allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7849585,
            "unit": "ns/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594828,
            "unit": "B/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75249,
            "unit": "allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8269,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "145072 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8269,
            "unit": "ns/op",
            "extra": "145072 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "145072 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "145072 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 889921,
            "unit": "ns/op\t  396610 B/op\t    6228 allocs/op",
            "extra": "1197 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 889921,
            "unit": "ns/op",
            "extra": "1197 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396610,
            "unit": "B/op",
            "extra": "1197 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6228,
            "unit": "allocs/op",
            "extra": "1197 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11513,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "103928 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11513,
            "unit": "ns/op",
            "extra": "103928 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "103928 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "103928 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 9040619,
            "unit": "ns/op\t 4980991 B/op\t   81257 allocs/op",
            "extra": "134 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 9040619,
            "unit": "ns/op",
            "extra": "134 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4980991,
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
            "value": 228.2,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5314447 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 228.2,
            "unit": "ns/op",
            "extra": "5314447 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5314447 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5314447 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6255,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6255,
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
            "extra": "4041496 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 298.2,
            "unit": "ns/op",
            "extra": "4041496 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4041496 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4041496 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 451.7,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2618085 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 451.7,
            "unit": "ns/op",
            "extra": "2618085 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2618085 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2618085 times\n4 procs"
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
          "id": "1981516d07a8e02c9b424eb2188988258417fd89",
          "message": "chore(deps): bump github.com/prometheus/common from 0.60.1 to 0.61.0 (#6783)\n\nBumps [github.com/prometheus/common](https://github.com/prometheus/common) from 0.60.1 to 0.61.0.\r\n- [Release notes](https://github.com/prometheus/common/releases)\r\n- [Changelog](https://github.com/prometheus/common/blob/main/RELEASE.md)\r\n- [Commits](https://github.com/prometheus/common/compare/v0.60.1...v0.61.0)\r\n\r\n---\r\nupdated-dependencies:\r\n- dependency-name: github.com/prometheus/common\r\n  dependency-type: direct:production\r\n  update-type: version-update:semver-minor\r\n...\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2024-12-05T14:34:07Z",
          "tree_id": "7aff9d183edb20aa0d78d31e93d192cbcbe96545",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/1981516d07a8e02c9b424eb2188988258417fd89"
        },
        "date": 1733409433335,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1351469,
            "unit": "ns/op\t  819447 B/op\t       5 allocs/op",
            "extra": "1149 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1351469,
            "unit": "ns/op",
            "extra": "1149 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819447,
            "unit": "B/op",
            "extra": "1149 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1149 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 8437,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "145566 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 8437,
            "unit": "ns/op",
            "extra": "145566 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "145566 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "145566 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 78.96,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15150114 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 78.96,
            "unit": "ns/op",
            "extra": "15150114 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15150114 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15150114 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22839,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53160 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22839,
            "unit": "ns/op",
            "extra": "53160 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53160 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53160 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 223166,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4993 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 223166,
            "unit": "ns/op",
            "extra": "4993 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
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
            "value": 2623436,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "459 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2623436,
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
            "value": 35351682,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 35351682,
            "unit": "ns/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
            "unit": "B/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9726,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "124725 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9726,
            "unit": "ns/op",
            "extra": "124725 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "124725 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "124725 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7457877,
            "unit": "ns/op\t 4594640 B/op\t   75248 allocs/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7457877,
            "unit": "ns/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594640,
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
            "value": 7965,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "149245 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7965,
            "unit": "ns/op",
            "extra": "149245 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "149245 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "149245 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 866250,
            "unit": "ns/op\t  396468 B/op\t    6226 allocs/op",
            "extra": "1238 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 866250,
            "unit": "ns/op",
            "extra": "1238 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396468,
            "unit": "B/op",
            "extra": "1238 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6226,
            "unit": "allocs/op",
            "extra": "1238 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11054,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "107596 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11054,
            "unit": "ns/op",
            "extra": "107596 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "107596 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "107596 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8606476,
            "unit": "ns/op\t 4981402 B/op\t   81258 allocs/op",
            "extra": "135 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8606476,
            "unit": "ns/op",
            "extra": "135 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981402,
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
            "value": 231.4,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5209414 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 231.4,
            "unit": "ns/op",
            "extra": "5209414 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5209414 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5209414 times\n4 procs"
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
            "value": 295.7,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "3947638 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 295.7,
            "unit": "ns/op",
            "extra": "3947638 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "3947638 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "3947638 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 470,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2652948 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 470,
            "unit": "ns/op",
            "extra": "2652948 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2652948 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2652948 times\n4 procs"
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
          "id": "41670b6fa1d1c5bdada854f8a883ef3c9db5b2ff",
          "message": "chore(deps): bump google.golang.org/api from 0.209.0 to 0.210.0\n\nBumps [google.golang.org/api](https://github.com/googleapis/google-api-go-client) from 0.209.0 to 0.210.0.\n- [Release notes](https://github.com/googleapis/google-api-go-client/releases)\n- [Changelog](https://github.com/googleapis/google-api-go-client/blob/main/CHANGES.md)\n- [Commits](https://github.com/googleapis/google-api-go-client/compare/v0.209.0...v0.210.0)\n\n---\nupdated-dependencies:\n- dependency-name: google.golang.org/api\n  dependency-type: direct:production\n  update-type: version-update:semver-minor\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>",
          "timestamp": "2024-12-05T14:34:36Z",
          "tree_id": "746210937f6fde80af247bbf44ede86150b46271",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/41670b6fa1d1c5bdada854f8a883ef3c9db5b2ff"
        },
        "date": 1733409464949,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1121351,
            "unit": "ns/op\t  819448 B/op\t       5 allocs/op",
            "extra": "898 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1121351,
            "unit": "ns/op",
            "extra": "898 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819448,
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
            "value": 6901,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "173466 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 6901,
            "unit": "ns/op",
            "extra": "173466 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "173466 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "173466 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.51,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15024932 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.51,
            "unit": "ns/op",
            "extra": "15024932 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15024932 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15024932 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22621,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52328 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22621,
            "unit": "ns/op",
            "extra": "52328 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52328 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52328 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 218192,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5575 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 218192,
            "unit": "ns/op",
            "extra": "5575 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5575 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5575 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2479880,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "404 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2479880,
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
            "value": 37201958,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "75 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 37201958,
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
            "value": 9487,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "124803 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9487,
            "unit": "ns/op",
            "extra": "124803 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "124803 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "124803 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7343514,
            "unit": "ns/op\t 4594434 B/op\t   75247 allocs/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7343514,
            "unit": "ns/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594434,
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
            "value": 7867,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "149256 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7867,
            "unit": "ns/op",
            "extra": "149256 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "149256 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "149256 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 884057,
            "unit": "ns/op\t  396440 B/op\t    6226 allocs/op",
            "extra": "1246 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 884057,
            "unit": "ns/op",
            "extra": "1246 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396440,
            "unit": "B/op",
            "extra": "1246 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6226,
            "unit": "allocs/op",
            "extra": "1246 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10864,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "108242 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10864,
            "unit": "ns/op",
            "extra": "108242 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "108242 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "108242 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8416954,
            "unit": "ns/op\t 4981052 B/op\t   81257 allocs/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8416954,
            "unit": "ns/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981052,
            "unit": "B/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81257,
            "unit": "allocs/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 227.5,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5286117 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 227.5,
            "unit": "ns/op",
            "extra": "5286117 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5286117 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5286117 times\n4 procs"
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
            "value": 295.9,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4011114 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 295.9,
            "unit": "ns/op",
            "extra": "4011114 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4011114 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4011114 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 453.2,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2663908 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 453.2,
            "unit": "ns/op",
            "extra": "2663908 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2663908 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2663908 times\n4 procs"
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
          "id": "cc92b41a3aadb9df792bcc2608198635b9abe898",
          "message": "chore(deps): bump google.golang.org/api from 0.209.0 to 0.210.0\n\nBumps [google.golang.org/api](https://github.com/googleapis/google-api-go-client) from 0.209.0 to 0.210.0.\n- [Release notes](https://github.com/googleapis/google-api-go-client/releases)\n- [Changelog](https://github.com/googleapis/google-api-go-client/blob/main/CHANGES.md)\n- [Commits](https://github.com/googleapis/google-api-go-client/compare/v0.209.0...v0.210.0)\n\n---\nupdated-dependencies:\n- dependency-name: google.golang.org/api\n  dependency-type: direct:production\n  update-type: version-update:semver-minor\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>",
          "timestamp": "2024-12-05T14:35:25Z",
          "tree_id": "c8a91a0c746ad0d6a375a4ab5922f9fde429cd86",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/cc92b41a3aadb9df792bcc2608198635b9abe898"
        },
        "date": 1733409515828,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1258870,
            "unit": "ns/op\t  819447 B/op\t       5 allocs/op",
            "extra": "1107 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1258870,
            "unit": "ns/op",
            "extra": "1107 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819447,
            "unit": "B/op",
            "extra": "1107 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1107 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 6927,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "172797 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 6927,
            "unit": "ns/op",
            "extra": "172797 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "172797 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "172797 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.06,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15157752 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.06,
            "unit": "ns/op",
            "extra": "15157752 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15157752 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15157752 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22647,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52962 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22647,
            "unit": "ns/op",
            "extra": "52962 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52962 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52962 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 221701,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5493 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 221701,
            "unit": "ns/op",
            "extra": "5493 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5493 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5493 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2534764,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "464 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2534764,
            "unit": "ns/op",
            "extra": "464 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
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
            "value": 39453901,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "75 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 39453901,
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
            "value": 10166,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "107569 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10166,
            "unit": "ns/op",
            "extra": "107569 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "107569 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "107569 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7495021,
            "unit": "ns/op\t 4594400 B/op\t   75247 allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7495021,
            "unit": "ns/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594400,
            "unit": "B/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75247,
            "unit": "allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8118,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "148332 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8118,
            "unit": "ns/op",
            "extra": "148332 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "148332 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "148332 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 866728,
            "unit": "ns/op\t  396429 B/op\t    6225 allocs/op",
            "extra": "1252 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 866728,
            "unit": "ns/op",
            "extra": "1252 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396429,
            "unit": "B/op",
            "extra": "1252 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6225,
            "unit": "allocs/op",
            "extra": "1252 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11377,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "104398 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11377,
            "unit": "ns/op",
            "extra": "104398 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "104398 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "104398 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8646116,
            "unit": "ns/op\t 4981062 B/op\t   81257 allocs/op",
            "extra": "133 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8646116,
            "unit": "ns/op",
            "extra": "133 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981062,
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
            "extra": "5253865 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 228.5,
            "unit": "ns/op",
            "extra": "5253865 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5253865 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5253865 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6192,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6192,
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
            "extra": "4062942 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 294.9,
            "unit": "ns/op",
            "extra": "4062942 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4062942 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4062942 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 449.3,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2696431 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 449.3,
            "unit": "ns/op",
            "extra": "2696431 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2696431 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2696431 times\n4 procs"
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
        "date": 1733411754220,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1044726,
            "unit": "ns/op\t  819441 B/op\t       5 allocs/op",
            "extra": "1026 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1044726,
            "unit": "ns/op",
            "extra": "1026 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819441,
            "unit": "B/op",
            "extra": "1026 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1026 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 8510,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "118012 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 8510,
            "unit": "ns/op",
            "extra": "118012 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "118012 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "118012 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.02,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14539680 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.02,
            "unit": "ns/op",
            "extra": "14539680 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14539680 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14539680 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 25564,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52855 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 25564,
            "unit": "ns/op",
            "extra": "52855 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52855 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52855 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 242244,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4154 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 242244,
            "unit": "ns/op",
            "extra": "4154 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4154 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4154 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2633954,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "452 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2633954,
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
            "value": 38612902,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "44 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 38612902,
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
            "value": 9508,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "124275 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9508,
            "unit": "ns/op",
            "extra": "124275 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "124275 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "124275 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7359690,
            "unit": "ns/op\t 4594950 B/op\t   75249 allocs/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7359690,
            "unit": "ns/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594950,
            "unit": "B/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75249,
            "unit": "allocs/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8065,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "150715 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8065,
            "unit": "ns/op",
            "extra": "150715 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "150715 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "150715 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 868727,
            "unit": "ns/op\t  396375 B/op\t    6225 allocs/op",
            "extra": "1269 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 868727,
            "unit": "ns/op",
            "extra": "1269 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396375,
            "unit": "B/op",
            "extra": "1269 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6225,
            "unit": "allocs/op",
            "extra": "1269 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10968,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "109023 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10968,
            "unit": "ns/op",
            "extra": "109023 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "109023 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "109023 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8448871,
            "unit": "ns/op\t 4981150 B/op\t   81258 allocs/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8448871,
            "unit": "ns/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981150,
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
            "value": 280.4,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5249210 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 280.4,
            "unit": "ns/op",
            "extra": "5249210 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5249210 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5249210 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6398,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6398,
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
            "extra": "4054074 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 296.8,
            "unit": "ns/op",
            "extra": "4054074 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4054074 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4054074 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 448.5,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2673877 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 448.5,
            "unit": "ns/op",
            "extra": "2673877 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2673877 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2673877 times\n4 procs"
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
          "id": "76062fc05245dfcfa5c3bf7275de4cab21ab4e65",
          "message": "ci: only run benchmarks on main and labeled PRs (#6785)",
          "timestamp": "2024-12-05T15:17:35Z",
          "tree_id": "b84764409c51bf9a13b06e9474a39d1906788a96",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/76062fc05245dfcfa5c3bf7275de4cab21ab4e65"
        },
        "date": 1733412054835,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1285651,
            "unit": "ns/op\t  819446 B/op\t       5 allocs/op",
            "extra": "919 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1285651,
            "unit": "ns/op",
            "extra": "919 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819446,
            "unit": "B/op",
            "extra": "919 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "919 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 23730,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "44493 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 23730,
            "unit": "ns/op",
            "extra": "44493 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "44493 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "44493 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 90.79,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13828644 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 90.79,
            "unit": "ns/op",
            "extra": "13828644 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13828644 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13828644 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22948,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "48180 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22948,
            "unit": "ns/op",
            "extra": "48180 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "48180 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "48180 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 229361,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5083 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 229361,
            "unit": "ns/op",
            "extra": "5083 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5083 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5083 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2528755,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "440 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2528755,
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
            "value": 40913795,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 40913795,
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
            "value": 9841,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "122289 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9841,
            "unit": "ns/op",
            "extra": "122289 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "122289 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "122289 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7756554,
            "unit": "ns/op\t 4594751 B/op\t   75248 allocs/op",
            "extra": "154 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7756554,
            "unit": "ns/op",
            "extra": "154 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594751,
            "unit": "B/op",
            "extra": "154 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75248,
            "unit": "allocs/op",
            "extra": "154 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8161,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "145670 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8161,
            "unit": "ns/op",
            "extra": "145670 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "145670 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "145670 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 888981,
            "unit": "ns/op\t  396502 B/op\t    6227 allocs/op",
            "extra": "1227 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 888981,
            "unit": "ns/op",
            "extra": "1227 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396502,
            "unit": "B/op",
            "extra": "1227 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6227,
            "unit": "allocs/op",
            "extra": "1227 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11628,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "101210 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11628,
            "unit": "ns/op",
            "extra": "101210 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "101210 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "101210 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8941638,
            "unit": "ns/op\t 4981096 B/op\t   81258 allocs/op",
            "extra": "129 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8941638,
            "unit": "ns/op",
            "extra": "129 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981096,
            "unit": "B/op",
            "extra": "129 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81258,
            "unit": "allocs/op",
            "extra": "129 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 230.1,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5183887 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 230.1,
            "unit": "ns/op",
            "extra": "5183887 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5183887 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5183887 times\n4 procs"
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
            "value": 302.8,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4004526 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 302.8,
            "unit": "ns/op",
            "extra": "4004526 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4004526 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4004526 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 537.4,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2659281 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 537.4,
            "unit": "ns/op",
            "extra": "2659281 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2659281 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2659281 times\n4 procs"
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
          "id": "eb85e6dea4c8aaf5497d6c06e4dbd1cb1f1418eb",
          "message": "configmap controller introduced\n\nSigned-off-by: Mattia Lavacca <lavacca.mattia@gmail.com>",
          "timestamp": "2024-12-05T17:32:21+01:00",
          "tree_id": "da27653998c7d07bdf54ff2ef01438bf74bd2724",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/eb85e6dea4c8aaf5497d6c06e4dbd1cb1f1418eb"
        },
        "date": 1733416541304,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1052102,
            "unit": "ns/op\t  819446 B/op\t       5 allocs/op",
            "extra": "1028 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1052102,
            "unit": "ns/op",
            "extra": "1028 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819446,
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
            "value": 7332,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "163984 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7332,
            "unit": "ns/op",
            "extra": "163984 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "163984 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "163984 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 80.48,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15150064 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 80.48,
            "unit": "ns/op",
            "extra": "15150064 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15150064 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15150064 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22548,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53338 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22548,
            "unit": "ns/op",
            "extra": "53338 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53338 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53338 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 219464,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5343 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 219464,
            "unit": "ns/op",
            "extra": "5343 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5343 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5343 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2571264,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "410 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2571264,
            "unit": "ns/op",
            "extra": "410 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
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
            "value": 34280822,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "30 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 34280822,
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
            "value": 9986,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "119912 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9986,
            "unit": "ns/op",
            "extra": "119912 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "119912 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "119912 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7405245,
            "unit": "ns/op\t 4594878 B/op\t   75254 allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7405245,
            "unit": "ns/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594878,
            "unit": "B/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75254,
            "unit": "allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8343,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "139948 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8343,
            "unit": "ns/op",
            "extra": "139948 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "139948 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "139948 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 870157,
            "unit": "ns/op\t  396707 B/op\t    6232 allocs/op",
            "extra": "1236 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 870157,
            "unit": "ns/op",
            "extra": "1236 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396707,
            "unit": "B/op",
            "extra": "1236 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6232,
            "unit": "allocs/op",
            "extra": "1236 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11365,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "104859 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11365,
            "unit": "ns/op",
            "extra": "104859 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "104859 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "104859 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8447292,
            "unit": "ns/op\t 4981450 B/op\t   81264 allocs/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8447292,
            "unit": "ns/op",
            "extra": "140 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981450,
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
            "value": 225.5,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5300438 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 225.5,
            "unit": "ns/op",
            "extra": "5300438 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5300438 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5300438 times\n4 procs"
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
            "value": 293.6,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "3867544 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 293.6,
            "unit": "ns/op",
            "extra": "3867544 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "3867544 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "3867544 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 450.3,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2728896 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 450.3,
            "unit": "ns/op",
            "extra": "2728896 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2728896 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2728896 times\n4 procs"
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
          "id": "7380bab46e95a27252ec8287662b1e6c80bdf013",
          "message": "test started\n\nSigned-off-by: Mattia Lavacca <lavacca.mattia@gmail.com>",
          "timestamp": "2024-12-05T18:21:03+01:00",
          "tree_id": "7e2de71af29c8d4109f11ccc9e3aed7100cd1d08",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/7380bab46e95a27252ec8287662b1e6c80bdf013"
        },
        "date": 1733419465927,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1196772,
            "unit": "ns/op\t  819447 B/op\t       5 allocs/op",
            "extra": "1087 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1196772,
            "unit": "ns/op",
            "extra": "1087 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819447,
            "unit": "B/op",
            "extra": "1087 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1087 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 8560,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "164786 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 8560,
            "unit": "ns/op",
            "extra": "164786 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "164786 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "164786 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.15,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15038785 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.15,
            "unit": "ns/op",
            "extra": "15038785 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15038785 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15038785 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22294,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52760 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22294,
            "unit": "ns/op",
            "extra": "52760 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52760 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52760 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 219290,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5812 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 219290,
            "unit": "ns/op",
            "extra": "5812 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5812 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5812 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2552702,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "415 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2552702,
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
            "value": 31106375,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "38 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 31106375,
            "unit": "ns/op",
            "extra": "38 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 9925,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "122442 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9925,
            "unit": "ns/op",
            "extra": "122442 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "122442 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "122442 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7305504,
            "unit": "ns/op\t 4594741 B/op\t   75254 allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7305504,
            "unit": "ns/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594741,
            "unit": "B/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75254,
            "unit": "allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8111,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "147000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8111,
            "unit": "ns/op",
            "extra": "147000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "147000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "147000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 853296,
            "unit": "ns/op\t  396626 B/op\t    6231 allocs/op",
            "extra": "1263 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 853296,
            "unit": "ns/op",
            "extra": "1263 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396626,
            "unit": "B/op",
            "extra": "1263 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6231,
            "unit": "allocs/op",
            "extra": "1263 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11185,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "107852 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11185,
            "unit": "ns/op",
            "extra": "107852 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "107852 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "107852 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8351984,
            "unit": "ns/op\t 4981412 B/op\t   81264 allocs/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8351984,
            "unit": "ns/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981412,
            "unit": "B/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81264,
            "unit": "allocs/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 227.9,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5307796 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 227.9,
            "unit": "ns/op",
            "extra": "5307796 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5307796 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5307796 times\n4 procs"
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
            "value": 298.7,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4011993 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 298.7,
            "unit": "ns/op",
            "extra": "4011993 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4011993 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4011993 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 452.6,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2672352 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 452.6,
            "unit": "ns/op",
            "extra": "2672352 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2672352 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2672352 times\n4 procs"
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
          "id": "e2640329f5d3dc600989ea7eb9c89d65f2e81ee0",
          "message": "restore id in ca cert\n\nSigned-off-by: Mattia Lavacca <lavacca.mattia@gmail.com>",
          "timestamp": "2024-12-05T18:24:23+01:00",
          "tree_id": "035132e10db738d6b93c8678ad5cd2c7a042c1fb",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/e2640329f5d3dc600989ea7eb9c89d65f2e81ee0"
        },
        "date": 1733419659298,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1455589,
            "unit": "ns/op\t  819442 B/op\t       5 allocs/op",
            "extra": "716 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1455589,
            "unit": "ns/op",
            "extra": "716 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819442,
            "unit": "B/op",
            "extra": "716 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "716 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 17713,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "61336 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 17713,
            "unit": "ns/op",
            "extra": "61336 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "61336 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "61336 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 86.23,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13026462 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 86.23,
            "unit": "ns/op",
            "extra": "13026462 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13026462 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13026462 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 23020,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "50869 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 23020,
            "unit": "ns/op",
            "extra": "50869 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "50869 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "50869 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 256708,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4472 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 256708,
            "unit": "ns/op",
            "extra": "4472 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4472 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4472 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2600055,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "433 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2600055,
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
            "value": 42711752,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 42711752,
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
            "value": 10549,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "114415 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10549,
            "unit": "ns/op",
            "extra": "114415 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "114415 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "114415 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 8032153,
            "unit": "ns/op\t 4594823 B/op\t   75254 allocs/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 8032153,
            "unit": "ns/op",
            "extra": "146 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594823,
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
            "value": 8754,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "136000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8754,
            "unit": "ns/op",
            "extra": "136000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "136000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "136000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 929552,
            "unit": "ns/op\t  396858 B/op\t    6235 allocs/op",
            "extra": "1191 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 929552,
            "unit": "ns/op",
            "extra": "1191 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396858,
            "unit": "B/op",
            "extra": "1191 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6235,
            "unit": "allocs/op",
            "extra": "1191 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 12319,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "96206 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 12319,
            "unit": "ns/op",
            "extra": "96206 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "96206 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "96206 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 9451424,
            "unit": "ns/op\t 4981280 B/op\t   81263 allocs/op",
            "extra": "127 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 9451424,
            "unit": "ns/op",
            "extra": "127 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981280,
            "unit": "B/op",
            "extra": "127 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81263,
            "unit": "allocs/op",
            "extra": "127 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 230.6,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5253590 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 230.6,
            "unit": "ns/op",
            "extra": "5253590 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5253590 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5253590 times\n4 procs"
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
            "value": 359.6,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "3994656 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 359.6,
            "unit": "ns/op",
            "extra": "3994656 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "3994656 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "3994656 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 453.2,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2518935 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 453.2,
            "unit": "ns/op",
            "extra": "2518935 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2518935 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2518935 times\n4 procs"
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
          "id": "092ed4cad4b29abcf8d0d3c5b33525b49602ee85",
          "message": "configmap for CA certificates complete\n\nSigned-off-by: Mattia Lavacca <lavacca.mattia@gmail.com>",
          "timestamp": "2024-12-06T13:44:45+01:00",
          "tree_id": "224415501cd9453f80581bbc5499c895cfd71171",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/092ed4cad4b29abcf8d0d3c5b33525b49602ee85"
        },
        "date": 1733489299058,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1335018,
            "unit": "ns/op\t  819442 B/op\t       5 allocs/op",
            "extra": "1042 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1335018,
            "unit": "ns/op",
            "extra": "1042 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819442,
            "unit": "B/op",
            "extra": "1042 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1042 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7308,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "164271 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7308,
            "unit": "ns/op",
            "extra": "164271 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "164271 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "164271 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 80.67,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15169645 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 80.67,
            "unit": "ns/op",
            "extra": "15169645 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15169645 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15169645 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22435,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53295 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22435,
            "unit": "ns/op",
            "extra": "53295 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53295 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53295 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 219159,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5610 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 219159,
            "unit": "ns/op",
            "extra": "5610 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5610 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5610 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2598739,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "427 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2598739,
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
            "value": 31654136,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "68 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 31654136,
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
            "value": 9809,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "120412 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9809,
            "unit": "ns/op",
            "extra": "120412 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "120412 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "120412 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7308799,
            "unit": "ns/op\t 4594931 B/op\t   75254 allocs/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7308799,
            "unit": "ns/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594931,
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
            "value": 8241,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "144459 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8241,
            "unit": "ns/op",
            "extra": "144459 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "144459 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "144459 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 854039,
            "unit": "ns/op\t  396672 B/op\t    6232 allocs/op",
            "extra": "1250 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 854039,
            "unit": "ns/op",
            "extra": "1250 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396672,
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
            "value": 11291,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "106360 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11291,
            "unit": "ns/op",
            "extra": "106360 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "106360 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "106360 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8384283,
            "unit": "ns/op\t 4981540 B/op\t   81264 allocs/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8384283,
            "unit": "ns/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981540,
            "unit": "B/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81264,
            "unit": "allocs/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 224.4,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5264815 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 224.4,
            "unit": "ns/op",
            "extra": "5264815 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5264815 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5264815 times\n4 procs"
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
            "value": 300.8,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4115122 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 300.8,
            "unit": "ns/op",
            "extra": "4115122 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4115122 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4115122 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 442,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2630365 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 442,
            "unit": "ns/op",
            "extra": "2630365 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2630365 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2630365 times\n4 procs"
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
          "id": "4da589718f7cc64b52c8c0c7d08958f5dc5810fd",
          "message": "ci: fix assigning duplicated instance names to plugins bound to different entity types (e.g. route and consumer group) (#6786)\n\n* ci: fix assigning duplicated instance names to plugins bound to different entity types (e.g. route and consumer group)\r\n\r\n* chore: refactor to use bytes.Buffer",
          "timestamp": "2024-12-06T12:46:24Z",
          "tree_id": "8601168441f540138b5905cb7b25bcbed1036810",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/4da589718f7cc64b52c8c0c7d08958f5dc5810fd"
        },
        "date": 1733489379266,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1280802,
            "unit": "ns/op\t  819442 B/op\t       5 allocs/op",
            "extra": "949 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1280802,
            "unit": "ns/op",
            "extra": "949 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819442,
            "unit": "B/op",
            "extra": "949 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "949 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 8525,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "159199 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 8525,
            "unit": "ns/op",
            "extra": "159199 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "159199 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "159199 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.45,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13488668 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.45,
            "unit": "ns/op",
            "extra": "13488668 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13488668 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13488668 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 23488,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52976 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 23488,
            "unit": "ns/op",
            "extra": "52976 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52976 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52976 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 221484,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5043 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 221484,
            "unit": "ns/op",
            "extra": "5043 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5043 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5043 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2585640,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "456 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2585640,
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
            "value": 34015846,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 34015846,
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
            "value": 9526,
            "unit": "ns/op\t    7504 B/op\t     175 allocs/op",
            "extra": "127946 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9526,
            "unit": "ns/op",
            "extra": "127946 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7504,
            "unit": "B/op",
            "extra": "127946 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 175,
            "unit": "allocs/op",
            "extra": "127946 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7367912,
            "unit": "ns/op\t 4594391 B/op\t   75247 allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7367912,
            "unit": "ns/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594391,
            "unit": "B/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75247,
            "unit": "allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7909,
            "unit": "ns/op\t    6240 B/op\t     159 allocs/op",
            "extra": "145198 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7909,
            "unit": "ns/op",
            "extra": "145198 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6240,
            "unit": "B/op",
            "extra": "145198 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 159,
            "unit": "allocs/op",
            "extra": "145198 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 846201,
            "unit": "ns/op\t  396317 B/op\t    6224 allocs/op",
            "extra": "1284 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 846201,
            "unit": "ns/op",
            "extra": "1284 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396317,
            "unit": "B/op",
            "extra": "1284 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6224,
            "unit": "allocs/op",
            "extra": "1284 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10873,
            "unit": "ns/op\t    7920 B/op\t     183 allocs/op",
            "extra": "111115 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10873,
            "unit": "ns/op",
            "extra": "111115 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 7920,
            "unit": "B/op",
            "extra": "111115 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "111115 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8369532,
            "unit": "ns/op\t 4981395 B/op\t   81259 allocs/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8369532,
            "unit": "ns/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981395,
            "unit": "B/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81259,
            "unit": "allocs/op",
            "extra": "144 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 229.7,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5244746 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 229.7,
            "unit": "ns/op",
            "extra": "5244746 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5244746 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5244746 times\n4 procs"
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
            "value": 295.3,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4056800 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 295.3,
            "unit": "ns/op",
            "extra": "4056800 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4056800 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4056800 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 454.8,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2681577 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 454.8,
            "unit": "ns/op",
            "extra": "2681577 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2681577 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2681577 times\n4 procs"
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
          "id": "a5941c2b2dafacf366690e93f0cd16e784664155",
          "message": "configmap for CA certificates complete\n\nSigned-off-by: Mattia Lavacca <lavacca.mattia@gmail.com>",
          "timestamp": "2024-12-06T13:47:18+01:00",
          "tree_id": "e289273c7a0f872970cf137f9837c3504b5e0699",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/a5941c2b2dafacf366690e93f0cd16e784664155"
        },
        "date": 1733489424681,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1174518,
            "unit": "ns/op\t  819441 B/op\t       5 allocs/op",
            "extra": "1029 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1174518,
            "unit": "ns/op",
            "extra": "1029 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819441,
            "unit": "B/op",
            "extra": "1029 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "1029 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7301,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "164716 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7301,
            "unit": "ns/op",
            "extra": "164716 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "164716 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "164716 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 78.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15210958 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 78.9,
            "unit": "ns/op",
            "extra": "15210958 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15210958 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15210958 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22577,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53198 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22577,
            "unit": "ns/op",
            "extra": "53198 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53198 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53198 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 234643,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4800 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 234643,
            "unit": "ns/op",
            "extra": "4800 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4800 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4800 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2584709,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "450 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2584709,
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
            "value": 41707210,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 41707210,
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
            "value": 10060,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "117637 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10060,
            "unit": "ns/op",
            "extra": "117637 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "117637 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "117637 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7490325,
            "unit": "ns/op\t 4594519 B/op\t   75253 allocs/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7490325,
            "unit": "ns/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4594519,
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
            "value": 8279,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "143000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8279,
            "unit": "ns/op",
            "extra": "143000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "143000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "143000 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 872492,
            "unit": "ns/op\t  396795 B/op\t    6234 allocs/op",
            "extra": "1212 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 872492,
            "unit": "ns/op",
            "extra": "1212 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396795,
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
            "value": 11536,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "104863 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11536,
            "unit": "ns/op",
            "extra": "104863 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "104863 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "104863 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 9030231,
            "unit": "ns/op\t 4981339 B/op\t   81263 allocs/op",
            "extra": "135 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 9030231,
            "unit": "ns/op",
            "extra": "135 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981339,
            "unit": "B/op",
            "extra": "135 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81263,
            "unit": "allocs/op",
            "extra": "135 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 238.4,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5323862 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 238.4,
            "unit": "ns/op",
            "extra": "5323862 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5323862 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5323862 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.6209,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.6209,
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
            "value": 294.7,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4080861 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 294.7,
            "unit": "ns/op",
            "extra": "4080861 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4080861 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4080861 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 537.8,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2678992 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 537.8,
            "unit": "ns/op",
            "extra": "2678992 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2678992 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2678992 times\n4 procs"
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
          "id": "e9d2474707522ae5b6cd08d2f48d01aa8d2bee05",
          "message": "configmap for CA certificates complete\n\nSigned-off-by: Mattia Lavacca <lavacca.mattia@gmail.com>",
          "timestamp": "2024-12-06T13:54:31+01:00",
          "tree_id": "1f7e6a55b075c35ad6d670718f7974a5301846d9",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/e9d2474707522ae5b6cd08d2f48d01aa8d2bee05"
        },
        "date": 1733489865419,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 1241613,
            "unit": "ns/op\t  819442 B/op\t       5 allocs/op",
            "extra": "872 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 1241613,
            "unit": "ns/op",
            "extra": "872 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819442,
            "unit": "B/op",
            "extra": "872 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "872 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 6884,
            "unit": "ns/op\t    7296 B/op\t      66 allocs/op",
            "extra": "174231 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 6884,
            "unit": "ns/op",
            "extra": "174231 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7296,
            "unit": "B/op",
            "extra": "174231 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "174231 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 79.88,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15182355 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 79.88,
            "unit": "ns/op",
            "extra": "15182355 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15182355 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15182355 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22436,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53914 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22436,
            "unit": "ns/op",
            "extra": "53914 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53914 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53914 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 224241,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4636 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 224241,
            "unit": "ns/op",
            "extra": "4636 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4636 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4636 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2546711,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "458 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2546711,
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
            "value": 34140172,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "64 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 34140172,
            "unit": "ns/op",
            "extra": "64 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
            "unit": "B/op",
            "extra": "64 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "64 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 10053,
            "unit": "ns/op\t    7736 B/op\t     181 allocs/op",
            "extra": "118830 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10053,
            "unit": "ns/op",
            "extra": "118830 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7736,
            "unit": "B/op",
            "extra": "118830 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 181,
            "unit": "allocs/op",
            "extra": "118830 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7540593,
            "unit": "ns/op\t 4595113 B/op\t   75255 allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7540593,
            "unit": "ns/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4595113,
            "unit": "B/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 75255,
            "unit": "allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8441,
            "unit": "ns/op\t    6472 B/op\t     165 allocs/op",
            "extra": "143844 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8441,
            "unit": "ns/op",
            "extra": "143844 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6472,
            "unit": "B/op",
            "extra": "143844 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 165,
            "unit": "allocs/op",
            "extra": "143844 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 879350,
            "unit": "ns/op\t  396724 B/op\t    6233 allocs/op",
            "extra": "1231 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 879350,
            "unit": "ns/op",
            "extra": "1231 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 396724,
            "unit": "B/op",
            "extra": "1231 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6233,
            "unit": "allocs/op",
            "extra": "1231 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11557,
            "unit": "ns/op\t    8152 B/op\t     189 allocs/op",
            "extra": "95793 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11557,
            "unit": "ns/op",
            "extra": "95793 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8152,
            "unit": "B/op",
            "extra": "95793 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 189,
            "unit": "allocs/op",
            "extra": "95793 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8617700,
            "unit": "ns/op\t 4981310 B/op\t   81263 allocs/op",
            "extra": "134 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8617700,
            "unit": "ns/op",
            "extra": "134 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4981310,
            "unit": "B/op",
            "extra": "134 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 81263,
            "unit": "allocs/op",
            "extra": "134 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 261.2,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5236641 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 261.2,
            "unit": "ns/op",
            "extra": "5236641 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5236641 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5236641 times\n4 procs"
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
            "value": 295,
            "unit": "ns/op\t     512 B/op\t       1 allocs/op",
            "extra": "4037038 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 295,
            "unit": "ns/op",
            "extra": "4037038 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 512,
            "unit": "B/op",
            "extra": "4037038 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "4037038 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 446.8,
            "unit": "ns/op\t     896 B/op\t       1 allocs/op",
            "extra": "2688888 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 446.8,
            "unit": "ns/op",
            "extra": "2688888 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 896,
            "unit": "B/op",
            "extra": "2688888 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2688888 times\n4 procs"
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
      }
    ]
  }
}