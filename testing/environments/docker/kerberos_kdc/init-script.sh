#!/bin/bash

add_principal_to_elastic_realm() {
    username=$1
    password=$2

    echo "Adding $username principal"
    kadmin.local -q "delete_principal -force $username@$REALM"
    echo ""
    kadmin.local -q "addprinc -e AES256-CTS-HMAC-SHA1-96:normal -pw $password $username@$REALM"
    echo ""
}

echo "==================================================================================="
echo "==== Kerberos KDC and Kadmin ======================================================"
echo "==================================================================================="

echo "REALM: $REALM"
echo "KADMIN_PRINCIPALS: $PRINCIPALS"
echo "KADMIN_PASSWORD: $PASSWORD"
echo ""

echo "==================================================================================="
echo "==== /etc/krb5.conf ==============================================================="
echo "==================================================================================="
KDC_KADMIN_SERVER=$(hostname -f)
tee /etc/krb5.conf <<EOF
[libdefaults]
	default_realm = $REALM
    allow_weak_crypto = true
    ignore_acceptor_hostname = true

[realms]
	$REALM = {
		kdc_ports = 88,750
		kadmind_port = 749
		kdc = $KDC_KADMIN_SERVER
		admin_server = $KDC_KADMIN_SERVER
	}
EOF
echo ""

echo "==================================================================================="
echo "==== /etc/krb5kdc/kdc.conf ========================================================"
echo "==================================================================================="
tee /etc/krb5kdc/kdc.conf <<EOF
[realms]
	$REALM = {
		acl_file = /etc/krb5kdc/kadm5.acl
		max_renewable_life = 7d 0h 0m 0s
		default_principal_flags = +preauth
	}
EOF
echo ""

echo "==================================================================================="
echo "==== /etc/krb5kdc/kadm5.acl ======================================================="
echo "==================================================================================="
tee /etc/krb5kdc/kadm5.acl <<EOF
elastic@$REALM *
HTTP/elasticsearch_kerberos@$REALM *
noPermissions@$REALM X
EOF
echo ""

echo "==================================================================================="
echo "==== Creating realm ==============================================================="
echo "==================================================================================="
MASTER_PASSWORD=$(tr -cd '[:alnum:]' < /dev/urandom | fold -w30 | head -n1)
# This command also starts the krb5-kdc and krb5-admin-server services
krb5_newrealm <<EOF
$MASTER_PASSWORD
$MASTER_PASSWORD
EOF
echo ""

echo "==================================================================================="
echo "==== Create the principals in the acl ============================================="
echo "==================================================================================="
for principal in $PRINCIPALS; do
  add_principal_to_elastic_realm $principal $PASSWORD
done

kadmin.local -q "list_principals"
cat /etc/hosts

sleep infinity
