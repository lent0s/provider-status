# ​ ​ ​ ​ ​ provider-status
​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​
*Network multithreaded microservice for communication provider*  
​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​
![picFile](https://raw.githubusercontent.com/lent0s/provider-status/main/doc/gears.png)

---

### Инструкция по запуску:

  1. Клонировать репозиторий командой:  
     `git clone https://github.com/lent0s/provider-status`
  2. Микросервис получает данные из [симулятора](#sim). Настройте его и запустите
  3. Настройте [config.ini](https://github.com/lent0s/provider-status/blob/main/cmd/config.ini)
  4. Запустите микросервис командой из директории проекта:  
     `go run main.go`
     
---

#### Описание основных файлов проекта:

​     .\main.go ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ // запуск сервера сбора данных  
​     .\cmd\config.ini    ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ // указание путей для *.data и задание хостов  
​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​\lib\countryISO3166_1.txt               ​ ​ // актуальный на 01.04.2023 список стран по стандарту ISO3166-1  

---

#### Особенности по работе с симулятором:

* *версия страницы вывода данных с Email по всем странам:​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ 
[status_pageFull.html](https://github.com/lent0s/provider-status/skillbox-diploma/status_pageFull.html)  
(для настройки входных данных используется файл
[main.js](https://github.com/lent0s/provider-status/skillbox-diploma/main.js))*
* *версия страницы вывода данных с Email по конкретной стране:​​ ​  ​ 
[status_page.html](https://github.com/lent0s/provider-status/skillbox-diploma/status_page.html)  
(для настройки входных данных используется файл
[mainWeb.js](https://github.com/lent0s/provider-status/skillbox-diploma/mainWeb.js))*  


​      http://host/ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ // выводит данные с записями Email для всех стран  
​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ /?c=X  ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​    // выводит данные с записями Email для Х страны  
​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ ​ // (Х - код страны по ISO3166-1)

---

## <a name="sim">Симулятор</a>  

Проект симулятора данных с github в отдельной директории 
[skillbox-diploma](https://github.com/lent0s/provider-status/tree/main/skillbox-diploma).  
В соответствии с файлом Readme в проекте симулятора запустите проект.  
Проект сгенерирует нужные файлы и продолжит работать для обращения к нему  
по API для получения дополнительных данных.

---

#### Описание задания в [файле](https://github.com/lent0s/provider-status/blob/main/doc/task.pdf)
