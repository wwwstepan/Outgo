# Outgo (Учёт расходов)
Outgo - expense accounting

Консольная программа под Windows, управляется исключительно цифровой клавиатурой, что помогает достигать высокой скорости ввода данных. 
Есть небольшая кастомизация &ndash; список статей расходов можно переопределить в файле JSON-формата. 
Управление хранением данных &ndash; самописное. Все записи о всех расходах при каждом запуске программы читаются полностью в память, при выходе полностью перезаписываются. Для целей этой программы при объемах, превышающих разумные в сотни раз, такой подход не сказываться на скорости, то есть задержка не заметна на глаз. Особое внимание уделено размеру. На одну запись о расходе (дата, сумма, статья расходов) требуется всего 8 байт, достигается это побитовой упаковкой значений.</p>

2021 год. Golang
