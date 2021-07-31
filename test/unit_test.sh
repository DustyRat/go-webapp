#!/bin/bash
trap clean_up 0 1 2 3 9 15
clean_up() {
  echo "Removing docker containers..."
  docker compose -f ./test/docker-compose.yaml down -v --rmi all --remove-orphans
}

echo "Starting mock database..."
docker compose -f ./test/docker-compose.yaml up -d mongo

echo "Running unit tests"
docker compose -f ./test/docker-compose.yaml up --build --abort-on-container-exit unit_tests
status=$?

if [ $status -eq 0 ]; then
  echo "
~~~888~~~ 888~~  ,d88~~\\ ~~~888~~~       888~-_        e      ,d88~~\\ ,d88~~\\ 888~~  888~-_
   888    888___ 8888       888          888   \\      d8b     8888    8888    888___ 888   \\
   888    888    \`Y88b      888          888    |    /Y88b    \`Y88b   \`Y88b   888    888    |
   888    888     \`Y88b,    888          888   /    /  Y88b    \`Y88b,  \`Y88b, 888    888    |
   888    888       8888    888          888_-~    /____Y88b     8888    8888 888    888   /
   888    888___ \\__88P'    888          888      /      Y88b \\__88P' \\__88P' 888___ 888_-~
                                                                                              "
else
  echo "
@@@@@@@  @@@@@@@@   @@@@@@   @@@@@@@     @@@@@@@@   @@@@@@   @@@  @@@       @@@@@@@@  @@@@@@@   @@@  @@@
@@@@@@@  @@@@@@@@  @@@@@@@   @@@@@@@     @@@@@@@@  @@@@@@@@  @@@  @@@       @@@@@@@@  @@@@@@@@  @@@  @@@
  @@!    @@!       !@@         @@!       @@!       @@!  @@@  @@!  @@!       @@!       @@!  @@@  @@!  @@!
  !@!    !@!       !@!         !@!       !@!       !@!  @!@  !@!  !@!       !@!       !@!  @!@  !@   !@
  @!!    @!!!:!    !!@@!!      @!!       @!!!:!    @!@!@!@!  !!@  @!!       @!!!:!    @!@  !@!  @!@  @!@
  !!!    !!!!!:     !!@!!!     !!!       !!!!!:    !!!@!!!!  !!!  !!!       !!!!!:    !@!  !!!  !!!  !!!
  !!:    !!:            !:!    !!:       !!:       !!:  !!!  !!:  !!:       !!:       !!:  !!!
  :!:    :!:           !:!     :!:       :!:       :!:  !:!  :!:   :!:      :!:       :!:  !:!  :!:  :!:
   ::     :: ::::  :::: ::      ::        ::       ::   :::   ::   :: ::::   :: ::::   :::: ::   ::   ::
   :     : :: ::   :: : :       :         :         :   : :  :    : :: : :  : :: ::   :: :  :   :::  :::
                                                                                                          "
fi

coverage_file="./test/coverage.out"
if [[ -f "$coverage_file" ]]; then
  coverage=$(go tool cover -func ${coverage_file} | grep total | awk '{print $3}')
  echo "Test Coverage: $coverage"

  echo "To view coverage run:
  go tool cover --html=$coverage_file"
fi

exit $status
