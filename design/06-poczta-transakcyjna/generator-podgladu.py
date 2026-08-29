#!/usr/bin/env python3
"""Wytwarza podglad.html z trzech szablonów w html/.

Podstawia dane przykładowe za zmienne {{ }} i osadza znak jako data URI,
żeby plik podglądu był samowystarczalny. Wysyłka produkcyjna używa cid:.
Uruchomienie: python3 generator-podgladu.py
"""

import base64
import html
import re
from pathlib import Path

KATALOG = Path(__file__).resolve().parent
KATALOG_HTML = KATALOG / "html"
PLIK_ZNAKU = KATALOG.parent / "03-marka" / "zastosowania" / "png" / "znak-poczty-jasny@2x.png"
PLIK_WYNIKOWY = KATALOG / "podglad.html"

SZABLONY = [
    ("list-01-aktywacja-konta.html", "List 1 · Aktywacja konta",
     "Danaco Console — kod aktywacji konta: 418 402"),
    ("list-02-uwierzytelnienie-logowania.html", "List 2 · Uwierzytelnienie logowania",
     "Danaco Console — kod logowania: 418 402"),
    ("list-03-autoresponder-brak-skrzynki.html", "List 3 · Autoresponder",
     "Adres noreply@danaco-group.pl nie przyjmuje korespondencji"),
]

# Dane przykładowe. Nie są danymi rzeczywistego Operatora ani rzeczywistej sprawy.
DANE_PRZYKLADOWE = {
    "code": "418 402",
    "expiry_minutes": "60",
    "requested_at": "29.08.2026, 17:22 CEST",
    "ip_address": "10.14.2.37",
    "device": "Windows 11 · Danaco Console 2.0",
    "activation_url": "https://console.danaco-group.pl/aktywacja/9f2ca71b",
    "recipient_address": "operator@przyklad.pl",
    "support_address": "support@danaco-group.pl",
    "original_subject": "Prośba o zwiększenie limitu przebiegów",
    "received_at": "29.08.2026, 17:22 CEST",
    "year": "2026",
}


def znak_data_uri() -> str:
    dane = base64.b64encode(PLIK_ZNAKU.read_bytes()).decode("ascii")
    return f"data:image/png;base64,{dane}"


def tresc_ciala(zrodlo: str) -> str:
    """Zwraca zawartość <body> szablonu — bez powłoki dokumentu."""
    dopasowanie = re.search(r"<body[^>]*>(.*)</body>", zrodlo, re.S)
    if not dopasowanie:
        raise ValueError("szablon bez znacznika body")
    return dopasowanie.group(1)


def podstaw(zrodlo: str, znak: str) -> str:
    wartosci = dict(DANE_PRZYKLADOWE, logo_src=znak)
    # {{original_subject}} pochodzi spoza rdzenia — ucieczka HTML jak w produkcji.
    wartosci["original_subject"] = html.escape(wartosci["original_subject"], quote=True)[:120]
    for klucz, wartosc in wartosci.items():
        zrodlo = zrodlo.replace("{{" + klucz + "}}", wartosc)
    pozostale = set(re.findall(r"\{\{(\w+)\}\}", zrodlo))
    if pozostale:
        raise ValueError(f"zmienne bez wartości: {sorted(pozostale)}")
    return zrodlo


