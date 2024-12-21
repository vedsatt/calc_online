# calc_online
## Описание
**Это простой сервер с калькулятором, который получает POST-запрос с телом:**
```
{"expression": "выражение, которое ввёл пользователь"}
```
И в ответ выдает HTTP-ответ с телом:
```
{"result": "результат выражения"}
```
## Установка
1. Выберите папку, куда хотите установить проект
2. В консоль введите команду для клонирования репозитория: `git clone git@github.com:vedsatt/calc_online.git`
3. Установите зависимости `go mod download`
4. Включите сервер `go run cmd/cmd.go`
6. Ввод запроса с выражением: `curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"ваше выражение\"}" http://127.0.0.1:8080/api/v1/calculate`
## Тестирование
**Безошибочный ввод:**
```
curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"2+2\"}" http://127.0.0.1:8080/api/v1/calculate
```
Ответ сервера:
```
{"result":4}
```
> Код статуса: 200
> 
**Неподдерживаемый метод запроса:**
```
curl http://127.0.0.1:8080/api/v1/calculate
```
Ответ сервера:
```
{"error":"invalid request method"}
```
> Код статуса: 405
> 
**Неподдерживаемое тело запроса:**
```
curl -X POST -H "Content-Type: application/json" -d "2+2" http://127.0.0.1:8080/api/v1/calculate
```
Ответ сервера:
```
{"error":"invalid request body"}
```
> Код статуса: 422
> 
**Ошибка в выражении:**
```
curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"2+*2\"}" http://127.0.0.1:8080/api/v1/calculate
```
Ответ сервера:
```
{"error":"the two operators are next to each other"}
```
> Код статуса: 422
