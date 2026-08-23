# Pomocnik osadzeń — jedyny kod tego pakietu liczący wektory.
#
# Rdzeń jest w Go, a modele osadzeń wydawane są jako wagi ONNX obsługiwane
# bibliotekami Pythona. Przepisanie inferencji transformera do Go byłoby drugą
# implementacją tej samej rzeczy i rozjeżdżałoby się z wagami przy każdym
# kolejnym wydaniu modelu, więc liczenie wektorów jest procesem obok rdzenia —
# tak samo jak rozpoznawanie mowy (`mowa/pomocnik.go`).
#
# Zlecenie przychodzi ścieżką pliku JSON w argumencie, a odpowiedź wraca jednym
# obiektem JSON na standardowym wyjściu. Rdzeń woła procesy wyłącznie przez
# `zewnetrzne/wolanie.go`, a ta droga nie podaje procesowi standardowego wejścia
# i nie dziedziczy środowiska. Diagnostyka biblioteki idzie na strumień
# diagnostyczny, bo ostrzeżenie wstawione w środek JSON-a uczyniłoby odpowiedź
# nieczytelną.
#
# Bez dziedziczenia środowiska nie ma `HOME` ani `HF_HOME`, więc biblioteka nie
# zna swojego katalogu pamięci podręcznej. Rdzeń podaje katalog modeli wprost
# w zleceniu; katalog ten leży pod katalogiem danych rdzenia, tam gdzie baza
# i magazyn biblioteki, więc model pobrany raz zostaje na dysku.
#
# Brak jest odpowiedzią, a nie wywróceniem: gdy biblioteki nie ma albo wag nie
# da się pobrać, pomocnik oddaje `{"ok": false, "brak": …}` z nazwą braku i wagą
# modelu do dociągnięcia, a kod wyjścia zostaje zerowy. Rdzeń zamienia to na
# odmowę nazywającą brak; sam ślad stosu Pythona nie powiedziałby, ile waży to,
# czego nie ma.
import json
import sys


def odpowiedz(tresc):
    """Wypisuje jeden obiekt JSON na standardowe wyjście i kończy pracę."""
    sys.stdout.write(json.dumps(tresc, ensure_ascii=False))
    sys.stdout.flush()
    sys.exit(0)


def brak(rodzaj, powod, wagaMb):
    odpowiedz({"ok": False, "brak": rodzaj, "powod": powod, "wagaMb": wagaMb})


def main():
    if len(sys.argv) < 2:
        brak("zlecenie", "pomocnik osadzeń uruchomiony bez pliku zlecenia", 0)
    try:
        with open(sys.argv[1], "r", encoding="utf-8") as plik:
            zlecenie = json.load(plik)
    except Exception as blad:  # noqa: BLE001 — powód idzie do Operatora w całości
        brak("zlecenie", "plik zlecenia nieczytelny: " + str(blad), 0)

    model = zlecenie.get("model") or ""
    katalog = zlecenie.get("katalogModeli") or None
    teksty = zlecenie.get("teksty") or []
    wagaMb = int(zlecenie.get("wagaMb") or 0)

    try:
        from fastembed import TextEmbedding
    except Exception as blad:  # noqa: BLE001
        brak(
            "biblioteka",
            "biblioteka fastembed nie stoi w tym interpreterze: " + str(blad),
            wagaMb,
        )

    try:
        silnik = TextEmbedding(model_name=model, cache_dir=katalog)
    except Exception as blad:  # noqa: BLE001
        # Tu ląduje najczęściej brak wag: pierwsze uruchomienie bez sieci albo
        # nazwa modelu, której wydanie biblioteki nie zna.
        brak("model", "modelu " + model + " nie da się przygotować: " + str(blad), wagaMb)

    # Sprawdzenie gotowości bez liczenia. Wykaz tekstów pusty znaczy pytanie
    # „czy silnik stoi", a nie „osadź nic": rdzeń pyta o to przed pierwszym
    # wskaźnikiem, żeby odmówić wcześnie i nazwać brak.
    if not teksty:
        odpowiedz({"ok": True, "model": model, "wymiar": 0, "wektory": []})

    try:
        wektory = [[float(x) for x in w] for w in silnik.embed(teksty)]
    except Exception as blad:  # noqa: BLE001
        brak("liczenie", "osadzanie nie powiodło się: " + str(blad), wagaMb)

    odpowiedz(
        {
            "ok": True,
            "model": model,
            "wymiar": len(wektory[0]) if wektory else 0,
            "wektory": wektory,
        }
    )


main()