STRONA = """<!doctype html>
<html lang="pl">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Danaco Console — poczta transakcyjna · podgląd</title>
<style>
@font-face {{ font-family:'Space Grotesk'; font-weight:600; font-display:swap;
  src:url('../zasoby/fonty/space-grotesk-latin-ext-600-normal.woff2') format('woff2'); }}
@font-face {{ font-family:'Space Grotesk'; font-weight:700; font-display:swap;
  src:url('../zasoby/fonty/space-grotesk-latin-ext-700-normal.woff2') format('woff2'); }}
@font-face {{ font-family:'IBM Plex Sans'; font-weight:400; font-display:swap;
  src:url('../zasoby/fonty/ibm-plex-sans-latin-ext-400-normal.woff2') format('woff2'); }}
@font-face {{ font-family:'IBM Plex Sans'; font-weight:600; font-display:swap;
  src:url('../zasoby/fonty/ibm-plex-sans-latin-ext-600-normal.woff2') format('woff2'); }}
@font-face {{ font-family:'IBM Plex Mono'; font-weight:400; font-display:swap;
  src:url('../zasoby/fonty/ibm-plex-mono-latin-ext-400-normal.woff2') format('woff2'); }}
@font-face {{ font-family:'IBM Plex Mono'; font-weight:600; font-display:swap;
  src:url('../zasoby/fonty/ibm-plex-mono-latin-ext-600-normal.woff2') format('woff2'); }}
:root {{ color-scheme: light only; }}
* {{ box-sizing:border-box; }}
body {{ margin:0; background:#F4F4F4; color:#181818;
  font-family:'IBM Plex Sans','Segoe UI',Helvetica,Arial,sans-serif; }}
.strona {{ max-width:1120px; margin:0 auto; padding:48px 24px 96px; }}
.naglowek {{ border-top:3px solid #181818; padding-top:24px; margin-bottom:48px; }}
.naglowek h1 {{ margin:0 0 8px; font-family:'Space Grotesk','IBM Plex Sans',sans-serif;
  font-size:40px; line-height:48px; font-weight:700; letter-spacing:-0.01em; }}
.naglowek p {{ margin:0; max-width:62ch; font-size:15px; line-height:24px; color:#4A4A4A; }}
.metryka {{ margin-top:20px; font-family:'IBM Plex Mono',Consolas,monospace;
  font-size:11px; line-height:18px; letter-spacing:0.14em; text-transform:uppercase; color:#616161; }}
.sekcja {{ margin-bottom:56px; }}
.pasek {{ display:flex; flex-wrap:wrap; gap:8px 24px; align-items:baseline;
  border-bottom:1px solid #E3E3E3; padding-bottom:12px; margin-bottom:24px; }}
.pasek h2 {{ margin:0; font-family:'Space Grotesk','IBM Plex Sans',sans-serif;
  font-size:24px; line-height:32px; font-weight:600; letter-spacing:-0.01em; }}
.temat {{ font-family:'IBM Plex Mono',Consolas,monospace; font-size:12px;
  line-height:20px; color:#616161; }}
.temat b {{ font-weight:400; color:#181818; }}
.ramka {{ border:1px solid #E3E3E3; background:#F4F4F4; }}
.stopka {{ border-top:1px solid #E3E3E3; padding-top:20px; font-size:12px;
  line-height:20px; color:#616161; max-width:76ch; }}
.stopka code {{ font-family:'IBM Plex Mono',Consolas,monospace; color:#181818; }}
</style>
</head>
<body>
<div class="strona">

<div class="naglowek">
  <h1>Poczta transakcyjna</h1>
  <p>Trzy szablony wiadomości wychodzących z adresu noreply@danaco-group.pl w ramach
     Danaco Console. Wariant jasny. Dane w podglądzie są przykładowe.</p>
  <p class="metryka">Karta 600 px · Space Grotesk 24/600 · IBM Plex Sans 15/24 · IBM Plex Mono 32/0,18 em</p>
</div>

{sekcje}

<div class="stopka">
  <p>Znak osadzony w podglądzie jako <code>data:</code> URI, żeby plik działał samodzielnie.
     W wysyłce produkcyjnej znak idzie jako załącznik <code>cid:danaco-lockup</code> —
     Gmail odrzuca schemat <code>data:</code> w obrazach.</p>
  <p>Kod źródłowy z komentarzem: <code>html/</code> · opracowanie: <code>opracowanie-techniczne.md</code></p>
</div>

</div>
</body>
</html>
"""

SEKCJA = """<div class="sekcja">
  <div class="pasek">
    <h2>{tytul}</h2>
    <span class="temat">Temat: <b>{temat}</b></span>
  </div>
  <div class="ramka">{tresc}</div>
</div>
"""


def main() -> None:
    znak = znak_data_uri()
    sekcje = []
    for nazwa, tytul, temat in SZABLONY:
        zrodlo = (KATALOG_HTML / nazwa).read_text(encoding="utf-8")
        tresc = podstaw(tresc_ciala(zrodlo), znak)
        sekcje.append(SEKCJA.format(tytul=tytul, temat=html.escape(temat), tresc=tresc))
    PLIK_WYNIKOWY.write_text(STRONA.format(sekcje="\n".join(sekcje)), encoding="utf-8")
    print(f"zapisano {PLIK_WYNIKOWY} ({PLIK_WYNIKOWY.stat().st_size} B)")


if __name__ == "__main__":
    main()
