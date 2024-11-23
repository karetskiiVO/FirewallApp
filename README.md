# FirewallApp

## Назначение

Утилита для установки black/white-списка

## Сборка

``` go build -o <appname> ```

## Запуск

``` <appname> <dev1 name> <dev2 name> <config file> ```

Где ```<dev1 name>``` и ```<dev2 name>``` - устройства которые соединяются программным мостом.

## Настройка

```<config file> ``` - json-файл с правилами из следующего вида:

* Сетевые:
    
    1. ```ip``` - фильтр на cетевом уровне:
        ```
        {
            "type"="ip",
            "source":"<source ip>",          # optional
            "destination":"<destination ip>" # optional
        }
        ```
    2. ```udp``` - фильтр на транспортном уровне над ```udp```
        ```
        {
            "type"="udp",
            "source":"<source port>",          # optional
            "destination":"<destination port>" # optional
        }
        ```
    3. ```tcp```- фильтр на транспортном уровне над ```tcp```
        ```
        {
            "type"="tcp",
            "source":"<source ip>",          # optional
            "destination":"<destination ip>" # optional
        }
        ```

* Логические

    1. ```and``` - логическое И правил
        ```
        {
            "type":"and",
            "rules":[
                <rule 1>,
                <rule 2>,
                ...
            ]
        }
        ```
    2. ```or``` - логическое ИЛИ правил
        ```
        {
            "type":"and",
            "rules":[
                <rule 1>,
                <rule 2>,
                ...
            ]
        }
        ```

### Общая структура файла конфигурации
```
    {
        "default": "drop", # что происходит с пакетами по умолчанию, либо "drop", либо "accept"
        "rules":[
            <rule 1>,
            <rule 2>,
            ...
        ]
    }
```


