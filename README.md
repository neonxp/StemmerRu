# Стеммер Портера для русского языка

[![Build Status](https://travis-ci.org/NeonXP/StemmerRu.svg?branch=master)](https://travis-ci.org/NeonXP/StemmerRu)
[![codecov](https://codecov.io/gh/NeonXP/StemmerRu/branch/master/graph/badge.svg)](https://codecov.io/gh/NeonXP/StemmerRu)

Стемминг - процесс получения основы слова из любой его формы. Иными словами, отсекает лишние суффиксы и окончания.

Самое очевидное применение - в полнотекстовом поиске, где нужно, чтобы слово находилось, даже если у него другое окончание.

Этот пакет - реализация [стеммера Портера](https://ru.wikipedia.org/wiki/Стемминг#Стеммер_Портера) для русского языка на Go.

## Использование

`основа := StemmerRu.StemWord("слово")`

Преобразует слово на входе в его основу на выходе

## Пример

```
package main

import (
	"fmt"
	"github.com/neonxp/StemmerRu"
)

func main() {
	fmt.Println(StemmerRu.StemWord("безмолвны") // выведет: безмолвн
	fmt.Println(StemmerRu.StemWord("безмолвные") // выведет: безмолвн
	fmt.Println(StemmerRu.StemWord("безмолвный") // выведет: безмолвн
	fmt.Println(StemmerRu.StemWord("безмолвным") // выведет: безмолвн
	fmt.Println(StemmerRu.StemWord("безмолвных") // выведет: безмолвн
}
```
