# Pomocnik przesiewu jest jedynym kodem tego pakietu liczącym oceny krzyżowego
# kodera pary pytanie-fragment w jednym przebiegu modelu.
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
        # Wyłączenie sieci jest rozstrzygnięciem: wagi stoją na dysku, więc
        # sieć pobrałaby drugą kopię.
        os.environ["HF_HUB_OFFLINE"] = "1"
        skad, miejscowe = katalog, True
    else:
        skad, miejscowe = model, False

    # Katalog pobrania podawany jest jawnie, bo droga wołania procesu nie
    # niesie HOME ani HF_HOME.
    pobranie = None if miejscowe else katalog
    try:
        tokenizator = AutoTokenizer.from_pretrained(skad, local_files_only=miejscowe,
                                                    cache_dir=pobranie)
        koder = AutoModelForSequenceClassification.from_pretrained(
            skad, local_files_only=miejscowe, cache_dir=pobranie)
    except Exception as blad:
        if miejscowe:
            # Odmowa na tej gałęzi różni się od odmowy przy pobieraniu: wagi
            # są już na dysku.
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
            # Krzyżowy koder oddaje jedną liczbę na parę; model o wielu
            # wyjściach jest innym klasyfikatorem.
            raise Brak("wagi", "model " + str(surowe.shape[-1]) +
                       "-wyjściowy nie jest krzyżowym koderem pary pytanie-fragment")
        wynik = surowe.view(-1)
        # Ocena krzyżowego kodera jest logitem; sigmoida sprowadza go do
        # przedziału (0, 1).
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

    # Wykaz tekstów pusty znaczy pytanie, czy koder stoi, a nie polecenie
    # przesiania niczego.
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
