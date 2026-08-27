# Pomocnik przesiewu — jedyny kod tego pakietu liczący oceny krzyżowego kodera.
#
# Osadzenia liczy `pomocnik_osadzen.py`, a to jest przebieg drugi. Różnica nie
# jest w modelu, tylko w tym, co model dostaje na wejście: osadzarka widzi
# pytanie i fragment OSOBNO i sprowadza każde z nich do wektora, więc spotykają
# się dopiero jako dwie liczby. Krzyżowy koder czyta pytanie RAZEM z fragmentem
# w jednym przebiegu i oddaje jedną ocenę ich dopasowania. Stąd bierze się
# kolejność inna niż z kosinusa — i stąd bierze się koszt: ocen jest tyle, ile
# kandydatów, a nie jedna na tekst raz na zawsze.
#
# Droga rozmowy jest ta sama co u osadzarki i to jest warunek, nie zbieg
# okoliczności: zlecenie przychodzi ścieżką pliku JSON w argumencie, odpowiedź
# wraca jednym obiektem JSON na standardowym wyjściu, a diagnostyka biblioteki
# idzie na strumień diagnostyczny. Rdzeń woła procesy wyłącznie przez
# `zewnetrzne/wolanie.go` i ta droga nie podaje procesowi standardowego wejścia.
#
# Katalog modeli bywa dwiema różnymi rzeczami i pomocnik je rozróżnia tak samo
# jak osadzarka. Pusty jest MIEJSCEM, do którego biblioteka dopiero pobierze
# wagi. Katalog, w którym wagi już leżą, jest samym MODELEM — wtedy biblioteka
# dostaje go wprost, pobieranie jest wyłączone zmienną `HF_HUB_OFFLINE`, a druga
# kopia tego, co stoi na dysku, nie powstaje. Rozstrzyga obecność pliku
# `model.safetensors`, bo tylko on jest tu wagami.
#
# Kształtu wag pomocnik nie odczytuje z osobnych deklaracji, inaczej niż
# osadzarka. Powód jest w samym modelu: krzyżowy koder jest klasyfikatorem pary,
# więc `config.json` niesie komplet — ustrój transformera wraz z głową oceniającą
# — a warstwy łączącej tokeny w wektor tu po prostu nie ma. Nie ma więc czego
# zgadywać ani skąd doczytywać.
#
# Brak jest odpowiedzią, a nie wywróceniem: gdy biblioteki nie ma albo wag nie
# da się wczytać, pomocnik oddaje `{"ok": false, "brak": …}` z nazwą braku i wagą
# modelu do dociągnięcia, a kod wyjścia zostaje zerowy.
import json
import os
import sys

# Nazwa pliku wag modelu stojącego. Wydanie `safetensors` jest jedynym, jakie
# publikują dziś wydawcy krzyżowych koderów tej rodziny.
PLIK_WAG = "model.safetensors"


def odpowiedz(tresc):
    """Wypisuje jeden obiekt JSON na standardowe wyjście i kończy pracę."""
    sys.stdout.write(json.dumps(tresc, ensure_ascii=False))
    sys.stdout.flush()
    sys.exit(0)


def brak(rodzaj, powod, wagaMb):
    odpowiedz({"ok": False, "brak": rodzaj, "powod": powod, "wagaMb": wagaMb})


class Brak(Exception):
    """Brak nazwany, podniesiony tam, gdzie się go rozpoznaje.

    Zamiana na odpowiedź dzieje się w jednym miejscu (`main`) — z tego samego
    powodu co w pomocniku osadzeń: droga wyjścia pominięta przy nowym powodzie
    przepuściłaby ślad stosu tam, gdzie miał iść nazwany brak.
    """

    def __init__(self, rodzaj, powod):
        super().__init__(powod)
        self.rodzaj = rodzaj
        self.powod = powod


def wagiStoja(katalog):
    """Mówi, czy w katalogu leżą wagi, czy jest on dopiero miejscem na nie."""
    return bool(katalog) and os.path.isfile(os.path.join(katalog, PLIK_WAG))


