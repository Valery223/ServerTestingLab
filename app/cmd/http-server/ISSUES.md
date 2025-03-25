## План на http серсер
- +Реализовать методы к /user 
    - GET, POST, OPTION +
    - Правильые статус коды ответов +
    - Хорошая маршрутизация +
- Наастроит логи
    - Иницилизацию slog
    - Middleware-log обертки над запросами 
    - Как-то сунуть свой логер в net/http Server(а то он свой использует)
        - "http: TLS handshake error from 127.0.0.1:58092: read tcp 127.0.0.1:8081->127.0.0.1:58092: i/o timeout"
- Тесты
    - Изучить http/httptest

- +Настройка сервера
    - Реализовать Grasefull shutdown
    - KeepALive,  TimeToReadHead etc
    - TLS (HTTPS)
        -генерирует сертификат и ключ: 
        ```bash
        openssl req -x509 -newkey rsa:4096 -nodes -keyout localhost.key -out localhost.crt -days 365 -subj "/CN=localhost"
        ```

- Отправка файлов и HTTP/2.0
    - Сделать index.html и соответствующие ему css файлы
    - Организовать отдачу их
    - Исправить доступ к защищенным ресурсам:
        - Проверить .secret и ../secretDir
        - Возможно ./
         