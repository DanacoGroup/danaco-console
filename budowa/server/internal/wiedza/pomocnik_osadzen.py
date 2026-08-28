# Pomocnik osadzeń jest jedynym kodem tego pakietu liczącym wektory: liczenie
# stoi obok rdzenia Go, bo modele wydawane są jako wagi ONNX obsługiwane
# bibliotekami Pythona. Pełne uzasadnienie stoi w dokumentacji architektury.
import json
import os
import sys

# Rozmieszczenie plików modelu stojącego, zgodne z eksportem biblioteki sentence-transformers przy zapisie.
UKLADY_WAG = ("onnx/model.onnx", "model.onnx")
# PLIK_MODULOW wylicza warstwy modelu po transformerze — stąd wiadomo, która
# warstwa składa tokeny w wektor i czy wynik jest normalizowany.
PLIK_MODULOW = "modules.json"
# PLIK_USTROJU niesie wymiar: `word_embedding_dimension` w warstwie łączącej
# i `hidden_size` w ustroju samego transformera.
PLIK_USTROJU = "config.json"


def odpowiedz(tresc):
    """Wypisuje jeden obiekt JSON na standardowe wyjście i kończy pracę."""
    sys.stdout.write(json.dumps(tresc, ensure_ascii=False))
    sys.stdout.flush()
    sys.exit(0)


def brak(rodzaj, powod, wagaMb):
    odpowiedz({"ok": False, "brak": rodzaj, "powod": powod, "wagaMb": wagaMb})


class Brak(Exception):
    """Brak nazwany, podniesiony tam, gdzie się go rozpoznaje.

    Zamiana na odpowiedź dzieje się w jednym miejscu (`main`), bo inaczej każdy
    nowy powód wymagałby własnej drogi wyjścia z funkcji, a droga pominięta
    przepuściłaby ślad stosu tam, gdzie miał iść nazwany brak.
    """

    def __init__(self, rodzaj, powod):
        super().__init__(powod)
        self.rodzaj = rodzaj
        self.powod = powod


def wagiWKatalogu(katalog):
    """Oddaje położenie pliku ONNX wewnątrz katalogu wag albo pusty napis.

    Pusty napis znaczy „w tym katalogu modelu nie ma" i kieruje katalog do
    biblioteki jako pamięć podręczną. Sprawdzany jest sam plik wag, a nie opis
    obok niego: katalog pamięci podręcznej też bywa pełen plików JSON po
    przerwanym pobieraniu, a bez wag nie ma czym liczyć.
    """
    if not katalog:
        return ""
    for uklad in UKLADY_WAG:
        if os.path.isfile(os.path.join(katalog, *uklad.split("/"))):
            return uklad
    return ""


def wczytajJson(sciezka, czego):
    """Czyta plik opisu albo podnosi brak nazywający ścieżkę, której zabrakło."""
    try:
        with open(sciezka, "r", encoding="utf-8") as plik:
            return json.load(plik)
    except Exception as blad:  # powód idzie do Operatora w całości
        raise Brak("wagi", czego + " (" + sciezka + "): " + str(blad)) from blad


