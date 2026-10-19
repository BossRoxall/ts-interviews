# -*- mode: Python -*-
# Local deployment of the rate limiter using the production overlay.
#
#   tilt up       build the selected service, deploy infra/deployments/production, port-forward
#   tilt down     remove everything Tilt created
#
# Options come from tilt_config.json, or override them on the command line:
#   tilt up -- --selected_language python --host_port 9090

config.define_string('selected_language', usage='Service to build and deploy: go, python, or typescript (node, nodejs)')
config.define_string('host_port', usage='Local port forwarded to the rate-limiter pod')
config.define_bool('loadgen', usage='Add the manual loadgen resource')
cfg = config.parse()

LANGUAGES = {
    'go': 'go',
    'python': 'python',
    'typescript': 'typescript',
    'node': 'typescript',
    'nodejs': 'typescript',
}
SELECTED = cfg.get('selected_language', 'go').strip().lower()
if SELECTED not in LANGUAGES:
    fail('selected_language %r is not one of: %s' % (SELECTED, ', '.join(sorted(LANGUAGES.keys()))))
SERVICE_DIR = 'services/rate-limiter-%s' % LANGUAGES[SELECTED]

HOST_PORT = cfg.get('host_port', '8080')
ENABLE_LOADGEN = cfg.get('loadgen', True)

print('Deploying %s from %s' % (LANGUAGES[SELECTED], SERVICE_DIR))

# Local clusters only. Developer kubeconfigs often hold production contexts too,
# so refuse anything that is not a local engine.
LOCAL_CONTEXTS = ['orbstack', 'docker-desktop', 'minikube', 'rancher-desktop', 'colima']
context = k8s_context()
if context not in LOCAL_CONTEXTS and not context.startswith('kind-'):
    fail('Refusing to deploy to kube context %r. Switch to a local cluster (%s, kind-*).'
         % (context, ', '.join(LOCAL_CONTEXTS)))
allow_k8s_contexts(context)

# Same image name as production, built from the selected service. Tilt replaces
# the :latest reference with its own build when it applies the manifests.
IMAGE = 'registry.gitlab.com/turionspace/interviews/rate-limiter'
docker_build(IMAGE, SERVICE_DIR)

# Deploy exactly what Argo CD deploys to production.
watch_file('infra/charts')
watch_file('infra/deployments/production')
k8s_yaml(local('kustomize build --enable-helm infra/deployments/production', quiet=True))

k8s_resource(
    'rate-limiter',
    port_forwards=['%s:8080' % HOST_PORT],
    labels=['app'],
)

if ENABLE_LOADGEN:
    local_resource(
        'loadgen',
        cmd='cd tools/loadgen && go run . -url http://localhost:%s/v1/check' % HOST_PORT,
        auto_init=False,
        trigger_mode=TRIGGER_MODE_MANUAL,
        resource_deps=['rate-limiter'],
        labels=['tools'],
    )
