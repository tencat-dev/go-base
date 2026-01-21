{ pkgs, lib, config, inputs, ... }:

{
  # https://devenv.sh/basics/
  env.GREET = "devenv";

  # https://devenv.sh/packages/
  packages = [ pkgs.git ];

  # https://devenv.sh/languages/
	languages.go = {
		enable = true;
		version = "1.25.6";
	};

  # https://devenv.sh/processes/
  # processes.dev.exec = "${lib.getExe pkgs.watchexec} -n -- ls -la";

  # https://devenv.sh/services/
  # services.postgres.enable = true;

  # https://devenv.sh/scripts/
  scripts.init.exec = ''
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
		go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
		go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
		go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
		go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
		go install github.com/google/wire/cmd/wire@latest
		go install github.com/stephenafamo/bob/gen/bobgen-psql@latest
		go install github.com/bufbuild/buf/cmd/buf@latest
  '';

  # https://devenv.sh/basics/
  enterShell = ''
    init         # Run scripts directly
  '';

  # https://devenv.sh/tasks/
  # tasks = {
  #   "myproj:setup".exec = "mytool build";
  #   "devenv:enterShell".after = [ "myproj:setup" ];
  # };

  # https://devenv.sh/tests/
  enterTest = ''
    echo "Running tests"
    git --version | grep --color=auto "${pkgs.git.version}"
  '';

  # https://devenv.sh/git-hooks/
  # git-hooks.hooks.shellcheck.enable = true;

  # See full reference at https://devenv.sh/reference/options/
}
