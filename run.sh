#!/bin/bash

docker build -t lamassu-virtual-dms . 

docker run -it \
    -e VDMS_DOMAIN=$VDMS_DOMAIN \
    -e VDMS_USERNAME=$VDMS_USERNAME \
    -e VDMS_PASSWORD=$VDMS_PASSWORD \
    -e VDMS_COUNTRY=$VDMS_COUNTRY \
    -e VDMS_STATE=$VDMS_STATE \
    -e VDMS_LOCALITY=$VDMS_LOCALITY \
    -e VDMS_ORGANIZATION=$VDMS_ORGANIZATION \
    -e VDMS_ORGANIZATION_UNIT=$VDMS_ORGANIZATION_UNIT \
    -v ${PWD}/device-certificates:/app/device-certificates \
    -v ${PWD}/dms-certificates:/app/dms-certificates \
    -v $VDMS_EST_SERVER_CERT:/app/downstream/tls.crt \
    lamassu-virtual-dms 