# Mothership Helm Chart

This chart is used to deploy the Mothership application to a Kubernetes cluster.   The chart assumes that the following 
services are already deployed in the cluster:

* CertManager (if ingress is enabled)
* Temporal

## Configuration

The following table lists the configurable parameters of the Mothership chart and their default values.

| Parameter                   | Description                                                                                                  | Default                                                                    |
|-----------------------------|--------------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------|
| `database.uri`              | uri is the connection string for the database.                                                               | `"postgres://postgres:password@localhost:5432/postgres"`                   |
| `ui.enabled`                | enables/disables the UI deployment.                                                                          | `true`                                                                     |
| `admin.enabled`             | enables/disables the admin API deployment.                                                                   | `true`                                                                     |
| `oidc.required_group`       | required_group is the group required for access to the application.                                          | `"releng"`                                                                 |
| `oidc.issuer`               | issuer is the URL of the OpenID Connect provider.                                                            | `"https://id.openela.org/realms/openela"`                                  |
| `oidc.client_id`            | client_id is the client id used to authenticate with the OpenID Connect provider.                            | `"mothership"`                                                             |
| `oidc.client_secret`        | client_secret is the client secret used to authenticate with the OpenID Connect provider.                    | `""`                                                                       |
| `github.public`             | public determines whether repositories are public or private.                                                | `"false"`                                                                  |
| `github.organization`       | organization is the GitHub organization used for the application.                                            | `"openela-main"`                                                           |
| `github.app_id`             | app_id is the GitHub App ID used for the application.                                                        | `"416803"`                                                                 |
| `github.private_key`        | private_key is the private key used to authenticate with the GitHub App.                                     | `""`                                                                       |
| `storage.endpoint`          | endpoint is the endpoint for the storage service.                                                            | `"https://ax8edlmsvvfp.compat.objectstorage.us-phoenix-1.oraclecloud.com"` |
| `storage.access_key`        | access_key is the access key used to authenticate with the storage service                                   | `""`                                                                       |
| `storage.secret_key`        | secret_key is the secret key used to authenticate with the storage service                                   | `""`                                                                       |
| `storage.region`            | region is the region used for the storage service.                                                           | `"us-phoenix-1"`                                                           |
| `storage.connection_string` | connection_string is the connection string used to access the storage service.                               | `"s3://mship-srpm1"`                                                       |
| `storage.path_style`        | path_style determines whether the storage service uses path style access which includes the bucket name.     | `"true"`                                                                   |
| `bugtracker.provider`       | provider is the provider used for the bug tracker.                                                           | `"github"`                                                                 |
| `bugtracker.repository`     | repository is the repository used for the bug tracker.                                                       | `"openela/issues"`                                                         |
| `bugtracker.use_forge_auth` | use_forge_auth determines whether the bug tracker uses Forge authentication.                                 | `"true"`                                                                   |
| `temporal.address`          | address is the address used to connect to the Temporal service.                                              | `"temporal-frontend.default.svc.cluster.local:7233"`                       |
| `temporal.namespace`        | namespace is the namespace used for the Temporal service.                                                    | `"mothership"`                                                             |
| `temporal.task_queue`       | task_queue is the task queue used for the Temporal service.                                                  | `"worker"`                                                                 |
| `image.repository`          | repository is a base url and has expectations for the image name to be set off that per application.         | `"ghcr.io/mstg"`                                                           |
| `image.pullPolicy`          | pullPolicy is set to always since default configuration is "latest".                                         | `"Always"`                                                                 |
| `image.tag`                 | tag assumes all services are using the same tag.                                                             | `"latest"`                                                                 |
| `ingress.enabled`           | enabled determines whether the ingress is enabled and if it is then it requires CertManager to be installed. | `false`                                                                    |
| `ingress.host`              | host is the host used for the ingress.                                                                       | `"imports.openela.org"`                                                    |
| `resources.requests.cpu`    | cpu is set for all containers.                                                                               | `"300m"`                                                                   |
| `resources.requests.memory` | memory is set for all containers.                                                                            | `"128Mi"`                                                                  |

---
_Documentation generated by [Frigate](https://frigate.readthedocs.io)._

## Installation

From the chart directory, run the following command:

```bash
helm upgrade --install --create-namespace --namespace mothership mothership .
```