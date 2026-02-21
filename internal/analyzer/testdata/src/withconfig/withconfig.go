package withconfig

import "log"

func example() {
    // first_letter отключён — ошибки нет
    log.Println("Hello world")

    // english отключён — ошибки нет
    log.Println("Привет мир")

    // emoji включён — ошибка
    log.Println("check status 🚀") // want `log message must not contain emoji`

    // кастомный keyword "cvv" в литерале — ошибка
    log.Println("processing cvv") // want `log message may expose sensitive data`

    // кастомный keyword "ssn" через имя переменной — ошибка
    ssnVar := "secret"
    _ = ssnVar
    log.Println("value " + ssnVar) // want `log message may expose sensitive data`
}