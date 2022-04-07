# Introduction

The goal of this project is to ease running OHDSI projects in Azure.

Some of the OHDSI projects include:

* [Common Data Model (CDM)](https://github.com/OHDSI/CommonDataModel), including [Vocabulary](https://github.com/OHDSI/Vocabulary-v5.0)
* [Atlas](https://github.com/OHDSI/Atlas)
* [WebApi](https://github.com/OHDSI/WebAPI)
* [Achilles](https://github.com/OHDSI/Achilles)
* [ETL-Synthea](https://github.com/OHDSI/ETL-Synthea)

## Overview

![Overview](/docs/media/azure_overview.png)

You can use [Azure DevOps pipelines](/pipelines/README.md/#pipelines) to manage your environment.  To create your environment, please review the [guide](/docs/creating_your_environment.md) for an overview.

Your administrator can follow the [administrative steps](/infra/README.md/#administrative-steps) to manage your [bootstrap resource group](/infra/terraform/bootstrap/README.md#setup-azure-bootstrap-resource-group), your [Azure AD Groups](/infra/terraform/bootstrap/README.md#setup-azure-ad-group), and your [Azure DevOps](/infra/terraform/bootstrap/README.md#setup-azure-devops) project.

You can also use [terraform to manage](/infra/README.md/#running-terraform) your Azure resources in the [OMOP resource group](/infra/terraform/omop/README.md).

You can host your [CDM in Azure SQL](/sql/README.md#cdm-notes).  You can [load your vocabularies](/docs/setup/setup_vocabulary.md) into Azure Storage so that the [Azure DevOps Vocabulary Release Pipeline](/pipelines/README.md/#vocabulary-release-pipeline) can populate your [Azure SQL CDM](/sql/README.md/#vocabulary-notes).

You can [setup Atlas and Webapi](/docs/setup/setup_atlas_webapi.md) using the [Broadsea Build Pipeline](/pipelines/README.md/#broadsea-build-pipeline) to build and push the [Broadsea webtools (for Atlas / WebApi)](/apps/broadsea-webtools/README.md) image into Azure Container Registry. You can then run the [Broadsea Release Pipeline](/pipelines/README.md/#broadsea-release-pipeline) to configure Atlas and WebApi in your Azure App Service.

You can also [setup Achilles and Synthea](/docs/setup/setup_achilles_synthea.md) using the [Broadsea Build Pipeline](/pipelines/README.md/#broadsea-build-pipeline) to build and push the [Broadsea Methods (for Achilles and Synthea)](/apps/broadsea-methods/README.md) image into Azure Container Registry.  You can then run the [Broadsea Release Pipeline](/pipelines/README.md/#broadsea-release-pipeline) to perform the following steps:

1. Run an [ETL job](/apps/broadsea-methods/README.md/#synthea-etl) and use [Synthea to generate synthetic patient data](/apps/broadsea-methods/README.md/#use-synthea-to-generate-synthetic-patient-data) as an optional step
2. Run [Achilles](/apps/broadsea-methods/README.md/#achilles) to characterize the CDM data in Azure SQL

## CDM Version

This setup supports a modified version of the CDM [v5.3.1](/sql/cdm/v5.3.1/) schema based on the [CDM v5.3.1 for SQL Server](https://github.com/OHDSI/CommonDataModel/tree/v5.3.1/Sql%20Server).

You can review more notes on the modifications in the [readme](/sql/README.md/#modifications-from-ohdsi).

# Getting Started

To get started, first clone the repository.

```console
git clone https://github.com/microsoft/OHDSIonAzure
```

You can work through the notes on [creating your environment](/docs/creating_your_environment.md) which will walk through how to set up OHDSI on Azure.

### TODO Review usage here

TODO: Guide users through getting your code up and running on their own system. In this section you can talk about:
1.	Installation process
2.	Software dependencies
3.	Latest releases
4.	API references

# Build and Test
TODO: Describe and show how to build your code and run the tests.

# Contribute
TODO: Explain how other users and developers can contribute to make your code better.

If you want to learn more about creating good readme files then refer the following [guidelines](https://docs.microsoft.com/en-us/azure/devops/repos/git/create-a-readme?view=azure-devops). You can also seek inspiration from the below readme files:
- [ASP.NET Core](https://github.com/aspnet/Home)
- [Visual Studio Code](https://github.com/Microsoft/vscode)
- [Chakra Core](https://github.com/Microsoft/ChakraCore)