def ksztaltWag(katalog):
    """Odczytuje z deklaracji przy wagach to, czego biblioteka nie odgadnie.

    Trzy wartości: sposób składania tokenów w jeden wektor, czy wynik jest
    normalizowany i jaki ma wymiar. Wszystkie trzy są ODCZYTANE — przyjęcie
    którejkolwiek z góry byłoby dopowiedzeniem o modelu, którego pomocnik nie
    zna, a skutek — wektory bez znaczenia — wychodzi dopiero na jakości
    wyszukiwania.
    """
    moduly = wczytajJson(os.path.join(katalog, PLIK_MODULOW),
                         "wagi leżą w katalogu, ale nie ma przy nich wykazu warstw modelu")
    if not isinstance(moduly, list):
        raise Brak("wagi", "wykaz warstw modelu w " + os.path.join(katalog, PLIK_MODULOW) +
                   " nie jest listą, więc nie ma z czego odczytać warstwy łączącej")

    podkatalogLaczenia = None
    normalizacja = False
    for warstwa in moduly:
        rodzaj = str((warstwa or {}).get("type") or "")
        if rodzaj.endswith(".Pooling"):
            podkatalogLaczenia = str((warstwa or {}).get("path") or "")
        elif rodzaj.endswith(".Normalize"):
            normalizacja = True
    if podkatalogLaczenia is None:
        raise Brak("wagi", "wykaz warstw w " + os.path.join(katalog, PLIK_MODULOW) +
                   " nie wymienia warstwy łączącej tokeny w jeden wektor")

    laczenie = wczytajJson(
        os.path.join(katalog, *podkatalogLaczenia.split("/"), PLIK_USTROJU),
        "wykaz warstw wskazuje warstwę łączącą, ale nie ma jej opisu")
    if laczenie.get("pooling_mode_cls_token"):
        sposob = "CLS"
    elif laczenie.get("pooling_mode_mean_tokens"):
        sposob = "MEAN"
    else:
        raise Brak("wagi", "opis warstwy łączącej w katalogu " + katalog +
                   " nie wskazuje ani składania po pierwszym tokenie, ani uśrednienia, "
                   "a biblioteka zna wyłącznie te dwa sposoby")

    wymiar = int(laczenie.get("word_embedding_dimension") or 0)
    if wymiar <= 0:
        ustroj = wczytajJson(os.path.join(katalog, PLIK_USTROJU),
                             "opis warstwy łączącej nie podaje wymiaru wektora i nie ma obok "
                             "ustroju transformera, z którego dałoby się go odczytać")
        wymiar = int(ustroj.get("hidden_size") or 0)
    if wymiar <= 0:
        raise Brak("wagi", "ani opis warstwy łączącej, ani ustrój transformera w katalogu " +
                   katalog + " nie podają wymiaru wektora, a wskaźnik zapisuje go przy "
                   "każdej pozycji")
    return sposob, normalizacja, wymiar


def silnikNaWagachStojacych(model, katalog, ukladWag):
    """Składa silnik na wagach, które już leżą w katalogu, bez pobierania.

    Nazwa modelu zostaje ta z nastawy, bo pod nią rdzeń zapisuje każdą pozycję
    wskaźnika i pod nią zawęża odczyt przy szukaniu. Gdy biblioteka zna tę nazwę
    z własnego wykazu, opis bierze się z wykazu; gdy nie zna, opis powstaje
    z deklaracji leżących przy wagach.
    """
    from fastembed import TextEmbedding
    from fastembed.common.model_description import ModelSource, PoolingType

    znane = {str(opis.get("model") or "").lower()
             for opis in TextEmbedding.list_supported_models()}
    if model.lower() not in znane:
        sposob, normalizacja, wymiar = ksztaltWag(katalog)
        # Źródło pobrania jest wymagane przez bibliotekę, ale nigdy nie zostanie użyte.
        TextEmbedding.add_custom_model(
            model=model,
            pooling=PoolingType(sposob),
            normalization=normalizacja,
            sources=ModelSource(hf=model),
            dim=wymiar,
            model_file=ukladWag,
        )
    return TextEmbedding(model_name=model, cache_dir=katalog,
                         specific_model_path=katalog)


def zbudujSilnik(model, katalog):
    """Oddaje silnik osadzeń: na wagach stojących albo pobierający własne."""
    try:
        from fastembed import TextEmbedding
    except Exception as blad:
        raise Brak("biblioteka",
                   "biblioteka fastembed nie stoi w tym interpreterze: " + str(blad)) from blad

    ukladWag = wagiWKatalogu(katalog)
    if ukladWag:
        # Odmowa tutaj jest innym brakiem niż przy pobieraniu: wagi już są na dysku.
        try:
            return silnikNaWagachStojacych(model, katalog, ukladWag)
        except Brak:
            raise
        except Exception as blad:
            raise Brak("wagi", "silnika na wagach z " + katalog +
                       " nie da się postawić: " + str(blad)) from blad
    try:
        return TextEmbedding(model_name=model, cache_dir=katalog)
    except Exception as blad:
        # Tu ląduje najczęściej brak wag: brak sieci albo nazwa nieznana bibliotece.
        raise Brak("model", "modelu " + model + " nie da się przygotować: " + str(blad)) from blad


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
        silnik = zbudujSilnik(model, katalog)
    except Brak as nazwany:
        brak(nazwany.rodzaj, nazwany.powod, wagaMb)

    # Wykaz tekstów pusty znaczy pytanie o gotowość silnika, nie żądanie osadzenia.
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
