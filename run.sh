sh build.sh
echo 
./db2struct --host localhost -d test -t test_table --package myGoPackage --struct testTable -p --user testUser
echo
