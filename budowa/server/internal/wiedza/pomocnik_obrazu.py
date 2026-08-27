# Pomocnik osi obrazu — jedyny kod tego pakietu porównujący zdanie z obrazem.
#
# Model tej osi jest DWUWEJŚCIOWY: ma osobną wieżę dla obrazu i osobną dla
# tekstu, a obie kończą w jednej przestrzeni, w której iloczyn skalarny znaczy
# „to zdanie opisuje ten obraz". Dlatego nie da się tu użyć osadzarki tekstu:
# jej wektor leży w przestrzeni, w której obrazu nie ma i nigdy nie było, a
# porównanie dałoby liczbę bez związku z czymkolwiek.
#
# Obrazy przychodzą ścieżkami plików, nie bajtami. Bajty biblioteki leżą
# w magazynie treści rdzenia i tamtą ścieżką czyta je podgląd modułu Library;
# przepisanie ich do zlecenia oznaczałoby drugi komplet obrazów w pliku JSON,
# rosnący w megabajtach na każde zapytanie.
#
# Droga rozmowy jest ta sama co u osadzarki i u przesiewu: zlecenie ścieżką
# pliku JSON w argumencie, odpowiedź jednym obiektem JSON na standardowym
# wyjściu, diagnostyka biblioteki na strumieniu diagnostycznym.
#
# Katalog modeli bywa dwiema rzeczami i pomocnik je rozróżnia tak samo jak
# osadzarka: pusty jest MIEJSCEM na wagi, a katalog z wagami jest samym MODELEM
# i wtedy pobieranie jest wyłączone. Rozstrzyga obecność `model.safetensors`.
#
# Obraz nieczytelny nie przerywa zapytania. Jeden plik uszkodzony albo w postaci,
# której biblioteka obrazu nie otwiera, nie ma prawa odebrać Operatorowi
# odpowiedzi o pozostałych — dostaje ocenę zerową i wraca w wykazie `pominiete`,
# żeby rdzeń wiedział, że nie porównał wszystkiego, o co prosił.
import json
import os
import sys

# Nazwa pliku wag modelu stojącego — patrz nagłówek.
PLIK_WAG = "model.safetensors"


def odpowiedz(tresc):
    """Wypisuje jeden obiekt JSON na standardowe wyjście i kończy pracę."""
    sys.stdout.write(json.dumps(tresc, ensure_ascii=False))
    sys.stdout.flush()
    sys.exit(0)


def brak(rodzaj, powod, wagaMb):
    odpowiedz({"ok": False, "brak": rodzaj, "powod": powod, "wagaMb": wagaMb})


class Brak(Exception):
    """Brak nazwany, podniesiony tam, gdzie się go rozpoznaje."""

    def __init__(self, rodzaj, powod):
        super().__init__(powod)
        self.rodzaj = rodzaj
        self.powod = powod


def wagiStoja(katalog):
    """Mówi, czy w katalogu leżą wagi, czy jest on dopiero miejscem na nie."""
    return bool(katalog) and os.path.isfile(os.path.join(katalog, PLIK_WAG))


def zbudujModel(model, katalog):
    """Oddaje przygotowywacz wejścia i model dwuwieżowy osi obrazu."""
    try:
        import torch
        from PIL import Image  # noqa: F401 — sprawdzenie obecności, użycie niżej
        from transformers import AutoModel, AutoProcessor
    except Exception as blad:
        raise Brak("biblioteka",
                   "bibliotek osi obrazu nie ma w tym interpreterze: " + str(blad)) from blad

    if wagiStoja(katalog):
        os.environ["HF_HUB_OFFLINE"] = "1"
        skad, miejscowe = katalog, True
    else:
        skad, miejscowe = model, False

    # Katalog pobrania podawany jest jawnie — powód ten sam co w pomocniku
    # osadzeń: bez `HOME` biblioteka nie zna swojego katalogu pamięci podręcznej,
    # a Operator wskazał ustawieniem miejsce, w którym wagi mają leżeć.
    pobranie = None if miejscowe else katalog
    try:
        przygotowywacz = AutoProcessor.from_pretrained(skad, local_files_only=miejscowe,
                                                       cache_dir=pobranie)
        siec = AutoModel.from_pretrained(skad, local_files_only=miejscowe,
                                         cache_dir=pobranie)
    except Exception as blad:
        if miejscowe:
            raise Brak("wagi", "modelu osi obrazu na wagach z " + katalog +
                       " nie da się postawić: " + str(blad)) from blad
        raise Brak("model", "modelu " + model + " nie da się przygotować: " +
                   str(blad)) from blad
    siec.eval()
    return torch, przygotowywacz, siec


