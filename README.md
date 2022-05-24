<a href="https://www.lamassu.io/">
    <img src="logo.png" alt="Lamassu logo" title="Lamassu" align="right" height="80" />
</a>

Lamassu Virtual DMS
=======
[![License: MPL 2.0](https://img.shields.io/badge/License-MPL%202.0-blue.svg)](http://www.mozilla.org/MPL/2.0/index.txt)

[Lamassu](https://www.lamassu.io) project is a Public Key Infrastructure (PKI) for the Internet of Things (IoT).

To launch Lamassu-Virtual-DMS follow the next steps:

1. Clone the repository and get into the directory: `https://github.com/lamassuiot/lamassu-virtual-dms.git && cd lamassu-virtual-dms`.
2. Change the configuration variables of the `config.json` file.

```
{
    "dms": {
        "device_store": "<DEVICES_STORE>", // Folder where device certificates are stored
        "dms_store": "<DMS_STORE>", // Folder where DMS certificates are stored
        "endpoint":"<DMS_SERVER>", // DMS server endpoint
        "dms_name":"<DEFAULT_DMS>", // DMS Name
        "common_name":"<DEFAULT_DMS>", // Common_name to create the CSR
        "country":"<COUNTRY>", // Country to create the CSR
        "locality":"<LOCALITY>", // Locality to create the CSR
        "organization":"<ORGANIZATION>", // Organization to create the CSR
        "organization_unit":"<ORGANIZATION_UNIT>", // Organization_unit to create the CSR
        "state":"<STATE>" // State to create the CSR


    },
    "devmanager":{
        "devcrt": "<DEV_CERTIFICATE>", // Public certificate to connect to the device-manager
        "addr": "<DEVMANAGER_SERVER>" //Device Manager Server Endpoint
    },
    "auth":{
        "endpoint":"<AUTH_SERVER>", // Authentication Server endpoint
        "username":"<PASSWORD>", // User name to connect to the authentication server
        "password":"<PASSWORD>" // Password to connect to the authentication server

    }
}

```
*Common_name and Dms_name have to have the same value

4. Run the Lamassu-Default-DMS UI:
    ```
    go run cmd/main.go
    ```
## Lamassu Virtual DMS operating modes
 
In the Lamassu-Virtual-DMS we have two pages:

1. Create the DMS, once the DMS is created, the Auto_Enroll of the devices is done.

<img src="CreateDMS.png" alt="Create DMS UI" title="Create DMS" />

2. Make the Auto_Enroll of the devices indicating the ID of a DMS.

<img src="AutoEnroll.png" alt="Auto Enroll Devices UI" title="Auto Enroll" />