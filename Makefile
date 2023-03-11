
plugin:
	cd plugins/mssql/; go build -gcflags "all=-N -l"   -buildmode=plugin -o ../../plugins_build/mssql.so;  cd -
	


plugin_prod:
	cd plugins/mssql/; go build    -buildmode=plugin -o ../../plugins_build/mssql.so;  cd -
	

build_debug:
	go build -gcflags "all=-N -l"  -o .

build_prod:
	go build    -o .

run_only:
	go run -gcflags "all=-N -l"  .

run: plugin run_only

build: plugin_prod build_prod 

commit: 
	git add . 
	git commit -m "Make by Make with gitall"
	git push

all: plugin build 


