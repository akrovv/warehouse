# Выбор пакетов
1. **Mock** /  **github.com/golang/mock**
**About**: Mock-фреймворк для golang .  
**Why**: Для создания фиктивных объектов (сервисов), которые использовались для тестирования handler'ов.    
**Where**: Весь код находится в **jsonrpc**.
2. **Postgres** / **github.com/lib/pq**
**About**: Библиотека для базы данных PostgreSQL.  
**Why**: Предоставляет драйвер для PostgreSQL.  
**Where**: Весь код находится в **adapters**.
3. **Viper** / **github.com/spf13/viper**
**About**: Библиотека для управления конфигурацией приложений.  
**Why**: Поддержка нескольких форматов файлов конфигурации (в моем случае .YML). Автоматическая привязка переменных окружения. Удобное разархивирование в пользовательские структуры через теги.  
**Where**: Весь код находится в **config**.
4. **ZapLogger** / **go.uber.org/zap**
**About**: Высокопроизводительная библиотека структурированного ведения журналов.  
**Why**: Более быстрое логирование среди остальных библиотек. Удобный формат вывода данных.  
**Where**: Весь код находится в **jsonrpc**.
5. **SQLMock** / **gopkg.in/DATA-DOG/go-sqlmock.v1**
**About**: Библиотека для симулирования поведения реальной БД на уровне sql/driver.  
**Why**: Тестирование функций в **adapters**.  
**Where**: Весь код находится в **jsonrpc**.