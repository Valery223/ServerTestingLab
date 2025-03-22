## План на http серсер
- Реализовать методы к /user
    - GET, POST, OPTION
    - Правильые статус коды ответов
    - Хорошая маршрутизация
- Наастроит логи
    - Иницилизацию slog
    - Middleware-log обертки над запросами 
- Тесты
    - Изучить http/httptest

- Настройка сервера
    - Реализовать Grasefull shutdown
    - KeepALive,  TimeToReadHead etc