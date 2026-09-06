# Danaco Console — przekrój pionowy: uwierzytelnienie

Opracowanie opisuje drogę uwierzytelnienia w postaci, w jakiej działa
w bieżącej budowie. Obejmuje komendy bramki dostępu, przechowanie sekretów,
ochronę przed zgadywaniem oraz zmianę adresu konta.

## Bramka dostępu

Bramką dostępu Operatora rządzi rodzina komend `auth.*` kontraktu:
`auth.register`, `auth.verify`, `auth.recover`, `auth.reset`, `auth.login`,
`auth.method.add`, `auth.method.remove`, `auth.password.reset`,
`auth.token.refresh`, `auth.email.change.start`, `auth.email.change.confirm`
i `auth.email.change.revoke`; zmiany stanu bramki ogłasza zdarzenie
`auth.changed`. Obsługa leży w plikach `budowa/server/internal/core/`
o przedrostku `adapter_modul_auth`, wpiętych przez `handlers_auth.go`.
Konto zakłada `auth.register`, a otwiera `auth.login`; potwierdzenie adresu
drogą z listu wykonuje `auth.verify`, które bramkę kanału przechodzi zawsze
(opracowanie [kanał](kanal.md)). Konto założone tam, gdzie nie było czym
wskazać konta nadawczego poczty, nosi znacznik pierwszego uruchomienia bez
poczty (`adapter_modul_auth_pierwsze_uruchomienie.go`).

## Przechowanie sekretów i token sesji

Postać, w jakiej bramka trzyma hasło i PIN, wyprowadza wyłącznie biblioteka
standardowa Go (`adapter_modul_auth_sekret.go`). Token sesji bramki jest
zapisem samoopisującym: niesie parametry wyprowadzenia, więc weryfikacja nie
zależy od nastaw spoza tokenu. Odświeżenie tokenu wykonuje
`auth.token.refresh`.

## Ochrona przed zgadywaniem

Dławik prób wejścia (`adapter_modul_auth_dlawik.go`) narzuca rosnącą zwłokę po
próbach nieudanych i zeruje ją pierwszym wejściem udanym; ogranicza wyłącznie
prędkość zgadywania sekretu, nie blokując konta. Odmowy bramki wracają kodami
kontraktu `not_authenticated` i `permission_denied`.

## Zmiana adresu konta

Zmianę adresu prowadzi trójka komend: `auth.email.change.start` otwiera
zmianę, `auth.email.change.confirm` potwierdza ją drogą z listu,
a `auth.email.change.revoke` wycofuje z listu ostrzegawczego. Obsługę niesie
`adapter_modul_auth_zmiana_adresu.go`, stan trzyma tabela kroku
`migracja_498_zmiana_adresu_konta.sql`, a okno ustawień konta w kliencie
wystawia tę drogę wraz z trasą wycofania.
