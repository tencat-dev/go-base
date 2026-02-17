{ pkgs, lib, config, inputs, ... }:

{
  # https://devenv.sh/basics/
  env.GOOSE_MIGRATION_DIR = "./migrations";

  # https://devenv.sh/packages/
  packages = [
  	pkgs.git
  	pkgs.yq
  	pkgs.protobuf
  ];

  # https://devenv.sh/languages/
	languages.go = {
		enable = true;
		version = "1.26.0";
	};

  # https://devenv.sh/processes/
  # processes.dev.exec = "${lib.getExe pkgs.watchexec} -n -- ls -la";

  # https://devenv.sh/services/
  # services.postgres.enable = true;

  # https://devenv.sh/scripts/
  scripts.init.exec = ''
    go install github.com/pressly/goose/v3/cmd/goose@latest
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
		go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
		go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
		go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
		go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest
		go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
		go install github.com/goforj/wire/cmd/wire@latest
		go install github.com/stephenafamo/bob/gen/bobgen-psql@latest
		go install github.com/bufbuild/buf/cmd/buf@latest
  '';

  scripts.goget.exec = ''
		go get -u ./...
  '';

  scripts.db-create.exec = ''
    goose -s create "$@" sql
  '';

  scripts.db.exec = ''
		GOOSE_DRIVER=$(yq -r '.data.database.driver' ./configs/config.yaml) \
		GOOSE_DBSTRING=$(yq -r '.data.database.dsn' ./configs/config.yaml) \
		goose -env=none "$@"
  '';

  scripts.db-gen.exec = ''
		PSQL_DSN=$(yq -r '.data.database.dsn' ./configs/config.yaml) \
		bobgen-psql -c ./bobgen.yaml
  '';

  scripts.wire-gen.exec = ''
		wire gen ./...
  '';

  scripts.buf-gen.exec = ''
  	set -e
    (cd api && buf generate --template buf.gen.api.yaml)
    (cd internal/conf && buf generate --template buf.gen.internal.yaml)
  '';

  scripts.proto-add.exec = ''
		kratos proto add "$@"
  '';

  scripts.proto-server.exec = ''
		kratos proto server "$@" -t internal/service
  '';

  scripts.gen.exec = ''
		wire-gen
		buf-gen
		db-gen
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
