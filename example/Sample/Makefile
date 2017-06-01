
publish:
	 @export COMMAND_CA=cert/ca-command.crt;																					\
	 export COMMAND_CRT=cert/cli.crt;																									\
	 export COMMAND_KEY=cert/cli.key;																									\
	 export COMMAND_ADDR=127.0.0.1:10000;																							\
	 ./warp publish -a sample -p $(platform) -r $(rollout) -v $(version) -u $(upgrades);

publish-root:
	 @export COMMAND_CA=cert/ca-command.crt;																					\
	 export COMMAND_CRT=cert/cli.crt;																									\
	 export COMMAND_KEY=cert/cli.key;																									\
	 export COMMAND_ADDR=127.0.0.1:10000;																							\
	 ./warp publish -a sample -p $(platform) -r $(rollout) -v $(version);
