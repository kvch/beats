#!/bin/sh

# setup role mapping for Kerberos user
curl -u elastic:changeme -H "Content-Type: application/json" -XPOST http://elasticsearch_kerberos.elastic:9200/_xpack/security/role_mapping/kerbrolemapping -d @- <<EOF
{
    "roles" : [ "superuser" ],
    "enabled": true,
    "rules" : {
    "field" : { "username" : "beats@ELASTIC" }
    }
}
EOF

# start filebeat in debug mode
./filebeat -e -d "*"
