package middleware

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

/*
Этот код реализует механизм логирования тела HTTP-запросов и ответов на сервере.
Его основная цель — перехватить данные, которые отправляются в запросе от клиента,
и данные, которые возвращает сервер в ответ, с последующей записью этой информации в лог.

### Логирование тела запроса:

1. **Чтение тела запроса:**
   Тело запроса передается через `r.Body`. Однако в Go его можно прочитать только один раз.
   После чтения содержимое `r.Body` становится недоступным для других частей программы.
   Поэтому, чтобы корректно логировать содержимое тела запроса и не нарушить его дальнейшую обработку,
   тело запроса сначала читается с помощью `io.ReadAll`. Это позволяет нам получить все данные, отправленные клиентом.

2. **Логирование тела запроса:**
   После успешного чтения содержимого оно преобразуется в строку и очищается от лишних пробелов и символов переноса строки.
   Если тело запроса не пустое, оно добавляется в строку лога для записи. Это полезно для отладки и мониторинга,
   так как позволяет видеть, какие данные были отправлены клиентом на сервер.

3. **Восстановление тела запроса:**
   Поскольку тело запроса было прочитано, его необходимо восстановить, чтобы последующие обработчики могли работать с ним.
   Для этого тело запроса помещается обратно в объект `ReadCloser` с помощью `io.NopCloser` и нового `strings.Reader`.
   Это важно, так как многие обработчики ожидают наличие тела запроса для выполнения основной логики обработки.

### Логирование ответа сервера:

1. **Обертка ResponseWriter:**
   Чтобы логировать ответ, создается специальная структура `responseCapture`, которая оборачивает оригинальный `ResponseWriter`.
   Эта структура используется для захвата тела ответа и статуса ответа, которые записываются сервером для отправки клиенту.

2. **Запись тела ответа:**
   Вместо того чтобы данные сразу отправлялись клиенту, они сначала сохраняются в буфер `responseCapture`,
   что позволяет нам захватить данные, которые отправляются клиенту. Это полезно для логирования,
   так как можно увидеть не только запрос, но и то, что сервер вернул в ответ на этот запрос.

3. **Логирование ответа:**
   После того как основной обработчик завершил свою работу и записал ответ в буфер,
   данные о статусе ответа и содержимом тела логируются. Это помогает отследить, какой статус и данные были отправлены клиенту,
   что может быть крайне полезно для отладки и анализа работы системы.
*/

type LogMux struct {
	h http.Handler
}

func NewLogMux(h http.Handler) http.Handler {
	return &LogMux{h: h}
}

type responseCapture struct {
	http.ResponseWriter
	body *bytes.Buffer
	code int
}

func (rc *responseCapture) WriteHeader(statusCode int) {
	rc.code = statusCode
	rc.ResponseWriter.WriteHeader(statusCode)
}

func (rc *responseCapture) Write(b []byte) (int, error) {
	rc.body.Write(b)
	return rc.ResponseWriter.Write(b)
}

func (m *LogMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	clientIP := r.RemoteAddr
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		clientIP = ip
	}

	rqMessageLog := fmt.Sprintf("Request: Method: %s Path: %s ClientIp: %s ", r.Method, r.URL.Path, clientIP)

	if r.Body != nil {
		bodyBytes, err := io.ReadAll(r.Body)
		if err == nil {
			body := strings.TrimSpace(string(bodyBytes))
			if len(body) > 0 {
				rqMessageLog += fmt.Sprintf("Body: %s\n", body)
			}
		}
		r.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
	}

	rc := &responseCapture{
		ResponseWriter: w,
		body:           bytes.NewBuffer(nil),
	}

	m.h.ServeHTTP(rc, r)

	responseLog := fmt.Sprintf("Response: StatusCode: %d Body: %s", rc.code, rc.body.String())

	duration := time.Since(startTime)

	fullLog := fmt.Sprintf("%s%s Duration: %v", rqMessageLog, responseLog, duration)
	log.Println("[middleware] " + fullLog)
}
