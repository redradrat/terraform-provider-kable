local lib = import 'lib/main.libsonnet';

// Final JSON Output
{
    test: {
       apiVersion: "v1",
       kind: "Test",
       metadata: { name: std.extVar("instanceName"), namespace: std.extVar("nameSelection") },
    },
}
