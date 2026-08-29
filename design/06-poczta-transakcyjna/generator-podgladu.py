#!/usr/bin/env python3
"""Wytwarza podglad.html z szablonów w html/ i text/.

Podstawia dane przykładowe za zmienne {{ }} i osadza oba warianty znaku
jako data URI, żeby plik podglądu był samowystarczalny. Wysyłka produkcyjna
używa cid:. Podgląd przenosi też arkusz <style> pierwszego szablonu, więc
tryb ciemny działa zgodnie z ustawieniem systemu.
Uruchomienie: python3 generator-podgladu.py
"""

import base64
import html
import re
from pathlib import Path

KATALOG = Path(__file__).resolve().parent
KATALOG_HTML = KATALOG / "html"
KATALOG_TEXT = KATALOG / "text"
KATALOG_ZNAKOW = KATALOG.parent / "03-marka" / "zastosowania" / "png"
PLIK_WYNIKOWY = KATALOG / "podglad.html"

SZABLONY = [
    ("list-01-aktywacja-konta", "List 1 · Aktywacja konta",
     "Danaco Console — kod aktywacji konta: 418 402"),
    ("list-02-uwierzytelnienie-logowania", "List 2 · Uwierzytelnienie logowania",
     "Danaco Console — kod logowania: 418 402"),
    ("list-03-autoresponder-brak-skrzynki", "List 3 · Autoresponder",
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


def znak_data_uri(nazwa: str) -> str:
    dane = base64.b64encode((KATALOG_ZNAKOW / nazwa).read_bytes()).decode("ascii")
    return f"data:image/png;base64,{dane}"


def czysty_temat(wartosc: str) -> str:
    """Sanityzacja tematu spoza rdzenia — wspólna dla obu części wiadomości."""
    bez_sterujacych = re.sub(r"[\r\n\t\x00-\x1f\x7f]", " ", wartosc)
    return bez_sterujacych.strip()[:120]


def podstaw(zrodlo: str, wartosci: dict, ucieczka_html: bool) -> str:
    dane = dict(wartosci)
    # Część HTML wymaga ucieczki znaczników; część tekstowa nie — tam ucieczka
    # wprowadziłaby encje widoczne dla odbiorcy jako &amp; i &quot;.
    temat = czysty_temat(dane["original_subject"])
    dane["original_subject"] = html.escape(temat, quote=True) if ucieczka_html else temat
    for klucz, wartosc in dane.items():
        zrodlo = zrodlo.replace("{{" + klucz + "}}", wartosc)
    pozostale = set(re.findall(r"\{\{(\w+)\}\}", zrodlo))
    if pozostale:
        raise ValueError(f"zmienne bez wartości: {sorted(pozostale)}")
    return zrodlo


def tresc_ciala(zrodlo: str) -> str:
    dopasowanie = re.search(r"<body[^>]*>(.*)</body>", zrodlo, re.S)
    if not dopasowanie:
        raise ValueError("szablon bez znacznika body")
    return dopasowanie.group(1)


def arkusz_szablonu(zrodlo: str) -> str:
    dopasowanie = re.search(r"<style>(.*?)</style>", zrodlo, re.S)
    if not dopasowanie:
        raise ValueError("szablon bez arkusza trybu ciemnego")
    return dopasowanie.group(1)


STRONA = """<!doctype html>
<html lang="pl">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<meta name="color-scheme" content="light dark">
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

/* Powłoka podglądu. Barwy z tej samej skali co listy. */
:root {{ --tlo:#F4F4F4; --powierzchnia:#FFFFFF; --atrament:#181818;
        --tekst:#4A4A4A; --etykieta:#616161; --obrys:#E3E3E3; }}
@media (prefers-color-scheme: dark) {{
  :root {{ --tlo:#0F0F0F; --powierzchnia:#131313; --atrament:#ECECEC;
          --tekst:#9E9E9E; --etykieta:#9E9E9E; --obrys:#2A2A2A; }}
}}
* {{ box-sizing:border-box; }}
body {{ margin:0; background:var(--tlo); color:var(--atrament);
  font-family:'IBM Plex Sans','Segoe UI',Helvetica,Arial,sans-serif; }}
.strona {{ max-width:1120px; margin:0 auto; padding:48px 24px 96px; }}
.naglowek {{ border-top:3px solid var(--atrament); padding-top:24px; margin-bottom:48px; }}
.naglowek h1 {{ margin:0 0 8px; font-family:'Space Grotesk','IBM Plex Sans',sans-serif;
  font-size:40px; line-height:48px; font-weight:700; letter-spacing:-0.01em; }}
.naglowek p {{ margin:0; max-width:62ch; font-size:15px; line-height:24px; color:var(--tekst); }}
.metryka {{ margin-top:20px; font-family:'IBM Plex Mono',Consolas,monospace;
  font-size:11px; line-height:18px; letter-spacing:0.14em; text-transform:uppercase;
  color:var(--etykieta) !important; }}
.sekcja {{ margin-bottom:56px; }}
.pasek {{ display:flex; flex-wrap:wrap; gap:8px 24px; align-items:baseline;
  border-bottom:1px solid var(--obrys); padding-bottom:12px; margin-bottom:24px; }}
.pasek h2 {{ margin:0; font-family:'Space Grotesk','IBM Plex Sans',sans-serif;
  font-size:24px; line-height:32px; font-weight:600; letter-spacing:-0.01em; }}
.temat {{ font-family:'IBM Plex Mono',Consolas,monospace; font-size:12px;
  line-height:20px; color:var(--etykieta); }}
.temat b {{ font-weight:400; color:var(--atrament); }}
.ramka {{ border:1px solid var(--obrys); background:var(--tlo); }}
.czesc {{ margin-top:24px; }}
.czesc h3 {{ margin:0 0 8px; font-family:'IBM Plex Mono',Consolas,monospace;
  font-size:11px; line-height:16px; letter-spacing:0.14em; text-transform:uppercase;
  font-weight:400; color:var(--etykieta); }}
.czesc pre {{ margin:0; padding:24px 28px; overflow-x:auto;
  border:1px solid var(--obrys); background:var(--powierzchnia); color:var(--atrament);
  font-family:'IBM Plex Mono',Consolas,monospace; font-size:12.5px; line-height:20px; }}
.stopka {{ border-top:1px solid var(--obrys); padding-top:20px; font-size:12px;
  line-height:20px; color:var(--etykieta); max-width:76ch; }}
.stopka code {{ font-family:'IBM Plex Mono',Consolas,monospace; color:var(--atrament); }}
</style>
<style>
/* Arkusz przeniesiony z szablonu listu — obsługa trybu ciemnego wewnątrz kart. */
{arkusz}
</style>
</head>
<body>
<div class="strona">

<div class="naglowek">
  <h1>Poczta transakcyjna</h1>
  <p>Trzy szablony wiadomości wychodzących z adresu noreply@danaco-group.pl w ramach
     Danaco Console. Każdy list w dwóch częściach: HTML oraz text/plain. Wariant barwny
     zależy od motywu systemu — przełącz motyw, aby zobaczyć drugi. Dane są przykładowe.</p>
  <p class="metryka">Karta 600 px · Space Grotesk 24/600 · IBM Plex Sans 15/24 · IBM Plex Mono 32/0,18 em</p>
</div>

{sekcje}

<div class="stopka">
  <p>Znak osadzony w podglądzie jako <code>data:</code> URI, żeby plik działał samodzielnie.
     W wysyłce produkcyjnej oba warianty znaku idą jako załączniki
     <code>cid:danaco-lockup</code> i <code>cid:danaco-lockup-dark</code> —
     Gmail odrzuca schemat <code>data:</code> w obrazach.</p>
  <p>Kod źródłowy z komentarzem: <code>html/</code> i <code>text/</code> ·
     opracowanie: <code>opracowanie-techniczne.md</code></p>
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
  <div class="czesc">
    <h3>Część text/plain</h3>
    <pre>{tekst}</pre>
  </div>
</div>
"""


def main() -> None:
    znaki = {
        "logo_src": znak_data_uri("znak-poczty-jasny@2x.png"),
        "logo_src_dark": znak_data_uri("znak-poczty-ciemny@2x.png"),
    }
    wartosci = dict(DANE_PRZYKLADOWE, **znaki)
    arkusz = ""
    sekcje = []
    for nazwa, tytul, temat in SZABLONY:
        zrodlo = (KATALOG_HTML / f"{nazwa}.html").read_text(encoding="utf-8")
        arkusz = arkusz or arkusz_szablonu(zrodlo)
        tresc = podstaw(tresc_ciala(zrodlo), wartosci, ucieczka_html=True)
        tekst = podstaw((KATALOG_TEXT / f"{nazwa}.txt").read_text(encoding="utf-8"),
                        wartosci, ucieczka_html=False)
        sekcje.append(SEKCJA.format(tytul=tytul, temat=html.escape(temat),
                                    tresc=tresc, tekst=html.escape(tekst)))
    PLIK_WYNIKOWY.write_text(
        STRONA.format(sekcje="\n".join(sekcje), arkusz=arkusz), encoding="utf-8")
    print(f"zapisano {PLIK_WYNIKOWY} ({PLIK_WYNIKOWY.stat().st_size} B)")


if __name__ == "__main__":
    main()
