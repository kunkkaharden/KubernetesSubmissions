# Google Cloud SQL vs. a Database on a Kubernetes PVC

This comparison covers two ways to run a relational database for an application:

1. **Database as a Service (DBaaS)**, such as Google Cloud SQL, where a cloud
   provider operates the database service.
2. **Self-managed database on Kubernetes storage**, where the database runs in a
   Pod and stores its data on a PersistentVolumeClaim (PVC).

A PVC is storage, not a database service. With a PVC, the team is still
responsible for installing, configuring, securing, upgrading, monitoring, and
recovering the database software.

## At a glance

| Area | Cloud SQL (DBaaS) | Database on a PVC |
|---|---|---|
| Initial setup | Provision an instance, configure network access and identity/credentials, then connect the application. Usually less database infrastructure to assemble, but cloud networking and permissions can take work. | Configure a database workload, persistent storage, credentials, Services, probes, resource requests, and often policies and disruption handling. More Kubernetes and database setup is the team's responsibility. |
| Initial cost | Usually no database server to purchase or install, but provisioning may take engineering time. | Can use existing cluster capacity and storage, but setup and operational engineering time are higher. |
| Ongoing cost | Recurring charges for the database instance and storage; backups, high availability, replicas, network traffic, and other options may add charges. | Persistent storage and cluster compute are billed or consume existing capacity. Staff time for operations is a significant part of the total cost. |
| Maintenance | Provider manages much of the underlying infrastructure and routine service operations. The team still manages database settings, users, schema, application compatibility, capacity choices, and its backup/recovery policy. | Team manages database and Kubernetes lifecycle, including patching, upgrades, capacity, monitoring, availability, and recovery procedures. |
| Availability and recovery | Managed availability and failover features can be enabled, depending on configuration and price. They reduce operational work but do not eliminate outages or the need to test recovery. | Depends on the cluster, storage system, database configuration, and operational procedures. A PVC preserves data across some Pod restarts; it does not by itself provide database high availability or disaster recovery. |
| Portability | Convenient within the provider, but provider-specific networking, identity, and features can make migration less direct. Database engines and supported features may also be constrained by the service. | Kubernetes manifests and standard database tooling can be portable, but storage classes, snapshots, and cluster behavior vary between environments. |

## Initialization and cost

### Cloud SQL

**Work to get started**

- Select the database engine, version, region, machine size, and storage.
- Configure private or public connectivity, firewall/network rules, and access
  controls such as service identities or database users.
- Store credentials securely and configure the application to connect, including
  TLS or a provider connector where appropriate.
- Decide on availability, maintenance windows, backup retention, and recovery
  options before production use.

Cloud consoles and infrastructure-as-code can make provisioning straightforward.
Private networking, IAM, secrets, and application connectivity are often the
parts that require the most initial coordination.

**Cost**

There is usually a recurring charge for the provisioned database capacity and
storage, even when utilization is low. High availability, read replicas, backup
retention, extra storage, and data transfer can increase the bill. Exact prices
depend on region, engine, configuration, and current provider pricing; estimate
them with the provider's pricing calculator for the expected workload.

### Database on a PVC

**Work to get started**

- Choose and configure a database image, version, credentials, and database
  settings.
- Create the Kubernetes workload and Service, and request a PVC through an
  appropriate StorageClass.
- Configure health checks, resource limits, access controls, updates, and
  scheduling behavior.
- Decide how to handle failures, upgrades, backups, restores, and storage
  expansion.

The application can connect through an in-cluster Service, and the setup can
reuse existing Kubernetes infrastructure. However, the cluster providing
storage does not configure or operate the database for the team.

**Cost**

The PVC consumes billable or finite storage, and the database consumes cluster
compute and memory. This can have a low incremental cash cost when capacity is
already available, but that does not make it operationally free: engineering
time, monitoring, backup storage, and recovery testing should be included in the
total cost.

## Maintenance and operational responsibility

### Cloud SQL

The provider operates the managed service infrastructure and offers
provider-managed maintenance and availability features. The user generally has
less work managing database hosts, storage systems, and service failover.

The application team remains responsible for choosing supported versions and
settings, testing upgrades against the application, managing schema changes and
access, monitoring performance and cost, and understanding maintenance or
failover impact. A managed service does not automatically make application
connections resilient or guarantee that a workload meets its recovery goals.

### Database on a PVC

The team owns both the database and its Kubernetes deployment. Typical
responsibilities include security updates and database upgrades, monitoring and
alerting, storage and compute sizing, replication or failover, backup
verification, and recovery after node, cluster, storage, or data failures.
Kubernetes can restart a failed Pod, but restarting a database is not the same
as restoring lost or corrupted data.

This approach provides more control over database configuration and avoids
relying on a particular managed database service. In return, reliable
operations require people with database and Kubernetes expertise and tested
procedures.

## Backups and restore

Backups should be planned separately from ordinary persistence. A PVC helps
data survive some Pod restarts or rescheduling events, but it is not a backup:
storage failure, accidental deletion, corruption, or a cluster/region incident
can affect both the database and its volume. For either option, define retention,
recovery-point and recovery-time objectives, access controls, and regularly test
restores.

| Method | Where it fits | Ease and limitations |
|---|---|---|
| Cloud SQL automated backups and point-in-time recovery | Cloud SQL | Can be enabled and managed through the provider's console, API, or infrastructure-as-code. Retention and point-in-time recovery depend on configuration, engine, and service options; storage and retention may affect cost. Restore is integrated, but usually creates or recovers an instance that the application must be pointed to and validated against. |
| Cloud SQL on-demand/manual backup | Cloud SQL | Useful before a risky change or for a specific recovery point. It is typically easy to initiate and manage through provider tooling, but still needs a retention and deletion policy and a tested restore procedure. |
| Database-native logical backup (for example, `pg_dump` or `mysqldump`) | Either option | Produces a portable database-level copy and can be sent to object storage. Requires a scheduled job, secure handling, retention, monitoring for failures, and restore testing. Large databases can take longer to restore than physical backups. |
| Volume snapshot | PVC-based database | May be convenient when the CSI driver and storage provider support snapshots. A snapshot of a live database may not be application-consistent unless the database is quiesced or the backup process coordinates with it. Snapshots on the same storage system may not protect against that system or region failing. |
| Database-native physical backup or replica-based backup | PVC-based database | Can support larger databases or shorter restore times, but setup and operations are more involved. Replication alone is not a backup: unwanted changes or corruption can also replicate. Keep recoverable copies outside the database's failure domain. |

For Cloud SQL, integrated backup and restore features are usually the easiest
starting point, with logical exports useful for portability or independent
copies. For a PVC-hosted database, a practical baseline is a scheduled,
database-native backup copied to independent object storage, plus tested
restore instructions. Snapshots can complement this, but should not be assumed
to replace database-aware backups.

## Which option to choose

Choose **Cloud SQL** when reducing database operations, getting integrated
backup/recovery options, and using managed availability features are worth the
ongoing service cost and provider-specific integration. It is often a better fit
for production workloads without a team prepared to operate databases.

Choose a **database on a PVC** when direct control, portability, learning, or
using existing cluster capacity are priorities, and the team can own
maintenance, availability, and tested backup/recovery. It can be economical for
small or non-critical workloads, but a single database Pod and PVC should not
be treated as a highly available or disaster-recovery design.

The decision should compare the total cost of ownership, not only the cloud
invoice: managed-service charges versus cluster resources, backup storage,
engineering time, expected downtime, and the cost of data loss.
