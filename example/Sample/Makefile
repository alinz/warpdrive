publish-root:
	@./warp publish																																		\
		--addr 127.0.0.1:10000																													\
		--ca-path cert/ca-command.crt																										\
		--cert-path cert/cli.crt																												\
		--key-path cert/cli.key																													\
		--app sample																																		\
		--platform $(platform)																													\
		--rollout $(rollout)																														\
		--version $(version);

publish:
	@./warp publish																																		\
		--addr 127.0.0.1:10000																													\
		--ca-path cert/ca-command.crt																										\
		--cert-path cert/cli.crt																												\
		--key-path cert/cli.key																													\
		--app sample																																		\
		--platform $(platform)																													\
		--rollout $(rollout)																														\
		--version $(version)																														\
		--root $(root)																																	\
		--upgrade $(upgrade)																														\
		--notes $(notes);
