#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""Pomocnik transkrypcji Danaco Console — warstwa wiersza poleceń.

Bierze plik dźwiękowy z dysku, rozpoznaje w nim mowę lokalnie i wypisuje na
standardowe wyjście jeden dokument JSON.

Ograniczenia wpisane w konstrukcję pomocnika:
  * Nagranie nie opuszcza maszyny — nie ma tu wywołania sieciowego ani importu
    biblioteki obsługującej sieć.
  * Praca wyłącznie na procesorze: `device="cpu"` i `compute_type="int8"` są
    wpisane na stałe w `silnik.py`. Wolniej, ale jednakowo na każdej maszynie.
  * Tryb `--wersja` niczego nie pobiera, tylko patrzy na dysk — sprawdzenie
    gotowości nie może zająć maszyny na kilka minut ani ściągnąć wag bez pytania.
  * Brak silnika i brak modelu to dwie osobne odpowiedzi z dwiema osobnymi
    naprawami, nie jedno wspólne niepowodzenie.

Umowa ze standardowym wyjściem: idzie tam wyłącznie JSON, bo rdzeń parsuje
stdout. Ostrzeżenia, ślady wyjątków i komunikaty bibliotek idą na stderr;
pojedynczy `print()` na stdout psuje odczyt po stronie rdzenia.
"""

from __future__ import annotations

import argparse
import json
import sys
from pathlib import Path

# Import „obok siebie” — pomocnik jest wołany po ścieżce pliku, nie jako
# zainstalowany pakiet, więc katalog skryptu musi być na ścieżce importu.
sys.path.insert(0, str(Path(__file__).resolve().parent))

import silnik as _silnik  # noqa: E402

DOMYSLNY_MODEL = "small"
DOMYSLNY_JEZYK = "pl"

# Stany wyniku transkrypcji. Rozróżniają przetworzenie udane, w którym mowy nie
# było, od przetworzenia, które się nie odbyło — to ostatnie idzie osobnym polem
# `blad` i tylko ono jest odmową.
STAN_ROZPOZNANO = "rozpoznano"
STAN_BEZ_MOWY = "bez_mowy"


def _wypisz(dane: dict) -> None:
    """Jedyne wyjście na stdout w całym pomocniku."""
    # ensure_ascii=False zostawia polskie znaki w postaci czytelnej; domyślne
    # escapowanie zamieniłoby je na sekwencje \uXXXX i odmowa w dzienniku
    # wymagałaby odkodowania, zanim dałoby się ją przeczytać.
    json.dump(dane, sys.stdout, ensure_ascii=False)
    sys.stdout.write("\n")


def _sprawdz_plik(sciezka: Path) -> str:
    """Zwraca odmowę trójczęściową albo pusty napis, gdy plik jest w porządku."""
    wykaz = " ".join(_silnik.PRZYJMOWANE_ROZSZERZENIA)
    if not sciezka.exists():
        return (
            f"Plik wejściowy nie istnieje: {sciezka}. "
            f"Pomocnik nie ma czego transkrybować i nie zwróci pustego tekstu "
            f"udającego ciszę. "
            f"Operator poprawi to podaniem istniejącej ścieżki w --plik."
        )
    if not sciezka.is_file():
        return (
            f"Ścieżka nie wskazuje pliku: {sciezka}. "
            f"To katalog albo obiekt specjalny, a rozpoznanie mowy wymaga "
            f"pojedynczego nagrania. "
            f"Operator poprawi to wskazaniem konkretnego pliku w --plik."
        )
    if sciezka.stat().st_size == 0:
        return (
            f"Plik wejściowy jest pusty (0 bajtów): {sciezka}. "
            f"Zerowej długości nagranie nie zawiera mowy, a wynik „” podany "
            f"jako udana transkrypcja byłby atrapą. "
            f"Operator poprawi to nagraniem zapisanym do końca — najczęściej "
            f"plik jest pusty, bo nagrywanie przerwano przed zamknięciem."
        )
    if sciezka.suffix.lower() not in _silnik.PRZYJMOWANE_ROZSZERZENIA:
        return (
            f"Format pliku nie jest przyjmowany: {sciezka} "
            f"(rozszerzenie „{sciezka.suffix or 'brak'}”). "
            f"Pomocnik przyjmuje wyłącznie kontenery, które potrafi otworzyć "
            f"pewnie, i nie udaje obsługi pozostałych. "
            f"Operator poprawi to konwersją nagrania do jednego z formatów: "
            f"{wykaz}."
        )
    return ""


def tryb_wersja(model: str, katalog_modeli: str) -> int:
    """Sprawdzenie gotowości. Zawsze kod 0 — brak jest odpowiedzią, nie awarią."""
    # Kod 0 obowiązuje także przy `gotowy=false`: rdzeń nie odróżniłby
    # zdiagnozowanego braku od pomocnika, który przerwał pracę bez odpowiedzi,
    # gdyby oba przypadki kończyły się kodem niezerowym.
    odpowiedz = {
        "gotowy": False,
        "python": _silnik.opis_python(),
        "silnik": "",
        "model": "",
        "powod": "",
    }

    if not _silnik.silnik_obecny():
        odpowiedz["powod"] = _silnik.powod_braku_silnika()
        _wypisz(odpowiedz)
        return 0

    odpowiedz["silnik"] = _silnik.wersja_silnika()

    # Obecność modelu rozpoznaje się po pliku na dysku, bez ładowania wag:
    # załadowanie modelu zajmuje kilkaset megabajtów pamięci i kilka sekund.
    # Ładowanie zostaje jako droga zapasowa dla nietypowych układów katalogów,
    # a jego wyjątek trafia do pola „powod”, nie na wyjście jako ślad stosu.
    if _silnik.sciezka_modelu(model, katalog_modeli) is not None:
        odpowiedz["gotowy"] = True
        odpowiedz["model"] = model
        _wypisz(odpowiedz)
        return 0

    try:
        from faster_whisper import WhisperModel

        WhisperModel(
            model,
            device="cpu",
            compute_type="int8",
            download_root=katalog_modeli or None,
            local_files_only=True,
        )
    except Exception as wyjatek:
        print(f"[transkrypcja] model niedostępny: {wyjatek!r}", file=sys.stderr)
        odpowiedz["powod"] = _silnik.powod_braku_modelu(model, katalog_modeli)
        _wypisz(odpowiedz)
        return 0

    odpowiedz["gotowy"] = True
    odpowiedz["model"] = model
    _wypisz(odpowiedz)
    return 0


def tryb_plik(plik: str, model: str, jezyk: str, katalog_modeli: str) -> int:
    """Transkrypcja jednego pliku. Kod 0 przy sukcesie, 1 przy odmowie."""
    sciezka = Path(plik).expanduser()

    odmowa = _sprawdz_plik(sciezka)
    if odmowa:
        _wypisz({"blad": odmowa})
        return 1

    # Obecność silnika sprawdzana jest przed rozpoznaniem mowy, ale po walidacji
    # ścieżki: gdy zawodzą obie rzeczy naraz, użyteczniejsza jest informacja
    # o pliku, który użytkownik widzi, niż o bibliotece, której nie widzi.
    if not _silnik.silnik_obecny():
        _wypisz({"blad": _silnik.powod_braku_silnika()})
        return 1

    try:
        tekst, trwanie_s = _silnik.transkrybuj(sciezka, model, jezyk, katalog_modeli)
    except Exception as wyjatek:
        # Ślad wyjątku idzie na stderr, na stdout trafia sama odmowa w JSON:
        # szczegół techniczny służy zgłoszeniu, JSON — odczytowi przez rdzeń.
        print(f"[transkrypcja] wyjątek silnika: {wyjatek!r}", file=sys.stderr)
        _wypisz(
            {
                "blad": (
                    f"Transkrypcja pliku {sciezka} nie powiodła się. "
                    f"Silnik faster-whisper przerwał pracę komunikatem: "
                    f"{wyjatek}. "
                    f"Operator zacznie od sprawdzenia gotowości poleceniem "
                    f"--wersja; jeśli brakuje modelu „{model}”, pobranie go "
                    f"usunie ten błąd."
                )
            }
        )
        return 1

    # Pusty tekst jest wynikiem, nie odmową: przetworzenie się odbyło i wykazało,
    # że nagranie nie zawiera mowy. Dlatego pole `stan` rozróżnia „rozpoznano"
    # od „przetworzono, mowy brak"; przypadek „nie przetworzono" nie pojawia się
    # w tym miejscu, bo idzie osobnym polem `blad`.
    _wypisz(
        {
            "tekst": tekst,
            "znakow": len(tekst),
            # trwanie_ms to długość nagrania, nie czas przetwarzania.
            "trwanie_ms": int(round(trwanie_s * 1000)),
            "model": model,
            "jezyk": jezyk,
            "stan": STAN_ROZPOZNANO if tekst.strip() else STAN_BEZ_MOWY,
        }
    )
    return 0


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(
        prog="transkrypcja.py",
        description="Lokalna transkrypcja nagrań (CPU, bez sieci).",
    )
    parser.add_argument("--wersja", action="store_true", help="sprawdź gotowość")
    parser.add_argument("--plik", default="", help="ścieżka do nagrania")
    parser.add_argument("--model", default=DOMYSLNY_MODEL, help="rozmiar modelu")
    parser.add_argument("--jezyk", default=DOMYSLNY_JEZYK, help="język nagrania")
    parser.add_argument(
        "--katalog-modeli",
        dest="katalog_modeli",
        default="",
        help="katalog wag modelu (pusty = domyślny katalog biblioteki)",
    )
    args = parser.parse_args(argv)

    if args.wersja:
        return tryb_wersja(args.model, args.katalog_modeli)
    if args.plik:
        return tryb_plik(args.plik, args.model, args.jezyk, args.katalog_modeli)

    # Brak trybu również kończy się odmową w JSON, a nie pomocą argparse na
    # stdout: rdzeń woła ten skrypt maszynowo i musi przeczytać każdą odpowiedź.
    _wypisz(
        {
            "blad": (
                "Nie wskazano trybu pracy pomocnika. "
                "Skrypt nie zgaduje, czy ma sprawdzić gotowość, czy coś "
                "transkrybować, bo te dwie rzeczy kosztują nieporównywalnie "
                "różnie. "
                "Operator wskaże tryb: --wersja albo --plik ŚCIEŻKA."
            )
        }
    )
    return 1


if __name__ == "__main__":
    sys.exit(main())
