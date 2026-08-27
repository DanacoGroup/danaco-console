# Uzasadnienia pracy twórczej

Dokument gromadzi uzasadnienia decyzji projektowych, które w kodzie źródłowym
nie mieszczą się w nagłówku komentarza. Każdy rozdział nosi nazwę pliku,
którego uzasadnienia dotyczą.

## adapter_modul_studio_ooxml.go

Rachunek nie korzysta z żadnej biblioteki obcej ponad standardową bibliotekę
Go. `.docx` jest archiwum ZIP z dokumentem XML w środku, a biblioteka wzorcowa
Go niesie oba potrzebne składniki — `archive/zip` i `encoding/xml`. Biblioteka
`unioffice` odpada na licencji handlowej, a LibreOffice i `pandoc` byłyby
procesem potomnym tam, gdzie proces potomny jest regresem: jedno binarium
serwera zamienia się w dwa, dochodzi koszt uruchomienia i rozjazd wersji
między maszynami.

Rozbiór dokumentu idzie przez drzewo węzłów, nie przez strumień tokenów
z ręcznym stanem. Postać akapitu OOXML leży w węźle `w:pPr`, który stoi PRZED
fragmentami tekstu, ale odnosi się do całego akapitu, a postać znaku leży
w `w:rPr` wewnątrz każdego fragmentu. Automat stanów musiałby pamiętać oba
konteksty naraz na każdej głębokości; drzewo pozwala je odczytać tam, gdzie
są. Ten sam rozbiór obsługuje ODF, bo i tam postać leży w węzłach obok treści.

OOXML mierzy odległości w twipach (1/1440 cala), stopień pisma w półpunktach,
a odstępy akapitowe w twipach, podczas gdy kontrakt platformy mówi
w milimetrach i punktach. Przeliczenie stoi w jednym miejscu rachunku, bo
przeliczenie rozsypane po wielu miejscach rozjeżdża się przy pierwszej
poprawce.

Bilans wniesienia dokumentu jest obowiązkowy, nie ozdobny: plik wejściowy
niesie rzeczy, których rachunek nie odczytuje (pola obliczane, wykresy,
kształty rysowane, osadzone obiekty innych programów), a przemilczenie ich
zamieniłoby dokument okaleczony w dokument wczytany bez uwag.

Zasięgi nagłówka i stopki — domyślny, pierwszej strony i stron parzystych —
idą osobno, bo nagłówek urzędowy bywa dla tych trzech przypadków różny.
Zasięg pierwszej strony dokłada do sekcji przełącznik `w:titlePg`, bez
którego dokument nagłówka pierwszej strony nie pokaże.
