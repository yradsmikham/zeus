# Zeus - A CLI for OHDSI on Azure

Welcome to Zeus!

The objective of this project is to significantly reduce the time and cost for deploying an instance of OHDSI CDM on Azure. Zeus is a command line interface (CLI) that allows for easy deployment and management of OHDSI CDM resources in Azure.

## Prerequisites

1. Terraform resources deployed in Azure
2. Login to Azure:

```
az login
az account set --subscription <subscription-id>
```

### Common Commands:

- Check for prerequisites (i.e. TF, az, azure devops ext, git), import pipelines to Azure DevOps (if they don't already exist)

```
zeus init --env dev --org https://dev.azure.com/<organization> --proj <project>
```

![zeus init](./docs/zeus_init.gif)

- Importing vocabulary files to Storage Account and inserting them into CDM:

```
zeus vocab upload --path /path/to/vocab/files --storage-account name-of-storage-account
zeus vocab import --env dev
```

- Deploy OHDSI applications to environments

```
zeus deploy broadsea-webtools --env dev
zeus deploy broadsea-methods --env dev
zeus deploy achilles --env dev
```

![zeus deploy](./docs/zeus_deploy.gif)