def zbudujKoder(model, katalog):
    """Oddaje tokenizator i model oceniający parę pytanie–fragment."""
    try:
        import torch
        from transformers import AutoModelForSequenceClassification, AutoTokenizer
    except Exception as blad:
        raise Brak("biblioteka",
                   "bibliotek przesiewu nie ma w tym interpreterze: " + str(blad)) from blad

    if wagiStoja(katalog):
        # Wyłączenie sieci jest tu rozstrzygnięciem, nie ostrożnością: wagi
        # stoją na dysku, więc każde sięgnięcie biblioteki po sieć byłoby
        # pobieraniem drugiej kopii tego, co Operator już ma.
        os.environ["HF_HUB_OFFLINE"] = "1"
        skad, miejscowe = katalog, True
    else:
        skad, miejscowe = model, False

    # Katalog pobrania podawany jest jawnie z tego samego powodu, dla którego
    # robi to pomocnik osadzeń: droga wołania procesu nie niesie `HOME` ani
    # `HF_HOME`, więc bez tego wagi lądowałyby w katalogu pamięci podręcznej
    # zależnym od tego, gdzie akurat stoi rdzeń — a Operator wskazał ustawieniem
    # miejsce, w którym mają leżeć.
    pobranie = None if miejscowe else katalog
    try:
        tokenizator = AutoTokenizer.from_pretrained(skad, local_files_only=miejscowe,
                                                    cache_dir=pobranie)
        koder = AutoModelForSequenceClassification.from_pretrained(
            skad, local_files_only=miejscowe, cache_dir=pobranie)
    except Exception as blad:
        if miejscowe:
            # Odmowa na tej gałęzi jest innym brakiem niż na gałęzi pobierania:
            # wagi są na dysku, więc odesłanie Operatora po pobranie kierowałoby
            # go po to, co już ma.
            raise Brak("wagi", "krzyżowego kodera na wagach z " + katalog +
                       " nie da się postawić: " + str(blad)) from blad
        raise Brak("model", "modelu " + model + " nie da się przygotować: " +
                   str(blad)) from blad
    koder.eval()
    return torch, tokenizator, koder


def oceny(torch, tokenizator, koder, pytanie, teksty, oknoTokenow):
    """Oddaje po jednej ocenie na tekst, w kolejności tekstów.

    Kolejność jest warunkiem, a nie wygodą: wołający wiąże ocenę z kandydatem
    pozycją, nie treścią — dwa fragmenty tego samego dokumentu bywają identyczne
    co do znaku.
    """
    with torch.no_grad():
        wejscie = tokenizator([pytanie] * len(teksty), teksty, padding=True,
                              truncation=True, max_length=oknoTokenow,
                              return_tensors="pt")
        surowe = koder(**wejscie).logits.float()
        if surowe.shape[-1] != 1:
            # Krzyżowy koder oddaje JEDNĄ liczbę na parę. Model o wielu wyjściach
            # jest klasyfikatorem czegoś innego i nie wiadomo, która jego kolumna
            # miałaby znaczyć dopasowanie — zgadnięcie dałoby ranking, który
            # wygląda jak wynik.
            raise Brak("wagi", "model " + str(surowe.shape[-1]) +
                       "-wyjściowy nie jest krzyżowym koderem pary pytanie-fragment")
        wynik = surowe.view(-1)
        # Ocena krzyżowego kodera jest logitem — liczbą bez ograniczenia z obu
        # stron. Sigmoida sprowadza ją do przedziału (0, 1), w którym rdzeń
        # przelicza trafność na setne kontraktu. Sama kolejność się przez to nie
        # zmienia, bo sigmoida jest rosnąca.
        return [float(x) for x in torch.sigmoid(wynik)]


def main():
    if len(sys.argv) < 2:
        brak("zlecenie", "pomocnik przesiewu uruchomiony bez pliku zlecenia", 0)
    try:
        with open(sys.argv[1], "r", encoding="utf-8") as plik:
            zlecenie = json.load(plik)
    except Exception as blad:  # noqa: BLE001 — powód idzie do Operatora w całości
        brak("zlecenie", "plik zlecenia nieczytelny: " + str(blad), 0)

    model = zlecenie.get("model") or ""
    katalog = zlecenie.get("katalogModeli") or None
    pytanie = zlecenie.get("pytanie") or ""
    teksty = zlecenie.get("teksty") or []
    wagaMb = int(zlecenie.get("wagaMb") or 0)
    oknoTokenow = int(zlecenie.get("oknoTokenow") or 512)

    try:
        torch, tokenizator, koder = zbudujKoder(model, katalog)
    except Brak as nazwany:
        brak(nazwany.rodzaj, nazwany.powod, wagaMb)

    # Wykaz tekstów pusty znaczy pytanie „czy koder stoi", a nie „przesiej nic":
    # rdzeń pyta o to, zanim przeczyta wskaźnik, żeby odmówić wcześnie i nazwać
    # brak. Ta sama umowa obowiązuje w pomocniku osadzeń.
    if not teksty:
        odpowiedz({"ok": True, "model": model, "oceny": []})

    try:
        wyliczone = oceny(torch, tokenizator, koder, pytanie, teksty, oknoTokenow)
    except Brak as nazwany:
        brak(nazwany.rodzaj, nazwany.powod, wagaMb)
    except Exception as blad:  # noqa: BLE001
        brak("liczenie", "przesiew nie powiódł się: " + str(blad), wagaMb)

    odpowiedz({"ok": True, "model": model, "oceny": wyliczone})


main()
