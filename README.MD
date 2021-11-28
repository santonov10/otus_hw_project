https://github.com/OtusGolang/final_project - описание проектной работы

https://github.com/OtusGolang/final_project/blob/master/03-image-previewer.md - ТЗ на сервис "Превьювер изображений"

запуск докер контейнера и развертывания сервера на порту 8080:
- make run

тесты:
- make test

конфиг в папке /configs/default.json:
- "http": 
  - "port": "8080" - порт, на котором будет запущен сервер
- "cacheImagesLru":
  - "capacity" : 10, - количество картинок которое будет хранится в LRU кэше
  - "dir" : "./imgCache" - папка в которой будут храниться кэш изображений

Кешируется каждый запрос (включая на невалидные страницы или на страницы не с изображениями).
Проверка работоспособности после запуска докера:
- http://localhost:8080/fill/300/200/github.com/OtusGolang/final_project/blob/master/03-image-previewer.md - закешируется и выдаст 502 т.к. это не картинка
- http://localhost:8080/fill/300/200/raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg - закешируется и выдаст 200
- http://localhost:8080/fill/2000/2000/raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg - новый кеш, т.к. другие размеры для сервиса
- http://localhost:8080/fill/50/50/raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg
- http://localhost:8080/fill/200/700/raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg
- http://localhost:8080/fill/256/126/raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg
- http://localhost:8080/fill/333/666/raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg
- http://localhost:8080/fill/500/500/raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg
- http://localhost:8080/fill/1024/252/raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg
- http://localhost:8080/fill/2000/1000/raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg
 
