window.BENCHMARK_DATA = {
  "lastUpdate": 1793030118875,
  "repoUrl": "https://github.com/Kong/kubernetes-ingress-controller",
  "entries": {
    "Go Benchmark": [
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
          "id": "373a115559d6efef3f342c6a85b9ff5d1cf05e3a",
          "message": "chore: do not fill plugin config auto fields (#7995)\n\n* chore: bump go-kong to 792e7bd7ce957538c3d9162183434678a55d2c51\n\n* chore: fix kongintegration test TestUpdateStrategyInMemory_PropagatesResourcesErrors\n\n* tests: add auto fields plugin kongintegration tests\n\n* chore: revert go-kong to 410aea1a8763b30a6e006f0b5b1a7d01d6a42134\n\n* chore: use FillPluginsDefaultsWithOpts instead of changing FillPluginsDefaults, revert go-kong bump",
          "timestamp": "2026-06-22T13:49:27+02:00",
          "tree_id": "d5151ad37bc87dcf79be9a7f9c4fa6be48d536cc",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/373a115559d6efef3f342c6a85b9ff5d1cf05e3a"
        },
        "date": 1782129171815,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 898978,
            "unit": "ns/op\t  819427 B/op\t       4 allocs/op",
            "extra": "1231 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 898978,
            "unit": "ns/op",
            "extra": "1231 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819427,
            "unit": "B/op",
            "extra": "1231 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1231 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7974,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "148737 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7974,
            "unit": "ns/op",
            "extra": "148737 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "148737 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "148737 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 80.03,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14991696 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 80.03,
            "unit": "ns/op",
            "extra": "14991696 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14991696 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14991696 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22244,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52975 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22244,
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
            "value": 229802,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4593 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 229802,
            "unit": "ns/op",
            "extra": "4593 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4593 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4593 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2333351,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "506 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2333351,
            "unit": "ns/op",
            "extra": "506 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
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
            "value": 29501742,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "43 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 29501742,
            "unit": "ns/op",
            "extra": "43 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
            "unit": "B/op",
            "extra": "43 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "43 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9620,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "123589 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9620,
            "unit": "ns/op",
            "extra": "123589 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "123589 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "123589 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6660068,
            "unit": "ns/op\t 4335505 B/op\t   63229 allocs/op",
            "extra": "180 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6660068,
            "unit": "ns/op",
            "extra": "180 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335505,
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
            "value": 7949,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "148522 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7949,
            "unit": "ns/op",
            "extra": "148522 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "148522 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "148522 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 845944,
            "unit": "ns/op\t  392908 B/op\t    6173 allocs/op",
            "extra": "1358 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 845944,
            "unit": "ns/op",
            "extra": "1358 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392908,
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
            "value": 11085,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "108170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11085,
            "unit": "ns/op",
            "extra": "108170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "108170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "108170 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7519365,
            "unit": "ns/op\t 4722108 B/op\t   69239 allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7519365,
            "unit": "ns/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722108,
            "unit": "B/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 237.9,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5463121 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 237.9,
            "unit": "ns/op",
            "extra": "5463121 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5463121 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5463121 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "118432717 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.1,
            "unit": "ns/op",
            "extra": "118432717 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "118432717 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "118432717 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 518.1,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2317479 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 518.1,
            "unit": "ns/op",
            "extra": "2317479 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2317479 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2317479 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 735.4,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1630891 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 735.4,
            "unit": "ns/op",
            "extra": "1630891 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1630891 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1630891 times\n4 procs"
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
          "id": "9c3ef379e017a29d8090f17c6a5cd862f5a67916",
          "message": "tests: fix TestDiagnosticsServer_Diffs (#8001)",
          "timestamp": "2026-06-22T14:03:30+02:00",
          "tree_id": "9de62e7f7e42624dd34363a700260297245d9872",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/9c3ef379e017a29d8090f17c6a5cd862f5a67916"
        },
        "date": 1782130010810,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 915166,
            "unit": "ns/op\t  819430 B/op\t       4 allocs/op",
            "extra": "1215 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 915166,
            "unit": "ns/op",
            "extra": "1215 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819430,
            "unit": "B/op",
            "extra": "1215 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1215 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7840,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "149936 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7840,
            "unit": "ns/op",
            "extra": "149936 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "149936 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "149936 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.55,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16437196 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.55,
            "unit": "ns/op",
            "extra": "16437196 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16437196 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16437196 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21742,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54974 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21742,
            "unit": "ns/op",
            "extra": "54974 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54974 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54974 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 209365,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4932 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 209365,
            "unit": "ns/op",
            "extra": "4932 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4932 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4932 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2357768,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "445 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2357768,
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
            "value": 37669479,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "34 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 37669479,
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
            "value": 9038,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "132814 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9038,
            "unit": "ns/op",
            "extra": "132814 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "132814 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "132814 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7530652,
            "unit": "ns/op\t 4335501 B/op\t   63229 allocs/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7530652,
            "unit": "ns/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335501,
            "unit": "B/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7970,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "138382 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7970,
            "unit": "ns/op",
            "extra": "138382 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "138382 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "138382 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 830828,
            "unit": "ns/op\t  392925 B/op\t    6173 allocs/op",
            "extra": "1388 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 830828,
            "unit": "ns/op",
            "extra": "1388 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392925,
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
            "value": 10429,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "114386 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10429,
            "unit": "ns/op",
            "extra": "114386 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "114386 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "114386 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7321669,
            "unit": "ns/op\t 4722128 B/op\t   69239 allocs/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7321669,
            "unit": "ns/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722128,
            "unit": "B/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 212.9,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5651299 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 212.9,
            "unit": "ns/op",
            "extra": "5651299 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5651299 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5651299 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 9.987,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "120127114 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 9.987,
            "unit": "ns/op",
            "extra": "120127114 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "120127114 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "120127114 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 515.9,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2302144 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 515.9,
            "unit": "ns/op",
            "extra": "2302144 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2302144 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2302144 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 733.1,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1637174 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 733.1,
            "unit": "ns/op",
            "extra": "1637174 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1637174 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1637174 times\n4 procs"
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
          "id": "6051211002ac194ef36e7b2144a513654db05017",
          "message": "chore: update changelog, 3.5.10, 3.4.17 (#8012)",
          "timestamp": "2026-06-23T11:33:01+02:00",
          "tree_id": "443c916b5d89719421be13e6cc366bd443b6a93b",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/6051211002ac194ef36e7b2144a513654db05017"
        },
        "date": 1782207341783,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 288550,
            "unit": "ns/op\t  819428 B/op\t       4 allocs/op",
            "extra": "3702 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 288550,
            "unit": "ns/op",
            "extra": "3702 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819428,
            "unit": "B/op",
            "extra": "3702 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "3702 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 4143,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "242660 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 4143,
            "unit": "ns/op",
            "extra": "242660 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "242660 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "242660 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 40.26,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "29436015 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 40.26,
            "unit": "ns/op",
            "extra": "29436015 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "29436015 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "29436015 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 12453,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "96336 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 12453,
            "unit": "ns/op",
            "extra": "96336 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "96336 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "96336 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 119869,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "8928 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 119869,
            "unit": "ns/op",
            "extra": "8928 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "8928 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "8928 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 1587367,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "724 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 1587367,
            "unit": "ns/op",
            "extra": "724 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "724 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "724 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 36752514,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "36 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 36752514,
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
            "value": 5703,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "202678 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 5703,
            "unit": "ns/op",
            "extra": "202678 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "202678 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "202678 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 3803212,
            "unit": "ns/op\t 4335502 B/op\t   63229 allocs/op",
            "extra": "315 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 3803212,
            "unit": "ns/op",
            "extra": "315 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335502,
            "unit": "B/op",
            "extra": "315 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "315 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 5048,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "230992 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 5048,
            "unit": "ns/op",
            "extra": "230992 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "230992 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "230992 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 522084,
            "unit": "ns/op\t  392908 B/op\t    6173 allocs/op",
            "extra": "2112 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 522084,
            "unit": "ns/op",
            "extra": "2112 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392908,
            "unit": "B/op",
            "extra": "2112 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "2112 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 6701,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "180474 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 6701,
            "unit": "ns/op",
            "extra": "180474 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "180474 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "180474 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 4168929,
            "unit": "ns/op\t 4722123 B/op\t   69239 allocs/op",
            "extra": "285 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 4168929,
            "unit": "ns/op",
            "extra": "285 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722123,
            "unit": "B/op",
            "extra": "285 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "285 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 113.6,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "10866612 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 113.6,
            "unit": "ns/op",
            "extra": "10866612 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "10866612 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "10866612 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 6.823,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "185300205 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 6.823,
            "unit": "ns/op",
            "extra": "185300205 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "185300205 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "185300205 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 302.7,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "3883995 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 302.7,
            "unit": "ns/op",
            "extra": "3883995 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "3883995 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "3883995 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 417.9,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "2872602 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 417.9,
            "unit": "ns/op",
            "extra": "2872602 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "2872602 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2872602 times\n4 procs"
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
          "id": "9b4d8a7b589d25190334749f50bd78ad99f20826",
          "message": "chore(ci): add option to integration and e2e tests to designate the workflow to run with distroless images (#8000)\n\n* chore(ci): add option to integration and e2e tests to set the clear stale pids option in chart\n\n* chore: bump kong chart to 3.4.0\n\n* chore: fix wrong chart value\n\n* chore(ci): also add the option for isolated\n\n* chore(ci): adjust for Kong Gateway >= 3.15 (EE only)\n\n* chore: remove lifecycle hooks for distroless images\n\n* chore: adjust all in one postgres manifests to run with distroless images\n\n* ci: add TEST_KONG_DISTROLESS_IMAGE env instead of TEST_KONG_CLEAR_STALE_PIDS\n\n* chore: don't use wait image in integration tests\n\n* tests(e2e): remove bash from wait and migration jobs\n\n* chore: fix e2e\n\n* chore: more e2e logs\n\n* chore: replace migrations list command in e2e tests for distroless\n\n* tests: run tests with enterprise images and fix e2e tests patch\n\n---------\n\nCo-authored-by: Jakub Warczarek <jakub.warczarek@konghq.com>",
          "timestamp": "2026-06-24T12:38:43+02:00",
          "tree_id": "ef4d19603c2d384cae56d34aa56f3db6f7bb8dbd",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/9b4d8a7b589d25190334749f50bd78ad99f20826"
        },
        "date": 1782297725508,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 809919,
            "unit": "ns/op\t  819430 B/op\t       4 allocs/op",
            "extra": "1348 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 809919,
            "unit": "ns/op",
            "extra": "1348 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819430,
            "unit": "B/op",
            "extra": "1348 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1348 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 6971,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "152444 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 6971,
            "unit": "ns/op",
            "extra": "152444 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "152444 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "152444 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 81.67,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14648697 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 81.67,
            "unit": "ns/op",
            "extra": "14648697 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14648697 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14648697 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 23813,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "59916 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 23813,
            "unit": "ns/op",
            "extra": "59916 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "59916 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "59916 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 250898,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4603 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 250898,
            "unit": "ns/op",
            "extra": "4603 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
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
            "value": 2554424,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "480 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2554424,
            "unit": "ns/op",
            "extra": "480 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
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
            "value": 43459774,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 43459774,
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
            "value": 9034,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "131943 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9034,
            "unit": "ns/op",
            "extra": "131943 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "131943 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "131943 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6383268,
            "unit": "ns/op\t 4335501 B/op\t   63229 allocs/op",
            "extra": "186 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6383268,
            "unit": "ns/op",
            "extra": "186 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335501,
            "unit": "B/op",
            "extra": "186 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "186 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7551,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "157874 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7551,
            "unit": "ns/op",
            "extra": "157874 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "157874 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "157874 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 838341,
            "unit": "ns/op\t  392909 B/op\t    6173 allocs/op",
            "extra": "1380 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 838341,
            "unit": "ns/op",
            "extra": "1380 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392909,
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
            "value": 10774,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "109648 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10774,
            "unit": "ns/op",
            "extra": "109648 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "109648 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "109648 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7394996,
            "unit": "ns/op\t 4722134 B/op\t   69239 allocs/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7394996,
            "unit": "ns/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722134,
            "unit": "B/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 206.2,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5912863 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 206.2,
            "unit": "ns/op",
            "extra": "5912863 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5912863 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5912863 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 9.616,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "125595315 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 9.616,
            "unit": "ns/op",
            "extra": "125595315 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "125595315 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "125595315 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 527.8,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2272957 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 527.8,
            "unit": "ns/op",
            "extra": "2272957 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2272957 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2272957 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 917,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1358534 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 917,
            "unit": "ns/op",
            "extra": "1358534 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1358534 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1358534 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "sandro.modarelli@konghq.com",
            "name": "Sandro Modarelli",
            "username": "sandromodarelli"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "f82c38933659fab01ebaa8cbcc231416840a931b",
          "message": "chore: Update CHANGELOG with new version links (#8013)\n\n* Update CHANGELOG with new version links\n\nAdded new version links for 3.5.10, 3.5.9, 3.4.17, 3.4.16, and 3.4.15.\n\n* Fix duplicate version entry in CHANGELOG\n\nRemoved duplicate entry for version 3.5.8.",
          "timestamp": "2026-06-24T11:35:59Z",
          "tree_id": "ad14ebcd3508df455954f25dd9e3b7efc847c8d5",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/f82c38933659fab01ebaa8cbcc231416840a931b"
        },
        "date": 1782303119216,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 766749,
            "unit": "ns/op\t  819426 B/op\t       4 allocs/op",
            "extra": "1454 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 766749,
            "unit": "ns/op",
            "extra": "1454 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819426,
            "unit": "B/op",
            "extra": "1454 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1454 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 8351,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "154016 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 8351,
            "unit": "ns/op",
            "extra": "154016 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "154016 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "154016 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.86,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15581944 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.86,
            "unit": "ns/op",
            "extra": "15581944 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15581944 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15581944 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22059,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53403 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22059,
            "unit": "ns/op",
            "extra": "53403 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53403 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53403 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 226509,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4821 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 226509,
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
            "value": 2381212,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "494 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2381212,
            "unit": "ns/op",
            "extra": "494 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "494 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "494 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 40576364,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "31 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 40576364,
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
            "value": 11685,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "120326 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 11685,
            "unit": "ns/op",
            "extra": "120326 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "120326 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "120326 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7104114,
            "unit": "ns/op\t 4335504 B/op\t   63229 allocs/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7104114,
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
            "value": 8244,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "147213 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8244,
            "unit": "ns/op",
            "extra": "147213 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "147213 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "147213 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 879381,
            "unit": "ns/op\t  392924 B/op\t    6173 allocs/op",
            "extra": "1357 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 879381,
            "unit": "ns/op",
            "extra": "1357 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392924,
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
            "value": 11812,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "100177 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11812,
            "unit": "ns/op",
            "extra": "100177 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "100177 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "100177 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8820946,
            "unit": "ns/op\t 4722109 B/op\t   69239 allocs/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8820946,
            "unit": "ns/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722109,
            "unit": "B/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "138 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 217.7,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5496621 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 217.7,
            "unit": "ns/op",
            "extra": "5496621 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5496621 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5496621 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10,
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
            "value": 532.9,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2209087 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 532.9,
            "unit": "ns/op",
            "extra": "2209087 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2209087 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2209087 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 915.3,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1432094 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 915.3,
            "unit": "ns/op",
            "extra": "1432094 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1432094 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1432094 times\n4 procs"
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
          "id": "e8ba35e45b154b780d4a4d507329adbf2d057ad7",
          "message": "chore: bump GKE e2e version to 1.35.6 (#8017)\n\nSigned-off-by: Jintao Zhang <zhangjintao9020@gmail.com>",
          "timestamp": "2026-07-02T16:48:55+08:00",
          "tree_id": "8f31e0f5dcf1e78967977b4cf3f777d7cb65bac1",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/e8ba35e45b154b780d4a4d507329adbf2d057ad7"
        },
        "date": 1782982339282,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 809880,
            "unit": "ns/op\t  819431 B/op\t       4 allocs/op",
            "extra": "1358 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 809880,
            "unit": "ns/op",
            "extra": "1358 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819431,
            "unit": "B/op",
            "extra": "1358 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1358 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7861,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "139633 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7861,
            "unit": "ns/op",
            "extra": "139633 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "139633 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "139633 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.41,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16522058 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.41,
            "unit": "ns/op",
            "extra": "16522058 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16522058 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16522058 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22181,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53232 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22181,
            "unit": "ns/op",
            "extra": "53232 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53232 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53232 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 211025,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4958 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 211025,
            "unit": "ns/op",
            "extra": "4958 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4958 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4958 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2662849,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "451 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2662849,
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
            "value": 36173354,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "78 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 36173354,
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
            "value": 9313,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "127564 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9313,
            "unit": "ns/op",
            "extra": "127564 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "127564 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "127564 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6550266,
            "unit": "ns/op\t 4335504 B/op\t   63229 allocs/op",
            "extra": "182 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6550266,
            "unit": "ns/op",
            "extra": "182 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335504,
            "unit": "B/op",
            "extra": "182 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "182 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7725,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "152377 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7725,
            "unit": "ns/op",
            "extra": "152377 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "152377 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "152377 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 839940,
            "unit": "ns/op\t  392917 B/op\t    6173 allocs/op",
            "extra": "1345 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 839940,
            "unit": "ns/op",
            "extra": "1345 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392917,
            "unit": "B/op",
            "extra": "1345 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1345 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10703,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "110512 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10703,
            "unit": "ns/op",
            "extra": "110512 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "110512 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "110512 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7510694,
            "unit": "ns/op\t 4722142 B/op\t   69239 allocs/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7510694,
            "unit": "ns/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722142,
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
            "value": 218.1,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5446885 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 218.1,
            "unit": "ns/op",
            "extra": "5446885 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5446885 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5446885 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 9.984,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "120153954 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 9.984,
            "unit": "ns/op",
            "extra": "120153954 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "120153954 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "120153954 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 529.2,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2249503 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 529.2,
            "unit": "ns/op",
            "extra": "2249503 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2249503 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2249503 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 755.3,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1571410 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 755.3,
            "unit": "ns/op",
            "extra": "1571410 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1571410 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1571410 times\n4 procs"
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
          "id": "e47c3406628ff4442f9893eb9171b8152335fb6c",
          "message": "chore(deps): bump golang from 1.26.4 to 1.26.5 (#8020)\n\nBumps golang from 1.26.4 to 1.26.5.\n\n---\nupdated-dependencies:\n- dependency-name: golang\n  dependency-version: 1.26.5\n  dependency-type: direct:production\n  update-type: version-update:semver-patch\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2026-07-15T02:51:11Z",
          "tree_id": "76d9db24e00e3abdde715798c6fe5f294fc40c99",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/e47c3406628ff4442f9893eb9171b8152335fb6c"
        },
        "date": 1784086113386,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 808409,
            "unit": "ns/op\t  819430 B/op\t       4 allocs/op",
            "extra": "1288 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 808409,
            "unit": "ns/op",
            "extra": "1288 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819430,
            "unit": "B/op",
            "extra": "1288 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1288 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7589,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "156900 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7589,
            "unit": "ns/op",
            "extra": "156900 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "156900 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "156900 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 65.67,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18160110 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 65.67,
            "unit": "ns/op",
            "extra": "18160110 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "18160110 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "18160110 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21674,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53479 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21674,
            "unit": "ns/op",
            "extra": "53479 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53479 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53479 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 307102,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4378 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 307102,
            "unit": "ns/op",
            "extra": "4378 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4378 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4378 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2763543,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "436 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2763543,
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
            "value": 39569113,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "62 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 39569113,
            "unit": "ns/op",
            "extra": "62 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 9919,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "118810 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9919,
            "unit": "ns/op",
            "extra": "118810 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "118810 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "118810 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7184726,
            "unit": "ns/op\t 4335504 B/op\t   63229 allocs/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7184726,
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
            "value": 8418,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "135872 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8418,
            "unit": "ns/op",
            "extra": "135872 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "135872 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "135872 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 941848,
            "unit": "ns/op\t  392910 B/op\t    6173 allocs/op",
            "extra": "1226 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 941848,
            "unit": "ns/op",
            "extra": "1226 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392910,
            "unit": "B/op",
            "extra": "1226 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1226 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11944,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "100143 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11944,
            "unit": "ns/op",
            "extra": "100143 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "100143 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "100143 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8362225,
            "unit": "ns/op\t 4722126 B/op\t   69239 allocs/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8362225,
            "unit": "ns/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722126,
            "unit": "B/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "142 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 215.5,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5519191 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 215.5,
            "unit": "ns/op",
            "extra": "5519191 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5519191 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5519191 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 9.493,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "126652208 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 9.493,
            "unit": "ns/op",
            "extra": "126652208 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "126652208 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "126652208 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 543.7,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2192646 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 543.7,
            "unit": "ns/op",
            "extra": "2192646 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2192646 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2192646 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 752.5,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1603606 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 752.5,
            "unit": "ns/op",
            "extra": "1603606 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1603606 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1603606 times\n4 procs"
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
          "id": "610aeffe89f6358c1778d2ca479566e9a760f064",
          "message": "Merge commit from fork\n\n* register translation failures if duplicate IDs found in translating CA certificates\n\n* modify the getCACert to pick the first created one when ID conflicts",
          "timestamp": "2026-07-16T10:46:28+08:00",
          "tree_id": "82e384df5e8d0403a7adb9a7974107221ded419d",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/610aeffe89f6358c1778d2ca479566e9a760f064"
        },
        "date": 1784170168961,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 445420,
            "unit": "ns/op\t  819431 B/op\t       4 allocs/op",
            "extra": "2422 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 445420,
            "unit": "ns/op",
            "extra": "2422 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819431,
            "unit": "B/op",
            "extra": "2422 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "2422 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 5030,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "243531 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 5030,
            "unit": "ns/op",
            "extra": "243531 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "243531 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "243531 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 64.98,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18951210 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 64.98,
            "unit": "ns/op",
            "extra": "18951210 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "18951210 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "18951210 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 15561,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "75906 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 15561,
            "unit": "ns/op",
            "extra": "75906 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "75906 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "75906 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 168180,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "6511 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 168180,
            "unit": "ns/op",
            "extra": "6511 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "6511 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "6511 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2057898,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "590 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2057898,
            "unit": "ns/op",
            "extra": "590 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "590 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "590 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 33551006,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "92 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 33551006,
            "unit": "ns/op",
            "extra": "92 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
            "unit": "B/op",
            "extra": "92 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "92 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 8188,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "128761 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 8188,
            "unit": "ns/op",
            "extra": "128761 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "128761 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "128761 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 5169954,
            "unit": "ns/op\t 4335502 B/op\t   63229 allocs/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 5169954,
            "unit": "ns/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335502,
            "unit": "B/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "226 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 6259,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "196920 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 6259,
            "unit": "ns/op",
            "extra": "196920 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "196920 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "196920 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 662933,
            "unit": "ns/op\t  392921 B/op\t    6173 allocs/op",
            "extra": "1734 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 662933,
            "unit": "ns/op",
            "extra": "1734 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392921,
            "unit": "B/op",
            "extra": "1734 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1734 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 8594,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "142149 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 8594,
            "unit": "ns/op",
            "extra": "142149 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "142149 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "142149 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 6039342,
            "unit": "ns/op\t 4722104 B/op\t   69239 allocs/op",
            "extra": "199 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 6039342,
            "unit": "ns/op",
            "extra": "199 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722104,
            "unit": "B/op",
            "extra": "199 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "199 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 160.4,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "7558970 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 160.4,
            "unit": "ns/op",
            "extra": "7558970 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "7558970 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "7558970 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 7.476,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "161203144 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 7.476,
            "unit": "ns/op",
            "extra": "161203144 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "161203144 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "161203144 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 433.6,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2942552 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 433.6,
            "unit": "ns/op",
            "extra": "2942552 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2942552 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2942552 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 681.8,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1610394 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 681.8,
            "unit": "ns/op",
            "extra": "1610394 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1610394 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1610394 times\n4 procs"
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
          "id": "b88ba4e3e50128d7ed63e17c6d8e682081fcfd41",
          "message": "test(helper): Use license directly read from env (#8032)",
          "timestamp": "2026-07-28T13:21:30+02:00",
          "tree_id": "788c33edbd2e5ef384d2043208967084d08ab420",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/b88ba4e3e50128d7ed63e17c6d8e682081fcfd41"
        },
        "date": 1785237904276,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 756774,
            "unit": "ns/op\t  819430 B/op\t       4 allocs/op",
            "extra": "1387 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 756774,
            "unit": "ns/op",
            "extra": "1387 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819430,
            "unit": "B/op",
            "extra": "1387 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1387 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7657,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "146206 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7657,
            "unit": "ns/op",
            "extra": "146206 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "146206 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "146206 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.81,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16572640 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.81,
            "unit": "ns/op",
            "extra": "16572640 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16572640 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16572640 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21778,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54490 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21778,
            "unit": "ns/op",
            "extra": "54490 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54490 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54490 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 207918,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5112 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 207918,
            "unit": "ns/op",
            "extra": "5112 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
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
            "value": 2324325,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "459 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2324325,
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
            "value": 34190362,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "39 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 34190362,
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
            "value": 10636,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "130101 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10636,
            "unit": "ns/op",
            "extra": "130101 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "130101 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "130101 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7151320,
            "unit": "ns/op\t 4335501 B/op\t   63229 allocs/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7151320,
            "unit": "ns/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335501,
            "unit": "B/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7592,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "156596 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7592,
            "unit": "ns/op",
            "extra": "156596 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "156596 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "156596 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 829902,
            "unit": "ns/op\t  392913 B/op\t    6173 allocs/op",
            "extra": "1390 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 829902,
            "unit": "ns/op",
            "extra": "1390 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392913,
            "unit": "B/op",
            "extra": "1390 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1390 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10528,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "115500 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10528,
            "unit": "ns/op",
            "extra": "115500 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "115500 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "115500 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7406835,
            "unit": "ns/op\t 4722110 B/op\t   69239 allocs/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7406835,
            "unit": "ns/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722110,
            "unit": "B/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "160 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 214.8,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5628900 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 214.8,
            "unit": "ns/op",
            "extra": "5628900 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5628900 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5628900 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 9.976,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "120266774 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 9.976,
            "unit": "ns/op",
            "extra": "120266774 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "120266774 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "120266774 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 522.4,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2319619 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 522.4,
            "unit": "ns/op",
            "extra": "2319619 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2319619 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2319619 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 828.8,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1615851 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 828.8,
            "unit": "ns/op",
            "extra": "1615851 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1615851 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1615851 times\n4 procs"
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
          "id": "44c5f058506316b479260639d7b51ec52cbf0139",
          "message": "chore: bump go-kong v0.67.1 (#8034)\n\n* chore: bump go-kong v0.67.1\n\nSigned-off-by: Jintao Zhang <zhangjintao9020@gmail.com>\n\n* test: cover json-typed plugin config field preservation\n\nSigned-off-by: Jintao Zhang <zhangjintao9020@gmail.com>\n\n---------\n\nSigned-off-by: Jintao Zhang <zhangjintao9020@gmail.com>",
          "timestamp": "2026-07-29T16:44:31Z",
          "tree_id": "5353facb7dcae43d124efa96ed211dbe127d018e",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/44c5f058506316b479260639d7b51ec52cbf0139"
        },
        "date": 1785345705040,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 813926,
            "unit": "ns/op\t  819431 B/op\t       4 allocs/op",
            "extra": "1312 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 813926,
            "unit": "ns/op",
            "extra": "1312 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819431,
            "unit": "B/op",
            "extra": "1312 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1312 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 6830,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "174828 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 6830,
            "unit": "ns/op",
            "extra": "174828 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "174828 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "174828 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 81.23,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14759619 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 81.23,
            "unit": "ns/op",
            "extra": "14759619 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14759619 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14759619 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 19896,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "60230 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 19896,
            "unit": "ns/op",
            "extra": "60230 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "60230 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "60230 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 295667,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4419 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 295667,
            "unit": "ns/op",
            "extra": "4419 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4419 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4419 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2478744,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "484 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2478744,
            "unit": "ns/op",
            "extra": "484 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
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
            "value": 42777509,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 42777509,
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
            "value": 9412,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "126939 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9412,
            "unit": "ns/op",
            "extra": "126939 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "126939 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "126939 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6473836,
            "unit": "ns/op\t 4335504 B/op\t   63229 allocs/op",
            "extra": "183 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6473836,
            "unit": "ns/op",
            "extra": "183 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335504,
            "unit": "B/op",
            "extra": "183 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "183 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7859,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "154936 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7859,
            "unit": "ns/op",
            "extra": "154936 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "154936 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "154936 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 835924,
            "unit": "ns/op\t  392908 B/op\t    6173 allocs/op",
            "extra": "1357 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 835924,
            "unit": "ns/op",
            "extra": "1357 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392908,
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
            "value": 10937,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "109258 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10937,
            "unit": "ns/op",
            "extra": "109258 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "109258 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "109258 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7534028,
            "unit": "ns/op\t 4722091 B/op\t   69239 allocs/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7534028,
            "unit": "ns/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722091,
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
            "value": 205.8,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5885600 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 205.8,
            "unit": "ns/op",
            "extra": "5885600 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5885600 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5885600 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 9.512,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "126136051 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 9.512,
            "unit": "ns/op",
            "extra": "126136051 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "126136051 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "126136051 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 516.3,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2320848 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 516.3,
            "unit": "ns/op",
            "extra": "2320848 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2320848 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2320848 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 742.3,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1600797 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 742.3,
            "unit": "ns/op",
            "extra": "1600797 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1600797 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1600797 times\n4 procs"
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
          "id": "5ba3f1d08cbaf40a21ec3723bfcfc427aff6709c",
          "message": "docs: backport 3.5 release notes to main (#8042)\n\nSigned-off-by: Jintao Zhang <zhangjintao9020@gmail.com>",
          "timestamp": "2026-07-30T11:37:37Z",
          "tree_id": "077a3272005959bd18585d66140d978c032f0f5d",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/5ba3f1d08cbaf40a21ec3723bfcfc427aff6709c"
        },
        "date": 1785415251419,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 889218,
            "unit": "ns/op\t  819431 B/op\t       4 allocs/op",
            "extra": "1153 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 889218,
            "unit": "ns/op",
            "extra": "1153 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819431,
            "unit": "B/op",
            "extra": "1153 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1153 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7798,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "128511 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7798,
            "unit": "ns/op",
            "extra": "128511 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "128511 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "128511 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 71.48,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16652432 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 71.48,
            "unit": "ns/op",
            "extra": "16652432 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16652432 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16652432 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21968,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21968,
            "unit": "ns/op",
            "extra": "54100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54100 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 215526,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4995 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 215526,
            "unit": "ns/op",
            "extra": "4995 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
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
            "value": 2396109,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "472 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2396109,
            "unit": "ns/op",
            "extra": "472 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
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
            "value": 29098408,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "79 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 29098408,
            "unit": "ns/op",
            "extra": "79 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
            "unit": "B/op",
            "extra": "79 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "79 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9186,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "129308 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9186,
            "unit": "ns/op",
            "extra": "129308 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "129308 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "129308 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6637700,
            "unit": "ns/op\t 4335504 B/op\t   63229 allocs/op",
            "extra": "181 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6637700,
            "unit": "ns/op",
            "extra": "181 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335504,
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
            "value": 7657,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "145579 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7657,
            "unit": "ns/op",
            "extra": "145579 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "145579 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "145579 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 832530,
            "unit": "ns/op\t  392921 B/op\t    6173 allocs/op",
            "extra": "1401 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 832530,
            "unit": "ns/op",
            "extra": "1401 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392921,
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
            "value": 10582,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "110790 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10582,
            "unit": "ns/op",
            "extra": "110790 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "110790 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "110790 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7493194,
            "unit": "ns/op\t 4722125 B/op\t   69239 allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7493194,
            "unit": "ns/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722125,
            "unit": "B/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 216.6,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5564922 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 216.6,
            "unit": "ns/op",
            "extra": "5564922 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5564922 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5564922 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 9.993,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "120062532 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 9.993,
            "unit": "ns/op",
            "extra": "120062532 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "120062532 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "120062532 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 522,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2319421 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 522,
            "unit": "ns/op",
            "extra": "2319421 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2319421 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2319421 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 749.7,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1615503 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 749.7,
            "unit": "ns/op",
            "extra": "1615503 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1615503 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1615503 times\n4 procs"
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
          "id": "3b68f6dc2980a1cbd9c7f9f0d8ebb9c2f6f7c732",
          "message": "docs: backport 3.4 release notes to main (#8041)\n\nSigned-off-by: Jintao Zhang <zhangjintao9020@gmail.com>",
          "timestamp": "2026-07-30T12:56:12Z",
          "tree_id": "65b84cc3c945a15d32da170d99e570c5c9ad7d92",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/3b68f6dc2980a1cbd9c7f9f0d8ebb9c2f6f7c732"
        },
        "date": 1785418406471,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 805241,
            "unit": "ns/op\t  819430 B/op\t       4 allocs/op",
            "extra": "1242 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 805241,
            "unit": "ns/op",
            "extra": "1242 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819430,
            "unit": "B/op",
            "extra": "1242 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1242 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7146,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "164569 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7146,
            "unit": "ns/op",
            "extra": "164569 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "164569 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "164569 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 81.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14636088 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 81.3,
            "unit": "ns/op",
            "extra": "14636088 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14636088 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14636088 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 20405,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "57576 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 20405,
            "unit": "ns/op",
            "extra": "57576 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "57576 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "57576 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 222775,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4742 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 222775,
            "unit": "ns/op",
            "extra": "4742 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4742 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4742 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2602722,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "477 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2602722,
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
            "value": 38215874,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 38215874,
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
            "value": 10121,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "101827 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10121,
            "unit": "ns/op",
            "extra": "101827 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "101827 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "101827 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6601027,
            "unit": "ns/op\t 4335504 B/op\t   63229 allocs/op",
            "extra": "177 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6601027,
            "unit": "ns/op",
            "extra": "177 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335504,
            "unit": "B/op",
            "extra": "177 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "177 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7944,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "151538 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7944,
            "unit": "ns/op",
            "extra": "151538 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "151538 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "151538 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 851877,
            "unit": "ns/op\t  392923 B/op\t    6173 allocs/op",
            "extra": "1338 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 851877,
            "unit": "ns/op",
            "extra": "1338 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392923,
            "unit": "B/op",
            "extra": "1338 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1338 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11024,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "108924 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11024,
            "unit": "ns/op",
            "extra": "108924 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "108924 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "108924 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7423181,
            "unit": "ns/op\t 4722106 B/op\t   69239 allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7423181,
            "unit": "ns/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722106,
            "unit": "B/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 205,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5859555 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 205,
            "unit": "ns/op",
            "extra": "5859555 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5859555 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5859555 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 9.872,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "121564110 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 9.872,
            "unit": "ns/op",
            "extra": "121564110 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "121564110 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "121564110 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 580.8,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2260335 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 580.8,
            "unit": "ns/op",
            "extra": "2260335 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2260335 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2260335 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 865.9,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1316796 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 865.9,
            "unit": "ns/op",
            "extra": "1316796 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1316796 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1316796 times\n4 procs"
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
          "id": "d81ffd3bfdb4005f07163b9a6422f19b6a185bc5",
          "message": "test: fix Kong version detection for nightly images (#8046)\n\nSigned-off-by: Jintao Zhang <zhangjintao9020@gmail.com>",
          "timestamp": "2026-08-02T11:47:07Z",
          "tree_id": "d1855a2dc3b8e2ac999f685822ba4752f00acc52",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/d81ffd3bfdb4005f07163b9a6422f19b6a185bc5"
        },
        "date": 1785673445961,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 831659,
            "unit": "ns/op\t  819426 B/op\t       4 allocs/op",
            "extra": "1232 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 831659,
            "unit": "ns/op",
            "extra": "1232 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819426,
            "unit": "B/op",
            "extra": "1232 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1232 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7123,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "167449 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7123,
            "unit": "ns/op",
            "extra": "167449 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "167449 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "167449 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 81.07,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14734992 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 81.07,
            "unit": "ns/op",
            "extra": "14734992 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14734992 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14734992 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 20017,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "60380 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 20017,
            "unit": "ns/op",
            "extra": "60380 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "60380 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "60380 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 223057,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4840 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 223057,
            "unit": "ns/op",
            "extra": "4840 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4840 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4840 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2442969,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "427 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2442969,
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
            "value": 39275564,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "40 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 39275564,
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
            "value": 8992,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "132712 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 8992,
            "unit": "ns/op",
            "extra": "132712 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "132712 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "132712 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6218918,
            "unit": "ns/op\t 4335502 B/op\t   63229 allocs/op",
            "extra": "190 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6218918,
            "unit": "ns/op",
            "extra": "190 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335502,
            "unit": "B/op",
            "extra": "190 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "190 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7404,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "157987 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7404,
            "unit": "ns/op",
            "extra": "157987 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "157987 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "157987 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 818078,
            "unit": "ns/op\t  392915 B/op\t    6173 allocs/op",
            "extra": "1366 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 818078,
            "unit": "ns/op",
            "extra": "1366 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392915,
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
            "value": 10668,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "114367 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10668,
            "unit": "ns/op",
            "extra": "114367 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "114367 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "114367 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7119723,
            "unit": "ns/op\t 4722118 B/op\t   69239 allocs/op",
            "extra": "168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7119723,
            "unit": "ns/op",
            "extra": "168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722118,
            "unit": "B/op",
            "extra": "168 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "168 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 204,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5850188 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 204,
            "unit": "ns/op",
            "extra": "5850188 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5850188 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5850188 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 9.556,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "125413808 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 9.556,
            "unit": "ns/op",
            "extra": "125413808 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "125413808 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "125413808 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 518.1,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2301138 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 518.1,
            "unit": "ns/op",
            "extra": "2301138 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2301138 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2301138 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 728.8,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1633724 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 728.8,
            "unit": "ns/op",
            "extra": "1633724 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1633724 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1633724 times\n4 procs"
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
          "id": "bb6b1b7aedff2a21eb915670f46111e28256e3a7",
          "message": "fix(ci): use correct OSS nightly Kong version (#8047)\n\nSigned-off-by: Jintao Zhang <zhangjintao9020@gmail.com>",
          "timestamp": "2026-08-04T09:47:23Z",
          "tree_id": "a85cc59d33413083ec7a36ddc019e82c84675da0",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/bb6b1b7aedff2a21eb915670f46111e28256e3a7"
        },
        "date": 1785839035961,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 707548,
            "unit": "ns/op\t  819430 B/op\t       4 allocs/op",
            "extra": "1501 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 707548,
            "unit": "ns/op",
            "extra": "1501 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819430,
            "unit": "B/op",
            "extra": "1501 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1501 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7688,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "154938 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7688,
            "unit": "ns/op",
            "extra": "154938 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "154938 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "154938 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 71.31,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16753071 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 71.31,
            "unit": "ns/op",
            "extra": "16753071 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16753071 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16753071 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22266,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53672 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22266,
            "unit": "ns/op",
            "extra": "53672 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53672 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53672 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 214395,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5024 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 214395,
            "unit": "ns/op",
            "extra": "5024 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5024 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5024 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2792916,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "505 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2792916,
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
            "value": 44000305,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "34 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 44000305,
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
            "value": 9512,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "127018 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9512,
            "unit": "ns/op",
            "extra": "127018 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "127018 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "127018 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6793088,
            "unit": "ns/op\t 4335499 B/op\t   63229 allocs/op",
            "extra": "177 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6793088,
            "unit": "ns/op",
            "extra": "177 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335499,
            "unit": "B/op",
            "extra": "177 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "177 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7809,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "147958 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7809,
            "unit": "ns/op",
            "extra": "147958 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "147958 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "147958 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 832651,
            "unit": "ns/op\t  392896 B/op\t    6173 allocs/op",
            "extra": "1372 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 832651,
            "unit": "ns/op",
            "extra": "1372 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392896,
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
            "value": 10875,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "108974 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10875,
            "unit": "ns/op",
            "extra": "108974 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "108974 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "108974 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7775975,
            "unit": "ns/op\t 4722093 B/op\t   69239 allocs/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7775975,
            "unit": "ns/op",
            "extra": "152 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722093,
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
            "value": 216.9,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5476160 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 216.9,
            "unit": "ns/op",
            "extra": "5476160 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5476160 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5476160 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.08,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "120358237 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.08,
            "unit": "ns/op",
            "extra": "120358237 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "120358237 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "120358237 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 525.5,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2282164 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 525.5,
            "unit": "ns/op",
            "extra": "2282164 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2282164 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2282164 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 754.7,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1591336 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 754.7,
            "unit": "ns/op",
            "extra": "1591336 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1591336 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1591336 times\n4 procs"
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
          "id": "63153302233cd07e2574753f50ff0ddab8766e7f",
          "message": "chore(deps): bump golang from `079e598` to `d52df9c` (#8021)\n\n* chore(deps): bump golang from `079e598` to `d52df9c`\n\nBumps golang from `079e598` to `d52df9c`.\n\n---\nupdated-dependencies:\n- dependency-name: golang\n  dependency-version: 1.26.5\n  dependency-type: direct:production\n  update-type: version-update:semver-patch\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\n\n* chore: bump golang.org deps\n\n---------\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>\nCo-authored-by: Patryk Małek <patryk.malek@konghq.com>",
          "timestamp": "2026-08-04T17:43:42+02:00",
          "tree_id": "a54d800357ff6c8d67531799824c8aa73b43e275",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/63153302233cd07e2574753f50ff0ddab8766e7f"
        },
        "date": 1785858427910,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 510842,
            "unit": "ns/op\t  819432 B/op\t       4 allocs/op",
            "extra": "1969 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 510842,
            "unit": "ns/op",
            "extra": "1969 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819432,
            "unit": "B/op",
            "extra": "1969 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1969 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 5445,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "212940 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 5445,
            "unit": "ns/op",
            "extra": "212940 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "212940 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "212940 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 61.52,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19715941 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 61.52,
            "unit": "ns/op",
            "extra": "19715941 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "19715941 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "19715941 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 17736,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "67548 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 17736,
            "unit": "ns/op",
            "extra": "67548 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "67548 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "67548 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 216485,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "6019 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 216485,
            "unit": "ns/op",
            "extra": "6019 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "6019 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "6019 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2570910,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "454 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2570910,
            "unit": "ns/op",
            "extra": "454 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "454 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "454 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 41837821,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "56 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 41837821,
            "unit": "ns/op",
            "extra": "56 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 8026,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "148106 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 8026,
            "unit": "ns/op",
            "extra": "148106 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "148106 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "148106 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 5517441,
            "unit": "ns/op\t 4335506 B/op\t   63229 allocs/op",
            "extra": "216 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 5517441,
            "unit": "ns/op",
            "extra": "216 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335506,
            "unit": "B/op",
            "extra": "216 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "216 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 6750,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "173652 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 6750,
            "unit": "ns/op",
            "extra": "173652 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "173652 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "173652 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 746898,
            "unit": "ns/op\t  392918 B/op\t    6173 allocs/op",
            "extra": "1576 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 746898,
            "unit": "ns/op",
            "extra": "1576 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392918,
            "unit": "B/op",
            "extra": "1576 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1576 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 9403,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "128454 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 9403,
            "unit": "ns/op",
            "extra": "128454 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "128454 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "128454 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 6457393,
            "unit": "ns/op\t 4722110 B/op\t   69239 allocs/op",
            "extra": "184 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 6457393,
            "unit": "ns/op",
            "extra": "184 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722110,
            "unit": "B/op",
            "extra": "184 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "184 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 164,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "7012744 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 164,
            "unit": "ns/op",
            "extra": "7012744 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "7012744 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "7012744 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 11.58,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 11.58,
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
            "value": 448.1,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2672690 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 448.1,
            "unit": "ns/op",
            "extra": "2672690 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2672690 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2672690 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 613.2,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1930492 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 613.2,
            "unit": "ns/op",
            "extra": "1930492 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1930492 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1930492 times\n4 procs"
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
          "id": "77a8eed357df74e4ea9cde3f1288780f7a632b79",
          "message": "chore(deps): bump google.golang.org/grpc from 1.79.3 to 1.82.1 (#8031)\n\nBumps [google.golang.org/grpc](https://github.com/grpc/grpc-go) from 1.79.3 to 1.82.1.\n- [Release notes](https://github.com/grpc/grpc-go/releases)\n- [Commits](https://github.com/grpc/grpc-go/compare/v1.79.3...v1.82.1)\n\n---\nupdated-dependencies:\n- dependency-name: google.golang.org/grpc\n  dependency-version: 1.82.1\n  dependency-type: direct:production\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2026-08-05T04:21:35Z",
          "tree_id": "0e58e0b2b38ace30f0a8f12abd60102825d873f5",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/77a8eed357df74e4ea9cde3f1288780f7a632b79"
        },
        "date": 1785905967111,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 794250,
            "unit": "ns/op\t  819426 B/op\t       4 allocs/op",
            "extra": "1340 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 794250,
            "unit": "ns/op",
            "extra": "1340 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819426,
            "unit": "B/op",
            "extra": "1340 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1340 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7260,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "170110 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7260,
            "unit": "ns/op",
            "extra": "170110 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "170110 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "170110 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 81.74,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14738904 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 81.74,
            "unit": "ns/op",
            "extra": "14738904 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14738904 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14738904 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 23367,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "49425 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 23367,
            "unit": "ns/op",
            "extra": "49425 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "49425 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "49425 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 230250,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4760 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 230250,
            "unit": "ns/op",
            "extra": "4760 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4760 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4760 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2584554,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "459 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2584554,
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
            "value": 43892381,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "24 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 43892381,
            "unit": "ns/op",
            "extra": "24 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
            "unit": "B/op",
            "extra": "24 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "24 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9768,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "121629 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9768,
            "unit": "ns/op",
            "extra": "121629 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "121629 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "121629 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7322453,
            "unit": "ns/op\t 4335502 B/op\t   63229 allocs/op",
            "extra": "166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7322453,
            "unit": "ns/op",
            "extra": "166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335502,
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
            "value": 9475,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "128061 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 9475,
            "unit": "ns/op",
            "extra": "128061 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "128061 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "128061 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 862933,
            "unit": "ns/op\t  392913 B/op\t    6173 allocs/op",
            "extra": "1326 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 862933,
            "unit": "ns/op",
            "extra": "1326 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392913,
            "unit": "B/op",
            "extra": "1326 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1326 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11247,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "104569 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11247,
            "unit": "ns/op",
            "extra": "104569 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "104569 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "104569 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7709218,
            "unit": "ns/op\t 4722112 B/op\t   69239 allocs/op",
            "extra": "156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7709218,
            "unit": "ns/op",
            "extra": "156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722112,
            "unit": "B/op",
            "extra": "156 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "156 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 204,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5909676 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 204,
            "unit": "ns/op",
            "extra": "5909676 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5909676 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5909676 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 9.863,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "121660780 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 9.863,
            "unit": "ns/op",
            "extra": "121660780 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "121660780 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "121660780 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 527.2,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2292032 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 527.2,
            "unit": "ns/op",
            "extra": "2292032 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2292032 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2292032 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 722.1,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1658446 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 722.1,
            "unit": "ns/op",
            "extra": "1658446 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1658446 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1658446 times\n4 procs"
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
          "id": "3826a10dab78fd028e92648181d368391ad1e1a5",
          "message": "Use envs in commands instead of using input vars directly (#8050)\n\n* Use envs in commands instead of using input vars directly\n\n* fix actionslint",
          "timestamp": "2026-08-06T07:59:49Z",
          "tree_id": "ae2e9ea5b130067de8dec3e39b96f5203f344529",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/3826a10dab78fd028e92648181d368391ad1e1a5"
        },
        "date": 1786005409848,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 867884,
            "unit": "ns/op\t  819427 B/op\t       4 allocs/op",
            "extra": "1309 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 867884,
            "unit": "ns/op",
            "extra": "1309 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819427,
            "unit": "B/op",
            "extra": "1309 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1309 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7247,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "158794 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7247,
            "unit": "ns/op",
            "extra": "158794 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "158794 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "158794 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 81.33,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14751438 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 81.33,
            "unit": "ns/op",
            "extra": "14751438 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14751438 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14751438 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 23296,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "59622 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 23296,
            "unit": "ns/op",
            "extra": "59622 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "59622 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "59622 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 254281,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4300 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 254281,
            "unit": "ns/op",
            "extra": "4300 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4300 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4300 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2533783,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "463 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2533783,
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
            "value": 42851751,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "73 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 42851751,
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
            "value": 9321,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "127569 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9321,
            "unit": "ns/op",
            "extra": "127569 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "127569 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "127569 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6328455,
            "unit": "ns/op\t 4335502 B/op\t   63229 allocs/op",
            "extra": "189 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6328455,
            "unit": "ns/op",
            "extra": "189 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335502,
            "unit": "B/op",
            "extra": "189 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "189 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7608,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "155610 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7608,
            "unit": "ns/op",
            "extra": "155610 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "155610 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "155610 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 831836,
            "unit": "ns/op\t  392915 B/op\t    6173 allocs/op",
            "extra": "1394 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 831836,
            "unit": "ns/op",
            "extra": "1394 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392915,
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
            "value": 10794,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "106968 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10794,
            "unit": "ns/op",
            "extra": "106968 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "106968 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "106968 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7286000,
            "unit": "ns/op\t 4722100 B/op\t   69239 allocs/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7286000,
            "unit": "ns/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722100,
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
            "value": 204,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5973717 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 204,
            "unit": "ns/op",
            "extra": "5973717 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5973717 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5973717 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 9.862,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "121668260 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 9.862,
            "unit": "ns/op",
            "extra": "121668260 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "121668260 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "121668260 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 523.3,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2254592 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 523.3,
            "unit": "ns/op",
            "extra": "2254592 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2254592 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2254592 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 848.8,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1621668 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 848.8,
            "unit": "ns/op",
            "extra": "1621668 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1621668 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1621668 times\n4 procs"
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
          "id": "5e2f05ed96972ad9199bf852a69c3762cf9d80a4",
          "message": "chore(deps): update kindest/node@only-patch docker tag to v1.33.12 (#7969)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2026-08-06T09:06:48Z",
          "tree_id": "224cffee8c7f366e588e01a08f8889a1b902f1cd",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/5e2f05ed96972ad9199bf852a69c3762cf9d80a4"
        },
        "date": 1786009438701,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 865875,
            "unit": "ns/op\t  819427 B/op\t       4 allocs/op",
            "extra": "1177 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 865875,
            "unit": "ns/op",
            "extra": "1177 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819427,
            "unit": "B/op",
            "extra": "1177 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1177 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7925,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "149018 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7925,
            "unit": "ns/op",
            "extra": "149018 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "149018 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "149018 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 71.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16803825 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 71.4,
            "unit": "ns/op",
            "extra": "16803825 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16803825 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16803825 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 21985,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "54506 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 21985,
            "unit": "ns/op",
            "extra": "54506 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "54506 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "54506 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 226705,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4912 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 226705,
            "unit": "ns/op",
            "extra": "4912 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4912 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4912 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2357104,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "520 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2357104,
            "unit": "ns/op",
            "extra": "520 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "520 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "520 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 41007638,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "78 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 41007638,
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
            "value": 9679,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "122294 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9679,
            "unit": "ns/op",
            "extra": "122294 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "122294 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "122294 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7084796,
            "unit": "ns/op\t 4335505 B/op\t   63229 allocs/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7084796,
            "unit": "ns/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335505,
            "unit": "B/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "171 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8096,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "148236 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8096,
            "unit": "ns/op",
            "extra": "148236 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "148236 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "148236 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 879841,
            "unit": "ns/op\t  392911 B/op\t    6173 allocs/op",
            "extra": "1266 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 879841,
            "unit": "ns/op",
            "extra": "1266 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392911,
            "unit": "B/op",
            "extra": "1266 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1266 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11282,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "106906 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11282,
            "unit": "ns/op",
            "extra": "106906 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "106906 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "106906 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8013369,
            "unit": "ns/op\t 4722107 B/op\t   69239 allocs/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8013369,
            "unit": "ns/op",
            "extra": "148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722107,
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
            "value": 215.5,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5610883 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 215.5,
            "unit": "ns/op",
            "extra": "5610883 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5610883 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5610883 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 9.977,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "120279806 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 9.977,
            "unit": "ns/op",
            "extra": "120279806 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "120279806 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "120279806 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 535.6,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2219088 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 535.6,
            "unit": "ns/op",
            "extra": "2219088 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2219088 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2219088 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 760.9,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1569770 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 760.9,
            "unit": "ns/op",
            "extra": "1569770 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1569770 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1569770 times\n4 procs"
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
          "id": "ccf23701c53ba3114c534652129f47c799d77a99",
          "message": "chore(deps): update kong/kong-gateway docker tag to v3.12.0.9 (#7940)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2026-08-06T09:06:50Z",
          "tree_id": "91adae04d24a19e86029d93deaeccc2a664abed0",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/ccf23701c53ba3114c534652129f47c799d77a99"
        },
        "date": 1786009464477,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 760631,
            "unit": "ns/op\t  819430 B/op\t       4 allocs/op",
            "extra": "1473 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 760631,
            "unit": "ns/op",
            "extra": "1473 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819430,
            "unit": "B/op",
            "extra": "1473 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1473 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 6930,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "168274 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 6930,
            "unit": "ns/op",
            "extra": "168274 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "168274 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "168274 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 71.57,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16615455 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 71.57,
            "unit": "ns/op",
            "extra": "16615455 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16615455 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16615455 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 25651,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "55208 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 25651,
            "unit": "ns/op",
            "extra": "55208 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "55208 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "55208 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 262653,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "3884 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 262653,
            "unit": "ns/op",
            "extra": "3884 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "3884 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "3884 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2288515,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "514 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2288515,
            "unit": "ns/op",
            "extra": "514 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
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
            "value": 41047147,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "79 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 41047147,
            "unit": "ns/op",
            "extra": "79 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
            "unit": "B/op",
            "extra": "79 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "79 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9205,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "125526 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9205,
            "unit": "ns/op",
            "extra": "125526 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "125526 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "125526 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6756655,
            "unit": "ns/op\t 4335498 B/op\t   63229 allocs/op",
            "extra": "177 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6756655,
            "unit": "ns/op",
            "extra": "177 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335498,
            "unit": "B/op",
            "extra": "177 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "177 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7739,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "152295 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7739,
            "unit": "ns/op",
            "extra": "152295 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "152295 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "152295 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 862252,
            "unit": "ns/op\t  392922 B/op\t    6173 allocs/op",
            "extra": "1316 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 862252,
            "unit": "ns/op",
            "extra": "1316 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392922,
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
            "value": 10694,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "109417 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10694,
            "unit": "ns/op",
            "extra": "109417 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "109417 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "109417 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7797373,
            "unit": "ns/op\t 4722091 B/op\t   69239 allocs/op",
            "extra": "153 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7797373,
            "unit": "ns/op",
            "extra": "153 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722091,
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
            "value": 215.5,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5537349 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 215.5,
            "unit": "ns/op",
            "extra": "5537349 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5537349 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5537349 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 9.968,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "120386595 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 9.968,
            "unit": "ns/op",
            "extra": "120386595 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "120386595 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "120386595 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 520.2,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2347837 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 520.2,
            "unit": "ns/op",
            "extra": "2347837 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2347837 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2347837 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 885.4,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1494499 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 885.4,
            "unit": "ns/op",
            "extra": "1494499 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1494499 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1494499 times\n4 procs"
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
          "id": "532e71d8d412f831f5a47cd141a39574011e3856",
          "message": "chore(deps): bump github.com/kong/go-kong to v0.78.0, go-database-reconciler to v1.31.6 and sync CHANGELOG (#8060)",
          "timestamp": "2026-08-07T14:03:27+02:00",
          "tree_id": "bdc770c5a875d430ff2891acb304decf21f5b425",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/532e71d8d412f831f5a47cd141a39574011e3856"
        },
        "date": 1786104431478,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 877666,
            "unit": "ns/op\t  819431 B/op\t       4 allocs/op",
            "extra": "1227 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 877666,
            "unit": "ns/op",
            "extra": "1227 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819431,
            "unit": "B/op",
            "extra": "1227 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1227 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7233,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "158714 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7233,
            "unit": "ns/op",
            "extra": "158714 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "158714 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "158714 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 81.92,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14693445 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 81.92,
            "unit": "ns/op",
            "extra": "14693445 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14693445 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14693445 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 20166,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "58803 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 20166,
            "unit": "ns/op",
            "extra": "58803 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "58803 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "58803 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 286908,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4321 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 286908,
            "unit": "ns/op",
            "extra": "4321 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4321 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4321 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2723320,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "380 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2723320,
            "unit": "ns/op",
            "extra": "380 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "380 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "380 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 34897738,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "52 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 34897738,
            "unit": "ns/op",
            "extra": "52 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
            "unit": "B/op",
            "extra": "52 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9507,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "127150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9507,
            "unit": "ns/op",
            "extra": "127150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "127150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "127150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7232417,
            "unit": "ns/op\t 4335498 B/op\t   63229 allocs/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7232417,
            "unit": "ns/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335498,
            "unit": "B/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7920,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "152503 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7920,
            "unit": "ns/op",
            "extra": "152503 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "152503 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "152503 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 887723,
            "unit": "ns/op\t  392906 B/op\t    6173 allocs/op",
            "extra": "1333 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 887723,
            "unit": "ns/op",
            "extra": "1333 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392906,
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
            "value": 11305,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "104504 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11305,
            "unit": "ns/op",
            "extra": "104504 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "104504 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "104504 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 9171194,
            "unit": "ns/op\t 4722141 B/op\t   69239 allocs/op",
            "extra": "128 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 9171194,
            "unit": "ns/op",
            "extra": "128 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722141,
            "unit": "B/op",
            "extra": "128 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "128 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 206.4,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5872036 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 206.4,
            "unit": "ns/op",
            "extra": "5872036 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5872036 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5872036 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10,
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
            "value": 530.6,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2253490 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 530.6,
            "unit": "ns/op",
            "extra": "2253490 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2253490 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2253490 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 748.2,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1615089 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 748.2,
            "unit": "ns/op",
            "extra": "1615089 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1615089 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1615089 times\n4 procs"
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
          "id": "12df2910083c97ee02cf767a9f7ab71d8511bd6c",
          "message": "chore(deps): update helm release kong to v3.4.1 (#8057)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2026-08-07T12:40:02Z",
          "tree_id": "4dc8329cc7d5e7a2e623148e9d2233ad172f39f9",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/12df2910083c97ee02cf767a9f7ab71d8511bd6c"
        },
        "date": 1786108542325,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 841990,
            "unit": "ns/op\t  819430 B/op\t       4 allocs/op",
            "extra": "1310 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 841990,
            "unit": "ns/op",
            "extra": "1310 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819430,
            "unit": "B/op",
            "extra": "1310 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1310 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 6373,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "163380 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 6373,
            "unit": "ns/op",
            "extra": "163380 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "163380 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "163380 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 81.52,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14641867 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 81.52,
            "unit": "ns/op",
            "extra": "14641867 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14641867 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14641867 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 19897,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "59974 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 19897,
            "unit": "ns/op",
            "extra": "59974 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "59974 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "59974 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 207938,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5036 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 207938,
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
            "value": 2771815,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "436 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2771815,
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
            "value": 37389196,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 37389196,
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
            "value": 9016,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "129609 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9016,
            "unit": "ns/op",
            "extra": "129609 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "129609 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "129609 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6423791,
            "unit": "ns/op\t 4335501 B/op\t   63229 allocs/op",
            "extra": "187 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6423791,
            "unit": "ns/op",
            "extra": "187 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335501,
            "unit": "B/op",
            "extra": "187 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "187 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7441,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "155632 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7441,
            "unit": "ns/op",
            "extra": "155632 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "155632 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "155632 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 826644,
            "unit": "ns/op\t  392929 B/op\t    6173 allocs/op",
            "extra": "1408 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 826644,
            "unit": "ns/op",
            "extra": "1408 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392929,
            "unit": "B/op",
            "extra": "1408 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1408 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10430,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "115464 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10430,
            "unit": "ns/op",
            "extra": "115464 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "115464 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "115464 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7311671,
            "unit": "ns/op\t 4722105 B/op\t   69239 allocs/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7311671,
            "unit": "ns/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722105,
            "unit": "B/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 204.2,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5892661 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 204.2,
            "unit": "ns/op",
            "extra": "5892661 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5892661 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5892661 times\n4 procs"
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
            "value": 531.2,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2240210 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 531.2,
            "unit": "ns/op",
            "extra": "2240210 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2240210 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2240210 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 735.9,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1615782 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 735.9,
            "unit": "ns/op",
            "extra": "1615782 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1615782 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1615782 times\n4 procs"
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
          "id": "4cd9348b209a209639ad8e49cb078d95505f4aae",
          "message": "chore(deps): bump actions/cache from 4.2.3 to 6.0.0 (#7987)\n\nBumps [actions/cache](https://github.com/actions/cache) from 4.2.3 to 6.0.0.\n- [Release notes](https://github.com/actions/cache/releases)\n- [Changelog](https://github.com/actions/cache/blob/main/RELEASES.md)\n- [Commits](https://github.com/actions/cache/compare/5a3ec84eff668545956fd18022155c47e93e2684...2c8a9bd7457de244a408f35966fab2fb45fda9c8)\n\n---\nupdated-dependencies:\n- dependency-name: actions/cache\n  dependency-version: 5.0.5\n  dependency-type: direct:production\n  update-type: version-update:semver-major\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2026-08-10T08:32:22Z",
          "tree_id": "c7b218cb7a4315d869e2a10d958b8003aaf7cfb6",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/4cd9348b209a209639ad8e49cb078d95505f4aae"
        },
        "date": 1786352960593,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 866781,
            "unit": "ns/op\t  819426 B/op\t       4 allocs/op",
            "extra": "1252 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 866781,
            "unit": "ns/op",
            "extra": "1252 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819426,
            "unit": "B/op",
            "extra": "1252 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1252 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7385,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "140644 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7385,
            "unit": "ns/op",
            "extra": "140644 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "140644 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "140644 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 81.67,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14647166 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 81.67,
            "unit": "ns/op",
            "extra": "14647166 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14647166 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14647166 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 19869,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "60146 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 19869,
            "unit": "ns/op",
            "extra": "60146 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "60146 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "60146 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 221428,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4810 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 221428,
            "unit": "ns/op",
            "extra": "4810 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4810 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4810 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 3017980,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "466 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 3017980,
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
            "value": 35674725,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "70 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 35674725,
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
            "value": 8908,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "130453 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 8908,
            "unit": "ns/op",
            "extra": "130453 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "130453 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "130453 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6438104,
            "unit": "ns/op\t 4335504 B/op\t   63229 allocs/op",
            "extra": "188 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6438104,
            "unit": "ns/op",
            "extra": "188 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335504,
            "unit": "B/op",
            "extra": "188 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "188 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7408,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "159153 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7408,
            "unit": "ns/op",
            "extra": "159153 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "159153 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "159153 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 833533,
            "unit": "ns/op\t  392919 B/op\t    6173 allocs/op",
            "extra": "1404 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 833533,
            "unit": "ns/op",
            "extra": "1404 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392919,
            "unit": "B/op",
            "extra": "1404 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1404 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10599,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "107780 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10599,
            "unit": "ns/op",
            "extra": "107780 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "107780 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "107780 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7487033,
            "unit": "ns/op\t 4722102 B/op\t   69239 allocs/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7487033,
            "unit": "ns/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722102,
            "unit": "B/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 202.4,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "6001537 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 202.4,
            "unit": "ns/op",
            "extra": "6001537 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "6001537 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "6001537 times\n4 procs"
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
            "value": 534.3,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2243559 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 534.3,
            "unit": "ns/op",
            "extra": "2243559 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2243559 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2243559 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 740.2,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1620058 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 740.2,
            "unit": "ns/op",
            "extra": "1620058 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1620058 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1620058 times\n4 procs"
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
          "id": "21fec8ba2661bd6b327385e566d6b443e12a25e5",
          "message": "chore(deps): bump ruby/setup-ruby from 1.254.0 to 1.313.0 (#7984)\n\nBumps [ruby/setup-ruby](https://github.com/ruby/setup-ruby) from 1.254.0 to 1.313.0.\n- [Release notes](https://github.com/ruby/setup-ruby/releases)\n- [Changelog](https://github.com/ruby/setup-ruby/blob/master/release.rb)\n- [Commits](https://github.com/ruby/setup-ruby/compare/2a7b30092b0caf9c046252510f9273b4875f3db9...89f90524b88a01fe6e0b732220432cc6142926af)\n\n---\nupdated-dependencies:\n- dependency-name: ruby/setup-ruby\n  dependency-version: 1.313.0\n  dependency-type: direct:production\n  update-type: version-update:semver-minor\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2026-08-10T08:34:10Z",
          "tree_id": "98add05e98d04152e407ba7fc19f1a518d3b425c",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/21fec8ba2661bd6b327385e566d6b443e12a25e5"
        },
        "date": 1786357039407,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 752689,
            "unit": "ns/op\t  819430 B/op\t       4 allocs/op",
            "extra": "1465 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 752689,
            "unit": "ns/op",
            "extra": "1465 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819430,
            "unit": "B/op",
            "extra": "1465 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1465 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7201,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "166455 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7201,
            "unit": "ns/op",
            "extra": "166455 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "166455 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "166455 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 81.62,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14714667 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 81.62,
            "unit": "ns/op",
            "extra": "14714667 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14714667 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14714667 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 19936,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "59588 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 19936,
            "unit": "ns/op",
            "extra": "59588 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "59588 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "59588 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 270325,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5180 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 270325,
            "unit": "ns/op",
            "extra": "5180 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5180 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5180 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2610448,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "488 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2610448,
            "unit": "ns/op",
            "extra": "488 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "488 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "488 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 43037515,
            "unit": "ns/op\t24010753 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 43037515,
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
            "value": 9052,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "131992 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9052,
            "unit": "ns/op",
            "extra": "131992 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "131992 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "131992 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6429364,
            "unit": "ns/op\t 4335504 B/op\t   63229 allocs/op",
            "extra": "186 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6429364,
            "unit": "ns/op",
            "extra": "186 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335504,
            "unit": "B/op",
            "extra": "186 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "186 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7535,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "158512 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7535,
            "unit": "ns/op",
            "extra": "158512 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "158512 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "158512 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 826160,
            "unit": "ns/op\t  392918 B/op\t    6173 allocs/op",
            "extra": "1392 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 826160,
            "unit": "ns/op",
            "extra": "1392 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392918,
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
            "value": 10610,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "113594 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10610,
            "unit": "ns/op",
            "extra": "113594 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "113594 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "113594 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7306527,
            "unit": "ns/op\t 4722093 B/op\t   69239 allocs/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7306527,
            "unit": "ns/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722093,
            "unit": "B/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 206,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5783395 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 206,
            "unit": "ns/op",
            "extra": "5783395 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5783395 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5783395 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 10.25,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 10.25,
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
            "value": 531.1,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2240526 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 531.1,
            "unit": "ns/op",
            "extra": "2240526 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2240526 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2240526 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 723.2,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1647561 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 723.2,
            "unit": "ns/op",
            "extra": "1647561 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1647561 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1647561 times\n4 procs"
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
          "id": "1fa68f8642daef18fe676d852381a37447b31daf",
          "message": "chore(deps): bump golang from 1.26.5 to 1.26.6 (#8073)\n\nBumps golang from 1.26.5 to 1.26.6.\n\n---\nupdated-dependencies:\n- dependency-name: golang\n  dependency-version: 1.26.6\n  dependency-type: direct:production\n  update-type: version-update:semver-patch\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2026-08-18T13:34:00Z",
          "tree_id": "5775f731a3343f5dff761596275820ef8dad1e9e",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/1fa68f8642daef18fe676d852381a37447b31daf"
        },
        "date": 1787064481754,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 832260,
            "unit": "ns/op\t  819430 B/op\t       4 allocs/op",
            "extra": "1256 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 832260,
            "unit": "ns/op",
            "extra": "1256 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819430,
            "unit": "B/op",
            "extra": "1256 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1256 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7970,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "148455 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7970,
            "unit": "ns/op",
            "extra": "148455 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "148455 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "148455 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.56,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16516326 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.56,
            "unit": "ns/op",
            "extra": "16516326 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16516326 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16516326 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22118,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53409 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22118,
            "unit": "ns/op",
            "extra": "53409 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53409 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53409 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 216768,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4890 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 216768,
            "unit": "ns/op",
            "extra": "4890 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4890 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4890 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2375918,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "512 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2375918,
            "unit": "ns/op",
            "extra": "512 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
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
            "value": 33776960,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "38 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 33776960,
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
            "value": 9504,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "130587 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9504,
            "unit": "ns/op",
            "extra": "130587 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "130587 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "130587 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 8036676,
            "unit": "ns/op\t 4335497 B/op\t   63229 allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 8036676,
            "unit": "ns/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335497,
            "unit": "B/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "147 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7606,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "144590 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7606,
            "unit": "ns/op",
            "extra": "144590 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "144590 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "144590 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 846078,
            "unit": "ns/op\t  392917 B/op\t    6173 allocs/op",
            "extra": "1370 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 846078,
            "unit": "ns/op",
            "extra": "1370 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392917,
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
            "value": 10673,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "114693 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10673,
            "unit": "ns/op",
            "extra": "114693 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "114693 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "114693 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7541514,
            "unit": "ns/op\t 4722121 B/op\t   69239 allocs/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7541514,
            "unit": "ns/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722121,
            "unit": "B/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "157 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 215.8,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5542322 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 215.8,
            "unit": "ns/op",
            "extra": "5542322 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5542322 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5542322 times\n4 procs"
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
            "value": 521.3,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2297839 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 521.3,
            "unit": "ns/op",
            "extra": "2297839 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2297839 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2297839 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 734.9,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1624808 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 734.9,
            "unit": "ns/op",
            "extra": "1624808 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1624808 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1624808 times\n4 procs"
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
          "id": "6822654d015651b93b4e52f2530aec5b58709c7f",
          "message": "chore(deps): update kindest/node@only-patch docker tag to v1.34.11 (#7974)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2026-08-27T11:12:28Z",
          "tree_id": "fd9e5ddb9da64f485ce54ce5ea0aead779af3fca",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/6822654d015651b93b4e52f2530aec5b58709c7f"
        },
        "date": 1787831349854,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 825215,
            "unit": "ns/op\t  819431 B/op\t       4 allocs/op",
            "extra": "1260 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 825215,
            "unit": "ns/op",
            "extra": "1260 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819431,
            "unit": "B/op",
            "extra": "1260 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1260 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7462,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "152839 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7462,
            "unit": "ns/op",
            "extra": "152839 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "152839 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "152839 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 81.66,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14624161 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 81.66,
            "unit": "ns/op",
            "extra": "14624161 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14624161 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14624161 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 19727,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "60525 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 19727,
            "unit": "ns/op",
            "extra": "60525 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "60525 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "60525 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 206929,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5160 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 206929,
            "unit": "ns/op",
            "extra": "5160 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
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
            "value": 2451635,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "476 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2451635,
            "unit": "ns/op",
            "extra": "476 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
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
            "value": 35929517,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "30 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 35929517,
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
            "value": 10531,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "128137 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10531,
            "unit": "ns/op",
            "extra": "128137 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "128137 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "128137 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6826170,
            "unit": "ns/op\t 4335502 B/op\t   63229 allocs/op",
            "extra": "166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6826170,
            "unit": "ns/op",
            "extra": "166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335502,
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
            "value": 7765,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "154137 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7765,
            "unit": "ns/op",
            "extra": "154137 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "154137 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "154137 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 832765,
            "unit": "ns/op\t  392911 B/op\t    6173 allocs/op",
            "extra": "1365 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 832765,
            "unit": "ns/op",
            "extra": "1365 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392911,
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
            "value": 10868,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "110096 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10868,
            "unit": "ns/op",
            "extra": "110096 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "110096 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "110096 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7354338,
            "unit": "ns/op\t 4722110 B/op\t   69239 allocs/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7354338,
            "unit": "ns/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722110,
            "unit": "B/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 204.4,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5954500 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 204.4,
            "unit": "ns/op",
            "extra": "5954500 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5954500 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5954500 times\n4 procs"
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
            "value": 519.7,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2312606 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 519.7,
            "unit": "ns/op",
            "extra": "2312606 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2312606 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2312606 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 791.9,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1615797 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 791.9,
            "unit": "ns/op",
            "extra": "1615797 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1615797 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1615797 times\n4 procs"
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
          "id": "34769e3af242f07c72d7ea2b591179940cd6f5f4",
          "message": "chore(deps): update dependency go-delve/delve to v1.27.2 (#8077)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2026-09-11T06:52:57Z",
          "tree_id": "58478f5a0a1197098e47fdc00ec9df01a8c7264f",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/34769e3af242f07c72d7ea2b591179940cd6f5f4"
        },
        "date": 1789111763744,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 791823,
            "unit": "ns/op\t  819426 B/op\t       4 allocs/op",
            "extra": "1350 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 791823,
            "unit": "ns/op",
            "extra": "1350 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819426,
            "unit": "B/op",
            "extra": "1350 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1350 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 9250,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "156160 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 9250,
            "unit": "ns/op",
            "extra": "156160 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "156160 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "156160 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.52,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16399724 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.52,
            "unit": "ns/op",
            "extra": "16399724 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16399724 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16399724 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22300,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53210 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22300,
            "unit": "ns/op",
            "extra": "53210 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53210 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53210 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 216510,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4828 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 216510,
            "unit": "ns/op",
            "extra": "4828 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4828 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4828 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2485711,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "470 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2485711,
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
            "value": 41483332,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "75 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 41483332,
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
            "value": 9676,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "121796 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9676,
            "unit": "ns/op",
            "extra": "121796 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "121796 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "121796 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6850518,
            "unit": "ns/op\t 4335506 B/op\t   63229 allocs/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6850518,
            "unit": "ns/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335506,
            "unit": "B/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "172 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8042,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "149205 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8042,
            "unit": "ns/op",
            "extra": "149205 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "149205 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "149205 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 860761,
            "unit": "ns/op\t  392923 B/op\t    6173 allocs/op",
            "extra": "1351 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 860761,
            "unit": "ns/op",
            "extra": "1351 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392923,
            "unit": "B/op",
            "extra": "1351 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1351 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11271,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "109450 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11271,
            "unit": "ns/op",
            "extra": "109450 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "109450 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "109450 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7762429,
            "unit": "ns/op\t 4722115 B/op\t   69239 allocs/op",
            "extra": "153 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7762429,
            "unit": "ns/op",
            "extra": "153 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722115,
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
            "value": 216,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5457618 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 216,
            "unit": "ns/op",
            "extra": "5457618 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5457618 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5457618 times\n4 procs"
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
            "value": 526.3,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2278210 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 526.3,
            "unit": "ns/op",
            "extra": "2278210 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2278210 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2278210 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 732.2,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1633722 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 732.2,
            "unit": "ns/op",
            "extra": "1633722 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1633722 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1633722 times\n4 procs"
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
          "id": "7805dd575f9de1f7b6ca5fe829e9e9412beafffb",
          "message": "bump GKE to 1.37.0 (#8087)",
          "timestamp": "2026-09-11T07:10:09Z",
          "tree_id": "aa07425d2a14ee3596101e810ab8d0acab1e0232",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/7805dd575f9de1f7b6ca5fe829e9e9412beafffb"
        },
        "date": 1789112793337,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 920778,
            "unit": "ns/op\t  819431 B/op\t       4 allocs/op",
            "extra": "1102 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 920778,
            "unit": "ns/op",
            "extra": "1102 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819431,
            "unit": "B/op",
            "extra": "1102 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1102 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7838,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "152691 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7838,
            "unit": "ns/op",
            "extra": "152691 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "152691 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "152691 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 74,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16552146 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 74,
            "unit": "ns/op",
            "extra": "16552146 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16552146 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16552146 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22097,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53946 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22097,
            "unit": "ns/op",
            "extra": "53946 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53946 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53946 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 233743,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4882 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 233743,
            "unit": "ns/op",
            "extra": "4882 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4882 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4882 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2666731,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "420 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2666731,
            "unit": "ns/op",
            "extra": "420 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "420 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "420 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 41223216,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "74 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 41223216,
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
            "value": 9615,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "126325 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9615,
            "unit": "ns/op",
            "extra": "126325 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "126325 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "126325 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6531508,
            "unit": "ns/op\t 4335505 B/op\t   63229 allocs/op",
            "extra": "182 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6531508,
            "unit": "ns/op",
            "extra": "182 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335505,
            "unit": "B/op",
            "extra": "182 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "182 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7477,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "161978 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7477,
            "unit": "ns/op",
            "extra": "161978 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "161978 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "161978 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 848482,
            "unit": "ns/op\t  392907 B/op\t    6173 allocs/op",
            "extra": "1374 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 848482,
            "unit": "ns/op",
            "extra": "1374 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392907,
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
            "value": 10382,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "115023 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10382,
            "unit": "ns/op",
            "extra": "115023 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "115023 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "115023 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7492729,
            "unit": "ns/op\t 4722107 B/op\t   69239 allocs/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7492729,
            "unit": "ns/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722107,
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
            "value": 213,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5590042 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 213,
            "unit": "ns/op",
            "extra": "5590042 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5590042 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5590042 times\n4 procs"
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
            "value": 513.9,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2335624 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 513.9,
            "unit": "ns/op",
            "extra": "2335624 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2335624 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2335624 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 746.5,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1606069 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 746.5,
            "unit": "ns/op",
            "extra": "1606069 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1606069 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1606069 times\n4 procs"
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
          "id": "5af2417fbb6aa42446c97e94069f609236c37c0d",
          "message": "fix(translator): set tags to the winning cert in merging certificates (#8085)\n\n* fix: set tags to the winning cert in merging certificates\n\n* add test cases and changelog\n\n* address comments",
          "timestamp": "2026-09-15T07:27:17Z",
          "tree_id": "2ad38ccc3e309aa2cd6ff60186c8b26bdbac318a",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/5af2417fbb6aa42446c97e94069f609236c37c0d"
        },
        "date": 1789459492085,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 848639,
            "unit": "ns/op\t  819430 B/op\t       4 allocs/op",
            "extra": "1299 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 848639,
            "unit": "ns/op",
            "extra": "1299 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819430,
            "unit": "B/op",
            "extra": "1299 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1299 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7752,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "146089 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7752,
            "unit": "ns/op",
            "extra": "146089 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "146089 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "146089 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.51,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16508883 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.51,
            "unit": "ns/op",
            "extra": "16508883 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16508883 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16508883 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 27936,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "52796 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 27936,
            "unit": "ns/op",
            "extra": "52796 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "52796 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "52796 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 227855,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5088 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 227855,
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
            "value": 2406676,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "454 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2406676,
            "unit": "ns/op",
            "extra": "454 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "454 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "454 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 31272160,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "72 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 31272160,
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
            "value": 9128,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "131542 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9128,
            "unit": "ns/op",
            "extra": "131542 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "131542 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "131542 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6550364,
            "unit": "ns/op\t 4335501 B/op\t   63229 allocs/op",
            "extra": "182 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6550364,
            "unit": "ns/op",
            "extra": "182 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335501,
            "unit": "B/op",
            "extra": "182 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "182 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7559,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "160670 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7559,
            "unit": "ns/op",
            "extra": "160670 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "160670 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "160670 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 841293,
            "unit": "ns/op\t  392915 B/op\t    6173 allocs/op",
            "extra": "1368 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 841293,
            "unit": "ns/op",
            "extra": "1368 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392915,
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
            "value": 10451,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "114538 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10451,
            "unit": "ns/op",
            "extra": "114538 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "114538 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "114538 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7508237,
            "unit": "ns/op\t 4722095 B/op\t   69239 allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7508237,
            "unit": "ns/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722095,
            "unit": "B/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 216,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5576548 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 216,
            "unit": "ns/op",
            "extra": "5576548 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5576548 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5576548 times\n4 procs"
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
            "value": 530,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2266927 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 530,
            "unit": "ns/op",
            "extra": "2266927 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2266927 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2266927 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 767.7,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1605120 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 767.7,
            "unit": "ns/op",
            "extra": "1605120 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1605120 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1605120 times\n4 procs"
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
          "id": "0d1f8fb63113bfe9614d1c03bce2123a4f56cf51",
          "message": "Fix(translator): split httproute match by path when they have PathPrefixRewrite filter (#8081)\n\n* split httproute match by path when they have PathPrefixRewrite filter\n\n* update changelog\n\n* fix routes for requestRedirect filter\n\n* update changelog\n\n* update changelog",
          "timestamp": "2026-09-16T02:51:38Z",
          "tree_id": "eae670d089199c0b6b30700a5995cc9bbc475b1f",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/0d1f8fb63113bfe9614d1c03bce2123a4f56cf51"
        },
        "date": 1789529399824,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 773153,
            "unit": "ns/op\t  819431 B/op\t       4 allocs/op",
            "extra": "1329 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 773153,
            "unit": "ns/op",
            "extra": "1329 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819431,
            "unit": "B/op",
            "extra": "1329 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1329 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7233,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "164640 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7233,
            "unit": "ns/op",
            "extra": "164640 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "164640 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "164640 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 81.57,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14712067 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 81.57,
            "unit": "ns/op",
            "extra": "14712067 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14712067 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14712067 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 26712,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "60369 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 26712,
            "unit": "ns/op",
            "extra": "60369 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "60369 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "60369 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 204732,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5186 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 204732,
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
            "value": 2490395,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "466 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2490395,
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
            "value": 38938729,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "26 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 38938729,
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
            "value": 9302,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "129308 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9302,
            "unit": "ns/op",
            "extra": "129308 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "129308 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "129308 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6444255,
            "unit": "ns/op\t 4335501 B/op\t   63229 allocs/op",
            "extra": "182 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6444255,
            "unit": "ns/op",
            "extra": "182 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335501,
            "unit": "B/op",
            "extra": "182 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "182 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7655,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "157358 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7655,
            "unit": "ns/op",
            "extra": "157358 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "157358 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "157358 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 830561,
            "unit": "ns/op\t  392917 B/op\t    6173 allocs/op",
            "extra": "1364 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 830561,
            "unit": "ns/op",
            "extra": "1364 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392917,
            "unit": "B/op",
            "extra": "1364 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1364 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10879,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "112449 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10879,
            "unit": "ns/op",
            "extra": "112449 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "112449 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "112449 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7259439,
            "unit": "ns/op\t 4722107 B/op\t   69239 allocs/op",
            "extra": "165 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7259439,
            "unit": "ns/op",
            "extra": "165 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722107,
            "unit": "B/op",
            "extra": "165 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "165 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 204.4,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5863335 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 204.4,
            "unit": "ns/op",
            "extra": "5863335 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5863335 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5863335 times\n4 procs"
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
            "value": 520.8,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2289415 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 520.8,
            "unit": "ns/op",
            "extra": "2289415 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2289415 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2289415 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 733.8,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1647577 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 733.8,
            "unit": "ns/op",
            "extra": "1647577 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1647577 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1647577 times\n4 procs"
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
          "id": "ee03b2f4fcd0fc1f857df7fe36b1880e9729ca45",
          "message": "chore(deps): Bump grpc to 1.83.1 (#8090)\n\n* chore(deps): bump google.golang.org/grpc from 1.82.1 to 1.83.1\n\nBumps [google.golang.org/grpc](https://github.com/grpc/grpc-go) from 1.82.1 to 1.83.1.\n- [Release notes](https://github.com/grpc/grpc-go/releases)\n- [Commits](https://github.com/grpc/grpc-go/compare/v1.82.1...v1.83.1)\n\n---\nupdated-dependencies:\n- dependency-name: google.golang.org/grpc\n  dependency-version: 1.83.1\n  dependency-type: direct:production\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\n\n* replace deprecated WithCredentialsJSON\n\n---------\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2026-09-17T02:20:16Z",
          "tree_id": "a526844dae685c2c277ce0e819b012906ae20775",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/ee03b2f4fcd0fc1f857df7fe36b1880e9729ca45"
        },
        "date": 1789613849244,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 465272,
            "unit": "ns/op\t  819431 B/op\t       4 allocs/op",
            "extra": "2319 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 465272,
            "unit": "ns/op",
            "extra": "2319 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819431,
            "unit": "B/op",
            "extra": "2319 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "2319 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 4862,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "245706 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 4862,
            "unit": "ns/op",
            "extra": "245706 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "245706 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "245706 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 63.34,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18891543 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 63.34,
            "unit": "ns/op",
            "extra": "18891543 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "18891543 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "18891543 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 15768,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "76069 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 15768,
            "unit": "ns/op",
            "extra": "76069 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "76069 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "76069 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 182757,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5953 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 182757,
            "unit": "ns/op",
            "extra": "5953 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5953 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5953 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2414901,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "457 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2414901,
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
            "value": 31380297,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "91 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 31380297,
            "unit": "ns/op",
            "extra": "91 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
            "unit": "B/op",
            "extra": "91 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "91 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 7233,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "161852 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 7233,
            "unit": "ns/op",
            "extra": "161852 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "161852 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "161852 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 5144525,
            "unit": "ns/op\t 4335501 B/op\t   63229 allocs/op",
            "extra": "235 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 5144525,
            "unit": "ns/op",
            "extra": "235 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335501,
            "unit": "B/op",
            "extra": "235 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "235 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 6102,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "197492 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 6102,
            "unit": "ns/op",
            "extra": "197492 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "197492 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "197492 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 659618,
            "unit": "ns/op\t  392912 B/op\t    6173 allocs/op",
            "extra": "1740 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 659618,
            "unit": "ns/op",
            "extra": "1740 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392912,
            "unit": "B/op",
            "extra": "1740 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1740 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 8803,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "138429 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 8803,
            "unit": "ns/op",
            "extra": "138429 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "138429 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "138429 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 6161149,
            "unit": "ns/op\t 4722101 B/op\t   69239 allocs/op",
            "extra": "188 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 6161149,
            "unit": "ns/op",
            "extra": "188 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722101,
            "unit": "B/op",
            "extra": "188 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "188 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 159.3,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "7359676 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 159.3,
            "unit": "ns/op",
            "extra": "7359676 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "7359676 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "7359676 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 7.935,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "151145322 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 7.935,
            "unit": "ns/op",
            "extra": "151145322 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "151145322 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "151145322 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 404.3,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2938011 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 404.3,
            "unit": "ns/op",
            "extra": "2938011 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2938011 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2938011 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 572.1,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "2056011 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 572.1,
            "unit": "ns/op",
            "extra": "2056011 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "2056011 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2056011 times\n4 procs"
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
          "id": "7d3925aca63bfc2e4357125a1f4bbcde5e8191b7",
          "message": "chore(deps): bump google.golang.org/grpc from 1.83.1 to 1.83.2 (#8091)\n\nBumps [google.golang.org/grpc](https://github.com/grpc/grpc-go) from 1.83.1 to 1.83.2.\n- [Release notes](https://github.com/grpc/grpc-go/releases)\n- [Commits](https://github.com/grpc/grpc-go/compare/v1.83.1...v1.83.2)\n\n---\nupdated-dependencies:\n- dependency-name: google.golang.org/grpc\n  dependency-version: 1.83.2\n  dependency-type: direct:production\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2026-09-17T06:03:37Z",
          "tree_id": "c12763c4bb14e6d8af8b1774bf83429f2153790f",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/7d3925aca63bfc2e4357125a1f4bbcde5e8191b7"
        },
        "date": 1789627252796,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 833113,
            "unit": "ns/op\t  819430 B/op\t       4 allocs/op",
            "extra": "1230 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 833113,
            "unit": "ns/op",
            "extra": "1230 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819430,
            "unit": "B/op",
            "extra": "1230 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1230 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7366,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "160912 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7366,
            "unit": "ns/op",
            "extra": "160912 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "160912 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "160912 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.58,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16509471 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.58,
            "unit": "ns/op",
            "extra": "16509471 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16509471 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16509471 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22939,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "51652 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22939,
            "unit": "ns/op",
            "extra": "51652 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "51652 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "51652 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 236029,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4528 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 236029,
            "unit": "ns/op",
            "extra": "4528 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4528 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4528 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2524456,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "440 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2524456,
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
            "value": 44698418,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "34 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 44698418,
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
            "value": 12246,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "96310 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 12246,
            "unit": "ns/op",
            "extra": "96310 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "96310 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "96310 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 7254131,
            "unit": "ns/op\t 4335502 B/op\t   63229 allocs/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 7254131,
            "unit": "ns/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335502,
            "unit": "B/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "163 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 8223,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "148789 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 8223,
            "unit": "ns/op",
            "extra": "148789 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "148789 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "148789 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 886807,
            "unit": "ns/op\t  392922 B/op\t    6173 allocs/op",
            "extra": "1293 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 886807,
            "unit": "ns/op",
            "extra": "1293 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392922,
            "unit": "B/op",
            "extra": "1293 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1293 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11553,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "101971 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11553,
            "unit": "ns/op",
            "extra": "101971 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "101971 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "101971 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 8898247,
            "unit": "ns/op\t 4722125 B/op\t   69239 allocs/op",
            "extra": "134 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 8898247,
            "unit": "ns/op",
            "extra": "134 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722125,
            "unit": "B/op",
            "extra": "134 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "134 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 219.5,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5465646 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 219.5,
            "unit": "ns/op",
            "extra": "5465646 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5465646 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5465646 times\n4 procs"
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
            "value": 529.8,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2272542 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 529.8,
            "unit": "ns/op",
            "extra": "2272542 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2272542 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2272542 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 923.8,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1420957 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 923.8,
            "unit": "ns/op",
            "extra": "1420957 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1420957 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1420957 times\n4 procs"
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
          "id": "5f2cbbc424a6fb985bd411bc1cf5ea3d8450e681",
          "message": "chore(deps): bump go.opentelemetry.io/otel/exporters/otlp/otlptrace (#8095)\n\nBumps [go.opentelemetry.io/otel/exporters/otlp/otlptrace](https://github.com/open-telemetry/opentelemetry-go) from 1.33.0 to 1.45.0.\n- [Release notes](https://github.com/open-telemetry/opentelemetry-go/releases)\n- [Changelog](https://github.com/open-telemetry/opentelemetry-go/blob/main/CHANGELOG.md)\n- [Commits](https://github.com/open-telemetry/opentelemetry-go/compare/v1.33.0...v1.45.0)\n\n---\nupdated-dependencies:\n- dependency-name: go.opentelemetry.io/otel/exporters/otlp/otlptrace\n  dependency-version: 1.45.0\n  dependency-type: indirect\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2026-09-18T08:30:53Z",
          "tree_id": "7b35d49bd27f5976f76d17e5b4d680131e4c9f82",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/5f2cbbc424a6fb985bd411bc1cf5ea3d8450e681"
        },
        "date": 1789722395392,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 632636,
            "unit": "ns/op\t  819430 B/op\t       4 allocs/op",
            "extra": "1708 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 632636,
            "unit": "ns/op",
            "extra": "1708 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819430,
            "unit": "B/op",
            "extra": "1708 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1708 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 6926,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "174298 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 6926,
            "unit": "ns/op",
            "extra": "174298 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "174298 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "174298 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 71.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16777723 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 71.8,
            "unit": "ns/op",
            "extra": "16777723 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16777723 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16777723 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 20162,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "58903 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 20162,
            "unit": "ns/op",
            "extra": "58903 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "58903 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "58903 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 214892,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4712 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 214892,
            "unit": "ns/op",
            "extra": "4712 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4712 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4712 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 3276779,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "414 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 3276779,
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
            "value": 63236146,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "16 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 63236146,
            "unit": "ns/op",
            "extra": "16 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
            "unit": "B/op",
            "extra": "16 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "16 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 9037,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "130972 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9037,
            "unit": "ns/op",
            "extra": "130972 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "130972 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "130972 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6173461,
            "unit": "ns/op\t 4335504 B/op\t   63229 allocs/op",
            "extra": "193 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6173461,
            "unit": "ns/op",
            "extra": "193 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335504,
            "unit": "B/op",
            "extra": "193 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "193 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7445,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "157740 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7445,
            "unit": "ns/op",
            "extra": "157740 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "157740 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "157740 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 843483,
            "unit": "ns/op\t  392909 B/op\t    6173 allocs/op",
            "extra": "1396 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 843483,
            "unit": "ns/op",
            "extra": "1396 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392909,
            "unit": "B/op",
            "extra": "1396 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1396 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 10606,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "112693 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10606,
            "unit": "ns/op",
            "extra": "112693 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "112693 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "112693 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7222742,
            "unit": "ns/op\t 4722104 B/op\t   69239 allocs/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7222742,
            "unit": "ns/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722104,
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
            "value": 186.8,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "6391353 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 186.8,
            "unit": "ns/op",
            "extra": "6391353 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "6391353 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "6391353 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 13.66,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "90305758 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 13.66,
            "unit": "ns/op",
            "extra": "90305758 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "90305758 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "90305758 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups",
            "value": 516.9,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2299972 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 516.9,
            "unit": "ns/op",
            "extra": "2299972 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2299972 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2299972 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 723.7,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1650748 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 723.7,
            "unit": "ns/op",
            "extra": "1650748 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1650748 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1650748 times\n4 procs"
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
          "id": "b6da030fde88a6a64f6ba42696ae666001a2e049",
          "message": "fix(test): Pass kong-effective-version to e2e in validating Kong image tests and fix GKEtest failures (#8089)\n\n* pass kong-effective-version to e2e\n\n* fix istio workflow\n\n* list only live pods\n\n* Use e2-standard-4 node to fix OOMs",
          "timestamp": "2026-09-20T02:51:20Z",
          "tree_id": "9a992bdc8f70b63238183f893f41bfbf23418440",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/b6da030fde88a6a64f6ba42696ae666001a2e049"
        },
        "date": 1789874889072,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 702795,
            "unit": "ns/op\t  819429 B/op\t       4 allocs/op",
            "extra": "1516 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 702795,
            "unit": "ns/op",
            "extra": "1516 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819429,
            "unit": "B/op",
            "extra": "1516 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1516 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7943,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "157047 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7943,
            "unit": "ns/op",
            "extra": "157047 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "157047 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "157047 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 70.63,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16169310 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 70.63,
            "unit": "ns/op",
            "extra": "16169310 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16169310 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16169310 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 20401,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "58231 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 20401,
            "unit": "ns/op",
            "extra": "58231 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "58231 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "58231 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 223462,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4688 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 223462,
            "unit": "ns/op",
            "extra": "4688 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4688 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4688 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2938147,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "402 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2938147,
            "unit": "ns/op",
            "extra": "402 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "402 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "402 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 45140191,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "62 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 45140191,
            "unit": "ns/op",
            "extra": "62 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
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
            "value": 9492,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "129975 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9492,
            "unit": "ns/op",
            "extra": "129975 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "129975 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "129975 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6460075,
            "unit": "ns/op\t 4335502 B/op\t   63229 allocs/op",
            "extra": "186 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6460075,
            "unit": "ns/op",
            "extra": "186 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335502,
            "unit": "B/op",
            "extra": "186 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "186 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7638,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "153501 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7638,
            "unit": "ns/op",
            "extra": "153501 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "153501 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "153501 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 881015,
            "unit": "ns/op\t  392935 B/op\t    6173 allocs/op",
            "extra": "1332 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 881015,
            "unit": "ns/op",
            "extra": "1332 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392935,
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
            "value": 10665,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "112689 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 10665,
            "unit": "ns/op",
            "extra": "112689 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "112689 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "112689 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7495199,
            "unit": "ns/op\t 4722106 B/op\t   69239 allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7495199,
            "unit": "ns/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722106,
            "unit": "B/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 186.3,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "6448328 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 186.3,
            "unit": "ns/op",
            "extra": "6448328 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "6448328 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "6448328 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 12.67,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 12.67,
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
            "value": 515.8,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2281629 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 515.8,
            "unit": "ns/op",
            "extra": "2281629 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2281629 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2281629 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 721.9,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1653775 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 721.9,
            "unit": "ns/op",
            "extra": "1653775 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1653775 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1653775 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "95225168+matteogastaldello@users.noreply.github.com",
            "name": "Matteo Gastaldello",
            "username": "matteogastaldello"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "40fe2b46cad79d6371f4b0efea4c9792e39071b4",
          "message": "chore(ci): use GitHub App tokens instead of the team-k8s-bot PAT (#8099)\n\n* chore(ci): mint GitHub App tokens instead of the team-k8s-bot PAT\n\n* chore(ci): mint a GitHub App token for backports",
          "timestamp": "2026-10-08T08:07:39Z",
          "tree_id": "1c940b22b6b7f138526e0e9c81fa9c267602f1c6",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/40fe2b46cad79d6371f4b0efea4c9792e39071b4"
        },
        "date": 1791449087672,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 809937,
            "unit": "ns/op\t  819430 B/op\t       4 allocs/op",
            "extra": "1354 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 809937,
            "unit": "ns/op",
            "extra": "1354 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819430,
            "unit": "B/op",
            "extra": "1354 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1354 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7498,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "158923 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7498,
            "unit": "ns/op",
            "extra": "158923 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "158923 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "158923 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 65.76,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18183972 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 65.76,
            "unit": "ns/op",
            "extra": "18183972 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "18183972 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "18183972 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 26968,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "38472 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 26968,
            "unit": "ns/op",
            "extra": "38472 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "38472 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "38472 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 236757,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4512 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 236757,
            "unit": "ns/op",
            "extra": "4512 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4512 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4512 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2672727,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "447 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2672727,
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
            "value": 45395791,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "64 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 45395791,
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
            "value": 11781,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "123589 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 11781,
            "unit": "ns/op",
            "extra": "123589 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "123589 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "123589 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6961186,
            "unit": "ns/op\t 4335501 B/op\t   63229 allocs/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6961186,
            "unit": "ns/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335501,
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
            "value": 7990,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "148527 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7990,
            "unit": "ns/op",
            "extra": "148527 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "148527 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "148527 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 901476,
            "unit": "ns/op\t  392915 B/op\t    6173 allocs/op",
            "extra": "1300 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 901476,
            "unit": "ns/op",
            "extra": "1300 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392915,
            "unit": "B/op",
            "extra": "1300 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1300 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11267,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "105751 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11267,
            "unit": "ns/op",
            "extra": "105751 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "105751 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "105751 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7921958,
            "unit": "ns/op\t 4722106 B/op\t   69239 allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7921958,
            "unit": "ns/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722106,
            "unit": "B/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "150 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 213.7,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5626297 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 213.7,
            "unit": "ns/op",
            "extra": "5626297 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5626297 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5626297 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 11.36,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "100000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 11.36,
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
            "value": 542.3,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2198769 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 542.3,
            "unit": "ns/op",
            "extra": "2198769 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2198769 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2198769 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 904,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1422068 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 904,
            "unit": "ns/op",
            "extra": "1422068 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1422068 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1422068 times\n4 procs"
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
          "id": "3a3ad8a6a0307b0e4fb271d41d08cdcd0af727c7",
          "message": "chore(deps): update kong docker tag to v3.9.4 (#8058)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2026-10-14T08:51:35Z",
          "tree_id": "bbc7605f76cbfc7040e0ae82ecf84080193212e6",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/3a3ad8a6a0307b0e4fb271d41d08cdcd0af727c7"
        },
        "date": 1791970174038,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 780580,
            "unit": "ns/op\t  819430 B/op\t       4 allocs/op",
            "extra": "1332 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 780580,
            "unit": "ns/op",
            "extra": "1332 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819430,
            "unit": "B/op",
            "extra": "1332 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1332 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7921,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "150517 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7921,
            "unit": "ns/op",
            "extra": "150517 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "150517 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "150517 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 72.63,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16532186 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 72.63,
            "unit": "ns/op",
            "extra": "16532186 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16532186 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16532186 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 22306,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "53283 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 22306,
            "unit": "ns/op",
            "extra": "53283 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "53283 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "53283 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 218170,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4879 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 218170,
            "unit": "ns/op",
            "extra": "4879 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4879 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4879 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 2341269,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "508 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 2341269,
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
            "value": 32263316,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "76 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 32263316,
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
            "value": 9797,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "121938 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9797,
            "unit": "ns/op",
            "extra": "121938 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "121938 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "121938 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6791263,
            "unit": "ns/op\t 4335500 B/op\t   63229 allocs/op",
            "extra": "175 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6791263,
            "unit": "ns/op",
            "extra": "175 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335500,
            "unit": "B/op",
            "extra": "175 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "175 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 7988,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "147056 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 7988,
            "unit": "ns/op",
            "extra": "147056 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "147056 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "147056 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 863177,
            "unit": "ns/op\t  392913 B/op\t    6173 allocs/op",
            "extra": "1326 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 863177,
            "unit": "ns/op",
            "extra": "1326 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392913,
            "unit": "B/op",
            "extra": "1326 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1326 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 11204,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "107162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 11204,
            "unit": "ns/op",
            "extra": "107162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "107162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "107162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7827878,
            "unit": "ns/op\t 4722093 B/op\t   69239 allocs/op",
            "extra": "153 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7827878,
            "unit": "ns/op",
            "extra": "153 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722093,
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
            "value": 228.1,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5211399 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 228.1,
            "unit": "ns/op",
            "extra": "5211399 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5211399 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5211399 times\n4 procs"
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
            "value": 606,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2327462 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 606,
            "unit": "ns/op",
            "extra": "2327462 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2327462 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2327462 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 874,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1327977 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 874,
            "unit": "ns/op",
            "extra": "1327977 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1327977 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1327977 times\n4 procs"
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
          "id": "39f668259a1c12b1662529ffcc58f6a843e7dd8c",
          "message": "bump golang version to 1.27.2 (#8110)\n\n* bump golang version to 1.27.2\n\n* bump mordenize and run fixes",
          "timestamp": "2026-10-16T02:28:46Z",
          "tree_id": "485b2d91e8722435f350765f82e3ce854f6616f9",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/39f668259a1c12b1662529ffcc58f6a843e7dd8c"
        },
        "date": 1792119979719,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 910877,
            "unit": "ns/op\t  819426 B/op\t       4 allocs/op",
            "extra": "1270 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 910877,
            "unit": "ns/op",
            "extra": "1270 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819426,
            "unit": "B/op",
            "extra": "1270 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1270 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7752,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "153409 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7752,
            "unit": "ns/op",
            "extra": "153409 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "153409 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "153409 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 59.69,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "20108136 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 59.69,
            "unit": "ns/op",
            "extra": "20108136 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "20108136 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "20108136 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 24368,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "49410 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 24368,
            "unit": "ns/op",
            "extra": "49410 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "49410 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "49410 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 204376,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5128 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 204376,
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
            "value": 1883583,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "565 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 1883583,
            "unit": "ns/op",
            "extra": "565 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "565 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "565 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 26006249,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "76 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 26006249,
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
            "value": 9645,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "142820 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 9645,
            "unit": "ns/op",
            "extra": "142820 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "142820 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "142820 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6862341,
            "unit": "ns/op\t 4335503 B/op\t   63229 allocs/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6862341,
            "unit": "ns/op",
            "extra": "169 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335503,
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
            "value": 6601,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "177734 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 6601,
            "unit": "ns/op",
            "extra": "177734 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "177734 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "177734 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 824343,
            "unit": "ns/op\t  392913 B/op\t    6173 allocs/op",
            "extra": "1395 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 824343,
            "unit": "ns/op",
            "extra": "1395 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392913,
            "unit": "B/op",
            "extra": "1395 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1395 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 9492,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "126644 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 9492,
            "unit": "ns/op",
            "extra": "126644 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "126644 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "126644 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7378168,
            "unit": "ns/op\t 4722109 B/op\t   69239 allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7378168,
            "unit": "ns/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722109,
            "unit": "B/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "162 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 244.3,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "4891402 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 244.3,
            "unit": "ns/op",
            "extra": "4891402 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "4891402 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4891402 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 1.004,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 1.004,
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
            "value": 628,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "1894261 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 628,
            "unit": "ns/op",
            "extra": "1894261 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "1894261 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1894261 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 996.4,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1320087 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 996.4,
            "unit": "ns/op",
            "extra": "1320087 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1320087 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1320087 times\n4 procs"
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
          "id": "39edf7f0a3492ba29aca80c504072835dbd77332",
          "message": "chore(deps): update dependency gke to v1.37.1 (#8112)\n\nCo-authored-by: renovate[bot] <29139614+renovate[bot]@users.noreply.github.com>",
          "timestamp": "2026-10-16T08:59:43Z",
          "tree_id": "d0f9a59d9786503c3e5ad83d85f182068413f9d8",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/39edf7f0a3492ba29aca80c504072835dbd77332"
        },
        "date": 1792143269586,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 847745,
            "unit": "ns/op\t  819431 B/op\t       4 allocs/op",
            "extra": "1255 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 847745,
            "unit": "ns/op",
            "extra": "1255 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819431,
            "unit": "B/op",
            "extra": "1255 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1255 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 7996,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "146408 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 7996,
            "unit": "ns/op",
            "extra": "146408 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "146408 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "146408 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 59.88,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "20009676 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 59.88,
            "unit": "ns/op",
            "extra": "20009676 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "20009676 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "20009676 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 23634,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "50691 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 23634,
            "unit": "ns/op",
            "extra": "50691 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "50691 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "50691 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 203489,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5034 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 203489,
            "unit": "ns/op",
            "extra": "5034 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5034 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5034 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 1924923,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "547 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 1924923,
            "unit": "ns/op",
            "extra": "547 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "547 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "547 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 26144515,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "76 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 26144515,
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
            "value": 10217,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "126751 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 10217,
            "unit": "ns/op",
            "extra": "126751 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "126751 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "126751 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 6563983,
            "unit": "ns/op\t 4335494 B/op\t   63229 allocs/op",
            "extra": "181 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 6563983,
            "unit": "ns/op",
            "extra": "181 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335494,
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
            "value": 6576,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "181503 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 6576,
            "unit": "ns/op",
            "extra": "181503 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "181503 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "181503 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 833058,
            "unit": "ns/op\t  392922 B/op\t    6173 allocs/op",
            "extra": "1394 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 833058,
            "unit": "ns/op",
            "extra": "1394 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392922,
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
            "value": 9406,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "129776 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 9406,
            "unit": "ns/op",
            "extra": "129776 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "129776 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "129776 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 7437525,
            "unit": "ns/op\t 4722086 B/op\t   69239 allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 7437525,
            "unit": "ns/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722086,
            "unit": "B/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 253.6,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "4819210 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 253.6,
            "unit": "ns/op",
            "extra": "4819210 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "4819210 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4819210 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 1.014,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 1.014,
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
            "value": 654.9,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "1863296 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 654.9,
            "unit": "ns/op",
            "extra": "1863296 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "1863296 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1863296 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 1074,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 1074,
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
          "id": "41fa3e371290498f91d0c01087816c73a3370cd7",
          "message": "chore(dependencies): bump `moby/go-archive` and `x.net` (#8114)\n\n* bump moby/go-archive\n\n* bump golang.org/x/net",
          "timestamp": "2026-10-19T08:21:59Z",
          "tree_id": "5713b038657d569a42b64f479e496be714b973b7",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/41fa3e371290498f91d0c01087816c73a3370cd7"
        },
        "date": 1792400260334,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 464004,
            "unit": "ns/op\t  819428 B/op\t       4 allocs/op",
            "extra": "2193 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 464004,
            "unit": "ns/op",
            "extra": "2193 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819428,
            "unit": "B/op",
            "extra": "2193 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "2193 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 3879,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "307263 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 3879,
            "unit": "ns/op",
            "extra": "307263 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "307263 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "307263 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 28.78,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "41603298 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 28.78,
            "unit": "ns/op",
            "extra": "41603298 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "41603298 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "41603298 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 13645,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "87980 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 13645,
            "unit": "ns/op",
            "extra": "87980 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "87980 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "87980 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 118320,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "8991 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 118320,
            "unit": "ns/op",
            "extra": "8991 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "8991 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "8991 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 1050695,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "999 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 1050695,
            "unit": "ns/op",
            "extra": "999 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "999 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "999 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 21998553,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "99 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 21998553,
            "unit": "ns/op",
            "extra": "99 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
            "unit": "B/op",
            "extra": "99 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "99 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 5126,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "223750 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 5126,
            "unit": "ns/op",
            "extra": "223750 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "223750 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "223750 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 4044496,
            "unit": "ns/op\t 4335493 B/op\t   63229 allocs/op",
            "extra": "313 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 4044496,
            "unit": "ns/op",
            "extra": "313 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335493,
            "unit": "B/op",
            "extra": "313 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "313 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 4238,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "262252 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 4238,
            "unit": "ns/op",
            "extra": "262252 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "262252 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "262252 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 530337,
            "unit": "ns/op\t  392910 B/op\t    6173 allocs/op",
            "extra": "2198 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 530337,
            "unit": "ns/op",
            "extra": "2198 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392910,
            "unit": "B/op",
            "extra": "2198 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "2198 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 6359,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "193557 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 6359,
            "unit": "ns/op",
            "extra": "193557 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "193557 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "193557 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 4500439,
            "unit": "ns/op\t 4722094 B/op\t   69239 allocs/op",
            "extra": "256 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 4500439,
            "unit": "ns/op",
            "extra": "256 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722094,
            "unit": "B/op",
            "extra": "256 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "256 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 137.7,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "8635214 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 137.7,
            "unit": "ns/op",
            "extra": "8635214 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "8635214 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "8635214 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 0.5943,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 0.5943,
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
            "value": 424.7,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "3466208 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 424.7,
            "unit": "ns/op",
            "extra": "3466208 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "3466208 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "3466208 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 519.1,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "2327008 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 519.1,
            "unit": "ns/op",
            "extra": "2327008 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "2327008 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2327008 times\n4 procs"
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
          "id": "319351b4295d37a2c7f0d88fdf4d7bf1b2f2cd25",
          "message": "chore(deps): bump docker/build-push-action from 6.18.0 to 7.3.0 (#8067)\n\nBumps [docker/build-push-action](https://github.com/docker/build-push-action) from 6.18.0 to 7.3.0.\n- [Release notes](https://github.com/docker/build-push-action/releases)\n- [Commits](https://github.com/docker/build-push-action/compare/263435318d21b8e681c14492fe198d362a7d2c83...53b7df96c91f9c12dcc8a07bcb9ccacbed38856a)\n\n---\nupdated-dependencies:\n- dependency-name: docker/build-push-action\n  dependency-version: 7.3.0\n  dependency-type: direct:production\n  update-type: version-update:semver-major\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2026-10-20T08:49:02Z",
          "tree_id": "a42f200fbb88ebdad0094237a7a6a1436eae103c",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/319351b4295d37a2c7f0d88fdf4d7bf1b2f2cd25"
        },
        "date": 1792488318879,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy",
            "value": 514645,
            "unit": "ns/op\t  819428 B/op\t       4 allocs/op",
            "extra": "2054 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - ns/op",
            "value": 514645,
            "unit": "ns/op",
            "extra": "2054 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - B/op",
            "value": 819428,
            "unit": "B/op",
            "extra": "2054 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "2054 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations",
            "value": 5573,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "212966 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - ns/op",
            "value": 5573,
            "unit": "ns/op",
            "extra": "212966 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "212966 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "212966 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert",
            "value": 52.13,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "22976091 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - ns/op",
            "value": 52.13,
            "unit": "ns/op",
            "extra": "22976091 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "22976091 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "22976091 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000",
            "value": 18833,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "75525 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - ns/op",
            "value": 18833,
            "unit": "ns/op",
            "extra": "75525 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "75525 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "75525 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000",
            "value": 147415,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "7965 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - ns/op",
            "value": 147415,
            "unit": "ns/op",
            "extra": "7965 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "7965 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "7965 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000",
            "value": 1383126,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "760 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - ns/op",
            "value": 1383126,
            "unit": "ns/op",
            "extra": "760 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "760 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "760 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000",
            "value": 20987680,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "91 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - ns/op",
            "value": 20987680,
            "unit": "ns/op",
            "extra": "91 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - B/op",
            "value": 24010752,
            "unit": "B/op",
            "extra": "91 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "91 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache",
            "value": 6258,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "190599 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - ns/op",
            "value": 6258,
            "unit": "ns/op",
            "extra": "190599 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "190599 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "190599 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache",
            "value": 5081377,
            "unit": "ns/op\t 4335490 B/op\t   63229 allocs/op",
            "extra": "236 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - ns/op",
            "value": 5081377,
            "unit": "ns/op",
            "extra": "236 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - B/op",
            "value": 4335490,
            "unit": "B/op",
            "extra": "236 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "236 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache",
            "value": 5135,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "230068 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - ns/op",
            "value": 5135,
            "unit": "ns/op",
            "extra": "230068 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "230068 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "230068 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache",
            "value": 645822,
            "unit": "ns/op\t  392902 B/op\t    6173 allocs/op",
            "extra": "1778 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - ns/op",
            "value": 645822,
            "unit": "ns/op",
            "extra": "1778 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - B/op",
            "value": 392902,
            "unit": "B/op",
            "extra": "1778 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1778 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache",
            "value": 7476,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "154148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - ns/op",
            "value": 7476,
            "unit": "ns/op",
            "extra": "154148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "154148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "154148 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache",
            "value": 5761169,
            "unit": "ns/op\t 4722106 B/op\t   69239 allocs/op",
            "extra": "206 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - ns/op",
            "value": 5761169,
            "unit": "ns/op",
            "extra": "206 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - B/op",
            "value": 4722106,
            "unit": "B/op",
            "extra": "206 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "206 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject",
            "value": 182.8,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "6621338 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - ns/op",
            "value": 182.8,
            "unit": "ns/op",
            "extra": "6621338 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "6621338 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "6621338 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol",
            "value": 1.009,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol - ns/op",
            "value": 1.009,
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
            "value": 501.4,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "2401167 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - ns/op",
            "value": 501.4,
            "unit": "ns/op",
            "extra": "2401167 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "2401167 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2401167 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers",
            "value": 694.6,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1737214 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - ns/op",
            "value": 694.6,
            "unit": "ns/op",
            "extra": "1737214 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1737214 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1737214 times\n4 procs"
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
          "id": "c32076b0499a6a4b8fd3c6057bdfb7dfa1444892",
          "message": "chore(deps): bump benchmark-action/github-action-benchmark (#8121)\n\nBumps [benchmark-action/github-action-benchmark](https://github.com/benchmark-action/github-action-benchmark) from 1.20.4 to 1.22.2.\n- [Release notes](https://github.com/benchmark-action/github-action-benchmark/releases)\n- [Changelog](https://github.com/benchmark-action/github-action-benchmark/blob/master/CHANGELOG.md)\n- [Commits](https://github.com/benchmark-action/github-action-benchmark/compare/d48d326b4ca9ba73ca0cd0d59f108f9e02a381c7...4322e5726e6334590d251fc4f92bec0efafc45dc)\n\n---\nupdated-dependencies:\n- dependency-name: benchmark-action/github-action-benchmark\n  dependency-version: 1.22.2\n  dependency-type: direct:production\n  update-type: version-update:semver-minor\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2026-10-21T02:38:30Z",
          "tree_id": "bd135893a7962e8f77459886473d66ebf885b872",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/c32076b0499a6a4b8fd3c6057bdfb7dfa1444892"
        },
        "date": 1792554518476,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate)",
            "value": 344206,
            "unit": "ns/op\t  819428 B/op\t       4 allocs/op",
            "extra": "3046 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate) - ns/op",
            "value": 344206,
            "unit": "ns/op",
            "extra": "3046 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate) - B/op",
            "value": 819428,
            "unit": "B/op",
            "extra": "3046 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "3046 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate)",
            "value": 3872,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "311850 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate) - ns/op",
            "value": 3872,
            "unit": "ns/op",
            "extra": "311850 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate) - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "311850 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate) - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "311850 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig)",
            "value": 28.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "41289410 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig) - ns/op",
            "value": 28.4,
            "unit": "ns/op",
            "extra": "41289410 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "41289410 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "41289410 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 13119,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "87069 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 13119,
            "unit": "ns/op",
            "extra": "87069 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "87069 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "87069 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 109346,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "10966 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 109346,
            "unit": "ns/op",
            "extra": "10966 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "10966 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "10966 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 1095291,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "1087 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 1095291,
            "unit": "ns/op",
            "extra": "1087 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "1087 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "1087 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 24673166,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 24673166,
            "unit": "ns/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 24010752,
            "unit": "B/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "58 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 5156,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "222753 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 5156,
            "unit": "ns/op",
            "extra": "222753 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "222753 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "222753 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 3654485,
            "unit": "ns/op\t 4335494 B/op\t   63229 allocs/op",
            "extra": "326 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 3654485,
            "unit": "ns/op",
            "extra": "326 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 4335494,
            "unit": "B/op",
            "extra": "326 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "326 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 4135,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "291978 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 4135,
            "unit": "ns/op",
            "extra": "291978 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "291978 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "291978 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 535300,
            "unit": "ns/op\t  392908 B/op\t    6173 allocs/op",
            "extra": "2166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 535300,
            "unit": "ns/op",
            "extra": "2166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 392908,
            "unit": "B/op",
            "extra": "2166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "2166 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 6050,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "191233 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 6050,
            "unit": "ns/op",
            "extra": "191233 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "191233 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "191233 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 4469030,
            "unit": "ns/op\t 4722081 B/op\t   69239 allocs/op",
            "extra": "259 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 4469030,
            "unit": "ns/op",
            "extra": "259 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 4722081,
            "unit": "B/op",
            "extra": "259 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "259 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject (github.com/kong/kubernetes-ingress-controller/v3/internal/util)",
            "value": 167.1,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "8596321 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject (github.com/kong/kubernetes-ingress-controller/v3/internal/util) - ns/op",
            "value": 167.1,
            "unit": "ns/op",
            "extra": "8596321 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject (github.com/kong/kubernetes-ingress-controller/v3/internal/util) - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "8596321 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject (github.com/kong/kubernetes-ingress-controller/v3/internal/util) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "8596321 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol (github.com/kong/kubernetes-ingress-controller/v3/internal/util)",
            "value": 0.5883,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol (github.com/kong/kubernetes-ingress-controller/v3/internal/util) - ns/op",
            "value": 0.5883,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol (github.com/kong/kubernetes-ingress-controller/v3/internal/util) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol (github.com/kong/kubernetes-ingress-controller/v3/internal/util) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels)",
            "value": 370.7,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "3262969 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels) - ns/op",
            "value": 370.7,
            "unit": "ns/op",
            "extra": "3262969 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels) - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "3262969 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "3262969 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels)",
            "value": 516.2,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "2287281 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels) - ns/op",
            "value": 516.2,
            "unit": "ns/op",
            "extra": "2287281 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels) - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "2287281 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "2287281 times\n4 procs"
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
          "id": "247f724d02a9b655bf3ffdea45e27ce7a0f4e2ac",
          "message": "bump golang to 1.27.3 (#8124)",
          "timestamp": "2026-10-21T09:36:07Z",
          "tree_id": "b6a4b171346e2e86956fe357ec2a72aea8488664",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/247f724d02a9b655bf3ffdea45e27ce7a0f4e2ac"
        },
        "date": 1792579044802,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate)",
            "value": 909476,
            "unit": "ns/op\t  819431 B/op\t       4 allocs/op",
            "extra": "1168 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate) - ns/op",
            "value": 909476,
            "unit": "ns/op",
            "extra": "1168 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate) - B/op",
            "value": 819431,
            "unit": "B/op",
            "extra": "1168 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1168 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate)",
            "value": 7491,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "137204 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate) - ns/op",
            "value": 7491,
            "unit": "ns/op",
            "extra": "137204 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate) - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "137204 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate) - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "137204 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig)",
            "value": 38.86,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "30784940 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig) - ns/op",
            "value": 38.86,
            "unit": "ns/op",
            "extra": "30784940 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "30784940 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "30784940 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 20956,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "56740 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 20956,
            "unit": "ns/op",
            "extra": "56740 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "56740 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "56740 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 195675,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5496 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 195675,
            "unit": "ns/op",
            "extra": "5496 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5496 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5496 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 2015115,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "591 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 2015115,
            "unit": "ns/op",
            "extra": "591 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "591 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "591 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 31110678,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "57 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 31110678,
            "unit": "ns/op",
            "extra": "57 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 24010752,
            "unit": "B/op",
            "extra": "57 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "57 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 8736,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "140116 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 8736,
            "unit": "ns/op",
            "extra": "140116 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "140116 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "140116 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 6790228,
            "unit": "ns/op\t 4335491 B/op\t   63229 allocs/op",
            "extra": "174 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 6790228,
            "unit": "ns/op",
            "extra": "174 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 4335491,
            "unit": "B/op",
            "extra": "174 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "174 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 7069,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "164182 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 7069,
            "unit": "ns/op",
            "extra": "164182 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "164182 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "164182 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 907411,
            "unit": "ns/op\t  392906 B/op\t    6173 allocs/op",
            "extra": "1275 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 907411,
            "unit": "ns/op",
            "extra": "1275 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 392906,
            "unit": "B/op",
            "extra": "1275 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1275 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 10266,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "116020 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 10266,
            "unit": "ns/op",
            "extra": "116020 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "116020 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "116020 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 7579978,
            "unit": "ns/op\t 4722081 B/op\t   69239 allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 7579978,
            "unit": "ns/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 4722081,
            "unit": "B/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "159 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject (github.com/kong/kubernetes-ingress-controller/v3/internal/util)",
            "value": 229.7,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "5231752 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject (github.com/kong/kubernetes-ingress-controller/v3/internal/util) - ns/op",
            "value": 229.7,
            "unit": "ns/op",
            "extra": "5231752 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject (github.com/kong/kubernetes-ingress-controller/v3/internal/util) - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "5231752 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject (github.com/kong/kubernetes-ingress-controller/v3/internal/util) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5231752 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol (github.com/kong/kubernetes-ingress-controller/v3/internal/util)",
            "value": 1.142,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol (github.com/kong/kubernetes-ingress-controller/v3/internal/util) - ns/op",
            "value": 1.142,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol (github.com/kong/kubernetes-ingress-controller/v3/internal/util) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol (github.com/kong/kubernetes-ingress-controller/v3/internal/util) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels)",
            "value": 604.8,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "1967794 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels) - ns/op",
            "value": 604.8,
            "unit": "ns/op",
            "extra": "1967794 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels) - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "1967794 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1967794 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels)",
            "value": 874.2,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1366016 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels) - ns/op",
            "value": 874.2,
            "unit": "ns/op",
            "extra": "1366016 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels) - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1366016 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1366016 times\n4 procs"
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
          "id": "3106e4319c4cc4094cbab26b0f24470fa9e008f6",
          "message": "chore(deps): bump go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc (#8094)\n\nBumps [go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc](https://github.com/open-telemetry/opentelemetry-go) from 1.33.0 to 1.45.0.\n- [Release notes](https://github.com/open-telemetry/opentelemetry-go/releases)\n- [Changelog](https://github.com/open-telemetry/opentelemetry-go/blob/main/CHANGELOG.md)\n- [Commits](https://github.com/open-telemetry/opentelemetry-go/compare/v1.33.0...v1.45.0)\n\n---\nupdated-dependencies:\n- dependency-name: go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc\n  dependency-version: 1.45.0\n  dependency-type: indirect\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2026-10-22T03:55:42Z",
          "tree_id": "e47e5e2884a5f0d3563fe0acee28ed62d2435be5",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/3106e4319c4cc4094cbab26b0f24470fa9e008f6"
        },
        "date": 1792645574936,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate)",
            "value": 945339,
            "unit": "ns/op\t  819431 B/op\t       4 allocs/op",
            "extra": "1184 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate) - ns/op",
            "value": 945339,
            "unit": "ns/op",
            "extra": "1184 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate) - B/op",
            "value": 819431,
            "unit": "B/op",
            "extra": "1184 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1184 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate)",
            "value": 8471,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "124316 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate) - ns/op",
            "value": 8471,
            "unit": "ns/op",
            "extra": "124316 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate) - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "124316 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate) - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "124316 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig)",
            "value": 59.82,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "20083999 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig) - ns/op",
            "value": 59.82,
            "unit": "ns/op",
            "extra": "20083999 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "20083999 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "20083999 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 24380,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "47460 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 24380,
            "unit": "ns/op",
            "extra": "47460 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "47460 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "47460 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 220238,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "5156 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 220238,
            "unit": "ns/op",
            "extra": "5156 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "5156 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5156 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 2148487,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "480 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 2148487,
            "unit": "ns/op",
            "extra": "480 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "480 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "480 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 26166163,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "75 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 26166163,
            "unit": "ns/op",
            "extra": "75 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 24010752,
            "unit": "B/op",
            "extra": "75 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "75 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 8465,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "142207 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 8465,
            "unit": "ns/op",
            "extra": "142207 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "142207 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "142207 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 6717504,
            "unit": "ns/op\t 4335495 B/op\t   63229 allocs/op",
            "extra": "178 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 6717504,
            "unit": "ns/op",
            "extra": "178 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 4335495,
            "unit": "B/op",
            "extra": "178 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "178 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 6835,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "173073 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 6835,
            "unit": "ns/op",
            "extra": "173073 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "173073 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "173073 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 836421,
            "unit": "ns/op\t  392908 B/op\t    6173 allocs/op",
            "extra": "1382 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 836421,
            "unit": "ns/op",
            "extra": "1382 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 392908,
            "unit": "B/op",
            "extra": "1382 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1382 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 9728,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "123774 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 9728,
            "unit": "ns/op",
            "extra": "123774 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "123774 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "123774 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 7534196,
            "unit": "ns/op\t 4722097 B/op\t   69239 allocs/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 7534196,
            "unit": "ns/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 4722097,
            "unit": "B/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "158 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject (github.com/kong/kubernetes-ingress-controller/v3/internal/util)",
            "value": 250.3,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "4853122 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject (github.com/kong/kubernetes-ingress-controller/v3/internal/util) - ns/op",
            "value": 250.3,
            "unit": "ns/op",
            "extra": "4853122 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject (github.com/kong/kubernetes-ingress-controller/v3/internal/util) - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "4853122 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject (github.com/kong/kubernetes-ingress-controller/v3/internal/util) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4853122 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol (github.com/kong/kubernetes-ingress-controller/v3/internal/util)",
            "value": 0.9764,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol (github.com/kong/kubernetes-ingress-controller/v3/internal/util) - ns/op",
            "value": 0.9764,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol (github.com/kong/kubernetes-ingress-controller/v3/internal/util) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol (github.com/kong/kubernetes-ingress-controller/v3/internal/util) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels)",
            "value": 597.3,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "1998842 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels) - ns/op",
            "value": 597.3,
            "unit": "ns/op",
            "extra": "1998842 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels) - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "1998842 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1998842 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels)",
            "value": 915.2,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1325748 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels) - ns/op",
            "value": 915.2,
            "unit": "ns/op",
            "extra": "1325748 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels) - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1325748 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1325748 times\n4 procs"
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
          "id": "98b3bb682dbc89c25496db1e56cd20d4ad79d3b8",
          "message": "chore(deps): bump actions/download-artifact from 4.3.0 to 8.0.2 (#8123)\n\nBumps [actions/download-artifact](https://github.com/actions/download-artifact) from 4.3.0 to 8.0.2.\n- [Release notes](https://github.com/actions/download-artifact/releases)\n- [Commits](https://github.com/actions/download-artifact/compare/d3f86a106a0bac45b974a628896c90dbdf5c8093...9000827ccba6bdab643e8b6fd33ac0654aef8333)\n\n---\nupdated-dependencies:\n- dependency-name: actions/download-artifact\n  dependency-version: 8.0.2\n  dependency-type: direct:production\n  update-type: version-update:semver-major\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2026-10-22T07:31:44Z",
          "tree_id": "4eb35b75dfa90cb5de3afc63f18fb111e42a2ef9",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/98b3bb682dbc89c25496db1e56cd20d4ad79d3b8"
        },
        "date": 1792656562295,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate)",
            "value": 953255,
            "unit": "ns/op\t  819432 B/op\t       4 allocs/op",
            "extra": "1138 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate) - ns/op",
            "value": 953255,
            "unit": "ns/op",
            "extra": "1138 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate) - B/op",
            "value": 819432,
            "unit": "B/op",
            "extra": "1138 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1138 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate)",
            "value": 8270,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "138028 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate) - ns/op",
            "value": 8270,
            "unit": "ns/op",
            "extra": "138028 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate) - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "138028 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate) - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "138028 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig)",
            "value": 59.65,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "20115279 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig) - ns/op",
            "value": 59.65,
            "unit": "ns/op",
            "extra": "20115279 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "20115279 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "20115279 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 28388,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "49486 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 28388,
            "unit": "ns/op",
            "extra": "49486 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "49486 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "49486 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 262368,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4129 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 262368,
            "unit": "ns/op",
            "extra": "4129 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4129 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4129 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 2062750,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "584 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 2062750,
            "unit": "ns/op",
            "extra": "584 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "584 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "584 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 35585966,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "46 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 35585966,
            "unit": "ns/op",
            "extra": "46 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 24010752,
            "unit": "B/op",
            "extra": "46 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "46 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 8908,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "138643 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 8908,
            "unit": "ns/op",
            "extra": "138643 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "138643 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "138643 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 7294553,
            "unit": "ns/op\t 4335494 B/op\t   63229 allocs/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 7294553,
            "unit": "ns/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 4335494,
            "unit": "B/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "164 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 7328,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "159138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 7328,
            "unit": "ns/op",
            "extra": "159138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "159138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "159138 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 883372,
            "unit": "ns/op\t  392915 B/op\t    6173 allocs/op",
            "extra": "1239 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 883372,
            "unit": "ns/op",
            "extra": "1239 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 392915,
            "unit": "B/op",
            "extra": "1239 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1239 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 10629,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "111625 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 10629,
            "unit": "ns/op",
            "extra": "111625 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "111625 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "111625 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 8453146,
            "unit": "ns/op\t 4722087 B/op\t   69239 allocs/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 8453146,
            "unit": "ns/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 4722087,
            "unit": "B/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "141 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject (github.com/kong/kubernetes-ingress-controller/v3/internal/util)",
            "value": 247.6,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "4891303 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject (github.com/kong/kubernetes-ingress-controller/v3/internal/util) - ns/op",
            "value": 247.6,
            "unit": "ns/op",
            "extra": "4891303 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject (github.com/kong/kubernetes-ingress-controller/v3/internal/util) - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "4891303 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject (github.com/kong/kubernetes-ingress-controller/v3/internal/util) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4891303 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol (github.com/kong/kubernetes-ingress-controller/v3/internal/util)",
            "value": 0.986,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol (github.com/kong/kubernetes-ingress-controller/v3/internal/util) - ns/op",
            "value": 0.986,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol (github.com/kong/kubernetes-ingress-controller/v3/internal/util) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol (github.com/kong/kubernetes-ingress-controller/v3/internal/util) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels)",
            "value": 647.7,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "1827880 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels) - ns/op",
            "value": 647.7,
            "unit": "ns/op",
            "extra": "1827880 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels) - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "1827880 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1827880 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels)",
            "value": 895.3,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1331521 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels) - ns/op",
            "value": 895.3,
            "unit": "ns/op",
            "extra": "1331521 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels) - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1331521 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1331521 times\n4 procs"
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
          "id": "76d2886de5bc8374aa23245e30f50768f27f5017",
          "message": "chore(deps): bump golang from `f241642` to `f20d83a` (#8135)\n\nBumps golang from `f241642` to `f20d83a`.\n\n---\nupdated-dependencies:\n- dependency-name: golang\n  dependency-version: 1.27.3\n  dependency-type: direct:production\n  update-type: version-update:semver-patch\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2026-10-26T15:13:47Z",
          "tree_id": "f236067d17677b1c67a178b81a3bb26ead0856fd",
          "url": "https://github.com/Kong/kubernetes-ingress-controller/commit/76d2886de5bc8374aa23245e30f50768f27f5017"
        },
        "date": 1793029966597,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkSanitizedCopy (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate)",
            "value": 939175,
            "unit": "ns/op\t  819430 B/op\t       4 allocs/op",
            "extra": "1203 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate) - ns/op",
            "value": 939175,
            "unit": "ns/op",
            "extra": "1203 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate) - B/op",
            "value": 819430,
            "unit": "B/op",
            "extra": "1203 times\n4 procs"
          },
          {
            "name": "BenchmarkSanitizedCopy (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "1203 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate)",
            "value": 7938,
            "unit": "ns/op\t    7600 B/op\t      66 allocs/op",
            "extra": "128961 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate) - ns/op",
            "value": 7938,
            "unit": "ns/op",
            "extra": "128961 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate) - B/op",
            "value": 7600,
            "unit": "B/op",
            "extra": "128961 times\n4 procs"
          },
          {
            "name": "BenchmarkGetPluginRelations (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/kongstate) - allocs/op",
            "value": 66,
            "unit": "allocs/op",
            "extra": "128961 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig)",
            "value": 59.67,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "20079961 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig) - ns/op",
            "value": 59.67,
            "unit": "ns/op",
            "extra": "20079961 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "20079961 times\n4 procs"
          },
          {
            "name": "BenchmarkDefaultContentToDBLessConfigConverter_Convert (github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/sendconfig) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "20079961 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 27894,
            "unit": "ns/op\t   24576 B/op\t       2 allocs/op",
            "extra": "48922 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 27894,
            "unit": "ns/op",
            "extra": "48922 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 24576,
            "unit": "B/op",
            "extra": "48922 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "48922 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 255651,
            "unit": "ns/op\t  245760 B/op\t       2 allocs/op",
            "extra": "4285 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 255651,
            "unit": "ns/op",
            "extra": "4285 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 245760,
            "unit": "B/op",
            "extra": "4285 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/10000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4285 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 1903839,
            "unit": "ns/op\t 2408448 B/op\t       2 allocs/op",
            "extra": "579 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 1903839,
            "unit": "ns/op",
            "extra": "579 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 2408448,
            "unit": "B/op",
            "extra": "579 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/100000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "579 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 26047419,
            "unit": "ns/op\t24010752 B/op\t       2 allocs/op",
            "extra": "76 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 26047419,
            "unit": "ns/op",
            "extra": "76 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 24010752,
            "unit": "B/op",
            "extra": "76 times\n4 procs"
          },
          {
            "name": "BenchmarkListHTTPRoutes/1000000 (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "76 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 8294,
            "unit": "ns/op\t    7888 B/op\t     183 allocs/op",
            "extra": "145034 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 8294,
            "unit": "ns/op",
            "extra": "145034 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 7888,
            "unit": "B/op",
            "extra": "145034 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 183,
            "unit": "allocs/op",
            "extra": "145034 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 6735324,
            "unit": "ns/op\t 4335491 B/op\t   63229 allocs/op",
            "extra": "177 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 6735324,
            "unit": "ns/op",
            "extra": "177 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 4335491,
            "unit": "B/op",
            "extra": "177 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_Without_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 63229,
            "unit": "allocs/op",
            "extra": "177 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 6694,
            "unit": "ns/op\t    6688 B/op\t     170 allocs/op",
            "extra": "175197 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 6694,
            "unit": "ns/op",
            "extra": "175197 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 6688,
            "unit": "B/op",
            "extra": "175197 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 170,
            "unit": "allocs/op",
            "extra": "175197 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 825565,
            "unit": "ns/op\t  392917 B/op\t    6173 allocs/op",
            "extra": "1419 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 825565,
            "unit": "ns/op",
            "extra": "1419 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 392917,
            "unit": "B/op",
            "extra": "1419 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 6173,
            "unit": "allocs/op",
            "extra": "1419 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 9696,
            "unit": "ns/op\t    8288 B/op\t     190 allocs/op",
            "extra": "122089 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 9696,
            "unit": "ns/op",
            "extra": "122089 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 8288,
            "unit": "B/op",
            "extra": "122089 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Small_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 190,
            "unit": "allocs/op",
            "extra": "122089 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store)",
            "value": 7754391,
            "unit": "ns/op\t 4722084 B/op\t   69239 allocs/op",
            "extra": "154 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - ns/op",
            "value": 7754391,
            "unit": "ns/op",
            "extra": "154 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - B/op",
            "value": 4722084,
            "unit": "B/op",
            "extra": "154 times\n4 procs"
          },
          {
            "name": "BenchmarkCacheStores_TakeSnapshot/Big_With_Missed_Cache (github.com/kong/kubernetes-ingress-controller/v3/internal/store) - allocs/op",
            "value": 69239,
            "unit": "allocs/op",
            "extra": "154 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject (github.com/kong/kubernetes-ingress-controller/v3/internal/util)",
            "value": 251.1,
            "unit": "ns/op\t     336 B/op\t       2 allocs/op",
            "extra": "4851159 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject (github.com/kong/kubernetes-ingress-controller/v3/internal/util) - ns/op",
            "value": 251.1,
            "unit": "ns/op",
            "extra": "4851159 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject (github.com/kong/kubernetes-ingress-controller/v3/internal/util) - B/op",
            "value": 336,
            "unit": "B/op",
            "extra": "4851159 times\n4 procs"
          },
          {
            "name": "BenchmarkFromK8sObject (github.com/kong/kubernetes-ingress-controller/v3/internal/util) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4851159 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol (github.com/kong/kubernetes-ingress-controller/v3/internal/util)",
            "value": 0.9779,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol (github.com/kong/kubernetes-ingress-controller/v3/internal/util) - ns/op",
            "value": 0.9779,
            "unit": "ns/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol (github.com/kong/kubernetes-ingress-controller/v3/internal/util) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkValidateProtocol (github.com/kong/kubernetes-ingress-controller/v3/internal/util) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "1000000000 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels)",
            "value": 612.9,
            "unit": "ns/op\t    1152 B/op\t       1 allocs/op",
            "extra": "1926400 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels) - ns/op",
            "value": 612.9,
            "unit": "ns/op",
            "extra": "1926400 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels) - B/op",
            "value": 1152,
            "unit": "B/op",
            "extra": "1926400 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumer_groups (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1926400 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels)",
            "value": 930.4,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "1286734 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels) - ns/op",
            "value": 930.4,
            "unit": "ns/op",
            "extra": "1286734 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels) - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "1286734 times\n4 procs"
          },
          {
            "name": "BenchmarkGetCombinations/consumers (github.com/kong/kubernetes-ingress-controller/v3/internal/util/rels) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1286734 times\n4 procs"
          }
        ]
      }
    ]
  }
}