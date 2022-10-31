WORKDIR=$(pwd)
cd .audit/gotest
./coverage_gen.bash
cd $WORKDIR
sonar-scanner -Dproject.settings=.audit/sonarqube/sonar-project.properties
