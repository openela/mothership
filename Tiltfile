# Load necessary extensions
load('ext://cert_manager', 'deploy_cert_manager')
load('ext://configmap', 'configmap_create')
load('ext://helm_remote', 'helm_remote')
load('ext://secret', 'secret_create_generic', 'secret_from_dict')
load('ext://ko', 'ko_build')

# Database credentials
DB_PASSWORD = '43254SSfsdAA32ds3232sds'
NAMESPACE = 'mothership'

# Deploy PostgreSQL
helm_remote(
    'postgresql',
    repo_url='https://charts.bitnami.com/bitnami',
    set=[
        "auth.database=mothership",
        "auth.postgresPassword="+DB_PASSWORD
    ],
    namespace=NAMESPACE,
)
k8s_resource('postgresql', port_forwards=[5432])

# Deploy Valkey (Redis alternative)
helm_remote(
    'valkey',
    repo_url='oci://registry-1.docker.io/bitnamicharts',
    set=["auth.password="+DB_PASSWORD],
    namespace=NAMESPACE,
)
k8s_resource('valkey-primary', port_forwards=[6379])

# Deploy Dex
k8s_yaml('k8s/dev/dex.yaml')
configmap_create('dex-config', from_file=['dex.yaml=./dev/dex.yaml'])
k8s_resource('dev-dex', port_forwards=[5556])

# Deploy Temporal
k8s_yaml('k8s/dev/temporal.yaml')
k8s_resource('dev-temporal', port_forwards=[7233, 8233])

# Deploy Cert Manager
deploy_cert_manager()

# Secrets Management
k8s_yaml(secret_from_dict(
    "db", namespace=NAMESPACE, 
    inputs={'uri': "postgres://postgres:"+DB_PASSWORD+"@postgresql/mothership?sslmode=disable"}
))

k8s_yaml(secret_from_dict("gh", namespace=NAMESPACE, inputs={'client_id': '', 'client_secret': ''}))
k8s_yaml(secret_from_dict("csrf", namespace=NAMESPACE, inputs={'secret': DB_PASSWORD * 2}))
k8s_yaml(secret_from_dict("redis", namespace=NAMESPACE, inputs={'password': DB_PASSWORD}))

# Load all Kubernetes manifests
k8s_yaml([
    "k8s/admin-api.yaml",
    "k8s/api.yaml",
    "k8s/namespace.yaml",
    "k8s/serviceaccount.yaml",
    "k8s/ui.yaml",
    "k8s/worker.yaml",
])

# Update DNS names in cert configuration
objects = read_yaml_stream('k8s/cert.yaml')
for obj in objects:
    obj['spec']['dnsNames'] = ['mothership.local']
k8s_yaml(encode_yaml_stream(objects))

# Build Docker images using `ko`
ko_build("github.com/openela/mothership/cmd/mship_server", "./cmd/mship_server")
ko_build("github.com/openela/mothership/cmd/mship_admin_server", "./cmd/mship_admin_server")
ko_build("github.com/openela/mothership/cmd/mship_worker_server", "./cmd/mship_worker_server")

# Build UI image (live rebuild enabled)
docker_build('ghcr.io/openela/mothership-ui', './ui',
    build_args={'node_env': 'development'},
    # entrypoint='',
    live_update=[
        sync('./ui', '/app'),
        run('cd /app && pnpm install --frozen-lockfile && pnpm run build', trigger=['./ui/package.json', './ui/pnpm-lock.lock']),

        run('touch /app/server.ts', trigger='./server.ts'),
])

# Assign local port-forwards for easier access
k8s_resource("mothership-api-deployment", port_forwards=6677, resource_deps=[
    "postgresql", "valkey-primary", "dev-dex"
])
k8s_resource("mothership-admin-api-deployment", port_forwards=6687, resource_deps=[
    "postgresql", "valkey-primary", "dev-dex"
])
k8s_resource("mothership-worker-deployment", port_forwards=9114, resource_deps=[
    "postgresql", "valkey-primary", "dev-temporal"
])
k8s_resource("mothership-ui-deployment", port_forwards=9111, resource_deps=[
    "mothership-api-deployment", "mothership-admin-api-deployment"
])
