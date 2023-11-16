package help

const Help = `			 mygossh

NAME
	mygossh is a smart ssh tool.It is developed by Go,compiled into a separate binary without any dependencies.

DESCRIPTION
		mygossh can do the follow things:
		1.runs cmd on the remote host.
		2.push a local file or path to the remote host.
		3.pull remote host file to local.

USAGE
	1.Single Mode
		remote-comand:
		mygossh -t cmd  -h host -P port(default 22) -u user(default root) -p passswrod [-f] command 

		Files-transfer:   
		<push file>   
		mygossh -t push  -h host -P port(default 22) -u user(default root) -p passswrod [-f] localfile  remotepath 

		<pull file> 
		mygossh -t pull -h host -P port(default 22) -u user(default root) -p passswrod [-f] remotefile localpath 

	2.Batch Mode
		Ssh-comand:
		mygossh -t cmd -i ip_filename -P port(default 22) -u user(default root) -p passswrod [-f] command 

		Files-transfer:   
		mygossh -t push -i ip_filename -P port(default 22) -u user(default root) -p passswrod [-f] localfile  remotepath 
		gosh -t pull -i ip_filename -P port(default 22) -u user(default root) -p passswrod [-f] remotefile localpath

EMAIL
    	email.tata@qq.com 
`