def wczytajObrazy(sciezki):
    """Otwiera obrazy, oddając wykaz otwartych wraz z pozycjami pominiętych."""
    from PIL import Image

    otwarte, pozycje, pominiete = [], [], []
    for numer, sciezka in enumerate(sciezki):
        try:
            with Image.open(sciezka) as obraz:
                # Kopia w RGB powstaje wewnątrz `with`, bo biblioteka czyta plik
                # leniwie i po zamknięciu uchwytu nie miałaby skąd wziąć pikseli.
                otwarte.append(obraz.convert("RGB"))
            pozycje.append(numer)
        except Exception as blad:  # noqa: BLE001 — jeden plik nie psuje zapytania
            pominiete.append({"obraz": sciezka, "powod": str(blad)})
    return otwarte, pozycje, pominiete


def podobienstwa(torch, przygotowywacz, siec, pytanie, obrazy):
    """Oddaje kosinus zdania z każdym obrazem, w kolejności obrazów.

    Model liczy się raz na obrazy i raz na zdanie, w jednym przebiegu: obie
    wieże są częścią tego samego bytu, a rozdzielenie wywołań nie oszczędziłoby
    niczego poza jednym przejściem przez wieżę tekstu.
    """
    with torch.no_grad():
        wejscie = przygotowywacz(images=obrazy, text=[pytanie], padding=True,
                                 truncation=True, return_tensors="pt")
        wynik = siec(**wejscie)
        wektoryObrazow = wynik.image_embeds
        wektorZdania = wynik.text_embeds
        wektoryObrazow = wektoryObrazow / wektoryObrazow.norm(dim=-1, keepdim=True)
        wektorZdania = wektorZdania / wektorZdania.norm(dim=-1, keepdim=True)
        return [float(x) for x in (wektoryObrazow @ wektorZdania.T).view(-1)]


def main():
    if len(sys.argv) < 2:
        brak("zlecenie", "pomocnik osi obrazu uruchomiony bez pliku zlecenia", 0)
    try:
        with open(sys.argv[1], "r", encoding="utf-8") as plik:
            zlecenie = json.load(plik)
    except Exception as blad:  # noqa: BLE001 — powód idzie do Operatora w całości
        brak("zlecenie", "plik zlecenia nieczytelny: " + str(blad), 0)

    model = zlecenie.get("model") or ""
    katalog = zlecenie.get("katalogModeli") or None
    pytanie = zlecenie.get("pytanie") or ""
    sciezki = zlecenie.get("obrazy") or []
    wagaMb = int(zlecenie.get("wagaMb") or 0)

    try:
        torch, przygotowywacz, siec = zbudujModel(model, katalog)
    except Brak as nazwany:
        brak(nazwany.rodzaj, nazwany.powod, wagaMb)

    # Wykaz obrazów pusty znaczy pytanie „czy model stoi", tak samo jak pusty
    # wykaz tekstów u osadzarki.
    if not sciezki:
        odpowiedz({"ok": True, "model": model, "oceny": [], "pominiete": []})

    otwarte, pozycje, pominiete = wczytajObrazy(sciezki)
    oceny = [0.0] * len(sciezki)
    if otwarte:
        try:
            wyliczone = podobienstwa(torch, przygotowywacz, siec, pytanie, otwarte)
        except Exception as blad:  # noqa: BLE001
            brak("liczenie", "porównanie zdania z obrazami nie powiodło się: " +
                 str(blad), wagaMb)
        for numer, ocena in zip(pozycje, wyliczone):
            oceny[numer] = ocena

    odpowiedz({"ok": True, "model": model, "oceny": oceny, "pominiete": pominiete})


main()
