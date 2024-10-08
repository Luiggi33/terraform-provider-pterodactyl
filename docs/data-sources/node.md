---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pterodactyl_node Data Source - pterodactyl"
subcategory: ""
description: |-
  The Pterodactyl node data source allows Terraform to read a nodes data from the Pterodactyl Panel API.
---

# pterodactyl_node (Data Source)

The Pterodactyl node data source allows Terraform to read a nodes data from the Pterodactyl Panel API.



<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `id` (Number) The ID of the node.
- `name` (String) The name of the node.
- `uuid` (String) The UUID of the node.

### Read-Only

- `behind_proxy` (Boolean) The behind proxy status of the node.
- `created_at` (String) The creation date of the node.
- `daemon_listen` (Number) The daemon listen of the node.
- `daemon_sftp` (Number) The daemon SFTP of the node.
- `description` (String) The description of the node.
- `disk` (Number) The disk of the node.
- `disk_overallocate` (Number) The disk overallocate of the node.
- `fqdn` (String) The FQDN of the node.
- `location_id` (Number) The location ID of the node.
- `maintenance_mode` (Boolean) The maintenance mode status of the node.
- `memory` (Number) The memory of the node.
- `memory_overallocate` (Number) The memory overallocate of the node.
- `public` (Boolean) The public status of the node.
- `scheme` (String) The scheme of the node.
- `updated_at` (String) The last update date of the node.
- `upload_size` (Number) The upload size of the node.
