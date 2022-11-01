## Distributed config
Тестовое для GoCloudCamp

### Getting Started
    docker-compose --project-name="gocloudcamp" up -d

Сервис взаимодействует с клиентами по протоколам HTTP и gRPC, порты по умолчанию 8080 и 50051, настраиваются в env файле. Клиентская библиотека для протокола gRPC в каталоге client, там же пример использования в файле grpc_test.go

##
### Методы gRPC сервера/клиента:
* SetConfig — создать/обновить конфиг
* GetConfig — получить определенную версию конфига
* UpdConfig — установить/сбросить признак использования
* DelConfig — удалить конфиг

##
### Аналогично с использованием HTTP протокола:
`curl -d "@data.json" -H "Content-Type: application/json" -X POST  http://localhost:8080/config`

`curl "http://localhost:8080/config?service=managed-k8s&version=3"`

`curl -X PUT "http://localhost:8080/config?service=managed-k8s&version=2&used=true"`

`curl -X DELETE "http://localhost:8080/config?service=managed-k8s&version=3"`