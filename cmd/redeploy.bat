cd C:\Users\user1337\Desktop\projects\buycars\tg-service\cmd

echo 1 - check changes in projects and rebuild them
echo 2 - rebuild all projects except db and rabbit
echo 3 - rebuild all containers

set /p user_input=Input the mode number:

if %user_input% equ 1 (
    REM ___________________________________USER_SERVICE__________________________________________________

    cd C:\Users\user1337\Desktop\projects\buycars\user-service\user_service

    git status | findstr "Changes to be committed"

    if %errorlevel% equ 0 (
        docker rm user_service
        docker image rm andrey48358424/user-service:latest
        docker build -t andrey48358424/user-service:latest .
        echo %git_status%
    ) else (
        echo No changes in user_service
    )

    REM ___________________________________TX_SERVICE__________________________________________________

    cd C:\Users\user1337\Desktop\projects\buycars\tx-service
    REM Проверяем, что переменная git_status не пуста

    git status | findstr "Changes to be committed"

    if %errorlevel% equ 0 (
        docker rm tx-service
        docker image rm andrey48358424/tx-service:latest
        docker build -t andrey48358424/tx-service:latest .
        echo %git_status%
    ) else (
        echo No changes in tx_service
    )

    REM ___________________________________TG_SERVICE__________________________________________________

    cd C:\Users\user1337\Desktop\projects\buycars\tg-service
    REM Проверяем, что переменная git_status не пуста

    git status | findstr "Changes to be committed"

    if %errorlevel% equ 0 (
        docker rm tg-service
        docker image rm andrey48358424/tg-service:latest
        docker build -t andrey48358424/tg-service:latest .
        echo %git_status%
    ) else (
        echo No changes in tg_service
    )

    echo Start building compose

    cd C:\Users\user1337\Desktop\projects\buycars\tg-service\cmd

    docker-compose up -d
)

if %user_input% equ 2 (
    REM ___________________________________USER_SERVICE__________________________________________________

    cd C:\Users\user1337\Desktop\projects\buycars\user-service\user_service

    docker rm user_service
    docker image rm andrey48358424/user-service:latest
    docker build -t andrey48358424/user-service:latest .

    REM ___________________________________TX_SERVICE__________________________________________________

    cd C:\Users\user1337\Desktop\projects\buycars\tx-service
    REM Проверяем, что переменная git_status не пуста

    docker rm tx-service
    docker image rm andrey48358424/tx-service:latest
    docker build -t andrey48358424/tx-service:latest .

    REM ___________________________________TG_SERVICE__________________________________________________

    cd C:\Users\user1337\Desktop\projects\buycars\tg-service
    REM Проверяем, что переменная git_status не пуста

    docker rm tg-service
    docker image rm andrey48358424/tg-service:latest
    docker build -t andrey48358424/tg-service:latest .


    echo Start building compose

    cd C:\Users\user1337\Desktop\projects\buycars\tg-service\cmd

    docker-compose up -d
)

if %user_input% equ 3 (

    docker-compose down

    docker-compose -f docker-compose.yaml rm -f

    REM ___________________________________USER_SERVICE__________________________________________________

    cd C:\Users\user1337\Desktop\projects\buycars\user-service\user_service

    docker image rm andrey48358424/user-service:latest
    docker build -t andrey48358424/user-service:latest .

    REM ___________________________________TX_SERVICE__________________________________________________

    cd C:\Users\user1337\Desktop\projects\buycars\tx-service
    REM Проверяем, что переменная git_status не пуста

    docker image rm andrey48358424/tx-service:latest
    docker build -t andrey48358424/tx-service:latest .

    REM ___________________________________TG_SERVICE__________________________________________________

    cd C:\Users\user1337\Desktop\projects\buycars\tg-service
    REM Проверяем, что переменная git_status не пуста

    docker image rm andrey48358424/tg-service:latest
    docker build -t andrey48358424/tg-service:latest .


    echo Start building compose

    cd C:\Users\user1337\Desktop\projects\buycars\tg-service\cmd

    docker-compose up -d
)


