package shortcut

import "time"

const (
	DateTemplate       string = "2006-01-02" // Шаблон даты
	DateTemplatePast   string = "1900-01-01" // Далекое прошлое
	DateTemplateFuture string = "2999-12-31" // Далекое будущее
	DateTimeTemplate   string = "2006-01-02 15:04:05"
)

// Преобразование строки в дату
func StrToDate(value string) time.Time {
	v, _ := time.Parse(DateTemplate, value)

	return v
}

func StrToDateTime(value string) time.Time {
	v, _ := time.Parse(DateTimeTemplate, value)

	return v
}

// Начало времен
func DatePast() time.Time {
	return StrToDate(DateTemplatePast)
}

// Конец времен
func DateFuture() time.Time {
	return StrToDate(DateTemplateFuture)
}

// Текущая дата
func DateCurrent() time.Time {
	current := time.Now().Format(DateTimeTemplate)

	return StrToDateTime(current)
}
