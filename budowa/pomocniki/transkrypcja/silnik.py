#!/usr/bin/env python3
# -*- coding: utf-8 -*-
# Silnik transkrypcji jest warstwą styku z biblioteką faster-whisper i nie wypisuje
# niczego na standardowe wyjście.
"""Silnik transkrypcji Danaco Console — warstwa styku z faster-whisper.

Moduł nie ma CLI i nie drukuje na stdout. Rozdzielenie jest celowe: stdout
pomocnika jest kanałem danych dla rdzenia i musi zawierać wyłącznie jeden
dokument JSON. Gdyby wykrywanie silnika albo ładowanie modelu miało prawo
cokolwiek dopisać do stdout, rdzeń dostałby śmieci zamiast odpowiedzi. Dlatego
tutaj są wyłącznie funkcje zwracające wartości — o tym, co i gdzie wypisać,
rozstrzyga transkrypcja.py.
"""

from __future__ import annotations

import importlib.util
import os
import platform
from pathlib import Path

# Wykaz przyjmowanych formatów jest zamknięty i krótki: obejmuje kontenery
# wychodzące z nagrywania w przeglądarce i z dyktafonów telefonów, nie wszystko,
# co czyta ffmpeg.
PRZYJMOWANE_ROZSZERZENIA = (".wav", ".ogg", ".m4a", ".webm")

# Nazwa repozytorium modelu jest stała: faster-whisper pobiera wagi z Hugging Face
# Hub, który zapisuje je w katalogu według tego wzorca nazwy.
_SZABLON_KATALOGU = "models--Systran--faster-whisper-{}"


def silnik_obecny() -> bool:
    """Czy moduł faster_whisper da się w ogóle zaimportować."""
    # Sprawdzenie idzie przez find_spec, nie import: import ciągnie ciężką bibliotekę natywną.
    return importlib.util.find_spec("faster_whisper") is not None


def wersja_silnika() -> str:
    """Wersja zainstalowanego pakietu faster-whisper albo pusty napis."""
    try:
        from importlib.metadata import version

        return "faster-whisper " + version("faster-whisper")
    except Exception:
        # Pusty napis zamiast zgadywania: bez metadanych pakietu wersja nie jest znana.
        return ""


def katalogi_cache(katalog_modeli: str = "") -> list[Path]:
    """Miejsca, w których model mógł wylądować, w kolejności pierwszeństwa."""
    if katalog_modeli:
        # Katalog podany wprost wygrywa i jest jedyny, bo download_root wskaże właśnie ten katalog.
        return [Path(katalog_modeli).expanduser()]
    kandydaci = []
    for zmienna in ("HF_HUB_CACHE", "HUGGINGFACE_HUB_CACHE"):
        wartosc = os.environ.get(zmienna)
        if wartosc:
            kandydaci.append(Path(wartosc))
    dom_hf = os.environ.get("HF_HOME")
    if dom_hf:
        kandydaci.append(Path(dom_hf) / "hub")
    kandydaci.append(Path.home() / ".cache" / "huggingface" / "hub")
    return kandydaci


def sciezka_modelu(model: str, katalog_modeli: str = "") -> Path | None:
    """Katalog pobranego modelu albo None, gdy modelu nie widać na dysku."""
    nazwa = _SZABLON_KATALOGU.format(model)
    for katalog in katalogi_cache(katalog_modeli):
        kandydat = katalog / nazwa
        # Sprawdzana jest obecność pliku wag, nie katalogu: przerwane pobieranie zostawia pustą skorupę.
        if kandydat.is_dir() and any(kandydat.rglob("model.bin")):
            return kandydat
        # Model wypakowany "płasko" (bez układu Huba) też jest poprawny.
        if katalog.is_dir() and (katalog / "model.bin").is_file():
            return katalog
    return None


def opis_python() -> str:
    """Wersja Pythona w postaci 3.12.3."""
    return platform.python_version()


def powod_braku_silnika() -> str:
    """Odmowa trójczęściowa dla przypadku: nie ma modułu faster_whisper."""
    # Zdanie rozpoznawcze jest dosłowne i niezmienne: Operator trafia po nim do właściwego akapitu README.
    return (
        "Silnik faster-whisper niedostępny w Pythonie. "
        "Pomocnik transkrypcji nie ma czym rozpoznać mowy, więc nie zgłasza "
        "gotowości — wynik zmyślony byłby gorszy niż jego brak. "
        "Operator naprawi to instalacją zależności: "
        "python -m pip install -r wymagania.txt "
        "(albo wprost: python -m pip install faster-whisper==1.2.1)."
    )


def powod_braku_modelu(model: str, katalog_modeli: str = "") -> str:
    """Odmowa trójczęściowa dla przypadku: silnik jest, modelu nie ma."""
    przeszukane = ", ".join(str(k) for k in katalogi_cache(katalog_modeli))
    return (
        f"Model „{model}” nie jest pobrany na tę maszynę. "
        f"Silnik faster-whisper działa, ale wag modelu nie ma w katalogu "
        f"pobrań ({przeszukane}), a pobrać ich teraz nie wolno — pomocnik "
        f"pracuje lokalnie i nie sięga sam do sieci w trakcie sprawdzenia. "
        f"Operator naprawi to jednorazowym pobraniem modelu przy dostępie "
        f"do sieci: python -c \"from faster_whisper import WhisperModel; "
        f"WhisperModel('{model}', device='cpu', compute_type='int8')\" "
        f"— albo wskazaniem gotowego katalogu przez --katalog-modeli."
    )


def transkrybuj(
    sciezka: Path, model: str, jezyk: str, katalog_modeli: str = ""
) -> tuple[str, float]:
    """Zwraca (tekst, długość nagrania w sekundach). Wyjątki idą w górę."""
    from faster_whisper import WhisperModel

    # CPU i int8 są wpisane na sztywno, bo pomocnik ma działać także bez karty graficznej.
    silnik = WhisperModel(
        model,
        device="cpu",
        compute_type="int8",
        download_root=katalog_modeli or None,
        local_files_only=True,
    )
    segmenty, info = silnik.transcribe(str(sciezka), language=jezyk or None)

    # Segmenty są generatorem; sklejenie idzie tutaj, żeby wywołujący nie znał leniwego strumienia.
    czesci = [s.text for s in segmenty]

    # Długość pochodzi z info.duration, czasu nagrania, nie z czasu przetwarzania materiału.
    return "".join(czesci).strip(), float(getattr(info, "duration", 0.0) or 0.0)
