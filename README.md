Микросервис для анализа доступности страниц

Создать папку в вашей /home/username

mkdir /home/username/demo-service

клонировать репозитарий https:/github.com/spa-nsk/demo-service в созданную папку 
редактировать файл Makefile, заменив username на имя вашей учетной записи

make build

make run

сделать GET запрос на поинт http://127.0.0.1:8080/sites?search=строка для поиска на яндекс
получить ответ в json карта в которой ключи это адреса страниц, а значения это худшее время доступа к старнице, при параллельных запросах с количеством из config.yaml.
